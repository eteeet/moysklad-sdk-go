Status: resolved
Stage: IMPLEMENTED
Type: task
Blocked by: ../../document-common/issues/01-rate.md; ../../document-common/issues/02-overhead.md; ../../document-positions/issues/01-json-contracts.md; ../../document-positions/issues/02-get-page.md

> Реализовано по команде пользователя выполнить все 11 tickets. Предложения и
> запреты запуска в исходном DRAFT ниже отражают этап подготовки; команда исполнения
> их заменила. Итоговый контракт — docs/move.md, результат — в Answer.

# 01: Move — Get/Create/Update

## What to build

Законченный срез Move (перемещение): разработчик получает документ, создаёт его с двумя складами и частично меняет его, передавая positions без коммерческих discount/VAT. Результат включает публичные model/requests/Service, JSON/HTTP-проверки и документацию. [Spec](../spec.md) — DRAFT; до реализации все proposed choices согласуются.

## Полный предлагаемый контракт

Все optional write поля — model.Field с omitzero; read References — указатели. DateTime/Attribute — существующие типы model.

| JSON-поля | Read | Create | Update |
| --- | --- | --- | --- |
| meta; id, accountId | Meta; string | — | — |
| created, updated, deleted | DateTime | — | — |
| name, description, code, externalCode | string | optional Field[string] | Field[string] |
| moment | DateTime | optional Field[DateTime] | Field[DateTime] |
| applicable, shared | bool | Field[bool] | Field[bool] |
| organization, sourceStore, targetStore | *Reference | обязательные Reference | Field[Reference] |
| owner, group, project, state, internalOrder, customerOrder | *Reference | Field[Reference] | Field[Reference] |
| demand, supply | *Reference | — | —, API readonly |
| productionTasks | []Reference | TODO write | TODO write |
| rate | *DocumentRate | Field[DocumentRateRequest] | Field[DocumentRateRequest] |
| overhead | *DocumentOverhead | Field[DocumentOverheadRequest] | Field[DocumentOverheadRequest] |
| positions | *Collection[positions.Position] | Field[[]positions.PositionRequest] | Field[[]positions.PositionRequest] |
| attributes | []Attribute | Field[[]Attribute] | Field[[]Attribute] |
| syncId | string | Field[string] | Field[string], write-once проверяет сервер |
| sum | float64 | — | — |
| printed, published | bool | — | — |
| files | *Collection[Reference], только metadata | — | — |

Go ID/AccountID/SyncID, sourceStore → SourceStore, targetStore → TargetStore. Нет agent/agentAccount/organizationAccount/store/contract/VAT полей документа. Минимальный Create — organization/sourceStore/targetStore, без name/positions.

Строки используют общие positions.Position/PositionRequest без собственной формы. Общая Position читает meta/id/accountId, assortment, quantity/price, read-only overhead и optional sourceSlot/targetSlot. Неиспользуемые discount/vat/vatEnabled дают zero value: у move это отсутствие поля, не «без НДС». В общем request эти Field существуют, но обычно не задаются; допустимость заданных проверяет сервер. Overhead в request отсутствует. Новая строка требует assortment и положительный quantity по API, matching ID применяется для replacement. Pack/things/serial/tracking/stock — TODO read/write. Таблица допускает service, вводный текст не перечисляет его; локального assortment validator нет до решения.

Общий DocumentOverhead/Request: sum float64, открытый distribution weight/volume/price; документ без positions не принимает overhead по API, но SDK не загружает его для проверки.

## Операции и семантика

- NewService использует transport.Client и отклоняет nil; Get GET entity/move/{id}, Create POST entity/move, Update PUT entity/move/{id}; возвращают (*Move, error).
- Get без query/body делает один запрос, читает positions.meta; rows получают явно через общий GetPage. Нет expand, href follow-up и автоматического обхода.
- Proposed local validation: UUID и полностью незаданный Update по UOM, непустой Meta.Href обязательных Create refs, не более 1000 embedded positions. Доступность refs, совпадение/различие складов и бизнес-условия оставлены серверу; нельзя добавлять недоказанные ограничения.
- Успех — HTTP 200 с непустым id; ошибки — nil-модель, контекст операции без двойного moysklad:, сохранённые APIError/errors.As/Is.
- Unset Field не отправляется; false/0/""/[]/null сохраняются. DocumentRate:null сбрасывает валюту на default с пересчётом. Допустимость null по каждому полю не выводится из наличия Field.
- Positions omitted сохраняет прежние строки; [] очищает по правилу replacement; совпавший ID обновляется, отсутствующие удаляются, новые добавляются. Positions:null не равно []; локального merge нет, embedded максимум 1000 по API. Attributes имеет собственный общий контракт.
- InternalOrder/customerOrder write и SyncID в Create/Update — proposed scope. Demand/supply запрещены в requests по readonly marker; productionTasks write остаётся TODO.

