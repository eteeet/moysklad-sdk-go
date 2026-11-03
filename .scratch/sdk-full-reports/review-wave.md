# Review захватки 5 и хвостов G12/G13 — 2026-11-03

BASE: `3778088`. HEAD review: `68ad1b0`. Исправления: `aa5c8ea`.
Scope: I01–I23, G12, G13; 21 коммит, 224 файла. Метод: независимые Standards и Spec reviewers
(Spec разделён: отчёты I01–I10+I23 и печать/интеграции I11–I22+G12/G13), проверка root по
источнику, make check.

## Standards

Жёстких нарушений нет: ResetResponseInfo, один `moysklad:`, memnet, без `shortcut:`, деньги,
запись без повторов; контракт 201 ослаблен только для коллекций с `Created`; same-origin Download
не тронут; SendNoRedirect — флаг в контексте одного вызова.

На суждение:
1. ADR 0001: «Services вызывают только Send», а print вызывает SendNoRedirect. — Исправлено: ADR
   описывает SendNoRedirect как обёртку Send с тем же limiter/retry/ResponseInfo.
2. ResetResponseInfo в ветке ошибки (task comments, eventfeed) — поведение то же, nested сбрасывает
   сам; прецеденты прошлых волн. Не менялось.
3. Duplicated Code: testAsync/testAsyncRejects шестью копиями в report/*; Interval*/Options у money и
   sales; разные имена PlotSeries/Plot. — Не исправлено: «каждый отчёт — собственная форма» и
   «без новых общих абстракций» (spec); тестовые хелперы в internal/documenttest — кандидат на
   отдельную чистку.
4. nested.Collection.Created ветвится на пару экспортируемых функций write. — Не менялось: так
   существующие 8 вызовов CreateManyNested/UpdateManyNested остались без правок.
5. Пакет `print` затеняет builtin print в client.go — builtin там не используется; не менялось.

## Spec

Отчёты I01–I10, I23:
- (a) Устаревшие «async (I23) не входит» в docs report-sales/money/profit/counterparty. — Исправлено.
- (a) «orders/plotseries нет в статусе» — не подтвердилось: строка `orders — plotseries` есть.
- (b) UOM у ByVariant (общая строка с ByProduct, где uom в таблице) и Context/Meta plotseries (общая
  обёртка ответа API, не поля данных) — оставлены; в отличие от inn/email/phone I07 это не поля
  сущности вне таблицы.
- (c) report/sales отвергал ответ без meta.href, хотя обязателен только series. — Исправлено, тест.

Печать и интеграции I11–I22, G12/G13:
- (c) G13 CreatePositions не проверял assortment/trackingCodes («Необходимо при создании»). —
  Исправлено, тесты отказа (мутация роняет оба).
- (c) I15 `method` без маркера readonly не писался. — Исправлено: Field в Create/Update, golden.
  `authorApplication` — метаданные решения, ставит сервер; только чтение.
- Остальное сверено: операции, обязательные поля, пути, 303/202 без скрытого GET, 201 ленты событий.

## Решения пользователя и допущения (не реализовано, ждёт)

- I07: inn/email/phone в показателях контрагентов — нет в таблице, есть в фильтрах, примерах и
  фикстуре официального SDK; сейчас не читаются.
- I22: модель application — по примеру (таблицы атрибутов нет).
- I04: byslot/current читает `stock`, фикстура официального SDK содержит `quantity`.
- I06: единицы money.movement не названы.
- I23: async кратких stock-отчётов, byoperation и Get показателей контрагента — по толкованию
  `_async.md` (раздел целиком), пример есть только для bystore.
- I19: создание событий — 201 по источнику (принимается и 200).
- O07 — needs-info без изменений (нужен ответ МойСклад или live-проверка с согласия).

## Verification

- make check (CHECK_P=2) на `68ad1b0` — PASS; на `aa5c8ea` — PASS.
- Live API не вызывался.
