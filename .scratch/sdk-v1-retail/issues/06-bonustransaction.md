Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md

# E06: BonusTransaction Get/List/Create/Update

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

agent/bonusProgram/transactionType Create; bonusValue, status, parentDocument и ограничения
доступной записи по source. Баллы — собственная числовая величина, не деньги. Permission model
Employee token vs solution token описать отдельно; JSON API наличие POST не гарантирует
разрешение Vendor token.

## Специальные проверки

WAIT_PROCESSING/COMPLETED/unknown, accrual/debit specifics, immutable/C-only fields;
синтетический отказ токена решения передаётся без скрытого fallback на personal token. Live в
Draft — отдельное согласие и не доказано unit tests.

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

[bonustransaction](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bonus_operation.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `0de1969`, code commit `a98771c`; общий тест опций — `410a255`.

| Поле | Read | Create | Update |
| --- | --- | --- | --- |
| agent, bonusProgram | `*Reference` | обяз. | Field |
| transactionType (EARNING/SPENDING) | `string` | обяз. | нет: тип менять нельзя |
| bonusValue | `int` (баллы) | Field | Field |
| name, description, code, externalCode, moment, executionDate, applicable, shared | значения | Field | Field |
| organization, owner, group, parentDocument | `*Reference` | Field | Field |
| transactionStatus, categoryType (readonly), id, meta, accountId, created, updated, updatedBy | | — | — |

- Перечисления открыты; статусы WAIT_PROCESSING, COMPLETED, CANCELED и неизвестный читаются (тест).
- executionDate источник описывает для начислений (отложенная обработка, тарифная опция); для
  списаний и баланс баллов решает сервер.
- Противоречия: bonusProgram без маркера в таблице атрибутов, но обязателен в таблице метода; пример Create
  передаёт readonly categoryType и не передаёт bonusProgram/transactionType; updatedBy назван атрибутом
  только для фильтрации, но фикстура его возвращает.
- Права: источник не описывает разрешения токенов. Синтетический отказ токену решения (WithAuth, 403)
  возвращается как APIError после одной попытки, без перехода на токен клиента (тест). Live в Draft не
  выполнялся.
- Review Standards: замечаний нет. Review Spec: неподтверждённое утверждение о токене сотрудника убрано,
  противоречия записаны, golden Update без недокументированного null.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
