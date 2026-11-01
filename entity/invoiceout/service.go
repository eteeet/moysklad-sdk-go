package invoiceout

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
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
		return nil, fmt.Errorf("moysklad: invoiceout service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-out
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*InvoiceOut, error) {
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: invoiceout id must be a UUID")
	}
	return s.do(ctx, "get invoiceout "+id, http.MethodGet, "entity/invoiceout/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-out
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль InvoiceOut.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*InvoiceOut, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create invoiceout: name must not be empty")
	}
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create invoiceout: organization.meta.href must not be empty")
	}
	if input.Agent.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create invoiceout: agent.meta.href must not be empty")
	}
	if rows, ok := input.Positions.Get(); ok && len(rows) > 1000 {
		return nil, fmt.Errorf("moysklad: create invoiceout: positions must not exceed 1000")
	}
	return s.do(ctx, "create invoiceout", http.MethodPost, "entity/invoiceout", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-out
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*InvoiceOut, error) {
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update invoiceout: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update invoiceout: no fields specified")
	}
	if rows, ok := input.Positions.Get(); ok && len(rows) > 1000 {
		return nil, fmt.Errorf("moysklad: update invoiceout: positions must not exceed 1000")
	}
	return s.do(ctx, "update invoiceout "+id, http.MethodPut, "entity/invoiceout/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*InvoiceOut, error) {
	var result InvoiceOut
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
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-out
func (s *Service) GetPositions(ctx context.Context, id string, page positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	reader, err := positions.NewService[Position](s.client, "invoiceout")
	if err != nil {
		return nil, err
	}
	return reader.GetPage(ctx, id, page, opts...)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// Позиции в строках — только то, что вернул сервер; строки позиций читает GetPositions.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoiceout
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[InvoiceOut], error) {
	return list.Entity[InvoiceOut](ctx, s.client, "invoiceout", options, nil, opts...)
}
