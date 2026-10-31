Status: needs-triage
Stage: DRAFT

# Документы SDK: карта spec и tickets

Подготовлено 2026-11-01: 7 spec, 17 tickets. Первая волна C01/C02/F01/P01/P02 утверждена
и разрешена к реализации после коммита накопленных изменений. Остальные 12 tickets
остаются DRAFT/needs-triage. Полные entity-контракты и P03 не разрешены этим запуском.
В колонке зависимостей указаны технические blockers; готовность не заменяет triage.

Следующий шаг — [краткий пакет рекомендованных контрактов](approval.md).
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

| ID | Ticket | Blocked by |
| --- | --- | --- |
| C01 | [DocumentRate / DocumentRateRequest](../document-common/issues/01-rate.md) | нет |
| C02 | [DocumentOverhead / DocumentOverheadRequest](../document-common/issues/02-overhead.md) | нет |
| P01 | [Core/Item/Commercial и embedded JSON](../document-positions/issues/01-json-contracts.md) | нет |
| P02 | [Generic reader одной страницы GetPage](../document-positions/issues/02-get-page.md) | нет |
| P03 | [Movement / MovementRequest](../document-positions/issues/03-movement-contracts.md) | P01 |
| F01 | [FinanceOperation / FinanceOperationRequest](../financial-documents/issues/01-finance-operation.md) | нет |
| F02 | [PaymentIn Get/Create/Update](../financial-documents/issues/02-paymentin.md) | C01, F01 |
| F03 | [PaymentOut Get/Create/Update](../financial-documents/issues/03-paymentout.md) | C01, F01 |
| F04 | [CashIn Get/Create/Update](../financial-documents/issues/04-cashin.md) | C01, F01 |
| F05 | [CashOut Get/Create/Update](../financial-documents/issues/05-cashout.md) | C01, F01 |
| I01 | [InvoiceIn Get/Create/Update + GetPositions](../invoice-documents/issues/01-invoicein.md) | C01, P01, P02 |
| I02 | [InvoiceOut Get/Create/Update + GetPositions](../invoice-documents/issues/02-invoiceout.md) | C01, P01, P02 |
| S01 | [Supply Get/Create/Update + GetPositions](../stock-documents/issues/01-supply.md) | C01, C02, P01, P02 |
| S02 | [Demand Get/Create/Update + GetPositions](../stock-documents/issues/02-demand.md) | C01, C02, P01, P02 |
| S03 | [PurchaseOrder Get/Create/Update + GetPositions](../stock-documents/issues/03-purchaseorder.md) | C01, P01, P02 |
| M01 | [Move Get/Create/Update + GetPositions](../stock-move/issues/01-move.md) | C01, C02, P01, P02, P03 |
| R01 | [PurchaseReturn Get/Create/Update + GetPositions](../purchase-return/issues/01-purchasereturn.md) | C01, P01, P02 |

### Параллельность после утверждения

Все 17 одновременно запускать нельзя: есть зависимости общих публичных типов.
Ниже технически допустимый порядок, не команда запуска всех семейств:

1. Независимо: C01, C02, P01, P02 и F01.
2. По готовности соответствующих зависимостей: P03, финансовые F02–F05,
   счета I01–I02, складские S01–S03 и возврат R01. Между entity tickets зависимостей нет.
3. M01 — после P03 и остальных его зависимостей. Не нужно ждать все tickets шага 2.

P02 не требует P01: reader принимает T any и уже существующий model.Collection[T].
P02 — обязательная зависимость позиционных entity tickets: они включают типизированный GetPositions. Финансовая группа не зависит от P01/P02/P03.
Reference не требует готовых services связанных сущностей, счетов, складов и ячеек.

Для первой волны разрешён [трёхпоточный запуск](wave-1.md); три model tickets выполняются
одним implementer последовательно, отдельными коммитами в одном worktree. P01/P02 —
в отдельных worktree с непересекающимися файлами. Для остальных согласуемых групп —
отдельные branch/worktree и implementer на ticket,
перенос resolved dependencies до старта зависимого ticket. Общие coverage/README
сводит исполнитель интеграции, чтобы параллельные изменения не перетирали друг друга.
Общие JSON-типы и HTTP reader — самостоятельные проверяемые срезы, не этапы-заглушки.

## Решения для triage

Это список предметных решений, а не предложение снова обсуждать состав групп.
Нормативные детали и предложения находятся в соответствующих spec.

1. Общий positions: утвердить exports Core/Item/Commercial/Movement, отдельные requests,
   `NewService[T any](client, documentType)`, GetPage, PageOptions с 0→1000,
   синтаксис documentType и минимальную проверку response envelope.
2. Общие типы: утвердить DocumentRate/Request с *Reference и *float64 при чтении,
   без default курса; DocumentOverhead/Request с sum float64 и открытым distribution.
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
5. Финансы: утвердить FinanceOperation с meta/linkedSum, float64 без арифметики;
   не обещать full replacement/очистку operations по одной лишь записи [].
   CashIn/CashOut: текст Update упоминает account-поля, отсутствующие в таблицах;
   proposed не добавлять их без подтверждения. Agent PaymentOut polymorphic,
   аналогичную поддержку у CashOut не выводить автоматически.
6. Demand: предложен Store Field[Reference] в Create/Update без local required-validator;
   сервер проверяет применимость. Утвердить это и reuse Address/Request
   без обещания FIASCodeRU; reserve и транспортные поля пока TODO.
7. PurchaseOrder: предложен документный PositionRequest с InTransit Field[float64]
   по официальным примерам записи; shipped и waitSum proposed read-only в scope SDK.
   Нельзя приписывать этим двум полям отсутствующий readonly marker таблиц.
8. Move: утвердить отдельный Movement без discount/VAT, ссылки на ячейки и scope связей.
9. PurchaseReturn: разрешить расхождения embedded price/vat/discount, запреты Update
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
  review; после интеграции — общий прогон. Сейчас Go-код и тесты реализации не менялись.
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
