package saleschannel

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/batch"
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
		return nil, fmt.Errorf("moysklad: saleschannel service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/saleschannel#3-poluchit-kanal-prodazh
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*SalesChannel, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: saleschannel id must be a UUID")
	}
	return s.do(ctx, "get saleschannel "+id, http.MethodGet, "entity/saleschannel/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/saleschannel#3-sozdat-kanal-prodazh
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*SalesChannel, error) {
	transport.ResetResponseInfo(opts)
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create saleschannel: name must not be empty")
	}
	if input.Type == "" {
		return nil, fmt.Errorf("moysklad: create saleschannel: type must not be empty")
	}
	return s.do(ctx, "create saleschannel", http.MethodPost, "entity/saleschannel", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/saleschannel#3-izmenit-kanal-prodazh
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*SalesChannel, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update saleschannel: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update saleschannel: no fields specified")
	}
	return s.do(ctx, "update saleschannel "+id, http.MethodPut, "entity/saleschannel/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*SalesChannel, error) {
	var result SalesChannel
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
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/saleschannel#3-poluchit-kanaly-prodazh
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[SalesChannel], error) {
	return list.Entity[SalesChannel](ctx, s.client, "saleschannel", options, nil, opts...)
}

// Delete — DELETE entity/saleschannel/{id}, не в корзину. Повторов нет; права и связи проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/saleschannel#3-udalit-kanal-prodazh
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete saleschannel: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete saleschannel "+id, "entity/saleschannel/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/saleschannel (до batch.MaxItems); ответ — по правилам batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/saleschannel#3-massovoe-sozdanie-i-obnovlenie-kanalov-prodazh
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[SalesChannel], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[SalesChannel](ctx, s.client, "saleschannel", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/saleschannel: meta по ID и поля Request. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/saleschannel#3-massovoe-sozdanie-i-obnovlenie-kanalov-prodazh
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[SalesChannel], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[SalesChannel](ctx, s.client, "saleschannel", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/saleschannel/delete, не в корзину. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/saleschannel#3-massovoe-udalenie-kanalov-prodazh
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "saleschannel", ids, opts)
}
