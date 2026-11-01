// Package turnover читает отчёт «Обороты»: report/turnover/all.
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

// Named — ссылка с именем: группа товара, единица измерения.
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
