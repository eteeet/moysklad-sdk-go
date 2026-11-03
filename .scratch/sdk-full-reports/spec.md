Status: needs-triage
Stage: DRAFT

# 5. Отчёты, печать и интеграции

## Problem Statement

Добавить оставшиеся отчёты, печать документов и этикеток, шаблоны и публикации, вебхуки, уведомления, ленту событий, задачи и оставшиеся служебные ресурсы.

## Solution

Завершённые вертикальные срезы; каждый ticket включает models, methods, тесты и собственное
описание полей. Матрица ниже определяет scope этой захватки.
[Общая карта](../sdk-full/map.md), [нормативные ограничения v1](../sdk-v1/spec.md),
[каталог JSON-контрактов](../sdk-v1/contracts.md), [backlog](../../docs/status-backlog.md).
Источник: api-remap-1.2-doc revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.

## User Stories

1. Как разработчик решения, я могу использовать «Report stock/all (расширенный)» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
2. Как разработчик решения, я могу использовать «Report stock/bystore» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
3. Как разработчик решения, я могу использовать «Краткие отчёты stock/all/current и bystore/current» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
4. Как разработчик решения, я могу использовать «stock/byoperation и byslot/current» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
5. Как разработчик решения, я могу использовать «Report byoperations stock/reserve/intransit» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
6. Как разработчик решения, я могу использовать «Report dashboard day/week/month» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
7. Как разработчик решения, я могу использовать «Report counterparty» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
8. Как разработчик решения, я могу использовать «Report money plotseries/byaccount» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
9. Как разработчик решения, я могу использовать «Report profit (5 разрезов)» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
10. Как разработчик решения, я могу использовать «Report orders/sales plotseries» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
11. Как разработчик решения, я могу использовать «Шаблоны печатных форм (read-only)» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
12. Как разработчик решения, я могу использовать «Печать документов export» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
13. Как разработчик решения, я могу использовать «Печать этикеток и ценников» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
14. Как разработчик решения, я могу использовать «Публикации документов» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
15. Как разработчик решения, я могу использовать «Webhook CRUD» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
16. Как разработчик решения, я могу использовать «WebhookStock CRUD» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
17. Как разработчик решения, я могу использовать «Notification: лента и прочтение» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
18. Как разработчик решения, я могу использовать «Notification settings GET/PUT» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
19. Как разработчик решения, я могу использовать «Лента событий документов (notes)» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
20. Как разработчик решения, я могу использовать «Task + типы задач и комментарии» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
21. Как разработчик решения, я могу использовать «NamedFilter (read-only)» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
22. Как разработчик решения, я могу использовать «Application GET» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
23. Как разработчик решения, я могу использовать «Async создания отчётов» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.

## Implementation Decisions и scope tickets

| ID | Ticket | Blocked by |
| --- | --- | --- |
| I01 | [Report stock/all (расширенный)](issues/01-stock-all.md) | нет |
| I02 | [Report stock/bystore](issues/02-stock-bystore.md) | нет |
| I03 | [Краткие отчёты stock/all/current и bystore/current](issues/03-stock-current.md) | нет |
| I04 | [stock/byoperation и byslot/current](issues/04-stock-operation-slot.md) | нет |
| I05 | [Report byoperations stock/reserve/intransit](issues/05-byoperations.md) | нет |
| I06 | [Report dashboard day/week/month](issues/06-dashboard.md) | нет |
| I07 | [Report counterparty](issues/07-report-counterparty.md) | нет |
| I08 | [Report money plotseries/byaccount](issues/08-money.md) | нет |
| I09 | [Report profit (5 разрезов)](issues/09-profit.md) | нет |
| I10 | [Report orders/sales plotseries](issues/10-orders-sales.md) | нет |
| I11 | [Шаблоны печатных форм (read-only)](issues/11-templates.md) | нет |
| I12 | [Печать документов export](issues/12-print.md) | I11 |
| I13 | [Печать этикеток и ценников](issues/13-labels.md) | I11, I12 |
| I14 | [Публикации документов](issues/14-publications.md) | I11 |
| I15 | [Webhook CRUD](issues/15-webhook.md) | нет |
| I16 | [WebhookStock CRUD](issues/16-webhookstock.md) | нет |
| I17 | [Notification: лента и прочтение](issues/17-notification.md) | нет |
| I18 | [Notification settings GET/PUT](issues/18-notification-settings.md) | нет |
| I19 | [Лента событий документов (notes)](issues/19-eventfeed.md) | нет |
| I20 | [Task + типы задач и комментарии](issues/20-task.md) | нет |
| I21 | [NamedFilter (read-only)](issues/21-namedfilter.md) | нет |
| I22 | [Application GET](issues/22-application.md) | нет |
| I23 | [Async создания отчётов](issues/23-reports-async.md) | I01, I02, I03, I04, I06, I07, I08, I09, I10 |

