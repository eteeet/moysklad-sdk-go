package salesreturn

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// Общая форма строки товарных документов; у возврата — cost (без основания), country, gtd, slot,
// pack, things.
type Position = positions.Position
type PositionRequest = positions.PositionRequest

// SalesReturn — возврат покупателя. Модель чтения без expand: суммы — копейки и только для
// чтения; Demand — основание (у возврата без основания nil), Losses, Payments, FactureOut —
// связанные документы, пишутся в Create/Update.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/sales-return
type SalesReturn struct {
	Meta                model.Meta                         `json:"meta"`
	ID                  string                             `json:"id"`
	AccountID           string                             `json:"accountId"`
	Name                string                             `json:"name"`
	Description         string                             `json:"description"`
	Code                string                             `json:"code"`
	ExternalCode        string                             `json:"externalCode"`
	SyncID              string                             `json:"syncId"`
	Created             model.DateTime                     `json:"created,omitzero"`
	Updated             model.DateTime                     `json:"updated,omitzero"`
	Deleted             model.DateTime                     `json:"deleted,omitzero"`
	Moment              model.DateTime                     `json:"moment,omitzero"`
	Applicable          bool                               `json:"applicable"`
	Shared              bool                               `json:"shared"`
	Printed             bool                               `json:"printed"`
	Published           bool                               `json:"published"`
	Organization        *model.Reference                   `json:"organization,omitzero"`
	OrganizationAccount *model.Reference                   `json:"organizationAccount,omitzero"`
	Agent               *model.Reference                   `json:"agent,omitzero"`
	AgentAccount        *model.Reference                   `json:"agentAccount,omitzero"`
	Contract            *model.Reference                   `json:"contract,omitzero"`
	Store               *model.Reference                   `json:"store,omitzero"`
	Owner               *model.Reference                   `json:"owner,omitzero"`
	Group               *model.Reference                   `json:"group,omitzero"`
	Project             *model.Reference                   `json:"project,omitzero"`
	State               *model.Reference                   `json:"state,omitzero"`
	SalesChannel        *model.Reference                   `json:"salesChannel,omitzero"`
	Attributes          []model.Attribute                  `json:"attributes"`
	Rate                *model.DocumentRate                `json:"rate,omitzero"`
	VATEnabled          bool                               `json:"vatEnabled"`
	VATIncluded         bool                               `json:"vatIncluded"`
	Sum                 float64                            `json:"sum"`
	VATSum              float64                            `json:"vatSum"`
	PayedSum            float64                            `json:"payedSum"`
	Positions           *model.Collection[Position]        `json:"positions,omitzero"`
	Files               *model.Collection[model.Reference] `json:"files,omitzero"`

	Demand     *model.Reference  `json:"demand,omitzero"`
	Losses     []model.Reference `json:"losses"`
	Payments   []model.Reference `json:"payments"`
	FactureOut *model.Reference  `json:"factureOut,omitzero"`
}
