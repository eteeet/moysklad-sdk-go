Status: needs-triage
Stage: DRAFT

# SDK v1: согласуемый профиль первой версии

## Problem Statement

В dev реализованы 6 справочников и 11 документов, но отсутствуют общие List/query,
контроль лимитов, Audit, нужные товарные/розничные сущности и единый клиент.
Цель — закончить выбранный профиль SDK несколькими крупными захватками, сохраняя
самостоятельные проверяемые и откатываемые tickets. Это не обещание покрыть весь JSON API.

## Solution и результат планирования

7 основных захваток, 49 implementation tickets; дополнительно 1 условная spec
с 3 tickets Delete/batch/async. [Карта и зависимости](map.md),
[покрытие исходных требований](coverage.md), [каталог контрактов](contracts.md).

Все документы имеют статус DRAFT/needs-triage; optional tickets — needs-info.
Создание spec/tickets авторизовано пользователем. Реализация, публикация,
коммиты planning, live API, credentials и cleanup worktree не авторизованы этой задачей.

## Baseline и источники

- Код: dev, проверенная отправная точка `acbe294`; 14 коммитов предыдущей захватки
  уже созданы, общие JSON test helpers вынесены в internal/documenttest.
- Нормативный source: api-remap-1.2-doc, revision
  `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`. Research сверил official HEAD с этим
  revision; runtime implementation перед стартом повторяет проверку deltas.
- [Transport research](research/transport.md),
  [инвентаризация сущностей](research/entities.md),
  [Audit/context/metadata/files/reports](research/audit-metadata-files-reports.md).
- [status-entities.md](../../docs/status-entities.md) и
  [status-technical.md](../../docs/status-technical.md) — входные списки, не готовая
  матрица обязательных методов. У PriceType, context, metadata, Audit и reports
  нет одинакового Get/Create/Update. Legend статусов уточняется R04.

## Пользовательский scope

Audit; Bundle/Variant/Service; Employee/Role; context/employee; BonusTransaction/
BonusProgram; Assortment; Currency; Files; PriceType; RetailStore; States; TaxRate;
CustomerOrder; Enter/Loss; RetailDemand/RetailSalesReturn/RetailShift/SalesReturn;
Turnover All/ByStore/ByOperations; POST security/token (одно требование, не два);
entity/metadata; коды ошибок; syncId; все headers/limits; query кроме expand;
List; context/application; фасад SDK.

### Рекомендуемая граница операций — пока предположение

Get/List/Create/Update там, где они действительно поддержаны и соответствуют ресурсу.
GetPositions — у товарных документов; files upload/download, PriceType ReplaceAll,
Employee security и metadata states — явно описанные специальные методы, не generic CRUD.

Вопрос о включении Delete/batch/async задан пользователю; явного ответа нет.
Принято рабочее предположение: они не входят в основной v1. Для них всё равно
созданы X01–X03 и отдельная spec. Если включить — readiness R01/R04 и critical path
пересчитываются; код optional не запускается просто по готовности зависимостей.

«Все лимиты» означает поддержку применимых ограничений реализованных операций.
Квота очереди async и batch processing timeout имеют реализацию в X02/X03;
до включения этих tickets их можно документировать, но нельзя отметить как
реализованные. Полная поддержка всех query ресурсов за пределами выбранного набора
также не объявляется готовыми typed services только благодаря raw url.Values.

## Architecture: interface и seam

- Одна shared transport policy обрабатывает options, headers, response metadata,
  weighted budget, concurrency, cancellation и bounded retries. Не вводить общий
  entity CRUD, который проверяет ID/200/Collection для всех ответов.
- Существующие standalone services сохраняются; их constructors принимают transport.
  Фасад собирает typed handles и не делает HTTP. NewClient не выпускает токен,
  не определяет App/Employee/Draft по Bearer строке и не читает context автоматически.
- T01 фиксирует публичный interface до массового retrofit. Рекомендация:
  optional per-call options и request-owned ResponseInfo при сохранении результатов
  Get/Create/Update. Возможное изменение method-value type при variadic options
  проверяется отдельно; «старые вызовы компилируются» не равно полной совместимости.
- Никакого LastResponse/global context/global singleton. Response metadata доступна
  при ошибках, raw headers clone сохраняет repeated/unknown/malformed values.
- JSON response, plain array, report envelope и binary stream — разные контракты.
  Stream ownership и lifetime semaphore фиксирует F01; обычный Do закрывает body.
- DateTime по умолчанию Москва. Report timezone header содержит дату с numeric zone,
  а не произвольный timezone ID; он не переключает глобальный decoder между запросами.

## Предметные ограничения

- Reference без expand; неизвестные MetaType/enum не отвергаются. Ресурс по ссылке
  не становится технической зависимостью, если загружается только Reference.
- Одна Position/PositionRequest для совместимых строк; новое поле добавляется по
  источнику. Смена типа/смысла JSON-поля — основание отдельной формы, не отсутствие
  части полей. Bundle components и report rows не подгоняются под Position.
