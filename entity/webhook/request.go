package webhook

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: url, action и entityType «Необходимо при создании» — пустые отклоняются до
// HTTP. Значения не проверяются: action, entityType и diffType — открытые строки, доступность
// url проверяет сервер. Ограничения уникальности и тарифа также проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhook#3-sozdat-vebhuk
type CreateRequest struct {
	URL        string              `json:"url"`
	Action     string              `json:"action"`
	EntityType string              `json:"entityType"`
	DiffType   model.Field[string] `json:"diffType,omitzero"`
	Enabled    model.Field[bool]   `json:"enabled,omitzero"`
}

// UpdateRequest: Enabled: model.Value(false) отключает вебхук («Отключить вебхук»).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhook#3-izmenit-vebhuk
type UpdateRequest struct {
	URL        model.Field[string] `json:"url,omitzero"`
	Action     model.Field[string] `json:"action,omitzero"`
	EntityType model.Field[string] `json:"entityType,omitzero"`
	DiffType   model.Field[string] `json:"diffType,omitzero"`
	Enabled    model.Field[bool]   `json:"enabled,omitzero"`
}
