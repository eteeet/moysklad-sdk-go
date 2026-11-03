// Package dashboard читает «Показатели» за день, неделю и месяц: report/dashboard/day, /week и
// /month. Ответ — один объект без meta; нужно право viewDashboard.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/dashboard
package dashboard

// Dashboard — показатели за период; источник отмечает sales, orders и money обязательными в
// ответе, SDK отклоняет ответ без любого из них.
type Dashboard struct {
	Sales  Sales  `json:"sales"`
	Orders Orders `json:"orders"`
	Money  Money  `json:"money"`
}

// Sales — продажи за период. Count — количество, Amount — сумма в копейках (источник называет
// её «Прибыль»), MovementAmount — изменение Amount к прошлому такому же периоду в копейках:
// пример источника даёт amount 0 и movementAmount -80000.
type Sales struct {
	Count          int     `json:"count"`
	Amount         float64 `json:"amount"`
	MovementAmount float64 `json:"movementAmount"`
}

// Orders — заказы за период; поля как у Sales (таблица источника повторяет описания продаж).
type Orders struct {
	Count          int     `json:"count"`
	Amount         float64 `json:"amount"`
	MovementAmount float64 `json:"movementAmount"`
}

// Money — деньги за период. Income, Outcome и Balance — копейки. Movement и TodayMovement —
// «дельта за период» и «за сегодня»: единицу источник не называет, ticket I06 считает их
// процентами, поэтому SDK не объявляет их деньгами и отдаёт число как пришло.
type Money struct {
	Income        float64 `json:"income"`
	Outcome       float64 `json:"outcome"`
	Balance       float64 `json:"balance"`
	TodayMovement float64 `json:"todayMovement"`
	Movement      float64 `json:"movement"`
}
