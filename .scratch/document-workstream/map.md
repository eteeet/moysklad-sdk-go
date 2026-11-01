Status: resolved
Stage: IMPLEMENTED

# Документы SDK: карта spec и tickets

Все 11 оставшихся entity tickets реализованы по команде пользователя 2026-11-01.
В исходных 17 tickets: **16 resolved, 1 wontfix/superseded (P03)**.
Общая Position/PositionRequest завершена отдельно; Cost/Shipped/InTransit добавлены
соответствующими entity tickets, Cost и InTransit доступны в общем request.
[Итоговые проверки и независимые review](review-entities.md): race/vet/gofmt/diff-check
PASS, Standards 0 / Spec 0; upstream/provenance сохранены.

Актуальные [коммиты, проверки и независимые review первой волны](review-wave-1.md).
[Утверждённый scope первой волны](approval.md) сохраняет основания решений.
Повторная сверка спорных полей — [исследование triage](triage-research.md).

## Spec и состав захваток

| Захватка | Spec | Tickets |
| --- | --- | --- |
| Общие вложенные типы документов | [DocumentRate, DocumentOverhead](../document-common/spec.md) | 2 |
| Универсальный positions | [JSON shapes и generic GetPage](../document-positions/spec.md) | 3 |
| Финансовые документы | [paymentin, paymentout, cashin, cashout](../financial-documents/spec.md), [поимённый контракт](../financial-documents/contract.md) | 5 |
| Позиционные счета | [invoicein, invoiceout](../invoice-documents/spec.md) | 2 |
| Позиционные складские документы | [supply, demand, purchaseorder](../stock-documents/spec.md) | 3 |
| Перемещение | [move](../stock-move/spec.md) | 1 |
| Возврат поставщику | [purchasereturn](../purchase-return/spec.md) | 1 |

## Все tickets и зависимости

Префиксы C/P/F/I/S/M/R — только короткие обозначения этой карты, не новые номера tracker.
В самом ticket Blocked by содержит относительные пути к файлам зависимостей.

| ID | Ticket | Blocked by | Status |
| --- | --- | --- | --- |
| C01 | [DocumentRate / DocumentRateRequest](../document-common/issues/01-rate.md) | нет | resolved |
| C02 | [DocumentOverhead / DocumentOverheadRequest](../document-common/issues/02-overhead.md) | нет | resolved |
| P01 | [Core/Item и общий JSON](../document-positions/issues/01-json-contracts.md) | нет | resolved |
| P02 | [Generic reader одной страницы GetPage](../document-positions/issues/02-get-page.md) | нет | resolved |
| P03 | [Отдельный контракт заменён общей формой](../document-positions/issues/03-movement-contracts.md) | нет | wontfix / superseded |
| F01 | [FinanceOperation / FinanceOperationRequest](../financial-documents/issues/01-finance-operation.md) | нет | resolved |
| F02 | [PaymentIn Get/Create/Update](../financial-documents/issues/02-paymentin.md) | C01, F01 | resolved |
| F03 | [PaymentOut Get/Create/Update](../financial-documents/issues/03-paymentout.md) | C01, F01 | resolved |
| F04 | [CashIn Get/Create/Update](../financial-documents/issues/04-cashin.md) | C01, F01 | resolved |
| F05 | [CashOut Get/Create/Update](../financial-documents/issues/05-cashout.md) | C01, F01 | resolved |
| I01 | [InvoiceIn Get/Create/Update + GetPositions](../invoice-documents/issues/01-invoicein.md) | C01, P01, P02 | resolved |
| I02 | [InvoiceOut Get/Create/Update + GetPositions](../invoice-documents/issues/02-invoiceout.md) | C01, P01, P02 | resolved |
| S01 | [Supply Get/Create/Update + GetPositions](../stock-documents/issues/01-supply.md) | C01, C02, P01, P02 | resolved |
| S02 | [Demand Get/Create/Update + GetPositions](../stock-documents/issues/02-demand.md) | C01, C02, P01, P02 | resolved |
| S03 | [PurchaseOrder Get/Create/Update + GetPositions](../stock-documents/issues/03-purchaseorder.md) | C01, P01, P02 | resolved |
| M01 | [Move Get/Create/Update + GetPositions](../stock-move/issues/01-move.md) | C01, C02, P01, P02 | resolved |
| R01 | [PurchaseReturn Get/Create/Update + GetPositions](../purchase-return/issues/01-purchasereturn.md) | C01, P01, P02 | resolved |

### Параллельность после утверждения

Все 17 одновременно запускать нельзя: есть зависимости общих публичных типов.
Ниже технически допустимый порядок, не команда запуска всех семейств:

1. Независимо: C01, C02, P01, P02 и F01.
2. По готовности соответствующих зависимостей: финансовые F02–F05,
   счета I01–I02, складские S01–S03, M01 и возврат R01. Отдельный P03 не нужен.
   S02/S03 расширяют общий positions полями Cost и Shipped/InTransit: их изменения
   в общих файлах согласовать при параллельной интеграции, не создавать entity-owned формы.

