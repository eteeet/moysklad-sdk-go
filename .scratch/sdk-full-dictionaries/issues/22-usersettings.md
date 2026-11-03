Status: resolved
Stage: APPROVED
Type: task
Blocked by: нет

# S22: UserSettings GET/PUT

[Spec](../spec.md), [карта](../../sdk-full/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Настройки пользователя: GET/PUT context/usersettings и metadata.

## Специальные проверки

Значения enum открытые; defaultScreen, printFormat, locale по таблице.

## Образец

`contexts` (GET context/*), тесты — `contexts/contexts_test.go`.

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

[usersettings](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_usersettings.md)

## Out of Scope

Expand, live API и не перечисленные What to build операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `30e5d76`, code commit `1f6fd24; fix dd628ca`; docs `20667cc`.

- mailFooter — string по примеру и фикстуре (таблица: Boolean). Default*: пример — Meta без обёртки, фикстура — {meta}; читаются обе формы, пишется форма источника. Метода context/usersettings/metadata источник не описывает — не добавлен.
- Review: [standards/spec 6b3e272](../review-6b3e272.md), [round 2 246383e](../review-246383e.md).
- Проверки: make check (gofmt, vet, go test -race -p 1, diff --check) — PASS.
