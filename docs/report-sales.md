# Показатели продаж и заказов (report/sales, report/orders)

Пакет `report/sales`, handle фасада — `c.Report().Sales()` (после подключения в `client.go`).
Источник: [md/reports/_report_sales_orders.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_sales_orders.md)
(revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`), раздел API
[Показатели продаж и заказов](https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-sales-orders#2-pokazateli-prodazh-i-zakazov).
Нужно право `viewDashboard`.

## API

| Метод `sales.Service` | Endpoint | Учитывается |
| --- | --- | --- |
| `Sales(ctx, Options)` | GET report/sales/plotseries | отгрузки, розничные продажи, полученные отчёты комиссионера |
| `Orders(ctx, Options)` | GET report/orders/plotseries | только заказы покупателей |

Оба возвращают `*sales.Plot`: таблица series общая для обоих разделов (date, quantity, sum), поэтому
один тип; `meta.type` различает графики (`salesplotseries`, `ordersplotseries`). Один GET, без
повторов. Async (I23) не входит.

## Параметры `Options`

| Поле | Query | Проверка SDK |
| --- | --- | --- |
| `MomentFrom model.DateTime` | momentFrom | обязателен (required в источнике): пустой — ошибка до HTTP |
| `MomentTo model.DateTime` | momentTo | обязателен |
| `Interval string` | interval | обязателен; `IntervalHour`, `IntervalDay`, `IntervalMonth`; иное значение SDK не отвергает |
| `Filter string` | filter | результат `query.Filter` |
| `Extra url.Values` | прочие | без повторов typed-полей, expand и async |

Filter: organization, store, project у обоих графиков; retailStore — только у Sales (тогда учитываются
только розничные продажи). Каждое поле допускает несколько значений, пустые сервер не принимает.
`transport.WithAcceptTimezone` смещает даты операций в зону клиента; ответная зона — в
`ResponseInfo.ContentTimezone`.

## Поля ответа

| JSON | Go type | Смысл |
| --- | --- | --- |
| context | `*model.ResponseContext` | контекст запроса |
| meta | `model.Meta` | href обязателен, иначе ошибка |
| series | `[]Point` | обязателен; `null` или отсутствие — ошибка, `[]` — успех |
| series[].date | `model.DateTime` | начало интервала |
| series[].quantity | `int` | количество документов (Int источника), не деньги |
| series[].sum | `float64` | сумма, копейки |

## Противоречия источника и решения

- Пример запроса продаж с фильтром разделяет условия `&` (`filter=retailStore=…&project=…`), что
  превращает их в отдельные query-параметры; общий синтаксис filter и пример заказов — `;`. SDK
  передаёт filter как собрал `query.Filter` — через `;`.
- `meta.href` первого примера заказов содержит `momentTo=2016-09-06 …` при запросе `2018-09-06 …` —
  опечатка примера, на SDK не влияет.
- Таблица называет quantity Int, все примеры — целые: `int`. Ответ `3.0` не прочитается; если живой API
  так отвечает, тип нужно сменить на `float64`.

## Тесты

`report/sales/*_test.go`: JSON обоими кодеками (`testdata/upstream/report_{orders,sales}_plotseries_list.json`
и примеры источника в `testdata/custom/sales/` с SHA-256 в README), даты, деньги 0/0.0/дробь/большое
значение, точный GET path/query/headers, plain/gzip, ошибки API, null/отсутствующий series, обязательные
параметры до HTTP (ResponseInfo сброшен, без вызовов), timezone-заголовки, одна отправка.
