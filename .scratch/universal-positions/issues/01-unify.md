Status: resolved
Stage: IMPLEMENTED
Blocked by: нет

# 01: Объединить формы позиций

Выполнить [согласованный контракт](../spec.md) одним implementer.

- [x] Общие Position/PositionRequest вместо Commercial и Movement, без aliases старых имён.
- [x] Поля и оба JSON codecs, перенос regression tests, read-only Overhead.
- [x] Публичная документация и правило общей формы.
- [x] Актуальные planning и зависимости, P03 заменён, будущие поля в общем positions.
- [x] Полные проверки и независимые Standards/Spec review, локальные code/docs commits.

## Answer

2026-11-01: код 1a9d400, публичная документация 7d43628; итоговый HEAD eba6f75
(косметическая правка подписи assertion). Position встраивает Item,
PositionRequest — ItemRequest. Overhead только чтение, optional ячейки omitzero,
в запросе все optional поля Field. Старые типы/aliases и 4 movement*.go файла удалены;
полезные проверки перенесены в общие model_test.go/request_test.go.

Red→green отдельно подтверждён для новых экспортных имён и перенесённых read/write
полей. Full race/vet/gofmt/diff-check PASS. Independent Standards 0, Spec 0 после
устранения одного P2 противоречия planning возврата. [Отчёт](../review.md).

P03 wontfix/superseded, M01 от него не зависит. Return/stock/invoice spec и tickets
используют общую форму; S02 Cost и S03 Shipped/InTransit запланированы в общий positions,
сейчас не реализованы. Остальные entity tickets не запущены. Planning/tracker/temp_plan
не коммитились по действующему правилу. Пять dirty worktree не затронуты.
