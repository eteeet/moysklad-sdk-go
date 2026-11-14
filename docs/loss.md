# Loss

`entity/loss`: модель `Loss`, `CreateRequest`, `UpdateRequest`, `NewService(client)`,
`Get/List/Create/Update` и `GetPositions`; позиции — aliases `Position`/`PositionRequest` общей формы.
[Общие правила методов, JSON, ошибок и позиций](documents.md).

Минимальный Create: `organization`, `store`. Positions в Create/Update — не больше 1000 позиций, при Update заменяют все позиции.

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
| `store` | `*model.Reference` | required | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `project` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `attributes` | `[]model.Attribute` | Field | Field |
| `rate` | `*model.DocumentRate` | Field | Field |
| `sum` | `float64` | — | — |
| `positions` | `*model.Collection[Position]` | Field | Field |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `salesReturn` | `*model.Reference` | Field | Field |
| `expenseItem` | `*model.Reference` | — | — |
| `inventory` | `*model.Reference` | Field | Field |

## Ограничения

- Позиции: reason, slot, pack, things; `GetPositions` с `fields=declaration` — прослеживаемость
  импортных товаров (readonly).
- salesReturn и inventory пишутся (связи документов); expenseItem — readonly.
- Пример Update источника передаёт sum, таблица помечает его readonly: SDK sum не пишет.
- Шаблон `New` и wizard — [documents-extra.md](documents-extra.md).

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/loss),
revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`.
