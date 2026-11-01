Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md

# F03: Turnover.All

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

GET report/turnover/all: own row type assortment + onPeriodStart/onPeriodEnd/income/outcome,
metric sums/cost/quantities из source. momentFrom/To/groupBy/filter/page, timezone header в
report semantics.

## Специальные проверки

Отсутствующие/nullable quantities не превращаются в invented0 там где семантика требует
nullable; деньги float64, counts не деньги. Meta/context/plain scalar fields, no top-level
entity ID requirement, default groupBy product.

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

BASE `bf97b07`, code commit `f7e1507`.

- `turnover.Service.All(ctx, AllOptions, opts...)` → `*model.Collection[Row]`, GET report/turnover/all.
- `AllOptions`: Limit/Offset (через `query.PageOptions`, 0 → 1000), MomentFrom/MomentTo, GroupBy
  (product по умолчанию сервера, variant), Filter (результат query.Filter), Extra без повтора typed-полей.
- `Row`: assortment (`Assortment{Meta, Name, Code, Article, Image (JSON как есть), ProductFolder, UOM}`) и
  четыре `Metrics{Sum, Quantity}`: Sum — себестоимость в копейках, Quantity — количество. Все показатели
  источник помечает «обязательное при ответе» — обычные float64, не nullable.
- Период по умолчанию (последний месяц; только momentTo — с начала года; только momentFrom — по сегодня)
  решает сервер — в godoc. Зона расчёта — `transport.WithAcceptTimezone`, ответ — `ResponseInfo.ContentTimezone`;
  DateTime моделей от неё не зависит.
- Фикстуры: upstream строка и полная страница из источника (`testdata/custom/turnover/all.json`, SHA-256).
- Review Standards: `Extra`, общая проверка страницы через PageOptions, точные значения метрик в тесте,
  полный префикс в ошибках; package doc описывает только реализованное. Review Spec: полный пример
  страницы; оговорка про московское время ReportOptions и сдвиг в зону клиента — в docs волны.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
