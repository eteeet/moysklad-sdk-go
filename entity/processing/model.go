package processing

import "github.com/eteeet/moysklad-sdk-go/model"

// Processing — техоперация. Модель чтения без expand: ProcessingSum — затраты на производство за
// единицу объёма в копейках (в таблице Int, в примерах бывает дробным — SDK держит float64);
// Quantity — объём производства, не деньги. Материалы — GetMaterials, продукты — GetProducts.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing
type Processing struct {
	Meta                model.Meta                         `json:"meta"`
	ID                  string                             `json:"id"`
	AccountID           string                             `json:"accountId"`
	Name                string                             `json:"name"`
	Description         string                             `json:"description"`
	Code                string                             `json:"code"`
	ExternalCode        string                             `json:"externalCode"`
	SyncID              string                             `json:"syncId"`
	Created             model.DateTime                     `json:"created,omitzero"`
	Updated             model.DateTime                     `json:"updated,omitzero"`
	Deleted             model.DateTime                     `json:"deleted,omitzero"`
	Moment              model.DateTime                     `json:"moment,omitzero"`
	Applicable          bool                               `json:"applicable"`
	Shared              bool                               `json:"shared"`
	Printed             bool                               `json:"printed"`
	Published           bool                               `json:"published"`
	Organization        *model.Reference                   `json:"organization,omitzero"`
	OrganizationAccount *model.Reference                   `json:"organizationAccount,omitzero"`
	Owner               *model.Reference                   `json:"owner,omitzero"`
	Group               *model.Reference                   `json:"group,omitzero"`
	Project             *model.Reference                   `json:"project,omitzero"`
	State               *model.Reference                   `json:"state,omitzero"`
	Attributes          []model.Attribute                  `json:"attributes"`
	ProcessingPlan      *model.Reference                   `json:"processingPlan,omitzero"`
	ProcessingSum       float64                            `json:"processingSum"`
	Quantity            float64                            `json:"quantity"`
	MaterialsStore      *model.Reference                   `json:"materialsStore,omitzero"`
	ProductsStore       *model.Reference                   `json:"productsStore,omitzero"`
	Materials           *model.Collection[Position]        `json:"materials,omitzero"`
	Products            *model.Collection[Position]        `json:"products,omitzero"`
	Files               *model.Collection[model.Reference] `json:"files,omitzero"`
	ProcessingOrder     *model.Reference                   `json:"processingOrder,omitzero"`
}

// Position — строка материалов (meta.type processingpositionmaterial) или продуктов
// (processingpositionresult). Таблицы источника для обоих наборов совпадают; это собственная форма,
// не positions.Position: цены, НДС и упаковки у строки нет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#4-materialy-tehoperacii
type Position struct {
	Meta       model.Meta       `json:"meta"`
	ID         string           `json:"id"`
	AccountID  string           `json:"accountId"`
	Assortment *model.Reference `json:"assortment,omitzero"`
	Quantity   float64          `json:"quantity"`
	Slot       *model.Reference `json:"slot,omitzero"`
}
