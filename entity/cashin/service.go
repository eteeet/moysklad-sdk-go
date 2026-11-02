package cashin

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
		return nil, fmt.Errorf("moysklad: cashin service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashin
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*CashIn, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: cashin id must be a UUID")
	}
	return s.do(ctx, "get cashin "+id, http.MethodGet, "entity/cashin/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashin
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль CashIn.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*CashIn, error) {
	transport.ResetResponseInfo(opts)
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
	transport.ResetResponseInfo(opts)
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

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// Позиции в строках — только то, что вернул сервер; строки позиций читает GetPositions.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashin
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[CashIn], error) {
	return list.Entity[CashIn](ctx, s.client, "cashin", options, nil, opts...)
}

// Delete — «Удалить Приходный ордер»: DELETE entity/cashin/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashin#3-udalit-prihodnyj-order
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete cashin: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete cashin "+id, "entity/cashin/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/cashin: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Приходных ордеров».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashin#3-massovoe-sozdanie-i-obnovlenie-prihodnyh-orderov
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[CashIn], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[CashIn](ctx, s.client, "cashin", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/cashin: элемент — meta по ID и поля Request.
// Источник: «Массовое создание и обновление Приходных ордеров».
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashin#3-massovoe-sozdanie-i-obnovlenie-prihodnyh-orderov
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[CashIn], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[CashIn](ctx, s.client, "cashin", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/cashin/delete, не в корзину; до
// batch.MaxItems элементов. Источник: «Массовое удаление Приходных ордеров».
// Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashin#3-massovoe-udalenie-prihodnyh-orderov
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "cashin", ids, opts)
}
