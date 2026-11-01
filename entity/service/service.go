package service

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

// Client — клиент API услуг. Сущность называется Service, поэтому клиент здесь не Service,
// как в других пакетах; конструктор — NewService, как у остальных.
type Client struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Client, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: service client requires a transport client")
	}
	return &Client{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/service#3-poluchit-uslugu
func (s *Client) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Service, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: service id must be a UUID")
	}
	return s.do(ctx, "get service "+id, http.MethodGet, "entity/service/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/service#3-sozdat-uslugu
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль Service.
func (s *Client) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Service, error) {
	transport.ResetResponseInfo(opts)
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create service: name must not be empty")
	}
	return s.do(ctx, "create service", http.MethodPost, "entity/service", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/service#3-izmenit-uslugu
func (s *Client) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Service, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update service: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update service: no fields specified")
	}
	return s.do(ctx, "update service "+id, http.MethodPut, "entity/service/"+id, input, opts)
}

func (s *Client) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Service, error) {
	var result Service
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
// Поиск — по name и code; фильтр штрихкода — поле barcode в единственном числе; вместе
// с фильтром по id фильтр archived сервер не учитывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/service
func (s *Client) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Service], error) {
	return list.Entity[Service](ctx, s.client, "service", options, nil, opts...)
}
