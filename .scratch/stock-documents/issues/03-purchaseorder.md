Status: resolved
Stage: IMPLEMENTED
Type: task
Blocked by: ../../document-common/issues/01-rate.md; ../../document-positions/issues/01-json-contracts.md; ../../document-positions/issues/02-get-page.md

> Реализовано по команде пользователя выполнить все 11 tickets. Предложения и
> запреты запуска в исходном DRAFT ниже отражают этап подготовки; команда исполнения
> их заменила. Итоговый контракт — docs/purchaseorder.md, результат — в Answer.

# 03: PurchaseOrder — Get/Create/Update

## What to build

Разработчик решения читает, создаёт и частично изменяет PurchaseOrder (заказ поставщику) одним запросом на метод, получает фактический ответ сервера и передаёт embedded positions. Законченный срез включает публичную модель/requests/Service, JSON/HTTP-проверки и документацию. Это proposed контракт для согласования, не команда реализации. [Spec](../spec.md).

## Предлагаемый точный scope

R — модель чтения; C/U — отдельные requests. Ниже полный core и конкретные поля. DateTime/Reference/Attribute — существующие model-типы; optional read Reference — указатель. Каждое необязательное поле записи — model.Field с omitzero. Не перечисленные поля в requests отсутствуют.

| JSON-поле | R: Go type | C | U |
| --- | --- | --- | --- |
| meta; id, accountId | Meta; string | — | — |
| created, updated, deleted | DateTime | — | — |
| name, code, externalCode, description | string | optional Field[string] | Field[string] |
| moment | DateTime | optional Field[DateTime] | Field[DateTime] |
| applicable, shared, vatEnabled, vatIncluded | bool | optional Field[bool] | Field[bool] |
| organization, agent | *Reference | обязательные Reference | Field[Reference], с парным account |
| organizationAccount, agentAccount | *Reference | optional Field[Reference] | Field[Reference] |
| owner, group, project, contract, state | *Reference | optional Field[Reference] | Field[Reference] |
| attributes | []Attribute | optional Field[[]Attribute] | Field[[]Attribute] |
| rate | *DocumentRate | optional Field[DocumentRateRequest] | Field[DocumentRateRequest] |
| positions | *model.Collection[Position] | optional Field[[]PositionRequest] | Field[[]PositionRequest] |
| syncId | string | optional Field[string] | Field[string], write-once проверяет сервер |
| sum, payedSum, vatSum | float64 | — | — |
| printed, published | bool | — | — |
| files | *model.Collection[Reference], только metadata | — | — |
| store | *Reference | optional Field[Reference] | Field[Reference] |
| deliveryPlannedMoment | DateTime | Field[DateTime] | Field[DateTime] |
| invoicedSum, shippedSum, waitSum | float64 | — | — |
| internalOrder | *Reference | proposed Field[Reference] | proposed Field[Reference] |
| customerOrders, invoicesIn, productionTasks, payments, supplies | []Reference | TODO write | TODO write |

DocumentRate/DocumentRateRequest — proposed общие экспорты dependency общей валюты; не копировать их в пакет сущности. Go names: ID/AccountID/SyncID/VATEnabled/VATIncluded/VATSum. PurchaseOrder использует aliases Position = positions.Position и PositionRequest = positions.PositionRequest. S03 добавит Shipped/InTransit float64 в общую Position, InTransit model.Field[float64] с omitzero — в общий PositionRequest. Shipped отсутствует в request как выбранный scope SDK. Собственная форма строки не создаётся; сейчас этих полей в коде нет.

Минимальный Create содержит только organization, agent. Required-at-response не означает required-at-Create. Встроенная новая строка требует assortment и положительный quantity по API; в write shape это Field, обязательность зависит от нового/matching элемента. Readonly meta/accountId/id (кроме matching ID строки), вычисленные суммы/даты и flags не попадают в requests.

## Поведение и ограничения

