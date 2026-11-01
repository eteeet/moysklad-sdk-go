Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-files-reports/issues/03-turnover-all.md

# F04: Turnover.ByStore

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

GET report/turnover/bystore с собственной формой детализации по складам; filter=product=<href>
либо variant=<href> обязателен, даты optional по методу. Не flatten или вычислять её из All на
клиенте.

## Специальные проверки

Exact path/query, product/variant/reference combinations, nested store metrics и timezone; не
переносить limit/offset/groupBy из All: в таблице ByStore они не подтверждены.

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

[report-turnover](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `f7e1507`, code commit `252bb73`.

- `turnover.Service.ByStore(ctx, DetailOptions, opts...)` → `*model.Collection[StoreRow]`, GET report/turnover/bystore.
- `DetailOptions{Assortment, MomentFrom, MomentTo, Filter}`: Assortment обязателен, его meta.type (product или
  variant) задаёт обязательный фильтр `product=`/`variant=`; повтор product/variant в Filter отклоняется
  до запроса. Период необязателен (умолчания сервера). limit, offset и groupBy не переносятся из All.
- `StoreRow{Assortment, StockByStore []StoreTurnover}`: показатели каждого склада как вернул сервер, без
  выведения из All. Противоречие: таблица называет stockByStore объектом, ответы — массив (SDK: массив).
- Тесты: upstream строка и полная страница из источника (`testdata/custom/turnover/bystore.json`), точные
  query для product и variant с периодом и store, зона расчёта, отказ без товара/модификации и при повторе.
- Review Standards: фикстура страницы, тест зоны, t.Run по codecs. Review Spec: проверка повтора
  product/variant в Filter, sparse-строка склада. Сброс ResponseInfo — в общем тесте волны.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
