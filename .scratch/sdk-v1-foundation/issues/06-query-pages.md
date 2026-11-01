Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/01-public-interface.md; ../../sdk-v1-foundation/issues/05-filter-order.md

# T06: Query options и получение одной страницы

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Limit/offset, search, namedfilter, fields (одно значение), report dates/groupBy/interval и
полный инвентарь прочих query names из research. Typed options только на поддерживающих
endpoints; общий url.Values сохраняет новые параметры. Обычные Collection[T], plain arrays и
report-specific ответы не смешиваются. Expand не добавляется.

## Специальные проверки

Endpoint-specific defaults: entity lists обычно 1000, audit 25/max100; GetPositions existing
Limit=0→1000 сохраняется. fields restrictions scoped, namedfilter href, exact query encoding.
Для неподдержанных ресурсов raw passthrough не объявляется готовым typed service. async не
разрешается общей Get options без специальной операции.

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

[транспортный research](../../sdk-v1/research/transport.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `a220c7c`, code commit `974f8c3`.

- `query.ListOptions` (limit/offset/filter/order/search/namedfilter/fields + Extra) и
  `query.ReportOptions` (momentFrom/momentTo в московском времени, groupBy, interval, filter + Extra).
  Нулевые поля не передаются; max limit задаёт endpoint (`Values(maxLimit)`: 1000, у audit 100);
  fields — одно значение; namedfilter — абсолютный href; expand/async и повтор typed-имён в Extra
  отклоняются; Extra копируется. Константы `Param*` — инвентарь имён из research.
- `model.DateTime.QueryValue`: «2006-01-02 15:04:05», миллисекунды только ненулевые
  (риск: принимают ли отчёты миллисекунды, не подтверждено — проверить в F03).
- `internal/list.Page[T]`: один GET коллекции {context, meta, rows}; голый массив, null rows,
  нет meta.href, не-200 — ошибка. Не общий CRUD: только коллекции. positions.GetPage использует его;
  его Limit=0→1000 сохранён (`query.MaxLimit`).
- Для T11/T12: публичные List подключают ListOptions и ограничивают fields/search по endpoint.
  Для A07: audit — maxLimit 100, default сервера 25.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review Spec: ReportOptions больше не блокирует search/order/fields в Extra, убраны лишние
  limit/offset и проверка периода, проверка namedfilter, тест отмены. Review Standards: единый
  список typed-имён (build), nonZero, query.MaxLimit в positions, тест без ветвления по имени.
