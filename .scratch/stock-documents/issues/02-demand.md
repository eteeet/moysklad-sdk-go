Status: resolved
Stage: IMPLEMENTED
Type: task
Blocked by: ../../document-common/issues/01-rate.md; ../../document-common/issues/02-overhead.md; ../../document-positions/issues/01-json-contracts.md; ../../document-positions/issues/02-get-page.md

> Реализовано по команде пользователя выполнить все 11 tickets. Предложения и
> запреты запуска в исходном DRAFT ниже отражают этап подготовки; команда исполнения
> их заменила. Итоговый контракт — docs/demand.md, результат — в Answer.

# 02: Demand — Get/Create/Update

## What to build

Разработчик решения читает, создаёт и частично изменяет Demand (отгрузка) одним запросом на метод, получает фактический ответ сервера и передаёт embedded positions. Законченный срез включает публичную модель/requests/Service, JSON/HTTP-проверки и документацию. Это proposed контракт для согласования, не команда реализации. [Spec](../spec.md).

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
| positions | *model.Collection[Position] | optional Field[[]positions.PositionRequest] | Field[[]positions.PositionRequest] |
| syncId | string | optional Field[string] | Field[string], write-once проверяет сервер |
| sum, payedSum, vatSum | float64 | — | — |
| printed, published | bool | — | — |
| files | *model.Collection[Reference], только metadata | — | — |
| store | *Reference | optional Field[Reference] | Field[Reference] |
| overhead | *DocumentOverhead | Field[DocumentOverheadRequest] | Field[DocumentOverheadRequest] |
| salesChannel, customerOrder | *Reference | proposed Field[Reference] | proposed Field[Reference] |
| shipmentAddress | string | Field[string] | Field[string] |
| shipmentAddressFull | *model.Address | Field[model.AddressRequest] | Field[model.AddressRequest] |
| factureOut | *Reference | TODO write | TODO write |
| returns, payments, productionTasks, invoicesOut | []Reference | TODO write | TODO write |

DocumentRate/DocumentRateRequest и DocumentOverhead/DocumentOverheadRequest — proposed общие экспорты соответствующих dependencies; не копировать их в пакет сущности. Go names: ID/AccountID/SyncID/VATEnabled/VATIncluded/VATSum. Demand использует aliases Position = positions.Position и PositionRequest = positions.PositionRequest. S02 добавит Cost float64 в копейках в общую positions.Position и Cost model.Field[float64] с omitzero в общий positions.PositionRequest. Cost — себестоимость только для услуг; применимость проверяет сервер без скрытого Get или локального определения типа assortment. Overhead уже находится в общей Position и отсутствует в request вместе с reserve. Собственная форма строки не создаётся; до S02 Cost в коде не реализован.

Proposed обязательные поля Create — organization/agent. Безопасный пример по разделу Create дополнительно содержит store через model.Value; отдельный wire case с applicable=false допускает отсутствие store до отправки HTTP. Его принятие живым API не проверено. Required-at-response не означает required-at-Create. Встроенная новая строка требует assortment и положительный quantity по API; в write shape это Field, обязательность зависит от нового/matching элемента. Readonly meta/accountId/id (кроме matching ID строки), вычисленные суммы/даты и flags не попадают в requests.

## Поведение и ограничения

