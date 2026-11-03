package store

import (
	"context"
	"encoding/hex"
	"encoding/json/jsontext"
	"encoding/json/v2"
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
		return nil, fmt.Errorf("moysklad: store service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-poluchit-sklad
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Store, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: store id must be a UUID")
	}
	return s.do(ctx, "get store "+id, http.MethodGet, "entity/store/"+id, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-sozdat-sklad
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль склада.
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Store, error) {
	transport.ResetResponseInfo(opts)
	if input.Name == "" {
		return nil, fmt.Errorf("moysklad: create store: name must not be empty")
	}
	return s.do(ctx, "create store", http.MethodPost, "entity/store", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-izmenit-sklad
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Store, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update store: id must be a UUID")
	}
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update store: no fields specified")
	}
	return s.do(ctx, "update store "+id, http.MethodPut, "entity/store/"+id, input, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Store, error) {
	var result Store
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

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-poluchit-sklady
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Store], error) {
	return list.Entity[Store](ctx, s.client, "store", options, nil, opts...)
}

// Delete — «Удалить Склад»: DELETE entity/store/{id}, не в корзину. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-udalit-sklad
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete store: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete store "+id, "entity/store/"+id, opts)
}

// CreateMany создаёт склады одним POST entity/store; правила ответа — пакет batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-massovoe-sozdanie-i-obnovlenie-skladov
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[Store], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[Store](ctx, s.client, "store", items, opts)
}

// UpdateMany изменяет склады одним POST entity/store: meta по ID и поля Request. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-massovoe-sozdanie-i-obnovlenie-skladov
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Store], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[Store](ctx, s.client, "store", items, opts)
}

// DeleteMany удаляет склады одним POST entity/store/delete, не в корзину. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-massovoe-udalenie-skladov
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "store", ids, opts)
}

