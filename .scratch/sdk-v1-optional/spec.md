Status: resolved
Stage: APPROVED

# 8. Опциональное расширение: Delete, batch и async

## Problem Statement

Закрыть неоднозначность «все лимиты» без скрытого добавления операций, не перечисленных явно в
v1. Эти tickets требуют решения пользователя.

## Solution

Завершённые вертикальные срезы в составе v1; каждый ticket включает models, methods,
тесты и собственное описание полей. Матрица ниже определяет scope этой захватки.
[Общая карта](../sdk-v1/map.md), [нормативные ограничения](../sdk-v1/spec.md),
[каталог JSON-контрактов](../sdk-v1/contracts.md). Текущая работа — планирование.

## User Stories

1. Как разработчик решения, я могу использовать «Delete выбранных ресурсов» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
2. Как разработчик решения, я могу использовать «Batch create/update/delete» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
3. Как разработчик решения, я могу использовать «Async submit/status/result/cancel» через
   публичный SDK
   и получать фактический ответ API с документированными ограничениями.

## Implementation Decisions и scope tickets

| ID | Ticket | Blocked by |
| --- | --- | --- |
| X01 | [Delete выбранных ресурсов](issues/01-delete.md) | T02, T03 |
| X02 | [Batch create/update/delete](issues/02-batch.md) | T02, T03, T08, T10 |
| X03 | [Async submit/status/result/cancel](issues/03-async.md) | T02, T03, T08, T09, F01 |

### X01: Delete выбранных ресурсов

По-resource capability matrix single Delete для поддержанных entity/metadata/files.
Метод/маршрут/status индивидуальны. Не смешивать delete, recycle bin и purge. No write retries
по умолчанию.

Проверки и особенности: 204/200 по endpoint, APIError body, absent JSON, UUID, permissions; не
генерировать Delete для assortment/context/audit/report и прочих read-only endpoints.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### X02: Batch create/update/delete

Resource-specific arrays max1000, mixed per-item results/errors и timeout5min. Explicit batch
methods, без autosplit и без иллюзии транзакционного all-or-nothing. Удаление POST .../delete
отдельно от single DELETE.

Проверки и особенности: 1000/1001, canceled/timeout ambiguous outcome, response item errors и
metadata, не повторять неизвестный результат; batch status не равен успеху всех элементов.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### X03: Async submit/status/result/cancel

Только documented async-enabled endpoints; Submit даже для GET является side effect. Own task
model, statuses и poll/result/cancel endpoint. ≤4 queued account tasks policy отдельно от HTTP
semaphore. async=true нельзя сочетать limit/offset.

Проверки и особенности: Submit не retry; poll bounded/context; result download origin/auth rules
F01 при ссылочном result; 202/Location/response shape из _async.md, cancel нельзя выдавать за
rollback. Distributed account queue полностью одним Client не контролируется.

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

- 2026-11-02, пользователь: X01–X03 входят в v1, объём — всё, что источник описывает для сервисов
  SDK; batch — раздельные CreateMany/UpdateMany/DeleteMany. Итоги — Answer tickets,
  [docs/write-async.md](../../docs/write-async.md).

## Out of Scope

Expand и implicit follow-up; все не перечисленные в tickets ресурсы/операции.
План не объявляет общую CRUD поддержку у ресурсов с другой формой interface.
