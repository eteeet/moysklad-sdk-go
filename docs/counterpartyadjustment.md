# CounterpartyAdjustment

`entity/counterpartyadjustment`: модель `CounterpartyAdjustment`, `CreateRequest`, `UpdateRequest`,
`NewService(client)`, `Get/List/Create/Update/Delete`, `CreateMany/UpdateMany/DeleteMany` и `Trash`.
Позиций и шаблона (`New`) у документа нет: ticket шаблон не называет, раздела шаблона в источнике
нет. [Общие правила методов, JSON, ошибок](documents.md), корзина —
[documents-extra.md](documents-extra.md).

Минимальный Create: `organization`, `agent` (контрагент или сотрудник).

```go
created, err := s.Create(ctx, counterpartyadjustment.CreateRequest{Organization: organization, Agent: counterparty})
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
| `created` | `model.DateTime` | — | — |
| `updated` | `model.DateTime` | — | — |
| `deleted` | `model.DateTime` | — | — |
| `moment` | `model.DateTime` | Field | Field |
| `applicable` | `bool` | Field | Field |
| `shared` | `bool` | Field | Field |
| `printed` | `bool` | — | — |
| `published` | `bool` | — | — |
| `organization` | `*model.Reference` | required | Field |
| `agent` | `*model.Reference` | required | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `attributes` | `[]model.Attribute` | Field | Field |
| `sum` | `float64` (копейки) | — | — |
| `files` | `*model.Collection[model.Reference]` | — | — |

`code`, `syncId`, `rate`, `state`, `project` в таблице атрибутов нет — SDK их не читает и не пишет.

## Сумма

- `sum` — «Сумма Корректировки взаиморасчетов в копейках», `float64`, как у остальных документов.
- Знак суммы источник не описывает: все примеры положительные (12.5, 12.0, 21300.0, 32, 123),
  что означает знак (увеличение или уменьшение долга), не сказано. SDK значение не меняет и не
  проверяет: отрицательная сумма из ответа читается как есть.
- **Противоречие источника:** таблица атрибутов помечает `sum` «Только для чтения», текст «Изменить»
  запрещает передавать такие поля, но примеры «Создать», «Массовое создание» и «Изменить» передают
  `sum`. SDK следует таблице: `sum` нет в `CreateRequest` и `UpdateRequest`, поэтому задать сумму
  корректировки через SDK нельзя. Если живой API принимает `sum`, поле нужно добавить отдельным
  решением.

## Корзина

`Trash(ctx, id)` — `POST entity/counterpartyadjustment/{id}/trash` по общему разделу «Удаление в
корзину»; `Delete` удаляет безвозвратно.

## Ограничения

- `files` маркера readonly не имеет, но, как у остальных документов SDK, не пишется (файлы — пакет
  `files`).
- Тип `agent` (контрагент или сотрудник) SDK не ограничивает. `CreateMany` обязательные поля не
  проверяет.
- Expand, метаданные и доп. поля (`metadata`) не входят.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/counterpartyadjustment),
revision [`0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_counterpartyadjustment.md).
