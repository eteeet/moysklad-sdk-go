// Package turnover читает отчёт «Обороты»: report/turnover/all, /bystore и /byoperations.
// Отчёт только для чтения; суммы себестоимости — копейки, количества — не деньги.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-turnover
package turnover

import (
	"encoding/json/jsontext"

	"github.com/eteeet/moysklad-sdk-go/model"
)

// Assortment — краткое представление товара или модификации в отчёте. Image источник описывает
// только как объект — JSON как пришёл.
type Assortment struct {
	Meta          model.Meta     `json:"meta"`
	Name          string         `json:"name"`
	Code          string         `json:"code"`
	Article       string         `json:"article"`
	Image         jsontext.Value `json:"image,omitzero"`
	ProductFolder *Named         `json:"productFolder,omitzero"`
	UOM           *Named         `json:"uom,omitzero"`
}

// Named — ссылка с именем: группа товара, единица измерения, склад, контрагент.
type Named struct {
	Meta model.Meta `json:"meta"`
	Name string     `json:"name"`
}

// Metrics — показатели: Sum — сумма себестоимости в копейках, Quantity — количество единиц.
type Metrics struct {
	Sum      float64 `json:"sum"`
	Quantity float64 `json:"quantity"`
}

// Row — строка отчёта «Обороты по товарам»: показатели на начало и конец периода, приход и расход.
type Row struct {
	Assortment    Assortment `json:"assortment"`
	OnPeriodStart Metrics    `json:"onPeriodStart"`
	OnPeriodEnd   Metrics    `json:"onPeriodEnd"`
	Income        Metrics    `json:"income"`
	Outcome       Metrics    `json:"outcome"`
}

// StoreRow — строка отчёта по товару с детализацией по складам. Таблица называет stockByStore
// объектом, ответы содержат массив — SDK читает массив.
type StoreRow struct {
	Assortment   Assortment      `json:"assortment"`
	StockByStore []StoreTurnover `json:"stockByStore"`
}

// StoreTurnover — показатели одного склада.
type StoreTurnover struct {
	Store         Named   `json:"store"`
	OnPeriodStart Metrics `json:"onPeriodStart"`
	OnPeriodEnd   Metrics `json:"onPeriodEnd"`
	Income        Metrics `json:"income"`
	Outcome       Metrics `json:"outcome"`
}

// OperationRow — строка отчёта по товару с детализацией по документам: одна операция с товаром
// на складе. Quantity — количество, Cost — себестоимость единицы в копейках, Sum — сумма
// себестоимостей в копейках.
type OperationRow struct {
	Assortment Assortment `json:"assortment"`
	Store      Named      `json:"store"`
	Operation  Operation  `json:"operation"`
	Quantity   float64    `json:"quantity"`
	Cost       float64    `json:"cost"`
	Sum        float64    `json:"sum"`
}

// Operation — документ оборота; тип открыт (Meta.Type: supply, demand, retaildemand …) и не
// требует пакета SDK для этого документа. Таблица называет moment объектом, ответы содержат
// дату — SDK читает DateTime.
type Operation struct {
	Meta        model.Meta     `json:"meta"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Moment      model.DateTime `json:"moment,omitzero"`
	Agent       *Named         `json:"agent,omitzero"`
}
