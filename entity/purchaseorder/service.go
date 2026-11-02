package purchaseorder

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
		return nil, fmt.Errorf("moysklad: purchaseorder service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/purchaseOrder
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*PurchaseOrder, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: purchaseorder id must be a UUID")
	}
	return s.do(ctx, "get purchaseorder "+id, http.MethodGet, "entity/purchaseorder/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/purchaseOrder
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль PurchaseOrder.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*PurchaseOrder, error) {
	transport.ResetResponseInfo(opts)
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create purchaseorder: organization.meta.href must not be empty")
	}
	if input.Agent.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create purchaseorder: agent.meta.href must not be empty")
	}
	if rows, ok := input.Positions.Get(); ok && len(rows) > 1000 {
		return nil, fmt.Errorf("moysklad: create purchaseorder: positions must not exceed 1000")
	}
	return s.do(ctx, "create purchaseorder", http.MethodPost, "entity/purchaseorder", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/purchaseOrder
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*PurchaseOrder, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update purchaseorder: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update purchaseorder: no fields specified")
	}
	if rows, ok := input.Positions.Get(); ok && len(rows) > 1000 {
		return nil, fmt.Errorf("moysklad: update purchaseorder: positions must not exceed 1000")
	}
	return s.do(ctx, "update purchaseorder "+id, http.MethodPut, "entity/purchaseorder/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*PurchaseOrder, error) {
	var result PurchaseOrder
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
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/purchaseOrder
func (s *Service) GetPositions(ctx context.Context, id string, page positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	transport.ResetResponseInfo(opts)
	reader, err := positions.NewService[Position](s.client, "purchaseorder")
	if err != nil {
		return nil, err
	}
	return reader.GetPage(ctx, id, page, opts...)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// Позиции в строках — только то, что вернул сервер; строки позиций читает GetPositions.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/purchaseorder
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[PurchaseOrder], error) {
	return list.Entity[PurchaseOrder](ctx, s.client, "purchaseorder", options, nil, opts...)
}

// Delete — «Удалить Заказ поставщику»: DELETE entity/purchaseorder/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/purchaseOrder#3-udalit-zakaz-postavshiku
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete purchaseorder: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete purchaseorder "+id, "entity/purchaseorder/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/purchaseorder: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Заказов поставщику».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/purchaseOrder#3-massovoe-sozdanie-i-obnovlenie-zakazov-postavshiku
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[PurchaseOrder], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[PurchaseOrder](ctx, s.client, "purchaseorder", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/purchaseorder: элемент — meta по ID и поля Request.
// Источник: «Массовое создание и обновление Заказов поставщику».
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/purchaseOrder#3-massovoe-sozdanie-i-obnovlenie-zakazov-postavshiku
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[PurchaseOrder], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[PurchaseOrder](ctx, s.client, "purchaseorder", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/purchaseorder/delete, не в корзину; до
// batch.MaxItems элементов. Источник: «Массовое удаление Заказов Поставщикам».
// Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/purchaseOrder#3-massovoe-udalenie-zakazov-postavshikam
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "purchaseorder", ids, opts)
}
