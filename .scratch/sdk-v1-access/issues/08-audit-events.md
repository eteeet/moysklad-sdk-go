Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-access/issues/07-audit-contexts.md

# A08: Audit events и lossless diff

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

EventsByContext и EventsByEntity с безопасным type/id path; типизированные event headers и
raw/lossless diff с typed accessors where justified. Неизвестные event/type/diff fields
сохраняются; не декодировать каждое старое значение в текущую entity model.

## Специальные проверки

Большие integer long, nested/null/removed/unknown fields, event-specific
registration/publication/email diff, отсутствие подтверждённых event pagination guarantees явно
отражено. Worker синхронизации, durable checkpoint и реплика не входят в SDK.

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

[audit](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `38787b5`, code commit `c272bd7`.

- `audit.Service.EventsByContext(ctx, contextID)` — GET audit/{id}/events; `EventsByEntity(ctx, type, id)` —
  GET entity/{type}/{id}/audit. Одна страница; параметров пагинации маршруты не описывают — SDK их не
  передаёт и не обещает полноты (в примере limit 25), nextHref можно пройти pages.All. Тип — безопасный
  сегмент; customentity отклоняется: элементы справочника адресуются с metadataId, маршрут их не покрывает.
- `Event`: типизированные заголовки (audit-контекст и entity как Reference, name, entityType, eventType,
  objectType, objectCount, source, uid, moment, supportAccess, additionalInfo), `Diff` сырым JSON,
  `Raw` — всё событие без потерь. Старые значения не раскладываются в текущие модели сущностей.
- Доступ к diff: `FieldChanges()` — «атрибут → oldValue/newValue» (удаление — только oldValue) со значениями
  jsontext.Value; `DecodeDiff` в `RegistrationDiff`, `PublicationDiff`, `EmailDiff`.
- Не входят: воркер синхронизации, durable checkpoint, реплика.
- Тесты: фикстуры обоими codecs, diff без потерь (большие целые, вложенность, null/отсутствие, неизвестное),
  Raw, типизированные diff, оба маршрута plain/gzip/403 с префиксом, валидация без запросов.
  Независимый golden записи не нужен: Event только читается, round-trip Raw проверен.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review Spec: Raw для неизвестных полей, customentity. Review Standards: gzip и префикс, единые
  метки операций, godoc типов diff.