### I01: Report stock/all (расширенный)

GET report/stock/all с параметрами отчёта и страницами.

Проверки и особенности: Вес 5 единиц лимита; stockMode/groupBy — открытые; деньги — копейки.

Источники: [stock](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_stock.md).

### I02: Report stock/bystore

GET report/stock/bystore.

Проверки и особенности: Вес 5 единиц; stockByStore — вложенный массив.

Источники: [stock](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_stock.md).

### I03: Краткие отчёты stock/all/current и bystore/current

GET report/stock/all/current и report/stock/bystore/current.

Проверки и особенности: Ответ — plain array, не Collection; параметры include/changedSince/stockType.

Источники: [stock](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_stock.md).

### I04: stock/byoperation и byslot/current

GET report/stock/byoperation (по документу) и report/stock/byslot/current.

Проверки и особенности: Формы двух endpoints раздельные: byslot/current — plain array со строковыми UUID assortmentId/storeId/slotId и stock, без meta и Reference; byoperation — по таблице, operation.id обязателен.

Источники: [stock](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_stock.md).

### I05: Report byoperations stock/reserve/intransit

GET report/byoperations/stock|reserve|intransit.

Проверки и особенности: Один пакет, три метода с общей формой строки, если источник подтверждает.

Источники: [byoperations](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_by_operations.md).

### I06: Report dashboard day/week/month

GET report/dashboard/day|week|month.

Проверки и особенности: Показатели sales/orders/money — копейки; movement — проценты.

Источники: [dashboard](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_dashboard.md).

### I07: Report counterparty

GET report/counterparty, POST выборочных показателей, GET по контрагенту.

Проверки и особенности: Выборочный запрос — POST с массивом ссылок, не запись.

Источники: [counterparty](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_counterparty.md).

### I08: Report money plotseries/byaccount

GET report/money/plotseries и report/money/byaccount.

Проверки и особенности: interval обязателен для plotseries; суммы — копейки.

Источники: [money](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_money.md).

### I09: Report profit (5 разрезов)

GET report/profit/byproduct|byvariant|byemployee|bycounterparty|bysaleschannel.

Проверки и особенности: byEmployee и byemployee встречаются в источнике — выбрать и задокументировать путь.

Источники: [profit](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_pnl.md).

### I10: Report orders/sales plotseries

GET report/orders/plotseries и report/sales/plotseries.

Проверки и особенности: Общая форма series, если источник подтверждает одинаковые поля.

Источники: [sales/orders](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_sales_orders.md).

### I11: Шаблоны печатных форм (read-only)

Стандартные и пользовательские шаблоны документов и этикеток: List/Get (entity/{type}/metadata/embeddedtemplate|customtemplate).

Проверки и особенности: Только чтение; тип шаблона — открытый.

Источники: [template](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_template.md).

### I12: Печать документов export

POST entity/{type}/{id}/export: одиночная и пакетная печать. Transport: режим без следования 303 для этого вызова, сохранение status и Location; явная загрузка результата отдельным вызовом без ослабления same-origin контракта JSON-вызовов.

Проверки и особенности: Сейчас Client сам следует 303 и делает скрытый GET — тест это исключает. Origin ссылки на файл источник не фиксирует: не считать любой Location поддержанным Download. 202 источник называет будущим, протокол опроса не задан — вернуть различимые результаты 202/303 без Wait.

Источники: [print](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_print.md), [workbook](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/workbook/_workbook_print_templates.md).

### I13: Печать этикеток и ценников

POST entity/{product|variant|...}/{id}/export — этикетки и ценники.

Проверки и особенности: organization, count, salePrice, template — обязательные поля запроса.

Источники: [labels](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_labels.md).

### I14: Публикации документов

List/Get/Create/Delete entity/{type}/{id}/publication.

Проверки и особенности: href публикации — readonly; template — ссылка.

Источники: [publication](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_publication.md).

### I15: Webhook CRUD

Вебхук: Get/List/Create/Update/Delete, массовые.

Проверки и особенности: action/entityType/diffType — открытые; url без проверки доступности.

Источники: [webhook](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_webhook.md).

### I16: WebhookStock CRUD

Вебхук остатков: Get/List/Create/Update/Delete, массовые.

Проверки и особенности: reportType/stockType — открытые.

Источники: [webhookstock](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_webhookstock.md).

### I17: Notification: лента и прочтение

