package product

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
		return nil, fmt.Errorf("moysklad: product service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#3-poluchit-tovar
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Product, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: product id must be a UUID")
	}

	return s.do(ctx, "get product "+id, http.MethodGet, "entity/product/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#3-sozdat-tovar
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль товара.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Product, error) {
	transport.ResetResponseInfo(opts)
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create product: name must not be empty")
	}
	return s.do(ctx, "create product", http.MethodPost, "entity/product", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#3-izmenit-tovar
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Product, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update product: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value(0) и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update product: no fields specified")
	}
	return s.do(ctx, "update product "+id, http.MethodPut, "entity/product/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Product, error) {
	var result Product
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
// fields=minimumStock не принимается: в модели Product нет поля неснижаемого остатка (TODO).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Product], error) {
	return list.Entity[Product](ctx, s.client, "product", options, nil, opts...)
}

// Delete — «Удалить Товар»: DELETE entity/product/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#3-udalit-tovar
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete product: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete product "+id, "entity/product/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/product: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление товаров».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#3-massovoe-sozdanie-i-obnovlenie-tovarov
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[Product], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[Product](ctx, s.client, "product", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/product: элемент — meta по ID и поля Request.
// Источник: «Массовое создание и обновление товаров».
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#3-massovoe-sozdanie-i-obnovlenie-tovarov
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Product], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[Product](ctx, s.client, "product", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/product/delete, не в корзину; до
// batch.MaxItems элементов. Источник: «Массовое удаление Товаров».
// Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#3-massovoe-udalenie-tovarov
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "product", ids, opts)
}
