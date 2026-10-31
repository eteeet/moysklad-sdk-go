Status: needs-triage
Stage: DRAFT

F01 (самостоятельные FinanceOperation/FinanceOperationRequest) утверждён отдельно
в первой волне 2026-11-01. Ниже DRAFT entity-контрактов F02–F05; их реализация не разрешена.

# DRAFT — PaymentIn, PaymentOut, CashIn, CashOut: Get, Create, Update

## Problem Statement

Разработчик решения не может получить, создать или частично изменить четыре финансовых
документа через SDK. Нужны PaymentIn (входящий платеж), PaymentOut (исходящий платеж),
CashIn (приходный ордер) и CashOut (расходный ордер), включая их суммы, валюту и связи
с оплачиваемыми операциями. InvoiceIn/InvoiceOut — отдельные документы с позициями;
они не входят в эту финансовую четвёрку.

## Solution

Предложение: добавить четыре пакета сущностей по образцу Product/UOM, каждый с моделью
чтения, CreateRequest, UpdateRequest и Service.Get/Create/Update. Предоставить общий
публичный JSON-контракт ссылки на оплаченную операцию с linkedSum; общий DocumentRate получить
из отдельного document-common. Контракты связей используют Reference, поэтому реализация
SDK для банковских счетов, операций-оснований и positions не блокирует эту группу.

Это DRAFT для согласования детального контракта, а не разрешение реализации.
Состав группы и существующие точки проверки согласованы; перечисленные ниже новые типы,
write scope и правила локальной валидации остаются proposals. Все tickets needs-triage.

## User Stories

