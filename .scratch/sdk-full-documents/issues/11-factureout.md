Status: resolved
Stage: APPROVED
Type: task
Blocked by: 02-document-template.md; 04-trash.md

# O11: FactureOut

[Spec](../spec.md), [карта](../../sdk-full/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Счёт-фактура выданный: Get/List/Create/Update/Delete, массовые, New, Trash.

## Специальные проверки

Без позиций; demands/returns/payments — ссылки; stateContractId.

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

[factureout](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_factureout.md)

## Out of Scope

Expand, live API и не перечисленные What to build операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `d21d84d`, code commit `2cd7740`.

- FactureOut: CRUD, массовые, New, Trash; авансовые поля.
- Review: [review-wave.md](../review-wave.md); исправления — `2e9e03d`, `3b62b18`, `341a4fd`.
- Проверки: make check CHECK_P=1 на `341a4fd` — PASS.
