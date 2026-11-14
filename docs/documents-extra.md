# Документы: позиции, шаблоны, автозаполнение, корзина и связи

Источник контрактов — `api-remap-1.2-doc` revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
Все методы принимают `opts ...transport.RequestOption` ([ADR 0002](adr/0002-baseline.md)).
Деньги — float64 в копейках; ссылки — `model.Reference` без expand. Записи не повторяются.

## Отдельные операции позиций

У customerorder, demand, enter, invoicein, invoiceout, loss, move, purchaseorder, purchasereturn,
retaildemand, retailsalesreturn, salesreturn, supply и internalorder (общая реализация —
`positions.Service[T]`). Inventory, commissionreportin/out, payroll и pricelist — те же методы со своей
формой строки (у commissionreportin ещё набор `…ReturnPositions`), prepayment и prepaymentreturn —
только чтение строк:

| Метод | Маршрут |
| --- | --- |
| `GetPosition(ctx, id, positionID)` | `GET entity/{type}/{id}/positions/{positionID}` |
| `CreatePositions(ctx, id, []PositionRequest)` — ответ по правилам batch | `POST …/positions` массивом |
| `UpdatePosition(ctx, id, positionID, PositionRequest)` — пустой запрос отклоняется | `PUT …/positions/{positionID}` |
| `DeletePosition(ctx, id, positionID)` | `DELETE …/positions/{positionID}` |
| `DeleteManyPositions(ctx, id, []positionID)` — 1–1000, ответ без тела | `POST …/positions/delete` |

meta.type строк — `{тип}position`, кроме invoicein/invoiceout (`invoiceposition`), retaildemand
(`demandposition`), retailsalesreturn (`salesreturnposition`) и pricelist (`pricelistrow`): так в
примерах источника. Список
строк — прежний `GetPositions`.

## Шаблон документа — `New(ctx, basis)`

`PUT entity/{type}/new`: без основания (`basis == nil`) — запрос без тела, с основанием — объект
`{"customerOrder": {meta}}` и т.п. Документ не создаётся; ответ — модель чтения документа (без id),
а не `CreateRequest`: прямого преобразования нет, нужные поля переносятся в `CreateRequest` вручную.
Есть у cashin, cashout, demand, enter, invoicein, invoiceout, loss, move, paymentin, paymentout,
purchaseorder, purchasereturn, retaildemand, retailsalesreturn, salesreturn, supply, а также
internalorder, inventory, facturein, factureout, retaildrawercashin, retaildrawercashout, processing и
processingorder. Какие основания допустимы для типа, проверяет сервер.

## Автозаполнение — пакет `wizard`, `c.Wizard()`

`Fill(ctx, document, actions, Template)` — `POST wizard/{тип}?action=evaluate_price,…`. Документ не
создаётся. Действия: `EvaluatePrice`, `EvaluateVAT`, `EvaluateDiscount`, `EvaluateCost`; нужно хотя
бы одно, запятые внутри значения отклоняются до запроса. `Template` — organization, agent, store,
vatEnabled, vatIncluded, rate, moment и positions (assortment, quantity, price, discount, vat,
vatEnabled); `Result` — те же поля и рассчитанные discountedPrice, sum, cost строк, исходный JSON —
`Raw`.

Ticket называет demand и salesreturn, источник — 17 типов документов и начисления зарплаты, поэтому
тип передаётся строкой и не ограничивается списком. Какие действия поддерживает тип (скидки — не у
inventory, move, internalorder, enter, loss, payroll; цены — не у inventory; себестоимость — только у
возврата без основания), проверяет сервер.

## Корзина — `Trash(ctx, id)`

`POST entity/{type}/{id}/trash`, ответ 200 или 204 без тела. У cashin, cashout, customerorder,
demand, enter, invoicein, invoiceout, loss, move, paymentin, paymentout, purchaseorder,
purchasereturn, retaildemand, retailsalesreturn, retailshift, salesreturn, supply и у всех новых
документов (internalorder, inventory, facturein, factureout, commissionreportin, commissionreportout,
counterpartyadjustment, pricelist, prepayment, prepaymentreturn, retaildrawercashin,
retaildrawercashout, payroll), а также у документов производства (processing, processingorder,
productiontask, productionstagecompletion). Восстановление из
корзины API не предоставляет; `Delete` удаляет безвозвратно.

