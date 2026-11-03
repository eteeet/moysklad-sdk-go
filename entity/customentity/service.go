package customentity

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: customentity service requires a client")
	}
	return &Service{client: client}, nil
}

// Create — POST entity/customentity. Не повторять после неоднозначного сбоя: возможен дубль.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/customentity#3-sozdat-spravochnik
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*CustomEntity, error) {
	transport.ResetResponseInfo(opts)
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create customentity: name must not be empty")
	}
	var result CustomEntity
	if err := s.do(ctx, "create customentity", http.MethodPost, "entity/customentity", input, &result, &result.ID, opts...); err != nil {
		return nil, err
	}
	return &result, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/customentity#3-izmenit-spravochnik
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*CustomEntity, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update customentity: id must be a UUID")
	}
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update customentity: no fields specified")
	}
	var result CustomEntity
	if err := s.do(ctx, "update customentity "+id, http.MethodPut, "entity/customentity/"+id, input, &result, &result.ID, opts...); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete — DELETE entity/customentity/{id}: справочник удаляется вместе с элементами. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/customentity#3-udalit-spravochnik
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete customentity: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete customentity "+id, "entity/customentity/"+id, opts)
}

// ListElements — GET entity/customentity/{id}: одна страница элементов без expand и автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/customentity#3-poluchit-elementy-spravochnika
func (s *Service) ListElements(ctx context.Context, id string, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Element], error) {
	transport.ResetResponseInfo(opts)
	operation := "list customentity " + id + " elements"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: customentity id must be a UUID", operation)
	}
	if options.Fields != "" {
		return nil, fmt.Errorf("moysklad: %s: fields=%q is not supported", operation, options.Fields)
	}
	params, err := options.Values(query.MaxLimit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[Element](ctx, s.client, operation, "entity/customentity/"+id, params, opts...)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/customentity#3-poluchit-element-spravochnika
func (s *Service) GetElement(ctx context.Context, id, elementID string, opts ...transport.RequestOption) (*Element, error) {
	transport.ResetResponseInfo(opts)
	operation := "get customentity " + id + " element " + elementID
	if !validID(id) || !validID(elementID) {
		return nil, fmt.Errorf("moysklad: %s: customentity and element ids must be UUIDs", operation)
	}
	var result Element
	if err := s.do(ctx, operation, http.MethodGet, "entity/customentity/"+id+"/"+elementID, nil, &result, &result.ID, opts...); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateElement — POST entity/customentity/{id}. Не повторять после неоднозначного сбоя.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/customentity#3-sozdat-element-spravochnika
func (s *Service) CreateElement(ctx context.Context, id string, input ElementCreateRequest, opts ...transport.RequestOption) (*Element, error) {
	transport.ResetResponseInfo(opts)
	operation := "create customentity " + id + " element"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: customentity id must be a UUID", operation)
	}
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: %s: name must not be empty", operation)
	}
	var result Element
	if err := s.do(ctx, operation, http.MethodPost, "entity/customentity/"+id, input, &result, &result.ID, opts...); err != nil {
		return nil, err
	}
	return &result, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/customentity#3-izmenit-element-spravochnika
func (s *Service) UpdateElement(ctx context.Context, id, elementID string, input ElementUpdateRequest, opts ...transport.RequestOption) (*Element, error) {
	transport.ResetResponseInfo(opts)
	operation := "update customentity " + id + " element " + elementID
	if !validID(id) || !validID(elementID) {
		return nil, fmt.Errorf("moysklad: %s: customentity and element ids must be UUIDs", operation)
	}
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	var result Element
	if err := s.do(ctx, operation, http.MethodPut, "entity/customentity/"+id+"/"+elementID, input, &result, &result.ID, opts...); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteElement — DELETE элемента. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/customentity#3-udalit-element-spravochnika
func (s *Service) DeleteElement(ctx context.Context, id, elementID string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	operation := "delete customentity " + id + " element " + elementID
	if !validID(id) || !validID(elementID) {
		return fmt.Errorf("moysklad: %s: customentity and element ids must be UUIDs", operation)
	}
	return write.Delete(ctx, s.client, operation, "entity/customentity/"+id+"/"+elementID, opts)
}

// do — один запрос; успех — HTTP 200 и непустой id ответа.
func (s *Service) do(ctx context.Context, operation, method, endpoint string, body, dst any, id *string, opts ...transport.RequestOption) error {
	resp, err := s.client.Send(ctx, method, endpoint, nil, body, dst, opts...)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if *id == "" {
		return fmt.Errorf("moysklad: %s: response has no id", operation)
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
