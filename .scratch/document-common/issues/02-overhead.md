Status: ready-for-agent
Stage: APPROVED
Blocked by: None

# 02: Публичные DocumentOverhead и DocumentOverheadRequest

## What to build

Разработчик читает и формирует накладные расходы документа общими типами,
без собственной копии модели в Supply, Demand или Move.
Срез завершён публичным JSON-контрактом, тестами и документацией.
Утверждённый контракт — [spec](../spec.md). Разрешена реализация первой волны.

## Acceptance criteria

- [ ] DocumentOverhead содержит Sum float64 и Distribution OverheadDistribution.
- [ ] OverheadDistribution — открытый именованный string с точными константами weight/volume/price;
  неизвестная строка сохраняется, default и локальная enum-валидация не добавляются.
- [ ] Request содержит Sum Field[float64] и Distribution Field[OverheadDistribution] с omitzero.
- [ ] Внешний Field[DocumentOverheadRequest] и вложенные Field сохраняют unset/null/{}/частичный
  объект, явный 0, дробное число, пустую/неизвестную строку.
- [ ] Оба JSON-кодека проверяют все поля и константы через независимые ожидания;
  суммы остаются в API-единицах, без округления и распределения на клиенте.
- [ ] Read-only overhead отдельной позиции не подменяется этим объектом.
- [ ] Документация объясняет серверную ответственность за распределение и допустимость null.
- [ ] Нет HTTP, expand, изменения существующих общих типов или новой test infrastructure.
- [ ] Пройдены gofmt, полный race/vet/diff-check и независимые Standards/Spec review.
- [ ] После разрешённой реализации — локальный commit, Status: resolved и ## Answer;
  предшествующие изменения планирования не включены в commit реализации.

## Blocked by

Нет. Не зависит от DocumentRate, positions или готовой сущности документа.

## Testing decisions

Публичный JSON-контракт, оба кодека, независимые значения всех констант и нулевой суммы.
См. полный [контракт](../spec.md) и образцы существующих общих типов SDK.

## Утверждение и ownership первой волны

Пользователь 2026-11-01 утвердил первую волну после исправления денежных типов и смысла nil курса.
JSON и HTTP seams этого ticket согласованы. Локальная реализация через TDD и commit разрешены;
независимые Standards/Spec review координирует исполнитель интеграции относительно planning BASE.
Не менять coverage.md/README и общие карты: документацию интегратор сведёт отдельным коммитом.
Файлы и трёхпоточный запуск закреплены в [handoff](../../document-workstream/wave-1.md).
