# Отчёт «По документам номенклатуры» (report/byoperations)

Пакет `report/byoperations`; handle фасада планируется как `c.Report().ByOperations()` (пока
`byoperations.NewService(transportClient)`). Источник:
[_report_by_operations.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_by_operations.md)
(revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`).

| Метод `byoperations.Service` | Endpoint | Ответ |
| --- | --- | --- |
| `Stock(ctx, Options, opts...)` | GET report/byoperations/stock | `*model.Collection[StockRow]` |
| `Reserve(ctx, Options, opts...)` | GET report/byoperations/reserve | `*model.Collection[ReserveRow]` |
| `InTransit(ctx, Options, opts...)` | GET report/byoperations/intransit | `*model.Collection[InTransitRow]` |

Каждый метод — один GET одной страницы, без автообхода по nextHref и без повторов, если caller
не включил `transport.Config.Retry`.

## Параметры (`Options`)

| Поле | Query | Правило |
| --- | --- | --- |
| `Assortment model.Reference` | `filter=assortment=<href>` | обязателен; meta.type `product`, `variant` или `consignment` и непустой href, иначе ошибка до HTTP |
| `Limit int` | `limit` | 0 — 1000; допустимо 1–1000 |
| `Offset int` | `offset` | ≥ 0; всегда передаётся |
| `Extra url.Values` | как есть | без `filter`, `limit`, `offset`, `expand`, `async`; имена `momentFrom`, `momentTo`, `groupBy`, `interval` тоже отклоняются (общий `query.ReportOptions`) — источник их для этого отчёта не описывает |

Источник описывает только фильтр `assortment=` с одной номенклатурой; другой фильтрации нет,
поэтому свободного поля `Filter` нет.

## Строки

Источник даёт трём отчётам разные поля (и разную обязательность `store`), поэтому общей строки
нет. Все поля только для чтения.

| JSON | Go type | Stock | Reserve | InTransit | Смысл |
| --- | --- | --- | --- | --- | --- |
| assortment | `*model.Reference` | обяз. | обяз. | обяз. | товар, модификация или партия |
| operation | `*model.Reference` | обяз. | обяз. | обяз. | документ; meta.type открыт |
| moment | `model.DateTime` | обяз. | обяз. | обяз. | дата документа |
| store | `*model.Reference` | обяз. | — | — | склад документа (у reserve/intransit не обязателен) |
| stock | `float64` | обяз. | | | остаток, количество |
| costPerUnit | `*float64` | опц. | | | себестоимость единицы, копейки |
| sumCost | `*float64` | опц. | | | сумма себестоимости, копейки |
| avgStockDays | `float64` | обяз. | | | дней на складе, не деньги |
| reserve | `float64` | | обяз. | | резерв, количество |
| inTransit | `float64` | | | обяз. | ожидание, количество |

- `costPerUnit` и `sumCost` сервер не присылает без права «Видеть себестоимость, цену закупки,
  прибыль товаров»: тогда поле `nil`, а не 0.
- Ссылки — только meta, как в источнике; expand вне SDK.

## Ошибки

Ошибки содержат один префикс `moysklad:` и операцию `get byoperations stock|reserve|intransit`;
ошибки API — `*errors.APIError`. Ответ без `meta.href` или с `rows` не массивом — ошибка.
Права «просмотр товаров» и «просмотр остатков» проверяет сервер.
