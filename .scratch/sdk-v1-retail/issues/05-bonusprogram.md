Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md

# E05: BonusProgram Get/List/Create/Update

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Штатная бонусная программа; required
name/active/allProducts/allAgents/earnRateRoublesToPoint/spendRatePointsToRouble/maxPaidRatePercents.
Диапазоны/зависимости и активная программа по source, не рассчитывать баланс/правила маркетинга
в SDK.

## Специальные проверки

Явные false и 0 обязательных полей не пропадают, arrays/refs/null/readonly, rate restrictions —
сервер. Не считать баллы копейками только из-за слова bonus.

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

[bonusprogram](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bonus_program.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `2114fad`, code commit `0de1969`.

| Поле | Read | Create | Update |
| --- | --- | --- | --- |
| name, active, allProducts, allAgents, earnRateRoublesToPoint, spendRatePointsToRouble, maxPaidRatePercents | значения | обяз., обычные значения (false и 0 уходят) | Field |
| agentTags | `[]string` | Field | Field |
| earnWhileRedeeming, welcomeBonusesEnabled, welcomeBonusesMode, welcomeBonusesValue, postponedBonusesDelayDays | значения | Field | Field |
| id, meta, accountId | readonly | — | — |

- Курсы, проценты, баллы и дни — `int`, не деньги; welcomeBonusesMode открыт.
- Сервер проверяет диапазоны, приветственные баллы (не пусто при welcomeBonusesEnabled, не отрицательно)
  и тарифную опцию отложенного начисления; тест отправляет 150% и −1 как есть. allProducts источник
  описывает «всегда true» — SDK передаёт значение вызывающего кода.
- Локально: непустое name (метод называет его обязательным), UUID и непустой Update.
- Review Standards: замечаний нет. Review Spec: чтение отсутствующих полей как 0/"" отмечено в godoc.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
