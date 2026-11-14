# Eventfeed

Пакет `eventfeed` — лента событий документа: `entity/{type}/{id}/notes`, meta.type
`eventnote`. `eventfeed.NewService(client)` принимает `*transport.Client` и отклоняет nil.
Handle фасада — `c.Entity().EventFeed()` (как `Entity().Files()`: вложенный
ресурс нескольких типов сущностей).

События контрагента (`entity/counterparty/{id}/notes`, meta.type `note`) — другой ресурс:
`counterparty.Service.ListNotes/GetNote/CreateNote/UpdateNote/DeleteNote`; этот пакет их не
дублирует. Комментарии задачи — `task.Service` ([task.md](task.md)).

## Методы

| Метод | HTTP | Ответ |
| --- | --- | --- |
| `List(ctx, documentType, documentID, query.PageOptions)` | GET `entity/{type}/{id}/notes?limit&offset` | `*model.Collection[Note]` |
| `Get(ctx, documentType, documentID, noteID)` | GET `…/notes/{noteID}` | `*Note` |
| `Create(ctx, documentType, documentID, NoteRequest)` | POST `…/notes` | `*Note` (объект или массив из одного) |
| `Update(ctx, documentType, documentID, noteID, NoteRequest)` | PUT `…/notes/{noteID}` | `*Note` (объект или массив из одного) |
| `Delete(ctx, documentType, documentID, noteID)` | DELETE `…/notes/{noteID}` | `error`; 200 или 204 |
| `CreateMany(ctx, documentType, documentID, []NoteRequest)` | POST `…/notes` с массивом | `[]batch.Result[Note]` |
| `UpdateMany(ctx, documentType, documentID, []batch.Update[NoteRequest])` | POST `…/notes` с массивом meta+полей | `[]batch.Result[Note]` |
| `DeleteMany(ctx, documentType, documentID, ids)` | POST `…/notes/delete` | `error`; 200/204 с пустым телом |

Источник перечисляет документы: `purchasereturn`, `commissionreportout`, `customerorder`,
`purchaseorder`, `supply`, `invoiceout`, `invoicein`. `documentType` — строка; до HTTP
проверяется только форма `[a-z][a-z0-9]*` (путь не подменить), допустимость типа решает
сервер. `documentID`, `noteID` и ids — форма UUID. List принимает только limit (1–1000,
0 — 1000) и offset; fields отклоняется. Все проверки — до HTTP, ResponseInfo сбрасывается.
Каждый вызов — одна отправка, POST/PUT/DELETE не повторяются.

## Поля

| JSON | Go type | Запись |
| --- | --- | --- |
| `meta` | `model.Meta` | — |
| `id` | `string` | — (readonly) |
| `accountId` | `string` | — (readonly) |
| `created` | `model.DateTime` | — (readonly) |
| `description` | `string` | `NoteRequest.Description`, необходим при создании |
| `author` | `*model.Reference` | — (readonly) |
| `authorApplication` | `*model.Reference` | — (readonly; есть у событий решения) |

Упоминания в тексте (`{{employee;uuid}}`) — часть строки `description`: SDK не разбирает и
не проверяет их. Пустой description в Create/Update отклоняется до HTTP; длину (4096) и лимит
5000 событий на документ проверяет сервер. Права (просмотр документа, автор/администратор
для изменения и удаления) — сервер.

## Ограничения и противоречия

- Create и массовое создание/обновление: источник указывает `Response 201` (у остальных
  записей SDK — 200). Для ленты событий SDK принимает 201 и 200 (`nested.Collection.Created`);
  у других коллекций 201 по-прежнему ошибка.
- Пример ответа одиночного Create — массив из одного события: принимается и объект, и такой
  массив; массив другого размера — ошибка.
- Массовое создание и обновление в источнике — один POST со смешанным массивом; SDK
  разделяет его на CreateMany и UpdateMany (как у остальных сервисов).
- В примере списка обе строки с одинаковым id — пример, не контракт.
- Expand (`author`, `authorApplication`) не поддерживается.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/eventfeed),
[revision 0e1b3d61](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_eventfeed.md).
