Status: needs-triage
Stage: DRAFT

# 4. Производство

## Problem Statement

Добавить производственные справочники и документы: техкарты, техпроцессы, этапы, техоперации, заказы на производство и производственные задания.

## Solution

Завершённые вертикальные срезы; каждый ticket включает models, methods, тесты и собственное
описание полей. Матрица ниже определяет scope этой захватки.
[Общая карта](../sdk-full/map.md), [нормативные ограничения v1](../sdk-v1/spec.md),
[каталог JSON-контрактов](../sdk-v1/contracts.md), [backlog](../../docs/status-backlog.md).
Источник: api-remap-1.2-doc revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.

## User Stories

1. Как разработчик решения, я могу использовать «ProcessingPlanFolder» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
2. Как разработчик решения, я могу использовать «ProcessingStage» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
3. Как разработчик решения, я могу использовать «ProcessingProcess + positions» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
4. Как разработчик решения, я могу использовать «ProcessingPlan + materials/products/stages» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
5. Как разработчик решения, я могу использовать «Processing + materials/products» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
6. Как разработчик решения, я могу использовать «ProcessingOrder + positions» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
7. Как разработчик решения, я могу использовать «ProductionTask + materials/products/stages» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
8. Как разработчик решения, я могу использовать «ProductionStageCompletion» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.

## Implementation Decisions и scope tickets

| ID | Ticket | Blocked by |
| --- | --- | --- |
| M01 | [ProcessingPlanFolder](issues/01-processingplanfolder.md) | нет |
| M02 | [ProcessingStage](issues/02-processingstage.md) | нет |
| M03 | [ProcessingProcess + positions](issues/03-processingprocess.md) | M02 |
| M04 | [ProcessingPlan + materials/products/stages](issues/04-processingplan.md) | M01, M03 |
| M05 | [Processing + materials/products](issues/05-processing.md) | M04, O02, O04 |
| M06 | [ProcessingOrder + positions](issues/06-processingorder.md) | M04, O01, O02, O04 |
| M07 | [ProductionTask + materials/products/stages](issues/07-productiontask.md) | M04, O04 |
| M08 | [ProductionStageCompletion](issues/08-productionstagecompletion.md) | M07, O04 |

### M01: ProcessingPlanFolder

Группа техкарт: Get/List/Create/Update/Delete, массовые, metadata.

Проверки и особенности: Родительская группа и pathName — как у productfolder.

Источники: [processingplanfolder](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_processingplanfolder.md).

### M02: ProcessingStage

Этап производства: Get/List/Create/Update/Delete, массовые.

Проверки и особенности: allPerformers/performers — флаг и ссылки.

Источники: [processingstage](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_processingstage.md).

### M03: ProcessingProcess + positions

Техпроцесс: Get/List/Create/Update/Delete, массовые; позиции (этапы) List/Get/Create/Update/Delete, массовое удаление.

Проверки и особенности: Порядок этапов и nextPositions — по таблице.

Источники: [processingprocess](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_processingprocess.md).

### M04: ProcessingPlan + materials/products/stages

Техкарта: Get/List/Create/Update/Delete, массовые; материалы, параметрические материалы, продукты и этапы — вложенные методы.

Проверки и особенности: cost — копейки; quantity — не деньги; параметрические материалы отдельной формой.

Источники: [processingplan](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_processingplan.md).

### M05: Processing + materials/products

Техоперация: Get/List/Create/Update/Delete, массовые, New, Trash; материалы и продукты.

Проверки и особенности: processingSum — копейки; processingPlan — ссылка; склады материалов/продуктов.

Источники: [processing](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_processing.md).

### M06: ProcessingOrder + positions

Заказ на производство: Get/List/Create/Update/Delete, массовые, New, Trash, позиции.

Проверки и особенности: processings — read-only ссылки.

Источники: [processingorder](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_processingorder.md).

### M07: ProductionTask + materials/products/stages

Производственное задание: Get/List/Create/Update/Delete, массовые, Trash; позиции, продукты, материалы, производственные этапы.

Проверки и особенности: Статусы этапов — открытые; enableHourAccounting и пр. флаги — явный false.

Источники: [productiontask](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_productionTask.md).

### M08: ProductionStageCompletion

Выполнение этапа: Get/List/Create/Update/Delete, массовые, Trash; материалы и продукты.

Проверки и особенности: productionStage обязателен; quantity — не деньги.

Источники: [productionstagecompletion](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_productionStageCompletion.md).

## Общий контракт реализации

- Перед кодом — field/method freeze: поимённая Write matrix, readonly/TODO,
  null/required/Create-only/Update-only; status needs-triage сам не разрешает реализацию.
- Образец пакета и тестов указан в каждом ticket; новый пакет повторяет его структуру
  (model.go, request.go, service.go, *_test.go), без новых общих абстракций.
- Delete, CreateMany/UpdateMany/DeleteMany — только если раздел источника их описывает,
  по образцу [optional X01/X02](../sdk-v1-optional/spec.md).
- Новый ресурс получает handle фасада `moysklad.Client`, строку в
  [status-entities](../../docs/status-entities.md) / [status-technical](../../docs/status-technical.md)
  и проходит `status_test.go`; пункт снимается из [status-backlog](../../docs/status-backlog.md).
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
- Проверки: `make check` (go test -race, go vet ./... ./testdata, gofmt, SHA-256 фикстур,
  статус сущностей, diff --check); независимые Standards/Spec review на фиксированном BASE.
- Коммит по ticket; shared dependency раньше зависимого, публичные docs отдельным
  завершающим commit волны. Planning не коммитить без отдельной команды.
- Live API, credentials, push/tag/release и worktree cleanup не входят.

## Open Decisions

- Материалы/продукты/этапы — собственные вложенные формы, не positions.Position.
- Документы производства используют O01/O02 там, где источник описывает позиции и шаблон, и O04 для Trash.

## Out of Scope

Expand; расчёт себестоимости и планирование — сервер. Все не перечисленные в tickets ресурсы/операции.
Этот DRAFT не является командой реализации или публикации.
