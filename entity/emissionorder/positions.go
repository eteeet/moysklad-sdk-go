package emissionorder

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Отдельные операции со строками: entity/emissionorder/{id}/positions[/{positionID}]. Удаления
// строк источник не описывает. Повторов нет.

func (s *Service) rows() *positions.Service[Position] {
	// NewService падает только на nil client и неверном коде документа — здесь их нет.
	reader, _ := positions.NewService[Position](s.client, "emissionorder")
	return reader
}

func (s *Service) collection(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "emissionorder", ParentID: id, Path: "positions", Item: "position", MetaType: "emissionorderposition"}
}

// GetPositions читает одну страницу строк без скрытого Get и автообхода; fields источник для
// этого документа не документирует.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/emissionorder#3-poluchit-pozicii-zakaza-kodov-markirovki
func (s *Service) GetPositions(ctx context.Context, id string, page positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	return s.rows().GetPage(ctx, id, page, opts...)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/emissionorder#3-poluchit-poziciyu-zakaza-kodov-markirovki
func (s *Service) GetPosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) (*Position, error) {
	return s.rows().GetPosition(ctx, id, positionID, opts...)
}

// CreatePositions добавляет строки к существующим одним POST массивом (одна строка — массив из
// одного элемента); assortment и положительный quantity проверяет сервер. Ответ — по правилам
// batch.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/emissionorder#3-sozdat-poziciyu
func (s *Service) CreatePositions(ctx context.Context, id string, items []PositionRequest, opts ...transport.RequestOption) ([]batch.Result[Position], error) {
	return nested.CreateMany[Position](ctx, s.collection(id), items, opts)
}

// UpdatePosition изменяет переданные поля одной строки; пустой запрос — ошибка до HTTP.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/emissionorder#3-izmenit-poziciyu
func (s *Service) UpdatePosition(ctx context.Context, id, positionID string, input PositionRequest, opts ...transport.RequestOption) (*Position, error) {
	return nested.Update[Position](ctx, s.collection(id), positionID, input, opts)
}
