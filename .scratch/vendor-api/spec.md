Status: needs-triage

# Vendor API 1.0: авторизация решения, контекст пользователя, жизненный цикл

## Problem Statement

Серверное решение МойСклад обязано работать с Vendor API: принимать запросы активации
и деактивации от МоегоСклада, проверять их JWT, узнавать пользователя iframe по
contextKey и сообщать статус установки. В SDK этого нет (этап 18 coverage), поэтому каждое
решение пишет свой клиент и свою проверку JWT. Legacy-решение doc2ms проверяло входящий
JWT только по подписи и сроку тем же ключом, которым подписывало собственные токены
пользователей, — это позволяло подделать активацию. Строгая проверка claims эту дыру не
закрывает: токен решения может содержать те же `iat`/`exp`/`jti`. Закрывает её только
разделение ключей (RFC 8725 §3.12): `secretKey` решения служит исключительно Vendor API.
SDK этого обеспечить не может и фиксирует правило в документации пакета. Нужен один
нейтральный, переиспользуемый контур в SDK без логики конкретного решения.

## Solution

Новый пакет `vendorapi` (не `vendor`: каталог `vendor/` в корне модуля зарезервирован Go).
Он содержит три части:

1. Исходящие запросы решения → МойСклад: `AppAuth` (подпись JWT на каждый запрос)
   и `Service` поверх существующего `transport.Client` с методами Context, UserContext,
   Status, SetStatus.
2. Входящие запросы МойСклад → решение: `VerifyRequest`/`VerifyToken` для JWT и типы тел
   активации, деактивации и ответа, константы путей.
3. Документация пакета и статусы в coverage.

SDK не хранит состояние: защита от повтора `jti`, хранение токенов, маршрутизация HTTP и
бизнес-решение о статусе остаются у приложения. Контракт и источники: [research.md](research.md).

## User Stories

1. Как разработчик решения, я хочу подписывать запросы к Vendor API одной функцией, чтобы не реализовывать JWT вручную.
2. Как разработчик решения, я хочу, чтобы каждый запрос и каждая повторная попытка получали новый `jti`, чтобы МойСклад не отклонял их как повтор.
3. Как разработчик решения, я хочу получить сотрудника по contextKey из iframe, чтобы знать аккаунт и права пользователя.
4. Как разработчик решения, я хочу обменять одноразовый токен пользователя на accountId, userId и роль.
5. Как разработчик решения, я хочу проверить JWT входящего запроса МоегоСклада с отказом на чужой алгоритм и `typ`, подпись, истёкший срок, отсутствующие или неверно типизированные claims и пустой ключ, чтобы активацию нельзя было подделать токеном с другим ключом или форматом.
6. Как разработчик решения, я хочу получить `jti` проверенного токена, чтобы самому отклонять повторы.
7. Как разработчик решения, я хочу типизированно читать тела активации и деактивации, включая отсутствующие `access` и `subscription`, чтобы не разбирать JSON вручную.
8. Как разработчик решения, я хочу отвечать на активацию типизированным статусом.
9. Как разработчик решения, я хочу узнать и выставить статус установки на аккаунте.
10. Как разработчик решения, я хочу получать ошибки Vendor API как `*sdkerrors.APIError`, как в JSON API.
11. Как разработчик решения, я хочу видеть в документации SDK точные границы поддержки Vendor API.
12. Как разработчик решения, я хочу видеть в документации запрет подписывать собственные токены `secretKey` решения, чтобы токен моего приложения не прошёл как запрос МоегоСклада.

## Implementation Decisions

### Исходящие запросы

- `DefaultBaseURL = "https://apps-api.moysklad.ru/api/vendor/1.0"`. Клиент создаётся
  обычным `transport.NewClient` с этим BaseURL и `Auth: vendorapi.AppAuth(...)`;
  отдельного transport нет. Limiter для JSON API к Vendor API не относится: лимиты
  Vendor API не документированы, по умолчанию limiter не задаётся.
- `AppAuth(appUID, secretKey string) (transport.Auth, error)`: пустые значения отклоняются.
  На каждый вызов Auth — новый JWT: header `{"alg":"HS256","typ":"JWT"}`, payload
  `sub`=appUid, `iat`=текущее время, `jti`=32 случайных байта (crypto/rand) в hex; `exp`
  не передаётся (сервер ограничит iat+300). Ошибка crypto/rand возвращается как ошибка Auth.
  Только стандартная библиотека, без JWT-зависимостей.
- `Service`, `NewService(*transport.Client)`; nil client отклоняется.

| Метод SDK | HTTP | Маршрут от BaseURL | Результат |
| --- | --- | --- | --- |
| `Context(ctx, contextKey)` | POST, без тела | `context/{contextKey}` | `*contexts.EmployeeContext` |
| `UserContext(ctx, token)` | POST `{"token"}` | `context/user` | `*UserContext` |
| `Status(ctx, appID, accountID)` | GET | `apps/{appId}/{accountId}/status` | `*AppStatus` |
| `SetStatus(ctx, appID, accountID, status)` | PUT `{"status"}` | `apps/{appId}/{accountId}/status` | `error` |

- Ответ `context/{contextKey}` — тот же Employee, что `context/employee`; переиспользуем
  `contexts.EmployeeContext`. Пустой contextKey или ключ с `/`, `?`, `#` отклоняются до HTTP.
- appId и accountId проверяются как UUID по образцу существующих сервисов.
- `SetStatus` принимает только Activating, SettingsRequired, Activated; иное — ошибка до HTTP.
- Последний параметр каждого метода — `opts ...transport.RequestOption` (ADR 0001).
  Первым действием метод вызывает `transport.ResetResponseInfo(opts)`, затем локальная
  валидация, затем `Client.Send`; `WithResponseInfo` заполняется и при APIError.
