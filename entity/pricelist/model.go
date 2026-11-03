package pricelist

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// PriceList — прайс-лист. Модель чтения без expand: columns — столбцы таблицы, priceType — тип
// цены (только чтение), строки — собственная форма Position (meta.type pricelistrow).
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/pricelist
type PriceList struct {
	Meta         model.Meta                         `json:"meta"`
	ID           string                             `json:"id"`
	AccountID    string                             `json:"accountId"`
	Name         string                             `json:"name"`
	Description  string                             `json:"description"`
	Code         string                             `json:"code"`
	ExternalCode string                             `json:"externalCode"`
	SyncID       string                             `json:"syncId"`
	Created      model.DateTime                     `json:"created,omitzero"`
	Updated      model.DateTime                     `json:"updated,omitzero"`
	Deleted      model.DateTime                     `json:"deleted,omitzero"`
	Moment       model.DateTime                     `json:"moment,omitzero"`
	Applicable   bool                               `json:"applicable"`
	Shared       bool                               `json:"shared"`
	Printed      bool                               `json:"printed"`
	Published    bool                               `json:"published"`
	Organization *model.Reference                   `json:"organization,omitzero"`
	Owner        *model.Reference                   `json:"owner,omitzero"`
	Group        *model.Reference                   `json:"group,omitzero"`
	State        *model.Reference                   `json:"state,omitzero"`
	Attributes   []model.Attribute                  `json:"attributes"`
	Columns      []Column                           `json:"columns"`
	PriceType    *model.PriceType                   `json:"priceType,omitzero"`
	Positions    *model.Collection[Position]        `json:"positions,omitzero"`
	Files        *model.Collection[model.Reference] `json:"files,omitzero"`
}

// Column — столбец прайс-листа: уникальное непустое name и процентная скидка (наценка) по
// умолчанию; 0 при чтении — скидка не задана или равна нулю.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/pricelist#3-stolbcy
type Column struct {
	Name               string `json:"name"`
	PercentageDiscount int    `json:"percentageDiscount,omitzero"`
}

// Cell — ячейка строки: column — name столбца, sum — числовое значение (таблица: Int; в примерах —
// цена в копейках, SDK держит float64 как остальные суммы). Оба поля необходимы при записи.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/pricelist#3-yachejki
type Cell struct {
	Column string  `json:"column"`
	Sum    float64 `json:"sum"`
}

// Position — строка прайс-листа (meta.type pricelistrow): товар, услуга, модификация или комплект
// и ячейки по столбцам. Это собственная форма, не positions.Position документов.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/pricelist#3-pozicii-prajs-lista
type Position struct {
	Meta       model.Meta       `json:"meta"`
	ID         string           `json:"id"`
	AccountID  string           `json:"accountId"`
	Assortment *model.Reference `json:"assortment,omitzero"`
	Cells      []Cell           `json:"cells"`
	Pack       *positions.Pack  `json:"pack,omitzero"`
}
