Status: resolved
Stage: APPROVED
Type: task
Blocked by: None

# T03: Каталог кодов ошибок и status-aware helpers

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Полный числовой каталог официального _errors.md (545 уникальных кодов в закреплённом revision):
именованные константы с provenance, стабильный int-compatible Code, распознавание status/code
без сравнения локализованных строк. Неизвестные коды сохраняются. Сохранить errors.As/Is и
единственный префикс.

## Специальные проверки

Все code в источнике покрыты inventory test; aliases/дубли описаны. 413 возвращает
APIError/детали, не делает autosplit или retry. HTML/empty/gzip errors и raw headers доступны.
Никаких новых retry решений по тексту ошибки.

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

BASE `8559a8f`, code commit `2bf89ca`.

- `errors/codes.go` сгенерирован `errors/gen_codes.go` из `_errors.md` revision 9a811cef
  (SHA-256 источника в заголовке): 545 нетипизированных констант `CodeNNNN` с сообщением
  источника в комментарии. Совместимы с `Detail.Code int`; неизвестные коды декодируются.
  Семантических алиасов нет; дубль сообщения у 1004 и 1083 описан в тесте.
- `APIError.HasCode`, `HasCode(err, code)` — все APIError цепочки, включая ветки errors.Join;
  `StatusCode(err)`. Код 0 не совпадает. Ничего не классифицируется по тексту.
- 413 — APIError с деталями и headers через ResponseInfo, один запрос без retry/autosplit.
- Исправлен двойной префикс `moysklad:` у склеенных ошибок transport (gzip/encoding/read/limit).
- Inventory test: provenance, 545 кодов, SHA-256 отсортированного списка id источника.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
- Review Standards/Spec: исправлены `&nbsp;`, хрупкий Comments[1], HasCode по всем веткам Join,
  префикс. Осталось для docs-коммита волны: запись `_errors.md` в docs/upstream.md.
