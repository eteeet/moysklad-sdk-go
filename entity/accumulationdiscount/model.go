package accumulationdiscount

import (
	"github.com/eteeet/moysklad-sdk-go/entity/discount"
	"github.com/eteeet/moysklad-sdk-go/model"
)

// AccumulationDiscount — накопительная скидка: уровни Levels (сумма накоплений в копейках → процент).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount
type AccumulationDiscount struct {
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
	Levels         []discount.Level  `json:"levels"`
}
