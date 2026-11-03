package aggregatepack

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Position — позиция транспортной упаковки (meta.type aggregatepackposition): товар, партия,
// модификация или комплект. Это собственная форма, не positions.Position документов.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/handling-unit#3-pozicii-transportnoj-upakovki
type Position struct {
	Meta       model.Meta       `json:"meta"`
	ID         string           `json:"id"`
	AccountID  string           `json:"accountId"`
	Assortment *model.Reference `json:"assortment,omitzero"`
	Quantity   float64          `json:"quantity"`
	Things     []string         `json:"things"`
}

// PositionRequest — новая позиция: assortment и положительный quantity (проверяет сервер).
type PositionRequest struct {
	Assortment model.Reference       `json:"assortment"`
	Quantity   float64               `json:"quantity"`
	Things     model.Field[[]string] `json:"things,omitzero"`
}

// PositionUpdateRequest: обязательных полей нет — только изменяемые.
type PositionUpdateRequest struct {
	Assortment model.Field[model.Reference] `json:"assortment,omitzero"`
	Quantity   model.Field[float64]         `json:"quantity,omitzero"`
	Things     model.Field[[]string]        `json:"things,omitzero"`
}

func (s *Service) positions(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "aggregatepack", ParentID: id, Path: "positions", Item: "position", MetaType: "aggregatepackposition"}
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/handling-unit#3-poluchit-pozicii-transportnoj-upakovki
func (s *Service) ListPositions(ctx context.Context, id string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	return nested.List[Position](ctx, s.positions(id), page, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/handling-unit#3-poluchit-poziciyu
func (s *Service) GetPosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) (*Position, error) {
	return nested.Get[Position](ctx, s.positions(id), positionID, opts)
}

// CreatePositions добавляет позиции к существующим одним POST массивом (одна позиция — массив
// из одного элемента); ответ — по правилам batch. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/handling-unit#3-sozdat-poziciyu
func (s *Service) CreatePositions(ctx context.Context, id string, items []PositionRequest, opts ...transport.RequestOption) ([]batch.Result[Position], error) {
	return nested.CreateMany[Position](ctx, s.positions(id), items, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/handling-unit#3-izmenit-poziciyu
func (s *Service) UpdatePosition(ctx context.Context, id, positionID string, input PositionUpdateRequest, opts ...transport.RequestOption) (*Position, error) {
	return nested.Update[Position](ctx, s.positions(id), positionID, input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/handling-unit#3-udalit-poziciyu
func (s *Service) DeletePosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.positions(id), positionID, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/handling-unit#3-massovoe-udalenie-pozicij
func (s *Service) DeleteManyPositions(ctx context.Context, id string, ids []string, opts ...transport.RequestOption) ([]batch.Deleted, error) {
	return nested.DeleteMany(ctx, s.positions(id), ids, opts)
}
