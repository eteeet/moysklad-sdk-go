Status: needs-triage
Stage: DRAFT

# 6. Файлы и семейство отчётов Turnover

## Problem Statement

Выделить два разных вида ответов: закрываемый binary stream и типизированные отчётные rows. Не
подгонять их под entity CRUD.

## Solution

Завершённые вертикальные срезы в составе v1; каждый ticket включает models, methods,
тесты и собственное описание полей. Матрица ниже определяет scope этой захватки.
[Общая карта](../sdk-v1/map.md), [нормативные ограничения](../sdk-v1/spec.md),
[каталог JSON-контрактов](../sdk-v1/contracts.md). Текущая работа — планирование.

## User Stories

1. Как разработчик решения, я могу использовать «Transport для binary download» через публичный
   SDK
   и получать фактический ответ API с документированными ограничениями.
2. Как разработчик решения, я могу использовать «Files List/Upload/Download» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
3. Как разработчик решения, я могу использовать «Turnover.All» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
4. Как разработчик решения, я могу использовать «Turnover.ByStore» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
5. Как разработчик решения, я могу использовать «Turnover.ByOperations» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.

## Implementation Decisions и scope tickets

| ID | Ticket | Blocked by |
| --- | --- | --- |
| F01 | [Transport для binary download](issues/01-download-transport.md) | T02, T08, T10 |
| F02 | [Files List/Upload/Download](issues/02-files.md) | T02, T06, T10, F01 |
| F03 | [Turnover.All](issues/03-turnover-all.md) | T02, T06 |
| F04 | [Turnover.ByStore](issues/04-turnover-bystore.md) | F03 |
| F05 | [Turnover.ByOperations](issues/05-turnover-byoperations.md) | F03 |

### F01: Transport для binary download

Отдельный stream/download path рядом с JSON Do: writer-based Download либо явный ReadCloser
ownership; окончательный interface зафиксировать. Auth только trusted origin, signed redirects
разрешаются с удалением credentials; запрет arbitrary URL/SSRF. ResponseInfo/headers и limiter
lifetime учитывают всё чтение.

Проверки и особенности: Large responses без 20-MB upload cap, cancellation/short read/redirect
chain/error body, close на всех путях, semaphore release, cross-origin Authorization stripping;
no auto retry после частично записанного output без explicit caller recovery.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### F02: Files List/Upload/Download

Files как owner-scoped entity/{type}/{id}/files; List metadata, POST upload filename/content
Base64, download по документированным meta.downloadHref. File metadata без вымышленного
entity/file CRUD или общего GetFile endpoint.

Проверки и особенности: ≤10 upload/request, ≤100 files/object server knowledge, encoded-body
bound отдельно от decoded file. Names/media type/no filename path traversal; signed external URL
secret не логируется. Никакого pre-GET для подсчёта текущих100 файлов. Images — отдельная
будущая feature, Delete file — X01.

Источники: [files](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_files.md).

### F03: Turnover.All

GET report/turnover/all: own row type assortment + onPeriodStart/onPeriodEnd/income/outcome,
metric sums/cost/quantities из source. momentFrom/To/groupBy/filter/page, timezone header в
report semantics.

Проверки и особенности: Отсутствующие/nullable quantities не превращаются в invented0 там где
семантика требует nullable; деньги float64, counts не деньги. Meta/context/plain scalar fields,
no top-level entity ID requirement, default groupBy product.

Источники: [report-turnover](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md).

### F04: Turnover.ByStore

GET report/turnover/bystore с собственной формой детализации по складам; filter=product=<href>
либо variant=<href> обязателен, даты optional по методу. Не flatten или вычислять её из All на
клиенте.

Проверки и особенности: Exact path/query, product/variant/reference combinations, nested store
metrics и timezone; не переносить limit/offset/groupBy из All: в таблице ByStore они не
подтверждены.

Источники: [report-turnover](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md).

### F05: Turnover.ByOperations

GET report/turnover/byoperations: own rows operation/assortment и numeric metrics; обязательный
filter product либо variant, по одному непустому значению каждого фильтра. Operation types
открыты; references не требуют готового SDK каждого вида документа.

Проверки и особенности: Таблица называет momentFrom/To required, но соседний текст и пример
допускают defaults: это нерешённое противоречие. Предлагается explicit range в SDK examples, без
заявления о серверной обязательности до freeze. Limit/offset не переносить из All. Read-only
report, DateTime и unknown operation type, не триGET/all с агрегацией; один server report
request.

Источники: [report-turnover](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md).

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

- Files и Reports не требуют переделывать все entity result types в generic CRUD.
- Timezone report header влияет на расчёт периода, а не глобальный DateTime decoder для всех
  запросов.

## Out of Scope

Expand и implicit follow-up; все не перечисленные в tickets ресурсы/операции.
План не объявляет общую CRUD поддержку у ресурсов с другой формой interface.