1. Как разработчик решения, я хочу получить PaymentIn по ID, чтобы прочитать входящий платеж.
2. Как разработчик решения, я хочу создать PaymentIn с organization и agent, чтобы остальные поля заполнил сервер.
3. Как разработчик решения, я хочу частично изменить PaymentIn, чтобы передать только выбранные значения.
4. Как разработчик решения, я хочу читать и записывать incomingDate и incomingNumber PaymentIn, чтобы сохранить реквизиты входящего документа.
5. Как разработчик решения, я хочу передать agentAccount и organizationAccount PaymentIn ссылками, чтобы указать существующие банковские счета.
6. Как разработчик решения, я хочу получить PaymentOut по ID, чтобы прочитать исходящий платеж.
7. Как разработчик решения, я хочу создать PaymentOut с organization, agent и expenseItem, чтобы явно выбрать статью расходов.
8. Как разработчик решения, я хочу передать в agent PaymentOut Counterparty, Employee или Organization, чтобы использовать документированные виды получателей.
9. Как разработчик решения, я хочу изменить банковские ссылки вместе со сторонами PaymentOut, чтобы сервер проверил их согласованность.
10. Как разработчик решения, я хочу читать и задавать accrualDate PaymentOut, чтобы сохранить дату начисления.
11. Как разработчик решения, я хочу получить CashIn по ID, чтобы прочитать приходный ордер.
12. Как разработчик решения, я хочу создать CashIn с organization и agent, чтобы не задавать необязательные суммы или имя.
13. Как разработчик решения, я хочу изменить CashIn одним PUT, чтобы не делать предварительный GET.
14. Как разработчик решения, я хочу читать и задавать paymentPurpose CashIn, чтобы сохранить основание ордера.
15. Как разработчик решения, я хочу получить CashOut по ID, чтобы прочитать расходный ордер.
16. Как разработчик решения, я хочу создать CashOut с organization, agent и expenseItem, чтобы явно связать расход со статьёй.
17. Как разработчик решения, я хочу читать и задавать accrualDate CashOut, чтобы сохранить дату начисления.
18. Как разработчик решения, я хочу читать и задавать noClosingDocs исходящих документов, чтобы передать признак отсутствия закрывающих документов.
19. Как разработчик решения, я хочу знать несовместимость noClosingDocs=true и непустого operations, чтобы обработать ограничение API.
20. Как разработчик решения, я хочу передавать sum и vatSum без преобразования единиц и округления SDK, чтобы сервер получил заданные числа.
21. Как разработчик решения, я хочу явно передать сумму 0 и флаг false, чтобы они не исчезли при сериализации.
22. Как разработчик решения, я хочу использовать общий DocumentRate с currency и необязательным value, чтобы описать валюту документа.
23. Как разработчик решения, я хочу прочитать отсутствие rate.value, чтобы SDK не добавлял в JSON несуществующее значение.
24. Как разработчик решения, я хочу отправить rate:null, чтобы запросить документированное возвращение к валюте по умолчанию с пересчётом сервером.
25. Как разработчик решения, я хочу прочитать operations как ссылки с linkedSum, чтобы знать распределение платежа.
26. Как разработчик решения, я хочу связать один платеж с несколькими допустимыми типами операций, чтобы не создавать отдельный тип SDK для каждого основания.
27. Как разработчик решения, я хочу опустить linkedSum в запросе, чтобы сохранить предусмотренный API вариант связи только по meta.
28. Как разработчик решения, я хочу явно задать linkedSum=0, чтобы отличать ноль от пропущенного поля.
29. Как разработчик решения, я хочу передать operations:[], чтобы SDK отправил именно пустую коллекцию без локального слияния.
30. Как разработчик решения, я хочу читать factureOut у входящих и factureIn у исходящих документов, чтобы видеть связанные счета-фактуры.
31. Как разработчик решения, я хочу передавать договор, проект, канал продаж и state ссылками, чтобы не ждать SDK всех связанных сущностей.
32. Как разработчик решения, я хочу читать и задавать owner, group и shared, чтобы управлять доступом к документу.
33. Как разработчик решения, я хочу использовать существующие attributes и helpers, чтобы передавать дополнительные поля без нового формата.
34. Как разработчик решения, я хочу задать syncId при создании, чтобы использовать серверный механизм повторного создания с тем же идентификатором.
35. Как разработчик решения, я хочу читать created, updated, deleted, printed и published, чтобы учитывать состояние документа.
36. Как разработчик решения, я хочу получать нулевые значения для отсутствующих полей чтения, чтобы разреженные ответы не считались ошибкой.
37. Как разработчик решения, я хочу различать отсутствие, null, пустую строку, пустой массив и значение request-поля, чтобы формировать точный JSON.
38. Как разработчик решения, я хочу использовать обе публичные JSON-библиотеки Go, чтобы не зависеть от выбранного кодека приложения.
39. Как разработчик решения, я хочу получать фактическую модель ответа Create/Update, чтобы видеть серверные значения и пересчёты.
40. Как разработчик решения, я хочу отменять запросы через context, чтобы ограничивать их время жизни.
41. Как разработчик решения, я хочу распознавать APIError и исходные ошибки через errors.As/errors.Is, чтобы обрабатывать ошибки программно.
42. Как разработчик решения, я хочу знать ограничения float64, чтобы не принимать транспортную модель за инструмент точной финансовой арифметики.
43. Как сопровождающий SDK, я хочу проверить общий allocation JSON отдельно от сущностей, чтобы его контракт оставался самостоятельным.
44. Как сопровождающий SDK, я хочу выполнить каждый entity ticket целиком в отдельном контексте, чтобы модель, запросы, HTTP и документация не расходились.

## Implementation Decisions

### Утверждённые основания

- Состав: только paymentin, paymentout, cashin, cashout. Счета invoicein/invoiceout
  могут встречаться в Reference.Meta.Type операций, но их SDK находится за границей группы.
- Сохраняются Product/UOM, отдельные request-типы, Field с omitzero, DateTime,
  открытый Reference/MetaType и существующие ошибки/transport.Client.
- Согласованные seams: публичная JSON-сериализация обоими encoding/json и encoding/json/v2;
  локальный HTTP-сервер с настоящим transport.Client. Для HTTP структура UOM:
  отдельные TestServiceGet, TestServiceCreate, TestServiceUpdate с таблицами сценариев.
- Expand отложен до самого конца SDK. Эта группа использует только ссылки и пустой query.

### Предлагаемый контракт

Полный перечень полей чтения, Create/Update, null, required и write-once вынесен в
[contract.md](contract.md). При согласовании или реализации любой из четырёх сущностей
прочитать его целиком вместе с этой spec; это нормативная часть DRAFT, а не приложение
с необязательными примерами. Проверенные основания и противоречия — в
[research.md](research.md).

| Пакет | Модель | Ресурс относительно API 1.2 | Минимальный Create JSON: только перечисленные поля |
| --- | --- | --- | --- |
| paymentin | PaymentIn | entity/paymentin | organization и agent: объекты с meta |
| paymentout | PaymentOut | entity/paymentout | organization, agent и expenseItem: объекты с meta |
| cashin | CashIn | entity/cashin | organization и agent: объекты с meta |
| cashout | CashOut | entity/cashout | organization, agent и expenseItem: объекты с meta |

