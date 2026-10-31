package expenseitem

import "github.com/eteeet/moysklad-sdk-go/model"

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/expenseitem#3-sozdat-statyu-rashodov
type CreateRequest struct {
	Name              string                       `json:"name"`
	Code              model.Field[string]          `json:"code,omitzero"`
	ExternalCode      model.Field[string]          `json:"externalCode,omitzero"`
	Description       model.Field[string]          `json:"description,omitzero"`
	Shared            model.Field[bool]            `json:"shared,omitzero"`
	Archived          model.Field[bool]            `json:"archived,omitzero"`
	OperatingExpenses model.Field[bool]            `json:"operatingExpenses,omitzero"`
	Owner             model.Field[model.Reference] `json:"owner,omitzero"`
	Group             model.Field[model.Reference] `json:"group,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/expenseitem#3-izmenit-statyu-rashodov
type UpdateRequest struct {
	Name              model.Field[string]          `json:"name,omitzero"`
	Code              model.Field[string]          `json:"code,omitzero"`
	ExternalCode      model.Field[string]          `json:"externalCode,omitzero"`
	Description       model.Field[string]          `json:"description,omitzero"`
	Shared            model.Field[bool]            `json:"shared,omitzero"`
	Archived          model.Field[bool]            `json:"archived,omitzero"`
	OperatingExpenses model.Field[bool]            `json:"operatingExpenses,omitzero"`
	Owner             model.Field[model.Reference] `json:"owner,omitzero"`
	Group             model.Field[model.Reference] `json:"group,omitzero"`
}
