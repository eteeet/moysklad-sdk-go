package processingstage

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: обязателен непустой name. Связь AllPerformers и Performers SDK не проверяет:
// по источнику только allPerformers=true очищает performers, только performers ставит
// allPerformers=false; итог при обоих полях определяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingstage#3-sozdat-etap
type CreateRequest struct {
	Name                 string                         `json:"name"`
	Description          model.Field[string]            `json:"description,omitzero"`
	ExternalCode         model.Field[string]            `json:"externalCode,omitzero"`
	Archived             model.Field[bool]              `json:"archived,omitzero"`
	Shared               model.Field[bool]              `json:"shared,omitzero"`
	AllPerformers        model.Field[bool]              `json:"allPerformers,omitzero"`
	DistributionRequired model.Field[bool]              `json:"distributionRequired,omitzero"`
	StandardHourCost     model.Field[float64]           `json:"standardHourCost,omitzero"`
	Performers           model.Field[[]model.Reference] `json:"performers,omitzero"`
	Owner                model.Field[model.Reference]   `json:"owner,omitzero"`
	Group                model.Field[model.Reference]   `json:"group,omitzero"`
}

// UpdateRequest: Performers заменяет весь список исполнителей.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingstage#3-izmenit-etap
type UpdateRequest struct {
	Name                 model.Field[string]            `json:"name,omitzero"`
	Description          model.Field[string]            `json:"description,omitzero"`
	ExternalCode         model.Field[string]            `json:"externalCode,omitzero"`
	Archived             model.Field[bool]              `json:"archived,omitzero"`
	Shared               model.Field[bool]              `json:"shared,omitzero"`
	AllPerformers        model.Field[bool]              `json:"allPerformers,omitzero"`
	DistributionRequired model.Field[bool]              `json:"distributionRequired,omitzero"`
	StandardHourCost     model.Field[float64]           `json:"standardHourCost,omitzero"`
	Performers           model.Field[[]model.Reference] `json:"performers,omitzero"`
	Owner                model.Field[model.Reference]   `json:"owner,omitzero"`
	Group                model.Field[model.Reference]   `json:"group,omitzero"`
}
