# PriceList

`entity/pricelist`: модель `PriceList`, `CreateRequest`, `UpdateRequest`, `NewService(client)`;
методы `Get/List/Create/Update/Delete`, `CreateMany/UpdateMany/DeleteMany` и `Trash(ctx, id)`.
Строки — собственная форма `Position` (meta.type `pricelistrow`), не `positions.Position`.
[Общие правила методов, JSON и ошибок](documents.md), [корзина](documents-extra.md).

Минимальный Create: `columns` — хотя бы один столбец с непустым `name` (проверяется до запроса;
уникальность имён проверяет сервер).

```go
doc, err := s.Create(ctx, pricelist.CreateRequest{
	Columns: []pricelist.ColumnRequest{{Name: "Опт"}, {Name: "Розница", PercentageDiscount: model.Value(10)}},
	Positions: model.Value(pricelist.PositionRows{Rows: []pricelist.PositionRequest{{
		Assortment: model.Value(product),
		Cells:      model.Value([]pricelist.Cell{{Column: "Опт", Sum: 90000}}),
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
| `applicable` | `bool` | Field | Field |
| `shared` | `bool` | Field | Field |
| `printed` | `bool` | — | — |
| `published` | `bool` | — | — |
| `organization` | `*model.Reference` | Field | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `attributes` | `[]model.Attribute` | Field | Field |
| `columns` | `[]Column` | required `[]ColumnRequest` | — (после создания не меняется) |
| `priceType` | `*model.PriceType` | Field | — |
| `positions` | `*model.Collection[Position]` | Field `PositionRows` | Field `PositionRows` |
| `files` | `*model.Collection[model.Reference]` | — | — |

Столбец: `name` (String, необходим), `percentageDiscount` (Int, скидка/наценка по умолчанию; в
запросе `model.Field[int]`, чтобы отличать 0 от «не задано»). Единицу источник не называет: таблица
говорит о проценте, а все примеры передают 1000 — SDK значение не пересчитывает. Ячейка `Cell`: `column` (name
столбца) и `sum` — оба необходимы.

## Позиции

| Метод | Маршрут |
| --- | --- |
| `GetPositions(ctx, id, positions.PageOptions) (*model.Collection[Position], error)` | `GET entity/pricelist/{id}/positions?limit&offset` |
| `GetPosition(ctx, id, positionID) (*Position, error)` | `GET …/positions/{positionID}` |
| `CreatePositions(ctx, id, []PositionRequest) ([]batch.Result[Position], error)` | `POST …/positions` массивом |
| `UpdatePosition(ctx, id, positionID, PositionRequest) (*Position, error)` — пустой запрос отклоняется | `PUT …/positions/{positionID}` |
| `DeletePosition(ctx, id, positionID) error` | `DELETE …/positions/{positionID}` |
| `DeleteManyPositions(ctx, id, []positionID) error` — 1–1000, ответ без тела | `POST …/positions/delete` |

`Position`: `meta`, `id`, `accountId`, `assortment`, `cells []Cell`, `pack *positions.Pack`.
`PositionRequest`: `id` (сопоставление строки в Update прайс-листа), `assortment` (нужен при
создании строки, проверяет сервер), `cells` (можно не передавать — сервер заполнит ячейки столбцов
со скидкой по умолчанию), `pack`. В Create/Update прайс-листа строки передаются как
`PositionRows{Rows: …}` → `"positions": {"rows": [...]}`, не больше 1000; при Update набор строк
заменяется полностью.

## Ограничения и противоречия источника

- `priceType` таблица помечает «Только для чтения», а пример Create передаёт его и примечание
  описывает поведение при несуществующем типе цен (сервер подставит тип цен продаж по умолчанию).
  SDK пишет `priceType` только в Create: исключение из правила «readonly важнее примеров» утверждено пользователем. Пример Update его не передаёт, текст «Изменить»
  запрещает readonly-поля — в `UpdateRequest` его нет.
- В примерах Create/Update `positions` — объект `{"rows": [...]}`, а не массив, как у других
  документов; SDK следует примерам.
- `Cell.sum` в таблице — Int («числовое значение ячейки»), в примерах — цены в копейках; SDK
  держит `float64`, как остальные суммы.
- Таблица не помечает `organization` обязательным, обязательны только `columns`.
- `DeleteManyPositions` — общий `positions.Service.DeleteManyPositions`; тип строки `pricelistrow`
  внесён в общую таблицу positions.
- Шаблона `New`, wizard и expand нет; Delete удаляет безвозвратно, Trash — в корзину.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/pricelist),
`md/documents/_pricelist.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
