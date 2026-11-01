# Documents review — 2026-11-02

BASE: `25a74857b58f06b7045cb5a8f79ff01dc9587499`.
HEAD: `34882fb05cb876bfb913c6cc55ddf54dbf90b1d8`.
Scope: D01–D05, общий test commit и завершающий docs commit; 7 commits, 55 files.
Метод: независимые Standards/Spec reviewers, проверка pinned первичных источников.

## Standards

Замечаний нет. Соблюдены RequestOption/ResponseInfo, общий transport seam,
Field/omitzero и отдельные read/request модели. Новые positions проверяются
обоими codecs; readonly поля отсутствуют в запросе. Дублирование приватных
do/validID допускается стандартами. Существенных baseline smells не обнаружено.
Регрессии positions/query/Bundle и прежних документов прошли.

## Spec

1. **P2 — неверно зафиксированы ограничения SalesReturn по основанию.**
   `.scratch/sdk-v1-documents/issues/05-salesreturn.md:66–67` утверждает, что
   источник не запрещает менять цену строки. Но pinned _sales_return.md:366
   прямо запрещает при Create по основанию менять price и ограничивает quantity
   количеством в исходном документе. В `entity/salesreturn/request.go:5–8`
   и `docs/salesreturn.md:56–59` эти ограничения явно не описаны.
   D05 требует фиксировать запреты сторон/основания/цен по источнику.
   Исправить Answer и описание Create: цена должна совпадать с основанием,
   quantity не больше исходного; проверяет сервер. Общий PositionRequest
   сохраняется, дополнительный GET или локальная бизнес-валидация не нужны.
   Это требование Create: без отдельного основания нельзя переносить его
   на Update как новое недоказанное ограничение.

2. **P2 — общий guide противоречит реализованному scope.**
   `docs/documents.md:82–84` объявляет List/search/filter и
   pack/things/tracking/stock вне срезов. Между тем этот же docs commit
   включает новые документы в общий guide и описывает fields=stock/declaration,
   а новые entity guides ссылаются на него как на общий контракт.
   Нарушен критерий tickets: README/guide/status отражают реальный scope.
   Актуализировать исключения и различить уже реализованные List/поля строк
   от действительно отложенных mutations, expand и файловых операций.
   Старый текст был в BASE, но в текущей волне применимость guide расширена.

Доказанных runtime defects в новых methods/models или общей пагинации
не обнаружено. Header/required/readonly матрицы D02–D05 сверены с источником;
fields allowlist и типы stock/declaration/tracking также проверены.

## Verification

- `go test -race -count=1 -p 2 ./...` — PASS.
- `go vet ./... ./testdata` — PASS.
- gofmt изменённых Go files и `git diff --check 25a7485...34882fb` — PASS.
- Дополнительные regression tests positions/query/Bundle и прежних 11 документов — PASS.
- Источники: api-remap-1.2-doc revision 9a811cef3dc4a3f3a46f070e60452361a06bf7cd,
  _general.md, _customerOrder.md, _enter.md, _loss.md, _sales_return.md,
  _retaildemand.md и _retail_sales_return.md.
- HEAD оставался 34882fb; рабочее дерево до отчёта чистое.
- Production code и статусы tickets не менялись; live API не использовался.

Standards: 0 замечаний. Spec: 2 замечания, максимум P2.

## Ответ и исправления

Оба замечания приняты; дефектов кода нет. Одним docs-коммитом поверх 34882fb, коммиты tickets не
переписывались.

| Замечание | Commit | Где |
| --- | --- | --- |
| P2, ограничения SalesReturn по основанию | `1576899` | godoc CreateRequest, docs/salesreturn.md, Answer D05 |
| P2, общий guide противоречит scope | `1576899` | docs/documents.md, README (Product) |

- SalesReturn: для Create по основанию перечислены совпадение контрагента/юрлица/валюты, соответствие
  строк, quantity не больше основания, запрет строк не из основания и запрет менять price (L366).
  Проверяет сервер, без Get основания и локальной валидации. На Update запрет price не перенесён:
  источник формулирует его только для Create; для Update остаются валюта и строки не из основания.
  Неверная фраза Answer D05 исправлена.
- documents.md: реализованное (List с filter/search/order, явный pages.All, поля строк общей формы,
  fields=stock/declaration у GetPositions) отделено от отложенного (Delete, batch, templates/new,
  wizard, expand, mutations позиций, файлы, автоматический обход, финансовые вычисления).
- README: устаревшая строка Product («нет List, пагинации, retry, rate limiting, syncId») заменена
  ссылкой на общие механизмы; minimumStock отмечен как отложенный.
- go test -race ./... (-p 2), go vet ./... ./testdata, gofmt, `git diff --check 25a7485 HEAD` — PASS.
