# Универсальная форма positions: реализация и review

Дата: 2026-11-01. BASE f00ddd9b59007313c746fc1bf13a9860f823feef.
Commits: 1a9d400 code, 7d43628 public docs, eba6f75 assertion label.
Итоговый HEAD eba6f75154929790c5b9d2cb1454dd80783b14db.
Diff: git diff f00ddd9...eba6f75 плюс рабочий planning/spec/tickets.

## Реализация и проверки

- Position/PositionRequest вместо Commercial и Movement, старых aliases нет.
- В read общий Overhead float64 и slots *Reference omitzero; в request слоты Field,
  Overhead отсутствует. Core/Item, reader T any и transport не менялись.
- Четыре movement*.go удалены, regression tests перенесены в общие model/request tests.
  Both codecs: overhead 0/0.0/дробь, missing/null/present slots, marshal nil/one/both,
  full flat fields, unknown, request zero/null/unset/sparse и consumer []/nil/nonempty.
- Red→green отдельно для Position/PositionRequest и добавления полей в общий тип.
- go test -race -count=1 ./... — PASS.
- go vet ./... ./testdata — PASS.
- gofmt и git diff --check — PASS.
- После косметической правки подписи assertion полный race/vet/diff-check повторён:
  PASS; оба независимых reviewer подтвердили 0 открытых для итогового HEAD eba6f75.
- В начале full-прогон был остановлен лимитом approval review, не исполнялся;
  после команды продолжения повторён штатно с разрешением и завершён успешно.

## Standards

Независимый read-only reviewer unify_review_standards: 0.
Production/test diff и рабочие документы соответствуют стандартам. Field/Reference,
публичные seams и независимые JSON-ожидания сохранены. Существенных smells нет.
Переданный хвост planning purchasereturn исправлен; повторная проверка подтверждена.

## Spec

Независимый read-only reviewer unify_review_spec: первоначально один P2 в planning
purchasereturn. Старые Open Decisions требовали отдельные формы и исключение Discount,
а readonly-критерии смешивали заголовок с общей строкой. Исправлено: форма утверждена,
у строки нет Meta/AccountID/Overhead, допустимость commercial Field проверяет сервер;
ограничения заголовка/основания документируются отдельно. Повторный Spec: 0 открытых.
Остальной код и planning соответствуют scope, scope creep не обнаружен.

Итого: Standards 0, Spec 0 открытых; один P2 исправлен до завершения.

## Границы и handoff

P03 wontfix/superseded, M01 от него не зависит. Обновлены действующие move/return/stock/
invoice spec/tickets, map/approval и temp_plan. P01 с прежним контрактом явно исторический;
исследования и старые review/handoff не переписываются. Cost/Shipped/InTransit — будущие
расширения общего positions в соответствующих entity tickets, не новые entity-owned формы.
Entity HTTP/expand/live API/.env/push/PR/VPS не выполнялись. Пять dirty worktree не затронуты.
Planning/tracker/temp_plan обновлены локально, не включены в code/docs commits без отдельной команды.
