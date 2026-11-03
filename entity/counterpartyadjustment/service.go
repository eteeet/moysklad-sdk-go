package counterpartyadjustment

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
		return nil, fmt.Errorf("moysklad: counterpartyadjustment service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/counterpartyadjustment#3-poluchit-korrektirovku-vzaimoraschetov
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*CounterpartyAdjustment, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get counterpartyadjustment: id must be a UUID")
	}
	return s.do(ctx, "get counterpartyadjustment "+id, http.MethodGet, "entity/counterpartyadjustment/"+id, nil, opts)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/counterpartyadjustment#3-poluchit-spisok-korrektirovok-vzaimoraschetov
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[CounterpartyAdjustment], error) {
	return list.Entity[CounterpartyAdjustment](ctx, s.client, "counterpartyadjustment", options, nil, opts...)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль документа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/counterpartyadjustment#3-sozdat-korrektirovku-vzaimoraschetov
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*CounterpartyAdjustment, error) {
	transport.ResetResponseInfo(opts)
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create counterpartyadjustment: organization meta.href must not be empty")
	}
	if input.Agent.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create counterpartyadjustment: agent meta.href must not be empty")
	}
	return s.do(ctx, "create counterpartyadjustment", http.MethodPost, "entity/counterpartyadjustment", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/counterpartyadjustment#3-izmenit-korrektirovku-vzaimoraschetov
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*CounterpartyAdjustment, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update counterpartyadjustment: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update counterpartyadjustment: no fields specified")
	}
	return s.do(ctx, "update counterpartyadjustment "+id, http.MethodPut, "entity/counterpartyadjustment/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*CounterpartyAdjustment, error) {
	var result CounterpartyAdjustment
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

// Delete — «Удалить Корректировку взаиморасчетов»: DELETE entity/counterpartyadjustment/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/counterpartyadjustment#3-udalit-korrektirovku-vzaimoraschetov
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete counterpartyadjustment: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete counterpartyadjustment "+id, "entity/counterpartyadjustment/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/counterpartyadjustment: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Корректировок взаиморасчетов».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и
// неизвестный исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/counterpartyadjustment#3-massovoe-sozdanie-i-obnovlenie-korrektirovok-vzaimoraschetov
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[CounterpartyAdjustment], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[CounterpartyAdjustment](ctx, s.client, "counterpartyadjustment", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/counterpartyadjustment: элемент — meta по ID и поля Request.
// Источник: «Массовое создание и обновление Корректировок взаиморасчетов».
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/counterpartyadjustment#3-massovoe-sozdanie-i-obnovlenie-korrektirovok-vzaimoraschetov
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[CounterpartyAdjustment], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[CounterpartyAdjustment](ctx, s.client, "counterpartyadjustment", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/counterpartyadjustment/delete, не в корзину; до
// batch.MaxItems элементов. Источник: «Массовое удаление Корректировок взаиморасчетов».
// Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/counterpartyadjustment#3-massovoe-udalenie-korrektirovok-vzaimoraschetov
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "counterpartyadjustment", ids, opts)
}
