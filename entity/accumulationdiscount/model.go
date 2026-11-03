package accumulationdiscount

import (
	"encoding/json/v2"

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

// UnmarshalJSON читает группы товаров по ключу productFolders (примеры, фикстуры) и productfolders
// (таблица источника); при обоих приоритет у productFolders.
func (d *AccumulationDiscount) UnmarshalJSON(data []byte) error {
	type plain AccumulationDiscount
	var v struct {
		plain
		Lower []model.Reference `json:"productfolders"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*d = AccumulationDiscount(v.plain)
	if d.ProductFolders == nil {
		d.ProductFolders = v.Lower
	}
	return nil
}
