Status: resolved
Stage: APPROVED

# 3. Товары, услуги, комплекты, модификации и цены

## Problem Statement

Закрыть товарную модель и единое чтение ассортимента; Currency/PriceType/TaxRate имеют
собственные методы и response shapes.

## Solution

Завершённые вертикальные срезы в составе v1; каждый ticket включает models, methods,
тесты и собственное описание полей. Матрица ниже определяет scope этой захватки.
[Общая карта](../sdk-v1/map.md), [нормативные ограничения](../sdk-v1/spec.md),
[каталог JSON-контрактов](../sdk-v1/contracts.md). Текущая работа — планирование.

## User Stories

1. Как разработчик решения, я могу использовать «Currency Get/List/Create/Update» через
   публичный SDK
   и получать фактический ответ API с документированными ограничениями.
2. Как разработчик решения, я могу использовать «PriceType List/Get/GetDefault/ReplaceAll» через
   публичный SDK
   и получать фактический ответ API с документированными ограничениями.
3. Как разработчик решения, я могу использовать «TaxRate Get/List/Create/Update» через публичный
   SDK
   и получать фактический ответ API с документированными ограничениями.
4. Как разработчик решения, я могу использовать «Service Get/List/Create/Update» через публичный
   SDK
   и получать фактический ответ API с документированными ограничениями.
5. Как разработчик решения, я могу использовать «Характеристики Variant: definitions и значения»
   через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
6. Как разработчик решения, я могу использовать «Bundle и components» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
7. Как разработчик решения, я могу использовать «Variant Get/List/Create/Update» через публичный
   SDK
   и получать фактический ответ API с документированными ограничениями.
8. Как разработчик решения, я могу использовать «Assortment List/search и mixed rows» через
   публичный SDK
   и получать фактический ответ API с документированными ограничениями.

## Implementation Decisions и scope tickets

| ID | Ticket | Blocked by |
| --- | --- | --- |
| C01 | [Currency Get/List/Create/Update](issues/01-currency.md) | T02, T06 |
| C02 | [PriceType List/Get/GetDefault/ReplaceAll](issues/02-pricetype.md) | T02 |
| C03 | [TaxRate Get/List/Create/Update](issues/03-taxrate.md) | T02, T06 |
| C04 | [Service Get/List/Create/Update](issues/04-service.md) | T02, T06 |
| C05 | [Характеристики Variant: definitions и значения](issues/05-characteristics.md) | A05 |
| C06 | [Bundle и components](issues/06-bundle.md) | T02, T06 |
| C07 | [Variant Get/List/Create/Update](issues/07-variant.md) | T02, T06, C05 |
| C08 | [Assortment List/search и mixed rows](issues/08-assortment.md) | T02, T06, C04, C06, C07 |

### C01: Currency Get/List/Create/Update

Полный Currency contract, name/code/isoCode Create, rate и system/custom currency различия. Курс
не денежная сумма.

Проверки и особенности: Required/read-only и ограничения rate, нет автоподстановки курса;
user/system fixtures и APIError.

Источники: [currency](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_currency.md).

### C02: PriceType List/Get/GetDefault/ReplaceAll

Существующий model.PriceType переиспользуется без создания несовместимого второго типа. GET
context/companysettings/pricetype возвращает plain array; GET /{id}, GET /default; POST полного
упорядоченного списка ReplaceAll. Никакого вымышленного одиночного PUT.

Проверки и особенности: 1…100, default первый, порядок сохраняется, excluded item удаляется
сервером. Exact full-list JSON, no hidden GET/read-modify-write; externalCode/name/readonly
shape по методу; возможная потеря конкурентных изменений явно документируется.

Источники: [pricetype](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_pricetypes.md).

### C03: TaxRate Get/List/Create/Update

Полный TaxRate contract и readonly/required полей из source. Не смешивать ставку TaxRate с
денежными суммами или Decimal; server business validation.

