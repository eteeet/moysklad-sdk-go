package retailsalesreturn

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// Общая форма строки товарных документов; у розничного возврата — cost (без основания), pack,
// things.
type Position = positions.Position
type PositionRequest = positions.PositionRequest

// RetailSalesReturn — розничный возврат. Модель чтения без expand: суммы — копейки, sum и vatSum
// только для чтения. Demand — основание, розничная продажа (не Demand-отгрузка); у возврата без
// основания nil.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retail-sales-return
type RetailSalesReturn struct {
	Meta                model.Meta                  `json:"meta"`
	ID                  string                      `json:"id"`
	AccountID           string                      `json:"accountId"`
	Name                string                      `json:"name"`
	Description         string                      `json:"description"`
	Code                string                      `json:"code"`
	ExternalCode        string                      `json:"externalCode"`
	SyncID              string                      `json:"syncId"`
	Created             model.DateTime              `json:"created,omitzero"`
	Updated             model.DateTime              `json:"updated,omitzero"`
	Deleted             model.DateTime              `json:"deleted,omitzero"`
	Moment              model.DateTime              `json:"moment,omitzero"`
	Applicable          bool                        `json:"applicable"`
	Shared              bool                        `json:"shared"`
	Printed             bool                        `json:"printed"`
	Published           bool                        `json:"published"`
	RetailShift         *model.Reference            `json:"retailShift,omitzero"`
	RetailStore         *model.Reference            `json:"retailStore,omitzero"`
	Demand              *model.Reference            `json:"demand,omitzero"`
	Organization        *model.Reference            `json:"organization,omitzero"`
	OrganizationAccount *model.Reference            `json:"organizationAccount,omitzero"`
	OrganizationBranch  *model.Reference            `json:"organizationBranch,omitzero"`
	Agent               *model.Reference            `json:"agent,omitzero"`
	AgentAccount        *model.Reference            `json:"agentAccount,omitzero"`
	Contract            *model.Reference            `json:"contract,omitzero"`
	Store               *model.Reference            `json:"store,omitzero"`
	Owner               *model.Reference            `json:"owner,omitzero"`
	Group               *model.Reference            `json:"group,omitzero"`
	State               *model.Reference            `json:"state,omitzero"`
	Attributes          []model.Attribute           `json:"attributes"`
	Rate                *model.DocumentRate         `json:"rate,omitzero"`
	TaxSystem           string                      `json:"taxSystem"`
	VATEnabled          bool                        `json:"vatEnabled"`
	VATIncluded         bool                        `json:"vatIncluded"`
	Sum                 float64                     `json:"sum"`
	VATSum              float64                     `json:"vatSum"`
	CashSum             float64                     `json:"cashSum"`
	NoCashSum           float64                     `json:"noCashSum"`
	QRSum               float64                     `json:"qrSum"`
	Positions           *model.Collection[Position] `json:"positions,omitzero"`
}
