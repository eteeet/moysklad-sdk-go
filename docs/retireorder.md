# RetireOrder

`entity/retireorder`: модель `RetireOrder`, `CreateRequest`, `UpdateRequest`, `NewService(client)`;
методы `Get/List/Create/Update`, `CreateMany/UpdateMany` и `Trash(ctx, id)`. Позиции — собственная
форма `Position`/`PositionRequest` (meta.type `retireorderposition`).
Handle фасада: `c.Entity().RetireOrder()`.
[Общие правила методов, JSON и ошибок](documents.md), [корзина](documents-extra.md).

Минимальный Create: `organization`, `retireOrderType`, `trackingType` — проверяются до запроса
(непустой meta.href, непустые строки). Значения enum SDK не проверяет.

```go
doc, err := s.Create(ctx, retireorder.CreateRequest{
	Organization: organization, RetireOrderType: "RETAIL_SALE", TrackingType: "SHOES",
	SupportingTransaction: model.Value("RECEIPT"), SupportingTransactionNumber: model.Value("123"),
	SupportingTransactionDate: model.Value(date),
	Positions: model.Value([]retireorder.PositionRequest{{
		ItemRequest:   positions.ItemRequest{Assortment: model.Value(product)},
		TrackingCodes: model.Value([]positions.TrackingCode{{CIS: cis, Type: "trackingcode"}}),
	}}),
})
```

## Поддержанные поля

| JSON | Read Go type | Create | Update |
| --- | --- | --- | --- |
| `meta` | `model.Meta` | — | — |
| `id` | `string` | — | — |
| `accountId` | `string` | — | — |
| `syncId` | `string` | Field | Field |
| `name` | `string` | Field | Field |
| `description` | `string` | Field | Field |
| `externalCode` | `string` | Field | Field |
| `created` | `model.DateTime` | — | — |
| `updated` | `model.DateTime` | — | — |
| `deleted` | `model.DateTime` | — | — |
| `moment` | `model.DateTime` | Field | Field |
| `shared` | `bool` | Field | Field |
| `printed` | `bool` | — | — |
| `published` | `bool` | — | — |
| `organization` | `*model.Reference` | required | Field |
| `organizationBranch` | `*model.Reference` | Field | Field |
| `agent` | `*model.Reference` | Field | Field |
| `destinationCountry` | `*model.Reference` | Field | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `rate` | `*model.DocumentRate` | Field (`DocumentRateRequest`) | Field |
| `sum` | `float64` (копейки) | Field | Field |
| `vatEnabled` | `bool` | Field | Field |
| `vatIncluded` | `bool` | Field | Field |
| `vatSum` | `float64` (копейки) | — | — |
| `trackingType` | `string` (открытый enum) | required | Field |
| `retireOrderType` | `string` (открытый enum) | required | Field |
| `documentState` | `string` (открытый enum) | — | — |
| `reasonDescription` | `string` | Field | Field |
| `primaryDocumentName` | `string` | Field | Field |
| `stateContractId` | `string` | Field | Field |
| `supportingTransaction` | `string` (открытый enum) | Field | Field |
| `supportingTransactionDate` | `model.DateTime` | Field | Field |
| `supportingTransactionNumber` | `string` | Field | Field |
| `positions` | `*model.Collection[Position]` | Field | Field |

## Позиции

- `GetPositions(ctx, id, positions.PageOptions) (*model.Collection[Position], error)` — одна страница
  `…/{id}/positions`, только limit/offset (`fields` отклоняется до запроса).
- `GetPosition`, `CreatePositions` (один POST массивом, ответ `[]batch.Result[Position]`),
  `UpdatePosition` (пустой request — ошибка до запроса), `DeletePosition`, `DeleteManyPositions`
  (POST `…/positions/delete`, 1–1000 id; ответ без тела — только error).
- Поля позиции: `meta`, `id`, `accountId` (только чтение), `assortment` (`model.Reference`),
  `quantity` (`float64`), `price` (`float64`, копейки), `vat` (`int`, процент), `vatEnabled`,
  `trackingCodes` (`[]positions.TrackingCode`: `cis`, открытый `type`, вложенные коды).
- Для новой позиции источник требует `assortment` и `trackingCodes`: `CreatePositions` отклоняет
  позицию без них до запроса. Позиции без `ID` в теле Create/Update документа и `UpdatePosition`
  проверяет сервер (одна форма `PositionRequest` служит и созданию, и изменению).
- Positions в Create/Update — не больше 1000 позиций (1001 отклоняется до запроса); в Update массив
  заменяет весь набор: позиция с `ID` — существующая, без — новая.

## Ограничения и противоречия

- Delete/DeleteMany документа в источнике нет — методов нет; есть только `Trash` (общий раздел
  `documents/_common_info.md`, «Удаление в корзину»).
- Update допустим только при `documentState` CREATED, CHECKED_NOT_OK, PROCESSING_ERROR — проверяет
  сервер, SDK не делает GET перед PUT.
- Допустимые комбинации `trackingType`/`retireOrderType`/`supportingTransaction`/`destinationCountry`
  и условно обязательные поля (`supportingTransaction*`, `stateContractId`, `reasonDescription`,
  `primaryDocumentName`) определяет «Честный знак» — проверяет сервер.
- `sum` и `meta` в таблице не помечены «Только для чтения»; `sum` пишется как Field (примеры
  запросов его не передают, сервер может пересчитать по позициям). `meta` в теле не пишется: адрес
  задаёт URL, в UpdateMany — `batch.Update.ID`.
- В таблице позиций нет `meta`, но ответы его содержат — модель его читает.
- Ticket упоминает `cisType`, в источнике такого поля нет — SDK его не добавляет.
- `attributes`, `files`, `code`, `project` в таблице атрибутов нет. Expand, metadata и шаблон `New`
  не входят в ticket.
- List: источник описывает только limit/offset; SDK передаёт `query.ListOptions` как у остальных
  документов, поддержку filter/search определяет сервер.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/retireorder),
`md/documents/_retireorder.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
Фикстуры — примеры ответов источника `testdata/custom/retireorder/read.json` («Получить Вывод из
оборота») и `positions.json` («Получить позиции Вывода из оборота»).
