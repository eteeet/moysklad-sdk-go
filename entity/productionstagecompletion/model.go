package productionstagecompletion

import "github.com/eteeet/moysklad-sdk-go/model"

// ProductionStageCompletion — выполнение этапа производства. Модель чтения без expand.
// ProcessingUnitCost, LabourUnitCost и StandardHourCost — деньги (копейки, как все суммы API);
// ProductionVolume — объём производства и StandardHourUnit — нормо-часы: не деньги.
// Materials и Products — метаданные наборов (meta.type productionstagecompletionmaterial и
// productionstagecompletionresult); строки читают GetMaterials и GetProducts.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion
type ProductionStageCompletion struct {
	Meta                 model.Meta                  `json:"meta"`
	ID                   string                      `json:"id"`
	AccountID            string                      `json:"accountId"`
	Name                 string                      `json:"name"`
	Description          string                      `json:"description"`
	ExternalCode         string                      `json:"externalCode"`
	Created              model.DateTime              `json:"created,omitzero"`
	Updated              model.DateTime              `json:"updated,omitzero"`
	Moment               model.DateTime              `json:"moment,omitzero"`
	Applicable           bool                        `json:"applicable"`
	Shared               bool                        `json:"shared"`
	Defect               bool                        `json:"defect"`
	EnableHourAccounting bool                        `json:"enableHourAccounting"`
	Owner                *model.Reference            `json:"owner,omitzero"`
	Group                *model.Reference            `json:"group,omitzero"`
	Performer            *model.Reference            `json:"performer,omitzero"`
	Service              *model.Reference            `json:"service,omitzero"`
	ProductionStage      *model.Reference            `json:"productionStage,omitzero"`
	ProductionVolume     float64                     `json:"productionVolume"`
	ProcessingUnitCost   float64                     `json:"processingUnitCost"`
	LabourUnitCost       float64                     `json:"labourUnitCost"`
	StandardHourCost     float64                     `json:"standardHourCost"`
	StandardHourUnit     float64                     `json:"standardHourUnit"`
	Materials            *model.Collection[Material] `json:"materials,omitzero"`
	Products             *model.Collection[Product]  `json:"products,omitzero"`
}

// Material — материал выполнения этапа (meta.type productionstagecompletionmaterial): товар,
// модификация или партия, затраченные на этапе. ConsumedQuantity — количество, не деньги.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-materialy-vypolneniya-etapa-proizvodstva
type Material struct {
	Meta             model.Meta       `json:"meta"`
	ID               string           `json:"id"`
	AccountID        string           `json:"accountId"`
	Assortment       *model.Reference `json:"assortment,omitzero"`
	ConsumedQuantity float64          `json:"consumedQuantity"`
	Things           []string         `json:"things"`
}

// Product — продукт выполнения этапа (meta.type productionstagecompletionresult): товар,
// модификация или партия, созданные на этапе; есть только у последнего этапа.
// ProducedQuantity — количество, не деньги.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-produkty-vypolneniya-etapa-proizvodstva
type Product struct {
	Meta             model.Meta       `json:"meta"`
	ID               string           `json:"id"`
	AccountID        string           `json:"accountId"`
	Assortment       *model.Reference `json:"assortment,omitzero"`
	ProducedQuantity float64          `json:"producedQuantity"`
	Things           []string         `json:"things"`
}
