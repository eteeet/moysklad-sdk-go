# Review захватки 8 — 2026-11-02

BASE: `7fb0c63b30dc784d296856c3ebd7eb4b1e1812d3`.
HEAD: `0fafaaa3e60c6b25129ee4086c98269e0d024075`.
Scope: X01 Delete, X02 batch, X03 async, сопутствующий assortment.All и docs.
Метод: независимые Standards и Spec reviewers; дополнительная проверка root,
полный make check и воспроизведение через публичные сервисы с настоящим transport.

## Standards

### S1. P2 — null элемента batch принимается как успешный результат

Место: `batch/batch.go:42`, `batch/batch.go:63`.

HTTP 200 с `[null]` проходит проверку количества элементов в internal/write.post.
CreateMany возвращает err=nil и Value!=nil с пустой моделью; DeleteMany возвращает
err=nil и Deleted без Info и Errors. Это нарушает контракт результата: записанный
объект либо ошибки элемента, сообщение об удалении либо ошибки. Вызывающий код
получает признак успеха при неизвестном исходе записи.

Подтверждено обоими JSON codecs и публичными Product.CreateMany/DeleteMany.
В docs/documents.md требуется проверять malformed/null/empty ответы; имеющиеся
тесты отвергают null всего массива, но не его элемента.

Исправление: отвергать null элементы и проверять это через публичные методы;
не добавлять необоснованные обязательные поля всем generic T.

Других существенных нарушений стандартов и baseline smells не обнаружено.
Общая логика находится в internal/write; копия validID разрешена coverage.md.

## Spec

### P1. P2 — отдельная policy очереди async отсутствует

Место: `async/async.go:5`, `internal/write/write.go:171`.

X03 в spec.md и issues/03-async.md требует «≤4 queued account tasks policy отдельно
от HTTP semaphore». Общий spec утверждает реализацию квоты в X02/X03. Фактически
код документирует серверный лимит и отправляет каждый ListAsync сразу; общей
policy жизненного цикла задач нет. HTTP semaphore освобождается после ответа 202,
хотя созданная задача продолжает занимать очередь. Полностью контролировать другие
процессы невозможно, но это не отменяет заявленного локального shared контракта.

Исправление: реализовать отдельную разделяемую policy с проверкой жизненного цикла
и неизвестного исхода submit либо явно согласовать изменение scope на server-only
и привести spec/master/ticket в соответствие. Запись реализации в Answer сама по
себе не подтверждает согласованное изменение требования.

### P2. P2 — Result принимает неожиданные успешные HTTP-статусы

Место: `async/async.go:152`.

Result делегирует транспортному Download, который принимает любой 2xx. Поэтому
HTTP 204 даёт n=0, err=nil — успешный пустой результат. Источник `_async.md`
описывает Response 200 либо 302 на загрузку; остальные lifecycle методы уже
проверяют свой конкретный статус. Это позволяет вызывающему коду принять отсутствие
результата за успешную загрузку.

Исправление: проверять конечный статус 200 перед записью тела результата,
покрыть 202/204 и сохранение ResponseInfo. Общий контракт Download для файлов
не требуется сужать ради одного async endpoint.

Capabilities Delete/CreateMany/UpdateMany/DeleteMany соответствуют проверенным
разделам pinned источника; необоснованного расширения методов не найдено.

## Verification

- `make check` — PASS: gofmt, go vet ./... ./testdata,
  go test -race -count=1 -p 2 ./..., tracked HEAD и working tree diff-check.
- Batch repro: `/private/tmp/review-wave8-result.go` — воспроизведён root.
  HTTP200 `[null]`: CreateMany err=nil, Value!=nil; DeleteMany err=nil, пустой Deleted.
- Async repro: `/private/tmp/moysklad-async-spec-repro.go` — воспроизведён root.
  HTTP204: Result n=0, err=nil.
- Repros используют custom RoundTripper без сети, live API не вызывался.
- До отчёта working tree был чистым, HEAD оставался 0fafaaa.
- Production code, tracker, commits reviewer не менял.

Standards: 1 finding P2. Spec: 2 findings P2. Захватка требует правок.

## Ответ и исправления

Согласен с S1 и P2; P1 закрыт решением пользователя.

- S1 — `0294292` fix: `batch.Result` отвергает `null` и не объект, `batch.Deleted` — ещё и объект
  без `info` и `errors` (свой UnmarshalJSON). Обязательные поля моделям не добавлены: `{}` у
  CreateMany/UpdateMany по-прежнему разбирается как объект. Тесты — оба кодека
  (`TestRejectsUnknownItems`) и публичные CreateMany/UpdateMany/DeleteMany всех сервисов
  (`TestBatchRequestErrors`, случай `[null,null]`).
- P2 — `6b1bc33` fix: `async.Result` проверяет итоговый статус 200 через ResponseInfo вызывающего
  кода или свой; 202/204 — ошибка, `Download` для файлов не менялся. Тело ответа не 200 могло уже
  попасть в w — это видно по n и записано в godoc. Тест `TestResultRequiresOK`: 202/204 с
  ResponseInfo и без, 200.
- P1 — решение пользователя 2026-11-02: лимит 4 задач — серверный. Учесть задачи других
  процессов и завершённые без опроса можно только скрытым GET async перед ListAsync, а он запрещён
  общим контрактом. Spec X03, issue, общий spec, coverage и карта приведены к этому; SDK
  документирует лимит и передаёт ошибку 61002.
- Docs: docs/write-async.md — null-элементы, 200 у Result, серверный лимит.
