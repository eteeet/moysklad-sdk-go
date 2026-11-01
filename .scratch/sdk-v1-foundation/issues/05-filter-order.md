Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/01-public-interface.md

# T05: Filter grammar и order

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Builders для documented filter operators, escape literal semicolon и одного URL-encoding;
Reference href, UUID, Attribute href, null/empty, LIKE, archived/isDeleted. Order по нескольким
полям с asc/desc. Open field names и raw escape hatch; endpoint capabilities не заменяются
универсальным enum.

## Специальные проверки

Повторы = дают OR, != дают AND, разные поля AND; unsupported комбинации не переосмыслять как
произвольный AST. Unicode, %, &, +, ;, backslash, nested href; report-specific grammar отдельно.
Ни hidden lookup для фильтров, ни double encoding.

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

BASE `53ffa95`, code commit `a220c7c`.

- Пакет `query`: Eq/NotEq/Less/LessOrEq/Greater/GreaterOrEq, Contains/StartsWith/EndsWith (~, ~=, =~),
  IsEmpty/NotEmpty, RawCondition; `Filter(...)` и `Order(Asc/Desc...)`. Поля открытые: имя, href доп.
  поля; значения — строка, UUID, href. Скрытых lookup нет.
- По _general.md: повтор = → OR, != → AND, разные поля → AND; отклоняются = с неравенством того же поля
  (пример sum=100;sum>99) и повтор <, <=, >, >= (будущая ошибка API). Повтор ~ разрешён — источник
  его не ограничивает. ';' в значении → `\;`; значение с `\` в конце отклоняется. Raw допускает `\;`.
- URL-кодирование один раз через url.Values: точный RawQuery в httptest (Unicode, %, &, +, ;, \, #, href).
- Пустые заглушки transport/{filter,order,query,pagination}.go удалены: место query — пакет query.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review Spec: исправлены правило комбинаций, хвостовой backslash, escape в raw, лишний запрет
  повторов ~ и полей order; добавлены isDeleted/одно поле order. Review Standards: префиксы в тестах,
  точный RawQuery, сообщения ошибок.
