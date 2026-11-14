Status: ready-for-agent
Type: task
Blocked by: —

# 01: Create документов без `Name`

`~/ETEEET/moysklad/AGENTS.md:18`: документам не задают `name`, МойСклад нумерует их сам («последний + 1»). Владелец подтвердил 2026-11-14, что правило распространяется на валидацию SDK. Пять сервисов документов всё ещё отказывают в Create при пустом `Name`:

- `entity/invoiceout/service.go:43`
- `entity/retailsalesreturn/service.go:44`
- `entity/processingorder/service.go:42`
- `entity/internalorder/service.go:47`
- `entity/commissionreportout/service.go:51`

Для invoicein то же самое сделано в v0.2.2 (`CHANGELOG.md:7,21-22`): живая проверка показала автонумерацию, `CreateRequest.Name` получил `omitzero`.

## What to build

В пяти сервисах Create не проверяет `Name`, а `CreateRequest.Name` не уходит в теле, когда он пустой (`json:"name,omitzero"`, комментарий как у `entity/invoicein/request.go:10`). Остальные проверки Create не меняются. Справочники, у которых имя обязательно, в тикет не входят.

## Acceptance criteria

- [ ] Тест на `httptest` для каждого из пяти сервисов: Create без `Name` отправляет тело без ключа `name`; Create с `Name` отправляет его как раньше.
- [ ] Живая проверка на `admin@documan` для каждого из пяти документов: Create без `name` → 200, номер присвоен автонумерацией; метка проверки — в `description`; результат с датой записан в тикет. Документ, для которого МойСклад не нумерует сам, оставляет проверку, и это записывается в тикет.
- [ ] `docs/documents.md:16-18` описывает проверки Create по коду (сейчас там ещё InvoiceIn и нет ProcessingOrder); `CHANGELOG.md` — запись к следующему выпуску.
- [ ] `make check` зелёный.
