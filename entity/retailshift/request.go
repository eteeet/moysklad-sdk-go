package retailshift

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы organization и retailStore. Store без значения сервер берёт из точки
// продаж; комиссии без значения — 0, и они не могут превышать безналичные платежи (проверяет
// сервер). SDK не создаёт и не ищет смену сам: открытая смена — только результат этого вызова.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retailshift#3-sozdat-roznichnuyu-smenu
type CreateRequest struct {
	Organization        model.Reference                `json:"organization"`
	RetailStore         model.Reference                `json:"retailStore"`
	Name                model.Field[string]            `json:"name,omitzero"`
	Description         model.Field[string]            `json:"description,omitzero"`
	ExternalCode        model.Field[string]            `json:"externalCode,omitzero"`
	SyncID              model.Field[string]            `json:"syncId,omitzero"`
	Moment              model.Field[model.DateTime]    `json:"moment,omitzero"`
	CloseDate           model.Field[model.DateTime]    `json:"closeDate,omitzero"`
	Shared              model.Field[bool]              `json:"shared,omitzero"`
	OrganizationAccount model.Field[model.Reference]   `json:"organizationAccount,omitzero"`
	OrganizationBranch  model.Field[model.Reference]   `json:"organizationBranch,omitzero"`
	Store               model.Field[model.Reference]   `json:"store,omitzero"`
	Owner               model.Field[model.Reference]   `json:"owner,omitzero"`
	Group               model.Field[model.Reference]   `json:"group,omitzero"`
	Acquire             model.Field[model.Reference]   `json:"acquire,omitzero"`
	QRAcquire           model.Field[model.Reference]   `json:"qrAcquire,omitzero"`
	Attributes          model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	BankComission       model.Field[float64]           `json:"bankComission,omitzero"`
	BankPercent         model.Field[float64]           `json:"bankPercent,omitzero"`
	QRBankComission     model.Field[float64]           `json:"qrBankComission,omitzero"`
	QRBankPercent       model.Field[float64]           `json:"qrBankPercent,omitzero"`
	Cheque              model.Field[ChequeRequest]     `json:"cheque,omitzero"`
}

// UpdateRequest. Смену закрывает явный CloseDate; фискальные данные — явный Cheque. syncId после
// заполнения не меняется (проверяет сервер). Null у organization и retailStore отправляется как есть.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retailshift#3-izmenit-roznichnuyu-smenu
type UpdateRequest struct {
	Organization        model.Field[model.Reference]   `json:"organization,omitzero"`
	RetailStore         model.Field[model.Reference]   `json:"retailStore,omitzero"`
	Name                model.Field[string]            `json:"name,omitzero"`
	Description         model.Field[string]            `json:"description,omitzero"`
	ExternalCode        model.Field[string]            `json:"externalCode,omitzero"`
	SyncID              model.Field[string]            `json:"syncId,omitzero"`
	Moment              model.Field[model.DateTime]    `json:"moment,omitzero"`
	CloseDate           model.Field[model.DateTime]    `json:"closeDate,omitzero"`
	Shared              model.Field[bool]              `json:"shared,omitzero"`
	OrganizationAccount model.Field[model.Reference]   `json:"organizationAccount,omitzero"`
	OrganizationBranch  model.Field[model.Reference]   `json:"organizationBranch,omitzero"`
	Store               model.Field[model.Reference]   `json:"store,omitzero"`
	Owner               model.Field[model.Reference]   `json:"owner,omitzero"`
	Group               model.Field[model.Reference]   `json:"group,omitzero"`
	Acquire             model.Field[model.Reference]   `json:"acquire,omitzero"`
	QRAcquire           model.Field[model.Reference]   `json:"qrAcquire,omitzero"`
	Attributes          model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	BankComission       model.Field[float64]           `json:"bankComission,omitzero"`
	BankPercent         model.Field[float64]           `json:"bankPercent,omitzero"`
	QRBankComission     model.Field[float64]           `json:"qrBankComission,omitzero"`
	QRBankPercent       model.Field[float64]           `json:"qrBankPercent,omitzero"`
	Cheque              model.Field[ChequeRequest]     `json:"cheque,omitzero"`
}