## Типизированное чтение позиций

Entity Service предоставляет GetPositions(ctx context.Context, id string, page positions.PageOptions)
(*model.Collection[Position], error). DocumentType и тип строки фиксированы пакетом; вызывающий
код не передаёт строку ресурса и type parameter. Фиксированный documentType — move.
Один GET через общий reader, без expand,
follow-up и auto-pagination. Ошибки, UUID/page validation и response envelope — как у P02.
Отдельные TestServiceGetPositions с таблицами проверяют фиксированный resource path, правильный
Position, limit/offset, plain/gzip, пустую страницу, ошибки и ровно один запрос через настоящий transport.
Position — публичный alias positions.Position; PositionRequest — alias positions.PositionRequest.
Технический blocker P02 обязателен. Реализация сущности остаётся needs-triage до отдельного разрешения.

## Acceptance criteria

- [x] Перед реализацией согласованы таблица, proposed exports, local validation и open decisions; dependencies resolved.
- [x] Model/requests/Service соответствуют таблице полностью, общие DocumentRate/Overhead/позиции не дублируются. Row overhead отсутствует в общем request; неиспользуемые Field пропускаются, допустимость заданных проверяет сервер.
- [x] Оба JSON-кодека проверяют значения каждого поля upstream move fixture, отсутствующие/null/unknown поля, References, DateTime, metadata без rows, sourceSlot/targetSlot. Overhead проверен как 0, 0.0 и дробное число обоими кодеками, без обещания точности integer >2^53.
- [x] Независимые golden: минимальный/полный Create, полный/одно поле Update, omitted/false/0/""/[]/null, пустой overhead object и частичный overhead, нет agent/VAT заголовка и readonly demand/supply/row overhead. В строках move коммерческие Field обычно не заданы и пропускаются; общий тип не запрещает их локально.
- [x] Proposed validation проверена без HTTP при пустом Meta.Href обязательных refs и 1001 позиции; 1000 позиций допустимы. Positions:null/omitted остаются отдельными состояниями.
- [x] entity/uom/service_test.go — обязательная структура: отдельные TestServiceGet/Create/Update с table cases, реальный transport + httptest.Server, method/path/query/headers/body, plain/gzip, один запрос, фактический server response.
- [x] Все три метода: context/auth/transport/APIError 4xx/5xx/429, errors.Is/As, unexpected 201/204, malformed/null JSON, пустой/отсутствующий id; без retry и побочных запросов.
- [x] Документация и coverage перечисляют поля/TODO, minimum Create, двухскладскую модель, общую форму строк и смысл отсутствующих VAT-полей, full replacement/1000, Get metadata и отдельный GetPage.
- [x] Пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata, git diff --check и предусмотренное review относительно BASE; исходные изменения сохранены отдельно.

## Blocked by / readiness

[DocumentRate](../../document-common/issues/01-rate.md), [Overhead](../../document-common/issues/02-overhead.md), [JSON positions](../../document-positions/issues/01-json-contracts.md). Нет зависимости от готовых Supply/Demand/store/slot services. [GetPage](../../document-positions/issues/02-get-page.md) — обязательная зависимость типизированного GetPositions этого ticket.

## Explicit Open Decisions

Утвердить writable internalOrder/customerOrder, SyncID в Create/Update, service assortment, TODO pack/things, local refs validation и server-null semantics. До этого DRAFT не ready.

## Sources и границы

[Поля/связи](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_move.md#L8-L60), [positions](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_move.md#L62-L89), [Create](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_move.md#L317-L323), [Update](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_move.md#L1432-L1444), [Create position](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_move.md#L2044-L2051), [replacement](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1031-L1063), [null](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187). Pin проверен командой исследования 2026-11-01.

List/query/expand/Delete/batch/templates, отдельные position mutations и transport features вне scope. Текущая работа — только DRAFT; Go-код, tests реализации, .env, живой API и commits не выполняются.

## Triage 2026-11-01

Пользователь поручил выполнить все 11 entity tickets SDK. Предложенные таблицы полей,
local validation и ограничения scope приняты к реализации этой командой. Общие зависимости resolved.
BASE: eba6f75154929790c5b9d2cb1454dd80783b14db; исходные изменения сохранены.

## Answer

Реализован `entity/move`: Move, отдельные CreateRequest/UpdateRequest,
Get/Create/Update и типизированный GetPositions, aliases общей формы позиции.
[Публичный контракт](../../../docs/move.md). JSON и HTTP проверки, required/readonly
scope и ограничения покрыты; общие типы не дублируются.

Полный race/vet/gofmt/diff-check PASS; независимые Standards 0 / Spec 0.
Upstream fixtures/provenance неизменны.
[Интеграционный отчёт](../../document-workstream/review-entities.md).
