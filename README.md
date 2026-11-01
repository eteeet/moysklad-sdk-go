# moysklad-sdk-go

SDK для МойСклад JSON API 1.2. Требуется Go 1.27.
Поддерживаемые возможности и план работ: [docs/coverage.md](docs/coverage.md).

Опции вызова и метаданные ответа, коды ошибок, выпуск токена, List с filter/order и обход
страниц, rate limiter, повторы и ограничения размера запроса: [docs/transport.md](docs/transport.md).

Сотрудники (`entity/employee`) и их права, роли (`entity/role`), контекст запроса (`contexts`),
метаданные и статусы (`metadata`), аудит (`audit`): [docs/access.md](docs/access.md).

Валюты, типы цен, ставки НДС, услуги, комплекты, модификации и их характеристики, ассортимент
(`entity/currency`, `pricetype`, `taxrate`, `service`, `bundle`, `variant`, `assortment`):
[docs/catalog.md](docs/catalog.md).

Точки продаж, розничные смены, продажи и возвраты, бонусные программы и операции
(`entity/retailstore`, `retailshift`, `retaildemand`, `retailsalesreturn`, `bonusprogram`,
`bonustransaction`): [docs/retail.md](docs/retail.md).
У всех справочников и документов ниже есть `List(ctx, query.ListOptions, opts...)`.

## Справочники

Помимо Product и UOM, поддержаны следующие согласованные срезы:

| Пакет | Методы | Документация и границы |
| --- | --- | --- |
| `entity/country` | `Get`, `Create`, `Update` | [Country](docs/country.md) |
| `entity/expenseitem` | `Get`, `Create`, `Update` | [ExpenseItem](docs/expenseitem.md) |
| `entity/counterparty` | `Get`, `Create`, `Update` | [Counterparty](docs/counterparty.md) |
| `entity/organization` | `Get`, `Create`, `Update` | [Organization](docs/organization.md) |

Общие `model.Address`, `AddressRequest` и открытый `CompanyType` описаны в
[документации реквизитов](docs/requisites.md). Счета и другие отложенные подресурсы,
Delete, expand и массовые операции в эти срезы не входят.

## Общие контракты документов и positions

Пакет `model` предоставляет `DocumentRate`/`DocumentRateRequest`,
`DocumentOverhead`/`DocumentOverheadRequest`, `FinanceOperation`/`FinanceOperationRequest`.
Это самостоятельные JSON-типы без HTTP и финансовых вычислений.
В присутствующем rate отсутствие Value означает курс 1; SDK не подставляет 1 в JSON.
[Документация общих типов](docs/document-common.md).

Общий пакет `positions` содержит универсальные Position/PositionRequest для строк
товарных документов. Core/Item и их request-формы — переиспользуемые части этих типов.
`positions.NewService[T](client, documentType)` предоставляет `GetPage(ctx, id, PageOptions)`:
одна страница `/positions` через существующий transport, с явными limit/offset и
`fields=stock`/`declaration` там, где источник их документирует, без expand, автообхода и скрытого Get. Используется существующий `model.Collection[T]`.
[Контракт позиций и пагинации](docs/positions.md).

Get/Create/Update документов и типизированные entity.GetPositions реализованы
в перечисленных ниже пакетах. Все деньги — float64 в копейках с ограничениями точности.

## Документы

| Пакет | Методы | Контракт полей и ограничения |
| --- | --- | --- |
| `entity/paymentin` | `Get/Create/Update` | [PaymentIn](docs/paymentin.md) |
| `entity/paymentout` | `Get/Create/Update` | [PaymentOut](docs/paymentout.md) |
| `entity/cashin` | `Get/Create/Update` | [CashIn](docs/cashin.md) |
| `entity/cashout` | `Get/Create/Update` | [CashOut](docs/cashout.md) |
| `entity/invoicein` | `Get/Create/Update` и `GetPositions` | [InvoiceIn](docs/invoicein.md) |
| `entity/invoiceout` | `Get/Create/Update` и `GetPositions` | [InvoiceOut](docs/invoiceout.md) |
| `entity/supply` | `Get/Create/Update` и `GetPositions` | [Supply](docs/supply.md) |
| `entity/demand` | `Get/Create/Update` и `GetPositions` | [Demand](docs/demand.md) |
| `entity/purchaseorder` | `Get/Create/Update` и `GetPositions` | [PurchaseOrder](docs/purchaseorder.md) |
| `entity/move` | `Get/Create/Update` и `GetPositions` | [Move](docs/move.md) |
| `entity/purchasereturn` | `Get/Create/Update` и `GetPositions` | [PurchaseReturn](docs/purchasereturn.md) |
| `entity/customerorder` | `Get/List/Create/Update` и `GetPositions` | [CustomerOrder](docs/customerorder.md) |
| `entity/enter` | `Get/List/Create/Update` и `GetPositions` | [Enter](docs/enter.md) |
| `entity/loss` | `Get/List/Create/Update` и `GetPositions` | [Loss](docs/loss.md) |
| `entity/salesreturn` | `Get/List/Create/Update` и `GetPositions` | [SalesReturn](docs/salesreturn.md) |

