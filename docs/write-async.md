# Удаление, массовые операции и асинхронные задачи

Источник: api-remap-1.2-doc, revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd` — `_general.md`,
`_async.md`, `_restrictions.md`, `_errors.md` и разделы сущностей. Все операции этой страницы
пишут или создают задачу: transport их не повторяет даже с `Config.Retry`.

## Delete

```go
err := c.Entity().Product().Delete(ctx, id)
err = c.Entity().Metadata().DeleteState(ctx, "customerorder", stateID)
err = c.Entity().Files().Delete(ctx, "product", productID, fileID)
```

- `Delete(ctx, id, opts...)` есть у всех сервисов `entity/*`, кроме assortment и pricetype: у
  них раздела «Удалить» в источнике нет. Роль — пользовательская (`entity/role`).
- Один `DELETE`, тела нет; успех — 200 или 204 (примеры источника — 200 с пустым телом). Неверный
  UUID отклоняется до запроса; права и связанные объекты проверяет сервер (APIError).
- Delete — не корзина. Удаление документа в корзину (`POST …/trash`) и восстановление из неё
  в SDK нет. Удаление позиций документов, компонентов комплекта, счетов контрагента и юрлица —
  тоже нет.
- `files.Delete`: по источнику сервер удаляет первый найденный файл с этим id у операции,
  номенклатуры, задачи или контрагента.

## Массовые операции

```go
res, err := c.Entity().Product().CreateMany(ctx, []product.CreateRequest{a, b})
res, err = c.Entity().Product().UpdateMany(ctx, []batch.Update[product.UpdateRequest]{{ID: id, Request: u}})
deleted, err := c.Entity().Product().DeleteMany(ctx, []string{id1, id2})
deleted, err = c.Entity().Assortment().DeleteMany(ctx, []model.Meta{row.Product.Meta, row.Service.Meta})

for _, r := range res {
	if r.Errors != nil { /* ошибка этого элемента */ }
	_ = r.Value // записанный объект
}
```

- Один вызов — один `POST` с массивом от 1 до `batch.MaxItems` (1000) элементов: 0 и 1001
  отклоняются до запроса, SDK массив не делит. Create и update в одном запросе не смешиваются.
- `err` — ошибка всего запроса: сеть, HTTP-статус, лимит, разбор. Ошибки отдельных элементов —
  в `Result.Errors`/`Deleted.Errors`; остальные элементы при этом могут быть записаны. Успех
  вызова не значит успех всех элементов.
- Массив разбирается только у HTTP 200. Статус ответа с ошибками отдельных элементов источник не
  указывает (фикстура `batch/product_batch_partial.json` — без статуса); при другом статусе исход
  элементов неизвестен.
- Порядок элементов ответа источник не описывает: в примерах товаров он совпадает с запросом,
  в примере удаления ассортимента — нет. Сверяйте объекты по id, удаления — по `Info`.
- Обработка ограничена 5 минутами (ошибка 1088). После таймаута, обрыва или отмены ctx часть
  элементов может быть записана: повторять запрос вслепую нельзя.
- До запроса проверяются только количество, UUID у UpdateMany/DeleteMany и непустой Request у
  UpdateMany. Обязательные поля и значения, которые одиночные Create/Update проверяют локально,
  в массовых методах проверяет сервер у каждого элемента.
- Элемент UpdateMany и DeleteMany — `meta` объекта: href, metadataHref, type, mediaType от
  `Config.BaseURL` (`transport.Client.Href`), как в примерах источника; у комплектов и услуг
  metadataHref — `entity/product/metadata`.

Набор методов — по разделам источника «Массовое создание и обновление …» и «Массовое удаление …»:

| Сервис | CreateMany | UpdateMany | DeleteMany |
| --- | --- | --- | --- |
| bonusprogram | — | — | да |
| employee | — | да | да |
| retailsalesreturn | да | да | — |
| role, pricetype | — | — | — |
| assortment | — | — | да, по `model.Meta` разных типов |
| остальные 29 сервисов `entity/*` | да | да | да |

Раздел «Массовое обновление Сотрудников» называет только обновление, хотя второй элемент его
примера — без meta (создание); SDK следует заголовку. Массовое создание характеристик и статусов
не входит.

## Асинхронные задачи

```go
task, err := c.Entity().Counterparty().ListAsync(ctx, query.ListOptions{Filter: f})
done, err := c.Async().Wait(ctx, task.ID, 5*time.Second)
if done.State == async.StateDone {
	_, err = c.Async().Result(ctx, task.ID, file) // JSON того же вида, что синхронный ответ
}
```

- Создают задачи только `ListAsync` контрагентов и ассортимента (новый endpoint, с beta-тегом):
  из запросов SDK источник разрешает async лишь им. Async-отчёты (остатки, прибыльность, деньги,
  показатели) в SDK не входят. Limit и Offset с async источник запрещает — ошибка до запроса.
- Создание задачи — побочный эффект и у GET: `async=true` transport не повторяет. Ответ 202 даёт
  `async.Submitted`: ID, адрес статуса (Content-Location, …/async/{id}) и результата (Location,
  …/async/{id}/result); другой ответ — ошибка. После сбоя сети или отмены ctx задача могла
  создаться: проверьте `c.Async().List` с фильтром по `request`, прежде чем создавать её снова.
- Async ассортимента идёт с тем же beta-тегом, что List: источник ссылается на раздел нового
  endpoint, но отдельно async с beta-тегом не подтверждает.
- `c.Async()`: `Get`, `List` (статусы за неделю), `Wait` (опрос Get до завершения; срок задаёт
  только ctx — без дедлайна опрос не кончится, пока задача не завершится), `Result` (скачивание
  через `transport.Client.Download`: redirect на ссылку без Authorization, ссылка живёт 5 минут) и
  `Cancel` (PUT, 204; останавливает подготовку результата, не откат; завершённая задача — 61007).
- Ошибки async — APIError с кодами 61000–61008: 61001 — limit/offset (SDK отклоняет их до
  запроса), 61002 — очередь заполнена, 61003 — результат удалён, 61006 — результат не готов.
- Статусы: PENDING, PROCESSING, DONE, ERROR, CANCEL, API_ERROR; набор открыт, неизвестный не
  считается завершением. У API_ERROR ошибки — в `Task.Errors`, а `Result` возвращает ошибку
  исходного запроса как APIError. Таблица называет `errors` типом Object, пример — массив; SDK
  следует примеру.
- В очереди аккаунта не больше 4 задач — лимит сервера на весь аккаунт; один Client и Limiter
  его не считают. Результат хранится 1 час (до `DeletionDate`); пока он готовится, возможны дубли
  и пропуски строк.
- Вебхуки задач не входят.
