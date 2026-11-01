Status: needs-triage
Stage: DRAFT
Type: task
Blocked by: ../../sdk-v1-foundation/issues/06-query-pages.md

# T07: Явный walker/iterator страниц

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Явный iterator по List/GetPositions/report pages, cancellation/early stop, bounded память;
обычные методы по-прежнему делают один запрос. Проверять nextHref same-origin и safe path либо
вычислять offsets по документированному контракту, не следовать слепо внешним ссылкам.

## Специальные проверки

Empty page, no total/size where optional, wrong/missing rows, repeated nextHref, no-progress,
unusual offsets, cancellation до следующей страницы. Не обещать snapshot-consistency или
exactly-once при меняющемся наборе. Array-only PriceType не пагинировать вымышленным limit.

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
