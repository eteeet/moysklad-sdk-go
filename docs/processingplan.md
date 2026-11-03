# ProcessingPlan

`entity/processingplan`: техкарта — модель `ProcessingPlan`, `CreateRequest`, `UpdateRequest`,
`NewService(client)`; методы `Get/List/Create/Update/Delete` и `CreateMany/UpdateMany/DeleteMany`.
Четыре вложенных набора собственной формы (не `positions.Position`): этапы `Stage`, материалы
`Material`, параметрические материалы `ParametricMaterial` и продукты `Product`.
Корзины, шаблона `New` и expand нет. Delete и DeleteMany удаляют безвозвратно.
[Общие правила методов, JSON и ошибок](documents.md).

Минимальный Create: непустые `name` и `products` — проверяются до запроса. Без
`processingProcess` сервер подставит «Основной техпроцесс».

```go
plan, err := s.Create(ctx, processingplan.CreateRequest{
	Name: "Пирог",
	Products: []processingplan.ProductRequest{{Assortment: model.Value(pie), Quantity: model.Value(1.0)}},
	Cost: model.Value(1000.0), // копейки
	Materials: model.Value([]processingplan.MaterialRequest{{
		Assortment: model.Value(flour), Quantity: model.Value(0.5),
	}}),
})
```

## Поддержанные поля

| JSON | Read Go type | Create | Update |
| --- | --- | --- | --- |
| `meta` | `model.Meta` | — | — |
| `id` | `string` | — | — |
| `accountId` | `string` | — | — |
| `name` | `string` | required | Field |
| `code` | `string` | Field | Field |
| `externalCode` | `string` | Field | Field |
| `archived` | `bool` | Field | Field |
| `shared` | `bool` | Field | Field |
| `updated` | `model.DateTime` | — | — |
| `pathName` | `string` | — | — |
| `cost` | `float64` | Field | Field |
| `costDistributionType` | `string` | — | — |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `parent` | `*model.Reference` | Field | Field |
| `processingProcess` | `*model.Reference` | Field | Field |
| `stages` | `*model.Collection[Stage]` | Field | Field |
| `materials` | `*model.Collection[Material]` | Field | Field |
| `parametricMaterials` | `*model.Collection[ParametricMaterial]` | Field | Field |
| `products` | `*model.Collection[Product]` | required | Field |

`cost` — стоимость производства, `float64` в копейках. `costDistributionType` — открытый enum
(`BY_PRICE`, `BY_PRODUCTION`), только чтение; `BY_PRODUCTION` возможен при двух и более продуктах.
`parent` — группа техкарт (`processingplanfolder`).

Наборы в теле Create/Update ни один не помечен «Только для чтения»: каждый — не больше 1000 строк
(1001 отклоняется до запроса), больше — методами наборов. В Update переданный набор заменяет весь
набор: строка с `ID` — существующая (id только адресует её), без — новая, не переданные удаляются.
Смена `processingProcess` удаляет материалы старого техпроцесса. Одновременно `cost` техкарты и
`cost` этапов сервер не принимает. В параметрической техкарте продукт один.

## Наборы

| Набор | Ресурс | meta.type | Методы |
| --- | --- | --- | --- |
| этапы | `…/{id}/stages` | `processingplanstages` | `GetStages`, `GetStage`, `UpdateStage` |
| материалы | `…/{id}/materials` | `processingplanmaterial` | `GetMaterials`, `GetMaterial`, `CreateMaterials`, `UpdateMaterial`, `DeleteMaterial`, `DeleteManyMaterials` |
| параметрические материалы | `…/{id}/parametricmaterials` | `processingplanparametricmaterial` (в теле массового удаления — `processingplanmaterial`) | `GetParametricMaterials`, `GetParametricMaterial`, `CreateParametricMaterials`, `UpdateParametricMaterial`, `DeleteParametricMaterial`, `DeleteManyParametricMaterials` |
| продукты | `…/{id}/products` | `processingplanresult` | `GetProducts`, `GetProduct`, `CreateProducts`, `UpdateProduct`, `DeleteProduct`, `DeleteManyProducts` |

