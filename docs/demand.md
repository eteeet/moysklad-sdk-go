# Demand

`entity/demand`: модель `Demand`, `CreateRequest`, `UpdateRequest`,
`NewService(client)` и `Get/Create/Update`. Дополнительно `GetPositions` и aliases
`Position`/`PositionRequest` общих форм.
[Общие правила методов, JSON, ошибок и позиций](documents.md).

Минимальный Create: `organization`, `agent`.

## Поддержанные поля

| JSON | Read Go type | Create | Update |
| --- | --- | --- | --- |
| `id` | `string` | — | — |
| `accountId` | `string` | — | — |
| `name` | `string` | Field | Field |
| `description` | `string` | Field | Field |
| `code` | `string` | Field | Field |
| `externalCode` | `string` | Field | Field |
| `syncId` | `string` | Field | Field |
| `created` | `model.DateTime` | — | — |
| `updated` | `model.DateTime` | — | — |
| `deleted` | `model.DateTime` | — | — |
| `moment` | `model.DateTime` | Field | Field |
| `applicable` | `bool` | Field | Field |
| `shared` | `bool` | Field | Field |
| `printed` | `bool` | — | — |
| `published` | `bool` | — | — |
| `organization` | `*model.Reference` | required | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `project` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `meta` | `model.Meta` | — | — |
| `attributes` | `[]model.Attribute` | Field | Field |
| `rate` | `*model.DocumentRate` | Field | Field |
| `sum` | `float64` | — | — |
| `agent` | `*model.Reference` | required | Field |
| `organizationAccount` | `*model.Reference` | Field | Field |
| `agentAccount` | `*model.Reference` | Field | Field |
| `contract` | `*model.Reference` | Field | Field |
| `vatSum` | `float64` | — | — |
| `positions` | `*model.Collection[Position]` | Field | Field |
| `vatEnabled` | `bool` | Field | Field |
| `vatIncluded` | `bool` | Field | Field |
| `payedSum` | `float64` | — | — |
| `store` | `*model.Reference` | Field | Field |
| `payments` | `[]model.Reference` | Field | Field |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `salesChannel` | `*model.Reference` | Field | Field |
| `customerOrder` | `*model.Reference` | Field | Field |
| `factureOut` | `*model.Reference` | Field | Field |
| `shipmentAddress` | `string` | Field | Field |
| `shipmentAddressFull` | `*model.Address` | Field | Field |
| `returns` | `[]model.Reference` | Field | Field |
| `productionTasks` | `[]model.Reference` | Field | Field |
| `invoicesOut` | `[]model.Reference` | Field | Field |
| `overhead` | `*model.DocumentOverhead` | Field | Field |
| `cargoName` | `string` | Field | Field |
| `carrier` | `*model.Reference` | Field | Field |
| `consignee` | `*model.Reference` | Field | Field |
| `consigneeBranch` | `*model.Reference` | Field | Field |
| `consignor` | `*model.Reference` | Field | Field |
| `consignorBranch` | `*model.Reference` | Field | Field |
| `goodPackQuantity` | `int` | Field | Field |
| `shippingInstructions` | `string` | Field | Field |
| `stateContractId` | `string` | Field | Field |
| `transportFacility` | `string` | Field | Field |
| `transportFacilityNumber` | `string` | Field | Field |
| `shippingDocumentNumber` | `string` | Field | Field |
| `shippingDocumentDate` | `model.DateTime` | Field | Field |

Required Reference передаются обычными значениями; остальные поля записи —
`model.Field[T]`. Вложенные write-формы: `DocumentRateRequest`, `DocumentOverheadRequest`,
`AddressRequest`, `[]PositionRequest`.
Readonly поля в requests отсутствуют.

## Ограничения и TODO

При смене organization/agent API требует соответствующие счета сторон; SDK
передаёт organizationAccount/agentAccount без локального сравнения или скрытого Get.
Files представлены только metadata коллекции, без обещания rows или файловых операций.
Store — Field в Create, условную обязательность проверяет сервер. Safe Create
содержит store; сценарий applicable=false без store проверен только как точный wire JSON,
принятие живым API не проверено. ShipmentAddressFull:null не очищает адрес;
shipmentAddress:"" очищает, структурированная форма имеет приоритет. Используются
Address/AddressRequest, но FIASCodeRU у Demand не обещается. Cost строки записывается для
услуг; SDK передаёт Field без Get для определения
типа assortment. Unset/0/дробь/null сохраняются, очистка null сервером не обещается.
Overhead распределяет сервер, наличие позиций не проверяется скрытым Get.
FactureOut, Returns, Payments, ProductionTasks, InvoicesOut пишутся как связи документов
([общее правило](documents-extra.md#связи-документов)).
Транспортные поля (раздел «Другие поля» источника) читаются и пишутся: cargoName, carrier,
consignee, consigneeBranch, consignor, consignorBranch, goodPackQuantity, shippingInstructions,
stateContractId, transportFacility, transportFacilityNumber, shippingDocumentNumber,
shippingDocumentDate. Маркеров readonly нет; carrier/consignee/consignor — контрагент или юрлицо,
*Branch — подразделение юрлица, тип ссылки проверяет сервер.
Строки — все поля таблицы позиций источника: cost, pack, slot, things, trackingCodes (общая форма,
[контракт позиций](positions.md)); declaration, trackingCodes_1162 и overhead — только чтение.

Источник полей: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/demand),
[зафиксированный revision 9a811cef](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md).
