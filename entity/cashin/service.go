package cashin

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
		return nil, fmt.Errorf("moysklad: cashin service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashin
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*CashIn, error) {
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: cashin id must be a UUID")
	}
	return s.do(ctx, "get cashin "+id, http.MethodGet, "entity/cashin/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashin
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль CashIn.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*CashIn, error) {
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create cashin: organization.meta.href must not be empty")
	}
	if input.Agent.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create cashin: agent.meta.href must not be empty")
	}
	return s.do(ctx, "create cashin", http.MethodPost, "entity/cashin", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashin
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*CashIn, error) {
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update cashin: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update cashin: no fields specified")
	}
	return s.do(ctx, "update cashin "+id, http.MethodPut, "entity/cashin/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*CashIn, error) {
	var result CashIn
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
