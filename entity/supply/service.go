package supply

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
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: supply service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/supply
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Supply, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: supply id must be a UUID")
	}
	return s.do(ctx, "get supply "+id, http.MethodGet, "entity/supply/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/supply
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль Supply.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Supply, error) {
	transport.ResetResponseInfo(opts)
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create supply: organization.meta.href must not be empty")
	}
	if input.Agent.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create supply: agent.meta.href must not be empty")
	}
	if input.Store.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create supply: store.meta.href must not be empty")
	}
	if rows, ok := input.Positions.Get(); ok && len(rows) > 1000 {
		return nil, fmt.Errorf("moysklad: create supply: positions must not exceed 1000")
	}
	return s.do(ctx, "create supply", http.MethodPost, "entity/supply", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/supply
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Supply, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update supply: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update supply: no fields specified")
	}
	if rows, ok := input.Positions.Get(); ok && len(rows) > 1000 {
		return nil, fmt.Errorf("moysklad: update supply: positions must not exceed 1000")
	}
	return s.do(ctx, "update supply "+id, http.MethodPut, "entity/supply/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Supply, error) {
	var result Supply
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

// GetPositions читает одну страницу без expand, скрытого Get и автоматического обхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/supply
func (s *Service) GetPositions(ctx context.Context, id string, page positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	transport.ResetResponseInfo(opts)
	reader, err := positions.NewService[Position](s.client, "supply")
	if err != nil {
		return nil, err
	}
	return reader.GetPage(ctx, id, page, opts...)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// Позиции в строках — только то, что вернул сервер; строки позиций читает GetPositions.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/supply
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Supply], error) {
	return list.Entity[Supply](ctx, s.client, "supply", options, nil, opts...)
}

// Delete — «Удалить Приемку»: DELETE entity/supply/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/supply#3-udalit-priemku
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete supply: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete supply "+id, "entity/supply/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/supply: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Приемок».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/supply#3-massovoe-sozdanie-i-obnovlenie-priemok
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[Supply], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[Supply](ctx, s.client, "supply", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/supply: элемент — meta по ID и поля Request.
// Источник: «Массовое создание и обновление Приемок».
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/supply#3-massovoe-sozdanie-i-obnovlenie-priemok
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Supply], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[Supply](ctx, s.client, "supply", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/supply/delete, не в корзину; до
// batch.MaxItems элементов. Источник: «Массовое удаление Приемок».
// Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/supply#3-massovoe-udalenie-priemok
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "supply", ids, opts)
}
