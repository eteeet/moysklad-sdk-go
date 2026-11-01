Status: needs-triage
Stage: DRAFT

# Каталог контрактов v1

Источник всех таблиц — закреплённый official revision
`9a811cef3dc4a3f3a46f070e60452361a06bf7cd`.
Ниже перечислены JSON-поля и явные markers источника, без копирования описаний.
Это обязательный вход для spec/tickets, а не механически сгенерированная write-модель.
Отсутствие readonly marker не доказывает возможность записи. Create/Update имеют
собственные ограничения; required в ответе не означает required при Create.
До реализации каждого entity ticket нормативная Write matrix должна перечислить
для каждого поля C/U, Go type, required/null, create-only/write-once и explicit TODO.
Противоречия фиксируются в ticket, готовность не выводится из таблицы автоматически.

## employee

Ticket: A01. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md).
Минимальный Create / особенность: **lastName**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md#L14).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `archived` | Boolean | response required |
| `attributes` | Array(Object) | — |
| `cashiers` | MetaArray | readonly, Reference, без expand |
| `created` | DateTime | readonly, response required |
| `description` | String(4096) | — |
| `email` | String(255) | — |
| `externalCode` | String(255) | response required |
| `firstName` | String(255) | — |
| `fullName` | String(255) | readonly |
| `group` | Meta | response required, Reference, без expand |
| `id` | UUID | readonly, response required |
| `image` | Object | — |
| `inn` | String(255) | — |
| `lastName` | String(255) | Create required, response required |
| `meta` | Meta | response required |
| `middleName` | String(255) | — |
| `name` | String(255) | readonly, response required |
| `owner` | Meta | response required, Reference, без expand |
| `phone` | String(255) | — |
| `position` | String(255) | — |
| `salary` | Object | — |
| `shared` | Boolean | response required |
| `shortFio` | String(255) | readonly |
| `uid` | String(255) | readonly |
| `updated` | DateTime | readonly, response required |

### Оклад

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md#L49).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `value` | Float | — |

### Кассир

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md#L55).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `employee` | Meta | readonly, response required, Reference, без expand |
| `id` | UUID | readonly, response required |
| `retailStore` | Meta | readonly, response required, Reference, без expand |

### Фотография сотрудника: структура и загрузка.

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md#L65).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `filename` | String(255) | response required |
| `meta` | Meta | response required |
| `miniature` | Meta | response required |
| `size` | Int | response required |
| `tiny` | Meta | response required |
| `title` | String(255) | response required |
| `updated` | DateTime | response required |

### Получить Сотрудников

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md#L94).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

### Метаданные Сотрудников

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md#L423).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `attributes` | Array(Object) | — |
| `createShared` | Boolean | — |

### Атрибуты сущности Сотрудник

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md#L795).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `authorizedHosts` | Array(String) | — |
| `authorizedIpNetmask` | String(255) | — |
| `authorizedIpNetwork` | String(255) | — |
| `email` | String(255) | — |
| `group` | Object | response required |
| `isActive` | Boolean | response required |
| `login` | String(255) | — |
| `role` | Object | — |

### Роль

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md#L817).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | response required |
| `permissions` | Array(Object) | — |

## role

Ticket: A02. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_custom_role.md).
Минимальный Create / особенность: **name; permissions — по методу**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_custom_role.md#L11).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `id` | UUID | readonly, response required |
| `meta` | Meta | response required |
| `name` | String(255) | Create required, response required |
| `permissions` | Array(Object) | response required |

## currency

Ticket: C01. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_currency.md).
Минимальный Create / особенность: **name, code, isoCode; rate!=0 при задании**.

### Атрибуты Сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_currency.md#L13).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `archived` | Boolean | response required |
| `code` | String(255) | Create required, response required |
| `default` | Boolean | readonly, response required |
| `fullName` | String(255) | — |
| `id` | UUID | readonly, response required |
| `indirect` | Boolean | response required |
| `isoCode` | String(255) | Create required, response required |
| `majorUnit` | Object | response required |
| `margin` | Double | response required |
| `meta` | Meta | response required |
| `minorUnit` | Object | response required |
| `multiplicity` | Int | response required |
| `name` | String(255) | Create required, response required |
| `rate` | Double | response required |
| `rateUpdateType` | String(255) | readonly, response required |
| `system` | Boolean | readonly |

## pricetype

Ticket: C02. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_pricetypes.md).
Минимальный Create / особенность: **полный упорядоченный список, 1…100 элементов**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_pricetypes.md#L7).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `externalCode` | String(255) | response required |
| `id` | UUID | readonly, response required |
| `meta` | Meta | readonly, response required |
| `name` | String(255) | Create required, response required |

## taxrate

Ticket: C03. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_taxrate.md).
Минимальный Create / особенность: **по таблице и методу Create**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_taxrate.md#L13).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | response required |
| `comment` | String | — |
| `group` | Meta | Reference, без expand |
| `id` | UUID | readonly, response required |
| `meta` | Meta | response required |
| `rate` | Number | Create required, response required |
| `owner` | Meta | Reference, без expand |
| `shared` | Boolean | — |
| `updated` | DateTime | Create required, response required |
| `archived` | Boolean | — |

## service

Ticket: C04. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_service.md).
Минимальный Create / особенность: **name**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_service.md#L15).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `archived` | Boolean | response required |
| `attributes` | Array(Object) | — |
| `barcodes` | Array(Object) | — |
| `buyPrice` | Object | — |
| `code` | String(255) | — |
| `description` | String(4096) | — |
| `discountProhibited` | Boolean | response required |
| `effectiveVat` | Int | readonly |
| `effectiveVatEnabled` | Boolean | readonly |
| `externalCode` | String(255) | response required |
| `files` | MetaArray | Reference, без expand |
| `group` | Meta | response required, Reference, без expand |
| `id` | UUID | readonly, response required |
| `meta` | Meta | response required |
| `minPrice` | Object | — |
| `name` | String(255) | Create required, response required |
| `owner` | Meta | Reference, без expand |
| `pathName` | String | readonly, response required |
| `paymentItemType` | Enum | — |
| `productFolder` | Meta | Reference, без expand |
| `salePrices` | Array(Object) | — |
| `shared` | Boolean | response required |
| `syncId` | UUID | readonly |
| `taxSystem` | Enum | — |
| `uom` | Meta | Reference, без expand |
| `updated` | DateTime | readonly, response required |
| `useParentVat` | Boolean | response required |
| `vat` | Int | — |
| `vatEnabled` | Boolean | — |