Проверки и особенности: Нулевые ставки, optional dates/refs если есть, sparse payloads,
HTTP/query и error paths.

Источники: [taxrate](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_taxrate.md).

### C04: Service Get/List/Create/Update

Полный Service core/read-write contract: name, prices/taxes/attributes/refs и документированные
fields; common Price и SalePrice. Meta files/image links читаются только как объявленные
metadata, без догрузки.

Проверки и особенности: Оба codecs, minimal/full Create/Update, nil/false/0/null,
list/search/barcodes по источнику, read-only exclusions.

Источники: [service](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_service.md).

### C05: Характеристики Variant: definitions и значения

Разделить definition в variant metadata и value в Variant.characteristics. Чтение definitions;
Create по POST metadata/characteristics и Get definition по ID; отдельного Update в источнике
нет, не добавлять его. Это необходимое дополнение к Variant, не новый общий metadata CRUD
framework.

Проверки и особенности: Точные endpoints и list/envelope shapes; name/id/value неоднозначности
из примеров зафиксировать, unknown values, существующие definitions не создаются скрытым
запросом.

Источники: [characteristic](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_characteristic.md), [variant](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md).

### C06: Bundle и components

Bundle Get/List/Create/Update, components read collection и explicit GetComponents/GetComponent
при подтверждённых endpoints; embedded components write с ID/quantity/assortment по своему
contract. Components не притворяются документными Position, если отличаются смыслом.

Проверки и особенности: Полная таблица bundle/components, required/minimal,
replacement/limits/null из source, product/service/variant refs без expand;
salePrices/attributes сохраняют свои merge rules. Отдельные mutations компонентов не добавляются
без фиксации scope.

Источники: [bundle](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bundle.md).

### C07: Variant Get/List/Create/Update

Полный Variant contract, product link и characteristics, prices/barcodes/taxes и ограничения
immutable fields. Metadata-related definitions отдельны; Product не загружается автоматически.

Проверки и особенности: Create required и update restrictions по источнику; unknown
characteristic refs, false/0/[], missing/null, upstream+custom, list/filter/search. Full entity
Product service не технический blocker references.

Источники: [variant](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md).

### C08: Assortment List/search и mixed rows

Один смешанный список product/service/bundle/variant/consignment, единые search/filter/order и
pagination; read union с known typed four и Consignment reference/common-read, unknown meta.type
+ raw JSON сохраняется. Без generic write resource.

Проверки и особенности: Явный выбор режима: legacy без beta header либо облегчённый
assortmentWithoutStock с own capability/default groupBy. Не включать beta молча; default и
supported filters закрепить до code. Нет скрытых stock report запросов, mixed rows включая
consignment и неизвестный тип, один List вместо склейки четырёх списков.

Источники: [assortment](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_assortment.md).

## Общий контракт реализации

- Перед кодом — field/method freeze: поимённая Write matrix, readonly/TODO,
  null/required/Create-only/Update-only; status needs-triage сам не разрешает реализацию.
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
- Проверки: go test -race -count=1 ./..., go vet ./... ./testdata, gofmt,
  git diff --check; независимые Standards/Spec review на фиксированном BASE.
- Коммит по ticket; shared dependency раньше зависимого, публичные docs отдельным
  завершающим commit волны. Planning не коммитить без отдельной команды.
- Live API, credentials, push/tag/release, worktree cleanup и production-код
  не входят в текущую задачу подготовки spec/tickets.

## Open Decisions

- Assortment нужен для общего поиска/смешанной пагинации; отдельные packages обслуживают запись
  и специализированное чтение.
- Legacy и beta assortment имеют разные query/stock semantics; это явный mode, не
  взаимозаменяемые ответы.
- Новые ресурсы с syncId включают поле только по собственной таблице, общие price types не
  дублируются.

## Out of Scope

Expand и implicit follow-up; все не перечисленные в tickets ресурсы/операции.
План не объявляет общую CRUD поддержку у ресурсов с другой формой interface.
