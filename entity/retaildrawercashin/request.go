package retaildrawercashin

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы organization, agent (сотрудник) и retailShift. Sum (копейки) таблица
// помечает «Только для чтения», но примеры Create, массового создания и Update его передают: SDK
// sum пишет — исключение из правила «пометка таблицы сильнее примеров» утверждено пользователем.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildrawercashin#3-sozdat-vnesenie-deneg
type CreateRequest struct {
	Organization       model.Reference                        `json:"organization"`
	Agent              model.Reference                        `json:"agent"`
	RetailShift        model.Reference                        `json:"retailShift"`
	OrganizationBranch model.Field[model.Reference]           `json:"organizationBranch,omitzero"`
	Name               model.Field[string]                    `json:"name,omitzero"`
	Description        model.Field[string]                    `json:"description,omitzero"`
	ExternalCode       model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID             model.Field[string]                    `json:"syncId,omitzero"`
	Moment             model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable         model.Field[bool]                      `json:"applicable,omitzero"`
	Shared             model.Field[bool]                      `json:"shared,omitzero"`
	Owner              model.Field[model.Reference]           `json:"owner,omitzero"`
	Group              model.Field[model.Reference]           `json:"group,omitzero"`
	State              model.Field[model.Reference]           `json:"state,omitzero"`
	Attributes         model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate               model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	Sum                model.Field[float64]                   `json:"sum,omitzero"`
}

// UpdateRequest. Sum пишется по исключению из CreateRequest. syncId после заполнения не меняется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildrawercashin#3-izmenit-vnesenie-deneg
type UpdateRequest struct {
	Organization       model.Field[model.Reference]           `json:"organization,omitzero"`
	Agent              model.Field[model.Reference]           `json:"agent,omitzero"`
	RetailShift        model.Field[model.Reference]           `json:"retailShift,omitzero"`
	OrganizationBranch model.Field[model.Reference]           `json:"organizationBranch,omitzero"`
	Name               model.Field[string]                    `json:"name,omitzero"`
	Description        model.Field[string]                    `json:"description,omitzero"`
	ExternalCode       model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID             model.Field[string]                    `json:"syncId,omitzero"`
	Moment             model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable         model.Field[bool]                      `json:"applicable,omitzero"`
	Shared             model.Field[bool]                      `json:"shared,omitzero"`
	Owner              model.Field[model.Reference]           `json:"owner,omitzero"`
	Group              model.Field[model.Reference]           `json:"group,omitzero"`
	State              model.Field[model.Reference]           `json:"state,omitzero"`
	Attributes         model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate               model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	Sum                model.Field[float64]                   `json:"sum,omitzero"`
}
