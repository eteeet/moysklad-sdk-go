Status: needs-triage
Stage: DRAFT
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md

# T08: Weighted rate limiter и concurrency

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Раздельные shared principal/account budgets: 45 weighted units/3s; ≤5 principal и ≤20 account
concurrent. Auth profile задаётся явно (Employee/App/Draft), clock/date и weights
конфигурируемы. Текущий Employee weight3, с 2026-12-01 weight4, Draft4, App1; только два точных
stock endpoints weight5.

## Специальные проверки

Fake clock, rolling-window границы, fairness/queue cancellation; acquire каждой попытки, release
после body close, не во время backoff. Один процесс не гарантирует global account limit других
процессов; sharing через общий coordinator описан. Headers rate budgets не интерпретировать
недоказанными weighted units.

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
