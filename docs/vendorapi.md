# Vendor API 1.0

Пакет `vendorapi` — Vendor API 1.0 серверного решения: запросы к МоемуСкладу и проверка входящих
запросов жизненного цикла. Это отдельный
API со своим BASE-ENDPOINT `vendorapi.DefaultBaseURL`
(`https://apps-api.moysklad.ru/api/vendor/1.0`) и своей авторизацией; JSON API 1.2 он не заменяет.
Отдельного transport нет: клиент — обычный `transport.Client`.

```go
auth, err := vendorapi.AppAuth(appUID, secretKey)
client, err := transport.NewClient(transport.Config{BaseURL: vendorapi.DefaultBaseURL, Auth: auth})
vendor, err := vendorapi.NewService(client)
employee, err := vendor.Context(ctx, contextKey)
```

`NewService(client)` отклоняет nil client. Limiter JSON API к Vendor API не относится: лимиты
Vendor API не документированы, по умолчанию limiter не задавайте.

## Подпись запросов

`AppAuth(appUID, secretKey)` возвращает `transport.Auth` и отклоняет пустые аргументы. Каждый
вызов Auth — каждая попытка запроса — подписывает новый JWT HS256:

- header `{"alg":"HS256","typ":"JWT"}`;
- payload `sub` = appUID, `iat` = текущее время в секундах, `jti` = 32 случайных байта
  (`crypto/rand`) в hex;
- `exp` не передаётся: МойСклад считает срок жизни `iat + 300` с.

Повтор `jti` МойСклад отклоняет, поэтому токен не кэшируется. Часы сервера решения должны быть
синхронизированы: допуск расхождения документация не задаёт.

