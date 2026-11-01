Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/06-query-pages.md

# D01: Расширения общей формы позиций

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Зафиксировать union новых документированных полей rows
customerorder/enter/loss/salesreturn/retaildemand/retailsalesreturn: типы, money vs counts,
writable markers, Reference/pack/slot/stock enrichment для нужных fields options. Реальный
type/meaning conflict требует отдельной формы, простое отсутствие поля — нет.

## Специальные проверки

Write matrix из contracts.md: reserve/shipped/cost/inTransit и прочие не добавляются по
аналогии. Явный readonly важнее примера, неизвестные refs/server constraints не валидируются
локальным GET. Оба codecs и прежние семь aliases не регрессируют.

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

[customerorder](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_customerOrder.md), [enter](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_enter.md), [loss](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_loss.md), [salesreturn](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_sales_return.md), [retaildemand](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retaildemand.md), [retailsalesreturn](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retail_sales_return.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `25a7485`, code commit `448ff35`.

| Поле строки | Read | Write | Источник / примечание |
| --- | --- | --- | --- |
| reserve | `float64` | Field | customerorder; количество, не деньги |
| shipped | `float64` | — | readonly (уже было) |
| taxSystem | `string` | Field | customerorder, enum открыт |
| reason | `string` | Field | enter, loss |
| slot | `*Reference` | Field | enter, loss, salesreturn |
| pack | `*Pack{ID, Quantity, UOM, Barcodes}` | `PackRequest{ID}` | ссылка таблиц на «Упаковки Товара»; примеры пишут только id |
| country | `*Reference` | Field | enter, salesreturn |
| gtd | `*GTD{Name}` | Field | enter, salesreturn (объект); в declaration — строка |
| things | `[]string` | Field | enter: Object(String) против Array(String) у остальных — массив |
| trackingCodes | `[]TrackingCode` (рекурсивно) | Field | retaildemand: cis/type необходимы при создании; есть в фикстурах demand/supply |
| trackingCodes_1162 | `[]TrackingCode1162` | — | readonly |
| declaration | `[]Declaration{GTD, RNPT, Country, Quantity}` | — | loss, retaildemand, `fields=declaration`, readonly |
| stock | `*Stock{Cost *float64, Quantity, Reserve, InTransit, Available}` | — | `fields=stock`, readonly; Cost — копейки, нет у возвратов без основания |
| cost | `float64` | Field | уже было; теперь и возвраты без основания |

- `query.PageOptions.Fields` — одно значение. `positions.GetPage` принимает только значения из
  `positionFields`: stock для 10 документов из общего раздела (_general L1122), declaration для loss
  и retaildemand. Список — в одном месте пакета positions, а не в каждом документе: он выписан из
  одного раздела источника. Компоненты комплекта fields отклоняют.
- Новые скалярные поля чтения — с omitzero; прежние семь документов не регрессировали, их фикстуры
  теперь дополнительно читают things/trackingCodes.
- Review Standards: trackingCodes записывается по маркеру «Необходимо при создании» (не по аналогии);
  противоречие things и устаревший комментарий Cost исправлены; формулировка ошибки fields едина;
  allowlist оставлен в positions (обоснование выше). Docs (positions.md, demand.md, status) — в
  завершающем docs-коммите волны.
- Review Spec: источник полей Pack указан; остальное — как Standards.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
