package task

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/metadata"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Типы задач — статусы сущности task (meta.type state): те же маршруты, что у
// metadata.Service с entityType "task". Отличие источника — при создании stateType не нужен,
// поэтому CreateState и CreateManyStates принимают свой StateCreateRequest.

// ListStates читает типы задач из GET entity/task/metadata (поле states); отдельного списка
// источник не описывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-poluchit-tipy-zadach
func (s *Service) ListStates(ctx context.Context, opts ...transport.RequestOption) ([]metadata.State, error) {
	m, err := s.meta.Get(ctx, "task", opts...)
	if err != nil {
		return nil, err
	}
	return m.States, nil
}

// CreateState создаёт тип задачи; name и color проверяются до запроса. Не повторять POST
// автоматически после неоднозначного сбоя: возможен дубль.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-sozdat-tip-zadachi
func (s *Service) CreateState(ctx context.Context, input StateCreateRequest, opts ...transport.RequestOption) (*metadata.State, error) {
	transport.ResetResponseInfo(opts)
	const operation = "create task state"
	if input.Name == "" || input.Color == nil {
		return nil, fmt.Errorf("moysklad: %s: name and color are required", operation)
	}
	var result metadata.State
	resp, err := s.client.Send(ctx, http.MethodPost, "entity/task/metadata/states", nil, input, &result, opts...)
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

// UpdateState меняет переданные поля типа задачи (источник: name и color).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-izmenit-tip-zadachi
func (s *Service) UpdateState(ctx context.Context, id string, input metadata.StateUpdateRequest, opts ...transport.RequestOption) (*metadata.State, error) {
	return s.meta.UpdateState(ctx, "task", id, input, opts...)
}

// DeleteState удаляет тип задачи: DELETE, успех — 200 или 204; повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-udalit-tip-zadachi
func (s *Service) DeleteState(ctx context.Context, id string, opts ...transport.RequestOption) error {
	return s.meta.DeleteState(ctx, "task", id, opts...)
}

// CreateManyStates и UpdateManyStates — один POST entity/task/metadata/states с массивом;
// правила ответа — пакет batch. Массового удаления типов источник не описывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-massovoe-sozdanie-i-obnovlenie-tipov-zadach
func (s *Service) CreateManyStates(ctx context.Context, items []StateCreateRequest, opts ...transport.RequestOption) ([]batch.Result[metadata.State], error) {
	transport.ResetResponseInfo(opts)
	return write.CreateManyNested[metadata.State](ctx, s.client, "create many task states", "entity/task/metadata/states", items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task#3-massovoe-sozdanie-i-obnovlenie-tipov-zadach
func (s *Service) UpdateManyStates(ctx context.Context, items []batch.Update[metadata.StateUpdateRequest], opts ...transport.RequestOption) ([]batch.Result[metadata.State], error) {
	return s.meta.UpdateManyStates(ctx, "task", items, opts...)
}
