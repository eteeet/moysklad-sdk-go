package customentity

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: единственное поле справочника — name; создание и изменение доступны только
// администратору (проверяет сервер).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/customentity#3-sozdat-spravochnik
type CreateRequest struct {
	Name string `json:"name"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/customentity#3-izmenit-spravochnik
type UpdateRequest struct {
	Name model.Field[string] `json:"name,omitzero"`
}

// ElementCreateRequest: обязателен только name.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/customentity#3-sozdat-element-spravochnika
type ElementCreateRequest struct {
	Name         string                         `json:"name"`
	Code         model.Field[string]            `json:"code,omitzero"`
	Description  model.Field[string]            `json:"description,omitzero"`
	ExternalCode model.Field[string]            `json:"externalCode,omitzero"`
	Attributes   model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	Shared       model.Field[bool]              `json:"shared,omitzero"`
	Owner        model.Field[model.Reference]   `json:"owner,omitzero"`
	Group        model.Field[model.Reference]   `json:"group,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/customentity#3-izmenit-element-spravochnika
type ElementUpdateRequest struct {
	Name         model.Field[string]            `json:"name,omitzero"`
	Code         model.Field[string]            `json:"code,omitzero"`
	Description  model.Field[string]            `json:"description,omitzero"`
	ExternalCode model.Field[string]            `json:"externalCode,omitzero"`
	Attributes   model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	Shared       model.Field[bool]              `json:"shared,omitzero"`
	Owner        model.Field[model.Reference]   `json:"owner,omitzero"`
	Group        model.Field[model.Reference]   `json:"group,omitzero"`
}
