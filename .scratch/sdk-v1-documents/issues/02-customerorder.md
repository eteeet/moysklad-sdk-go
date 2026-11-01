Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md; ../../sdk-v1-documents/issues/01-position-contracts.md

# D02: CustomerOrder Get/List/Create/Update/GetPositions

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Полный CustomerOrder header/positions scope, required organization/agent, reserve/shipped и
shipment/payment links по source. Header arrays, attrs, prices и positions не имеют одной общей
merge semantics.

## Специальные проверки

Minimal/full/golden, readonly fields, typed fixed reader, list/search/filter, exact 1000/1001 и
accounts/null/positions replacement, server-responses вместо input.

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

[customerorder](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_customerOrder.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `448ff35`, code commit `c4dbe13`.

| Группа | Read | Create | Update |
| --- | --- | --- | --- |
| organization, agent | `*Reference` | обяз. | Field |
| organizationAccount, agentAccount, contract, owner, group, project, state, store, salesChannel | `*Reference` | Field | Field |
| name, description, code, externalCode, syncId, moment, deliveryPlannedMoment, applicable, shared | значения | Field | Field |
| attributes, rate, taxSystem, vatEnabled, vatIncluded | значения | Field | Field |
| shipmentAddress, shipmentAddressFull | `string`, `*Address` | Field | Field |
| positions | `*Collection[Position]` (meta) | Field, ≤1000 | Field, ≤1000, заменяет все строки |
| id, meta, accountId, created, updated, deleted, printed, published | | — | — |
| sum, vatSum, payedSum, shippedSum, invoicedSum, reservedSum | копейки | — | — |
| files | Collection | — (файлы — отдельный ресурс) | — |
| purchaseOrders, demands, payments, productionTasks, invoicesOut, moves, prepayments | `[]Reference` | — | — |

- Строки — общая форма D01: reserve пишется, shipped только читается, pack/taxSystem.
- GetPositions — одна страница, `fields=stock` (остатки на текущий момент, общий раздел источника).
- Правила сервера только в godoc: при смене organization/agent нужен и счёт; syncId после заполнения
  не меняется; адресу null не поддерживается, удаление — shipmentAddress "".
- Union-поля строк D01, которых нет у заказа (cost, inTransit, slot, reason …), доступны через общий тип;
  принимает ли их заказ, решает API — решение D01, не дополнение по аналогии.
- Review Standards: раздельные ошибки organization/agent; в тесте лимитов — единый префикс и операция,
  локальные ошибки GetPositions и Update. Дублирование do/validID — конвенция пакетов.
- Review Spec: правило null адреса в godoc UpdateRequest; Files в тесте помечен «не пишется», а не readonly.
  Docs/status — в завершающем docs-коммите волны.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
