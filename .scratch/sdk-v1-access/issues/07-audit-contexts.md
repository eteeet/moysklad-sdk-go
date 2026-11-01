Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md

# A07: Audit contexts и filters

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Audit.ListContexts и GetFilters по official routes; фильтры
moment/uid/source/entityType/eventType/application по documented scope; typed AuditContext и
FilterCatalog, metadata страницы. Default limit25, max100 именно contexts.

## Специальные проверки

Не обещать endless changelog/cursor/retention/exactly-once; timestamp tie/duplicates возможны.
Отсутствие прав/неизвестный source/entityType, контекст Employee vs App, exact filters, одна
страница без чтения events.

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

[audit](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md), [general](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `1c6f83a`, code commit `38787b5`.

- Пакет `audit`: `ListContexts(ctx, ContextsOptions, opts...)` — одна страница GET audit без чтения
  событий (limit 0 → default сервера 25, max 100, offset, filter; других параметров маршрут не описывает);
  `GetFilters` — GET audit/metadata/filters. Типы `audit.Context` и `audit.FilterCatalog` (в ticket —
  AuditContext/FilterCatalog; имя пакета уже даёт префикс).
- Фильтры (godoc): moment >= <= в МСК через DateTime.QueryValue, employee href, uid, source, eventType,
  entityType (код или href пользовательского справочника), application (href установки, не appUid).
- Противоречия источника записаны: events — Array(Object) в таблице, объект с meta в примерах;
  info — Object в таблице, строка в примере списка, объект в фикстуре (хранится сырым);
  ключи каталога entityType/eventType (таблица) vs entitytype/eventtype (пример) — читаются оба.
- Критерии успеха: страница — Collection (meta.href, rows); каталог — ключ source.
- Не обещается курсор, хранение, exactly-once и порядок при равных moment (godoc пакета).
- Тесты: фикстуры обоими codecs, открытые значения, info строкой/объектом, null, оба регистра ключей,
  точный query для всех полей фильтра, страница без context, 403 с ResponseInfo, max 100, локальные ошибки.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review Spec: оба регистра ключей, objectType убран из фильтров, href справочника, info, Extra убран.
  Review Standards: критерий source у каталога, FilterCatalog, doc типов, wantErr, t.Run.
