Status: needs-triage

# 3. Документы и операции документов

## Problem Statement

Добавить общие операции документов (отдельные позиции, шаблон new, автозаполнение, корзина, запись связей) и оставшиеся торговые и финансовые документы.

## Solution

Завершённые вертикальные срезы; каждый ticket включает models, methods, тесты и собственное
описание полей. Матрица ниже определяет scope этой захватки.
[Общая карта](../sdk-full/map.md), [нормативные ограничения v1](../sdk-v1/spec.md),
[каталог JSON-контрактов](../sdk-v1/contracts.md), [backlog](../../docs/status-backlog.md).
Источник: api-remap-1.2-doc revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.

## User Stories

1. Как разработчик решения, я могу использовать «Отдельные операции позиций» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
2. Как разработчик решения, я могу использовать «Шаблон документа PUT entity/{type}/new» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
3. Как разработчик решения, я могу использовать «Автозаполнение wizard/demand и wizard/salesreturn» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
4. Как разработчик решения, я могу использовать «Удаление в корзину» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
5. Как разработчик решения, я могу использовать «Запись связанных документов» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
6. Как разработчик решения, я могу использовать «Оставшиеся поля позиций и transport-поля demand» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
7. Как разработчик решения, я могу использовать «agentAccount/organizationAccount у cashin/cashout» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
8. Как разработчик решения, я могу использовать «InternalOrder + positions» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
9. Как разработчик решения, я могу использовать «Inventory + positions» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
10. Как разработчик решения, я могу использовать «FactureIn» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
11. Как разработчик решения, я могу использовать «FactureOut» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
12. Как разработчик решения, я могу использовать «CommissionReportIn + 2 вида позиций» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
13. Как разработчик решения, я могу использовать «CommissionReportOut + positions» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
14. Как разработчик решения, я могу использовать «CounterpartyAdjustment» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
15. Как разработчик решения, я могу использовать «PriceList + positions» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
16. Как разработчик решения, я могу использовать «Prepayment (read + Delete)» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
17. Как разработчик решения, я могу использовать «PrepaymentReturn (read + Delete)» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
18. Как разработчик решения, я могу использовать «RetailDrawerCashIn» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
19. Как разработчик решения, я могу использовать «RetailDrawerCashOut» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
20. Как разработчик решения, я могу использовать «Payroll + positions» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.

## Implementation Decisions и scope tickets

| ID | Ticket | Blocked by |
| --- | --- | --- |
| O01 | [Отдельные операции позиций](issues/01-position-mutations.md) | нет |
| O02 | [Шаблон документа PUT entity/{type}/new](issues/02-document-template.md) | нет |
| O03 | [Автозаполнение wizard/demand и wizard/salesreturn](issues/03-wizard.md) | нет |
| O04 | [Удаление в корзину](issues/04-trash.md) | нет |
| O05 | [Запись связанных документов](issues/05-linked-documents.md) | нет |
| O06 | [Оставшиеся поля позиций и transport-поля demand](issues/06-position-fields.md) | нет |
| O07 | [agentAccount/organizationAccount у cashin/cashout](issues/07-cash-accounts.md) | нет |
| O08 | [InternalOrder + positions](issues/08-internalorder.md) | O01, O02, O04 |
| O09 | [Inventory + positions](issues/09-inventory.md) | O01, O02, O04 |
| O10 | [FactureIn](issues/10-facturein.md) | O02, O04 |
| O11 | [FactureOut](issues/11-factureout.md) | O02, O04 |
| O12 | [CommissionReportIn + 2 вида позиций](issues/12-commissionreportin.md) | O01, O04 |
| O13 | [CommissionReportOut + positions](issues/13-commissionreportout.md) | O01, O04 |
| O14 | [CounterpartyAdjustment](issues/14-counterpartyadjustment.md) | O04 |
| O15 | [PriceList + positions](issues/15-pricelist.md) | O01, O04 |
| O16 | [Prepayment (read + Delete)](issues/16-prepayment.md) | O04 |
| O17 | [PrepaymentReturn (read + Delete)](issues/17-prepaymentreturn.md) | O04 |
| O18 | [RetailDrawerCashIn](issues/18-retaildrawercashin.md) | O02, O04 |
| O19 | [RetailDrawerCashOut](issues/19-retaildrawercashout.md) | O02, O04 |
| O20 | [Payroll + positions](issues/20-payroll.md) | O01, O04 |

### O01: Отдельные операции позиций

Generic GetPosition/CreatePositions/UpdatePosition/DeletePosition/DeletePositions в `positions`; typed-методы во всех 13 документах с позициями.

Проверки и особенности: ≤1000 позиций в Create; ID документа и позиции проверяются до HTTP; POST массивом возвращает массив; никаких повторов записи.

Источники: [customerorder](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_customerOrder.md), [common](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_common_info.md).

### O02: Шаблон документа PUT entity/{type}/new

New(ctx, base) у существующих документов с шаблоном: cashin, cashout, demand, invoicein, invoiceout, loss, enter, move, paymentin, paymentout, purchaseorder, purchasereturn, retaildemand, retailsalesreturn, salesreturn, supply. Новые документы добавляют New сами.

Проверки и особенности: Пустое тело и тело с основанием; ответ без id; ничего не сохраняется. enter — по общей таблице шаблонов (основание inventory, первые 500 позиций), отдельного раздела в _enter.md нет.

Источники: [common](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_common_info.md), [demand](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_demand.md).

### O03: Автозаполнение wizard/demand и wizard/salesreturn

Запросы автозаполнения: цены, скидки, НДС, себестоимость (POST wizard/{type}).

Проверки и особенности: Ответ — заполненные позиции без сохранения; типы запроса разделены.

Источники: [autofill](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_autofill.md).

### O04: Удаление в корзину

