package contract

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
		return nil, fmt.Errorf("moysklad: contract service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contract#3-poluchit-dogovor
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Contract, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: contract id must be a UUID")
	}
	return s.do(ctx, "get contract "+id, http.MethodGet, "entity/contract/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contract#3-sozdat-novyj-dogovor
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Contract, error) {
	transport.ResetResponseInfo(opts)
	if input.OwnAgent.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create contract: ownAgent meta.href must not be empty")
	}
	if input.Agent.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create contract: agent meta.href must not be empty")
	}
	return s.do(ctx, "create contract", http.MethodPost, "entity/contract", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contract#3-izmenit-dogovor
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Contract, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update contract: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update contract: no fields specified")
	}
	return s.do(ctx, "update contract "+id, http.MethodPut, "entity/contract/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Contract, error) {
	var result Contract
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
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contract#3-poluchit-spisok-dogovorov
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Contract], error) {
	return list.Entity[Contract](ctx, s.client, "contract", options, nil, opts...)
}

// Delete — DELETE entity/contract/{id}, не в корзину. Повторов нет; права и связи проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contract#3-udalit-dogovor
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete contract: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete contract "+id, "entity/contract/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/contract (до batch.MaxItems); ответ — по правилам batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contract#3-massovoe-sozdanie-i-obnovlenie-dogovorov
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[Contract], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[Contract](ctx, s.client, "contract", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/contract: meta по ID и поля Request. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contract#3-massovoe-sozdanie-i-obnovlenie-dogovorov
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Contract], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[Contract](ctx, s.client, "contract", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/contract/delete, не в корзину. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contract#3-massovoe-udalenie-dogovorov
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "contract", ids, opts)
}
