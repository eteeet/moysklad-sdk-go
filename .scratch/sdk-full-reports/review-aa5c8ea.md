# Review c3bace0…aa5c8ea

- BASE: `c3bace0997f7b8b4b08ea6fd8caf6ceff9279621`.
- HEAD: `aa5c8eac53a7f33db55639253886eed04ca9559d`.
- Diff: `git diff c3bace0997f7b8b4b08ea6fd8caf6ceff9279621...aa5c8eac53a7f33db55639253886eed04ca9559d`.
- Scope: M01–M08, I01–I23, G12/G13; 331 files, 28149 additions / 159 deletions.
- Standards и Spec проверены независимыми агентами; печать, интеграции и transport дополнительно проверены основным агентом. Реализация не изменялась.

## Standards

Жёстких нарушений не найдено.

### На суждение, P2: source-response fixtures производства частично обходят SHA-256

`entity/productiontask/nested_test.go:24` объявляет фикстуры примерами ответа источника, но в `testdata/custom/productiontask` нет README с хешами. Аналогичный пробел у M01/M02/M05/M06/M08. У M03 `testdata/custom/processingprocess/README.md:1–5` прямо называет `read.json` неизменённым примером источника, однако таблицы SHA-256 нет. M04 оформлен полностью.

`TestCustomManifests` в `fixtures_test.go:75–92` проверяет только перечисленные в существующих README таблицы. Отсутствующий manifest не приводит к падению проверки.

Основание: общий контракт производства (`.scratch/sdk-full-production/spec.md:125–126,135`) требует SHA-256 provenance; README обещает проверку SHA-256 в `make check`. Строгое нарушение спорно: формулировка нормы отдельно называет upstream, а эти файлы лежат в custom. Рекомендация: добавить manifests для ответов, скопированных из источника, и обнаружение пропущенных manifests; не требовать source-хеши для самостоятельно составленных golden.

## Spec

### P2: корректный assortment href без Meta.Type отклоняется до HTTP

`report/byoperations/service.go:28–30` требует `Meta.Type` из списка product/variant/consignment. Поэтому `Options{Assortment: model.Reference{Meta: model.Meta{Href: productHref}}}` не позволяет запросить ни Stock, ни Reserve, ни InTransit. Тип отсутствует в отправляемом запросе: `values()` использует только `Meta.Href` для `filter=assortment=<href>`.

Источник `md/reports/_report_by_operations.md:12,38,43–45` требует фильтр со ссылкой на товар/модификацию/партию и не требует отдельного type. Контракт `.scratch/sdk-full-reports/spec.md:293,298` не допускает локальное отклонение неизвестного enum/type и оставляет бизнес-ограничения API. Требование `Meta.Type` не следует из ticket I05.

Исправление: не делать `Meta.Type` обязательным для URL-фильтра; проверить href-only reference через публичные методы и настоящий transport. Сохранить проверку отсутствующего href. Актуализировать godoc и соответствующие отрицательные тесты.

Воспроизведение без изменения SDK:

```sh
go test -race -overlay=/private/tmp/moysklad-review-byoperations-href-only/overlay.json -run '^TestReviewHrefOnlyAssortment$' -count=1 -v ./report/byoperations
```

Один и тот же валидный product href: `Type="product"` — PASS, один HTTP-запрос; `Type=""` — FAIL, `calls=0`, `assortment must be a product, variant or consignment reference`.

Других подтверждённых замечаний по Spec не найдено. Ранее утверждённые readonly-исключения не переоткрывались. Ограничение print.Download на same-origin предусмотрено I12 и явно описано; origin реального Location без live API не проверялся.

## Проверки

- `make check` на HEAD — PASS полностью: форматирование, vet, race-тесты, manifests и diff-check.
- Дополнительный regression repro выше — FAIL, ожидаемо подтверждает замечание Spec.
- Live API и credentials не использовались. Имеющиеся изменения AGENTS.md и untracked planning не изменялись, коммитов и push нет.

Итог: Standards — 0 жёстких нарушений, 1 замечание P2 на суждение; Spec — 1 подтверждённое замечание P2.

## Answer

- Spec P2 (byoperations) — исправлено в `e594825`: обязателен только непустой `meta.href`,
  `Meta.Type` до HTTP не проверяется (в запрос не уходит, тип проверяет сервер). Отрицательный случай
  `bundle` заменён положительным «href only» через Stock/Reserve/InTransit и настоящий transport;
  godoc и docs/report-byoperations.md обновлены. Других проверок `Meta.Type` в новом коде нет.
- Standards P2 (manifests) — исправлено в `0c6de42`: `TestCustomManifests` требует запись SHA-256 для
  каждой `testdata/custom/*/*.json`, кроме `*.golden.json`; без README или строки — ошибка (проверено
  удалением README: тест падает). Добавлены записи для 51 файла, происхождение проверено поиском по
  источнику revision 0e1b3d61: 36 — дословные копии (файл и строки), 5 — те же значения с другими
  отступами (строки списков, branch с исправленной запятой), 10 — составленные для теста фикстуры v1
  (их README это уже описывали). Шире scope review: проверка поймала и захватку 3, и v1.
- make check на `0c6de42` — PASS.