- Proposal: в каждом пакете Service/NewService и Get(ctx,id), Create(ctx,CreateRequest),
  Update(ctx,id,UpdateRequest); возвращают указатель на свою модель и error.
  NewService принимает существующий *transport.Client и отклоняет nil.
- Get: GET ресурса /{id}, тело отсутствует. Create: POST основного ресурса.
  Update: PUT ресурса /{id}; минимальный непустой пример — только applicable:false.
  Ни один метод не делает предварительных или follow-up-запросов.
- Proposal: обязательные Create.Organization/Create.Agent и исходящий Create.ExpenseItem —
  обычные Reference. Остальные поля Create и все поля Update — Field соответствующего типа.
  Name необязателен: Product.Create.Name нельзя механически переносить на документы.
- Proposal: все перечисленные write-поля включаются; files и расширенные объекты связей
  исключены явно. SyncID предлагается в Create и Update: API запрещает изменение после заполнения,
  а не первоначальное заполнение существующего документа. Проверяет сервер; SDK не делает Get.
- Proposal: общий model.FinanceOperation содержит Meta и LinkedSum float64;
  отдельный model.FinanceOperationRequest содержит Meta и LinkedSum Field[float64].
  Это ссылки на операции, без вложенных моделей и без дискриминированного union.
  Имена и числовое представление требуют утверждения.
- DocumentRate/DocumentRateRequest — внешний контракт document-common; не копировать его в финансовые пакеты.
  Read currency — *Reference, read value — nullable/optional число с различием «нет» и 0;
  write currency/value должны поддерживать частичный JSON через согласованный DocumentRateRequest.
- Proposal: Sum, VATSum и LinkedSum — float64 по существующему Price и документальному Float.
  JSON 80000 и 80000.0 читается одинаково; SDK не пересчитывает, не масштабирует и не округляет
  суммы. Сохранение исходного текстового написания числа и десятичная точность сверх float64
  не обещаются. Валютный курс — отдельная величина, не денежная сумма.
- Сервер контролирует допустимость сторон, счетов, валют, noClosingDocs и операций.
  При смене сторон банковских платежей документация требует соответствующие account-поля;
  SDK передаёт точно заданный запрос и ошибку сервера, не подбирает счета.
- Proposal: локально проверить UUID маршрута; у Create только непустой Meta.Href обязательных
  ссылок; у Update — хотя бы одно заданное Field. Не проверять типы/доступность ссылок,
  суммы, null, noClosingDocs, сочетания счетов или наличие expenseItem по результату запроса.
  Эти новые правила required-reference validation ещё не утверждены.
- Proposal по образцу UOM: успех только HTTP 200 и непустой ID; при ошибке nil результата.
  Сохранить контекст операции и errors.As/errors.Is, существующие context, gzip и заголовки.
  Retry, новые транспортные механизмы и сравнение ID ответа с ID маршрута не добавлять.

### Зависимости и разбиение

| Ticket | Что даёт | Blocked by |
| --- | --- | --- |
| [01 — FinanceOperation](issues/01-finance-operation.md) | Самостоятельный read/write JSON-контракт операций | нет технических зависимостей; требуется triage |
| [02 — PaymentIn](issues/02-paymentin.md) | Полный согласованный Get/Create/Update с документацией | 01; document-common/01-rate |
| [03 — PaymentOut](issues/03-paymentout.md) | Полный согласованный Get/Create/Update с документацией | 01; document-common/01-rate |
| [04 — CashIn](issues/04-cashin.md) | Полный согласованный Get/Create/Update с документацией | 01; document-common/01-rate |
| [05 — CashOut](issues/05-cashout.md) | Полный согласованный Get/Create/Update с документацией | 01; document-common/01-rate |

Зависимость на DocumentRate — [внешний ticket](../document-common/issues/01-rate.md).
Entity tickets независимы друг от друга; positions, SDK банковских счетов и реализация
оснований не являются blockers. Needs-triage сохраняется и для ticket без blockers.

## Testing Decisions

- Общие типы и entity models/requests проверять публичной сериализацией обоими JSON-кодеками.
  Проверять значения каждого поля из contract, а не только отсутствие ошибки декодирования.
- Независимые golden JSON для минимального/полного Create и полного/частичного Update;
  zero/unset, false, 0, пустая строка, [], nil slice как [], null и SyncID в Create/Update.
  Required refs, readonly и excluded поля должны иметь явно проверенный request-scope.
