Status: needs-triage
Stage: DRAFT

# 1. Транспорт, query, лимиты и List

## Problem Statement

Один transport управляет request options, response metadata, бюджетом и отменой. Сущности
получают List и доступ к headers без дублирования этих механизмов.

## Solution

Завершённые вертикальные срезы в составе v1; каждый ticket включает models, methods,
тесты и собственное описание полей. Матрица ниже определяет scope этой захватки.
[Общая карта](../sdk-v1/map.md), [нормативные ограничения](../sdk-v1/spec.md),
[каталог JSON-контрактов](../sdk-v1/contracts.md). Текущая работа — планирование.

## User Stories

1. Как разработчик решения, я могу использовать «Утвердить общий interface options/response
   metadata» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
2. Как разработчик решения, я могу использовать «Все request/response headers и интеграция
   options» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
3. Как разработчик решения, я могу использовать «Каталог кодов ошибок и status-aware helpers»
   через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
4. Как разработчик решения, я могу использовать «POST security/token» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
5. Как разработчик решения, я могу использовать «Filter grammar и order» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
6. Как разработчик решения, я могу использовать «Query options и получение одной страницы» через
   публичный SDK
   и получать фактический ответ API с документированными ограничениями.
7. Как разработчик решения, я могу использовать «Явный walker/iterator страниц» через публичный
   SDK
   и получать фактический ответ API с документированными ограничениями.
8. Как разработчик решения, я могу использовать «Weighted rate limiter и concurrency» через
   публичный SDK
   и получать фактический ответ API с документированными ограничениями.
9. Как разработчик решения, я могу использовать «Ограниченные retries/backoff» через публичный
   SDK
   и получать фактический ответ API с документированными ограничениями.
10. Как разработчик решения, я могу использовать «Размеры запросов и endpoint limits» через
    публичный SDK
   и получать фактический ответ API с документированными ограничениями.
11. Как разработчик решения, я могу использовать «List шести существующих справочников» через
    публичный SDK
   и получать фактический ответ API с документированными ограничениями.
12. Как разработчик решения, я могу использовать «List одиннадцати существующих документов»
    через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
13. Как разработчик решения, я могу использовать «SyncID Product и единый контракт syncId» через
    публичный SDK
   и получать фактический ответ API с документированными ограничениями.

## Implementation Decisions и scope tickets

| ID | Ticket | Blocked by |
| --- | --- | --- |
| T01 | [Утвердить общий interface options/response metadata](issues/01-public-interface.md) | нет |
| T02 | [Все request/response headers и интеграция options](issues/02-headers.md) | T01 |
| T03 | [Каталог кодов ошибок и status-aware helpers](issues/03-error-codes.md) | нет |
| T04 | [POST security/token](issues/04-security-token.md) | T01, T03 |
| T05 | [Filter grammar и order](issues/05-filter-order.md) | T01 |
| T06 | [Query options и получение одной страницы](issues/06-query-pages.md) | T01, T05 |
| T07 | [Явный walker/iterator страниц](issues/07-iterator.md) | T06 |
| T08 | [Weighted rate limiter и concurrency](issues/08-rate-concurrency.md) | T02 |
| T09 | [Ограниченные retries/backoff](issues/09-safe-retry.md) | T03, T08 |
| T10 | [Размеры запросов и endpoint limits](issues/10-size-limits.md) | T01, T03 |
| T11 | [List шести существующих справочников](issues/11-list-dictionaries.md) | T02, T06 |
| T12 | [List одиннадцати существующих документов](issues/12-list-documents.md) | T02, T06 |
| T13 | [SyncID Product и единый контракт syncId](issues/13-product-syncid.md) | нет |

### T01: Утвердить общий interface options/response metadata

Зафиксировать RequestOption/RequestOptions, client defaults и получение ResponseInfo; прототипы
вызовов Get/List/Create/Update/GetPositions/отчёта/download. Рекомендация: сохранить
существующие результаты и вызовы, добавить optional per-call options и WithResponseInfo(&info).
Никакого global LastResponse. Альтернативу Result[T] сравнить явно; это решение, не уже принятый
interface.

