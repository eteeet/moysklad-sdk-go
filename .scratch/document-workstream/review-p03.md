# P03: Movement JSON — реализация и review

Дата: 2026-11-01. BASE: 1ee21152e94270149e049c473e34898a58fdb7ee.
Код: 519b7a6 (`feat(positions): add movement JSON contracts`).
[Ticket](../document-positions/issues/03-movement-contracts.md),
[spec](../document-positions/spec.md), [публичный контракт](../../docs/positions.md).

## Реализация и проверки

- Movement встраивает Item; Price/Overhead float64, optional SourceSlot/TargetSlot ссылки.
- MovementRequest встраивает ItemRequest; Fields Price/SourceSlot/TargetSlot с omitzero.
  Overhead и коммерческих полей в request нет. ID — только matching existing-row selector.
- Отдельный red→green каждого типа: тест сначала не компилируется из-за undefined
  Movement/MovementRequest, после реализации проходит обоими JSON codecs.
- Independent literal JSON: все read/write поля, overhead 0/0.0/дробь, unknown fields,
  absent/present/null slots, sparse ID request, explicit zero/null, consumer unset/[]/nil/nonempty.
- go test -race -count=1 ./... — PASS.
- go vet ./... ./testdata — PASS.
- gofmt новых Go-файлов и git diff --check — PASS.

HTTP reader/transport/Core/Item/Field не менялись. Entity Move, Pack/Things, expand,
live API, .env, push/PR/VPS не входили в scope. Остальные 11 entity tickets не запущены.

## Standards

Независимый read-only reviewer: p03_review_standards. Standards 0.
Проверены commit 519b7a6 относительно BASE и рабочая документация/spec/ticket.
Нарушений документированных правил и существенных smells не найдено. Типы следуют
существующим Core/Item и request-формам; комментарии краткие, ссылки закреплены на
revision и раздел. Тесты используют публичный JSON-контракт, оба codecs и независимые
ожидания. Повторяется только принятая структура табличных тестов; дополнительная
абстракция не требуется.

## Spec

Независимый read-only reviewer: p03_review_spec. Spec 0.
Movement/MovementRequest соответствуют полям, embedding и JSON tags. Тесты покрывают
все read-поля, overhead 0/0.0/дробь, optional refs, unknown fields, independent request
JSON и состояния consumer container. Request исключает readonly и коммерческие поля.
Документация отражает float64, отсутствие гарантии очистки через null и Pack/Things
вне покрытия. Scope creep нет. Reviewer тесты не перезапускал; полный прогон выполнен
основным implementer.

Итого: Standards 0, Spec 0; блокирующих замечаний нет.

После review пользователь выявил пропущенный случай: nil SourceSlot/TargetSlot
кодировались как null, а не пропускались по convention SDK. Исправлено bf2eaac
с red→green тестом обоими codecs. Изначальные нулевые отчёты не обнаружили этот пробел.
[Последующие согласованные исправления](../review-followups/spec.md).

## Сохранение документации

Публичные README/coverage/positions guide сохранены отдельным docs-коммитом f26ec32.
Tracker/spec/map и temp_plan обновлены локально; по правилу отдельной команды
на их коммит не было, поэтому они не включаются в implementation/docs commits P03.
