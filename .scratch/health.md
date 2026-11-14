# Здоровье документов

Пересборка 2026-11-11, `dev` @ `3fc370b`; сверка 2026-11-14, `dev` @ `159992b`: изменились `GIT_FLOW.md`, `INFRA.md`, `~/ETEEET/moysklad/AGENTS.md` и добавлен тикет `trash-content-type/01`, код не менялся. Только действующая цепочка; история (закрытые спеки и тикеты, research, ревью, handoff, снимки) не сверяется.
Цепочка: 1 md → 2 ADR и glossary → 3 открытые спеки → 4 открытые тикеты и бэклог → 5 docs и README → 6 код.
Строка: `файл:строка — с каким звеном расходится — предложение — кто решает`. «Механика» — однозначная правка вниз; «владелец» — выбор или правка верхнего звена.

## 1. md

Находок нет.

## 2. ADR и glossary

Находок нет.

## 3. Открытые спеки

Находок нет.

## 4. Открытые тикеты и бэклог

Находок нет.

## 5. docs и README

Находок нет.

## 6. Код

- K3. `origin/main` = `5f3fb93 Initial commit`, последний выпуск `v0.2.2` — с `GIT_FLOW.md:20` (владелец перематывает `main` на тег) — `git push origin 'v0.2.2^{}:main'` — владелец.

## Открытые тикеты и спеки

- `.scratch/document-name-autonumber/issues/01-create-without-name.md` (ready-for-agent) — оставить: заведён по K1, сверен с `~/ETEEET/moysklad/AGENTS.md:16,18`; включает `docs/documents.md:16-18` (бывший фрагмент D9).
- `.scratch/ci/issues/01-github-actions.md` (ready-for-agent) — оставить: заведён по T3 после «да» владельца на смену способа CI, сверен с `INFRA.md:17-24,31` и `GIT_FLOW.md:13,19`.
- `.scratch/adr-0002-code-refs/issues/01-assortment-comment.md` (ready-for-agent) — оставить: заведён по A5, сверен, согласован.
- `.scratch/trash-content-type/issues/01-trash-content-type.md` (ready-for-agent) — оставить: работа не сделана; выпуск исправления ждёт `ci/01`.
- `.scratch/sdk-full-documents/issues/07-cash-accounts.md` (ready-for-agent) — оставить: работа не сделана.
- `.scratch/sdk-full-documents/spec.md` (needs-triage) — оставить: открыт O07.
- `.scratch/sdk-full/map.md` (needs-triage) — оставить: в карте открыт O07.

## Покрытие

- md: `~/ETEEET/AGENTS.md`, `GIT_FLOW.md`, `INFRA.md`, `~/ETEEET/moysklad/AGENTS.md`, `AGENTS.md`, `docs/agents/domain.md`, `docs/agents/triage-labels.md` — согласованы; `docs/agents/issue-tracker.md` — согласован.
- ADR: `docs/adr/0002-baseline.md` — согласован; `0001` заменён им. Glossary: общий — docs согласованы с ним; своего `GLOSSARY.md` нет.
- Открытые спеки: `.scratch/sdk-full-documents/spec.md`, `.scratch/sdk-full/map.md` — согласованы.
- Открытые тикеты: `07-cash-accounts`, `trash-content-type/01` и новые тикеты — согласованы. Бэклог: `docs/status-backlog.md` — согласован.
- `README.md`, `docs/*.md`, `testdata/readme.md` — согласованы.

Не проверено: внешние ссылки на документацию МойСклад, поведение живого API (O07; автонумерация — в `document-name-autonumber/01`), ссылки на API над пакетами в коде (`~/ETEEET/AGENTS.md`, «Комментарии в коде»).
