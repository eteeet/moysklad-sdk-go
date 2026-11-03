package aggregatepack

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: обязателен непустой barcodes. ChildrenList без бета-заголовка сервер не обрабатывает.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/handling-unit#3-sozdat-transportnuyu-upakovku
type CreateRequest struct {
	Moment       model.Field[model.DateTime]    `json:"moment,omitzero"`
	Barcodes     []model.Barcode                `json:"barcodes"`
	ChildrenList model.Field[[]model.Reference] `json:"childrenList,omitzero"`
	Shared       model.Field[bool]              `json:"shared,omitzero"`
	Owner        model.Field[model.Reference]   `json:"owner,omitzero"`
	Group        model.Field[model.Reference]   `json:"group,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/handling-unit#3-izmenit-transportnuyu-upakovku
type UpdateRequest struct {
	Moment       model.Field[model.DateTime]    `json:"moment,omitzero"`
	Barcodes     model.Field[[]model.Barcode]   `json:"barcodes,omitzero"`
	ChildrenList model.Field[[]model.Reference] `json:"childrenList,omitzero"`
	Shared       model.Field[bool]              `json:"shared,omitzero"`
	Owner        model.Field[model.Reference]   `json:"owner,omitzero"`
	Group        model.Field[model.Reference]   `json:"group,omitzero"`
}
