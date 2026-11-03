# Inventory

`entity/inventory`: модель `Inventory`, `CreateRequest`, `UpdateRequest`, `NewService(client)`,
`Get/List/Create/Update/Delete`, `CreateMany/UpdateMany/DeleteMany`, `New` (шаблон), `Trash` (корзина)
и операции строк `GetPositions/GetPosition/CreatePositions/UpdatePosition/DeletePosition/DeleteManyPositions`.
Строки — собственные `Position`/`PositionRequest` (meta.type `inventoryposition`), не общая форма.
[Общие правила методов, JSON, ошибок и позиций](documents.md).

Минимальный Create: `organization`, `store`. Positions в Create/Update — не больше 1000 строк, при
Update заменяют все строки.

```go
s, _ := inventory.NewService(client)
doc, err := s.Create(ctx, inventory.CreateRequest{
	Organization: org, Store: store,
	Positions: model.Value(inventory.PositionRows{Rows: []inventory.PositionRequest{{
		ItemRequest: positions.ItemRequest{Assortment: model.Value(product), Quantity: model.Value(140.0)},
	}}}),
})
```

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
| `shared` | `bool` | Field | Field |
| `printed` | `bool` | — | — |
| `published` | `bool` | — | — |
| `organization` | `*model.Reference` | required | Field |
| `store` | `*model.Reference` | required | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `attributes` | `[]model.Attribute` | Field | Field |
| `sum` | `float64` | — | — |
| `positions` | `*model.Collection[Position]` | Field[`PositionRows`] | Field[`PositionRows`] |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `enters` | `[]model.Reference` | Field | Field |
| `losses` | `[]model.Reference` | Field | Field |

У инвентаризации в таблице нет `applicable`, `rate`, `project` и полей НДС — в модели их нет.

## Позиции

| JSON | Read Go type | Request |
| --- | --- | --- |
| `meta`, `id`, `accountId` | `positions.Core` (через `positions.Item`) | `id` — Field (сопоставление строки в embedded Update) |
| `assortment` | `model.Reference` | Field |
| `quantity` | `float64` | Field |
| `price` | `float64` (копейки) | Field |
| `pack` | `*positions.Pack` | Field[`positions.PackRequest`] |
| `calculatedQuantity` | `float64` (расчётный остаток) | Field |
| `correctionAmount` | `float64` (разница с фактическим) | — |
| `correctionSum` | `float64` (избыток/недостача, копейки) | — |

`calculatedQuantity` маркера readonly не имеет и пишется (его передаёт и пример Update);
`correctionAmount` и `correctionSum` — «Только для чтения», в `PositionRequest` их нет.
`GetPositions` читает одну страницу (Limit 1…1000, 0 — 1000); `fields` источник не документирует,
непустой `PageOptions.Fields` отклоняется до HTTP. `CreatePositions` — один POST массивом, ответ
`[]batch.Result[Position]`; `UpdatePosition` с пустым запросом отклоняется до HTTP;
`DeleteManyPositions` — POST `…/positions/delete`, ответ без тела.

## Решения и противоречия источника

- `enters` и `losses` — `[]model.Reference`, хотя источник описывает их в единственном числе и без
  примеров: если сервер вернёт объект, а не массив, Get и List документа не разберутся.
- `enters` и `losses` описаны только в разделе «Связи с другими документами», в основной таблице их
  нет. По общему правилу привязки документов (`_common_info.md`) поля связей пишутся, если таблица не
  помечает их readonly, поэтому SDK их читает и пишет. Ticket O09 называет enter/loss read-only — это
  противоречит источнику, решение принято по источнику. Описание называет поле «ссылкой» в
  единственном числе, но имя во множественном числе по общему правилу означает коллекцию: SDK
  использует массив ссылок; примеров с этими полями в источнике нет.
- Текст раздела позиций называет `positions` массивом, но все примеры запросов передают объект
  `{"rows":[…]}`: SDK следует примерам (`PositionRows`), как у прайс-листа.
- `New` (PUT `entity/inventory/new`): источник описывает шаблон только без основания; непустой basis
  передаётся как есть, его принятие определяет сервер.
- `Trash` — POST `…/{id}/trash` (общий раздел «Удаление в корзину»), в отличие от `Delete`.
- Пересчёт расчётного остатка (`rpc/inventory/{id}/recalcCalculatedQuantity`, раздел
  «Пересчет расчетного остатка в инвентаризации») в ticket O09 не входит и не реализован.
- Expand, метаданные и доп. поля (`/metadata`) не входят.

Источник: [`_inventory.md`](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_inventory.md),
[официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/inventory),
revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`. Фикстуры `testdata/custom/inventory/read.json`
и `positions.json` — примеры ответов «Получить Инвентаризацию» и «Получить позиции».
