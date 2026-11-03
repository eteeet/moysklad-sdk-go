package inventory

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Отдельные операции со строками: entity/inventory/{id}/positions[/{positionID}]. Форма строки
// своя, поэтому создание и изменение идут через nested с PositionRequest этого пакета; чтение
// страницы и массовое удаление (ответ без тела) — через generic positions.Service. Повторов нет.
// Пересчёт расчётного остатка (rpc/inventory/{id}/recalcCalculatedQuantity) в ticket O09 не входит.

func (s *Service) rows() *positions.Service[Position] {
	// NewService падает только на nil client и неверном коде документа — здесь их нет.
	reader, _ := positions.NewService[Position](s.client, "inventory")
	return reader
}

func (s *Service) positions(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "inventory", ParentID: id, Path: "positions", Item: "position", MetaType: "inventoryposition"}
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-poluchit-poziciyu
func (s *Service) GetPosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) (*Position, error) {
	return nested.Get[Position](ctx, s.positions(id), positionID, opts)
}

// CreatePositions добавляет строки к существующим одним POST массивом (одна строка — массив из
// одного элемента); assortment и положительный quantity проверяет сервер. Ответ — по правилам
// batch.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-sozdat-poziciyu
func (s *Service) CreatePositions(ctx context.Context, id string, items []PositionRequest, opts ...transport.RequestOption) ([]batch.Result[Position], error) {
	return nested.CreateMany[Position](ctx, s.positions(id), items, opts)
}

// UpdatePosition изменяет переданные поля строки; пустой запрос — ошибка до HTTP.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-izmenit-poziciyu
func (s *Service) UpdatePosition(ctx context.Context, id, positionID string, input PositionRequest, opts ...transport.RequestOption) (*Position, error) {
	return nested.Update[Position](ctx, s.positions(id), positionID, input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-udalit-poziciyu
func (s *Service) DeletePosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.positions(id), positionID, opts)
}

// DeleteManyPositions — POST …/positions/delete; ответ без тела, результата по элементам нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-massovoe-udalenie-pozicij
func (s *Service) DeleteManyPositions(ctx context.Context, id string, positionIDs []string, opts ...transport.RequestOption) error {
	return s.rows().DeletePositions(ctx, id, positionIDs, opts...)
}
