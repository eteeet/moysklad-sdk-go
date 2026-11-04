Status: resolved
Type: research

# Vendor API 1.0: проверенный контракт

Проверка 2026-11-04. Источник — [moysklad/api-vendor-1.0-doc](https://github.com/moysklad/api-vendor-1.0-doc),
master = `be7188e50e29b4ab9acb7e85d0cda08455d0b9fc` (совпадает с `doc-version` сайта
https://dev.moysklad.ru/doc/api/vendor/1.0/, файлы сверены побайтно). Официальные SDK
(java/python/typescript/php) Vendor API не реализуют; есть только демо-приложения.
С валидным JWT живые запросы не выполнялись.

## Проверенные факты

| Факт | Источник |
| --- | --- |
| Решение → МойСклад: `Authorization: Bearer`, JWT HS256 на secretKey; payload `sub`=appUid (обяз.), `iat` (обяз.), `exp` (необяз.), `jti` (обяз., повтор отклоняется); срок не больше 300 с: `effectiveExp = exp ? min(exp, iat+300) : iat+300` | [_authentication.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_authentication.md) |
| МойСклад → решение: JWT HS256 тем же secretKey; header `alg`=HS256 и `typ`=JWT обязательны; payload только `iat`, `exp`, `jti` (NumericDate — JSON Number, без `sub`); решение обязано отклонять после `exp` и при повторе `jti` | [_authentication.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_authentication.md) |
| Ошибки — формат JSON API 1.2: `{"errors":[{"error","code"}]}`, заголовки `X-Lognex-Auth`, `X-Lognex-Auth-Message`; без `Accept-Encoding: gzip` — 415 (проверено вживую на 401-ответе) | [_authentication.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_authentication.md) |
| Базовый URL запросов решения: `https://apps-api.moysklad.ru/api/vendor/1.0` | [_moysklad_endpoints.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_moysklad_endpoints.md) |
| `POST context/{contextKey}`, тело пустое → Employee в форме remap `context/employee`; 200/403/404; contextKey живёт ≤5 мин | [_moysklad_endpoints.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_moysklad_endpoints.md), [_user_context.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/developer_guide/_user_context.md) |
| `POST context/user` `{"token"}` → `{accountId,userId,userUid,role}`; токен одноразовый, 60 с; 400 (3008)/401/404 (3007) | [_moysklad_endpoints.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_moysklad_endpoints.md) |
| `GET apps/{appId}/{accountId}/status` → `{status, cause?, subscription?, access?[]}`; 200/404 (2004). `PUT …/status` `{"status"}` → пусто; 200 (и если статус уже целевой)/404/409 | [_moysklad_endpoints.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_moysklad_endpoints.md) |
| Статусы: Activating, ActivationFailed, SettingsRequired, Activated, Deactivating, DeactivationFailed, Suspended; решение может выставить только Activating/SettingsRequired/Activated | [_moysklad_endpoints.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_moysklad_endpoints.md) |
| Входящие запросы: базовый путь `{endpointBase}/api/moysklad/vendor/1.0`; заголовок `X_Lognex_RequestId` (не меняется при повторе); обработка должна быть идемпотентной | [_vendor_endpoints.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_vendor_endpoints.md) |
| Активация `PUT apps/{appId}/{accountId}`: `appUid`, `accountName`, `cause` (Install/Resume/TariffChanged/Autoprolongation), `access[]?` {resource, scope[], permissions?, access_token}, `subscription?` {tariffId, trial, tariffName?, expiryMoment?, notForResale, partner}, `additional.fiscalApi?`. Ответ `{"status"}` + `Content-Type: application/json` | [_activation.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_activation.md), [_vendor_endpoints.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_vendor_endpoints.md) |
| Деактивация `DELETE` того же URL с JSON-телом: `cause` (Uninstall/Suspend), в примерах ещё `appUid`, `accountName`; токен уже аннулирован; ответ пустой: 200 / 204 (уже отключено) | [_deactivation.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_deactivation.md) |
| Коды ответа решения: 4xx и 551 → Failed (при активации токен удаляется), прочие 5xx и таймаут 10 с → повтор: 3 мин каждые 10 с (активация/деактивация), 24 ч каждые 5 мин (TariffChanged, Autoprolongation, события) | [_retry.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_retry.md) |
| `PUT …/event` с `cause:"PermissionsChanged"` — только при `<additionalEvents><updatePermissions/>` в дескрипторе | [_additional_events.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_additional_events.md) |
| Хост-окно добавляет к `sourceUrl`: `contextKey`, `appUid`, `appId`, `userLocale`; `useContextKey="false"` отключает | [_user_context.md](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/developer_guide/_user_context.md) |

## Противоречия

1. `access[].permissions` описан как Array, в примере — объект; в `GET …/status` `scope`/`permissions` описаны String, в примере — массив и объект.
2. DELETE: спецификация описывает только `cause`, примеры содержат `appUid`, `accountName`.
3. URL событий: шаблон `/api/vendor/1.0/…/event`, пример и демо — `/api/moysklad/vendor/1.0/…/event`.
4. В таблице статусов у Activated пустой cause, в примере `Activated` + `cause:"Install"`.
5. Допуск расхождения часов, формат `jti` и окно `exp − iat` для входящих JWT не описаны.
6. `resource` в примерах встречается со слешем в конце и без.
7. Заголовок `X_Lognex_RequestId` с подчёркиваниями: nginx по умолчанию его отбрасывает (`underscores_in_headers off`) — общее знание, не из документации.

## Вне документации, но важно для SDK

- Каталог `vendor/` в корне Go-модуля зарезервирован под вендоринг: пакет называем `vendorapi`.
- transport вызывает Auth на каждую попытку запроса (`transport/request.go`), поэтому подпись с новым `jti` на каждый запрос совместима с RetryPolicy.
