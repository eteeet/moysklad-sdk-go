package paymentin

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
		return nil, fmt.Errorf("moysklad: paymentin service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payment-in
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*PaymentIn, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: paymentin id must be a UUID")
	}
	return s.do(ctx, "get paymentin "+id, http.MethodGet, "entity/paymentin/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payment-in
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль PaymentIn.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*PaymentIn, error) {
	transport.ResetResponseInfo(opts)
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create paymentin: organization.meta.href must not be empty")
	}
	if input.Agent.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create paymentin: agent.meta.href must not be empty")
	}
	return s.do(ctx, "create paymentin", http.MethodPost, "entity/paymentin", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payment-in
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*PaymentIn, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update paymentin: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update paymentin: no fields specified")
	}
	return s.do(ctx, "update paymentin "+id, http.MethodPut, "entity/paymentin/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*PaymentIn, error) {
	var result PaymentIn
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
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// Позиции в строках — только то, что вернул сервер; строки позиций читает GetPositions.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/paymentin
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[PaymentIn], error) {
	return list.Entity[PaymentIn](ctx, s.client, "paymentin", options, nil, opts...)
}

// Delete — «Удалить Входящий платеж»: DELETE entity/paymentin/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payment-in#3-udalit-vhodyashij-platezh
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete paymentin: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete paymentin "+id, "entity/paymentin/"+id, opts)
}