- FinanceOperation: meta-only request и linkedSum=0, дробное значение, null как wire
  (без обещания принятия API), неизвестный Meta.Type и игнорирование лишних полей ответа.
- Read DocumentRate с отсутствующим value, числом 0, 1 и дробным значением; outer rate absent/null;
  write rate:null и частичный DocumentRateRequest. Собственные JSON-тесты DocumentRate выполняет его ticket;
  entity tests доказывают корректное включение в requests и декодирование ответов.
- HTTP-тесты по UOM: точный метод/путь и пустой query, тело GET отсутствует,
  Content-Type/Accept/Authorization/User-Agent/Accept-Encoding, plain/gzip,
  один вызов, фактические данные ответа, отличающиеся от входных.
- Для каждого метода: nil/отменённый context; nil client; malformed UUID Get/Update;
  Create без обязательных ссылок и Update без полей по утверждённому решению triage;
  APIError 400/404/429/5xx, ошибка transport/auth, errors.As/errors.Is,
  неожиданный статус, malformed JSON, null, пустой/отсутствующий ID. Retry отсутствует.
- PaymentIn проверяет incomingDate/incomingNumber и банковские ссылки; PaymentOut —
  accrualDate, polymorphic agent, expenseItem/noClosingDocs и банковские ссылки;
  CashIn — отсутствие банковских полей; CashOut — accrualDate и expenseItem/noClosingDocs.
- Ограничения 22004, noClosingDocs и типов операций проверять как передачу синтетической
  APIError, если она нужна для сценария; local server не доказывает бизнес-поведение живого API.
- Upstream fixtures сохранять без изменений с существующим provenance.
  Fixtures этой группы содержат расширенные поля операций и files.rows:
  тестировать поддержанные meta/linkedSum, остальные поля игнорировать. Наличие их в fixture
  не означает поддержку expand; собственные маленькие unexpanded cases обязательны.
- После будущей реализации: форматирование, go test -race -count=1 ./...,
  go vet ./... ./testdata, git diff --check и независимые Standards/Spec reviews
  с фиксированным BASE. Подготовка текущего DRAFT не запускает implementation pipeline.

## Out of Scope

- InvoiceIn/InvoiceOut, positions и получение строк/пагинация; остальные финансовые,
  розничные и складские сущности. Ссылка на них не расширяет scope.
- Files: metadata коллекции, rows, upload/download и изменение вложений.
- Банковские счета как сущности/коллекции, их создание, изменение и получение;
  разрешены только document.agentAccount/organizationAccount как Reference у paymentin/out.
- Расширенные объекты любых ссылок, expand и query-параметр expand; общая поддержка expand
  отложена до самого конца SDK.
- List, фильтры, поиск, сортировка, пагинация, Delete, массовые операции, /new templates,
  обработка metadata и изменение определений attributes.
- Локальная имитация write-once syncId; запись accountId, id, meta, created, updated, deleted,
  printed, published; недокументированные поля сверх явных таблиц contract.
- Финансовая арифметика, новый Money/Decimal API, конвертация валют, округление,
  локальная имитация серверного распределения linkedSum или очистки operations.
- Общий CRUD, генерация сущностей, рефакторинг transport/Product/UOM, retry/rate limiting.
- Живой API, .env, операции с реальным аккаунтом, реализация кода, commits и push в этой сессии.

## Further Notes

- Официальный master проверен 2026-11-01: HEAD
  9a811cef3dc4a3f3a46f070e60452361a06bf7cd, совпадает с revision прежнего исследования.
  Источник контрактов — официальный api-remap-1.2-doc; fixtures — вспомогательное покрытие.
- Пользователь разрешил сразу записать DRAFT spec/tickets без интервью; это отменяет
  промежуточные вопросы skills to-spec/to-tickets для текущего синтеза, но не утверждает
  новые API/design decisions. Triage должен согласовать contract, null-policy, имена
  FinanceOperation/Request, float64 и required-reference validation перед реализацией.
- Отдельно разрешить расхождения cashin/cashout account-полей, исходящий agent CashOut
  и семантику полного замещения/очистки operations; proposals записаны в research/contract.
- Spec и каждый ticket сохраняют needs-triage до решения. Готовность DocumentRate или общего
  allocation не делает неподтверждённый entity ticket ready-for-agent автоматически.
