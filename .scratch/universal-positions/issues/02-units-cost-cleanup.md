Status: resolved
Stage: IMPLEMENTED
Blocked by: нет

# 02: Discount units, Cost planning и дублирующий тест

Пользователь разрешил только три согласованные темы. BASE eba6f75.

- [x] Coverage и guide: Price/Overhead — деньги в копейках, Discount — процент
  скидки/наценки; отрицательное значение задаёт наценку. Тип float64 не меняется.
- [x] S02 запланирует Cost float64 в общей Position и Cost Field[float64] с omitzero
  в общей PositionRequest, для услуг. Согласовать это в stock-documents/spec и
  document-positions/spec, закрепить JSON/HTTP cases будущего S02 и первичные источники.
  Cost сейчас в код не добавляется; применимость проверит сервер без скрытого Get.
- [x] Удалить только повторный absent fields из TestPositionOptionalFieldsJSON;
  оставить этот case в TestPositionJSON и все уникальные проверки.
- [x] Локальные race/vet/gofmt/diff-check, независимые Standards/Spec review.

Production Go, расширения SDK, cleanup worktree и коммит planning вне scope.
Публичные docs и тест сохраняются локальным коммитом по действующему workflow.

## Answer

2026-11-01: публичные docs и удаление одного дублирующего case — 59b16e7.
Cost read/write планируется S02; текущая задача production Go не меняла.
Уточнены stock-documents/spec, S02 ticket, document-positions/spec и общий contract
universal-positions/spec; добавлены будущие JSON/HTTP cases и pinned primary sources.

Own delta (eba6f75 + удаление одного case) прошла полный race/vet/gofmt в изолированном
snapshot /private/tmp/moysklad-doc-fixes.YNlCMm. Diff-check PASS. Independent Standards 0,
Spec 0. [Отчёт](../review-units-cost.md).

Важно: общий рабочий каталог параллельно изменяется внешним процессом: добавлены
document packages и Cost/Shipped/InTransit в positions. Полный main-прогон упал на
JSON golden этих новых полей; текущая задача их не меняла и не исправляла.
Planning не коммитился; worktree cleanup не выполнялся.
