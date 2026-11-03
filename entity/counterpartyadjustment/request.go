package counterpartyadjustment

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы organization и agent (контрагент или сотрудник). Sum (копейки) таблица
// атрибутов помечает «Только для чтения», но примеры создания и изменения его передают: SDK sum
// пишет — исключение из правила «пометка таблицы сильнее примеров» утверждено пользователем, без
// него сумму корректировки не задать. Code, syncId, rate, state и project таблица не содержит.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/counterpartyadjustment#3-sozdat-korrektirovku-vzaimoraschetov
type CreateRequest struct {
	Organization model.Reference                `json:"organization"`
	Agent        model.Reference                `json:"agent"`
	Name         model.Field[string]            `json:"name,omitzero"`
	Description  model.Field[string]            `json:"description,omitzero"`
	ExternalCode model.Field[string]            `json:"externalCode,omitzero"`
	Moment       model.Field[model.DateTime]    `json:"moment,omitzero"`
	Applicable   model.Field[bool]              `json:"applicable,omitzero"`
	Shared       model.Field[bool]              `json:"shared,omitzero"`
	Owner        model.Field[model.Reference]   `json:"owner,omitzero"`
	Group        model.Field[model.Reference]   `json:"group,omitzero"`
	Attributes   model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	Sum          model.Field[float64]           `json:"sum,omitzero"`
}

// UpdateRequest. Текст раздела запрещает передавать поля «Только для чтения», но пример изменения
// передаёт sum — SDK его пишет по исключению из CreateRequest.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/counterpartyadjustment#3-izmenit-korrektirovku-vzaimoraschetov
type UpdateRequest struct {
	Organization model.Field[model.Reference]   `json:"organization,omitzero"`
	Agent        model.Field[model.Reference]   `json:"agent,omitzero"`
	Name         model.Field[string]            `json:"name,omitzero"`
	Description  model.Field[string]            `json:"description,omitzero"`
	ExternalCode model.Field[string]            `json:"externalCode,omitzero"`
	Moment       model.Field[model.DateTime]    `json:"moment,omitzero"`
	Applicable   model.Field[bool]              `json:"applicable,omitzero"`
	Shared       model.Field[bool]              `json:"shared,omitzero"`
	Owner        model.Field[model.Reference]   `json:"owner,omitzero"`
	Group        model.Field[model.Reference]   `json:"group,omitzero"`
	Attributes   model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	Sum          model.Field[float64]           `json:"sum,omitzero"`
}
