# SalesReturn

`entity/salesreturn`: модель `SalesReturn`, `CreateRequest`, `UpdateRequest`, `NewService(client)`,
`Get/List/Create/Update` и `GetPositions`; позиции — aliases `Position`/`PositionRequest` общей формы.
[Общие правила методов, JSON, ошибок и позиций](documents.md).

Минимальный Create: `organization`, `store`, `agent`; `demand` — только для возврата по основанию. Positions в Create/Update — не больше 1000 позиций, при Update заменяют все позиции.

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
| `applicable` | `bool` | Field | Field |
| `shared` | `bool` | Field | Field |
| `printed` | `bool` | — | — |
| `published` | `bool` | — | — |
| `organization` | `*model.Reference` | required | Field |
| `organizationAccount` | `*model.Reference` | Field | Field |
| `agent` | `*model.Reference` | required | — |
| `agentAccount` | `*model.Reference` | Field | — |
| `contract` | `*model.Reference` | Field | Field |
| `store` | `*model.Reference` | required | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `project` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `salesChannel` | `*model.Reference` | Field | Field |
| `attributes` | `[]model.Attribute` | Field | Field |
| `rate` | `*model.DocumentRate` | Field | Field |
| `vatEnabled` | `bool` | Field | Field |
| `vatIncluded` | `bool` | Field | Field |
| `sum` | `float64` | — | — |
| `vatSum` | `float64` | — | — |
| `payedSum` | `float64` | — | — |
| `positions` | `*model.Collection[Position]` | Field | Field |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `demand` | `*model.Reference` | Field | — |
| `losses` | `[]model.Reference` | Field | Field |
| `payments` | `[]model.Reference` | Field | Field |
| `factureOut` | `*model.Reference` | Field | Field |

## Ограничения

- Create по основанию: контрагент, юрлицо и валюта совпадают с основанием; позиции соответствуют
  его позициям, quantity не больше, чем в основании, позиций не из основания нет, price менять
  нельзя. Проверяет сервер; SDK не делает Get основания.
- Update не содержит agent, agentAccount и demand — источник запрещает их менять. Валюту,
  отличную от основания, и позиции, которых нет в основании, отклоняет сервер. Запрет price
  источник формулирует только для Create по основанию.
- Позиции: cost (возврат без основания), country, gtd, slot, pack, things. `GetPositions` принимает
  `fields=stock`; у возврата без основания stock приходит без cost.
- losses, payments, factureOut пишутся как связи документов; payedSum — только для чтения.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/sales-return),
revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`.
