Status: resolved
Stage: APPROVED
Type: task
Blocked by: нет

# S01: Store Get/List/Create/Update + zones/slots

[Spec](../spec.md), [карта](../../sdk-full/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Склад: Get/List/Create/Update/Delete, массовые. Вложенные zones и slots: List/Get/Create/Update/Delete и массовые create/update/delete — отдельные методы сервиса склада.

## Специальные проверки

Address/addressFull как у Organization; parent (группа) и pathName записываются — маркера readonly нет, пример PUT передаёт оба; slot ссылается на zone; ID склада в пути вложенных ресурсов проверяется до HTTP.

## Образец

`entity/uom` (service_test.go: TestServiceGet/Create/Update) + batch/delete как в `entity/country`.

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

[store](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_store.md)

## Out of Scope

Expand, live API и не перечисленные What to build операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `30e5d76`, code commit `998d089; refactor d0f7f5e; fix 3af1c7b`; docs `20667cc`.

- Склад Get/List/Create/Update/Delete и массовые; зоны и ячейки — List/Get/Create/Update/Delete и массовые через internal/nested. parent и pathName пишутся. Противоречия: slot.zone readonly в таблице, но источник даёт отдельные примеры создания и изменения ячейки с зоной — пишется в Create и Update: исключение утверждено пользователем 2026-11-03 (после review dd628ca, c81f2ff убирал запись; возвращено коммитом f2f5f09), для остальных полей правило «readonly сильнее примеров» действует; Create/Update зон и ячеек в примерах возвращают массив из одного объекта — принимаются объект и такой массив.
- Review: [standards/spec 6b3e272](../review-6b3e272.md), [round 2 246383e](../review-246383e.md), [dd628ca](../review-dd628ca.md).
- Проверки: make check (gofmt, vet, go test -race -p 1, diff --check) — PASS.
