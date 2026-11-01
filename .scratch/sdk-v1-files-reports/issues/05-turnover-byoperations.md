Status: needs-triage
Stage: DRAFT
Type: task
Blocked by: ../../sdk-v1-files-reports/issues/03-turnover-all.md

# F05: Turnover.ByOperations

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

GET report/turnover/byoperations: own rows operation/assortment и numeric metrics; обязательный
filter product либо variant, по одному непустому значению каждого фильтра. Operation types
открыты; references не требуют готового SDK каждого вида документа.

## Специальные проверки

Таблица называет momentFrom/To required, но соседний текст и пример допускают defaults: это
нерешённое противоречие. Предлагается explicit range в SDK examples, без заявления о серверной
обязательности до freeze. Limit/offset не переносить из All. Read-only report, DateTime и
unknown operation type, не триGET/all с агрегацией; один server report request.

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
