package contract

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: обязательны name (номер), ownAgent (юрлицо) и agent (контрагент) — по списку
// «Обязательные для создания Договора поля»; в общей таблице у name маркера нет. Соответствие счетов сторонам
// и тип вознаграждения проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contract#3-sozdat-novyj-dogovor
type CreateRequest struct {
	Name                string                                 `json:"name"`
	Description         model.Field[string]                    `json:"description,omitzero"`
	Code                model.Field[string]                    `json:"code,omitzero"`
	ExternalCode        model.Field[string]                    `json:"externalCode,omitzero"`
	Archived            model.Field[bool]                      `json:"archived,omitzero"`
	Moment              model.Field[model.DateTime]            `json:"moment,omitzero"`
	Sum                 model.Field[float64]                   `json:"sum,omitzero"`
	ContractType        model.Field[string]                    `json:"contractType,omitzero"`
	RewardType          model.Field[string]                    `json:"rewardType,omitzero"`
	RewardPercent       model.Field[int]                       `json:"rewardPercent,omitzero"`
	OwnAgent            model.Reference                        `json:"ownAgent"`
	Agent               model.Reference                        `json:"agent"`
	OrganizationAccount model.Field[model.Reference]           `json:"organizationAccount,omitzero"`
	AgentAccount        model.Field[model.Reference]           `json:"agentAccount,omitzero"`
	Rate                model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	State               model.Field[model.Reference]           `json:"state,omitzero"`
	Attributes          model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Shared              model.Field[bool]                      `json:"shared,omitzero"`
	Owner               model.Field[model.Reference]           `json:"owner,omitzero"`
	Group               model.Field[model.Reference]           `json:"group,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contract#3-izmenit-dogovor
type UpdateRequest struct {
	Name                model.Field[string]                    `json:"name,omitzero"`
	Description         model.Field[string]                    `json:"description,omitzero"`
	Code                model.Field[string]                    `json:"code,omitzero"`
	ExternalCode        model.Field[string]                    `json:"externalCode,omitzero"`
	Archived            model.Field[bool]                      `json:"archived,omitzero"`
	Moment              model.Field[model.DateTime]            `json:"moment,omitzero"`
	Sum                 model.Field[float64]                   `json:"sum,omitzero"`
	ContractType        model.Field[string]                    `json:"contractType,omitzero"`
	RewardType          model.Field[string]                    `json:"rewardType,omitzero"`
	RewardPercent       model.Field[int]                       `json:"rewardPercent,omitzero"`
	OwnAgent            model.Field[model.Reference]           `json:"ownAgent,omitzero"`
	Agent               model.Field[model.Reference]           `json:"agent,omitzero"`
	OrganizationAccount model.Field[model.Reference]           `json:"organizationAccount,omitzero"`
	AgentAccount        model.Field[model.Reference]           `json:"agentAccount,omitzero"`
	Rate                model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	State               model.Field[model.Reference]           `json:"state,omitzero"`
	Attributes          model.Field[[]model.Attribute]         `json:"attributes,omitzero"`
	Shared              model.Field[bool]                      `json:"shared,omitzero"`
	Owner               model.Field[model.Reference]           `json:"owner,omitzero"`
	Group               model.Field[model.Reference]           `json:"group,omitzero"`
}
