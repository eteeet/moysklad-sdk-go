package productionstagecompletion

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: раздел «Создать» требует productionStage (непустой meta.href проверяется до
// запроса) и productionVolume (передаётся всегда, значение проверяет сервер). productionStage и
// defect после создания не меняются — их нет в UpdateRequest. applicable в таблице только для
// чтения. Материалы и продукты в теле документа не пишутся: таблица описывает их как MetaArray
// без формы записи, а строки меняются ресурсами …/materials и …/products.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-sozdat-vypolnenie-etapa-proizvodstva
type CreateRequest struct {
	ProductionStage      model.Reference              `json:"productionStage"`
	ProductionVolume     float64                      `json:"productionVolume"`
	Name                 model.Field[string]          `json:"name,omitzero"`
	Description          model.Field[string]          `json:"description,omitzero"`
	ExternalCode         model.Field[string]          `json:"externalCode,omitzero"`
	Moment               model.Field[model.DateTime]  `json:"moment,omitzero"`
	Shared               model.Field[bool]            `json:"shared,omitzero"`
	Defect               model.Field[bool]            `json:"defect,omitzero"`
	EnableHourAccounting model.Field[bool]            `json:"enableHourAccounting,omitzero"`
	Owner                model.Field[model.Reference] `json:"owner,omitzero"`
	Group                model.Field[model.Reference] `json:"group,omitzero"`
	Performer            model.Field[model.Reference] `json:"performer,omitzero"`
	Service              model.Field[model.Reference] `json:"service,omitzero"`
	ProcessingUnitCost   model.Field[float64]         `json:"processingUnitCost,omitzero"`
	LabourUnitCost       model.Field[float64]         `json:"labourUnitCost,omitzero"`
	StandardHourCost     model.Field[float64]         `json:"standardHourCost,omitzero"`
	StandardHourUnit     model.Field[float64]         `json:"standardHourUnit,omitzero"`
}

// UpdateRequest — поля, кроме только для чтения и неизменяемых после создания (productionStage,
// defect). При enableHourAccounting=true labourUnitCost вычисляет сервер, и его изменение
// отклоняется (проверяет сервер).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-izmenit-vypolnenie-etapa-proizvodstva
type UpdateRequest struct {
	ProductionVolume     model.Field[float64]         `json:"productionVolume,omitzero"`
	Name                 model.Field[string]          `json:"name,omitzero"`
	Description          model.Field[string]          `json:"description,omitzero"`
	ExternalCode         model.Field[string]          `json:"externalCode,omitzero"`
	Moment               model.Field[model.DateTime]  `json:"moment,omitzero"`
	Shared               model.Field[bool]            `json:"shared,omitzero"`
	EnableHourAccounting model.Field[bool]            `json:"enableHourAccounting,omitzero"`
	Owner                model.Field[model.Reference] `json:"owner,omitzero"`
	Group                model.Field[model.Reference] `json:"group,omitzero"`
	Performer            model.Field[model.Reference] `json:"performer,omitzero"`
	Service              model.Field[model.Reference] `json:"service,omitzero"`
	ProcessingUnitCost   model.Field[float64]         `json:"processingUnitCost,omitzero"`
	LabourUnitCost       model.Field[float64]         `json:"labourUnitCost,omitzero"`
	StandardHourCost     model.Field[float64]         `json:"standardHourCost,omitzero"`
	StandardHourUnit     model.Field[float64]         `json:"standardHourUnit,omitzero"`
}

// MaterialCreateRequest — новый материал: assortment обязателен (проверяется до запроса);
// consumedQuantity сервер игнорирует, если переданы things (только для серийного учёта).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-dobavit-material-vypolneniya-etapa-proizvodstva
type MaterialCreateRequest struct {
	Assortment       model.Reference       `json:"assortment"`
	ConsumedQuantity model.Field[float64]  `json:"consumedQuantity,omitzero"`
	Things           model.Field[[]string] `json:"things,omitzero"`
}

// MaterialUpdateRequest. При things сервер выставляет consumedQuantity по их числу.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-izmenit-material-vypolneniya-etapa-proizvodstva
type MaterialUpdateRequest struct {
	Assortment       model.Field[model.Reference] `json:"assortment,omitzero"`
	ConsumedQuantity model.Field[float64]         `json:"consumedQuantity,omitzero"`
	Things           model.Field[[]string]        `json:"things,omitzero"`
}

// ProductUpdateRequest. Товар или модификацию сменить нельзя, партию — можно (проверяет сервер).
// При things сервер выставляет producedQuantity по их числу.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion#3-izmenit-produkt-vypolneniya-etapa-proizvodstva
type ProductUpdateRequest struct {
	Assortment       model.Field[model.Reference] `json:"assortment,omitzero"`
	ProducedQuantity model.Field[float64]         `json:"producedQuantity,omitzero"`
	Things           model.Field[[]string]        `json:"things,omitzero"`
}
