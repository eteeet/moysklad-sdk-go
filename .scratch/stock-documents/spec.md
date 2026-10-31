Status: needs-triage
Stage: DRAFT

# Supply, Demand, PurchaseOrder — Get/Create/Update

## Problem Statement

Разработчику решения нужны чтение, создание и частичное изменение Supply (приёмка), Demand (отгрузка) и PurchaseOrder (заказ поставщику). Сейчас SDK покрывает справочники, но не документы и их positions. Сходство документов не определяет обязательность склада, права изменения связей или строк.

## Solution

Предлагаются три самостоятельных законченных среза: модель чтения, отдельные CreateRequest/UpdateRequest, Get/Create/Update, публичные JSON-контракты, локальные HTTP-тесты и документация каждого документа. Все решения ниже proposed до отдельного согласования. Этот DRAFT не разрешает реализацию.

Общий модуль positions отвечает за формы коллекций, коммерческие строки и отдельное постраничное чтение. Обычный Get документа получает metadata positions одним запросом; строки читаются явным GetPage общего reader. Expand отложен до окончания остальных возможностей SDK.

## User Stories

1. Как разработчик решения, я хочу получить документ по UUID, чтобы прочитать его состояние и ссылки.
2. Как разработчик решения, я хочу отличать metadata positions от загруженных строк, чтобы не принять неполный документ за полный.
3. Как разработчик решения, я хочу читать страницы positions отдельно, чтобы обработать документы больше 1000 строк.
4. Как разработчик решения, я хочу создать Supply с Organization, Counterparty и складом.
5. Как разработчик решения, я хочу передать Demand со складом либо явно запросить непроведённый документ без него, не встречая безусловного локального запрета SDK.
6. Как разработчик решения, я хочу создать PurchaseOrder без обязательного склада.
7. Как разработчик решения, я хочу задать позиции при создании документа, чтобы записать документ одним запросом.
8. Как разработчик решения, я хочу изменить одно поле, сохранив остальные без предварительного Get.
9. Как разработчик решения, я хочу передавать false, 0, пустые строки, null и [], сохранив их смысл.
10. Как разработчик решения, я хочу явно заменить все positions и понимать риск удаления отсутствующих строк.
11. Как разработчик решения, я хочу использовать Reference на ещё не реализованные справочники и документы.
12. Как разработчик решения, я хочу видеть фактический ответ сервера и разбирать APIError через errors.As.

## Implementation Decisions

### Предлагаемый общий контракт

Маршруты для каждого `entity/{supply|demand|purchaseorder}`: Get — GET `/{id}`; Create — POST на коллекцию; Update — PUT `/{id}`. Результат — указатель на модель и error. Service использует существующий transport.Client, NewService отклоняет nil. Успех — HTTP 200 и непустой id; модель при ошибке nil. Get не имеет query/body, не делает follow-up. Ошибки и UUID-проверка соответствуют образцу UOM; новые универсальные механизмы транспорта не вводятся.

Таблица — полный предлагаемый core всех трёх документов. R — чтение; C — Create; U — Update. В requests необязательные значения представлены model.Field с omitzero; обязательные Reference — обычные значения. Метки API «обязательно при ответе» не становятся обязательностью Create. Nullable-контракт описан отдельно.

| Поля JSON | R: тип Go | C | U |
| --- | --- | --- | --- |
| meta; id, accountId | model.Meta; string | — | — |
| created, updated, deleted | model.DateTime | — | — |
| name, code, externalCode, description | string | optional Field[string] | Field[string] |
| moment | model.DateTime | optional Field[DateTime] | Field[DateTime] |
| applicable, shared, vatEnabled, vatIncluded | bool | optional Field[bool] | Field[bool] |
| organization, agent | *model.Reference | обязательные Reference | Field[Reference], с парным account |
| organizationAccount, agentAccount | *model.Reference | optional Field[Reference] | Field[Reference] |
| owner, group, project, contract, state | *model.Reference | optional Field[Reference] | Field[Reference] |
| attributes | []model.Attribute | optional Field[[]Attribute] | Field[[]Attribute] |
| rate | *model.DocumentRate из общего ticket | optional Field[DocumentRateRequest] | Field[DocumentRateRequest] |
| positions | *model.Collection[Position] | optional Field[[]PositionRequest] | Field[[]PositionRequest], полная замена |
| syncId | string | optional Field[string] | Field[string], write-once проверяет сервер |
| sum, payedSum, vatSum | float64 | — | — |
| printed, published | bool | — | — |
| files | *model.Collection[model.Reference], только meta | — | — |

Имена Go следуют действующим соглашениям SDK: ID, AccountID, SyncID, VATEnabled, VATIncluded, VATSum. DocumentRate/DocumentRateRequest — proposed имена общего контракта; итоговые экспорты определяет его ticket. Для JSON DateTime/Reference/Attribute применяются существующие типы. Файлы читаются только как metadata, без promise о rows.

