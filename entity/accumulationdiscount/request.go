package accumulationdiscount

import (
	"github.com/eteeet/moysklad-sdk-go/entity/discount"
	"github.com/eteeet/moysklad-sdk-go/model"
)

// CreateRequest: name, active, allProducts и allAgents обязательны и передаются всегда, false — тоже.
// Таблица называет коллекцию групп productfolders, примеры — productFolders: SDK следует примерам.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-sozdat-nakopitelnuyu-skidku
type CreateRequest struct {
	Name           string                         `json:"name"`
	Active         bool                           `json:"active"`
	AllProducts    bool                           `json:"allProducts"`
	AllAgents      bool                           `json:"allAgents"`
	AgentTags      model.Field[[]string]          `json:"agentTags,omitzero"`
	Assortment     model.Field[[]model.Reference] `json:"assortment,omitzero"`
	ProductFolders model.Field[[]model.Reference] `json:"productFolders,omitzero"`
	Levels         model.Field[[]discount.Level]  `json:"levels,omitzero"`
}

// UpdateRequest: массивы передаются целиком; пустой массив или null удаляет значение.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-izmenit-nakopitelnuyu-skidku
type UpdateRequest struct {
	Name           model.Field[string]            `json:"name,omitzero"`
	Active         model.Field[bool]              `json:"active,omitzero"`
	AllProducts    model.Field[bool]              `json:"allProducts,omitzero"`
	AllAgents      model.Field[bool]              `json:"allAgents,omitzero"`
	AgentTags      model.Field[[]string]          `json:"agentTags,omitzero"`
	Assortment     model.Field[[]model.Reference] `json:"assortment,omitzero"`
	ProductFolders model.Field[[]model.Reference] `json:"productFolders,omitzero"`
	Levels         model.Field[[]discount.Level]  `json:"levels,omitzero"`
}
