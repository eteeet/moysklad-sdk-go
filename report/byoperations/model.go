// Package byoperations читает отчёты «По документам номенклатуры»: report/byoperations/stock,
// /reserve и /intransit. Отчёты только для чтения; у каждого своя строка — источник даёт им
// разные поля. Остатки, резервы, ожидания и дни на складе — количества, не деньги.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-by-operations
package byoperations

import "github.com/eteeet/moysklad-sdk-go/model"

// StockRow — документ, формирующий остаток номенклатуры. CostPerUnit и SumCost — копейки; без
// права «Видеть себестоимость, цену закупки, прибыль товаров» сервер их не присылает — nil,
// а не 0. Stock — остаток, AvgStockDays — дни на складе.
type StockRow struct {
	Assortment   *model.Reference `json:"assortment,omitzero"`
	Operation    *model.Reference `json:"operation,omitzero"`
	Moment       model.DateTime   `json:"moment,omitzero"`
	Store        *model.Reference `json:"store,omitzero"`
	Stock        float64          `json:"stock"`
	CostPerUnit  *float64         `json:"costPerUnit,omitzero"`
	SumCost      *float64         `json:"sumCost,omitzero"`
	AvgStockDays float64          `json:"avgStockDays"`
}

// ReserveRow — документ, формирующий резерв. Store источник не отмечает обязательным.
type ReserveRow struct {
	Assortment *model.Reference `json:"assortment,omitzero"`
	Operation  *model.Reference `json:"operation,omitzero"`
	Moment     model.DateTime   `json:"moment,omitzero"`
	Store      *model.Reference `json:"store,omitzero"`
	Reserve    float64          `json:"reserve"`
}

// InTransitRow — документ, формирующий ожидание. Store источник не отмечает обязательным.
type InTransitRow struct {
	Assortment *model.Reference `json:"assortment,omitzero"`
	Operation  *model.Reference `json:"operation,omitzero"`
	Moment     model.DateTime   `json:"moment,omitzero"`
	Store      *model.Reference `json:"store,omitzero"`
	InTransit  float64          `json:"inTransit"`
}
