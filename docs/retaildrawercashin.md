# RetailDrawerCashIn

`entity/retaildrawercashin`: модель `RetailDrawerCashIn`, `CreateRequest`, `UpdateRequest`,
`NewService(client)`; методы `Get/List/Create/Update/Delete`, `CreateMany/UpdateMany/DeleteMany`,
`New(ctx, basis)` (шаблон `PUT entity/retaildrawercashin/new`) и `Trash(ctx, id)`. Позиций нет.
[Общие правила методов, JSON и ошибок](documents.md), [шаблоны и корзина](documents-extra.md).

Минимальный Create: `organization`, `agent` (сотрудник, совершивший внесение) и `retailShift` —
пустой `meta.href` любого из них отклоняется до запроса. Шаблон на основе розничной смены:
`New(ctx, map[string]any{"retailShift": shiftRef})`.

```go
doc, err := s.Create(ctx, retaildrawercashin.CreateRequest{
	Organization: org, Agent: employee, RetailShift: shift,
	Description:  model.Value("Внесение через API"),
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
| `organizationBranch` | `*model.Reference` | Field | Field |
| `agent` | `*model.Reference` | required | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `attributes` | `[]model.Attribute` | Field | Field |
| `rate` | `*model.DocumentRate` | Field | Field |
| `sum` | `float64` (копейки) | Field | Field |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `retailShift` | `*model.Reference` | required | Field |

## Ограничения и противоречия источника

- `sum` таблица атрибутов помечает «Только для чтения», а примеры Create, массового создания и
  Update передают его (`"sum": 100500`, `"sum": 700`). SDK пишет `sum` в Create и Update:
  исключение из правила «readonly важнее примеров» утверждено пользователем, без него сумму внесения не задать.
- `retailShift` описан в разделе «Связи с другими документами» с пометкой «Необходимое» и в
  списке обязательных полей Create; пишется и в Update (маркера «Только для чтения» нет).
- `code` в таблице внесения нет (у выплаты есть) — SDK его не читает и не пишет.
- Delete удаляет безвозвратно, Trash — в корзину (восстановление только в интерфейсе).
- Expand не поддержан.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retaildrawercashin),
`md/documents/_retaildrawercashin.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
