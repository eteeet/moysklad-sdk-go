Status: resolved
Blocked by: None

# 02: Проверка входящего JWT и типы активации/деактивации

## What to build

Решение проверяет JWT запроса МоегоСклада через `VerifyRequest`/`VerifyToken`, читает
тела активации и деактивации и отвечает типизированным статусом. Контракт: [spec](../spec.md).

## Acceptance criteria

- [x] `VerifyToken`/`VerifyRequest` по spec в порядке проверок spec: пустой `secretKey` — ошибка конфигурации; `alg` строго `HS256`, `typ` строго `JWT`; constant-time подпись до разбора payload; `iat`/`exp` — JSON-числа, `jti` — непустая строка; истечение при `now >= exp`; ошибки через `errors.Is`, секрет не попадает в текст ошибки.
- [x] Эталонные входящие токены собирает и подписывает тестовый helper (не `AppAuth`).
- [x] Отрицательные тесты: пустой ключ; alg none/HS384/HS512/RS256; typ отсутствует/null/иной; чужой ключ; изменённый payload; `now == exp` и `now > exp`; нет iat/exp/jti; null в каждом claim; `"exp":"123"`; `jti` пустой или число; header/payload не объект; битый base64; 2 или 4 части; нет Bearer или иная схема.
- [x] `ActivationRequest`, `Access`, `Subscription`, `Additional`, `DeactivationRequest`, `ActivationResponse`, `Cause`, `Status` по spec; неизвестные значения сохраняются.
- [x] Примеры документации (Install с access, с permissions, TariffChanged без access, Uninstall, Suspend) дословно в `testdata/custom/vendorapi`; оба codecs; отсутствие access/subscription читается без ошибки.
- [x] `ActivationPattern` работает с `http.ServeMux` для PUT и DELETE (тест).
- [x] Документация: проверка `jti` на повтор — обязанность приложения; `secretKey` решения не подписывает собственные токены приложения (RFC 8725 §3.12); `AccessToken` не логировать; `underscores_in_headers` для nginx.
- [x] `make check` проходит.

## Blocked by

Нет; можно параллельно с 01.

## Answer

`VerifyToken`/`VerifyRequest`, типы жизненного цикла, `Secret`, `ActivationPattern`, 10 фикстур — коммит a26af61. Сверх spec по ревью: алфавит base64url (CR/LF), канонический base64, отказ на `crit`, повтор имён claims, `bearer` без учёта регистра.
