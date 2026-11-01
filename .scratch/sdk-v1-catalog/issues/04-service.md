Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md

# C04: Service Get/List/Create/Update

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Полный Service core/read-write contract: name, prices/taxes/attributes/refs и документированные
fields; common Price и SalePrice. Meta files/image links читаются только как объявленные
metadata, без догрузки.

## Специальные проверки

Оба codecs, minimal/full Create/Update, nil/false/0/null, list/search/barcodes по источнику,
read-only exclusions.

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

[service](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_service.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `bc01138`, code commit `bdfa7fa`.

| Поле | Create | Update | Примечание |
| --- | --- | --- | --- |
| name | обяз. | + | непустое |
| code, externalCode, description, archived, shared | + | + | |
| vat, vatEnabled, useParentVat, discountProhibited | + | + | vat — процент |
| paymentItemType, taxSystem | + | + | открытые строки |
| owner, group, productFolder, uom | + | + | Reference, null очищает |
| buyPrice, minPrice, salePrices | + | + | деньги в копейках |
| barcodes | + | + | `model.Barcode`: объект с одним форматом |
| attributes | + | + | |
| mod__tasnif__uz | + | + | `model.TasnifUZ`, только с `WithAcceptRegion(RegionUZ)` |
| syncId | + | — | сервер пишет один раз |
| id, meta, accountId, updated, pathName, files, effectiveVat… | — | — | readonly; files — Collection без догрузки |

- Противоречия источника записаны в godoc: minPrice (таблица) / minprice (примеры) — SDK следует таблице
  и фикстуре; «услуга без закупочной цены» при buyPrice в таблице — поле читается и пишется.
- Слияние attributes и salePrices по метаданным/типу цены — серверное; что делает пустой массив
  штрихкодов, источник не описывает, SDK передаёт его как есть.
- List: поиск по name и code, фильтр `barcode` в единственном числе, archived не учитывается вместе с id —
  в godoc; кодирование проверяет listtest.
- Клиент — `service.Client` (имя Service занято сущностью), конструктор `NewService`, как в других пакетах.
- Тесты: фикстура и сжатые формы обоими codecs, golden full/update, null/false/0, readonly, HTTP-контракт
  documenttest + LocalErrors, List, валидация без запросов.
- Review Standards: Fields встраивался в обе requests → отдельные плоские CreateRequest/UpdateRequest;
  TasnifUZ вынесен в model/tasnif.go с пояснением omitempty; имя клиента задокументировано.
- Review Spec: противоречия minPrice/buyPrice, регион UZ, заметки List, слияние attributes, неподтверждённое
  утверждение о пустых массивах и о защите SyncID от дублей — исправлено в godoc.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Регистрация в общих тестах опций и сброса ResponseInfo (`entity/options_test.go`) — `cc61f6a` (найдено в review C07).
