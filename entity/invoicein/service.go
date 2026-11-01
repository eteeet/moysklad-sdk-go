package invoicein

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
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
func (s *Service) Get(ctx context.Context, id string) (*InvoiceIn, error) {
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: invoicein id must be a UUID")
	}
	return s.do(ctx, "get invoicein "+id, http.MethodGet, "entity/invoicein/"+id, nil)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-in
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль InvoiceIn.
func (s *Service) Create(ctx context.Context, input CreateRequest) (*InvoiceIn, error) {
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
	return s.do(ctx, "create invoicein", http.MethodPost, "entity/invoicein", input)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-in
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest) (*InvoiceIn, error) {
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
	return s.do(ctx, "update invoicein "+id, http.MethodPut, "entity/invoicein/"+id, input)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any) (*InvoiceIn, error) {
	req, err := s.client.NewRequest(ctx, method, endpoint, nil, body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	var result InvoiceIn
	resp, err := s.client.Do(req, &result)
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
func (s *Service) GetPositions(ctx context.Context, id string, page positions.PageOptions) (*model.Collection[Position], error) {
	reader, err := positions.NewService[Position](s.client, "invoicein")
	if err != nil {
		return nil, err
	}
	return reader.GetPage(ctx, id, page)
}
