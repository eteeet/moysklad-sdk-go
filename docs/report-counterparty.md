# Показатели контрагентов (report/counterparty)

Пакет `report/counterparty`; handle фасада планируется как `c.Report().Counterparty()` (пока
`counterparty.NewService(transportClient)`). Источник:
[_report_counterparty.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_counterparty.md)
(revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`).

| Метод `counterparty.Service` | Endpoint | Ответ |
| --- | --- | --- |
| `List(ctx, ListOptions, opts...)` | GET report/counterparty | `*model.Collection[Row]` — одна страница |
| `Selected(ctx, []model.Reference, opts...)` | POST report/counterparty | `*model.Collection[Row]` |
| `Get(ctx, id, opts...)` | GET report/counterparty/{id} | `*Row` |

Нужны право `viewCompanyCRM` и опция CRM тарифа; `profit` — ещё `viewProductCostAndProfit`
(всё проверяет сервер). Каждый метод — одна отправка.

- `Selected` — выборочные показатели: POST с телом
  `{"counterparties":[{"counterparty":{"meta":…}}, …]}`. Это чтение, не запись: ничего не создаёт и
  не изменяет. Transport не повторяет POST даже при `Config.Retry` — повтор решает вызывающий.
  До запроса проверяется: список не пуст, у каждой ссылки непустой `meta.href`; тип и существование
  контрагентов проверяет сервер. Ссылка передаётся как есть (href, metadataHref, type, mediaType).
- `Get` — id должен быть UUID, иначе ошибка до запроса.

## Параметры `ListOptions`

| Поле | Query | Ограничения |
| --- | --- | --- |
| `Limit` | limit | 0 — 1000; 1–1000 |
| `Offset` | offset | ≥ 0 |
| `Filter` | filter | результат `query.Filter`; см. ниже |
| `Extra` | любые | без повторов limit/offset/filter, без expand и async |

Фильтруемые атрибуты (операторы проверяет сервер): averageReceipt, balance, demandsCount, demandsSum,
discountsSum, firstDemandDate, lastDemandDate, lastEventDate, lastEventText, profit, returnsCount,
returnsSum, updated, id, counterparty.name, counterparty.companyType, counterparty.description,
counterparty.email, counterparty.inn, counterparty.phone.

## Поля `Row` (все только для чтения, все обязательны в ответе)

| JSON | Go type | Смысл |
| --- | --- | --- |
| meta | `model.Meta` | метаданные отчёта по контрагенту |
| counterparty | `Counterparty` | контрагент, см. ниже |
| averageReceipt | `float64` | средний чек, копейки |
| balance | `float64` | баланс, копейки |
| bonusBalance | `float64` | баллы — не деньги |
| demandsCount | `int` | количество продаж |
| demandsSum | `float64` | сумма продаж, копейки |
| discountsSum | `float64` | сумма скидок, копейки |
| firstDemandDate | `model.DateTime` | дата первой продажи; null — нулевой DateTime |
| lastDemandDate | `model.DateTime` | дата последней продажи; null — нулевой |
| lastEventDate | `model.DateTime` | дата последнего события; null — нулевой |
| lastEventText | `string` | текст последнего события; null — "" |
| profit | `float64` | прибыль, копейки |
| returnsCount | `int` | количество возвратов |
| returnsSum | `float64` | сумма возвратов, копейки |
| updated | `model.DateTime` | момент последнего изменения контрагента |

### `Counterparty`

| JSON | Go type | Смысл |
| --- | --- | --- |
| meta | `model.Meta` | метаданные контрагента |
| id | `string` | UUID |
| name | `string` | наименование |
| externalCode | `string` | внешний код |
| companyType | `model.CompanyType` | открытый enum: неизвестное значение сохраняется |

## Асинхронно

«Отчет Показатели контрагентов» входит в список async источника, кроме выборочных показателей. [_async.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/_async.md) (revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`).

| Метод | Endpoint |
| --- | --- |
| `ListAsync(ctx, ListOptions)` | GET report/counterparty?async=true |
| `GetAsync(ctx, id)` | GET report/counterparty/{id}?async=true |

Каждый метод — один GET с теми же параметрами, что у синхронного, плюс `async=true`; ответ 202 с
Content-Location (…/async/{id}) и Location (…/async/{id}/result) — `*async.Submitted` (ID,
StatusHref, ResultHref), другой статус или 202 без этих адресов — ошибка. Создание задачи —
побочный эффект: transport не повторяет `async=true` даже с `Config.Retry`; после сбоя сети задача
могла создаться — проверьте `c.Async().List`. Результат: `c.Async().Wait`, затем
`c.Async().Result(ctx, task.ID, w)` — JSON того же вида, что ответ синхронного метода (SDK его не
разбирает). `limit`/`offset` источник у async запрещает: ненулевые `Limit`/`Offset` и `limit`/`offset` в `Extra` — ошибка
до запроса (`ResponseInfo` сброшен); id `GetAsync` — UUID.

- `Selected` (POST, выборочные показатели) async не получает: источник исключает его явно.
- `GetAsync`: раздел «Показатели контрагента» входит в раздел отчёта, исключён только выборочный;
  отдельно async одного контрагента источник не подтверждает.

## Противоречия и неизвестное

- `inn`, `email`, `phone` контрагента не перечислены в таблице «Контрагент», но описаны как фильтры
  `counterparty.inn/email/phone`, `inn` есть в примерах ответов, а фикстура официального SDK
  (`testdata/upstream/report_counterparty.json`) содержит все три. По правилу полей SDK (поле не из
  таблицы не добавляется) они не читаются; решение — за пользователем.
- Таблица называет `bonusBalance`, даты и `lastEventText` обязательными, но примеры содержат `null`
  и строки без `bonusBalance`; null и отсутствие читаются как нулевое значение.
- Ссылка «Подробнее тут» у `counterparty` ведёт на раздел кода системы налогообложения товара —
  опечатка источника.
- У «Показатели контрагента» (Get) источник перечисляет limit и offset; к одному объекту они
  неприменимы, SDK их не передаёт.
- Параметры запроса и лимит числа ссылок у `Selected` источник не описывает.
- Выборочные показатели (POST) async не поддерживают — см. раздел «Асинхронно».

## Ошибки

Один префикс `moysklad:` и операция: `get counterparty report` (List), `select counterparty report`
(Selected), `get counterparty report {id}` (Get). Ошибки API — `*errors.APIError`; статус, отличный
от 200, ответ без `meta.href` или с `rows` null — ошибка.

Фикстуры: `testdata/custom/reportcounterparty` (примеры ответов с SHA-256, golden тела Selected).