### Регионоспецифичные Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_service.md#L53).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `mod\_\_tasnif\_\_uz` | Object | — |

### Цены продажи

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_service.md#L106).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `value` | Float | response required |
| `currency` | Meta | response required, Reference, без expand |
| `priceType` | Object | response required |

### Закупочная цена

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_service.md#L115).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `value` | Float | response required |
| `currency` | Meta | response required, Reference, без expand |

### Минимальная цена

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_service.md#L122).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `value` | Float | response required |
| `currency` | Meta | response required, Reference, без expand |

### Группа Услуги

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_service.md#L129).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | response required |

### Модуль для Узбекистана с информацией из справочника ТАСНИФ

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_service.md#L137).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `ikpu` | String | — |
| `packCode` | String | — |
| `barcodeTasnif` | String | — |

### Получить список Услуг

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_service.md#L151).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

## bundle

Ticket: C06. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bundle.md).
Минимальный Create / особенность: **name; components — по методу**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bundle.md#L8).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `archived` | Boolean | response required |
| `article` | String(255) | — |
| `attributes` | Array(Object) | — |
| `barcodes` | Array(Object) | — |
| `code` | String(255) | — |
| `components` | MetaArray | Reference, без expand |
| `country` | Meta | Reference, без expand |
| `description` | String(4096) | — |
| `discountProhibited` | Boolean | response required |
| `effectiveVat` | Int | readonly |
| `effectiveVatEnabled` | Boolean | readonly |
| `externalCode` | String(255) | response required |
| `files` | MetaArray | Reference, без expand |
| `group` | Meta | response required, Reference, без expand |
| `id` | UUID | readonly, response required |
| `images` | MetaArray | Reference, без expand |
| `meta` | Meta | response required |
| `minPrice` | Object | — |
| `name` | String(255) | Create required, response required |
| `overhead` | Object | — |
| `owner` | Meta | Reference, без expand |
| `partialDisposal` | Boolean | — |
| `pathName` | String | readonly, response required |
| `paymentItemType` | Enum | — |
| `productFolder` | Meta | Reference, без expand |
| `salePrices` | Array(Object) | — |
| `shared` | Boolean | response required |
| `syncId` | UUID | readonly |
| `taxSystem` | Enum | — |
| `tnved` | String(255) | — |
| `trackingType` | Enum | — |
| `uom` | Meta | Reference, без expand |
| `updated` | DateTime | readonly, response required |
| `useParentVat` | Boolean | response required |
| `vat` | Int | — |
| `vatEnabled` | Boolean | — |
| `volume` | Float | — |
| `weight` | Float | — |

### Регионоспецифичные Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bundle.md#L53).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `mod\_\_marking\_\_uz` | Object | — |
| `mod\_\_tasnif\_\_uz` | Object | — |

### Модуль маркирования продукции для Узбекистана

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bundle.md#L121).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `trackingType` | Enum | — |

### Модуль для Узбекистана с информацией из справочника ТАСНИФ

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bundle.md#L138).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `ikpu` | String | — |
| `packCode` | String | — |
| `barcodeTasnif` | String | — |

### Дополнительные расходы

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bundle.md#L168).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `value` | Float | response required |
| `currency` | Meta | response required, Reference, без expand |

### Минимальная цена

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bundle.md#L176).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `value` | Float | response required |
| `currency` | Meta | response required, Reference, без expand |

### Компоненты Комплекта

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bundle.md#L185).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `assortment` | Meta | response required, Reference, без expand |
| `id` | UUID | readonly, response required |
| `quantity` | Float | readonly, response required |

### Изображение: структура и загрузка.

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bundle.md#L228).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `filename` | String(255) | response required |
| `meta` | Meta | response required |
| `miniature` | Meta | response required |
| `size` | Int | response required |
| `tiny` | Meta | response required |
| `title` | String(255) | response required |
| `updated` | DateTime | response required |

### Получить список комплектов

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bundle.md#L258).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

### Получить компоненты Комплекта

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bundle.md#L1743).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

## variant

Ticket: C07. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md).
Минимальный Create / особенность: **product, characteristics — сверить требования метода**.

### Атрибуты Сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L23).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `archived` | Boolean | response required |
| `article` | String(255) | — |
| `barcodes` | Array(Object) | — |
| `buyPrice` | Object | — |
| `characteristics` | Array(Object) | Create required, response required |
| `code` | String(255) | — |
| `description` | String(4096) | — |
| `discountProhibited` | Boolean | response required |
| `externalCode` | String(255) | response required |
| `id` | UUID | readonly, response required |
| `images` | MetaArray | response required, Reference, без expand |
| `meta` | Meta | response required |
| `minPrice` | Object | — |
| `minimumStock` | Object | — |
| `name` | String(255) | response required |
| `packs` | Array(Object) | — |
| `product` | Meta | Create required, response required, Reference, без expand |
| `salePrices` | Array(Object) | — |
| `things` | Array(String) | readonly |
| `updated` | DateTime | readonly, response required |

### Метаданные Модификаций

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L86).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `id` | UUID | response required |
| `meta` | Meta | response required |
| `name` | String(255) | response required |
| `value` | String(255) | Create required, response required |

### Метаданные Модификаций

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L97).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | response required |
| `characteristics` | Array(Object) | response required |

### Характеристики модификации

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L105).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `id` | UUID | response required |
| `meta` | Meta | response required |
| `name` | String(255) | response required |
| `required` | Boolean | response required |
| `type` | String(255) | response required |

### Изображение: структура и загрузка.

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L117).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `filename` | String(255) | response required |
| `meta` | Meta | response required |
| `miniature` | Meta | response required |
| `size` | Int | response required |
| `tiny` | Meta | response required |
| `title` | String(255) | response required |
| `updated` | DateTime | response required |

### Цены продажи

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L144).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `value` | Float | response required |
| `currency` | Meta | response required, Reference, без expand |
| `priceType` | Object | response required |

### Минимальная цена

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L152).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `value` | Float | response required |
| `currency` | Meta | response required, Reference, без expand |

### Атрибуты Неснижаемого остатка с типом ALL_WAREHOUSE_SUM

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L179).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `type` | String | — |
| `inherited` | Boolean | — |
| `quantity` | Double | — |

### Атрибуты Неснижаемого остатка с типом ALL_WAREHOUSE_SAME

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L187).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `type` | String | readonly |
| `inherited` | Boolean | — |
| `quantity` | Double | readonly |

### Атрибуты Неснижаемого остатка с типом WAREHOUSE_VARIED

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L195).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `type` | String | — |
| `inherited` | Boolean | — |
| `storebalances` | MetaArray | — |

