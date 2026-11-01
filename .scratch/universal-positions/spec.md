Status: resolved
Stage: IMPLEMENTED

# Одна форма positions для всех товарных документов

Пользователь 2026-11-01 утвердил переход и имена Position/PositionRequest.
BASE: f00ddd9b59007313c746fc1bf13a9860f823feef.

## Контракт и критерии

- Публичные позиции товарных документов имеют одну форму positions.Position и
  PositionRequest. Commercial/CommercialRequest и Movement/MovementRequest удаляются
  без deprecated aliases: SDK не выпущен, совместимость не требуется.
- Position встраивает Item, сохраняет Price/Discount float64, VAT int/VATEnabled bool;
  добавляет read-only Overhead float64 в копейках и optional SourceSlot/TargetSlot
  *model.Reference с omitzero. Нет масштабирования денег или exactness >2^53.
- PositionRequest встраивает ItemRequest и содержит Price/Discount/VAT/VATEnabled и
  SourceSlot/TargetSlot через model.Field с omitzero. Overhead только чтение.
  Unset не отправляется; допустимость заданного поля/значения определяет сервер.
- Нули отсутствующих полей — zero values. У move vat=0/vatEnabled=false означают
  отсутствие поля, а не «без НДС». Read-модель не служит write request.
- Четыре movement*.go файла удаляются; их полезные тесты переносятся на общий тип:
  все поля/unknown/flat JSON, overhead 0/0.0/дробь, optional refs absent/null/present,
  marshal nil/каждая/обе ссылки обоими codecs. Request проверяет совместную запись
  коммерческих полей и slots, omission/zero/null/sparse и отсутствие readonly полей.
  Consumer Field[[]PositionRequest] сохраняет unset/null/[]/nil/nonempty без merge.
- Публичные README/coverage/positions guide не упоминают Movement. В docs фиксируется:
  отдельный тип строки возможен только при реальном конфликте типа или смысла одного
  JSON-поля в разных документах. Различие набора полей само по себе не конфликт.
- Обновить действующие spec/tickets positions, map/approval, stock-move, purchase-return,
  stock-documents и invoice-documents. P03 отменён/заменён, M01 больше от него не зависит.
  PurchaseReturn использует общий PositionRequest; бизнес-ограничения остаются в docs.
  Cost demand и Shipped/InTransit purchaseorder добавятся в общую форму в entity tickets,
  сейчас этой правкой не разрабатываются. S02 добавит Cost float64 в Position и
  Cost Field[float64] с omitzero в PositionRequest, в копейках только для услуг;
  применимость проверяет сервер. InTransit write добавится в общий PositionRequest в S03.
- Generic reader T any и typed entity.GetPositions остаются; entity aliases на общий
  Position/PositionRequest допустимы и не создают отдельную форму. Нет собственных
  entity structs строк только ради дополнительных полей.
- Gofmt, go test -race -count=1 ./..., go vet ./... ./testdata, diff-check и независимые
  Standards/Spec review относительно BASE; код и публичные docs — локальные коммиты.
  Spec/tracker/temp_plan остаются локальными до отдельной команды на planning-коммит.

## Границы

Entity HTTP не реализуется, новые поля Cost/Shipped/InTransit сейчас не разрабатываются.
Core/Item/Field/Reference/transport и алгоритм reader не рефакторятся. Expand, live API,
.env, push/PR/VPS исключены. Пять dirty worktree остаются сохранёнными до отдельного выбора.