## Связи документов

Общие сведения источника: поля раздела «Связи с другими документами» пишутся — в поле передаётся meta
привязываемого документа (коллекция — массив meta, можно несколько). Поэтому они есть в Create и
Update как `model.Field[model.Reference]` или `model.Field[[]model.Reference]`:

| Документ | Поля записи |
| --- | --- |
| customerorder | purchaseOrders, demands, payments, productionTasks, invoicesOut, moves, prepayments |
| demand | customerOrder, factureOut, returns, payments, productionTasks, invoicesOut |
| enter, loss | inventory (у loss и salesReturn) |
| invoicein, invoiceout | purchaseOrder и supplies; customerOrder и demands |
| move | customerOrder, internalOrder, productionTasks |
| purchaseorder | internalOrder, customerOrders, invoicesIn, productionTasks, payments, supplies |
| purchasereturn, salesreturn | factureIn/factureOut, payments; losses — у salesreturn |
| supply | purchaseOrder, factureIn, invoicesIn, payments, productionTask, returns |
| paymentin/out, cashin/out | operations, factureOut/factureIn (раньше) |
| internalorder | moves, purchaseOrders, productionTasks |
| inventory | enters, losses |
| facturein; factureout | supplies, payments; demands, returns, payments |
| commissionreportin/out | payments |

Только для чтения по маркеру основной таблицы: payments у invoicein и invoiceout, demand и supply у
move. supply у purchasereturn и demand у salesreturn задаются только в Create (основание). payedSum
в разделе связей — сумма, только чтение. Пустой массив уходит как `[]`, null — как `null`; что
делает с ними сервер (отвязка), источник не описывает. Совместимость типа привязываемого документа
и права проверяет сервер.

## Поля строк и транспортные поля demand

Общая форма `positions.Position`/`PositionRequest` покрывает все поля таблиц позиций источника у
документов SDK: у invoicein и invoiceout — pack; у move — pack, things, sourceSlot, targetSlot; у
supply, demand, purchasereturn — pack, slot, things и коды маркировки по таблице. serial-полей у
строк этих документов источник не даёт.

Demand читает и пишет транспортные поля («Другие поля»): cargoName, carrier, consignee(Branch),
consignor(Branch), goodPackQuantity, shippingInstructions, stateContractId, transportFacility(Number),
shippingDocumentNumber, shippingDocumentDate ([demand](demand.md)).

## Маркер «Только для чтения» против примеров

По общему правилу SDK маркер таблицы важнее примеров; исключения — только с решением пользователя.
Утверждённые исключения (поле пишется только в тех методах, где его передаёт пример источника):

| Документ | Поле | Create | Update |
| --- | --- | --- | --- |
| internalorder | `positions` | да | да |
| commissionreportout | `positions`, `attributes` | да | да |
| commissionreportout | `shared` | да | нет |
| counterpartyadjustment | `sum` | да | да |
| retaildrawercashin, retaildrawercashout | `sum` | да | да |
| pricelist | `priceType` | да | нет |

Вне этого списка правило действует: например, `shared` и `syncId` внутреннего заказа, `sum` и
`vatSum` отчётов комиссионера и внутреннего заказа не пишутся, хотя примеры их передают.
Справочники: исключение `slot.zone` — [dictionaries](dictionaries.md).

Связи, которые ticket называл «только для чтения», а таблица не помечает (moves/purchaseOrders
внутреннего заказа, enters/losses инвентаризации), пишутся по общему правилу связей. У payroll
`payments` только читаются: источник помечает допустимые типы связи readonly.

## Допущения источника

- payments в разделе «Связи» пишутся по общему правилу, хотя раздел «Привязка платежей» общих
  сведений описывает привязку только через `operations` платежа; что сервер делает с payments в
  теле документа, источник не подтверждает.
- Копейки: таблицы не называют единицу у facturein.sum, factureout.vatSum, retaildrawercashout.sum
  («в установленной валюте») и reward строк отчётов комиссионера — SDK считает их копейками, как все
  суммы; это допущение, а не контракт источника.
