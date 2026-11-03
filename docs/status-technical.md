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
- [x] Массовые операции: ≤1000 объектов до HTTP, таймаут 5 мин — неизвестный исход, без повторов
- [x] Async: ≤4 задачи в очереди на аккаунт — документировано, считает сервер

## Query-параметры

- [x] limit/offset в GetPositions: 0→1000, 1…1000, offset≥0
- [x] limit/offset в List
- [x] filter: операторы, ссылки, ID, доп. поля, «подобие»
- [x] order
- [x] search: контекстный поиск
- [x] namedfilter: применение сохранённого фильтра
- [x] fields: одно значение; в List справочников и документов пока не принимается
- [x] fields=minimumStock у Product и Variant (model.MinimumStock)
- [x] fields=stock и fields=declaration у GetPositions — по документам источника
- [x] Параметры отчётов: momentFrom/momentTo, groupBy, interval (query.ReportOptions)
- [x] Обход страниц: pages.All по nextHref
- [x] Нет скрытых запросов по href/nextHref
- [ ] expand: limit≤100, глубина ≤3, в Create/Update — **последним, только по команде**

## Методы ресурсов

- [x] Get/Create/Update: 16 справочников, 18 документов
- [x] GetPositions: generic GetPage и typed в 13 документах
- [x] List: 20 справочников со страницами, assortment, pricetype и характеристики (массив целиком), 22 документа
- [x] Beta-режим: assortment List и All с assortmentWithoutStock на каждой странице, ListLegacy без него
- [x] Вложенные коллекции справочников: зоны и ячейки склада, счета, контакты и события контрагента,
  счета юрлица, компоненты комплекта (с мутациями), кассиры, элементы пользовательского справочника
- [x] Delete и массовое удаление (POST …/delete); у ассортимента — по meta разных типов
- [x] Удаление в корзину: Trash у 37 документов ([documents-extra](documents-extra.md))
- [x] Массовое создание/обновление (POST массивом): CreateMany/UpdateMany по разделам источника
- [x] Отдельные операции позиций: Get, Create, Update, Delete, массовое удаление
- [x] Позиции: pack, things, slot, trackingCodes, gtd, declaration, остатки (fields=stock)
- [x] Шаблон документа: PUT entity/{type}/new — New у 24 документов
- [x] Автозаполнение: POST wizard/{тип} — пакет wizard
- [x] Печать (export без следования 303), публикации, шаблоны печатных форм ([print](print.md), [publication](publication.md))
- [x] Файлы: список, загрузка base64, скачивание по downloadHref; изображения — пакет images
- [x] Metadata: entity/{type}/metadata, states; шаблоны customtemplate/embeddedtemplate — пакет print
- [x] Доп. поля: CRUD описаний и массовые операции, тип справочник (customEntityMeta)
- [x] Настройки справочников, компании и пользователя (GET/частичный PUT)
- [x] Async: ListAsync counterparty/assortment, *Async отчётов, статус, список, ожидание, результат, отмена
- [x] Серверные решения: context/application

## Качество и выпуск

- [x] JSON-тесты обоими codecs, golden, upstream с SHA-256 provenance
- [x] HTTP-тесты: httptest и настоящий transport
- [x] go vet, `go test -race`
- [ ] Live read-only проверка на тестовом аккаунте (только с согласия)
- [x] Примеры сценариев на фасаде: контекст, ассортимент, позиции заказа, аудит, обороты, файлы
- [x] Фасад SDK: `moysklad.NewClient`, группы Entity/Report/Context/Audit/Security
- [x] `make check`: gofmt, vet, race-тесты, SHA-256 фикстур, статус сущностей, diff --check
- [ ] CI workflow (для v1 не делается: только `make check`), проверки безопасности
- [x] CHANGELOG: предложение v1.0.0 и совместимость с 7442803
- [ ] Тег и релиз — отдельной командой
- [ ] Vendor API: отдельный контур

## Замечания ревью до sdk-v1

Закрыты коммитами dbe1e01…acbe294 (см. `.scratch/sdk-v1/coverage.md`).

- [x] Закоммитить 11 документов по тикетам
- [x] Дубли codecs_test.go → internal/documenttest
- [x] PositionRequest: Cost/InTransit после VATEnabled, комментарий Cost в model.go
- [x] Docs: условие валюты paymentin/paymentout, лишние write-формы, TODO purchaseorder
- [x] Docs: правило «пометка read-only важнее примеров»
- [x] Удалить устаревшие worktree core-* (stash -u, branches сохранены)
