package processingorder

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Отдельные операции со строками: entity/processingorder/{id}/positions[/{positionID}]. Повторов
// нет. Раздела «Создать позицию» источник не содержит, поэтому CreatePositions нет: строки
// передаются в Create/Update документа.

func (s *Service) rows() *positions.Service[Position] {
	// NewService падает только на nil client и неверном коде документа — здесь их нет.
	reader, _ := positions.NewService[Position](s.client, "processingorder")
	return reader
}

// GetPositions читает одну страницу строк без скрытого Get и автообхода; fields источник для этого
// документа не документирует — отклоняется до запроса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processingorder#3-poluchit-pozicii-zakaza-na-proizvodstvo
func (s *Service) GetPositions(ctx context.Context, id string, page positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	return s.rows().GetPage(ctx, id, page, opts...)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processingorder#3-poluchit-poziciyu
func (s *Service) GetPosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) (*Position, error) {
	return s.rows().GetPosition(ctx, id, positionID, opts...)
}

// UpdatePosition изменяет переданные поля одной строки; пустой запрос — ошибка до HTTP.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processingorder#3-izmenit-poziciyu
func (s *Service) UpdatePosition(ctx context.Context, id, positionID string, input PositionRequest, opts ...transport.RequestOption) (*Position, error) {
	rows := nested.Collection{Client: s.client, Parent: "processingorder", ParentID: id, Path: "positions", Item: "position", MetaType: "processingorderposition"}
	return nested.Update[Position](ctx, rows, positionID, input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processingorder#3-udalit-poziciyu
func (s *Service) DeletePosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) error {
	return s.rows().DeletePosition(ctx, id, positionID, opts...)
}

// DeleteManyPositions — POST …/positions/delete с meta строк (1–1000); ответ без тела, результата
// по элементам нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processingorder#3-massovoe-udalenie-pozicij
func (s *Service) DeleteManyPositions(ctx context.Context, id string, positionIDs []string, opts ...transport.RequestOption) error {
	return s.rows().DeleteManyPositions(ctx, id, positionIDs, opts...)
}
