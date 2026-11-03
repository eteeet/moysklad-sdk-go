package specialpricediscount

import (
	"encoding/json/v2"

	"github.com/eteeet/moysklad-sdk-go/entity/discount"
	"github.com/eteeet/moysklad-sdk-go/model"
)

// SpecialPriceDiscount — спеццена: фиксированный процент Discount или тип цены SpecialPrice
// (UsePriceType = true); значение спеццены — копейки.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount
type SpecialPriceDiscount struct {
	Meta           model.Meta             `json:"meta"`
	ID             string                 `json:"id"`
	AccountID      string                 `json:"accountId"`
	Name           string                 `json:"name"`
	Active         bool                   `json:"active"`
	AllProducts    bool                   `json:"allProducts"`
	AllAgents      bool                   `json:"allAgents"`
	UsePriceType   bool                   `json:"usePriceType"`
	AgentTags      []string               `json:"agentTags"`
	Assortment     []model.Reference      `json:"assortment"`
	ProductFolders []model.Reference      `json:"productFolders"`
	Discount       float64                `json:"discount"`
	SpecialPrice   *discount.SpecialPrice `json:"specialPrice,omitzero"`
}

// UnmarshalJSON читает группы товаров по ключу productFolders (примеры, фикстуры) и productfolders
// (таблица источника); при обоих приоритет у productFolders.
func (d *SpecialPriceDiscount) UnmarshalJSON(data []byte) error {
	type plain SpecialPriceDiscount
	var v struct {
		plain
		Lower []model.Reference `json:"productfolders"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*d = SpecialPriceDiscount(v.plain)
	if d.ProductFolders == nil {
		d.ProductFolders = v.Lower
	}
	return nil
}
