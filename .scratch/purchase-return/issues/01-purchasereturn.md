Status: resolved
Stage: IMPLEMENTED
Type: task
Blocked by: ../../document-common/issues/01-rate.md; ../../document-positions/issues/01-json-contracts.md; ../../document-positions/issues/02-get-page.md

> Реализовано по команде пользователя выполнить все 11 tickets. Предложения и
> запреты запуска в исходном DRAFT ниже отражают этап подготовки; команда исполнения
> их заменила. Итоговый контракт — docs/purchasereturn.md, результат — в Answer.

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
- Rate — model.DocumentRate/DocumentRateRequest. Positions читается через
  *model.Collection[positions.Position], запись Create/Update — Field[[]positions.PositionRequest].
  Собственные формы строки не создаются. Matching ID — selector существующей строки
  embedded Update, не назначение ID новой Create строки.
- Общий PositionRequest содержит коммерческие Field и ячейки; unset не отправляется,
  допустимость заданного проверяет сервер. Overhead/Meta/AccountID в request отсутствуют.
- Ограничения основания и расхождения price/vat/discount сохраняются в документации:
  Discount помечен readonly в таблице, хотя есть в write examples. Его Field в общем
  request не обещает изменение. Quantity-only отдельного PUT не переносится на embedded Update.
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
Position — публичный alias positions.Position.
Технический blocker P02 обязателен. Реализация сущности остаётся needs-triage до отдельного разрешения.

## Acceptance criteria

- [x] Все open decisions spec согласованы до реализации; DRAFT переведён отдельным triage,
  обе зависимости resolved. Готовый Supply или SDK связанных сущностей не требуются.
- [x] Model/requests/методы покрывают полную таблицу spec; readonly/TODO заголовка отсутствуют
  в request документа, общие Rate/позиции не дублируются. В общем request строки нет
  Meta/AccountID/Overhead; допустимость заданных коммерческих Field проверяет сервер.
  Запрет Update сторон выражен в request заголовка, не в отдельной форме позиции.
- [x] Оба JSON-кодека проверяют значения каждого поддержанного поля upstream fixture
  с provenance, optional/unknown/null, даты/ссылки/readonly Discount и metadata без rows.
- [x] Независимые golden: Create с основанием/без него и полный, полный/одно поле Update,
  все zero/unset/null/[] states, matching ID, частичный rate и отсутствие запрещённых полей.
- [x] HTTP-структура строго entity/uom/service_test.go: отдельные TestServiceGet/Create/Update
  с table cases, настоящий transport.Client + httptest.Server. Exact method/path/query/
  headers/body, plain/gzip, одна отправка, фактический server response без merge/input fallback.
- [x] Проверены local validation без HTTP, 1000/1001 позиций, nil/cancelled context,
  auth/transport/APIError 4xx/5xx/429, errors.As/Is и operation context без двойного префикса;
  unexpected status, malformed/null JSON, пустой/отсутствующий ID, nil-модель при ошибке.
- [x] Синтетическая ошибка ограничений основания передаётся без retry; локальный тест
  не выдаётся за доказательство правил настоящего API. Get не вызывает GetPage.
- [x] Документация и coverage фиксируют минимум Create, два режима, запреты Update,
  full replacement/1000, общую форму строк и серверные ограничения, отдельный reader, TODO и неоднозначности.
- [x] Пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata,
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

## Triage 2026-11-01

Пользователь поручил выполнить все 11 entity tickets SDK. Предложенные таблицы полей,
local validation и ограничения scope приняты к реализации этой командой. Общие зависимости resolved.
BASE: eba6f75154929790c5b9d2cb1454dd80783b14db; исходные изменения сохранены.

## Answer

Реализован `entity/purchasereturn`: PurchaseReturn, отдельные CreateRequest/UpdateRequest,
Get/Create/Update и типизированный GetPositions, aliases общей формы позиции.
[Публичный контракт](../../../docs/purchasereturn.md). JSON и HTTP проверки, required/readonly
scope и ограничения покрыты; общие типы не дублируются.

Полный race/vet/gofmt/diff-check PASS; независимые Standards 0 / Spec 0.
Upstream fixtures/provenance неизменны.
[Интеграционный отчёт](../../document-workstream/review-entities.md).
