Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md

# T08: Weighted rate limiter и concurrency

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Раздельные shared principal/account budgets: 45 weighted units/3s; ≤5 principal и ≤20 account
concurrent. Auth profile задаётся явно (Employee/App/Draft), clock/date и weights
конфигурируемы. Текущий Employee weight3, с 2026-12-01 weight4, Draft4, App1; только два точных
stock endpoints weight5.

## Специальные проверки

Fake clock, rolling-window границы, fairness/queue cancellation; acquire каждой попытки, release
после body close, не во время backoff. Один процесс не гарантирует global account limit других
процессов; sharing через общий coordinator описан. Headers rate budgets не интерпретировать
недоказанными weighted units.

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

BASE `de7248f`, code commit `7fb0c21`.

- `transport.Limiter` (`NewLimiter(LimiterConfig)`) подключается явно через `Config.Limiter`; nil —
  без ограничений SDK. Профиль обязателен (App/DraftApp/Employee/Unknown=4 ед.), часы `Clock` и
  `Weight` подменяются. `DefaultWeight` — расписание _restrictions.md; часовой пояс дат источник не
  задаёт — SDK берёт полночь по Москве; вес 5 только у GET report/stock/all и /bystore (путь относительно
  BaseURL, включая /api/remap/1.2).
- Скользящее окно 45 ед./3 с на principal, очередь ожидания бюджета по каналу; ≤5 principal и
  `Account` ≤20 (делится между Limiter одного процесса; между процессами нужна внешняя координация —
  в godoc). Порядок: слот principal → бюджет → слот аккаунта, ожидание бюджета не держит слот аккаунта.
  Вес вне [1, Budget] прижимается. X-RateLimit-* не используются.
- Send берёт слоты на каждый вызов и отпускает после чтения body. Вызов с WithAuth не расходует бюджет
  клиента; `WithLimiter` учитывает его в бюджете своего principal.
- Решение: без равномерного pacing и запаса окна — ровно правило источника; 429 обрабатывает T09.
- Тесты на фейковых часах без sleep: веса и границы дат, окно ±1 нс, отмена в очереди, слоты аккаунта,
  clamp веса, путь с префиксом, WithAuth/WithLimiter.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Для docs-коммита: термины principal/Profile в meta/GLOSSARY.md. Для F01: stream держит слоты до Close.

## Comments

Ревью 46e9d06 (Standards P1): бюджет списывался до ожидания слота аккаунта. Исправлено в `6a07889`: окно проверяется без списания, единицы списываются после получения слота аккаунта; регрессия `TestLimiterBudgetChargedAtSend` (на старом коде падает).
