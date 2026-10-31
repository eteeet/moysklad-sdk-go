Status: needs-triage
Stage: DRAFT
Type: task
Blocked by: ../../document-common/issues/01-rate.md; ../../document-common/issues/02-overhead.md; ../../document-positions/issues/01-json-contracts.md; ../../document-positions/issues/02-get-page.md

# 01: Supply — Get/Create/Update

## What to build

Разработчик решения читает, создаёт и частично изменяет Supply (приёмка) одним запросом на метод, получает фактический ответ сервера и передаёт embedded positions. Законченный срез включает публичную модель/requests/Service, JSON/HTTP-проверки и документацию. Это proposed контракт для согласования, не команда реализации. [Spec](../spec.md).

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
| positions | *model.Collection[Position] | optional Field[[]positions.CommercialRequest] | Field[[]positions.CommercialRequest] |
| syncId | string | optional Field[string] | Field[string], write-once проверяет сервер |
| sum, payedSum, vatSum | float64 | — | — |
| printed, published | bool | — | — |
| files | *model.Collection[Reference], только metadata | — | — |
| store | *Reference | обязательный Reference | Field[Reference] |
| incomingDate; incomingNumber | DateTime; string | optional Field[DateTime]; Field[string] | те же Field |
| overhead | *DocumentOverhead | Field[DocumentOverheadRequest] | Field[DocumentOverheadRequest] |
| purchaseOrder | *Reference | proposed Field[Reference] | proposed Field[Reference] |
| factureIn, productionTask | *Reference | TODO write | TODO write |
| invoicesIn, payments, returns | []Reference | TODO write | TODO write |

DocumentRate/DocumentRateRequest и DocumentOverhead/DocumentOverheadRequest — proposed общие экспорты соответствующих dependencies; не копировать их в пакет сущности. Go names: ID/AccountID/SyncID/VATEnabled/VATIncluded/VATSum. Supply Position в entity package встраивает positions.Commercial и добавляет read-only Overhead float64 (API Int, копейки). Position write использует positions.CommercialRequest без overhead.

Минимальный Create содержит только organization, agent, store. Required-at-response не означает required-at-Create. Встроенная новая строка требует assortment и положительный quantity по API; в write shape это Field, обязательность зависит от нового/matching элемента. Readonly meta/accountId/id (кроме matching ID строки), вычисленные суммы/даты и flags не попадают в requests.

## Поведение и ограничения

- Get — GET entity/supply/{id}; Create — POST entity/supply; Update — PUT entity/supply/{id}; результаты (*Supply, error). NewService использует transport.Client и отклоняет nil.
- Get без query/body делает ровно один запрос и читает positions.meta без rows. Не добавляет expand и не переходит по href. [Общий reader](../../document-positions/issues/02-get-page.md) — отдельная явная операция, готовность полного чтения группы проверяется после него.
- Update отправляет только заданные поля. При смене organization/agent нужны organizationAccount/agentAccount соответственно. SDK не узнаёт старые значения, не заполняет и не проверяет равенство ссылок.
- Proposed local validation: UUID по UOM, Update без заданных Field, непустой Meta.Href обязательных Create References и не более 1000 embedded positions. Типы/доступность ссылок, matching account pairs и налоговые/складские ограничения проверяет сервер; SDK не загружает прежнее состояние.
- Field unset отсутствует; Value(false/0/""/[]) сохраняется. Null кодируется как null, не означает автоматического server clearing каждого поля. contract:null удаляет необязательную связь; rate:null сбрасывает валюту на default с пересчётом.
- Positions omitted не изменяет набор; [] очищает по правилу полной замены; matching ID обновляется, отсутствующий удаляется, новый добавляется. Это server replacement, без local merge. positions:null не интерпретируется как []; его server-смысл не подтверждён.
- Embedded write — до 1000 строк по API. Отдельные POST/PUT/DELETE позиций и обход всех страниц не входят в этот ticket. Attributes не наследует replacement-rule positions.
- overhead документа использует общий model.DocumentOverhead/DocumentOverheadRequest: sum float64, открытое строковое distribution с известными weight/volume/price; запись требует наличия positions по API. SDK не читает документ для проверки.
- purchaseOrder — proposed записываемая Reference. Не добавлять dependency на готовый PurchaseOrder; factureIn/productionTask и массивы invoicesIn/payments/returns в этом срезе только читаются.
- Assortment: product/service/variant/consignment References. country/gtd/pack/slot/things/trackingCodes — explicit TODO read/write; серийный/маркированный сценарий не считается поддержанным.

