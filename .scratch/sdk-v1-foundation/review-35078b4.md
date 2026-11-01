# Foundation повторное review — 2026-11-01

BASE: `46e9d068dda6401fe6cb1306eb6c7af55139bc41`.
HEAD: `35078b498d514eeeb31b8874fe302bab1702f22a`.
Scope: четыре fix commits и ответ на первое review, 5 commits.
Метод: независимые Standards/Spec reviewers и публичные воспроизведения.

## Standards

Новых подтверждённых нарушений нет. Limiter списывает бюджет после получения
account slot, не удерживает account slot во время ожидания бюджета; cancellation
освобождает занятые slots. При росте Weight во время ожидания слот возвращается
перед повторным ожиданием бюджета. Косметических smell-замечаний нет.

Ограничение о создании WithResponseInfo заново на каждый вызов внесено в ADR/godoc.
Неполнота исходного исправления оценивается ниже по оси Spec.

## Spec

**P2 — ResponseInfo исправлен лишь для прямого вызова конструктора опции.**
`transport/options.go:66–71`. ADR 0001 требует: «Если ответа нет, info обнуляется».
Сброс при создании WithResponseInfo оставляет дефект при композиции публичного
RequestOption, даже если новую опцию создают для каждого вызова:

```go
func infoOption(info *transport.ResponseInfo) transport.RequestOption {
    return func(o *transport.RequestOptions) {
        transport.WithResponseInfo(info)(o)
    }
}
```

`svc.Get(ctx, "invalid", infoOption(&info))` возвращает UUID error до применения
опции, сохраняя прежние StatusCode=503, Attempts=3. Публичное воспроизведение:
`/private/tmp/foundation_r2_spec_info.go`.
Ограничение про повторное использование здесь неприменимо: wrapper свежий.
Нужна обработка опций и сброс metadata на входе публичной операции до локальной
валидации. Конструктор опции сам по себе этот контракт не обеспечивает.
Добавить regression test со свежей составной опцией.

Остальные три исходных замечания закрыты: rate budget списывается после account
slot; async определяется по итоговому decoded query, включая повторные значения;
minimumStock явно отклоняется и отложен в docs, как допускало первое review.
Оба первоначальных P1 закрыты.

## Verification

- Публичный async repro: literal/encoded value/encoded key — по одной отправке.
- Публичный metadata repro: fresh composed option сохраняет 503/3 при UUID error.
- Targeted transport regression tests с race, в том числе limiter, прошли.
- gofmt для изменённых Go files и `git diff --check 46e9d06...35078b4` — PASS.
- Первый полный race запуск заблокирован sandbox: localhost bind operation not
  permitted; запущен повтор с разрешением на локальные httptest listeners.
- Повтор `go test -race -count=1 -p 1 ./...` — PASS; `go vet ./... ./testdata` — PASS.
- Production code и статусы tickets не менялись, live API не использовался.

Standards: 0 новых замечаний. Spec: 1 замечание, максимум P2.

## Ответ и исправление

Замечание принято: сброс в конструкторе `WithResponseInfo` (`bc2dad4`) чинил прямой вызов, но не
составную опцию, которая применяется только в Send, после валидации.

Исправлено в `5867c94`:
- `transport.ResetResponseInfo(opts)` применяет опции к временной структуре и обнуляет найденный
  ResponseInfo при любой композиции; сброс из конструктора убран, механизм один.
- Вызывается первой строкой, до локальной проверки, во всех публичных операциях с опциями:
  Get/Create/Update и GetPositions services, `positions.GetPage`, `internal/list.Entity` (List),
  `security.CreateToken` — 61 место; однострочные List сразу делегируют в `list.Entity`.
- ADR 0001: «публичная операция обнуляет info на входе»; ограничение про переиспользование
  опции снято.
- Регрессия `TestLocalErrorResetsResponseInfo`: 17 services × Get/Create/Update/List и 7 GetPositions,
  прямая и свежая составная опция; `TestCreateTokenLocalErrorResetsInfo` — составная опция.
  На 35078b4 — 76 падений, после исправления — PASS.
- `pages.All` не менялся: его ResponseInfo по контракту описывает последнюю загруженную страницу.
