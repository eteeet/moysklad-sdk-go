Status: resolved
Blocked by: 01

# 04: Обмен одноразового токена пользователя

## What to build

`Service.UserContext(ctx, token, opts ...transport.RequestOption)`: POST `context/user` с `{"token"}` → `UserContext{AccountID,
UserID, UserUID, Role}`.

## Acceptance criteria

- [x] Пустой токен отклоняется до HTTP; токен передаётся только в теле.
- [x] Пример ответа документации; оба codecs.
- [x] HTTP-тесты: 200, 400 (3008), 401, 404 (3007) как APIError.
- [x] Повторов нет: токен одноразовый. Тест: с RetryPolicy и ответом 5xx ровно одна попытка.
- [x] ADR 0001: `ResetResponseInfo(opts)` до валидации — тест: пустой токен обнуляет заранее заполненный `ResponseInfo`; при APIError он содержит status и headers.
- [x] Документация и coverage обновлены; `make check` проходит.

## Blocked by

01 (Service и AppAuth).

## Answer

`Service.UserContext`, фикстура `user-context.json`; ровно одна попытка (RetryPolicy POST не повторяет; `WithoutRetry` убран по ревью как непроверяемый), обязательны все четыре поля ответа, токен только в теле.
