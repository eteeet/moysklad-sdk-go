package counterparty

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/async"
	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: counterparty service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-poluchit-kontragenta
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Counterparty, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: counterparty id must be a UUID")
	}
	return s.do(ctx, "get counterparty "+id, http.MethodGet, "entity/counterparty/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-sozdat-kontragenta
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль Counterparty.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Counterparty, error) {
	transport.ResetResponseInfo(opts)
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create counterparty: name must not be empty")
	}
	return s.do(ctx, "create counterparty", http.MethodPost, "entity/counterparty", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-izmenit-kontragenta
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Counterparty, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update counterparty: id must be a UUID")
	}
	// Нулевые Field — «не передано»; false, пустые значения и Null считаются заданными.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update counterparty: no fields specified")
	}
	return s.do(ctx, "update counterparty "+id, http.MethodPut, "entity/counterparty/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Counterparty, error) {
	var result Counterparty
	resp, err := s.client.Send(ctx, method, endpoint, nil, body, &result, opts...)
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

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Counterparty], error) {
	return list.Entity[Counterparty](ctx, s.client, "counterparty", options, nil, opts...)
}

// ListAsync создаёт асинхронную задачу «Получить список Контрагентов»: GET entity/counterparty
// с async=true и фильтрами options. Limit и Offset источник у async запрещает, fields List не
// принимает. Ответ — созданная задача; весь список — через async.Service.Result. Повторов нет;
// исход после сбоя — по правилам пакета async.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime
func (s *Service) ListAsync(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*async.Submitted, error) {
	transport.ResetResponseInfo(opts)
	const operation = "list counterparty async"
	if options.Fields != "" {
		return nil, fmt.Errorf("moysklad: %s: fields=%q is not supported", operation, options.Fields)
	}
	params, err := options.Values(query.MaxLimit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return write.Submit(ctx, s.client, operation, "entity/counterparty", params, opts)
}

// Delete — «Удалить Контрагента»: DELETE entity/counterparty/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-udalit-kontragenta
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete counterparty: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete counterparty "+id, "entity/counterparty/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/counterparty: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Контрагентов».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-massovoe-sozdanie-i-obnovlenie-kontragentov
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[Counterparty], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[Counterparty](ctx, s.client, "counterparty", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/counterparty: элемент — meta по ID и поля Request.
// Источник: «Массовое создание и обновление Контрагентов».
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-massovoe-sozdanie-i-obnovlenie-kontragentov
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Counterparty], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[Counterparty](ctx, s.client, "counterparty", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/counterparty/delete, не в корзину; до
// batch.MaxItems элементов. Источник: «Массовое удаление Контрагентов».
// Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-massovoe-udalenie-kontragentov
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "counterparty", ids, opts)
}
