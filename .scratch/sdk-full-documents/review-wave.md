# Review волны документов — 2026-11-03

BASE: `d21d84d`. HEAD до исправлений: `b02da9d`. Scope: O01–O06, O08–O20 (O07 — needs-info).
Метод: независимые Standards и Spec reviewers (subagents), проверка root, make check CHECK_P=1.

## Standards (P1 нет)

- P2 Устаревшие docs: documents.md, coverage.md, write-async.md, positions.md утверждали, что корзины,
  операций строк и шаблонов нет.
- P2 wizard: VatEnabled/VatIncluded/EvaluateVat против VAT* остального SDK.
- P2 Деньги 0/0.0/дробь не проверены у facturein, retaildrawercashin/out, prepayment(return), payroll,
  pricelist, wizard.
- P3 ID строк payroll и commissionreportin; дубль DeleteManyReturnPositions; пример основания в godoc New
  (7 пакетов); DeletePositions против DeleteManyPositions; rows() у pricelist; счётчики status-technical;
  links_test проверял Create только Contains.

## Spec (P1 нет)

- P2 ID строк payroll и commissionreportin (Update пересоздаёт строки).
- P2 inventory (массив) и pricelist ({"rows"}) по-разному решают одно противоречие источника.
- P2 sum у counterpartyadjustment и retaildrawercashin/out не задать: маркер readonly против примеров.
- P2 write-async.md: корзина/позиции «нет», счётчик «38 сервисов» без новых.
- P3 status-technical счётчики; payments против «Привязки платежей»; pricelist.md устарел; единица
  percentageDiscount; пример payroll с 0.0; допущение копеек; лимит строк в CreateMany
  commissionreportin; имя файла в комментарии prepaymentreturn; форма enters/losses.

## Answer

Исправлено в `2e9e03d`, `3b62b18`, `341a4fd`:

- ID в PositionRequest payroll и обоих наборах commissionreportin; тесты Update с существующей строкой.
- nested.DeleteRows — общий POST …/delete с пустым ответом; positions.Service.DeleteManyPositions
  (переименован), commissionreportin возврат комиссионеру; pricelist rows()/collection(id).
- inventory: positions — {"rows": […]} по всем примерам, как у pricelist; golden и docs.
- wizard: VATEnabled, VATIncluded, EvaluateVAT.
- godoc New: пример основания из списка документа; у supply пример убран (основания нет).
- links_test: значение связи в Create сравнивается целиком.
- Тесты денег 0/0.0/дробь обоими кодеками в 8 пакетах.
- Docs: documents.md, coverage.md, write-async.md (49 сервисов, prepayment(return) без массовых),
  positions.md, status-technical (Trash 31, New 22), pricelist.md, payroll.md, inventory.md,
  commissionreportin.md, documents-extra (допущения: payments, копейки). Лимит строк элементов
  CreateMany/UpdateMany проверяет сервер — записано в documents.md (как у прежних документов).

Не исправлено — решение пользователя: запись sum (counterpartyadjustment, retaildrawercashin/out),
positions/attributes/shared (internalorder, commissionreportout), priceType (pricelist) — маркер
readonly против примеров; список в docs/documents-extra.md.
Вне волны: VatEnabled в productfolder и DefaultCompanyVatEnabled в organization — прежний код.

## Verification

- make check CHECK_P=1 на `341a4fd` — PASS: gofmt, go vet ./... ./testdata, go test -race -count=1 -p 1 ./...,
  diff-check. Предыдущий прогон на `b02da9d` падал только сетевыми флейками (connect timed out,
  7 пакетов) — при повторе зелёные. Live API не вызывался.

## Решение пользователя: readonly против примеров

Пользователь разрешил все четыре группы как конкретные исключения; для остальных полей правило
«readonly важнее примеров» сохраняется. Поле пишется там, где его передаёт пример источника:
internalorder `positions` (Create/Update), commissionreportout `positions`, `attributes`
(Create/Update) и `shared` (только Create), counterpartyadjustment и retaildrawercashin/out `sum`
(Create/Update), pricelist `priceType` (только Create). Строки в теле: лимит 1000 до запроса,
`ID` у строки commissionreportout. Godoc, тесты (readonly-списки, golden, 1001), docs и CHANGELOG
обновлены.
