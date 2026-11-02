package invoicein

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
		return nil, fmt.Errorf("moysklad: invoicein service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-in
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*InvoiceIn, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: invoicein id must be a UUID")
	}
	return s.do(ctx, "get invoicein "+id, http.MethodGet, "entity/invoicein/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-in
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль InvoiceIn.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*InvoiceIn, error) {
	transport.ResetResponseInfo(opts)
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create invoicein: name must not be empty")
	}
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create invoicein: organization.meta.href must not be empty")
	}
	if input.Agent.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create invoicein: agent.meta.href must not be empty")
	}
	if rows, ok := input.Positions.Get(); ok && len(rows) > 1000 {
		return nil, fmt.Errorf("moysklad: create invoicein: positions must not exceed 1000")
	}
	return s.do(ctx, "create invoicein", http.MethodPost, "entity/invoicein", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-in
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*InvoiceIn, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update invoicein: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update invoicein: no fields specified")
	}
	if rows, ok := input.Positions.Get(); ok && len(rows) > 1000 {
		return nil, fmt.Errorf("moysklad: update invoicein: positions must not exceed 1000")
	}
	return s.do(ctx, "update invoicein "+id, http.MethodPut, "entity/invoicein/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*InvoiceIn, error) {
	var result InvoiceIn
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
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-in
func (s *Service) GetPositions(ctx context.Context, id string, page positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	transport.ResetResponseInfo(opts)
	reader, err := positions.NewService[Position](s.client, "invoicein")
	if err != nil {
		return nil, err
	}
	return reader.GetPage(ctx, id, page, opts...)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// Позиции в строках — только то, что вернул сервер; строки позиций читает GetPositions.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoicein
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[InvoiceIn], error) {
	return list.Entity[InvoiceIn](ctx, s.client, "invoicein", options, nil, opts...)
}

// Delete — «Удалить Счет поставщика»: DELETE entity/invoicein/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-in#3-udalit-schet-postavshika
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete invoicein: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete invoicein "+id, "entity/invoicein/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/invoicein: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Счетов поставщика».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-in#3-massovoe-sozdanie-i-obnovlenie-schetov-postavshika
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[InvoiceIn], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[InvoiceIn](ctx, s.client, "invoicein", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/invoicein: элемент — meta по ID и поля Request.
// Источник: «Массовое создание и обновление Счетов поставщика».
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-in#3-massovoe-sozdanie-i-obnovlenie-schetov-postavshika
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[InvoiceIn], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[InvoiceIn](ctx, s.client, "invoicein", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/invoicein/delete, не в корзину; до
// batch.MaxItems элементов. Источник: «Массовое удаление Счетов поставщиков».
// Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-in#3-massovoe-udalenie-schetov-postavshikov
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "invoicein", ids, opts)
}
