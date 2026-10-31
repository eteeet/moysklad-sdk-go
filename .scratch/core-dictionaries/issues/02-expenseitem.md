Status: resolved
Blocked by: None

# 02: ExpenseItem — Get, Create, Update

## What to build

Разработчик решения получает, создаёт и частично изменяет ExpenseItem с полным согласованным
набором: 13 полей чтения и 9 записи. Срез включает модели, три метода, тесты и документацию.

Контракт: [общая spec](../spec.md), разделы ExpenseItem, общих правил и Testing Decisions.
Образцы — Product/UOM; Country при наличии полезен как пример, но не является зависимостью.

## Acceptance criteria

- [x] Пакет expenseitem предоставляет ExpenseItem, Service, NewService с контрактом Product/UOM:
  существующий *transport.Client, возврат *Service и error, отказ для nil client.
- [x] Модель читает meta (model.Meta), id/accountId/name/description/code/externalCode (string),
  updated (model.DateTime), owner/group (*model.Reference), shared/archived/operatingExpenses (bool).
  Все 13 полей из таблицы атрибутов поддержаны; code сохраняется строкой.
- [x] CreateRequest содержит обязательный Name string; description, code, externalCode,
  owner, group, shared, archived, operatingExpenses используют Field с omitzero.
  В UpdateRequest все девять полей используют Field. Readonly-полей нет.
- [x] Get — GET entity/expenseitem/{id}, Create — POST entity/expenseitem,
  Update — PUT entity/expenseitem/{id}; без лишних query, тела GET и дополнительных запросов.
  Возвращаются (*ExpenseItem, error) и фактические данные сервера.
- [x] Defaults archived=false и operatingExpenses=true оставлены серверу; SDK не добавляет
  их в минимальный Create. Отсутствующее operatingExpenses при чтении даёт нулевой bool,
  а явный operatingExpenses=false отправляется без потери.
- [x] Предустановленные примеры без accountId/owner/group/shared декодируются.
  Запрет Update предустановленных UOM/Country не переносится на ExpenseItem:
  в проверенной документации ExpenseItem такого запрета нет, происхождение локально не определяется.
- [x] Проверены значения всех поддержанных полей upstream, оба JSON-кодека и неизвестные поля.
  Upstream не меняется: отсутствующие там operatingExpenses/owner/group/shared покрыты небольшими
  отдельными случаями; официальный пример «Возврат» или аналогичный хранится с источником/revision.
- [x] Проверены operatingExpenses=true/false/пропуск, archived=true/false, shared=false,
  ссылки, даты и самостоятельный строковый code с ведущими нулями.
- [x] Requests проверены обоими JSON-кодеками: минимальный Create, полные Create/Update
  по независимым golden, одно поле Update, пропуски, false/true, "" и null.
  Поддержка кодирования null не является обещанием серверной очистки.
- [x] Проверка UUID не ограничивает версию; Create отклоняет только Name=="" без trim;
  незаданный Update отклоняется, явные false/""/null считаются заданными.
- [x] HTTP-тесты используют настоящий transport и проверяют метод/маршрут/заголовки,
  plain/gzip, отсутствие лишних запросов, минимальный/полный request и серверные данные ответа.
- [x] Для трёх методов покрыты context, auth/transport/APIError 4xx/5xx/429,
  errors.As/errors.Is, контекст без двойного moysklad:, неожиданный статус,
  malformed/null JSON, отсутствующий/пустой ID; успех — HTTP 200 и непустой ID, ошибки возвращают nil.
  Предварительного Get, merge или retry нет.
- [x] Документация отражает готовый ExpenseItem и TODO. Копии do/validID допустимы;
  обнаруженные общие исправления рассматриваются во всех применимых копиях по правилу coverage.
- [x] Пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check и встроенное ревью Standards/Spec с фиксированным исходным BASE.

## Blocked by

Нет. Не требует завершения Country или других tickets этой группы.

## Выполнение и границы

- Один основной агент, одна новая сессия, implement/TDD на весь ticket без учебных остановок.
- Прочитать инструкции и всю spec; записать исходный HEAD для ревью, проверить dirty diff
  и сохранить все ранее существовавшие изменения планирования и чужие изменения.
- Не включать исходные изменения планирования в commit реализации; свои изменения документации
  отделять от них. Reviewer только сообщает замечания; исправляет основной агент.
- Не добавлять List/Delete/expand/batch, общие транспортные механизмы, фасад, живой API или .env.
- Новые контракты требуют согласования. После полного ticket — локальный commit,
  Status: resolved и ## Answer с SHA, проверками и результатами ревью; отчёт и остановка.
  Push — только по отдельной команде.

## Sources

[Официальная документация ExpenseItem, зафиксированный revision](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_expenseitem.md).
Upstream уже имеется и не редактируется; provenance закреплён механизмом SDK.

## Answer

Реализованы ExpenseItem, отдельные CreateRequest/UpdateRequest и Get/Create/Update,
все 13 полей чтения и 9 записи, оба JSON-кодека, HTTP-тесты, независимые golden
и официальный пример «Возврат» с provenance. Документация: `docs/expenseitem.md`.

- Commit реализации: `95f1b4ad0bd87e599d9391236a3ae61ccd1c36f1`; включён в `dev` как `09f7dc8`.
- Проверки: gofmt; `go test -race -count=1 ./...`; `go vet ./... ./testdata`;
  `git diff --cached --check` — пройдены.
- Независимые ревью: Standards — 0 замечаний; Spec — 0 замечаний.
- Живой API не использовался; исходные изменения планирования не включены в commit.
