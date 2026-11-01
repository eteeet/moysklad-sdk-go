package bonustransaction

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы agent, bonusProgram и transactionType (EARNING или SPENDING, открытое
// значение). ExecutionDate источник описывает для операций начисления: откладывает обработку, дата
// в будущем требует тарифной опции; для списания и баланс баллов решает сервер. Противоречия
// источника: таблица атрибутов не помечает bonusProgram «Необходимо при создании», таблица метода —
// помечает (SDK требует); пример Create передаёт readonly categoryType и не передаёт bonusProgram и
// transactionType. Права токенов источник не описывает: допускает ли сервер запись токеном решения
// (Vendor), выясняется ответом — SDK передаёт отказ как APIError и не пробует другой токен.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonustransaction#3-sozdat-bonusnuyu-operaciyu
type CreateRequest struct {
	Agent           model.Reference              `json:"agent"`
	BonusProgram    model.Reference              `json:"bonusProgram"`
	TransactionType string                       `json:"transactionType"`
	Name            model.Field[string]          `json:"name,omitzero"`
	Description     model.Field[string]          `json:"description,omitzero"`
	Code            model.Field[string]          `json:"code,omitzero"`
	ExternalCode    model.Field[string]          `json:"externalCode,omitzero"`
	Moment          model.Field[model.DateTime]  `json:"moment,omitzero"`
	ExecutionDate   model.Field[model.DateTime]  `json:"executionDate,omitzero"`
	Applicable      model.Field[bool]            `json:"applicable,omitzero"`
	Shared          model.Field[bool]            `json:"shared,omitzero"`
	Organization    model.Field[model.Reference] `json:"organization,omitzero"`
	Owner           model.Field[model.Reference] `json:"owner,omitzero"`
	Group           model.Field[model.Reference] `json:"group,omitzero"`
	ParentDocument  model.Field[model.Reference] `json:"parentDocument,omitzero"`
	BonusValue      model.Field[int]             `json:"bonusValue,omitzero"`
}

// UpdateRequest. Тип операции менять нельзя — transactionType здесь нет; transactionStatus и
// categoryType только для чтения.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/bonustransaction#3-izmenit-bonusnuyu-operaciyu
type UpdateRequest struct {
	Agent          model.Field[model.Reference] `json:"agent,omitzero"`
	BonusProgram   model.Field[model.Reference] `json:"bonusProgram,omitzero"`
	Name           model.Field[string]          `json:"name,omitzero"`
	Description    model.Field[string]          `json:"description,omitzero"`
	Code           model.Field[string]          `json:"code,omitzero"`
	ExternalCode   model.Field[string]          `json:"externalCode,omitzero"`
	Moment         model.Field[model.DateTime]  `json:"moment,omitzero"`
	ExecutionDate  model.Field[model.DateTime]  `json:"executionDate,omitzero"`
	Applicable     model.Field[bool]            `json:"applicable,omitzero"`
	Shared         model.Field[bool]            `json:"shared,omitzero"`
	Organization   model.Field[model.Reference] `json:"organization,omitzero"`
	Owner          model.Field[model.Reference] `json:"owner,omitzero"`
	Group          model.Field[model.Reference] `json:"group,omitzero"`
	ParentDocument model.Field[model.Reference] `json:"parentDocument,omitzero"`
	BonusValue     model.Field[int]             `json:"bonusValue,omitzero"`
}
