// Package money читает отчёт «Деньги»: движение денежных средств (report/money/plotseries) и
// остатки по кассам и счетам (report/money/byaccount). Нужны права viewDashboard и
// viewMoneyDashboard (проверяет сервер). Суммы — копейки.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-money
package money

import "github.com/eteeet/moysklad-sdk-go/model"

// PlotSeries — движение денежных средств за период: Credit — доход, Debit — расход за весь
// период, Series — показатели по интервалам; суммы в копейках. Context и Meta таблица атрибутов
// не перечисляет, но они есть в каждом примере ответа; ответ — не коллекция (нет rows и size).
type PlotSeries struct {
	Context *model.ResponseContext `json:"context,omitzero"`
	Meta    model.Meta             `json:"meta"`
	Credit  float64                `json:"credit"`
	Debit   float64                `json:"debit"`
	Series  []Point                `json:"series"`
}

// Point — показатели одного интервала: Date — его начало, Credit — доход, Debit — расход,
// Balance — доход минус расход; суммы в копейках.
type Point struct {
	Date    model.DateTime `json:"date"`
	Credit  float64        `json:"credit"`
	Debit   float64        `json:"debit"`
	Balance float64        `json:"balance"`
}

// AccountRow — текущий остаток одной кассы или счёта. Balance — в валюте счёта, в копейках
// (минимальных единицах); SDK валюту не пересчитывает. Account приходит только у счетов: у
// остатка кассы (одна на организацию) он nil.
type AccountRow struct {
	Account      *Named  `json:"account,omitzero"`
	Organization Named   `json:"organization"`
	Balance      float64 `json:"balance"`
}

// Named — ссылка с именем: организация (Name — наименование) или её счёт (Name — номер счёта,
// Meta.Type account).
type Named struct {
	Meta model.Meta `json:"meta"`
	Name string     `json:"name"`
}
