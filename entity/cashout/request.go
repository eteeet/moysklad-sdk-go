package cashout

import "github.com/eteeet/moysklad-sdk-go/model"

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashout
type CreateRequest struct {
	Name           model.Field[string]                          `json:"name,omitzero"`
	Description    model.Field[string]                          `json:"description,omitzero"`
	Code           model.Field[string]                          `json:"code,omitzero"`
	ExternalCode   model.Field[string]                          `json:"externalCode,omitzero"`
	SyncID         model.Field[string]                          `json:"syncId,omitzero"`
	Moment         model.Field[model.DateTime]                  `json:"moment,omitzero"`
	Applicable     model.Field[bool]                            `json:"applicable,omitzero"`
	Shared         model.Field[bool]                            `json:"shared,omitzero"`
	Organization   model.Reference                              `json:"organization"`
	Owner          model.Field[model.Reference]                 `json:"owner,omitzero"`
	Group          model.Field[model.Reference]                 `json:"group,omitzero"`
	Project        model.Field[model.Reference]                 `json:"project,omitzero"`
	State          model.Field[model.Reference]                 `json:"state,omitzero"`
	Attributes     model.Field[[]model.Attribute]               `json:"attributes,omitzero"`
	Rate           model.Field[model.DocumentRateRequest]       `json:"rate,omitzero"`
	Sum            model.Field[float64]                         `json:"sum,omitzero"`
	Agent          model.Reference                              `json:"agent"`
	Contract       model.Field[model.Reference]                 `json:"contract,omitzero"`
	VATSum         model.Field[float64]                         `json:"vatSum,omitzero"`
	PaymentPurpose model.Field[string]                          `json:"paymentPurpose,omitzero"`
	Operations     model.Field[[]model.FinanceOperationRequest] `json:"operations,omitzero"`
	SalesChannel   model.Field[model.Reference]                 `json:"salesChannel,omitzero"`
	ExpenseItem    model.Reference                              `json:"expenseItem"`
	FactureIn      model.Field[model.Reference]                 `json:"factureIn,omitzero"`
	NoClosingDocs  model.Field[bool]                            `json:"noClosingDocs,omitzero"`
	AccrualDate    model.Field[model.DateTime]                  `json:"accrualDate,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashout
type UpdateRequest struct {
	Name           model.Field[string]                          `json:"name,omitzero"`
	Description    model.Field[string]                          `json:"description,omitzero"`
	Code           model.Field[string]                          `json:"code,omitzero"`
	ExternalCode   model.Field[string]                          `json:"externalCode,omitzero"`
	SyncID         model.Field[string]                          `json:"syncId,omitzero"`
	Moment         model.Field[model.DateTime]                  `json:"moment,omitzero"`
	Applicable     model.Field[bool]                            `json:"applicable,omitzero"`
	Shared         model.Field[bool]                            `json:"shared,omitzero"`
	Organization   model.Field[model.Reference]                 `json:"organization,omitzero"`
	Owner          model.Field[model.Reference]                 `json:"owner,omitzero"`
	Group          model.Field[model.Reference]                 `json:"group,omitzero"`
	Project        model.Field[model.Reference]                 `json:"project,omitzero"`
	State          model.Field[model.Reference]                 `json:"state,omitzero"`
	Attributes     model.Field[[]model.Attribute]               `json:"attributes,omitzero"`
	Rate           model.Field[model.DocumentRateRequest]       `json:"rate,omitzero"`
	Sum            model.Field[float64]                         `json:"sum,omitzero"`
	Agent          model.Field[model.Reference]                 `json:"agent,omitzero"`
	Contract       model.Field[model.Reference]                 `json:"contract,omitzero"`
	VATSum         model.Field[float64]                         `json:"vatSum,omitzero"`
	PaymentPurpose model.Field[string]                          `json:"paymentPurpose,omitzero"`
	Operations     model.Field[[]model.FinanceOperationRequest] `json:"operations,omitzero"`
	SalesChannel   model.Field[model.Reference]                 `json:"salesChannel,omitzero"`
	ExpenseItem    model.Field[model.Reference]                 `json:"expenseItem,omitzero"`
	FactureIn      model.Field[model.Reference]                 `json:"factureIn,omitzero"`
	NoClosingDocs  model.Field[bool]                            `json:"noClosingDocs,omitzero"`
	AccrualDate    model.Field[model.DateTime]                  `json:"accrualDate,omitzero"`
}
