Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/03-error-codes.md; ../../sdk-v1-foundation/issues/08-rate-concurrency.md; ../../sdk-v1-foundation/issues/10-size-limits.md

# X02: Batch create/update/delete

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Resource-specific arrays max1000, mixed per-item results/errors и timeout5min. Explicit batch
methods, без autosplit и без иллюзии транзакционного all-or-nothing. Удаление POST .../delete
отдельно от single DELETE.

## Специальные проверки

1000/1001, canceled/timeout ambiguous outcome, response item errors и metadata, не повторять
неизвестный результат; batch status не равен успеху всех элементов.

## Acceptance criteria

- [ ] Scope/публичный interface и field/method matrix утверждены; зависимости resolved
  и их код присутствует в рабочем дереве. Readonly, required response/Create, null,
  immutable и выбранные TODO разделены по источникам; противоречия не замолчаны.
- [ ] Реализован весь What to build, без неоговорённых методов и hidden requests.
- [ ] Специальные проверки этого ticket реализованы через публичный interface.
- [ ] JSON обоими codecs и независимые golden; upstream/provenance не изменены.
  Для transport/iterator/CI tickets вместо entity JSON применяются их собственные
  публичные seams и контрактные tests, а не фиктивные модели.
- [ ] Реальные transport/httptest tests покрывают успех и ошибки; context и
  metadata сохраняются, caller-controlled retry/walker не подменяет один вызов.
- [ ] go test -race, go vet ./... ./testdata, gofmt и diff-check PASS; независимые
  Standards/Spec review, исправления и финальный BASE/результат записаны в Answer.
- [ ] README/guide/status отражают реальный scope, ограничения и неизвестные факты.
  Planning без команды не коммитится; один code commit по ticket.

## Источники

[транспортный research](../../sdk-v1/research/transport.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `c0c5d1c`, коммит `6a3bbbc`. Решение пользователя 2026-11-02: раздельные методы
`CreateMany(ctx, []CreateRequest)`, `UpdateMany(ctx, []batch.Update[UpdateRequest])`,
`DeleteMany(ctx, ids)`; результат по элементам (Value или Errors), err — только ошибка всего запроса.

- Matrix по разделам «Массовое создание и обновление …», «Массовое обновление Сотрудников» и
  «Массовое удаление …»: все три метода у 29 сервисов; bonusprogram — DeleteMany; employee —
  UpdateMany и DeleteMany; retailsalesreturn — CreateMany и UpdateMany; role, pricetype — нет;
  assortment — `DeleteMany(ctx, []model.Meta)` для объектов разных типов (нашёл review Spec).
  Матрица закреплена `TestBatchOnlyWhereDocumented`.
- Пакет `batch`: MaxItems=1000, Update, Result (UnmarshalJSON: errors → Errors, иначе T), Deleted.
  `internal/write`: один POST, 0 и 1001 элементов — ошибка до запроса, без разбиения; UpdateMany и
  DeleteMany — meta {href, metadataHref, type, mediaType} от BaseURL, у bundle и service metadataHref —
  entity/product/metadata, как в примерах. Нужен публичный `transport.Client.Href`: BaseURL приватный.
- До запроса проверяются только количество, UUID и непустой Request; обязательные поля и значения
  одиночных Create/Update массовые методы оставляют серверу (ошибки элементов). Решение записано в godoc.
- Массив разбирается только у HTTP 200; статус ответа с ошибками отдельных элементов источник не
  указывает — при другом статусе исход неизвестен. Порядок элементов ответа источник не описывает
  (пример удаления ассортимента — в другом порядке). 1088, обрыв и отмена — неизвестный исход, POST не
  повторяется. Записано в doc пакета batch и docs/write-async.md.
- Тесты: `entity/batch_test.go` — форма запроса и meta, результаты по элементам, опции и ResponseInfo у
  каждого метода; APIError фикстуры batch_global_error, 503 без повтора, 201, короткий массив, null,
  объект; отмена до отправки и во время запроса, 1088 при статусе 500 — один вызов; 0/1001/плохой UUID/
  пустой Request без запроса, 1000 — одним запросом; DeleteMany ассортимента. `batch/batch_test.go` —
  фикстура product_batch_partial и пример удаления обоими кодеками, `{"errors":[]}`, не объект.
  `transport` — TestHref.
- Review Standards: комментарий UpdateMany уточнён, `ref`→`objectMeta`, `want`→`sent`, godoc
  перенесены по строкам, тест `errors:[]`. Оставлено: копия validID в internal/write (правило
  coverage.md); meta склеивается с JSON Request — у UpdateRequest нет поля meta. Review Spec:
  DeleteMany ассортимента, metadataHref bundle/service, тесты отмены во время запроса и 1088,
  оговорки статуса и порядка ответа.
- Расхождение источника: пример «Массовое обновление Сотрудников» содержит элемент без meta
  (создание); SDK следует заголовку — только UpdateMany. Не проверено: живой API, статус ответа
  с ошибками элементов.
