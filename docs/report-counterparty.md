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
- Async-режим (I23) не входит; выборочные показатели async не поддерживают.

## Ошибки

Один префикс `moysklad:` и операция: `get counterparty report` (List), `select counterparty report`
(Selected), `get counterparty report {id}` (Get). Ошибки API — `*errors.APIError`; статус, отличный
от 200, ответ без `meta.href` или с `rows` null — ошибка.

Фикстуры: `testdata/custom/reportcounterparty` (примеры ответов с SHA-256, golden тела Selected).
