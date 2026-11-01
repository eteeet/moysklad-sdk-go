package paymentin

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

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
func (s *Service) Get(ctx context.Context, id string) (*PaymentIn, error) {
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: paymentin id must be a UUID")
	}
	return s.do(ctx, "get paymentin "+id, http.MethodGet, "entity/paymentin/"+id, nil)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payment-in
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль PaymentIn.
func (s *Service) Create(ctx context.Context, input CreateRequest) (*PaymentIn, error) {
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create paymentin: organization.meta.href must not be empty")
	}
	if input.Agent.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create paymentin: agent.meta.href must not be empty")
	}
	return s.do(ctx, "create paymentin", http.MethodPost, "entity/paymentin", input)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payment-in
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest) (*PaymentIn, error) {
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update paymentin: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update paymentin: no fields specified")
	}
	return s.do(ctx, "update paymentin "+id, http.MethodPut, "entity/paymentin/"+id, input)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any) (*PaymentIn, error) {
	req, err := s.client.NewRequest(ctx, method, endpoint, nil, body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	var result PaymentIn
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
