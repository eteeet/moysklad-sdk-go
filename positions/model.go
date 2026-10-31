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

// Price задаётся в копейках без масштабирования.
type Commercial struct {
	Item
	Price      float64 `json:"price"`
	Discount   float64 `json:"discount"`
	VAT        int     `json:"vat"`
	VATEnabled bool    `json:"vatEnabled"`
}