### Поля и правила конкретных документов

| Документ / JSON-поле | R | C / U | Источник / ограничение |
| --- | --- | --- | --- |
| Supply: store | *Reference | C: обязательный Reference; U: Field[Reference] | [Create supply](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L452-L460) |
| Supply: incomingDate, incomingNumber | DateTime, string | optional Field соответствующего типа в C/U | [таблица supply](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L8-L45) |
| Supply: overhead | *DocumentOverhead | Field[DocumentOverheadRequest] в C/U | sum float64, distribution weight/volume/price; без positions задавать нельзя |
| Supply: purchaseOrder | *Reference | proposed Field[Reference] в C/U | отдельная связь, не prerequisite реализации PurchaseOrder |
| Supply: factureIn, productionTask | *Reference | TODO write | таблица связей не доказывает универсальную запись |
| Supply: invoicesIn, payments, returns | []Reference | TODO write | [связи](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L47-L64) |
| Demand: store | *Reference | C/U: Field[Reference] | proposed: условную обязательность оставляем серверу; Create и таблица расходятся |
| Demand: overhead | *DocumentOverhead | Field[DocumentOverheadRequest] в C/U | общее overhead, наличие positions — server constraint |
| Demand: salesChannel, customerOrder | *Reference | proposed Field[Reference] в C/U | остальные связи не выводятся из этих двух |
| Demand: shipmentAddress, shipmentAddressFull | string, *model.Address | Field[string], Field[model.AddressRequest] в C/U | 9 официальных частей; FIASCodeRU не обещается |
| Demand: factureOut | *Reference | TODO write | read only в предлагаемом scope |
| Demand: returns, payments, productionTasks, invoicesOut | []Reference | TODO write | [таблицы Demand](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L11-L84) |
| PurchaseOrder: store | *Reference | optional Field[Reference] в C/U | необязателен при Create |
| PurchaseOrder: deliveryPlannedMoment | DateTime | Field[DateTime] в C/U | [таблица PurchaseOrder](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L9-L58) |
| PurchaseOrder: invoicedSum, shippedSum, waitSum | float64 | — | waitSum proposed read-only; таблица не маркирует его read-only |
| PurchaseOrder: internalOrder | *Reference | proposed Field[Reference] в C/U | запись связи требует утверждения |
| PurchaseOrder: customerOrders, invoicesIn, productionTasks, payments, supplies | []Reference | TODO write | читаются ссылки; без автоматической загрузки |

Минимальный Create Supply — organization, agent, store; PurchaseOrder — organization, agent. У Demand proposed обязательные Create поля — organization/agent; store — Field[Reference], безусловный локальный запрет пропуска не вводится. Безопасный пример Demand по разделу Create содержит store; принятие сервером applicable=false без store живым API не проверено. name и positions не обязательны. Смена organization/agent требует organizationAccount/agentAccount соответственно; SDK передаёт пару, но не выполняет предварительный Get и не сравнивает прежнее значение. Proposed: локально проверять UUID, полностью незаданный Update, непустой Meta.Href обязательных Create References и предел 1000 embedded positions; типы/доступность ссылок, условную обязательность Demand.Store и пары счетов проверяет сервер. [Demand Create](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L557-L565), [PurchaseOrder Create](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L356-L363), [парные accounts](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L2191-L2196).

### Positions и значения записи

Position документов использует proposed positions.Commercial: core meta/id/accountId, assortment Reference, quantity/price/discount float64, vat int, vatEnabled bool. Supply добавляет read-only overhead float64; Demand — read-only cost и overhead float64; PurchaseOrder — read shipped/inTransit float64. Supply/Demand используют positions.CommercialRequest. PurchaseOrder.PositionRequest встраивает CommercialRequest и добавляет InTransit Field[float64] с omitzero; поле документное, не добавляется в общий CommercialRequest. Overhead/cost/shipped/reserve в записи отсутствуют. Все поля запроса строки — Field, ID служит только matching при полной замене. Required новой строки — assortment и положительный quantity по API; response-required price/discount/vat не становятся required Create. Запись inTransit подтверждается embedded Create и отдельными примерами; readonly shipped — выбранный scope SDK, не readonly marker таблицы.

Assortment поддерживает Reference на product/service/variant/consignment; Demand также bundle. Полные объекты Product и остальных справочников не загружаются. Country/gtd/pack/slot/things/trackingCodes/trackingCodes_1162/declaration и stock enrichment — explicit TODO; эти поля не присутствуют в предлагаемых request types. Общие позиции не принадлежат entity tickets. [Supply positions](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L66-L96), [Demand](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L86-L117), [PurchaseOrder](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L60-L86).

Omitted field не отправляется; false/0/"" не исчезают. Null сохраняется как JSON-состояние Field, но server-допустимость не предполагается для каждого поля: contract:null удаляет необязательную связь; rate:null возвращает валюту по умолчанию с пересчётом. У Demand shipmentAddressFull:null не очищает адрес, shipmentAddress:"" очищает; при двух формах приоритет у структурированной. [Null](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187), [адрес](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L178-L197).

