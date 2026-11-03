package processingplanfolder

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: обязателен только name. accountId, id, pathName и updated таблица помечает
// «только для чтения» — их нет. Родительской группы нет: её нет в таблице атрибутов.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplanfolder#3-sozdat-novuyu-gruppu-tehkart
type CreateRequest struct {
	Name         string                       `json:"name"`
	Description  model.Field[string]          `json:"description,omitzero"`
	Code         model.Field[string]          `json:"code,omitzero"`
	ExternalCode model.Field[string]          `json:"externalCode,omitzero"`
	Archived     model.Field[bool]            `json:"archived,omitzero"`
	Shared       model.Field[bool]            `json:"shared,omitzero"`
	Owner        model.Field[model.Reference] `json:"owner,omitzero"`
	Group        model.Field[model.Reference] `json:"group,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplanfolder#3-izmenit-gruppu-tehkart
type UpdateRequest struct {
	Name         model.Field[string]          `json:"name,omitzero"`
	Description  model.Field[string]          `json:"description,omitzero"`
	Code         model.Field[string]          `json:"code,omitzero"`
	ExternalCode model.Field[string]          `json:"externalCode,omitzero"`
	Archived     model.Field[bool]            `json:"archived,omitzero"`
	Shared       model.Field[bool]            `json:"shared,omitzero"`
	Owner        model.Field[model.Reference] `json:"owner,omitzero"`
	Group        model.Field[model.Reference] `json:"group,omitzero"`
}
