package contentcard

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
		return nil, fmt.Errorf("moysklad: contentcard service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contentcard#3-poluchit-kartochku-kontenta
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*ContentCard, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: contentcard id must be a UUID")
	}
	return s.do(ctx, "get contentcard "+id, http.MethodGet, "entity/contentcard/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contentcard#3-sozdat-kartochku-kontenta
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*ContentCard, error) {
	transport.ResetResponseInfo(opts)
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create contentcard: name must not be empty")
	}
	if input.Description == "" {
		return nil, fmt.Errorf("moysklad: create contentcard: description must not be empty")
	}
	if input.CardContentName == "" {
		return nil, fmt.Errorf("moysklad: create contentcard: cardContentName must not be empty")
	}
	if input.Assortment.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create contentcard: assortment meta.href must not be empty")
	}
	if input.SalePlatform.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create contentcard: salePlatform meta.href must not be empty")
	}
	if len(input.SalesChannels) == 0 {
		return nil, fmt.Errorf("moysklad: create contentcard: salesChannels must not be empty")
	}
	return s.do(ctx, "create contentcard", http.MethodPost, "entity/contentcard", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contentcard#3-izmenit-kartochku-kontenta
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*ContentCard, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update contentcard: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update contentcard: no fields specified")
	}
	return s.do(ctx, "update contentcard "+id, http.MethodPut, "entity/contentcard/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*ContentCard, error) {
	var result ContentCard
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
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contentcard#3-poluchit-kartochki-kontenta
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[ContentCard], error) {
	return list.Entity[ContentCard](ctx, s.client, "contentcard", options, nil, opts...)
}

// Delete — DELETE entity/contentcard/{id}, не в корзину. Повторов нет; права и связи проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contentcard#3-udalit-kartochku-kontenta
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete contentcard: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete contentcard "+id, "entity/contentcard/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/contentcard (до batch.MaxItems); ответ — по правилам batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contentcard#3-massovoe-sozdanie-i-obnovlenie-kartochek-kontenta
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[ContentCard], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[ContentCard](ctx, s.client, "contentcard", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/contentcard: meta по ID и поля Request. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contentcard#3-massovoe-sozdanie-i-obnovlenie-kartochek-kontenta
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[ContentCard], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[ContentCard](ctx, s.client, "contentcard", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/contentcard/delete, не в корзину. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contentcard#3-massovoe-udalenie-kartochek-kontenta
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "contentcard", ids, opts)
}
