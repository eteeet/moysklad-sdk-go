Status: needs-triage
Blocked by: None

# 01: AppAuth и контекст пользователя по contextKey

## What to build

Решение подписывает запросы к Vendor API через `vendorapi.AppAuth` и получает сотрудника
iframe через `Service.Context`. Первый законченный срез: подпись, сервис, один метод,
тесты и документация. Контракт: [spec](../spec.md), [research](../research.md).

## Acceptance criteria

- [x] `DefaultBaseURL`, `AppAuth`, `Service`, `NewService` по spec; nil client и пустые appUID/secretKey отклоняются.
- [x] Исходящий JWT разобран в тесте независимо (crypto/hmac, без `VerifyToken`): header `{"alg":"HS256","typ":"JWT"}`, `sub`=appUID, `iat` в пределах времени теста, `jti` — 64 hex-символа, `exp` отсутствует.
- [x] Два последовательных запроса получают разные `jti` (повтор попытки проверяется в 03 на GET).
- [x] `Context(ctx, contextKey, opts ...transport.RequestOption)` выполняет POST `context/{contextKey}` без тела и возвращает `*contexts.EmployeeContext`; неверный ключ отклоняется до HTTP.
- [x] ADR 0001: `transport.ResetResponseInfo(opts)` до валидации — тест: ошибка неверного ключа обнуляет заранее заполненный `ResponseInfo`; при 403 `ResponseInfo` содержит status и headers.
- [x] POST не повторяется: с RetryPolicy и ответом 5xx ровно одна попытка.
- [x] Пример ответа из документации в `testdata/custom/vendorapi` с provenance; оба codecs.
- [x] HTTP-тесты: маршрут, заголовки, gzip, 200, 403/404 как APIError, отмена context.
- [x] `docs/vendorapi.md`, README и coverage (этап 18) отражают только реализованное; в docs обязательное правило: `secretKey` решения только для Vendor API, свои токены — отдельным ключом (RFC 8725 §3.12).
- [x] `make check` проходит.

## Blocked by

Нет.
