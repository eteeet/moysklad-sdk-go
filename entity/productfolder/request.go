package productfolder

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: обязателен только name. archived, pathName, effectiveVat и effectiveVatEnabled
// таблица помечает «только для чтения» — их нет. ProductFolder — родительская группа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/productfolder#3-sozdat-novuyu-gruppu-tovarov
type CreateRequest struct {
	Name          string                       `json:"name"`
	Description   model.Field[string]          `json:"description,omitzero"`
	Code          model.Field[string]          `json:"code,omitzero"`
	ExternalCode  model.Field[string]          `json:"externalCode,omitzero"`
	ProductFolder model.Field[model.Reference] `json:"productFolder,omitzero"`
	Vat           model.Field[int]             `json:"vat,omitzero"`
	VatEnabled    model.Field[bool]            `json:"vatEnabled,omitzero"`
	UseParentVat  model.Field[bool]            `json:"useParentVat,omitzero"`
	TaxSystem     model.Field[string]          `json:"taxSystem,omitzero"`
	Shared        model.Field[bool]            `json:"shared,omitzero"`
	Owner         model.Field[model.Reference] `json:"owner,omitzero"`
	Group         model.Field[model.Reference] `json:"group,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/productfolder#3-izmenit-gruppu-tovarov
type UpdateRequest struct {
	Name          model.Field[string]          `json:"name,omitzero"`
	Description   model.Field[string]          `json:"description,omitzero"`
	Code          model.Field[string]          `json:"code,omitzero"`
	ExternalCode  model.Field[string]          `json:"externalCode,omitzero"`
	ProductFolder model.Field[model.Reference] `json:"productFolder,omitzero"`
	Vat           model.Field[int]             `json:"vat,omitzero"`
	VatEnabled    model.Field[bool]            `json:"vatEnabled,omitzero"`
	UseParentVat  model.Field[bool]            `json:"useParentVat,omitzero"`
	TaxSystem     model.Field[string]          `json:"taxSystem,omitzero"`
	Shared        model.Field[bool]            `json:"shared,omitzero"`
	Owner         model.Field[model.Reference] `json:"owner,omitzero"`
	Group         model.Field[model.Reference] `json:"group,omitzero"`
}