### Атрибуты Неснижаемого остатка с типом WAREHOUSE_VARIED

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L203).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly |
| `id` | UUID | readonly |
| `meta` | Meta | readonly |
| `store` | Meta | — |
| `quantity` | Double | — |

### Упаковки Модификации

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L224).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `barcodes` | Array(Object) | response required |
| `id` | UUID | readonly, response required |
| `parentpack` | Meta | response required, Reference, без expand |

### Получить список Модификаций

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L234).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

### Создать Модификацию

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L540).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `id` | UUID | — |
| `name` | String(255) | — |
| `value` | String(255) | Create required, response required |

### Запросы - Метаданные Модификаций

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L1409).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | response required |
| `characteristics` | Array(Object) | response required |

### Изменить Модификацию

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md#L1621).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `id` | UUID | — |
| `name` | String(255) | — |
| `value` | String(255) | Create required, response required |

## assortment

Ticket: C08. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_assortment.md).
Минимальный Create / особенность: **read-only List; Create/Update не добавлять**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_assortment.md#L212).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | response required |
| `uniqueCodeRules` | Object | response required |
| `barcodeRules` | Object | response required |
| `createdShared` | Boolean | response required |

### Настройки уникальности кода для сущностей справочника

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_assortment.md#L222).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `checkUniqueCode` | Boolean | response required |
| `fillUniqueCode` | Boolean | response required |

### Настройки правил штрихкодов для сущностей справочника

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_assortment.md#L229).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `fillEAN13Barcode` | Boolean | response required |
| `weightBarcode` | Boolean | response required |
| `weightBarcodePrefix` | Int | response required |

## characteristic

Ticket: C05. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_characteristic.md).
Минимальный Create / особенность: **metadata resource; по методу, не обычный entity CRUD**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_characteristic.md#L8).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `id` | UUID | readonly, response required |
| `meta` | Meta | readonly, response required |
| `name` | String(255) | Create required, response required |
| `required` | Boolean | readonly, response required |
| `type` | String(255) | readonly, response required |

## customerorder

Ticket: D02. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_customerOrder.md).
Минимальный Create / особенность: **organization, agent**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_customerOrder.md#L6).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `agent` | Meta | Create required, response required, Reference, без expand |
| `agentAccount` | Meta | Reference, без expand |
| `applicable` | Boolean | response required |
| `attributes` | Array(Object) | — |
| `code` | String(255) | — |
| `contract` | Meta | Reference, без expand |
| `created` | DateTime | readonly, response required |
| `deleted` | DateTime | readonly |
| `deliveryPlannedMoment` | DateTime | — |
| `description` | String(4096) | — |
| `externalCode` | String(255) | response required |
| `files` | MetaArray | response required, Reference, без expand |
| `group` | Meta | response required, Reference, без expand |
| `id` | UUID | readonly, response required |
| `invoicedSum` | Float | readonly, response required |
| `meta` | Meta | response required |
| `moment` | DateTime | response required |
| `name` | String(255) | response required |
| `organization` | Meta | Create required, response required, Reference, без expand |
| `organizationAccount` | Meta | Reference, без expand |
| `owner` | Meta | Reference, без expand |
| `payedSum` | Float | readonly, response required |
| `positions` | MetaArray | response required, Reference, без expand |
| `printed` | Boolean | readonly, response required |
| `project` | Meta | Reference, без expand |
| `published` | Boolean | readonly, response required |
| `rate` | Object | response required |
| `reservedSum` | Float | readonly, response required |
| `salesChannel` | Meta | Reference, без expand |
| `shared` | Boolean | response required |
| `shipmentAddress` | String(255) | — |
| `shipmentAddressFull` | Object | — |
| `shippedSum` | Float | readonly, response required |
| `state` | Meta | Reference, без expand |
| `store` | Meta | Reference, без expand |
| `sum` | Float | readonly, response required |
| `syncId` | UUID | — |
| `taxSystem` | Enum | — |
| `updated` | DateTime | readonly, response required |
| `vatEnabled` | Boolean | response required |
| `vatIncluded` | Boolean | — |
| `vatSum` | Float | readonly |

### Позиции Заказа покупателя

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_customerOrder.md#L80).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `assortment` | Meta | response required, Reference, без expand |
| `discount` | Float | response required |
| `id` | UUID | readonly, response required |
| `pack` | Object | — |
| `price` | Float | response required |
| `quantity` | Float | response required |
| `reserve` | Float | — |
| `shipped` | Float | readonly, response required |
| `taxSystem` | Enum | — |
| `vat` | Int | response required |
| `vatEnabled` | Boolean | response required |

### Атрибуты сущности Адрес доставки

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_customerOrder.md#L109).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `addInfo` | String(255) | — |
| `apartment` | String(30) | — |
| `city` | String(255) | — |
| `comment` | String(255) | — |
| `country` | Meta | — |
| `house` | String(30) | — |
| `postalCode` | String(6) | — |
| `region` | Meta | — |
| `street` | String(255) | — |

### Получить список Заказов покупателей

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_customerOrder.md#L132).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

### Получить позиции Заказа покупателя

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_customerOrder.md#L2377).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

## enter

Ticket: D03. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_enter.md).
Минимальный Create / особенность: **organization, store**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_enter.md#L7).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `applicable` | Boolean | response required |
| `attributes` | Array(Object) | — |
| `code` | String(255) | — |
| `created` | DateTime | readonly, response required |
| `deleted` | DateTime | readonly |
| `description` | String(4096) | — |
| `externalCode` | String(255) | response required |
| `files` | MetaArray | response required, Reference, без expand |
| `group` | Meta | response required, Reference, без expand |
| `id` | UUID | readonly, response required |
| `meta` | Meta | response required |
| `moment` | DateTime | response required |
| `name` | String(255) | response required |
| `organization` | Meta | Create required, response required, Reference, без expand |
| `overhead` | Object | — |
| `owner` | Meta | Reference, без expand |
| `positions` | MetaArray | response required, Reference, без expand |
| `printed` | Boolean | readonly, response required |
| `project` | Meta | Reference, без expand |
| `published` | Boolean | readonly, response required |
| `rate` | Object | response required |
| `shared` | Boolean | response required |
| `state` | Meta | Reference, без expand |
| `store` | Meta | Create required, response required, Reference, без expand |
| `sum` | Float | readonly, response required |
| `syncId` | UUID | — |
| `updated` | DateTime | readonly, response required |

### Накладные расходы

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_enter.md#L41).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `sum` | Float | response required |
| `distribution` | Enum | response required |

