package paymentout

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
		return nil, fmt.Errorf("moysklad: paymentout service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payment-out
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*PaymentOut, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: paymentout id must be a UUID")
	}
	return s.do(ctx, "get paymentout "+id, http.MethodGet, "entity/paymentout/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payment-out
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль PaymentOut.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*PaymentOut, error) {
	transport.ResetResponseInfo(opts)
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create paymentout: organization.meta.href must not be empty")
	}
	if input.Agent.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create paymentout: agent.meta.href must not be empty")
	}
	if input.ExpenseItem.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create paymentout: expenseItem.meta.href must not be empty")
	}
	return s.do(ctx, "create paymentout", http.MethodPost, "entity/paymentout", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payment-out
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*PaymentOut, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update paymentout: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update paymentout: no fields specified")
	}
	return s.do(ctx, "update paymentout "+id, http.MethodPut, "entity/paymentout/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*PaymentOut, error) {
	var result PaymentOut
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
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/paymentout
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[PaymentOut], error) {
	return list.Entity[PaymentOut](ctx, s.client, "paymentout", options, nil, opts...)
}

// Delete — «Удалить Исходящий платеж»: DELETE entity/paymentout/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payment-out#3-udalit-ishodyashij-platezh
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete paymentout: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete paymentout "+id, "entity/paymentout/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/paymentout: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Исходящих платежей».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payment-out#3-massovoe-sozdanie-i-obnovlenie-ishodyashih-platezhej
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[PaymentOut], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[PaymentOut](ctx, s.client, "paymentout", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/paymentout: элемент — meta по ID и поля Request.
// Источник: «Массовое создание и обновление Исходящих платежей».
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payment-out#3-massovoe-sozdanie-i-obnovlenie-ishodyashih-platezhej
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[PaymentOut], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[PaymentOut](ctx, s.client, "paymentout", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/paymentout/delete, не в корзину; до
// batch.MaxItems элементов. Источник: «Массовое удаление Исходящих платежей».
// Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payment-out#3-massovoe-udalenie-ishodyashih-platezhej
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "paymentout", ids, opts)
}
