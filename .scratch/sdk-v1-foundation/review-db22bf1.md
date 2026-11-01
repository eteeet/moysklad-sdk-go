# Foundation повторное review — 2026-11-01

BASE: `35078b498d514eeeb31b8874fe302bab1702f22a`.
HEAD: `db22bf15dfc6b70a36ae19b259cb6eb1cd93eb8c`.
Scope: fix `5867c94` и docs `db22bf1`.
Метод: независимые Standards/Spec reviewers, публичные воспроизведения.

## Standards

Подтверждённых замечаний нет. ResetResponseInfo выполняется до локальной валидации
во всех 17 services, typed GetPositions, generic positions.GetPage, List через
list.Entity и security.CreateToken. ADR согласован с реализацией.
Новых smell-замечаний нет.

## Spec

Замечаний нет. Предыдущий P2 закрыт: прежний сценарий со свежей составной опцией
возвращает UUID error с нулевым ResponseInfo. WithResponseInfo больше не меняет
состояние при создании; повторное использование значения опции также корректно.
Регрессии охватывают 17 services × Get/Create/Update/List, семь GetPositions
и составную опцию security.CreateToken.

Штатные опции и их композиция корректны. Повторное применение опций во входном
сбросе и Send не нарушает согласованный контракт: однократность применения
произвольной stateful RequestOption не обещана.

## Verification

- `go test -race -count=1 -p 1 ./...` — PASS.
- `go vet ./... ./testdata` — PASS.
- gofmt изменённых Go files и `git diff --check 35078b4...db22bf1` — PASS.
- Прежний публичный repro `/private/tmp/foundation_r2_spec_info.go` — дефект устранён.
- `/private/tmp/foundation_r3_reused_info.go` — повторное использование опции PASS.
- HEAD в начале и конце проверки — db22bf1; рабочее дерево до отчёта чистое.
- Production code и статусы tickets не менялись; live API не использовался.

Все четыре замечания первого review Foundation и уточнение второго review закрыты
в проверенных срезах. Это результат проверки Foundation, а не всей SDK v1.

Standards: 0 замечаний. Spec: 0 замечаний.
