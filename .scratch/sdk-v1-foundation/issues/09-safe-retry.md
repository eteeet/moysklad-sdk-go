Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/03-error-codes.md; ../../sdk-v1-foundation/issues/08-rate-concurrency.md

# T09: Ограниченные retries/backoff

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Operation-aware retry policy: только подтверждённые reads GET/HEAD по умолчанию;
POST/PUT/security-token/async-submit — без автоматического повтора. 429/временные 5xx/transport
failures с ограничением attempts, context, jitter и корректными server hints. Окончательные
metadata и APIError сохраняются.

## Специальные проверки

X-Lognex durations в миллисекундах, стандартный Retry-After — seconds/date; не обрезать server
minimum своим cap. Каждая попытка новый replay body/headers и limiter admission; close
промежуточных body до ожидания. Nil/unreplayable stream, decode error после 2xx и canceled
context не повторяются.

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

BASE `7fb0c21`, code commit `cd465f9`.

- `Config.Retry *RetryPolicy`; nil — без повторов. Повтор только GET/HEAD, не async (ни в query,
  ни в строке endpoint), не при `WithoutRetry` (его ставит security.CreateToken). Поводы: 429, 500,
  502, 503, 504 с чистой APIError и временные сбои соединения (EOF, reset, refused, pipe, timeout);
  TLS/неверный URL, отмена ctx, decode после 2xx, неизвестная кодировка — нет. 1073 — как 429, бюджет
  Limiter не меняется.
- Пауза = max(подсказка, full-jitter экспоненты): X-Lognex-Retry-After/Reset в мс, Retry-After в секундах
  или HTTP-date; MaxDelay ограничивает только экспоненту, подсказку не обрезает. Пауза больше MaxWait или
  дедлайна ctx — повтора нет. Defaults 3 попытки, 500 мс, 30 с, 1 мин — выбор SDK.
- Каждая попытка заново собирает запрос и проходит Limiter; слоты свободны во время паузы; body закрывает Do.
- ResponseInfo — последней отправленной попытки (без ответа — StatusCode 0, Header nil), Attempts — число
  отправленных. Отмена в паузе или в ожидании Limiter возвращает последнюю APIError вместе с ctx.Err().
- Тесты на часах без sleep: решения, паузы и подсказки, обрыв (с Connection: close — иначе net/http сам
  повторяет GET на переиспользованном соединении), отмены, слоты в паузе, параллельные вызовы под -race.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review Spec/Standards: исправлены info после финального обрыва, async в endpoint, узкий список временных
  ошибок, кодировка на 5xx, потеря APIError при отмене, имя sendOnce, единый список статусов.