- Get — GET entity/demand/{id}; Create — POST entity/demand; Update — PUT entity/demand/{id}; результаты (*Demand, error). NewService использует transport.Client и отклоняет nil.
- Get без query/body делает ровно один запрос и читает positions.meta без rows. Не добавляет expand и не переходит по href. [Общий reader](../../document-positions/issues/02-get-page.md) — отдельная явная операция, готовность полного чтения группы проверяется после него.
- Update отправляет только заданные поля. При смене organization/agent нужны organizationAccount/agentAccount соответственно. SDK не узнаёт старые значения, не заполняет и не проверяет равенство ссылок.
- Proposed local validation: UUID по UOM, Update без заданных Field, непустой Meta.Href Create.Organization/Create.Agent и не более 1000 embedded positions. Store не проверяется локально как required. Типы/доступность ссылок, matching account pairs и налоговые/складские ограничения проверяет сервер; SDK не загружает прежнее состояние.
- Field unset отсутствует; Value(false/0/""/[]) сохраняется. Null кодируется как null, не означает автоматического server clearing каждого поля. contract:null удаляет необязательную связь; rate:null сбрасывает валюту на default с пересчётом.
- Positions omitted не изменяет набор; [] очищает по правилу полной замены; matching ID обновляется, отсутствующий удаляется, новый добавляется. Это server replacement, без local merge. positions:null не интерпретируется как []; его server-смысл не подтверждён.
- Embedded write — до 1000 строк по API. Отдельные POST/PUT/DELETE позиций и обход всех страниц не входят в этот ticket. Attributes не наследует replacement-rule positions.
- store противоречив: раздел Create перечисляет его обязательным; строка таблицы требует его при applicable=true. Proposed Field[Reference] не запрещает пропуск; первый пример использует store, а условную обязательность оставляем серверу. Успех Create без store при applicable=false живым API не подтверждён.
- overhead документа — общий DocumentOverhead/DocumentOverheadRequest; sum float64, distribution weight/volume/price; необходимость positions остаётся server constraint.
- shipmentAddressFull использует existing Address/AddressRequest. Официальные части: addInfo/apartment/city/comment/country/house/postalCode/region/street. Дополнительный FIASCodeRU существующего типа не обещает поддержки Demand.
- При двух адресах server предпочитает структурированный; только строковый отражается в addInfo. shipmentAddressFull:null не очищает, shipmentAddress:"" очищает. SDK не синхронизирует две формы.
- Assortment: product/service/variant/consignment/bundle References. Reserve есть в примере, но отсутствует в таблице — TODO. pack/slot/things/trackingCodes/trackingCodes_1162/declaration и stock — TODO read/write.
- Транспортные поля cargoName/carrier/consignee/consigneeBranch/consignor/consignorBranch/goodPackQuantity/shippingInstructions/stateContractId/transportFacility/transportFacilityNumber/shippingDocumentNumber/shippingDocumentDate — explicit TODO read/write.

## Типизированное чтение позиций

Entity Service предоставляет GetPositions(ctx context.Context, id string, page positions.PageOptions)
(*model.Collection[Position], error). DocumentType и тип строки фиксированы пакетом; вызывающий
код не передаёт строку ресурса и type parameter. Фиксированный documentType — demand.
Один GET через общий reader, без expand,
follow-up и auto-pagination. Ошибки, UUID/page validation и response envelope — как у P02.
Отдельные TestServiceGetPositions с таблицами проверяют фиксированный resource path, правильный
Position, limit/offset, plain/gzip, пустую страницу, ошибки и ровно один запрос через настоящий transport.
Position — alias positions.Position. S02 добавляет Cost float64 в общую Position и Cost Field[float64] с omitzero в общий PositionRequest, не в entity-owned тип.
Технический blocker P02 обязателен. Реализация сущности остаётся needs-triage до отдельного разрешения.

## Acceptance criteria

