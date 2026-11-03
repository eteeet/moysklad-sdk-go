// Package stock читает отчёт «Остатки»: расширенный (report/stock/all), по складам
// (report/stock/bystore), краткие (report/stock/all/current, /bystore/current, /byslot/current) и
// по документу (report/stock/byoperation). Отчёт только для чтения; price, salePrice и cost —
// копейки, остаток, резерв, ожидание и доступно — количества, не деньги.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-stock
package stock

import (
	"github.com/eteeet/moysklad-sdk-go/images"
	"github.com/eteeet/moysklad-sdk-go/model"
)

// Named — объект отчёта с метаданными и именем: единица измерения.
type Named struct {
	Meta model.Meta `json:"meta"`
	Name string     `json:"name"`
}

// Folder — группа товара, модификации или партии; PathName — наименование родительской группы.
type Folder struct {
	Meta     model.Meta `json:"meta"`
	Name     string     `json:"name"`
	PathName string     `json:"pathName"`
}

// Row — строка расширенного отчёта: товар, модификация или партия (Meta.Type открыт). Price —
// себестоимость, SalePrice — цена продажи, обе в копейках. Stock — остаток, Reserve — резерв,
// InTransit — ожидание, Quantity — доступно. Таблица называет image типом Meta, а пример ответа
// содержит объект изображения (meta, title, filename, size, updated, miniature, tiny) — SDK
// читает объект. Таблица называет stockDays обязательным, пример его почти не содержит — 0.
type Row struct {
	Meta         model.Meta    `json:"meta"`
	Name         string        `json:"name"`
	Code         string        `json:"code"`
	Article      string        `json:"article"`
	ExternalCode string        `json:"externalCode"`
	Folder       *Folder       `json:"folder,omitzero"`
	Image        *images.Image `json:"image,omitzero"`
	UOM          Named         `json:"uom"`
	Stock        float64       `json:"stock"`
	InTransit    float64       `json:"inTransit"`
	Reserve      float64       `json:"reserve"`
	Quantity     float64       `json:"quantity"`
	Price        float64       `json:"price"`
	SalePrice    float64       `json:"salePrice"`
	StockDays    int           `json:"stockDays"`
}

// StoreRow — строка отчёта по складам: позиция (Meta) и остатки на каждом складе. Таблица называет
// stockByStore объектом, ответы содержат массив — SDK читает массив.
type StoreRow struct {
	Meta         model.Meta   `json:"meta"`
	StockByStore []StoreStock `json:"stockByStore"`
}

// StoreStock — остатки одного склада: Meta и Name — склад.
type StoreStock struct {
	Meta      model.Meta `json:"meta"`
	Name      string     `json:"name"`
	Stock     float64    `json:"stock"`
	InTransit float64    `json:"inTransit"`
	Reserve   float64    `json:"reserve"`
}

// Amounts — значение краткого отчёта. Ответ содержит одно поле, выбранное stockType (stock по
// умолчанию); остальные не заданы (IsZero). С withRecalculate пересчитываемая позиция приходит с
// null (IsNull). Значения — количества, не деньги.
type Amounts struct {
	Stock     model.Field[float64] `json:"stock,omitzero"`
	FreeStock model.Field[float64] `json:"freeStock,omitzero"`
	Quantity  model.Field[float64] `json:"quantity,omitzero"`
	Reserve   model.Field[float64] `json:"reserve,omitzero"`
	InTransit model.Field[float64] `json:"inTransit,omitzero"`
}

// CurrentRow — строка report/stock/all/current: UUID товара, модификации или партии и значение.
type CurrentRow struct {
	AssortmentID string `json:"assortmentId"`
	Amounts
}

// StoreCurrentRow — строка report/stock/bystore/current. StoreID null (резерв или ожидание
// заказа без склада) читается как "".
type StoreCurrentRow struct {
	AssortmentID string `json:"assortmentId"`
	StoreID      string `json:"storeId"`
	Amounts
}

// SlotRow — строка report/stock/byslot/current: UUID позиции, склада и ячейки и остаток; без meta.
type SlotRow struct {
	AssortmentID string  `json:"assortmentId"`
	StoreID      string  `json:"storeId"`
	SlotID       string  `json:"slotId"`
	Stock        float64 `json:"stock"`
}

// OperationRow — строка отчёта по документу: документ (Meta.Type открыт) и остатки его позиций.
type OperationRow struct {
	Meta      model.Meta          `json:"meta"`
	Positions []OperationPosition `json:"positions"`
}

// OperationPosition — остатки позиции документа; Meta и Name — товар, модификация или комплект.
// Cost — себестоимость в копейках; у комплекта InTransit, Reserve и Quantity всегда 0.
type OperationPosition struct {
	Meta      model.Meta `json:"meta"`
	Name      string     `json:"name"`
	Stock     float64    `json:"stock"`
	Cost      float64    `json:"cost"`
	InTransit float64    `json:"inTransit"`
	Reserve   float64    `json:"reserve"`
	Quantity  float64    `json:"quantity"`
}
