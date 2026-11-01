package taxrate

import "github.com/eteeet/moysklad-sdk-go/model"

// TaxRate — ставка НДС. Rate — процент, а не денежная сумма. Owner, Group и Shared есть только
// у пользовательских ставок.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/taxrate
type TaxRate struct {
	Meta      model.Meta       `json:"meta"`
	ID        string           `json:"id"`
	AccountID string           `json:"accountId"`
	Updated   model.DateTime   `json:"updated,omitzero"`
	Rate      float64          `json:"rate"`
	Comment   string           `json:"comment"`
	Archived  bool             `json:"archived"`
	Shared    bool             `json:"shared"`
	Owner     *model.Reference `json:"owner,omitzero"`
	Group     *model.Reference `json:"group,omitzero"`
}