Проверки и особенности: Компилируемые consumer examples; показать совместимость обычных вызовов
и отдельно изменения method-value types при variadic options. Параллельные запросы не делят
mutable response pointer. Одна shared policy instance на transport, несколько фасадов могут
использовать её.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### T02: Все request/response headers и интеграция options

Весь инвентарь headers из research/transport.md: typed request options, repeated values, clone
raw response headers, parsed optional values с сохранением malformed/raw. Протянуть per-call
options и metadata во все 17 существующих services и positions. Управление auth остаётся явным;
новая мета не теряется на APIError.

Проверки и особенности: Repeated WebHook-DisableByPrefix без CSV; beta feature одно значение;
удалённый WebHook-Disable не рекламируется как рабочий. Region/Accept-Language, report date with
numeric timezone, content/deprecation/auth/download headers. Defaults/overrides,
malformed/unknown/negative/overflow значения, cancellation и race. DateTime не перенастраивается
глобально по report header.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### T03: Каталог кодов ошибок и status-aware helpers

Полный числовой каталог официального _errors.md (545 уникальных кодов в закреплённом revision):
именованные константы с provenance, стабильный int-compatible Code, распознавание status/code
без сравнения локализованных строк. Неизвестные коды сохраняются. Сохранить errors.As/Is и
единственный префикс.

Проверки и особенности: Все code в источнике покрыты inventory test; aliases/дубли описаны. 413
возвращает APIError/детали, не делает autosplit или retry. HTML/empty/gzip errors и raw headers
доступны. Никаких новых retry решений по тексту ошибки.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### T04: POST security/token

Отдельный security service: CreateToken(ctx, options), POST security/token с Basic Auth, без
тела; HTTP 200 и непустой access_token. В документации явно: выпуск отзывает предыдущие токены
сотрудника. Не вызывать из NewClient и не переключать Auth автоматически.

Проверки и особенности: Exact method/path/headers/empty body; malformed/empty token,
auth/API/context errors; одна отправка даже при 429/5xx/потере ответа. Токен не попадает в
diagnostics/fixtures/logs. Live выпуск только по отдельной команде.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### T05: Filter grammar и order

Builders для documented filter operators, escape literal semicolon и одного URL-encoding;
Reference href, UUID, Attribute href, null/empty, LIKE, archived/isDeleted. Order по нескольким
полям с asc/desc. Open field names и raw escape hatch; endpoint capabilities не заменяются
универсальным enum.

Проверки и особенности: Повторы = дают OR, != дают AND, разные поля AND; unsupported комбинации
не переосмыслять как произвольный AST. Unicode, %, &, +, ;, backslash, nested href;
report-specific grammar отдельно. Ни hidden lookup для фильтров, ни double encoding.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### T06: Query options и получение одной страницы

Limit/offset, search, namedfilter, fields (одно значение), report dates/groupBy/interval и
полный инвентарь прочих query names из research. Typed options только на поддерживающих
endpoints; общий url.Values сохраняет новые параметры. Обычные Collection[T], plain arrays и
report-specific ответы не смешиваются. Expand не добавляется.

Проверки и особенности: Endpoint-specific defaults: entity lists обычно 1000, audit 25/max100;
GetPositions existing Limit=0→1000 сохраняется. fields restrictions scoped, namedfilter href,
exact query encoding. Для неподдержанных ресурсов raw passthrough не объявляется готовым typed
service. async не разрешается общей Get options без специальной операции.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### T07: Явный walker/iterator страниц

Явный iterator по List/GetPositions/report pages, cancellation/early stop, bounded память;
обычные методы по-прежнему делают один запрос. Проверять nextHref same-origin и safe path либо
вычислять offsets по документированному контракту, не следовать слепо внешним ссылкам.

Проверки и особенности: Empty page, no total/size where optional, wrong/missing rows, repeated
nextHref, no-progress, unusual offsets, cancellation до следующей страницы. Не обещать
snapshot-consistency или exactly-once при меняющемся наборе. Array-only PriceType не
пагинировать вымышленным limit.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### T08: Weighted rate limiter и concurrency

