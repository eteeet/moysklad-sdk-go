Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-release/issues/01-facade.md

# R03: CI и проверки выпуска

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Go1.27 toolchain, tests/race/vet/gofmt/diff check, manifest хешей fixtures и проверка public
examples; migration notes/API compatibility и CHANGELOG/v1 version proposal. Существующие
команды сохраняются; network/live интеграция не обязательна для unit CI.

## Специальные проверки

Чистый checkout воспроизводит checks, pin toolchain/workflow permissions, нет секретов и live
writes. При локальных TCP resource flakes допустим -p1 с явным reporting. Publish/tag/release не
выполняются без отдельной команды.

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

## Answer

BASE `ca8cac0` (R02). Решение пользователя: только `make check`, без GitHub workflow; CHANGELOG
предлагает v1.0.0, тег не ставится.

- `make check`: gofmt, `go vet ./... ./testdata`, `go test -race -count=1 -p $(CHECK_P) ./...`
  (по умолчанию 2; при локальных сбоях TCP допустим `CHECK_P=1` с явной записью в отчёте),
  `git diff --check` всего дерева HEAD от пустого дерева и незакоммиченных изменений отслеживаемых
  файлов (неотслеживаемые файлы не проверяются). Сеть, токены
  и live-запись не нужны; существующие цели fixtures/example не менялись.
- Toolchain: `toolchain go1.27.3` в go.mod при `go 1.27`.
- `testdata/upstream/** -whitespace` в `.gitattributes`: upstream синхронизируется байт в байт
  (как фикстуры ассортимента), иначе проверка всего дерева падает на пустых строках в конце файлов.
- `fixtures_test.go`: блок syncfixtures docs/upstream.md перечисляет ровно файлы testdata/upstream
  с верным SHA-256; копии источника в testdata/custom совпадают с SHA-256 своих README (добавлены
  хеши predefined.json country, expenseitem, uom). Изменение байта или лишний upstream-файл — FAIL.
  Golden и `read.json` — наши независимые файлы, хеша нет.
- Публичные примеры проверяются `go vet`/`go test ./...`: R02-сценарии и `Example*` компилируются.
- CHANGELOG.md: v1.0.0 — добавленное, совместимость с 7442803 (сравнение экспортированного API
  `go doc -all`: удалений нет; opts у 17 сервисов и `positions.GetPage` ломают method value и
  интерфейсы со старой сигнатурой; `positions.PageOptions` — alias `query.PageOptions`, литерал без
  имён полей не компилируется), что не входит в v1.

Review Standards: манифест upstream читается только из блока syncfixtures (формат строки пишет
генератор testdata), порядок Rel/хеша в WalkDir; status-technical и README — в docs-коммите.
Review Spec: хеши трёх неизменённых копий источника, литерал PageOptions без имён в CHANGELOG,
точная формулировка diff --check. Оба ревьюера видели сбои `dial tcp 127.0.0.1 … operation timed
out` под параллельной нагрузкой нескольких прогонов (около 1700 TIME_WAIT); те же пакеты отдельно
проходят — это ресурс машины, не логика.

## Источники

[транспортный research](../../sdk-v1/research/transport.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Comments

- 2026-11-02, решение пользователя: CI — только цель `make check`, без GitHub workflow.
  CHANGELOG предлагает версию v1.0.0; тег не ставится.
