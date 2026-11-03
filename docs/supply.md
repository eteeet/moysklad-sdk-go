# Supply

`entity/supply`: модель `Supply`, `CreateRequest`, `UpdateRequest`,
`NewService(client)` и `Get/Create/Update`. Дополнительно `GetPositions` и aliases
`Position`/`PositionRequest` общих форм.
[Общие правила методов, JSON, ошибок и позиций](documents.md).

Минимальный Create: `organization`, `agent`, `store`.

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
| `store` | `*model.Reference` | required | Field |
| `payments` | `[]model.Reference` | Field | Field |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `incomingDate` | `model.DateTime` | Field | Field |
| `incomingNumber` | `string` | Field | Field |
| `purchaseOrder` | `*model.Reference` | Field | Field |
| `factureIn` | `*model.Reference` | Field | Field |
| `productionTask` | `*model.Reference` | Field | Field |
| `invoicesIn` | `[]model.Reference` | Field | Field |
| `returns` | `[]model.Reference` | Field | Field |
| `overhead` | `*model.DocumentOverhead` | Field | Field |

Required Reference передаются обычными значениями; остальные поля записи —
`model.Field[T]`. Вложенные write-формы: `DocumentRateRequest`, `DocumentOverheadRequest`,
`[]PositionRequest`.
Readonly поля в requests отсутствуют.

## Ограничения и TODO

При смене organization/agent API требует соответствующие счета сторон; SDK
передаёт organizationAccount/agentAccount без локального сравнения или скрытого Get.
Files представлены только metadata коллекции, без обещания rows или файловых операций.
Overhead распределяет сервер; без positions его задавать нельзя по API.
FactureIn, ProductionTask, InvoicesIn, Payments, Returns пишутся как связи документов
([общее правило](documents-extra.md#связи-документов)).
Pack, things, slot (кроме SourceSlot/TargetSlot общей формы), serial/tracking/stock
и остальные не перечисленные поля позиций остаются TODO.

Источник полей: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/supply),
[зафиксированный revision 9a811cef](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md).
