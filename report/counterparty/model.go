// Package counterparty читает отчёт «Показатели контрагентов»: список report/counterparty,
// выборочные показатели (POST со ссылками на контрагентов — чтение, не запись) и показатели
// одного контрагента. Нужны право viewCompanyCRM и опция CRM тарифа; profit — ещё
// viewProductCostAndProfit (проверяет сервер). Суммы — копейки, счётчики и баллы — не деньги.
// https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-counterparty
package counterparty

import "github.com/eteeet/moysklad-sdk-go/model"

// Row — показатели одного контрагента. Деньги в копейках: AverageReceipt, Balance, DemandsSum,
// DiscountsSum, Profit, ReturnsSum. BonusBalance — баллы, DemandsCount и ReturnsCount —
// количества. Даты без события приходят null — нулевой DateTime; null LastEventText — "".
type Row struct {
	Meta            model.Meta     `json:"meta"`
	Counterparty    Counterparty   `json:"counterparty"`
	AverageReceipt  float64        `json:"averageReceipt"`
	Balance         float64        `json:"balance"`
	BonusBalance    float64        `json:"bonusBalance"`
	DemandsCount    int            `json:"demandsCount"`
	DemandsSum      float64        `json:"demandsSum"`
	DiscountsSum    float64        `json:"discountsSum"`
	FirstDemandDate model.DateTime `json:"firstDemandDate,omitzero"`
	LastDemandDate  model.DateTime `json:"lastDemandDate,omitzero"`
	LastEventDate   model.DateTime `json:"lastEventDate,omitzero"`
	LastEventText   string         `json:"lastEventText"`
	Profit          float64        `json:"profit"`
	ReturnsCount    int            `json:"returnsCount"`
	ReturnsSum      float64        `json:"returnsSum"`
	Updated         model.DateTime `json:"updated,omitzero"`
}

// Counterparty — краткое представление контрагента в отчёте; CompanyType открыт. inn, email и phone
// в таблице «Контрагент» источника нет (есть в фильтрах и примерах) — по правилу полей SDK их не
// читает; решение о них — за пользователем.
type Counterparty struct {
	Meta         model.Meta        `json:"meta"`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	ExternalCode string            `json:"externalCode"`
	CompanyType  model.CompanyType `json:"companyType"`
}