Trash(ctx, id) — POST entity/{type}/{id}/trash у существующих документов (общий раздел источника). Новые документы добавляют Trash в своих tickets.

Проверки и особенности: Отличие от Delete описано в guide; ответ без тела.

Источники: [common](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_common_info.md).

### O05: Запись связанных документов

Запись payments, supplies, demands, factureOut, returns, productionTasks, invoicesOut там, где источник разрешает запись (сейчас read-only).

Проверки и особенности: Только поля без readonly-маркера; противоречия таблицы и примеров — в Answer.

Источники: [common](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_common_info.md).

### O06: Оставшиеся поля позиций и transport-поля demand

pack, things, slot, serial/tracking у invoicein, invoiceout, move; transport-поля и shippingDocumentNumber/Date у demand.

Проверки и особенности: Общая форма positions.Position расширяется, не копируется.

Источники: [demand](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_demand.md), [invoicein](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_invoice_in.md).

### O07: agentAccount/organizationAccount у cashin/cashout

Решить расхождение: текст Update требует счета, таблица их не содержит.

Проверки и особенности: живая проверка на `admin@documan` (`~/ETEEET/moysklad/AGENTS.md:16`) — принимает ли Update счета.

Источники: [cashin](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_cashin.md), [cashout](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_cashout.md).

### O08: InternalOrder + positions

Внутренний заказ: Get/List/Create/Update/Delete, массовые, New, Trash, позиции.

Проверки и особенности: moves/purchaseOrders — read-only ссылки.

Источники: [internalorder](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_internalOrder.md).

### O09: Inventory + positions

Инвентаризация: Get/List/Create/Update/Delete, массовые, New, Trash, позиции.

Проверки и особенности: calculatedQuantity/correctionAmount — по таблице; enter/loss — read-only.

Источники: [inventory](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_inventory.md).

### O10: FactureIn

Счёт-фактура полученный: Get/List/Create/Update/Delete, массовые, New, Trash.

Проверки и особенности: Без позиций; supplies/payments — ссылки; incomingNumber/Date.

Источники: [facturein](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_facturein.md).

### O11: FactureOut

Счёт-фактура выданный: Get/List/Create/Update/Delete, массовые, New, Trash.

Проверки и особенности: Без позиций; demands/returns/payments — ссылки; stateContractId.

Источники: [factureout](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_factureout.md).

### O12: CommissionReportIn + 2 вида позиций

Полученный отчёт комиссионера: Get/List/Create/Update/Delete, массовые, Trash; позиции «реализовано» и «возврат» — два отдельных набора методов.

Проверки и особенности: commissionPeriodStart/End; rewardPercent — процент; returnToCommissionerPositions.

Источники: [commissionreportin](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_commissionreportin.md).

### O13: CommissionReportOut + positions

Выданный отчёт комиссионера: Get/List/Create/Update/Delete, массовые, Trash, позиции.

Проверки и особенности: rewardType/rewardPercent; contract обязателен.

Источники: [commissionreportout](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_commissionreportout.md).

### O14: CounterpartyAdjustment

Корректировка взаиморасчётов: Get/List/Create/Update/Delete, массовые, Trash.

Проверки и особенности: sum — копейки, знак суммы — по источнику.

Источники: [counterpartyadjustment](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_counterpartyadjustment.md).

### O15: PriceList + positions

Прайс-лист: Get/List/Create/Update/Delete, массовые, Trash, позиции.

Проверки и особенности: columns (name, percentageDiscount) и cells позиций — собственная форма позиции.

Источники: [pricelist](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_pricelist.md).

### O16: Prepayment (read + Delete)

Предоплата: Get/List/Delete, Trash, GetPositions, GetPosition. Create/Update в источнике нет.

Проверки и особенности: Суммы cashSum/noCashSum/qrSum — копейки; retailShift — ссылка.

Источники: [prepayment](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_prepayment.md).

### O17: PrepaymentReturn (read + Delete)

Возврат предоплаты: Get/List/Delete, Trash, GetPositions/GetPosition.

Проверки и особенности: Без записи; prepayment — ссылка.

Источники: [prepaymentreturn](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_prepayment_return.md).

### O18: RetailDrawerCashIn

Внесение денег: Get/List/Create/Update/Delete, массовые, New, Trash.

Проверки и особенности: retailShift обязателен; sum — копейки.

Источники: [retaildrawercashin](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_retaildrawercashin.md).

### O19: RetailDrawerCashOut

Выплата денег: Get/List/Create/Update/Delete, массовые, New, Trash.

Проверки и особенности: retailShift обязателен; sum — копейки.

Источники: [retaildrawercashout](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_retaildrawercashout.md).

### O20: Payroll + positions

Начисление зарплаты: Get/List/Create/Update/Delete, массовые, Trash, позиции.

Проверки и особенности: Позиции — employee/sum/type, собственная форма; связь с paymentout.

Источники: [payroll](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_payroll.md).

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
  завершающим commit волны.

## Open Decisions

- O01 вводит мутации позиций в `positions` и typed-методы во всех существующих документах; новые документы этой и других захваток добавляют их сразу.
- Шаблон `PUT entity/{type}/new` — метод New(ctx, base) у документов, где источник его описывает; возвращает модель документа без ID.
- O04 добавляет Trash в существующие документы; каждый новый документный ticket (O08–O20, G12/G13, M05–M08) включает Trash сам и зависит от O04.
- Файлы документов покрывает пакет `files`; отдельного ticket нет.
- O07 (счета у cashin/cashout) — расхождение источников: needs-info до решения.

## Out of Scope

Expand; печать и публикации (захватка 5); производственные документы (захватка 4). Все не перечисленные в tickets ресурсы/операции.
Этот DRAFT не является командой реализации или публикации.
