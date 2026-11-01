Status: needs-triage
Stage: DRAFT

# 5. Розница и штатные бонусы

## Problem Statement

Закрыть розничные документы и бонусную модель; permission/fiscal ограничения оставляются API,
SDK не имитирует кассу или расчёт бонусов.

## Solution

Завершённые вертикальные срезы в составе v1; каждый ticket включает models, methods,
тесты и собственное описание полей. Матрица ниже определяет scope этой захватки.
[Общая карта](../sdk-v1/map.md), [нормативные ограничения](../sdk-v1/spec.md),
[каталог JSON-контрактов](../sdk-v1/contracts.md). Текущая работа — планирование.

## User Stories

1. Как разработчик решения, я могу использовать «RetailStore Get/List/Create/Update» через
   публичный SDK
   и получать фактический ответ API с документированными ограничениями.
2. Как разработчик решения, я могу использовать «RetailShift Get/List/Create/Update» через
   публичный SDK
   и получать фактический ответ API с документированными ограничениями.
3. Как разработчик решения, я могу использовать «RetailDemand
   Get/List/Create/Update/GetPositions» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
4. Как разработчик решения, я могу использовать «RetailSalesReturn
   Get/List/Create/Update/GetPositions» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
5. Как разработчик решения, я могу использовать «BonusProgram Get/List/Create/Update» через
   публичный SDK
   и получать фактический ответ API с документированными ограничениями.
6. Как разработчик решения, я могу использовать «BonusTransaction Get/List/Create/Update» через
   публичный SDK
   и получать фактический ответ API с документированными ограничениями.

## Implementation Decisions и scope tickets

| ID | Ticket | Blocked by |
| --- | --- | --- |
| E01 | [RetailStore Get/List/Create/Update](issues/01-retailstore.md) | T02, T06 |
| E02 | [RetailShift Get/List/Create/Update](issues/02-retailshift.md) | T02, T06 |
| E03 | [RetailDemand Get/List/Create/Update/GetPositions](issues/03-retaildemand.md) | T02, T06, D01 |
| E04 | [RetailSalesReturn Get/List/Create/Update/GetPositions](issues/04-retailsalesreturn.md) | T02, T06, D01 |
| E05 | [BonusProgram Get/List/Create/Update](issues/05-bonusprogram.md) | T02, T06 |
| E06 | [BonusTransaction Get/List/Create/Update](issues/06-bonustransaction.md) | T02, T06 |

### E01: RetailStore Get/List/Create/Update

Точка продаж: name/organization/store/priceType Create, полный documented core/read-write flags.
Различать metadata кассиров и специальные cashier endpoints; отдельный cashier administration не
включается молча.

Проверки и особенности: allowCreateProducts/controlShippingStock conflict и тариф проверяет
сервер, no hidden config reads; tax/price/payment settings exact JSON, readonly secure fields не
отправляются.

Источники: [retailstore](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md).

### E02: RetailShift Get/List/Create/Update

Смена: required organization/retailStore. Не добавлять GetPositions к смене без соответствующего
API; attached collections остаются typed references/meta. Открытие/закрытие/фискальные side
effects только документированным Update и с явным input.

Проверки и особенности: Активные/закрытые смены, required flags/refs, sparse server data, C/U
readonly rules, APIError и отсутствие автоматически созданной смены.

Источники: [retailshift](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retailshift.md).

### E03: RetailDemand Get/List/Create/Update/GetPositions

Продажа: минимальный Create только retailShift согласно методу; сервер проверяет активность и
moment. Полный documented cash/noCash/QR/prepayment/bonus/fiscal scope зафиксировать отдельной
Write matrix, не требовать все server response поля на Create.

Проверки и особенности: Минимальный Create и full/partial, sums как float64 в копейках, typed
pages, immutable fiscal/read-only fields, server permissions и ошибки, ни автоматической
фискализации, ни бонусного начисления.

Источники: [retaildemand](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retaildemand.md).

### E04: RetailSalesReturn Get/List/Create/Update/GetPositions

Розничный возврат с based/unbased mode; demand обозначает RetailDemand, не обычную Demand.
Name/org/store/retailStore/retailShift/agent и условные cashSum/noCashSum/demand — по методу,
full field scope без предположений о текущем чеке.

Проверки и особенности: Based/unbased minimal bodies, cash/noCash explicit0, active shift/moment
проверки сервером, собственные update restrictions. Не объявлять все поля изменяемыми по
похожести с SalesReturn.

Источники: [retailsalesreturn](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retail_sales_return.md).

### E05: BonusProgram Get/List/Create/Update

Штатная бонусная программа; required
name/active/allProducts/allAgents/earnRateRoublesToPoint/spendRatePointsToRouble/maxPaidRatePercents.
Диапазоны/зависимости и активная программа по source, не рассчитывать баланс/правила маркетинга
в SDK.

Проверки и особенности: Явные false и 0 обязательных полей не пропадают,
arrays/refs/null/readonly, rate restrictions — сервер. Не считать баллы копейками только из-за
слова bonus.

Источники: [bonusprogram](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bonus_program.md).

### E06: BonusTransaction Get/List/Create/Update

agent/bonusProgram/transactionType Create; bonusValue, status, parentDocument и ограничения
доступной записи по source. Баллы — собственная числовая величина, не деньги. Permission model
Employee token vs solution token описать отдельно; JSON API наличие POST не гарантирует
разрешение Vendor token.

Проверки и особенности: WAIT_PROCESSING/COMPLETED/unknown, accrual/debit specifics,
immutable/C-only fields; синтетический отказ токена решения передаётся без скрытого fallback на
personal token. Live в Draft — отдельное согласие и не доказано unit tests.

Источники: [bonustransaction](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bonus_operation.md).

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

- Референсы не требуют предварительного Get связанных сущностей; порядок E01→E02→E03/E04 полезен
  для интеграционного сценария.
- Розничные документы и бонусные операции не получают автоматических write retries.

## Out of Scope

Expand и implicit follow-up; все не перечисленные в tickets ресурсы/операции.
План не объявляет общую CRUD поддержку у ресурсов с другой формой interface.
