package productiontask

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы organization, materialsStore и productsStore. Таблица атрибутов
// помечает «Необходимо при создании» только organization, но раздел «Создать Производственное
// задание» называет обязательными все три — SDK проверяет все три до запроса. Без productionStart
// задание создаётся не начатым. ProductionRows — не больше 200 строк, Products — не больше 1000.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-sozdat-proizvodstvennoe-zadanie
type CreateRequest struct {
	Organization          model.Reference                     `json:"organization"`
	MaterialsStore        model.Reference                     `json:"materialsStore"`
	ProductsStore         model.Reference                     `json:"productsStore"`
	Name                  model.Field[string]                 `json:"name,omitzero"`
	Description           model.Field[string]                 `json:"description,omitzero"`
	Code                  model.Field[string]                 `json:"code,omitzero"`
	ExternalCode          model.Field[string]                 `json:"externalCode,omitzero"`
	Moment                model.Field[model.DateTime]         `json:"moment,omitzero"`
	DeliveryPlannedMoment model.Field[model.DateTime]         `json:"deliveryPlannedMoment,omitzero"`
	ProductionStart       model.Field[model.DateTime]         `json:"productionStart,omitzero"`
	Applicable            model.Field[bool]                   `json:"applicable,omitzero"`
	Awaiting              model.Field[bool]                   `json:"awaiting,omitzero"`
	Reserve               model.Field[bool]                   `json:"reserve,omitzero"`
	Shared                model.Field[bool]                   `json:"shared,omitzero"`
	Owner                 model.Field[model.Reference]        `json:"owner,omitzero"`
	Group                 model.Field[model.Reference]        `json:"group,omitzero"`
	State                 model.Field[model.Reference]        `json:"state,omitzero"`
	Attributes            model.Field[[]model.Attribute]      `json:"attributes,omitzero"`
	ProductionRows        model.Field[[]ProductionRowRequest] `json:"productionRows,omitzero"`
	Products              model.Field[[]ProductRequest]       `json:"products,omitzero"`
	// Связи с другими документами: основная таблица не помечает их «Только для чтения» — SDK их пишет.
	CustomerOrders         model.Field[[]model.Reference] `json:"customerOrders,omitzero"`
	Demands                model.Field[[]model.Reference] `json:"demands,omitzero"`
	EmissionOrders         model.Field[[]model.Reference] `json:"emissionOrders,omitzero"`
	InternalOrders         model.Field[[]model.Reference] `json:"internalOrders,omitzero"`
	Moves                  model.Field[[]model.Reference] `json:"moves,omitzero"`
	ProductionTasks        model.Field[[]model.Reference] `json:"productionTasks,omitzero"`
	ProductionTaskSupplies model.Field[[]model.Reference] `json:"productionTaskSupplies,omitzero"`
	PurchaseOrders         model.Field[[]model.Reference] `json:"purchaseOrders,omitzero"`
	Supplies               model.Field[[]model.Reference] `json:"supplies,omitzero"`
}

// UpdateRequest. ProductionRows и Products заменяют весь свой набор: строка с ID — существующая,
// без — новая, отсутствующие удаляются. Передача productionStart начинает производство или меняет
// его дату (дату позже выполненного этапа сервер отклоняет).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask#3-izmenit-proizvodstvennoe-zadanie
type UpdateRequest struct {
	Organization          model.Field[model.Reference]        `json:"organization,omitzero"`
	MaterialsStore        model.Field[model.Reference]        `json:"materialsStore,omitzero"`
	ProductsStore         model.Field[model.Reference]        `json:"productsStore,omitzero"`
	Name                  model.Field[string]                 `json:"name,omitzero"`
	Description           model.Field[string]                 `json:"description,omitzero"`
	Code                  model.Field[string]                 `json:"code,omitzero"`
	ExternalCode          model.Field[string]                 `json:"externalCode,omitzero"`
	Moment                model.Field[model.DateTime]         `json:"moment,omitzero"`
	DeliveryPlannedMoment model.Field[model.DateTime]         `json:"deliveryPlannedMoment,omitzero"`
	ProductionStart       model.Field[model.DateTime]         `json:"productionStart,omitzero"`
	Applicable            model.Field[bool]                   `json:"applicable,omitzero"`
	Awaiting              model.Field[bool]                   `json:"awaiting,omitzero"`
	Reserve               model.Field[bool]                   `json:"reserve,omitzero"`
	Shared                model.Field[bool]                   `json:"shared,omitzero"`
	Owner                 model.Field[model.Reference]        `json:"owner,omitzero"`
	Group                 model.Field[model.Reference]        `json:"group,omitzero"`
	State                 model.Field[model.Reference]        `json:"state,omitzero"`
	Attributes            model.Field[[]model.Attribute]      `json:"attributes,omitzero"`
	ProductionRows        model.Field[[]ProductionRowRequest] `json:"productionRows,omitzero"`
	Products              model.Field[[]ProductRequest]       `json:"products,omitzero"`
	// Связи с другими документами: основная таблица не помечает их «Только для чтения» — SDK их пишет.
	CustomerOrders         model.Field[[]model.Reference] `json:"customerOrders,omitzero"`
	Demands                model.Field[[]model.Reference] `json:"demands,omitzero"`
	EmissionOrders         model.Field[[]model.Reference] `json:"emissionOrders,omitzero"`
	InternalOrders         model.Field[[]model.Reference] `json:"internalOrders,omitzero"`
	Moves                  model.Field[[]model.Reference] `json:"moves,omitzero"`
	ProductionTasks        model.Field[[]model.Reference] `json:"productionTasks,omitzero"`
	ProductionTaskSupplies model.Field[[]model.Reference] `json:"productionTaskSupplies,omitzero"`
	PurchaseOrders         model.Field[[]model.Reference] `json:"purchaseOrders,omitzero"`
	Supplies               model.Field[[]model.Reference] `json:"supplies,omitzero"`
}