// Зоны и ячейки — вложенные коллекции entity/store/{id}/zones и …/slots. Список принимает
// только limit и offset.

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-poluchit-zony-sklada
func (s *Service) ListZones(ctx context.Context, storeID string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Zone], error) {
	return listNested[Zone](ctx, s, storeID, "zones", page, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-poluchit-zonu-sklada
func (s *Service) GetZone(ctx context.Context, storeID, zoneID string, opts ...transport.RequestOption) (*Zone, error) {
	return oneNested[Zone](ctx, s, http.MethodGet, storeID, "zones", zoneID, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-sozdat-zonu-sklada
func (s *Service) CreateZone(ctx context.Context, storeID string, input ZoneCreateRequest, opts ...transport.RequestOption) (*Zone, error) {
	if input.Name == "" {
		transport.ResetResponseInfo(opts)
		return nil, fmt.Errorf("moysklad: create store %s zone: name must not be empty", storeID)
	}
	return oneNested[Zone](ctx, s, http.MethodPost, storeID, "zones", "", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-izmenit-zonu-sklada
func (s *Service) UpdateZone(ctx context.Context, storeID, zoneID string, input ZoneUpdateRequest, opts ...transport.RequestOption) (*Zone, error) {
	return oneNested[Zone](ctx, s, http.MethodPut, storeID, "zones", zoneID, input, opts)
}

// DeleteZone — DELETE entity/store/{id}/zones/{zoneID}. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-udalit-zonu-sklada
func (s *Service) DeleteZone(ctx context.Context, storeID, zoneID string, opts ...transport.RequestOption) error {
	return deleteNested(ctx, s, storeID, "zones", zoneID, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-massovoe-sozdanie-i-obnovlenie-zon-sklada
func (s *Service) CreateManyZones(ctx context.Context, storeID string, items []ZoneCreateRequest, opts ...transport.RequestOption) ([]batch.Result[Zone], error) {
	return createManyNested[Zone](ctx, s, storeID, "zones", items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-massovoe-sozdanie-i-obnovlenie-zon-sklada
func (s *Service) UpdateManyZones(ctx context.Context, storeID string, items []batch.Update[ZoneUpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Zone], error) {
	return updateManyNested[Zone](ctx, s, storeID, "zones", "storezone", items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-massovoe-udalenie-zon-sklada
func (s *Service) DeleteManyZones(ctx context.Context, storeID string, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	return deleteManyNested(ctx, s, storeID, "zones", "storezone", ids, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-poluchit-yachejki-sklada
func (s *Service) ListSlots(ctx context.Context, storeID string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Slot], error) {
	return listNested[Slot](ctx, s, storeID, "slots", page, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-poluchit-yachejku-sklada
func (s *Service) GetSlot(ctx context.Context, storeID, slotID string, opts ...transport.RequestOption) (*Slot, error) {
	return oneNested[Slot](ctx, s, http.MethodGet, storeID, "slots", slotID, nil, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-sozdat-yachejku-sklada
func (s *Service) CreateSlot(ctx context.Context, storeID string, input SlotCreateRequest, opts ...transport.RequestOption) (*Slot, error) {
	if input.Name == "" {
		transport.ResetResponseInfo(opts)
		return nil, fmt.Errorf("moysklad: create store %s slot: name must not be empty", storeID)
	}
	return oneNested[Slot](ctx, s, http.MethodPost, storeID, "slots", "", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-izmenit-yachejku-sklada
func (s *Service) UpdateSlot(ctx context.Context, storeID, slotID string, input SlotUpdateRequest, opts ...transport.RequestOption) (*Slot, error) {
	return oneNested[Slot](ctx, s, http.MethodPut, storeID, "slots", slotID, input, opts)
}

// DeleteSlot — DELETE entity/store/{id}/slots/{slotID}. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-udalit-yachejku-sklada
func (s *Service) DeleteSlot(ctx context.Context, storeID, slotID string, opts ...transport.RequestOption) error {
	return deleteNested(ctx, s, storeID, "slots", slotID, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-massovoe-sozdanie-i-obnovlenie-yacheek-sklada
func (s *Service) CreateManySlots(ctx context.Context, storeID string, items []SlotCreateRequest, opts ...transport.RequestOption) ([]batch.Result[Slot], error) {
	return createManyNested[Slot](ctx, s, storeID, "slots", items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-massovoe-sozdanie-i-obnovlenie-yacheek-sklada
func (s *Service) UpdateManySlots(ctx context.Context, storeID string, items []batch.Update[SlotUpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Slot], error) {
	return updateManyNested[Slot](ctx, s, storeID, "slots", "slot", items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-massovoe-udalenie-yacheek-sklada
func (s *Service) DeleteManySlots(ctx context.Context, storeID string, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	return deleteManyNested(ctx, s, storeID, "slots", "slot", ids, opts)
}

func listNested[T any](ctx context.Context, s *Service, storeID, collection string, page query.PageOptions, opts []transport.RequestOption) (*model.Collection[T], error) {
	transport.ResetResponseInfo(opts)
	operation := "list store " + storeID + " " + collection
	if !validID(storeID) {
		return nil, fmt.Errorf("moysklad: %s: store id must be a UUID", operation)
	}
	if page.Fields != "" {
		return nil, fmt.Errorf("moysklad: %s: fields=%q is not supported", operation, page.Fields)
	}
	params, err := page.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[T](ctx, s.client, operation, "entity/store/"+storeID+"/"+collection, params, opts...)
}

// oneNested выполняет Get (GET), Create (POST, id пуст) или Update (PUT) зоны или ячейки.
// Примеры Create и Update источника возвращают массив из одного объекта, Get — объект:
// для записи принимаются обе формы, другой размер массива — ошибка.
func oneNested[T any](ctx context.Context, s *Service, method, storeID, collection, id string, body any, opts []transport.RequestOption) (*T, error) {
	transport.ResetResponseInfo(opts)
	verb := map[string]string{http.MethodGet: "get", http.MethodPost: "create", http.MethodPut: "update"}[method]
	operation := verb + " store " + storeID + " " + collection[:len(collection)-1]
	endpoint := "entity/store/" + storeID + "/" + collection
	if method != http.MethodPost {
		operation += " " + id
		endpoint += "/" + id
		if !validID(id) {
			return nil, fmt.Errorf("moysklad: %s: store and %s ids must be UUIDs", operation, collection[:len(collection)-1])
		}
	}
	if !validID(storeID) {
		return nil, fmt.Errorf("moysklad: %s: store id must be a UUID", operation)
	}
	if method == http.MethodPut && reflect.ValueOf(body).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	var raw jsontext.Value
	resp, err := s.client.Send(ctx, method, endpoint, nil, body, &raw, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if method != http.MethodGet && raw.Kind() == '[' {
		var items []jsontext.Value
		if err := json.Unmarshal(raw, &items); err != nil || len(items) != 1 {
			return nil, fmt.Errorf("moysklad: %s: response array must contain exactly one object", operation)
		}
		raw = items[0]
	}
	var head struct {
		ID string `json:"id"`
	}
	var result T
	if err := json.Unmarshal(raw, &head); err != nil || head.ID == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("moysklad: %s: decode response: %w", operation, err)
	}
	return &result, nil
}

func deleteNested(ctx context.Context, s *Service, storeID, collection, id string, opts []transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	operation := "delete store " + storeID + " " + collection[:len(collection)-1] + " " + id
	if !validID(storeID) || !validID(id) {
		return fmt.Errorf("moysklad: %s: store and %s ids must be UUIDs", operation, collection[:len(collection)-1])
	}
	return write.Delete(ctx, s.client, operation, "entity/store/"+storeID+"/"+collection+"/"+id, opts)
}

func createManyNested[T, C any](ctx context.Context, s *Service, storeID, collection string, items []C, opts []transport.RequestOption) ([]batch.Result[T], error) {
	transport.ResetResponseInfo(opts)
	operation := "create many store " + storeID + " " + collection
	if !validID(storeID) {
		return nil, fmt.Errorf("moysklad: %s: store id must be a UUID", operation)
	}
	return write.CreateManyNested[T](ctx, s.client, operation, "entity/store/"+storeID+"/"+collection, items, opts)
}

func updateManyNested[T, U any](ctx context.Context, s *Service, storeID, collection string, metaType model.MetaType, items []batch.Update[U], opts []transport.RequestOption) ([]batch.Result[T], error) {
	transport.ResetResponseInfo(opts)
	operation := "update many store " + storeID + " " + collection
	if !validID(storeID) {
		return nil, fmt.Errorf("moysklad: %s: store id must be a UUID", operation)
	}
	return write.UpdateManyNested[T](ctx, s.client, operation, "entity/store/"+storeID+"/"+collection, metaType, items, opts)
}

func deleteManyNested(ctx context.Context, s *Service, storeID, collection string, metaType model.MetaType, ids []string, opts []transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	operation := "delete many store " + storeID + " " + collection
	if !validID(storeID) {
		return nil, fmt.Errorf("moysklad: %s: store id must be a UUID", operation)
	}
	return write.DeleteManyNested(ctx, s.client, operation, "entity/store/"+storeID+"/"+collection, metaType, ids, opts)
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
