package inventory

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы organization и store. Positions — не больше 1000 строк; дальше —
// отдельные операции позиций. enters и losses — связи документов: основная таблица их не
// помечает «Только для чтения», поэтому они пишутся по общему правилу привязки документов
// (ticket O09 называет enter/loss read-only — решение по источнику).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-sozdat-inventarizaciyu
type CreateRequest struct {
	Organization model.Reference                `json:"organization"`
	Store        model.Reference                `json:"store"`
	Name         model.Field[string]            `json:"name,omitzero"`
	Description  model.Field[string]            `json:"description,omitzero"`
	Code         model.Field[string]            `json:"code,omitzero"`
	ExternalCode model.Field[string]            `json:"externalCode,omitzero"`
	SyncID       model.Field[string]            `json:"syncId,omitzero"`
	Moment       model.Field[model.DateTime]    `json:"moment,omitzero"`
	Shared       model.Field[bool]              `json:"shared,omitzero"`
	Owner        model.Field[model.Reference]   `json:"owner,omitzero"`
	Group        model.Field[model.Reference]   `json:"group,omitzero"`
	State        model.Field[model.Reference]   `json:"state,omitzero"`
	Attributes   model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	Positions    model.Field[[]PositionRequest] `json:"positions,omitzero"`
	Enters       model.Field[[]model.Reference] `json:"enters,omitzero"`
	Losses       model.Field[[]model.Reference] `json:"losses,omitzero"`
}

// UpdateRequest. Positions заменяют все строки: лишние удаляются, строки с ID меняются. Пример
// Update источника передаёт positions объектом {"rows":[…]}; SDK, как у остальных документов,
// передаёт массив. syncId после заполнения не меняется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory#3-izmenit-inventarizaciyu
type UpdateRequest struct {
	Organization model.Field[model.Reference]   `json:"organization,omitzero"`
	Store        model.Field[model.Reference]   `json:"store,omitzero"`
	Name         model.Field[string]            `json:"name,omitzero"`
	Description  model.Field[string]            `json:"description,omitzero"`
	Code         model.Field[string]            `json:"code,omitzero"`
	ExternalCode model.Field[string]            `json:"externalCode,omitzero"`
	SyncID       model.Field[string]            `json:"syncId,omitzero"`
	Moment       model.Field[model.DateTime]    `json:"moment,omitzero"`
	Shared       model.Field[bool]              `json:"shared,omitzero"`
	Owner        model.Field[model.Reference]   `json:"owner,omitzero"`
	Group        model.Field[model.Reference]   `json:"group,omitzero"`
	State        model.Field[model.Reference]   `json:"state,omitzero"`
	Attributes   model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	Positions    model.Field[[]PositionRequest] `json:"positions,omitzero"`
	Enters       model.Field[[]model.Reference] `json:"enters,omitzero"`
	Losses       model.Field[[]model.Reference] `json:"losses,omitzero"`
}
