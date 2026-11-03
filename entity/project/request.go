package project

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: обязателен только name.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/project#3-sozdat-proekt
type CreateRequest struct {
	Name         string                         `json:"name"`
	Description  model.Field[string]            `json:"description,omitzero"`
	Code         model.Field[string]            `json:"code,omitzero"`
	ExternalCode model.Field[string]            `json:"externalCode,omitzero"`
	Archived     model.Field[bool]              `json:"archived,omitzero"`
	Attributes   model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	Shared       model.Field[bool]              `json:"shared,omitzero"`
	Owner        model.Field[model.Reference]   `json:"owner,omitzero"`
	Group        model.Field[model.Reference]   `json:"group,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/project#3-izmenit-proekt
type UpdateRequest struct {
	Name         model.Field[string]            `json:"name,omitzero"`
	Description  model.Field[string]            `json:"description,omitzero"`
	Code         model.Field[string]            `json:"code,omitzero"`
	ExternalCode model.Field[string]            `json:"externalCode,omitzero"`
	Archived     model.Field[bool]              `json:"archived,omitzero"`
	Attributes   model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	Shared       model.Field[bool]              `json:"shared,omitzero"`
	Owner        model.Field[model.Reference]   `json:"owner,omitzero"`
	Group        model.Field[model.Reference]   `json:"group,omitzero"`
}
