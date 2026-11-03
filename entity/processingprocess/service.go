package processingprocess

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// maxPositions — лимит позиций техпроцесса по источнику («от 1 до 100 позиций»).
const maxPositions = 100

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: processingprocess service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-poluchit-tehprocess
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*ProcessingProcess, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get processingprocess: id must be a UUID")
	}
	return s.do(ctx, "get processingprocess "+id, http.MethodGet, "entity/processingprocess/"+id, nil, opts)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-poluchit-spisok-tehprocessov
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[ProcessingProcess], error) {
	return list.Entity[ProcessingProcess](ctx, s.client, "processingprocess", options, nil, opts...)
}

// Create: name и 1–100 позиций с processingstage проверяются до запроса. Не повторять POST
// автоматически после неоднозначного сбоя: возможен дубль.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-sozdat-tehprocess
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*ProcessingProcess, error) {
	transport.ResetResponseInfo(opts)
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create processingprocess: name must not be empty")
	}
	if err := checkNew(input.Positions); err != nil {
		return nil, fmt.Errorf("moysklad: create processingprocess: %w", err)
	}
	return s.do(ctx, "create processingprocess", http.MethodPost, "entity/processingprocess", input, opts)
}

// Update: пустой request и Positions вне 1–100 (или новая строка без processingstage)
// отклоняются до запроса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-izmenit-tehprocess
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*ProcessingProcess, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update processingprocess: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update processingprocess: no fields specified")
	}
	if rows, ok := input.Positions.Get(); ok || input.Positions.IsNull() {
		if err := checkRows(rows); err != nil {
			return nil, fmt.Errorf("moysklad: update processingprocess %s: %w", id, err)
		}
	}
	return s.do(ctx, "update processingprocess "+id, http.MethodPut, "entity/processingprocess/"+id, input, opts)
}

// Delete — DELETE entity/processingprocess/{id}, безвозвратно. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-udalit-tehprocess
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete processingprocess: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete processingprocess "+id, "entity/processingprocess/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/processingprocess (до batch.MaxItems); обязательные
// поля элементов проверяет сервер, ответ — по правилам batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-massovoe-sozdanie-i-obnovlenie-tehprocessov
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[ProcessingProcess], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[ProcessingProcess](ctx, s.client, "processingprocess", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/processingprocess: meta по ID и поля Request;
// до запроса проверяются только UUID и непустой Request. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-massovoe-sozdanie-i-obnovlenie-tehprocessov
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[ProcessingProcess], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[ProcessingProcess](ctx, s.client, "processingprocess", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/processingprocess/delete. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-massovoe-udalenie-tehprocessov
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "processingprocess", ids, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*ProcessingProcess, error) {
	var result ProcessingProcess
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

func checkCount(n int) error {
	if n == 0 || n > maxPositions {
		return fmt.Errorf("positions must contain 1 to %d elements, got %d", maxPositions, n)
	}
	return nil
}

func checkNew(rows []PositionRequest) error {
	if err := checkCount(len(rows)); err != nil {
		return err
	}
	for i, row := range rows {
		if row.ProcessingStage.Meta.Href == "" {
			return fmt.Errorf("positions[%d]: processingstage meta.href must not be empty", i)
		}
	}
	return nil
}

// Строка без meta — новая позиция: processingstage обязателен.
func checkRows(rows []PositionUpdateRequest) error {
	if err := checkCount(len(rows)); err != nil {
		return err
	}
	for i, row := range rows {
		stage, ok := row.ProcessingStage.Get()
		if row.Meta.IsZero() && (!ok || stage.Meta.Href == "") {
			return fmt.Errorf("positions[%d]: new position needs processingstage meta.href", i)
		}
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
