Status: needs-triage
Stage: DRAFT

# InvoiceIn и InvoiceOut: отдельная захватка счетов

DRAFT от 2026-11-01. Счета отделены от финансовой четвёрки
paymentin/paymentout/cashin/cashout по явной команде пользователя.
Подготовка разрешена; подробные поля, имена новых типов и validation ниже
предлагаются к утверждению, реализация пока не разрешена этой spec.

## Problem Statement

Разработчик решения хочет читать, создавать и частично обновлять счета поставщика
и покупателя, включая товарные позиции. Наличие банковских и кассовых документов
не должно смешивать их контракты со счетами и не является зависимостью этой группы.

## Solution

Два самостоятельных пакета invoicein/invoiceout с моделями чтения, отдельными
CreateRequest/UpdateRequest и Service по образцу Product/UOM.
Валюта берётся из общих типов документов, позиции — из универсального positions.
Строки позиций читаются отдельным постраничным reader общего модуля, без expand.

## User Stories

1. Как разработчик решения, я хочу получить InvoiceIn по ID, чтобы прочитать счёт поставщика.
2. Как разработчик решения, я хочу получить InvoiceOut по ID, чтобы прочитать счёт покупателю.
3. Как разработчик решения, я хочу создать счёт с номером, юрлицом и контрагентом.
4. Как разработчик решения, я хочу различать обязательные поля Create и частичные поля Update.
5. Как разработчик решения, я хочу явно отправить applicable=false, чтобы не потерять заданный флаг.
6. Как разработчик решения, я хочу передать moment и планируемую дату оплаты через существующий DateTime.
7. Как разработчик решения, я хочу читать суммы и НДС без пересчёта единиц SDK.
8. Как разработчик решения, я хочу читать оплату и отгрузку как серверные вычисленные суммы.
9. Как разработчик решения, я хочу передавать rate без подстановки валюты или курса клиентом.
10. Как разработчик решения, я хочу передавать ссылки на счета сторон, чтобы корректно сменить organization или agent.
11. Как разработчик решения, я хочу читать метаданные positions, даже если строк в Get нет.
12. Как разработчик решения, я хочу отдельно получить страницу позиций без expand.
13. Как разработчик решения, я хочу передать позиции в Create без дополнительного HTTP-запроса.
14. Как разработчик решения, я хочу заменить позиции через Update, зная о полной замене коллекции.
15. Как разработчик решения, я хочу различать незаданные позиции и явно заданный пустой массив.
16. Как разработчик решения, я хочу сохранить явные quantity:0, price:0 и vatEnabled:false.
17. Как разработчик решения, я хочу использовать разные assortment meta.type без заранее реализованных связанных сущностей.
18. Как разработчик решения, я хочу передать attributes существующими helpers.
19. Как разработчик решения, я хочу прочитать входящий номер и дату только у InvoiceIn.
20. Как разработчик решения, я хочу связать InvoiceIn с purchaseOrder и InvoiceOut с customerOrder.
21. Как разработчик решения, я хочу передать канал продаж InvoiceOut отдельно от InvoiceIn.
22. Как разработчик решения, я хочу задать syncId при создании, сохранив ограничение его изменения.
23. Как разработчик решения, я хочу получить фактический объект ответа, а не эхо запроса.
24. Как разработчик решения, я хочу распознавать APIError и отменять операцию через context.
25. Как сопровождающий SDK, я хочу одинаковую структуру HTTP-тестов в новых сущностях.
26. Как разработчик решения, я хочу видеть отложенные поля и операции явно в документации.

## Implementation Decisions

### Публичные методы и валидация — proposal

NewService принимает существующий *transport.Client, возвращает *Service и error,
отклоняет nil. Get(ctx,id), Create(ctx,input), Update(ctx,id,input) возвращают
(*InvoiceIn,error) либо (*InvoiceOut,error). Ресурсы — entity/invoicein и entity/invoiceout;
Get/Update используют /{id}. Один HTTP-запрос на вызов, без query/expand и тела Get.

В Create обязательны Name string, Organization model.Reference и Agent model.Reference.
Это подтверждено разделами «Создать», хотя в таблице атрибутов name не помечен
необходимым при создании. Предлагаемая локальная validation: Name!="", без trim;
у двух обязательных ссылок Meta.Href!="", без новой URL- или предметной проверки.
Остальные записываемые поля Create используют Field с omitzero.
В Update все изменяемые поля — Field; полностью незаданный запрос отклоняется.
Все Get/Update проверяют форму UUID без ограничения версии.

Успех — HTTP 200 и непустой ID; при любой ошибке nil-модель.
Не проверять совпадение ID маршрута/ответа и не добавлять контроль полноты всех полей.
Сохранить errors.As/errors.Is и контекст без двойного moysklad:.
Не добавлять автоматических повторов, предварительного Get, merge коллекций или follow-up.

