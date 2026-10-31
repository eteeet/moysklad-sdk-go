package expenseitem

import "github.com/eteeet/moysklad-sdk-go/model"

// У предустановленных ExpenseItem accountId, owner, group и shared могут отсутствовать.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/expenseitem
type ExpenseItem struct {
	Meta              model.Meta     `json:"meta"`
	ID                string         `json:"id"`
	AccountID         string         `json:"accountId"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Code              string         `json:"code"`
	ExternalCode      string         `json:"externalCode"`
	Shared            bool           `json:"shared"`
	Archived          bool           `json:"archived"`
	OperatingExpenses bool           `json:"operatingExpenses"`
	Updated           model.DateTime `json:"updated,omitzero"`

	Owner *model.Reference `json:"owner,omitzero"`
	Group *model.Reference `json:"group,omitzero"`
}
