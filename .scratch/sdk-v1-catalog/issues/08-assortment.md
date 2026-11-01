Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md; ../../sdk-v1-catalog/issues/04-service.md; ../../sdk-v1-catalog/issues/06-bundle.md; ../../sdk-v1-catalog/issues/07-variant.md

# C08: Assortment List/search и mixed rows

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Один смешанный список product/service/bundle/variant/consignment, единые search/filter/order и
pagination; read union с known typed four и Consignment reference/common-read, unknown meta.type
+ raw JSON сохраняется. Без generic write resource.

## Специальные проверки

Явный выбор режима: legacy без beta header либо облегчённый assortmentWithoutStock с own
capability/default groupBy. Не включать beta молча; default и supported filters закрепить до
code. Нет скрытых stock report запросов, mixed rows включая consignment и неизвестный тип, один
List вместо склейки четырёх списков.

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

## Источники

[assortment](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_assortment.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `9eccf22`, code commit `df29cc6`; общий тест опций — `cc61f6a`.

| Метод | Endpoint / режим | Остатки |
| --- | --- | --- |
| `List(ctx, Options)` | GET entity/assortment + `X-Lognex-Remap-Beta-Feature: assortmentWithoutStock` | нет |
| `ListLegacy(ctx, Options)` (Deprecated) | GET entity/assortment без заголовка | stock, reserve, inTransit, quantity |

- Решение пользователя: два явных метода. List ставит beta-тег последним, поверх опций вызова — отступление
  от порядка ADR 0001 записано в godoc: иначе чужой тег молча переключил бы режим.
- `Options{query.ListOptions; GroupBy}`; GroupBy — открытый тип (product, variant, consignment), умолчание —
  без параметра. fields не принимается, groupBy через Extra отклоняется.
- Фильтры, умолчания (archived — только неархивные, withSubFolders — true, type=consignment только с
  groupBy=consignment), недоступные фильтры, LIKE-комбинации и сортировки нового endpoint — в godoc List;
  legacy-фильтры и поиск как filter search= — в godoc ListLegacy. Проверяет сервер.
- Противоречие: groupBy=variant — «товары и модификации», но пример без groupBy содержит услуги и
  комплекты; SDK передаёт значение как есть.
- `Row`: union по meta.type — Product, Service, Bundle, Variant (модели C04, C06, C07, product),
  `Consignment` (общая модель чтения, 13 полей); неизвестный тип — только Type и Raw. Запись Row —
  Raw байт в байт. Остатки — `Balances`.
- Скрытых запросов нет: остатки для List — отчёт report/stock/all/current по id (вне этого ticket).
- Фикстуры: примеры List и groupBy=consignment из источника (`testdata/custom/assortment`, SHA-256 в
  README), legacy — upstream `assortment.json`.
- Review Standards: типизированный GroupBy, Balances вместо дублирования четырёх полей, `json:"-"` у Raw,
  комментарий о записи Row, t.Run по codecs. Дублирование с list.Entity (нужен groupBy) — осознанно.
- Review Spec: поддерживаемые фильтры и умолчания, minimumBalance и LIKE-комбинации, противоречие groupBy,
  поиск legacy; listtest List — на строке нового endpoint без остатков.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
