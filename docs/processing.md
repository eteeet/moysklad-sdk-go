# Processing

`entity/processing`: модель `Processing`, `CreateRequest`, `UpdateRequest`, `NewService(client)`;
методы `Get/List/Create/Update/Delete`, `CreateMany/UpdateMany/DeleteMany`, шаблон `New(ctx, basis)`
и `Trash(ctx, id)`. Два набора строк одной собственной формы (не `positions.Position`): материалы и
продукты — `Position`/`PositionRequest`.
[Общие правила методов, JSON и ошибок](documents.md), [шаблон и корзина](documents-extra.md).

Минимальный Create: `organization`, `materialsStore`, `productsStore` — проверяются до запроса
(непустой meta.href). Дальше сервер требует один из двух вариантов.

```go
// С техкартой: продукты, материалы и затраты сервер берёт из техкарты.
doc, err := s.Create(ctx, processing.CreateRequest{
	Organization: organization, MaterialsStore: store, ProductsStore: store,
	ProcessingPlan: model.Value(plan),
	Quantity:       model.Value(3.0), // по умолчанию 1
})

// Без техкарты: products и processingSum обязательны.
doc, err = s.Create(ctx, processing.CreateRequest{
	Organization: organization, MaterialsStore: store, ProductsStore: store,
	ProcessingSum: model.Value(5000.0), // копейки за единицу объёма производства
	Products: model.Value([]processing.PositionRequest{{
		Assortment: model.Value(product), Quantity: model.Value(10.0),
	}}),
})
```

## Создание и изменение: что делает сервер

Поля `quantity`, `processingSum`, `products`, `materials` связаны с техкартой. Техоперация должна
содержать хотя бы один продукт. SDK эти условия не проверяет — их проверяет сервер.

- **С техкартой** (`processingPlan`): незаданные `products`, `materials`, `processingSum` создаются
  по техкарте. `quantity` по умолчанию 1; другое значение пересчитывает количества продуктов и
  материалов кратно ему. Свои продукты/материалы — передать заполненные `Products`/`Materials`.
  Чтобы не брать материалы техкарты, передайте пустой набор явно:
  `Materials: model.Value([]processing.PositionRequest{})` (уходит `"materials":[]`); незаданное
  поле сервер возьмёт из техкарты.
- **Без техкарты**: обязательны `products` и `processingSum`; `quantity` по умолчанию 1.
- **Update**: новая техкарта заменяет продукты, материалы и затраты значениями техкарты, если не
  переданы `products`, `materials`, `processingSum`. Пустые материалы — явный пустой массив. Без
  `quantity` объём не меняется; количества строк сервер пересчитывает кратно объёму.

Materials и Products передаются в теле Create/Update массивом строк (так в примерах «Создать»).
Не больше 1000 строк в каждом наборе (1001 отклоняется до запроса); в Update массив заменяет весь
набор: строка с `ID` — существующая (id только адресует её), без — новая.

## Поддержанные поля

| JSON | Read Go type | Create | Update |
| --- | --- | --- | --- |
| `meta` | `model.Meta` | — | — |
| `id` | `string` | — | — |
| `accountId` | `string` | — | — |
| `name` | `string` | Field | Field |
| `description` | `string` | Field | Field |
| `code` | `string` | Field | Field |
| `externalCode` | `string` | Field | Field |
| `syncId` | `string` | — | — |
| `created` | `model.DateTime` | — | — |
| `updated` | `model.DateTime` | — | — |
| `deleted` | `model.DateTime` | — | — |
| `moment` | `model.DateTime` | Field | Field |
| `applicable` | `bool` | Field | Field |
| `shared` | `bool` | Field | Field |
| `printed` | `bool` | — | — |
| `published` | `bool` | — | — |
| `organization` | `*model.Reference` | required | Field |
| `organizationAccount` | `*model.Reference` | Field | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `project` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `attributes` | `[]model.Attribute` | Field | Field |
| `processingPlan` | `*model.Reference` | Field | Field |
| `processingSum` | `float64` | Field | Field |
| `quantity` | `float64` | Field | Field |
| `materialsStore` | `*model.Reference` | required | Field |
| `productsStore` | `*model.Reference` | required | Field |
| `materials` | `*model.Collection[Position]` | Field | Field |
| `products` | `*model.Collection[Position]` | Field | Field |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `processingOrder` | `*model.Reference` | Field | Field |

`processingSum` — затраты на производство за единицу объёма в копейках, `float64` (таблица: Int,
примеры передают `50.0`). `quantity` — объём производства, не деньги. `processingPlan` — ссылка на
техкарту (`entity/processingplan`), `processingOrder` — на заказ на производство.

## Материалы и продукты

| Набор | Ресурс | meta.type | Методы |
| --- | --- | --- | --- |
| материалы | `…/{id}/materials` | `processingpositionmaterial` | `GetMaterials`, `GetMaterial`, `CreateMaterial`, `UpdateMaterial`, `DeleteMaterial`, `DeleteManyMaterials` |
| продукты | `…/{id}/products` | `processingpositionresult` | `GetProducts`, `GetProduct`, `CreateProduct`, `UpdateProduct`, `DeleteProduct`, `DeleteManyProducts` |

- `GetMaterials`/`GetProducts(ctx, id, positions.PageOptions) (*model.Collection[Position], error)` —
  одна страница, только limit/offset (`fields` отклоняется до запроса).
- `CreateMaterial`/`CreateProduct(ctx, id, PositionRequest) (*Position, error)` — один POST одной
  строкой-объектом, как в примере источника; ответ — объект или массив из одного объекта.
  `assortment` (непустой meta.href) и `quantity` обязательны — проверяются до запроса.
- `UpdateMaterial`/`UpdateProduct` — обязательных полей нет, пустой request отклоняется до запроса.
- `DeleteManyMaterials`/`DeleteManyProducts` — POST `…/delete` с meta строк (1–1000); источник
  отвечает 200 с пустым телом, поэтому возвращается только error.
- Поля строки: `meta`, `id`, `accountId` (только чтение), `assortment`, `quantity`, `slot` (ячейка
  склада). Удаляя продукт, помните: техоперация должна сохранить хотя бы один продукт (проверяет сервер).

## Шаблон

`New(ctx, basis)` — PUT `entity/processing/new`; основание по источнику —
`{"processingOrder": ref}` или `{"processingPlan": ref}`, `nil` — без основания. Ответ — `Processing`
без id; `materials`/`products` в нём приходят только с `rows` (без meta).

## Ограничения и противоречия источника

- `syncId` таблица помечает «Только для чтения» — SDK его не пишет.
- `processingSum` в таблице Int, в примерах дробный (`50.0`, `150.0`) — SDK держит `float64`.
- `processingOrder` (раздел «Связи с другими документами») не помечен readonly и пишется по общему
  правилу привязки документов.
- `files` пишется отдельным API файлов — в запросах документа его нет.
- Пример «Массовое создание» передаёт `products`/`materials` формой `{meta, rows}`, примеры «Создать»
  и «Изменить» — массивом; SDK отправляет массив.
- Пример изменения документа адресует существующую строку через `meta`; SDK — через `id`
  (общее правило строк в Update).
- Поле `things` у строк есть в OpenAPI-фикстуре `testdata/upstream/processing.json`, но не в таблице
  источника — SDK его не читает.
- Метаданные (`entity/processing/metadata`) — общий пакет `metadata`; «Отдельное доп. поле», expand
  и массовое создание строк ticket не называет — их нет. Delete удаляет безвозвратно, Trash — в корзину.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/processing),
`md/documents/_processing.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
