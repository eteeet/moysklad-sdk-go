Status: needs-triage
Stage: DRAFT
Type: wayfinder:map

# SDK v1 — карта захваток и тикетов

## Destination

Закончить выбранный профиль SDK по [общей spec](spec.md), а не весь JSON API.
[Coverage исходного запроса](coverage.md), [каталог контрактов](contracts.md).
Код не менялся; planning не коммитится.

## Захватки

| № | Захватка | Tickets | Когда запускать |
| --- | --- | --- | --- |
| 1 | [Транспорт, query, лимиты и List](../sdk-v1-foundation/spec.md) | 13 | Начать T01; transport tests/источники можно готовить независимо |
| 2 | [Сотрудники, права, metadata и Audit](../sdk-v1-access/spec.md) | 8 | После T02/T06; Audit не ждёт товарные/розничные services |
| 3 | [Товары, услуги, комплекты, модификации и цены](../sdk-v1-catalog/spec.md) | 8 | После transport options/List; C01–C04 независимы |
| 4 | [Заказ покупателя, складские операции и возврат](../sdk-v1-documents/spec.md) | 5 | После T02/T06 и D01; D02–D05 независимы |
| 5 | [Розница и штатные бонусы](../sdk-v1-retail/spec.md) | 6 | Розничный сценарий E01→E02→E03/E04; references не требуют services |
| 6 | [Файлы и семейство отчётов Turnover](../sdk-v1-files-reports/spec.md) | 5 | F01→F02; F03→F04/F05 независимо документов |
| 7 | [Фасад, примеры и готовность v1](../sdk-v1-release/spec.md) | 4 | После основных capabilities; naming/consumer examples проектируются раньше |
| 8 | [Опциональное расширение: Delete, batch и async](../sdk-v1-optional/spec.md) | 3 | Только после явного включения в scope |

В основной v1 — 49 tickets в 7 захватках; optional — ещё 3 tickets в восьмой spec.
Для реализации «крупная захватка» означает один согласованный пакет работ,
а не один большой code commit. Коммит остаётся по ticket.

## Порядок и параллельность

1. Утвердить T01 и общий scope. Затем T02–T10; безопасно отделять query/parser и
   error/token от limiter/retry. T11/T12 и T13 завершают первый пакет.
2. После T02/T06 можно параллельно вести Access/Audit, Catalog и file/report work.
   Не ждать весь CRUD ради Audit. Общие transport/model файлы имеют одного владельца.
3. D01 фиксирует общую форму строк, затем новые документы. RetailStore/Shift
   дают удобный порядок интеграционных сценариев; это не технические blockers
   для Reference. BonusProgram/Transaction также независимы по services.
4. Facade assembly/examples/CI/status завершают v1. Optional X01–X03 добавляют
   зависимости только после решения о включении; их отсутствие не скрывается.

## Все tickets

### Транспорт, query, лимиты и List

| ID | Ticket | Blocked by | Status |
| --- | --- | --- | --- |
| T01 | [Утвердить общий interface options/response metadata](../sdk-v1-foundation/issues/01-public-interface.md) | нет | resolved |
| T02 | [Все request/response headers и интеграция options](../sdk-v1-foundation/issues/02-headers.md) | T01 | resolved |
| T03 | [Каталог кодов ошибок и status-aware helpers](../sdk-v1-foundation/issues/03-error-codes.md) | нет | resolved |
| T04 | [POST security/token](../sdk-v1-foundation/issues/04-security-token.md) | T01, T03 | resolved |
| T05 | [Filter grammar и order](../sdk-v1-foundation/issues/05-filter-order.md) | T01 | resolved |
| T06 | [Query options и получение одной страницы](../sdk-v1-foundation/issues/06-query-pages.md) | T01, T05 | resolved |
| T07 | [Явный walker/iterator страниц](../sdk-v1-foundation/issues/07-iterator.md) | T06 | resolved |
| T08 | [Weighted rate limiter и concurrency](../sdk-v1-foundation/issues/08-rate-concurrency.md) | T02 | resolved |
| T09 | [Ограниченные retries/backoff](../sdk-v1-foundation/issues/09-safe-retry.md) | T03, T08 | resolved |
| T10 | [Размеры запросов и endpoint limits](../sdk-v1-foundation/issues/10-size-limits.md) | T01, T03 | resolved |
| T11 | [List шести существующих справочников](../sdk-v1-foundation/issues/11-list-dictionaries.md) | T02, T06 | resolved |
| T12 | [List одиннадцати существующих документов](../sdk-v1-foundation/issues/12-list-documents.md) | T02, T06 | resolved |
| T13 | [SyncID Product и единый контракт syncId](../sdk-v1-foundation/issues/13-product-syncid.md) | нет | resolved |

