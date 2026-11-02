// Package async — асинхронные задачи JSON API: статус, список, ожидание, результат и отмена.
//
// Задачу создают только ListAsync сервисов counterparty и assortment — из запросов SDK источник
// разрешает async лишь им; остальные async-отчёты в SDK не входят. Создание задачи — побочный
// эффект и у GET: transport его не повторяет. В очереди аккаунта не больше 4 задач — это лимит
// сервера на весь аккаунт (ошибка 61002), один Client его не контролирует. Если ListAsync
// вернул ошибку сети или отмены ctx после отправки, задача могла создаться и занять место в
// очереди: проверьте List с фильтром по request, прежде чем создавать её снова. Результат
// хранится 1 час (до DeletionDate). Пока результат готовится, в нём возможны дубли и пропуски
// строк. Ошибки async-запросов — коды 61000–61008.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#2-asinhronnyj-obmen
package async

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// State — статус задачи; набор открыт.
type State string

const (
	StatePending    State = "PENDING"    // в очереди
	StateProcessing State = "PROCESSING" // выполняется
	StateDone       State = "DONE"       // результат готов
	StateError      State = "ERROR"      // внутренняя ошибка; источник советует запустить задачу заново
	StateCancel     State = "CANCEL"     // отменена
	StateAPIError   State = "API_ERROR"  // ошибка API в Errors
)

// Finished — задача завершена: DONE, ERROR, CANCEL или API_ERROR. Неизвестный статус не
// считается завершённым.
func (s State) Finished() bool {
	return s == StateDone || s == StateError || s == StateCancel || s == StateAPIError
}

// Task — асинхронная задача; все поля только для чтения. ResultURL и DeletionDate приходят
// у DONE, Errors — у API_ERROR. Таблица называет errors типом Object, пример — массив ошибок
// API; SDK следует примеру.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-asinhronnaya-zadacha
type Task struct {
	Meta         model.Meta         `json:"meta"`
	ID           string             `json:"id"`
	AccountID    string             `json:"accountId"`
	Owner        *model.Reference   `json:"owner,omitzero"`
	State        State              `json:"state"`
	Request      string             `json:"request"`
	ResultURL    string             `json:"resultUrl,omitempty"`
	DeletionDate model.DateTime     `json:"deletionDate,omitzero"`
	Errors       []sdkerrors.Detail `json:"errors,omitempty"`
}

// Submitted — созданная задача из ответа 202: ID, адрес статуса (Content-Location) и адрес
// результата (Location).
type Submitted struct {
	ID         string
	StatusHref string
	ResultHref string
}

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: async service requires a client")
	}
	return &Service{client: client}, nil
}

// Get читает статус задачи.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-poluchenie-statusa-asinhronnoj-zadachi
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Task, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get async task: id must be a UUID")
	}
	operation := "get async task " + id
	var result Task
	resp, err := s.client.Send(ctx, http.MethodGet, "async/"+id, nil, nil, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result.ID == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
	}
	return &result, nil
}

// List читает одну страницу статусов задач за последнюю неделю. Фильтры по источнику:
// deletionDate, request и state; сортировка — deletionDate и request.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-statusy-asinhronnyh-zadach
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Task], error) {
	transport.ResetResponseInfo(opts)
	params, err := options.Values(query.MaxLimit)
	if err != nil {
		return nil, fmt.Errorf("list async tasks: %w", err)
	}
	return list.Page[Task](ctx, s.client, "list async tasks", "async", params, opts...)
}

// Wait опрашивает статус через Get каждые interval, пока задача не завершится (State.Finished),
// и возвращает последний статус; первый Get — сразу. Срок ожидания задаёт только ctx: без
// дедлайна опрос длится, пока задача не завершится. ERROR, CANCEL
// и API_ERROR — завершение, а не ошибка Wait: их разбирает вызывающий код. ResponseInfo
// описывает последний Get.
func (s *Service) Wait(ctx context.Context, id string, interval time.Duration, opts ...transport.RequestOption) (*Task, error) {
	transport.ResetResponseInfo(opts)
	if interval <= 0 {
		return nil, fmt.Errorf("moysklad: wait async task: interval must be positive")
	}
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("moysklad: wait async task %s: %w", id, ctx.Err())
		case <-timer.C:
		}
		task, err := s.Get(ctx, id, opts...)
		if err != nil || task.State.Finished() {
			return task, err
		}
		timer.Reset(interval)
	}
}

// Result копирует результат задачи в w через transport.Client.Download: GET async/{id}/result,
// redirect на ссылку загрузки (действует 5 минут) без Authorization. Результат — JSON того же
// вида, что ответ синхронного запроса. Ошибки — APIError: у API_ERROR — ошибка исходного запроса,
// иначе 61003 (результат удалён после DeletionDate), 61004, 61005, 61006 (ещё не готов).
// Успех — только итоговый HTTP 200, как в источнике: другой 2xx — ошибка (тело такого ответа
// могло уже попасть в w, n это показывает).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-poluchenie-rezultata-vypolneniya-asinhronnoj-zadachi
func (s *Service) Result(ctx context.Context, id string, w io.Writer, opts ...transport.RequestOption) (int64, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return 0, fmt.Errorf("moysklad: get async task result: id must be a UUID")
	}
	// Статус нужен и без WithResponseInfo вызывающего кода; его указатель, если есть, сохраняем.
	var options transport.RequestOptions
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}
	if options.ResponseInfo == nil {
		options.ResponseInfo = new(transport.ResponseInfo)
		opts = append(slices.Clip(opts), transport.WithResponseInfo(options.ResponseInfo))
	}
	n, err := s.client.Download(ctx, "async/"+id+"/result", w, opts...)
	if err != nil {
		return n, fmt.Errorf("get async task %s result: %w", id, err)
	}
	if status := options.ResponseInfo.StatusCode; status != http.StatusOK {
		return n, fmt.Errorf("moysklad: get async task %s result: unexpected HTTP status %d", id, status)
	}
	return n, nil
}

// Cancel отменяет задачу в PENDING или PROCESSING: PUT async/{id}/cancel, успех — 204; уже
// завершённая задача — ошибка 61007. Отмена останавливает подготовку результата, это не откат
// изменений. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-otmena-asinhronnoj-zadachi
func (s *Service) Cancel(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: cancel async task: id must be a UUID")
	}
	operation := "cancel async task " + id
	resp, err := s.client.Send(ctx, http.MethodPut, "async/"+id+"/cancel", nil, nil, nil, opts...)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	return nil
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
