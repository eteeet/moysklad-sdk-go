Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md

# A03: GET context/employee

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Отдельный context/employee.Get(ctx) без ID; собственная форма контекста и permissions.
Переиспользовать Employee/Permission types только при совпадении JSON contract, не путать с
ResponseContext в коллекциях.

## Специальные проверки

Basic/Employee/App responses/permission errors, sparse fields и неизвестные permissions, нет
обязательного caller id и дополнительных запросов.

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

BASE `8f5942d`, code commit `842adec`.

- `contexts.NewService(client).Employee(ctx, opts...)`: один GET context/employee без id, query и тела.
  Пакет `contexts`, а не `context/employee.Get` из ticket: имя `context` занято стандартной библиотекой.
- `EmployeeContext` встраивает `employee.Employee` (JSON-контракт совпадает; Salary в таблице контекста нет,
  остаётся nil) и добавляет `code` и `permissions`. Не путается с model.ResponseContext.
- Права — общий `model.Permissions` (введён здесь, переиспользуется A02): значение — флаг, области по
  действиям или, если форма другая, Raw без потерь; неизвестные права, действия (done) и области открыты.
  Противоречия источника записаны в godoc: Array(Object) в таблицах vs объект в примерах, VIEW vs view,
  customerOrder vs customerorder; ключи хранятся как пришли.
- Критерий успеха контекста — HTTP 200 и id, который таблица помечает «Обязательное при ответе»;
  не общий ID-gate сущностей.
- Тесты: пример источника обоими codecs (73 права, task.done), открытые множества и нестроковое значение,
  HTTP через documenttest, Basic и Bearer, 401/403 как APIError, локальные ошибки. Фикстура из
  _general.md с SHA-256 в README.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review Spec: lossless права через model.Permissions, константа AUTHOR, 401/403 и Basic, противоречия.
  Review Standards: критерий id обоснован источником, SHA-256 provenance, общий get[T], тесты.