- Readonly marker таблицы сильнее противоречащего write example; исключения утверждает
  пользователь и записывает contracts.md (bundle components quantity). Отсутствие marker
  само по себе не доказывает writable. Полная Write matrix C/U/null/required/immutable
  перед code обязательна; каталог полей не заменяет эту работу.
- Деньги — float64 в копейках; баллы, quantity, процент и курс отдельно.
  Raw Attribute/Audit/Permissions JSON сохраняет long >2^53 без map[string]any.
- syncId только у поддерживающих сущностей; server write-once, не глобальный
  Idempotency-Key. Автоматические retries записи не включаются по одному syncId.
- Бизнес-проверки и permissions выполняет API; SDK не делает pre-GET ради проверки
  склада, основания, смены, тарифов, остатков или числа уже загруженных файлов.

## Assortment: зачем отдельный ресурс

Это единый смешанный List/search по Product/Service/Bundle/Variant и Consignment.
Он нужен для общего поиска и одной пагинации без четырёх отдельных запросов,
клиентской склейки результатов и попытки воспроизвести общую релевантность.
Сами четыре entity packages нужны для записи и специализированного чтения.

C08 сохраняет known typed rows + fallback unknown meta.type/raw JSON. Legacy и
облегчённый beta assortmentWithoutStock имеют разные query/stock/default groupBy;
mode фиксируется явно. Beta header не включается молча, остатки автоматически
отдельным stock report не подгружаются. Полный Consignment CRUD не входит в scope.

## Limits и retries

- 45 weighted units/3s; principal concurrency ≤5, account ≤20. Текущий Employee
  weight3, с 2026-12-01 weight4; Draft4, App1. Точные два stock endpoints weight5.
  Профиль задаётся явно, clock/weights конфигурируемы и тестируются на границе даты.
- Один процесс не контролирует другие процессы/токены Account. Shared policy может
  координировать clients одного приложения; для distributed enforcement нужна
  внешняя shared координация. Не обещать её встроенным локальным semaphore.
- Request budget учитывает encoded body/Base64 и URL+headers/auth, response download
  не получает выдуманный общий 20-MB cap. Неопределённые KB/MB единицы отмечены как
  source uncertainty, SDK cap называется своим именем и конфигурируется.
- Default retry только подтверждённых reads. security/token и async submission
  side effects не retry даже при GET для async. Write retries требуют отдельного
  доказанного contract и явного opt-in; cancellation/decode errors не повторяются.
- X-Lognex delays — milliseconds; standard Retry-After — seconds/HTTP-date.
  Невалидные hints не дают unbounded waits; server minimum не обрезается своим cap.

## Testing Decisions и процесс завершения

JSON — v1/v2, независимые golden, upstream SHA-256 provenance, sparse/null/unknown,
readonly и required разные для Create/Update; public HTTP seam — transport+httptest.
Для limiter/retry — детерминированное время, cancellation и race; не sleep-based tests.
Для file stream — close/cancellation/redirect/auth/partial output и limiter lifetime.
Для reports/contexts — own success criteria, без generic обязательного ID.

Каждый ticket — отдельный code commit; общая dependency раньше dependent.
Каждая захватка проходит fixed-BASE независимые Standards/Spec reviews и integration
race/vet/gofmt/diff-check. При разделении накопленных изменений проверяется именно
снимок каждого коммита, а не только финальное рабочее дерево. Docs — отдельный
последний commit волны, planning — только по отдельной команде.

## Definition of Done v1

- Все 49 основных tickets resolved, либо user-approved исключения явно изменили scope.
- Таблица coverage для каждого исходного требования указывает реально реализованный
  метод, ограничения и доказательство проверки; optional gates не скрыты чекбоксами.
- Фасад и standalone entry points согласованы; examples компилируются; clean checkout
  проходит CI/race/vet; публичные docs/CHANGELOG отражают scope и изменения interface.
- Нет незаявленных live checks, выдачи токенов, автоматической фискализации или
  бонусного начисления. Runtime secrets не сохраняются в fixtures/logs.
- Release artifact готов к выпуску. Реальный tag/publish, live API и worktree cleanup
  требуют отдельной команды; их отсутствие не подменяется выдуманной отметкой PASS.

## Out of Scope

Expand — последним и только по отдельной команде; полный Vendor API, остальные
не перечисленные сущности/отчёты; реплика/outbox/sync engine на базе Audit;
fiscal engine, расчёт лояльности, автоматически подключаемые внешние сервисы.
Images/cashier lifecycle/employee password activation и metadata definition mutations
не следуют автоматически из поддержки Files/Employee/Metadata.

## Open Decisions до реализации

1. Включаются ли X01–X03 Delete/batch/async? Пока рекомендованный основной scope — без них.
2. T01: per-call options/ResponseInfo и совместимость public method types; R01: naming фасада.
3. C08: default legacy mode или явно выбранный beta; никакого автоматического переключения.
4. Источники: Role.permissions Array vs Object; Turnover.ByOperations даты required
   в таблице vs optional defaults в тексте. Разрешить в ticket freeze, не выдумывать
   серверное поведение. В примерах Turnover рекомендован explicit period.
