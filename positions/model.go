package positions

import "github.com/eteeet/moysklad-sdk-go/model"

type Core struct {
	Meta      model.Meta `json:"meta"`
	ID        string     `json:"id"`
	AccountID string     `json:"accountId"`
}

type Item struct {
	Core
	Assortment model.Reference `json:"assortment"`
	Quantity   float64         `json:"quantity"`
}

// Price, Overhead и Cost — в копейках; Overhead доступен только для чтения.
type Position struct {
	Item
	Price      float64          `json:"price"`
	Discount   float64          `json:"discount"`
	VAT        int              `json:"vat"`
	VATEnabled bool             `json:"vatEnabled"`
	Cost       float64          `json:"cost"`
	Shipped    float64          `json:"shipped"`
	InTransit  float64          `json:"inTransit"`
	Overhead   float64          `json:"overhead"`
	SourceSlot *model.Reference `json:"sourceSlot,omitzero"`
	TargetSlot *model.Reference `json:"targetSlot,omitzero"`
}
