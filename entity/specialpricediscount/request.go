package specialpricediscount

import (
	"github.com/eteeet/moysklad-sdk-go/entity/discount"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// CreateRequest: name, active, allProducts, allAgents и usePriceType обязательны и передаются всегда.
// specialPrice — по таблице (priceType, value); пример создания передаёт meta типа цены вместо них.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-sozdat-specialnuyu-cenu
type CreateRequest struct {
	Name           string                                    `json:"name"`
	Active         bool                                      `json:"active"`
	AllProducts    bool                                      `json:"allProducts"`
	AllAgents      bool                                      `json:"allAgents"`
	UsePriceType   bool                                      `json:"usePriceType"`
	AgentTags      model.Field[[]string]                     `json:"agentTags,omitzero"`
	Assortment     model.Field[[]model.Reference]            `json:"assortment,omitzero"`
	ProductFolders model.Field[[]model.Reference]            `json:"productFolders,omitzero"`
	Discount       model.Field[float64]                      `json:"discount,omitzero"`
	SpecialPrice   model.Field[discount.SpecialPriceRequest] `json:"specialPrice,omitzero"`
}

// UpdateRequest: массивы передаются целиком; пустой массив или null удаляет значение.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/discount#3-izmenit-specialnuyu-cenu
type UpdateRequest struct {
	Name           model.Field[string]                       `json:"name,omitzero"`
	Active         model.Field[bool]                         `json:"active,omitzero"`
	AllProducts    model.Field[bool]                         `json:"allProducts,omitzero"`
	AllAgents      model.Field[bool]                         `json:"allAgents,omitzero"`
	UsePriceType   model.Field[bool]                         `json:"usePriceType,omitzero"`
	AgentTags      model.Field[[]string]                     `json:"agentTags,omitzero"`
	Assortment     model.Field[[]model.Reference]            `json:"assortment,omitzero"`
	ProductFolders model.Field[[]model.Reference]            `json:"productFolders,omitzero"`
	Discount       model.Field[float64]                      `json:"discount,omitzero"`
	SpecialPrice   model.Field[discount.SpecialPriceRequest] `json:"specialPrice,omitzero"`
}