- Get — GET entity/purchaseorder/{id}; Create — POST entity/purchaseorder; Update — PUT entity/purchaseorder/{id}; результаты (*PurchaseOrder, error). NewService использует transport.Client и отклоняет nil.
- Get без query/body делает ровно один запрос и читает positions.meta без rows. Не добавляет expand и не переходит по href. [Общий reader](../../document-positions/issues/02-get-page.md) — отдельная явная операция, готовность полного чтения группы проверяется после него.
- Update отправляет только заданные поля. При смене organization/agent нужны organizationAccount/agentAccount соответственно. SDK не узнаёт старые значения, не заполняет и не проверяет равенство ссылок.
- Proposed local validation: UUID по UOM, Update без заданных Field, непустой Meta.Href обязательных Create References и не более 1000 embedded positions. Типы/доступность ссылок, matching account pairs и налоговые/складские ограничения проверяет сервер; SDK не загружает прежнее состояние.
- Field unset отсутствует; Value(false/0/""/[]) сохраняется. Null кодируется как null, не означает автоматического server clearing каждого поля. contract:null удаляет необязательную связь; rate:null сбрасывает валюту на default с пересчётом.
- Positions omitted не изменяет набор; [] очищает по правилу полной замены; matching ID обновляется, отсутствующий удаляется, новый добавляется. Это server replacement, без local merge. positions:null не интерпретируется как []; его server-смысл не подтверждён.
- Embedded write — до 1000 строк по API. Отдельные POST/PUT/DELETE позиций и обход всех страниц не входят в этот ticket. Attributes не наследует replacement-rule positions.
- store необязателен в Create; минимальный golden не должен изобретать склад, name или positions.
- deliveryPlannedMoment — DateTime чтения/записи; отсутствие поля отличается от переданного null на JSON seam, допустимость null оставлена серверу.
- shipped и inTransit читаются как float64. В таблице у них нет readonly marker. Запись inTransit показана в embedded Create и отдельных Create/Update position примерах: InTransit Field[float64] в общем positions.PositionRequest в S03. Shipped, возвращаемый 0 несмотря на переданные 20/10, proposed только чтение; не объявлять запретом API.
- waitSum не имеет read-only marker в таблице; proposed только чтение как производной суммы. Не заявлять это ограничение запретом API.
- Assortment: product/service/variant/consignment References; pack, serial/tracking и stock enrichment — TODO. internalOrder write proposed; массивы связей читаются без отдельной загрузки.

## Типизированное чтение позиций

Entity Service предоставляет GetPositions(ctx context.Context, id string, page positions.PageOptions)
(*model.Collection[Position], error). DocumentType и тип строки фиксированы пакетом; вызывающий
код не передаёт строку ресурса и type parameter. Фиксированный documentType — purchaseorder.
Один GET через общий reader, без expand,
follow-up и auto-pagination. Ошибки, UUID/page validation и response envelope — как у P02.
Отдельные TestServiceGetPositions с таблицами проверяют фиксированный resource path, правильный
Position, limit/offset, plain/gzip, пустую страницу, ошибки и ровно один запрос через настоящий transport.
Position — alias positions.Position. S03 добавляет Shipped/InTransit в общую Position и InTransit Field в общий PositionRequest.
Технический blocker P02 обязателен. Реализация сущности остаётся needs-triage до отдельного разрешения.

## Acceptance criteria

