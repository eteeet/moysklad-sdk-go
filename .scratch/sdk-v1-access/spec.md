Status: resolved
Stage: APPROVED

# 2. Сотрудники, права, metadata и Audit

## Problem Statement

Получить контекст аккаунта, права и изменения через отдельные typed ресурсы; Audit доступен до
завершения розничных документов.

## Solution

Завершённые вертикальные срезы в составе v1; каждый ticket включает models, methods,
тесты и собственное описание полей. Матрица ниже определяет scope этой захватки.
[Общая карта](../sdk-v1/map.md), [нормативные ограничения](../sdk-v1/spec.md),
[каталог JSON-контрактов](../sdk-v1/contracts.md). Текущая работа — планирование.

## User Stories

1. Как разработчик решения, я могу использовать «Employee: Get/List/Create/Update» через
   публичный SDK
   и получать фактический ответ API с документированными ограничениями.
2. Как разработчик решения, я могу использовать «Role и права Employee» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
3. Как разработчик решения, я могу использовать «GET context/employee» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
4. Как разработчик решения, я могу использовать «GET context/application» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
5. Как разработчик решения, я могу использовать «Общий каталог и metadata сущностей» через
   публичный SDK
   и получать фактический ответ API с документированными ограничениями.
6. Как разработчик решения, я могу использовать «States: чтение, Create/Update» через публичный
   SDK
   и получать фактический ответ API с документированными ограничениями.
7. Как разработчик решения, я могу использовать «Audit contexts и filters» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
8. Как разработчик решения, я могу использовать «Audit events и lossless diff» через публичный
   SDK
   и получать фактический ответ API с документированными ограничениями.

## Implementation Decisions и scope tickets

| ID | Ticket | Blocked by |
| --- | --- | --- |
| A01 | [Employee: Get/List/Create/Update](issues/01-employee.md) | T02, T06 |
| A02 | [Role и права Employee](issues/02-roles-security.md) | A01 |
| A03 | [GET context/employee](issues/03-employee-context.md) | T02 |
| A04 | [GET context/application](issues/04-application-context.md) | T02 |
| A05 | [Общий каталог и metadata сущностей](issues/05-entity-metadata.md) | T02, T06 |
| A06 | [States: чтение, Create/Update](issues/06-states.md) | A05 |
| A07 | [Audit contexts и filters](issues/07-audit-contexts.md) | T02, T06 |
| A08 | [Audit events и lossless diff](issues/08-audit-events.md) | A07 |

### A01: Employee: Get/List/Create/Update

Полный основной Employee contract; required lastName; отдельные request types, references и
typed nested salary/cashier where documented. Фото и credential lifecycle не добавляются этим
тикетом.

Проверки и особенности: Поля по contracts.md, sparse/null, false/zero, exact required Create,
list/query/headers и errors; Employee rights не угадываются по имени роли.

Источники: [employee](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md).

### A02: Role и права Employee

Role Get/List/Create/Update и Employee GetSecurity/UpdateSecurity, чтение documented role
presets. Permissions — неоднородный JSON object по примерам, несмотря на Array(Object) в таблице
Role; known typed values + lossless unknown jsontext.Value, без map[string]any/float64 для
сырого JSON.

Проверки и особенности: Приоритет admin/App permissions проверяет сервер; не активировать
Employee, не сбрасывать пароль и не присваивать admin/owner скрытыми вызовами. Bool/scope/object
task permissions, неизвестные ключи/enum, exact security routes, отсутствие потерянных values.

Источники: [role](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_custom_role.md), [employee](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md).

### A03: GET context/employee

Отдельный context/employee.Get(ctx) без ID; собственная форма контекста и permissions.
Переиспользовать Employee/Permission types только при совпадении JSON contract, не путать с
ResponseContext в коллекциях.

Проверки и особенности: Basic/Employee/App responses/permission errors, sparse fields и
неизвестные permissions, нет обязательного caller id и дополнительных запросов.

Источники: [employee](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md).

### A04: GET context/application

Отдельный context/application.Get(ctx) с App token, meta/id/appUid; не выдавать за Vendor API,
не реализовывать установку/активацию решения. entity/application/{installationId} — отдельная
будущая возможность, не alias контекста.

Проверки и особенности: Один GET без path id/query/body; App auth и APIError, отсутствующие поля
не заполняются. NewClient не делает этот запрос.

Источники: [general](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md).

### A05: Общий каталог и metadata сущностей

GET entity/metadata с filter=type=..., GET entity/{type}/metadata, чтение definitions attributes
и отдельного attribute metadata, states/reference lists по фактическому endpoint. Open entity
types safe segment; named supported routes, без expand.

Проверки и особенности: Global response не обязательно Collection; формы
attributes/meta/rows/plain arrays сверить по источнику. Не навязывать ID. Разные metadata
поддерживаемых типов, unknown definitions/types и auth errors. Создание/редактирование
definitions — отдельно от чтения и не входит скрыто.

Источники: [general](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md).

### A06: States: чтение, Create/Update

Динамический entity/{type}/metadata/states resource, чтение через metadata и отдельный GET
только при наличии documented route. Create/Update, required name/color/stateType, ARGB color
integer, открытый stateType; schema не обычный entity/state.

Проверки и особенности: Admin requirements передаются серверу; точные пути,
Regular/Successful/Unsuccessful/unknown, zero color/Field states. Не придумывать states.List
endpoint, если список находится в metadata; Delete зависит от X01.

Источники: [states](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_states.md).

### A07: Audit contexts и filters

Audit.ListContexts и GetFilters по official routes; фильтры
moment/uid/source/entityType/eventType/application по documented scope; typed AuditContext и
FilterCatalog, metadata страницы. Default limit25, max100 именно contexts.

Проверки и особенности: Не обещать endless changelog/cursor/retention/exactly-once; timestamp
tie/duplicates возможны. Отсутствие прав/неизвестный source/entityType, контекст Employee vs
App, exact filters, одна страница без чтения events.

Источники: [audit](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md), [general](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md).

### A08: Audit events и lossless diff

EventsByContext и EventsByEntity с безопасным type/id path; типизированные event headers и
raw/lossless diff с typed accessors where justified. Неизвестные event/type/diff fields
сохраняются; не декодировать каждое старое значение в текущую entity model.

Проверки и особенности: Большие integer long, nested/null/removed/unknown fields, event-specific
registration/publication/email diff, отсутствие подтверждённых event pagination guarantees явно
отражено. Worker синхронизации, durable checkpoint и реплика не входят в SDK.

Источники: [audit](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md).

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

- Metadata в v1 читается; States отдельно имеют Create/Update. Definition CRUD не маскируется
  словом metadata.
- Противоречие Role.permissions table vs example требует записанного JSON-контракта в A02 перед
  кодом.
- Audit не зависит от готовых services всех entityType: references и diff остаются открытыми.

## Out of Scope

Expand и implicit follow-up; все не перечисленные в tickets ресурсы/операции.
План не объявляет общую CRUD поддержку у ресурсов с другой формой interface.
