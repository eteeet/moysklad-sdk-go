# ProcessingOrder

`entity/processingorder`: модель `ProcessingOrder`, `CreateRequest`, `UpdateRequest`,
`NewService(client)`; методы `Get/List/Create/Update/Delete`, `CreateMany/UpdateMany/DeleteMany`,
`New` (шаблон на основе техкарты), `Trash` (корзина) и операции строк
`GetPositions/GetPosition/UpdatePosition/DeletePosition/DeleteManyPositions`. Строки — собственная
форма `Position`/`PositionRequest` (meta.type `processingorderposition`), не `positions.Position`:
цен, НДС и скидок у строки нет. `CreatePositions` нет — источник не описывает «Создать позицию».
[Общие правила методов, JSON и ошибок](documents.md), [корзина и шаблоны](documents-extra.md).

Минимальный Create: непустое `name`, `organization`, `store`, `processingPlan` и `positions`
(1–1000 строк) — проверяются до запроса. Update всегда требует `processingPlan` (даже если техкарта
не меняется) — пустой meta.href отклоняется до запроса.

```go
s, _ := processingorder.NewService(client)
doc, err := s.Create(ctx, processingorder.CreateRequest{
	Name: "ЗП-1", Organization: org, Store: store, ProcessingPlan: plan,
	Positions: []processingorder.PositionRequest{{Assortment: model.Value(product), Quantity: model.Value(111.0)}},
	Quantity: model.Value(5.0),
})
doc, err = s.Update(ctx, doc.ID, processingorder.UpdateRequest{ProcessingPlan: plan, Name: model.Value("ЗП-2")})
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
| `store` | `*model.Reference` | required | Field |
| `processingPlan` | `*model.Reference` | required | required |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `project` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `attributes` | `[]model.Attribute` | — | — |
| `quantity` | `float64` | Field | Field |
| `positions` | `*model.Collection[Position]` | required | Field |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `processings` | `[]model.Reference` | Field | Field |

Денежных полей в таблице нет: `quantity` документа (объём производства), `quantity` и `reserve` строки —
количества `float64`, не копейки. `UpdateMany` проверяет до запроса только UUID и непустой Request;
`processingPlan` у каждого элемента проверяет сервер.

## Позиции

| Метод | Запрос |
| --- | --- |
| `GetPositions(ctx, id, positions.PageOptions) (*model.Collection[Position], error)` | GET `…/{id}/positions` — одна страница, только limit/offset (`fields` отклоняется до запроса) |
| `GetPosition(ctx, id, positionID)` | GET `…/positions/{positionID}` |
| `UpdatePosition(ctx, id, positionID, PositionRequest)` | PUT `…/positions/{positionID}`; пустой запрос отклоняется до запроса |
| `DeletePosition(ctx, id, positionID) error` | DELETE `…/positions/{positionID}` |
| `DeleteManyPositions(ctx, id, positionIDs) error` | POST `…/positions/delete` с meta строк (1–1000); ответ без тела |

- Поля строки: `meta`, `id`, `accountId` (только чтение), `assortment`, `pack`
  (`positions.Pack`/`positions.PackRequest`), `quantity`, `reserve`.
- Строки добавляются только в составе документа: `Positions` в Create (обязательно, 1–1000) или
  Update (не больше 1000, 1001 отклоняется до запроса). В Update массив заменяет весь набор: строка с
  `ID` — существующая, без — новая. Источник: позиции с новыми товарами в заказ добавить нельзя
  (строки соответствуют материалам техкарты) — проверяет сервер.

## Ограничения и противоречия источника

- `processings` (раздел «Связи с другими документами») таблица не помечает «Только для чтения», и
  по общему правилу привязки документов (`_common_info.md`) такие коллекции пишутся. Ticket M06
  называет processings read-only ссылками — это противоречит источнику; решение принято по
  источнику: `Processings` есть в Create и Update (`Value([]…)` отвязывает все).
- `attributes` таблица помечает «Только для чтения», хотя примеры Create и Update «с доп. полями»
  передают `attributes`. По правилу «readonly важнее примеров» SDK их не пишет; исключение
  требует решения пользователя.
- `shared` и `syncId` помечены «Только для чтения» — не пишутся, хотя пример Create передаёт `shared`.
- `name` и `store` помечены «Необходимо при создании», а раздел «Создать» и минимальный пример
  называют только `organization`, `processingPlan` и `positions` (сервер в примере сам присваивает
  имя). SDK следует таблице: пустые `name`/`store` отклоняются до запроса.
- Примеры ответов содержат `sum`, которого нет в таблице атрибутов, — в модель не входит.
- `files` маркера readonly не имеет, но, как у остальных документов SDK, не пишется (файлы — пакет
  `files`).
- `New` (PUT `entity/processingorder/new`): источник показывает основание `processingPlan`; ответ
  без `id`, строки — `positions.rows` без meta. nil basis отправляет пустое тело — источник такого
  примера не даёт.
- `Trash` — POST `…/{id}/trash` (общий раздел «Удаление в корзину»), в отличие от `Delete`.
- Metadata (`entity/processingorder/metadata`, доп. поля, статусы) покрыта общим пакетом `metadata`
  (`metadata.Service.Get(ctx, "processingorder")`); ticket её не называет, отдельных тестов в пакете нет.
- Expand и подбор строк не входят.

Источник: [`_processingorder.md`](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_processingorder.md),
[официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processingorder),
revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`. Фикстуры `testdata/custom/processingorder`:
`read.json` — «Получить Заказ на производство», `list-row.json` — первая строка «Получить список»
(с `processings`), `positions.json` — «Получить позиции», `template.json` — шаблон на основе техкарты;
`create-full.golden.json` и `update.golden.json` — независимые golden запросов.
