Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/01-public-interface.md; ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/03-error-codes.md; ../../sdk-v1-foundation/issues/04-security-token.md; ../../sdk-v1-foundation/issues/05-filter-order.md; ../../sdk-v1-foundation/issues/06-query-pages.md; ../../sdk-v1-foundation/issues/07-iterator.md; ../../sdk-v1-foundation/issues/08-rate-concurrency.md; ../../sdk-v1-foundation/issues/09-safe-retry.md; ../../sdk-v1-foundation/issues/10-size-limits.md; ../../sdk-v1-foundation/issues/11-list-dictionaries.md; ../../sdk-v1-foundation/issues/12-list-documents.md; ../../sdk-v1-foundation/issues/13-product-syncid.md; ../../sdk-v1-access/issues/01-employee.md; ../../sdk-v1-access/issues/02-roles-security.md; ../../sdk-v1-access/issues/03-employee-context.md; ../../sdk-v1-access/issues/04-application-context.md; ../../sdk-v1-access/issues/05-entity-metadata.md; ../../sdk-v1-access/issues/06-states.md; ../../sdk-v1-access/issues/07-audit-contexts.md; ../../sdk-v1-access/issues/08-audit-events.md; ../../sdk-v1-catalog/issues/01-currency.md; ../../sdk-v1-catalog/issues/02-pricetype.md; ../../sdk-v1-catalog/issues/03-taxrate.md; ../../sdk-v1-catalog/issues/04-service.md; ../../sdk-v1-catalog/issues/05-characteristics.md; ../../sdk-v1-catalog/issues/06-bundle.md; ../../sdk-v1-catalog/issues/07-variant.md; ../../sdk-v1-catalog/issues/08-assortment.md; ../../sdk-v1-documents/issues/01-position-contracts.md; ../../sdk-v1-documents/issues/02-customerorder.md; ../../sdk-v1-documents/issues/03-enter.md; ../../sdk-v1-documents/issues/04-loss.md; ../../sdk-v1-documents/issues/05-salesreturn.md; ../../sdk-v1-retail/issues/01-retailstore.md; ../../sdk-v1-retail/issues/02-retailshift.md; ../../sdk-v1-retail/issues/03-retaildemand.md; ../../sdk-v1-retail/issues/04-retailsalesreturn.md; ../../sdk-v1-retail/issues/05-bonusprogram.md; ../../sdk-v1-retail/issues/06-bonustransaction.md; ../../sdk-v1-files-reports/issues/01-download-transport.md; ../../sdk-v1-files-reports/issues/02-files.md; ../../sdk-v1-files-reports/issues/03-turnover-all.md; ../../sdk-v1-files-reports/issues/04-turnover-bystore.md; ../../sdk-v1-files-reports/issues/05-turnover-byoperations.md

# R01: Единая точка входа SDK

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Root moysklad.NewClient(Config) создаёт typed handles над одним transport/policy.
Рекомендованные пути: client.Product, client.CustomerOrder, client.Context.Employee,
client.Context.Application, client.Audit, client.Metadata, client.States, client.Files,
client.Reports.Turnover, client.Security. Точный naming принять до реализации; конструктор не
делает HTTP и не получает токен.

## Специальные проверки

Каждый существующий и новый service доступен, все используют один HTTPClient/auth/limiter.
Старые standalone constructors работают; no global singleton/auth credentials leakage.
Method/API compatibility examples и zero hidden network at startup.

## Acceptance criteria

- [ ] Scope/публичный interface и field/method matrix утверждены; зависимости resolved
  и их код присутствует в рабочем дереве. Readonly, required response/Create, null,
  immutable и выбранные TODO разделены по источникам; противоречия не замолчаны.
- [ ] Реализован весь What to build, без неоговорённых методов и hidden requests.
- [ ] Специальные проверки этого ticket реализованы через публичный interface.
- [ ] JSON обоими codecs и независимые golden; upstream/provenance не изменены.
  Для transport/iterator/CI tickets вместо entity JSON применяются их собственные
  публичные seams и контрактные tests, а не фиктивные модели.
- [ ] Реальные transport/httptest tests покрывают успех и ошибки; context и
  metadata сохраняются, caller-controlled retry/walker не подменяет один вызов.
- [ ] go test -race, go vet ./... ./testdata, gofmt и diff-check PASS; независимые
  Standards/Spec review, исправления и финальный BASE/результат записаны в Answer.
- [ ] README/guide/status отражают реальный scope, ограничения и неизвестные факты.
  Planning без команды не коммитится; один code commit по ticket.

## Источники

[транспортный research](../../sdk-v1/research/transport.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `5b9723c`. `moysklad.NewClient(Config)` (`Config = transport.Config`) создаёт только transport.Client:
HTTP нет, токен не выпускается, вид токена не определяется. Группы повторяют путь API (решение
пользователя, по образцу go-moysklad):

- `c.Entity()` — `EntityGroup`: 36 пакетов `entity/*` (`Product()`, `CustomerOrder()`, … `UOM()`,
  `Service()` → `*service.Client`, имя типа унаследовано от пакета), `Metadata()` со статусами,
  `Files()`. Отдельного States нет.
- `c.Report()` — `ReportGroup`: `Turnover()`.
- `c.Context()` — `contexts.Service` с методами `Employee(ctx)`/`Application(ctx)`.
- `c.Audit()`, `c.Security()`, `c.Transport()` — общий transport для Download, HeaderValues и
  standalone `NewService`, которые продолжают работать.

Handle не хранит собственного состояния и создаётся при каждом вызове: указатели разные, transport
один. Нулевое значение Client не используется: получение handle паникует сообщением NewService.
Глобального состояния нет.

Тесты через публичный interface и httptest: каждый каталог `entity/*` имеет handle; все 44 handle
держат один transport.Client; NewClient и получение handles не делают запросов; по одной группе
(Context, Audit, Report.Turnover, Entity.Files, Entity.Metadata) запрос уходит с Authorization
клиента; facade и standalone на `c.Transport()` и два клиента с разными токенами не смешиваются;
токен не виден в `%v/%+v/%#v` клиента, групп, transport и handle; ошибка NewClient. ExampleNewClient
компилируется (без Output, без сети) и показывает группы и прежний standalone путь.

Review Standards: `handle` → `must`, комментарии о нулевом значении и разных указателях,
`Entity`/`Report` → `EntityGroup`/`ReportGroup`, поведенческий тест групп, standalone в примере.
Godoc у методов групп оставлен: это публичный фасад. Review Spec: Answer/статус, пример
совместимости, проверка утечки по transport и handle; README/status — в docs-коммите волны.
`Transport()` оставлен как мост к standalone сервисам. Проверки: go test -race ./..., go vet
./... ./testdata, gofmt, diff --check — PASS.

## Comments

- 2026-11-02, решение пользователя: фасад делается через группы, в том числе сущности
  (группа сущностей, группа отчётов, группа контекстов — по образцу go-moysklad
  `client.Entity().Product()`). Отдельного handle States нет: статусы — методы Metadata.
