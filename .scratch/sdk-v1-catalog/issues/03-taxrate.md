Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md

# C03: TaxRate Get/List/Create/Update

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Полный TaxRate contract и readonly/required полей из source. Не смешивать ставку TaxRate с
денежными суммами или Decimal; server business validation.

## Специальные проверки

Нулевые ставки, optional dates/refs если есть, sparse payloads, HTTP/query и error paths.

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

[taxrate](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_taxrate.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `d9aced4`, code commit `bc01138`.

| Поле | Create | Update | Примечание |
| --- | --- | --- | --- |
| rate | обяз. (`*float64`, 0 допустим) | + | процент, не деньги; у предустановленных ставок менять не даёт сервер (L526) |
| comment, archived | + | + | |
| shared, owner, group | + | + | только у пользовательских ставок |
| id, meta, accountId, updated | — | — | readonly |

- Противоречие: таблица помечает updated «Необходимо при создании», текст L234 — «единственное поле … rate»;
  SDK следует тексту (в contracts.md пометка осталась как в источнике).
- Shared у предустановленной ставки отсутствует и читается как false — как у других моделей SDK;
  отличить можно по Owner == nil.
- search (по rate и comment), filter и order передаются общими ListOptions; listtest проверяет кодирование.
- Тесты: фикстура и разреженная системная ставка обоими codecs, golden full/zero/update, null, защита
  readonly, HTTP-контракт documenttest + LocalErrors, List, валидация без запросов.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review: обе оси одним агентом с раздельными разделами (малый ticket); дефектов нет.
- Регистрация в общих тестах опций и сброса ResponseInfo (`entity/options_test.go`) — `cc61f6a` (найдено в review C07).
