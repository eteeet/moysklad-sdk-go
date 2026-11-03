// Package nested — запросы к вложенной коллекции entity/{parent}/{id}/{collection}: зоны и
// ячейки склада, счета, контактные лица и события контрагента, счета юрлица.
package nested

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

// Collection описывает коллекцию одного родителя: Parent — сущность («store»), Path — сегмент
// коллекции («zones»), Item — имя элемента в тексте ошибок («zone»), MetaType — тип meta
// элемента в массовых запросах («storezone»).
type Collection struct {
	Client   *transport.Client
	Parent   string
	ParentID string
	Path     string
	Item     string
	MetaType model.MetaType
}

func (c Collection) endpoint() string {
	return "entity/" + c.Parent + "/" + c.ParentID + "/" + c.Path
}

func (c Collection) prefix() string {
	return c.Parent + " " + c.ParentID + " " + c.Item
}

// List читает одну страницу коллекции: только limit и offset.
func List[T any](ctx context.Context, c Collection, page query.PageOptions, opts []transport.RequestOption) (*model.Collection[T], error) {
	transport.ResetResponseInfo(opts)
	operation := "list " + c.Parent + " " + c.ParentID + " " + c.Path
	if !validID(c.ParentID) {
		return nil, fmt.Errorf("moysklad: %s: %s id must be a UUID", operation, c.Parent)
	}
	if page.Fields != "" {
		return nil, fmt.Errorf("moysklad: %s: fields=%q is not supported", operation, page.Fields)
	}
	params, err := page.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[T](ctx, c.Client, operation, c.endpoint(), params, opts...)
}

// Get — GET элемента; ответ — объект.
func Get[T any](ctx context.Context, c Collection, id string, opts []transport.RequestOption) (*T, error) {
	transport.ResetResponseInfo(opts)
	operation := "get " + c.prefix() + " " + id
	if err := c.check(operation, id); err != nil {
		return nil, err
	}
	return one[T](ctx, c, operation, http.MethodGet, c.endpoint()+"/"+id, nil, opts)
}

// Create — POST коллекции. Примеры Create источника возвращают массив из одного объекта:
// принимаются и объект, и такой массив; другой размер массива — ошибка.
func Create[T any](ctx context.Context, c Collection, body any, opts []transport.RequestOption) (*T, error) {
	transport.ResetResponseInfo(opts)
	operation := "create " + c.prefix()
	if !validID(c.ParentID) {
		return nil, fmt.Errorf("moysklad: %s: %s id must be a UUID", operation, c.Parent)
	}
	return one[T](ctx, c, operation, http.MethodPost, c.endpoint(), body, opts)
}

// Update — PUT элемента; пустой body (все Field нулевые) — ошибка до запроса. Ответ — объект
// или массив из одного объекта, как у Create.
func Update[T any](ctx context.Context, c Collection, id string, body any, opts []transport.RequestOption) (*T, error) {
	transport.ResetResponseInfo(opts)
	operation := "update " + c.prefix() + " " + id
	if err := c.check(operation, id); err != nil {
		return nil, err
	}
	if reflect.ValueOf(body).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	return one[T](ctx, c, operation, http.MethodPut, c.endpoint()+"/"+id, body, opts)
}

// Delete — DELETE элемента. Повторов нет.
func Delete(ctx context.Context, c Collection, id string, opts []transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	operation := "delete " + c.prefix() + " " + id
	if err := c.check(operation, id); err != nil {
		return err
	}
	return write.Delete(ctx, c.Client, operation, c.endpoint()+"/"+id, opts)
}

// CreateMany — один POST коллекции с массивом новых элементов; правила ответа — пакет batch.
func CreateMany[T, C any](ctx context.Context, c Collection, items []C, opts []transport.RequestOption) ([]batch.Result[T], error) {
	transport.ResetResponseInfo(opts)
	operation := "create many " + c.Parent + " " + c.ParentID + " " + c.Path
	if !validID(c.ParentID) {
		return nil, fmt.Errorf("moysklad: %s: %s id must be a UUID", operation, c.Parent)
	}
	return write.CreateManyNested[T](ctx, c.Client, operation, c.endpoint(), items, opts)
}

// UpdateMany — один POST коллекции: элемент — meta по ID (тип MetaType) и поля запроса.
func UpdateMany[T, U any](ctx context.Context, c Collection, items []batch.Update[U], opts []transport.RequestOption) ([]batch.Result[T], error) {
	transport.ResetResponseInfo(opts)
	operation := "update many " + c.Parent + " " + c.ParentID + " " + c.Path
	if !validID(c.ParentID) {
		return nil, fmt.Errorf("moysklad: %s: %s id must be a UUID", operation, c.Parent)
	}
	return write.UpdateManyNested[T](ctx, c.Client, operation, c.endpoint(), c.MetaType, items, opts)
}

// DeleteMany — один POST {коллекция}/delete с meta элементов по ids.
func DeleteMany(ctx context.Context, c Collection, ids []string, opts []transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	operation := "delete many " + c.Parent + " " + c.ParentID + " " + c.Path
	if !validID(c.ParentID) {
		return nil, fmt.Errorf("moysklad: %s: %s id must be a UUID", operation, c.Parent)
	}
	return write.DeleteManyNested(ctx, c.Client, operation, c.endpoint(), c.MetaType, ids, opts)
}

func (c Collection) check(operation, id string) error {
	if !validID(c.ParentID) || !validID(id) {
		return fmt.Errorf("moysklad: %s: %s and %s ids must be UUIDs", operation, c.Parent, c.Item)
	}
	return nil
}

func one[T any](ctx context.Context, c Collection, operation, method, endpoint string, body any, opts []transport.RequestOption) (*T, error) {
	var raw jsontext.Value
	resp, err := c.Client.Send(ctx, method, endpoint, nil, body, &raw, opts...)
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
	if err := json.Unmarshal(raw, &head); err != nil || head.ID == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
	}
	var result T
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("moysklad: %s: decode response: %w", operation, err)
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
