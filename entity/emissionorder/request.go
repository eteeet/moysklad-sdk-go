package emissionorder

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: необходимы organization, trackingType и emissionType (пустые отклоняются до
// запроса); значения enum не проверяются. Positions — не больше 1000 строк; лимит источника
// «не больше 10 позиций» и совпадение trackingType строк и заказа проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/emissionorder#3-sozdat-zakaz-kodov-markirovki
type CreateRequest struct {
	Organization    model.Reference                `json:"organization"`
	TrackingType    string                         `json:"trackingType"`
	EmissionType    string                         `json:"emissionType"`
	Name            model.Field[string]            `json:"name,omitzero"`
	Description     model.Field[string]            `json:"description,omitzero"`
	ExternalCode    model.Field[string]            `json:"externalCode,omitzero"`
	Moment          model.Field[model.DateTime]    `json:"moment,omitzero"`
	Shared          model.Field[bool]              `json:"shared,omitzero"`
	Owner           model.Field[model.Reference]   `json:"owner,omitzero"`
	Group           model.Field[model.Reference]   `json:"group,omitzero"`
	State           model.Field[model.Reference]   `json:"state,omitzero"`
	Positions       model.Field[[]PositionRequest] `json:"positions,omitzero"`
	ProductionTasks model.Field[[]model.Reference] `json:"productionTasks,omitzero"`
}

// UpdateRequest. Positions заменяют все строки: строка с ID — существующая, без — новая.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/emissionorder#3-izmenit-zakaz-kodov-markirovki
type UpdateRequest struct {
	Organization    model.Field[model.Reference]   `json:"organization,omitzero"`
	TrackingType    model.Field[string]            `json:"trackingType,omitzero"`
	EmissionType    model.Field[string]            `json:"emissionType,omitzero"`
	Name            model.Field[string]            `json:"name,omitzero"`
	Description     model.Field[string]            `json:"description,omitzero"`
	ExternalCode    model.Field[string]            `json:"externalCode,omitzero"`
	Moment          model.Field[model.DateTime]    `json:"moment,omitzero"`
	Shared          model.Field[bool]              `json:"shared,omitzero"`
	Owner           model.Field[model.Reference]   `json:"owner,omitzero"`
	Group           model.Field[model.Reference]   `json:"group,omitzero"`
	State           model.Field[model.Reference]   `json:"state,omitzero"`
	Positions       model.Field[[]PositionRequest] `json:"positions,omitzero"`
	ProductionTasks model.Field[[]model.Reference] `json:"productionTasks,omitzero"`
}
