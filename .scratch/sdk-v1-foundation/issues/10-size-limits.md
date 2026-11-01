Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/01-public-interface.md; ../../sdk-v1-foundation/issues/03-error-codes.md

# T10: Размеры запросов и endpoint limits

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Preflight encoded outbound body ≤20 MB и request metadata ≤8 KB (URL/auth/headers), depth10 и
ограничение массивов в применимых операциях. Не вводить cap успешного download/report response.
Отдельно перечислить client caps при неуточнённой API единице MB/KB и способ override.

## Специальные проверки

Base64 expansion, UTF-8 bytes, repeated headers, boundary max/max+1, streaming unknown length и
cancellation. API limits batch1000/5min и async queue4 остаются gates X02/X03, если эти операции
не включены; не заявлять реализацию отсутствующих методов. Не применять 1000 к произвольному
вложенному массиву без API основания.

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

BASE `3ada0c1`, code commit `7b2b9a9`.

- Preflight в NewRequest/Send (каждая попытка): закодированное тело ≤ `MaxBodySize` (default 20 000 000),
  метаданные ≤ `MaxHeaderSize` (default 8 000: request line, Host, Content-Length — и у записи без тела —
  и каждое значение каждого заголовка, включая Authorization), вложенность ≤ `MaxDepth` (default 10).
  0 — default, отрицательное — без проверки. Ответы не ограничиваются.
- Client caps при неуточнённых единицах (перечислить в docs): MB/KB — десятичные (строже); глубина — только
  объекты, корень = 1 (мягче, чтобы не отвергать принимаемое сервером; иначе ответит 2006); оценка
  метаданных — HTTP/1.1, без cookies Jar, своего RoundTripper и HPACK. Do с собственным запросом не проверяет.
- Base64/UTF-8 считаются по байтам json.Marshal; границы max/max+1 для тела (custom и default),
  метаданных (custom и default 8000), глубины. Исходящих потоков нет — streaming/unknown length неприменимы;
  отказ preflight не доходит до Limiter и сети.
- Массивы: общего лимита 1000 нет; остаются проверки операций (positions ≤1000 в 7 документах).
  batch1000/5 мин и очередь async 4 — gates X02/X03, не реализованы.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review Spec: глубина по объектам вместо объектов+массивов, Content-Length без тела, default 8000, допущения
  в godoc. Review Standards: лимиты вычисляются в NewClient, checkHeaderSize — свободная функция, без
  затенения max.
