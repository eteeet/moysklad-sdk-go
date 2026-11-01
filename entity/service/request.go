package service

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимо и достаточно непустое name. Каждый штрихкод — объект с одним
// форматом. По источнику доп. поля и цены продажи сервер сливает по метаданным и типу цены;
// что делает пустой массив штрихкодов, источник не описывает. Files задаются отдельным ресурсом.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/service#3-sozdat-uslugu
type CreateRequest struct {
	Name               string              `json:"name"`
	Code               model.Field[string] `json:"code,omitzero"`
	ExternalCode       model.Field[string] `json:"externalCode,omitzero"`
	Description        model.Field[string] `json:"description,omitzero"`
	Archived           model.Field[bool]   `json:"archived,omitzero"`
	Shared             model.Field[bool]   `json:"shared,omitzero"`
	VAT                model.Field[int]    `json:"vat,omitzero"`
	VATEnabled         model.Field[bool]   `json:"vatEnabled,omitzero"`
	UseParentVAT       model.Field[bool]   `json:"useParentVat,omitzero"`
	DiscountProhibited model.Field[bool]   `json:"discountProhibited,omitzero"`
	PaymentItemType    model.Field[string] `json:"paymentItemType,omitzero"`
	TaxSystem          model.Field[string] `json:"taxSystem,omitzero"`

	Owner         model.Field[model.Reference] `json:"owner,omitzero"`
	Group         model.Field[model.Reference] `json:"group,omitzero"`
	ProductFolder model.Field[model.Reference] `json:"productFolder,omitzero"`
	UOM           model.Field[model.Reference] `json:"uom,omitzero"`

	BuyPrice   model.Field[model.Price]       `json:"buyPrice,omitzero"`
	MinPrice   model.Field[model.Price]       `json:"minPrice,omitzero"`
	SalePrices model.Field[[]model.SalePrice] `json:"salePrices,omitzero"`
	Barcodes   model.Field[[]model.Barcode]   `json:"barcodes,omitzero"`
	Attributes model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	TasnifUZ   model.Field[model.TasnifUZ]    `json:"mod__tasnif__uz,omitzero"`

	// SyncID — UUID клиента; сервер записывает его один раз при создании, в UpdateRequest его нет.
	// По общему разделу о syncId повторный Create с тем же значением вернёт созданную услугу,
	// но SDK сам Create не повторяет.
	SyncID model.Field[string] `json:"syncId,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/service#3-izmenit-uslugu
type UpdateRequest struct {
	Name               model.Field[string] `json:"name,omitzero"`
	Code               model.Field[string] `json:"code,omitzero"`
	ExternalCode       model.Field[string] `json:"externalCode,omitzero"`
	Description        model.Field[string] `json:"description,omitzero"`
	Archived           model.Field[bool]   `json:"archived,omitzero"`
	Shared             model.Field[bool]   `json:"shared,omitzero"`
	VAT                model.Field[int]    `json:"vat,omitzero"`
	VATEnabled         model.Field[bool]   `json:"vatEnabled,omitzero"`
	UseParentVAT       model.Field[bool]   `json:"useParentVat,omitzero"`
	DiscountProhibited model.Field[bool]   `json:"discountProhibited,omitzero"`
	PaymentItemType    model.Field[string] `json:"paymentItemType,omitzero"`
	TaxSystem          model.Field[string] `json:"taxSystem,omitzero"`

	Owner         model.Field[model.Reference] `json:"owner,omitzero"`
	Group         model.Field[model.Reference] `json:"group,omitzero"`
	ProductFolder model.Field[model.Reference] `json:"productFolder,omitzero"`
	UOM           model.Field[model.Reference] `json:"uom,omitzero"`

	BuyPrice   model.Field[model.Price]       `json:"buyPrice,omitzero"`
	MinPrice   model.Field[model.Price]       `json:"minPrice,omitzero"`
	SalePrices model.Field[[]model.SalePrice] `json:"salePrices,omitzero"`
	Barcodes   model.Field[[]model.Barcode]   `json:"barcodes,omitzero"`
	Attributes model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	TasnifUZ   model.Field[model.TasnifUZ]    `json:"mod__tasnif__uz,omitzero"`
}
