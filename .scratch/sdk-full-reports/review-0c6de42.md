# Review aa5c8ea…0c6de42

- BASE: `aa5c8eac53a7f33db55639253886eed04ca9559d`.
- HEAD: `0c6de42bbc8fd7862627b3adada930263978cd4e`.
- Diff: `git diff aa5c8ea...0c6de42` — 39 files, 337 additions / 12 deletions.
- Commits: `e594825` (byoperations href-only), `0c6de42` (custom SHA-256 manifests).
- Независимые Standards/Spec review; реализация не менялась.

## Standards

Новых замечаний нет. Предыдущее advisory P2 закрыто:

- Добавлены SHA-256, revision и происхождение фикстур, включая productiontask и processingprocess; преобразования примеров отмечены явно.
- `TestCustomManifests` проверяет наличие manifest entry для каждой текущей `custom/*/*.json`, кроме независимых `.golden.json`.
- Отсутствующий README больше не позволяет молча пропустить фикстуру ответа.

## Spec

Новых замечаний нет. P2 закрыт:

- `byoperations.Options` требует только непустой `Meta.Href`; `Meta.Type` не блокирует запрос.
- Положительный href-only case покрывает Stock/Reserve/InTransit через настоящий transport.
- Проверки отсутствующего href, Limit/Offset, Extra, ошибок и ResponseInfo сохранены; godoc и guide согласованы.
- Исходный внешний repro теперь PASS: один product href проходит как с `Type="product"`, так и с `Type=""`.

## Что осталось

Статусы всех 91 issue-файла совпадают с таблицами карты:

| Scope | resolved | needs-info |
| --- | ---: | ---: |
| dictionaries | 27 | 0 |
| catalog | 13 | 0 |
| documents | 19 | 1 |
| production | 8 | 0 |
| reports | 23 | 0 |
| Всего | 90 | 1 |

- O07 — `agentAccount/organizationAccount` у cashin/cashout: таблицы и fixtures не содержат счетов, текст Update их упоминает. Текущий SDK следует таблицам. Для закрытия вопроса нужен ответ МойСклад либо отдельная согласованная live-проверка. Другие tickets от O07 не зависят.
- За пределами 91 tickets: expand, раскрытые объекты и ограничения запросов; fields в List справочников/документов; UOM expand. Только отдельной командой.
- Документационные хвосты: revision в status-*; исторические этапы/TODO coverage.md; устаревшая фраза «G12/G13 ждут O01/O04» в map.md:43. Фактические строки G12/G13 уже resolved. В status-backlog строка S16 упоминает fsrarId/utmUrl как закрытое, хотя решение S16 и dictionaries.md прямо оставляют их неподдержанными. Это отложенные поля, не новая ошибка исправлений.
- Live read-only и Product.Create/Update проверки не выполнялись; они требуют отдельного согласия. Тег/релиз — отдельная команда. CI и Vendor API не входят в выбранный v1.

## Проверки

- Исходный repro: `go test -race -overlay=/private/tmp/moysklad-review-byoperations-href-only/overlay.json -run '^TestReviewHrefOnlyAssortment$' -count=1 -v ./report/byoperations` — PASS.
- Первый make check остановился в TestDownloadGzip из-за sandbox: bind loopback port, operation not permitted. Повтор с разрешёнными локальными портами — PASS полностью, включая vet, race-тесты, manifests и diff-check. Лог: `/private/tmp/moysklad-review-0c6de42-check.log`.
- Live API не использовался. Существующие изменения AGENTS.md и planning не менялись; коммитов и push нет.

Итог review исправлений: Standards — 0 замечаний; Spec — 0 замечаний.
