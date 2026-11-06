# InvoiceIn

`entity/invoicein`: модель `InvoiceIn`, `CreateRequest`, `UpdateRequest`,
`NewService(client)` и `Get/Create/Update`. Дополнительно `GetPositions` и aliases
`Position`/`PositionRequest` общих форм.
[Общие правила методов, JSON, ошибок и позиций](documents.md).

Минимальный Create: `organization`, `agent` (`name` необязателен: без него номер присваивает МойСклад).

## Поддержанные поля

| JSON | Read Go type | Create | Update |
| --- | --- | --- | --- |
| `id` | `string` | — | — |
| `accountId` | `string` | — | — |
| `name` | `string` | string (пусто не уходит — автонумерация МС) | Field |
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
| `paymentPlannedMoment` | `model.DateTime` | Field | Field |
| `shippedSum` | `float64` | — | — |
| `incomingDate` | `model.DateTime` | Field | Field |
| `incomingNumber` | `string` | Field | Field |
| `purchaseOrder` | `*model.Reference` | Field | Field |
| `supplies` | `[]model.Reference` | Field | Field |
| `files` | — | Field (`[]files.Upload`) | — |

Required Reference передаются обычными значениями; остальные поля записи —
`model.Field[T]`. Вложенные write-формы: `DocumentRateRequest`, `[]PositionRequest`.
Readonly поля в requests отсутствуют.

## Ограничения и TODO

При смене organization/agent API требует соответствующие счета сторон; SDK
передаёт organizationAccount/agentAccount без локального сравнения или скрытого Get.
Files при чтении исключены; при создании `Files` прикрепляет файлы тем же запросом
(`[]files.Upload`: имя и содержимое Base64); добавить к существующему документу — пакет `files`. Supplies пишутся как связи документов; payments источник помечает «Только для чтения».
Строки — все поля таблицы позиций источника, включая pack (общая форма); things, slot и
кодов маркировки у строк счёта источник не даёт.

Источник полей: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/invoice-in),
[зафиксированный revision 9a811cef](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_invoice_in.md).
