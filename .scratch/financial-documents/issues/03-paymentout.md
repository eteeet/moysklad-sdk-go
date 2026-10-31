Status: needs-triage
Stage: DRAFT
Blocked by: 01-finance-operation.md; ../../document-common/issues/01-rate.md

# 03: PaymentOut — полный Get/Create/Update slice

**What to build:** разработчик решения получает PaymentOut, создаёт исходящий платеж
с обязательной статьёй расходов, частично меняет его и читает фактическое состояние
из ответа сервера. Slice включает read/request models, Service, JSON/HTTP-тесты и docs.

Это DRAFT; write scope, numeric/null-policy и validation требуют triage до реализации.

## Контракт и зависимости

- Прочитать [spec](../spec.md), полную таблицу PaymentOut и общие правила
  [contract](../contract.md), [исследование](../research.md).
- Blockers: [01 — FinanceOperation](01-finance-operation.md),
  [document-common/01 — DocumentRate](../../document-common/issues/01-rate.md).
  Реализация начинается после утверждения и Status: resolved обоих, с их кодом в рабочем
  дереве. PaymentIn/CashIn/CashOut, positions, SDK банковских счетов не являются blockers.
- Proposal package paymentout: PaymentOut, CreateRequest, UpdateRequest, Service/NewService.
  Get(ctx,id): GET entity/paymentout/{id}; Create(ctx,input): POST entity/paymentout;
  Update(ctx,id,input): PUT entity/paymentout/{id}; (*PaymentOut,error). Через существующий
  transport.Client, query пустой, GET без тела, ровно один запрос без pre-GET/follow-up.
- Read, полный proposed scope: meta; id/accountId/name/description/code/externalCode/
  syncId/paymentPurpose — string; created/updated/deleted/moment/accrualDate — DateTime;
  applicable/shared/printed/published/noClosingDocs — bool; sum/vatSum — float64;
  attributes — []Attribute; agent/agentAccount/organization/organizationAccount/owner/group/
  contract/project/salesChannel/state/expenseItem/factureIn — *Reference;
  rate — *DocumentRate; operations — []FinanceOperation. References/DateTime omitzero, как Product/UOM.
  Files, incomingDate/incomingNumber отсутствуют в scope.
- Proposal Create write scope: accrualDate, agent, agentAccount, applicable, attributes,
  code, contract, description, expenseItem, externalCode, group, noClosingDocs, moment, name,
  organization, organizationAccount, owner, paymentPurpose, project, rate, salesChannel,
  shared, state, sum, syncId, vatSum, factureIn, operations. Обязательные Organization,
  Agent, ExpenseItem — обычные Reference; остальные поля Field[T]. Минимальный JSON —
  только три эти ссылки с meta; Name/sum/moment необязательны, expenseItem default не описан.
- Update: тот же scope, включая syncId, все поля Field[T]; минимально applicable:false.
  Readonly accountId/id/meta/created/updated/deleted/printed/published в requests отсутствуют.
  SyncID доступен в Create/Update; запрет изменения после заполнения проверяет сервер без предварительного Get.
- Agent может ссылаться на Counterparty/Employee/Organization; AgentAccount — счет
  Counterparty/Organization. Не ограничивать Reference закрытым union или локальным enum.
- При смене сторон API требует соответствующие account-поля; несовместимая валюта
  organizationAccount/rate даёт 22004. NoClosingDocs=true несовместим с непустым operations.
  Proposal: сервер проверяет оба ограничения, SDK не подбирает счета и не вводит валидатор.
- Общие DocumentRateRequest/FinanceOperationRequest приходят из blockers. Outgoing operations:
  salesreturn, supply, invoicein, purchaseorder, commissionreportout, payroll.
  Meta+linkedSum, linkedSum optional write/explicit 0; без expanded объектов или загрузки оснований.
- Float64 Sum/VATSum/LinkedSum — proposal без масштабирования/округления/арифметики;
  не обещать decimal-precision. Существующий DateTime/Attribute; Field+omitzero сохраняет
  unset/null/false/0/""/[] и nil slice как []. Contract:null и rate:null документированы;
  rate:null запускает server recalculation. Другие null/очистка operations через []
  не объявляются гарантированными. Attributes [] не очищает всю коллекцию.
- Proposal local validation: nil client, UUID формы Get/Update без проверки версии,
  Create.Meta.Href обязательных refs не пуст, Update хотя бы одно Field. Name==""
  не запрещать как в Product; без проверки business constraints/types/null/sum.
- Proposal UOM response/error rules: HTTP 200 + непустой ID, при ошибке nil результата,
  errors.As/errors.Is и контекст операции без повторного moysklad:. Retry отсутствует.

## Acceptance criteria

- [ ] Новые решения triage утверждены; blockers resolved и доступны в рабочем дереве.
- [ ] Read всех полей проверен по значениям encoding/json и encoding/json/v2;
  polymorphic agent каждого документированного типа, sparse поля, DateTime, false/0,
  expenseItem/noClosingDocs/accrualDate, references и operations с linkedSum.
- [ ] Независимые golden JSON minimal/full Create, full/partial Update; обязательные
  три ссылки, отсутствие defaults, явные noClosingDocs=false, sum/vatSum=0, дроби,
  "", [], nil slice как [], null/unset, SyncID в Create/Update и отсутствие readonly/excluded полей.
- [ ] HTTP-тесты построены как entity/uom/service_test.go: отдельные TestServiceGet,
  TestServiceCreate, TestServiceUpdate с таблицами, локальный HTTP и реальный transport.Client.
  Точные method/path/query/GET-body/headers, plain/gzip, серверные данные отличаются от input,
  один запрос без retry/preliminary GET.
- [ ] Minimal Create отправляет только organization/agent/expenseItem; minimal Update
  передаёт false. Проверен полный body по независимому golden.
- [ ] Nil/отменённый context, local validation по triage; APIError 400/404/429/5xx,
  synthetic noClosingDocs conflict и 22004; auth/transport wrapping, errors.As/errors.Is,
  unexpected status, malformed JSON, null/empty/missing ID. Synthetic fixture не доказывает
  реальное business-поведение сервера.
- [ ] Upstream fixtures неизменены с provenance; собственные unexpanded cases проверяют
  accrualDate и отсутствующие поля. Expanded operation fields/files.rows игнорируются.
- [ ] Docs/coverage фиксируют полный scope, типы agent, expenseItem required Create,
  ограничения счетов/валют/noClosingDocs, float64/null/operations и TODO files/expand.
- [ ] Форматирование, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check; независимые Standards/Spec reviews с фиксированным BASE.

## Out of Scope

Остальные сущности, invoice SDK, positions, банковские CRUD/collections, файлы, expand,
List/Delete/batch/metadata/templates, финансовая арифметика, транспортные доработки,
живой API/.env. Сейчас только DRAFT, без реализации/commit/push.

## Источники

[Official PaymentOut, master SHA проверен 2026-11-01](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_out.md),
[общие DocumentRate/operations](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_common_info.md#L476-L699),
[null](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187).