### Сотрудники, права, metadata и Audit

| ID | Ticket | Blocked by | Status |
| --- | --- | --- | --- |
| A01 | [Employee: Get/List/Create/Update](../sdk-v1-access/issues/01-employee.md) | T02, T06 | resolved |
| A02 | [Role и права Employee](../sdk-v1-access/issues/02-roles-security.md) | A01 | resolved |
| A03 | [GET context/employee](../sdk-v1-access/issues/03-employee-context.md) | T02 | resolved |
| A04 | [GET context/application](../sdk-v1-access/issues/04-application-context.md) | T02 | resolved |
| A05 | [Общий каталог и metadata сущностей](../sdk-v1-access/issues/05-entity-metadata.md) | T02, T06 | resolved |
| A06 | [States: чтение, Create/Update](../sdk-v1-access/issues/06-states.md) | A05 | resolved |
| A07 | [Audit contexts и filters](../sdk-v1-access/issues/07-audit-contexts.md) | T02, T06 | resolved |
| A08 | [Audit events и lossless diff](../sdk-v1-access/issues/08-audit-events.md) | A07 | resolved |

### Товары, услуги, комплекты, модификации и цены

| ID | Ticket | Blocked by | Status |
| --- | --- | --- | --- |
| C01 | [Currency Get/List/Create/Update](../sdk-v1-catalog/issues/01-currency.md) | T02, T06 | resolved |
| C02 | [PriceType List/Get/GetDefault/ReplaceAll](../sdk-v1-catalog/issues/02-pricetype.md) | T02 | resolved |
| C03 | [TaxRate Get/List/Create/Update](../sdk-v1-catalog/issues/03-taxrate.md) | T02, T06 | resolved |
| C04 | [Service Get/List/Create/Update](../sdk-v1-catalog/issues/04-service.md) | T02, T06 | resolved |
| C05 | [Характеристики Variant: definitions и значения](../sdk-v1-catalog/issues/05-characteristics.md) | A05 | resolved |
| C06 | [Bundle и components](../sdk-v1-catalog/issues/06-bundle.md) | T02, T06 | resolved |
| C07 | [Variant Get/List/Create/Update](../sdk-v1-catalog/issues/07-variant.md) | T02, T06, C05 | resolved |
| C08 | [Assortment List/search и mixed rows](../sdk-v1-catalog/issues/08-assortment.md) | T02, T06, C04, C06, C07 | resolved |

### Заказ покупателя, складские операции и возврат

| ID | Ticket | Blocked by | Status |
| --- | --- | --- | --- |
| D01 | [Расширения общей формы позиций](../sdk-v1-documents/issues/01-position-contracts.md) | T06 | resolved |
| D02 | [CustomerOrder Get/List/Create/Update/GetPositions](../sdk-v1-documents/issues/02-customerorder.md) | T02, T06, D01 | resolved |
| D03 | [Enter Get/List/Create/Update/GetPositions](../sdk-v1-documents/issues/03-enter.md) | T02, T06, D01 | resolved |
| D04 | [Loss Get/List/Create/Update/GetPositions](../sdk-v1-documents/issues/04-loss.md) | T02, T06, D01 | resolved |
| D05 | [SalesReturn Get/List/Create/Update/GetPositions](../sdk-v1-documents/issues/05-salesreturn.md) | T02, T06, D01 | resolved |

### Розница и штатные бонусы

| ID | Ticket | Blocked by | Status |
| --- | --- | --- | --- |
| E01 | [RetailStore Get/List/Create/Update](../sdk-v1-retail/issues/01-retailstore.md) | T02, T06 | resolved |
| E02 | [RetailShift Get/List/Create/Update](../sdk-v1-retail/issues/02-retailshift.md) | T02, T06 | resolved |
| E03 | [RetailDemand Get/List/Create/Update/GetPositions](../sdk-v1-retail/issues/03-retaildemand.md) | T02, T06, D01 | resolved |
| E04 | [RetailSalesReturn Get/List/Create/Update/GetPositions](../sdk-v1-retail/issues/04-retailsalesreturn.md) | T02, T06, D01 | resolved |
| E05 | [BonusProgram Get/List/Create/Update](../sdk-v1-retail/issues/05-bonusprogram.md) | T02, T06 | resolved |
| E06 | [BonusTransaction Get/List/Create/Update](../sdk-v1-retail/issues/06-bonustransaction.md) | T02, T06 | resolved |

### Файлы и семейство отчётов Turnover

