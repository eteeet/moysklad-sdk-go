Status: needs-triage
Stage: DRAFT

# 4. Заказ покупателя, складские операции и возврат

## Problem Statement

Продолжить вертикальные срезы документов поверх общих типов и List, сохранив одну универсальную
форму позиций там, где JSON не конфликтует.

## Solution

Завершённые вертикальные срезы в составе v1; каждый ticket включает models, methods,
тесты и собственное описание полей. Матрица ниже определяет scope этой захватки.
[Общая карта](../sdk-v1/map.md), [нормативные ограничения](../sdk-v1/spec.md),
[каталог JSON-контрактов](../sdk-v1/contracts.md). Текущая работа — планирование.

## User Stories

1. Как разработчик решения, я могу использовать «Расширения общей формы позиций» через публичный
   SDK
   и получать фактический ответ API с документированными ограничениями.
2. Как разработчик решения, я могу использовать «CustomerOrder
   Get/List/Create/Update/GetPositions» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
3. Как разработчик решения, я могу использовать «Enter Get/List/Create/Update/GetPositions»
   через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
4. Как разработчик решения, я могу использовать «Loss Get/List/Create/Update/GetPositions» через
   публичный SDK
   и получать фактический ответ API с документированными ограничениями.
5. Как разработчик решения, я могу использовать «SalesReturn
   Get/List/Create/Update/GetPositions» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.

## Implementation Decisions и scope tickets

| ID | Ticket | Blocked by |
| --- | --- | --- |
| D01 | [Расширения общей формы позиций](issues/01-position-contracts.md) | T06 |
| D02 | [CustomerOrder Get/List/Create/Update/GetPositions](issues/02-customerorder.md) | T02, T06, D01 |
| D03 | [Enter Get/List/Create/Update/GetPositions](issues/03-enter.md) | T02, T06, D01 |
| D04 | [Loss Get/List/Create/Update/GetPositions](issues/04-loss.md) | T02, T06, D01 |
| D05 | [SalesReturn Get/List/Create/Update/GetPositions](issues/05-salesreturn.md) | T02, T06, D01 |

### D01: Расширения общей формы позиций

Зафиксировать union новых документированных полей rows
customerorder/enter/loss/salesreturn/retaildemand/retailsalesreturn: типы, money vs counts,
writable markers, Reference/pack/slot/stock enrichment для нужных fields options. Реальный
type/meaning conflict требует отдельной формы, простое отсутствие поля — нет.

Проверки и особенности: Write matrix из contracts.md: reserve/shipped/cost/inTransit и прочие не
добавляются по аналогии. Явный readonly важнее примера, неизвестные refs/server constraints не
валидируются локальным GET. Оба codecs и прежние семь aliases не регрессируют.

Источники: [customerorder](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_customerOrder.md), [enter](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_enter.md), [loss](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_loss.md), [salesreturn](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_sales_return.md), [retaildemand](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retaildemand.md), [retailsalesreturn](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retail_sales_return.md).

### D02: CustomerOrder Get/List/Create/Update/GetPositions

Полный CustomerOrder header/positions scope, required organization/agent, reserve/shipped и
shipment/payment links по source. Header arrays, attrs, prices и positions не имеют одной общей
merge semantics.

Проверки и особенности: Minimal/full/golden, readonly fields, typed fixed reader,
list/search/filter, exact 1000/1001 и accounts/null/positions replacement, server-responses
вместо input.

Источники: [customerorder](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_customerOrder.md).

### D03: Enter Get/List/Create/Update/GetPositions

Оприходование с organization/store, без выдуманного agent; собственный read/write header и поля
строк, shared types там где совпадают.

Проверки и особенности: Минимальный Create из двух required refs, fields и read-only sums, nil
dates, stock/positions queries только где документированы, одна отправка.

Источники: [enter](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_enter.md).

### D04: Loss Get/List/Create/Update/GetPositions

Списание с organization/store, documented rows и links. Template /new и wizard не добавляются
автоматически.

Проверки и особенности: Required refs/readonly/query, embedded max1000, price/cost/quantity без
локальной арифметики, unknown refs и APIErrors.

Источники: [loss](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_loss.md).

### D05: SalesReturn Get/List/Create/Update/GetPositions

Возврат покупателя по Demand и без основания: organization/store/agent обязательны, demand
required только для based. Update запреты сторон/основания/цен по source, без копирования
PurchaseReturn запретов по аналогии.

Проверки и особенности: Два режима Create, full/partial Update, position matching ID/quantity
constraints оставлены серверу, типы row/read-only и source contradictions фиксируются. Existing
Demand service не нужен для отправки Reference.

Источники: [salesreturn](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_sales_return.md).

## Общий контракт реализации

- Перед кодом — field/method freeze: поимённая Write matrix, readonly/TODO,
  null/required/Create-only/Update-only; status needs-triage сам не разрешает реализацию.
- JSON: encoding/json v1 и v2, Field+omitzero для optional request, независимые
  golden, неизменённые upstream с SHA-256 provenance; raw JSON long не проходит float64.
- Dates — существующий DateTime; деньги float64 в копейках, counts/проценты/баллы/курс
  не объявляются деньгами. Reference без expand, unknown enum/type не отвергается.
- HTTP: настоящий transport+httptest, exact method/path/query/headers/body,
  plain/gzip, фактический server response, ошибки/auth/context/cancellation,
  errors.Is/As и единый moysklad:. Обычный вызов — одна отправка, если caller
  явно не включил retry/walker; каждый дополнительный запрос доказуем тестами.
- Бизнес-ограничения и permissions оставлены API; SDK не выполняет pre-GET ради
  local validation. Exception — только явный пользовательский multi-call operation.
- Проверки: go test -race -count=1 ./..., go vet ./... ./testdata, gofmt,
  git diff --check; независимые Standards/Spec review на фиксированном BASE.
- Коммит по ticket; shared dependency раньше зависимого, публичные docs отдельным
  завершающим commit волны. Planning не коммитить без отдельной команды.
- Live API, credentials, push/tag/release, worktree cleanup и production-код
  не входят в текущую задачу подготовки spec/tickets.

## Open Decisions

- Полный field freeze header/positions каждого ticket до реализации; специальные запреты Update
  выражены в request типах.
- Отдельные position mutations, templates/new/wizard не входят по одной лишь поддержке
  GetPositions.

## Out of Scope

Expand и implicit follow-up; все не перечисленные в tickets ресурсы/операции.
План не объявляет общую CRUD поддержку у ресурсов с другой формой interface.
