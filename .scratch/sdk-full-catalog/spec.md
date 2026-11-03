Status: needs-triage
Stage: DRAFT

# 2. Товары, скидки и маркировка

## Problem Statement

Довести модель Product до полной таблицы источника, добавить изображения, карточки контента, все виды скидок и сущности маркировки.

## Solution

Завершённые вертикальные срезы; каждый ticket включает models, methods, тесты и собственное
описание полей. Матрица ниже определяет scope этой захватки.
[Общая карта](../sdk-full/map.md), [нормативные ограничения v1](../sdk-v1/spec.md),
[каталог JSON-контрактов](../sdk-v1/contracts.md), [backlog](../../docs/status-backlog.md).
Источник: api-remap-1.2-doc revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.

## User Stories

1. Как разработчик решения, я могу использовать «Product: штрихкоды и упаковки» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
2. Как разработчик решения, я могу использовать «Product: alcohol, отслеживание и прочие поля» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
3. Как разработчик решения, я могу использовать «Неснижаемый остаток, storebalances и fields=minimumStock» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
4. Как разработчик решения, я могу использовать «Images: product/variant/bundle» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
5. Как разработчик решения, я могу использовать «ContentCard Get/List/Create/Update + images» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
6. Как разработчик решения, я могу использовать «Discount: общий список и округление копеек» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
7. Как разработчик решения, я могу использовать «AccumulationDiscount CRUD» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
8. Как разработчик решения, я могу использовать «PersonalDiscount CRUD» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
9. Как разработчик решения, я могу использовать «SpecialPriceDiscount CRUD» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
10. Как разработчик решения, я могу использовать «TrackingCode List + массовые» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
11. Как разработчик решения, я могу использовать «AggregatePack + positions» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
12. Как разработчик решения, я могу использовать «EmissionOrder Get/List/Create/Update + positions» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
13. Как разработчик решения, я могу использовать «RetireOrder Get/List/Create/Update + positions» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.

## Implementation Decisions и scope tickets

| ID | Ticket | Blocked by |
| --- | --- | --- |
| G01 | [Product: штрихкоды и упаковки](issues/01-barcodes-packs.md) | нет |
| G02 | [Product: alcohol, отслеживание и прочие поля](issues/02-product-attributes.md) | G01 |
| G03 | [Неснижаемый остаток, storebalances и fields=minimumStock](issues/03-minimumstock.md) | нет |
| G04 | [Images: product/variant/bundle](issues/04-images.md) | нет |
| G05 | [ContentCard Get/List/Create/Update + images](issues/05-contentcard.md) | G04 |
| G06 | [Discount: общий список и округление копеек](issues/06-discount.md) | нет |
| G07 | [AccumulationDiscount CRUD](issues/07-accumulationdiscount.md) | G06 |
| G08 | [PersonalDiscount CRUD](issues/08-personaldiscount.md) | G06 |
| G09 | [SpecialPriceDiscount CRUD](issues/09-specialpricediscount.md) | G06 |
| G10 | [TrackingCode List + массовые](issues/10-trackingcode.md) | нет |
| G11 | [AggregatePack + positions](issues/11-aggregatepack.md) | нет |
| G12 | [EmissionOrder Get/List/Create/Update + positions](issues/12-emissionorder.md) | O01, O04 |
| G13 | [RetireOrder Get/List/Create/Update + positions](issues/13-retireorder.md) | O01, O04 |

### G01: Product: штрихкоды и упаковки

Product.barcodes через существующую model.Barcode; Product.packs — собственная форма (uom, quantity, barcodes) в read/write. Variant packs (parentpack) уже есть — не трогать.

Проверки и особенности: Тип штрихкода — ключ объекта (ean13, gtin, …), неизвестный тип сохраняется; pack.uom — ссылка; packs у Service/Bundle не добавлять — в источнике их нет.

Источники: [product](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_product.md), [barcode](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/workbook/_workbook_barcode.md).

### G02: Product: alcohol, отслеживание и прочие поля

Оставшиеся поля Product по таблице: alcoholic, trackingType, ppeType, tnved, paymentItemType, weight/volume, things, minPrice, buyPrice и т.д.; закрыть TODO coverage.md.

Проверки и особенности: Readonly-маркер таблицы сильнее примера; деньги — копейки; enum открытые.

Источники: [product](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_product.md).

