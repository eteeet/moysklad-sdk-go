# Публикации документов

Пакет `publication`. Источник — `api-remap-1.2-doc` revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`:
[md/documents/_publication.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_publication.md).
`publication.NewService(client)` отклоняет nil client. Все методы принимают
`opts ...transport.RequestOption`, первой строкой сбрасывают `ResponseInfo` и проверяют аргументы до HTTP.

| Метод `publication.Service` | Endpoint | Успех |
| --- | --- | --- |
| `List(ctx, type, id)` | GET entity/{type}/{id}/publication | 200, meta.href, rows — `*model.Collection[Publication]` |
| `Get(ctx, type, id, publicationID)` | GET entity/{type}/{id}/publication/{publicationId} | 200, meta.href |
| `Create(ctx, type, id, CreateRequest)` | POST entity/{type}/{id}/publication | 201 — новая, 200 — уже существующая |
| `Delete(ctx, type, id, publicationID)` | DELETE entity/{type}/{id}/publication/{publicationId} | 204 (или 200) |

- `type` — строка формы `^[a-z][a-z0-9]*$` (как в `metadata`), `id` и `publicationID` — UUID.
  Источник перечисляет типы с публикациями (заказы, счета, отгрузка, приёмка, платежи и ордера,
  перемещение, оприходование, списание, счета-фактуры, возвраты, внесение/выплата, розница, договор,
  розничная смена, заказ на производство, отчёты комиссионера, инвентаризация, техоперация) —
  проверяет сервер. Нужны права на чтение и печать типа.
- Create отличает новую публикацию от существующей только статусом: он виден в
  `transport.ResponseInfo`. Запись и удаление не повторяются.
- Параметров списка источник не описывает.

## Поля

| JSON | Go type | Запись |
| --- | --- | --- |
| meta | `model.Meta` (type `operationpublication`) | нет |
| template | `model.Reference` — шаблон печати ([print.md](print.md)) | `CreateRequest.Template`, обязательно (meta.href) |
| href | `string` — адрес страницы публикации | нет, только чтение |

## Противоречия источника

- meta.type списка в примере — тип документа (`demand`), а не publication.
- `template` помечен `+Expand`; expand в SDK нет, Template — ссылка.
- Параметр id для Create описан как «id сущности, по которой получить Публикации» — копия описания List.

Фикстуры — `testdata/custom/publication` (README с SHA-256), тело запроса — `create.golden.json`.
