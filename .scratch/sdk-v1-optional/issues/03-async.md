Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/03-error-codes.md; ../../sdk-v1-foundation/issues/08-rate-concurrency.md; ../../sdk-v1-foundation/issues/09-safe-retry.md; ../../sdk-v1-files-reports/issues/01-download-transport.md

# X03: Async submit/status/result/cancel

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Только documented async-enabled endpoints; Submit даже для GET является side effect. Own task
model, statuses и poll/result/cancel endpoint. ≤4 queued account tasks — серверный лимит аккаунта:
SDK его документирует и возвращает ошибку 61002, локальной policy очереди нет (решение пользователя
2026-11-02, review-0fafaaa P1). async=true нельзя сочетать limit/offset.

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

## Answer

BASE `6a3bbbc`, коммит `791490f`. Решение пользователя 2026-11-02: ListAsync у counterparty и
assortment (единственные async-endpoints источника, которые есть в SDK) и сервис задач `c.Async()`:
Get, List, Wait, Result(w), Cancel. Async-отчёты в SDK нет — вне scope.

- Пакет `async`: State (открытый набор, Finished), Task (errors — массивом, как в примере; таблица
  называет Object), Submitted (ID, StatusHref, ResultHref), Service. `internal/write.Submit`: limit/
  offset — ошибка до запроса (61001), GET с async=true, только 202, Content-Location …/async/{id} и
  Location …/async/{id}/result, иначе ошибка. transport не повторяет async=true даже у GET.
- Wait — опрос Get с первым запросом сразу, срок задаёт только ctx (оговорено в godoc). Result —
  `transport.Client.Download` (правила F01: redirect без Authorization). Cancel — PUT, только 204
  (источник), «не откат». Лимит 4 задач — серверный на аккаунт, Client его не считает; после сбоя
  задача могла создаться — проверка через List по request. Assortment async идёт с beta-тегом, как
  List; источник ссылается на новый endpoint, но отдельно это не подтверждает.
- Фикстуры `testdata/custom/async`: три примера `_async.md` байт в байт с SHA-256 (README);
  `task_done.json` с пробелами в конце строк — исключение в .gitattributes.
- Фасад `c.Async()`, строка async в status-entities («прочее» 5/8), status_test.
- Тесты `async/*_test.go`: фикстуры обоими кодеками (deletionDate с миллисекундами), Finished, Get/List,
  Wait (до API_ERROR, дедлайн ctx, ошибка Get), Result (redirect другого origin без Authorization,
  APIError 1016 и 61006, w не тронут), Cancel (204; 200/202/503 — ошибка, один вызов при Retry),
  ошибки до запроса; ListAsync: query, beta-тег, 202 и разбор адресов, 200/503 без повтора, ответы без
  адресов или с чужим Location, limit/offset/fields/async в Extra.
- Review Standards: проверка limit/offset перенесена в Submit, beta-опции — `withBeta`. Оставлено:
  Submit в `internal/write` (публичный пакет async не должен давать создавать задачи у любого
  endpoint). Review Spec: оговорки неизвестного исхода и бесконечного Wait, Cancel только 204,
  проверка Location, тесты 61006 и миллисекунд, коды 61000–61008 в godoc.
- Не проверено: живой API, async ассортимента с beta-тегом, вебхуки задач (не входят).
- Пользовательское review 0fafaaa: P1 — лимит очереди решением пользователя серверный, spec и
  ticket приведены к нему; P2 — `Result` требует итоговый HTTP 200 (`6b1bc33`), 202/204 — ошибка.
