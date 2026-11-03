package factureout

import "github.com/eteeet/moysklad-sdk-go/model"

// FactureOut — счёт-фактура выданный. Модель чтения без expand и без позиций: sum и vatSum —
// копейки, readonly; AdvancePaymentVat — ставка НДС в процентах, не деньги. Demands (отгрузки),
// Returns (возвраты поставщикам) и Payments (входящие платежи) — связи с документами-основаниями.
// AdvancePaymentVat, PaymentPurpose и VATSum сервер возвращает только у счёта-фактуры с платежами.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/factureout
type FactureOut struct {
	Meta              model.Meta                         `json:"meta"`
	ID                string                             `json:"id"`
	AccountID         string                             `json:"accountId"`
	Name              string                             `json:"name"`
	Description       string                             `json:"description"`
	Code              string                             `json:"code"`
	ExternalCode      string                             `json:"externalCode"`
	SyncID            string                             `json:"syncId"`
	Created           model.DateTime                     `json:"created,omitzero"`
	Updated           model.DateTime                     `json:"updated,omitzero"`
	Deleted           model.DateTime                     `json:"deleted,omitzero"`
	Moment            model.DateTime                     `json:"moment,omitzero"`
	Applicable        bool                               `json:"applicable"`
	Shared            bool                               `json:"shared"`
	Printed           bool                               `json:"printed"`
	Published         bool                               `json:"published"`
	Organization      *model.Reference                   `json:"organization,omitzero"`
	Agent             *model.Reference                   `json:"agent,omitzero"`
	Contract          *model.Reference                   `json:"contract,omitzero"`
	Owner             *model.Reference                   `json:"owner,omitzero"`
	Group             *model.Reference                   `json:"group,omitzero"`
	State             *model.Reference                   `json:"state,omitzero"`
	Attributes        []model.Attribute                  `json:"attributes"`
	Rate              *model.DocumentRate                `json:"rate,omitzero"`
	Sum               float64                            `json:"sum"`
	Files             *model.Collection[model.Reference] `json:"files,omitzero"`
	StateContractID   string                             `json:"stateContractId"`
	AdvancePaymentVat int                                `json:"advancePaymentVat"`
	PaymentPurpose    string                             `json:"paymentPurpose"`
	VATSum            float64                            `json:"vatSum"`
	Demands           []model.Reference                  `json:"demands"`
	Payments          []model.Reference                  `json:"payments"`
	Returns           []model.Reference                  `json:"returns"`
	Consignee         *model.Reference                   `json:"consignee,omitzero"`
	ConsigneeBranch   *model.Reference                   `json:"consigneeBranch,omitzero"`
	PaymentNumber     string                             `json:"paymentNumber"`
	PaymentDate       model.DateTime                     `json:"paymentDate,omitzero"`
}
