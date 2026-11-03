package commissionreportout

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// CreateRequest: необходимы organization, agent, contract (договор комиссии с этим agent),
// commissionPeriodStart, commissionPeriodEnd и name — name помечен «Необходимо при создании» в
// таблице атрибутов, хотя список обязательных полей раздела «Создать» его не называет.
// Счета organizationAccount/agentAccount нужны, если у стороны несколько счетов (проверяет сервер).
// Таблица помечает attributes, shared, positions и syncId «Только для чтения». Пример Create передаёт
// первые три: SDK их пишет — исключение из правила «пометка таблицы сильнее примеров» утверждено
// пользователем; syncId не пишется. Positions — не больше 1000 строк; дальше — CreatePositions.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportout#3-sozdat-vydannyj-otchet-komissionera
type CreateRequest struct {
	Organization          model.Reference                        `json:"organization"`
	Agent                 model.Reference                        `json:"agent"`
	Contract              model.Reference                        `json:"contract"`
	CommissionPeriodStart model.DateTime                         `json:"commissionPeriodStart"`
	CommissionPeriodEnd   model.DateTime                         `json:"commissionPeriodEnd"`
	Name                  string                                 `json:"name"`
	Description           model.Field[string]                    `json:"description,omitzero"`
	Code                  model.Field[string]                    `json:"code,omitzero"`
	ExternalCode          model.Field[string]                    `json:"externalCode,omitzero"`
	Moment                model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable            model.Field[bool]                      `json:"applicable,omitzero"`
	OrganizationAccount   model.Field[model.Reference]           `json:"organizationAccount,omitzero"`
	AgentAccount          model.Field[model.Reference]           `json:"agentAccount,omitzero"`
	Owner                 model.Field[model.Reference]           `json:"owner,omitzero"`
	Group                 model.Field[model.Reference]           `json:"group,omitzero"`
	Project               model.Field[model.Reference]           `json:"project,omitzero"`
	SalesChannel          model.Field[model.Reference]           `json:"salesChannel,omitzero"`
	State                 model.Field[model.Reference]           `json:"state,omitzero"`
	Rate                  model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	VATEnabled            model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded           model.Field[bool]                      `json:"vatIncluded,omitzero"`
	RewardType            model.Field[string]                    `json:"rewardType,omitzero"`
	RewardPercent         model.Field[int]                       `json:"rewardPercent,omitzero"`
	Shared                model.Field[bool]                      `json:"shared,omitzero"`
	Payments              model.Field[[]model.Reference]         `json:"payments,omitzero"`
	Attributes            model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Positions             model.Field[[]PositionRequest]         `json:"positions,omitzero"`
}

// UpdateRequest. Positions и attributes пишутся по тому же исключению, что в CreateRequest (их
// передаёт пример Update); positions заменяют все строки: лишние удаляются, строки с ID меняются.
// Shared пример Update не передаёт — в UpdateRequest его нет.
// При смене organization или agent источник требует передать и organizationAccount/agentAccount.
// Примеры источника передают sum и vatSum, но таблица помечает их readonly — SDK их не пишет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportout#3-izmenit-vydannyj-otchet-komissionera
type UpdateRequest struct {
	Organization          model.Field[model.Reference]           `json:"organization,omitzero"`
	Agent                 model.Field[model.Reference]           `json:"agent,omitzero"`
	Contract              model.Field[model.Reference]           `json:"contract,omitzero"`
	CommissionPeriodStart model.Field[model.DateTime]            `json:"commissionPeriodStart,omitzero"`
	CommissionPeriodEnd   model.Field[model.DateTime]            `json:"commissionPeriodEnd,omitzero"`
	Name                  model.Field[string]                    `json:"name,omitzero"`
	Description           model.Field[string]                    `json:"description,omitzero"`
	Code                  model.Field[string]                    `json:"code,omitzero"`
	ExternalCode          model.Field[string]                    `json:"externalCode,omitzero"`
	Moment                model.Field[model.DateTime]            `json:"moment,omitzero"`
	Applicable            model.Field[bool]                      `json:"applicable,omitzero"`
	OrganizationAccount   model.Field[model.Reference]           `json:"organizationAccount,omitzero"`
	AgentAccount          model.Field[model.Reference]           `json:"agentAccount,omitzero"`
	Owner                 model.Field[model.Reference]           `json:"owner,omitzero"`
	Group                 model.Field[model.Reference]           `json:"group,omitzero"`
	Project               model.Field[model.Reference]           `json:"project,omitzero"`
	SalesChannel          model.Field[model.Reference]           `json:"salesChannel,omitzero"`
	State                 model.Field[model.Reference]           `json:"state,omitzero"`
	Rate                  model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	VATEnabled            model.Field[bool]                      `json:"vatEnabled,omitzero"`
	VATIncluded           model.Field[bool]                      `json:"vatIncluded,omitzero"`
	RewardType            model.Field[string]                    `json:"rewardType,omitzero"`
	RewardPercent         model.Field[int]                       `json:"rewardPercent,omitzero"`
	Payments              model.Field[[]model.Reference]         `json:"payments,omitzero"`
	Attributes            model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Positions             model.Field[[]PositionRequest]         `json:"positions,omitzero"`
}

// PositionRequest — строка в Create/Update документа, CreatePositions и UpdatePosition. Для новой
// строки сервер требует assortment и положительный quantity; в UpdatePosition обязательных полей
// нет. ID сопоставляет существующую строку в Update документа; для новой строки его не задают.
type PositionRequest struct {
	ID         model.Field[string]                `json:"id,omitzero"`
	Assortment model.Field[model.Reference]       `json:"assortment,omitzero"`
	Quantity   model.Field[float64]               `json:"quantity,omitzero"`
	Price      model.Field[float64]               `json:"price,omitzero"`
	VAT        model.Field[int]                   `json:"vat,omitzero"`
	VATEnabled model.Field[bool]                  `json:"vatEnabled,omitzero"`
	Reward     model.Field[float64]               `json:"reward,omitzero"`
	Pack       model.Field[positions.PackRequest] `json:"pack,omitzero"`
}
