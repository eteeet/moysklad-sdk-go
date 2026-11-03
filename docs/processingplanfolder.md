# ProcessingPlanFolder

`entity/processingplanfolder`: группа техкарт. `NewService(client)` принимает
`*transport.Client`, возвращает `(*Service, error)` и отклоняет nil client.

## API

| Метод | Запрос |
| --- | --- |
| `Get(ctx, id) (*ProcessingPlanFolder, error)` | `GET entity/processingplanfolder/{id}` |
| `List(ctx, query.ListOptions) (*model.Collection[ProcessingPlanFolder], error)` | `GET entity/processingplanfolder`, одна страница, без expand |
| `Create(ctx, CreateRequest) (*ProcessingPlanFolder, error)` | `POST entity/processingplanfolder` |
| `Update(ctx, id, UpdateRequest) (*ProcessingPlanFolder, error)` | `PUT entity/processingplanfolder/{id}` |
| `Delete(ctx, id) error` | `DELETE entity/processingplanfolder/{id}`, не в корзину |
| `CreateMany(ctx, []CreateRequest) ([]batch.Result[ProcessingPlanFolder], error)` | `POST entity/processingplanfolder` с массивом |
| `UpdateMany(ctx, []batch.Update[UpdateRequest]) ([]batch.Result[ProcessingPlanFolder], error)` | `POST entity/processingplanfolder` с массивом, meta по ID |
| `DeleteMany(ctx, ids) ([]batch.Deleted, error)` | `POST entity/processingplanfolder/delete` |
| `metadata.Service.Get(ctx, "processingplanfolder")` | `GET entity/processingplanfolder/metadata` — общий пакет `metadata`, своих методов у пакета нет |

Успех Get/Create/Update — HTTP 200 и непустой `id`; возвращается фактический ответ сервера.
id проверяется как UUID до HTTP; Update без единого заданного поля отклоняется.
Ошибки начинаются с `moysklad:` один раз и содержат операцию; `errors.As` находит
`*sdkerrors.APIError`, `errors.Is` сохраняет ошибки context, авторизации и транспорта.
Каждый метод отправляет один запрос без автоматических повторов; после сбоя Create
с неизвестным исходом проверь аккаунт перед повторной отправкой. Массовые операции —
по правилам [write-async.md](write-async.md), List — [transport.md](transport.md).

Минимальный Create:

```go
folder, err := svc.Create(ctx, processingplanfolder.CreateRequest{Name: "Группа 1"})
```

## Поля

| JSON | Read Go type | Create | Update |
| --- | --- | --- | --- |
| `meta` | `model.Meta` | — | — |
| `id` | `string` | — | — |
| `accountId` | `string` | — | — |
| `updated` | `model.DateTime` | — | — |
| `name` | `string` | required | Field |
| `description` | `string` | Field | Field |
| `code` | `string` | Field | Field |
| `externalCode` | `string` | Field | Field |
| `archived` | `bool` | Field | Field |
| `pathName` | `string` | — | — |
| `shared` | `bool` | Field | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |

`required` — обычное значение; пустой `Name` отклоняется до HTTP, без trim.
`Field` — `model.Field[T]` с `omitzero`: незаданное не отправляется, `model.Value(false)`,
`model.Value("")` и `model.Null[T]()` уходят в JSON. Кодирование null не обещает очистку поля сервером.
`accountId`, `id`, `pathName`, `updated` таблица помечает «Только для чтения» — в запросах их нет.

## Родительская группа и ограничения

- `pathName` — имя родительской группы, только чтение (как у productfolder).
- Ссылки на родительскую группу в таблице атрибутов нет, и ни один пример её не передаёт и не
  возвращает. Раздел «Изменить Группу техкарт» говорит, что для смены `pathName` нужно обновить
  поле `processingplanfolder`, — это единственное упоминание, без типа и без признака
  чтения/записи. SDK это поле не читает и не пишет (у productfolder `productFolder` есть в
  таблице, здесь — нет). Перенос группы в другую родительскую через SDK недоступен.
- Метаданные источника содержат только `meta`; доп. полей и статусов в примере нет.
- Expand (`owner`, `group`) не поддерживается.

JSON проверен через `encoding/json` и `encoding/json/v2`; HTTP — настоящим transport и
сервером в памяти. Фикстуры: `testdata/custom/processingplanfolder/read.json` и `metadata.json` —
примеры «Получить Группу техкарт» и «Метаданные Групп техкарт» без изменений;
`create-full.golden.json`, `update.golden.json` — синтетические ожидаемые тела запросов.
Живой аккаунт не использовался.

Источник: [официальная документация](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_processingplanfolder.md),
revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
