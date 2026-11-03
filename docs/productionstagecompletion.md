# ProductionStageCompletion

`entity/productionstagecompletion`: модель `ProductionStageCompletion`, `CreateRequest`, `UpdateRequest`,
`NewService(client)`; методы `Get/List/Create/Update/Delete`, `CreateMany/UpdateMany/DeleteMany`
и `Trash(ctx, id)`. Два набора строк собственной формы (не `positions.Position`): материалы —
`Material`, продукты — `Product`.
[Общие правила методов, JSON и ошибок](documents.md), [корзина](documents-extra.md).

Минимальный Create: `productionStage` (непустой meta.href проверяется до запроса) и
`productionVolume` (передаётся всегда, значение проверяет сервер). Производство по
производственному заданию должно быть начато — иначе сервер отклоняет создание.

```go
doc, err := s.Create(ctx, productionstagecompletion.CreateRequest{
	ProductionStage:  stage, // ссылка на производственный этап (productionstage)
	ProductionVolume: 5,
	Performer:        model.Value(counterparty),
	Service:          model.Value(service),
})
```

## Поддержанные поля

| JSON | Read Go type | Create | Update |
| --- | --- | --- | --- |
| `meta` | `model.Meta` | — | — |
| `id` | `string` | — | — |
| `accountId` | `string` | — | — |
| `name` | `string` | Field | Field |
| `description` | `string` | Field | Field |
| `externalCode` | `string` | Field | Field |
| `created` | `model.DateTime` | — | — |
| `updated` | `model.DateTime` | — | — |
| `moment` | `model.DateTime` | Field | Field |
| `applicable` | `bool` | — | — |
| `shared` | `bool` | Field | Field |
| `defect` | `bool` | Field | — |
| `enableHourAccounting` | `bool` | Field | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `performer` | `*model.Reference` | Field | Field |
| `service` | `*model.Reference` | Field | Field |
| `productionStage` | `*model.Reference` | required | — |
| `productionVolume` | `float64` | required (always sent) | Field |
| `processingUnitCost` | `float64` | Field | Field |
| `labourUnitCost` | `float64` | Field | Field |
| `standardHourCost` | `float64` | Field | Field |
| `standardHourUnit` | `float64` | Field | Field |
| `materials` | `*model.Collection[Material]` | — | — |
| `products` | `*model.Collection[Product]` | — | — |

Деньги — `processingUnitCost`, `labourUnitCost`, `standardHourCost`: `float64` в копейках.
Не деньги: `productionVolume` (объём производства), `standardHourUnit` (нормо-часы),
`consumedQuantity`/`producedQuantity` строк — дробные количества без пересчёта.
`productionStage` и `defect` после создания не меняются — только в `CreateRequest`.

Особенности (проверяет сервер): `service` задают, только если исполнитель — контрагент (для
контрагента не из списка исполнителей этапа — обязательно); при `enableHourAccounting=true`
`labourUnitCost = standardHourCost * standardHourUnit` вычисляется сервером, и его изменение
отклоняется. Create/Update сдвигают дату старта производства, если она позже даты этапа.

List: одна страница, `fields` не принимается. Источник называет фильтры
`productionTask=<href задания>` и `assortment=<href>` (материал или продукт) — передаются через
`query.ListOptions.Filter`. DeleteMany удаляет в порядке `ids`: выполнение последнего этапа
должно идти раньше выполнений, от которых оно зависит.

## Материалы и продукты

| Набор | Ресурс | meta.type | Методы |
| --- | --- | --- | --- |
| материалы | `…/{id}/materials` | `productionstagecompletionmaterial` | `GetMaterials`, `CreateMaterial`, `UpdateMaterial` |
| продукты | `…/{id}/products` | `productionstagecompletionresult` | `GetProducts`, `UpdateProduct` |

- `GetMaterials(ctx, id, query.PageOptions) (*model.Collection[Material], error)` и
  `GetProducts(ctx, id, query.PageOptions) (*model.Collection[Product], error)` — одна страница,
  только limit/offset (`fields` отклоняется до запроса).
- `CreateMaterial(ctx, id, MaterialCreateRequest) (*Material, error)` — один POST объектом;
  `assortment` обязателен (до запроса), `consumedQuantity` сервер игнорирует при `things`.
- `UpdateMaterial(ctx, id, materialID, MaterialUpdateRequest)` и
  `UpdateProduct(ctx, id, productID, ProductUpdateRequest)` — PUT строки; пустой request
  отклоняется до запроса. У продукта товар/модификацию сменить нельзя, партию — можно.
- Поля строки: `meta`, `id`, `accountId` (только чтение), `assortment`, `consumedQuantity` или
  `producedQuantity`, `things`. Продукты есть только у последнего этапа.

## Ограничения и противоречия источника

- Источник не описывает получение одной строки, удаление строк, массовые операции строк и
  добавление продуктов — этих методов нет. Шаблона `New`, metadata и expand нет (ticket их не называет).
- `materials` и `products` в таблице не помечены «Только для чтения», но описаны как MetaArray без
  формы записи в теле документа и без примеров; SDK их в Create/Update не пишет — строки меняются
  методами наборов.
- Таблица не помечает `productionStage` и `productionVolume` «Обязательное при создании»; раздел
  «Создать» называет оба обязательными — SDK следует разделу «Создать».
- Примеры списка содержат `standardHour` (нет в таблице) — SDK его не читает; поле таблицы —
  `standardHourUnit`.
- Второй пример «Изменить материал» отвечает meta.type `productionstagecompletionresult`, href
  `…/products/…` и `entity/materials/…` у assortment; ответ примера добавления с серийными номерами кладёт `things`
  внутрь `assortment`. Модель следует таблицам строк (`things` — поле строки).
- Пример запроса списка продуктов использует href с завершающим `/` (`…/products/`); SDK
  запрашивает `…/products` без него.
- Единица денежных полей в источнике не названа; SDK следует общему правилу API (копейки).
- Создание с серийными номерами сервер не поддерживает; `things` принимают добавление и
  изменение строк.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionStageCompletion),
`md/documents/_productionStageCompletion.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
