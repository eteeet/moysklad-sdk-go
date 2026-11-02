# Files/Reports review — 2026-11-02

BASE: `5a703f2c24d77ee19e58853c1d27d42076aa6ef2`.
HEAD: `daadf36d993ef8ab3ba1b7680d8aea8e5b142e96`.
Scope: F01–F05, общий test commit и завершающий docs commit; 7 commits, 33 files.
Метод: независимые Standards/Spec reviewers и публичные воспроизведения.
Итог: нужны исправления двух дефектов Download.

## Standards

1. **P2 — подпись невалидного redirect попадает в текст ошибки.**
   `transport/download.go:78–83` очищает только внешний url.Error.URL.
   Вложенная Err сохраняет полный Location. Для
   `https://storage.invalid/%zz?signature=top-secret` ошибка содержит подпись дважды.
   Нарушены требование F02 «signed external URL secret не логируется» и
   опубликованная гарантия скрывать query подписанной ссылки.
   Очищать представление всей цепочки ошибок redirect, сохраняя предусмотренные
   errors.Is/As; добавить regression через публичный Download.

2. **P2 — ResponseInfo теряет полученный redirect-ответ.**
   `transport/download.go:83`: при отказе CheckRedirect http.Client.Do возвращает
   одновременно response и error, но downloadOnce подменяет response на nil.
   Полученный 302 с X-Source: api превращается в StatusCode=0, Header=nil, Attempts=1.
   Нарушен ADR 0001: status и headers последнего полученного ответа сохраняются
   при ошибках. Возвращать доступный response вместе с ошибкой; проверить metadata.

Других подтверждённых Standards defects и убедительных baseline smells не найдено.

## Spec

1. **P1 — утечка подписи malformed redirect.**
   `transport/download.go:78–83`. F02 требует не раскрывать secret подписанного URL,
   Download обещает исключать query из текста ошибки. Вложенная url.Error.Err
   обходит текущую очистку. Публичное воспроизведение с Location выше возвращает
   signature=top-secret дважды. Исправить санитизацию всей отображаемой цепочки,
   добавить regression Download/Files.Download, сохранить классификацию ошибок.
   Это тот же дефект, что Standards №1; оценки приоритета независимы.

2. **P2 — потеря status/headers при ошибке CheckRedirect.**
   `transport/download.go:83`. ADR 0001 и F01 Answer требуют metadata последнего
   ответа redirect chain. Полученный 302 доступен из http.Client.Do, но теряется
   в downloadOnce. Публичное воспроизведение подтверждает нулевой ResponseInfo.
   Это тот же дефект, что Standards №2.

По Files и Turnover иных замечаний нет: owner paths, Upload Base64/1–10 файлов,
отсутствие pre-GET, encoded body limit, отдельные endpoints/row shapes/filter/page
соответствуют freeze. Противоречия дат ByOperations, stockByStore и operation.moment
зафиксированы; пять новых fixtures совпадают с pinned source и SHA-256 provenance.

## Verification

- `go test -race -count=1 -p 2 ./...` — PASS.
- `go vet ./... ./testdata` — PASS.
- gofmt изменённых Go files и `git diff --check 5a703f2...daadf36` — PASS.
- `/private/tmp/files-reports-standards-repro.go` — оба дефекта воспроизведены root
  и независимыми reviewers через публичные NewClient/Download и fake RoundTripper.
- Обрыв, writer error, limiter lifetime, отсутствие retry после partial output,
  headers/auth redirect и общая регрессия JSON transport просмотрены.
- HEAD оставался daadf36; рабочее дерево до отчёта чистое.
- Production code и статусы tickets не менялись; live API не использовался.

Standards: 2 замечания, максимум P2. Spec: 2 замечания, максимум P1.
Всего 2 уникальных дефекта.

## Ответ и исправления

Оба дефекта приняты. Второй принят с уточнением: ответ можно сохранить только при отказе CheckRedirect.
Первое исправление (c555b69 и ответ d6a4b6c) не прошло повторное review
([review-d6a4b6c](review-d6a4b6c.md)). Оба коммита не были запушены; они удалены из истории
и заменены одним коммитом `f2e2bf8` поверх daadf36. Коммиты tickets не переписывались.

| Замечание | Commit | Где |
| --- | --- | --- |
| P1/P2, подпись неразбираемого Location в тексте ошибки | `f2e2bf8` | transport/download.go, тесты transport и files |
| P2, ResponseInfo теряет 302 при отказе CheckRedirect | `f2e2bf8` | transport/download.go, тест transport |

- Подпись: если текст вложенной url.Error.Err содержит `?`, ошибка оборачивается в redactedError. Его текст
  обрезается на первом `?`, потому что границу query в тексте не угадать: в ней бывают пробелы и кавычки.
  Unwrap сохраняет errors.Is/As, а Timeout/Temporary пробрасываются, так что классификация повторов не меняется.
- ResponseInfo: downloadOnce возвращает ответ, который http.Client.Do отдаёт вместе с ошибкой
  CheckRedirect, как уже делает Client.Do для JSON. Тест: StatusCode 302, X-Source, Attempts 1, без повтора.
- Не исправляется: при неразбираемом Location net/http ответа не возвращает, и StatusCode
  остаётся 0. Ограничение записано в godoc Download, docs/files-reports.md, ADR 0001 и
  Answer F01. Обёртку Transport ради этого случая не вводим.
- go test -race -count=1 -p 2 ./..., go vet ./... ./testdata, gofmt, `git diff --check daadf36 f2e2bf8` — PASS
  в чистом worktree на f2e2bf8.
