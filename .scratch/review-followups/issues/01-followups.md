Status: needs-info
Stage: APPROVED
Blocked by: нет

# 01: Закрыть согласованные замечания review

[Контракт и критерии](../spec.md). Один implementer; разные code-коммиты
для Movement omission, SyncID и косметического rename. Worktree cleanup после
read-only audit, без удаления branches и без --force.

## Acceptance criteria

- [x] Movement omitzero и обе codecs, nil/одна/обе ссылки.
- [x] SyncID Update обеих сущностей, JSON/HTTP-тесты, без скрытого Get/retry.
- [x] validID rename без изменения поведения.
- [x] Обновлены публичные описания SyncID и закрыт coverage TODO.
- [ ] Вспомогательные worktree удалены безопасно; основной worktree сохранён.
- [x] Полные проверки и независимые Standards/Spec review пройдены.

## Comments

Код: bf2eaac (Movement), 4a37c77 (SyncID), d488f56 (rename).
Публичная документация: f00ddd9. Red→green подтверждён для Movement и каждого
UpdateRequest. Полный race/vet/gofmt/diff-check пройден, независимые Standards 0,
Spec 0 замечаний к коду. [Отчёт](../review.md).
Worktree cleanup: 3 clean временных worktree удалены, branches сохранены;
5 dirty ожидают решения пользователя о сохранении. [Audit](../worktree-audit.md).

## Answer

Кодовые и документационные правки завершены. Ticket needs-info только из-за cleanup:
пять worktree содержат local modified/untracked файлы, поэтому не удалены и не сброшены.
Пользователю предложено stash -u с отдельными метками либо сохранение worktree до
отдельного разбора. Ответ не получен, resolved не установлен.
