package retaildrawercashin

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы organization, agent (сотрудник) и retailShift. sum таблица помечает
// «Только для чтения», хотя примеры Create и массового создания его передают: SDK sum не пишет.
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
}

// UpdateRequest. Пример Update источника передаёт sum, таблица помечает его readonly — SDK sum не
// пишет. syncId после заполнения не меняется.
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
}