[Общий контракт](docs/documents.md): обязательные поля Create различаются по
документам; optional поля и Update используют Field. Get читает metadata positions,
строки — отдельным GetPositions. Embedded positions ограничены 1000 и полностью
заменяют набор при Update. Expand, Delete/batch и отдельные mutations строк отложены.

## Product

`product.NewService(client)` предоставляет `Get(ctx, id)`, `Create(ctx, CreateRequest)` и `Update(ctx, id, UpdateRequest)`.
Все три метода возвращают `(*Product, error)`; ошибки API доступны через `errors.As` как `*errors.APIError`
из пакета `github.com/eteeet/moysklad-sdk-go/errors` (импортируй с алиасом `sdkerrors`).

Create требует непустой Name. Update требует UUID и хотя бы одно заданное поле.
SDK не повторяет запись автоматически. После ошибки создания с неизвестным результатом проверь аккаунт
перед повторной отправкой; `CreateRequest.SyncID` (UUID клиента) защищает повторный Create от дубля.
`syncId` в API у товара только для чтения и задаётся при создании, поэтому в UpdateRequest его нет.

Поддерживаемые поля записи:

- `name`, `code`, `externalCode`, `description`, `article`, `archived`, `shared`.
- `weight`, `volume`, `vat`, `vatEnabled`, `useParentVat`, `discountProhibited`.
- `owner`, `group`, `productFolder`, `uom`, `supplier`, `country` — ссылки без expand.
- `buyPrice`, `minPrice`, `salePrices`, `attributes`.

В модели чтения дополнительно есть `meta`, `id`, `accountId`, `updated`, `syncId`, `pathName`,
`effectiveVat`, `effectiveVatEnabled`, `isSerialTrackable`.
Запросы используют отдельные модели: эти дополнительные поля в них не включены.

## Примеры

`.env` не загружается библиотекой. Makefile загружает его только для запуска примеров.
Храни в нём `MOYSKLAD_TOKEN`, а для GET — `MOYSKLAD_PRODUCT_ID`. Не коммить `.env`.

```sh
make example
```

Следующая команда действительно создаёт товар. Используй тестовый аккаунт; не сохраняй разрешение записи в `.env`:

```sh
MOYSKLAD_ALLOW_WRITE=1 make example-create
```

Пример печатает `MOYSKLAD_EXAMPLE_PRODUCT_ID` и `MOYSKLAD_EXAMPLE_EXTERNAL_CODE`.
Сохрани обе строки в `.env`, затем запусти:

```sh
MOYSKLAD_ALLOW_WRITE=1 make example-update
```

Перед PUT пример делает GET и проверяет, что ID, имя и externalCode соответствуют переданному маркеру `sdk-example-…`.
Это проверка соглашения об именовании, не доказательство происхождения товара.
Затем изменяет описание и устанавливает `archived: true`. Рабочий `MOYSKLAD_PRODUCT_ID` не используется.
Архивация не удаляет товар. После неоднозначного сбоя создания маркер из stderr поможет найти товар вручную.

Код примеров: [Get](examples/product_get/main.go), [Create](examples/product_create/main.go), [Update](examples/product_update/main.go).

## UOM

