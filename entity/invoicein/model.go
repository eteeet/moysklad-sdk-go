package invoicein

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// Общая форма строки для всех товарных документов.
type Position = positions.Position
type PositionRequest = positions.PositionRequest

// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-in
type InvoiceIn struct {
	ID                   string                      `json:"id"`
	AccountID            string                      `json:"accountId"`
	Name                 string                      `json:"name"`
	Description          string                      `json:"description"`
	Code                 string                      `json:"code"`
	ExternalCode         string                      `json:"externalCode"`
	SyncID               string                      `json:"syncId"`
	Created              model.DateTime              `json:"created,omitzero"`
	Updated              model.DateTime              `json:"updated,omitzero"`
	Deleted              model.DateTime              `json:"deleted,omitzero"`
	Moment               model.DateTime              `json:"moment,omitzero"`
	Applicable           bool                        `json:"applicable"`
	Shared               bool                        `json:"shared"`
	Printed              bool                        `json:"printed"`
	Published            bool                        `json:"published"`
	Organization         *model.Reference            `json:"organization,omitzero"`
	Owner                *model.Reference            `json:"owner,omitzero"`
	Group                *model.Reference            `json:"group,omitzero"`
	Project              *model.Reference            `json:"project,omitzero"`
	State                *model.Reference            `json:"state,omitzero"`
	Meta                 model.Meta                  `json:"meta"`
	Attributes           []model.Attribute           `json:"attributes"`
	Rate                 *model.DocumentRate         `json:"rate,omitzero"`
	Sum                  float64                     `json:"sum"`
	Agent                *model.Reference            `json:"agent,omitzero"`
	OrganizationAccount  *model.Reference            `json:"organizationAccount,omitzero"`
	AgentAccount         *model.Reference            `json:"agentAccount,omitzero"`
	Contract             *model.Reference            `json:"contract,omitzero"`
	VATSum               float64                     `json:"vatSum"`
	Positions            *model.Collection[Position] `json:"positions,omitzero"`
	VATEnabled           bool                        `json:"vatEnabled"`
	VATIncluded          bool                        `json:"vatIncluded"`
	PayedSum             float64                     `json:"payedSum"`
	Store                *model.Reference            `json:"store,omitzero"`
	Payments             []model.Reference           `json:"payments"`
	PaymentPlannedMoment model.DateTime              `json:"paymentPlannedMoment,omitzero"`
	ShippedSum           float64                     `json:"shippedSum"`
	IncomingDate         model.DateTime              `json:"incomingDate,omitzero"`
	IncomingNumber       string                      `json:"incomingNumber"`
	PurchaseOrder        *model.Reference            `json:"purchaseOrder,omitzero"`
	Supplies             []model.Reference           `json:"supplies"`
}
