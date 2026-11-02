package uom

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
		return nil, fmt.Errorf("moysklad: uom service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/uom#3-poluchit-edinicu-izmereniya
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*UOM, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: uom id must be a UUID")
	}
	return s.do(ctx, "get uom "+id, http.MethodGet, "entity/uom/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/uom#3-sozdat-edinicu-izmereniya
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль UOM.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*UOM, error) {
	transport.ResetResponseInfo(opts)
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create uom: name must not be empty")
	}
	return s.do(ctx, "create uom", http.MethodPost, "entity/uom", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/uom#3-izmenit-edinicu-izmereniya
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*UOM, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update uom: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update uom: no fields specified")
	}
	return s.do(ctx, "update uom "+id, http.MethodPut, "entity/uom/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*UOM, error) {
	var result UOM
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
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/uom
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[UOM], error) {
	return list.Entity[UOM](ctx, s.client, "uom", options, nil, opts...)
}

// Delete — «Удалить Единицу измерения»: DELETE entity/uom/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/uom#3-udalit-edinicu-izmereniya
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete uom: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete uom "+id, "entity/uom/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/uom: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление единиц измерения».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/uom#3-massovoe-sozdanie-i-obnovlenie-edinic-izmereniya
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[UOM], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[UOM](ctx, s.client, "uom", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/uom: элемент — meta по ID и поля Request.
// Источник: «Массовое создание и обновление единиц измерения».
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/uom#3-massovoe-sozdanie-i-obnovlenie-edinic-izmereniya
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[UOM], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[UOM](ctx, s.client, "uom", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/uom/delete, не в корзину; до
// batch.MaxItems элементов. Источник: «Массовое удаление Единиц измерения».
// Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/uom#3-massovoe-udalenie-edinic-izmereniya
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "uom", ids, opts)
}