### Позиции Оприходования

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_enter.md#L56).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `assortment` | Meta | response required, Reference, без expand |
| `country` | Meta | Reference, без expand |
| `gtd` | Object | — |
| `id` | UUID | readonly, response required |
| `overhead` | Int | readonly, response required |
| `pack` | Object | — |
| `price` | Float | response required |
| `quantity` | Float | response required |
| `reason` | String(255) | — |
| `slot` | Meta | Reference, без expand |
| `things` | Object(String) | — |

### Получать Оприходования

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_enter.md#L88).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

### Создать Оприходования

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_enter.md#L278).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `organization` | Meta | — |
| `store` | Meta | — |

### Получить позиции Оприходования

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_enter.md#L1278).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

## loss

Ticket: D04. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_loss.md).
Минимальный Create / особенность: **organization, store**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_loss.md#L6).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `applicable` | Boolean | response required |
| `attributes` | Array(Object) | — |
| `code` | String(255) | — |
| `created` | DateTime | readonly, response required |
| `deleted` | DateTime | readonly |
| `description` | String(4096) | — |
| `expenseItem` | Meta | readonly, Reference, без expand |
| `externalCode` | String(255) | response required |
| `files` | MetaArray | response required, Reference, без expand |
| `group` | Meta | response required, Reference, без expand |
| `id` | UUID | readonly, response required |
| `meta` | Meta | response required |
| `moment` | DateTime | response required |
| `name` | String(255) | response required |
| `organization` | Meta | Create required, response required, Reference, без expand |
| `owner` | Meta | Reference, без expand |
| `positions` | MetaArray | response required, Reference, без expand |
| `printed` | Boolean | readonly, response required |
| `project` | Meta | Reference, без expand |
| `published` | Boolean | readonly, response required |
| `rate` | Object | response required |
| `shared` | Boolean | response required |
| `state` | Meta | Reference, без expand |
| `store` | Meta | Create required, response required, Reference, без expand |
| `sum` | Float | readonly, response required |
| `syncId` | UUID | — |
| `updated` | DateTime | readonly, response required |

### Позиции Списания

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_loss.md#L49).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `assortment` | Meta | response required, Reference, без expand |
| `declaration` | Array(Object) | readonly |
| `id` | UUID | readonly, response required |
| `pack` | Object | — |
| `price` | Float | response required |
| `quantity` | Float | response required |
| `reason` | String(255) | — |
| `slot` | Meta | Reference, без expand |
| `things` | Array(String) | — |

### Информация о прослеживаемости импортных товаров

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_loss.md#L80).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `gtd` | String | readonly |
| `rnpt` | String | readonly |
| `country` | Meta | readonly |
| `quantity` | Float | readonly |

### Получить Списания

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_loss.md#L91).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

### Получить позиции Списания

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_loss.md#L1506).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

## salesreturn

Ticket: D05. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_sales_return.md).
Минимальный Create / особенность: **organization, store, agent; demand только для возврата по
основанию**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_sales_return.md#L6).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `agent` | Meta | response required, Reference, без expand |
| `agentAccount` | Meta | Reference, без expand |
| `applicable` | Boolean | response required |
| `attributes` | Array(Object) | — |
| `code` | String(255) | — |
| `contract` | Meta | Reference, без expand |
| `created` | DateTime | readonly, response required |
| `deleted` | DateTime | readonly |
| `description` | String(4096) | — |
| `externalCode` | String(255) | response required |
| `files` | MetaArray | response required, Reference, без expand |
| `group` | Meta | response required, Reference, без expand |
| `id` | UUID | readonly, response required |
| `meta` | Meta | response required |
| `moment` | DateTime | response required |
| `name` | String(255) | response required |
| `organization` | Meta | Create required, response required, Reference, без expand |
| `organizationAccount` | Meta | Reference, без expand |
| `owner` | Meta | Reference, без expand |
| `positions` | MetaArray | response required, Reference, без expand |
| `printed` | Boolean | readonly, response required |
| `project` | Meta | Reference, без expand |
| `published` | Boolean | readonly, response required |
| `rate` | Object | response required |
| `salesChannel` | Meta | Reference, без expand |
| `shared` | Boolean | response required |
| `state` | Meta | Reference, без expand |
| `store` | Meta | Create required, response required, Reference, без expand |
| `sum` | Float | readonly, response required |
| `syncId` | UUID | — |
| `updated` | DateTime | readonly, response required |
| `vatEnabled` | Boolean | response required |
| `vatIncluded` | Boolean | — |
| `vatSum` | Float | readonly |

### Позиции Возврата покупателя

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_sales_return.md#L58).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `assortment` | Meta | response required, Reference, без expand |
| `cost` | Int | — |
| `country` | Meta | Reference, без expand |
| `discount` | Float | response required |
| `gtd` | Object | — |
| `id` | UUID | readonly, response required |
| `pack` | Object | — |
| `price` | Float | response required |
| `quantity` | Float | response required |
| `slot` | Meta | Reference, без expand |
| `things` | Array(String) | — |
| `vat` | Int | response required |
| `vatEnabled` | Boolean | response required |

### Получить Возвраты покупателей

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_sales_return.md#L92).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

### Получить Позиции

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_sales_return.md#L1732).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

## retailstore

Ticket: E01. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md).
Минимальный Create / особенность: **name, organization, store, priceType**.

