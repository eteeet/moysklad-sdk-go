Status: resolved
Stage: IMPLEMENTED
Blocked by: 01-finance-operation.md; ../../document-common/issues/01-rate.md

> Реализовано по команде пользователя выполнить все 11 tickets. Предложения и
> запреты запуска в исходном DRAFT ниже отражают этап подготовки; команда исполнения
> их заменила. Итоговый контракт — docs/cashout.md, результат — в Answer.

# 05: CashOut — полный Get/Create/Update slice

**What to build:** разработчик решения получает CashOut, создаёт расходный ордер
с обязательной статьёй расходов, меняет выбранные поля и получает сохранённое состояние.
Slice включает модель/requests, Service, JSON/HTTP-тесты и документацию.

Это DRAFT: поимённый scope, числовое представление/null и validation ожидают triage.

## Контракт и зависимости

- Прочитать [spec](../spec.md), полную CashOut-таблицу/общие правила
  [contract](../contract.md), [исследование](../research.md).
- Blockers: [01 — FinanceOperation](01-finance-operation.md),
  [document-common/01 — DocumentRate](../../document-common/issues/01-rate.md).
  Реализация после утверждения и Status: resolved обоих, с их кодом в рабочем дереве.
  PaymentIn/PaymentOut/CashIn, positions, банковский SDK и модели оснований не blockers.
- Proposal package cashout: CashOut, CreateRequest, UpdateRequest, Service/NewService.
  Get(ctx,id): GET entity/cashout/{id}, Create(ctx,input): POST entity/cashout,
  Update(ctx,id,input): PUT entity/cashout/{id}; (*CashOut,error).
  Существующий transport.Client; пустой query, GET без тела, ровно один запрос.
- Read, полный proposed scope: meta; id/accountId/name/description/code/externalCode/
  syncId/paymentPurpose — string; created/updated/deleted/moment/accrualDate — DateTime;
  applicable/shared/printed/published/noClosingDocs — bool; sum/vatSum — float64;
  attributes — []Attribute; agent/organization/owner/group/contract/project/salesChannel/
  state/expenseItem/factureIn — *Reference; rate — *DocumentRate; operations — []FinanceOperation.
  References/DateTime omitzero как Product/UOM. Files исключён; account-ссылок и
  incomingDate/incomingNumber в CashOut нет.
- Proposal Create write scope: accrualDate, agent, applicable, attributes, code, contract,
  description, expenseItem, externalCode, group, noClosingDocs, moment, name, organization,
  owner, paymentPurpose, project, rate, salesChannel, shared, state, sum, syncId, vatSum,
  factureIn, operations. Organization/Agent/ExpenseItem — ordinary required Reference;
  остальные Field[T]. Minimal JSON — только три эти ссылки с meta. Name/sum/moment
  необязательны, defaults expenseItem не описаны, SDK ничего не подставляет.
- Update: тот же scope, включая syncId, все поля Field[T]; минимально applicable:false.
  Readonly accountId/id/meta/created/updated/deleted/printed/published в requests отсутствуют.
  SyncID доступен в обоих requests; проверка write-once остаётся серверу.
- CashOut.agent описан Counterparty. Polymorphic agent PaymentOut нельзя автоматически
  объявить документированным CashOut-контрактом; Reference остаётся открытым без validator.
  CashOut.Update повторяет требование account-полей, которых нет в атрибутах.
  Proposal — не включать их, явно описать расхождение, не проверять account-пары локально.
  Валютная проверка банковского organizationAccount прямо не применяется к ордерам.
- NoClosingDocs=true несовместим с непустым operations. Proposal — сервер проверяет
  это, SDK передаёт APIError без локального бизнес-валидатора.
- DocumentRateRequest/FinanceOperationRequest из blockers. Operations только meta+optional
  linkedSum; outgoing типы salesreturn, supply, invoicein, purchaseorder,
  commissionreportout, payroll. Никаких typed unions, expanded refs или загрузки оснований.
