# Товары: поля, остатки, изображения, скидки и маркировка

Источник контрактов — `api-remap-1.2-doc` revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
Все методы принимают `opts ...transport.RequestOption` ([ADR 0002](adr/0002-baseline.md)).
Деньги — float64 в копейках, проценты — нет; ссылки — `model.Reference` без expand. Записи не
повторяются.

## Product

- Штрихкоды — `model.Barcode`; упаковки — `Pack` (uom, quantity, не больше одного штрихкода).
  В Update barcodes и packs заменяют коллекцию целиком: упаковка с `ID` — существующая (id только
  адресует её), без — новая. Без barcodes в Create сервер создаёт один случайный EAN13.
- Поля: alcoholic, minimumBalance, partialDisposal, paymentItemType, ppeType, taxSystem, trackingType,
  tnved, things, weighed, onTap, isSerialTrackable; variantsCount, images и files — только чтение
  (изображения — пакет `images`, файлы — `files`). Модули `mod__tasnif__uz`, `mod__marking__uz`,
  `mod__nct__kz` (`model.NctKZ`). Несочетаемые признаки проверяет сервер.

## Неснижаемый остаток

- `model.MinimumStock` приходит в List с `fields=minimumStock` (Get параметра fields не принимает) и
  в ответе Create/Update, где поле передано (Product и Variant):
  ALL_WAREHOUSE_SUM и ALL_WAREHOUSE_SAME — quantity, WAREHOUSE_VARIED — метаданные storebalances.
  Источник пишет коллекцию то `storebalances`, то `storeBalances` — читаются оба ключа.
- Meta строки storebalances и id кода маркировки или упаковки в запросах только адресуют
  существующий объект: readonly-маркер таблицы запрещает менять значение, а не выбирать по нему.
- Запись — `MinimumStockRequest` в Create/Update; ALL_WAREHOUSE_SAME только читается.
- `ListStoreBalances`, `GetStoreBalance`, `CreateStoreBalance`, `UpdateStoreBalance`,
  `DeleteStoreBalance`, `DeleteManyStoreBalances` — `entity/{product,variant}/{id}/storebalances`.
  Разделов этих методов источник не даёт, только маршруты: поля — по таблице (store, quantity).

## Изображения — пакет `images`, `c.Entity().Images()`

| Метод | Маршрут |
| --- | --- |
| `List(ctx, type, id, page)`; `page.Fields = images.PermanentHref` — постоянные ссылки | `GET entity/{type}/{id}/images` |
| `Add(ctx, type, id, Upload)` — ответ: весь список | `POST …/images` |
| `Replace(ctx, type, id, []Item)` — `images.Keep`, `images.New`; не переданные удаляются | `POST …/images` массивом |
| `Delete`, `DeleteMany` | `DELETE …/images/{id}`, `POST …/images/delete` |
| `Download(ctx, image, w)` — по `meta.downloadHref`; 302 на временную ссылку идёт transport | `GET download/{id}` |

Типы владельцев — только product, bundle, variant, contentcard. Пустые filename и content
отклоняются до запроса (источник требует оба непустыми).

## Карточка контента, партия, транспортная упаковка

- `entity/contentcard`: Create требует name, description, cardContentName, assortment, salePlatform и
  непустой salesChannels; assortment после создания не меняется.
- `entity/consignment`: источник противоречит себе (таблица — label и expiryDate, метод — label и
  assortment, вводный текст — label или expiryDate); SDK требует assortment и хотя бы одно из label и
  expiryDate, остальное решает сервер. name и images — чтение.
- `entity/aggregatepack`: Create требует непустой barcodes; childrenList — только с
  `transport.WithBetaFeature(aggregatepack.BetaChildrenList)`. Позиции можно передать в теле
  Create/Update (маркера readonly нет; как Update сочетает их с прежними, источник не описывает)
  или отдельными методами — своя форма:
  `ListPositions`, `GetPosition`, `CreatePositions` (POST массивом), `UpdatePosition`,
  `DeletePosition`, `DeleteManyPositions`.

## Скидки

- `entity/discount`: `List` — все виды (общий тип по `Meta.Type`, исходный JSON строки — `Raw`),
  `UpdateRoundOff` и `UpdateMany`
  округления копеек (name, active).
- `entity/accumulationdiscount`, `personaldiscount`, `specialpricediscount`: Get/Create/Update/Delete;
  списка у вида нет — общий `discount.List`. name, active, allProducts, allAgents (и usePriceType у
  спеццены) обязательны и передаются всегда, false — тоже. В Update пустой массив или null удаляет
  значение.
- Таблица называет коллекцию групп `productfolders`, примеры и фикстуры — `productFolders`: SDK пишет
  productFolders, а читает оба ключа во всех типах скидок. specialPrice в запросе — по таблице
  (priceType, value); пример создания передаёт meta типа цены, сервером форма записи не проверена.

## Коды маркировки — пакет `trackingcode`, `c.Entity().TrackingCode()`

`List` (limit до 100, codetype gs1/tag_1162/all), `Save` — один POST массивом (с id — изменение,
без — создание; cis и type — Field, при изменении по id их можно не передавать; вложенные коды
заменяются), `DeleteMany` (вложенные удаляются вместе) —
`entity/{type}/{id}/positions/{positionID}/trackingCodes`.

## Не входит

Expand. Заказ кодов маркировки и вывод из оборота — [emissionorder](emissionorder.md), [retireorder](retireorder.md).
