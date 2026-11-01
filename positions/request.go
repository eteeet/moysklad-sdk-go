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

// Price задаётся в копейках; допустимость Null определяет API документа.
type PositionRequest struct {
	ItemRequest
	Price      model.Field[float64] `json:"price,omitzero"`
	Discount   model.Field[float64] `json:"discount,omitzero"`
	VAT        model.Field[int]     `json:"vat,omitzero"`
	VATEnabled model.Field[bool]    `json:"vatEnabled,omitzero"`
	// Cost — себестоимость (услуги Demand; возвратов без основания), InTransit — количество
	// в пути PurchaseOrder.
	Cost       model.Field[float64]         `json:"cost,omitzero"`
	InTransit  model.Field[float64]         `json:"inTransit,omitzero"`
	SourceSlot model.Field[model.Reference] `json:"sourceSlot,omitzero"`
	TargetSlot model.Field[model.Reference] `json:"targetSlot,omitzero"`

	// Поля строк заказа покупателя, оприходования, списания, возвратов и розницы; какие из них
	// принимает конкретный документ, определяет API. TrackingCodes: cis и type необходимы при
	// создании (раздел кодов маркировки retaildemand); trackingCodes_1162 только для чтения.
	Reserve       model.Field[float64]         `json:"reserve,omitzero"`
	TaxSystem     model.Field[string]          `json:"taxSystem,omitzero"`
	Reason        model.Field[string]          `json:"reason,omitzero"`
	Slot          model.Field[model.Reference] `json:"slot,omitzero"`
	Pack          model.Field[PackRequest]     `json:"pack,omitzero"`
	Country       model.Field[model.Reference] `json:"country,omitzero"`
	GTD           model.Field[GTD]             `json:"gtd,omitzero"`
	Things        model.Field[[]string]        `json:"things,omitzero"`
	TrackingCodes model.Field[[]TrackingCode]  `json:"trackingCodes,omitzero"`
}

// PackRequest ссылается на упаковку товара по её id.
type PackRequest struct {
	ID string `json:"id"`
}
