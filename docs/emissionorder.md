# EmissionOrder

`entity/emissionorder`: модель `EmissionOrder`, `CreateRequest`, `UpdateRequest`, `NewService(client)`;
методы `Get/List/Create/Update`, `CreateMany/UpdateMany` и `Trash(ctx, id)`. Строки — собственная
форма `Position` (meta.type `emissionorderposition`) и `PositionRequest` (= `positions.ItemRequest`).
Handle фасада: `c.Entity().EmissionOrder()`.
[Общие правила методов, JSON и ошибок](documents.md), [корзина](documents-extra.md).

Минимальный Create: `organization`, `trackingType`, `emissionType` — проверяются до запроса
(непустой meta.href, непустые строки). Значения enum SDK не проверяет.

```go
doc, err := s.Create(ctx, emissionorder.CreateRequest{
	Organization: organization, TrackingType: "MILK", EmissionType: "LOCAL",
	Positions: model.Value([]emissionorder.PositionRequest{{
		Assortment: model.Value(product), Quantity: model.Value(3.0),
	}}),
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
| `externalCode` | `string` | Field | Field |
| `created` | `model.DateTime` | — | — |
| `updated` | `model.DateTime` | — | — |
| `moment` | `model.DateTime` | Field | Field |
| `shared` | `bool` | Field | Field |
| `printed` | `bool` | — | — |
| `published` | `bool` | — | — |
| `organization` | `*model.Reference` | required | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `trackingType` | `string` (открытый enum) | required | Field |
| `emissionType` | `string` (открытый enum) | required | Field |
| `documentState` | `string` (открытый enum) | — | — |
| `positions` | `*model.Collection[Position]` | Field | Field |
| `productionTasks` | `[]model.Reference` | Field | Field |

`productionTasks` — «Связи с другими документами» (производственные задания); маркера «Только для
чтения» нет, поэтому пишется.

## Позиции

- `GetPositions(ctx, id, positions.PageOptions) (*model.Collection[Position], error)` — одна страница
  `…/{id}/positions`, только limit/offset (`fields` отклоняется до запроса).
- `GetPosition`, `CreatePositions` (один POST массивом, ответ `[]batch.Result[Position]`),
  `UpdatePosition` (пустой request — ошибка до запроса).
- Поля строки: `meta`, `id`, `accountId`, `status` (только чтение; `status` — открытый enum
  `EMISSION_*`), `assortment` (`model.Reference`), `quantity` (`float64`, не деньги).
- Positions в Create/Update — не больше 1000 строк (1001 отклоняется до запроса); в Update массив
  заменяет весь набор: строка с `ID` — существующая, без — новая.

## Ограничения и противоречия

- Delete/DeleteMany документа и удаления строк в источнике нет — методов нет; есть только `Trash`
  (общий раздел `documents/_common_info.md`, «Удаление в корзину»).
- Лимиты источника «не больше 10 позиций», `quantity` до 1000 шт., положительный `quantity` и
  совпадение `trackingType` строк и заказа проверяет сервер.
- Ticket называет `cisType` открытой строкой, но в таблицах emissionorder такого поля нет — SDK его
  не добавляет.
- Таблица «Создать» описывает `emissionType` как «Состояние документов маркировки» со ссылкой на
  `documentState`; по таблице атрибутов это «Способ ввода в оборот» (LOCAL, FOREIGN, …). SDK следует
  таблице атрибутов.
- Пример «Создать позицию» с одним объектом в теле возвращает массив; SDK всегда отправляет массив
  (`CreatePositions`).
- `Метаданные` (`…/metadata`), expand и шаблон `New` не входят в ticket.
- List: источник описывает только limit/offset; SDK передаёт `query.ListOptions` как у остальных
  документов, поддержку filter/search определяет сервер.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/emissionorder),
`md/documents/_emissionorder.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
Фикстура чтения — `testdata/upstream/emission_order.json`.
