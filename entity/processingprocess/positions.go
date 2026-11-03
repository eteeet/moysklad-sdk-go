package processingprocess

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Позиции (этапы) — entity/processingprocess/{id}/positions, meta.type processingprocessposition.
// Лимит — 100 позиций на техпроцесс (проверяет сервер). Повторов записи нет.

func (s *Service) positions(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "processingprocess", ParentID: id, Path: "positions", Item: "position", MetaType: "processingprocessposition"}
}

// GetPositions читает одну страницу позиций: только limit и offset (fields отклоняется до запроса).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-poluchit-pozicii-tehprocessa
func (s *Service) GetPositions(ctx context.Context, id string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	transport.ResetResponseInfo(opts)
	operation := "get processingprocess " + id + " positions"
	if !validID(id) {
		return nil, fmt.Errorf("moysklad: %s: processingprocess id must be a UUID", operation)
	}
	if page.Fields != "" {
		return nil, fmt.Errorf("moysklad: %s: fields=%q is not supported", operation, page.Fields)
	}
	params, err := page.Values()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return list.Page[Position](ctx, s.client, operation, "entity/processingprocess/"+id+"/positions", params, opts...)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-poluchit-otdelnuyu-poziciyu-tehprocessa
func (s *Service) GetPosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) (*Position, error) {
	return nested.Get[Position](ctx, s.positions(id), positionID, opts)
}

// CreatePositions добавляет позиции одним POST массивом в порядке среза (одна позиция — массив из
// одного элемента); processingstage каждой проверяется до запроса, ответ — по правилам batch.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-sozdat-pozicii-tehprocessa
func (s *Service) CreatePositions(ctx context.Context, id string, items []PositionRequest, opts ...transport.RequestOption) ([]batch.Result[Position], error) {
	for i, item := range items {
		if item.ProcessingStage.Meta.Href == "" {
			transport.ResetResponseInfo(opts)
			return nil, fmt.Errorf("moysklad: create many processingprocess %s positions: items[%d]: processingstage meta.href must not be empty", id, i)
		}
	}
	return nested.CreateMany[Position](ctx, s.positions(id), items, opts)
}

// UpdatePosition: смена этапа изменяет позицию, а не пересоздаёт её (так в источнике). Пустой
// request и заданный Meta отклоняются до запроса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-izmenit-poziciyu-tehprocessa
func (s *Service) UpdatePosition(ctx context.Context, id, positionID string, input PositionUpdateRequest, opts ...transport.RequestOption) (*Position, error) {
	if !input.Meta.IsZero() {
		transport.ResetResponseInfo(opts)
		return nil, fmt.Errorf("moysklad: update processingprocess %s position %s: meta addresses rows only in UpdateRequest.Positions", id, positionID)
	}
	return nested.Update[Position](ctx, s.positions(id), positionID, input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-udalit-poziciyu-tehprocessa
func (s *Service) DeletePosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.positions(id), positionID, opts)
}

// DeleteManyPositions — POST …/positions/delete с meta позиций (1–1000 id). Источник отвечает 200
// без тела, поэтому результата по элементам нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess#3-massovoe-udalenie-pozicij-tehprocessa
func (s *Service) DeleteManyPositions(ctx context.Context, id string, positionIDs []string, opts ...transport.RequestOption) error {
	return nested.DeleteRows(ctx, s.positions(id), "delete many processingprocess "+id+" positions", positionIDs, opts)
}
