Status: needs-info
Stage: DRAFT
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/03-error-codes.md; ../../sdk-v1-foundation/issues/08-rate-concurrency.md; ../../sdk-v1-foundation/issues/09-safe-retry.md; ../../sdk-v1-files-reports/issues/01-download-transport.md

# X03: Async submit/status/result/cancel

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Только documented async-enabled endpoints; Submit даже для GET является side effect. Own task
model, statuses и poll/result/cancel endpoint. ≤4 queued account tasks policy отдельно от HTTP
semaphore. async=true нельзя сочетать limit/offset.

## Специальные проверки

Submit не retry; poll bounded/context; result download origin/auth rules F01 при ссылочном
result; 202/Location/response shape из _async.md, cancel нельзя выдавать за rollback.
Distributed account queue полностью одним Client не контролируется.

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
