package processing

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
		return nil, fmt.Errorf("moysklad: processing service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-poluchit-tehoperaciyu
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Processing, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get processing: id must be a UUID")
	}
	return s.do(ctx, "get processing "+id, http.MethodGet, "entity/processing/"+id, nil, opts)
}

// Create создаёт техоперацию с техкартой или без неё — правила в CreateRequest. Не повторять POST
// автоматически после неоднозначного сбоя: возможен дубль документа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-sozdat-tehoperaciyu
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Processing, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create processing"
	for _, ref := range []struct {
		name string
		ref  model.Reference
	}{{"organization", input.Organization}, {"materialsStore", input.MaterialsStore}, {"productsStore", input.ProductsStore}} {
		if ref.ref.Meta.Href == "" {
			return nil, fmt.Errorf("moysklad: %s: %s meta.href must not be empty", operation, ref.name)
		}
	}
	if err := checkRows(input.Materials, input.Products); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/processing", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-izmenit-tehoperaciyu
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Processing, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update processing"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	if err := checkRows(input.Materials, input.Products); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/processing/"+id, input, opts)
}

// Лимит строк, сохраняемых вместе с документом, — 1000 на каждый набор.
func checkRows(materials, products model.Field[[]PositionRequest]) error {
	if rows, ok := materials.Get(); ok && len(rows) > 1000 {
		return fmt.Errorf("materials must not exceed 1000")
	}
	if rows, ok := products.Get(); ok && len(rows) > 1000 {
		return fmt.Errorf("products must not exceed 1000")
	}
	return nil
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-poluchit-spisok-tehoperacij
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Processing], error) {
	return list.Entity[Processing](ctx, s.client, "processing", options, nil, opts...)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Processing, error) {
	var result Processing
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

// Delete — «Удалить Техоперацию»: DELETE entity/processing/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-udalit-tehoperaciyu
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete processing: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete processing "+id, "entity/processing/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/processing: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Техопераций».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-massovoe-sozdanie-i-obnovlenie-tehoperacij
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[Processing], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[Processing](ctx, s.client, "processing", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/processing: элемент — meta по ID и поля Request.
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-massovoe-sozdanie-i-obnovlenie-tehoperacij
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Processing], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[Processing](ctx, s.client, "processing", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/processing/delete, не в корзину; до
// batch.MaxItems элементов. Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-massovoe-udalenie-tehoperacij
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "processing", ids, opts)
}
