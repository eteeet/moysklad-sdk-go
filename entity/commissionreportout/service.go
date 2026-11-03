package commissionreportout

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
		return nil, fmt.Errorf("moysklad: commissionreportout service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportout#3-poluchit-vydannyj-otchet-komissionera
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*CommissionReportOut, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get commissionreportout: id must be a UUID")
	}
	return s.do(ctx, "get commissionreportout "+id, http.MethodGet, "entity/commissionreportout/"+id, nil, opts)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль документа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportout#3-sozdat-vydannyj-otchet-komissionera
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*CommissionReportOut, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create commissionreportout"
	for _, ref := range []struct {
		name string
		ref  model.Reference
	}{{"organization", input.Organization}, {"agent", input.Agent}, {"contract", input.Contract}} {
		if ref.ref.Meta.Href == "" {
			return nil, fmt.Errorf("moysklad: %s: %s meta.href must not be empty", operation, ref.name)
		}
	}
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: %s: name must not be empty", operation)
	}
	if input.CommissionPeriodStart.IsZero() || input.CommissionPeriodEnd.IsZero() {
		return nil, fmt.Errorf("moysklad: %s: commissionPeriodStart and commissionPeriodEnd must be set", operation)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/commissionreportout", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportout#3-izmenit-vydannyj-otchet-komissionera
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*CommissionReportOut, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update commissionreportout"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/commissionreportout/"+id, input, opts)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportout#3-poluchit-vydannye-otchety-komissionera
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[CommissionReportOut], error) {
	return list.Entity[CommissionReportOut](ctx, s.client, "commissionreportout", options, nil, opts...)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*CommissionReportOut, error) {
	var result CommissionReportOut
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

// Delete — «Удалить Выданный отчет комиссионера»: DELETE entity/commissionreportout/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportout#3-udalit-vydannyj-otchet-komissionera
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete commissionreportout: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete commissionreportout "+id, "entity/commissionreportout/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/commissionreportout: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Выданных отчетов комиссионера».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportout#3-massovoe-sozdanie-i-obnovlenie-vydannyh-otchetov-komissionera
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[CommissionReportOut], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[CommissionReportOut](ctx, s.client, "commissionreportout", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/commissionreportout: элемент — meta по ID и поля Request.
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportout#3-massovoe-sozdanie-i-obnovlenie-vydannyh-otchetov-komissionera
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[CommissionReportOut], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[CommissionReportOut](ctx, s.client, "commissionreportout", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/commissionreportout/delete, не в корзину; до
// batch.MaxItems элементов. Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportout#3-massovoe-udalenie-vydannyh-otchetov-komissionera
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "commissionreportout", ids, opts)
}
