Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md

# C01: Currency Get/List/Create/Update

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Полный Currency contract, name/code/isoCode Create, rate и system/custom currency различия. Курс
не денежная сумма.

## Специальные проверки

Required/read-only и ограничения rate, нет автоподстановки курса; user/system fixtures и
APIError.

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

[currency](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_currency.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `8eda31b`, code commit `1e29d0d`. Источник _currency.md (L11–31, L167–330, L558–563).

| Поле | Create | Update | Примечание |
| --- | --- | --- | --- |
| name, code, isoCode | обяз. (пользовательская) | непустые, если переданы | системная: code или isoCode |
| system | + (системная форма) | — | readonly в таблице, но раздел системной валюты требует — исключение |
| rateUpdateType | + (системная форма) | — | то же исключение; auto или manual |
| margin | + | + | текст упоминает его у системной формы; без пометки |
| rate | ≠ 0 | ≠ 0 | курс, не деньги; без автоподстановки |
| fullName, archived, multiplicity, indirect, majorUnit, minorUnit | + | + | без пометки, сервер проверяет системные ограничения |
| id, meta, default | — | — | readonly |

- Ограничения системной валюты и валюты учёта (L562–563) проверяет сервер; pre-GET нет. null для полей
  источник не описывает: Field позволяет, решает сервер (кроме rate/name/code/isoCode — отклоняются).
- Тесты: upstream RUB (system, default), примеры источника пользовательской USD и системной с autokурсом
  (README с SHA-256), sparse/unknown, golden custom/system/update, null в запросе, защита readonly-полей,
  HTTP-контракт documenttest, List, валидация обеих форм без запросов.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review Spec: единое исключение для system/rateUpdateType, фикстуры, матрица, margin в custom golden убран.
  Review Standards: sparse/null/readonly тесты, имя checkValues.
- Регистрация в общих тестах опций и сброса ResponseInfo (`entity/options_test.go`) — `cc61f6a` (найдено в review C07).
