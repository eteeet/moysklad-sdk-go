Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md; ../../sdk-v1-documents/issues/01-position-contracts.md

# D03: Enter Get/List/Create/Update/GetPositions

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Оприходование с organization/store, без выдуманного agent; собственный read/write header и поля
строк, shared types там где совпадают.

## Специальные проверки

Минимальный Create из двух required refs, fields и read-only sums, nil dates, stock/positions
queries только где документированы, одна отправка.

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

[enter](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_enter.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `c4dbe13`, code commit `9b27979`.

| Группа | Read | Create | Update |
| --- | --- | --- | --- |
| organization, store | `*Reference` | обяз. | Field |
| owner, group, project, state | `*Reference` | Field | Field |
| name, description, code, externalCode, syncId, moment, applicable, shared, attributes, rate | значения | Field | Field |
| overhead `{sum, distribution}` | `*DocumentOverhead` | Field | Field |
| positions | `*Collection[Position]` | Field, ≤1000 | Field, ≤1000, заменяет все строки |
| id, meta, accountId, created, updated, deleted, printed, published, sum, files | | — | — |
| inventory | `*Reference` | — | — |

- Контрагента у оприходования нет. Строки: reason, country, gtd, things, slot, pack; overhead строки readonly.
- GetPositions — одна страница; fields у позиций источник не документирует (отклоняется до запроса).
- inventory: маркера readonly нет, но и в примерах записи нет — SDK не пишет (в godoc).
- Review Standards: замечаний нет (дублирование do/validID — конвенция пакетов). Review Spec: формулировка
  об inventory исправлена; нулевые даты проверяются в тесте.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
