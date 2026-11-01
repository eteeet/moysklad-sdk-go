package supply

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
		return nil, fmt.Errorf("moysklad: supply service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/supply
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Supply, error) {
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: supply id must be a UUID")
	}
	return s.do(ctx, "get supply "+id, http.MethodGet, "entity/supply/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/supply
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль Supply.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Supply, error) {
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
	reader, err := positions.NewService[Position](s.client, "supply")
	if err != nil {
		return nil, err
	}
	return reader.GetPage(ctx, id, page, opts...)
}
