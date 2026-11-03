package internalorder

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы organization и непустое name — так помечает таблица атрибутов
// («Необходимо при создании»), хотя раздел «Создать» перечисляет только organization.
// shared, syncId и positions таблица помечает «Только для чтения», поэтому их нет в запросах,
// хотя примеры Create/Update их передают: строки пишутся отдельными методами позиций.
// moves, purchaseOrders и productionTasks — связи документов без readonly-маркера: пишутся
// (ticket O08 называет moves/purchaseOrders read-only — решение по источнику).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-sozdat-vnutrennij-zakaz
type CreateRequest struct {
	Organization          model.Reference                        `json:"organization"`
	Name                  string                                 `json:"name"`
	Description           model.Field[string]                    `json:"description,omitzero"`
	Code                  model.Field[string]                    `json:"code,omitzero"`
	ExternalCode          model.Field[string]                    `json:"externalCode,omitzero"`
	Moment                model.Field[model.DateTime]            `json:"moment,omitzero"`
	DeliveryPlannedMoment model.Field[model.DateTime]            `json:"deliveryPlannedMoment,omitzero"`
	Applicable            model.Field[bool]                      `json:"applicable,omitzero"`
	Store                 model.Field[model.Reference]           `json:"store,omitzero"`
	Owner                 model.Field[model.Reference]           `json:"owner,omitzero"`
	Group                 model.Field[model.Reference]           `json:"group,omitzero"`
	Project               model.Field[model.Reference]           `json:"project,omitzero"`
	State                 model.Field[model.Reference]           `json:"state,omitzero"`
	Attributes            model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate                  model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	VATEnabled            model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded           model.Field[bool]                      `json:"vatIncluded,omitzero"`
	Moves                 model.Field[[]model.Reference]         `json:"moves,omitzero"`
	PurchaseOrders        model.Field[[]model.Reference]         `json:"purchaseOrders,omitzero"`
	ProductionTasks       model.Field[[]model.Reference]         `json:"productionTasks,omitzero"`
}

// UpdateRequest: обязательных полей нет; sum, vatSum, shared, syncId и positions — readonly.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder#3-izmenit-vnutrennij-zakaz
type UpdateRequest struct {
	Organization          model.Field[model.Reference]           `json:"organization,omitzero"`
	Name                  model.Field[string]                    `json:"name,omitzero"`
	Description           model.Field[string]                    `json:"description,omitzero"`
	Code                  model.Field[string]                    `json:"code,omitzero"`
	ExternalCode          model.Field[string]                    `json:"externalCode,omitzero"`
	Moment                model.Field[model.DateTime]            `json:"moment,omitzero"`
	DeliveryPlannedMoment model.Field[model.DateTime]            `json:"deliveryPlannedMoment,omitzero"`
	Applicable            model.Field[bool]                      `json:"applicable,omitzero"`
	Store                 model.Field[model.Reference]           `json:"store,omitzero"`
	Owner                 model.Field[model.Reference]           `json:"owner,omitzero"`
	Group                 model.Field[model.Reference]           `json:"group,omitzero"`
	Project               model.Field[model.Reference]           `json:"project,omitzero"`
	State                 model.Field[model.Reference]           `json:"state,omitzero"`
	Attributes            model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate                  model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	VATEnabled            model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded           model.Field[bool]                      `json:"vatIncluded,omitzero"`
	Moves                 model.Field[[]model.Reference]         `json:"moves,omitzero"`
	PurchaseOrders        model.Field[[]model.Reference]         `json:"purchaseOrders,omitzero"`
	ProductionTasks       model.Field[[]model.Reference]         `json:"productionTasks,omitzero"`
}
