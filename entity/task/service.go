package task

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/metadata"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type Service struct {
	client *transport.Client
	meta   *metadata.Service // типы задач — общий ресурс статусов entity/task/metadata/states
}

func NewService(client *transport.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("moysklad: task service requires a client")
	}
	meta, err := metadata.NewService(client)
	if err != nil {
		return nil, err
	}
	return &Service{client: client, meta: meta}, nil
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-poluchit-zadachu
func (s *Service) Get(ctx context.Context, id string, opts ...transport.RequestOption) (*Task, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: get task: id must be a UUID")
	}
	return s.do(ctx, "get task "+id, http.MethodGet, "entity/task/"+id, nil, opts)
}

// List читает одну страницу без expand и автообхода; следующие страницы — pages.All. Без
// фильтров сервер отдаёт только активные задачи (done=false), доступные сотруднику; фильтры
// assignee, author, done, state.name и др. — query.Filter.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-poluchit-zadachi
func (s *Service) List(ctx context.Context, options query.ListOptions, opts ...transport.RequestOption) (*model.Collection[Task], error) {
	return list.Entity[Task](ctx, s.client, "task", options, nil, opts...)
}

// Create требует description и assignee.meta.href до запроса. Не повторять POST автоматически
// после неоднозначного сбоя: возможен дубль задачи.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-sozdat-zadachu
func (s *Service) Create(ctx context.Context, input CreateRequest, opts ...transport.RequestOption) (*Task, error) {
	transport.ResetResponseInfo(opts)
	if input.Description == "" {
		return nil, fmt.Errorf("moysklad: create task: description must not be empty")
	}
	if input.Assignee.Meta.Href == "" {
		return nil, fmt.Errorf("moysklad: create task: assignee must not be empty")
	}
	return s.do(ctx, "create task", http.MethodPost, "entity/task", input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-izmenit-zadachu
func (s *Service) Update(ctx context.Context, id string, input UpdateRequest, opts ...transport.RequestOption) (*Task, error) {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update task: id must be a UUID")
	}
	// Нулевые Field — «не передано»; Value(false) и Null не нулевые.
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: update task: no fields specified")
	}
	return s.do(ctx, "update task "+id, http.MethodPut, "entity/task/"+id, input, opts)
}

// Delete — DELETE entity/task/{id}. Повторов нет; права проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-udalit-zadachu
func (s *Service) Delete(ctx context.Context, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !validID(id) {
		return fmt.Errorf("moysklad: delete task: id must be a UUID")
	}
	return write.Delete(ctx, s.client, "delete task "+id, "entity/task/"+id, opts)
}

// CreateMany создаёт задачи одним POST entity/task; обязательные поля элементов проверяет
// сервер. Ответ — по правилам пакета batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-massovoe-sozdanie-i-obnovlenie-zadach
func (s *Service) CreateMany(ctx context.Context, items []CreateRequest, opts ...transport.RequestOption) ([]batch.Result[Task], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateMany[Task](ctx, s.client, "task", items, opts)
}

// UpdateMany изменяет задачи одним POST entity/task: элемент — meta по ID и поля Request.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-massovoe-sozdanie-i-obnovlenie-zadach
func (s *Service) UpdateMany(ctx context.Context, items []batch.Update[UpdateRequest], opts ...transport.RequestOption) ([]batch.Result[Task], error) {
	transport.ResetResponseInfo(opts)
	return write.UpdateMany[Task](ctx, s.client, "task", items, opts)
}

// DeleteMany удаляет задачи одним POST entity/task/delete; ответ — batch.Deleted по элементам.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-massovoe-udalenie-zadach
func (s *Service) DeleteMany(ctx context.Context, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	transport.ResetResponseInfo(opts)
	return write.DeleteMany(ctx, s.client, "task", ids, opts)
}

func (s *Service) do(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*Task, error) {
	var result Task
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
