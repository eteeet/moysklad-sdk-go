// Package write — общие запросы записи сервисов: удаление и массовые операции.
package write

import (
	"context"
	"encoding/hex"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Delete — один DELETE endpoint. Успех — 200 или 204, тело ответа не читается. DELETE не
// повторяется: transport повторяет только GET и HEAD.
func Delete(ctx context.Context, client *transport.Client, operation, endpoint string, opts []transport.RequestOption) error {
	resp, err := client.Send(ctx, http.MethodDelete, endpoint, nil, nil, nil, opts...)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	return nil
}

// CreateMany — один POST entity/{entity} с массивом новых объектов.
func CreateMany[T, C any](ctx context.Context, client *transport.Client, entity string, items []C, opts []transport.RequestOption) ([]batch.Result[T], error) {
	operation := "create many " + entity
	if err := checkCount(len(items)); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	return post[batch.Result[T]](ctx, client, operation, "entity/"+entity, items, len(items), opts)
}

// UpdateMany — один POST entity/{entity}: каждый элемент — meta объекта по ID и поля Request.
// До запроса проверяются только UUID и непустой Request; остальные проверки Update — на сервере.
func UpdateMany[T, U any](ctx context.Context, client *transport.Client, entity string, items []batch.Update[U], opts []transport.RequestOption) ([]batch.Result[T], error) {
	operation := "update many " + entity
	if err := checkCount(len(items)); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	body := make([]jsontext.Value, len(items))
	for i, item := range items {
		if !validID(item.ID) {
			return nil, fmt.Errorf("moysklad: %s: items[%d]: id must be a UUID", operation, i)
		}
		if reflect.ValueOf(item.Request).IsZero() {
			return nil, fmt.Errorf("moysklad: %s: items[%d]: no fields specified", operation, i)
		}
		fields, err := json.Marshal(item.Request)
		if err != nil {
			return nil, fmt.Errorf("moysklad: %s: items[%d]: %w", operation, i, err)
		}
		if len(fields) < 2 || fields[0] != '{' {
			return nil, fmt.Errorf("moysklad: %s: items[%d]: request must encode to a JSON object", operation, i)
		}
		meta, err := json.Marshal(objectMeta(client, entity, item.ID))
		if err != nil {
			return nil, fmt.Errorf("moysklad: %s: items[%d]: %w", operation, i, err)
		}
		// {"meta":{...},<поля Request>}: meta первым, поля — как в одиночном Update.
		object := append([]byte(`{"meta":`), meta...)
		if string(fields) != "{}" {
			object = append(append(object, ','), fields[1:]...)
		} else {
			object = append(object, '}')
		}
		body[i] = object
	}
	return post[batch.Result[T]](ctx, client, operation, "entity/"+entity, body, len(items), opts)
}

// DeleteMany — один POST entity/{entity}/delete с массивом meta объектов по ids.
func DeleteMany(ctx context.Context, client *transport.Client, entity string, ids []string, opts []transport.RequestOption) ([]batch.Deleted, error) {
	operation := "delete many " + entity
	if err := checkCount(len(ids)); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	metas := make([]model.Meta, len(ids))
	for i, id := range ids {
		if !validID(id) {
			return nil, fmt.Errorf("moysklad: %s: ids[%d]: id must be a UUID", operation, i)
		}
		metas[i] = objectMeta(client, entity, id)
	}
	return deleteMetas(ctx, client, operation, entity, metas, opts)
}

// DeleteMetas — один POST entity/{entity}/delete с готовыми meta (ассортимент: объекты разных
// типов). У каждой meta нужны href и type.
func DeleteMetas(ctx context.Context, client *transport.Client, entity string, metas []model.Meta, opts []transport.RequestOption) ([]batch.Deleted, error) {
	operation := "delete many " + entity
	if err := checkCount(len(metas)); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	for i, meta := range metas {
		if meta.Href == "" || meta.Type == "" {
			return nil, fmt.Errorf("moysklad: %s: metas[%d]: href and type must not be empty", operation, i)
		}
	}
	return deleteMetas(ctx, client, operation, entity, metas, opts)
}

func deleteMetas(ctx context.Context, client *transport.Client, operation, entity string, metas []model.Meta, opts []transport.RequestOption) ([]batch.Deleted, error) {
	type item struct {
		Meta model.Meta `json:"meta"`
	}
	body := make([]item, len(metas))
	for i, meta := range metas {
		body[i] = item{meta}
	}
	return post[batch.Deleted](ctx, client, operation, "entity/"+entity+"/delete", body, len(metas), opts)
}

func checkCount(n int) error {
	if n == 0 || n > batch.MaxItems {
		return fmt.Errorf("items must contain 1 to %d elements, got %d", batch.MaxItems, n)
	}
	return nil
}

// objectMeta — meta объекта, как в примерах массовых запросов источника. У комплектов и услуг
// метаданные общие с товарами: metadataHref — entity/product/metadata.
func objectMeta(client *transport.Client, entity, id string) model.Meta {
	metadata := entity
	if entity == "bundle" || entity == "service" {
		metadata = "product"
	}
	return model.Meta{
		Href:         client.Href("entity/" + entity + "/" + id),
		MetadataHref: client.Href("entity/" + metadata + "/metadata"),
		Type:         model.MetaType(entity),
		MediaType:    "application/json",
	}
}

// post отправляет массив из sent элементов и ждёт HTTP 200 с массивом той же длины. POST
// transport не повторяет.
func post[R any](ctx context.Context, client *transport.Client, operation, endpoint string, body any, sent int, opts []transport.RequestOption) ([]R, error) {
	var results []R
	resp, err := client.Send(ctx, http.MethodPost, endpoint, nil, body, &results, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if len(results) != sent {
		return nil, fmt.Errorf("moysklad: %s: response has %d items for %d sent", operation, len(results), sent)
	}
	return results, nil
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}
