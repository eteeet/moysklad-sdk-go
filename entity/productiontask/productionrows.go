package productiontask

import (
	"context"

	"github.com/eteeet/moysklad-sdk-go/internal/nested"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// ProductionRow — позиция производственного задания (техкарта в задании), meta.type
// productionrow: собственная форма, не positions.Position. ProductionVolume — объём производства,
// не деньги.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-pozicii-proizvodstvennogo-zadaniya
type ProductionRow struct {
	Meta              model.Meta       `json:"meta"`
	ID                string           `json:"id"`
	AccountID         string           `json:"accountId"`
	Name              string           `json:"name"`
	ExternalCode      string           `json:"externalCode"`
	ProcessingPlan    *model.Reference `json:"processingPlan,omitzero"`
	ProductionVariant *model.Reference `json:"productionVariant,omitzero"`
	ProductionVolume  float64          `json:"productionVolume"`
	Updated           model.DateTime   `json:"updated,omitzero"`
}

// ProductionRowRequest — позиция в CreateRequest/UpdateRequest документа. Отдельного ресурса
// создания позиции источник не описывает: позиции добавляются только телом документа. updated
// в таблице без маркера «Только для чтения», но это момент обновления, который ставит сервер, —
// SDK его не пишет.
type ProductionRowRequest struct {
	// ID сопоставляет существующую позицию в Update документа; для новой позиции его не задают.
	ID                model.Field[string]          `json:"id,omitzero"`
	Name              model.Field[string]          `json:"name,omitzero"`
	ExternalCode      model.Field[string]          `json:"externalCode,omitzero"`
	ProcessingPlan    model.Field[model.Reference] `json:"processingPlan,omitzero"`
	ProductionVariant model.Field[model.Reference] `json:"productionVariant,omitzero"`
	ProductionVolume  model.Field[float64]         `json:"productionVolume,omitzero"`
}

// ProductionRowUpdate — тело UpdateProductionRow: источник разрешает менять только объём
// производства.
type ProductionRowUpdate struct {
	ProductionVolume model.Field[float64] `json:"productionVolume,omitzero"`
}

func (s *Service) productionRows(id string) nested.Collection {
	return nested.Collection{Client: s.client, Parent: "productiontask", ParentID: id, Path: "productionrows", Item: "production row", MetaType: "productionrow"}
}

// GetProductionRows читает одну страницу позиций: только limit и offset, без автообхода.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-poluchit-spisok-pozicij-proizvodstvennogo-zadaniya
func (s *Service) GetProductionRows(ctx context.Context, id string, page query.PageOptions, opts ...transport.RequestOption) (*model.Collection[ProductionRow], error) {
	return getPage[ProductionRow](ctx, s.client, "get productiontask "+id+" production rows", id, "entity/productiontask/"+id+"/productionrows", page, "", opts)
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-poluchit-otdel-nuyu-poziciyu-proizvodstvennogo-zadaniya
func (s *Service) GetProductionRow(ctx context.Context, id, rowID string, opts ...transport.RequestOption) (*ProductionRow, error) {
	return nested.Get[ProductionRow](ctx, s.productionRows(id), rowID, opts)
}

// UpdateProductionRow меняет объём производства позиции; пустой input отклоняется до запроса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-izmenit-pozicii-proizvodstvennogo-zadaniya
func (s *Service) UpdateProductionRow(ctx context.Context, id, rowID string, input ProductionRowUpdate, opts ...transport.RequestOption) (*ProductionRow, error) {
	return nested.Update[ProductionRow](ctx, s.productionRows(id), rowID, input, opts)
}

// DeleteProductionRow удаляет позицию; позицию с начатыми этапами сервер не удаляет. Повторов нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-udalit-otdel-nuyu-poziciyu-proizvodstvennogo-zadaniya
func (s *Service) DeleteProductionRow(ctx context.Context, id, rowID string, opts ...transport.RequestOption) error {
	return nested.Delete(ctx, s.productionRows(id), rowID, opts)
}
