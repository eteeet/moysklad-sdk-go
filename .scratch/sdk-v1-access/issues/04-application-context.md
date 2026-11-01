Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md

# A04: GET context/application

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Отдельный context/application.Get(ctx) с App token, meta/id/appUid; не выдавать за Vendor API,
не реализовывать установку/активацию решения. entity/application/{installationId} — отдельная
будущая возможность, не alias контекста.

## Специальные проверки

Один GET без path id/query/body; App auth и APIError, отсутствующие поля не заполняются.
NewClient не делает этот запрос.

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

[general](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `842adec`, code commit `c7b17f3`.

- `contexts.Service.Application(ctx, opts...)`: один GET context/application; `ApplicationContext` —
  meta, id (id установки) и appUid. Не Vendor API, не установка/активация; NewClient не вызывает;
  entity/application/{id} не алиасится.
- Тесты: пример источника, HTTP-контракт documenttest (успех, ошибки, пустой/отсутствующий id),
  локальные ошибки. Общий `get[T]` с A03.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review: обе оси вместе с A03; дублирование с Employee убрано через get[T].
