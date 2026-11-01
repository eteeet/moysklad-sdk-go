# CustomerOrder

`entity/customerorder`: модель `CustomerOrder`, `CreateRequest`, `UpdateRequest`, `NewService(client)`,
`Get/List/Create/Update` и `GetPositions`; строки — aliases `Position`/`PositionRequest` общей формы.
[Общие правила методов, JSON, ошибок и позиций](documents.md).

Минимальный Create: `organization`, `agent`. Positions в Create/Update — не больше 1000 строк, при Update заменяют все строки.

## Поддержанные поля

| JSON | Read Go type | Create | Update |
| --- | --- | --- | --- |
| `meta` | `model.Meta` | — | — |
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
| `deliveryPlannedMoment` | `model.DateTime` | Field | Field |
| `applicable` | `bool` | Field | Field |
| `shared` | `bool` | Field | Field |
| `printed` | `bool` | — | — |
| `published` | `bool` | — | — |
| `organization` | `*model.Reference` | required | Field |
| `organizationAccount` | `*model.Reference` | Field | Field |
| `agent` | `*model.Reference` | required | Field |
| `agentAccount` | `*model.Reference` | Field | Field |
| `contract` | `*model.Reference` | Field | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `project` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `store` | `*model.Reference` | Field | Field |
| `salesChannel` | `*model.Reference` | Field | Field |
| `attributes` | `[]model.Attribute` | Field | Field |
| `rate` | `*model.DocumentRate` | Field | Field |
| `taxSystem` | `string` | Field | Field |
| `vatEnabled` | `bool` | Field | Field |
| `vatIncluded` | `bool` | Field | Field |
| `sum` | `float64` | — | — |
| `vatSum` | `float64` | — | — |
| `payedSum` | `float64` | — | — |
| `shippedSum` | `float64` | — | — |
| `invoicedSum` | `float64` | — | — |
| `reservedSum` | `float64` | — | — |
| `shipmentAddress` | `string` | Field | Field |
| `shipmentAddressFull` | `*model.Address` | Field | Field |
| `positions` | `*model.Collection[Position]` | Field | Field |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `purchaseOrders` | `[]model.Reference` | — | — |
| `demands` | `[]model.Reference` | — | — |
| `payments` | `[]model.Reference` | — | — |
| `productionTasks` | `[]model.Reference` | — | — |
| `invoicesOut` | `[]model.Reference` | — | — |
| `moves` | `[]model.Reference` | — | — |
| `prepayments` | `[]model.Reference` | — | — |

## Ограничения

- Строки: reserve пишется, shipped только читается; pack, taxSystem. `GetPositions` принимает
  `fields=stock` (остатки и себестоимость на текущий момент).
- При смене organization/agent сервер требует и соответствующий счёт; SDK не сравнивает и не делает Get.
- syncId после заполнения не меняется. Адресу null не поддерживается: удалить — `shipmentAddress: ""`.
- Связанные документы (purchaseOrders, demands, payments, productionTasks, invoicesOut, moves,
  prepayments) и суммы — только для чтения. Files — только metadata коллекции.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/customerOrder),
revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`.
