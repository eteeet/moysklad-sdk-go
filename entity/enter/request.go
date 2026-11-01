package enter

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы organization и store. Positions — не больше 1000 строк; дальше
// источник предлагает отдельный ресурс позиций (в SDK v1 не входит).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/enter#3-sozdat-oprihodovaniya
type CreateRequest struct {
	Organization model.Reference                            `json:"organization"`
	Store        model.Reference                            `json:"store"`
	Name         model.Field[string]                        `json:"name,omitzero"`
	Description  model.Field[string]                        `json:"description,omitzero"`
	Code         model.Field[string]                        `json:"code,omitzero"`
	ExternalCode model.Field[string]                        `json:"externalCode,omitzero"`
	SyncID       model.Field[string]                        `json:"syncId,omitzero"`
	Moment       model.Field[model.DateTime]                `json:"moment,omitzero"`
	Applicable   model.Field[bool]                          `json:"applicable,omitzero"`
	Shared       model.Field[bool]                          `json:"shared,omitzero"`
	Owner        model.Field[model.Reference]               `json:"owner,omitzero"`
	Group        model.Field[model.Reference]               `json:"group,omitzero"`
	Project      model.Field[model.Reference]               `json:"project,omitzero"`
	State        model.Field[model.Reference]               `json:"state,omitzero"`
	Attributes   model.Field[[]model.Attribute]             `json:"attributes,omitzero"`
	Rate         model.Field[model.DocumentRateRequest]     `json:"rate,omitzero"`
	Positions    model.Field[[]PositionRequest]             `json:"positions,omitzero"`
	Overhead     model.Field[model.DocumentOverheadRequest] `json:"overhead,omitzero"`
}

// UpdateRequest. Positions заменяют все строки: лишние удаляются, строки с ID меняются.
// syncId после заполнения не меняется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/enter#3-izmenit-oprihodovanie
type UpdateRequest struct {
	Organization model.Field[model.Reference]               `json:"organization,omitzero"`
	Store        model.Field[model.Reference]               `json:"store,omitzero"`
	Name         model.Field[string]                        `json:"name,omitzero"`
	Description  model.Field[string]                        `json:"description,omitzero"`
	Code         model.Field[string]                        `json:"code,omitzero"`
	ExternalCode model.Field[string]                        `json:"externalCode,omitzero"`
	SyncID       model.Field[string]                        `json:"syncId,omitzero"`
	Moment       model.Field[model.DateTime]                `json:"moment,omitzero"`
	Applicable   model.Field[bool]                          `json:"applicable,omitzero"`
	Shared       model.Field[bool]                          `json:"shared,omitzero"`
	Owner        model.Field[model.Reference]               `json:"owner,omitzero"`
	Group        model.Field[model.Reference]               `json:"group,omitzero"`
	Project      model.Field[model.Reference]               `json:"project,omitzero"`
	State        model.Field[model.Reference]               `json:"state,omitzero"`
	Attributes   model.Field[[]model.Attribute]             `json:"attributes,omitzero"`
	Rate         model.Field[model.DocumentRateRequest]     `json:"rate,omitzero"`
	Positions    model.Field[[]PositionRequest]             `json:"positions,omitzero"`
	Overhead     model.Field[model.DocumentOverheadRequest] `json:"overhead,omitzero"`
}
