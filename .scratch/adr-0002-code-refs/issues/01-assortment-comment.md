Status: ready-for-agent
Type: task
Blocked by: —

# 01: комментарий assortment ссылается на заменённый ADR 0001

ADR 0001 заменён `docs/adr/0002-baseline.md`, и исключение из порядка заголовков для beta-тега assortment теперь записано там. Комментарий над `List` в `entity/assortment/service.go:54-56` по-прежнему называет его «исключение из порядка ADR 0001».

## What to build

В комментарии заменить ссылку на ADR 0001 ссылкой на ADR 0002. Поведение кода не меняется.

## Acceptance criteria

- [ ] В `entity/assortment/service.go` нет упоминаний ADR 0001; комментарий ссылается на ADR 0002.
- [ ] `grep -rn 'ADR 0001' --include='*.go' .` ничего не находит.
- [ ] `make check` зелёный.
