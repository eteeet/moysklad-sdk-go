package invoicein

import "github.com/eteeet/moysklad-sdk-go/model"

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-in
type CreateRequest struct {
	Name                 string                                 `json:"name"`
	Description          model.Field[string]                    `json:"description,omitzero"`
	Code                 model.Field[string]                    `json:"code,omitzero"`
	ExternalCode         model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID               model.Field[string]                    `json:"syncId,omitzero"`
	Moment               model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable           model.Field[bool]                      `json:"applicable,omitzero"`
	Shared               model.Field[bool]                      `json:"shared,omitzero"`
	Organization         model.Reference                        `json:"organization"`
	Owner                model.Field[model.Reference]           `json:"owner,omitzero"`
	Group                model.Field[model.Reference]           `json:"group,omitzero"`
	Project              model.Field[model.Reference]           `json:"project,omitzero"`
	State                model.Field[model.Reference]           `json:"state,omitzero"`
	Attributes           model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate                 model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	Agent                model.Reference                        `json:"agent"`
	OrganizationAccount  model.Field[model.Reference]           `json:"organizationAccount,omitzero"`
	AgentAccount         model.Field[model.Reference]           `json:"agentAccount,omitzero"`
	Contract             model.Field[model.Reference]           `json:"contract,omitzero"`
	Positions            model.Field[[]PositionRequest]         `json:"positions,omitzero"`
	VATEnabled           model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded          model.Field[bool]                      `json:"vatIncluded,omitzero"`
	Store                model.Field[model.Reference]           `json:"store,omitzero"`
	PaymentPlannedMoment model.Field[model.DateTime]            `json:"paymentPlannedMoment,omitzero"`
	IncomingDate         model.Field[model.DateTime]            `json:"incomingDate,omitzero"`
	IncomingNumber       model.Field[string]                    `json:"incomingNumber,omitzero"`
	PurchaseOrder        model.Field[model.Reference]           `json:"purchaseOrder,omitzero"`
	Supplies             model.Field[[]model.Reference]         `json:"supplies,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-in
type UpdateRequest struct {
	Name                 model.Field[string]                    `json:"name,omitzero"`
	Description          model.Field[string]                    `json:"description,omitzero"`
	Code                 model.Field[string]                    `json:"code,omitzero"`
	ExternalCode         model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID               model.Field[string]                    `json:"syncId,omitzero"`
	Moment               model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable           model.Field[bool]                      `json:"applicable,omitzero"`
	Shared               model.Field[bool]                      `json:"shared,omitzero"`
	Organization         model.Field[model.Reference]           `json:"organization,omitzero"`
	Owner                model.Field[model.Reference]           `json:"owner,omitzero"`
	Group                model.Field[model.Reference]           `json:"group,omitzero"`
	Project              model.Field[model.Reference]           `json:"project,omitzero"`
	State                model.Field[model.Reference]           `json:"state,omitzero"`
	Attributes           model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate                 model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	Agent                model.Field[model.Reference]           `json:"agent,omitzero"`
	OrganizationAccount  model.Field[model.Reference]           `json:"organizationAccount,omitzero"`
	AgentAccount         model.Field[model.Reference]           `json:"agentAccount,omitzero"`
	Contract             model.Field[model.Reference]           `json:"contract,omitzero"`
	Positions            model.Field[[]PositionRequest]         `json:"positions,omitzero"`
	VATEnabled           model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded          model.Field[bool]                      `json:"vatIncluded,omitzero"`
	Store                model.Field[model.Reference]           `json:"store,omitzero"`
	PaymentPlannedMoment model.Field[model.DateTime]            `json:"paymentPlannedMoment,omitzero"`
	IncomingDate         model.Field[model.DateTime]            `json:"incomingDate,omitzero"`
	IncomingNumber       model.Field[string]                    `json:"incomingNumber,omitzero"`
	PurchaseOrder        model.Field[model.Reference]           `json:"purchaseOrder,omitzero"`
	Supplies             model.Field[[]model.Reference]         `json:"supplies,omitzero"`
}
