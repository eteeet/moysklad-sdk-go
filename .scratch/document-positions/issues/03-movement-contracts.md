Status: needs-triage
Stage: DRAFT
Blocked by: 01-json-contracts.md

# 03: Самостоятельный JSON-контракт Movement

## What to build

Разработчик решения читает и записывает согласованные поля позиции Movement через
отдельные public JSON types общего positions, без коммерческих discount/vat.
Проверяемый результат — сериализация реального consumer container; HTTP заголовка move
и reader в этот JSON-срез не включаются.

Все имена и поля PROPOSED: [spec](../spec.md). Реализация — после утверждения DRAFT,
перевода ticket в ready-for-agent и завершения 01.

## Acceptance criteria

- [ ] Movement встраивает Item из 01 и содержит Price float64, Overhead float64,
  SourceSlot/TargetSlot *model.Reference. JSON плоский; денежные значения в копейках.
  API Int overhead отличен от document-level Overhead.Sum (API Float).
- [ ] MovementRequest встраивает ItemRequest и содержит Price Field[float64],
  SourceSlot/TargetSlot Field[model.Reference] с omitzero.
- [ ] Request не содержит Meta/AccountID/Overhead/Discount/VAT/VATEnabled.
  ID наследуется только как matching selector existing-row embedded Update;
  Create новой строки не обещает назначение client-generated ID.
- [ ] Оба JSON codecs проверяют значения всех read полей, дробные quantity/price,
  overhead, optional refs absent/present и неизвестное поле; без обязательного expand.
  Отдельные cases проверяют overhead в формах 0, 0.0 и дробного числа обоими кодеками;
  точность integer >2^53 для денежных float64 не обещается.
- [ ] Независимые golden JSON обоими codecs проверяют все writable поля, минимальную
  новую строку, sparse matching request, явные 0 и null, omission optional slots и
  отсутствие readonly/коммерческих полей.
- [ ] Consumer container с Positions Field[[]MovementRequest] различает omission,
  empty/nil [] и непустой массив; full replacement не превращается в merge.
  Кодирование null не объявляется поддерживаемой очисткой slots или positions.
- [ ] Документация объясняет отдельный flavour без скидок/НДС и расширяемость
  transport через T any. Pack/Things/прочие movement поля не считаются покрытыми.
- [ ] Тестовая точка — public JSON model/requests и container, по образцу UOM codecs
  и самостоятельных общих типов; private model structure не закрепляется.
- [ ] Existing Core/Item, Field/Reference и transport не рефакторятся ради Movement;
  полная реализация entity move не объявлена готовой.
- [ ] После реализации пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check и отдельные Standards/Spec review относительно зафиксированного BASE.

## Blocked by

[01: Общий JSON-контракт](01-json-contracts.md) — Movement/MovementRequest используют
Item/ItemRequest. Reader 02, DocumentRate, document-common и move HTTP-сервис не блокируют
этот самостоятельный JSON-контракт. Entity move позже использует готовый 03.

## Границы

После утверждения — один основной implementer и законченный JSON-срез.
Сейчас только DRAFT. Не добавлять entity move HTTP, mutations позиций, auto-pagination,
expand, pack/serial/tracking, live API, .env, push или commits подготовки.
Новые решения контракта согласовать отдельно.

## Sources

- [Movement positions: отдельные поля и overhead read-only](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_move.md#L62-L89).
- [Общие правила embedded positions](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1031-L1037).

## Comments
