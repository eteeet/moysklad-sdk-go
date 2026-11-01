package cashin

import "github.com/eteeet/moysklad-sdk-go/model"

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashin
type CashIn struct {
	ID             string                   `json:"id"`
	AccountID      string                   `json:"accountId"`
	Name           string                   `json:"name"`
	Description    string                   `json:"description"`
	Code           string                   `json:"code"`
	ExternalCode   string                   `json:"externalCode"`
	SyncID         string                   `json:"syncId"`
	Created        model.DateTime           `json:"created,omitzero"`
	Updated        model.DateTime           `json:"updated,omitzero"`
	Deleted        model.DateTime           `json:"deleted,omitzero"`
	Moment         model.DateTime           `json:"moment,omitzero"`
	Applicable     bool                     `json:"applicable"`
	Shared         bool                     `json:"shared"`
	Printed        bool                     `json:"printed"`
	Published      bool                     `json:"published"`
	Organization   *model.Reference         `json:"organization,omitzero"`
	Owner          *model.Reference         `json:"owner,omitzero"`
	Group          *model.Reference         `json:"group,omitzero"`
	Project        *model.Reference         `json:"project,omitzero"`
	State          *model.Reference         `json:"state,omitzero"`
	Meta           model.Meta               `json:"meta"`
	Attributes     []model.Attribute        `json:"attributes"`
	Rate           *model.DocumentRate      `json:"rate,omitzero"`
	Sum            float64                  `json:"sum"`
	Agent          *model.Reference         `json:"agent,omitzero"`
	Contract       *model.Reference         `json:"contract,omitzero"`
	VATSum         float64                  `json:"vatSum"`
	PaymentPurpose string                   `json:"paymentPurpose"`
	Operations     []model.FinanceOperation `json:"operations"`
	SalesChannel   *model.Reference         `json:"salesChannel,omitzero"`
	FactureOut     *model.Reference         `json:"factureOut,omitzero"`
}
