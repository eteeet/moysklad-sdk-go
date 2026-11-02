Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-release/issues/01-facade.md

# R02: Примеры основных сценариев

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Компилируемые примеры: поиск Assortment, CustomerOrder→страница positions, Audit contexts→явные
events, Turnover, Files и client initialization для personal/App token. Запись защищена явным
opt-in; token issuance не часть обычного getting-started.

## Специальные проверки

Examples unit tests с httptest и компиляция CI, no real credentials/live writes;
pagination/close/cancellation показаны корректно. Нет автоматического выполнения examples
Makefile в тестовом прогоне.

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

## Answer

BASE `6706923` (R01). Примеры — программы `examples/<сценарий>` на фасаде; `main` читает окружение,
сценарий — функция с `*moysklad.Client` и `io.Writer`, её тест идёт через настоящий transport и
httptest. Реальных токенов и запросов к API в тестах нет; Makefile примеры в тестах не запускает.

- `examples/internal/env`: `Client()` по `MOYSKLAD_TOKEN` (персональный токен сотрудника и токен
  решения — одинаково Bearer; выпуск токена security/token в примеры не входит), `MOYSKLAD_BASE_URL`,
  `Context` с отменой по Ctrl+C и таймаутом.
- `getting_started`: context/employee, с `MOYSKLAD_APP=1` — context/application; SDK вид токена
  не определяет.
- `assortment_search`: List с Search и явный обход `pages.All`. Находка: следующим страницам нужна
  та же `transport.WithBetaFeature(assortment.BetaFeature)`, иначе GET по nextHref уходит в legacy
  endpoint. Пример передаёт опцию, godoc `assortment.List` дополнен; тест проверяет тег на второй
  странице и остановку обхода при отмене ctx.
- `customerorder_positions`: Get заказа и страницы GetPositions через pages.All.
- `audit_events`: страница контекстов и по одному явному GET событий на контекст.
- `turnover`: All за явный период (противоречие источника о датах).
- `files`: List, Download первого файла с закрытием на всех путях и удалением частичного файла,
  Upload только с `MOYSKLAD_ALLOW_WRITE=1`.

Существующие `product_get/create/update` и цели Makefile не менялись.

Тесты: успех и ошибка API (401 со статусом через errors) каждого сценария, страницы по nextHref,
отмена ctx до следующей страницы, opt-in записи до создания клиента, сохранение существующего файла
и удаление частичного.

Review Standards: исправлены комментарий find, godoc Require, перезапись существующего файла
(`O_EXCL`), синхронизация в тестах. Switch имени строки ассортимента оставлен в примере: метода
у Row нет, общий helper — отдельное решение. Ловушка beta-тега в pages.All записана как находка
для отдельного решения пользователя. Review Spec: turnover и события аудита обходят все страницы
через pages.All, «первые 10 контекстов в порядке сервера» вместо «последних», тесты ошибок.
README и статусы — в docs-коммите волны. Проверки: go test -race ./..., go vet ./... ./testdata,
gofmt, diff --check — PASS.

## Источники

[транспортный research](../../sdk-v1/research/transport.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.