`uom.NewService(client)` создаёт обёртку с методами `Get(ctx, id)`,
`Create(ctx, CreateRequest)` и `Update(ctx, id, UpdateRequest)`.
Все три метода возвращают `(*UOM, error)`.
Get/Update проверяют переданный ID по форме UUID, без ограничения версии. Ошибки доступны через
`errors.As`/`errors.Is` по тому же контракту, что и у Product; автоматических повторов нет.

Модель чтения содержит `meta`, `id`, `accountId`, `updated`, `name`, `description`,
`code`, `externalCode`, `owner`, `group`, `shared`. Ссылки представлены без expand,
`code` остаётся строкой, включая ведущие нули и буквенные значения.

Get читает предустановленные и пользовательские UOM. У предустановленных UOM
могут отсутствовать `accountId`, `owner`, `group`, `shared`: тогда поля имеют
нулевые значения Go. `Shared == false` не определяет происхождение UOM.

Create создаёт пользовательскую UOM. Поля записи: `name`, `code`, `externalCode`,
`description`, `shared`, `owner`, `group`. Обязательный Name — string; остальные поля
используют `model.Field` с тремя состояниями: не задано, значение, явный null.
Create отклоняет пустой Name; дополнительные ограничения значений проверяет API.
После ошибки создания с неизвестным результатом проверь аккаунт перед повторной отправкой.

UpdateRequest использует `model.Field` для всех семи полей записи, включая Name.
Update требует хотя бы одно заданное поле и отправляет только заданные поля,
включая явные false, пустую строку и null; допустимость значения проверяет API.
Предварительного Get и локального слияния нет. Изменять разрешено только пользовательские UOM;
при отказе сервера, в том числе для предустановленной UOM, доступен `*sdkerrors.APIError`.

Пока нет Delete, expand и массовых операций UOM.
Проверки Get/Create/Update используют локальный HTTP-сервер, без обращения к живому аккаунту.

[Документация UOM](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/uom).

## DateTime

`model.DateTime` — именованный тип поверх `time.Time`; преобразование в обе стороны выполняется явно.

- Чтение: формат API с секундами и необязательной дробной частью; дополнительно RFC 3339 из официальных фикстур.
- Запись: `ГГГГ-ММ-ДД ЧЧ:мм:сс.ммм` в часовом поясе `Europe/Moscow`.
- При записи разряды меньше миллисекунды отбрасываются, а не округляются.
- `null` и пустая строка читаются как нулевое значение; нулевое значение записывается как `null`.
- Тег `omitzero` пропускает нулевое поле, в том числе нулевой момент времени в зоне Москвы.
- Поддерживаются `encoding/json` и `encoding/json/v2` через методы `MarshalJSON`/`UnmarshalJSON`.

[Формат даты и времени МойСклад](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-format-daty-i-vremeni).

## Часовой пояс в контейнерах

SDK загружает `Europe/Moscow` один раз. Эта зона учитывает исторические изменения московского времени.
Если загрузить базу часовых поясов невозможно, используется фиксированное смещение `MSK +03:00`.
В этом режиме исторические даты могут соответствовать неверному моменту времени.

Для минимальных контейнеров без базы часовых поясов добавь в главную программу:

```go
import _ "time/tzdata"
```

