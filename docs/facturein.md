# FactureIn

`entity/facturein`: модель `FactureIn`, `CreateRequest`, `UpdateRequest`, `NewService(client)`,
`Get/List/Create/Update/Delete`, `CreateMany/UpdateMany/DeleteMany`, `New` (шаблон) и `Trash`.
Позиций у документа нет. [Общие правила методов, JSON, ошибок](documents.md), шаблоны, корзина и
связи — [documents-extra.md](documents-extra.md).

Минимальный Create SDK: `organization`. Счёт-фактура создаётся только на основании — ровно один
документ в `supplies` (приёмка) или `payments` (исходящий платёж); это проверяет сервер.
`incomingNumber` и `incomingDate` из основания не заполняются.

```go
created, err := s.Create(ctx, facturein.CreateRequest{
	Organization:   organization,
	Supplies:       model.Value([]model.Reference{supply}),
	IncomingNumber: model.Value("356124365"),
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
| `organization` | `*model.Reference` | required | Field |
| `agent` | `*model.Reference` | Field | Field |
| `contract` | `*model.Reference` | Field | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `attributes` | `[]model.Attribute` | Field | Field |
| `rate` | `*model.DocumentRate` | Field | Field |
| `sum` | `float64` (копейки) | — | — |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `supplies` | `[]model.Reference` | Field | Field |
| `payments` | `[]model.Reference` | Field | Field |
| `incomingNumber` | `string` | Field | Field |
| `incomingDate` | `model.DateTime` | Field | Field |

## Шаблон и корзина

- `New(ctx, basis)` — `PUT entity/facturein/new`. Основания по общей таблице шаблонов — приёмка и
  исходящий платёж; пример раздела передаёт их массивом связей: `{"supplies": []model.Reference{…}}`
  или `{"payments": […]}`. Ответ без id, документ не создаётся.
- `Trash(ctx, id)` — `POST entity/facturein/{id}/trash`; `Delete` удаляет безвозвратно.

## Ограничения и противоречия источника

- `organization` помечен «Необходимо при создании», а пример «только необходимые поля» его не
  передаёт (видимо, берётся из основания; источник этого не говорит). SDK следует таблице: пустой organization отклоняется до
  запроса. `CreateMany` обязательные поля не проверяет.
- `agent` и `contract` маркера «Только для чтения» не имеют — пишутся как Field.
- `files` маркера readonly не имеет, но, как у остальных документов SDK, не пишется (файлы — пакет
  `files`).
- Пример шаблона запрашивает основание `supplies`, а ответ показывает `payments` — SDK тело ответа
  не интерпретирует.
- Наличие и единственность основания, совместимость типов связей и отвязка пустым массивом `[]` —
  на стороне сервера; источник последнее не описывает.
- Expand, метаданные и доп. поля (`metadata`) не входят.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/facturein),
revision [`0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_facturein.md).
