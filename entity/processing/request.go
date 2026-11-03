package processing

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы organization, materialsStore и productsStore (проверяются до запроса).
// Дальше — один из двух вариантов, его проверяет сервер:
//   - с техкартой (ProcessingPlan): незаданные Products, Materials и ProcessingSum сервер берёт из
//     техкарты; Quantity по умолчанию 1, иное значение пересчитывает количества строк кратно ему;
//     чтобы не брать материалы техкарты, передайте Materials: model.Value([]PositionRequest{}).
//   - без техкарты: обязательны Products (хотя бы одна строка) и ProcessingSum; Quantity по
//     умолчанию 1.
//
// Materials и Products — не больше 1000 строк каждый. ProcessingSum — копейки. syncId в таблице
// только для чтения — SDK его не пишет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-sozdat-tehoperaciyu
type CreateRequest struct {
	Organization        model.Reference                `json:"organization"`
	MaterialsStore      model.Reference                `json:"materialsStore"`
	ProductsStore       model.Reference                `json:"productsStore"`
	Name                model.Field[string]            `json:"name,omitzero"`
	Description         model.Field[string]            `json:"description,omitzero"`
	Code                model.Field[string]            `json:"code,omitzero"`
	ExternalCode        model.Field[string]            `json:"externalCode,omitzero"`
	Moment              model.Field[model.DateTime]    `json:"moment,omitzero"`
	Applicable          model.Field[bool]              `json:"applicable,omitzero"`
	Shared              model.Field[bool]              `json:"shared,omitzero"`
	OrganizationAccount model.Field[model.Reference]   `json:"organizationAccount,omitzero"`
	Owner               model.Field[model.Reference]   `json:"owner,omitzero"`
	Group               model.Field[model.Reference]   `json:"group,omitzero"`
	Project             model.Field[model.Reference]   `json:"project,omitzero"`
	State               model.Field[model.Reference]   `json:"state,omitzero"`
	Attributes          model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	ProcessingPlan      model.Field[model.Reference]   `json:"processingPlan,omitzero"`
	ProcessingSum       model.Field[float64]           `json:"processingSum,omitzero"`
	Quantity            model.Field[float64]           `json:"quantity,omitzero"`
	Materials           model.Field[[]PositionRequest] `json:"materials,omitzero"`
	Products            model.Field[[]PositionRequest] `json:"products,omitzero"`
	ProcessingOrder     model.Field[model.Reference]   `json:"processingOrder,omitzero"`
}

// UpdateRequest. Новая ProcessingPlan заменяет продукты, материалы и затраты значениями техкарты,
// если Products, Materials и ProcessingSum не переданы; пустые материалы — Value(пустой слайс).
// Без Quantity объём не меняется; количества строк сервер пересчитывает кратно объёму.
// Materials и Products заменяют весь свой набор: строка с ID — существующая, без — новая.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing#3-izmenit-tehoperaciyu
type UpdateRequest struct {
	Organization        model.Field[model.Reference]   `json:"organization,omitzero"`
	MaterialsStore      model.Field[model.Reference]   `json:"materialsStore,omitzero"`
	ProductsStore       model.Field[model.Reference]   `json:"productsStore,omitzero"`
	Name                model.Field[string]            `json:"name,omitzero"`
	Description         model.Field[string]            `json:"description,omitzero"`
	Code                model.Field[string]            `json:"code,omitzero"`
	ExternalCode        model.Field[string]            `json:"externalCode,omitzero"`
	Moment              model.Field[model.DateTime]    `json:"moment,omitzero"`
	Applicable          model.Field[bool]              `json:"applicable,omitzero"`
	Shared              model.Field[bool]              `json:"shared,omitzero"`
	OrganizationAccount model.Field[model.Reference]   `json:"organizationAccount,omitzero"`
	Owner               model.Field[model.Reference]   `json:"owner,omitzero"`
	Group               model.Field[model.Reference]   `json:"group,omitzero"`
	Project             model.Field[model.Reference]   `json:"project,omitzero"`
	State               model.Field[model.Reference]   `json:"state,omitzero"`
	Attributes          model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	ProcessingPlan      model.Field[model.Reference]   `json:"processingPlan,omitzero"`
	ProcessingSum       model.Field[float64]           `json:"processingSum,omitzero"`
	Quantity            model.Field[float64]           `json:"quantity,omitzero"`
	Materials           model.Field[[]PositionRequest] `json:"materials,omitzero"`
	Products            model.Field[[]PositionRequest] `json:"products,omitzero"`
	ProcessingOrder     model.Field[model.Reference]   `json:"processingOrder,omitzero"`
}

// PositionRequest — строка материалов или продуктов в Create/Update документа, CreateMaterial,
// CreateProduct, UpdateMaterial и UpdateProduct. CreateMaterial/CreateProduct требуют Assortment и
// Quantity (проверяются до запроса); в Update строки обязательных полей нет.
type PositionRequest struct {
	// ID сопоставляет существующую строку в Update документа; для новой строки его не задают.
	ID         model.Field[string]          `json:"id,omitzero"`
	Assortment model.Field[model.Reference] `json:"assortment,omitzero"`
	Quantity   model.Field[float64]         `json:"quantity,omitzero"`
	Slot       model.Field[model.Reference] `json:"slot,omitzero"`
}
