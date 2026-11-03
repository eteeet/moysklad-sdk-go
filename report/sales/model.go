// Package sales читает показатели продаж и заказов: report/sales/plotseries и
// report/orders/plotseries. Оба графика возвращают одинаковую форму series (date, quantity, sum),
// поэтому у них общий тип ответа Plot; meta.type различает их (salesplotseries, ordersplotseries).
// Нужно право viewDashboard.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-sales-orders#2-pokazateli-prodazh-i-zakazov
package sales

import "github.com/eteeet/moysklad-sdk-go/model"

// Plot — ответ графика: точки за период с шагом Options.Interval. Series обязателен в ответе,
// пустой массив допустим.
type Plot struct {
	Context *model.ResponseContext `json:"context,omitzero"`
	Meta    model.Meta             `json:"meta"`
	Series  []Point                `json:"series"`
}

// Point — показатели одного интервала: Date — его начало, Quantity — количество документов
// (Int источника), Sum — сумма в копейках.
type Point struct {
	Date     model.DateTime `json:"date"`
	Quantity int            `json:"quantity"`
	Sum      float64        `json:"sum"`
}
