# CommissionReportOut

`entity/commissionreportout`: модель `CommissionReportOut`, `CreateRequest`, `UpdateRequest`,
`NewService(client)`; методы `Get/List/Create/Update/Delete`, `CreateMany/UpdateMany/DeleteMany`
и `Trash(ctx, id)`. Позиции — собственная форма `Position`/`PositionRequest`
(meta.type `commissionreportoutposition`), не `positions.Position`.
[Общие правила методов, JSON и ошибок](documents.md), [корзина](documents-extra.md).

Минимальный Create: `organization`, `agent`, `contract` (договор комиссии с этим `agent`),
`commissionPeriodStart`, `commissionPeriodEnd` и `name` — проверяются до запроса (непустой
meta.href, ненулевая дата, непустое имя). Счета нужны, если у стороны несколько счетов, — проверяет
сервер.

```go
doc, err := s.Create(ctx, commissionreportout.CreateRequest{
	Organization: organization, Agent: agent, Contract: contract,
	CommissionPeriodStart: start, CommissionPeriodEnd: end, Name: "00001",
	RewardType: model.Value("PercentOfSales"), RewardPercent: model.Value(15),
	Positions: model.Value([]commissionreportout.PositionRequest{{
		Assortment: model.Value(product), Quantity: model.Value(2.0), Price: model.Value(13205.0),
	}}),
})
```

## Поддержанные поля

| JSON | Read Go type | Create | Update |
| --- | --- | --- | --- |
| `meta` | `model.Meta` | — | — |
| `id` | `string` | — | — |
| `accountId` | `string` | — | — |
| `name` | `string` | required | Field |
| `description` | `string` | Field | Field |
| `code` | `string` | Field | Field |
| `externalCode` | `string` | Field | Field |
| `syncId` | `string` | — | — |
| `created` | `model.DateTime` | — | — |
| `updated` | `model.DateTime` | — | — |
| `deleted` | `model.DateTime` | — | — |
| `moment` | `model.DateTime` | Field | Field |
| `applicable` | `bool` | Field | Field |
| `shared` | `bool` | Field | — |
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
| `positions` | `*model.Collection[Position]` | Field | Field |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `payments` | `[]model.Reference` | Field | Field |

Суммы и `price`/`reward` позиций — `float64` в копейках; `rewardPercent` — целый процент (таблица:
Int), не деньги; `rewardType` — открытый enum `PercentOfSales`, `None`.

## Позиции

- `GetPositions(ctx, id, positions.PageOptions) (*model.Collection[Position], error)` — одна страница
  `…/{id}/positions`, только limit/offset (`fields` отклоняется до запроса).
- `GetPosition`, `CreatePositions` (один POST массивом, ответ `[]batch.Result[Position]`),
  `UpdatePosition` (пустой request — ошибка до запроса), `DeletePosition`, `DeleteManyPositions`
  (POST `…/positions/delete`, 1–1000 id; ответ без тела — только error).
- Поля позиции: `meta`, `accountId` (только чтение), `id`, `assortment`, `quantity`, `price`, `vat`,
  `vatEnabled`, `reward`, `pack`. Для новой позиции сервер требует `assortment` и положительный
  `quantity`.
- Positions в Create/Update — не больше 1000 позиций (1001 отклоняется до запроса); в Update массив
  заменяет весь набор: позиция с `ID` — существующая (id только адресует её), без — новая.

## Ограничения и противоречия источника

- Таблица помечает `positions`, `attributes`, `shared` и `syncId` «Только для чтения», а текст
  «Изменить» запрещает передавать такие поля. При этом текст раздела позиций разрешает передавать
  позиции в теле документа, пример Create передаёт `positions`, `attributes` и `shared`, пример
  Update — `positions` и `attributes`. SDK пишет эти поля там, где их передаёт пример: `positions` и
  `attributes` — в Create и Update, `shared` — только в Create; исключение из правила «readonly важнее примеров» утверждено пользователем. `syncId` не пишется.
- `name` таблица помечает «Необходимо при создании», а список обязательных полей раздела «Создать»
  его не называет. SDK следует таблице: пустое имя отклоняется до запроса.
- Примеры передают `sum`/`vatSum` (readonly) — SDK их не пишет.
- Пример списка содержит `returnToCommissionerPositions`, которого нет в таблице атрибутов, —
  в модель не входит.
- `payments` (раздел «Связи с другими документами») не помечен readonly и пишется по общему правилу
  привязки документов; `linkedSum` элемента SDK не читает.
- Шаблона `New`, wizard и expand нет; Delete удаляет безвозвратно, Trash — в корзину.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/commissionreportout),
`md/documents/_commissionreportout.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
