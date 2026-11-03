package commissionreportin

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
		return nil, fmt.Errorf("moysklad: commissionreportin service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-poluchit-otchet-komissionera
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*CommissionReportIn, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get commissionreportin: id must be a UUID")
	}
	return s.do(ctx, "get commissionreportin "+id, http.MethodGet, "entity/commissionreportin/"+id, nil, opts)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль документа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-sozdat-poluchennyj-otchet-komissionera
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*CommissionReportIn, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create commissionreportin"
	for _, ref := range []struct {
		name string
		ref  model.Reference
	}{{"organization", input.Organization}, {"agent", input.Agent}, {"contract", input.Contract}} {
		if ref.ref.Meta.Href == "" {
			return nil, fmt.Errorf("moysklad: %s: %s meta.href must not be empty", operation, ref.name)
		}
	}
	if input.CommissionPeriodStart.IsZero() || input.CommissionPeriodEnd.IsZero() {
		return nil, fmt.Errorf("moysklad: %s: commissionPeriodStart and commissionPeriodEnd must be set", operation)
	}
	if err := checkRows(input.Positions, input.ReturnToCommissionerPositions); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/commissionreportin", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-izmenit-poluchennyj-otchet-komissionera
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*CommissionReportIn, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update commissionreportin"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	if err := checkRows(input.Positions, input.ReturnToCommissionerPositions); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/commissionreportin/"+id, input, opts)
}

// Лимит строк, сохраняемых вместе с документом, — 1000 на каждый набор.
func checkRows(sold model.Field[[]PositionRequest], returned model.Field[[]ReturnPositionRequest]) error {
	if rows, ok := sold.Get(); ok && len(rows) > 1000 {
		return fmt.Errorf("positions must not exceed 1000")
	}
	if rows, ok := returned.Get(); ok && len(rows) > 1000 {
		return fmt.Errorf("returnToCommissionerPositions must not exceed 1000")
	}
	return nil
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-poluchit-otchety-komissionera
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[CommissionReportIn], error) {
	return list.Entity[CommissionReportIn](ctx, s.client, "commissionreportin", options, nil, opts...)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*CommissionReportIn, error) {
	var result CommissionReportIn
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

// Delete — «Удалить Полученный отчет комиссионера»: DELETE entity/commissionreportin/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-udalit-poluchennyj-otchet-komissionera
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete commissionreportin: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete commissionreportin "+id, "entity/commissionreportin/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/commissionreportin: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Полученных отчетов комиссионера».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-massovoe-sozdanie-i-obnovlenie-poluchennyh-otchetov-komissionera
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[CommissionReportIn], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[CommissionReportIn](ctx, s.client, "commissionreportin", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/commissionreportin: элемент — meta по ID и поля Request.
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-massovoe-sozdanie-i-obnovlenie-poluchennyh-otchetov-komissionera
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[CommissionReportIn], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[CommissionReportIn](ctx, s.client, "commissionreportin", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/commissionreportin/delete, не в корзину; до
// batch.MaxItems элементов. Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin#3-massovoe-udalenie-poluchennyh-otchetov-komissionera
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "commissionreportin", ids, opts)
}
