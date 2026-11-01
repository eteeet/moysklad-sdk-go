# Общий модуль positions

`positions` не принадлежит Supply или другому документу. Он предоставляет базовые
JSON-формы позиций и generic reader одной страницы; новые формы могут использовать
тот же reader без изменения SDK. Это не полное покрытие всех полей всех документов.

## JSON-формы

| Тип | Поля чтения |
| --- | --- |
| `Core` | Meta, ID, AccountID |
| `Item`, встраивает Core | Assortment model.Reference, Quantity float64 |
| `Position`, встраивает Item | Price/Discount/Overhead/Cost float64, Shipped/InTransit float64, VAT int, VATEnabled bool, SourceSlot/TargetSlot *model.Reference |
| `Position`, поля D01 | Reserve float64, TaxSystem/Reason string, Slot/Country *model.Reference, Pack *Pack, GTD *GTD, Things []string, TrackingCodes []TrackingCode, TrackingCodes1162 []TrackingCode1162, Declaration []Declaration, Stock *Stock |

Встраивание даёт плоский JSON. Неизвестные поля игнорируются, отсутствующие значения
читаются как нулевые; ID и обязательность полей не проверяются самой JSON-моделью.
Quantity/Price могут быть дробными. Discount — процент скидки или наценки:
отрицательное значение задаёт наценку (например, -10 — наценка 10%), это не деньги.
VAT=0 с VATEnabled=false означает «без НДС», с true — ставку 0%, если документ
поддерживает эти поля.
У move значения `vat=0` и `vatEnabled=false` при чтении означают отсутствие поля,
а не «без НДС». Деньги остаются в копейках, без арифметики SDK.

`PositionRequest` встраивает `ItemRequest` и использует Field с omitzero для Price,
Discount, VAT, VATEnabled, Cost, InTransit, SourceSlot и TargetSlot. CoreRequest/ItemRequest
— общие части.
В request нет Meta/AccountID. ID разрешён только как selector существующей строки
embedded Update, не как назначение ID новой позиции. Sparse request не достраивается
из старого документа.

Строки всех товарных документов — одна форма `positions.Position`/`PositionRequest`,
включая move и purchasereturn. Различие набора полей не требует отдельного типа:
отсутствующее поле чтения даёт zero value, незаданный Field не отправляется.
Допустимость заданных полей и бизнес-ограничения проверяет сервер; JSON-тип не обещает,
что любое поле можно записать в любой документ. В частности, ограничения возврата
по основанию сохраняются, а коммерческие поля для move обычно не задают.
Отдельный тип строки вводится только при реальном конфликте: одно JSON-поле имеет
разный тип или смысл в разных документах. Core/Item — части формы, не варианты по
документам.

Overhead — read-only float64 в копейках, включая API Int: принимаются 0, 0.0 и дроби;
в PositionRequest его нет. Точность денежных integer за пределами 2^53 не гарантируется.
SourceSlot/TargetSlot — ссылки без expand; отсутствие или null в ответе дают nil.
При кодировании модели чтения nil-ссылки пропускаются через omitzero.
Кодирование Null в request не обещает допустимость очистки ячеек сервером.

Cost читается и записывается через Field для услуг Demand и возвратов без основания; Shipped
только читается, InTransit читается и записывается через Field для PurchaseOrder.
Поля строк заказа покупателя, оприходования, списания, возвратов и розницы (D01): reserve,
taxSystem, reason, slot, pack (`PackRequest{ID}`), country, gtd (`GTD{Name}`), things,
trackingCodes пишутся через Field; trackingCodes_1162, declaration и stock только читаются.
Таблица enter называет things Object(String), остальные — Array(String): SDK использует массив.
Stock.Cost — указатель: у возвратов без основания сервер его не возвращает. Прочие
не перечисленные поля не считаются покрытыми. Generic reader остаётся T any: это не закрытый
union
и не требование реализовать каждую будущую строку через наследование Core/Item.

## Одна страница

`NewService[T any](client *transport.Client, documentType string)` возвращает
`(*Service[T], error)`. Nil client отклоняется. documentType — открытый код одного
сегмента `[a-z][a-z0-9]*`: нет allowlist ресурсов и скрытой нормализации.
`PageOptions.Fields` — одно значение и только документированное для кода: `stock` у demand,
customerorder, invoiceout, purchaseorder, invoicein, supply, purchasereturn, salesreturn,
retaildemand, retailsalesreturn; `declaration` у loss и retaildemand. Иное — ошибка до HTTP.
Проверка синтаксиса не гарантирует поддержку endpoint сервером.

```go
reader, err := positions.NewService[positions.Position](client, "supply")
if err != nil {
	return err
}
page, err := reader.GetPage(ctx, documentID, positions.PageOptions{Limit: 1000, Offset: 0})
if err != nil {
	return err
}
// page.Meta и page.Rows — данные ответа, SDK следующую страницу не запрашивает.
```

Этот generic вход требует выбрать подходящий T. Типизированные entity.GetPositions
с фиксированными ресурсом и Position реализованы для одиннадцати товарных документов.
Они используют общий reader без expand; [контракты документов](documents.md).

GetPage делает один GET `entity/{documentType}/{documentID}/positions`, без тела,
с явными limit/offset. documentID проверяется по форме UUID без ограничения версии.
Limit=0 означает 1000; 1…1000 допустимы, остальные значения — ошибка до HTTP.
Offset>=0; отрицательный offset отклоняется.

Возвращается `*model.Collection[T]`. Успех — HTTP 200, непустой meta.href и non-null
rows-array; [] — успешная пустая страница. Нет требования ID в T или top-level ID.
Context, metadata, href, строки и серверные числа страницы не пересчитываются;
nextHref/previousHref не запускают скрытый запрос. Размер страницы не сравнивается
с total size, offset или len(rows).

При ошибке результат nil; APIError распознаётся через errors.As, исходные auth/transport/
context errors — errors.Is. Нет retry, auto-pagination и перехода по ссылкам.
Транспорт, заголовки и gzip — существующий transport.Client.

## Embedded запись

Для поля документа используется `model.Field[[]PositionRequest]` с omitzero,
а не response-объект Collection. Это самостоятельный JSON-контракт, используемый
HTTP Create/Update всех семи товарных документов.

- Unset: поле не отправляется, при Update существующие строки остаются.
- Value пустого/nil slice: передаётся [].
- Value непустого slice: передаётся весь новый набор. При Update отсутствующие старые
  строки удаляются, matching ID обновляются, новые строки добавляются сервером.
- Null: передаётся null; он не объявляется эквивалентом [] или допустимой очисткой.

Embedded write имеет предел API 1000 строк. SDK не должен читать/сливать старые строки,
делить запрос или повторно запрашивать ответ для заполнения Rows. Требования и допустимость
конкретных полей/количеств определяет документ; модель JSON не является бизнес-валидатором.

Обычный Get заголовка документа возвращает positions metadata; наличие meta не означает,
что Rows прочитаны. Для строк нужна явная страница. Этот модуль не добавляет List
документов,
отдельные HTTP Get/Create/Update/Delete строк, batch, templates, walker, stock enrichment
или expand. Expand отложен до конца разработки всего остального SDK, если понадобится.

Первичные источники на revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`:
[общие правила](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1022-L1063),
[страница позиций](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L2468-L2525),
[позиции перемещения](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_move.md#L62-L89).
