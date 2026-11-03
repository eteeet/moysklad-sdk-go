package accumulationdiscount

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: accumulationdiscount service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-poluchit-nakopitelnuyu-skidku
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*AccumulationDiscount, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: accumulationdiscount id must be a UUID")
	}
	return s.do(ctx, "get accumulationdiscount "+id, http.MethodGet, "entity/accumulationdiscount/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-sozdat-nakopitelnuyu-skidku
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*AccumulationDiscount, error) {
	transport.ResetResponseInfo(opts)
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create accumulationdiscount: name must not be empty")
	}
	return s.do(ctx, "create accumulationdiscount", http.MethodPost, "entity/accumulationdiscount", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-izmenit-nakopitelnuyu-skidku
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*AccumulationDiscount, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update accumulationdiscount: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update accumulationdiscount: no fields specified")
	}
	return s.do(ctx, "update accumulationdiscount "+id, http.MethodPut, "entity/accumulationdiscount/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*AccumulationDiscount, error) {
	var result AccumulationDiscount
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

// Delete — DELETE entity/accumulationdiscount/{id}, не в корзину. Повторов нет; права и связи проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-udalit-nakopitelnuyu-skidku
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete accumulationdiscount: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete accumulationdiscount "+id, "entity/accumulationdiscount/"+id, opts)
}
