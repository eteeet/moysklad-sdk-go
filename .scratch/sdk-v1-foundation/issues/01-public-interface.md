Status: resolved
Stage: APPROVED
Type: task
Blocked by: None

# T01: Утвердить общий interface options/response metadata

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Зафиксировать RequestOption/RequestOptions, client defaults и получение ResponseInfo; прототипы
вызовов Get/List/Create/Update/GetPositions/отчёта/download. Рекомендация: сохранить
существующие результаты и вызовы, добавить optional per-call options и WithResponseInfo(&info).
Никакого global LastResponse. Альтернативу Result[T] сравнить явно; это решение, не уже принятый
interface.

## Специальные проверки

Компилируемые consumer examples; показать совместимость обычных вызовов и отдельно изменения
method-value types при variadic options. Параллельные запросы не делят mutable response pointer.
Одна shared policy instance на transport, несколько фасадов могут использовать её.

## Acceptance criteria

- [ ] Scope/публичный interface и field/method matrix утверждены; зависимости resolved
  и их код присутствует в рабочем дереве. Readonly, required response/Create, null,
  immutable и выбранные TODO разделены по источникам; противоречия не замолчаны.
- [ ] Реализован весь What to build, без неоговорённых методов и hidden requests.
- [ ] Специальные проверки этого ticket реализованы через публичный interface.
- [ ] JSON обоими codecs и независимые golden; upstream/provenance не изменены.
  Для transport/iterator/CI tickets вместо entity JSON применяются их собственные
  публичные seams и контрактные tests, а не фиктивные модели.
- [ ] Реальные transport/httptest tests покрывают успех и ошибки; context и
  metadata сохраняются, caller-controlled retry/walker не подменяет один вызов.
- [ ] go test -race, go vet ./... ./testdata, gofmt и diff-check PASS; независимые
  Standards/Spec review, исправления и финальный BASE/результат записаны в Answer.
- [ ] README/guide/status отражают реальный scope, ограничения и неизвестные факты.
  Planning без команды не коммитится; один code commit по ticket.

## Источники

[транспортный research](../../sdk-v1/research/transport.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `7442803`, code commit `ba5d2bc`. Решение пользователя: per-call `opts ...RequestOption`
+ `WithResponseInfo(&info)`, без `Result[T]`. ADR: [0001](../../../docs/adr/0001-request-options.md).

- transport: `RequestOption`, `RequestOptions`, `ResponseInfo`, `WithHeader`, `WithResponseInfo`,
  `Config.Header` (client defaults), `Client.Send` — seam для limiter/retry.
  Защищённые headers и дубли key по регистру отклоняются до отправки.
- Прототип — UOM Get/Create/Update; остальные 16 services и positions переводит T02.
  List/GetPositions/Report/Download зафиксированы формой в ADR, компилируются в своих tickets.
- Consumer examples: `entity/uom/example_test.go`, включая слом method value.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review Standards: стандарты не нарушены; по замечаниям unexported `newRequestOptions`,
  `NewRequest` снова без opts. Review Spec: исправлены ADR (Download не через Send,
  shared policy в Client, attempts в T09), добавлены тесты shared Client и дублей header.

## Comments

Ревью 46e9d06 (Standards P2 / Spec P2): локальная ошибка оставляла ResponseInfo прошлого вызова. Исправлено в `bc2dad4`: обнуление в `WithResponseInfo` при создании опции; ADR 0001 дополнен.
