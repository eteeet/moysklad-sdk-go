# Foundation review — 2026-11-01

BASE: `74428035ed407cb06a8a80420e09edc4829057da`.
HEAD: `46e9d068dda6401fe6cb1306eb6c7af55139bc41`.
Scope: T01–T13 и последний docs commit; 14 commits, 101 files.
Метод: независимые Standards/Spec reviewers, дополнительные публичные воспроизведения.
Итог: требуются исправления; 4 уникальных замечания (2 P1, 2 P2).

## Standards

1. **P1 — rate budget истекает до фактической отправки.**
   `transport/limiter.go:171–179`: admission резервируется до ожидания общего
   `Account` semaphore. Если аккаунт занят другим principal дольше Window, после
   освобождения отправляется старый запрос и одновременно разрешается новый бюджет
   того же principal. Нарушен rolling-window контракт из общей spec, раздел Limits.
   Публичное воспроизведение с fake clock: Budget=1/Window=3s, обе фактические
   отправки в `12:00:04`, интервал `0s`.
   Согласовать admission rate budget с получением account slot, сохранив правило
   «ожидание бюджета не занимает account slot»; добавить интеграционный regression test.

2. **P2 — локальная ошибка сохраняет metadata предыдущего вызова.**
   `entity/uom/service.go:28–31`: валидация возвращает ошибку до Client.Send, где
   обнуляется ResponseInfo. Аналогично устроены остальные services/List/GetPositions.
   Нарушен ADR 0001: «Если ответа нет, info обнуляется».
   Публичный repro: начальные StatusCode=503/Attempts=3, затем Get с invalid UUID —
   ошибка, но info остаётся 503/3 с прежними headers.
   Сбрасывать metadata на входе публичного вызова, включая preflight errors.

Документированные контракты нарушены; косметических smell-замечаний нет.

## Spec

1. **P1 — async submission повторяется при URL-encoding параметра.**
   `transport/retry.go:73–77`. T09: «POST/PUT/security-token/async-submit — без
   автоматического повтора». Проверка подстроки сырого RawQuery пропускает
   `?async=%74rue` и `?%61sync=true`, хотя сервер получает async=true.
   После неоднозначного сбоя это может создать вторую async task.
   Публичное воспроизведение с ответом 429: буквальное async=true — 1 отправка;
   оба encoded-варианта — 2. Классифицировать итоговый decoded query после
   объединения endpoint/query, включая корректную обработку повторных параметров.

2. **P2 — fields=minimumStock принимается, но результат теряется.**
   `entity/product/service.go:81–84`. T11: «Покрыть supported … fields по endpoint».
   List объявляет minimumStock поддержанным, однако Product не содержит поля или
   raw fallback. Ответ с minimumStock `{type: "ALL_WAREHOUSE_SUM", quantity: 7}`
   успешно декодируется, но запрошенные данные недоступны caller.
   Добавить read-модель и HTTP/JSON проверки либо явно отложить поддержку и
   отклонять этот параметр.

3. **P2 — ResponseInfo сохраняет чужой ответ при локальной ошибке.**
   `entity/product/service.go:28–30`, `internal/list/page.go:43–48`;
   аналогично остальные services/positions/security.
   ADR 0001: «Если ответа нет, info обнуляется».
   После успешного List ошибочный Get или List с Limit=1001 оставляет StatusCode=200,
   Attempts=1 и headers прежнего успеха. Это то же замечание, что Standards №2;
   оно независимо обнаружено обеими проверками.

## Verification

- `go test -race -count=1 -p 1 ./...` — PASS.
- `go vet ./... ./testdata` — PASS; `go test ./testdata` — PASS (no test files).
- gofmt для существующих изменённых Go files — PASS; `git diff --check` — PASS.
- В начале и конце проверки HEAD оставался 46e9d06, рабочее дерево было чистым.
- Публичные воспроизведения: `/private/tmp/sdk_review_limiter.go`,
  `/private/tmp/sdk_review_async.go`, `/private/tmp/sdk_review_info.go`,
  `/private/tmp/foundation-spec-repro.go`.
- Live API не использовался. Production code и статусы tickets не изменялись.

Standards: 2 замечания, самое серьёзное — rate budget (P1).
Spec: 3 замечания, самое серьёзное — retry async submission (P1).

## Ответ и исправления

Все замечания приняты; Spec №3 совпадает со Standards №2. Отдельные fix-коммиты поверх 46e9d06,
коммиты tickets не переписывались.

| Замечание | Commit | Регрессия |
| --- | --- | --- |
| Standards P1, бюджет до слота аккаунта | `6a07889` | `TestLimiterBudgetChargedAtSend` |
| Standards P2 / Spec P2, ResponseInfo при локальной ошибке | `bc2dad4` | `TestLocalErrorResetsResponseInfo`, `TestCreateTokenLocalErrorResetsInfo` |
| Spec P1, async при URL-encoding | `ebe0426` | `TestRetryableRequest`, `TestSendRetryAsyncEncoded` |
| Spec P2, minimumStock | `1fac290` | `TestServiceListRejectsMinimumStock` |

Каждая регрессия проверена на коде 46e9d06 — падает, на исправленном — проходит.
Ограничение: ResponseInfo обнуляется при создании опции; переиспользование одного значения опции
между вызовами не поддерживается для метаданных (ADR 0001).
