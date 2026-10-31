# Общий модуль positions

`positions` не принадлежит Supply или другому документу. Он предоставляет базовые
JSON-формы позиций и generic reader одной страницы; новые формы могут использовать
тот же reader без изменения SDK. Это не полное покрытие всех полей всех документов.

## JSON-формы

| Тип | Поля чтения |
| --- | --- |
| `Core` | Meta, ID, AccountID |
| `Item`, встраивает Core | Assortment model.Reference, Quantity float64 |
| `Commercial`, встраивает Item | Price/Discount float64, VAT int, VATEnabled bool |

Встраивание даёт плоский JSON. Неизвестные поля игнорируются, отсутствующие значения
читаются как нулевые; ID и обязательность полей не проверяются самой JSON-моделью.
Quantity/Price могут быть дробными, Discount — отрицательным. VAT=0 с VATEnabled=false
означает «без НДС», с true — ставку 0%. Деньги остаются в копейках, без арифметики SDK.

Отдельные `CoreRequest`, `ItemRequest`, `CommercialRequest` используют Field с omitzero.
В request нет Meta/AccountID. ID разрешён только как selector существующей строки
embedded Update, не как назначение ID новой позиции. Sparse request не достраивается
из старого документа.

Overhead/cost, shipped/inTransit, ячейки, pack/serial/tracking/GTD и другие специфичные
поля не добавлены в Commercial автоматически. Movement и типы отдельных документов
разрабатываются своими tickets; общий reader может использовать собственный T.

## Одна страница

`NewService[T any](client *transport.Client, documentType string)` возвращает
`(*Service[T], error)`. Nil client отклоняется. documentType — открытый код одного
сегмента `[a-z][a-z0-9]*`: нет allowlist ресурсов и скрытой нормализации.
Проверка синтаксиса не гарантирует поддержку endpoint сервером.

```go
reader, err := positions.NewService[positions.Commercial](client, "supply")
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
с фиксированными ресурсом и Position предусмотрены в entity spec, но пока не реализованы.
Их появление не меняет общий reader и не требует expand.

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
а не response-объект Collection. В текущем модуле это самостоятельный JSON-контракт;
HTTP Create/Update документа пока не реализованы.

- Unset: поле не отправляется, при Update существующие строки остаются.
- Value пустого/nil slice: передаётся [].
- Value непустого slice: передаётся весь новый набор. При Update отсутствующие старые
  строки удаляются, matching ID обновляются, новые строки добавляются сервером.
- Null: передаётся null; он не объявляется эквивалентом [] или допустимой очисткой.

Embedded write имеет предел API 1000 строк. SDK не должен читать/сливать старые строки,
делить запрос или повторно запрашивать ответ для заполнения Rows. Требования и допустимость
конкретных полей/количеств определяет документ; модель JSON не является бизнес-валидатором.

Обычный Get заголовка документа возвращает positions metadata; наличие meta не означает,
что Rows прочитаны. Для строк нужна явная страница. Этот модуль не добавляет List документов,
отдельные HTTP Get/Create/Update/Delete строк, batch, templates, walker, stock enrichment
или expand. Expand отложен до конца разработки всего остального SDK, если понадобится.

Первичные источники на revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`:
[общие правила](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1022-L1063),
[страница позиций](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L2468-L2525).
