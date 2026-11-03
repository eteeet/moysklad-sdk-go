Status: needs-info
Stage: DRAFT
Type: task
Blocked by: нет

# O07: agentAccount/organizationAccount у cashin/cashout

[Spec](../spec.md), [карта](../../sdk-full/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Решить расхождение: текст Update требует счета, таблица их не содержит.

## Специальные проверки

needs-info: нужен ответ МойСклад или live-проверка с согласия.

## Образец

`entity/paymentin` / `entity/cashin` (FinanceOperation, без позиций).

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

[cashin](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_cashin.md), [cashout](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_cashout.md)

## Out of Scope

Expand, live API и не перечисленные What to build операции.
Этот DRAFT не является командой реализации или публикации.

## Notes 2026-11-03

- Дополнительный источник: официальная фикстура testdata/upstream/cash_in.json (и cash_out.json) не
  содержит agentAccount/organizationAccount — совпадает с таблицей атрибутов; текст «Изменить» о счетах,
  вероятно, перенесён из paymentin/paymentout.
- SDK уже следует таблице (счета не читаются и не пишутся), противоречие записано в docs/cashin.md и
  docs/cashout.md. Остаётся needs-info: закрыть может только ответ МойСклад или live-проверка с согласия.
