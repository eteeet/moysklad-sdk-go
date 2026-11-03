# Отчёт «Прибыльность» (report/profit)

Пакет `report/profit`, handle фасада — `c.Report().Profit()` (после подключения в `client.go`).
Источник: [md/reports/_report_pnl.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_pnl.md)
(revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`), раздел API
[Отчёт Прибыльность](https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-pnl#2-otchet-pribylnost).
Нужно право на просмотр отчёта «Прибыльность».

## API

| Метод `profit.Service` | Endpoint | Строка |
| --- | --- | --- |
| `ByProduct(ctx, Options)` | GET report/profit/byproduct | `AssortmentRow` |
| `ByVariant(ctx, Options)` | GET report/profit/byvariant | `AssortmentRow` |
| `ByEmployee(ctx, Options)` | GET report/profit/byemployee | `EmployeeRow` |
| `ByCounterparty(ctx, Options)` | GET report/profit/bycounterparty | `CounterpartyRow` |
| `BySalesChannel(ctx, Options)` | GET report/profit/bysaleschannel | `SalesChannelRow` |

Каждый метод — один GET одной страницы (`*model.Collection[T]`), без автообхода `nextHref`.
Отчёт только для чтения: запись, expand и async (I23) не входят.

## Параметры `Options`

| Поле | Query | Проверка SDK |
| --- | --- | --- |
| `Limit int` | limit | 0 → 1000; вне 0..1000 — ошибка до HTTP |
| `Offset int` | offset | отрицательный — ошибка до HTTP |
| `MomentFrom model.DateTime` | momentFrom | необязателен; московское время |
| `MomentTo model.DateTime` | momentTo | необязателен |
| `Filter string` | filter | результат `query.Filter` |
| `Extra url.Values` | прочие | без limit, offset, momentFrom, momentTo, filter, groupBy, interval, expand, async |

Период по умолчанию решает сервер: без границ — последний месяц, только MomentTo — с начала года,
только MomentFrom — по текущий день. `transport.WithAcceptTimezone` задаёт зону расчёта, ответная — в
`ResponseInfo.ContentTimezone`.

Поля filter (все разрезы): product (`=`/`!=`, несколько значений), productFolder (`=`/`!=`, несколько),
withSubFolders, agentTag, entityType (`demand`, `retaildemand`), counterparty, organization, store,
project, retailStore (кроме bysaleschannel — его таблица поле не называет), supplier, salesChannel
(можно повторять). Остальные поля — одно значение. product вместе с productFolder сервер не принимает;
SDK это не проверяет (бизнес-ограничение остаётся API).

## Поля ответа

Все поля только читаются. Деньги — `float64` в копейках.

`AssortmentRow` (byproduct, byvariant — таблицы полей совпадают):

| JSON | Go type | Смысл |
| --- | --- | --- |
| assortment | `Assortment` | meta, name, code, article, uom `*Named`, image `jsontext.Value` |
| sellQuantity, returnQuantity | `float64` | количество, не деньги |
| sellPrice, returnPrice | `float64` | средняя цена, копейки |
| sellCost, returnCost | `float64` | себестоимость единицы, копейки |
| sellSum, returnSum, sellCostSum, returnCostSum | `float64` | суммы, копейки |
| profit | `float64` | прибыль, копейки |
| margin, salesMargin | `float64` | рентабельность товара и продаж, доля |

`EmployeeRow` (employee `Named`), `CounterpartyRow` (counterparty `Named`), `SalesChannelRow`
(salesChannel `SalesChannel`: meta, name, type — открытый enum) встраивают `Totals`:

| JSON | Go type | Смысл |
| --- | --- | --- |
| salesCount, returnCount | `float64` | количество продаж и возвратов, не деньги |
| salesAvgCheck, returnAvgCheck | `float64` | средний чек, копейки |
| sellSum, sellCostSum, returnSum, returnCostSum | `float64` | суммы, копейки |
| profit | `float64` | прибыль, копейки |
| margin, salesMargin | `float64` | доля |

## Асинхронно

«Отчет Прибыльность» входит в список async источника целиком. [_async.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/_async.md) (revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`).

| Метод | Endpoint |
| --- | --- |
| `ByProductAsync(ctx, Options)` | GET report/profit/byproduct?async=true |
| `ByVariantAsync(ctx, Options)` | GET report/profit/byvariant?async=true |
| `ByEmployeeAsync(ctx, Options)` | GET report/profit/byemployee?async=true |
| `ByCounterpartyAsync(ctx, Options)` | GET report/profit/bycounterparty?async=true |
| `BySalesChannelAsync(ctx, Options)` | GET report/profit/bysaleschannel?async=true |

Каждый метод — один GET с теми же параметрами, что у синхронного, плюс `async=true`; ответ 202 с
Content-Location (…/async/{id}) и Location (…/async/{id}/result) — `*async.Submitted` (ID,
StatusHref, ResultHref), другой статус или 202 без этих адресов — ошибка. Создание задачи —
побочный эффект: transport не повторяет `async=true` даже с `Config.Retry`; после сбоя сети задача
могла создаться — проверьте `c.Async().List`. Результат: `c.Async().Wait`, затем
`c.Async().Result(ctx, task.ID, w)` — JSON того же вида, что ответ синхронного метода (SDK его не
разбирает). `limit`/`offset` источник у async запрещает: ненулевые `Limit`/`Offset` и `limit`/`offset` в `Extra` — ошибка
до запроса (`ResponseInfo` сброшен).

## Противоречия источника и решения

- Путь сотрудников: запрос раздела — `report/profit/byemployee`, `meta.href` примера ответа —
  `report/profit/byEmployee`. SDK отправляет `byemployee` (как запрос раздела, upstream-фикстура
  `report_profit_by_employee_list.json` и остальные разрезы в нижнем регистре).
- salesCount и returnCount — Int в таблице, но примеры передают `1.0` и `0.0`; Go-int их не прочитал бы,
  поэтому `float64`.
- Таблица assortment по модификациям не называет uom, пример ответа его содержит — `UOM` читается
  необязательным полем у обоих разрезов.
- `meta.type` примера по покупателям — `salesbyCounterparty`, upstream-фикстуры — `salesbycounterparty`;
  SDK не проверяет meta.type.
- Пометку «в копейках» таблицы ставят только себестоимостям; цены, суммы, прибыль и средние чеки SDK
  тоже считает копейками (общий контракт API, примеры согласуются).

## Тесты

`report/profit/*_test.go`: JSON обоими кодеками (`testdata/upstream/report_profit_*.json` и примеры
источника в `testdata/custom/profit/` с SHA-256 в README), деньги 0/0.0/дробь/большое значение, точный
GET path/query/headers, plain/gzip, ошибки API, локальные ошибки (ResponseInfo сброшен, без HTTP), одна
отправка.
