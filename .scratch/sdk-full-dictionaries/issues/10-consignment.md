Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-full-catalog/issues/01-barcodes-packs.md

# S10: Consignment Get/List/Create/Update

[Spec](../spec.md), [карта](../../sdk-full/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Партия: Get/List/Create/Update/Delete, массовые, attributes.

## Специальные проверки

assortment — обязательная ссылка; label; barcodes — существующая model.Barcode, порядок после G01 только ради одного владельца формы штрихкодов.

## Образец

`entity/uom` (service_test.go: TestServiceGet/Create/Update) + batch/delete как в `entity/country`.

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

[consignment](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_consignment.md)

## Out of Scope

Expand, live API и не перечисленные What to build операции.
Этот DRAFT не является командой реализации или публикации.

## Comments

Не начат: заблокирован G01 (sdk-full-catalog), который ещё не resolved.

## Answer

BASE `f2f5f09`, code commit `d89f11a`; docs `670f379`.

- Три версии обязательных полей в источнике (таблица, метод, вводный текст); SDK требует assortment и хотя бы одно из label и expiryDate (review 2a76d13, `7de204f`). barcodes — model.Barcode.
- Review: [sdk-full-catalog/review-catalog-wave.md](../../sdk-full-catalog/review-catalog-wave.md).
- Исправления по review — `2a76d13`.
- Проверки: make check — PASS.