Сам SDK не включает базу принудительно. Подробнее: [time/tzdata](https://pkg.go.dev/time/tzdata).

## Необязательные поля запросов

`model.Field[T]` используется в моделях запросов:

- Нулевое значение — поле не отправляется благодаря тегу `json:"name,omitzero"`.
- `model.Null[T]()` — отправляется явный `null`.
- `model.Value(v)` — отправляется значение, включая `false`, `0`, `""` и пустой массив.

Без `omitzero` незаданное поле возвращает ошибку кодирования. `omitempty` не заменяет этот тег:
в json v2 он может скрыть явный `null` или пустой массив.
Обязательные поля CreateRequest остаются обычными типами; в UpdateRequest изменяемые поля используют `Field`.
Допустимость `null` определяется конкретным полем API: [поддержка null](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-podderzhka-null).

`IsZero()` проверяет состояние «не передано», `IsNull()` — явный `null`.
`Get()` возвращает `(значение, true)` только для состояния Value.
Внутреннее кодирование использует json v2 даже при внешнем json v1:
nil-слайс становится `[]`, nil-map — `{}`; `[]byte` имеет специальное представление строкой Base64.
`Value(DateTime{})`, nil-указатель и nil-интерфейс записываются как `null`; после JSON roundtrip это состояние Null.

При чтении `null` становится Null, остальные значения — Value. Отсутствующее поле не меняется:
для независимых запросов декодируй в новую структуру. Ошибка разбора значения внутри Field сохраняет его прежнее состояние,
но атомарность всего внешнего Unmarshal не гарантируется. При любой ошибке не используй частично декодированный результат.

[Правила кодирования JSON v2](https://pkg.go.dev/encoding/json/v2).

## Цены

`model.Price` используется для закупочной и минимальной цены, `model.SalePrice` — для цены продажи.
`model.PriceType` содержит тип цены; для запроса можно заполнить только его `Meta`.
Эти типы используются и при чтении, и внутри моделей запросов.

- `Value` — `float64` в копейках (сотых долях валюты), без пересчёта в рубли или другую валюту.
- Нулевая цена записывается как `"value":0`, а не пропускается.
- Незаполненные `Currency` и `PriceType` пропускаются через `omitzero`.
- Точная десятичная арифметика не предоставляется; `float64` имеет ограничения двоичного представления дробей.

[Документация цен](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#5-ceny-prodazhi).

## Дополнительные поля

`model.Attribute` — общая модель чтения и элемента запроса. Здесь Value и File также используют Field,
чтобы различать отсутствие, null и значение. Value хранит `jsontext.Value`, не промежуточный float64.
Неизвестные типы и исходное значение сохраняются; тип проверяется только при вызове геттера.

- `AsString()`, `AsBool()`, `AsLong()`, `AsDouble()`, `AsTime()`, `AsReference()` возвращают значение и ошибку.
- AsLong читает int64 напрямую: значения за пределами ±2⁵³ не теряют точность.
- AsReference возвращает `AttributeReference{Meta, Name}`; справочники распознаются по объекту со ссылкой, без закрытого списка типов.
- `AttrValue(meta, value)` возвращает `(Attribute, error)`; Type определяет сервер по метаданным.
- `AttrNull(meta)` сбрасывает value, `AttrFileNull(meta)` — file.
- `AttrFile(meta, filename, content)` принимает байты; content кодируется в Base64.
- Для даты передавай `model.DateTime`. SDK не меняет точность исходного значения; по контракту API сервер сохраняет дату доп. поля с точностью до минуты.

Для записи используй конструкторы, а не повторную отправку прочитанного Attribute вместе с read-only полями.
[Работа с дополнительными полями](https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-rabota-s-dopolnitelnymi-polyami).

## Семантика Update

Неуказанные поля не отправляются. `Value(false)`, `Value(0)` и `Value("")` — явные значения.
`Null[Reference]()` очищает ссылку только там, где сервер разрешает null, например productFolder.

Сервер сливает `salePrices` по типу цены (`priceType`), `attributes` — по метаданным поля (`meta`).
Переданные элементы создаются или перезаписываются, непереданные не меняются.
Пустые массивы, включая `Value([]SalePrice{})`, не очищают ранее заданные значения.
Для очистки конкретного доп. поля передавай элемент с value/file: null.
SDK не выполняет локальное слияние: возвращает ответ сервера.
[Изменить Товар](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/product#3-izmenit-tovar).

## Тесты и ограничения

```sh
go test -race -count=1 ./...
go vet ./... ./testdata
```

Официальные ответы лежат в `testdata/upstream`; версия и хеши — в [docs/upstream.md](docs/upstream.md).
`make fixtures` обновляет снимок, `make fixtures-commit` обновляет и коммитит только снимок и его описание.
Golden-тела запросов лежат в `testdata/custom/product`; они заданы независимо от сериализации SDK.
Сравнение канонизирует JSON без округления целых чисел за пределами ±2⁵³.

Срез проверен локальными HTTP-тестами, не реальной записью в аккаунт.
List, обход страниц, повторы и лимиты — общие механизмы ([docs/transport.md](docs/transport.md)).
Для Product пока нет expand, Delete, полного фасада и специальных операций товара; fields=minimumStock
отложен. Штрихкоды, алкоголь, упаковки, отслеживание и остальные неподдерживаемые поля остаются в TODO.
