Status: wontfix
Stage: SUPERSEDED
Blocked by: нет

# 03: Отдельная форма позиции — заменена универсальной

## Решение

2026-11-01 пользователь отменил отдельный контракт P03 и утвердил одну форму
`positions.Position` / `PositionRequest` для всех товарных документов.
Этот ticket больше не запускается и не блокирует M01.

## Answer

Предыдущая реализация 519b7a6, документация f26ec32 и regression fix bf2eaac
сохраняются в Git-истории. Четыре movement*.go файла удаляются; overhead и ссылки
на ячейки, включая полезные JSON regression tests, переносятся в общую форму.
Deprecated aliases не нужны: SDK не выпущен и потребителей старого типа нет.

Активный контракт — [spec positions](../spec.md),
[согласованная замена](../../universal-positions/spec.md)
и [ticket объединения](../../universal-positions/issues/01-unify.md).
Замена завершена: код 1a9d400, публичная документация 7d43628;
entity Move остаётся needs-triage.
