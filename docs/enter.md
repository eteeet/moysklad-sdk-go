# Enter

`entity/enter`: модель `Enter`, `CreateRequest`, `UpdateRequest`, `NewService(client)`,
`Get/List/Create/Update` и `GetPositions`; строки — aliases `Position`/`PositionRequest` общей формы.
[Общие правила методов, JSON, ошибок и позиций](documents.md).

Минимальный Create: `organization`, `store`. Контрагента у оприходования нет. Positions в Create/Update — не больше 1000 строк, при Update заменяют все строки.

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
| `overhead` | `*model.DocumentOverhead` | Field | Field |
| `inventory` | `*model.Reference` | — | — |

## Ограничения

- Строки: reason, country, gtd, things, slot, pack; overhead строки только читается.
  fields у позиций источник не документирует.
- overhead документа `{sum, distribution}` распределяет сервер.
- inventory — ссылка, маркера readonly нет, в примерах записи нет: SDK её не пишет.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/enter),
revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`.
