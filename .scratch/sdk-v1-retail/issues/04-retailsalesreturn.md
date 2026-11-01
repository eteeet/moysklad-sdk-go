Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md; ../../sdk-v1-documents/issues/01-position-contracts.md

# E04: RetailSalesReturn Get/List/Create/Update/GetPositions

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Розничный возврат с based/unbased mode; demand обозначает RetailDemand, не обычную Demand.
Name/org/store/retailStore/retailShift/agent и условные cashSum/noCashSum/demand — по методу,
full field scope без предположений о текущем чеке.

## Специальные проверки

Based/unbased minimal bodies, cash/noCash explicit0, active shift/moment проверки сервером,
собственные update restrictions. Не объявлять все поля изменяемыми по похожести с SalesReturn.

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

[retailsalesreturn](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retail_sales_return.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `64d76cd`, code commit `2114fad`.

| Группа | Read | Create | Update |
| --- | --- | --- | --- |
| name, organization, store, retailStore, retailShift, agent | значения / `*Reference` | обяз. | Field (кроме agent) |
| demand (RetailDemand) | `*Reference` | Field — только по основанию | нет: запрещено |
| agentAccount | `*Reference` | Field | нет: запрещено |
| cashSum, noCashSum | копейки | Field (без основания обязательны; Value(0) уходит) | Field |
| qrSum | копейки | Field | Field |
| organizationAccount, contract, owner, group, state | `*Reference` | Field | Field |
| description, code, externalCode, syncId, moment, applicable, shared, attributes, rate, taxSystem, vatEnabled, vatIncluded | значения | Field | Field |
| positions | `*Collection[Position]` | Field, ≤1000 | Field, ≤1000, заменяет все строки |
| organizationBranch, sum, vatSum, printed, published, id, meta, accountId, created, updated, deleted | | — | — |

- Минимальные тела: по основанию (demand без позиций) и без основания (cashSum и noCashSum = 0) — в тестах.
- Сервер проверяет: активность смены и moment, совпадение контрагента/юрлица/валюты с основанием, строки
  основания с количеством не больше исходного; на Update — валюту и строки не из основания, счёт при смене
  organization, неизменность syncId. Запреты Update взяты из источника, не по аналогии с SalesReturn.
- Противоречие: agent метод называет обязательным, но примеры записи его не передают — SDK следует методу.
- Строки — общая форма D01: cost (без основания), pack, things; GetPositions — `fields=stock`
  (у возврата без основания stock без cost).
- Review Standards: имя обязательно по методу; исключение для общего guide документов — в docs-коммите
  волны. Review Spec: противоречие agent, минимальные тела, комментарий stock.cost, syncId — исправлено.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
