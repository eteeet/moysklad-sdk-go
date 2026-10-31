package counterparty

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
		return nil, fmt.Errorf("moysklad: counterparty service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-poluchit-kontragenta
func (s *Service) Get(ctx context.Context, id string) (*Counterparty, error) {
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: counterparty id must be a UUID")
	}
	return s.do(ctx, "get counterparty "+id, http.MethodGet, "entity/counterparty/"+id, nil)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-sozdat-kontragenta
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль Counterparty.
func (s *Service) Create(ctx context.Context, input CreateRequest) (*Counterparty, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create counterparty: name must not be empty")
	}
	return s.do(ctx, "create counterparty", http.MethodPost, "entity/counterparty", input)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/counterparty#3-izmenit-kontragenta
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest) (*Counterparty, error) {
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update counterparty: id must be a UUID")
	}
	// Нулевые Field — «не передано»; false, пустые значения и Null считаются заданными.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update counterparty: no fields specified")
	}
	return s.do(ctx, "update counterparty "+id, http.MethodPut, "entity/counterparty/"+id, input)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any) (*Counterparty, error) {
	req, err := s.client.NewRequest(ctx, method, endpoint, nil, body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	var result Counterparty
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
