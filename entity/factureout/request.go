package factureout

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: organization обязателен по таблице атрибутов («Необходимо при создании»), хотя
// пример «только необходимые поля» его не передаёт — SDK следует таблице. Счёт-фактура создаётся
// только на основании: ровно один документ в Demands (отгрузка), Returns (возврат поставщику) или
// Payments (входящий платёж, приходный ордер); наличие и единственность основания проверяет сервер.
// paymentNumber и paymentDate из основания не заполняются — их нужно передать.
// AdvancePaymentVat (0–100, ставка из справочника) и PaymentPurpose сервер учитывает только при
// основании-платеже, при отгрузке и возврате игнорирует; диапазон и ставку проверяет сервер
// (ошибки 3002, 3008, 3005).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/factureout#3-sozdat-schet-fakturu
type CreateRequest struct {
	Organization      model.Reference                        `json:"organization"`
	Name              model.Field[string]                    `json:"name,omitzero"`
	Description       model.Field[string]                    `json:"description,omitzero"`
	Code              model.Field[string]                    `json:"code,omitzero"`
	ExternalCode      model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID            model.Field[string]                    `json:"syncId,omitzero"`
	Moment            model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable        model.Field[bool]                      `json:"applicable,omitzero"`
	Shared            model.Field[bool]                      `json:"shared,omitzero"`
	Agent             model.Field[model.Reference]           `json:"agent,omitzero"`
	Contract          model.Field[model.Reference]           `json:"contract,omitzero"`
	Owner             model.Field[model.Reference]           `json:"owner,omitzero"`
	Group             model.Field[model.Reference]           `json:"group,omitzero"`
	State             model.Field[model.Reference]           `json:"state,omitzero"`
	Attributes        model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate              model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	StateContractID   model.Field[string]                    `json:"stateContractId,omitzero"`
	AdvancePaymentVat model.Field[int]                       `json:"advancePaymentVat,omitzero"`
	PaymentPurpose    model.Field[string]                    `json:"paymentPurpose,omitzero"`
	Demands           model.Field[[]model.Reference]         `json:"demands,omitzero"`
	Payments          model.Field[[]model.Reference]         `json:"payments,omitzero"`
	Returns           model.Field[[]model.Reference]         `json:"returns,omitzero"`
	Consignee         model.Field[model.Reference]           `json:"consignee,omitzero"`
	ConsigneeBranch   model.Field[model.Reference]           `json:"consigneeBranch,omitzero"`
	PaymentNumber     model.Field[string]                    `json:"paymentNumber,omitzero"`
	PaymentDate       model.Field[model.DateTime]            `json:"paymentDate,omitzero"`
}

// UpdateRequest. sum, vatSum, files и прочие readonly-поля не пишутся; syncId после заполнения не
// меняется. AdvancePaymentVat и PaymentPurpose при Update сервер обрабатывает, только если они
// переданы. Пустой массив связей уходит как [], что с ним делает сервер, источник не описывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/factureout#3-izmenit-schet-fakturu-vydannyj
type UpdateRequest struct {
	Organization      model.Field[model.Reference]           `json:"organization,omitzero"`
	Name              model.Field[string]                    `json:"name,omitzero"`
	Description       model.Field[string]                    `json:"description,omitzero"`
	Code              model.Field[string]                    `json:"code,omitzero"`
	ExternalCode      model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID            model.Field[string]                    `json:"syncId,omitzero"`
	Moment            model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable        model.Field[bool]                      `json:"applicable,omitzero"`
	Shared            model.Field[bool]                      `json:"shared,omitzero"`
	Agent             model.Field[model.Reference]           `json:"agent,omitzero"`
	Contract          model.Field[model.Reference]           `json:"contract,omitzero"`
	Owner             model.Field[model.Reference]           `json:"owner,omitzero"`
	Group             model.Field[model.Reference]           `json:"group,omitzero"`
	State             model.Field[model.Reference]           `json:"state,omitzero"`
	Attributes        model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate              model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	StateContractID   model.Field[string]                    `json:"stateContractId,omitzero"`
	AdvancePaymentVat model.Field[int]                       `json:"advancePaymentVat,omitzero"`
	PaymentPurpose    model.Field[string]                    `json:"paymentPurpose,omitzero"`
	Demands           model.Field[[]model.Reference]         `json:"demands,omitzero"`
	Payments          model.Field[[]model.Reference]         `json:"payments,omitzero"`
	Returns           model.Field[[]model.Reference]         `json:"returns,omitzero"`
	Consignee         model.Field[model.Reference]           `json:"consignee,omitzero"`
	ConsigneeBranch   model.Field[model.Reference]           `json:"consigneeBranch,omitzero"`
	PaymentNumber     model.Field[string]                    `json:"paymentNumber,omitzero"`
	PaymentDate       model.Field[model.DateTime]            `json:"paymentDate,omitzero"`
}
