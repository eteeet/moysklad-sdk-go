Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md

# T11: List шести существующих справочников

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Добавить typed List + общие options в Product/UOM/Country/ExpenseItem/Counterparty/Organization.
Покрыть supported filter/order/search/namedfilter/fields по endpoint, полный Collection
context/meta/rows. Текущие Get/Create/Update сохраняют прежний вызов.

## Специальные проверки

Табличные exact HTTP tests на каждый package, plain/gzip, empty array, invalid/null envelope,
query encoding, server error/metadata и один запрос. Не добавлять отсутствующие методы по
шаблону. Existing golden и read models не регрессируют.

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

BASE `7b2b9a9`, code commit `9acae45`.

- `List(ctx, query.ListOptions, opts...) (*model.Collection[T], error)` у Product, UOM, Country,
  ExpenseItem, Counterparty, Organization; общий `internal/list.Entity` (одна страница, без expand).
- Матрица по источнику: limit/offset/filter/order/search/namedfilter — общие для списков сущностей
  (_general.md), поля filter/order определяет API; fields — только Product `minimumStock`, у остальных
  отклоняется до запроса; archived (true;false — всё) — у Product, ExpenseItem, Counterparty,
  Organization, у UOM/Country поля нет.
- Тесты (`internal/listtest`, 17 пакетов): точные method/path/query/headers, plain/gzip, пустой массив,
  null/malformed/голый массив/204, APIError с ResponseInfo, ровно один запрос; неверные опции и
  неподдержанный fields — ноль запросов; локальные ошибки (nil/отменённый ctx, Auth, транспорт) через
  documenttest.LocalErrors. Строка — upstream-фикстура сущности.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS; старые
  Get/Create/Update тесты и golden не менялись.
- Review: обе оси одним прогоном на T11+T12 (раздельные агенты Standards и Spec); добавлены локальные
  ошибки и фильтры только по поддержанным полям.
