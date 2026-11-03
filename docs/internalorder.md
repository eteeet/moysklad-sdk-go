# InternalOrder

`entity/internalorder`: модель `InternalOrder`, `CreateRequest`, `UpdateRequest`, `NewService(client)`,
`Get/List/Create/Update/Delete`, `CreateMany/UpdateMany/DeleteMany`, `New` (шаблон), `Trash` (корзина)
и операции строк `GetPositions/GetPosition/CreatePositions/UpdatePosition/DeletePosition/DeleteManyPositions`.
Строки — aliases `Position`/`PositionRequest` общей формы (meta.type `internalorderposition`).
[Общие правила методов, JSON, ошибок и позиций](documents.md).

Минимальный Create: `organization` и непустое `name`. Строки можно передать сразу в теле документа
(`Positions`, не больше 1000 — 1001 отклоняется до запроса) или добавить потом `CreatePositions`.

```go
s, _ := internalorder.NewService(client)
doc, err := s.Create(ctx, internalorder.CreateRequest{Organization: org, Name: "ВЗ-1",
	Positions: model.Value([]internalorder.PositionRequest{{
		ItemRequest: positions.ItemRequest{Assortment: model.Value(product), Quantity: model.Value(2.0)},
	}}),
})
```

## Поддержанные поля

| JSON | Read Go type | Create | Update |
| --- | --- | --- | --- |
| `meta` | `model.Meta` | — | — |
| `id` | `string` | — | — |
| `accountId` | `string` | — | — |
| `name` | `string` | required | Field |
| `description` | `string` | Field | Field |
| `code` | `string` | Field | Field |
| `externalCode` | `string` | Field | Field |
| `syncId` | `string` | — | — |
| `created` | `model.DateTime` | — | — |
| `updated` | `model.DateTime` | — | — |
| `deleted` | `model.DateTime` | — | — |
| `moment` | `model.DateTime` | Field | Field |
| `deliveryPlannedMoment` | `model.DateTime` | Field | Field |
| `applicable` | `bool` | Field | Field |
| `shared` | `bool` | — | — |
| `printed` | `bool` | — | — |
| `published` | `bool` | — | — |
| `organization` | `*model.Reference` | required | Field |
| `store` | `*model.Reference` | Field | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `project` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `attributes` | `[]model.Attribute` | Field | Field |
| `rate` | `*model.DocumentRate` | Field | Field |
| `sum` | `float64` | — | — |
| `vatEnabled` | `bool` | Field | Field |
| `vatIncluded` | `bool` | Field | Field |
| `vatSum` | `float64` | — | — |
| `positions` | `*model.Collection[Position]` | Field | Field |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `moves` | `[]model.Reference` | Field | Field |
| `purchaseOrders` | `[]model.Reference` | Field | Field |
| `productionTasks` | `[]model.Reference` | Field | Field |

## Позиции

Таблица строк: `assortment`, `quantity`, `price` (копейки), `vat`, `vatEnabled`, `pack`; общая форма
`positions.Position`/`PositionRequest` покрывает её, остальные поля общей формы остаются нулевыми, а
их приём определяет сервер. `GetPositions` читает одну страницу (Limit 1…1000, 0 — 1000); `fields`
источник для внутреннего заказа не документирует, непустой `PageOptions.Fields` отклоняется до HTTP.
`CreatePositions` — один POST массивом, ответ `[]batch.Result[Position]`; `DeleteManyPositions` —
POST `…/positions/delete`, ответ без тела. В Update документа `Positions` заменяют весь набор: строка с
`ID` — существующая (id только адресует её), без — новая, лишние удаляются.

## Решения и противоречия источника

- `positions` в основной таблице помечен «Только для чтения», хотя текст раздела позиций и примеры
  Create, Update и массового создания передают массив строк в составе документа. SDK пишет
  `positions` в Create и Update: исключение из правила «readonly важнее примеров» утверждено пользователем, без него строки нельзя сохранить вместе с документом.
- `shared` и `syncId` помечены «Только для чтения» — в запросах их нет, хотя примеры Create/Update
  передают `shared`.
- `name` помечен «Необходимо при создании», а раздел «Создать» перечисляет только `organization`.
  SDK следует таблице: пустое `name` отклоняется до HTTP.
- `moves` и `purchaseOrders` — в основной таблице без readonly-маркера, `productionTasks` — только в
  разделе «Связи с другими документами». По общему правилу привязки документов (`_common_info.md`)
  все три пишутся. Ticket O08 называет moves/purchaseOrders read-only — это противоречит источнику,
  решение принято по источнику.
- Пример Create передаёт `discount` у строки и `sum`/`vatSum` у документа: discount нет в таблице
  строк, sum/vatSum — readonly; SDK их не задаёт (discount общей формы остаётся доступным, приём
  определяет сервер).
- `New` (PUT `entity/internalorder/new`): источник показывает только пустое тело; основания для
  внутреннего заказа не описаны, непустой basis передаётся как есть.
- `Trash` — POST `…/{id}/trash` (общий раздел «Удаление в корзину»), в отличие от `Delete`.
- Expand, метаданные и доп. поля (`/metadata`), wizard не входят.

Источник: [`_internalOrder.md`](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_internalOrder.md),
[официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/internalOrder),
revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`. Фикстуры `testdata/custom/internalorder/read.json`
и `positions.json` — примеры ответов «Получить Внутренний заказ» и «Получить позиции Внутреннего заказа».
