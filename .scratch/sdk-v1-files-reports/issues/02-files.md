Status: resolved
Stage: APPROVED
Type: task
Blocked by: ../../sdk-v1-foundation/issues/02-headers.md; ../../sdk-v1-foundation/issues/06-query-pages.md; ../../sdk-v1-foundation/issues/10-size-limits.md; ../../sdk-v1-files-reports/issues/01-download-transport.md

# F02: Files List/Upload/Download

[Spec](../spec.md), [карта](../../sdk-v1/map.md),
[общий контракт v1](../../sdk-v1/spec.md),
[каталог полей](../../sdk-v1/contracts.md).

## What to build

Files как owner-scoped entity/{type}/{id}/files; List metadata, POST upload filename/content
Base64, download по документированным meta.downloadHref. File metadata без вымышленного
entity/file CRUD или общего GetFile endpoint.

## Специальные проверки

≤10 upload/request, ≤100 files/object server knowledge, encoded-body bound отдельно от decoded
file. Names/media type/no filename path traversal; signed external URL secret не логируется.
Никакого pre-GET для подсчёта текущих100 файлов. Images — отдельная будущая feature, Delete file
— X01.

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

[files](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_files.md)

## Out of Scope

Expand, live API и не перечисленные What to build lifecycle операции.
Этот DRAFT не является командой реализации или публикации.

## Answer

BASE `423de8d`, code commit `bf97b07`.

| Метод | Endpoint | Примечание |
| --- | --- | --- |
| `List(ctx, type, id, query.PageOptions)` | GET entity/{type}/{id}/files | одна страница; fields не принимается |
| `Upload(ctx, type, id, []Upload)` | POST entity/{type}/{id}/files | 1–10 файлов, `{filename, content}` (Base64), ответ — весь список файлов объекта |
| `Download(ctx, File, w)` | GET meta.downloadHref | через `transport.Client.Download` |

- `File`: meta (с downloadHref), title, filename, size (байты), created, createdBy; miniature и tiny — только у
  изображений. Файлового CRUD, GetFile, удаления (X01) и изображений нет.
- Лимит 100 файлов на объект — сервер, без предварительного Get. Предел тела — закодированный JSON
  (MaxBodySize, Base64 больше файла на треть), отдельно от размера файла; тест.
- Имя файла: непустое, без разделителей пути, NUL, «.» и «..»; длину (String(255)) проверяет сервер.
  Тип содержимого SDK не проверяет: источник его для загрузки не описывает. nil Content уходит как "".
- Подпись временной ссылки не попадает в текст ошибки (transport, F01); логирования в SDK нет.
- Фикстуры: примеры List и Upload из источника (`testdata/custom/files`, SHA-256 в README).
- Review Standards: ошибки Download с операцией, тесты Download на уровне сервиса (LocalErrors, 404).
  Review Spec: nil Content, неподтверждённое скачивание миниатюр убрано из godoc, длина имени и media type
  записаны как решения сервера.
- Проверки: go test -race ./..., go vet ./... ./testdata, gofmt, diff --check — PASS.
