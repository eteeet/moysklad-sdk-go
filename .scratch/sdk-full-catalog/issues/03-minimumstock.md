Status: resolved
Stage: APPROVED
Type: task
Blocked by: нет

# G03: Неснижаемый остаток, storebalances и fields=minimumStock

[Spec](../spec.md), [карта](../../sdk-full/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Модель minimumStock; снять отказ fields=minimumStock в List Product и Variant; storebalances Product и Variant: List/Get/Create/Update/Delete/DeleteMany.

## Специальные проверки

Формы all/store, значения quantity; fields передаётся одним значением.

## Образец

`entity/bundle` (GetComponents/GetComponent) и `entity/variant` (характеристики).

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

[product](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_product.md), [variant](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_variant.md)

## Out of Scope

Expand, live API и не перечисленные What to build операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `f2f5f09`, code commit `512b519`; docs `670f379`.

- model.MinimumStock читает storebalances и storeBalances (источник пишет оба); fields=minimumStock в List Product/Variant; storebalances — только маршруты в источнике, поля по таблице; Create принимает объект и массив из одного.
- Review: [sdk-full-catalog/review-catalog-wave.md](../../sdk-full-catalog/review-catalog-wave.md).
- Исправления по review — `2a76d13`.
- Проверки: make check — PASS.
