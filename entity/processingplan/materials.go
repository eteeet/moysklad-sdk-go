package processingplan

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Материалы — entity/processingplan/{id}/materials. Повторов нет.

func (s *Service) materials(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "processingplan", ParentID: id, Path: "materials", Item: "material", MetaType: "processingplanmaterial"}
}

// GetMaterials читает одну страницу материалов: только limit и offset, без автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-poluchit-materialy-tehkarty
func (s *Service) GetMaterials(ctx context.Context, id string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Material], error) {
	return nested.List[Material](ctx, s.materials(id), page, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-poluchit-material
func (s *Service) GetMaterial(ctx context.Context, id, materialID string, opts ...transport.RequestOption) (*Material, error) {
	return nested.Get[Material](ctx, s.materials(id), materialID, opts)
}

// CreateMaterials добавляет материалы к существующим одним POST массивом (один материал — массив
// из одного элемента); ответ — по правилам batch. Не повторять после неоднозначного сбоя.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-sozdat-material
func (s *Service) CreateMaterials(ctx context.Context, id string, items []MaterialRequest, opts ...transport.RequestOption) ([]batch.Result[Material], error) {
	return nested.CreateMany[Material](ctx, s.materials(id), items, opts)
}

// UpdateMaterial: обязательных полей нет; пустой input отклоняется до запроса. Пример ответа
// источника — массив из одного объекта, принимается и объект.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-izmenit-material
func (s *Service) UpdateMaterial(ctx context.Context, id, materialID string, input MaterialRequest, opts ...transport.RequestOption) (*Material, error) {
	return nested.Update[Material](ctx, s.materials(id), materialID, input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-udalit-material
func (s *Service) DeleteMaterial(ctx context.Context, id, materialID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.materials(id), materialID, opts)
}

// DeleteManyMaterials — POST …/materials/delete с meta материалов (1–1000 id). Источник отвечает
// 200 с пустым телом, поэтому результата по элементам нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-massovoe-udalenie-materialov
func (s *Service) DeleteManyMaterials(ctx context.Context, id string, materialIDs []string, opts ...transport.RequestOption) error {
	return nested.DeleteRows(ctx, s.materials(id), "delete processingplan "+id+" materials", materialIDs, opts)
}
