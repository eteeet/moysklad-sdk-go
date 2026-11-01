package retailstore

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
		return nil, fmt.Errorf("moysklad: retailstore service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/retailstore#3-poluchit-tochku-prodazh
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*RetailStore, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get retailstore: id must be a UUID")
	}
	return s.do(ctx, "get retailstore "+id, http.MethodGet, "entity/retailstore/"+id, nil, opts)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль точки продаж.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/retailstore#3-sozdat-tochku-prodazh
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*RetailStore, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create retailstore"
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: %s: name must not be empty", operation)
	}
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: organization meta.href must not be empty", operation)
	}
	if input.Store.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: store meta.href must not be empty", operation)
	}
	if input.PriceType.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: priceType meta.href must not be empty", operation)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/retailstore", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/retailstore#3-izmenit-tochku-prodazh
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*RetailStore, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update retailstore"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/retailstore/"+id, input, opts)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Поиск — параметр Search; архивные точки сервер по умолчанию не отдаёт: см. filter archived.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/retailstore#3-poluchit-tochki-prodazh
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[RetailStore], error) {
	return list.Entity[RetailStore](ctx, s.client, "retailstore", options, nil, opts...)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*RetailStore, error) {
	var result RetailStore
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