### Точки продаж

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md#L12).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `acquire` | Meta | response required, Reference, без expand |
| `active` | Boolean | response required |
| `address` | String(255) | — |
| `addressFull` | Object | — |
| `allowCreateProducts` | Boolean | response required |
| `allowCustomPrice` | Boolean | readonly, response required |
| `allowDeleteReceiptPositions` | Boolean | response required |
| `allowSellTobaccoWithoutMRC` | Boolean | response required |
| `archived` | Boolean | response required |
| `authTokenAttached` | Boolean | readonly, response required |
| `bankPercent` | Double | — |
| `cashiers` | Meta | response required, Reference, без expand |
| `controlCashierChoice` | Boolean | response required |
| `controlShippingStock` | Boolean | response required |
| `createAgentsTags` | Array(Object) | — |
| `createCashInOnRetailShiftClosing` | Boolean | response required |
| `createOrderWithState` | Meta | Reference, без expand |
| `createPaymentInOnRetailShiftClosing` | Boolean | response required |
| `customerOrderStates` | Array(Object) | Reference, без expand |
| `defaultTaxSystem` | Enum | response required |
| `demandPrefix` | String(255) | — |
| `description` | String(4096) | — |
| `discountEnable` | Boolean | response required |
| `discountMaxPercent` | Int | — |
| `enableReturnsWithNoReason` | Boolean | response required |
| `environment` | Object | readonly, response required |
| `externalCode` | String(255) | readonly, response required |
| `filterAgentsTags` | Array(String) | — |
| `fiscalType` | Enum | response required |
| `group` | Meta | response required, Reference, без expand |
| `id` | UUID | readonly, response required |
| `idQR` | String(255) | — |
| `issueOrders` | Boolean | response required |
| `lastOperationNames` | Array(Object) | readonly, response required |
| `markingSellingMode` | Enum | response required |
| `marksCheckMode` | Enum | response required |
| `masterRetailStores` | MetaArray | Reference, без expand |
| `meta` | Meta | response required |
| `minionToMasterType` | Enum | response required |
| `name` | String(255) | Create required, response required |
| `ofdEnabled` | Boolean | readonly, response required |
| `onlyInStock` | Boolean | response required |
| `orderTaxSystem` | Enum | response required |
| `orderToState` | Meta | Reference, без expand |
| `organization` | Meta | Create required, response required, Reference, без expand |
| `organizationBranch` | Meta | Reference, без expand |
| `owner` | Meta | Reference, без expand |
| `priceType` | Object | Create required, response required |
| `printAlways` | Boolean | response required |
| `priorityOfdSend` | Enum | response required |
| `productFolders` | MetaArray | Reference, без expand |
| `qrAcquire` | Meta | Reference, без expand |
| `qrBankPercent` | Double | — |
| `qrPayEnabled` | Boolean | response required |
| `qrTerminalId` | String(255) | — |
| `receiptTemplate` | Meta | Reference, без expand |
| `requiredFio` | Boolean | response required |
| `requiredPhone` | Boolean | response required |
| `requiredEmail` | Boolean | response required |
| `requiredBirthdate` | Boolean | response required |
| `requiredSex` | Boolean | response required |
| `requiredDiscountCardNumber` | Boolean | response required |
| `reservePrepaidGoods` | Boolean | response required |
| `returnFromClosedShiftEnabled` | Boolean | response required |
| `sellReserves` | Boolean | response required |
| `sendMarksForCheck` | Boolean | response required |
| `sendMarksToChestnyZnakOnCloud` | Boolean | response required |
| `syncAgents` | Boolean | response required |
| `shared` | Boolean | response required |
| `showBeerOnTap` | Boolean | response required |
| `state` | Object | readonly |
| `store` | Meta | Create required, response required, Reference, без expand |
| `tobaccoMrcControlType` | Enum | response required |
| `updated` | DateTime | readonly, response required |

### Окружение

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md#L163).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `device` | String(255) | — |
| `os` | String(255) | — |
| `software` | Object | — |
| `chequePrinter` | Object | — |
| `paymentTerminal` | String(255) | — |

### Атрибуты сущности ПО

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md#L173).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `name` | String(255) | Create required, response required |
| `vendor` | String(255) | — |
| `version` | String(255) | — |

### Атрибуты сущности ККТ

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md#L181).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `driver` | Object | — |
| `firmwareVersion` | String(255) | — |
| `fiscalDataVersion` | String(255) | — |
| `fiscalMemory` | Object | — |
| `name` | String(255) | Create required, response required |
| `serial` | String(255) | — |
| `vendor` | String(255) | — |

### Атрибуты сущности Драйвер

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md#L193).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `name` | String(255) | — |
| `version` | String(255) | — |

### Атрибуты сущности Фискальный накопитель

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md#L200).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `fiscalDataVersion` | String(255) | — |
| `fiscalValidityDate` | DateTime | — |

### Атрибуты сущности Статус

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md#L207).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `sync` | Object | — |
| `lastCheckMoment` | DateTime | — |
| `fiscalMemory` | Object | — |
| `paymentTerminal` | Object | — |

### Атрибуты сущности Синхронизация

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md#L216).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `message` | String(255) | — |
| `lastAttempMoment` | DateTime | Create required |

### Атрибуты сущности Фискальная Память

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md#L223).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `error` | Object | — |
| `notSendDocCount` | Int | — |
| `notSendFirstDocMoment` | DateTime | — |

### Атрибуты сущности Ошибка

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md#L231).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `сode` | String(255) | response required |
| `message` | String(255) | response required |

### Атрибуты сущности Платежный Терминал

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md#L237).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `acquiringType` | String(255) | — |

### Атрибуты сущности Адрес

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md#L243).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `addInfo` | String(255) | — |
| `apartment` | String(30) | — |
| `city` | String(255) | — |
| `comment` | String(255) | — |
| `country` | Meta | — |
| `house` | String(30) | — |
| `postalCode` | String(6) | — |
| `region` | Meta | — |
| `street` | String(255) | — |

### Последние операции

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md#L269).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `entity` | String(255) | readonly |
| `name` | String(255) | readonly |

### Кассир

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md#L283).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `employee` | Meta | readonly, response required, Reference, без expand |
| `id` | UUID | readonly, response required |
| `meta` | Meta | response required |
| `retailStore` | Meta | readonly, response required, Reference, без expand |

### Получить точки продаж

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md#L296).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

## retailshift

Ticket: E02. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retailshift.md).
Минимальный Create / особенность: **organization, retailStore**.

### Атрибуты смены

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retailshift.md#L9).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `acquire` | Meta | response required, Reference, без expand |
| `attributes` | Array(Object) | — |
| `bankComission` | Double | response required |
| `bankPercent` | Double | response required |
| `cheque` | Object | — |
| `closeDate` | DateTime | — |
| `created` | DateTime | readonly, response required |
| `deleted` | DateTime | readonly |
| `description` | String(4096) | — |
| `externalCode` | String(255) | response required |
| `files` | MetaArray | response required, Reference, без expand |
| `group` | Meta | response required, Reference, без expand |
| `id` | UUID | readonly, response required |
| `meta` | Meta | readonly, response required |
| `moment` | DateTime | response required |
| `name` | String(255) | response required |
| `operations` | Array(Object) | readonly, response required, Reference, без expand |
| `organization` | Meta | Create required, response required, Reference, без expand |
| `organizationAccount` | Meta | Reference, без expand |
| `organizationBranch` | Meta | Reference, без expand |
| `owner` | Meta | Reference, без expand |
| `paymentOperations` | Array(Object) | readonly, Reference, без expand |
| `printed` | Boolean | readonly, response required |
| `proceedsCash` | Float | readonly, response required |
| `proceedsNoCash` | Float | readonly, response required |
| `published` | Boolean | readonly, response required |
| `qrAcquire` | Meta | response required, Reference, без expand |
| `qrBankComission` | Double | response required |
| `qrBankPercent` | Double | response required |
| `receivedCash` | Float | readonly, response required |
| `receivedNoCash` | Float | readonly, response required |
| `retailStore` | Meta | Create required, response required, Reference, без expand |
| `shared` | Boolean | response required |
| `store` | Meta | response required, Reference, без expand |
| `syncId` | UUID | — |
| `updated` | DateTime | response required |