- [x] До реализации согласованы таблица полей, proposed exports и открытые вопросы ниже; Status меняется отдельным triage, DRAFT не считается ready.
- [x] Dependencies resolved; общий JSON/DocumentRate/Overhead не дублируется. Нет зависимости от готовой другой entity: связи задаются Reference.
- [x] Публичные model/CreateRequest/UpdateRequest и три метода полностью покрывают таблицу; readonly/TODO отсутствуют в requests. Сохраняются absence/null/false/0/""/[] без нормализации.
- [x] JSON-тесты через encoding/json и encoding/json/v2 проверяют значения каждого поддержанного поля неизменённой upstream fixture, отсутствующие/null/unknown поля, даты, References и metadata positions без rows.
- [x] Независимые golden покрывают минимальный Create, полный Create/Update, Update одного поля, paired accounts, false/0/""/[], contract:null и rate:null, omitted/[] positions и matching ID. Дробные quantity/price, отрицательный discount, vat=0 с vatEnabled=false/true не теряются.
- [x] Cost в общей Position/PositionRequest проверен обоими JSON-кодеками: чтение 0/0.0/дроби, запись unset/0/дроби/null. HTTP body для позиции услуги сохраняет заданный Cost, omission не подставляет default; один запрос без Get для проверки типа assortment. Null кодируется, но не объявляется поддерживаемой очисткой сервером.
- [x] Proposed validation проверена без HTTP при пустых обязательных Meta.Href и 1001 позиции; 1000 позиций допускаются, positions:null/omitted не превращаются в [] или local merge.
- [x] Wire case с applicable=false без store доходит до HTTP с точным телом без SDK defaults. Отдельный safe Create golden содержит store. Эти проверки не доказывают принятие сценария настоящим API.
- [x] HTTP-тесты строго по структуре entity/uom/service_test.go: отдельные TestServiceGet/Create/Update с table cases; реальный transport.Client + httptest.Server, точные method/path/query/headers/body, plain/gzip, один запрос, серверный ответ отличается от input.
- [x] Каждый метод проверяет context cancellation, auth/transport/APIError 4xx/5xx/429, errors.As/Is, отсутствие double moysklad:, неожиданные 201/204, malformed/null JSON, пустой/отсутствующий id. Успех — 200 с id, ошибка — nil-модель, без retry.
- [x] Документация сущности и coverage перечисляют реально поддержанные поля, TODO, minimum Create, account coupling, full replacement/1000 и разделение Get/GetPage; пример чтения документа не обещает rows.
- [x] Пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata, git diff --check и предусмотренное репозиторием review; исходные чужие/плановые изменения сохранены отдельно.

## Blocked by

- [01-rate.md](../../document-common/issues/01-rate.md)
- [02-overhead.md](../../document-common/issues/02-overhead.md)
- [01-json-contracts.md](../../document-positions/issues/01-json-contracts.md)

[GetPage](../../document-positions/issues/02-get-page.md) — обязательная зависимость этого entity ticket: типизированный GetPositions вызывает общий reader. Ни Supply, ни Demand, ни PurchaseOrder, ни services accounts/складов/assortment не блокируют этот ticket.

## Explicit Open Decisions

Утвердить store при applicable=false, reuse Address без promise о FIASCodeRU, write salesChannel/customerOrder, TODO reserve и транспортных полей; SyncID в Create/Update. Также подтвердить предложенные JSON/HTTP seams и local validation. До решения спорные поля не добавляются по аналогии.

## Sources и границы

- [Cost: таблица, только для услуг](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L94), [embedded Create с Cost](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L1263-L1282), [Create позиции услуги с Cost](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L3855-L3876). Cost не помечен readonly; read-only сужение отменено до запуска S02.

- [Официальные поля и связи](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L11-L84), [минимальный Create](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L557-L565), [Update](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L2545-L2550), [positions](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L86-L117).
- [Null и адресные исключения](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187), [positions replacement](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1031-L1063).
- [Адрес Demand](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L178-L197), [reserve в примере](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L1184-L1249).

Revision 9a811cef3dc4a3f3a46f070e60452361a06bf7cd проверен командой исследования 2026-11-01. List/query/expand/Delete/batch/templates, дополнительные position mutations и transport features — вне scope. Go-код, tests реализации, .env, живой API и commits сейчас не выполняются. Завершение DRAFT — сохранённый reviewable ticket с явными blockers/open decisions.

## Triage 2026-11-01

Пользователь поручил выполнить все 11 entity tickets SDK. Предложенные таблицы полей,
local validation и ограничения scope приняты к реализации этой командой. Общие зависимости resolved.
BASE: eba6f75154929790c5b9d2cb1454dd80783b14db; исходные изменения сохранены.

## Answer

Реализован `entity/demand`: Demand, отдельные CreateRequest/UpdateRequest,
Get/Create/Update и типизированный GetPositions, aliases общей формы позиции.
[Публичный контракт](../../../docs/demand.md). JSON и HTTP проверки, required/readonly
scope и ограничения покрыты; общие типы не дублируются.
Cost доступен в общей PositionRequest для услуг; unset/0/дробь/null передаются без Get.

Полный race/vet/gofmt/diff-check PASS; независимые Standards 0 / Spec 0.
Upstream fixtures/provenance неизменны.
[Интеграционный отчёт](../../document-workstream/review-entities.md).
