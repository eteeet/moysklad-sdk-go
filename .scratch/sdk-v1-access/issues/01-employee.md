Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md

# A01: Employee: Get/List/Create/Update

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Полный основной Employee contract; required lastName; отдельные request types, references и
typed nested salary/cashier where documented. Фото и credential lifecycle не добавляются этим
тикетом.

## Специальные проверки

Поля по contracts.md, sparse/null, false/zero, exact required Create, list/query/headers и
errors; Employee rights не угадываются по имени роли.

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

[employee](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `db22bf1`, code commit `8f5942d`.

Write matrix (источник _employee.md, L14–44, Create/Update):

| Поле | C | U | Тип | Примечание |
| --- | --- | --- | --- | --- |
| lastName | required | required | string | «В теле запроса обязательно … lastName» и при изменении |
| firstName, middleName, description, externalCode, email, phone, inn, position | + | + | Field[string] | email — только админ |
| shared | + | + | Field[bool] | |
| archived | + | + | Field[bool] | только админ |
| salary | + | + | Field[Salary{Value float64}] | копейки; только админ |
| owner, group | + | + | Field[Reference] | только админ, без expand |
| attributes | + | + | Field[[]Attribute] | |
| id, accountId, meta, name, fullName, shortFio, uid, created, updated, cashiers | — | — | read | readonly |
| image | — | — | read Image | загрузка фото вне scope |

- Права «только администратор» проверяет сервер. null для полей источник не описывает: Field позволяет
  его отправить, но golden его не используют — непроверенный контракт, решает сервер.
- Роль и права сотрудника в модели отсутствуют и не выводятся по имени роли — отдельный ресурс security.
- `employee.NewService`: Get/List/Create/Update; List через list.Entity (archived, без fields).
- Тесты: fixture обоими codecs (оклад, кассиры, фото, доп. поля), sparse, golden Create min/full и
  Update full, HTTP Get/Create/Update/List, локальные ошибки, валидация без запросов; добавлен в
  entity/options_test (опции, сброс ResponseInfo).
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review Spec: Update требует lastName (исправлено), null убран из golden. Review Standards: без ID ticket
  в комментарии, порядок пар в options_test.
