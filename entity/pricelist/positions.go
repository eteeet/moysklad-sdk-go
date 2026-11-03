package pricelist

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Отдельные операции со строками: entity/pricelist/{id}/positions[/{positionID}], meta.type
// строки — pricelistrow. Повторов нет.

func (s *Service) rows() *positions.Service[Position] {
	// NewService падает только на nil client и неверном коде документа — здесь их нет.
	reader, _ := positions.NewService[Position](s.client, "pricelist")
	return reader
}

func (s *Service) collection(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "pricelist", ParentID: id, Path: "positions", Item: "position", MetaType: "pricelistrow"}
}

// GetPositions читает одну страницу строк без скрытого Get и автообхода (общий reader страниц
// документов, тип строки — свой); fields не поддержан.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/pricelist#3-poluchit-pozicii-prajs-lista
func (s *Service) GetPositions(ctx context.Context, id string, page positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	transport.ResetResponseInfo(opts)
	return s.rows().GetPage(ctx, id, page, opts...)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/pricelist#3-poluchit-poziciyu
func (s *Service) GetPosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) (*Position, error) {
	return nested.Get[Position](ctx, s.collection(id), positionID, opts)
}

// CreatePositions добавляет строки к существующим одним POST массивом (одна строка — массив из
// одного элемента); assortment проверяет сервер. Ответ — по правилам batch.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/pricelist#3-sozdat-poziciyu
func (s *Service) CreatePositions(ctx context.Context, id string, items []PositionRequest, opts ...transport.RequestOption) ([]batch.Result[Position], error) {
	return nested.CreateMany[Position](ctx, s.collection(id), items, opts)
}

// UpdatePosition: пустой запрос отклоняется до HTTP.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/pricelist#3-izmenit-poziciyu
func (s *Service) UpdatePosition(ctx context.Context, id, positionID string, input PositionRequest, opts ...transport.RequestOption) (*Position, error) {
	return nested.Update[Position](ctx, s.collection(id), positionID, input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/pricelist#3-udalit-poziciyu
func (s *Service) DeletePosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.collection(id), positionID, opts)
}

// DeleteManyPositions — POST …/positions/delete с meta строк (pricelistrow); 1–1000 строк.
// Источник отвечает 200 с пустым телом, поэтому результата по элементам нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/pricelist#3-massovoe-udalenie-pozicij
func (s *Service) DeleteManyPositions(ctx context.Context, id string, positionIDs []string, opts ...transport.RequestOption) error {
	return s.rows().DeleteManyPositions(ctx, id, positionIDs, opts...)
}
