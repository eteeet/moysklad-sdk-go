package taxrate

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
		return nil, fmt.Errorf("moysklad: taxrate service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/taxrate#3-poluchit-stavku-nds
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*TaxRate, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: taxrate id must be a UUID")
	}
	return s.do(ctx, "get taxrate "+id, http.MethodGet, "entity/taxrate/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/taxrate#3-sozdat-stavku-nds
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль TaxRate.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*TaxRate, error) {
	transport.ResetResponseInfo(opts)
	if input.Rate == nil {
		return nil, fmt.Errorf("moysklad: create taxrate: rate is required")
	}
	return s.do(ctx, "create taxrate", http.MethodPost, "entity/taxrate", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/taxrate#3-izmenit-stavku-nds
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*TaxRate, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update taxrate: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update taxrate: no fields specified")
	}
	return s.do(ctx, "update taxrate "+id, http.MethodPut, "entity/taxrate/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*TaxRate, error) {
	var result TaxRate
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
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/taxrate
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[TaxRate], error) {
	return list.Entity[TaxRate](ctx, s.client, "taxrate", options, nil, opts...)
}

// Delete — «Удалить Ставку НДС»: DELETE entity/taxrate/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/taxrate#3-udalit-stavku-nds
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete taxrate: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete taxrate "+id, "entity/taxrate/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/taxrate: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Ставок НДС».
// Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и неизвестный
// исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/taxrate#3-massovoe-sozdanie-i-obnovlenie-stavok-nds
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[TaxRate], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[TaxRate](ctx, s.client, "taxrate", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/taxrate: элемент — meta по ID и поля Request.
// Источник: «Массовое создание и обновление Ставок НДС».
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/taxrate#3-massovoe-sozdanie-i-obnovlenie-stavok-nds
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[TaxRate], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[TaxRate](ctx, s.client, "taxrate", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/taxrate/delete, не в корзину; до
// batch.MaxItems элементов. Источник: «Массовое удаление Ставок НДС».
// Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/taxrate#3-massovoe-udalenie-stavok-nds
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "taxrate", ids, opts)
}
