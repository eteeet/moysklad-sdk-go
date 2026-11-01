package country

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: country service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/country#3-poluchit-stranu
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Country, error) {
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: country id must be a UUID")
	}
	return s.do(ctx, "get country "+id, http.MethodGet, "entity/country/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/country#3-sozdat-stranu
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль Country.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Country, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create country: name must not be empty")
	}
	return s.do(ctx, "create country", http.MethodPost, "entity/country", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/country#3-izmenit-stranu
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Country, error) {
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update country: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update country: no fields specified")
	}
	return s.do(ctx, "update country "+id, http.MethodPut, "entity/country/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Country, error) {
	var result Country
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
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/country
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Country], error) {
	return list.Entity[Country](ctx, s.client, "country", options, nil, opts...)
}