List/Get/Delete notification, markasread, markasreadall.

Проверки и особенности: Типы уведомлений — открытый meta.type с raw; детали типов по _notification_types.

Источники: [notification](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/notification/_notification.md), [types](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/notification/_notification_types.md).

### I18: Notification settings GET/PUT

GET/PUT notification/settings.

Проверки и особенности: Группы и каналы — открытые ключи; enabled — явный false.

Источники: [settings](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/notification/_notification_settings.md).

### I19: Лента событий документов (notes)

List/Get/Create/Update/Delete и массовые create/update/delete entity/{type}/{id}/notes для документов из источника.

Проверки и особенности: Упоминания в тексте — строка без разбора.

Источники: [eventfeed](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_eventfeed.md).

### I20: Task + типы задач и комментарии

Задача: Get/List/Create/Update/Delete, массовые; комментарии; типы задач CRUD.

Проверки и особенности: assignee/operation/agent — ссылки; dueToDate — DateTime.

Источники: [task](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_task.md).

### I21: NamedFilter (read-only)

Сохранённые фильтры: List (свои и другого пользователя), Get.

Проверки и особенности: Применение уже есть в query; только чтение.

Источники: [namedfilter](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_named_filter.md).

### I22: Application GET

GET entity/application/{id}.

Проверки и особенности: Только чтение; appUid — строка.

Источники: [general](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/_general.md).

### I23: Async создания отчётов

ListAsync (создание async-задачи) для отчётов из списка _async.md: Остатки, Прибыльность, Деньги, Показатели продаж и заказов, Показатели контрагентов, Показатели (dashboard).

Проверки и особенности: Только endpoints, разрешённые источником; выборочные показатели контрагентов (POST) — без async; результат через существующий пакет async, по образцу ListAsync counterparty/assortment.

Источники: [async](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/_async.md).

## Общий контракт реализации

- Перед кодом — field/method freeze: поимённая Write matrix, readonly/TODO,
  null/required/Create-only/Update-only; status needs-triage сам не разрешает реализацию.
- Образец пакета и тестов указан в каждом ticket; новый пакет повторяет его структуру
  (model.go, request.go, service.go, *_test.go), без новых общих абстракций.
- Delete, CreateMany/UpdateMany/DeleteMany — только если раздел источника их описывает,
  по образцу [optional X01/X02](../sdk-v1-optional/spec.md).
- Новый ресурс получает handle фасада `moysklad.Client`, строку в
  [status-entities](../../docs/status-entities.md) / [status-technical](../../docs/status-technical.md)
  и проходит `status_test.go`; пункт снимается из [status-backlog](../../docs/status-backlog.md).
- JSON: encoding/json v1 и v2, Field+omitzero для optional request, независимые
  golden, неизменённые upstream с SHA-256 provenance; raw JSON long не проходит float64.
- Dates — существующий DateTime; деньги float64 в копейках, counts/проценты/баллы/курс
  не объявляются деньгами. Reference без expand, unknown enum/type не отвергается.
- HTTP: настоящий transport+httptest, exact method/path/query/headers/body,
  plain/gzip, фактический server response, ошибки/auth/context/cancellation,
  errors.Is/As и единый moysklad:. Обычный вызов — одна отправка, если caller
  явно не включил retry/walker; каждый дополнительный запрос доказуем тестами.
- Бизнес-ограничения и permissions оставлены API; SDK не выполняет pre-GET ради
  local validation. Exception — только явный пользовательский multi-call operation.
- Проверки: `make check` (go test -race, go vet ./... ./testdata, gofmt, SHA-256 фикстур,
  статус сущностей, diff --check); независимые Standards/Spec review на фиксированном BASE.
- Коммит по ticket; shared dependency раньше зависимого, публичные docs отдельным
  завершающим commit волны. Planning не коммитить без отдельной команды.
- Live API, credentials, push/tag/release и worktree cleanup не входят.

## Open Decisions

- Каждый отчёт — собственная форма ответа, без общего типа на все отчёты.
- stock/all и stock/bystore весят 5 единиц лимита — уже в limiter, тест подтверждает вес.
- Печать: 303 + Location на файл возвращается вызывающему без следования redirect; 202 источник называет будущим, протокол опроса не задан — различимый результат без Wait (см. I12).
- application — только GET entity/application/{id}; context/application уже есть.
- Async создания отчётов — отдельный I23 после синхронных методов.

## Out of Scope

Expand; приём входящих вебхуков (HTTP-сервер) — не SDK клиента. Все не перечисленные в tickets ресурсы/операции.
Этот DRAFT не является командой реализации или публикации.
