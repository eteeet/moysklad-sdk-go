package salesreturn

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы organization, store и agent. Возврат по основанию задаёт Demand;
// без основания Demand не передают. При Create по основанию сервер требует: контрагент, юрлицо
// и валюта совпадают с основанием; строки соответствуют его строкам, quantity не больше, чем
// в основании, строк не из основания нет, price не меняется. SDK это не проверяет и не делает
// Get основания. cost строки задаётся у возврата без основания. Positions — не больше 1000 строк.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/sales-return#3-sozdat-vozvrat-pokupatelya
type CreateRequest struct {
	Organization        model.Reference                        `json:"organization"`
	Store               model.Reference                        `json:"store"`
	Agent               model.Reference                        `json:"agent"`
	AgentAccount        model.Field[model.Reference]           `json:"agentAccount,omitzero"`
	Demand              model.Field[model.Reference]           `json:"demand,omitzero"`
	Name                model.Field[string]                    `json:"name,omitzero"`
	Description         model.Field[string]                    `json:"description,omitzero"`
	Code                model.Field[string]                    `json:"code,omitzero"`
	ExternalCode        model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID              model.Field[string]                    `json:"syncId,omitzero"`
	Moment              model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable          model.Field[bool]                      `json:"applicable,omitzero"`
	Shared              model.Field[bool]                      `json:"shared,omitzero"`
	OrganizationAccount model.Field[model.Reference]           `json:"organizationAccount,omitzero"`
	Contract            model.Field[model.Reference]           `json:"contract,omitzero"`
	Owner               model.Field[model.Reference]           `json:"owner,omitzero"`
	Group               model.Field[model.Reference]           `json:"group,omitzero"`
	Project             model.Field[model.Reference]           `json:"project,omitzero"`
	State               model.Field[model.Reference]           `json:"state,omitzero"`
	SalesChannel        model.Field[model.Reference]           `json:"salesChannel,omitzero"`
	Attributes          model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate                model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	VATEnabled          model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded         model.Field[bool]                      `json:"vatIncluded,omitzero"`
	Positions           model.Field[[]PositionRequest]         `json:"positions,omitzero"`
}

// UpdateRequest. Источник запрещает менять agent, agentAccount и demand — их здесь нет.
// При смене organization нужен и organizationAccount (проверяет сервер); валюту, отличную от
// основания, и строки, которых нет в основании, сервер отклоняет. Positions заменяют все строки.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/sales-return#3-izmenit-vozvrat-pokupatelya
type UpdateRequest struct {
	Organization        model.Field[model.Reference]           `json:"organization,omitzero"`
	Store               model.Field[model.Reference]           `json:"store,omitzero"`
	Name                model.Field[string]                    `json:"name,omitzero"`
	Description         model.Field[string]                    `json:"description,omitzero"`
	Code                model.Field[string]                    `json:"code,omitzero"`
	ExternalCode        model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID              model.Field[string]                    `json:"syncId,omitzero"`
	Moment              model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable          model.Field[bool]                      `json:"applicable,omitzero"`
	Shared              model.Field[bool]                      `json:"shared,omitzero"`
	OrganizationAccount model.Field[model.Reference]           `json:"organizationAccount,omitzero"`
	Contract            model.Field[model.Reference]           `json:"contract,omitzero"`
	Owner               model.Field[model.Reference]           `json:"owner,omitzero"`
	Group               model.Field[model.Reference]           `json:"group,omitzero"`
	Project             model.Field[model.Reference]           `json:"project,omitzero"`
	State               model.Field[model.Reference]           `json:"state,omitzero"`
	SalesChannel        model.Field[model.Reference]           `json:"salesChannel,omitzero"`
	Attributes          model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate                model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	VATEnabled          model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded         model.Field[bool]                      `json:"vatIncluded,omitzero"`
	Positions           model.Field[[]PositionRequest]         `json:"positions,omitzero"`
}
