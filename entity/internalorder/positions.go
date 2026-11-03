package internalorder

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Отдельные операции со строками: entity/internalorder/{id}/positions[/{positionID}]. Это
// единственный способ записи строк: positions в таблице документа помечен readonly. Повторов нет.

func (s *Service) rows() *positions.Service[Position] {
	// NewService падает только на nil client и неверном коде документа — здесь их нет.
	reader, _ := positions.NewService[Position](s.client, "internalorder")
	return reader
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-poluchit-poziciyu
func (s *Service) GetPosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) (*Position, error) {
	return s.rows().GetPosition(ctx, id, positionID, opts...)
}

// CreatePositions добавляет строки к существующим одним POST массивом; assortment и
// положительный quantity проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-sozdat-pozicii-vnutrennego-zakaza
func (s *Service) CreatePositions(ctx context.Context, id string, items []PositionRequest, opts ...transport.RequestOption) ([]batch.Result[Position], error) {
	return s.rows().CreatePositions(ctx, id, items, opts...)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-izmenit-poziciyu
func (s *Service) UpdatePosition(ctx context.Context, id, positionID string, input PositionRequest, opts ...transport.RequestOption) (*Position, error) {
	return s.rows().UpdatePosition(ctx, id, positionID, input, opts...)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-udalit-poziciyu
func (s *Service) DeletePosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) error {
	return s.rows().DeletePosition(ctx, id, positionID, opts...)
}

// DeleteManyPositions — POST …/positions/delete; ответ без тела, результата по элементам нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-massovoe-udalenie-pozicij
func (s *Service) DeleteManyPositions(ctx context.Context, id string, positionIDs []string, opts ...transport.RequestOption) error {
	return s.rows().DeleteManyPositions(ctx, id, positionIDs, opts...)
}
