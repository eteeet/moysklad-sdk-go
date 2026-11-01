Status: resolved
Stage: IMPLEMENTED
Blocked by: 01-finance-operation.md; ../../document-common/issues/01-rate.md

> Реализовано по команде пользователя выполнить все 11 tickets. Предложения и
> запреты запуска в исходном DRAFT ниже отражают этап подготовки; команда исполнения
> их заменила. Итоговый контракт — docs/cashin.md, результат — в Answer.

# 04: CashIn — полный Get/Create/Update slice

**What to build:** разработчик решения получает CashIn по ID, создаёт приходный ордер
с organization/agent и меняет выбранные поля одним PUT, получая фактический ответ сервера.
Slice включает модель, requests, Service, публичные JSON/HTTP-тесты и документацию.

Это DRAFT. Scope, числа/null и validation — proposals; реализация до triage не разрешена.

## Контракт и зависимости

- Прочитать [spec](../spec.md), полную CashIn-таблицу и общие правила
  [contract](../contract.md), [исследование](../research.md).
- Blockers: [01 — FinanceOperation](01-finance-operation.md) и
  [внешний document-common/01 — DocumentRate](../../document-common/issues/01-rate.md).
  Реализация после утверждения и Status: resolved обоих, с их кодом в рабочем дереве.
  PaymentIn/PaymentOut/CashOut, positions и SDK банковских счетов не блокируют CashIn.
- Proposal package cashin: CashIn, CreateRequest, UpdateRequest, Service/NewService.
  Get(ctx,id): GET entity/cashin/{id}; Create(ctx,input): POST entity/cashin;
  Update(ctx,id,input): PUT entity/cashin/{id}; (*CashIn,error).
  Существующий transport.Client, пустой query, GET без тела, один запрос.
- Read, полный proposed scope: meta; id/accountId/name/description/code/externalCode/
  syncId/paymentPurpose — string; created/updated/deleted/moment — DateTime;
  applicable/shared/printed/published — bool; sum/vatSum — float64; attributes — []Attribute;
  agent/organization/owner/group/contract/project/salesChannel/state/factureOut — *Reference;
  rate — *DocumentRate, operations — []FinanceOperation. References/DateTime omitzero по Product/UOM.
  Files целиком исключён; банковские поля, incomingDate/incomingNumber, accrualDate,
  expenseItem и noClosingDocs отсутствуют.
- Proposal Create write scope: agent, applicable, attributes, code, contract, description,
  externalCode, group, moment, name, organization, owner, paymentPurpose, project, rate,
  salesChannel, shared, state, sum, syncId, vatSum, factureOut, operations.
  Organization/Agent — обычные required Reference; остальные — Field[T]. Минимальный
  JSON содержит только эти две ссылки с meta. Name/sum/moment необязательны Create,
  серверные defaults не подставляются SDK.
- Update: тот же scope, включая syncId, все поля Field[T]; минимально applicable:false.
  Нет readonly accountId/id/meta/created/updated/deleted/printed/published в requests;
  SyncID доступен в обоих requests; проверка write-once остаётся серверу.
- Agent описан Counterparty. Account-ссылки отсутствуют в таблице CashIn; требование
  CashIn.Update менять organizationAccount/agentAccount вместе со сторонами противоречит
  этому. Proposal — не добавлять банковские поля по повторённой фразе, не валидировать
  account-пары локально, сохранить caveat в docs. Проверка валюты банковского счета
  прямо не применяется к cash ордерам.
- DocumentRateRequest и FinanceOperationRequest из blockers. Incoming operations:
  customerorder, purchasereturn, demand, invoiceout, commissionreportin, retailshift;
  только meta+optional linkedSum, открытый MetaType, без expanded references.
- Sum/VATSum/LinkedSum — float64 proposal, без округления/конвертации/арифметики.
  Field+omitzero сохраняет unset/null/false/0/""/[] и nil slice как []; DateTime/Attribute
  существующие. Contract:null и rate:null документированы; rate:null пересчитывает сервер.
  Для прочих null и operations:[] не обещать серверную очистку/замещение. Attributes []
  не очищает всю коллекцию, reset отдельных значений — value:null через existing helpers.
- Proposal local validation: nil client, форма UUID без ограничения версии Get/Update,
  непустой Meta.Href обеих required Create ссылок, непустой по Field Update. Без Name==""
  validator, проверки типов ссылок, business constraints или local merge/pre-GET.
- Proposal UOM response/error rules: HTTP 200 + непустой ID, nil результата при ошибке;
  errors.As/errors.Is и operation context без дублирования moysklad:. Retry отсутствует.

## Acceptance criteria

- [x] Scope/новые решения утверждены triage; blockers resolved и доступны в рабочем дереве.
- [x] Каждый read-field проверен по значениям encoding/json и encoding/json/v2;
  sparse/null даты, references, rate с отсутствующим value, operations с linkedSum,
  неизвестные поля не мешают декодированию.
- [x] Независимые golden minimal/full Create, full/partial Update покрывают обе required
  ссылки, отсутствие defaults/readonly/excluded, false, sum/vatSum=0, дроби, "", [],
  nil slice как [], null/unset, SyncID в Create и Update. Нет банковских request-полей.
- [x] HTTP структура entity/uom/service_test.go: отдельные TestServiceGet/Create/Update
  с table cases и настоящим transport.Client на локальный сервер. Точные method/path,
  пустой query, отсутствие тела GET, заголовки, plain/gzip, один запрос, server data != input.
- [x] Minimal Create — только organization/agent; minimal Update applicable:false;
  full body соответствует независимому golden и не содержит выдуманных defaults.
- [x] Проверены nil/отменённый context, validation по triage; APIError 400/404/429/5xx,
  auth/transport errors, errors.As/errors.Is, unexpected status, malformed JSON,
  null/empty/missing ID и отсутствие retry. Не имитировать недоказанное account-правило.
- [x] Upstream fixtures остаются неизменными с provenance; собственный маленький
  unexpanded JSON дополняет scope. Expanded поля operations/files.rows игнорируются,
  наличие их в fixture не объявляется поддержкой expand.
- [x] Docs/coverage отражают полный scope, required Create refs, account-противоречие,
  отсутствие account/expenseItem/noClosingDocs, null/float64/operations ограничения и TODO files.
- [x] Форматирование, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check; независимые Standards/Spec reviews с фиксированным BASE.

## Out of Scope

Остальные сущности/invoice SDK, positions, банковские счета/ссылки CashIn, файлы,
expand/query expand, List/Delete/batch/metadata/templates, финансовые вычисления,
транспортные доработки и live API/.env. Текущий DRAFT — без кода/commit/push.

## Источники

[Official CashIn, master SHA проверен 2026-11-01](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashin.md),
[общие операции/DocumentRate](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_common_info.md#L476-L699),
[null](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187).

## Triage 2026-11-01

Пользователь поручил выполнить все 11 entity tickets SDK. Предложенные таблицы полей,
local validation и ограничения scope приняты к реализации этой командой. Общие зависимости resolved.
BASE: eba6f75154929790c5b9d2cb1454dd80783b14db; исходные изменения сохранены.

## Answer

Реализован `entity/cashin`: CashIn, отдельные CreateRequest/UpdateRequest,
Get/Create/Update.
[Публичный контракт](../../../docs/cashin.md). JSON и HTTP проверки, required/readonly
scope и ограничения покрыты; общие типы не дублируются.

Полный race/vet/gofmt/diff-check PASS; независимые Standards 0 / Spec 0.
Upstream fixtures/provenance неизменны.
[Интеграционный отчёт](../../document-workstream/review-entities.md).