- `Get…s(ctx, id, query.PageOptions)` — одна страница, только limit/offset (`fields` отклоняется до запроса).
- `Create…s` — один POST массивом, ответ `[]batch.Result[…]`; повторов нет.
- `Update…` принимают тот же request, что и тело техкарты (`StageRequest`, `MaterialRequest`,
  `ParametricMaterialRequest`, `ProductRequest`); пустой отклоняется до запроса. Ответ — объект или
  массив из одного объекта (так в примерах материалов и продуктов).
- `DeleteMany…` — POST `…/delete` с meta строк (1–1000); источник отвечает 200 с пустым телом,
  поэтому возвращается только error.

| Строка | Поля чтения | Request | Числа |
| --- | --- | --- | --- |
| `Stage` | `meta`, `id`, `accountId`, `enableHourAccounting`, `cost`, `labourCost`, `standardHour`, `standardHourCost`, `processingProcessPosition` | все, кроме `meta`, `accountId`, `standardHourCost` | `cost`, `labourCost`, `standardHourCost` — копейки; `standardHour` — нормо-часы |
| `Material` | `meta`, `id`, `accountId`, `assortment`, `product`, `quantity`, `processingProcessPosition`, `materialProcessingPlan` | все, кроме `meta`, `accountId`, `materialProcessingPlan` | `quantity` — количество, не деньги |
| `ParametricMaterial` | `meta`, `id`, `accountId`, `assortment`, `quantity`, `paramFeatureCharacteristic`, `paramQuantityCharacteristic`, `paramFeatures[]{characteristicValue, assortment}`, `paramQuantities[]{characteristicValue, quantity}`, `processingProcessPosition` | все, кроме `meta`, `accountId`, `characteristicValue` | `quantity` — количество |
| `Product` | `meta`, `id`, `accountId`, `assortment`, `product`, `quantity` | все, кроме `meta`, `accountId` | `quantity` — количество |

Все поля строк в request — `model.Field`; `ID` нужен только в Update техкарты. Для новой строки
сервер требует товар и количество. Этапы нельзя создавать и удалять: они строго соответствуют
этапам позиций техпроцесса; при `enableHourAccounting=true` сервер не принимает `labourCost`.
Без `processingProcessPosition` материал привязывается к первой позиции техпроцесса.

Параметрический материал — отдельная форма: материал выбирается по `paramFeatureCharacteristic`
(значения — `paramFeatures`), норма — по `paramQuantityCharacteristic` (значения —
`paramQuantities`). Конкретное значение без характеристики допустимо только для одного из двух
(`assortment` или `quantity`). Материал, не зависящий от характеристик, — обычный
`MaterialRequest` в `/materials`.

## Ограничения и противоречия источника

- `characteristicValue` таблица параметрических материалов помечает «Только для чтения», а примеры
  Create передают его в каждом элементе `paramFeatures`/`paramQuantities`. По правилу SDK поле не
  пишется (`ParamFeatureRequest` — только `assortment`, `ParamQuantityRequest` — только `quantity`),
  поэтому сопоставить материал или норму значению характеристики через SDK нельзя. Исключение
  требует решения владельца SDK.
- Таблица `stages` не помечает readonly, примеров этапов в теле техкарты нет — SDK пишет их по
  таблице; как сервер сопоставляет этапы без `ID`, источник не описывает.
- Для материалов и продуктов текст Create называет обязательными `product` и `quantity`, а примеры
  передают `assortment`: обязательность не проверяется до запроса, оба поля — `Field`.
- Путь параметрических материалов в запросах — `parametricmaterials`, в части href примеров —
  `parametricmaterial`; SDK использует путь запросов.
- meta.type этапов — `processingplanstages`, в новых примерах — `processingplanstage`; у строк
  параметрических материалов примеры дают и `processingplanparametricmaterial`, и
  `processingplanmaterial`. Тело массового удаления параметрических материалов — как в примере
  (`processingplanmaterial`).
- Пример «Получить продукт» отвечает массивом из одного объекта: `GetProduct` принимает и массив,
  и объект; массив другого размера — ошибка.
- Отдельного раздела «Получить параметрический материал» нет — пример лежит в разделе списка.
- Metadata техкарт ticket не называет — отдельного метода нет.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/processingplan),
`md/dictionaries/_processingplan.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
