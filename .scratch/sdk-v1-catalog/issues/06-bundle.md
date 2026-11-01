Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md

# C06: Bundle и components

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Bundle Get/List/Create/Update, components read collection и explicit GetComponents/GetComponent
при подтверждённых endpoints; embedded components write с ID/quantity/assortment по своему
contract. Components не притворяются документными Position, если отличаются смыслом.

## Специальные проверки

Полная таблица bundle/components, required/minimal, replacement/limits/null из source,
product/service/variant refs без expand; salePrices/attributes сохраняют свои merge rules.
Отдельные mutations компонентов не добавляются без фиксации scope.

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

[bundle](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bundle.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `1d90849`, code commit `efabaa8`.

| Поле | Create | Update | Примечание |
| --- | --- | --- | --- |
| name | обяз. | + | |
| components | обяз., непустой | + | `ComponentRequest{ID, Assortment, Quantity}`; слияние по ID |
| code, externalCode, description, article, archived, shared | + | + | |
| vat, vatEnabled, useParentVat, discountProhibited | + | + | |
| paymentItemType, taxSystem, trackingType, partialDisposal, tnved, weight, volume | + | + | enum открыты |
| owner, group, productFolder, uom, country | + | + | Reference |
| minPrice, overhead, salePrices | + | + | копейки; salePrices сливаются по типу цены |
| barcodes | + | + | заменяются целиком; при Create пусто/нет → случайный EAN13 |
| attributes | + | + | сливаются |
| mod__tasnif__uz, mod__marking__uz | + | + | только с `WithAcceptRegion(RegionUZ)` |
| syncId | + | — | |
| id, meta, accountId, updated, pathName, effectiveVat*, images, files | — | — | readonly / не пишутся в v1 |

| Метод | Endpoint |
| --- | --- |
| Get / List / Create / Update | entity/bundle[/{id}] |
| `GetComponents(id, query.PageOptions)` | GET entity/bundle/{id}/components (limit/offset) |
| `GetComponent(id, componentID)` | GET entity/bundle/{id}/components/{componentID} |

- Component — отдельный тип (assortment + quantity), не Position. Добавление/изменение/удаление
  отдельных компонентов в v1 не входит: состав меняется через Update.
- Противоречия: quantity помечен readonly, но примеры записи его передают — SDK пишет; minPrice (таблица) /
  minprice (пример Update) — по таблице; assortment «товар/услуга/партия» в таблице и «модификация» во
  вступлении — тип открыт.
- Null: семантики источник не описывает; Field.Null отправляется как есть (как у остальных сущностей),
  факт не подтверждён.
- Review Standards: локальные проверки лимита 50 и assortment убраны (бизнес-ограничения — серверу),
  MaxComponents удалён; тип страницы вынесен в `query.PageOptions` с префиксом `moysklad:` и тестом.
- Review Spec: ID у ComponentRequest для изменения существующих компонентов; компоненты больше не
  зависят от пакета positions; null отмечен как неподтверждённый.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Регистрация в общих тестах опций и сброса ResponseInfo (`entity/options_test.go`) — `cc61f6a` (найдено в review C07).
