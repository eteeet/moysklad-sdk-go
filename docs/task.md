# Task

Пакет `entity/task`: задачи (`entity/task`), комментарии задачи
(`entity/task/{id}/notes`, meta.type `tasknote`) и типы задач
(`entity/task/metadata/states`, meta.type `state`). `task.NewService(client)` принимает
`*transport.Client` и отклоняет nil. Предлагаемый handle фасада — `c.Entity().Task()`.
Создание, изменение и удаление требуют тарифной опции CRM; опцию и права проверяет сервер.

## Методы

| Метод | HTTP | Ответ |
| --- | --- | --- |
| `Get(ctx, id)` | GET `entity/task/{id}` | `*Task` |
| `List(ctx, query.ListOptions)` | GET `entity/task` | `*model.Collection[Task]` |
| `Create(ctx, CreateRequest)` | POST `entity/task` | `*Task` |
| `Update(ctx, id, UpdateRequest)` | PUT `entity/task/{id}` | `*Task` |
| `Delete(ctx, id)` | DELETE `entity/task/{id}` | `error`; 200 или 204 |
| `CreateMany(ctx, []CreateRequest)` / `UpdateMany(ctx, []batch.Update[UpdateRequest])` | POST `entity/task` | `[]batch.Result[Task]` |
| `DeleteMany(ctx, ids)` | POST `entity/task/delete` | `[]batch.Deleted` |
| `ListComments(ctx, taskID, CommentsOptions)` | GET `entity/task/{id}/notes` | `*model.Collection[Comment]` |
| `GetComment(ctx, taskID, commentID)` | GET `…/notes/{commentID}` | `*Comment` |
| `CreateComment(ctx, taskID, CommentRequest)` | POST `…/notes` | `*Comment` (объект или массив из одного) |
| `CreateManyComments(ctx, taskID, []CommentRequest)` | POST `…/notes` с массивом | `[]batch.Result[Comment]` |
| `UpdateComment(ctx, taskID, commentID, CommentRequest)` | PUT `…/notes/{commentID}` | `*Comment` |
| `DeleteComment(ctx, taskID, commentID)` | DELETE `…/notes/{commentID}` | `error` |
| `DeleteManyComments(ctx, taskID, ids)` | POST `…/notes/delete` | `error`; 200/204 с пустым телом |
| `ListStates(ctx)` | GET `entity/task/metadata` (поле `states`) | `[]metadata.State` |
| `CreateState(ctx, StateCreateRequest)` | POST `entity/task/metadata/states` | `*metadata.State` |
| `UpdateState(ctx, id, metadata.StateUpdateRequest)` | PUT `…/states/{id}` | `*metadata.State` |
| `DeleteState(ctx, id)` | DELETE `…/states/{id}` | `error` |
| `CreateManyStates(ctx, []StateCreateRequest)` / `UpdateManyStates(ctx, []batch.Update[metadata.StateUpdateRequest])` | POST `…/states` | `[]batch.Result[metadata.State]` |

Типы задач — методы `task.Service`, а не отдельный сервис: их маршруты лежат под
`entity/task/metadata`, а модель и маршруты совпадают с общими статусами
`metadata.Service` (entityType `task`). ListStates, UpdateState, DeleteState и
UpdateManyStates делегируют `metadata.Service`; CreateState и CreateManyStates свои, потому
что у типа задачи stateType при создании не нужен (у статусов документов — нужен).
`c.Entity().Metadata().CreateState("task", …)` тоже работает, но требует stateType.

Все проверки входа (UUID, обязательные поля, пустой Update, limit) — до HTTP,
ResponseInfo сбрасывается. Каждый вызов — одна отправка; POST/PUT/DELETE не повторяются.
Expand нет.

## Поля задачи

| JSON | Go type | Create | Update |
| --- | --- | --- | --- |
| `meta` | `model.Meta` | — | — |
| `id` | `string` | — | — |
| `accountId` | `string` | — | — |
| `description` | `string` | required `string` | Field |
| `assignee` | `*model.Reference` | required `model.Reference` | Field |
| `agent` | `*model.Reference` | Field | Field |
| `operation` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `dueToDate` | `model.DateTime` | Field | Field |
| `done` | `bool` | Field | Field |
| `author` | `*model.Reference` | — | — |
| `authorApplication` | `*model.Reference` | — | — |
| `implementer` | `*model.Reference` | — | — |
| `created` | `model.DateTime` | — | — |
| `updated` | `model.DateTime` | — | — |
| `completed` | `model.DateTime` | — | — |
| `notes` | `*model.Collection[Comment]` | — | — |
| `files` | `*model.Collection[model.Reference]` | — | — |

Create отклоняет пустой description и assignee без meta.href до HTTP. Field — `omitzero`:
незаданные не отправляются, `model.Value(false)` у done уходит явным false, `model.Null`
— null (принятие null решает сервер). Agent (контрагент или юрлицо) и operation (документ)
— ссылки; источник допускает привязку только к одному из них, SDK это не проверяет.
Files только читаются как metadata коллекции (файлы — пакет `files`); notes не записываются
(комментарии — отдельные методы). Без фильтров List отдаёт только активные задачи
(`done=false`), доступные сотруднику; фильтры assignee, author, done, state.name и поля
таблицы — `query.Filter`.

## Комментарии

| JSON | Go type | Запись |
| --- | --- | --- |
| `meta` | `model.Meta` | — |
| `id`, `accountId` | `string` | — (есть в примерах ответа, не в таблице) |
| `text` | `string` | `CommentRequest.Text`, необходим при создании |
| `moment` | `model.DateTime` | — |
| `author`, `authorApplication` | `*model.Reference` | — |
| `files` | `*model.Collection[model.Reference]` | — (есть в примере создания) |

Упоминания `{{employee;uuid}}` — часть строки text, без разбора. `CommentsOptions`: Limit
1–100 (0 — default сервера 25), Offset, Filter (`query.Filter` по updated и updatedBy).
С задачей хранится не больше 100 комментариев; остальные читает ListComments.

## Типы задач

`metadata.State`: `meta`, `id`, `accountId`, `name`, `color` (ARGB целым числом:
rgb(162, 198, 23) → 10667543), `stateType` (открытая строка), `entityType` (всегда task).
`StateCreateRequest`: `Name` и `Color *int64` необходимы (nil color отклоняется, 0 —
допустимое значение), `StateType` — Field, по умолчанию не отправляется. Массового
удаления типов источник не описывает — метода нет.

## Противоречия источника

- Таблица атрибутов комментария называет текст `description`; таблицы Create/Update и все
  примеры — `text`. SDK использует `text`.
- Примеры Get/Update комментария не содержат `id`, пример создания (2025) и фикстура
  официального SDK — содержат. Успех GetComment/UpdateComment требует id (общий критерий
  `internal/nested`); ответ без id — ошибка.
- Таблица списка комментариев: default limit 25, максимум 100; пример ответа показывает
  `limit: 1000`. SDK ограничивает Limit по таблице и не подставляет default.
- Параметры updatedBy/updatedFrom/updatedTo источник называет параметрами фильтрации
  выборки, а не отдельными query-параметрами; SDK передаёт их через Filter.
- Таблица типа задачи не помечает stateType «только для чтения», а описание создания и
  изменения называет только name и color: stateType оставлен Field.
- Пример ответа Create комментария — массив из одного объекта: принимается и объект.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/task),
[revision 0e1b3d61](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_task.md).
