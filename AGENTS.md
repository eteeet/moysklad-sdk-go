# moysklad-sdk-go

Общие правила — `~/ETEEET/AGENTS.md` и `~/ETEEET/moysklad/AGENTS.md`: прочитай их первыми (Claude загрузил их сам).

SDK для JSON API МойСклад на Go.

## Официальные источники МойСклад
- GitHub: https://github.com/orgs/moysklad/repositories — следи за изменениями их SDK: java (новый), python, ts.
- `openspec` и документация API могут отставать от реального API.

## Тесты
- Контракт теста — официальная документация API; `openspec` и upstream-фикстуры могут от неё отставать.
- Upstream-фикстуры используются без изменений; версия и происхождение фиксируются через `make fixtures` (`testdata/readme.md`).
- Тест проверяет значения значимых полей, а не только то, что JSON декодировался.
- Тело запроса сверяется с независимо заданным golden-файлом; `false`, `0`, пустые значения и `null` закрепляются там, где API их допускает.
- HTTP-тест проверяет метод, путь, тело, декодирование ответа и ошибки; совместимость JSON-кодеков — по образцу существующих тестов.
- Тикет на новый сервис называет образец структуры HTTP-тестов. Для простых Get/Create/Update — `entity/uom/service_test.go`: отдельные `TestServiceGet`, `TestServiceCreate`, `TestServiceUpdate` с табличными сценариями внутри.

## Agent skills

### Issue tracker
Локальные Markdown-файлы в `.scratch/`. См. `docs/agents/issue-tracker.md`.

### Triage labels
Стандартные пять статусов. См. `docs/agents/triage-labels.md`.

### Domain docs
Общий и локальный glossary, локальные ADR. См. `docs/agents/domain.md`.
