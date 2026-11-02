# Повторное review захватки 8 — 2026-11-03

BASE: `0fafaaa3e60c6b25129ee4086c98269e0d024075`.
HEAD: `565cfa3e2ab23a0263f7f09d4aed436c96e12834`.
Scope: исправления 0294292, 6b1bc33 и ответ/docs 565cfa3.
Метод: независимые Standards/Spec reviewers и проверка root.

## Standards

### S1. P2 — Result вызывает panic при допустимой nil-опции

Место: `async/async.go:157–158`.

Новый цикл безусловно вызывает `opt(&options)`. Вызов
`Result(ctx, validID, io.Discard, nil)` поэтому вызывает panic до HTTP.
Публичный контракт `transport/options.go:8`: «nil игнорируется»;
существующий newRequestOptions проверяет opt != nil. Раньше Result передавал
опции прямо в Download и соблюдал этот контракт.

Исправление: пропускать nil и добавить regression test публичного Result,
в том числе с nil рядом с WithResponseInfo. Ожидается обычная загрузка и
сохранение метаданных вызывающего кода.

Предыдущий S1 закрыт: null и не объекты отвергаются; Deleted также отвергает
пустой результат без info и errors. Необоснованных обязательных полей у generic T
не появилось. Других нарушений стандартов и существенных baseline smells нет.

## Spec

### P1. P2 — nil RequestOption нарушает общий контракт публичной операции

Место: `async/async.go:157–158`.

Это та же регрессия, независимо найденная Spec reviewer: общий interface ADR 0001
использует transport.RequestOption, чей контракт разрешает nil. Исправление
проверки статуса не должно превращать допустимый вызов в panic.

Предыдущие замечания закрыты:

- Batch отвергает null элементы, включая публичные массовые методы.
- Result возвращает ошибку для итоговых 202/204 и сохраняет caller ResponseInfo.
  Возможная запись тела неожиданного 2xx до проверки статуса явно оговорена godoc;
  n/error сохраняют контракт потоковой загрузки. Общий Download не изменён.
- Решение о серверном контроле очереди вместо локальной policy согласованно
  записано в optional spec, X03, master spec, coverage, map и guide.

Других пропущенных требований или scope creep в исправлениях не найдено.

## Verification

- `make check` — PASS: gofmt, go vet ./... ./testdata,
  go test -race -count=1 -p 2 ./..., diff-check полного tracked HEAD и working tree.
- `git diff --check 0fafaaa...565cfa3` — PASS.
- Finding подтверждён прямым чтением безусловного вызова nil-функции и контракта
  RequestOption. Отдельный go run repro не включён в результат make check.
- До отчёта HEAD оставался 565cfa3; tracked files не изменены.
- Существующий untracked .DS_Store не читался и не изменялся.
- Production code, tracker и commits reviewer не менял; live API не использовался.

Standards: 1 finding P2. Spec: 1 finding P2. Одна уникальная регрессия; требуется правка.

## Ответ и исправления

Согласен: S1 и P1 — одна регрессия из `6b1bc33`. Других мест с собственным применением
опций нет: в production-коде цикл есть только в `transport.newRequestOptions` и здесь.

- `ba1d2f1` fix: `async.Result` пропускает nil-опцию, как `newRequestOptions`. Тест
  `TestResultNilOption`: публичный Result только с `nil` даёт обычную загрузку, а с `nil`
  рядом с `WithResponseInfo` заполняет ResponseInfo вызывающего кода статусом 200. Без
  исправления тест падает с panic.
- Docs и CHANGELOG не меняются: поведение соответствует уже описанному контракту, и
  регрессия не попала в релиз.
