# Отчёт «Остатки» (report/stock)

Пакет `report/stock`. `stock.NewService(client)` принимает `*transport.Client`, возвращает
`(*stock.Service, error)` и отклоняет nil client. Отчёт только для чтения. Фасад: `c.Report().Stock()`.

Источник: [md/reports/_report_stock.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_stock.md)
(revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`).

## Методы

| Метод `stock.Service` | Endpoint | Ответ |
| --- | --- | --- |
| `All(ctx, AllOptions)` | GET report/stock/all | `*model.Collection[Row]` — одна страница |
| `ByStore(ctx, ByStoreOptions)` | GET report/stock/bystore | `*model.Collection[StoreRow]` — одна страница |
| `AllCurrent(ctx, CurrentOptions)` | GET report/stock/all/current | `[]CurrentRow` — JSON-массив |
| `ByStoreCurrent(ctx, CurrentOptions)` | GET report/stock/bystore/current | `[]StoreCurrentRow` — JSON-массив |
| `BySlotCurrent(ctx, SlotOptions)` | GET report/stock/byslot/current | `[]SlotRow` — JSON-массив |
| `ByOperation(ctx, operationID)` | GET report/stock/byoperation?operation.id= | `*model.Collection[OperationRow]` |

Каждый вызов — один GET без автообхода `nextHref`, повторов и предварительных запросов.
Ошибки — с одним префиксом `moysklad:` и операцией (`get stock all` …); `errors.As` распознаёт
`*sdkerrors.APIError`, `errors.Is` сохраняет ошибки context, авторизации и транспорта. Ошибка
параметров возвращается до HTTP, `ResponseInfo` при этом сброшен. Collection-ответ требует HTTP 200,
`meta.href` и массив `rows`; ответ кратких отчётов — HTTP 200 и JSON-массив (`[]` — успех без
строк, `null` и объект — ошибка).

**Лимит:** GET report/stock/all и report/stock/bystore весят 5 единиц в `transport.DefaultWeight`
независимо от профиля; тест `TestLimiterWeight` подтверждает это через публичный
`transport.LimiterConfig`. Краткие отчёты и byoperation весят как обычный запрос.

## Параметры

| Options | Поле Go | Query | Примечание |
| --- | --- | --- | --- |
| `AllOptions`, `ByStoreOptions` | `Limit`, `Offset` | limit, offset | `query.PageOptions`: 0 → 1000, 1–1000, offset ≥ 0; всегда передаются |
| | `GroupBy string` | groupBy | открыт: `GroupByProduct`, `GroupByVariant` (по умолчанию сервера), `GroupByConsignment` |
| | `Filter string` | filter | результат `query.Filter`; `stockMode` — открытая строка (`StockMode*`) |
| | `Order string` | order | результат `query.Order` |
| `AllOptions` | `IncludeRelated bool` | includeRelated=true | требует в Filter условие `product` или `variant` (`=` или `!=`), иначе ошибка до HTTP |
| `CurrentOptions` | `StockType string` | stockType | открыт: `StockTypeStock` (по умолчанию), `FreeStock`, `Quantity`, `Reserve`, `InTransit` |
| | `Include string` | include | `IncludeZeroLines` |
| | `ChangedSince model.DateTime` | changedSince | московское время; ограничение «не старше 24 часов и не в будущем» проверяет сервер |
| | `WithRecalculate bool` | withRecalculate=true | только ByStoreCurrent; в AllCurrent — ошибка до HTTP |
| `CurrentOptions`, `SlotOptions` | `AssortmentIDs`, `StoreIDs []string` | filter=assortmentId=a,b;storeId=c | каждое значение — UUID (форма, без версии) |
| `SlotOptions` | — | — | хотя бы один из AssortmentIDs/StoreIDs обязателен (источник), иначе ошибка до HTTP |
| все Options | `Extra url.Values` | как есть | без повторов typed-полей, limit/offset/order/includeRelated, expand и async |
| `ByOperation` | `operationID string` | operation.id | обязателен, UUID |

Фильтры All: archived, consignment, inTransitOnly, moment, product, productFolder, withSubFolders,
quantityMode, reserveOnly, search, soldByWeight, stockDaysFrom/To, stockMode, store, supplier,
variant, доп. поля и характеристики. Фильтры ByStore: consignment, moment, product, productFolder,
search, soldByWeight, stockMode, store, supplier, variant, доп. поля (только `=`, без пустых значений,
одно из product/variant/consignment/productFolder/search). Совместимость условий проверяет сервер.

## Поля ответа (запись не поддерживается)

| JSON | Go type | Примечание |
| --- | --- | --- |
| **Row** (all) | | |
| meta | `model.Meta` | товар, модификация или партия; `Meta.Type` открыт |
| name, code, article, externalCode | `string` | |
| folder | `*Folder` (meta, name, pathName) | nil — не пришло |
| image | `*images.Image` | nil — не пришло |
| uom | `Named` (meta, name) | |
| stock, inTransit, reserve, quantity | `float64` | количества |
| price, salePrice | `float64` | копейки |
| stockDays | `int` | дни на складе |
| **StoreRow** (bystore) | | |
| meta | `model.Meta` | позиция |
| stockByStore | `[]StoreStock` (meta, name, stock, inTransit, reserve) | массив складов |
| **CurrentRow** / **StoreCurrentRow** | | |
| assortmentId | `string` | UUID |
| storeId (только bystore) | `string` | `null` (резерв/ожидание заказа без склада) → `""` |
| stock, freeStock, quantity, reserve, inTransit | `model.Field[float64]` (встроенный `Amounts`) | приходит одно поле по stockType; незаданное — `IsZero`, `null` при withRecalculate — `IsNull` |
| **SlotRow** (byslot/current) | | |
| assortmentId, storeId, slotId | `string` | UUID, без meta и Reference |
| stock | `float64` | количество |
| **OperationRow** (byoperation) | | |
| meta | `model.Meta` | документ, тип открыт |
| positions | `[]OperationPosition` (meta, name, stock, cost, inTransit, reserve, quantity) | cost — копейки; у комплекта inTransit/reserve/quantity всегда 0 |

Неизвестные поля ответа игнорируются.

## Асинхронно

«Отчет Остатки» входит в список async источника целиком (раздел `## Отчет Остатки`, в нём расширенный, краткие, по складам и по документам). [_async.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/_async.md) (revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`).

| Метод | Endpoint | Синхронный аналог |
| --- | --- | --- |
| `AllAsync(ctx, AllOptions)` | GET report/stock/all?async=true | `All` |
| `ByStoreAsync(ctx, ByStoreOptions)` | GET report/stock/bystore?async=true | `ByStore` |
| `AllCurrentAsync(ctx, CurrentOptions)` | GET report/stock/all/current?async=true | `AllCurrent` |
| `ByStoreCurrentAsync(ctx, CurrentOptions)` | GET report/stock/bystore/current?async=true | `ByStoreCurrent` |
| `BySlotCurrentAsync(ctx, SlotOptions)` | GET report/stock/byslot/current?async=true | `BySlotCurrent` |
| `ByOperationAsync(ctx, operationID)` | GET report/stock/byoperation?operation.id=…&async=true | `ByOperation` |

Каждый метод — один GET с теми же параметрами, что у синхронного, плюс `async=true`; ответ 202 с
Content-Location (…/async/{id}) и Location (…/async/{id}/result) — `*async.Submitted` (ID,
StatusHref, ResultHref), другой статус или 202 без этих адресов — ошибка. Создание задачи —
побочный эффект: transport не повторяет `async=true` даже с `Config.Retry`; после сбоя сети задача
могла создаться — проверьте `c.Async().List`. Результат: `c.Async().Wait`, затем
`c.Async().Result(ctx, task.ID, w)` — JSON того же вида, что ответ синхронного метода (SDK его не
разбирает). `limit`/`offset` источник у async запрещает: ненулевые `Limit`/`Offset` у `AllAsync`/`ByStoreAsync` и `limit`/`offset` в `Extra` любого метода —
ошибка до запроса (`ResponseInfo` сброшен). Остальные проверки — как у синхронного метода.

- Решение по scope: источник перечисляет отчёты ссылками на раздел целиком и исключает только
  выборочные показатели контрагентов, поэтому async получили все GET-endpoints раздела, включая
  краткие и byoperation. Отдельно async у кратких отчётов и byoperation источник не подтверждает
  (пример — только bystore); если сервер откажет, это будет APIError.
- Вес лимита async-запроса all/bystore источник не описывает: `transport.DefaultWeight` считает его
  как синхронный (5 единиц).

## Противоречия и ограничения источника

- `image` расширенного отчёта: таблица — `Meta`, пример ответа и фикстура официального SDK — объект
  изображения (meta, title, filename, size, updated, miniature, tiny). SDK читает объект `images.Image`.
- `stockDays` таблица называет обязательным, но в примерах он есть лишь у партии — без поля 0.
- `stockByStore`: таблица — Object, ответ — массив. SDK читает массив.
- Таблица «Параметры» all/bystore не перечисляет filter и order, но разделы фильтрации и сортировки
  их описывают — оба поддержаны.
- Ответы кратких отчётов описаны только примерами (таблицы полей нет); набор полей значения взят из
  описания stockType.
- byslot/current: в документации значение — `stock`; фикстура официального SDK
  `testdata/upstream/stock_current_by_slot.json` содержит `quantity`. SDK читает `stock` (тикет I04 и
  документация); `stockType` и `include` для byslot источник не описывает — только через Extra.
- byoperation: limit/offset источник не описывает (ответ содержит `limit: 1000`), SDK их не передаёт.
  Позицию таблица называет «ссылкой на позицию», пример содержит meta товара/модификации.
- Список документов byoperation в тексте раздела и в списке метода расходится (счёт покупателю,
  возврат покупателя); SDK тип документа не проверяет — сервер определяет его по UUID.
- Дата начала веса 5 («с февраля 2026») без дня: `DefaultWeight` применяет вес всегда.

Примеры ответов источника — `testdata/custom/stock/*.json` с SHA-256 в README.