### Информация о смене ККТ

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retailshift.md#L52).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `start` | Object | — |
| `end` | Object | — |

### Информация об открытии смены ККТ

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retailshift.md#L58).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `fnNumber` | String | — |
| `kktRegNumber` | String | — |
| `fiscalDocSign` | String | — |
| `shiftNumber` | String | — |
| `fiscalDocNumber` | String | — |
| `time` | DateTime | — |

### Информация о закрытии смены ККТ

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retailshift.md#L69).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `chequesTotal` | Int | — |
| `fiscalDocNumber` | String | — |
| `fiscalDocSign` | String | — |
| `fiscalDocsTotal` | Int | — |
| `fnNumber` | String | — |
| `kktRegNumber` | String | — |
| `shiftNumber` | String | — |
| `time` | DateTime | — |

### Получить Розничные смены

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retailshift.md#L85).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

### Создать розничную смену

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retailshift.md#L495).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `organization` | Meta | — |
| `retailStore` | Meta | — |

## retaildemand

Ticket: E03. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retaildemand.md).
Минимальный Create / особенность: **retailShift; активность смены проверяет сервер**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retaildemand.md#L6).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `advancePaymentSum` | Float | response required |
| `agent` | Meta | Create required, response required, Reference, без expand |
| `agentAccount` | Meta | Reference, без expand |
| `applicable` | Boolean | response required |
| `attributes` | Array(Object) | — |
| `cashSum` | Float | response required |
| `checkNumber` | Int | — |
| `checkSum` | Float | — |
| `cheque` | Object | readonly |
| `code` | String(255) | — |
| `created` | DateTime | readonly, response required |
| `customerOrder` | Meta | Reference, без expand |
| `deleted` | DateTime | readonly |
| `description` | String(4096) | — |
| `documentNumber` | Int | — |
| `externalCode` | String(255) | response required |
| `files` | MetaArray | response required, Reference, без expand |
| `giftCards` | Array(Object) | — |
| `group` | Meta | response required, Reference, без expand |
| `id` | UUID | readonly, response required |
| `meta` | Meta | response required |
| `moment` | DateTime | response required |
| `name` | String(255) | response required |
| `noCashSum` | Float | response required |
| `organization` | Meta | readonly, response required, Reference, без expand |
| `organizationAccount` | Meta | Reference, без expand |
| `organizationBranch` | Meta | readonly, Reference, без expand |
| `owner` | Meta | Reference, без expand |
| `payedSum` | Float | readonly, response required |
| `positions` | MetaArray | response required, Reference, без expand |
| `prepaymentCashSum` | Float | response required |
| `prepaymentNoCashSum` | Float | response required |
| `prepaymentQrSum` | Float | response required |
| `printed` | Boolean | readonly, response required |
| `published` | Boolean | readonly, response required |
| `qrSum` | Float | response required |
| `rate` | Object | response required |
| `retailShift` | Meta | Create required, response required, Reference, без expand |
| `retailStore` | Meta | response required, Reference, без expand |
| `sessionNumber` | Int | — |
| `shared` | Boolean | response required |
| `state` | Meta | Reference, без expand |
| `store` | Meta | Create required, response required, Reference, без expand |
| `sum` | Float | readonly, response required |
| `syncId` | UUID | — |
| `taxSystem` | Enum | response required |
| `updated` | DateTime | readonly, response required |
| `vatEnabled` | Boolean | response required |
| `vatIncluded` | Boolean | — |
| `vatSum` | Float | readonly |

### Позиции Розничной продажи

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retaildemand.md#L103).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `assortment` | Meta | response required, Reference, без expand |
| `cost` | Int | — |
| `declaration` | Array(Object) | readonly |
| `discount` | Float | response required |
| `id` | UUID | readonly, response required |
| `pack` | Object | — |
| `price` | Float | response required |
| `quantity` | Float | response required |
| `things` | Array(String) | — |
| `trackingCodes` | Array(Object) | — |
| `trackingCodes_1162` | Array(Object) | — |
| `vat` | Int | response required |
| `vatEnabled` | Boolean | response required |

### Подарочные сертификаты

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retaildemand.md#L136).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `name` | String | response required |
| `paymentSum` | Float | readonly, response required |

### Данные фискального чека

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retaildemand.md#L145).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `fnNumber` | String | — |
| `fiscalDocSign` | String | — |
| `fiscalDocNumber` | String | — |
| `kktRegNumber` | String | — |
| `time` | DateTime | — |

### Информация о прослеживаемости импортных товаров

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retaildemand.md#L160).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `gtd` | String | readonly |
| `rnpt` | String | readonly |
| `country` | Meta | readonly |
| `quantity` | Float | readonly |

### Коды маркировки товаров и транспортных упаковок

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retaildemand.md#L170).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `cis` | String | Create required, response required |
| `type` | Enum | Create required, response required |
| `trackingCodes` | Array(Object) | — |

### Коды маркировки товаров и транспортных упаковок в формате тега 1162

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retaildemand.md#L186).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `cis_1162` | String | readonly, response required |
| `type` | Enum | readonly, response required |
| `trackingCodes_1162` | Array(Object) | — |

### Получить Розничные продажи

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retaildemand.md#L215).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

### Получить позиции

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retaildemand.md#L2640).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

## retailsalesreturn