- Ошибки — существующий APIError transport. Повторы определяет `RetryPolicy` клиента:
  она повторяет только GET/HEAD, поэтому `Status` может повторяться, а POST `Context` и
  `UserContext` — нет. Для `UserContext` это требование: токен одноразовый, повтор
  после потерянного ответа вернул бы ошибку вместо контекста.

### Входящие запросы

- `VerifyToken(token, secretKey string, now time.Time) (*TokenClaims, error)`, порядок проверок:
  1. пустой `secretKey` — ошибка (не `ErrInvalidToken`: это ошибка конфигурации);
  2. ровно три base64url-части без padding;
  3. header — JSON-объект, `alg` строго `HS256` (отказ на `none`, HS384/512, RS*),
     `typ` строго `JWT` (во входящих он обязателен: research, таблица JWT-header); отказ на отсутствие и null;
  4. подпись HMAC-SHA256 сравнивается constant-time до разбора payload;
  5. payload — JSON-объект; `iat` и `exp` — JSON-числа (NumericDate; дробные секунды допустимы по RFC 7519 §2),
     `jti` — непустая строка; отсутствие, `null`, строка вместо числа и иные типы — отказ;
  6. истечение: `!now.Before(exp)`, то есть `now >= exp` — `ErrTokenExpired` (RFC 7519 §4.1.4).
  Допуска на расхождение часов нет (источник не задаёт). Прочие claims игнорируются.
  `TokenClaims{IssuedAt, ExpiresAt time.Time; ID string}`.
- `VerifyRequest(r *http.Request, secretKey string, now time.Time) (*TokenClaims, error)`:
  извлекает `Authorization: Bearer` и вызывает VerifyToken.
- Ошибки проверки — сигнальные значения (`ErrInvalidToken`, `ErrTokenExpired`) через
  `errors.Is`, без раскрытия секрета в тексте.
- Типы тел: `ActivationRequest{AppUID, AccountName, Cause, Access []Access,
  Subscription *Subscription, Additional *Additional}`,
  `Access{Resource, Scope []string, Permissions json.RawMessage, AccessToken}`,
  `Subscription{TariffID, Trial, TariffName, ExpiryMoment *time.Time, NotForResale, Partner}`,
  `DeactivationRequest{Cause, AppUID, AccountName}`, `ActivationResponse{Status}`.
  `Permissions` — `json.RawMessage`: источник противоречит себе (Array в таблице,
  объект в примере), SDK не выбирает форму за приложение.
- Открытые строковые типы `Cause` и `Status` с константами из research; неизвестные
  значения сохраняются.
- Константы путей для `net/http.ServeMux`: `ActivationPattern =
  "/api/moysklad/vendor/1.0/apps/{appId}/{accountId}"` (PUT активация, DELETE
  деактивация). Обработчик HTTP SDK не предоставляет.
- `AccessToken` — секрет: документация пакета предупреждает не логировать тела запросов.
- Документация пакета обязательным правилом требует: `secretKey` решения — только для
  Vendor API (`AppAuth` и `VerifyToken`); токены, которые решение выпускает своим
  пользователям или сервисам, подписываются отдельным ключом (RFC 8725 §3.12).

### Тесты

- Направления тестируются независимо; круга AppAuth → VerifyToken нет: исходящий токен
  намеренно без `exp`, входящий verifier требует `exp`.
- Исходящий (`AppAuth`): тест сам разбирает токен и проверяет подпись через crypto/hmac;
  header `{"alg":"HS256","typ":"JWT"}`, `sub`=appUID, `iat` в пределах времени теста,
  `jti` — 64 hex-символа, `exp` отсутствует.
- Входящий (`VerifyToken`): эталонные токены собирает и подписывает тест (свой helper на
  crypto/hmac, не AppAuth). Положительный случай и отрицательные: пустой ключ; alg
  none/HS384/HS512/RS256; typ отсутствует/null/иной; чужой ключ; изменённый payload;
  `now == exp` и `now > exp`; нет iat/exp/jti; null; `"exp":"123"`; `jti` пустой или число;
  payload не объект; битый base64; 2 или 4 части; нет Bearer.
- Два последовательных запроса через AppAuth получают разные `jti`. Новый `jti` на каждой
  попытке проверяется на GET `Status` с включённой RetryPolicy (5xx → 200). POST не
  повторяется: тест `UserContext` с RetryPolicy и ответом 5xx видит ровно одну попытку.
- Каждый метод Service: ошибка локальной валидации обнуляет заранее заполненный
  `ResponseInfo`; при APIError `ResponseInfo` содержит status и headers ответа.
- Примеры тел активации/деактивации/статуса — дословно из документации в
  `testdata/custom/vendorapi` со ссылкой на файл и SHA источника; проверка обоими codecs.
- HTTP: локальный сервер и настоящий transport — метод, маршрут, отсутствие тела у Context,
  заголовки, gzip, 200/4xx, APIError.

## Out of Scope

- Loyalty API, кнопки (`button`, `button/complete`), сценарии, Fiscal API (кроме чтения
  `additional.fiscalApi` в теле активации).
- Событие `PermissionsChanged`: только константа Cause; путь `…/event` в источнике
  противоречив.
- Хранение `jti`, токенов и статусов, HTTP-обработчик, дескриптор XML.
- Live-проверка на реальном решении — только после отдельного согласия.

## Further Notes

- `X_Lognex_RequestId` содержит подчёркивания: за nginx нужен `underscores_in_headers on`.
  Это заметка для приложений в docs пакета, не код SDK.
