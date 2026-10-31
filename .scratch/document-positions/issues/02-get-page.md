Status: ready-for-agent
Stage: APPROVED
Blocked by: None

# 02: Generic reader одной страницы позиций

## What to build

Разработчик решения читает одну страницу позиций любого document type с поддерживаемым
API endpoint через общий positions.Service[T any].GetPage и настоящий transport.Client;
получает model.Collection[T] с metadata и строками, управляет следующей страницей сам.
Срез закончен от public API до локального HTTP-теста и обработки ошибок.

Interface и HTTP seam P02 утверждены пользователем: [spec](../spec.md).
Ticket ready-for-agent; entity endpoints в эту волну не входят.

## Acceptance criteria

- [ ] `NewService[T](client *transport.Client, documentType string)` возвращает
  (*Service[T], error), nil client отклоняется; T any не требует Core, assortment или ID.
- [ ] documentType — открытая строка одного безопасного API route segment; согласованный
  syntax `[a-z][a-z0-9]*`, без allowlist первых документов и скрытой нормализации.
  Invalid type отклоняется локально; несуществующий, но синтаксически корректный код
  проходит к серверу и возвращает обычную APIError при его отказе.
- [ ] GetPage(ctx, documentID, PageOptions) выполняет один GET
  entity/{documentType}/{documentID}/positions без тела; query содержит только limit/offset.
  Limit=0 означает 1000; 1…1000 допустимы, остальные ошибочны; Offset>=0.
  Default отправляет limit=1000&offset=0.
- [ ] documentID проверяется как UUID по образцу UOM, без ограничения версии.
  Invalid type/ID/page, nil client/context и отмена context проверены без лишних HTTP calls.
- [ ] HTTP 200 возвращает *model.Collection[T]: Context, CollectionMeta, Rows и href
  сохранены из ответа. Rows=[] — успешная пустая страница; top-level id и ID в T не нужны.
- [ ] Минимальная проверка envelope требует collection object, meta.href!=empty
  и non-null rows-array. Missing meta.href/rows, rows:null, {}, null, malformed JSON и
  пустое тело дают error+nil; 201/204 — неожиданный status. Size/offset/len(rows) не
  подвергаются бизнес-проверке; nextHref не обязателен.
- [ ] Plain/gzip, exact path/query, отсутствие тела, стандартные Accept/Accept-Encoding/
  Authorization/User-Agent и отсутствие Content-Type у GET проверены локальным HTTP-сервером.
- [ ] Отдельные public tests GetPage/GetPageErrors/GetPageValidation/ErrorWrapping
  используют структуру табличных UOM TestServiceGet и настоящий transport.Client,
  без моков private do или нового транспортного framework.
- [ ] Проверены default, Limit=1/1000, invalid -1/1001, Offset=0/>0/<0, допустимый
  UUID нестандартной версии, несколько явных независимых страниц и ответ вне первой
  группы документов. Пользовательский T без Core/assortment/id доказывает generic контракт.
- [ ] 4xx/5xx доступны как APIError через errors.As; исходные auth/transport errors и
  context.Canceled через errors.Is; operation context включает type/ID, moysklad:
  не повторяется. Любая ошибка возвращает nil вместо частичного результата.
- [ ] Счётчик HTTP calls остаётся 1 при nextHref/metadata links/неполных строках/ошибках;
  неизвестные поля не ломают чтение. Автоматических повторов или followups нет.
- [ ] Документация описывает отдельную страницу, ручное управление limit/offset,
  существующие metadata header positions и границы typed поддержки. Наличие generic
  reader не означает полную реализацию полей всех документов.
- [ ] После реализации пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check и отдельные Standards/Spec review относительно зафиксированного BASE.

## Blocked by

Нет. Reader использует существующие model.Collection и transport.Client; T any не
зависит от Core/Commercial из 01. Общая готовность чтения строк требует завершения 02,
entity tickets с типизированным GetPositions зависят от 02. JSON-срез 01 независим.

## Testing Decisions

согласованный seam — публичный GetPage через локальный HTTP-сервер и transport.Client.
Структурный prior art — UOM TestServiceGet, TestServiceGetErrors, validation и error
wrapping. Страница возвращает metadata+rows, поэтому проверка UOM «response has ID»
сюда не переносится. Фикстура page endpoint и synthetic boundary cases отделены
от upstream fixtures заголовка документа; существующий provenance сохраняется.

## Границы

После утверждения — один основной implementer и законченный reader со всеми ошибками.
Ticket утверждён пользователем; разрешена первая волна. Не добавлять document Get/Create/Update, отдельные mutations
позиций, batch, auto-pagination walker, expand, live API, .env, push или commits подготовки.
Общий transport не рефакторить. Новые решения контракта согласовать отдельно.

## Sources

- [Первичные факты и ссылки](../research.md).
- [Пример endpoint, envelope и limit/offset](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L2466-L2525).

## Comments

## Утверждение и ownership первой волны

Пользователь 2026-11-01 утвердил первую волну после исправления денежных типов и смысла nil курса.
JSON и HTTP seams этого ticket согласованы. Локальная реализация через TDD и commit разрешены;
независимые Standards/Spec review координирует исполнитель интеграции относительно planning BASE.
Не менять coverage.md/README и общие карты: документацию интегратор сведёт отдельным коммитом.
Файлы и трёхпоточный запуск закреплены в [handoff](../../document-workstream/wave-1.md).
