Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md

# C02: PriceType List/Get/GetDefault/ReplaceAll

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Существующий model.PriceType переиспользуется без создания несовместимого второго типа. GET
context/companysettings/pricetype возвращает plain array; GET /{id}, GET /default; POST полного
упорядоченного списка ReplaceAll. Никакого вымышленного одиночного PUT.

## Специальные проверки

1…100, default первый, порядок сохраняется, excluded item удаляется сервером. Exact full-list
JSON, no hidden GET/read-modify-write; externalCode/name/readonly shape по методу; возможная
потеря конкурентных изменений явно документируется.

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

[pricetype](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_pricetypes.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `1e29d0d`, code commit `d9aced4`.

- Пакет `entity/pricetype` на context/companysettings/pricetype: `List` (голый массив), `Get(id)`,
  `GetDefault`, `ReplaceAll([]ReplaceItem)` — POST полного упорядоченного списка. Одиночных PUT/POST нет.
- Ответы — существующий `model.PriceType`. Запись — `ReplaceItem{Meta, Name, ExternalCode *string}`:
  model.PriceType не годится (meta не опускается у новых, externalCode "" опускается).
- ReplaceAll локально: 1…100 элементов (ограничение API), непустое (не из пробелов) имя у всех, meta с href
  у существующих; externalCode — указатель: null не отправить («нельзя удалить»), "" — сервер сгенерирует.
  `Existing(p)` передаёт meta, имя и внешний код. Уникальность имён, удаление исключённых и выбор первого
  типом по умолчанию — сервер. Pre-GET нет; потеря параллельных изменений описана в godoc.
- Критерии успеха: непустой массив (все типы удалить нельзя) и id у одиночного ответа.
- Тесты: фикстура, golden полного списка обоими codecs, HTTP-таблица (plain/gzip, 403, 412, null,
  malformed, неверная форма, 204, один запрос), LocalErrors, валидация без запросов.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review Spec: null externalCode невозможен, Existing с externalCode, href у meta, имя у всех, doc List.
  Review Standards: ReplaceItem, omitzero, читаемая таблица тестов, t.Helper.
- Регистрация в общих тестах опций и сброса ResponseInfo (`entity/options_test.go`) — `cc61f6a` (найдено в review C07).
