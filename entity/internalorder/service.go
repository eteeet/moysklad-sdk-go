package internalorder

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
		return nil, fmt.Errorf("moysklad: internalorder service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-poluchit-vnutrennij-zakaz
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*InternalOrder, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get internalorder: id must be a UUID")
	}
	return s.do(ctx, "get internalorder "+id, http.MethodGet, "entity/internalorder/"+id, nil, opts)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль документа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-sozdat-vnutrennij-zakaz
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*InternalOrder, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create internalorder"
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: organization meta.href must not be empty", operation)
	}
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: %s: name must not be empty", operation)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/internalorder", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-izmenit-vnutrennij-zakaz
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*InternalOrder, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update internalorder"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/internalorder/"+id, input, opts)
}

// GetPositions читает одну страницу строк без скрытого Get и автообхода; fields источник для
// внутреннего заказа не документирует — непустой page.Fields отклоняется до HTTP.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-poluchit-pozicii-vnutrennego-zakaza
func (s *Service) GetPositions(ctx context.Context, id string, page positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	return s.rows().GetPage(ctx, id, page, opts...)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-poluchit-vnutrennie-zakazy
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[InternalOrder], error) {
	return list.Entity[InternalOrder](ctx, s.client, "internalorder", options, nil, opts...)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*InternalOrder, error) {
	var result InternalOrder
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

// Delete — «Удалить Внутренний заказ»: DELETE entity/internalorder/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-udalit-vnutrennij-zakaz
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete internalorder: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete internalorder "+id, "entity/internalorder/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/internalorder: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Внутренних заказов».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-massovoe-sozdanie-i-obnovlenie-vnutrennih-zakazov
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[InternalOrder], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[InternalOrder](ctx, s.client, "internalorder", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/internalorder: элемент — meta по ID и поля Request.
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-massovoe-sozdanie-i-obnovlenie-vnutrennih-zakazov
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[InternalOrder], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[InternalOrder](ctx, s.client, "internalorder", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/internalorder/delete, не в корзину; до
// batch.MaxItems элементов. Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-massovoe-udalenie-vnutrennih-zakazov
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "internalorder", ids, opts)
}
