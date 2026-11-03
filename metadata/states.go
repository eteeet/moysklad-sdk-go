package metadata

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/write"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Типы статуса; набор открыт.
const (
	StateRegular      = "Regular"
	StateSuccessful   = "Successful"
	StateUnsuccessful = "Unsuccessful"
)

// StateCreateRequest: name, color и stateType необходимы при создании. Color — ARGB целым
// числом из 4 байт (rgb(162, 198, 23) → 10667543); указатель отличает «не задан» от 0.
// Диапазон проверяет сервер. Таблица называет Regular значением по умолчанию, но помечает
// stateType необходимым — SDK следует пометке.
type StateCreateRequest struct {
	Name      string `json:"name"`
	Color     *int64 `json:"color"`
	StateType string `json:"stateType"`
}

// StateUpdateRequest: одно или несколько полей name, color, stateType. null источник не
// описывает: Field позволяет его отправить, решает сервер.
type StateUpdateRequest struct {
	Name      model.Field[string] `json:"name,omitzero"`
	Color     model.Field[int64]  `json:"color,omitzero"`
	StateType model.Field[string] `json:"stateType,omitzero"`
}

// CreateState создаёт статус типа entityType. Отдельного GET статуса API не описывает:
// статусы читаются в Get(entityType).States. Права администратора проверяет сервер.
// Не повторять POST автоматически после неоднозначного сбоя: возможен дубль статуса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/states#3-sozdat-status
func (s *Service) CreateState(ctx context.Context, entityType string, input StateCreateRequest, opts ...transport.RequestOption) (*State, error) {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: create state: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	operation := "create " + entityType + " state"
	if input.Name == "" || input.Color == nil || input.StateType == "" {
		return nil, fmt.Errorf("moysklad: %s: name, color and stateType are required", operation)
	}
	return s.state(ctx, operation, http.MethodPost, "entity/"+entityType+"/metadata/states", input, opts)
}

// UpdateState меняет переданные поля статуса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/states#3-izmenit-status
func (s *Service) UpdateState(ctx context.Context, entityType, id string, input StateUpdateRequest, opts ...transport.RequestOption) (*State, error) {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: update state: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: update state: id %q must be a UUID", id)
	}
	operation := "update " + entityType + " state " + id
	if reflect.ValueOf(input).IsZero() {
		return nil, fmt.Errorf("moysklad: %s: no fields specified", operation)
	}
	return s.state(ctx, operation, http.MethodPut, "entity/"+entityType+"/metadata/states/"+id, input, opts)
}

// DeleteState удаляет статус типа entityType: DELETE, успех — 200 или 204; повторов нет.
// Права администратора проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/states#3-udalit-status
func (s *Service) DeleteState(ctx context.Context, entityType, id string, opts ...transport.RequestOption) error {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return fmt.Errorf("moysklad: delete state: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	if !validID(id) {
		return fmt.Errorf("moysklad: delete state: id %q must be a UUID", id)
	}
	return write.Delete(ctx, s.client, "delete "+entityType+" state "+id, "entity/"+entityType+"/metadata/states/"+id, opts)
}

// state — критерий успеха записи статуса: id.
func (s *Service) state(ctx context.Context, operation, method, endpoint string, body any, opts []transport.RequestOption) (*State, error) {
	var result State
	if err := s.do(ctx, operation, method, endpoint, nil, body, &result, opts); err != nil {
		return nil, err
	}
	if result.ID == "" {
		return nil, fmt.Errorf("moysklad: %s: response has no id", operation)
	}
	return &result, nil
}

// CreateManyStates и UpdateManyStates — один POST entity/{type}/metadata/states с массивом:
// элемент без meta создаёт статус, с meta (тип state) — изменяет. До запроса проверяются только
// тип сущности, размер массива, UUID и непустой запрос; поля — сервер (nil Color уходит как null). Правила ответа — пакет
// batch; повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/states#3-massovoe-sozdanie-i-obnovlenie-statusov
func (s *Service) CreateManyStates(ctx context.Context, entityType string, items []StateCreateRequest, opts ...transport.RequestOption) ([]batch.Result[State], error) {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: create many states: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	return write.CreateManyNested[State](ctx, s.client, "create many "+entityType+" states", "entity/"+entityType+"/metadata/states", items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/states#3-massovoe-sozdanie-i-obnovlenie-statusov
func (s *Service) UpdateManyStates(ctx context.Context, entityType string, items []batch.Update[StateUpdateRequest], opts ...transport.RequestOption) ([]batch.Result[State], error) {
	transport.ResetResponseInfo(opts)
	if !entityTypePattern.MatchString(entityType) {
		return nil, fmt.Errorf("moysklad: update many states: entity type %q must match [a-z][a-z0-9]*", entityType)
	}
	return write.UpdateManyNested[State](ctx, s.client, "update many "+entityType+" states", "entity/"+entityType+"/metadata/states", "state", items, opts)
}
