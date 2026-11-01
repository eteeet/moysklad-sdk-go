Status: needs-triage
Stage: DRAFT
Type: task
Blocked by: ../../sdk-v1-release/issues/02-examples.md; ../../sdk-v1-release/issues/03-ci-release.md

# R04: Приёмка и статус v1

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Проверить coverage-manifest каждого исходного требования пользователя, закрыть или явно вынести
optional gates, обновить README/status-entities/status-technical/coverage без фиктивных checked
для неподдержанных methods. v1 — выбранный профиль SDK, не весь JSON API.

## Специальные проверки

По тикету коммит, dependencies shared-types перед dependent, docs отдельным последним commit
волны; planning не коммитится без команды. Независимые Standards/Spec review каждой волны;
evidence links и unresolved limitations. Cleanup core worktrees — существующий отдельный ticket,
не blocker SDK.

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
