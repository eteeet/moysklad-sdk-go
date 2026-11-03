Status: resolved
Stage: APPROVED
Type: task
Blocked by: нет

# O02: Шаблон документа PUT entity/{type}/new

[Spec](../spec.md), [карта](../../sdk-full/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

New(ctx, base) у существующих документов с шаблоном: cashin, cashout, demand, invoicein, invoiceout, loss, enter, move, paymentin, paymentout, purchaseorder, purchasereturn, retaildemand, retailsalesreturn, salesreturn, supply. Новые документы добавляют New сами.

## Специальные проверки

Пустое тело и тело с основанием; ответ без id; ничего не сохраняется. enter — по общей таблице шаблонов (основание inventory, первые 500 позиций), отдельного раздела в _enter.md нет.

## Образец

`entity/customerorder` (документ с позициями, GetPositions, batch) и `internal/documenttest`.

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

[common](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_common_info.md), [demand](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_demand.md)

## Out of Scope

Expand, live API и не перечисленные What to build операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `d21d84d`, code commit `692f322`.

- New у 16 документов через write.Template; godoc-примеры оснований исправлены по review (`3b62b18`).
- Review: [review-wave.md](../review-wave.md); исправления — `2e9e03d`, `3b62b18`, `341a4fd`.
- Проверки: make check CHECK_P=1 на `341a4fd` — PASS.
