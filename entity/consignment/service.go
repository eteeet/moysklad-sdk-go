package consignment

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
		return nil, fmt.Errorf("moysklad: consignment service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/consignment#3-poluchit-partiyu
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Consignment, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: consignment id must be a UUID")
	}
	return s.do(ctx, "get consignment "+id, http.MethodGet, "entity/consignment/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/consignment#3-sozdat-partiyu
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Consignment, error) {
	transport.ResetResponseInfo(opts)
	if input.Label == "" {
		return nil, fmt.Errorf("moysklad: create consignment: label must not be empty")
	}
	if input.Assortment.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create consignment: assortment meta.href must not be empty")
	}
	return s.do(ctx, "create consignment", http.MethodPost, "entity/consignment", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/consignment#3-izmenit-partiyu
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Consignment, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update consignment: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update consignment: no fields specified")
	}
	return s.do(ctx, "update consignment "+id, http.MethodPut, "entity/consignment/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Consignment, error) {
	var result Consignment
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
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/consignment#3-poluchit-spisok-partij
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Consignment], error) {
	return list.Entity[Consignment](ctx, s.client, "consignment", options, nil, opts...)
}

// Delete — DELETE entity/consignment/{id}, не в корзину. Повторов нет; права и связи проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/consignment#3-udalit-partiyu
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete consignment: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete consignment "+id, "entity/consignment/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/consignment (до batch.MaxItems); ответ — по правилам batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/consignment#3-massovoe-sozdanie-i-obnovlenie-partij
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[Consignment], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[Consignment](ctx, s.client, "consignment", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/consignment: meta по ID и поля Request. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/consignment#3-massovoe-sozdanie-i-obnovlenie-partij
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Consignment], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[Consignment](ctx, s.client, "consignment", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/consignment/delete, не в корзину. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/consignment#3-massovoe-udalenie-partij
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "consignment", ids, opts)
}
