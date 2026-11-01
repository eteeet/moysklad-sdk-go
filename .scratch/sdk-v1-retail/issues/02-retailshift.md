Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md

# E02: RetailShift Get/List/Create/Update

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Смена: required organization/retailStore. Не добавлять GetPositions к смене без соответствующего
API; attached collections остаются typed references/meta. Открытие/закрытие/фискальные side
effects только документированным Update и с явным input.

## Специальные проверки

Активные/закрытые смены, required flags/refs, sparse server data, C/U readonly rules, APIError и
отсутствие автоматически созданной смены.

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

[retailshift](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retailshift.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `69c8611`, code commit `c476624`.

| Группа | Read | Create | Update |
| --- | --- | --- | --- |
| organization, retailStore | `*Reference` | обяз. | Field |
| organizationAccount, organizationBranch, store, owner, group, acquire, qrAcquire | `*Reference` | Field | Field |
| name, description, externalCode, syncId, moment, closeDate, shared, attributes | значения | Field | Field |
| bankComission, qrBankComission, bankPercent, qrBankPercent | `float64` | Field | Field |
| cheque {start, end} | `*Cheque` | `ChequeRequest` (Field) | `ChequeRequest` |
| proceedsCash/NoCash, receivedCash/NoCash (копейки), printed, published, meta, id, accountId, created, deleted, files | | — | — |
| operations | `[]Reference` | — | — |
| paymentOperations | `[]FinanceOperation` (linkedSum) | — | — |
| updated | `DateTime` | — (нет маркера readonly, но запись не описана) | — |

- Позиций и GetPositions у смены нет. Открытие — Create, закрытие — явный closeDate/cheque в Update; SDK
  смену сам не создаёт и не ищет. Открытая смена: closeDate и cheque.end отсутствуют (тест).
- Store без значения сервер берёт из точки; комиссии не больше безналичных платежей — сервер.
- syncId после заполнения не меняется (сервер); Null у organization/retailStore отправляется как есть.
- Фикстура содержит agentAccount, contract, vatEnabled, vatIncluded, которых нет в таблице смены, — не
  моделируются.
- Review Standards: отдельный `ChequeRequest` на Field — 0 и "" больше не теряются (golden с chequesTotal 0).
- Review Spec: paymentOperations — `model.FinanceOperation` с linkedSum; единица комиссий не утверждается.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
