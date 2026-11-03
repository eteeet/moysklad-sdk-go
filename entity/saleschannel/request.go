package saleschannel

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: обязательны name и type; значение type сервер проверяет сам.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/saleschannel#3-sozdat-kanal-prodazh
type CreateRequest struct {
	Name         string                       `json:"name"`
	Type         string                       `json:"type"`
	Description  model.Field[string]          `json:"description,omitzero"`
	Code         model.Field[string]          `json:"code,omitzero"`
	ExternalCode model.Field[string]          `json:"externalCode,omitzero"`
	Archived     model.Field[bool]            `json:"archived,omitzero"`
	Shared       model.Field[bool]            `json:"shared,omitzero"`
	Owner        model.Field[model.Reference] `json:"owner,omitzero"`
	Group        model.Field[model.Reference] `json:"group,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/saleschannel#3-izmenit-kanal-prodazh
type UpdateRequest struct {
	Name         model.Field[string]          `json:"name,omitzero"`
	Type         model.Field[string]          `json:"type,omitzero"`
	Description  model.Field[string]          `json:"description,omitzero"`
	Code         model.Field[string]          `json:"code,omitzero"`
	ExternalCode model.Field[string]          `json:"externalCode,omitzero"`
	Archived     model.Field[bool]            `json:"archived,omitzero"`
	Shared       model.Field[bool]            `json:"shared,omitzero"`
	Owner        model.Field[model.Reference] `json:"owner,omitzero"`
	Group        model.Field[model.Reference] `json:"group,omitzero"`
}