### G03: Неснижаемый остаток, storebalances и fields=minimumStock

Модель minimumStock; снять отказ fields=minimumStock в List Product и Variant; storebalances Product и Variant: List/Get/Create/Update/Delete/DeleteMany.

Проверки и особенности: Формы all/store, значения quantity; fields передаётся одним значением.

Источники: [product](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_product.md), [variant](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_variant.md).

### G04: Images: product/variant/bundle

Изображения: List, Add (base64), POST замены/упорядочивания списка, Delete, DeleteMany, временная ссылка, fields=downloadPermanentHref (entity/{type}/{id}/images).

Проверки и особенности: Лимит размера тела — существующий preflight; скачивание — transport.Download; filename обязателен.

Источники: [images](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_images.md), [workbook images](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/workbook/_workbook_images.md).

### G05: ContentCard Get/List/Create/Update + images

Карточка контента: Get/List/Create/Update/Delete, массовые; изображения через G04.

Проверки и особенности: Связь с assortment — ссылка; поля по таблице.

Источники: [contentcard](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_contentcard.md).

### G06: Discount: общий список и округление копеек

GET entity/discount (все типы, открытый meta.type) и roundoffdiscount: Update и массовое.

Проверки и особенности: Список возвращает разнородные объекты — общий тип с raw для неизвестных.

Источники: [discount](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_discount.md).

### G07: AccumulationDiscount CRUD

Накопительная скидка: Get/Create/Update/Delete.

Проверки и особенности: levels: amount (копейки) и discount (процент); productFolders/assortment — ссылки.

Источники: [discount](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_discount.md).

### G08: PersonalDiscount CRUD

Персональная скидка: Get/Create/Update/Delete.

Проверки и особенности: agentTags/allAgents/allProducts — явный false.

Источники: [discount](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_discount.md).

### G09: SpecialPriceDiscount CRUD

Скидка на спеццену: Get/Create/Update/Delete.

Проверки и особенности: specialPrice: value (копейки) + priceType; discount — процент; usePriceType.

Источники: [discount](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_discount.md).

### G10: TrackingCode List + массовые

Коды маркировки: List, массовое создание/обновление и удаление.

Проверки и особенности: Иерархия кодов (trackingCodes внутри) — рекурсивная модель; тип кода — открытый.

Источники: [trackingcode](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_tracking_code.md).

### G11: AggregatePack + positions

Транспортная упаковка: Get/List/Create/Update/Delete, массовые; позиции: List/Get/Create/Update/Delete, массовое удаление.

Проверки и особенности: Вложенность упаковок; positions — собственная форма, не positions.Position документов.

Источники: [aggregatepack](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_handling_unit.md).

### G12: EmissionOrder Get/List/Create/Update + positions

Заказ кодов маркировки: Get/List/Create/Update, массовые, Trash; позиции List/Get/Create/Update.

Проверки и особенности: Delete документа в источнике нет; trackingType/cisType — открытые string.

Источники: [emissionorder](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_emissionorder.md).

### G13: RetireOrder Get/List/Create/Update + positions

Вывод из оборота: Get/List/Create/Update, CreateMany/UpdateMany, Trash; позиции List/Get/Create/Update/Delete, массовое удаление.

Проверки и особенности: Допустимость Update по documentState проверяет сервер; причины вывода — открытый enum.

Источники: [retireorder](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_retireorder.md).

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

- Уточнение правила readonly (review 2a76d13): маркер запрещает менять значение поля; передача id/meta
  для адресации существующего объекта допустима, когда её задаёт метод API (PackRequest.ID,
  StoreBalanceRequest.Meta, trackingcode.Request.ID) — отдельного исключения не требует.

- Штрихкоды — существующая model.Barcode (уже у Service, Bundle, Variant). Упаковки не общая форма: Product — uom+quantity+barcodes, Variant — parentpack+barcodes (уже есть); у Service и Bundle packs нет. Consignment получает barcodes в S10.
- Изображения — отдельный nested-ресурс; поле images в моделях остаётся meta коллекции.
- emissionorder и retireorder — документы с позициями; ждут O01 (мутации позиций) и O04 (Trash).

## Out of Scope

Expand; бизнес-логика скидок (расчёт) — только сервер. Все не перечисленные в tickets ресурсы/операции.
Этот DRAFT не является командой реализации или публикации.
