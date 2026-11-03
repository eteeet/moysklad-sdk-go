package pricelist

import (
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

// CreateRequest: необходимы columns (хотя бы один столбец с непустым name); после создания
// столбцы не меняются, поэтому их нет в UpdateRequest. priceType таблица помечает «Только для
// чтения», хотя пример Create его передаёт: SDK priceType не пишет.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/pricelist#3-sozdat-prajs-list
type CreateRequest struct {
	Columns      []ColumnRequest                `json:"columns"`
	Name         model.Field[string]            `json:"name,omitzero"`
	Description  model.Field[string]            `json:"description,omitzero"`
	Code         model.Field[string]            `json:"code,omitzero"`
	ExternalCode model.Field[string]            `json:"externalCode,omitzero"`
	SyncID       model.Field[string]            `json:"syncId,omitzero"`
	Moment       model.Field[model.DateTime]    `json:"moment,omitzero"`
	Applicable   model.Field[bool]              `json:"applicable,omitzero"`
	Shared       model.Field[bool]              `json:"shared,omitzero"`
	Organization model.Field[model.Reference]   `json:"organization,omitzero"`
	Owner        model.Field[model.Reference]   `json:"owner,omitzero"`
	Group        model.Field[model.Reference]   `json:"group,omitzero"`
	State        model.Field[model.Reference]   `json:"state,omitzero"`
	Attributes   model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	Positions    model.Field[PositionRows]      `json:"positions,omitzero"`
}

// UpdateRequest. Positions заменяют все строки: лишние удаляются, строки с ID меняются.
// syncId после заполнения не меняется.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/pricelist#3-izmenit-prajs-list
type UpdateRequest struct {
	Name         model.Field[string]            `json:"name,omitzero"`
	Description  model.Field[string]            `json:"description,omitzero"`
	Code         model.Field[string]            `json:"code,omitzero"`
	ExternalCode model.Field[string]            `json:"externalCode,omitzero"`
	SyncID       model.Field[string]            `json:"syncId,omitzero"`
	Moment       model.Field[model.DateTime]    `json:"moment,omitzero"`
	Applicable   model.Field[bool]              `json:"applicable,omitzero"`
	Shared       model.Field[bool]              `json:"shared,omitzero"`
	Organization model.Field[model.Reference]   `json:"organization,omitzero"`
	Owner        model.Field[model.Reference]   `json:"owner,omitzero"`
	Group        model.Field[model.Reference]   `json:"group,omitzero"`
	State        model.Field[model.Reference]   `json:"state,omitzero"`
	Attributes   model.Field[[]model.Attribute] `json:"attributes,omitzero"`
	Positions    model.Field[PositionRows]      `json:"positions,omitzero"`
}

// ColumnRequest — новый столбец: name необходим и уникален в прайс-листе (уникальность проверяет
// сервер); percentageDiscount — скидка (наценка) по умолчанию, ею заполняются ячейки без значения.
type ColumnRequest struct {
	Name string `json:"name"`
	// Единицу источник не называет: таблица — процент, примеры — 1000; SDK передаёт как есть.
	PercentageDiscount model.Field[int] `json:"percentageDiscount,omitzero"`
}

// PositionRows — positions в теле Create/Update: примеры источника передают объект {"rows": [...]},
// а не массив, как у других документов. Не больше 1000 строк.
type PositionRows struct {
	Rows []PositionRequest `json:"rows"`
}

// PositionRequest — строка прайс-листа. В CreatePositions assortment необходим (проверяет
// сервер); cells можно не передавать — тогда сервер заполнит ячейки столбцов со скидкой по
// умолчанию. ID сопоставляет существующую строку в Update прайс-листа; для новой строки его не
// задают. В UpdatePosition обязательных полей нет.
type PositionRequest struct {
	ID         model.Field[string]                `json:"id,omitzero"`
	Assortment model.Field[model.Reference]       `json:"assortment,omitzero"`
	Cells      model.Field[[]Cell]                `json:"cells,omitzero"`
	Pack       model.Field[positions.PackRequest] `json:"pack,omitzero"`
}