| ID | Ticket | Blocked by | Status |
| --- | --- | --- | --- |
| F01 | [Transport для binary download](../sdk-v1-files-reports/issues/01-download-transport.md) | T02, T08, T10 | resolved |
| F02 | [Files List/Upload/Download](../sdk-v1-files-reports/issues/02-files.md) | T02, T06, T10, F01 | resolved |
| F03 | [Turnover.All](../sdk-v1-files-reports/issues/03-turnover-all.md) | T02, T06 | resolved |
| F04 | [Turnover.ByStore](../sdk-v1-files-reports/issues/04-turnover-bystore.md) | F03 | resolved |
| F05 | [Turnover.ByOperations](../sdk-v1-files-reports/issues/05-turnover-byoperations.md) | F03 | resolved |

### Фасад, примеры и готовность v1

| ID | Ticket | Blocked by | Status |
| --- | --- | --- | --- |
| R01 | [Единая точка входа SDK](../sdk-v1-release/issues/01-facade.md) | T01, T02, T03, T04, T05, T06, T07, T08, T09, T10, T11, T12, T13, A01, A02, A03, A04, A05, A06, A07, A08, C01, C02, C03, C04, C05, C06, C07, C08, D01, D02, D03, D04, D05, E01, E02, E03, E04, E05, E06, F01, F02, F03, F04, F05 | needs-triage |
| R02 | [Примеры основных сценариев](../sdk-v1-release/issues/02-examples.md) | R01 | needs-triage |
| R03 | [CI и проверки выпуска](../sdk-v1-release/issues/03-ci-release.md) | R01 | needs-triage |
| R04 | [Приёмка и статус v1](../sdk-v1-release/issues/04-acceptance.md) | R02, R03 | needs-triage |

### Опциональное расширение: Delete, batch и async

| ID | Ticket | Blocked by | Status |
| --- | --- | --- | --- |
| X01 | [Delete выбранных ресурсов](../sdk-v1-optional/issues/01-delete.md) | T02, T03 | needs-info |
| X02 | [Batch create/update/delete](../sdk-v1-optional/issues/02-batch.md) | T02, T03, T08, T10 | needs-info |
| X03 | [Async submit/status/result/cancel](../sdk-v1-optional/issues/03-async.md) | T02, T03, T08, T09, F01 | needs-info |

## Decisions so far

- 2026-11-01, пользователь: T01 — per-call `opts ...RequestOption` + `WithResponseInfo(&info)`, не
  Result[T] ([ADR 0001](../../docs/adr/0001-request-options.md)); X01–X03 не входят в v1; работа
  захватками с остановкой после каждой; dirty worktree core-* — stash -u и удалить.
- Захватка 1 (T01–T13) resolved: итоги и BASE в Answer каждого ticket. Отклонение T13: SyncID
  Product без Update — по read-only маркеру источника.
- Ревью foundation 46e9d06, 35078b4 и db22bf1: замечания исправлены, последнее — 0 замечаний.
- Захватка 2 (A01–A08) resolved. JSON-контракт прав записан в model.Permissions (A02/A03);
  противоречия источника (events/info аудита, attributes массивом/объектом, id описаний полей,
  stateType по умолчанию) записаны в godoc и Answer.

- Пользователь запросил план/spec/tickets для перечисленного профиля v1.
  Это авторизация planning, не запуска реализации или commit planning.
- Источник закреплён на official revision 9a811cef; сводка исследований — research/.
- Рабочее предположение по вопросу scope: Get/List/Create/Update и явно названные
  специальные методы; Delete/batch/async пока вынесены в optional needs-info.
- Audit ListContexts/GetFilters и events имеют реальные отдельные маршруты;
  GET audit/{id} не придуман. Результаты events не получают выдуманную pagination.
- PriceType редактируется только полным ordered списком; Assortment read-only
  mixed list полезен даже при наличии четырёх пишущих entity packages.
- Две неоднозначности источника требуют contract freeze: Role.permissions
  Array vs Object и обязательность периода Turnover.ByOperations.

## Not yet specified / перед code

T01 фиксирует точный публичный interface response metadata/options; C08 — default
режим assortment; entity tickets утверждают поимённые Write matrix с каталогом полей.
Эти открытые решения видимы, tickets не помечены ready автоматически.

## Проверки планирования

Каждое требование сопоставлено ticket; относительные ссылки и Blocked by проверяются,
циклов зависимостей быть не должно. Выполнение planning не является PASS тестов SDK
или приёмкой ещё не реализованных возможностей.

Проверено при подготовке 2026-11-01: все относительные ссылки и dependency paths
существуют; номера tickets уникальны, DAG без циклов; исходные требования покрыты.
Изменений Go-кода и новых коммитов нет. Specs остаются DRAFT, main tickets needs-triage,
optional tickets needs-info. Runtime tests SDK для Markdown planning не запускались.
