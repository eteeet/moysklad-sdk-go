package factureout

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
		return nil, fmt.Errorf("moysklad: factureout service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/factureout#3-poluchit-schet-fakturu-vydannyj
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*FactureOut, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get factureout: id must be a UUID")
	}
	return s.do(ctx, "get factureout "+id, http.MethodGet, "entity/factureout/"+id, nil, opts)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// Удалённые в корзину документы сервер по умолчанию не отдаёт: см. filter isDeleted.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/factureout#3-poluchit-vydannye-scheta-faktury
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[FactureOut], error) {
	return list.Entity[FactureOut](ctx, s.client, "factureout", options, nil, opts...)
}

// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль документа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/factureout#3-sozdat-schet-fakturu
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*FactureOut, error) {
	transport.ResetResponseInfo(opts)
	if input.Organization.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create factureout: organization meta.href must not be empty")
	}
	return s.do(ctx, "create factureout", http.MethodPost, "entity/factureout", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/factureout#3-izmenit-schet-fakturu-vydannyj
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*FactureOut, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update factureout: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false), Value("") и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update factureout: no fields specified")
	}
	return s.do(ctx, "update factureout "+id, http.MethodPut, "entity/factureout/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*FactureOut, error) {
	var result FactureOut
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

// Delete — «Удалить Счет-фактуру выданный»: DELETE entity/factureout/{id}, не в корзину.
// Повторов нет; права и связанные объекты проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/factureout#3-udalit-schet-fakturu-vydannyj
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete factureout: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete factureout "+id, "entity/factureout/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/factureout: массив до batch.MaxItems элементов.
// Источник: «Массовое создание и обновление Счетов-фактур выданных».
// Обязательные поля и основание проверяет сервер у каждого элемента; ответ, ошибки элементов и
// неизвестный исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/factureout#3-massovoe-sozdanie-i-obnovlenie-schetov-faktur-vydannyh
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[FactureOut], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[FactureOut](ctx, s.client, "factureout", items, opts)
}

// UpdateMany изменяет объекты одним POST entity/factureout: элемент — meta по ID и поля Request.
// Источник: «Массовое создание и обновление Счетов-фактур выданных».
// До запроса проверяются только UUID и непустой Request, остальные проверки Update делает
// сервер у каждого элемента. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/factureout#3-massovoe-sozdanie-i-obnovlenie-schetov-faktur-vydannyh
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[FactureOut], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[FactureOut](ctx, s.client, "factureout", items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/factureout/delete, не в корзину; до
// batch.MaxItems элементов. Источник: «Массовое удаление Счетов-фактур выданных».
// Ответ — batch.Deleted по каждому элементу. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/factureout#3-massovoe-udalenie-schetov-faktur-vydannyh
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "factureout", ids, opts)
}
