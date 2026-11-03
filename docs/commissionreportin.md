# CommissionReportIn

`entity/commissionreportin`: модель `CommissionReportIn`, `CreateRequest`, `UpdateRequest`,
`NewService(client)`; методы `Get/List/Create/Update/Delete`, `CreateMany/UpdateMany/DeleteMany`
и `Trash(ctx, id)`. Два набора строк собственной формы (не `positions.Position`):
«реализовано комиссионером» — `Position`/`PositionRequest`, «возврат на склад комиссионера» —
`ReturnPosition`/`ReturnPositionRequest`.
[Общие правила методов, JSON и ошибок](documents.md), [корзина](documents-extra.md).

Минимальный Create: `organization`, `agent`, `contract` (договор комиссии с этим `agent`),
`commissionPeriodStart`, `commissionPeriodEnd` — проверяются до запроса (непустой meta.href,
ненулевая дата). Счета `organizationAccount`/`agentAccount` нужны, если у стороны несколько счетов, —
проверяет сервер.

```go
doc, err := s.Create(ctx, commissionreportin.CreateRequest{
	Organization: organization, Agent: agent, Contract: contract,
	CommissionPeriodStart: start, CommissionPeriodEnd: end,
	RewardType: model.Value("PercentOfSales"), RewardPercent: model.Value(15),
	Positions: model.Value([]commissionreportin.PositionRequest{{
		Assortment: model.Value(product), Quantity: model.Value(2.0), Price: model.Value(13205.0), Reward: model.Value(123.0),
	}}),
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
| `syncId` | `string` | — | — |
| `created` | `model.DateTime` | — | — |
| `updated` | `model.DateTime` | — | — |
| `deleted` | `model.DateTime` | — | — |
| `moment` | `model.DateTime` | Field | Field |
| `applicable` | `bool` | Field | Field |
| `shared` | `bool` | Field | Field |
| `printed` | `bool` | — | — |
| `published` | `bool` | — | — |
| `organization` | `*model.Reference` | required | Field |
| `organizationAccount` | `*model.Reference` | Field | Field |
| `agent` | `*model.Reference` | required | Field |
| `agentAccount` | `*model.Reference` | Field | Field |
| `contract` | `*model.Reference` | required | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `project` | `*model.Reference` | Field | Field |
| `salesChannel` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `attributes` | `[]model.Attribute` | Field | Field |
| `rate` | `*model.DocumentRate` | Field | Field |
| `vatEnabled` | `bool` | Field | Field |
| `vatIncluded` | `bool` | Field | Field |
| `sum` | `float64` | — | — |
| `vatSum` | `float64` | — | — |
| `payedSum` | `float64` | — | — |
| `commitentSum` | `float64` | — | — |
| `commissionPeriodStart` | `model.DateTime` | required | Field |
| `commissionPeriodEnd` | `model.DateTime` | required | Field |
| `rewardType` | `string` | Field | Field |
| `rewardPercent` | `int` | Field | Field |
| `commissionOverhead` | `*CommissionOverhead` | Field | Field |
| `positions` | `*model.Collection[Position]` | Field | Field |
| `returnToCommissionerPositions` | `*model.Collection[ReturnPosition]` | Field | Field |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `payments` | `[]model.Reference` | Field | Field |

Суммы (`sum`, `vatSum`, `payedSum`, `commitentSum`, `commissionOverhead.sum`, `price`, `reward`) —
`float64` в копейках; `rewardPercent` — целый процент (таблица: Int), не деньги. `rewardType` —
открытый enum: `PercentOfSales`, `None`. `commissionOverhead.sum` и `reward` строки могут быть
отрицательными; прочие расходы задаются только при заданных позициях (проверяет сервер).

## Позиции

| Набор | Ресурс | meta.type | Методы |
| --- | --- | --- | --- |
| реализовано | `…/{id}/positions` | `commissionreportinposition` | `GetPositions`, `GetPosition`, `CreatePositions`, `UpdatePosition`, `DeletePosition`, `DeleteManyPositions` |
| возврат | `…/{id}/returntocommissionerpositions` | `commissionreportinreturnedposition` | `GetReturnPositions`, `GetReturnPosition`, `CreateReturnPositions`, `UpdateReturnPosition`, `DeleteReturnPosition`, `DeleteManyReturnPositions` |

- `GetPositions(ctx, id, positions.PageOptions) (*model.Collection[Position], error)` и
  `GetReturnPositions(ctx, id, positions.PageOptions) (*model.Collection[ReturnPosition], error)` —
  одна страница, только limit/offset (`fields` отклоняется до запроса).
- `CreatePositions`/`CreateReturnPositions` — один POST массивом, ответ `[]batch.Result[…]`.
  Для новой строки сервер требует `assortment` и положительный `quantity`.
- `UpdatePosition`/`UpdateReturnPosition` принимают тот же request; пустой отклоняется до запроса.
- `DeleteManyPositions`/`DeleteManyReturnPositions` — POST `…/delete` с meta строк (1–1000);
  источник отвечает 200 с пустым телом, поэтому возвращается только error.
- Поля строки: `meta`, `id`, `accountId` (только чтение), `assortment`, `quantity`, `price`, `vat`,
  `vatEnabled`, `reward`; у «реализовано» ещё `pack` (`positions.Pack`/`positions.PackRequest`).
- Positions и ReturnToCommissionerPositions в Create/Update — не больше 1000 строк каждый
  (1001 отклоняется до запроса); в Update массив заменяет весь набор. `PositionRequest` не имеет
  `id`: строки массива в Update — новые; адресная правка — `UpdatePosition`.

## Ограничения и противоречия источника

- `syncId` таблица помечает «Только для чтения» (у других документов он пишется) — SDK его не пишет.
- Примеры Create/Update передают `sum` и `vatSum`, таблица помечает их readonly — SDK их не пишет.
- `payments` (раздел «Связи с другими документами») не помечен readonly и пишется по общему правилу
  привязки документов; `linkedSum` элемента SDK не читает (как у других документов, `model.Reference`).
- `commissionOverhead` в примерах отсутствует; форма `{sum}` — по таблице «Прочие расходы».
- `rewardPercent` в таблице Int — SDK держит `int`, как `contract`.
- Шаблона `New`, wizard и expand нет (ticket их не называет); Delete удаляет безвозвратно,
  Trash — в корзину.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportin),
`md/documents/_commissionreportin.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
