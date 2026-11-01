package currency

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
		return nil, fmt.Errorf("moysklad: currency service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/currency#3-poluchit-valyutu
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Currency, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: currency id must be a UUID")
	}
	return s.do(ctx, "get currency "+id, http.MethodGet, "entity/currency/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/currency#3-sozdat-novuyu-valyutu
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль Currency.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Currency, error) {
	transport.ResetResponseInfo(opts)
	if err := input.validate(); err != nil {
		return nil, fmt.Errorf("moysklad: create currency: %w", err)
	}
	return s.do(ctx, "create currency", http.MethodPost, "entity/currency", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/currency#3-izmenit-valyutu
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Currency, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update currency: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update currency: no fields specified")
	}
	if err := checkValues(input.Rate, input.Name, input.Code, input.ISOCode); err != nil {
		return nil, fmt.Errorf("moysklad: update currency: %w", err)
	}
	return s.do(ctx, "update currency "+id, http.MethodPut, "entity/currency/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Currency, error) {
	var result Currency
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
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/currency
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Currency], error) {
	return list.Entity[Currency](ctx, s.client, "currency", options, nil, opts...)
}
