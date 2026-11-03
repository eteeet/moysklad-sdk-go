Status: resolved
Stage: APPROVED
Type: task
Blocked by: 01-stock-all.md; 02-stock-bystore.md; 03-stock-current.md; 04-stock-operation-slot.md; 06-dashboard.md; 07-report-counterparty.md; 08-money.md; 09-profit.md; 10-orders-sales.md

# I23: Async создания отчётов

[Spec](../spec.md), [карта](../../sdk-full/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

ListAsync (создание async-задачи) для отчётов из списка _async.md: Остатки, Прибыльность, Деньги, Показатели продаж и заказов, Показатели контрагентов, Показатели (dashboard).

## Специальные проверки

Только endpoints, разрешённые источником; выборочные показатели контрагентов (POST) — без async; результат через существующий пакет async, по образцу ListAsync counterparty/assortment.

## Образец

`async` (Wait/Result) и `transport.Client.Download`.

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

[async](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/_async.md)

## Out of Scope

Expand, live API и не перечисленные What to build операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `3778088`, code commit `1857634`.

- 20 методов *Async в шести пакетах отчётов через write.Submit; Selected без async.
- Review: [../review-wave.md](../review-wave.md); исправления — `aa5c8ea`.
- Проверки: make check на `aa5c8ea` — PASS.
