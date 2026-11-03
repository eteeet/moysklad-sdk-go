package personaldiscount

import "github.com/eteeet/moysklad-sdk-go/model"

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
