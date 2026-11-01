# Статус технических возможностей SDK

Источник: `_general.md`, `_restrictions.md`, `_async.md`, `_errors.md`
(revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`), [coverage.md](coverage.md).
Сущности — [status-entities.md](status-entities.md).

## Транспорт

- [x] Config: BaseURL, HTTPClient, Timeout
- [x] Auth: Bearer token, Basic
- [x] Получение токена: POST security/token (без повторов, Basic per call)
- [x] Заголовки Accept, Content-Type, User-Agent
- [x] gzip: Accept-Encoding и распаковка ответа
- [x] Сброс Authorization при redirect на чужой origin
- [x] APIError: status, errors[], code; errors.Is/As; префикс `moysklad:` один раз
- [x] Коды ошибок `_errors.md` как константы (545), HasCode/StatusCode
- [x] Ответ 413: APIError с деталями, без разбиения и повтора

## Общие типы и JSON

- [x] encoding/json/v2 и совместимость с v1
- [x] Field[T]: unset/null/value, omitzero
- [x] Meta, Reference, Collection[T], ResponseContext
- [x] DateTime: формат МойСклад, Москва, RFC 3339
- [x] Деньги: float64 в копейках, vat — int
- [x] Attribute: значения доп. полей, сброс через null
- [x] DocumentRate, DocumentOverhead, FinanceOperation
- [x] positions: Position/PositionRequest для всех товарных документов
- [x] syncId: Create и Update (Counterparty, Organization, документы)
- [x] syncId в Product: чтение и Create (в Update read-only по API)
- [ ] Раскрытые объекты вместо Reference (после expand)

## Заголовки

- [x] X-RateLimit-Limit, X-RateLimit-Remaining (разбор; единицы не уточнены)
- [x] X-Lognex-Retry-TimeInterval, X-Lognex-Reset, X-Lognex-Retry-After (мс)
- [x] X-Lognex-Accept-Timezone, X-Lognex-Content-Timezone
- [x] X-Lognex-Accept-Region, X-Lognex-Content-Region
- [x] Accept-Language
- [x] X-Lognex-WebHook-DisableByPrefix (старый X-Lognex-WebHook-Disable отключён API)
- [x] X-Lognex-Remap-Beta-Feature
- [x] X-Lognex-API-Version-Deprecated (сырое значение в ResponseInfo.Header)
- [x] X-Lognex-Download-Expiration-Seconds
- [x] X-Lognex-Auth, X-Lognex-Auth-Message (в ResponseInfo, в том числе при APIError)

## Лимиты

- [x] 429/5xx: RetryPolicy только для GET/HEAD, подсказки сервера, отмена через context
- [x] Rate limiter: 45 единиц за 3 с; вес 2→3→4 для логина/токена пользователя; 4 для решения-черновика; 5 для stock/all и stock/bystore
- [x] Параллельность: ≤5 от пользователя/решения, ≤20 от аккаунта (в пределах процесса)
- [x] Размер: ≤8 КБ заголовки, ≤20 МБ тело, вложенность 10 — preflight с выбором единиц SDK
- [x] ≤1000 позиций в embedded positions: проверка до HTTP
- [ ] Массовые операции: ≤1000 объектов, таймаут 5 мин
- [ ] Async: ≤4 задачи в очереди на аккаунт

## Query-параметры

- [x] limit/offset в GetPositions: 0→1000, 1…1000, offset≥0
- [x] limit/offset в List
- [x] filter: операторы, ссылки, ID, доп. поля, «подобие»
- [x] order
- [x] search: контекстный поиск
- [x] namedfilter: применение сохранённого фильтра
- [x] fields: одно значение, только где документирован (Product: minimumStock); fields=stock у positions — D01
- [x] Параметры отчётов: momentFrom/momentTo, groupBy, interval (query.ReportOptions)
- [x] Обход страниц: pages.All по nextHref
- [x] Нет скрытых запросов по href/nextHref
- [ ] expand: limit≤100, глубина ≤3, в Create/Update — **последним, только по команде**

## Методы ресурсов

- [x] Get/Create/Update: 6 справочников, 11 документов
- [x] GetPositions: generic GetPage и typed в 7 документах
- [x] List: 6 справочников, 11 документов
- [ ] Delete и массовое удаление (POST …/delete)
- [ ] Удаление в корзину
- [ ] Массовое создание/обновление (POST массивом)
- [ ] Отдельные операции позиций: Get, Create, Update, Delete
- [ ] Позиции: pack, things, slot, trackingCodes, gtd, остатки (fields=stock)
- [ ] Шаблон документа: PUT entity/{type}/new
- [ ] Автозаполнение: wizard/demand, wizard/salesreturn
- [ ] Печать (export), публикации, шаблоны печатных форм
- [ ] Файлы и изображения: загрузка base64, скачивание
- [ ] Metadata: entity/{type}/metadata, states, customtemplate/embeddedtemplate
- [ ] Доп. поля: CRUD описаний, типы файл/справочник
- [ ] Async: async=true, статус, результат, отмена
- [ ] Серверные решения: context/application

## Качество и выпуск

- [x] JSON-тесты обоими codecs, golden, upstream с SHA-256 provenance
- [x] HTTP-тесты: httptest и настоящий transport
- [x] go vet, `go test -race`
- [ ] Live read-only проверка на тестовом аккаунте (только с согласия)
- [ ] Примеры для документов (сейчас только product)
- [ ] Фасад SDK: единая точка входа
- [ ] CI, проверки безопасности
- [ ] Версионирование, CHANGELOG, релиз
- [ ] Vendor API: отдельный контур

## Открытые замечания последнего ревью

- [x] Закоммитить 11 документов по тикетам
- [ ] Дубли codecs_test.go → internal/documenttest
- [ ] PositionRequest: Cost/InTransit после VATEnabled, комментарий Cost в model.go
- [ ] Docs: условие валюты paymentin/paymentout, лишние write-формы, TODO purchaseorder
- [ ] Docs: правило «пометка read-only важнее примеров»
- [x] Удалить устаревшие worktree core-* (stash -u, branches сохранены)
