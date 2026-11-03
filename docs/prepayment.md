# Prepayment

`entity/prepayment`: модель `Prepayment`, `NewService(client)`; методы `Get/List`, `Delete`,
`Trash(ctx, id)`, `GetPositions(ctx, id, positions.PageOptions) (*model.Collection[Position], error)`
и `GetPosition(ctx, id, positionID) (*Position, error)`. Строки — alias `Position` общей формы
`positions.Position` (meta.type `prepaymentposition`).
[Общие правила методов, JSON, ошибок и позиций](documents.md).

Записи нет: источник не описывает Create, Update, массовые операции, шаблон и запись строк —
сервис их не предлагает. Предоплаты создаёт касса (розничная смена).

## Поля (только чтение)

| JSON | Read Go type |
| --- | --- |
| `meta` | `model.Meta` |
| `id`, `accountId` | `string` |
| `name`, `description`, `code`, `externalCode`, `syncId` | `string` |
| `created`, `updated`, `deleted`, `moment` | `model.DateTime` |
| `applicable`, `shared`, `printed`, `published` | `bool` |
| `organization`, `organizationBranch`, `agent`, `owner`, `group`, `state` | `*model.Reference` |
| `retailShift`, `retailStore`, `customerOrder` | `*model.Reference` |
| `returns` | `[]model.Reference` |
| `attributes` | `[]model.Attribute` |
| `rate` | `*model.DocumentRate` |
| `sum`, `cashSum`, `noCashSum`, `qrSum`, `vatSum` | `float64` (копейки) |
| `vatEnabled`, `vatIncluded` | `bool` |
| `taxSystem` | `string` (GENERAL_TAX_SYSTEM … PATENT_BASED) |
| `positions` | `*model.Collection[Position]` |
| `files` | `*model.Collection[model.Reference]` |

Строки: `assortment`, `quantity`, `price` (копейки), `discount`, `vat`, `vatEnabled`, `pack`;
остальные поля общей формы остаются нулевыми.

## Ограничения

- `GetPositions` — одна страница (limit 1–1000, offset), `fields` не поддержан.
- Пример ответа содержит `version` и `documents`, которых нет в таблице атрибутов, — SDK их не
  читает.
- Таблица помечает `organization` и `retailShift` «Необходимо при создании», но операции
  создания в источнике нет.
- Массового удаления в источнике нет — `DeleteMany` не реализован.
- `Trash` — общее правило документов (`_common_info.md`, «Удаление в корзину»); в
  `_prepayment.md` отдельного раздела нет. Delete удаляет безвозвратно.
- Expand не поддержан.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/prepayment),
`md/documents/_prepayment.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
