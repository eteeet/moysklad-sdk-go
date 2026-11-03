package webhook

import "github.com/eteeet/moysklad-sdk-go/model"

// Webhook — вебхук. Action, EntityType, DiffType и Method — открытые строки: значения вне
// перечислений источника не отвергаются. DiffType источник указывает только для action UPDATE;
// пустая строка — поле не пришло. AuthorApplication — ссылка на решение, создавшее вебхук;
// nil — поле не пришло.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhook#2-vebhuki
type Webhook struct {
	Meta              model.Meta       `json:"meta"`
	ID                string           `json:"id"`
	AccountID         string           `json:"accountId"`
	EntityType        string           `json:"entityType"`
	URL               string           `json:"url"`
	Method            string           `json:"method"`
	Enabled           bool             `json:"enabled"`
	Action            string           `json:"action"`
	DiffType          string           `json:"diffType"`
	AuthorApplication *model.Reference `json:"authorApplication,omitzero"`
}