- [x] До реализации согласованы таблица полей, proposed exports и открытые вопросы ниже; Status меняется отдельным triage, DRAFT не считается ready.
- [x] Dependencies resolved; общий JSON/DocumentRate не дублируется. Нет зависимости от готовой другой entity: связи задаются Reference.
- [x] Публичные model/CreateRequest/UpdateRequest и три метода полностью покрывают таблицу; readonly/TODO отсутствуют в requests. Сохраняются absence/null/false/0/""/[] без нормализации.
- [x] JSON-тесты через encoding/json и encoding/json/v2 проверяют значения каждого поддержанного поля неизменённой upstream fixture, отсутствующие/null/unknown поля, даты, References и metadata positions без rows.
- [x] Независимые golden покрывают минимальный Create, полный Create/Update, Update одного поля, paired accounts, false/0/""/[], contract:null и rate:null, omitted/[] positions и matching ID. Дробные quantity/price, отрицательный discount, vat=0 с vatEnabled=false/true не теряются.
- [x] S03 добавил Shipped/InTransit в общую positions.Position и InTransit Field в общий PositionRequest, без собственной формы строки. Оба кодека и HTTP body проверяют omission/0/дробь/null и matching ID. Shipped отсутствует в request.
- [x] Proposed validation проверена без HTTP при пустых обязательных Meta.Href и 1001 позиции; 1000 позиций допускаются, positions:null/omitted не превращаются в [] или local merge.
- [x] HTTP-тесты строго по структуре entity/uom/service_test.go: отдельные TestServiceGet/Create/Update с table cases; реальный transport.Client + httptest.Server, точные method/path/query/headers/body, plain/gzip, один запрос, серверный ответ отличается от input.
- [x] Каждый метод проверяет context cancellation, auth/transport/APIError 4xx/5xx/429, errors.As/Is, отсутствие double moysklad:, неожиданные 201/204, malformed/null JSON, пустой/отсутствующий id. Успех — 200 с id, ошибка — nil-модель, без retry.
- [x] Документация сущности и coverage перечисляют реально поддержанные поля, TODO, minimum Create, account coupling, full replacement/1000 и разделение Get/GetPage; пример чтения документа не обещает rows.
- [x] Пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata, git diff --check и предусмотренное репозиторием review; исходные чужие/плановые изменения сохранены отдельно.

## Blocked by

- [01-rate.md](../../document-common/issues/01-rate.md)
- [01-json-contracts.md](../../document-positions/issues/01-json-contracts.md)

[GetPage](../../document-positions/issues/02-get-page.md) — обязательная зависимость этого entity ticket: типизированный GetPositions вызывает общий reader. Ни Supply, ни Demand, ни PurchaseOrder, ни services accounts/складов/assortment не блокируют этот ticket.

## Explicit Open Decisions

Общая форма строк утверждена; уточнить server-ограничения inTransit write; shipped и waitSum read-only scope; internalOrder write и scope массивов связей; SyncID в Create/Update. Также подтвердить предложенные JSON/HTTP seams и local validation. До решения спорные поля не добавляются по аналогии.

## Sources и границы

- [Официальные поля и связи](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L9-L58), [минимальный Create](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L356-L363), [Update](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L2145-L2150), [positions](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L60-L86).
- [Null и адресные исключения](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187), [positions replacement](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1031-L1063).
- [Неоднозначный shipped/inTransit write example](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L2846-L2910).
- [Embedded Create с inTransit](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L712-L758), [повторный triage фактов и предложений](../../document-workstream/triage-research.md#2-purchaseorder-shipped-intransit-waitsum).

Revision 9a811cef3dc4a3f3a46f070e60452361a06bf7cd проверен командой исследования 2026-11-01. List/query/expand/Delete/batch/templates, дополнительные position mutations и transport features — вне scope. Go-код, tests реализации, .env, живой API и commits сейчас не выполняются. Завершение DRAFT — сохранённый reviewable ticket с явными blockers/open decisions.

## Triage 2026-11-01

Пользователь поручил выполнить все 11 entity tickets SDK. Предложенные таблицы полей,
local validation и ограничения scope приняты к реализации этой командой. Общие зависимости resolved.
BASE: eba6f75154929790c5b9d2cb1454dd80783b14db; исходные изменения сохранены.

## Answer

Реализован `entity/purchaseorder`: PurchaseOrder, отдельные CreateRequest/UpdateRequest,
Get/Create/Update и типизированный GetPositions, aliases общей формы позиции.
[Публичный контракт](../../../docs/purchaseorder.md). JSON и HTTP проверки, required/readonly
scope и ограничения покрыты; общие типы не дублируются.
Shipped/InTransit читаются общей Position; InTransit доступен в общем request.

Полный race/vet/gofmt/diff-check PASS; независимые Standards 0 / Spec 0.
Upstream fixtures/provenance неизменны.
[Интеграционный отчёт](../../document-workstream/review-entities.md).
