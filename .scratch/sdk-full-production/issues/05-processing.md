Status: resolved
Stage: APPROVED
Type: task
Blocked by: 04-processingplan.md; ../../sdk-full-documents/issues/02-document-template.md; ../../sdk-full-documents/issues/04-trash.md

# M05: Processing + materials/products

[Spec](../spec.md), [карта](../../sdk-full/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Техоперация: Get/List/Create/Update/Delete, массовые, New, Trash; материалы и продукты.

## Специальные проверки

processingSum — копейки; processingPlan — ссылка; склады материалов/продуктов.

## Образец

`entity/customerorder` (документ с позициями, GetPositions, batch) и `internal/documenttest`.

## Acceptance criteria

- [ ] Scope/публичный interface и field/method matrix утверждены; зависимости resolved
  и их код присутствует в рабочем дереве. Readonly, required response/Create, null,
  immutable и выбранные TODO разделены по источникам; противоречия не замолчаны.
- [ ] Реализован весь What to build, без неоговорённых методов и hidden requests.
- [ ] Специальные проверки этого ticket реализованы через публичный interface.
- [ ] JSON обоими codecs и независимые golden; upstream/provenance не изменены.
  Для transport/report/print tickets вместо entity JSON применяются их собственные
  публичные seams и контрактные tests, а не фиктивные модели.
- [ ] Реальные transport/httptest tests покрывают успех и ошибки; context и
  metadata сохраняются, caller-controlled retry/walker не подменяет один вызов.
- [ ] go test -race, go vet ./... ./testdata, gofmt и diff-check PASS; независимые
  Standards/Spec review, исправления и финальный BASE/результат записаны в Answer.
- [ ] Фасад, status-entities/status-technical, status-backlog и guide отражают реальный
  scope, ограничения и неизвестные факты. Planning без команды не коммитится; один code
  commit по ticket.

## Источники

[processing](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_processing.md)

## Out of Scope

Expand, live API и не перечисленные What to build операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `c3bace0`, code commit `6ca29e0`.

- CRUD, массовые, New, Trash, материалы и продукты; processingSum — копейки; строки в Update по id.
- Review: [review-wave.md](../review-wave.md); исправления — `3778088`.
- Проверки: make check на `3778088` — PASS.
