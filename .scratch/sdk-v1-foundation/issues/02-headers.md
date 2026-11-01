Status: needs-triage
Stage: DRAFT
Type: task
Blocked by: ../../sdk-v1-foundation/issues/01-public-interface.md

# T02: Все request/response headers и интеграция options

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Весь инвентарь headers из research/transport.md: typed request options, repeated values, clone
raw response headers, parsed optional values с сохранением malformed/raw. Протянуть per-call
options и metadata во все 17 существующих services и positions. Управление auth остаётся явным;
новая мета не теряется на APIError.

## Специальные проверки

Repeated WebHook-DisableByPrefix без CSV; beta feature одно значение; удалённый WebHook-Disable
не рекламируется как рабочий. Region/Accept-Language, report date with numeric timezone,
content/deprecation/auth/download headers. Defaults/overrides,
malformed/unknown/negative/overflow значения, cancellation и race. DateTime не перенастраивается
глобально по report header.

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
