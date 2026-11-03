package payroll

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
		return nil, fmt.Errorf("moysklad: payroll service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payroll#3-poluchit-nachisleniya-zarplaty
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Payroll, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get payroll: id must be a UUID")
	}
	return s.do(ctx, "get payroll "+id, http.MethodGet, "entity/payroll/"+id, nil, opts)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль документа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payroll#3-sozdat-nachislenie-zarplaty
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Payroll, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create payroll"
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: organization meta.href must not be empty", operation)
	}
	if input.StartPayrollPeriod.IsZero() || input.EndPayrollPeriod.IsZero() {
		return nil, fmt.Errorf("moysklad: %s: startPayrollPeriod and endPayrollPeriod must be set", operation)
	}
	if rows, ok := input.Positions.Get(); ok && len(rows) > 1000 {
		return nil, fmt.Errorf("moysklad: %s: positions must not exceed 1000", operation)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/payroll", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payroll#3-izmenit-nachislenie-zarplaty
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Payroll, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update payroll"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	if rows, ok := input.Positions.Get(); ok && len(rows) > 1000 {
		return nil, fmt.Errorf("moysklad: %s: positions must not exceed 1000", operation)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/payroll/"+id, input, opts)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payroll#3-poluchit-spisok-nachisleniya-zarplat
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Payroll], error) {
	return list.Entity[Payroll](ctx, s.client, "payroll", options, nil, opts...)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Payroll, error) {
	var result Payroll
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

// Delete — «Удалить Начисления зарплаты»: DELETE entity/payroll/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payroll#3-udalit-nachisleniya-zarplaty
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete payroll: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete payroll "+id, "entity/payroll/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/payroll: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Начисления зарплат».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payroll#3-massovoe-sozdanie-i-obnovlenie-nachisleniya-zarplat
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[Payroll], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[Payroll](ctx, s.client, "payroll", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/payroll: элемент — meta по ID и поля Request.
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payroll#3-massovoe-sozdanie-i-obnovlenie-nachisleniya-zarplat
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Payroll], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[Payroll](ctx, s.client, "payroll", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/payroll/delete, не в корзину; до
// batch.MaxItems элементов. Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payroll#3-massovoe-udalenie-nachisleniya-zarplat
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "payroll", ids, opts)
}