P02 не требует P01: reader принимает T any и уже существующий model.Collection[T].
P02 — обязательная зависимость позиционных entity tickets: они включают типизированный GetPositions. Финансовая группа не зависит от positions.
Reference не требует готовых services связанных сущностей, счетов, складов и ячеек.

Для первой волны разрешён [трёхпоточный запуск](wave-1.md); три model tickets выполняются
одним implementer последовательно, отдельными коммитами в одном worktree. P01/P02 —
в отдельных worktree с непересекающимися файлами. Для остальных согласуемых групп —
отдельные branch/worktree и implementer на ticket,
перенос resolved dependencies до старта зависимого ticket. Общие coverage/README
сводит исполнитель интеграции, чтобы параллельные изменения не перетирали друг друга.
Общие JSON-типы и HTTP reader — самостоятельные проверяемые срезы, не этапы-заглушки.

## Утверждённые контракты (история triage)

Ниже сохранены решения и предложения этапа triage. Команда пользователя выполнить
все 11 tickets утвердила предложенные контракты; реализация завершена.
Уточнение S02: Cost записывается для услуг через общий PositionRequest.
Актуальные поля и границы — в docs/documents.md и guides сущностей.

1. Core/Item/Position, общий PositionRequest, NewService/GetPage/PageOptions,
   синтаксис documentType и response envelope утверждены и реализованы P01/P02.
   Различия набора полей не требуют отдельных форм; допустимость записи проверяет сервер.
   Отдельная форма возможна только при реальном конфликте типа/смысла одного JSON-поля.
2. DocumentRate/Request и DocumentOverhead/Request реализованы. Value=nil в присутствующем
   rate означает курс 1 без подстановки в JSON; распределение открытое.
   Все денежные поля, включая cost/overhead строки, читаются float64 в копейках независимо от API Int/Float; строковый overhead не является объектом расходов.
3. Общая proposed validation: required Create refs только по непустому Meta.Href,
   UUID маршрута, непустой Update и максимум 1000 embedded positions. Сочетания счетов,
   доступность ссылок и бизнес-ограничения оставлены серверу. Invoice Create дополнительно
   требует Name; у остальных документов Name необязателен. Это не перенос требования UOM.
4. SyncID предложен в Create/Update; запрет изменения после заполнения проверяет сервер.
   Подтвердить выбранный первый scope: часть связанных массивов
   только чтение; pack/serial/tracking/stock и отдельные HTTP mutations позиций — TODO.
   Сужения SDK не объявляются запретами API. Financial files исключены полностью;
   складские/возврат читают только files metadata; счета files пока не включают.
5. FinanceOperation с meta/linkedSum реализован F01; entity-контракты F02–F05 ещё не утверждены;
   не обещать full replacement/очистку operations по одной лишь записи [].
   CashIn/CashOut: текст Update упоминает account-поля, отсутствующие в таблицах;
   proposed не добавлять их без подтверждения. Agent PaymentOut polymorphic,
   аналогичную поддержку у CashOut не выводить автоматически.
6. Demand: предложен Store Field[Reference] в Create/Update без local required-validator;
   сервер проверяет применимость. Утвердить это и reuse Address/Request
   без обещания FIASCodeRU; reserve и транспортные поля пока TODO.
7. PurchaseOrder: S03 добавит InTransit Field[float64] в общий positions.PositionRequest
   по официальным примерам записи; shipped и waitSum proposed read-only в scope SDK.
   Нельзя приписывать этим двум полям отсутствующий readonly marker таблиц.
8. Move использует общие Position/PositionRequest; overhead read-only, ячейки optional.
   Неиспользуемые Field не задают; нули VAT/VATEnabled в read не означают «без НДС».
   Entity-контракт Move и scope его связей требуют отдельного утверждения.
9. PurchaseReturn использует общие Position/PositionRequest; документировать расхождения embedded price/vat/discount, запреты Update
   сторон и применение ограничений основания к возврату без supply. Quantity-only
   отдельного PUT строки не переносить на embedded Update без основания.

Ни один ticket не переводится в ready-for-agent только потому, что его зависимости готовы.
После решений обновить spec/tickets, затем получить команду реализации выбранной группы.

## Общие границы и проверки

- Expand — только в конце всего SDK, после остальных возможностей, если понадобится.
  В этих захватках его нет; обычный Get получает positions metadata, строки — GetPage.
- Встроенные positions Update полностью заменяют набор: omitted сохраняет, [] запрашивает
  пустой набор, непустой массив заменяет; null не объявляется эквивалентом []. Нет merge/GET.
- Отдельные request types, Field с omitzero, существующие Reference/DateTime/Attribute,
  настоящий transport.Client, один HTTP-запрос, фактический ответ сервера, nil при ошибке.
- HTTP-тесты по entity/uom/service_test.go: отдельные TestServiceGet/Create/Update
  с табличными сценариями; reader имеет собственный GetPage test. JSON — оба кодека,
  независимые golden, upstream fixtures без изменения и с provenance.
