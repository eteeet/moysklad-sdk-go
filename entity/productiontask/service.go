package productiontask

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
		return nil, fmt.Errorf("moysklad: productiontask service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-poluchit-proizvodstvennoe-zadanie
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*ProductionTask, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get productiontask: id must be a UUID")
	}
	return s.do(ctx, "get productiontask "+id, http.MethodGet, "entity/productiontask/"+id, nil, opts)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль документа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-sozdat-proizvodstvennoe-zadanie
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*ProductionTask, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create productiontask"
	for _, ref := range []struct {
		name string
		ref  model.Reference
	}{{"organization", input.Organization}, {"materialsStore", input.MaterialsStore}, {"productsStore", input.ProductsStore}} {
		if ref.ref.Meta.Href == "" {
			return nil, fmt.Errorf("moysklad: %s: %s meta.href must not be empty", operation, ref.name)
		}
	}
	if err := checkRows(input.ProductionRows, input.Products); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/productiontask", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-izmenit-proizvodstvennoe-zadanie
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*ProductionTask, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update productiontask"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	if err := checkRows(input.ProductionRows, input.Products); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/productiontask/"+id, input, opts)
}

// Источник: «Максимально допустимое число позиций у производственного задания — 200»; продуктов
// вместе с документом сохраняется не больше 1000.
func checkRows(rows model.Field[[]ProductionRowRequest], products model.Field[[]ProductRequest]) error {
	if items, ok := rows.Get(); ok && len(items) > 200 {
		return fmt.Errorf("productionRows must not exceed 200")
	}
	if items, ok := products.Get(); ok && len(items) > 1000 {
		return fmt.Errorf("products must not exceed 1000")
	}
	return nil
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-poluchit-spisok-proizvodstvennyh-zadanij
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[ProductionTask], error) {
	return list.Entity[ProductionTask](ctx, s.client, "productiontask", options, nil, opts...)
}

// Delete — DELETE entity/productiontask/{id}, не в корзину. Сервер удаляет и связанные выполненные
// этапы производства. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-udalit-proizvodstvennoe-zadanie
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete productiontask: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete productiontask "+id, "entity/productiontask/"+id, opts)
}

// CreateMany создаёт задания одним POST entity/productiontask: массив до batch.MaxItems элементов.
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-massovoe-sozdanie-i-obnovlenie-proizvodstvennyh-zadanij
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[ProductionTask], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[ProductionTask](ctx, s.client, "productiontask", items, opts)
}

// UpdateMany изменяет задания одним POST entity/productiontask: элемент — meta по ID и поля
// Request. До запроса проверяются только UUID и непустой Request. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-massovoe-sozdanie-i-obnovlenie-proizvodstvennyh-zadanij
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[ProductionTask], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[ProductionTask](ctx, s.client, "productiontask", items, opts)
}

// DeleteMany удаляет задания по ids одним POST entity/productiontask/delete, не в корзину; до
// batch.MaxItems элементов. Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-massovoe-udalenie-proizvodstvennyh-zadanij
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "productiontask", ids, opts)
}

// Trash — удаление документа в корзину: POST entity/productiontask/{id}/trash. В отличие от
// Delete документ восстанавливается из корзины (только в интерфейсе); права и настройку корзины
// проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/common-info#3-udalenie-v-korzinu
func (s *Service) Trash(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	return write.Trash(ctx, s.client, "productiontask", id, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*ProductionTask, error) {
	var result ProductionTask
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

// getPage читает одну страницу вложенного набора: только limit и offset (fields отклоняется до
// запроса), filter — фиксированный фильтр ресурса или "".
func getPage[T any](ctx context.Context, client *transport.Client, operation, id, endpoint string, page query.PageOptions, filter string, opts []transport.RequestOption) (*model.Collection[T], error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	if page.Fields != "" {
		return nil, fmt.Errorf("moysklad: %s: fields=%q is not supported", operation, page.Fields)
	}
	params, err := page.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if filter != "" {
		params.Set("filter", filter)
	}
	return list.Page[T](ctx, client, operation, endpoint, params, opts...)
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
