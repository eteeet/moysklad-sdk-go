Status: needs-triage
Stage: DRAFT
Type: task
Blocked by: None

# T03: Каталог кодов ошибок и status-aware helpers

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Полный числовой каталог официального _errors.md (545 уникальных кодов в закреплённом revision):
именованные константы с provenance, стабильный int-compatible Code, распознавание status/code
без сравнения локализованных строк. Неизвестные коды сохраняются. Сохранить errors.As/Is и
единственный префикс.

## Специальные проверки

Все code в источнике покрыты inventory test; aliases/дубли описаны. 413 возвращает
APIError/детали, не делает autosplit или retry. HTML/empty/gzip errors и raw headers доступны.
Никаких новых retry решений по тексту ошибки.

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
