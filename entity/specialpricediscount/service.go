package specialpricediscount

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
		return nil, fmt.Errorf("moysklad: specialpricediscount service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-poluchit-specialnuyu-cenu
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*SpecialPriceDiscount, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: specialpricediscount id must be a UUID")
	}
	return s.do(ctx, "get specialpricediscount "+id, http.MethodGet, "entity/specialpricediscount/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-sozdat-specialnuyu-cenu
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*SpecialPriceDiscount, error) {
	transport.ResetResponseInfo(opts)
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create specialpricediscount: name must not be empty")
	}
	return s.do(ctx, "create specialpricediscount", http.MethodPost, "entity/specialpricediscount", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-izmenit-specialnuyu-cenu
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*SpecialPriceDiscount, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update specialpricediscount: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update specialpricediscount: no fields specified")
	}
	return s.do(ctx, "update specialpricediscount "+id, http.MethodPut, "entity/specialpricediscount/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*SpecialPriceDiscount, error) {
	var result SpecialPriceDiscount
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

// Delete — DELETE entity/specialpricediscount/{id}, не в корзину. Повторов нет; права и связи проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-udalit-specialnuyu-cenu
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete specialpricediscount: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete specialpricediscount "+id, "entity/specialpricediscount/"+id, opts)
}
