# SDK v1: transport, headers, limits, query, retry, List и facade

Дата проверки: 2026-11-01. Только исследование и рекомендации; production code не изменялся.

## Источники и воспроизводимость

- Официальный source: [api-remap-1.2-doc](https://github.com/moysklad/api-remap-1.2-doc/tree/9a811cef3dc4a3f3a46f070e60452361a06bf7cd). GitHub API `GET /repos/moysklad/api-remap-1.2-doc/commits?per_page=1` при проверке вернул тот же HEAD `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`, дата commit `2026-10-26T14:19:35Z`; отличие pinned от latest отсутствует на момент проверки.
- Ниже `[G]`, `[R]`, `[E]`, `[A]` — соответствующие первичные файлы pinned revision. Публичные URL разделов приведены для удобства; pinned GitHub links фиксируют проверенный контракт.
- [G — общий контракт](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md), [R — ограничения](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_restrictions.md), [E — каталог ошибок](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_errors.md), [A — async](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_async.md).
- Локальные временные файлы источника: `/private/tmp/moysklad-transport-source/md/`; это не зависимость SDK.

## Инвентаризация headers

Направление: request — SDK → JSON API; response — JSON API → SDK. В таблице перечислены все именованные API headers, найденные в полном `md/` pinned source, плюс стандартные headers, существенные для SDK. Наличие header в workbook example отдельно отмечено и не приравнено к полной спецификации.

| Header | Направление | Проверенный контракт / источник |
|---|---|---|
| `Authorization` | request | `Basic <base64(login:password)>` или `Bearer <token>`; [аутентификация](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-autentifikaciya). Секрет, не писать в logs/metadata. |
| `Accept` | request | Для JSON API поддерживается `application/json;charset=utf-8`; иначе error 1062. [E, error_1062](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_errors.md#L63). Binary download требует отдельного контракта, не навязывать ему JSON Accept. |
| `Accept-Encoding` | request | Обязателен `gzip`; отсутствующий/другой encoding → 415 без body. [сжатие](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-szhatie-soderzhimogo-otvetov). |
| `Content-Encoding` | response | `gzip`, если применено сжатие; сервер иногда возвращает несжатый body. Тот же раздел. |
| `Content-Type` | request/response | JSON requests используют `application/json`; отсутствие правильного типа может дать error 2001. Response определяется операцией, файлы не JSON. [E, ошибки формата](https://dev.moysklad.ru/doc/api/remap/1.2/#/errors#3-oshibki-formata). |
| `User-Agent` | request | Упомянут в общем лимите заголовков; API-specific формат и обязательное значение не установлены. [R](https://dev.moysklad.ru/doc/api/remap/1.2/#/restrictions#3-limity-i-ogranicheniya-v-json-api). Текущий SDK уже выставляет собственный User-Agent. |
| `X-Lognex-Auth` | response | Расширенный код ошибки аутентификации; отдельного enum значений общий раздел не задаёт. [обработка ошибок](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-obrabotka-oshibok). |
| `X-Lognex-Auth-Message` | response | Сообщение auth error; не использовать текст для классификации. Тот же раздел. |
| `X-Lognex-API-Version-Deprecated` | response | Дата отключения версии API; формат даты не нормирован общим разделом, сохранять raw string. Тот же раздел. |
| `X-RateLimit-Limit` | response | Количество запросов для равномерного расходования за интервал до 429. Не обещать, что значение всегда 45 либо всегда выражено в weighted units: это не уточнено новым разделом weights. [G, headers](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-obrabotka-oshibok). |
| `X-RateLimit-Remaining` | response | Число запросов до 429; такая же неоднозначность после изменения weights. Тот же раздел. |
| `X-Lognex-Retry-TimeInterval` | response | Интервал **в миллисекундах** для вышеуказанного бюджета. Тот же раздел. |
| `X-Lognex-Reset` | response | **Duration в миллисекундах**, не Unix timestamp; до сброса ограничения, `0` при отсутствии установленного ограничения. Тот же раздел. |
| `X-Lognex-Retry-After` | response | **Duration в миллисекундах** до сброса. Тот же раздел. |
| `Location` | response | Redirect target для 301/303/308; при async 202 — URL результата. [G](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-obrabotka-oshibok), [A](https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime). |
| `Content-Location` | response | При async 202 — URL статуса задачи, отличный от URL результата. [A](https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime). |
| `X-Lognex-Download-Expiration-Seconds` | request | Время жизни временного download URL в **секундах**; default 60, max 86400. Нижняя граница/поведение invalid не описаны; не выдумывать. [ссылки на файлы](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-ssylki-na-fajly). |
| `X-Lognex-Accept-Region` | request | `ru` default, `uz`, `kz`; для получения/изменения региональных полей non-Russian Account нужен этот header. [региональные headers](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-regionalnye-zagolovki). |
| `X-Lognex-Content-Region` | response | Код региона ответа. Тот же раздел. |
| `Accept-Language` | request | HTTP language negotiation; сейчас `ru`, `en`; влияет только на error messages. [язык](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-zagolovok-accept-language). |
| `X-Lognex-Remap-Beta-Feature` | request | **Один** tag, повторные headers и CSV не поддерживаются; empty/unknown tag молча игнорируется. Контракт beta может меняться. [beta](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-zagolovok-x-lognex-remap-beta-feature). Известный tag `assortmentWithoutStock` см. [assortment](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/assortment#2-assortiment); tag transport не должен ограничивать закрытым enum. |
| `X-Lognex-WebHook-DisableByPrefix` | request | URL prefix обязательно содержит scheme и полное domain; для нескольких prefix — **несколько header values**, CSV запрещён. Отменяет соответствующие webhooks в контексте запроса. [webhook](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhook#4-zagolovok-vremennogo-otklyucheniya-x-lognex-webhook-disablebyprefix-cherez-api). |
| `X-Lognex-WebHook-Disable` | request, удалён | Объявлено прекращение работы **2026-10-01**. Не предлагать как работающий typed option v1 на 2026-11-01. [официальное уведомление](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/notifications/deprecate-header-x-lognex-webhook-disable.md). |
| `X-Lognex-Accept-Timezone` | request | Для money plotseries, sales/orders plots, turnover docs: текущая дата клиента с numeric zone `+0700`, example `Wed, 16 Aug 2017 23:07:01 +0700`; операции смещаются в timezone клиента. Docs пишут «RFC 3522», что выглядит опечаткой: использовать явно заданный example-format, не timezone ID `Europe/Moscow` и не только offset. [money headers](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_money.md#L41-L46). |
| `X-Lognex-Content-Timezone` | response | Дата клиента в понимании сервера, тот же report-specific формат. Тот же источник. |
| `Lognex-Pretty-Print-JSON` | request | `true` включает formatted JSON. [workbook first steps](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/workbook/_workbook_first_steps.md#L60-L85). Иные allowed values не описаны. |
| `X-Lognex-Format-Millisecond` | request | Только example `true` в workbook; современный G и так обещает миллисекунды для response dates. Сохранить raw support, не обещать недокументированные `false`, default и применимость. [example](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/workbook/_workbook_first_steps.md#L1186-L1198). |
| `Cache-Control` | request, incidental | Есть `no-cache` в workbook examples; специальный API контракт отсутствует. [sorting example](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/workbook/_workbook_sorting.md#L80-L100). |
| `Retry-After` | response, стандарт HTTP | В MoySklad source не обещан. Если сервер/proxy прислал — поддержать standard: integer **seconds** или HTTP-date. Никогда не разбирать как миллисекунды. [RFC 9110 §10.2.3](https://www.rfc-editor.org/rfc/rfc9110.html#name-retry-after). |

Рекомендация: typed defaults + request-scoped `http.Header` с сохранением repeated values; response metadata содержит clone полного response header и typed parsed fields. Raw значения сохранять даже при parse error. Приоритет: defaults < client defaults < request overrides; body/method-specific headers задаёт transport. Не позволять incidental request options незаметно уничтожить auth. Отдельный observer может получать metadata каждой попытки; entity method должен давать доступ к metadata финальной попытки и retry count. Общие response headers `Date`/`Content-Length` и будущие headers доступны через raw map без обязательного enum.

## Общие limits и weighted limiter

Все факты этой таблицы: [R](https://dev.moysklad.ru/doc/api/remap/1.2/#/restrictions#3-limity-i-ogranicheniya-v-json-api), [вес и даты](https://dev.moysklad.ru/doc/api/remap/1.2/#/restrictions#3-limit-zaprosov-za-3-sekundnyj-period).

| Scope | Лимит |
|---|---|
| Principal (Employee или App) | Бюджет 45 units / 3 seconds; с 2026-02-16 отдельный бюджет каждого Employee и каждого App внутри Account. |
| Principal concurrency | Не более 5 одновременно выполняемых requests от одного Employee или App. |
| Account concurrency | Не более 20 одновременно выполняемых requests в сумме. |
| HTTP request metadata | Не более 8 КБ; docs явно включают URL, User-Agent, Authorization. Это шире простого `sum(Header)`; точный wire accounting, KB=1000/1024 не определены. |
| Outbound request | Не более 20 МБ отправляемых данных. Точная единица MB не раскрыта. Не путать с decoded file size. |
| Async | Не более 4 задач в очереди на Account. Не заменяет HTTP semaphore; в v1 implementation только если async включён в scope. |
| Массив объектов/positions/materials/products | Не более 1000 элементов в одном массиве request; превышение → 413. |
| JSON nesting | Max 10, error 2006. [E](https://dev.moysklad.ru/doc/api/remap/1.2/#/errors#3-oshibki-formata). |
| Batch future policy | Create/update/delete list max 1000; processing max 5 минут. После timeout результат записи неизвестен, не считать откатом. [G](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-sozdanie-i-obnovlenie-neskolkih-obuektov). |
| Автоотключение API | В течение последнего часа более 200 идентичных requests с ошибкой/мин; более 200 requests с 429/мин; более 100 PUT к одной сущности/мин. Не использовать эти пороги как рабочий request budget. [R](https://dev.moysklad.ru/doc/api/remap/1.2/#/restrictions#3-avtomaticheskoe-otklyuchenie-dostupa). |

| Auth / endpoint | Вес | Effective date / следствие |
|---|---:|---|
| Public/private App token | 1 | 45 обычных requests / 3s. |
| Draft или moderation App token | 4 | С 2026-04-13, floor(45/4)=11 / 3s. |
| Basic / Employee token | 2 | С 2026-05-12, 22 / 3s. |
| Basic / Employee token | 3 | С 2026-09-01, **текущая фаза** на дату проверки, 15 / 3s. |
| Basic / Employee token | 4 | С 2026-12-01, будущая фаза, 11 / 3s. |
| `GET /report/stock/all`, `GET /report/stock/bystore` | 5 | Независимо auth; отдельное правило endpoint, не `max(5, authWeight)` или произведение. Уведомление говорит «с февраля 2026», точный день для weight=5 не указан. [уведомление](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/notifications/rate-limits.md). **Не распространять автоматически на `/current`, `/bystore/current`, document reports и весь prefix `/report/stock`.** |

Рекомендации implementation:

1. Limiter знает principal profile, clock и endpoint classifier; Bearer string сам по себе не отличает Employee/App/Draft. Явная config нужна для корректного веса; unknown profile — консервативные 4 units с возможностью override, не считать каждый Bearer публичным App.
2. Единый weighted rolling-window/paced scheduler на principal, 45 units/3s. Если token bucket, одного bucket capacity 45 недостаточно доказать соответствие sliding-window server — допускается burst около границы. Для безопасного default равномерное планирование по весам и проверка накопленного окна.
3. Общий limiter/semaphore instance допускает sharing между clients. Account semaphore отдельно от principal semaphore; один Client не контролирует другие процессы и токены того же Account. Не обещать distributed global enforcement без внешнего shared coordinator.
4. Каждая попытка, включая retry, заново получает weighted budget и concurrency slot; не возвращать spent units после server error. Semaphore держать до чтения/закрытия response body; во время backoff slot освобождён. Context cancellation прекращает queue wait и sleep.
5. 429 + error 1073 — concurrency saturation: headers не показывают concurrency limits. Не превращать такой ответ в «нулевой remaining budget»; общий cooldown можно применить, но token bucket и concurrency — разные ограничения. [G](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-obrabotka-oshibok).
6. До clarification units в `X-RateLimit-*` не использовать Remaining как достоверное количество weighted tokens; хранить raw, соблюдать server cooldown, разрешить explicit configured limiter.

### Размеры, files и streams

- Файлы передаются как JSON с Base64 `content`; max 10 файлов/request, max 100 файлов/object. [files](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/files#2-fajly). Error 1079 задаёт max file size 20 МБ; 1080 — 100/object; 1082 — 10/request. [E](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_errors.md#L80-L83).
- Независимые ограничения: decoded file 20 МБ **и** весь encoded outbound JSON 20 МБ. Base64 увеличивает body; нельзя принимать 20-МБ file как гарантированно допустимый JSON upload. Для точного preflight нужны byte length после JSON encoding и file-specific limits.
- Images: product/bundle/variant max 10; content card max 15. [images](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_images.md#L1-L10).
- Streaming outbound, если будет добавлен: unknown length проверять bounded reader с max+1 bytes и отдельной replay factory; нельзя бесконечно читать stream и автоматически повторять unreplayable body. Сам streamed transfer не отменяет request-size limit.
- Download stream и большие JSON responses: общий предел outbound request 20 МБ не является response-size limit. Не вводить общий 20-МБ cap на успешные ответы. Raw response/stream API должен ясно передавать обязанность закрыть body; текущий `Do` сам закрывает body и для download неприменим.
- Exact KB/MB accounting остаётся неуточнённым фактом документации. В v1 допустим conservatively documented SDK cap в decimal bytes с override; не выдавать выбранные bytes за точное API правило. Field-level `String(N)`, entity-count/tariff limits и endpoint-specific image restrictions — capability schema/model validation, не универсальная transport эвристика.

## Query parameters без expand

Transport должен пропускать весь `url.Values` losslessly. Typed builders не должны обещать поддержку параметра любым endpoint. Общая grammar + endpoint-specific capabilities, raw escape hatch. `expand` намеренно вне v1; не нужен ни builder, ни скрытая догрузка.

| Query name | Контракт / scope / источник |
|---|---|
| `limit`, `offset` | Обычно `limit` 1..1000 default 1000, `offset` default 0; отсутствие значения отличается от explicit invalid 0. [product parameters](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#3-poluchit-spisok-tovarov). `audit` default 25, limit 1..100. [audit](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md#L90-L98). Не универсальные для любой response shape. |
| `filter` | Grammar ниже; supported fields/operators определяет endpoint и mode. [G filter](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-filtraciya-vyborki-s-pomoshyu-parametra-filter). |
| `order` | `field[,asc|desc]` через `;`; default asc. Типы numeric/string/datetime/bool/UUID; конкретные поля в [sorting matrix](https://dev.moysklad.ru/doc/api/remap/1.2/#/workbook/workbook-sorting#2-sortirovka). Ссылочные поля/arrays не универсально сортируемы. |
| `search` | Contextual search, URL-encoded string, endpoint-specific поля, результаты по релевантности; это самостоятельный query, не всегда `filter=search=...`. [G search](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-kontekstnyj-poisk). |
| `fields` | **Одно** значение, не CSV. Документированные values: `stock` (document positions), `declaration` (demand/loss/retaildemand positions), `minimumStock` (product/variant), `downloadPermanentHref` (images, платный тариф). [G fields](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-chto-takoe-fields), [product](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#5-nesnizhaemyj-ostatok), [images](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/images#4-poluchit-postoyannuyu-ssylku-na-izobrazhenie-tovara-komplekta-modifikacii-ili-kartochki-kontenta). Unsupported/multiple → error 3042. G сообщает size≤100 и связывает игнорирование с указанным expand; эту формулировку не расширять до недоказанных правил всех direct `/positions` endpoints. |
| `namedfilter` | href saved filter данного entity type; server applies saved UI filter. Конкретные ограничения и сочетание с `filter` в [named-filter](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/named-filter#3-primenenie-sohranennogo-filtra). |
| `owner` | href Employee для списка saved filters другого Employee; требуется admin. Не универсальный query для обычных entity lists. [named-filter owner](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_named_filter.md#L90-L100). |
| `groupBy` | Legacy assortment: `product`, `variant` default, `consignment`; новый beta assortment имеет иной default (`product`). Stock reports: `product`, `variant` default, `consignment`; turnover: `product` default, `variant`. [legacy](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/assortment-legacy#2-assortiment-deprecated), [new assortment](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/assortment#2-assortiment), [stock](https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-stock#3-poluchit-rasshirennyj-otchet-ob-ostatkah), [turnover](https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-turnover#2-otchet-oboroty). |
| `includeRelated` | Extended stock report, boolean; требует product/variant filter; true включает связанные variant/consignment, принудительно `groupBy=consignment`, явный groupBy игнорируется. [stock parameters](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_stock.md#L182-L191). |
| `include` | Краткие current stock reports: `zeroLines` включает нулевые остатки; никогда не участвовавшие в операциях товары всё равно исключены. [stock include](https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-stock#4-parametr-include). |
| `changedSince` | Краткие current stock: datetime selection изменений до now, возвращает текущие quantities, не delta; автоматически включает zero lines; не future. Docs одновременно пишут max lookback 24h и отсутствие ограничения minimum — противоречие, до уточнения не принимать как произвольно глубокую history. [stock changedSince](https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-stock#4-parametr-changedsince). |
| `stockType` | Краткие stock: одно из `stock` default, `freeStock`, `quantity`, `reserve`, `inTransit`. [stock type](https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-stock#4-parametr-stocktype). В старом workbook ассортимент/stock параметры встречаются иначе — не переносить автоматически. |
| `withRecalculate` | Только `/report/stock/bystore/current`: `true` включает пересчитываемые позиции со значением null. [stock recalc](https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-stock#4-parametr-withrecalculate). |
| `momentFrom`, `momentTo` | Report range; обязательность определяется report: money/sales plots required, некоторые pnl/turnover optional. Это query names, в отличие от `filter=moment=...`. [money](https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-money#3-dvizhenie-denezhnyh-sredstv), [turnover](https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-turnover#2-otchet-oboroty). |
| `interval` | Money/sales/orders plotseries: required `hour`, `day`, `month`. [sales/orders](https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-sales-orders#2-pokazateli-prodazh-i-zakazov). |
| `async` | Boolean false default; true создаёт async task, запрещает offset/limit, поддержан только явно перечисленными endpoints. До решения scope v1 — future capability, не превращать любой GET в async. [A](https://dev.moysklad.ru/doc/api/remap/1.2/#/async#3-vypolnenie-zaprosa-v-asinhronnom-rezhime). |
| `action` | Wizard autofill: CSV `evaluate_price`, `evaluate_discount`, `evaluate_vat`, `evaluate_cost`; body prerequisites зависят от action/document; не generic create/update option. [autofill](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/autofill#2-avtozapolnenie). |
| `miniature` | `true` в download URLs image miniatures; следует сохранять server URL. Не вводить как option любого entity GET. [files](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_files.md#L121-L126). |

В официальных параметрических таблицах дополнительно встречаются `updatedFrom`, `updatedTo`, `updatedBy` для task, `tags` для counterparty, `incomingDate` для commissionreportin: описания прямо отправляют к grammar `filter`. Не объявлять их подтверждёнными отдельными top-level query parameters только по имени строки таблицы. [task table](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_task.md#L827-L834), [counterparty table](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_counterparty.md#L950-L954), [commissionreportin](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_commissionreportin.md#L124).

Не query options SDK: `id`, `type`, `entityType`, `positionID`, `idFile`, `idImage`, `accountId`, `retailStoreId` и аналогичные placeholders — части path; `requestId` — callback URL webhook example, `raw` — Postman documentation links, `filename` в async workbook — filename команды, `organization/project/store` в одном ошибочном meta.href example plots не отменяют нормативное размещение в `filter`; `limi` — typo workbook. Источник: полная сверка URL examples и таблиц pinned `md/`.

### Filter grammar и escaping

Все общие правила: [G filter](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-filtraciya-vyborki-s-pomoshyu-parametra-filter), [LIKE](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-operator-filtracii-podobie), [custom attributes](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#5-dostupnye-operatory-dlya-filtracii-dop-polej).

- Form: `field operator value`, conditions joined `;`; escape буквального `;` в textual value — **два символа backslash и semicolon** (`\;`), затем URL encoding всего query через `url.Values`. Не подставлять SQL/JSON escaping и не double URL-encode.
- Operators: `=`, `!=`, `<`, `<=`, `>`, `>=`, `~` (contains), `~=` (starts with), `=~` (ends with). Report-specific grammar может дополнительно описывать `!~`; это не повод включать `!~` в универсальную grammar без capability.
- Повтор `=` одного поля — OR; повтор `!=` — AND; разные поля — AND. `=` совместно с неравенством того же поля запрещён. Для нескольких одинаковых сравнений одного поля сейчас учитывается первое, docs предупреждают о будущем error — builder должен отклонять ambiguous repetition.
- Empty `field=` — отсутствует/null, для strings также empty string; `field!=` — присутствует. Reference filters передают full href; ID filters передают UUID; custom attribute filters используют URL самого attribute как key.
- Документы default исключают корзину: `filter=isDeleted=true/false`, оба значения — всё. Сущности с archived default исключают архив; `archived=true;archived=false` включает оба. Это **filter fields**, не самостоятельные `?isDeleted=`/`?archived=`.
- Новый beta assortment (`assortmentWithoutStock`) не поддерживает часть legacy filters, в частности `filter=search`, `stockMode`, `stockMoment`, `stockStore`, `quantityMode`; поддерживает отдельный `search`. Barcode `!=` запрещён. [beta differences](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/assortment#3-otlichiya-ot-starogo-endpointa).
- Не обещать generic boolean AST с произвольными parentheses/AND/OR: API grammar этого не задаёт. Отдельные report filters имеют собственную multi-value/combination grammar.
- Backslash escaping собственного backslash и всех operator-like bytes в values документация не определяет исчерпывающе. Builder гарантирует documented `;` case, сохраняет raw escape hatch и не выдумывает escaping protocol.

## Retry/backoff: безопасная политика

Рекомендация, не заявленный алгоритм сервера:

- Default automatic retries только для read operations с известной семантикой: обычный GET/HEAD, **исключая** async task creation даже при GET. Write POST/PUT и будущий DELETE/batch требуют operation-specific opt-in. HTTP PUT формально idempotent по [RFC 9110 §9.2.2](https://www.rfc-editor.org/rfc/rfc9110.html#name-idempotent-methods), но API update может создавать side effects, менять files/positions, отправлять webhooks и затрагивать параллельное состояние. Без доказательства endpoint semantics безопаснее default no-write-retry.
- Retry candidates: 429, transient 500/502/503/504, временные transport failures. Нельзя автоматически повторять canceled/deadline context, invalid request, auth/permission/validation errors, JSON decode errors после успешного 2xx, unexpected content encoding. 1088 timeout batch и error 1065 concurrent update не превращать в unconditional retry.
- `syncId` документирован как deduplication create только для entities, **у которых есть это поле**. Повтор того же create и syncId возвращает прежнюю entity; не общий Idempotency-Key для любого POST. Idempotency-Key header в source не найден. [G syncId](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-naznachenie-polya-syncid).
- `POST security/token`: никогда automatic retry; успешная попытка отзывает предыдущие токены, а response может потеряться. Ошибка доставки не доказывает, что token не создан.
- Bounded retries: например max 3 attempts total, initial 500ms, cap 30s, full jitter для exponential component. Эти числа — выбранные SDK defaults, не API limits. Delay = max(valid server hints, jittered exponential); max hints: `X-Lognex-Retry-After` ms, `X-Lognex-Reset` ms, стандартный `Retry-After` seconds/date. Не обрезать server minimum cap 30s; если ожидание превышает context/deadline/budget — завершить без новой попытки. Невалидные/negative/overflow hints не превращать в бесконечное ожидание.
- Retry idempotence задаётся metadata операции, а не `body != nil` или global `RetryAll`. Replayable bytes — необходимое, но недостаточное условие. `NewRequest` сейчас создаёт bytes.Reader, `Request.GetBody` доступен; каждую попытку clone request + fresh GetBody + clone headers. Custom requests/streams без replay factory не повторять.
- Read/close промежуточный response до новой попытки; semaphore release до sleep. Каждая попытка проходит limiter заново. После исчерпания вернуть последнюю status/header/error и attempts count; `errors.As(*APIError)` сохраняется.
- Go HTTP Transport сам имеет ограниченные retry rules; не маркировать writes как idempotent через произвольный `Idempotency-Key` ради включения implicit retries. [net/http Transport](https://pkg.go.dev/net/http#Transport).

## POST security/token

Проверенный endpoint: `POST /api/remap/1.2/security/token`, Basic authentication, без JSON body в официальном example; required gzip. Успех **200**, object `{"access_token":"..."}`. Все ранее созданные токены Employee отзываются. [G token](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#4-poluchenie-novogo-tokena).

Рекомендация v1: отдельный `security.Service.CreateToken(ctx, credentials, options)` или `GenerateToken`, явные Basic credentials per call. Return typed `{AccessToken string}` + response metadata; не возвращать anonymous string без возможности видеть response. Не менять автоматически Auth существующего client, не сохранять token, не логировать response body. Reject malformed successful response (missing/empty token). Bearer-configured client не должен случайно отправить Bearer на token generation вместо требуемого Basic.

## Error constants и совместимость

- Envelope `errors: []Detail`; fields: `error`, `parameter`, integer `code`, `error_message`, `moreInfo`, `line`, `column`, `dependencies`, `meta`. `code` может отсутствовать: тогда классифицировать HTTP status. [G errors](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-obrabotka-oshibok).
- Полный pinned `_errors.md` содержит **545 unique numeric codes**; номера с пропусками, категории не фиксированные continuous ranges. Каталог охватывает и отсутствующие в v1 operations/entities. Для требования «все code constants» генерировать из всех `id="error_N"`, а не выбирать только 10xx или реализованные entities. [E](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_errors.md).
- Stable canonical names: `Code1000`, `Code1001`, … для **существующих** номеров. Официальных English symbol names нет; translation-derived names субъективны, могут конфликтовать и привязывать API к меняющемуся message. Optional curated semantic aliases допустимы с canonical numeric symbols (`CodeConcurrentRequests = Code1073`), без необходимости naming всех 545 вручную.
- Сохранить `Detail.Code int`, добавить **untyped integer constants**: source compatibility существующих struct literals/comparisons остаётся. Новый distinct `type Code int` в existing field ломает assignments; это отдельное breaking API decision. Unknown future integers не отвергать при decoding.
- Helpers `HasCode(err, code)` / `APIError.HasCode(code)` проверяют все details и wrapped errors; HTTP status и domain code раздельны. Header auth code не тот же namespace, что `Detail.Code`. Не вводить universal mapping code→HTTP status: такого контракта нет.
- Acceptance: generated code set равен всем 545 IDs pinned source (no omissions/duplicates), неизвестный code decoding работает, empty/malformed error body сохраняет HTTP status; классификация не зависит от locale и message. Из особо значимых: 1062 Accept, 1065 parallel update, 1073 concurrency, 1074 maintenance, 1079 file size, 1088 batch timeout, 1095 temporary App unavailable, 2006 depth, 2007 array length, 3042 fields.

## List, walker, GetPositions и facade

Это архитектурные рекомендации по inspected SDK, не отдельный API контракт:

- Сейчас transport `NewRequest(ctx,method,endpoint,url.Values,body)` + `Do(req,dst)(*http.Response,error)` уже поддерживает query/header access на низком уровне; entity `do` теряет metadata и передаёт nil query. `query.go/filter.go/order.go/pagination.go` — stubs. `Do` закрывает response body; на уровне entity его нельзя трактовать как usable download stream.
- Retrofit обязателен для всех 17 существующих packages: `cashin`, `cashout`, `counterparty`, `country`, `demand`, `expenseitem`, `invoicein`, `invoiceout`, `move`, `organization`, `paymentin`, `paymentout`, `product`, `purchaseorder`, `purchasereturn`, `supply`, `uom`; и каждого нового сервиса в v1. Headers/options/metadata нужны Get/Create/Update/List, не только новым методам. Пока API не выпущен, выбрать единый result contract один раз; предпочтительно typed `Result[T]{Value T, Response Metadata}` либо typed result pointer + metadata — не context mutable pointer side channel.
- Root facade собирает typed handles и один shared transport, limiter, retry policy. Например `client.Entity.Product`/`client.Reports.Stock`/`client.Context.Employee`/`client.Security`; точные имена синхронизировать с существующим design. Handles стабильны, не создают HTTP traffic при получении, не создают каждый свой limiter. Не global service registry по string с unsafe assertions. Нижние `NewService(*transport.Client)` остаются доступны.
- `List` читает **одну страницу**; принимает endpoint-specific options; возвращает `model.Collection[T]` и metadata. Не осуществляет скрытого Walk/expand/GetPositions.
- Walker/iterator — отдельная явная операция, sequential по умолчанию. Проверяет context перед каждой страницей/callback, прекращает обход при callback error, не накапливает всё. Переходит по `meta.nextHref` через существующий same-origin URL validator; raw query не double-encode и не подмешивать исходный offset поверх nextHref. Защита repeated/cyclic nextHref, invalid URL, отсутствие progress. Если server nextHref отсутствует — завершить; не выдумывать новое окно только по `len(rows)==limit`.
- `meta.size` — размер выборки на стороне server, не snapshot guarantee; concurrent mutations могут дать gaps/duplicates. Не обещать exactly-once; deterministic endpoint-supported `order` помогает, но не создаёт snapshot. Async docs явно предупреждают о duplicates/gaps при concurrent changes. [A](https://dev.moysklad.ru/doc/api/remap/1.2/#/async#2-asinhronnyj-obmen).
- `GetPositions` уже существует у document services и ходит непосредственно `/entity/{type}/{id}/positions` через `positions.Service[T].GetPage`. Сохранить имя, direct request и typed Position, добавить общий options/metadata и отдельный positions walker; не заменить на Get+expand и не переименовывать без причины. `positions.PageOptions` zero/default сейчас означает 1000; согласовать с общим PageOptions без второго несовместимого контракта.
- Нельзя навязать `Collection[T]` всем lists/reports: стандартные entity/positions — `{context?,meta,rows}`; current stock — **голый array**; money plots — object с `series`, `credit/debit`; dashboard — object `sales/orders/money`; context/employee и context/company — одиночный object. Только paginated collection имеет WalkPages. [current stock](https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-stock#3-kratkij-otchet-ob-ostatkah), [money](https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/report-money#3-dvizhenie-denezhnyh-sredstv), [dashboard](https://dev.moysklad.ru/doc/api/remap/1.2/#/reports/dashboard#2-pokazateli), [context](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-kontekst-zaprosa-sotrudnika).

## Нерешённые факты и решения до coding

1. Scope Delete/batch/async ещё требует решения пользователя. Выше их limits/errors/retry traps зафиксированы как future policy; это не авторизация реализовать эти методы.
2. Выбрать public response contract, options precedence и facade naming до массового retrofit — эти решения затрагивают весь SDK.
3. Уточнить в API/поддержке единицы `X-RateLimit-*` после weighted limits, точные bytes MB/KB, статус legacy `X-Lognex-Format-Millisecond`, contradictory changedSince lookback, minimum download expiry. До этого сохранять raw headers/query и маркировать conservative SDK policy отдельно от server contract.
4. Согласовать defaults unknown auth profile и automatic retries writes. Рекомендация: unknown=4 units, no-write-auto-retry; явный config для public/private App и checked idempotent operations.
5. Полная per-entity filter/order/search capability matrix должна жить с каждым endpoint ticket; общие transport builders реализуют grammar. Наличие поля в модели или параметра в общем G не доказывает, что конкретный endpoint это поддерживает.
