# Синхронизация фикстур

Из корня проекта:

```sh
make fixtures
make fixtures-commit
```

По умолчанию скачивается текущий `HEAD` официального репозитория.
Для конкретной версии:

```sh
make fixtures FIXTURES_REF=77a778f41da3c0d7b82de9c2204e11b61a32d563
```

`fixtures-commit` запускает тесты SDK и коммитит только `testdata/upstream/`
и `docs/upstream.md`. Посторонние staged changes блокируют коммит.
Push не выполняется.

Утилита сохраняет commit и SHA-256 файлов в `docs/upstream.md`.
Локальные изменения фикстур не перезаписываются; удалённые upstream-файлы
нужно разбирать вручную. Свои фикстуры хранятся в `testdata/custom/`.
