package inventory

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
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: inventory service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-poluchit-inventarizaciyu
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Inventory, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get inventory: id must be a UUID")
	}
	return s.do(ctx, "get inventory "+id, http.MethodGet, "entity/inventory/"+id, nil, opts)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль документа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-sozdat-inventarizaciyu
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Inventory, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create inventory"
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: organization meta.href must not be empty", operation)
	}
	if input.Store.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: store meta.href must not be empty", operation)
	}
	if rows, ok := input.Positions.Get(); ok && len(rows) > 1000 {
		return nil, fmt.Errorf("moysklad: %s: positions must not exceed 1000", operation)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/inventory", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-izmenit-inventarizaciyu
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Inventory, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update inventory"
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
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/inventory/"+id, input, opts)
}

// GetPositions читает одну страницу строк без скрытого Get и автообхода; fields источник для
// инвентаризации не документирует — непустой page.Fields отклоняется до HTTP.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-poluchit-pozicii
func (s *Service) GetPositions(ctx context.Context, id string, page positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	return s.rows().GetPage(ctx, id, page, opts...)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-poluchit-inventarizacii
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Inventory], error) {
	return list.Entity[Inventory](ctx, s.client, "inventory", options, nil, opts...)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Inventory, error) {
	var result Inventory
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

// Delete — «Удалить Инвентаризацию»: DELETE entity/inventory/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-udalit-inventarizaciyu
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete inventory: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete inventory "+id, "entity/inventory/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/inventory: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Инвентаризаций».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-massovoe-sozdanie-i-obnovlenie-inventarizacij
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[Inventory], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[Inventory](ctx, s.client, "inventory", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/inventory: элемент — meta по ID и поля Request.
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-massovoe-sozdanie-i-obnovlenie-inventarizacij
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Inventory], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[Inventory](ctx, s.client, "inventory", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/inventory/delete, не в корзину; до
// batch.MaxItems элементов. Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-massovoe-udalenie-inventarizacij
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "inventory", ids, opts)
}