### Поля общего среза

| JSON | Поля Go | Тип чтения | Запись |
| --- | --- | --- | --- |
| meta | Meta | model.Meta | нет |
| id, accountId | ID, AccountID | string | нет |
| created, updated, deleted | Created, Updated, Deleted | model.DateTime | нет |
| name, description, code, externalCode | Name, Description, Code, ExternalCode | string | Create/Update; Name обязателен в Create |
| moment, paymentPlannedMoment | Moment, PaymentPlannedMoment | model.DateTime | Create/Update |
| applicable, shared, vatEnabled, vatIncluded | Applicable, Shared, VATEnabled, VATIncluded | bool | Create/Update |
| organization, agent | Organization, Agent | *model.Reference | обязательные Reference в Create; Field[Reference] в Update |
| organizationAccount, agentAccount, owner, group, state, project, store, contract | OrganizationAccount, AgentAccount, Owner, Group, State, Project, Store, Contract | *model.Reference | Create/Update через Field[Reference] |
| rate | Rate | *model.DocumentRate | Create/Update через Field[model.DocumentRateRequest] |
| attributes | Attributes | []model.Attribute | Create/Update через Field[[]model.Attribute] |
| positions | Positions | *model.Collection[positions.Commercial] | Create/Update через Field[[]positions.CommercialRequest] |
| sum, vatSum, payedSum, shippedSum | Sum, VATSum, PayedSum, ShippedSum | float64 | нет |
| printed, published | Printed, Published | bool | нет |
| syncId | SyncID | string | Create/Update через Field[string], write-once проверяет сервер |
| payments | Payments | []model.Reference | только чтение связей в этом срезе |

Даты, nullable ссылки и вложенные объекты оформляются с omitzero как в Product/UOM.
Отсутствующие поля чтения дают нулевые значения, неизвестные поля не ломают чтение.
Суммы остаются в единицах API, без округления и новой финансовой арифметики.
Meta/ID/AccountID/Created/Updated/Deleted, суммы и Printed/Published не входят в requests.
SyncID доступен в модели чтения и обоих requests. Запрет изменения после заполнения
проверяет сервер; SDK не загружает старое значение и не блокирует первоначальный Update.

### Различия счетов

| Документ | JSON | Go / тип чтения | Запись |
| --- | --- | --- | --- |
| InvoiceIn | incomingDate | IncomingDate model.DateTime | Create/Update Field[DateTime] |
| InvoiceIn | incomingNumber | IncomingNumber string | Create/Update Field[string] |
| InvoiceIn | purchaseOrder | PurchaseOrder *model.Reference | Create/Update Field[Reference] |
| InvoiceIn | supplies | Supplies []model.Reference | только чтение связей в этом срезе |
| InvoiceOut | customerOrder | CustomerOrder *model.Reference | Create/Update Field[Reference] |
| InvoiceOut | salesChannel | SalesChannel *model.Reference | Create/Update Field[Reference] |
| InvoiceOut | demands | Demands []model.Reference | только чтение связей в этом срезе |

Запись ссылок purchaseOrder/customerOrder/salesChannel — proposal по первичной
таблице и разделу связей. Не расширять автоматически запись всех связанных коллекций.
Полноценные SDK-пакеты Order/Store/Project/Account/Payment не нужны для Reference.

### Позиции и связанные счета

Коммерческие поля позиции задаются общей spec positions: identity, assortment,
quantity, price, discount, vat, vatEnabled. Упаковки, серийные номера, tracking,
stock enrichment и все дополнительные поля, не названные здесь, остаются TODO.
InvoiceIn допускает ассортимент из его документированного списка; InvoiceOut
может ссылаться также на комплект. SDK передаёт существующие Reference без
закрытого enum meta.type; сервер проверяет применимость конкретного объекта.

Обычный Get не раскрывает rows. Его Positions содержит meta и допускает отсутствие Rows.
Страница строк — positions.Service[positions.Commercial].GetPage с documentType
invoicein либо invoiceout. Entity-пакеты не дублируют reader и не добавляют expand.
Reader — самостоятельный ticket и обязательная зависимость entity tickets с GetPositions.

При embedded Update переданный positions полностью заменяет набор.
Незаданный Field не отправляется; Value([]) отправляет []; Null сохраняется
как wire null, допустимость проверяет сервер. SDK не выполняет локальный merge,
проверку существования строк, загрузку предыдущих позиций или арифметику.
ID в request позиции — selector существующей строки Update по общей spec positions,
не назначение новых ID при Create. Proposed: Create/Update отклоняют embedded
массив больше 1000 строк до HTTP; предел 1000 проходит. Для omitted/null/[]
нет скрытого заполнения, split или дополнительных запросов.

