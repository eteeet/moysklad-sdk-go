Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md; ../../sdk-v1-documents/issues/01-position-contracts.md

# D05: SalesReturn Get/List/Create/Update/GetPositions

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Возврат покупателя по Demand и без основания: organization/store/agent обязательны, demand
required только для based. Update запреты сторон/основания/цен по source, без копирования
PurchaseReturn запретов по аналогии.

## Специальные проверки

Два режима Create, full/partial Update, position matching ID/quantity constraints оставлены
серверу, типы row/read-only и source contradictions фиксируются. Existing Demand service не
нужен для отправки Reference.

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

[salesreturn](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_sales_return.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `d403c7e`, code commit `826adc3`; общий тест опций — `42771d9`.

| Группа | Read | Create | Update |
| --- | --- | --- | --- |
| organization, store, agent | `*Reference` | обяз. | organization, store — Field; agent — нет |
| agentAccount, demand | `*Reference` | Field (demand — только по основанию) | нет: источник запрещает менять |
| organizationAccount, contract, owner, group, project, state, salesChannel | `*Reference` | Field | Field |
| name, description, code, externalCode, syncId, moment, applicable, shared, attributes, rate, vatEnabled, vatIncluded | значения | Field | Field |
| positions | `*Collection[Position]` | Field, ≤1000 | Field, ≤1000, заменяет все строки |
| id, meta, accountId, created, updated, deleted, printed, published, sum, vatSum, payedSum, files | | — | — |
| losses, payments, factureOut | ссылки | — | — |

- Два режима Create одним типом: с Demand (по основанию) и без; golden на оба. Совпадение контрагента,
  юрлица, валюты и строк с основанием проверяет сервер, Demand-сервис для ссылки не нужен.
- «Запреты цен» из ticket — в источнике это запрет валюты, отличной от основания; цену строки источник
  менять не запрещает. Запреты PurchaseReturn по аналогии не перенесены.
- Строки: cost (без основания), country, gtd, slot, pack, things; `GetPositions(fields=stock)`, stock без
  cost у возврата без основания (`Stock.Cost` nil). `Position.Cost` общей формы — float64: отсутствие
  у возврата по основанию читается как 0 (смена на указатель сломала бы прочие документы).
- Union-поля строк D01 (reserve, inTransit, slots …) доступны через общий тип; принимает ли их документ,
  решает API — решение D01.
- Review Standards: неиспользуемый параметр тестового хелпера убран. Review Spec: ограничение stock в
  Get/List/Update (только с expand) — в godoc GetPositions.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
