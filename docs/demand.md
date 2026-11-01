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
| `payments` | `[]model.Reference` | — | — |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `salesChannel` | `*model.Reference` | Field | Field |
| `customerOrder` | `*model.Reference` | Field | Field |
| `factureOut` | `*model.Reference` | — | — |
| `shipmentAddress` | `string` | Field | Field |
| `shipmentAddressFull` | `*model.Address` | Field | Field |
| `returns` | `[]model.Reference` | — | — |
| `productionTasks` | `[]model.Reference` | — | — |
| `invoicesOut` | `[]model.Reference` | — | — |
| `overhead` | `*model.DocumentOverhead` | Field | Field |

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
FactureOut, Returns, Payments, ProductionTasks, InvoicesOut — read-only в этом scope; запись
TODO.
Transport fields cargoName/carrier/consignee/consigneeBranch/consignor/consignorBranch/
goodPackQuantity/shippingInstructions/stateContractId/transportFacility/transportFacilityNumber/
shippingDocumentNumber/shippingDocumentDate — TODO и read, и write.
Поля строк общей формы (pack, things, trackingCodes, stock через `GetPositions` с
`fields=stock` и др.) описаны в [контракте позиций](positions.md); остальные не
перечисленные поля позиций остаются TODO.

Источник полей: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/demand),
[зафиксированный revision 9a811cef](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md).