> **Обязательное правило.** `secretKey` решения служит только Vendor API. Токены, которые
> решение выпускает своим сотрудникам, фронтенду или сервисам, подписывайте отдельным
> ключом ([RFC 8725 §3.12](https://www.rfc-editor.org/rfc/rfc8725#section-3.12)). Иначе
> собственный токен решения с теми же claims пройдёт проверку как запрос МоегоСклада, и
> активацию можно будет подделать. SDK этого запретить не может.

## Методы

| Метод | Запрос | Результат |
| --- | --- | --- |
| `Context(ctx, contextKey, opts...)` | POST `context/{contextKey}`, тело пустое | `*contexts.EmployeeContext` |
| `UserContext(ctx, token, opts...)` | POST `context/user` `{"token"}` | `*vendorapi.UserContext` |
| `Status(ctx, appID, accountID, opts...)` | GET `apps/{appId}/{accountId}/status` | `*vendorapi.AppStatus` |
| `SetStatus(ctx, appID, accountID, status, opts...)` | PUT `apps/{appId}/{accountId}/status` `{"status"}` | `error` |

Общее для всех методов: при ошибке результат nil, `errors.As` распознаёт `*sdkerrors.APIError`,
`errors.Is` сохраняет ошибки context, авторизации и транспорта. `transport.WithResponseInfo`
сбрасывается до проверки аргументов и заполняется и при APIError. `RetryPolicy` повторяет
только GET/HEAD: из четырёх методов повторяется только `Status`, каждая попытка — с новым `jti`.

### Context

`Context` возвращает сотрудника, открывшего iframe или виджет решения. contextKey приходит
GET-параметром при загрузке окна. Ответ по форме совпадает с `context/employee` JSON API,
поэтому тип — `contexts.EmployeeContext` ([access](access.md)). Признак администратора аккаунта —
`Permissions["admin"].Actions["view"] == model.ScopeAll`. Пустой contextKey, `.`, `..`, `user` (маршрут
`UserContext`) или ключ с `/`, `?`, `#` отклоняются до HTTP.

Успех — HTTP 200 и непустой `id`. Документированные ошибки: 403 — решение не авторизовано по
этому contextKey, 404 — ключ не найден или истёк.

### UserContext

Рекомендуемый источник сотрудника для главного окна, чатов, модальных окон и виджетов: виджет
получает у хост-окна одноразовый токен (живёт 60 с) и передаёт его на сервер решения, сервер
обменивает его этим вызовом. Пустой токен отклоняется до HTTP; токен идёт только в теле. Запрос
выполняется ровно один раз: `RetryPolicy` POST не повторяет, а после потерянного ответа токен уже
израсходован, повтор вернул бы 404.

`UserContext{AccountID, UserID, UserUID, Role}`; `Role` — открытый строковый тип: `RoleAdmin`,
`RoleCashier`, `RoleWorker`, `RoleIndividual`. Успех — HTTP 200 и все четыре поля непустые: источник помечает их обязательными. Ошибки: 400 (3008) — пустой
токен, 401 — ошибка JWT решения (токен не израсходован, обмен можно повторить), 404 (3007) —
токен использован, истёк или выдан другому решению. Токен не храните в `localStorage`,
`sessionStorage`, cookie и не пишите в логи.

### Status и SetStatus

`appID` и `accountID` проверяются как UUID до HTTP. `Status` возвращает
`AppStatus{Status, Cause, Subscription, Access}`; успех — HTTP 200 и непустой `status`, 404
(2004) — решение не установлено на аккаунте. `Access` заполнен только в `Activated` и
`SettingsRequired` и токена не содержит; `resource` в примере источника — со слешем в конце. `scope` таблица источника
называет String, а пример — массивом; SDK читает массив, как в примерах активации.

`SetStatus` сообщает итог активации, если решение ответило на PUT `Activating` или
`SettingsRequired`; после ответа `Activated` вызывать его не нужно. Допустимы только
`StatusActivating`, `StatusSettingsRequired`, `StatusActivated`, иначе ошибка до HTTP. Успех —
HTTP 200 с пустым телом, в том числе если статус уже целевой; 404 — не установлено, 409 —
перехода по жизненному циклу нет. PUT не повторяется.

## Входящие запросы МоегоСклада

МойСклад вызывает решение по `{endpointBase}` из дескриптора: PUT `ActivationPattern` —
активация (Install, Resume) и смена подписки (TariffChanged, Autoprolongation), DELETE того же
пути — деактивация. Обработчик пишет решение; пакет даёт проверку JWT, типы тел и шаблон для
`http.ServeMux`. Если в `endpointBase` есть путь (`https://host/app`), снимите его
`http.StripPrefix("/app", mux)`: шаблон начинается с `/api/moysklad/…`.

```go
mux.HandleFunc("PUT "+vendorapi.ActivationPattern, func(w http.ResponseWriter, r *http.Request) {
	claims, err := vendorapi.VerifyRequest(r, secretKey, time.Now())
	if err != nil { /* 401 */ }
	// решение: отклонить повтор claims.ID до exp, разобрать vendorapi.ActivationRequest,
	// r.PathValue("appId"), r.PathValue("accountId"); ответ vendorapi.ActivationResponse
})
```

`VerifyToken(token, secretKey, now)` проверяет по порядку:

1. пустой `secretKey` — ошибка конфигурации (не `ErrInvalidToken`);
2. только символы base64url и точки (`\r`, `\n` и пробелы — отказ: `encoding/base64` их
   пропускает), ровно три части без padding, канонические (`Strict`);
3. header — JSON-объект с `alg` строго `HS256`, без `crit` (расширения SDK не поддерживает,
   RFC 7515 §4.1.11); `typ` необязателен, а если есть — строка `JWT` без учёта регистра
   (RFC 7519 §5.1, RFC 7515 §4.1.9). Документация называет `typ` обязательным, но живой
   МойСклад 2026-11-05 строгую проверку не прошёл. Ошибка называет `alg` и `typ` токена;
4. подпись HMAC-SHA256 на `secretKey`, сравнение constant-time — до разбора payload;
5. payload — JSON-объект без повторов имён; `iat` и `exp` — JSON-числа от 0 до 9999 года
   (дробные секунды допустимы), `jti` — непустая строка; `null` и иные типы — отказ;
6. `now >= exp` — `ErrTokenExpired`: токен действителен строго до `exp` (RFC 7519 §4.1.4).

Допуска на расхождение часов нет, прочие claims игнорируются. `iat` не сравнивается ни с
`now`, ни с `exp`: окно `exp − iat` источник не задаёт, поэтому токен с далёким `exp`
проходит. Строже — проверяйте `ExpiresAt.Sub(IssuedAt)` в решении. Ошибки 2–5 — `ErrInvalidToken`
(`errors.Is`); секрет в текст ошибки не попадает. `VerifyRequest` требует ровно один заголовок
`Authorization` со схемой Bearer (регистр схемы не важен, RFC 7235) и вызывает `VerifyToken`.

Результат — `TokenClaims{IssuedAt, ExpiresAt, ID}`. **Повтор `jti` отклоняет решение**:
SDK состояния не хранит. Достаточно помнить `ID` до `ExpiresAt`. МойСклад повторяет
запросы при 5xx и таймауте 10 с, поэтому обработка должна быть идемпотентной; id запроса
при повторе не меняется — для логов берите `vendorapi.RequestID(r)`, не читайте заголовок сами.
Документация называет его `X_Lognex_RequestId`, а живой МойСклад (2026-11-05) шлёт
`X-Lognex-RequestId` через дефис; `RequestID` читает оба, дефис первым. Для дефисного имени
`underscores_in_headers` в nginx не нужен.

### Типы тел

| Тип | Тело |
| --- | --- |
| `ActivationRequest{AppUID, AccountName, Cause, Access, Subscription, Additional}` | PUT: активация и смена подписки |
| `DeactivationRequest{Cause, AppUID, AccountName}` | DELETE деактивации; ответ — пустое тело, 200 или 204 |
| `ActivationResponse{Status}` | ответ 200 на PUT с `Content-Type: application/json` |

Имена `ActivationRequest` и `ActivationPattern` повторяют раздел «Активация» источника, хотя по
glossary TariffChanged и Autoprolongation — не активация, а смена подписки без токена.

- `Access{Resource, Scope, Permissions, AccessToken}` не приходит при `TariffChanged`,
  `Autoprolongation` и у решений без доступа к API. `Permissions` — `jsontext.Value`, только
  при scope `custom`: источник описывает его массивом, а пример — объект.
- `Subscription{TariffID, Trial, TariffName, ExpiryMoment, NotForResale, Partner}` и
  `Additional{FiscalAPI{ID, Token}}` необязательны (nil).
- `Cause` (`CauseInstall`, `Resume`, `TariffChanged`, `Autoprolongation`, `Uninstall`, `Suspend`,
  `PermissionsChanged`) и `Status` (`StatusActivating`, `ActivationFailed`, `SettingsRequired`,
  `Activated`, `Deactivating`, `DeactivationFailed`, `Suspended`) — открытые строковые типы:
  неизвестные значения сохраняются. Ответить на активацию можно только `Activating`,
  `SettingsRequired` или `Activated`.
- `AccessToken` и `FiscalAPI.Token` — тип `vendorapi.Secret`: `fmt` с любым глаголом (и ошибочным
  `%d`), `slog` (текстовый и JSON) и `json.Marshal` выводят `<redacted>`. Значение достаётся
  явно — `string(s)`. Поэтому `ActivationRequest` после `json.Marshal` токенов не содержит:
  сохраняйте их отдельно.

## Вне пакета

Событие `PermissionsChanged` (путь `…/event` в источнике противоречив), Loyalty API, кнопки,
сценарии и Fiscal API не реализованы. Хранение `jti`, токенов и статусов, HTTP-обработчик и
дескриптор XML остаются у решения.

Контракт сверён с [официальной документацией на revision
be7188e50e29b4ab9acb7e85d0cda08455d0b9fc](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_moysklad_endpoints.md),
[аутентификацией](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_authentication.md)
и [эндпоинтами решения](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_vendor_endpoints.md).
Фикстуры — примеры документации: `testdata/custom/vendorapi`. Live-проверки не было.
