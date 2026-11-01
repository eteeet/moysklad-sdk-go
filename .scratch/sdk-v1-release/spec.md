Status: needs-triage
Stage: DRAFT

# 7. Фасад, примеры и готовность v1

## Problem Statement

Собрать стабильный публичный interface SDK, проверить согласованность всех семейств и
подготовить воспроизводимый первый выпуск.

## Solution

Завершённые вертикальные срезы в составе v1; каждый ticket включает models, methods,
тесты и собственное описание полей. Матрица ниже определяет scope этой захватки.
[Общая карта](../sdk-v1/map.md), [нормативные ограничения](../sdk-v1/spec.md),
[каталог JSON-контрактов](../sdk-v1/contracts.md). Текущая работа — планирование.

## User Stories

1. Как разработчик решения, я могу использовать «Единая точка входа SDK» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
2. Как разработчик решения, я могу использовать «Примеры основных сценариев» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
3. Как разработчик решения, я могу использовать «CI и проверки выпуска» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
4. Как разработчик решения, я могу использовать «Приёмка и статус v1» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.

## Implementation Decisions и scope tickets

| ID | Ticket | Blocked by |
| --- | --- | --- |
| R01 | [Единая точка входа SDK](issues/01-facade.md) | T01, T02, T03, T04, T05, T06, T07, T08, T09, T10, T11, T12, T13, A01, A02, A03, A04, A05, A06, A07, A08, C01, C02, C03, C04, C05, C06, C07, C08, D01, D02, D03, D04, D05, E01, E02, E03, E04, E05, E06, F01, F02, F03, F04, F05 |
| R02 | [Примеры основных сценариев](issues/02-examples.md) | R01 |
| R03 | [CI и проверки выпуска](issues/03-ci-release.md) | R01 |
| R04 | [Приёмка и статус v1](issues/04-acceptance.md) | R02, R03 |

### R01: Единая точка входа SDK

Root moysklad.NewClient(Config) создаёт typed handles над одним transport/policy.
Рекомендованные пути: client.Product, client.CustomerOrder, client.Context.Employee,
client.Context.Application, client.Audit, client.Metadata, client.States, client.Files,
client.Reports.Turnover, client.Security. Точный naming принять до реализации; конструктор не
делает HTTP и не получает токен.

Проверки и особенности: Каждый существующий и новый service доступен, все используют один
HTTPClient/auth/limiter. Старые standalone constructors работают; no global singleton/auth
credentials leakage. Method/API compatibility examples и zero hidden network at startup.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### R02: Примеры основных сценариев

Компилируемые примеры: поиск Assortment, CustomerOrder→страница positions, Audit contexts→явные
events, Turnover, Files и client initialization для personal/App token. Запись защищена явным
opt-in; token issuance не часть обычного getting-started.

Проверки и особенности: Examples unit tests с httptest и компиляция CI, no real credentials/live
writes; pagination/close/cancellation показаны корректно. Нет автоматического выполнения
examples Makefile в тестовом прогоне.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### R03: CI и проверки выпуска

Go1.27 toolchain, tests/race/vet/gofmt/diff check, manifest хешей fixtures и проверка public
examples; migration notes/API compatibility и CHANGELOG/v1 version proposal. Существующие
команды сохраняются; network/live интеграция не обязательна для unit CI.

Проверки и особенности: Чистый checkout воспроизводит checks, pin toolchain/workflow
permissions, нет секретов и live writes. При локальных TCP resource flakes допустим -p1 с явным
reporting. Publish/tag/release не выполняются без отдельной команды.

Источники: [транспортный research](../sdk-v1/research/transport.md).

### R04: Приёмка и статус v1

Проверить coverage-manifest каждого исходного требования пользователя, закрыть или явно вынести
optional gates, обновить README/status-entities/status-technical/coverage без фиктивных checked
для неподдержанных methods. v1 — выбранный профиль SDK, не весь JSON API.

Проверки и особенности: По тикету коммит, dependencies shared-types перед dependent, docs
отдельным последним commit волны; planning не коммитится без команды. Независимые Standards/Spec
review каждой волны; evidence links и unresolved limitations. Cleanup core worktrees —
существующий отдельный ticket, не blocker SDK.

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

- Стабильный interface утверждается до v1.0 tag; полный Vendor API, расширение на все сущности и
  все отчёты не входят.
- Опциональный live read-only smoke тест возможен только по отдельному согласию; его отсутствие
  не выдаётся за выполненный прогон.

## Out of Scope

Expand и implicit follow-up; все не перечисленные в tickets ресурсы/операции.
План не объявляет общую CRUD поддержку у ресурсов с другой формой interface.
