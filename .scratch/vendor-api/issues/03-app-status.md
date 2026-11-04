Status: resolved
Blocked by: 01, 02

# 03: Статус установки решения

## What to build

`Service.Status` и `Service.SetStatus` по spec: решение узнаёт и выставляет статус
установки на аккаунте.

## Acceptance criteria

- [x] `Status(ctx, appID, accountID, opts ...transport.RequestOption)` и `SetStatus(ctx, appID, accountID, status, opts ...transport.RequestOption)`: GET и PUT `apps/{appId}/{accountId}/status`; UUID проверяются до HTTP.
- [x] ADR 0001: `ResetResponseInfo(opts)` до валидации — тест: неверный UUID и недопустимый статус обнуляют заранее заполненный `ResponseInfo`; при APIError он содержит status и headers.
- [x] GET `Status` с RetryPolicy (5xx → 200): каждая попытка несёт новый `jti`, подпись каждой проверена независимым HMAC.
- [x] `SetStatus` допускает только Activating, SettingsRequired, Activated.
- [x] `AppStatus` читает пример документации (status, cause, subscription, access с permissions как jsontext.Value); оба codecs.
- [x] HTTP-тесты: 200, 404 (2004) и 409 как APIError; тело PUT `{"status":…}`.
- [x] Документация и coverage обновлены; `make check` проходит.

## Blocked by

01 (Service и AppAuth) и 02 (`Status`, `Access`, `Subscription`).

## Answer

`Service.Status`/`SetStatus`, `AppStatus` (переиспользует `Access` и `Subscription`), фикстура `status.json`; новый `jti` на каждой попытке GET проверен.
