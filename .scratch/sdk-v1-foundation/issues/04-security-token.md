Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/01-public-interface.md; ../../sdk-v1-foundation/issues/03-error-codes.md

# T04: POST security/token

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Отдельный security service: CreateToken(ctx, options), POST security/token с Basic Auth, без
тела; HTTP 200 и непустой access_token. В документации явно: выпуск отзывает предыдущие токены
сотрудника. Не вызывать из NewClient и не переключать Auth автоматически.

## Специальные проверки

Exact method/path/headers/empty body; malformed/empty token, auth/API/context errors; одна
отправка даже при 429/5xx/потере ответа. Токен не попадает в diagnostics/fixtures/logs. Live
выпуск только по отдельной команде.

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

BASE `2bf89ca`, code commit `53ffa95`.

- `security.Service.CreateToken(ctx, username, password, opts...)` — сознательно явные Basic
  credentials per call вместо `CreateToken(ctx, options)` (research). POST security/token без тела,
  только 200 + непустой access_token; Basic ставится последним `transport.WithAuth`, поэтому
  работает на клиенте с Bearer и не подменяется caller. Auth клиента не меняется, NewClient не вызывает.
- `Token` скрывает секрет в fmt (`%v/%+v/%#v`) и slog (LogValuer); json.Marshal раскрывает — описано.
- transport: `WithAuth` (ADR 0001); 307/308 для не-GET/HEAD не повторяются (ErrUseLastResponse);
  `Idempotency-Key`/`X-Idempotency-Key` защищены — иначе net/http повторил бы запись.
- Тесты: plain/gzip, точные method/path/headers/пустое тело, 401/429/5xx/201/empty/missing/malformed/null
  — по одной отправке, потерянный ответ (hijack) и отмена во время запроса без повтора, redirect.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Для T08/T09: WithAuth меняет principal вызова — limiter не должен списывать его с бюджета
  Bearer клиента; security/token нужен явный признак «никогда не повторять», а не только «POST».
