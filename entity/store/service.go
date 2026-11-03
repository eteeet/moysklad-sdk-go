package store

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
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
// только limit и offset; Create и Update принимают ответ-объект и массив из одного объекта
// (форма примеров источника).

func (s *Service) zones(storeID string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "store", ParentID: storeID, Path: "zones", Item: "zone", MetaType: "storezone"}
}

func (s *Service) slots(storeID string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "store", ParentID: storeID, Path: "slots", Item: "slot", MetaType: "slot"}
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-poluchit-zony-sklada
func (s *Service) ListZones(ctx context.Context, storeID string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Zone], error) {
	return nested.List[Zone](ctx, s.zones(storeID), page, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-poluchit-zonu-sklada
func (s *Service) GetZone(ctx context.Context, storeID, zoneID string, opts ...transport.RequestOption) (*Zone, error) {
	return nested.Get[Zone](ctx, s.zones(storeID), zoneID, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-sozdat-zonu-sklada
func (s *Service) CreateZone(ctx context.Context, storeID string, input ZoneCreateRequest, opts ...transport.RequestOption) (*Zone, error) {
	if input.Name == "" {
		transport.ResetResponseInfo(opts)
		return nil, fmt.Errorf("moysklad: create store %s zone: name must not be empty", storeID)
	}
	return nested.Create[Zone](ctx, s.zones(storeID), input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-izmenit-zonu-sklada
func (s *Service) UpdateZone(ctx context.Context, storeID, zoneID string, input ZoneUpdateRequest, opts ...transport.RequestOption) (*Zone, error) {
	return nested.Update[Zone](ctx, s.zones(storeID), zoneID, input, opts)
}

// DeleteZone — DELETE entity/store/{id}/zones/{zoneID}. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-udalit-zonu-sklada
func (s *Service) DeleteZone(ctx context.Context, storeID, zoneID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.zones(storeID), zoneID, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-massovoe-sozdanie-i-obnovlenie-zon-sklada
func (s *Service) CreateManyZones(ctx context.Context, storeID string, items []ZoneCreateRequest, opts ...transport.RequestOption) ([]batch.Result[Zone], error) {
	return nested.CreateMany[Zone](ctx, s.zones(storeID), items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-massovoe-sozdanie-i-obnovlenie-zon-sklada
func (s *Service) UpdateManyZones(ctx context.Context, storeID string, items []batch.Update[ZoneUpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Zone], error) {
	return nested.UpdateMany[Zone](ctx, s.zones(storeID), items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-massovoe-udalenie-zon-sklada
func (s *Service) DeleteManyZones(ctx context.Context, storeID string, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	return nested.DeleteMany(ctx, s.zones(storeID), ids, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-poluchit-yachejki-sklada
func (s *Service) ListSlots(ctx context.Context, storeID string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Slot], error) {
	return nested.List[Slot](ctx, s.slots(storeID), page, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-poluchit-yachejku-sklada
func (s *Service) GetSlot(ctx context.Context, storeID, slotID string, opts ...transport.RequestOption) (*Slot, error) {
	return nested.Get[Slot](ctx, s.slots(storeID), slotID, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-sozdat-yachejku-sklada
func (s *Service) CreateSlot(ctx context.Context, storeID string, input SlotCreateRequest, opts ...transport.RequestOption) (*Slot, error) {
	if input.Name == "" {
		transport.ResetResponseInfo(opts)
		return nil, fmt.Errorf("moysklad: create store %s slot: name must not be empty", storeID)
	}
	return nested.Create[Slot](ctx, s.slots(storeID), input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-izmenit-yachejku-sklada
func (s *Service) UpdateSlot(ctx context.Context, storeID, slotID string, input SlotUpdateRequest, opts ...transport.RequestOption) (*Slot, error) {
	return nested.Update[Slot](ctx, s.slots(storeID), slotID, input, opts)
}

// DeleteSlot — DELETE entity/store/{id}/slots/{slotID}. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-udalit-yachejku-sklada
func (s *Service) DeleteSlot(ctx context.Context, storeID, slotID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.slots(storeID), slotID, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-massovoe-sozdanie-i-obnovlenie-yacheek-sklada
func (s *Service) CreateManySlots(ctx context.Context, storeID string, items []SlotCreateRequest, opts ...transport.RequestOption) ([]batch.Result[Slot], error) {
	return nested.CreateMany[Slot](ctx, s.slots(storeID), items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-massovoe-sozdanie-i-obnovlenie-yacheek-sklada
func (s *Service) UpdateManySlots(ctx context.Context, storeID string, items []batch.Update[SlotUpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Slot], error) {
	return nested.UpdateMany[Slot](ctx, s.slots(storeID), items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/store#3-massovoe-udalenie-yacheek-sklada
func (s *Service) DeleteManySlots(ctx context.Context, storeID string, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	return nested.DeleteMany(ctx, s.slots(storeID), ids, opts)
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
