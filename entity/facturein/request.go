package facturein

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: organization обязателен по таблице атрибутов («Необходимо при создании»), хотя
// пример «только необходимые поля» его не передаёт — SDK следует таблице. Счёт-фактура создаётся
// только на основании: ровно один документ в Supplies (приёмка) или Payments (исходящий платёж);
// наличие и единственность основания проверяет сервер. incomingNumber и incomingDate из основания
// не заполняются — их нужно передать.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/facturein#3-sozdat-schet-fakturu
type CreateRequest struct {
	Organization   model.Reference                        `json:"organization"`
	Name           model.Field[string]                    `json:"name,omitzero"`
	Description    model.Field[string]                    `json:"description,omitzero"`
	Code           model.Field[string]                    `json:"code,omitzero"`
	ExternalCode   model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID         model.Field[string]                    `json:"syncId,omitzero"`
	Moment         model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable     model.Field[bool]                      `json:"applicable,omitzero"`
	Shared         model.Field[bool]                      `json:"shared,omitzero"`
	Agent          model.Field[model.Reference]           `json:"agent,omitzero"`
	Contract       model.Field[model.Reference]           `json:"contract,omitzero"`
	Owner          model.Field[model.Reference]           `json:"owner,omitzero"`
	Group          model.Field[model.Reference]           `json:"group,omitzero"`
	State          model.Field[model.Reference]           `json:"state,omitzero"`
	Attributes     model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate           model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	Supplies       model.Field[[]model.Reference]         `json:"supplies,omitzero"`
	Payments       model.Field[[]model.Reference]         `json:"payments,omitzero"`
	IncomingNumber model.Field[string]                    `json:"incomingNumber,omitzero"`
	IncomingDate   model.Field[model.DateTime]            `json:"incomingDate,omitzero"`
}

// UpdateRequest. sum, files и прочие readonly-поля не пишутся; syncId после заполнения не
// меняется. Пустой массив связей уходит как [], что с ним делает сервер, источник не описывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/facturein#3-izmenit-schet-fakturu-poluchennyj
type UpdateRequest struct {
	Organization   model.Field[model.Reference]           `json:"organization,omitzero"`
	Name           model.Field[string]                    `json:"name,omitzero"`
	Description    model.Field[string]                    `json:"description,omitzero"`
	Code           model.Field[string]                    `json:"code,omitzero"`
	ExternalCode   model.Field[string]                    `json:"externalCode,omitzero"`
	SyncID         model.Field[string]                    `json:"syncId,omitzero"`
	Moment         model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable     model.Field[bool]                      `json:"applicable,omitzero"`
	Shared         model.Field[bool]                      `json:"shared,omitzero"`
	Agent          model.Field[model.Reference]           `json:"agent,omitzero"`
	Contract       model.Field[model.Reference]           `json:"contract,omitzero"`
	Owner          model.Field[model.Reference]           `json:"owner,omitzero"`
	Group          model.Field[model.Reference]           `json:"group,omitzero"`
	State          model.Field[model.Reference]           `json:"state,omitzero"`
	Attributes     model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Rate           model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	Supplies       model.Field[[]model.Reference]         `json:"supplies,omitzero"`
	Payments       model.Field[[]model.Reference]         `json:"payments,omitzero"`
	IncomingNumber model.Field[string]                    `json:"incomingNumber,omitzero"`
	IncomingDate   model.Field[model.DateTime]            `json:"incomingDate,omitzero"`
}
