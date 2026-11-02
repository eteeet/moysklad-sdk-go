# Files/Reports повторное review — 2026-11-02

BASE: `daadf36d993ef8ab3ba1b7680d8aea8e5b142e96`.
HEAD: `5b9723cd7800d1aa41e30cc1b3028c604d910997`.
Scope: пересобранный fix f2e2bf8 и ответ 5b9723c.
Прежние c555b69/d6a4b6c заменены; общий merge-base с прежним срезом — daadf36.
Метод: независимые Standards/Spec reviewers и повторение всех публичных воспроизведений.

## Standards

Замечаний нет. Подпись скрывается для обычного query, query с пробелом и кавычкой.
Отказ CheckRedirect сохраняет 302, headers, Attempts=1 и исходную ошибку через errors.Is.
Timeout/Temporary делегируются исходной ошибке; Unwrap сохраняет errors.Is/As.
Новых убедительных baseline smells не обнаружено.

## Spec

Замечаний нет. Все прежние findings закрыты:

- redactedError обрезает вложенный текст на первом вопросительном знаке;
  подпись не остаётся после пробела или кавычки внутри query.
- Классификация timeout сохраняется; Config.Retry с MaxAttempts=2 выполняет
  две попытки. Через errors.As(net.Error) Timeout=true, errors.Is также работает.
- Полученный ответ при отказе CheckRedirect сохраняется для ResponseInfo.
- StatusCode=0 для неразбираемого Location соответствует принятому исключению
  в ADR/godoc/guide, Transport wrapper не требуется.

Regression tests покрывают все ранее выявленные случаи; новых нарушений F01/F02 нет.

## Verification

- `go test -race -count=1 -p 2 ./...` — PASS.
- `go vet ./... ./testdata` — PASS.
- gofmt изменённых Go files и `git diff --check daadf36...5b9723c` — PASS.
- Root повторил публичные воспроизведения без сети:
  `/private/tmp/files-reports-standards-repro.go`,
  `/private/tmp/download_review_redaction_boundaries.go`,
  `/private/tmp/files-reports-spec-rereview.go`.
- Результаты: leaked=false для трёх query; timeout calls=2, Timeout=true;
  CheckRedirect info.StatusCode=302 и headers сохранены.
- HEAD оставался 5b9723c; рабочее дерево до отчёта чистое.
- Production code и статусы tickets не менялись; live API не использовался.

Все замечания Files/Reports закрыты в проверенном срезе.
Standards: 0 замечаний. Spec: 0 замечаний.
