# Catalog review — 2026-11-02

BASE: `8eda31b6346b2b9f2e675e4806995ac4dc3d72ff`.
HEAD: `2b91a89bb555634ae7822a4e5ea0a30e27f5a070`.
Scope: C01–C08, общий test commit и завершающий docs commit; 10 commits, 70 files.
После остановки ревью начато заново, два независимых reviewers получили новый контекст.
Production code и статусы tickets не менялись; live API не использовался.

## Standards

**P2 — ListLegacy молча переключается на новый режим.**
`entity/assortment/service.go:77–79` передаёт beta header из Config.Header и per-call
options без изменений. Требование C08: «Явный выбор режима: legacy без beta header»;
общая spec указывает, что режимы не взаимозаменяемы.
При assortmentWithoutStock сервер выберет новый endpoint без остатков и с другими
доступными фильтрами. Публичное воспроизведение подтвердило отправку заголовка и
отсутствие ошибки для обоих источников настройки.
Нужно гарантировать legacy либо отклонять конфликт до запроса; regression tests
должны охватить client defaults и per-call options.

Других однозначных нарушений реализации и обоснованных smell-замечаний нет.
Дублирование приватных do/validID допускается правилами проекта.

## Spec

1. **P2 — ListLegacy не фиксирует legacy-режим.**
   `entity/assortment/service.go:77–79`. C08 требует явного выбора режима, freeze
   обещает GET entity/assortment без beta header. Config.Header или WithBetaFeature
   сохраняет assortmentWithoutStock. Pinned _assortment.md:9–15 подтверждает,
   что именно этот header выбирает новый endpoint. Это то же замечание, что Standards.

2. **P3 — документация обещает неподтверждённую семантику Service barcodes.**
   `docs/catalog.md:24` объединяет Service с Bundle/Variant утверждением
   «штрихкоды … заменяются целиком». C04 Answer явно признаёт неизвестность
   поведения пустого массива штрихкодов; соответствующие разделы pinned Service
   правила полного замещения не описывают. Это нарушает критерий отражать
   ограничения и неизвестные факты в публичной документации.
   Разделить описания и сохранить оговорку для Service. Фактическое поведение API
   этим ревью не установлено.

## Коллизия требований: Bundle Quantity

`entity/bundle/request.go:6–12` пишет Quantity, обосновывая это примерами.
Pinned _bundle.md:190 помечает quantity read-only; master spec:90 и
docs/documents.md:63 отдают приоритет этому marker.
Вместе с тем C06 What-to-build прямо требует embedded write ID/quantity/assortment,
и freeze это повторяет. Реализация следует конкретному C06, но расходится с общим
нормативным правилом. Комментарий/Answer не доказывают отдельного пользовательского
решения об исключении. Требуется явное разрешение коллизии и запись в контракте.
Автоматически удалять Quantity нельзя: это нарушит прямое требование C06.
Это не утверждение о том, что API отвергает quantity, и не однозначный runtime defect.

Currency system/rateUpdateType отличается: их запись при создании системной валюты
явно разрешена нормативным текстом источника, а не только примерами.

## Verification

- `go test -race -count=1 -p 1 ./...` — PASS.
- `go vet ./... ./testdata` — PASS; gofmt изменённых Go files — PASS.
- `git diff --check 8eda31b...2b91a89` — FAIL:
  `testdata/custom/assortment/list.json:394`, trailing whitespace.
  Это tooling-замечание, не отдельный review finding. Fixture скопирована из
  официального источника без изменений и имеет SHA-256 provenance; любое изменение
  требует согласованного обновления описания происхождения, не правки upstream.
- Публичные воспроизведения без сети:
  `/private/tmp/catalog_restart_modes.go`, `/private/tmp/catalog-review-repro.go`.
  ListLegacy с Config.Header и per-call beta отправляет assortmentWithoutStock.
- Повторно проверены models/requests/routes/query, JSON обоих codecs, PriceType
  replace-all, Variant characteristics/packs, shared PageOptions и совместимость
  positions. Других подтверждённых дефектов не обнаружено.
- HEAD оставался 2b91a89, рабочее дерево до отчёта чистое.

Standards: 1 замечание, максимум P2. Spec: 2 замечания, максимум P2.
Всего 2 уникальных implementation/docs findings; отдельно коллизия Quantity.

## Ответ и исправления

Все замечания приняты; Standards P2 и Spec №1 — один дефект. Отдельные коммиты поверх 2b91a89,
коммиты tickets не переписывались.

| Замечание | Commit | Регрессия / запись |
| --- | --- | --- |
| P2, ListLegacy молча уходит в новый режим | `8bde787` | `TestListLegacyRejectsBetaMode`, `TestHeaderValues` |
| P3, docs обещают замену штрихкодов Service | `a9ebfba` | docs/catalog.md |
| Коллизия Quantity | `a9ebfba` | contracts.md (bundle), spec.md, docs/catalog.md, Answer C06 |
| diff --check, пробелы в фикстуре | `b1f3e9c` | `.gitattributes` |

- ListLegacy: до запроса проверяет итоговый `X-Lognex-Remap-Beta-Feature` (Config.Header с заменой
  опциями вызова, значения через запятую); assortmentWithoutStock — ошибка без отправки, info обнулён.
  Другие теги и замена тега опцией вызова проходят. Итоговый заголовок считает новый
  `transport.Client.HeaderValues` тем же слиянием, что и отправка. Регрессия на 2b91a89 падает
  (4 подслучая), на исправленном — проходит.
- Штрихкоды: у Bundle и Variant — замена целиком (у Variant правила создания не подтверждены), у Service
  поведение неизвестно.
- Quantity: решение пользователя 2026-11-02 — оставить запись как исключение из правила readonly;
  удаление нарушило бы прямое требование C06.
- Фикстура не менялась: файл совпадает с источником байт в байт, SHA-256 в README прежний;
  whitespace-проверка для `testdata/custom/assortment/*.json` отключена через `.gitattributes`.
- go test -race ./... (-p 2), go vet ./... ./testdata, gofmt, `git diff --check 8eda31b HEAD` — PASS.
