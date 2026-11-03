# ProcessingProcess

`entity/processingprocess`: модель `ProcessingProcess`, `CreateRequest`, `UpdateRequest`,
`NewService(client)`; методы `Get/List/Create/Update/Delete`, `CreateMany/UpdateMany/DeleteMany`.
Позиции (этапы) — собственная форма `Position`/`PositionRequest`/`PositionUpdateRequest`
(не `positions.Position`). Справочник: `Trash` и шаблона `New` нет.
[Общие правила методов, JSON и ошибок](documents.md).

Минимальный Create: `name` и `positions` — 1–100 позиций, у каждой непустой `processingstage`
(meta.href); проверяется до запроса.

```go
p, err := s.Create(ctx, processingprocess.CreateRequest{
	Name: "Изготавливаем двигатель",
	Positions: []processingprocess.PositionRequest{
		{ProcessingStage: cutting}, // этапы уходят в порядке среза
		{ProcessingStage: assembly},
	},
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
| `externalCode` | `string` | Field | Field |
| `archived` | `bool` | Field | Field |
| `shared` | `bool` | Field | Field |
| `updated` | `model.DateTime` | — | — |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `positions` | `*model.Collection[Position]` | required `[]PositionRequest` | Field `[]PositionUpdateRequest` |

## Позиции

Ресурс `entity/processingprocess/{id}/positions`, meta.type `processingprocessposition`.

| Метод | Запрос |
| --- | --- |
| `GetPositions(ctx, id, query.PageOptions)` | GET, одна страница, только limit/offset (`fields` отклоняется до запроса) |
| `GetPosition(ctx, id, positionID)` | GET |
| `CreatePositions(ctx, id, []PositionRequest)` | POST массивом, ответ `[]batch.Result[Position]` |
| `UpdatePosition(ctx, id, positionID, PositionUpdateRequest)` | PUT; пустой request и заданный `Meta` отклоняются до запроса |
| `DeletePosition(ctx, id, positionID)` | DELETE |
| `DeleteManyPositions(ctx, id, positionIDs)` | POST `…/positions/delete`, 1–1000 id; источник отвечает 200 без тела — возвращается только error |

| JSON | Read Go type | PositionRequest (новая) | PositionUpdateRequest |
| --- | --- | --- | --- |
| `meta` | `model.Meta` | — | Field: только в `UpdateRequest.Positions`, адрес существующей позиции |
| `id` | `string` | — | — |
| `accountId` | `string` | — | — |
| `processingstage` | `*model.Reference` | required | Field |
| `nextPositions` | `[]model.Reference` | Field | Field |

## Порядок этапов и nextPositions

- Порядок позиций в теле Create/Update и CreatePositions — порядок среза; SDK его не меняет.
- В `UpdateRequest.Positions` строка с `Meta` (из `Position.Meta`) изменяет существующую позицию:
  смена `processingstage` позицию не пересоздаёт (в отличие от UI); строка без `Meta` — новая
  позиция, для неё `processingstage` проверяется до запроса. Массив — 1–100 строк, `Null` отклоняется.
  Что сервер делает с позициями, не переданными в массиве, источник не описывает.
- `nextPositions` — ссылки на следующие позиции того же техпроцесса (meta позиций). Поле не помечено
  «только для чтения» и пишется. Как задать `nextPositions` у позиций, которых ещё нет (Create),
  источник не описывает: практически — `UpdatePosition` или `Update` с `Meta` позиций после создания.

## Ограничения и противоречия источника

- `positions` таблица помечает «Необходимо при создании» и не помечает readonly: позиции пишутся в
  Create (обязательно) и Update. Лимит 1–100 позиций — из текста источника, проверяется до запроса
  в Create/Update; общий лимит при CreatePositions проверяет сервер.
- `nextPositions` есть в таблице «Позиции Техпроцесса», но отсутствует в таблице «Атрибуты позиции
  Техпроцесса» и во всех примерах. Тип в таблице — MetaArray; SDK использует массив `{meta}`
  (`[]model.Reference`), как у MetaArray `performers` этапа производства в примерах источника.
  Поле добавлено в changelog 22-04-2024.
- Пример Update адресует существующую позицию её `meta`, хотя `meta`/`id` позиции помечены readonly;
  SDK следует примеру (`PositionUpdateRequest.Meta`), `id` позиции не пишет.
- Ответ массового удаления позиций в источнике без тела; `DeleteManyPositions` не возвращает
  результат по элементам.
- Expand, metadata и корзины нет (ticket их не называет). `Delete` удаляет безвозвратно.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingprocess),
`md/dictionaries/_processingprocess.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
