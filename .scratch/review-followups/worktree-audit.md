Status: resolved

# Worktree audit 2026-11-01

Основной worktree `/Users/av/ETEEET/moysklad/moysklad-sdk-go` не затронут.
Для всех восьми вспомогательных worktree выполнены status --porcelain=v1
--untracked-files=all, git cherry dev <branch> и просмотр состава изменений.
Все branch-коммиты уже находятся в dev либо имеют patch-equivalent коммиты
(git cherry показывает только «-», ни одного «+»). Branches не удаляются.

## Удалены без --force

| Worktree | Branch | Проверка |
| --- | --- | --- |
| /private/tmp/moysklad-wave1.t9Nefz/model | feat/document-model-foundation-20261101 | clean, HEAD уже в dev |
| /private/tmp/moysklad-wave1.t9Nefz/positions-json | feat/positions-json-20261101 | clean, оба commit patch-equivalent |
| /private/tmp/moysklad-wave1.t9Nefz/positions-reader | feat/positions-reader-20261101 | clean, commit patch-equivalent |

Для каждого выполнен git worktree remove по точному пути. Branches сохранены;
восстановление checkout возможно через git worktree add <новый путь> <branch>.
После удаления git worktree prune --dry-run --verbose не нашёл устаревших записей,
поэтому реальный prune не нужен и не запускался.

## Оставлены до выбора пользователя

| Worktree | Branch |
| --- | --- |
| /Users/av/ETEEET/moysklad/.worktrees/core-counterparty | feat/core-counterparty |
| /Users/av/ETEEET/moysklad/.worktrees/core-country | feat/core-country |
| /Users/av/ETEEET/moysklad/.worktrees/core-expenseitem | feat/core-expenseitem |
| /Users/av/ETEEET/moysklad/.worktrees/core-organization | feat/core-organization |
| /Users/av/ETEEET/moysklad/.worktrees/core-requisites | feat/core-requisites |

Во всех пяти modified docs/coverage.md и docs/temp_plan.md; untracked
.scratch/core-dictionaries/{handoff.md,spec.md,issues/01-country.md,
issues/02-expenseitem.md,issues/03-shared-requisites.md,issues/04-counterparty.md,
issues/05-organization.md}. Их содержимое не сбрасывалось и не удалялось.

Пользователю задан выбор: сохранить каждый worktree через git stash -u
с отдельной меткой и затем удалить clean worktree либо оставить пять worktree
до отдельного разбора. Ответ ещё не получен. Без ответа dirty worktree не удаляются.

## Решение 2026-11-01

Пользователь выбрал stash -u и удаление. Все пять worktree сохранены в stash
(метки `worktree-cleanup core-<name> 2026-11-01`) и удалены без --force; branches остались.
