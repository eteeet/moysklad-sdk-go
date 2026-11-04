Status: needs-triage
Blocked by: None

# 02: Проверка входящего JWT и типы активации/деактивации

## What to build

Решение проверяет JWT запроса МоегоСклада через `VerifyRequest`/`VerifyToken`, читает
тела активации и деактивации и отвечает типизированным статусом. Контракт: [spec](../spec.md).

## Acceptance criteria

- [ ] `VerifyToken`/`VerifyRequest` по spec в порядке проверок spec: пустой `secretKey` — ошибка конфигурации; `alg` строго `HS256`, `typ` строго `JWT`; constant-time подпись до разбора payload; `iat`/`exp` — JSON-числа, `jti` — непустая строка; истечение при `now >= exp`; ошибки через `errors.Is`, секрет не попадает в текст ошибки.
- [ ] Эталонные входящие токены собирает и подписывает тестовый helper (не `AppAuth`).
- [ ] Отрицательные тесты: пустой ключ; alg none/HS384/HS512/RS256; typ отсутствует/null/иной; чужой ключ; изменённый payload; `now == exp` и `now > exp`; нет iat/exp/jti; null в каждом claim; `"exp":"123"`; `jti` пустой или число; header/payload не объект; битый base64; 2 или 4 части; нет Bearer или иная схема.
- [ ] `ActivationRequest`, `Access`, `Subscription`, `Additional`, `DeactivationRequest`, `ActivationResponse`, `Cause`, `Status` по spec; неизвестные значения сохраняются.
- [ ] Примеры документации (Install с access, с permissions, TariffChanged без access, Uninstall, Suspend) дословно в `testdata/custom/vendorapi`; оба codecs; отсутствие access/subscription читается без ошибки.
- [ ] `ActivationPattern` работает с `http.ServeMux` для PUT и DELETE (тест).
- [ ] Документация: проверка `jti` на повтор — обязанность приложения; `secretKey` решения не подписывает собственные токены приложения (RFC 8725 §3.12); `AccessToken` не логировать; `underscores_in_headers` для nginx.
- [ ] `make check` проходит.

## Blocked by

Нет; можно параллельно с 01.
