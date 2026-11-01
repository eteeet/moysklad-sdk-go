Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md

# A05: Общий каталог и metadata сущностей

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

GET entity/metadata с filter=type=..., GET entity/{type}/metadata, чтение definitions attributes
и отдельного attribute metadata, states/reference lists по фактическому endpoint. Open entity
types safe segment; named supported routes, без expand.

## Специальные проверки

Global response не обязательно Collection; формы attributes/meta/rows/plain arrays сверить по
источнику. Не навязывать ID. Разные metadata поддерживаемых типов, unknown definitions/types и
auth errors. Создание/редактирование definitions — отдельно от чтения и не входит скрыто.

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

[general](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `c272bd7`, code commit `bee5782`.

| Метод | Маршрут | Критерий успеха |
| --- | --- | --- |
| Catalog(types) | GET entity/metadata[?filter=type=a;type=b] | объект «тип → {meta}» |
| Get(type) | GET entity/{type}/metadata | meta.href |
| GetCustomEntity(id) | GET entity/customentity/{id}/metadata | meta.href |
| GetCompanySettings | GET context/companysettings/metadata | meta.href |
| ListAttributes(type, page) | GET entity/{type}/metadata/attributes (limit ≤ 1000) | коллекция |
| GetAttribute(type, id) | GET entity/{type}/metadata/attributes/{id} | meta.href |

- `Entity`, `CompanySettings` хранят весь ответ в `Raw` (tags контрагента и др.); известные поля —
  createShared, states, attributes. `attributes` принимает обе формы источника: массив (38 примеров) и
  коллекцию {meta[, rows]} (product, фикстуры) — `AttributeList`.
- Противоречия: id описания доп. поля «обязателен» в таблице, но отсутствует в примере customerorder —
  критерий meta.href, ID не навязывается. Catalog хранит только meta: других полей источник не описывает.
- Тип сущности — безопасный сегмент; пользовательские справочники — через GetCustomEntity и
  GetCompanySettings. Характеристики модификаций — C05. Создание/изменение описаний доп. полей не входит.
- Тесты: 7 upstream-фикстур metadata и companysettings обоими codecs, Raw без потерь, массив attributes
  без id, общий HTTP-контракт documenttest для всех 6 методов, 401/403/404 с кодами, точный filter
  каталога, валидация без запросов.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review Spec: форма-массив attributes (было падение декодирования), без требования id, маршруты
  справочников, ошибки сервера. Review Standards: gzip/401/403/заголовки через documenttest, LocalErrors
  на всех методах, сообщения с %q, единый do.

## Comments

Ревью 3306518 (P2): AttributeList писался как {"Meta","Rows"}. Исправлено в `d446807`: MarshalJSON в форме чтения, round-trip тесты.