Ticket: E04. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retail_sales_return.md).
Минимальный Create / особенность: **name, organization, store, retailStore, retailShift, agent;
demand для based, cashSum/noCashSum для unbased**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retail_sales_return.md#L6).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `agent` | Meta | Create required, response required, Reference, без expand |
| `agentAccount` | Meta | Reference, без expand |
| `applicable` | Boolean | response required |
| `attributes` | Array(Object) | — |
| `cashSum` | Float | response required |
| `code` | String(255) | — |
| `contract` | Meta | Reference, без expand |
| `created` | DateTime | readonly, response required |
| `deleted` | DateTime | readonly |
| `demand` | Meta | Reference, без expand |
| `description` | String(4096) | — |
| `externalCode` | String(255) | response required |
| `group` | Meta | response required, Reference, без expand |
| `id` | UUID | readonly, response required |
| `meta` | Meta | response required |
| `moment` | DateTime | response required |
| `name` | String(255) | response required |
| `noCashSum` | Float | response required |
| `organization` | Meta | Create required, response required, Reference, без expand |
| `organizationAccount` | Meta | Reference, без expand |
| `organizationBranch` | Meta | readonly, Reference, без expand |
| `owner` | Meta | Reference, без expand |
| `positions` | MetaArray | response required, Reference, без expand |
| `printed` | Boolean | readonly, response required |
| `published` | Boolean | readonly, response required |
| `qrSum` | Float | response required |
| `rate` | Object | response required |
| `retailShift` | Meta | Create required, response required, Reference, без expand |
| `retailStore` | Meta | Create required, response required, Reference, без expand |
| `shared` | Boolean | response required |
| `state` | Meta | Reference, без expand |
| `store` | Meta | Create required, response required, Reference, без expand |
| `sum` | Float | readonly, response required |
| `syncId` | UUID | — |
| `taxSystem` | Enum | — |
| `updated` | DateTime | readonly, response required |
| `vatEnabled` | Boolean | response required |
| `vatIncluded` | Boolean | — |
| `vatSum` | Float | readonly |

### Позиции Розничного возврата

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retail_sales_return.md#L65).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `assortment` | Meta | response required, Reference, без expand |
| `cost` | Int | — |
| `discount` | Float | response required |
| `id` | UUID | readonly, response required |
| `pack` | Object | — |
| `price` | Float | response required |
| `quantity` | Float | response required |
| `things` | Array(String) | — |
| `vat` | Int | response required |
| `vatEnabled` | Boolean | response required |

### Получить Розничные возвраты

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retail_sales_return.md#L95).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

### Получить позиции

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retail_sales_return.md#L2490).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

## bonusprogram

Ticket: E05. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bonus_program.md).
Минимальный Create / особенность: **name, active, allProducts, allAgents,
earnRateRoublesToPoint, spendRatePointsToRouble, maxPaidRatePercents**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bonus_program.md#L8).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `active` | Boolean | response required |
| `agentTags` | Array(String) | response required |
| `allAgents` | Boolean | response required |
| `allProducts` | Boolean | response required |
| `earnRateRoublesToPoint` | Int | — |
| `earnWhileRedeeming` | Boolean | response required |
| `id` | UUID | readonly, response required |
| `maxPaidRatePercents` | Int | — |
| `meta` | Meta | response required |
| `name` | String(255) | — |
| `postponedBonusesDelayDays` | Int | — |
| `spendRatePointsToRouble` | Int | — |
| `welcomeBonusesEnabled` | Boolean | response required |
| `welcomeBonusesMode` | Enum | — |
| `welcomeBonusesValue` | Int | — |

### Получить все Бонусные программы

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bonus_program.md#L39).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

## bonustransaction

Ticket: E06. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bonus_operation.md).
Минимальный Create / особенность: **agent, bonusProgram, transactionType**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bonus_operation.md#L8).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `agent` | Meta | Create required, response required, Reference, без expand |
| `applicable` | Boolean | response required |
| `bonusProgram` | Meta | Reference, без expand |
| `bonusValue` | Int | — |
| `categoryType` | Enum | readonly |
| `code` | String(255) | — |
| `created` | DateTime | readonly, response required |
| `description` | String(4096) | — |
| `executionDate` | DateTime | — |
| `externalCode` | String(255) | response required |
| `group` | Meta | response required, Reference, без expand |
| `id` | UUID | readonly, response required |
| `meta` | Meta | response required |
| `moment` | DateTime | — |
| `name` | String(255) | — |
| `organization` | Meta | Reference, без expand |
| `owner` | Meta | Reference, без expand |
| `parentDocument` | Meta | Reference, без expand |
| `shared` | Boolean | response required |
| `transactionStatus` | Enum | readonly |
| `transactionType` | Enum | Create required, response required |
| `updated` | DateTime | response required |
| `updatedBy` | UID | — |

### Получить Бонусные операции

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bonus_operation.md#L47).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

### Создать Бонусную операцию

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bonus_operation.md#L277).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `agent` | Meta | Create required, response required, Reference, без expand |
| `bonusProgram` | Meta | Create required, response required, Reference, без expand |
| `transactionType` | Enum | Create required, response required |

## states

Ticket: A06. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_states.md).
Минимальный Create / особенность: **name, color, stateType**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_states.md#L10).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `accountId` | UUID | readonly, response required |
| `color` | Int | Create required, response required |
| `entityType` | String(255) | readonly, response required |
| `id` | UUID | readonly, response required |
| `meta` | Meta | readonly, response required |
| `name` | String(255) | Create required, response required |
| `stateType` | Enum | Create required, response required |

## files

Ticket: F02. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_files.md).
Минимальный Create / особенность: **filename, content в upload; owner/type/id отдельным
маршрутом**.

### Работа с Файлами в рамках отдельных Операции, Номенклатуры, Задачи или Контрагента

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_files.md#L6).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `filename` | String(255) | — |
| `content` | String | — |

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_files.md#L31).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `created` | DateTime | readonly, response required |
| `createdBy` | Meta | response required, Reference, без expand |
| `filename` | String(255) | response required |
| `meta` | Meta | response required |
| `miniature` | Meta | — |
| `size` | Int | response required |
| `tiny` | Meta | — |
| `title` | String(255) | response required |

### Получить список Файлов Операции, Номенклатуры, Задачи или Контрагента

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_files.md#L46).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

## audit

Ticket: A07/A08. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md).
Минимальный Create / особенность: **read-only; own pagination/envelopes**.

### Атрибуты сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md#L33).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `entityType` | Enum | readonly, response required |
| `eventType` | Enum | readonly, response required |
| `events` | Array(Object) | readonly, response required |
| `id` | UUID | readonly, response required |
| `info` | Object | readonly |
| `meta` | Meta | readonly, response required |
| `moment` | DateTime | readonly, response required |
| `objectCount` | Int | readonly, response required |
| `objectType` | Enum | response required |
| `source` | Enum | readonly, response required |
| `supportAccess` | Boolean | readonly, response required |
| `uid` | String(255) | readonly, response required |

### Получить Контексты

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md#L86).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