- После реализации каждого ticket: race/vet/gofmt/diff-check и независимые Standards/Spec
  review; после интеграции — общий прогон.
- Service отложен, variant/bundle после документов. List/delete/batch и транспортные
  расширения не добавляются автоматически. Нет .env/live API/push/PR/VPS.
- Spec/tickets и temp_plan не коммитятся без отдельной команды пользователя.

## Decisions so far

- 2026-11-01: финансовая группа исправлена: paymentin/paymentout/cashin/cashout.
  InvoiceIn/InvoiceOut остаются отдельной позиционной группой; остальные согласованные
  группы сохранены. Подготовлены DRAFT, реализация не начата.
- 2026-11-01: повторный triage уточнил DRAFT: Demand.Store optional Field,
  PurchaseOrder.InTransit writable в document-owned request, SyncID в Create/Update.
  Подготовлен [пакет для утверждения](approval.md); это рекомендации, не одобрение пользователя.
- 2026-11-01: пользователь утвердил первую волну с исправлениями: все деньги float64,
  nil Value присутствующего rate означает курс 1 без подстановки в JSON; entity GetPositions
  типизирован, зависит от P02; SyncID Update старых Counterparty/Organization — отдельный TODO.
  Разрешил сначала коммит всех накопленных изменений, затем реализацию пяти tickets.
- Первичные источники закреплены на official revision
  9a811cef3dc4a3f3a46f070e60452361a06bf7cd, master повторно подтверждён 2026-11-01.
  Исследования: [positions](../document-positions/research.md),
  [финансы](../financial-documents/research.md), [складские документы](../stock-documents/research.md).

### Завершение первой волны

Planning commit a280428 сохранён до начала Go-изменений. Реализация: C01 8611789,
C02 7e6b812, F01 02026e0, P01 20e00ad (GoDoc fix fc8e040), P02 1b28603.
Полный интеграционный race/vet/gofmt/diff-check пройден; Standards 0, Spec 0.
Документацию README/coverage и отдельные guides свёл интегратор после слияния.
Первая волна принята пользователем. Следующей командой разрешён P03; entity tickets
не запущены и требуют отдельной команды/triage.

### Пересмотр общей формы

2026-11-01 пользователь заменил отдельный P03 одной Position/PositionRequest.
[Согласованный контракт](../universal-positions/spec.md); overhead и ссылки на ячейки
перенесены в общий тип, P03 больше не блокирует M01. Cost/Shipped/InTransit добавятся
в общий тип entity tickets, а не в собственные формы документов.
Это не утверждает entity-контракты автоматически. Следующий рекомендуемый triage —
paymentin/paymentout/cashin/cashout; четыре tickets независимы после утверждения.

### Завершение entity tickets

- 2026-11-01 [F02 PaymentIn](../financial-documents/issues/02-paymentin.md): Get/Create/Update реализованы; JSON/HTTP/race/vet и Standards/Spec PASS.
- 2026-11-01 [F03 PaymentOut](../financial-documents/issues/03-paymentout.md): Get/Create/Update реализованы; JSON/HTTP/race/vet и Standards/Spec PASS.
- 2026-11-01 [F04 CashIn](../financial-documents/issues/04-cashin.md): Get/Create/Update реализованы; JSON/HTTP/race/vet и Standards/Spec PASS.
- 2026-11-01 [F05 CashOut](../financial-documents/issues/05-cashout.md): Get/Create/Update реализованы; JSON/HTTP/race/vet и Standards/Spec PASS.
- 2026-11-01 [I01 InvoiceIn](../invoice-documents/issues/01-invoicein.md): Get/Create/Update + GetPositions реализованы; JSON/HTTP/race/vet и Standards/Spec PASS.
- 2026-11-01 [I02 InvoiceOut](../invoice-documents/issues/02-invoiceout.md): Get/Create/Update + GetPositions реализованы; JSON/HTTP/race/vet и Standards/Spec PASS.
- 2026-11-01 [S01 Supply](../stock-documents/issues/01-supply.md): Get/Create/Update + GetPositions реализованы; JSON/HTTP/race/vet и Standards/Spec PASS.
- 2026-11-01 [S02 Demand](../stock-documents/issues/02-demand.md): Get/Create/Update + GetPositions реализованы; JSON/HTTP/race/vet и Standards/Spec PASS.
- 2026-11-01 [S03 PurchaseOrder](../stock-documents/issues/03-purchaseorder.md): Get/Create/Update + GetPositions реализованы; JSON/HTTP/race/vet и Standards/Spec PASS.
- 2026-11-01 [M01 Move](../stock-move/issues/01-move.md): Get/Create/Update + GetPositions реализованы; JSON/HTTP/race/vet и Standards/Spec PASS.
- 2026-11-01 [R01 PurchaseReturn](../purchase-return/issues/01-purchasereturn.md): Get/Create/Update + GetPositions реализованы; JSON/HTTP/race/vet и Standards/Spec PASS.

Коммиты/push/live API и cleanup worktree не выполнялись. Review-followups cleanup
остаётся отдельной задачей. Исходные planning/temp_plan изменения сохранены.