Раздельные shared principal/account budgets: 45 weighted units/3s; ≤5 principal и ≤20 account
concurrent. Auth profile задаётся явно (Employee/App/Draft), clock/date и weights
конфигурируемы. Текущий Employee weight3, с 2026-12-01 weight4, Draft4, App1; только два точных
stock endpoints weight5.

Проверки и особенности: Fake clock, rolling-window границы, fairness/queue cancellation; acquire
каждой попытки, release после body close, не во время backoff. Один процесс не гарантирует
global account limit других процессов; sharing через общий coordinator описан. Headers rate
budgets не интерпретировать недоказанными weighted units.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### T09: Ограниченные retries/backoff

Operation-aware retry policy: только подтверждённые reads GET/HEAD по умолчанию;
POST/PUT/security-token/async-submit — без автоматического повтора. 429/временные 5xx/transport
failures с ограничением attempts, context, jitter и корректными server hints. Окончательные
metadata и APIError сохраняются.

Проверки и особенности: X-Lognex durations в миллисекундах, стандартный Retry-After —
seconds/date; не обрезать server minimum своим cap. Каждая попытка новый replay body/headers и
limiter admission; close промежуточных body до ожидания. Nil/unreplayable stream, decode error
после 2xx и canceled context не повторяются.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### T10: Размеры запросов и endpoint limits

Preflight encoded outbound body ≤20 MB и request metadata ≤8 KB (URL/auth/headers), depth10 и
ограничение массивов в применимых операциях. Не вводить cap успешного download/report response.
Отдельно перечислить client caps при неуточнённой API единице MB/KB и способ override.

Проверки и особенности: Base64 expansion, UTF-8 bytes, repeated headers, boundary max/max+1,
streaming unknown length и cancellation. API limits batch1000/5min и async queue4 остаются gates
X02/X03, если эти операции не включены; не заявлять реализацию отсутствующих методов. Не
применять 1000 к произвольному вложенному массиву без API основания.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### T11: List шести существующих справочников

Добавить typed List + общие options в Product/UOM/Country/ExpenseItem/Counterparty/Organization.
Покрыть supported filter/order/search/namedfilter/fields по endpoint, полный Collection
context/meta/rows. Текущие Get/Create/Update сохраняют прежний вызов.

Проверки и особенности: Табличные exact HTTP tests на каждый package, plain/gzip, empty array,
invalid/null envelope, query encoding, server error/metadata и один запрос. Не добавлять
отсутствующие методы по шаблону. Existing golden и read models не регрессируют.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### T12: List одиннадцати существующих документов

Typed List и request options для PaymentIn/Out, CashIn/Out, InvoiceIn/Out,
Supply/Demand/PurchaseOrder/Move/PurchaseReturn. Фильтры archived/isDeleted там, где поддержаны.
Positions в List/Get — только фактический ответ, без expand/GetPositions follow-up.

Проверки и особенности: Все 11 exact resources, collection/meta/context, query/headers/errors и
no hidden requests. Незаданные query/headers прежних вызовов не получают неоговорённых defaults.
Старые Get/Create/Update tests сохранены.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### T13: SyncID Product и единый контракт syncId

Добавить SyncID в Product read/Create/Update с Field[string] в requests. Проверить единообразие
уже реализованных Counterparty/Organization/11 документов и будущих entities только там, где
поле есть в API. Server write-once и resource-specific create dedup, не общий Idempotency-Key.

Проверки и особенности: Оба codecs, unset/value/empty/null, SyncID-only Update, exact body,
APIError, ни pre-GET, ни auto retry. Не добавлять syncId к ресурсам без этого поля.

Источники: [транспортный research](../sdk-v1/research/transport.md).

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

- Общий request/response interface утверждается T01 до массового retrofit.
- Default policy: явная конфигурация limiter/retry, без silent Auth-profile inference.
- Delete/batch/async — открытое решение X01–X03, не автоматическое следствие фразы «все лимиты».

## Out of Scope

Expand и implicit follow-up; все не перечисленные в tickets ресурсы/операции.
План не объявляет общую CRUD поддержку у ресурсов с другой формой interface.
