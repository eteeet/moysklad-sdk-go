package consignment

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: label и assortment обязательны по тексту метода. Таблица помечает «необходимо при
// создании» и expiryDate, но метод и пример его не требуют — SDK его не проверяет.
// Barcodes — model.Barcode, как у товара.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/consignment#3-sozdat-partiyu
type CreateRequest struct {
	Label        string                         `json:"label"`
	Assortment   model.Reference                `json:"assortment"`
	ExpiryDate   model.Field[model.DateTime]    `json:"expiryDate,omitzero"`
	Description  model.Field[string]            `json:"description,omitzero"`
	Code         model.Field[string]            `json:"code,omitzero"`
	ExternalCode model.Field[string]            `json:"externalCode,omitzero"`
	Archived     model.Field[bool]              `json:"archived,omitzero"`
	Barcodes     model.Field[[]model.Barcode]   `json:"barcodes,omitzero"`
	Attributes   model.Field[[]model.Attribute] `json:"attributes,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/consignment#3-izmenit-partiyu
type UpdateRequest struct {
	Label        model.Field[string]            `json:"label,omitzero"`
	Assortment   model.Field[model.Reference]   `json:"assortment,omitzero"`
	ExpiryDate   model.Field[model.DateTime]    `json:"expiryDate,omitzero"`
	Description  model.Field[string]            `json:"description,omitzero"`
	Code         model.Field[string]            `json:"code,omitzero"`
	ExternalCode model.Field[string]            `json:"externalCode,omitzero"`
	Archived     model.Field[bool]              `json:"archived,omitzero"`
	Barcodes     model.Field[[]model.Barcode]   `json:"barcodes,omitzero"`
	Attributes   model.Field[[]model.Attribute] `json:"attributes,omitzero"`
}
