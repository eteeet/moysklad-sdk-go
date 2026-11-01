Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md; ../../sdk-v1-documents/issues/01-position-contracts.md

# D04: Loss Get/List/Create/Update/GetPositions

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Списание с organization/store, documented rows и links. Template /new и wizard не добавляются
автоматически.

## Специальные проверки

Required refs/readonly/query, embedded max1000, price/cost/quantity без локальной арифметики,
unknown refs и APIErrors.

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

[loss](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_loss.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `9b27979`, code commit `d403c7e`.

| Группа | Read | Create | Update |
| --- | --- | --- | --- |
| organization, store | `*Reference` | обяз. | Field |
| owner, group, project, state, salesReturn | `*Reference` | Field | Field |
| name, description, code, externalCode, syncId, moment, applicable, shared, attributes, rate | значения | Field | Field |
| positions | `*Collection[Position]` | Field, ≤1000 | Field, ≤1000, заменяет все строки |
| id, meta, accountId, created, updated, deleted, printed, published, sum, files | | — | — |
| expenseItem (readonly по таблице), inventory | `*Reference` | — | — |

- Строки: reason, slot, pack, things; declaration — `GetPositions(fields=declaration)`, readonly.
- Цены, себестоимость и количество не пересчитываются; шаблон PUT /new и wizard не добавлены.
- Противоречие: пример Update передаёт sum, таблица — readonly; SDK sum не пишет (godoc).
- Review Standards: замечаний нет. Review Spec: inventory и противоречие sum записаны; нулевые даты в тесте.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
