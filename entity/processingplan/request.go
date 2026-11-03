package processingplan

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: обязательны непустые name и products (проверяются до запроса). Без
// processingProcess сервер подставит «Основной техпроцесс». Cost — копейки; одновременно с cost
// этапов сервер его не принимает. Наборы в теле — не больше 1000 строк каждый; больше —
// отдельными методами наборов. В параметрической техкарте продукт один (проверяет сервер).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-sozdat-tehkartu
type CreateRequest struct {
	Name                string                                   `json:"name"`
	Products            []ProductRequest                         `json:"products"`
	Code                model.Field[string]                      `json:"code,omitzero"`
	ExternalCode        model.Field[string]                      `json:"externalCode,omitzero"`
	Archived            model.Field[bool]                        `json:"archived,omitzero"`
	Shared              model.Field[bool]                        `json:"shared,omitzero"`
	Cost                model.Field[float64]                     `json:"cost,omitzero"`
	Owner               model.Field[model.Reference]             `json:"owner,omitzero"`
	Group               model.Field[model.Reference]             `json:"group,omitzero"`
	Parent              model.Field[model.Reference]             `json:"parent,omitzero"`
	ProcessingProcess   model.Field[model.Reference]             `json:"processingProcess,omitzero"`
	Stages              model.Field[[]StageRequest]              `json:"stages,omitzero"`
	Materials           model.Field[[]MaterialRequest]           `json:"materials,omitzero"`
	ParametricMaterials model.Field[[]ParametricMaterialRequest] `json:"parametricMaterials,omitzero"`
}

// UpdateRequest. Переданный набор заменяет весь набор: строка с ID — существующая, без ID — новая,
// не переданные удаляются. Смена processingProcess удаляет материалы, привязанные к старому
// техпроцессу. Этапы создаются и удаляются только вместе с этапами техпроцесса.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan#3-izmenit-tehkartu
type UpdateRequest struct {
	Name                model.Field[string]                      `json:"name,omitzero"`
	Products            model.Field[[]ProductRequest]            `json:"products,omitzero"`
	Code                model.Field[string]                      `json:"code,omitzero"`
	ExternalCode        model.Field[string]                      `json:"externalCode,omitzero"`
	Archived            model.Field[bool]                        `json:"archived,omitzero"`
	Shared              model.Field[bool]                        `json:"shared,omitzero"`
	Cost                model.Field[float64]                     `json:"cost,omitzero"`
	Owner               model.Field[model.Reference]             `json:"owner,omitzero"`
	Group               model.Field[model.Reference]             `json:"group,omitzero"`
	Parent              model.Field[model.Reference]             `json:"parent,omitzero"`
	ProcessingProcess   model.Field[model.Reference]             `json:"processingProcess,omitzero"`
	Stages              model.Field[[]StageRequest]              `json:"stages,omitzero"`
	Materials           model.Field[[]MaterialRequest]           `json:"materials,omitzero"`
	ParametricMaterials model.Field[[]ParametricMaterialRequest] `json:"parametricMaterials,omitzero"`
}

// StageRequest — этап в теле техкарты и в UpdateStage; обязательных полей нет. Cost и LabourCost —
// копейки; StandardHour — нормо-часы. При EnableHourAccounting=true сервер не принимает LabourCost.
type StageRequest struct {
	// ID сопоставляет существующий этап в Update техкарты; в UpdateStage его не задают.
	ID                        model.Field[string]          `json:"id,omitzero"`
	EnableHourAccounting      model.Field[bool]            `json:"enableHourAccounting,omitzero"`
	Cost                      model.Field[float64]         `json:"cost,omitzero"`
	LabourCost                model.Field[float64]         `json:"labourCost,omitzero"`
	StandardHour              model.Field[float64]         `json:"standardHour,omitzero"`
	ProcessingProcessPosition model.Field[model.Reference] `json:"processingProcessPosition,omitzero"`
}

// MaterialRequest — материал в теле техкарты, CreateMaterials и UpdateMaterial. Для нового
// материала сервер требует товар (assortment или product) и quantity; без processingProcessPosition
// материал привязывается к первой позиции техпроцесса.
type MaterialRequest struct {
	// ID сопоставляет существующий материал в Update техкарты; для нового его не задают.
	ID                        model.Field[string]          `json:"id,omitzero"`
	Assortment                model.Field[model.Reference] `json:"assortment,omitzero"`
	Product                   model.Field[model.Reference] `json:"product,omitzero"`
	Quantity                  model.Field[float64]         `json:"quantity,omitzero"`
	ProcessingProcessPosition model.Field[model.Reference] `json:"processingProcessPosition,omitzero"`
}

// ParametricMaterialRequest — параметрический материал в теле техкарты,
// CreateParametricMaterials и UpdateParametricMaterial. Без характеристики допустим только один из
// вариантов: конкретный Assortment или конкретный Quantity (проверяет сервер). Материал, не
// зависящий от характеристик, — MaterialRequest.
type ParametricMaterialRequest struct {
	// ID сопоставляет существующий материал в Update техкарты; для нового его не задают.
	ID                          model.Field[string]                 `json:"id,omitzero"`
	Assortment                  model.Field[model.Reference]        `json:"assortment,omitzero"`
	Quantity                    model.Field[float64]                `json:"quantity,omitzero"`
	ParamFeatureCharacteristic  model.Field[model.Reference]        `json:"paramFeatureCharacteristic,omitzero"`
	ParamQuantityCharacteristic model.Field[model.Reference]        `json:"paramQuantityCharacteristic,omitzero"`
	ParamFeatures               model.Field[[]ParamFeatureRequest]  `json:"paramFeatures,omitzero"`
	ParamQuantities             model.Field[[]ParamQuantityRequest] `json:"paramQuantities,omitzero"`
	ProcessingProcessPosition   model.Field[model.Reference]        `json:"processingProcessPosition,omitzero"`
}

// ParamFeatureRequest — элемент paramFeatures. characteristicValue таблица помечает «Только для
// чтения», хотя примеры Create его передают, — SDK его не пишет (см. docs/processingplan.md).
type ParamFeatureRequest struct {
	Assortment model.Reference `json:"assortment"`
}

// ParamQuantityRequest — элемент paramQuantities; characteristicValue не пишется, как в
// ParamFeatureRequest.
type ParamQuantityRequest struct {
	Quantity float64 `json:"quantity"`
}

// ProductRequest — продукт в теле техкарты, CreateProducts и UpdateProduct. Для нового продукта
// сервер требует товар (assortment или product) и quantity. Смена единственного продукта
// параметрической техкарты сбрасывает её параметрические материалы.
type ProductRequest struct {
	// ID сопоставляет существующий продукт в Update техкарты; для нового его не задают.
	ID         model.Field[string]          `json:"id,omitzero"`
	Assortment model.Field[model.Reference] `json:"assortment,omitzero"`
	Product    model.Field[model.Reference] `json:"product,omitzero"`
	Quantity   model.Field[float64]         `json:"quantity,omitzero"`
}
