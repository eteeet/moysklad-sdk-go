package productiontask

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Material — материал производственного этапа (meta.type productiontaskmaterial). PlanQuantity —
// количество, не деньги.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-materialy-proizvodstvennogo-etapa
type Material struct {
	Meta         model.Meta       `json:"meta"`
	ID           string           `json:"id"`
	AccountID    string           `json:"accountId"`
	Assortment   *model.Reference `json:"assortment,omitzero"`
	PlanQuantity float64          `json:"planQuantity"`
}

// MaterialRequest — материал в CreateMaterial, UpdateMaterial и StageUpdateRequest.Materials.
// Для нового материала сервер требует assortment и planQuantity; в UpdateMaterial обязательных
// полей нет.
type MaterialRequest struct {
	// ID сопоставляет существующий материал в UpdateStage; для нового материала его не задают.
	ID           model.Field[string]          `json:"id,omitzero"`
	Assortment   model.Field[model.Reference] `json:"assortment,omitzero"`
	PlanQuantity model.Field[float64]         `json:"planQuantity,omitzero"`
}

func (s *Service) materials(stageID string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "productionstage", ParentID: stageID, Path: "materials", Item: "material", MetaType: "productiontaskmaterial"}
}

// GetMaterials читает одну страницу материалов этапа stageID: только limit и offset.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-poluchit-materialy-proizvodstvennogo-etapa
func (s *Service) GetMaterials(ctx context.Context, stageID string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[Material], error) {
	return getPage[Material](ctx, s.client, "get productionstage "+stageID+" materials", stageID, "entity/productionstage/"+stageID+"/materials", page, "", opts)
}

// CreateMaterial добавляет один материал к этапу: POST объектом. Пример ответа источника —
// массив из одного объекта, принимается и объект. Не повторять после неоднозначного сбоя.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-dobavit-material-k-proizvodstvennomu-etapu
func (s *Service) CreateMaterial(ctx context.Context, stageID string, input MaterialRequest, opts ...transport.RequestOption) (*Material, error) {
	transport.ResetResponseInfo(opts)
	if !newRow(input.Assortment, input.PlanQuantity) {
		return nil, fmt.Errorf("moysklad: create productionstage %s material: assortment meta.href and planQuantity must be set", stageID)
	}
	return nested.Create[Material](ctx, s.materials(stageID), input, opts)
}

// newRow — для новой строки источник требует assortment и planQuantity («Для успешного создания
// необходимо…»); положительность количества проверяет сервер.
func newRow(assortment model.Field[model.Reference], quantity model.Field[float64]) bool {
	ref, ok := assortment.Get()
	return ok && ref.Meta.Href != "" && !quantity.IsZero() && !quantity.IsNull()
}

// UpdateMaterial: пустой input отклоняется до запроса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-izmenit-material-proizvodstvennogo-etapa
func (s *Service) UpdateMaterial(ctx context.Context, stageID, materialID string, input MaterialRequest, opts ...transport.RequestOption) (*Material, error) {
	return nested.Update[Material](ctx, s.materials(stageID), materialID, input, opts)
}

// DeleteMaterial: при резерве последний материал задания и материал этапа с выполнением сервер не
// удаляет. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-udalit-material-proizvodstvennogo-etapa
func (s *Service) DeleteMaterial(ctx context.Context, stageID, materialID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.materials(stageID), materialID, opts)
}
