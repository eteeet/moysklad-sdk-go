# Vendor API 1.0

Пакет `vendorapi` — запросы серверного решения к МоемуСкладу по Vendor API 1.0. Это отдельный
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
> приложение выпускает своим пользователям, фронтенду или сервисам, подписывайте отдельным
> ключом ([RFC 8725 §3.12](https://www.rfc-editor.org/rfc/rfc8725#section-3.12)). Иначе
> собственный токен приложения с теми же claims пройдёт проверку как запрос МоегоСклада, и
> активацию можно будет подделать. SDK этого запретить не может.

## Методы

| Метод | Запрос | Результат |
| --- | --- | --- |
| `Context(ctx, contextKey, opts...)` | POST `context/{contextKey}`, тело пустое | `*contexts.EmployeeContext` |

`Context` возвращает сотрудника, открывшего iframe или виджет решения. contextKey приходит
GET-параметром при загрузке окна. Ответ по форме совпадает с `context/employee` JSON API,
поэтому тип — `contexts.EmployeeContext` ([access](access.md)). Признак администратора аккаунта —
`Permissions["admin"].Actions["view"] == model.ScopeAll`. Пустой contextKey или ключ с `/`, `?`,
`#` отклоняются до HTTP.

Успех — HTTP 200 и непустой `id`. Документированные ошибки: 403 — решение не авторизовано по
этому contextKey, 404 — ключ не найден или истёк. При ошибке результат nil, `errors.As`
распознаёт `*sdkerrors.APIError`, `errors.Is` сохраняет ошибки context, авторизации и транспорта.
`transport.WithResponseInfo` сбрасывается до проверки аргументов и заполняется и при APIError.
`RetryPolicy` повторяет только GET/HEAD, поэтому POST `Context` выполняется ровно один раз.

## Вне пакета

Активация и деактивация (входящие запросы МоегоСклада), статус установки, обмен одноразового
токена, Loyalty API, кнопки, сценарии и Fiscal API пока не реализованы; план —
`.scratch/vendor-api/`. Хранение `jti`, токенов и статусов, HTTP-обработчик и дескриптор XML
остаются у приложения.

Контракт сверён с [официальной документацией на revision
be7188e50e29b4ab9acb7e85d0cda08455d0b9fc](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_moysklad_endpoints.md)
и [аутентификацией](https://github.com/moysklad/api-vendor-1.0-doc/blob/be7188e50e29b4ab9acb7e85d0cda08455d0b9fc/md/vendor_api/_authentication.md).
Фикстура ответа: `testdata/custom/vendorapi`. Live-проверки не было.
