Status: resolved
Stage: APPROVED
Type: task
Blocked by: нет

# O01: Отдельные операции позиций

[Spec](../spec.md), [карта](../../sdk-full/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Generic GetPosition/CreatePositions/UpdatePosition/DeletePosition/DeletePositions в `positions`; typed-методы во всех 13 документах с позициями.

## Специальные проверки

≤1000 позиций в Create; ID документа и позиции проверяются до HTTP; POST массивом возвращает массив; никаких повторов записи.

## Образец

`positions.Service[T].GetPage` и `entity/customerorder` GetPositions.

## Acceptance criteria

- [ ] Scope/публичный interface и field/method matrix утверждены; зависимости resolved
  и их код присутствует в рабочем дереве. Readonly, required response/Create, null,
  immutable и выбранные TODO разделены по источникам; противоречия не замолчаны.
- [ ] Реализован весь What to build, без неоговорённых методов и hidden requests.
- [ ] Специальные проверки этого ticket реализованы через публичный interface.
- [ ] JSON обоими codecs и независимые golden; upstream/provenance не изменены.
  Для transport/report/print tickets вместо entity JSON применяются их собственные
  публичные seams и контрактные tests, а не фиктивные модели.
- [ ] Реальные transport/httptest tests покрывают успех и ошибки; context и
  metadata сохраняются, caller-controlled retry/walker не подменяет один вызов.
- [ ] go test -race, go vet ./... ./testdata, gofmt и diff-check PASS; независимые
  Standards/Spec review, исправления и финальный BASE/результат записаны в Answer.
- [ ] Фасад, status-entities/status-technical, status-backlog и guide отражают реальный
  scope, ограничения и неизвестные факты. Planning без команды не коммитится; один code
  commit по ticket.

## Источники

[customerorder](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_customerOrder.md), [common](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_common_info.md)

## Out of Scope

Expand, live API и не перечисленные What to build операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `d21d84d`, code commit `5d258ce`.

- Операции строк в positions.Service и typed-методы 13 документов; DeletePositions переименован в DeleteManyPositions и вынесен в nested.DeleteRows (`2e9e03d`).
- Review: [review-wave.md](../review-wave.md); исправления — `2e9e03d`, `3b62b18`, `341a4fd`.
- Проверки: make check CHECK_P=1 на `341a4fd` — PASS.
