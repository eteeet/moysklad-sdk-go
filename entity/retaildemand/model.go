package retaildemand

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// Общая форма строки товарных документов; у розничной продажи — cost (услуги), pack, things,
// trackingCodes, trackingCodes_1162 и declaration (readonly).
type Position = positions.Position
type PositionRequest = positions.PositionRequest

// RetailDemand — розничная продажа. Модель чтения без expand: все суммы — копейки; sum, vatSum,
// payedSum и фискальный Cheque только для чтения. Organization и OrganizationBranch таблица
// помечает readonly — их задаёт смена.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildemand
type RetailDemand struct {
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
	RetailShift         *model.Reference                   `json:"retailShift,omitzero"`
	RetailStore         *model.Reference                   `json:"retailStore,omitzero"`
	Organization        *model.Reference                   `json:"organization,omitzero"`
	OrganizationAccount *model.Reference                   `json:"organizationAccount,omitzero"`
	OrganizationBranch  *model.Reference                   `json:"organizationBranch,omitzero"`
	Agent               *model.Reference                   `json:"agent,omitzero"`
	AgentAccount        *model.Reference                   `json:"agentAccount,omitzero"`
	Store               *model.Reference                   `json:"store,omitzero"`
	Owner               *model.Reference                   `json:"owner,omitzero"`
	Group               *model.Reference                   `json:"group,omitzero"`
	State               *model.Reference                   `json:"state,omitzero"`
	CustomerOrder       *model.Reference                   `json:"customerOrder,omitzero"`
	Attributes          []model.Attribute                  `json:"attributes"`
	Rate                *model.DocumentRate                `json:"rate,omitzero"`
	TaxSystem           string                             `json:"taxSystem"`
	VATEnabled          bool                               `json:"vatEnabled"`
	VATIncluded         bool                               `json:"vatIncluded"`
	Sum                 float64                            `json:"sum"`
	VATSum              float64                            `json:"vatSum"`
	PayedSum            float64                            `json:"payedSum"`
	CashSum             float64                            `json:"cashSum"`
	NoCashSum           float64                            `json:"noCashSum"`
	QRSum               float64                            `json:"qrSum"`
	PrepaymentCashSum   float64                            `json:"prepaymentCashSum"`
	PrepaymentNoCashSum float64                            `json:"prepaymentNoCashSum"`
	PrepaymentQRSum     float64                            `json:"prepaymentQrSum"`
	AdvancePaymentSum   float64                            `json:"advancePaymentSum"`
	CheckSum            float64                            `json:"checkSum"`
	CheckNumber         int                                `json:"checkNumber"`
	DocumentNumber      int                                `json:"documentNumber"`
	SessionNumber       int                                `json:"sessionNumber"`
	GiftCards           []GiftCard                         `json:"giftCards"`
	Cheque              *Cheque                            `json:"cheque,omitzero"`
	Positions           *model.Collection[Position]        `json:"positions,omitzero"`
	Files               *model.Collection[model.Reference] `json:"files,omitzero"`
}

// GiftCard — подарочный сертификат в оплате; PaymentSum (копейки) только для чтения.
type GiftCard struct {
	Name       string  `json:"name"`
	PaymentSum float64 `json:"paymentSum"`
}

// GiftCardRequest — сертификат при записи: только номер.
type GiftCardRequest struct {
	Name string `json:"name"`
}

// Cheque — данные фискального чека, только для чтения.
type Cheque struct {
	FNNumber        string         `json:"fnNumber,omitempty"`
	FiscalDocSign   string         `json:"fiscalDocSign,omitempty"`
	FiscalDocNumber string         `json:"fiscalDocNumber,omitempty"`
	KKTRegNumber    string         `json:"kktRegNumber,omitempty"`
	Time            model.DateTime `json:"time,omitzero"`
}
