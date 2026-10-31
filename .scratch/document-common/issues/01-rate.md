Status: ready-for-agent
Stage: APPROVED
Blocked by: None

# 01: Публичные DocumentRate и DocumentRateRequest

## What to build

Разработчик читает и формирует вложенный rate любого документа через общие типы
model.DocumentRate/DocumentRateRequest. Срез включает публичный JSON-контракт,
оба JSON-кодека, независимые ожидаемые тела и документацию.
Утверждённый контракт — [spec](../spec.md). Разрешена реализация первой волны.

## Acceptance criteria

- [ ] DocumentRate.Currency — *Reference, Value — *float64; nil Value внутри присутствующего rate означает курс 1, но SDK не подставляет 1 в JSON; над полем краткий комментарий об этом.
- [ ] DocumentRateRequest содержит Currency Field[Reference] и Value Field[float64] с omitzero.
- [ ] Внешний Field[DocumentRateRequest] различает unset/null/{}/значение; внутренние поля
  сохраняют null, 0, 1, дробный курс и пропуски.
- [ ] Оба кодека проверяют значения вложенной Meta, optional/null/unknown поля и независимый JSON.
- [ ] Модели не вычисляют валюту и не требуют готовой сущности Currency или HTTP-операций.
- [ ] Документация различает wire-контракт и серверный пересчёт при rate:null.
- [ ] Существующие Price/Field/Reference/DateTime не меняются; expand отсутствует.
- [ ] Пройдены gofmt, полный race/vet/diff-check и независимые Standards/Spec review
  относительно BASE, зафиксированного до правок.
- [ ] После разрешённой реализации — локальный commit, Status: resolved и ## Answer;
  исходное планирование и чужие изменения сохраняются отдельно.

## Blocked by

Нет. После утверждения DRAFT ticket можно выполнять независимо от positions,
общих накладных расходов и всех пакетов документов.

## Testing decisions

Публичные JSON-seams по образцу общих Address/Field, оба кодека;
ожидаемый JSON задаётся независимо от сериализации SDK. Полный контракт — в spec.

## Утверждение и ownership первой волны

Пользователь 2026-11-01 утвердил первую волну после исправления денежных типов и смысла nil курса.
JSON и HTTP seams этого ticket согласованы. Локальная реализация через TDD и commit разрешены;
независимые Standards/Spec review координирует исполнитель интеграции относительно planning BASE.
Не менять coverage.md/README и общие карты: документацию интегратор сведёт отдельным коммитом.
Файлы и трёхпоточный запуск закреплены в [handoff](../../document-workstream/wave-1.md).
