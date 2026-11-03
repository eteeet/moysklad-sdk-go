package personaldiscount

import (
	"encoding/json/v2"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// PersonalDiscount — персональная скидка; процент для контрагента задаётся в Counterparty.Discounts.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount
type PersonalDiscount struct {
	Meta           model.Meta        `json:"meta"`
	ID             string            `json:"id"`
	AccountID      string            `json:"accountId"`
	Name           string            `json:"name"`
	Active         bool              `json:"active"`
	AllProducts    bool              `json:"allProducts"`
	AllAgents      bool              `json:"allAgents"`
	AgentTags      []string          `json:"agentTags"`
	Assortment     []model.Reference `json:"assortment"`
	ProductFolders []model.Reference `json:"productFolders"`
}

// UnmarshalJSON читает группы товаров по ключу productFolders (примеры, фикстуры) и productfolders
// (таблица источника); при обоих приоритет у productFolders.
func (d *PersonalDiscount) UnmarshalJSON(data []byte) error {
	type plain PersonalDiscount
	var v struct {
		plain
		Lower []model.Reference `json:"productfolders"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*d = PersonalDiscount(v.plain)
	if d.ProductFolders == nil {
		d.ProductFolders = v.Lower
	}
	return nil
}