Positions omitted сохраняет прежний набор; positions:[] очищает по правилу full replacement (вывод из partial Update + replacement). Совпавший ID обновляется, отсутствующий удаляется, новый добавляется; локального merge нет. positions:null не приравнивается к []; смысл null для коллекции не подтверждён. Embedded write ограничен API 1000 строками; отдельные mutations /positions — другой scope. Attributes использует действующий Attribute-контракт, а не правило replacement positions. [Общее правило](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1031-L1063).

### Зависимости и ownership

Supply/Demand/PurchaseOrder зависят от [общего DocumentRate](../document-common/issues/01-rate.md) и [positions JSON](../document-positions/issues/01-json-contracts.md). Supply/Demand дополнительно зависят от [DocumentOverhead](../document-common/issues/02-overhead.md). Между тремя entity tickets зависимостей нет. [GetPage](../document-positions/issues/02-get-page.md) — обязательная зависимость entity tickets с типизированным GetPositions. Итоговая группа проверяется совместно после всех этих tickets.

## Testing Decisions

Два публичных seam: JSON encoding/decoding через encoding/json и encoding/json/v2; Service через реальный transport.Client и локальный httptest.Server. Структура HTTP-тестов обязательна по `entity/uom/service_test.go`: отдельные TestServiceGet/Create/Update с табличными сценариями. Самодостаточные entity tickets уточняют проверки отличий.

Проверяются значения всех поддержанных полей, metadata без rows, отсутствующие/unknown/null поля, дробные quantity/price, отрицательный discount; независимые golden минимальных и полных requests, partial Update, false/0/""/[]/null, отсутствие readonly полей. HTTP-тесты фиксируют метод/путь/query/body/заголовки, plain/gzip, ровно один запрос, серверный ответ, ошибки context/auth/transport/APIError, errors.Is/As, malformed/null JSON, неожиданный status и отсутствующий ID. Это проверка SDK, не доказательство поведения живого API.

## Out of Scope

List документов, query/filter/search, expand, Delete, batch, metadata/templates, файлы кроме metadata, отдельные mutations позиций, serial/tracking/stock enrichment, транспортные retry/rate limiting и новые services связанных сущностей. Demand transport fields cargoName/carrier/consignee/consigneeBranch/consignor/consignorBranch/goodPackQuantity/shippingInstructions/stateContractId/transportFacility/transportFacilityNumber/shippingDocumentNumber/shippingDocumentDate — TODO полностью, read и write.

## Further Notes

Официальный revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd` повторно проверен командой исследования 2026-11-01 через GitHub API. [Research](research.md) фиксирует первичные источники и неопределённости. Go-код, тесты реализации, .env, живой API и commits этим документом не разрешены.

### Explicit Open Decisions

1. Утвердить таблицы scope, proposed имена экспортов и локальную валидацию обязательных ссылок.
2. Demand store: proposed Field[Reference] в Create/Update, без локального required-validator; условие applicable проверяет сервер. Расхождение документации сохраняется, live-проверка не проведена.
3. PurchaseOrder: proposed запись inTransit через документный PositionRequest; shipped и waitSum остаются read-only в scope SDK, хотя readonly marker таблиц нет.
4. Demand reserve отсутствует в таблице, но есть в примерах; proposed TODO.
5. Утвердить write support purchaseOrder/customerOrder/internalOrder/salesChannel и оставление остальных связей read-only в scope.
6. Подтвердить reuse Address/AddressRequest без promise о demand FIASCodeRU и typed overhead из общего ticket.
7. SyncID: proposed Field[string] в Create/Update; API запрещает изменение после заполнения, проверку оставляем серверу.
8. Утвердить выбранные JSON/HTTP seams; получить отдельную команду реализации после согласования.

## Типизированное чтение позиций — согласованное уточнение

У каждого entity Service добавить GetPositions(ctx context.Context, id string, page positions.PageOptions)
(*model.Collection[Position], error). Он вызывает общий reader с фиксированными documentType
и Position; вызывающий код не выбирает строку ресурса и generic type. Один HTTP GET, без expand,
автообхода и скрытого Get заголовка. Ошибки/envelope/UUID/page validation остаются контрактом P02.
Move экспортирует alias Position = positions.Movement; InvoiceIn/InvoiceOut/PurchaseReturn —
alias Position = positions.Commercial. Supply/Demand/PurchaseOrder используют собственные Position.
Тестовый seam — публичный GetPositions через transport.Client + httptest.Server;
отдельный TestServiceGetPositions с таблицами правильного path/type/query, plain/gzip, пустой страницы,
ошибок и одной отправки. P02 — обязательная зависимость всего entity ticket.
Это уточнение не разрешает реализацию entity, остальные open decisions остаются needs-triage.
