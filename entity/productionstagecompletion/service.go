package productionstagecompletion

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

const entity = "productionstagecompletion"

type Service struct {
	client *transport.Client
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: productionstagecompletion service requires a client")
	}
	return &Service{client: client}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-poluchit-vypolnenie-etapa-proizvodstva
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*ProductionStageCompletion, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get productionstagecompletion: id must be a UUID")
	}
	return s.do(ctx, "get productionstagecompletion "+id, http.MethodGet, "entity/"+entity+"/"+id, nil, opts)
}

// Create: производство по производственному заданию должно быть начато (проверяет сервер).
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль документа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-sozdat-vypolnenie-etapa-proizvodstva
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*ProductionStageCompletion, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create productionstagecompletion"
	if input.ProductionStage.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: %s: productionStage meta.href must not be empty", operation)
	}
	return s.do(ctx, operation, http.MethodPost, "entity/"+entity, input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-izmenit-vypolnenie-etapa-proizvodstva
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*ProductionStageCompletion, error) {
	transport.ResetResponseInfo(opts)
	const operation = "update productionstagecompletion"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: id must be a UUID", operation)
	}
	// Нулевые Field — «не передано»; Value(false), Value(0.0) и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	return s.do(ctx, operation+" "+id, http.MethodPut, "entity/"+entity+"/"+id, input, opts)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All. Источник
// называет фильтры productionTask=<href> и assortment=<href> — передаются через options.Filter.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-poluchit-spisok-vypolnenij-etapov-proizvodstva
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[ProductionStageCompletion], error) {
	return list.Entity[ProductionStageCompletion](ctx, s.client, entity, options, nil, opts...)
}

// Delete — DELETE entity/productionstagecompletion/{id}, не в корзину. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-udalit-vypolnenie-etapa-proizvodstva
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete productionstagecompletion: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete productionstagecompletion "+id, "entity/"+entity+"/"+id, opts)
}

// CreateMany создаёт объекты одним POST entity/productionstagecompletion: массив до batch.MaxItems
// элементов. Обязательные поля проверяет сервер у каждого элемента; ответ, ошибки элементов и
// неизвестный исход при сбое — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-massovoe-sozdanie-i-obnovlenie-vypolnenij-etapov-proizvodstva
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[ProductionStageCompletion], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[ProductionStageCompletion](ctx, s.client, entity, items, opts)
}

// UpdateMany изменяет объекты одним POST entity/productionstagecompletion: элемент — meta по ID и
// поля Request. До запроса проверяются только UUID и непустой Request. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-massovoe-sozdanie-i-obnovlenie-vypolnenij-etapov-proizvodstva
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[ProductionStageCompletion], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[ProductionStageCompletion](ctx, s.client, entity, items, opts)
}

// DeleteMany удаляет объекты по ids одним POST entity/productionstagecompletion/delete, не в
// корзину; до batch.MaxItems элементов. Сервер удаляет в порядке ids: выполнение последнего
// этапа должно идти раньше выполнений этапов, от которых оно зависит. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-massovoe-udalenie-vypolnenij-etapov-proizvodstva
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, entity, ids, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*ProductionStageCompletion, error) {
	var result ProductionStageCompletion
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
