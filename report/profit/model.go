// Package profit читает отчёт «Прибыльность»: report/profit/byproduct, /byvariant, /byemployee,
// /bycounterparty и /bysaleschannel. Отчёт только для чтения; цены, суммы, себестоимости, прибыль и
// средние чеки — копейки; количества, счётчики и рентабельность — не деньги.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-pnl#2-otchet-pribylnost
package profit

import (
	"encoding/json/jsontext"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// Assortment — краткое представление товара, услуги, комплекта или модификации. Таблица по
// модификациям не называет uom, но её пример его содержит — UOM читается в обоих отчётах. Image
// источник описывает только как объект — JSON как пришёл.
type Assortment struct {
	Meta    model.Meta     `json:"meta"`
	Name    string         `json:"name"`
	Code    string         `json:"code"`
	Article string         `json:"article"`
	UOM     *Named         `json:"uom,omitzero"`
	Image   jsontext.Value `json:"image,omitzero"`
}

// Named — ссылка с именем: сотрудник, покупатель, единица измерения.
type Named struct {
	Meta model.Meta `json:"meta"`
	Name string     `json:"name"`
}

// SalesChannel — краткое представление канала продаж. Type — открытый enum (MESSENGER,
// SOCIAL_NETWORK, MARKETPLACE, ECOMMERCE, CLASSIFIED_ADS, DIRECT_SALES, OTHER …): неизвестное
// значение не отвергается.
type SalesChannel struct {
	Meta model.Meta `json:"meta"`
	Name string     `json:"name"`
	Type string     `json:"type"`
}

// AssortmentRow — строка отчётов по товарам (ByProduct) и по модификациям (ByVariant): таблицы
// полей этих разделов совпадают. Деньги в копейках: SellPrice и ReturnPrice — средняя цена,
// SellCost и ReturnCost — себестоимость единицы, *Sum и *CostSum — суммы, Profit — прибыль.
// SellQuantity и ReturnQuantity — количества; Margin (рентабельность товара) и SalesMargin
// (рентабельность продаж) — доли, не деньги.
type AssortmentRow struct {
	Assortment     Assortment `json:"assortment"`
	SellQuantity   float64    `json:"sellQuantity"`
	SellPrice      float64    `json:"sellPrice"`
	SellCost       float64    `json:"sellCost"`
	SellSum        float64    `json:"sellSum"`
	SellCostSum    float64    `json:"sellCostSum"`
	ReturnQuantity float64    `json:"returnQuantity"`
	ReturnPrice    float64    `json:"returnPrice"`
	ReturnCost     float64    `json:"returnCost"`
	ReturnSum      float64    `json:"returnSum"`
	ReturnCostSum  float64    `json:"returnCostSum"`
	Profit         float64    `json:"profit"`
	Margin         float64    `json:"margin"`
	SalesMargin    float64    `json:"salesMargin"`
}

// Totals — показатели отчётов по сотрудникам, покупателям и каналам продаж. Деньги в копейках:
// SalesAvgCheck и ReturnAvgCheck — средние чеки, *Sum и *CostSum — суммы, Profit — прибыль.
// SalesCount и ReturnCount — количества продаж и возвратов: таблица называет их Int, а примеры
// источника передают 1.0 и 0.0, поэтому SDK читает float64. Margin и SalesMargin — доли.
type Totals struct {
	SalesCount     float64 `json:"salesCount"`
	SalesAvgCheck  float64 `json:"salesAvgCheck"`
	SellSum        float64 `json:"sellSum"`
	SellCostSum    float64 `json:"sellCostSum"`
	ReturnCount    float64 `json:"returnCount"`
	ReturnAvgCheck float64 `json:"returnAvgCheck"`
	ReturnSum      float64 `json:"returnSum"`
	ReturnCostSum  float64 `json:"returnCostSum"`
	Profit         float64 `json:"profit"`
	Margin         float64 `json:"margin"`
	SalesMargin    float64 `json:"salesMargin"`
}

// EmployeeRow — строка отчёта по сотрудникам.
type EmployeeRow struct {
	Employee Named `json:"employee"`
	Totals
}

// CounterpartyRow — строка отчёта по покупателям.
type CounterpartyRow struct {
	Counterparty Named `json:"counterparty"`
	Totals
}

// SalesChannelRow — строка отчёта по каналам продаж.
type SalesChannelRow struct {
	SalesChannel SalesChannel `json:"salesChannel"`
	Totals
}
