package processingplan

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Параметрические материалы — entity/processingplan/{id}/parametricmaterials (так в запросах
// источника; href в части примеров — …/parametricmaterial). Повторов нет.

// meta.type строки в теле массового удаления — как в примере источника (processingplanmaterial),
// хотя Get возвращает processingplanparametricmaterial.
func (s *Service) parametric(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "processingplan", ParentID: id, Path: "parametricmaterials", Item: "parametric material", MetaType: "processingplanmaterial"}
}

// GetParametricMaterials читает одну страницу параметрических материалов: только limit и offset.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-poluchit-parametricheskie-materialy-tehkarty
func (s *Service) GetParametricMaterials(ctx context.Context, id string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[ParametricMaterial], error) {
	return nested.List[ParametricMaterial](ctx, s.parametric(id), page, opts)
}

// GetParametricMaterial — пример «отдельного параметрического материала» в разделе списка.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-poluchit-parametricheskie-materialy-tehkarty
func (s *Service) GetParametricMaterial(ctx context.Context, id, materialID string, opts ...transport.RequestOption) (*ParametricMaterial, error) {
	return nested.Get[ParametricMaterial](ctx, s.parametric(id), materialID, opts)
}

// CreateParametricMaterials добавляет параметрические материалы одним POST массивом; ответ — по
// правилам batch. Не повторять после неоднозначного сбоя.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-sozdat-parametricheskij-material
func (s *Service) CreateParametricMaterials(ctx context.Context, id string, items []ParametricMaterialRequest, opts ...transport.RequestOption) ([]batch.Result[ParametricMaterial], error) {
	return nested.CreateMany[ParametricMaterial](ctx, s.parametric(id), items, opts)
}

// UpdateParametricMaterial: пустой input отклоняется до запроса; переданный массив paramFeatures
// или paramQuantities заменяет прежний.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-izmenit-parametricheskij-material
func (s *Service) UpdateParametricMaterial(ctx context.Context, id, materialID string, input ParametricMaterialRequest, opts ...transport.RequestOption) (*ParametricMaterial, error) {
	return nested.Update[ParametricMaterial](ctx, s.parametric(id), materialID, input, opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-udalit-parametricheskij-material
func (s *Service) DeleteParametricMaterial(ctx context.Context, id, materialID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.parametric(id), materialID, opts)
}

// DeleteManyParametricMaterials — POST …/parametricmaterials/delete с meta материалов (1–1000 id);
// ответ без тела, результата по элементам нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-massovoe-udalenie-parametricheskih-materialov
func (s *Service) DeleteManyParametricMaterials(ctx context.Context, id string, materialIDs []string, opts ...transport.RequestOption) error {
	return nested.DeleteRows(ctx, s.parametric(id), "delete processingplan "+id+" parametric materials", materialIDs, opts)
}
