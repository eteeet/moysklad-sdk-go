# ProductionTask

`entity/productiontask`: модель `ProductionTask`, `CreateRequest`, `UpdateRequest`, `NewService(client)`;
методы `Get/List/Create/Update/Delete`, `CreateMany/UpdateMany/DeleteMany` и `Trash(ctx, id)`.
Четыре вложенных набора собственной формы (не `positions.Position`): позиции (техкарты задания),
продукты, производственные этапы и материалы этапа.
[Общие правила методов, JSON и ошибок](documents.md), [корзина](documents-extra.md).
Шаблона `New` нет: источник его не описывает.

Минимальный Create: `organization`, `materialsStore`, `productsStore` — проверяются до запроса
(непустой meta.href). Без `productionStart` задание создаётся не начатым; передача `productionStart`
в Create или Update начинает производство (в Update — меняет дату; дату позже выполненного этапа
сервер отклоняет).

```go
task, err := s.Create(ctx, productiontask.CreateRequest{
	Organization: organization, MaterialsStore: store, ProductsStore: store,
	Reserve: model.Value(false),
	ProductionRows: model.Value([]productiontask.ProductionRowRequest{{
		ProcessingPlan: model.Value(plan), ProductionVolume: model.Value(10.0),
	}}),
})
```

Метаданные (`entity/productiontask/metadata`, доп. поля) — общий пакет `metadata`:
`metadata.Service.Get(ctx, "productiontask")`, `ListAttributes`/`GetAttribute`.

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
| `created` | `model.DateTime` | — | — |
| `updated` | `model.DateTime` | — | — |
| `deleted` | `model.DateTime` | — | — |
| `moment` | `model.DateTime` | Field | Field |
| `deliveryPlannedMoment` | `model.DateTime` | Field | Field |
| `productionStart` | `model.DateTime` | Field | Field |
| `productionEnd` | `model.DateTime` | — | — |
| `applicable` | `bool` | Field | Field |
| `awaiting` | `bool` | Field | Field |
| `reserve` | `bool` | Field | Field |
| `shared` | `bool` | Field | Field |
| `printed` | `bool` | — | — |
| `published` | `bool` | — | — |
| `organization` | `*model.Reference` | required | Field |
| `materialsStore` | `*model.Reference` | required | Field |
| `productsStore` | `*model.Reference` | required | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `attributes` | `[]model.Attribute` | Field | Field |
| `productionRows` | `*model.Collection[ProductionRow]` | Field | Field |
| `products` | `*model.Collection[Product]` | Field | Field |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `customerOrders`, `demands`, `emissionOrders`, `internalOrders`, `moves`, `productionTasks`, `productionTaskSupplies`, `purchaseOrders`, `supplies` | `[]model.Reference` | Field | Field |

Денежных полей у документа нет. Флаги `applicable`, `awaiting`, `reserve`, `shared` — `model.Field[bool]`:
`model.Value(false)` отправляется явно, незаданный Field не отправляется.

`productionRows` в Create/Update — не больше 200 строк, `products` — не больше 1000 (больше
отклоняется до запроса). В Update массив заменяет весь набор: строка с `ID` — существующая (id только
адресует её), без — новая, отсутствующие удаляются.

## Вложенные наборы

| Набор | Ресурс | meta.type | Методы |
| --- | --- | --- | --- |
| позиции | `entity/productiontask/{id}/productionrows` | `productionrow` | `GetProductionRows`, `GetProductionRow`, `UpdateProductionRow`, `DeleteProductionRow` |
| продукты | `entity/productiontask/{id}/products` | `productiontaskresult` | `GetProducts`, `GetProduct`, `CreateProducts`, `UpdateProduct`, `DeleteProduct`, `DeleteManyProducts` |
| этапы | `entity/productionstage?filter=productionTask=<href задания>`, `entity/productionstage/{stageID}` | `productionstage` | `GetStages`, `UpdateStage` |
| материалы этапа | `entity/productionstage/{stageID}/materials` | `productiontaskmaterial` | `GetMaterials`, `CreateMaterial`, `UpdateMaterial`, `DeleteMaterial` |

