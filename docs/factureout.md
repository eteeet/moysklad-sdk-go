# FactureOut

`entity/factureout`: модель `FactureOut`, `CreateRequest`, `UpdateRequest`, `NewService(client)`,
`Get/List/Create/Update/Delete`, `CreateMany/UpdateMany/DeleteMany`, `New` (шаблон) и `Trash`.
Позиций у документа нет. [Общие правила методов, JSON, ошибок](documents.md), шаблоны, корзина и
связи — [documents-extra.md](documents-extra.md).

Минимальный Create SDK: `organization`. Счёт-фактура создаётся только на основании — ровно один
документ в `demands` (отгрузка), `returns` (возврат поставщику) или `payments` (входящий платёж,
приходный ордер); это проверяет сервер. `paymentNumber` и `paymentDate` из основания не заполняются.

```go
created, err := s.Create(ctx, factureout.CreateRequest{
	Organization:      organization,
	Payments:          model.Value([]model.Reference{paymentIn}),
	AdvancePaymentVat: model.Value(20),
	PaymentPurpose:    model.Value("Аванс по договору №123"),
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
| `stateContractId` | `string` | Field | Field |
| `advancePaymentVat` | `int` (процент) | Field | Field |
| `paymentPurpose` | `string` | Field | Field |
| `vatSum` | `float64` (копейки) | — | — |
| `demands` | `[]model.Reference` | Field | Field |
| `payments` | `[]model.Reference` | Field | Field |
| `returns` | `[]model.Reference` | Field | Field |
| `consignee` | `*model.Reference` | Field | Field |
| `consigneeBranch` | `*model.Reference` | Field | Field |
| `paymentNumber` | `string` | Field | Field |
| `paymentDate` | `model.DateTime` | Field | Field |

## Шаблон и корзина

- `New(ctx, basis)` — `PUT entity/factureout/new`. Основания — отгрузка, возврат поставщику и
  входящий платёж; примеры раздела передают их массивом связей: `{"demands": […]}`,
  `{"returns": […]}` или `{"payments": […]}`. Ответ без id, документ не создаётся.
- `Trash(ctx, id)` — `POST entity/factureout/{id}/trash`; `Delete` удаляет безвозвратно.

## Авансовые поля

`advancePaymentVat` и `paymentPurpose` сервер учитывает только при основании-платеже, при отгрузке
и возврате игнорирует и не возвращает; `vatSum` (readonly) он рассчитывает как
`sum * vatRate / (100 + vatRate)` или берёт из платежа. Ставку (0–100, из справочника; ошибки
3002, 3008, 3005) проверяет сервер. В модели чтения отсутствие этих полей даёт нули.

## Ограничения и противоречия источника

- `organization` помечен «Необходимо при создании», а пример «только необходимые поля» его не
  передаёт (видимо, берётся из основания; источник этого не говорит). SDK следует таблице: пустой organization отклоняется до
  запроса. `CreateMany` обязательные поля не проверяет.
- Текст «Создать» называет основанием и приходный ордер, а таблица связей описывает `payments` как
  входящие платежи; тип привязываемого документа SDK не ограничивает.
- `agent` и `contract` маркера «Только для чтения» не имеют — пишутся как Field.
- `files` маркера readonly не имеет, но, как у остальных документов SDK, не пишется (файлы — пакет
  `files`).
- Наличие и единственность основания, совместимость типов связей и отвязка пустым массивом `[]` —
  на стороне сервера; источник последнее не описывает.
- Expand, метаданные и доп. поля (`metadata`) не входят.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/factureout),
revision [`0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_factureout.md).
