package webhookstock

import "github.com/eteeet/moysklad-sdk-go/model"

// CreateRequest: url, stockType и reportType «Необходимо при создании» — пустые отклоняются до
// HTTP. Значения не проверяются: stockType и reportType — открытые строки, url проверяет сервер,
// как и уникальность сочетания и ограничения тарифа.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhookstock#3-sozdat-vebhuk-na-izmenenie-ostatkov
type CreateRequest struct {
	URL        string            `json:"url"`
	StockType  string            `json:"stockType"`
	ReportType string            `json:"reportType"`
	Enabled    model.Field[bool] `json:"enabled,omitzero"`
}

// UpdateRequest: Enabled: model.Value(false) отключает вебхук («Отключить вебхук на изменение остатков»).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhookstock#3-izmenit-vebhuk-na-izmenenie-ostatkov
type UpdateRequest struct {
	URL        model.Field[string] `json:"url,omitzero"`
	StockType  model.Field[string] `json:"stockType,omitzero"`
	ReportType model.Field[string] `json:"reportType,omitzero"`
	Enabled    model.Field[bool]   `json:"enabled,omitzero"`
}
