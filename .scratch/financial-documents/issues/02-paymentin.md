Status: needs-triage
Stage: DRAFT
Blocked by: 01-finance-operation.md; ../../document-common/issues/01-rate.md

# 02: PaymentIn — полный Get/Create/Update slice

**What to build:** разработчик решения получает PaymentIn по ID, создаёт его с organization
и agent, частично меняет выбранные поля и получает фактическое состояние из ответа сервера.
Slice включает модель, оба request-типа, Service, публичные JSON/HTTP-тесты и документацию.

Это DRAFT; поимённый write scope, новые типы и validation требуют triage до реализации.

## Контракт и зависимости

- Прочитать [spec](../spec.md), полную таблицу PaymentIn и общие правила
  [contract](../contract.md), [расхождения источников](../research.md).
- Blockers: [01 — FinanceOperation](01-finance-operation.md) и
  [внешний document-common/01 — DocumentRate](../../document-common/issues/01-rate.md).
  Начинать реализацию после утверждения и Status: resolved обоих blockers, с их кодом
  в рабочем дереве. Не зависит от PaymentOut/CashIn/CashOut, positions или SDK банковских счетов.
- Proposal package paymentin: PaymentIn, CreateRequest, UpdateRequest, Service/NewService.
  Get(ctx,id): GET entity/paymentin/{id}, Create(ctx,input): POST entity/paymentin,
  Update(ctx,id,input): PUT entity/paymentin/{id}. Возврат (*PaymentIn,error).
  Пустой query; GET без тела; ровно один запрос через существующий transport.Client.
- Read, полный proposed scope: meta; id/accountId/name/description/code/externalCode/
  syncId/incomingNumber — string; created/updated/deleted/moment/incomingDate — DateTime;
  applicable/shared/printed/published — bool; sum/vatSum — float64; attributes — []Attribute;
  agent/agentAccount/organization/organizationAccount/owner/group/contract/project/
  salesChannel/state/factureOut — *Reference; rate — *DocumentRate; operations — []FinanceOperation.
  Ссылки и DateTime используют omitzero по образцу Product/UOM. Files исключён.
- Proposal Create write scope: agent, agentAccount, applicable, attributes, code, contract,
  description, externalCode, group, incomingDate, incomingNumber, moment, name, organization,
  organizationAccount, owner, paymentPurpose, project, rate, salesChannel, shared, state,
  sum, syncId, vatSum, factureOut, operations. PaymentPurpose — string в read-модели.
  Organization и Agent — обязательные обычные Reference, остальные поля Field[T].
  Минимальный JSON содержит только organization и agent с meta. Name/sum/moment
  необязательны; серверные defaults SDK не подставляет.
- Update имеет тот же scope, включая syncId, все поля Field[T]. Минимальный Update —
  applicable:false. Readonly accountId/id/meta/created/updated/deleted/printed/published
  отсутствуют в обоих request-типах; SyncID предлагается в обоих requests; запрет изменения после заполнения проверяет сервер.
- Agent документация описывает как Counterparty; account-поля — только Reference.
  При смене organization/agent API требует соответствующий account; организация с валютным
  счетом должна согласовать rate.currency с ним, иначе 22004. SDK не загружает счета,
  не проверяет сочетания локально, передаёт APIError. Банковский CRUD не нужен.
- DocumentRateRequest и FinanceOperationRequest брать из blockers. Operations — meta+optional
  linkedSum, incoming типы customerorder/purchasereturn/demand/invoiceout/commissionreportin/
  retailshift; неизвестный meta.type не отвергается локально. Никаких expanded объектов.
- Sum/VATSum/LinkedSum — предложенный float64 без масштабирования/округления/арифметики;
  DateTime — существующий контракт, включая отсутствие/null. Ведущие нули incomingNumber
  сохраняются. Field с omitzero сохраняет false/0/""/[]/null/unset и nil slice как [].
- Contract:null и rate:null документированы; rate:null пересчитывает сервер. Для прочих
  null документация не даёт всеобщей гарантии очистки. [] operations отправляется точно,
  но не объявляется доказанной командой полной очистки. Attributes — существующий merge
  по metadata; сброс отдельного значения через value:null, не очистка всей коллекции [].
- Proposal validation: nil client; UUID формы без ограничения версии; Create обязательные
  Reference.Meta.Href != ""; Update хотя бы одно заданное Field. Нет проверки Name=="",
  typed refs, сумм/коллекций/валют и прочих бизнес-правил. Решение ещё требует утверждения.
- Proposal response/errors по UOM: HTTP 200 + непустой ID, при ошибке nil; контекст операции,
  errors.As/errors.Is, один moysklad:; никаких retry/pre-GET/follow-up по ссылкам.

## Acceptance criteria

- [ ] Triage утвердил scope/числа/null/required refs/SyncID в Create/Update; blockers resolved.
- [ ] Все перечисленные read-поля проверены по значениям обоими JSON-кодеками;
  read sparse/null даты, неизвестные поля, references и operation linkedSum декодируются.
- [ ] Create/Update JSON проверен независимо заданными golden: minimal/full Create,
  full/partial Update, unset, false, sum/vatSum=0, дроби, "", [], nil slice, null,
  attributes values, SyncID в Create и Update. Readonly/excluded поля отсутствуют в requests.
- [ ] Отдельные TestServiceGet/Create/Update с table cases по entity/uom/service_test.go;
  используется локальный HTTP-сервер и настоящий transport.Client. Проверены точные
  method/path/пустой query, GET без тела, заголовки, plain/gzip, один запрос,
  ответ с данными, отличающимися от input.
- [ ] HTTP minimal Create содержит только две required ссылки; minimal Update false
  не считается пустым. Full HTTP body совпадает с независимым golden.
- [ ] Проверены nil/отменённый context, локальная validation по утверждённому triage,
  transport/auth errors, APIError 400/404/429/5xx и 22004 как synthetic server error;
  errors.As/errors.Is, неожиданный статус, malformed JSON, null/empty/missing ID, без retry.
- [ ] Использованы неизменённые upstream fixtures с существующим provenance;
  собственный маленький unexpanded case дополняет incomingDate/incomingNumber,
  paymentPurpose, false/0 и отсутствующие поля. Expanded operation/files поля fixture
  игнорируются и не заявляются поддержанными.
- [ ] Документация сущности и coverage перечисляют полный read/write scope, readonly,
  account References, SyncID в Create/Update, null/float64/operations ограничения и TODO files/expand.
- [ ] Форматирование, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check; независимые Standards/Spec reviews с фиксированным BASE.

## Out of Scope

Остальные три сущности, invoicein/out SDK, positions, файлы, банковские счета как сущности,
expand/query expand, List/Delete/batch/metadata/templates, Money/Decimal и расчёты,
универсальный CRUD и транспортные доработки, живой API/.env. Текущий DRAFT — без кода/commit/push.

## Источники

[Official PaymentIn, master SHA проверен 2026-11-01](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_in.md),
[DocumentRate/operations](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_common_info.md#L476-L699),
[null](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187).
