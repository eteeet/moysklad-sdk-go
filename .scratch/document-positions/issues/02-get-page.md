Status: resolved
Stage: IMPLEMENTED
Blocked by: None

# 02: Generic reader одной страницы позиций

## What to build

Разработчик решения читает одну страницу позиций любого document type с поддерживаемым
API endpoint через общий positions.Service[T any].GetPage и настоящий transport.Client;
получает model.Collection[T] с metadata и строками, управляет следующей страницей сам.
Срез закончен от public API до локального HTTP-теста и обработки ошибок.

Interface и HTTP seam P02 утверждены пользователем: [spec](../spec.md).
Ticket завершён; entity endpoints в эту волну не входят.

## Acceptance criteria

- [x] `NewService[T](client *transport.Client, documentType string)` возвращает
  (*Service[T], error), nil client отклоняется; T any не требует Core, assortment или ID.
- [x] documentType — открытая строка одного безопасного API route segment; согласованный
  syntax `[a-z][a-z0-9]*`, без allowlist первых документов и скрытой нормализации.
  Invalid type отклоняется локально; несуществующий, но синтаксически корректный код
  проходит к серверу и возвращает обычную APIError при его отказе.
- [x] GetPage(ctx, documentID, PageOptions) выполняет один GET
  entity/{documentType}/{documentID}/positions без тела; query содержит только limit/offset.
  Limit=0 означает 1000; 1…1000 допустимы, остальные ошибочны; Offset>=0.
  Default отправляет limit=1000&offset=0.
- [x] documentID проверяется как UUID по образцу UOM, без ограничения версии.
  Invalid type/ID/page, nil client/context и отмена context проверены без лишних HTTP calls.
- [x] HTTP 200 возвращает *model.Collection[T]: Context, CollectionMeta, Rows и href
  сохранены из ответа. Rows=[] — успешная пустая страница; top-level id и ID в T не нужны.
- [x] Минимальная проверка envelope требует collection object, meta.href!=empty
  и non-null rows-array. Missing meta.href/rows, rows:null, {}, null, malformed JSON и
  пустое тело дают error+nil; 201/204 — неожиданный status. Size/offset/len(rows) не
  подвергаются бизнес-проверке; nextHref не обязателен.
- [x] Plain/gzip, exact path/query, отсутствие тела, стандартные Accept/Accept-Encoding/
  Authorization/User-Agent и отсутствие Content-Type у GET проверены локальным HTTP-сервером.
- [x] Отдельные public tests GetPage/GetPageErrors/GetPageValidation/ErrorWrapping
  используют структуру табличных UOM TestServiceGet и настоящий transport.Client,
  без моков private do или нового транспортного framework.
- [x] Проверены default, Limit=1/1000, invalid -1/1001, Offset=0/>0/<0, допустимый
  UUID нестандартной версии, несколько явных независимых страниц и ответ вне первой
  группы документов. Пользовательский T без Core/assortment/id доказывает generic контракт.
- [x] 4xx/5xx доступны как APIError через errors.As; исходные auth/transport errors и
  context.Canceled через errors.Is; operation context включает type/ID, moysklad:
  не повторяется. Любая ошибка возвращает nil вместо частичного результата.
- [x] Счётчик HTTP calls остаётся 1 при nextHref/metadata links/неполных строках/ошибках;
  неизвестные поля не ломают чтение. Автоматических повторов или followups нет.
- [x] Документация описывает отдельную страницу, ручное управление limit/offset,
  существующие metadata header positions и границы typed поддержки. Наличие generic
  reader не означает полную реализацию полей всех документов.
- [x] После реализации пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check и отдельные Standards/Spec review относительно зафиксированного BASE.

## Blocked by

Нет. Reader использует существующие model.Collection и transport.Client; T any не
зависит от Core/Position из 01. Общая готовность чтения строк требует завершения 02,
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

## Answer

Реализованы `positions.Service[T any]`, `NewService[T]`, `GetPage` и `PageOptions`.
Reader использует существующие `transport.Client` и `model.Collection[T]`, выполняет
один GET с явными limit/offset, проверяет UUID/страницу и минимальный envelope ответа.
Серверные Context/Meta/Rows сохраняются без пересчёта; ID и общие позиции в T не требуются.
GoDoc описывает отдельную страницу, ручное управление offset/limit, metadata заголовка
документа и границы generic/typed поддержки. P01 и P03 не изменялись.

TDD через публичный GetPage и локальный HTTP-сервер:

- Первый red до production: `go test -count=1 ./positions -run '^TestServiceGetPage$'`
  завершился с `positions: no non-test Go files`; после минимального reader — green.
- Следующие red→green: nil client возвращал service; invalid documentType принимался;
  validation выполняла 13 HTTP calls вместо 0; 201/204 и неверные envelopes принимались.
- Финальные тесты проверяют plain/gzip, заголовки и отсутствие GET body/Content-Type,
  безопасные открытые коды, custom T без Core/ID, неизвестные поля и неполные строки,
  default/границы страницы, UUID без ограничения версии, пустую страницу, две явные
  независимые страницы, invalid envelopes и частичный decode, APIError/errors.Is,
  nil/canceled/expired context и отмену во время HTTP. Retry/followups отсутствуют.

Проверки пройдены: gofmt; `go test -count=1 ./positions`;
`go test -race -count=1 ./...`; `go vet ./... ./testdata`; `git diff --check`.
Только локальные ответы: live API не проверялся. Independent Standards/Spec review
на этапе implementer не запускались; итог интеграции и независимых review приведён ниже.

### Итог интеграции

Включено в dev, commit реализации: 1b2860389302361a81228d8704be10cfa7f0df99.
Независимые Standards: 0 замечаний; Spec: 0 замечаний. Проверен полный срез
от planning BASE a2804283e7aafe3c6ee41c17a4c2c79f7e01f999 до HEAD 1b28603,
включая централизованные docs. [Отчёты](../../document-workstream/review-wave-1.md).
Общий go test -race -count=1 ./..., go vet ./... ./testdata, gofmt и diff-check пройдены.
Публичная документация сведена интегратором после слияния; live API и push не выполнялись.
