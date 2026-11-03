// Package write — общие запросы сервисов с побочным эффектом: удаление, массовые операции
// и создание асинхронной задачи.
package write

import (
	"context"
	"encoding/hex"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/async"
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
	return post[batch.Result[T]](ctx, client, operation, "entity/"+entity, items, len(items), 0, opts)
}

// UpdateMany — один POST entity/{entity}: каждый элемент — meta объекта по ID и поля Request.
// До запроса проверяются только UUID и непустой Request; остальные проверки Update — на сервере.
func UpdateMany[T, U any](ctx context.Context, client *transport.Client, entity string, items []batch.Update[U], opts []transport.RequestOption) ([]batch.Result[T], error) {
	return updateMany[T](ctx, client, "update many "+entity, "entity/"+entity, items, func(id string) model.Meta {
		return objectMeta(client, entity, id)
	}, 0, opts)
}

// CreateManyNested — CreateMany вложенной коллекции: один POST endpoint
// (entity/{entity}/{id}/{collection}) с массивом новых объектов.
func CreateManyNested[T, C any](ctx context.Context, client *transport.Client, operation, endpoint string, items []C, opts []transport.RequestOption) ([]batch.Result[T], error) {
	return createManyNested[T](ctx, client, operation, endpoint, items, 0, opts)
}

// CreateManyNestedCreated — CreateManyNested для коллекций, где источник отвечает 201 (лента
// событий документа): принимаются 200 и 201.
func CreateManyNestedCreated[T, C any](ctx context.Context, client *transport.Client, operation, endpoint string, items []C, opts []transport.RequestOption) ([]batch.Result[T], error) {
	return createManyNested[T](ctx, client, operation, endpoint, items, http.StatusCreated, opts)
}

func createManyNested[T, C any](ctx context.Context, client *transport.Client, operation, endpoint string, items []C, created int, opts []transport.RequestOption) ([]batch.Result[T], error) {
	if err := checkCount(len(items)); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	return post[batch.Result[T]](ctx, client, operation, endpoint, items, len(items), created, opts)
}

// UpdateManyNested — UpdateMany вложенной коллекции endpoint: meta элемента — href
// endpoint/{id} и metaType без metadataHref, как в примерах источника.
func UpdateManyNested[T, U any](ctx context.Context, client *transport.Client, operation, endpoint string, metaType model.MetaType, items []batch.Update[U], opts []transport.RequestOption) ([]batch.Result[T], error) {
	return updateMany[T](ctx, client, operation, endpoint, items, func(id string) model.Meta {
		return nestedMeta(client, endpoint, metaType, id)
	}, 0, opts)
}

// UpdateManyNestedCreated — UpdateManyNested для коллекций, где источник отвечает 201: принимаются
// 200 и 201.
func UpdateManyNestedCreated[T, U any](ctx context.Context, client *transport.Client, operation, endpoint string, metaType model.MetaType, items []batch.Update[U], opts []transport.RequestOption) ([]batch.Result[T], error) {
	return updateMany[T](ctx, client, operation, endpoint, items, func(id string) model.Meta {
		return nestedMeta(client, endpoint, metaType, id)
	}, http.StatusCreated, opts)
}

// DeleteManyNested — один POST endpoint/delete с meta вложенных объектов по ids.
func DeleteManyNested(ctx context.Context, client *transport.Client, operation, endpoint string, metaType model.MetaType, ids []string, opts []transport.RequestOption) ([]batch.Deleted, error) {
	if err := checkCount(len(ids)); err != nil {
		return nil, fmt.Errorf("moysklad: %s: %w", operation, err)
	}
	type item struct {
		Meta model.Meta `json:"meta"`
	}
	body := make([]item, len(ids))
	for i, id := range ids {
		if !validID(id) {
			return nil, fmt.Errorf("moysklad: %s: ids[%d]: id must be a UUID", operation, i)
		}
		body[i] = item{nestedMeta(client, endpoint, metaType, id)}
	}
	return post[batch.Deleted](ctx, client, operation, endpoint+"/delete", body, len(ids), 0, opts)
}

func nestedMeta(client *transport.Client, endpoint string, metaType model.MetaType, id string) model.Meta {
	return model.Meta{Href: client.Href(endpoint + "/" + id), Type: metaType, MediaType: "application/json"}
}

