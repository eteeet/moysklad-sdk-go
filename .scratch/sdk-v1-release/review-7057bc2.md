# Release review — 2026-11-02

BASE: `5b9723cd7800d1aa41e30cc1b3028c604d910997`.
HEAD: `7057bc2b3e7a040d4dc79b2c54836a0edf3d9a3e`.
Scope: R01–R04 и завершающий docs commit; 5 commits, 48 files.
Метод: независимые Standards/Spec reviewers и полный make check.
Итог: блокирующих замечаний нет; одно необязательное heuristic замечание.

## Standards

Нарушений документированных стандартов нет.

**Возможный Duplicated Code, низкий приоритет, не blocker.**
В API-error tests повторяется настройка httptest server с HTTP 401 и создание клиента
через env.ClientFor:

- `examples/getting_started/main_test.go:42`;
- `examples/customerorder_positions/main_test.go:47`;
- `examples/audit_events/main_test.go:53`;
- `examples/turnover/main_test.go:40`.

При развитии общей error fixture можно вынести создание тестового клиента с заданным
ответом в examples/internal/exampletest, сохранив assertions сценариев локальными.
Это эвристика smell baseline, не обязательное требование рефакторинга.

Facade Middle Man/Message Chains не отмечаются: утверждённый R01 interface требует
именно этого слоя. Приватные do/validID исключены по правилу docs/coverage.md.
Общий transport, отсутствие HTTP в конструкторе, upload opt-in, O_EXCL/cleanup
частичного файла, неизменность upstream и SHA-256 provenance проверены.

## Spec

Замечаний нет.

- R01: все выбранные services доступны через facade и делят один transport;
  standalone constructors сохранены, скрытых HTTP и автоматического token issuance нет.
- R02: примеры соответствуют сценариям, запись требует opt-in, ошибки передаются caller.
  Assortment walker явно получает beta-tag на последующие страницы.
- R03: make check, toolchain, SHA-256 manifests и CHANGELOG соответствуют scope.
  Отсутствие CI workflow разрешено пользователем; tag/publish не выполняются.
- R04: статусы методов сверены с реализацией; ограничения и optional gates видимы.

## Verification

- `make check` (CHECK_P=2) — PASS: gofmt, vet, полный race-прогон, fixtures/status/examples,
  git diff --check всего tracked HEAD от пустого дерева и изменений относительно HEAD.
- HEAD оставался 7057bc2; отслеживаемые файлы при старте проверки без изменений.
- Существующий untracked `.scratch/code-simplification/` не относится к срезу,
  не изучался и не изменялся. Неотслеживаемые файлы не проверяются git diff --check.
- Production code и статусы tickets reviewer не менял; live API, tag/publish не выполнялись.

Standards: 0 нарушений, 1 необязательная эвристика низкого приоритета.
Spec: 0 замечаний.

## Ответ и исправления

Duplicated Code — принято. Повтор шире указанного: связка httptest-сервер + `env.ClientFor`
стояла в 10 тестах шести примеров, тело 401 — в четырёх.

- `acfd562` test: `examples/internal/exampletest.Client(t, handler)` поднимает сервер, закрывает
  его в `t.Cleanup` и собирает клиент через `env.ClientFor`; `exampletest.Unauthorized` — общий
  ответ 401/1056. Href в ответах строятся из `r.Host`, поэтому хелпер возвращает только клиент.
  Assertions сценариев остались в тестах; пакет `testing` в `env` (и в бинарники примеров)
  не попал. Итог: −71 строка в тестах примеров, поведение тестов прежнее.
- Проверка: `make check` (CHECK_P=2) — PASS.
