package salesreturn

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: salesreturn service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/sales-return#3-poluchit-vozvrat-pokupatelya
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*SalesReturn, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get salesreturn: id must be a UUID")
	}
	return s.do(ctx, "get salesreturn "+id, http.MethodGet, "entity/salesreturn/"+id, nil, opts)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль документа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/sales-return#3-sozdat-vozvrat-pokupatelya
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*SalesReturn, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create salesreturn"
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: organization meta.href must not be empty", operation)
	}
	if input.Store.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: store meta.href must not be empty", operation)
	}
	if input.Agent.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: agent meta.href must not be empty", operation)
	}
	if rows, ok := input.Positions.Get(); ok && len(rows) > 1000 {
		return nil, fmt.Errorf("moysklad: %s: positions must not exceed 1000", operation)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/salesreturn", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/sales-return#3-izmenit-vozvrat-pokupatelya
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*SalesReturn, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update salesreturn"
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
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/salesreturn/"+id, input, opts)
}

// GetPositions читает одну страницу строк без скрытого Get и автообхода; page.Fields="stock" —
// остатки и себестоимость на момент документа (cost нет у возврата без основания). Get, List
// и Update отдают stock только вместе с expand=positions, который в SDK не входит.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/sales-return#3-poluchit-pozicii
func (s *Service) GetPositions(ctx context.Context, id string, page positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	transport.ResetResponseInfo(opts)
	reader, err := positions.NewService[Position](s.client, "salesreturn")
	if err != nil {
		return nil, err
	}
	return reader.GetPage(ctx, id, page, opts...)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/sales-return#3-poluchit-vozvraty-pokupatelej
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[SalesReturn], error) {
	return list.Entity[SalesReturn](ctx, s.client, "salesreturn", options, nil, opts...)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*SalesReturn, error) {
	var result SalesReturn
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

// Delete — «Удалить Возврат покупателя»: DELETE entity/salesreturn/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/sales-return#3-udalit-vozvrat-pokupatelya
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete salesreturn: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete salesreturn "+id, "entity/salesreturn/"+id, opts)
}
