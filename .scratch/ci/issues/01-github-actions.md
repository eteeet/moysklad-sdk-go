Status: ready-for-agent
Type: task
Blocked by: —

# 01: CI на GitHub Actions по INFRA.md

В репозитории нет `.github/workflows/`, хотя `~/ETEEET/INFRA.md:17-24` требует CI у каждого репозитория с кодом. Без CI не выполняются `GIT_FLOW.md:13` (ветка сливается в `dev`, когда CI ветки зелёный) и `GIT_FLOW.md:19` (тег ставится на коммит с зелёным CI). Владелец дал «да» на смену способа CI 2026-11-14 (`INFRA.md:31`). Выпуск исправления `.scratch/trash-content-type/issues/01-trash-content-type.md` ждёт этого тикета.

## What to build

Workflow `.github/workflows/ci.yml`, который запускает `make check`:

- на push любой ветки, кроме `dev` и `main`, и на pull request;
- на push в `dev`.

SDK образ не собирает, поэтому тег `v*` и push в `main` ничего не запускают. Go берётся из `go.mod` (`actions/setup-go`, `go-version-file`). Права — `contents: read`. `make check` не требует сети и токенов; `.env` нужен только примерам (`make example*`) и в CI не используется. Образец — `ms_documan/.github/workflows/ci.yml`, без Postgres и без сборки образа.

## Acceptance criteria

- [ ] `.github/workflows/ci.yml` запускается на push ветки, на PR и на push в `dev`; на `main` и тегах не запускается.
- [ ] Прогон на ветке этого тикета зелёный: `gh run list --branch <ветка>` показывает успешный `make check`.
- [ ] Строки «CI для v1 не делается» в `docs/status-backlog.md`, `docs/status-technical.md` и этап 17 в `docs/coverage.md` отмечены сделанными.
- [ ] В отчёте отдельной строкой — изменение способа CI (`INFRA.md:31`).