- Sum/VATSum/LinkedSum — float64 proposal без конвертации, округления и арифметики.
  Field+omitzero сохраняет unset/null/false/0/""/[] и nil slice как []; DateTime/Attribute
  существующие. Contract:null и rate:null документированы; rate:null пересчитывает сервер.
  Другие null/operations:[] не объявляются гарантией очистки/замещения; attributes []
  не очищает всю коллекцию, reset отдельных значений — существующий value:null.
- Proposal local validation: nil client; UUID форма Get/Update без ограничения версии;
  непустой Meta.Href трёх required Create refs; хотя бы одно Field Update. Без Name==""
  validator, проверки сумм/types/null/операций или автоматического получения связанных данных.
- Proposal UOM response/errors: HTTP 200 + непустой ID, nil при ошибке; operation context,
  errors.As/errors.Is, один moysklad:. Без retry/preliminary GET/follow-up.

## Acceptance criteria

- [x] Новые решения triage утверждены; blockers resolved, их код в рабочем дереве.
- [x] Все read-поля проверены по значениям encoding/json и encoding/json/v2;
  accrualDate/expenseItem/noClosingDocs, sparse/null даты, refs, rate без value,
  operation linkedSum, неизвестные поля. Не объявлены недоказанные agent-типы CashOut.
- [x] Независимые golden minimal/full Create и full/partial Update: required refs,
  отсутствие defaults/readonly/accounts, noClosingDocs=false, sum/vatSum=0, дробь,
  "", [], nil slice как [], null/unset, SyncID в Create/Update.
- [x] HTTP образец entity/uom/service_test.go: отдельные TestServiceGet/Create/Update
  с table cases, локальный сервер и настоящий transport.Client. Проверены точные
  method/path, пустой query, GET без тела, заголовки, plain/gzip, один запрос,
  фактические server data отличаются от input.
- [x] Minimal Create содержит только organization/agent/expenseItem; minimal Update
  передаёт applicable:false; full body соответствует независимому golden.
- [x] Nil/отменённый context, утверждённая local validation, APIError 400/404/429/5xx,
  synthetic noClosingDocs conflict, transport/auth wrapping, errors.As/errors.Is,
  unexpected status, malformed JSON, null/empty/missing ID и отсутствие retry.
  Local server не доказывает бизнес-поведение реального API или account-правило CashOut.
- [x] Upstream fixtures неизменены с provenance; собственный unexpanded case покрывает
  accrualDate/отсутствующие поля. Expanded operation/files.rows данные игнорируются.
- [x] Docs/coverage перечисляют scope, required expenseItem, noClosingDocs restriction,
  CashOut.agent/account неопределённости, float64/null/operations и files/expand TODO.
- [x] Форматирование, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check; независимые Standards/Spec reviews с фиксированным BASE.

## Out of Scope

Остальные сущности/invoice SDK, positions, account-поля/банковский SDK, файлы, expand,
List/Delete/batch/metadata/templates, новая финансовая арифметика, транспортные доработки,
живой API/.env. Текущая работа только DRAFT, без кода/commit/push.

## Источники

[Official CashOut, master SHA проверен 2026-11-01](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashout.md),
[общие операции/DocumentRate](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_common_info.md#L476-L699),
[null](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187).

## Triage 2026-11-01

Пользователь поручил выполнить все 11 entity tickets SDK. Предложенные таблицы полей,
local validation и ограничения scope приняты к реализации этой командой. Общие зависимости resolved.
BASE: eba6f75154929790c5b9d2cb1454dd80783b14db; исходные изменения сохранены.

## Answer

Реализован `entity/cashout`: CashOut, отдельные CreateRequest/UpdateRequest,
Get/Create/Update.
[Публичный контракт](../../../docs/cashout.md). JSON и HTTP проверки, required/readonly
scope и ограничения покрыты; общие типы не дублируются.

Полный race/vet/gofmt/diff-check PASS; независимые Standards 0 / Spec 0.
Upstream fixtures/provenance неизменны.
[Интеграционный отчёт](../../document-workstream/review-entities.md).
