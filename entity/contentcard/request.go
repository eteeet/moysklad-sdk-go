package contentcard

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: name, description, cardContentName, assortment, salePlatform и непустой
// salesChannels необходимы при создании. Assortment после создания не меняется — его нет в UpdateRequest.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contentcard#3-sozdat-kartochku-kontenta
type CreateRequest struct {
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	CardContentName string            `json:"cardContentName"`
	Assortment      model.Reference   `json:"assortment"`
	SalePlatform    model.Reference   `json:"salePlatform"`
	SalesChannels   []model.Reference `json:"salesChannels"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/contentcard#3-izmenit-kartochku-kontenta
type UpdateRequest struct {
	Name            model.Field[string]            `json:"name,omitzero"`
	Description     model.Field[string]            `json:"description,omitzero"`
	CardContentName model.Field[string]            `json:"cardContentName,omitzero"`
	SalePlatform    model.Field[model.Reference]   `json:"salePlatform,omitzero"`
	SalesChannels   model.Field[[]model.Reference] `json:"salesChannels,omitzero"`
}
