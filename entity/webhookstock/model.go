package webhookstock

import "github.com/eteeet/moysklad-sdk-go/model"

// WebhookStock — вебхук на изменение остатков. StockType и ReportType — открытые строки:
// значения вне перечислений источника ([stock], [all, bystore]) не отвергаются.
// AuthorApplication — ссылка на решение, создавшее вебхук; nil — поле не пришло.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhookstock#2-vebhuk-na-izmenenie-ostatkov
type WebhookStock struct {
	Meta              model.Meta       `json:"meta"`
	ID                string           `json:"id"`
	AccountID         string           `json:"accountId"`
	StockType         string           `json:"stockType"`
	ReportType        string           `json:"reportType"`
	URL               string           `json:"url"`
	Enabled           bool             `json:"enabled"`
	AuthorApplication *model.Reference `json:"authorApplication,omitzero"`
}
