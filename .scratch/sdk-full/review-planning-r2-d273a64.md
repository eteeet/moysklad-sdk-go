# Повторное review planning SDK — 2026-11-03

HEAD: `30e5d769fb27ec7a19ffe47394ffd796120aff0b` (без новых commits).
Scope: untracked planning, 98 Markdown-файлов и 91 ticket; review-* исключены.
Snapshot SHA-256 (sorted path + NUL + bytes + NUL):
`d273a64aca7d7f4bad771b83799105db378570e75e4a1c555007bf0e82b744a2`.
Метод: независимые Standards/Spec reviewers и дополнительная проверка root.
Предыдущий отчёт: [review-planning-30e5d76.md](review-planning-30e5d76.md).

## Standards

### S1. P2 — синхронизировать scope polling печати

`sdk-full-reports/spec.md:310` оставляет цепочку 202 + Location → опрос → 303
и реализацию по образцу async Wait. Исправленный I12 в spec.md:187 и
issues/12-print.md:18 прямо устанавливает: протокол опроса неизвестен, вернуть
различимые результаты 202/303 без Wait. Старый workflow также сохранился в
docs/status-backlog.md:137. Исполнитель получает несовместимые требования.

Нужно привести Open Decisions и backlog к исправленному I12. Это оставшаяся
часть P7 предыдущего review.

Прежние Standards S1/S2 закрыты: новые документные tickets включают Trash и
зависят от O04; S10 явно владеет barcodes Consignment и зависит от G01.
Новые S26/S27/I23 согласованы с основными разделами spec и таблицами.

## Spec

### P1. P2 — P7 закрыт частично: обещание polling осталось в Open Decisions

Независимо найдено то же противоречие в sdk-full-reports/spec.md:310 и
docs/status-backlog.md:137. Источники _print.md:8–9 и _labels.md:8–9 называют
вариант 202 будущим, протокол polling не описан. Новый I12 корректен; старые
обещания нужно удалить либо обозначить неизвестным будущим контрактом.

Остальные прежние замечания закрыты:

- S16: POST коллекции accounts, сохранение непереданных полей, неизменяемость валюты.
- G13: Update/UpdateMany; S01: запись parent/pathName.
- G01: отдельные формы Product/Variant packs и существующая Barcode; packs
  Service/Bundle не добавляются.
- G03: storebalances; I23: async разрешённых отчётов с исключением выборочного POST.
- G04/S19/I19: потерянные операции images/components/notes добавлены.
- I04: plain array и строковые UUID byslot/current.
- S26/S27: Organization fields и массовые операции states/characteristics.
- Root подтвердил O02: enter добавлен с оговоркой общего источника и 500 позиций;
  O16: read-only GetPosition предоплаты добавлен.

Новых подтверждённых расхождений в S26/S27/I23 не найдено.

## Verification

- 98 planning files / 91 ticket: local links, dependency targets, отсутствие циклов,
  статусы и таблицы зависимостей map/spec/issues проверены Standards reviewer.
- Snapshot SHA-256 до и после проверки совпал.
- Go tests не запускались: рассматривается только planning, Go code не менялся.
- HEAD оставался 30e5d76; spec, tickets, tracker statuses и code reviewer не менял.
- Создан только этот отчёт; .DS_Store не читались и не изменялись, live API не вызывался.

Standards: 1 finding P2. Spec: 1 finding P2. Один уникальный остаток; остальные
13 замечаний предыдущего review закрыты. DRAFT этим review не утверждается.
