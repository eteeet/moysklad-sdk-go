# Транспорт, списки, лимиты и ошибки

Источник контрактов — `api-remap-1.2-doc` revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`.
Решение по интерфейсу опций — [ADR 0001](adr/0001-request-options.md).

## Опции вызова и метаданные ответа

Каждый публичный метод принимает последним аргументом `opts ...transport.RequestOption`.
Обычные вызовы `svc.Get(ctx, id)` не меняются; меняется только тип method value.

| Опция | Заголовок / действие |
| --- | --- |
| `WithHeader(key, value)` | любой заголовок; повтор key — несколько значений, не CSV |
| `WithAcceptRegion(transport.RegionKZ)` | `X-Lognex-Accept-Region` |
| `WithAcceptLanguage("en")` | `Accept-Language` (язык сообщений ошибок) |
| `WithAcceptTimezone(now)` | `X-Lognex-Accept-Timezone`, дата с числовой зоной; DateTime моделей не меняет |
| `WithBetaFeature(tag)` | `X-Lognex-Remap-Beta-Feature`, одно значение |
| `WithWebHookDisableByPrefix(p...)` | `X-Lognex-WebHook-DisableByPrefix`, по значению на префикс |
| `WithDownloadExpiration(ttl)` | `X-Lognex-Download-Expiration-Seconds`, целые секунды (`Client.Download`, [файлы](files-reports.md)) |
| `WithAuth(auth)`, `WithLimiter(l)` | авторизация и бюджет другого principal для одного вызова |
| `WithoutRetry()` | запрет повтора вызова |
| `WithResponseInfo(&info)` | статус, копия заголовков и число попыток последней отправки |

`X-Lognex-WebHook-Disable` отключён API 2026-10-01 и не поддерживается. Pretty-print,
`X-Lognex-Format-Millisecond` и `Cache-Control` задаются через `WithHeader`.
`Authorization`, `Accept`, `Accept-Encoding`, `Content-Type`, `Content-Length`, `Host`,
`Idempotency-Key` задаёт transport: попытка их задать — ошибка до отправки.
Заголовки для всех запросов клиента — `transport.Config.Header`; опция вызова заменяет key целиком.

`ResponseInfo` разбирает `X-RateLimit-*`, `X-Lognex-Retry-TimeInterval`, `X-Lognex-Reset`,
`X-Lognex-Retry-After` (миллисекунды), стандартный `Retry-After` (секунды или HTTP-date) и
`X-Lognex-Content-Timezone`. Некорректное значение даёт `ok=false`, сырое остаётся в `Header`.
Единицы `X-RateLimit-*` после ввода весов источник не уточняет: SDK не считает их бюджетом.

## Вызов без следования redirect

`Client.Send` следует redirect как net/http: POST с ответом 303 превращается в скрытый GET по
Location (307/308 у записи не выполняются, остаются APIError). Для печати
([print.md](print.md)) есть `Client.SendNoRedirect(ctx, method, endpoint, body, opts...)`:

- одна отправка, redirect не выполняется: ответ 3xx возвращается без ошибки, статус и заголовки
  (Location) — в `*http.Response` и `ResponseInfo`; 2xx тоже успех; тело не декодируется;
- ответы вне 2xx/3xx — `APIError`, как у Send; Limiter, Retry (только GET/HEAD), заголовки и
  проверки размера — те же;
- режим действует только на этот вызов (значение в контексте запроса), Client его не хранит;
- Location не проверяется и не запрашивается. Загрузка по нему — отдельный вызов `Download`,
  same-origin контракт которого не ослаблен: ссылка чужого origin отклоняется до отправки;
- неразбираемый Location net/http отвергает без ответа (StatusCode 0); query ссылки из текста
  ошибки удаляется, как у Download.

## Ошибки

`errors/codes.go` содержит все 545 кодов `_errors.md` как `CodeNNNN` (нетипизированные int).
`sdkerrors.HasCode(err, sdkerrors.Code1073)` ищет код во всех `APIError` цепочки,
`sdkerrors.StatusCode(err)` — HTTP-статус. Классифицировать ошибки по тексту не нужно.

## Токен

`security.Service.CreateToken(ctx, login, password)` выполняет `POST security/token` с Basic
даже на клиенте с Bearer. Выпуск отзывает прежние токены сотрудника; запрос не повторяется,
Auth клиента не меняется. `Token` не раскрывается через fmt и slog, но раскрывается json.Marshal.

## Списки, фильтры и страницы

`List(ctx, query.ListOptions, opts...)` есть у 6 справочников и 11 документов; читает одну
страницу `model.Collection[T]` без expand и скрытых запросов. Нулевые поля опций не передаются.
`fields` в List пока не принимается: у Product `minimumStock` отложен до модели неснижаемого
остатка, иначе запрошенные данные терялись бы при декодировании.
`fields=stock` у списков документов работает только вместе с `expand=positions` (вне SDK),
поэтому List его не принимает; строки позиций с остатками — GetPositions.
Документы в корзине и архивные сущности сервер по умолчанию не отдаёт:
`query.Eq("isDeleted", "true"), query.Eq("isDeleted", "false")` отбирают все.

`query.Filter(...)` собирает filter: повтор `=` — OR, `!=` — AND, разные поля — AND;
отклоняются `=` вместе с неравенством того же поля и повтор `<`, `<=`, `>`, `>=`.
`;` в значении экранируется как `\;`. `query.Order(query.Asc("name"), query.Desc("moment"))`.
`query.ReportOptions` — период в московском времени, groupBy, interval для отчётов.

`pages.All(ctx, client, first, opts...)` явно обходит страницы по `meta.nextHref` той же
коллекции, требуя рост offset. Снимка нет: при изменении набора возможны пропуски и дубли.

## Лимиты и повторы

Без настройки SDK не ограничивает частоту и не повторяет запросы.

- `transport.NewLimiter(LimiterConfig{Profile: ...})` в `Config.Limiter`: окно 45 единиц за 3 с
  на principal, до 5 параллельных запросов principal и 20 — `transport.Account`. Профиль
  обязателен: App 1, DraftApp 4, Employee по расписанию (сейчас 3, с 2026-12-01 4), Unknown 4;
  `report/stock/all` и `/bystore` — 5. Даты считаются полуночью по Москве — источник зону не
  указывает. Бюджет списывается при получении слота аккаунта, то есть по времени отправки.
  Один процесс не видит другие процессы и токены аккаунта.
- `Config.Retry = &transport.RetryPolicy{}`: только GET/HEAD (не async) при 429, 500, 502, 503,
  504 и временных сбоях соединения; 3 попытки, full jitter от 500 мс до 30 с, подсказка сервера
  не обрезается, пауза больше `MaxWait` (1 мин) или дедлайна — без повтора. Запись, токен и
  async не повторяются никогда.

## Ограничения размера запроса

Проверяются до отправки в `NewRequest`/`Send`; `Config.MaxBodySize`, `MaxHeaderSize`,
`MaxDepth`: 0 — значение по умолчанию, отрицательное отключает проверку.

| Лимит API | Выбор SDK при неуточнённой единице |
| --- | --- |
| 20 МБ данных в запросе | 20 000 000 байт закодированного JSON (Base64 и UTF-8 по факту) |
| 8 КБ в заголовках | 8 000 байт: request line, Host, Content-Length, все значения заголовков, оценка HTTP/1.1 |
| вложенность 10 (ошибка 2006) | 10 уровней объектов, корень — 1; массивы не считаются |

Размер ответа не ограничивается. Лимит 1000 элементов проверяется только у positions документов.
Массовые операции (1000 объектов, 5 минут) и очередь async (4 задачи на аккаунт, Client её не
считает) — [write-async.md](write-async.md).
