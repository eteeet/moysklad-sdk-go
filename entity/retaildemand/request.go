package retaildemand

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: по методу необходима только retailShift; таблица дополнительно помечает agent и
// store «Необходимо при создании» — SDK следует методу и передаёт их, если заданы. Смена должна
// быть активной, а moment — позже moment смены (проверяет сервер). Суммы оплат должны в сумме
// совпадать с суммой продажи, смешанная оплата с QR и аванс вместе с предоплатой недопустимы —
// сервер распределяет остаток и отвечает ошибкой; SDK суммы не считает. sum и payedSum примеры
// передают, но таблица помечает readonly — их здесь нет. Positions — не больше 1000 строк.
// SDK не фискализирует продажу и не начисляет бонусы.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildemand#3-sozdat-roznichnuyu-prodazhu
type CreateRequest struct {
	RetailShift         model.Reference                        `json:"retailShift"`
	Name                model.Field[string]                    `json:"name,omitzero"`
	Description         model.Field[string]                    `json:"description,omitzero"`
	Code                model.Field[string]                    `json:"code,omitzero"`
	ExternalCode        model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID              model.Field[string]                    `json:"syncId,omitzero"`
	Moment              model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable          model.Field[bool]                      `json:"applicable,omitzero"`
	Shared              model.Field[bool]                      `json:"shared,omitzero"`
	OrganizationAccount model.Field[model.Reference]           `json:"organizationAccount,omitzero"`
	Agent               model.Field[model.Reference]           `json:"agent,omitzero"`
	AgentAccount        model.Field[model.Reference]           `json:"agentAccount,omitzero"`
	Store               model.Field[model.Reference]           `json:"store,omitzero"`
	Owner               model.Field[model.Reference]           `json:"owner,omitzero"`
	Group               model.Field[model.Reference]           `json:"group,omitzero"`
	State               model.Field[model.Reference]           `json:"state,omitzero"`
	CustomerOrder       model.Field[model.Reference]           `json:"customerOrder,omitzero"`
	Attributes          model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate                model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	TaxSystem           model.Field[string]                    `json:"taxSystem,omitzero"`
	VATEnabled          model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded         model.Field[bool]                      `json:"vatIncluded,omitzero"`
	CashSum             model.Field[float64]                   `json:"cashSum,omitzero"`
	NoCashSum           model.Field[float64]                   `json:"noCashSum,omitzero"`
	QRSum               model.Field[float64]                   `json:"qrSum,omitzero"`
	PrepaymentCashSum   model.Field[float64]                   `json:"prepaymentCashSum,omitzero"`
	PrepaymentNoCashSum model.Field[float64]                   `json:"prepaymentNoCashSum,omitzero"`
	PrepaymentQRSum     model.Field[float64]                   `json:"prepaymentQrSum,omitzero"`
	AdvancePaymentSum   model.Field[float64]                   `json:"advancePaymentSum,omitzero"`
	CheckSum            model.Field[float64]                   `json:"checkSum,omitzero"`
	CheckNumber         model.Field[int]                       `json:"checkNumber,omitzero"`
	DocumentNumber      model.Field[int]                       `json:"documentNumber,omitzero"`
	SessionNumber       model.Field[int]                       `json:"sessionNumber,omitzero"`
	GiftCards           model.Field[[]GiftCardRequest]         `json:"giftCards,omitzero"`
	Positions           model.Field[[]PositionRequest]         `json:"positions,omitzero"`
}

// UpdateRequest. При смене agent нужен и agentAccount (проверяет сервер). Positions заменяют все
// строки. organization и organizationBranch таблица помечает readonly, хотя текст метода
// упоминает смену organization — SDK их не пишет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildemand#3-izmenit-roznichnuyu-prodazhu
type UpdateRequest struct {
	RetailShift         model.Field[model.Reference]           `json:"retailShift,omitzero"`
	Name                model.Field[string]                    `json:"name,omitzero"`
	Description         model.Field[string]                    `json:"description,omitzero"`
	Code                model.Field[string]                    `json:"code,omitzero"`
	ExternalCode        model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID              model.Field[string]                    `json:"syncId,omitzero"`
	Moment              model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable          model.Field[bool]                      `json:"applicable,omitzero"`
	Shared              model.Field[bool]                      `json:"shared,omitzero"`
	OrganizationAccount model.Field[model.Reference]           `json:"organizationAccount,omitzero"`
	Agent               model.Field[model.Reference]           `json:"agent,omitzero"`
	AgentAccount        model.Field[model.Reference]           `json:"agentAccount,omitzero"`
	Store               model.Field[model.Reference]           `json:"store,omitzero"`
	Owner               model.Field[model.Reference]           `json:"owner,omitzero"`
	Group               model.Field[model.Reference]           `json:"group,omitzero"`
	State               model.Field[model.Reference]           `json:"state,omitzero"`
	CustomerOrder       model.Field[model.Reference]           `json:"customerOrder,omitzero"`
	Attributes          model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate                model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	TaxSystem           model.Field[string]                    `json:"taxSystem,omitzero"`
	VATEnabled          model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded         model.Field[bool]                      `json:"vatIncluded,omitzero"`
	CashSum             model.Field[float64]                   `json:"cashSum,omitzero"`
	NoCashSum           model.Field[float64]                   `json:"noCashSum,omitzero"`
	QRSum               model.Field[float64]                   `json:"qrSum,omitzero"`
	PrepaymentCashSum   model.Field[float64]                   `json:"prepaymentCashSum,omitzero"`
	PrepaymentNoCashSum model.Field[float64]                   `json:"prepaymentNoCashSum,omitzero"`
	PrepaymentQRSum     model.Field[float64]                   `json:"prepaymentQrSum,omitzero"`
	AdvancePaymentSum   model.Field[float64]                   `json:"advancePaymentSum,omitzero"`
	CheckSum            model.Field[float64]                   `json:"checkSum,omitzero"`
	CheckNumber         model.Field[int]                       `json:"checkNumber,omitzero"`
	DocumentNumber      model.Field[int]                       `json:"documentNumber,omitzero"`
	SessionNumber       model.Field[int]                       `json:"sessionNumber,omitzero"`
	GiftCards           model.Field[[]GiftCardRequest]         `json:"giftCards,omitzero"`
	Positions           model.Field[[]PositionRequest]         `json:"positions,omitzero"`
}