- Списки (`GetProductionRows`, `GetProducts`, `GetStages`, `GetMaterials`) принимают
  `query.PageOptions` — одна страница, только limit/offset (`fields` отклоняется до запроса).
  `GetStages(ctx, id, page)` ставит обязательный фильтр `productionTask=<href задания>`; href строится
  от BaseURL клиента.
- Позиция: `meta`, `id`, `accountId` (только чтение), `name`, `externalCode`, `processingPlan`,
  `productionVariant`, `productionVolume` (объём, не деньги), `updated`. `ProductionRowRequest` —
  позиция в теле документа; `UpdateProductionRow` принимает `ProductionRowUpdate` только с
  `productionVolume` («Изменить можно только объём производства»).
- Продукт: `meta`, `id`, `accountId`, `assortment`, `planQuantity` (не деньги), `productionRow`
  (только чтение). `CreateProducts` — один POST массивом, ответ `[]batch.Result[Product]`;
  `DeleteManyProducts` — POST `…/products/delete` с meta (1–1000), источник отвечает 200 с пустым
  телом, поэтому возвращается только error.
- Этап: `stage`, `productionRow`, `materialStore`, `files`, `instruction`, `orderingPosition`,
  `totalQuantity`, `completedQuantity`, `availableQuantity`, `blockedQuantity`, `skippedQuantity` —
  только чтение; `StageUpdateRequest` пишет `enableHourAccounting` (`model.Field[bool]`, явный false
  уходит), `processingUnitCost`, `labourUnitCost`, `standardHourCost`, `standardHourUnit`,
  `plannedEndDate` и `materials` (заменяет все материалы этапа). При `enableHourAccounting=true`
  сервер считает `labourUnitCost = standardHourCost * standardHourUnit` и не принимает
  `labourUnitCost`. Этапы создаёт сервер; создать или удалить этап нельзя.
- `processingUnitCost`, `labourUnitCost`, `standardHourCost` — денежные `float64` без пересчёта;
  `standardHourUnit` (нормо-часы) и количества этапа — не деньги.
- Материал: `meta`, `id`, `accountId`, `assortment`, `planQuantity`. `CreateMaterial` — POST одним
  объектом; ответ — объект или массив из одного объекта (как в примере источника).
- Пустой request в `Update*` отклоняется до запроса. `CreateMaterial` и `CreateProducts` требуют
  у новой строки `assortment` с непустым meta.href и `planQuantity` (раздел «Создать»: «Для
  успешного создания необходимо…») — без них отклоняется до запроса; положительность количества
  проверяет сервер.

## Ограничения и противоречия источника

- Таблица атрибутов помечает «Необходимо при создании» только `organization`, раздел «Создать
  Производственное задание» — ещё `materialsStore` и `productsStore`. SDK требует все три.
- Раздел «Создать продукт» требует `productionRow`, а таблица продукта помечает его «Только для
  чтения» — SDK `productionRow` не пишет (ни в `CreateProducts`, ни в продуктах тела документа).
  Если сервер откажет без `productionRow`, нужно решение об исключении.
- Пример «Создать продукт» отправляет POST на `…/{id}/product`, а получение, изменение, удаление,
  массовое удаление и meta.href коллекции — `…/{id}/products`. SDK отправляет на `…/products`.
- Источник описывает создание позиций только в теле документа: отдельного «Создать позицию» нет,
  поэтому метода нет. Текст говорит «для дальнейшего пополнения позиций — специальный ресурс», но
  такой операции ресурс позиций не описывает.
- Поле `updated` позиции в таблице без маркера «Только для чтения», но это момент обновления,
  который ставит сервер, — SDK его не пишет.
- Поля статуса у производственного этапа источник не описывает (ни в таблице, ни в примерах) —
  в модели его нет; неизвестные поля ответа игнорируются.
- Единицу денежных полей этапа источник не называет; SDK передаёт значение как есть.
- Получения одного этапа и одного материала, создания/удаления этапа и массовых операций материалов
  источник не описывает — методов нет. Expand вне SDK.
- Delete удаляет безвозвратно (с выполненными этапами задания), Trash — в корзину.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/productionTask),
`md/documents/_productionTask.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
