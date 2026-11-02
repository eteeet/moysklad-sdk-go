Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/03-error-codes.md

# X01: Delete выбранных ресурсов

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

По-resource capability matrix single Delete для поддержанных entity/metadata/files.
Метод/маршрут/status индивидуальны. Не смешивать delete, recycle bin и purge. No write retries
по умолчанию.

## Специальные проверки

204/200 по endpoint, APIError body, absent JSON, UUID, permissions; не генерировать Delete для
assortment/context/audit/report и прочих read-only endpoints.

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

BASE `30c9318`, коммит `c0c5d1c`. Решение пользователя 2026-11-02: захватка 8 входит в v1, объём —
«всё по источнику», интерфейс `Delete(ctx, id, opts...) error`.

- Matrix по разделам «### Удалить …» revision 9a811cef: Delete у всех 34 сервисов `entity/*`, кроме
  assortment и pricetype (раздела нет); role — пользовательская роль. Плюс `metadata.DeleteState`
  (DELETE entity/{type}/metadata/states/{id}) и `files.Delete` (DELETE entity/{type}/{id}/files/{fileId}).
  Позиции документов, компоненты комплекта, счета контрагента/юрлица, события и корзина (POST …/trash)
  не входят; Delete — не корзина.
- `internal/write.Delete`: один DELETE без тела, успех 200 или 204 (примеры — 200 с пустым телом),
  тело не читается; не 2xx — APIError transport. DELETE не повторяется (retry только GET/HEAD).
- UUID, тип сущности и id файла проверяются до запроса; ResponseInfo сбрасывается первым.
- Тесты `entity/delete_test.go`: опции, ResponseInfo и 401 у каждого Delete (общая проверка T02);
  200/204 — успех, 201/403(APIError 1016)/503 — ошибка, ровно один вызов при Retry{MaxAttempts:3};
  неверный id у всех 34 сервисов, files и states — без запроса; у assortment, pricetype, contexts,
  audit, turnover нет Delete.
- Review Standards: doc `internal/write` сужен до удаления (массовые операции — X02); `s, _ :=` в
  тестах оставлен по образцу options_test (NewService падает только на nil client); проверка read-only
  расширена на contexts/audit/turnover. Review Spec: проверка UUID у всех 34 сервисов, случай 403.
  Docs, README и статусы — в docs-коммите волны.
- Не проверено: живой API; якоря godoc построены транслитерацией заголовков (схема совпадает с 120
  существующими якорями).
