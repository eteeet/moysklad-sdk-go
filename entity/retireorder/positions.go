package retireorder

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Отдельные операции со строками: entity/retireorder/{id}/positions[/{positionID}]. Повторов нет.

func (s *Service) rows() *positions.Service[Position] {
	// NewService падает только на nil client и неверном коде документа — здесь их нет.
	reader, _ := positions.NewService[Position](s.client, "retireorder")
	return reader
}

func (s *Service) collection(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "retireorder", ParentID: id, Path: "positions", Item: "position", MetaType: "retireorderposition"}
}

// GetPositions читает одну страницу строк без скрытого Get и автообхода; fields источник для
// этого документа не документирует.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retireorder#3-poluchit-pozicii-vyvoda-iz-oborota
func (s *Service) GetPositions(ctx context.Context, id string, page positions.PageOptions, opts ...transport.RequestOption) (*model.Collection[Position], error) {
	return s.rows().GetPage(ctx, id, page, opts...)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retireorder#3-poluchit-poziciyu
func (s *Service) GetPosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) (*Position, error) {
	return s.rows().GetPosition(ctx, id, positionID, opts...)
}

// CreatePositions добавляет строки к существующим одним POST массивом (одна строка — массив из
// одного элемента). assortment и trackingCodes таблица помечает «Необходимо при создании» —
// строка без них отклоняется до запроса. Ответ — по правилам batch.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retireorder#3-sozdat-poziciyu
func (s *Service) CreatePositions(ctx context.Context, id string, items []PositionRequest, opts ...transport.RequestOption) ([]batch.Result[Position], error) {
	transport.ResetResponseInfo(opts)
	for i, item := range items {
		ref, ok := item.Assortment.Get()
		codes, _ := item.TrackingCodes.Get()
		if !ok || ref.Meta.Href == "" || len(codes) == 0 {
			return nil, fmt.Errorf("moysklad: create retireorder %s positions: items[%d]: assortment meta.href and trackingCodes must be set", id, i)
		}
	}
	return nested.CreateMany[Position](ctx, s.collection(id), items, opts)
}

// UpdatePosition изменяет переданные поля одной строки; пустой запрос — ошибка до HTTP.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retireorder#3-izmenit-poziciyu
func (s *Service) UpdatePosition(ctx context.Context, id, positionID string, input PositionRequest, opts ...transport.RequestOption) (*Position, error) {
	return nested.Update[Position](ctx, s.collection(id), positionID, input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retireorder#3-udalit-poziciyu
func (s *Service) DeletePosition(ctx context.Context, id, positionID string, opts ...transport.RequestOption) error {
	return s.rows().DeletePosition(ctx, id, positionID, opts...)
}

// DeleteManyPositions — POST …/positions/delete; ответ без тела, результата по элементам нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retireorder#3-massovoe-udalenie-pozicij
func (s *Service) DeleteManyPositions(ctx context.Context, id string, positionIDs []string, opts ...transport.RequestOption) error {
	return s.rows().DeleteManyPositions(ctx, id, positionIDs, opts...)
}
