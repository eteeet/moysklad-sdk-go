package taxrate

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: обязателен только rate («единственное поле, которое обязательно должно
// присутствовать»); указатель отличает ставку 0 от отсутствия. Таблица помечает и updated как
// «необходимо при создании» — это противоречит тексту, SDK его не передаёт.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/taxrate#3-sozdat-stavku-nds
type CreateRequest struct {
	Rate     *float64                     `json:"rate"`
	Comment  model.Field[string]          `json:"comment,omitzero"`
	Archived model.Field[bool]            `json:"archived,omitzero"`
	Shared   model.Field[bool]            `json:"shared,omitzero"`
	Owner    model.Field[model.Reference] `json:"owner,omitzero"`
	Group    model.Field[model.Reference] `json:"group,omitzero"`
}

// UpdateRequest: rate предустановленных ставок сервер менять не даёт.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/taxrate#3-izmenit-stavku-nds
type UpdateRequest struct {
	Rate     model.Field[float64]         `json:"rate,omitzero"`
	Comment  model.Field[string]          `json:"comment,omitzero"`
	Archived model.Field[bool]            `json:"archived,omitzero"`
	Shared   model.Field[bool]            `json:"shared,omitzero"`
	Owner    model.Field[model.Reference] `json:"owner,omitzero"`
	Group    model.Field[model.Reference] `json:"group,omitzero"`
}
