Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/08-rate-concurrency.md; ../../sdk-v1-foundation/issues/10-size-limits.md

# F01: Transport для binary download

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Отдельный stream/download path рядом с JSON Do: writer-based Download либо явный ReadCloser
ownership; окончательный interface зафиксировать. Auth только trusted origin, signed redirects
разрешаются с удалением credentials; запрет arbitrary URL/SSRF. ResponseInfo/headers и limiter
lifetime учитывают всё чтение.

## Специальные проверки

Large responses без 20-MB upload cap, cancellation/short read/redirect chain/error body, close
на всех путях, semaphore release, cross-origin Authorization stripping; no auto retry после
частично записанного output без explicit caller recovery.

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

BASE `5a703f2`, code commit `423de8d`; общий тест опций — `b6e8188`.

Интерфейс: `transport.Client.Download(ctx, href string, w io.Writer, opts ...RequestOption) (int64, error)`.
Выбран writer, а не ReadCloser: тело закрывается внутри на всех путях, слоты limiter держатся до конца
копирования, у вызывающего кода нет обязанности закрыть. Решение записано в ADR 0001.

- href — относительный endpoint или абсолютный URL того же origin, что BaseURL; чужой origin отклоняется
  до отправки (нет SSRF). Redirect на подписанный URL другого origin разрешён, Authorization удаляется
  (общий CheckRedirect); query подписанного URL в текст ошибки не попадает.
- Без ограничения размера ответа; Config.Timeout к телу не применяется (иначе большие файлы обрывались
  бы через 6 минут) — срок задаёт ctx. gzip распаковывается явно, независимо от DisableCompression.
- Повторы по Config.Retry — только пока в w ничего не записано; после частичной записи — ошибка и число
  байт, решение за вызывающим кодом. WithoutRetry отключает. Ответ не 2xx — APIError (тело ≤ 1 МБ).
- ResponseInfo — последний ответ цепочки redirect (заголовки хранилища), в том числе отклонённого
  CheckRedirect; Attempts — число попыток. При неразбираемом Location net/http ответа не отдаёт,
  StatusCode 0 (review-daadf36). Текст вложенной ошибки такого Location обрезается
  на первом `?` с сохранением Timeout/Temporary и errors.Is/As (review-d6a4b6c, f2e2bf8).
- Тесты: 21 МБ без лимита, same-origin абсолютный href, цепочка redirect (auth внутри API, без auth на
  хранилище), ошибка хранилища, чужой href, утечка подписи, повтор 503 и отказ от повтора после
  частичной записи, обрыв тела, отмена, ошибка writer, gzip и gzip-ошибка при DisableCompression,
  освобождение слотов после успеха, ошибки и отмены.
- Review Standards: общий `acquire` для Send и Download, повторы как обещал ADR, Attempts при сетевой
  ошибке, тест gzip, строгая проверка префикса. Review Spec: ADR, таймаут, явный gzip, цепочка redirect,
  ResponseInfo после redirect, ошибки writer. Docs/status — в завершающем коммите волны.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
