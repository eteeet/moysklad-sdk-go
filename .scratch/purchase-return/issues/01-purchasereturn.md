Status: needs-triage
Stage: DRAFT
Type: task
Blocked by: ../../document-common/issues/01-rate.md; ../../document-positions/issues/01-json-contracts.md; ../../document-positions/issues/02-get-page.md

# 01: PurchaseReturn — Get/Create/Update

## What to build

Законченный срез возврата поставщику: публичная модель, отдельные requests документа
и строк, Service.Get/Create/Update, JSON/HTTP-проверки, документация и coverage.
Нормативные таблицы всех полей и правил — [полная spec](../spec.md).
Сейчас это DRAFT, не разрешение реализации. Перед началом прочитать spec целиком.

## Предлагаемый контракт

- Модель PurchaseReturn читает все поля таблицы spec: core, стороны/счета/склад,
  rate, attributes, positions metadata, readonly суммы/даты/flags и ссылки связей.
- Create требует Organization/Agent/Store Reference. Supply Field[Reference]
  задаётся для возврата по основанию, без основания отсутствует. Name необязателен.
  Прочие разрешённые поля — Field с omitzero; SyncID доступен Create/Update, write-once проверяет сервер.
- Update использует Field и исключает Agent/AgentAccount/Supply,
  readonly и TODO поля; SyncID разрешён как Field[string] без local write-once validation. Organization/OrganizationAccount передаются явно вместе
  при необходимости; старое состояние не загружается для проверки.
- Rate — model.DocumentRate/DocumentRateRequest. Read Positions —
  *model.Collection[positions.Commercial]. CommercialRequest не является write-моделью возврата.
- CreatePositionRequest: Assortment Field[Reference], Quantity/Price Field[float64],
  VAT Field[int], VATEnabled Field[bool]. Нет ID/Meta/AccountID/Discount.
- UpdatePositionRequest: встраивает positions.ItemRequest и добавляет Price/VAT/
  VATEnabled через Field, JSON плоский; нет Meta/AccountID/Discount.
  Matching ID — selector существующей строки, не назначение ID новой.
- Price/VAT/VATEnabled embedded Update — PROPOSED по таблице и примеру; ограничения
  основания требуют triage. Quantity-only отдельного PUT позиции не переносить
  автоматически на embedded Update. Discount read-only по таблице, write — TODO.
- Get GET entity/purchasereturn/{id}, Create POST entity/purchasereturn,
  Update PUT entity/purchasereturn/{id}; (*PurchaseReturn,error), 200+непустой ID.
  transport.Client, context/errors.As/Is/gzip по UOM; nil-модель при ошибке.
- Proposed validation: nil client, UUID без ограничения версии, обязательные Create
  Meta.Href!="", непустой по Field Update, максимум 1000 embedded позиций.
  Не проверять основание/счета/остатки дополнительным HTTP и не добавлять бизнес-валидаторы.
- Get без query/body делает один запрос, rows не обещает. Строки читает явный общий
  GetPage с documentType=purchasereturn; никаких expand/followups/retry.
- Omitted/[]/null/false/0/"" различаются; positions — full replacement без local merge.
  Кодируемый null не обещает серверного clearing. Rate:null не отменяет ограничения валюты основания.

## Типизированное чтение позиций

Entity Service предоставляет GetPositions(ctx context.Context, id string, page positions.PageOptions)
(*model.Collection[Position], error). DocumentType и тип строки фиксированы пакетом; вызывающий
код не передаёт строку ресурса и type parameter. Фиксированный documentType — purchasereturn.
Один GET через общий reader, без expand,
follow-up и auto-pagination. Ошибки, UUID/page validation и response envelope — как у P02.
Отдельные TestServiceGetPositions с таблицами проверяют фиксированный resource path, правильный
Position, limit/offset, plain/gzip, пустую страницу, ошибки и ровно один запрос через настоящий transport.
Position — публичный alias positions.Commercial.
Технический blocker P02 обязателен. Реализация сущности остаётся needs-triage до отдельного разрешения.

## Acceptance criteria

- [ ] Все open decisions spec согласованы до реализации; DRAFT переведён отдельным triage,
  обе зависимости resolved. Готовый Supply или SDK связанных сущностей не требуются.
- [ ] Model/requests/методы покрывают полную таблицу spec; readonly/TODO отсутствуют в
  записи, общие Rate/позиции не дублируются. Запрет Update сторон выражен в публичном request.
- [ ] Оба JSON-кодека проверяют значения каждого поддержанного поля upstream fixture
  с provenance, optional/unknown/null, даты/ссылки/readonly Discount и metadata без rows.
- [ ] Независимые golden: Create с основанием/без него и полный, полный/одно поле Update,
  все zero/unset/null/[] states, matching ID, частичный rate и отсутствие запрещённых полей.
- [ ] HTTP-структура строго entity/uom/service_test.go: отдельные TestServiceGet/Create/Update
  с table cases, настоящий transport.Client + httptest.Server. Exact method/path/query/
  headers/body, plain/gzip, одна отправка, фактический server response без merge/input fallback.
- [ ] Проверены local validation без HTTP, 1000/1001 позиций, nil/cancelled context,
  auth/transport/APIError 4xx/5xx/429, errors.As/Is и operation context без двойного префикса;
  unexpected status, malformed/null JSON, пустой/отсутствующий ID, nil-модель при ошибке.
- [ ] Синтетическая ошибка ограничений основания передаётся без retry; локальный тест
  не выдаётся за доказательство правил настоящего API. Get не вызывает GetPage.
- [ ] Документация и coverage фиксируют минимум Create, два режима, запреты Update,
  full replacement/1000, различия request строк, отдельный reader, TODO и неоднозначности.
- [ ] Пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check и независимые Standards/Spec review с фиксированным BASE.

## Blocked by / readiness

- [DocumentRate](../../document-common/issues/01-rate.md).
- [JSON positions](../../document-positions/issues/01-json-contracts.md).

[GetPage](../../document-positions/issues/02-get-page.md) — обязательная зависимость entity ticket с GetPositions. Supply, другие документы и services справочников
не являются blockers. Требуются решения price/vat/discount, unbased constraints,
local validation и SyncID из spec; готовность зависимостей не заменяет triage.

## Sources и границы

[Первичные источники и Testing Decisions](../spec.md#further-notes),
[исследование](../../stock-documents/research.md). Официальный pin проверен 2026-11-01.

List/query/expand/Delete/batch/templates, отдельные mutations позиций, pack/slot/serial/
tracking/stock, новые transport features, live API и .env вне scope. Сейчас сохраняется
только DRAFT; Go-код, tests реализации, commits подготовки и push не выполняются.
