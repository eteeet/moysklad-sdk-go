package processingplan

import "github.com/eteeet/moysklad-sdk-go/model"

// ProcessingPlan — техкарта. Модель чтения без expand: Cost — стоимость производства, float64 в
// копейках; CostDistributionType — открытый enum (BY_PRICE, BY_PRODUCTION), только чтение.
// Наборы — метаданные коллекций (строки — GetStages, GetMaterials, GetParametricMaterials,
// GetProducts).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#2-tehkarta
type ProcessingPlan struct {
	Meta                 model.Meta                            `json:"meta"`
	ID                   string                                `json:"id"`
	AccountID            string                                `json:"accountId"`
	Name                 string                                `json:"name"`
	Code                 string                                `json:"code"`
	ExternalCode         string                                `json:"externalCode"`
	Archived             bool                                  `json:"archived"`
	Shared               bool                                  `json:"shared"`
	Updated              model.DateTime                        `json:"updated,omitzero"`
	PathName             string                                `json:"pathName"`
	Cost                 float64                               `json:"cost"`
	CostDistributionType string                                `json:"costDistributionType"`
	Owner                *model.Reference                      `json:"owner,omitzero"`
	Group                *model.Reference                      `json:"group,omitzero"`
	Parent               *model.Reference                      `json:"parent,omitzero"`
	ProcessingProcess    *model.Reference                      `json:"processingProcess,omitzero"`
	Stages               *model.Collection[Stage]              `json:"stages,omitzero"`
	Materials            *model.Collection[Material]           `json:"materials,omitzero"`
	ParametricMaterials  *model.Collection[ParametricMaterial] `json:"parametricMaterials,omitzero"`
	Products             *model.Collection[Product]            `json:"products,omitzero"`
}

// Stage — этап техкарты (meta.type processingplanstages; новые примеры источника —
// processingplanstage). Cost, LabourCost и StandardHourCost — копейки; StandardHour — нормо-часы,
// не деньги.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#4-etapy-tehkarty
type Stage struct {
	Meta                      model.Meta       `json:"meta"`
	ID                        string           `json:"id"`
	AccountID                 string           `json:"accountId"`
	EnableHourAccounting      bool             `json:"enableHourAccounting"`
	Cost                      float64          `json:"cost"`
	LabourCost                float64          `json:"labourCost"`
	StandardHour              float64          `json:"standardHour"`
	StandardHourCost          float64          `json:"standardHourCost"`
	ProcessingProcessPosition *model.Reference `json:"processingProcessPosition,omitzero"`
}

// Material — материал техкарты (meta.type processingplanmaterial). Quantity — количество, не
// деньги. MaterialProcessingPlan — техкарта материала, только чтение.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#4-materialy-tehkarty
type Material struct {
	Meta                      model.Meta       `json:"meta"`
	ID                        string           `json:"id"`
	AccountID                 string           `json:"accountId"`
	Assortment                *model.Reference `json:"assortment,omitzero"`
	Product                   *model.Reference `json:"product,omitzero"`
	Quantity                  float64          `json:"quantity"`
	ProcessingProcessPosition *model.Reference `json:"processingProcessPosition,omitzero"`
	MaterialProcessingPlan    *model.Reference `json:"materialProcessingPlan,omitzero"`
}

// ParametricMaterial — параметрический материал техкарты: собственная форма, не Material.
// Материал зависит от ParamFeatureCharacteristic (значения — ParamFeatures), норма — от
// ParamQuantityCharacteristic (значения — ParamQuantities). Если от характеристики не зависит
// материал, сервер отдаёт Assortment и ParamFeatureCharacteristic = null; если норма — Quantity и
// ParamQuantityCharacteristic = null. meta.type в примерах источника — и
// processingplanparametricmaterial, и processingplanmaterial.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#4-parametricheskie-materialy-tehkarty
type ParametricMaterial struct {
	Meta                        model.Meta       `json:"meta"`
	ID                          string           `json:"id"`
	AccountID                   string           `json:"accountId"`
	Assortment                  *model.Reference `json:"assortment,omitzero"`
	Quantity                    float64          `json:"quantity"`
	ParamFeatureCharacteristic  *model.Reference `json:"paramFeatureCharacteristic,omitzero"`
	ParamQuantityCharacteristic *model.Reference `json:"paramQuantityCharacteristic,omitzero"`
	ParamFeatures               []ParamFeature   `json:"paramFeatures"`
	ParamQuantities             []ParamQuantity  `json:"paramQuantities"`
	ProcessingProcessPosition   *model.Reference `json:"processingProcessPosition,omitzero"`
}

// ParamFeature — материал для значения характеристики.
type ParamFeature struct {
	CharacteristicValue string           `json:"characteristicValue"`
	Assortment          *model.Reference `json:"assortment,omitzero"`
}

// ParamQuantity — норма расхода (количество, не деньги) для значения характеристики.
type ParamQuantity struct {
	CharacteristicValue string  `json:"characteristicValue"`
	Quantity            float64 `json:"quantity"`
}

// Product — продукт техкарты (meta.type processingplanresult). Quantity — количество, не деньги.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#4-produkty-tehkarty
type Product struct {
	Meta       model.Meta       `json:"meta"`
	ID         string           `json:"id"`
	AccountID  string           `json:"accountId"`
	Assortment *model.Reference `json:"assortment,omitzero"`
	Product    *model.Reference `json:"product,omitzero"`
	Quantity   float64          `json:"quantity"`
}