## Типизированное чтение позиций

Entity Service предоставляет GetPositions(ctx context.Context, id string, page positions.PageOptions)
(*model.Collection[Position], error). DocumentType и тип строки фиксированы пакетом; вызывающий
код не передаёт строку ресурса и type parameter. Фиксированный documentType — supply.
Один GET через общий reader, без expand,
follow-up и auto-pagination. Ошибки, UUID/page validation и response envelope — как у P02.
Отдельные TestServiceGetPositions с таблицами проверяют фиксированный resource path, правильный
Position, limit/offset, plain/gzip, пустую страницу, ошибки и ровно один запрос через настоящий transport.
Position расширяет positions.Commercial полем Overhead float64.
Технический blocker P02 обязателен. Реализация сущности остаётся needs-triage до отдельного разрешения.

## Acceptance criteria

- [ ] До реализации согласованы таблица полей, proposed exports и открытые вопросы ниже; Status меняется отдельным triage, DRAFT не считается ready.
- [ ] Dependencies resolved; общий JSON/DocumentRate/Overhead не дублируется. Нет зависимости от готовой другой entity: связи задаются Reference.
- [ ] Публичные model/CreateRequest/UpdateRequest и три метода полностью покрывают таблицу; readonly/TODO отсутствуют в requests. Сохраняются absence/null/false/0/""/[] без нормализации.
- [ ] JSON-тесты через encoding/json и encoding/json/v2 проверяют значения каждого поддержанного поля неизменённой upstream fixture, отсутствующие/null/unknown поля, даты, References и metadata positions без rows.
- [ ] Независимые golden покрывают минимальный Create, полный Create/Update, Update одного поля, paired accounts, false/0/""/[], contract:null и rate:null, omitted/[] positions и matching ID. Дробные quantity/price, отрицательный discount, vat=0 с vatEnabled=false/true не теряются.
- [ ] Proposed validation проверена без HTTP при пустых обязательных Meta.Href и 1001 позиции; 1000 позиций допускаются, positions:null/omitted не превращаются в [] или local merge.
- [ ] HTTP-тесты строго по структуре entity/uom/service_test.go: отдельные TestServiceGet/Create/Update с table cases; реальный transport.Client + httptest.Server, точные method/path/query/headers/body, plain/gzip, один запрос, серверный ответ отличается от input.
- [ ] Каждый метод проверяет context cancellation, auth/transport/APIError 4xx/5xx/429, errors.As/Is, отсутствие double moysklad:, неожиданные 201/204, malformed/null JSON, пустой/отсутствующий id. Успех — 200 с id, ошибка — nil-модель, без retry.
- [ ] Документация сущности и coverage перечисляют реально поддержанные поля, TODO, minimum Create, account coupling, full replacement/1000 и разделение Get/GetPage; пример чтения документа не обещает rows.
- [ ] Пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata, git diff --check и предусмотренное репозиторием review; исходные чужие/плановые изменения сохранены отдельно.

## Blocked by

- [01-rate.md](../../document-common/issues/01-rate.md)
- [02-overhead.md](../../document-common/issues/02-overhead.md)
- [01-json-contracts.md](../../document-positions/issues/01-json-contracts.md)

[GetPage](../../document-positions/issues/02-get-page.md) — обязательная зависимость этого entity ticket: типизированный GetPositions вызывает общий reader. Ни Supply, ни Demand, ни PurchaseOrder, ни services accounts/складов/assortment не блокируют этот ticket.

## Explicit Open Decisions

Утвердить write purchaseOrder, выбранный scope links и overhead; SyncID в Create/Update; границу required-reference/account-pair validation. Также подтвердить предложенные JSON/HTTP seams и local validation. До решения спорные поля не добавляются по аналогии.

## Sources и границы

- [Официальные поля и связи](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L8-L64), [минимальный Create](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L452-L460), [Update](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L2191-L2196), [positions](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L66-L96).
- [Null и адресные исключения](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187), [positions replacement](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1031-L1063).

Revision 9a811cef3dc4a3f3a46f070e60452361a06bf7cd проверен командой исследования 2026-11-01. List/query/expand/Delete/batch/templates, дополнительные position mutations и transport features — вне scope. Go-код, tests реализации, .env, живой API и commits сейчас не выполняются. Завершение DRAFT — сохранённый reviewable ticket с явными blockers/open decisions.
