Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md; ../../sdk-v1-documents/issues/01-position-contracts.md

# E03: RetailDemand Get/List/Create/Update/GetPositions

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Продажа: минимальный Create только retailShift согласно методу; сервер проверяет активность и
moment. Полный documented cash/noCash/QR/prepayment/bonus/fiscal scope зафиксировать отдельной
Write matrix, не требовать все server response поля на Create.

## Специальные проверки

Минимальный Create и full/partial, sums как float64 в копейках, typed pages, immutable
fiscal/read-only fields, server permissions и ошибки, ни автоматической фискализации, ни
бонусного начисления.

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

[retaildemand](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retaildemand.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `c476624`, code commit `64d76cd`.

| Группа | Read | Create | Update |
| --- | --- | --- | --- |
| retailShift | `*Reference` | обяз. (по методу) | Field |
| agent, store | `*Reference` | Field (таблица: «Необходимо при создании» — противоречит методу) | Field |
| agentAccount, organizationAccount, owner, group, state, customerOrder | `*Reference` | Field | Field |
| name, description, code, externalCode, syncId, moment, applicable, shared, attributes, rate, taxSystem, vatEnabled, vatIncluded | значения | Field | Field |
| cashSum, noCashSum, qrSum, prepaymentCashSum, prepaymentNoCashSum, prepaymentQrSum, advancePaymentSum, checkSum | копейки | Field | Field |
| checkNumber, documentNumber, sessionNumber | `int` | Field | Field |
| giftCards | `[]GiftCard{Name, PaymentSum}` | `[]GiftCardRequest{Name}` | то же |
| positions | `*Collection[Position]` | Field, ≤1000 | Field, ≤1000, заменяет все строки |
| organization, organizationBranch (readonly), retailStore (записи нет в источнике), cheque, sum, vatSum, payedSum, printed, published, id, meta, accountId, created, updated, deleted, files | | — | — |

- Противоречия: обязательность agent/store (таблица против метода — SDK следует методу); sum и payedSum в
  примерах записи при маркере readonly; organization readonly при тексте Update о смене organization.
- Правила оплаты (сумма полей = сумме продажи, без смешанной оплаты с QR, без аванса вместе с предоплатой,
  распределение остатка) — сервер; тест отправляет несходящуюся оплату как есть. Активность смены и moment
  позже смены — сервер. SDK не фискализирует и не начисляет бонусы.
- Строки — общая форма D01: cost (услуги, Int в таблице, float64 копейки), pack, things, trackingCodes;
  trackingCodes_1162 и declaration только читаются. GetPositions — `fields=stock` или `declaration`.
- syncId после заполнения не меняется — проверяет сервер (как у прочих документов).
- Review Standards: замечаний нет. Review Spec: матрица и противоречия записаны здесь; docs — в
  завершающем коммите волны.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
