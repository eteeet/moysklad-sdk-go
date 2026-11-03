# Повторное review memnet — 2026-11-03

BASE: `c9399716f063284840341b2410957e7ce9edb460`.
HEAD: `c3bace0997f7b8b4b08ea6fd8caf6ceff9279621`.
Scope: c3bace0, два файла internal/memnet; working AGENTS.md diff остался прежним.
Метод: независимые Standards/Spec reviewers, полный make check и исходный repro root.

## Standards

Замечаний нет. Close сохраняет регистрацию адреса и закрывает done через sync.Once.
Оба конца net.Pipe в ветке ошибки закрываются; ошибка имеет форму net.OpError
с ECONNREFUSED. Tombstone обоснован запретом переиспользовать выданные адреса.
Новый regression test занимает тот же адрес настоящим TCP listener и проверяет
ошибку клиента при нулевых обращениях к TCP handler.

## Spec

Замечаний нет. Прежний P2 закрыт: запрос к закрытому memnet-адресу больше не
попадает в настоящий DialContext. Исходный детерминированный repro теперь проходит.
HTTP serialization и работа настоящих http.Client/Transport/Server сохраняются.
Прежние замечания wizard/docs были закрыты в review-c939971.

## Verification

- `make check` — PASS: gofmt, go vet ./... ./testdata,
  go test -race -count=1 -p 2 ./..., tracked HEAD и working tree diff-check.
- Исходный repro через `/private/tmp/memnet-review-overlay.json` — PASS:
  `go test -race -count=1 -overlay /private/tmp/memnet-review-overlay.json
  -run '^TestReviewClosedAddressDoesNotReachTCP$' ./internal/memnet`.
  Исходники SDK overlay не меняет; listener локальный, управляемый тестом.
- `git diff --check c939971...c3bace0` — PASS. HEAD оставался c3bace0.
- Code, planning, tracker, working AGENTS.md и commits reviewer не менял.
  Создан только этот отчёт; .DS_Store не читались/не изменялись; live API не вызывался.

Все замечания review-5de2fa7 и review-c939971 закрыты.
Standards: 0 findings. Spec: 0 findings.
