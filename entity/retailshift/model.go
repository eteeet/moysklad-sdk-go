package retailshift

import "github.com/eteeet/moysklad-sdk-go/model"

// RetailShift — розничная смена. Модель чтения без expand. Выручка и полученные суммы — копейки;
// для комиссий эквайера (bankComission, qrBankComission) источник единицу не уточняет, SDK их не
// пересчитывает; bankPercent и qrBankPercent — проценты. Позиций у смены нет: связанные продажи
// и возвраты — ссылки Operations, платежи — PaymentOperations с linkedSum; всё только для чтения.
// Поле updated таблица не помечает readonly, но и не описывает записи — SDK его не пишет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retailshift
type RetailShift struct {
	Meta                model.Meta                         `json:"meta"`
	ID                  string                             `json:"id"`
	AccountID           string                             `json:"accountId"`
	Name                string                             `json:"name"`
	Description         string                             `json:"description"`
	ExternalCode        string                             `json:"externalCode"`
	SyncID              string                             `json:"syncId"`
	Created             model.DateTime                     `json:"created,omitzero"`
	Updated             model.DateTime                     `json:"updated,omitzero"`
	Deleted             model.DateTime                     `json:"deleted,omitzero"`
	Moment              model.DateTime                     `json:"moment,omitzero"`
	CloseDate           model.DateTime                     `json:"closeDate,omitzero"`
	Shared              bool                               `json:"shared"`
	Printed             bool                               `json:"printed"`
	Published           bool                               `json:"published"`
	Organization        *model.Reference                   `json:"organization,omitzero"`
	OrganizationAccount *model.Reference                   `json:"organizationAccount,omitzero"`
	OrganizationBranch  *model.Reference                   `json:"organizationBranch,omitzero"`
	RetailStore         *model.Reference                   `json:"retailStore,omitzero"`
	Store               *model.Reference                   `json:"store,omitzero"`
	Owner               *model.Reference                   `json:"owner,omitzero"`
	Group               *model.Reference                   `json:"group,omitzero"`
	Acquire             *model.Reference                   `json:"acquire,omitzero"`
	QRAcquire           *model.Reference                   `json:"qrAcquire,omitzero"`
	Attributes          []model.Attribute                  `json:"attributes"`
	BankComission       float64                            `json:"bankComission"`
	BankPercent         float64                            `json:"bankPercent"`
	QRBankComission     float64                            `json:"qrBankComission"`
	QRBankPercent       float64                            `json:"qrBankPercent"`
	ProceedsCash        float64                            `json:"proceedsCash"`
	ProceedsNoCash      float64                            `json:"proceedsNoCash"`
	ReceivedCash        float64                            `json:"receivedCash"`
	ReceivedNoCash      float64                            `json:"receivedNoCash"`
	Cheque              *Cheque                            `json:"cheque,omitzero"`
	Operations          []model.Reference                  `json:"operations"`
	PaymentOperations   []model.FinanceOperation           `json:"paymentOperations"`
	Files               *model.Collection[model.Reference] `json:"files,omitzero"`
}

// Cheque — данные ККТ об открытии и закрытии смены.
type Cheque struct {
	Start *ChequeStart `json:"start,omitzero"`
	End   *ChequeEnd   `json:"end,omitzero"`
}

// ChequeStart — открытие смены ККТ.
type ChequeStart struct {
	FNNumber        string         `json:"fnNumber,omitempty"`
	KKTRegNumber    string         `json:"kktRegNumber,omitempty"`
	FiscalDocSign   string         `json:"fiscalDocSign,omitempty"`
	ShiftNumber     string         `json:"shiftNumber,omitempty"`
	FiscalDocNumber string         `json:"fiscalDocNumber,omitempty"`
	Time            model.DateTime `json:"time,omitzero"`
}

// ChequeEnd — закрытие смены ККТ: поля открытия и количества чеков и фискальных документов.
type ChequeEnd struct {
	ChequeStart
	ChequesTotal    int `json:"chequesTotal,omitempty"`
	FiscalDocsTotal int `json:"fiscalDocsTotal,omitempty"`
}

// ChequeRequest — данные ККТ при записи; незаданный Field не отправляется, 0 и "" — отправляются.
type ChequeRequest struct {
	Start model.Field[ChequeStartRequest] `json:"start,omitzero"`
	End   model.Field[ChequeEndRequest]   `json:"end,omitzero"`
}

type ChequeStartRequest struct {
	FNNumber        model.Field[string]         `json:"fnNumber,omitzero"`
	KKTRegNumber    model.Field[string]         `json:"kktRegNumber,omitzero"`
	FiscalDocSign   model.Field[string]         `json:"fiscalDocSign,omitzero"`
	ShiftNumber     model.Field[string]         `json:"shiftNumber,omitzero"`
	FiscalDocNumber model.Field[string]         `json:"fiscalDocNumber,omitzero"`
	Time            model.Field[model.DateTime] `json:"time,omitzero"`
}

type ChequeEndRequest struct {
	ChequeStartRequest
	ChequesTotal    model.Field[int] `json:"chequesTotal,omitzero"`
	FiscalDocsTotal model.Field[int] `json:"fiscalDocsTotal,omitzero"`
}
