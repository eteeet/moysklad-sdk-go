Status: needs-triage
Stage: DRAFT
Type: task
Blocked by: ../../sdk-v1-foundation/issues/03-error-codes.md; ../../sdk-v1-foundation/issues/08-rate-concurrency.md

# T09: Ограниченные retries/backoff

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Operation-aware retry policy: только подтверждённые reads GET/HEAD по умолчанию;
POST/PUT/security-token/async-submit — без автоматического повтора. 429/временные 5xx/transport
failures с ограничением attempts, context, jitter и корректными server hints. Окончательные
metadata и APIError сохраняются.

## Специальные проверки

X-Lognex durations в миллисекундах, стандартный Retry-After — seconds/date; не обрезать server
minimum своим cap. Каждая попытка новый replay body/headers и limiter admission; close
промежуточных body до ожидания. Nil/unreplayable stream, decode error после 2xx и canceled
context не повторяются.

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