### Атрибуты сущности Событие

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md#L244).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `additionalInfo` | String(4096) | readonly |
| `audit` | Enum | readonly, response required |
| `diff` | Object | readonly, response required |
| `entity` | Meta | readonly |
| `entityType` | Enum | readonly, response required |
| `eventType` | Enum | readonly, response required |
| `moment` | DateTime | readonly, response required |
| `name` | String(255) | readonly, response required |
| `objectCount` | Int | readonly, response required |
| `objectType` | Enum | response required |
| `source` | Enum | readonly, response required |
| `supportAccess` | Boolean | readonly, response required |
| `uid` | String(255) | readonly, response required |

### Событие регистрации

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md#L272).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `account` | String(255) | readonly, response required |
| `country` | String(255) | readonly, response required |

### События публикации документов

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md#L279).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `templateName` | String(255) | readonly, response required |
| `publicationHref` | URL | readonly, response required |

### События отправки писем

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md#L286).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `senderEmail` | String(255) | readonly, response required |
| `targetEmail` | String(255) | readonly, response required |
| `subjectEmail` | String(255) | readonly, response required |
| `text` | String(255) | readonly, response required |

### События удаления сущностей

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md#L295).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `attributeName` | String(255) | readonly, response required |
| `oldValue` | Boolean | readonly, response required |

### События обновления сущностей, перемещения/восстановления из корзины, перемещение/восстановление из архива

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md#L302).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `attributeName` | String(255) | readonly, response required |
| `oldValue` | Boolean | readonly, response required |
| `newValue` | Boolean | readonly, response required |

### Получить События по Контексту

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md#L311).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

### Получить События по Сущности

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/audit/_audit.md#L393).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | — |
| `context` | Meta | — |
| `rows` | Array(Object) | — |

## report-turnover

Ticket: F03/F04/F05. [Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md).
Минимальный Create / особенность: **read-only; All/ByStore/ByOperations имеют разные options и
rows**.

### Атрибуты объекта отчета Обороты по товарам

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md#L10).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `assortment` | Object | response required |
| `onPeriodStart` | Object | response required |
| `onPeriodEnd` | Object | response required |
| `income` | Object | response required |
| `outcome` | Object | response required |

### Структура объекта assortment в отчете по товарам

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md#L20).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `article` | String(255) | — |
| `code` | String(255) | — |
| `image` | Object | — |
| `meta` | Meta | response required |
| `name` | String(255) | response required |
| `productFolder` | Object | — |
| `uom` | Object | — |

### Структура объекта показатели (onPeriodStart, onPeriodEnd, income, outcome)

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md#L32).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `sum` | Float | response required |
| `quantity` | Float | response required |

### Атрибуты доступные для фильтрации отчета Обороты по товарам

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md#L42).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `agent` | Object | — |
| `agentTag` | String(255) | — |
| `contract` | Object | — |
| `organization` | Object | — |
| `product` | Object | — |
| `project` | Object | — |
| `retailStore` | Object | — |
| `store` | Object | — |
| `supplier` | Object | — |
| `type` | Enum | — |
| `variant` | Object | — |
| `withoutturnover` | Boolean | — |
| `archived` | Boolean | — |

### Атрибуты объекта отчета Обороты по товару с детализацией по складам

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md#L269).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `assortment` | Object | response required |
| `stockByStore` | Object | response required |

### Структура объекта assortment в отчете по товару с детализацией по складам

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md#L276).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `article` | String(255) | — |
| `code` | String(255) | — |
| `image` | Object | — |
| `meta` | Meta | response required |
| `name` | String(255) | response required |
| `productFolder` | Object | — |
| `uom` | Object | — |

### Структура объекта детализация оборотов по складам

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md#L288).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `store` | Object | response required |
| `onPeriodStart` | Object | response required |
| `onPeriodEnd` | Object | response required |
| `income` | Object | response required |
| `outcome` | Object | response required |

### Структура объекта показатели (onPeriodStart, onPeriodEnd, income, outcome) в отчете по товару с детализацией по складам

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md#L298).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `sum` | Float | response required |
| `quantity` | Float | response required |

### Атрибуты доступные для фильтрации отчета обороты по товару и его модификациям с детализацией по складам

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md#L310).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `agent` | Object | — |
| `agentTag` | String(255) | — |
| `contract` | Object | — |
| `organization` | Object | — |
| `product` | Object | — |
| `project` | Object | — |
| `retailStore` | Object | — |
| `store` | Object | — |
| `type` | Enum | — |
| `variant` | Object | — |

### Атрибуты объекта отчета Обороты по товару с детализацией по документам

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md#L530).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `assortment` | Object | response required |
| `store` | Object | response required |
| `operation` | Object | response required |
| `quantity` | Float | response required |
| `cost` | Float | response required |
| `sum` | Float | response required |

### Структура объекта assortment в отчете по товару с детализацией по документам

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md#L541).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | response required |
| `name` | String(255) | response required |
| `code` | String(255) | — |
| `article` | String(255) | — |
| `productFolder` | Object | — |
| `uom` | Object | — |

### Структура объекта operation

[Таблица источника](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/reports/_report_turnover.md#L552).

| JSON | API type | Явные markers |
| --- | --- | --- |
| `meta` | Meta | response required |
| `name` | String(255) | response required |
| `description` | String(255) | — |
| `moment` | Object | response required |
| `agent` | Object | — |

## context/employee

A03: только GET context/employee без ID. Собственный read-only context, не ResponseContext.
[Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L2898).

| JSON | API type | Scope |
| --- | --- | --- |
| `accountId` | UUID | read-only |
| `archived` | Boolean | read-only |
| `attributes` | Array(Object) | read-only |
| `cashiers` | MetaArray | read-only |
| `code` | String(255) | read-only |
| `created` | DateTime | read-only |
| `description` | String(4096) | read-only |
| `email` | String(255) | read-only |
| `externalCode` | String(255) | read-only |
| `firstName` | String(255) | read-only |
| `fullName` | String(255) | read-only |
| `group` | Meta | read-only |
| `id` | UUID | read-only |
| `image` | Object | read-only |
| `inn` | String(255) | read-only |
| `lastName` | String(255) | read-only |
| `meta` | Meta | read-only |
| `middleName` | String(255) | read-only |
| `name` | String(255) | read-only |
| `owner` | Meta | read-only |
| `permissions` | Object | read-only |
| `phone` | String(255) | read-only |
| `position` | String(255) | read-only |
| `shared` | Boolean | read-only |
| `shortFio` | String(255) | read-only |
| `uid` | String(255) | read-only |
| `updated` | DateTime | read-only |

## context/application

A04: GET context/application; read meta (model.Meta), id и appUid (string).
Caller ID отсутствует. GET entity/application/{installationId} — другой endpoint.
[Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L3618).
