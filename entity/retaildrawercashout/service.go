package retaildrawercashout

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
		return nil, fmt.Errorf("moysklad: retaildrawercashout service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildrawercashout#3-poluchit-vyplatu-deneg
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*RetailDrawerCashOut, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get retaildrawercashout: id must be a UUID")
	}
	return s.do(ctx, "get retaildrawercashout "+id, http.MethodGet, "entity/retaildrawercashout/"+id, nil, opts)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль документа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildrawercashout#3-sozdat-vyplatu-deneg
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*RetailDrawerCashOut, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create retaildrawercashout"
	for _, required := range []struct {
		name string
		ref  model.Reference
	}{{"organization", input.Organization}, {"agent", input.Agent}, {"retailShift", input.RetailShift}} {
		if required.ref.Meta.Href == "" {
			return nil, fmt.Errorf("moysklad: %s: %s meta.href must not be empty", operation, required.name)
		}
	}
	return s.do(ctx, operation, http.MethodPost, "entity/retaildrawercashout", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildrawercashout#3-izmenit-vyplatu-deneg
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*RetailDrawerCashOut, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update retaildrawercashout"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/retaildrawercashout/"+id, input, opts)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildrawercashout#3-poluchit-vyplaty-deneg
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[RetailDrawerCashOut], error) {
	return list.Entity[RetailDrawerCashOut](ctx, s.client, "retaildrawercashout", options, nil, opts...)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*RetailDrawerCashOut, error) {
	var result RetailDrawerCashOut
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

// Delete — «Удалить Выплату денег»: DELETE entity/retaildrawercashout/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildrawercashout#3-udalit-vyplatu-deneg
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete retaildrawercashout: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete retaildrawercashout "+id, "entity/retaildrawercashout/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/retaildrawercashout: массив до batch.MaxItems
// элементов. Источник: «Массовое создание и обновление Выплат денег».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildrawercashout#3-massovoe-sozdanie-i-obnovlenie-vyplat-deneg
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[RetailDrawerCashOut], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[RetailDrawerCashOut](ctx, s.client, "retaildrawercashout", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/retaildrawercashout: элемент — meta по ID и поля
// Request. Источник: «Массовое создание и обновление Выплат денег».
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildrawercashout#3-massovoe-sozdanie-i-obnovlenie-vyplat-deneg
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[RetailDrawerCashOut], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[RetailDrawerCashOut](ctx, s.client, "retaildrawercashout", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/retaildrawercashout/delete, не в корзину; до
// batch.MaxItems элементов. Источник: «Массовое удаление Выплат денег».
// Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildrawercashout#3-massovoe-udalenie-vyplat-deneg
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "retaildrawercashout", ids, opts)
}
