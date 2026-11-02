# Files/Reports повторное review — 2026-11-02

BASE: `daadf36d993ef8ab3ba1b7680d8aea8e5b142e96`.
HEAD: `d6a4b6cfde2c43bc9d30b1eed90e56424cf3a37d`.
Scope: fix c555b69 и ответ d6a4b6c.
Метод: независимые Standards/Spec reviewers и публичные воспроизведения root.

## Standards

**P2 — подпись malformed Location всё ещё раскрывается.**
`transport/download.go:110`: regex `\?[^\s"]*` прекращает удаление на пробеле
или кавычке внутри query. Location
`https://storage.invalid/%zz?pad=x signature=top-secret` возвращает ошибку
с signature=top-secret дважды; аналогично `?pad="&signature=top-secret`.
Нарушены F02 «signed external URL secret не логируется» и гарантия docs/godoc
скрывать query, включая вложенную ошибку неразбираемого Location.
Санитизировать полное значение URL в сообщении, включая пробелы и экранированные
кавычки; добавить regression через публичный Download.

Полученный response при отказе CheckRedirect теперь сохраняется корректно.
Принятое ограничение StatusCode=0 для неразбираемого Location документировано
и не считается замечанием. Других Standards findings нет.

## Spec

1. **P1 — санитизация пропускает часть query.**
   `transport/download.go:110`: те же воспроизводимые URL выше оставляют подпись
   в Error(). Это незакрытая гарантия F02, а не возражение против принятого
   ограничения net/http. Полная очистка нужна и для query с пробелами/кавычками.
   Это то же замечание, что Standards; приоритеты reviewers независимы.

2. **P2 — redactedError меняет timeout classification и отключает retry.**
   `transport/download.go:85–86,114–117`. F01/ADR требуют общий retry; исправление
   обещает сохранять классификацию ошибок. При поддерживаемом custom RoundTripper,
   который возвращает net.Error с Timeout=true и URL в тексте, redactedError
   скрывает Timeout от внешнего url.Error. errors.As на net.Error находит внешний
   url.Error с Timeout=false. Retry{MaxAttempts:2} выполняет одну попытку.
   Сохранять timeout/temporary semantics при редактировании текста и добавить
   regression с такой сетевой ошибкой.

## Исправленное

- Обычный malformed Location с query без пробелов/кавычек больше не раскрывает подпись.
- При отказе CheckRedirect сохраняются HTTP 302, headers и Attempts=1;
  errors.Is исходного sentinel работает, лишнего retry нет.
- Для неразбираемого Location net/http не отдаёт response: StatusCode=0 принят
  пользователем, отражён в godoc, ADR и guide. Обёртка Transport не требуется.

## Verification

- `go test -race -count=1 -p 2 ./...` — PASS.
- `go vet ./... ./testdata` — PASS; gofmt изменённых Go files — PASS.
- `git diff --check daadf36...d6a4b6c` — PASS.
- Прежний repro `/private/tmp/files-reports-standards-repro.go` подтверждает
  исправление простого случая и metadata CheckRedirect.
- Новые публичные воспроизведения без сети:
  `/private/tmp/download_review_redaction_boundaries.go`,
  `/private/tmp/files-reports-standards-fix-repro.go`,
  `/private/tmp/files-reports-spec-rereview.go`.
- Root повторил оба дефекта: secret присутствует в ошибке; calls=1, Timeout=false.
- HEAD оставался d6a4b6c; рабочее дерево до отчёта чистое.
- Production code и статусы tickets не менялись; live API не использовался.

Standards: 1 замечание, максимум P2. Spec: 2 замечания, максимум P1.
Всего 2 уникальных дефекта.

## Ответ и исправления

Оба замечания приняты. Исправление c555b69 и ответ d6a4b6c не были запушены, поэтому удалены из
истории. Вместо них сделан один коммит `f2e2bf8` поверх daadf36, который закрывает и это review, и
[review-daadf36](review-daadf36.md).

| Замечание | Где в f2e2bf8 |
| --- | --- |
| P1/P2, query с пробелом или кавычкой не скрывается | `redactedError.Error` обрезает текст на первом `?`; regex удалён |
| P2, redactedError прячет Timeout и отключает retry | `redactedError.Timeout`/`Temporary` пробрасывают значения исходной ошибки |

- Regression через публичный Download: неразбираемый Location с query plain, с пробелом
  и с кавычкой — подписи в тексте нет, errors.As(*url.Error) работает. Свой RoundTripper
  возвращает net.Error с Timeout=true и подписанной ссылкой в тексте. При MaxAttempts 2 выполняются
  две попытки, Timeout виден через errors.As(net.Error) и errors.Is работает, подписи в тексте нет.
- На daadf36 оба новых теста падают.
- go test -race -count=1 -p 2 ./..., go vet ./... ./testdata, gofmt, `git diff --check daadf36 f2e2bf8` — PASS
  в чистом worktree на f2e2bf8.
