# Payroll

`entity/payroll`: модель `Payroll`, `CreateRequest`, `UpdateRequest`, `NewService(client)`; методы
`Get/List/Create/Update/Delete`, `CreateMany/UpdateMany/DeleteMany` и `Trash(ctx, id)`. Позиции —
собственная форма `Position`/`PositionRequest` (meta.type `payrollposition`): сотрудник, начисление
по окладу и сдельная оплата.
[Общие правила методов, JSON и ошибок](documents.md), [корзина](documents-extra.md).

Минимальный Create: `organization`, `startPayrollPeriod`, `endPayrollPeriod` — проверяются до
запроса (непустой meta.href, ненулевая дата). Кратность периода 24 часам проверяет сервер.

```go
doc, err := s.Create(ctx, payroll.CreateRequest{
	Organization: organization, StartPayrollPeriod: start, EndPayrollPeriod: end,
	Positions: model.Value([]payroll.PositionRequest{{
		Employee: model.Value(employee), BaseSalary: model.Value(202000.0), PieceworkSalary: model.Value(1000.0),
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
| `externalCode` | `string` | Field | Field |
| `syncId` | `string` | Field | Field |
| `created` | `model.DateTime` | — | — |
| `updated` | `model.DateTime` | — | — |
| `moment` | `model.DateTime` | Field | Field |
| `startPayrollPeriod` | `model.DateTime` | required | Field |
| `endPayrollPeriod` | `model.DateTime` | required | Field |
| `applicable` | `bool` | Field | Field |
| `shared` | `bool` | Field | Field |
| `printed` | `bool` | — | — |
| `published` | `bool` | — | — |
| `organization` | `*model.Reference` | required | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `attributes` | `[]model.Attribute` | Field | Field |
| `sum` | `float64` | — | — |
| `positions` | `*model.Collection[Position]` | Field | Field |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `payments` | `[]model.Reference` | — | — |

`sum` — копейки; `baseSalary`/`pieceworkSalary` позиций — `float64` в тех же единицах (в примерах
источника `sum` равна их сумме по позициям).

## Позиции

- `GetPositions(ctx, id, positions.PageOptions) (*model.Collection[Position], error)` — одна страница
  `…/{id}/positions`, только limit/offset (`fields` отклоняется до запроса).
- `GetPosition`, `CreatePositions` (один POST массивом, ответ `[]batch.Result[Position]`),
  `UpdatePosition` (пустой request — ошибка до запроса), `DeletePosition`, `DeleteManyPositions`
  (POST `…/positions/delete`, 1–1000 id; ответ без тела — только error).
- Поля позиции: `meta`, `id`, `accountId` (только чтение), `employee`, `baseSalary`, `pieceworkSalary`.
  Для новой позиции источник требует все три поля (положительные суммы, округляются до целого) —
  проверяет сервер.
- Positions в Create/Update — не больше 1000 позиций (1001 отклоняется до запроса); в Update массив
  заменяет весь набор: позиция с `ID` — существующая (id только адресует её), без — новая.

## Ограничения и противоречия

- Ticket описывает позиции как «employee/sum/type», источник — `employee`, `baseSalary`,
  `pieceworkSalary`. SDK следует таблице позиций источника.
- `payments` — связь с исходящими платежами и расходными ордерами (`paymentout`, `cashout`);
  разрешённые типы помечены «Только для чтения», поэтому SDK `payments` только читает
  (`linkedSum` элемента не читает).
- У payroll нет `code`, `project`, `rate`, `deleted`: их нет в таблице атрибутов.
- `syncId` в таблице не помечен readonly и пишется в Create/Update; после заполнения его менять
  запрещает API.
- Шаблона `New`, wizard и expand нет; Delete удаляет безвозвратно, Trash — в корзину.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/payroll),
`md/documents/_payroll.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
