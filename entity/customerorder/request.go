package customerorder

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы organization и agent. Positions — не больше 1000 строк; дальше
// источник предлагает отдельный ресурс позиций (в SDK v1 не входит). Адрес: при передаче обоих
// shipmentAddress игнорируется; null адресу не поддерживается, удалить — shipmentAddress "".
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/customerOrder#3-sozdat-zakaz-pokupatelya
type CreateRequest struct {
	Organization          model.Reference                        `json:"organization"`
	Agent                 model.Reference                        `json:"agent"`
	Name                  model.Field[string]                    `json:"name,omitzero"`
	Description           model.Field[string]                    `json:"description,omitzero"`
	Code                  model.Field[string]                    `json:"code,omitzero"`
	ExternalCode          model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID                model.Field[string]                    `json:"syncId,omitzero"`
	Moment                model.Field[model.DateTime]            `json:"moment,omitzero"`
	DeliveryPlannedMoment model.Field[model.DateTime]            `json:"deliveryPlannedMoment,omitzero"`
	Applicable            model.Field[bool]                      `json:"applicable,omitzero"`
	Shared                model.Field[bool]                      `json:"shared,omitzero"`
	OrganizationAccount   model.Field[model.Reference]           `json:"organizationAccount,omitzero"`
	AgentAccount          model.Field[model.Reference]           `json:"agentAccount,omitzero"`
	Contract              model.Field[model.Reference]           `json:"contract,omitzero"`
	Owner                 model.Field[model.Reference]           `json:"owner,omitzero"`
	Group                 model.Field[model.Reference]           `json:"group,omitzero"`
	Project               model.Field[model.Reference]           `json:"project,omitzero"`
	State                 model.Field[model.Reference]           `json:"state,omitzero"`
	Store                 model.Field[model.Reference]           `json:"store,omitzero"`
	SalesChannel          model.Field[model.Reference]           `json:"salesChannel,omitzero"`
	Attributes            model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate                  model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	TaxSystem             model.Field[string]                    `json:"taxSystem,omitzero"`
	VATEnabled            model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded           model.Field[bool]                      `json:"vatIncluded,omitzero"`
	ShipmentAddress       model.Field[string]                    `json:"shipmentAddress,omitzero"`
	ShipmentAddressFull   model.Field[model.AddressRequest]      `json:"shipmentAddressFull,omitzero"`
	Positions             model.Field[[]PositionRequest]         `json:"positions,omitzero"`
}

// UpdateRequest. Positions заменяют все строки заказа: лишние удаляются, строки с ID меняются.
// При смене organization или agent источник требует передать и organizationAccount или
// agentAccount, иначе ошибка (проверяет сервер). syncId после заполнения не меняется. Адресу
// null не поддерживается: Null не удаляет его, удалить — ShipmentAddress Value("").
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/customerOrder#3-izmenit-zakaz-pokupatelya
type UpdateRequest struct {
	Organization          model.Field[model.Reference]           `json:"organization,omitzero"`
	Agent                 model.Field[model.Reference]           `json:"agent,omitzero"`
	Name                  model.Field[string]                    `json:"name,omitzero"`
	Description           model.Field[string]                    `json:"description,omitzero"`
	Code                  model.Field[string]                    `json:"code,omitzero"`
	ExternalCode          model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID                model.Field[string]                    `json:"syncId,omitzero"`
	Moment                model.Field[model.DateTime]            `json:"moment,omitzero"`
	DeliveryPlannedMoment model.Field[model.DateTime]            `json:"deliveryPlannedMoment,omitzero"`
	Applicable            model.Field[bool]                      `json:"applicable,omitzero"`
	Shared                model.Field[bool]                      `json:"shared,omitzero"`
	OrganizationAccount   model.Field[model.Reference]           `json:"organizationAccount,omitzero"`
	AgentAccount          model.Field[model.Reference]           `json:"agentAccount,omitzero"`
	Contract              model.Field[model.Reference]           `json:"contract,omitzero"`
	Owner                 model.Field[model.Reference]           `json:"owner,omitzero"`
	Group                 model.Field[model.Reference]           `json:"group,omitzero"`
	Project               model.Field[model.Reference]           `json:"project,omitzero"`
	State                 model.Field[model.Reference]           `json:"state,omitzero"`
	Store                 model.Field[model.Reference]           `json:"store,omitzero"`
	SalesChannel          model.Field[model.Reference]           `json:"salesChannel,omitzero"`
	Attributes            model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate                  model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	TaxSystem             model.Field[string]                    `json:"taxSystem,omitzero"`
	VATEnabled            model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded           model.Field[bool]                      `json:"vatIncluded,omitzero"`
	ShipmentAddress       model.Field[string]                    `json:"shipmentAddress,omitzero"`
	ShipmentAddressFull   model.Field[model.AddressRequest]      `json:"shipmentAddressFull,omitzero"`
	Positions             model.Field[[]PositionRequest]         `json:"positions,omitzero"`
}
