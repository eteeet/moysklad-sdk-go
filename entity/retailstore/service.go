package retailstore

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
		return nil, fmt.Errorf("moysklad: retailstore service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/retailstore#3-poluchit-tochku-prodazh
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*RetailStore, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get retailstore: id must be a UUID")
	}
	return s.do(ctx, "get retailstore "+id, http.MethodGet, "entity/retailstore/"+id, nil, opts)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль точки продаж.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/retailstore#3-sozdat-tochku-prodazh
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*RetailStore, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create retailstore"
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: %s: name must not be empty", operation)
	}
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: organization meta.href must not be empty", operation)
	}
	if input.Store.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: store meta.href must not be empty", operation)
	}
	if input.PriceType.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: priceType meta.href must not be empty", operation)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/retailstore", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/retailstore#3-izmenit-tochku-prodazh
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*RetailStore, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update retailstore"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/retailstore/"+id, input, opts)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Поиск — параметр Search; архивные точки сервер по умолчанию не отдаёт: см. filter archived.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/retailstore#3-poluchit-tochki-prodazh
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[RetailStore], error) {
	return list.Entity[RetailStore](ctx, s.client, "retailstore", options, nil, opts...)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*RetailStore, error) {
	var result RetailStore
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

// Delete — «Удалить точку продаж»: DELETE entity/retailstore/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/retailstore#3-udalit-tochku-prodazh
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete retailstore: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete retailstore "+id, "entity/retailstore/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/retailstore: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Точек продаж».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/retailstore#3-massovoe-sozdanie-i-obnovlenie-tochek-prodazh
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[RetailStore], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[RetailStore](ctx, s.client, "retailstore", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/retailstore: элемент — meta по ID и поля Request.
// Источник: «Массовое создание и обновление Точек продаж».
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/retailstore#3-massovoe-sozdanie-i-obnovlenie-tochek-prodazh
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[RetailStore], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[RetailStore](ctx, s.client, "retailstore", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/retailstore/delete, не в корзину; до
// batch.MaxItems элементов. Источник: «Массовое удаление Точек продаж».
// Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/retailstore#3-massovoe-udalenie-tochek-prodazh
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "retailstore", ids, opts)
}
