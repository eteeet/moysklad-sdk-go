package positions

import "github.com/eteeet/moysklad-sdk-go/model"

// ID сопоставляет существующую строку embedded Update; для новой строки его не задают.
type CoreRequest struct {
	ID model.Field[string] `json:"id,omitzero"`
}

type ItemRequest struct {
	CoreRequest
	Assortment model.Field[model.Reference] `json:"assortment,omitzero"`
	Quantity   model.Field[float64]         `json:"quantity,omitzero"`
}

// CommercialRequest не задаёт правила возврата: для него документ определяет свою форму.
// Price задаётся в копейках; допустимость Null определяет API документа.
type CommercialRequest struct {
	ItemRequest
	Price      model.Field[float64] `json:"price,omitzero"`
	Discount   model.Field[float64] `json:"discount,omitzero"`
	VAT        model.Field[int]     `json:"vat,omitzero"`
	VATEnabled model.Field[bool]    `json:"vatEnabled,omitzero"`
}
