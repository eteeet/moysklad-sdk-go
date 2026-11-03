Status: resolved
Stage: APPROVED
Type: task
Blocked by: нет

# G01: Product: штрихкоды и упаковки

[Spec](../spec.md), [карта](../../sdk-full/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Product.barcodes через существующую model.Barcode; Product.packs — собственная форма (uom, quantity, barcodes) в read/write. Variant packs (parentpack) уже есть — не трогать.

## Специальные проверки

Тип штрихкода — ключ объекта (ean13, gtin, …), неизвестный тип сохраняется; pack.uom — ссылка; packs у Service/Bundle не добавлять — в источнике их нет.

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

[product](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_product.md), [barcode](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/workbook/_workbook_barcode.md)

## Out of Scope

Expand, live API и не перечисленные What to build операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `f2f5f09`, code commit `8b0af26`; docs `670f379`.

- Barcodes — model.Barcode; Pack{ID, UOM, Quantity, Barcodes}; packs при Update заменяют коллекцию; PackRequest.ID адресует существующую упаковку, без ID — новая (исправлено по review 2a76d13, `7de204f`).
- Review: [sdk-full-catalog/review-catalog-wave.md](../../sdk-full-catalog/review-catalog-wave.md).
- Исправления по review — `2a76d13`.
- Проверки: make check — PASS.
