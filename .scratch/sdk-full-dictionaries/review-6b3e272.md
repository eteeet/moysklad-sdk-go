# Review 30e5d76..6b3e272 — 2026-11-03

Независимые Standards и Spec reviewers (субагенты), только чтение.

Standards: P1 нет. P2 — status-backlog и docs counterparty/organization устарели (закрыто
docs-коммитом 20667cc). P3 — godoc у Get/CreateMany/UpdateMany/DeleteMany internal/nested, ссылка на
ticket в коде store, мёртвые имена ModRequisites* в тесте, «копейки» у demandSumCorrection,
store не использует documenttest.HTTP (покрытие есть через counterparty/organization) — исправлено в
3af1c7b, кроме последнего (оставлено: поведение покрыто).

Spec: P1 — у Organization не было mod__requisites__uz/kz (поля есть в таблице с экранированными `\_`) —
исправлено в 3af1c7b. P2 — contract name обязателен по списку обязательных полей; rewardPercent — Int;
godoc agent контактного лица — исправлено в 3af1c7b; устаревшие docs — 20667cc. P3 — формулировки godoc
organizationbranch и Discount — исправлено; два типа Account — записано в Answer S12.