func updateMany[T, U any](ctx context.Context, client *transport.Client, operation, endpoint string, items []batch.Update[U],
	metaOf func(id string) model.Meta, created int, opts []transport.RequestOption) ([]batch.Result[T], error) {
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
		meta, err := json.Marshal(metaOf(item.ID))
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
	return post[batch.Result[T]](ctx, client, operation, endpoint, body, len(items), created, opts)
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
	return post[batch.Deleted](ctx, client, operation, "entity/"+entity+"/delete", body, len(metas), 0, opts)
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
func post[R any](ctx context.Context, client *transport.Client, operation, endpoint string, body any, sent, created int, opts []transport.RequestOption) ([]R, error) {
	var results []R
	resp, err := client.Send(ctx, http.MethodPost, endpoint, nil, body, &results, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != created {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if len(results) != sent {
		return nil, fmt.Errorf("moysklad: %s: response has %d items for %d sent", operation, len(results), sent)
	}
	return results, nil
}

// Submit создаёт асинхронную задачу: GET endpoint с params и async=true, успех — 202 с
// Content-Location (статус) и Location (результат). limit и offset источник с async запрещает —
// ошибка до запроса. Создание задачи — побочный эффект: transport не повторяет запрос с async=true.
func Submit(ctx context.Context, client *transport.Client, operation, endpoint string, params url.Values, opts []transport.RequestOption) (*async.Submitted, error) {
	if params.Has("limit") || params.Has("offset") {
		return nil, fmt.Errorf("moysklad: %s: limit and offset are not allowed with async", operation)
	}
	params.Set("async", "true")
	resp, err := client.Send(ctx, http.MethodGet, endpoint, params, nil, nil, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	// Content-Location — …/async/{id}, Location — …/async/{id}/result.
	status, result := resp.Header.Get("Content-Location"), resp.Header.Get("Location")
	statusURL, err1 := url.Parse(status)
	resultURL, err2 := url.Parse(result)
	id := ""
	if err1 == nil && err2 == nil {
		id = path.Base(statusURL.Path)
	}
	if !validID(id) || path.Base(path.Dir(statusURL.Path)) != "async" || resultURL.Path != statusURL.Path+"/result" {
		return nil, fmt.Errorf("moysklad: %s: response has no async task Content-Location and Location", operation)
	}
	return &async.Submitted{ID: id, StatusHref: status, ResultHref: result}, nil
}

// Проверяем форму UUID, не ограничивая его версию.
func validID(id string) bool {
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return false
	}
	_, err := hex.DecodeString(id[:8] + id[9:13] + id[14:18] + id[19:23] + id[24:])
	return err == nil
}

// Template — PUT entity/{entity}/new: шаблон документа без сохранения. basis — объект основания
// ({"customerOrder": {meta}} …); nil — пустое тело, шаблон только со стандартными полями. Ответ —
// JSON-объект документа без id; null или не объект — ошибка. PUT не повторяется.
func Template[T any](ctx context.Context, client *transport.Client, entity string, basis map[string]any, opts []transport.RequestOption) (*T, error) {
	operation := "get " + entity + " template"
	var body any
	if basis != nil {
		body = basis
	}
	var result *T
	resp, err := client.Send(ctx, http.MethodPut, "entity/"+entity+"/new", nil, body, &result, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	if result == nil {
		return nil, fmt.Errorf("moysklad: %s: response is not an object", operation)
	}
	return result, nil
}

// Trash — POST entity/{entity}/{id}/trash: удаление документа в корзину (восстановление — только
// в интерфейсе). Успех — 200 или 204, тело не читается; права и настройку корзины проверяет
// сервер. POST не повторяется.
func Trash(ctx context.Context, client *transport.Client, entity, id string, opts []transport.RequestOption) error {
	operation := "trash " + entity + " " + id
	if !validID(id) {
		return fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	resp, err := client.Send(ctx, http.MethodPost, "entity/"+entity+"/"+id+"/trash", nil, nil, nil, opts...)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("moysklad: %s: unexpected HTTP status %d", operation, resp.StatusCode)
	}
	return nil
}
