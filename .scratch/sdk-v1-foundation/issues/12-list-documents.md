Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md

# T12: List одиннадцати существующих документов

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Typed List и request options для PaymentIn/Out, CashIn/Out, InvoiceIn/Out,
Supply/Demand/PurchaseOrder/Move/PurchaseReturn. Фильтры archived/isDeleted там, где поддержаны.
Positions в List/Get — только фактический ответ, без expand/GetPositions follow-up.

## Специальные проверки

Все 11 exact resources, collection/meta/context, query/headers/errors и no hidden requests.
Незаданные query/headers прежних вызовов не получают неоговорённых defaults. Старые
Get/Create/Update tests сохранены.

## Acceptance criteria

- [ ] Scope/публичный interface и field/method matrix утверждены; зависимости resolved
  и их код присутствует в рабочем дереве. Readonly, required response/Create, null,
  immutable и выбранные TODO разделены по источникам; противоречия не замолчаны.
- [ ] Реализован весь What to build, без неоговорённых методов и hidden requests.
- [ ] Специальные проверки этого ticket реализованы через публичный interface.
- [ ] JSON обоими codecs и независимые golden; upstream/provenance не изменены.
  Для transport/iterator/CI tickets вместо entity JSON применяются их собственные
  публичные seams и контрактные tests, а не фиктивные модели.
- [ ] Реальные transport/httptest tests покрывают успех и ошибки; context и
  metadata сохраняются, caller-controlled retry/walker не подменяет один вызов.
- [ ] go test -race, go vet ./... ./testdata, gofmt и diff-check PASS; независимые
  Standards/Spec review, исправления и финальный BASE/результат записаны в Answer.
- [ ] README/guide/status отражают реальный scope, ограничения и неизвестные факты.
  Planning без команды не коммитится; один code commit по ticket.

## Источники

[транспортный research](../../sdk-v1/research/transport.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `9acae45`, code commit `71b3390`.

- `List` у PaymentIn/Out, CashIn/Out, InvoiceIn/Out, Supply, Demand, PurchaseOrder, Move,
  PurchaseReturn через `internal/list.Entity`; позиции — только фактический ответ, без expand и
  GetPositions. isDeleted=true;isDeleted=false отбирает и документы из корзины.
- Решение: fields=stock в List документов не принимается. По _general.md остатки лежат в positions и
  в списке приходят только с expand=positions (и limit ≤ 100); expand вне v1. Строки с остатками —
  GetPositions (fields=stock там — D01). search у PurchaseReturn не подтверждён таблицей параметров,
  но общий раздел контекстного поиска применим ко всем спискам — передаётся как есть.
- Тесты (`internal/listtest`, 17 пакетов): точные method/path/query/headers, plain/gzip, пустой массив,
  null/malformed/голый массив/204, APIError с ResponseInfo, ровно один запрос; неверные опции и
  неподдержанный fields — ноль запросов; локальные ошибки (nil/отменённый ctx, Auth, транспорт) через
  documenttest.LocalErrors. Строка — upstream-фикстура сущности.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS; старые
  Get/Create/Update тесты и golden не менялись.
