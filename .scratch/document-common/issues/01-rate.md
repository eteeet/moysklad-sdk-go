Status: resolved
Stage: IMPLEMENTED
Blocked by: None

# 01: Публичные DocumentRate и DocumentRateRequest

## What to build

Разработчик читает и формирует вложенный rate любого документа через общие типы
model.DocumentRate/DocumentRateRequest. Срез включает публичный JSON-контракт,
оба JSON-кодека, независимые ожидаемые тела и документацию.
Утверждённый контракт — [spec](../spec.md). Разрешена реализация первой волны.

## Acceptance criteria

- [x] DocumentRate.Currency — *Reference, Value — *float64; nil Value внутри присутствующего rate означает курс 1, но SDK не подставляет 1 в JSON; над полем краткий комментарий об этом.
- [x] DocumentRateRequest содержит Currency Field[Reference] и Value Field[float64] с omitzero.
- [x] Внешний Field[DocumentRateRequest] различает unset/null/{}/значение; внутренние поля
  сохраняют null, 0, 1, дробный курс и пропуски.
- [x] Оба кодека проверяют значения вложенной Meta, optional/null/unknown поля и независимый JSON.
- [x] Модели не вычисляют валюту и не требуют готовой сущности Currency или HTTP-операций.
- [x] Документация различает wire-контракт и серверный пересчёт при rate:null.
- [x] Существующие Price/Field/Reference/DateTime не меняются; expand отсутствует.
- [x] Пройдены gofmt, полный race/vet/diff-check и независимые Standards/Spec review
  относительно BASE, зафиксированного до правок.
- [x] После разрешённой реализации — локальный commit, Status: resolved и ## Answer;
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

## Answer

Реализованы DocumentRate/DocumentRateRequest в model/document_rate.go и публичные
JSON-тесты в model/document_rate_test.go. Оба кодека проверяют полную Meta валюты,
0/1/целый/дробный курс, пропущенные/null/unknown поля, внешний rate и вложенные Field.
Независимые JSON literals подтверждают отсутствие подстановки курса и точные request-состояния.
GoDoc описывает nil value и серверный пересчёт при внешнем rate:null.

TDD red: `go test ./model -run '^TestDocumentRateJSON$' -count=1` упал до реализации
с `undefined: model.DocumentRate`; отдельный request-срез упал с
`undefined: model.DocumentRateRequest`. После минимальной реализации оба среза green.

Проверки 2026-11-01: targeted `go test ./model -run '^TestDocumentRate' -count=1`,
`go test -race -count=1 ./...`, `go vet ./... ./testdata`, gofmt и `git diff --check` прошли.
Независимые Standards/Spec review координирует интегратор относительно BASE
`a2804283e7aafe3c6ee41c17a4c2c79f7e01f999`; на этапе implementer Status сохранялся claimed до review.

### Итог интеграции

Включено в dev, commit реализации: 86117897dc7a5beeaa73e29f46ccca738eec2f2d.
Независимые Standards: 0 замечаний; Spec: 0 замечаний. Проверен полный срез
от planning BASE a2804283e7aafe3c6ee41c17a4c2c79f7e01f999 до HEAD 1b28603,
включая централизованные docs. [Отчёты](../../document-workstream/review-wave-1.md).
Общий go test -race -count=1 ./..., go vet ./... ./testdata, gofmt и diff-check пройдены.
Публичная документация сведена интегратором после слияния; live API и push не выполнялись.
