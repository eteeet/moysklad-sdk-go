# ProcessingStage

`entity/processingstage`: модель `ProcessingStage`, `CreateRequest`, `UpdateRequest`,
`NewService(client)`; методы `Get/List/Create/Update/Delete` и `CreateMany/UpdateMany/DeleteMany`.
[Общие правила методов, JSON и ошибок](documents.md).

- `Get(ctx, id)`, `Create(ctx, CreateRequest)`, `Update(ctx, id, UpdateRequest)` →
  `(*ProcessingStage, error)`; успех — HTTP 200 и непустой `id`, возвращается ответ сервера.
- `List(ctx, query.ListOptions)` — одна страница `GET entity/processingstage` без expand и автообхода.
- `Delete(ctx, id)` — `DELETE entity/processingstage/{id}`, безвозвратно (корзины у справочника нет).
- `CreateMany`, `UpdateMany` — один `POST entity/processingstage` массивом (до `batch.MaxItems`),
  `DeleteMany` — `POST entity/processingstage/delete`; ответ по элементам по правилам пакета `batch`.
- id проверяется как UUID до запроса; пустой `UpdateRequest` отклоняется до запроса; записи не повторяются.

Минимальный Create: непустой `name` (проверяется до запроса, без trim).

```go
stage, err := s.Create(ctx, processingstage.CreateRequest{
	Name:          "Сборка",
	AllPerformers: model.Value(false),
	Performers:    model.Value([]model.Reference{employee}),
})
```

## Поддержанные поля

| JSON | Read Go type | Create | Update |
| --- | --- | --- | --- |
| `meta` | `model.Meta` | — | — |
| `id` | `string` | — | — |
| `accountId` | `string` | — | — |
| `updated` | `model.DateTime` | — | — |
| `name` | `string` | required | Field |
| `description` | `string` | Field | Field |
| `externalCode` | `string` | Field | Field |
| `archived` | `bool` | Field | Field |
| `shared` | `bool` | Field | Field |
| `allPerformers` | `bool` | Field | Field |
| `distributionRequired` | `bool` | Field | Field |
| `standardHourCost` | `float64` | Field | Field |
| `performers` | `[]model.Reference` | Field | Field |
| `materialStore` | `*model.Reference` | — | — |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |

Field — `model.Field[T]` с `omitzero`: не задано — поле не уходит; `model.Value(false)` и
`model.Value(0.0)` уходят явно.

## Исполнители: allPerformers и performers

`allPerformers` — флаг «на этап можно назначить любого сотрудника», `performers` — ссылки на
сотрудников (`meta.type` `employee`). По источнику:

| allPerformers | performers | Смысл |
| --- | --- | --- |
| `true` | пусто | любой активный сотрудник |
| `false` | пусто | никто |
| `false` | непусто | только сотрудники из списка |

- Только `allPerformers=true` в запросе — сервер очищает `performers`; только `performers` —
  сервер ставит `allPerformers=false`.
- SDK отправляет оба поля как заданы и не проверяет их сочетание: исход `allPerformers=true` вместе с
  непустым `performers` определяет сервер (источник его не описывает).
- `performers` в Update заменяет весь список; `model.Value([]model.Reference{})` (и nil-слайс в
  `model.Value`) отправляет `[]`.

## Ограничения и противоречия

- `standardHourCost` — «Стоимость нормо-часа», `Double`; единицы источник не называет. SDK читает
  и пишет её как остальные суммы API — `float64` в копейках, без пересчёта.
- `materialStore` помечен «Только для чтения» — только в модели; в примерах источника его нет.
- `group` помечен «Обязательное при ответе» и `+Expand`, но не «Только для чтения» — пишется.
- Expand, metadata и шаблона в пакете нет; ticket их не называет.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingstage),
`md/dictionaries/_processingstage.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
