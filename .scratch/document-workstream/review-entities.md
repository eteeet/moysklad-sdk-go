# Завершение 11 entity tickets документов

Дата: 2026-11-01. Пользователь поручил выполнить все 11 оставшихся tickets SDK.
BASE: `eba6f75154929790c5b9d2cb1454dd80783b14db`, ветка dev, изменения в рабочем дереве.
В рамках этой задачи агент не выполнял commit/push, live API или изменение временных worktree.
Во время работы HEAD внешне переместился на `59b16e7` (уточнение единиц Discount и
удаление дублирующего test case); эти изменения сохранены. Фиксированный BASE review
остаётся `eba6f75`, итоговые проверки выполнены на текущем рабочем дереве.
Начальные изменения planning и docs/temp_plan.md сохранены. Нормативное уточнение
S02 Cost write, появившееся в рабочем дереве во время реализации, включено в результат.

## Реализация

PaymentIn/PaymentOut/CashIn/CashOut, InvoiceIn/InvoiceOut, Supply/Demand/PurchaseOrder,
Move и PurchaseReturn: отдельные read/Create/Update модели и entity Service.
Все семь товарных документов имеют aliases общих Position/PositionRequest и фиксированный
GetPositions через существующий reader. Общая Position дополнена Cost/Shipped/InTransit;
общая PositionRequest — Cost для услуг Demand и InTransit для PurchaseOrder.
Overhead/Shipped только читаются. Отдельного production CRUD/transport framework нет.

Зафиксированы различия required Create, optional Demand.Store, запреты
PurchaseReturn.Update, readonly/TODO и серверные ограничения без скрытого Get.
Field сохраняет omission/null/false/0/""/[], APIError/errors.Is/As и один запрос.
JSON-тесты обоими codecs, прямое сравнение Go-моделей, raw Attribute long >2^53,
независимые request golden, неизменённые upstream и полные unexpanded samples.
HTTP-тесты через настоящий transport и httptest: методы/пути/query/headers/body,
plain/gzip, фактический ответ, 1000/1001, context/auth/transport/APIError,
ошибочные статусы/JSON/envelope, отсутствие retry/follow-up.

## Standards

Первый проход: 3 heuristic замечания к тестам, без нарушений production-архитектуры.
Исправлены: повторная сериализация expected/actual заменена reflect.DeepEqual с
канонизацией только сырых Attribute values; Success задан явным bool; полные inputs
вынесены в constructors тестового пакета. При повторном проходе исправлен ещё P3:
t.Fatal из handler goroutine заменён t.Error + return в JSONEqual.
Итог независимого Standards review: **0 открытых замечаний**.

## Spec

Исправлены Cost write и его HTTP/JSON-контракты, noClosingDocs=true + operations
synthetic conflict для PaymentOut/CashOut, polymorphic PaymentOut.Agent
(counterparty/employee/organization). Устаревшие ограничения docs/positions.md обновлены.
Итог независимого Spec review: **0 открытых замечаний**, расширения scope не обнаружены.

## Проверки

- PASS: `go test -race -count=1 ./...`.
- PASS: `go vet ./... ./testdata`.
- PASS: gofmt изменённых/новых Go-файлов и `git diff --check`.
- PASS: SHA-256 всех 161 upstream fixture соответствует docs/upstream.md.
- Реальный Account не использовался; HTTP tests проверяют SDK и не доказывают
  бизнес-поведение сервера. Неоднозначности и отложенные поля описаны в docs/*.md.

Публичные контракты: [общие правила](../../docs/documents.md),
[перечень документов](../../README.md#документы). Tickets и их результаты — в карте.
