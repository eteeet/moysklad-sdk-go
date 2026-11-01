Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md

# E01: RetailStore Get/List/Create/Update

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Точка продаж: name/organization/store/priceType Create, полный documented core/read-write flags.
Различать metadata кассиров и специальные cashier endpoints; отдельный cashier administration не
включается молча.

## Специальные проверки

allowCreateProducts/controlShippingStock conflict и тариф проверяет сервер, no hidden config
reads; tax/price/payment settings exact JSON, readonly secure fields не отправляются.

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

[retailstore](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `2f21bcf`, code commit `69c8611`; общий тест опций — `410a255`.

| Группа | Read | Create | Update |
| --- | --- | --- | --- |
| name, organization, store, priceType | значения / ссылки | обяз. | Field |
| 29 флагов (active, allowCreateProducts, controlShippingStock, required*, sync*, qrPayEnabled …) | `bool` | Field | Field |
| address, description, demandPrefix, idQR, qrTerminalId | `string` | Field | Field |
| 8 перечислений (defaultTaxSystem, orderTaxSystem, fiscalType, markingSellingMode, marksCheckMode, minionToMasterType, priorityOfdSend, tobaccoMrcControlType) | открытые `string` | Field | Field |
| bankPercent, qrBankPercent, discountMaxPercent | `float64` (проценты) | Field | Field |
| acquire, qrAcquire, createOrderWithState, orderToState, organizationBranch, owner, group, receiptTemplate | `*Reference` | Field | Field |
| addressFull | `*Address` | Field | Field |
| createAgentsTags, filterAgentsTags | `[]string` | Field | Field |
| customerOrderStates | `[]Reference` | Field | Field |
| masterRetailStores, productFolders | `*Links` (коллекция или массив) | Field `[]Reference` | Field |
| externalCode, allowCustomPrice, authTokenAttached, ofdEnabled, environment, state, lastOperationNames, id, meta, accountId, updated | readonly | — | — |
| cashiers | `*Collection` (meta) | — (администрирование кассиров не входит) | — |

- environment и state — `jsontext.Value` без потерь: readonly, SDK их не интерпретирует; вложенные таблицы
  источника (ПО, ККТ, ФН, синхронизация) не типизированы — осознанное решение.
- Противоречия: createAgentsTags — Array(Object) в таблице, строки в описании и примерах (SDK: строки);
  discountMaxPercent — Int в таблице, 100.0 в фикстуре (float64); masterRetailStores/productFolders —
  объект коллекции в фикстуре и массив в примере ответа Create (`Links` читает обе формы); пример Create
  передаёт readonly allowCustomPrice и authTokenAttached (не пишутся).
- Конфликт allowCreateProducts/controlShippingStock, тариф и валюта счёта организации — сервер; тест
  отправляет конфликт как есть. Скрытых чтений настроек нет.
- Review Standards: замечаний нет. Review Spec: `Links` для двух форм, причина исключения cashiers
  исправлена (не readonly, а вне scope), golden расширен (перечисления, ссылки, null).
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
