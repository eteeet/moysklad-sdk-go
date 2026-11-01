package bundle

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
		return nil, fmt.Errorf("moysklad: bundle service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle#3-poluchit-komplekt
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Bundle, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: bundle id must be a UUID")
	}
	return s.do(ctx, "get bundle "+id, http.MethodGet, "entity/bundle/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle#3-sozdat-komplekt
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль Bundle.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Bundle, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create bundle"
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: %s: name must not be empty", operation)
	}
	if len(input.Components) == 0 {
		return nil, fmt.Errorf("moysklad: %s: components must not be empty", operation)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/bundle", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle#3-izmenit-komplekt
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Bundle, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update bundle"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/bundle/"+id, input, opts)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Фильтр штрихкода — поле barcode в единственном числе.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle#3-poluchit-spisok-komplektov
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Bundle], error) {
	return list.Entity[Bundle](ctx, s.client, "bundle", options, nil, opts...)
}

// GetComponents читает одну страницу компонентов (источник объявляет только limit и offset)
// без скрытого Get комплекта и автообхода. Добавление, изменение и удаление отдельных
// компонентов в SDK v1 не входят: состав меняется через Update.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle#3-poluchit-komponenty-komplekta
func (s *Service) GetComponents(ctx context.Context, id string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Component], error) {
	transport.ResetResponseInfo(opts)
	operation := "get bundle " + id + " components"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: bundle id must be a UUID", operation)
	}
	if page.Fields != "" {
		return nil, fmt.Errorf("moysklad: %s: fields=%q is not supported", operation, page.Fields)
	}
	params, err := page.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[Component](ctx, s.client, operation, "entity/bundle/"+id+"/components", params, opts...)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bundle#3-poluchit-komponent
func (s *Service) GetComponent(ctx context.Context, id, componentID string, opts ...transport.RequestOption) (*Component, error) {
	transport.ResetResponseInfo(opts)
	operation := "get bundle " + id + " component " + componentID
	if !validID(id) || !validID(componentID) {
		return nil, fmt.Errorf("moysklad: %s: bundle and component ids must be UUIDs", operation)
	}
	var result Component
	if err := s.send(ctx, operation, http.MethodGet, "entity/bundle/"+id+"/components/"+componentID, nil, &result, opts); err != nil {
		return nil, err
	}
	if result.ID == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
	}
	return &result, nil
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Bundle, error) {
	var result Bundle
	if err := s.send(ctx, operation, method, endpoint, body, &result, opts); err != nil {
		return nil, err
	}
	if result.ID == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
	}
	return &result, nil
}

func (s *Service) send(ctx context.Context, operation, method, endpoint string, body, dst any, opts []transport.RequestOption) error {
	resp, err := s.client.Send(ctx, method, endpoint, nil, body, dst, opts...)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	return nil
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
