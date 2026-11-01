Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md; ../../sdk-v1-catalog/issues/05-characteristics.md

# C07: Variant Get/List/Create/Update

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Полный Variant contract, product link и characteristics, prices/barcodes/taxes и ограничения
immutable fields. Metadata-related definitions отдельны; Product не загружается автоматически.

## Специальные проверки

Create required и update restrictions по источнику; unknown characteristic refs, false/0/[],
missing/null, upstream+custom, list/filter/search. Full entity Product service не технический
blocker references.

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

[variant](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `efabaa8`, code commit `9eccf22`.

| Поле | Create | Update | Примечание |
| --- | --- | --- | --- |
| product | обяз. | + | Reference; товар не загружается. Смену источник не описывает — решает сервер |
| characteristics | обяз., непустой | + | `CharacteristicValue{ID, Name, Value}`; Update заменяет набор целиком |
| name, code, externalCode, description, article, archived, discountProhibited | + | + | |
| buyPrice, minPrice, salePrices | + | + | копейки; salePrices сливаются по типу цены, без своих — цены товара |
| barcodes | + | + | заменяются целиком |
| packs | + | + | `PackRequest{Barcodes, ParentPack}`; один штрихкод — проверяет сервер |
| id, meta, accountId, updated, things, images | — | — | readonly |
| minimumStock | — | — | отложен, как у Product; List отклоняет fields=minimumStock |

- Чтение значений характеристик — `Characteristic` из C05; описания — отдельные методы C05.
- List: search только по name, фильтр `productid` (= и !=), barcode в единственном числе — в godoc, тест productid.
- Противоречия: правила штрихкодов при создании в разделе variant скопированы из bundle («при создании
  комплекта») — отмечено как неподтверждённое; «taxes» из ticket — в источнике у модификации нет НДС/налогов.
- Null: семантики источник не описывает; Field.Null отправляется как есть.
- Review Standards: проверки префикса и операции в ошибках валидации; регистрация в общем тесте опций и
  сброса ResponseInfo — отдельным test-коммитом волны вместе с остальными сервисами каталога (пропущены и в C01–C06).
- Review Spec: исправлен комментарий о штрихкодах, отмечено отсутствие налогов, тесты неизвестной
  характеристики (один запрос, APIError) и фильтра productid.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