При смене organization и agent официальная документация требует одновременно
organizationAccount и agentAccount соответственно. SDK предоставляет оба поля,
передаёт запрос как задано; предложено оставить проверку сочетаний серверу.
Не имитировать эти серверные ограничения локальным HTTP-тестом.

### Открытые решения перед утверждением

1. Публичные имена и shapes общих positions/DocumentRate.
2. Минимальная локальная validation обязательных Reference через Meta.Href.
3. SyncID в обоих requests без local write-once validation; связанные массивы — выбранный readonly scope.
4. Первоначальный отказ от pack и остальных расширений позиции.
5. Финальное утверждение таблиц полей и точек проверки этой spec.

## Testing Decisions

Утверждённые ранее seams переиспользуются: публичные JSON-модели обоими кодеками
и локальный HTTP-сервер с настоящим transport.Client. Новые mock interfaces не нужны.
Структура HTTP-тестов каждого ticket — entity/uom/service_test.go:
отдельные TestServiceGet/Create/Update с таблицами сценариев внутри.
Общий positions reader имеет свои публичные HTTP-тесты.

- Неизменённые upstream invoice_in/invoice_out: значения каждого поддержанного поля,
  nullable/absent/unknown поля, даты, account-ссылки и readonly суммы.
- Независимые golden минимального и полного Create/Update; один Field Update,
  false/0/""/[]/null, rate без value, позиции с ID и без ID, отсутствие readonly/TODO-полей.
- Оба JSON-кодека: все общие числовые/булевы/optional состояния; большие integer
  значения Attribute канонизируются без округления за пределами 2^53.
- HTTP: GET/POST/PUT, точный путь и заголовки, нет query/тела Get/лишних запросов;
  plain/gzip, данные сервера отличаются от input; HTTP 200+ID.
- Nil/cancelled context; nil client; локальная validation; transport/auth errors;
  APIError 4xx/5xx/429, errors.As/errors.Is; неожиданный статус, malformed/null JSON,
  отсутствие/пустой ID, nil-модель при любой ошибке.
- Проверить одну отправку positions как задано; не объявлять поведение сервера
  полной замены проверенным синтетической реализацией серверного merge.
- Каждый ticket завершается gofmt, полным race/vet/diff-check и независимыми Standards/Spec review.

## Out of Scope

Финансовая четвёрка paymentin/paymentout/cashin/cashout; files и их metadata/операции;
упаковки и расширения позиций; запись связанных arrays payments/supplies/demands;
List документов, фильтры/поиск/сортировка, auto-pagination walker, отдельная запись
позиционных endpoints, Delete, batch, new/templates и автосоздание по основанию.
Expand — только после всего остального SDK, если понадобится.
Нет live API, .env, новых универсальных transport frameworks или изменения Product/UOM.
Планирование не коммитится без команды; реализация требует утверждения DRAFT.

## Further Notes

Контракты проверены по official revision 9a811cef3dc4a3f3a46f070e60452361a06bf7cd,
тот же master повторно подтверждён 2026-11-01. Исходный dirty temp_plan сохраняется.
Два entity tickets независимы после общей валюты и JSON positions.
Для end-to-end чтения строк к завершению захватки нужен общий GetPage reader.

- [InvoiceIn, поля и создание](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_invoice_in.md#L348).
- [InvoiceOut, поля и создание](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_invoice_out.md#L476).
- [Смена сторон InvoiceIn](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_invoice_in.md#L1958).
- [Смена сторон InvoiceOut](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_invoice_out.md#L2010).
- [Полная замена positions](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1031).
- [Общие позиции](../document-positions/spec.md).
- [Общая валюта](../document-common/spec.md).

## Типизированное чтение позиций — согласованное уточнение

У каждого entity Service добавить GetPositions(ctx context.Context, id string, page positions.PageOptions)
(*model.Collection[Position], error). Он вызывает общий reader с фиксированными documentType
и Position; вызывающий код не выбирает строку ресурса и generic type. Один HTTP GET, без expand,
автообхода и скрытого Get заголовка. Ошибки/envelope/UUID/page validation остаются контрактом P02.
Move экспортирует alias Position = positions.Movement; InvoiceIn/InvoiceOut/PurchaseReturn —
alias Position = positions.Commercial. Supply/Demand/PurchaseOrder используют собственные Position.
Тестовый seam — публичный GetPositions через transport.Client + httptest.Server;
отдельный TestServiceGetPositions с таблицами правильного path/type/query, plain/gzip, пустой страницы,
ошибок и одной отправки. P02 — обязательная зависимость всего entity ticket.
Это уточнение не разрешает реализацию entity, остальные open decisions остаются needs-triage.
