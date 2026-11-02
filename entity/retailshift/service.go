package retailshift

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
		return nil, fmt.Errorf("moysklad: retailshift service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retailshift#3-poluchit-roznichnuyu-smenu
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*RetailShift, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get retailshift: id must be a UUID")
	}
	return s.do(ctx, "get retailshift "+id, http.MethodGet, "entity/retailshift/"+id, nil, opts)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможна вторая смена.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retailshift#3-sozdat-roznichnuyu-smenu
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*RetailShift, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create retailshift"
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: organization meta.href must not be empty", operation)
	}
	if input.RetailStore.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: retailStore meta.href must not be empty", operation)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/retailshift", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retailshift#3-izmenit-roznichnuyu-smenu
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*RetailShift, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update retailshift"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/retailshift/"+id, input, opts)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину смены сервер по умолчанию не отдаёт: см. filter isDeleted.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retailshift#3-poluchit-roznichnye-smeny
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[RetailShift], error) {
	return list.Entity[RetailShift](ctx, s.client, "retailshift", options, nil, opts...)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*RetailShift, error) {
	var result RetailShift
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

// Delete — «Удалить Розничную смену»: DELETE entity/retailshift/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retailshift#3-udalit-roznichnuyu-smenu
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete retailshift: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete retailshift "+id, "entity/retailshift/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/retailshift: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Розничных смен».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retailshift#3-massovoe-sozdanie-i-obnovlenie-roznichnyh-smen
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[RetailShift], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[RetailShift](ctx, s.client, "retailshift", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/retailshift: элемент — meta по ID и поля Request.
// Источник: «Массовое создание и обновление Розничных смен».
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retailshift#3-massovoe-sozdanie-i-obnovlenie-roznichnyh-smen
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[RetailShift], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[RetailShift](ctx, s.client, "retailshift", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/retailshift/delete, не в корзину; до
// batch.MaxItems элементов. Источник: «Массовое удаление Розничных смен».
// Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retailshift#3-massovoe-udalenie-roznichnyh-smen
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "retailshift", ids, opts)
}
