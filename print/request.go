package print

import "github.com/eteeet/moysklad-sdk-go/model"

// ExportRequest — печать одного документа по шаблону. Оба поля обязательны. Extension —
// xls, pdf, html или ods по источнику; строка открыта, допустимость проверяет сервер.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/print
type ExportRequest struct {
	Template  model.Reference `json:"template"`
	Extension string          `json:"extension"`
}

// ExportSetRequest — печать комплекта: непустой список шаблонов с числом копий. Комплект
// печатается в pdf, extension не передаётся; template вне списка запрещён — отдельный тип.
type ExportSetRequest struct {
	Templates []TemplateCopies `json:"templates"`
}

// TemplateCopies — шаблон комплекта. Count обязателен; диапазон 1…10 проверяет сервер.
// Для части документов допустимы шаблоны связанных документов (см. источник).
type TemplateCopies struct {
	Template model.Reference `json:"template"`
	Count    int             `json:"count"`
}

// LabelRequest — печать этикеток и ценников товара, услуги, комплекта или модификации. Все поля
// обязательны; Count — число ценников/термоэтикеток, максимум 1000 проверяет сервер. Template —
// шаблон этикеток (entity/assortment/metadata/embeddedtemplate|customtemplate).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/labels
type LabelRequest struct {
	Organization model.Reference `json:"organization"`
	Count        int             `json:"count"`
	SalePrice    LabelSalePrice  `json:"salePrice"`
	Template     model.Reference `json:"template"`
}

// LabelSalePrice — цена продажи на этикетке: обязательный тип цены.
type LabelSalePrice struct {
	PriceType model.Reference `json:"priceType"`
}
