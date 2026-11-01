# CashOut

`entity/cashout`: модель `CashOut`, `CreateRequest`, `UpdateRequest`,
`NewService(client)` и `Get/Create/Update`.
[Общие правила методов, JSON, ошибок и позиций](documents.md).

Минимальный Create: `organization`, `agent`, `expenseItem`.

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
| `sum` | `float64` | Field | Field |
| `agent` | `*model.Reference` | required | Field |
| `contract` | `*model.Reference` | Field | Field |
| `vatSum` | `float64` | Field | Field |
| `paymentPurpose` | `string` | Field | Field |
| `operations` | `[]model.FinanceOperation` | Field | Field |
| `salesChannel` | `*model.Reference` | Field | Field |
| `expenseItem` | `*model.Reference` | required | Field |
| `factureIn` | `*model.Reference` | Field | Field |
| `noClosingDocs` | `bool` | Field | Field |
| `accrualDate` | `model.DateTime` | Field | Field |

Required Reference передаются обычными значениями; остальные поля записи —
`model.Field[T]`. Вложенные write-формы: `DocumentRateRequest`, `[]FinanceOperationRequest`.
Readonly поля в requests отсутствуют.

## Ограничения и TODO

Operations — плоские meta + linkedSum. В read LinkedSum — float64, в write —
Field[float64], поэтому 0 отличается от omission. Неизвестный meta.type сохраняется.
Документированные операции: salesreturn, supply, invoicein, purchaseorder,
commissionreportout, payroll.
SDK не проверяет распределение, типы/суммы операций или соответствие сторон; []
передаётся точно без обещания очистки/полного замещения operations сервером.
Банковских agentAccount/organizationAccount в срезе нет. Текст Update требует их
одновременной смены со сторонами, но атрибутные таблицы их не содержат: это расхождение
официальных источников, SDK не добавляет поля по аналогии. Валютная проверка банковского
счёта Organization прямо не применяется к кассовым ордерам.
Agent документирован как Counterparty; открытая Reference не вводит локальный enum.
NoClosingDocs=true несовместим с непустым operations; это проверяет сервер.
Files целиком и expanded-поля operations исключены; наличие их в upstream fixture
не означает поддержки expand.

Источник полей: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/cashout),
[зафиксированный revision 9a811cef](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashout.md).
