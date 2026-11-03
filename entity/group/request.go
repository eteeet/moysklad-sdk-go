package group

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: обязателен name; новый отдел встаёт в конец списка, Index = 0 — в начало.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/group#3-sozdat-otdel
type CreateRequest struct {
	Name  string           `json:"name"`
	Index model.Field[int] `json:"index,omitzero"`
}

// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/group#3-izmenit-otdel
type UpdateRequest struct {
	Name  model.Field[string] `json:"name,omitzero"`
	Index model.Field[int]    `json:"index,omitzero"`
}
