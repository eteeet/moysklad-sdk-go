# Печать: шаблоны, документы, этикетки и ценники

Пакет `print`. Источник — `api-remap-1.2-doc` revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`:
[md/dictionaries/_template.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_template.md),
[md/documents/_print.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_print.md),
[md/workbook/_workbook_print_templates.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/workbook/_workbook_print_templates.md),
[md/dictionaries/_labels.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_labels.md).
`print.NewService(client)` отклоняет nil client. Все методы принимают `opts ...transport.RequestOption`,
первой строкой сбрасывают `ResponseInfo` и проверяют аргументы до HTTP; повторов записи нет.

## API

| Метод `print.Service` | Endpoint | Результат |
| --- | --- | --- |
| `ListEmbeddedTemplates(ctx, type)` | GET entity/{type}/metadata/embeddedtemplate | `*model.Collection[Template]` |
| `GetEmbeddedTemplate(ctx, type, id)` | GET entity/{type}/metadata/embeddedtemplate/{id} | `*Template` |
| `ListCustomTemplates(ctx, type)` | GET entity/{type}/metadata/customtemplate | `*model.Collection[Template]` |
| `GetCustomTemplate(ctx, type, id)` | GET entity/{type}/metadata/customtemplate/{id} | `*Template` |
| `Export(ctx, type, id, ExportRequest)` | POST entity/{type}/{id}/export | `*Result` |
| `ExportSet(ctx, type, id, ExportSetRequest)` | POST entity/{type}/{id}/export | `*Result` |
| `ExportLabels(ctx, type, id, LabelRequest)` | POST entity/{type}/{id}/export | `*Result` |
| `Download(ctx, Result, w)` | GET Location через `transport.Client.Download` | `int64` |

- `type` — строка, форма `^[a-z][a-z0-9]*$` проверяется до HTTP (как в `metadata`); допустимость
  типа решает сервер. Шаблоны этикеток и ценников — `type = "assortment"`; этикетки печатаются для
  product, service, bundle, variant. `id` — UUID любой версии.
- Шаблоны только читаются: изменить или удалить их через JSON API нельзя. Скрытые формы сервер не
  отдаёт. Параметров списка (limit/offset) источник не описывает — List читает одну страницу без них.
  Успех Get — HTTP 200 и непустой id, List — meta.href и массив rows.

## Шаблон (`Template`)

| JSON | Go type | Запись |
| --- | --- | --- |
| meta | `model.Meta` | нет |
| id | `string` (UUID) | нет |
| name | `string` | нет |
| type | `string`, открытый: `entity` — документ, `mxtemplate` — ценники и этикетки | нет |
| content | `string` (URL скачивания шаблона) | нет |

Ссылка на шаблон в запросах — `model.Reference{Meta: template.Meta}`.

## Запросы печати

| JSON | Go type | Обязательно |
| --- | --- | --- |
| `ExportRequest.template` | `model.Reference` | да, meta.href |
| `ExportRequest.extension` | `string` (xls, pdf, html, ods — проверяет сервер) | да |
| `ExportSetRequest.templates[]` | `[]TemplateCopies` | непустой |
| `TemplateCopies.template` | `model.Reference` | да |
| `TemplateCopies.count` | `int`, 1…10 проверяет сервер | ≥ 1 |
| `LabelRequest.organization` | `model.Reference` | да |
| `LabelRequest.count` | `int`, максимум 1000 проверяет сервер | ≥ 1 |
| `LabelRequest.salePrice.priceType` | `model.Reference` | да |
| `LabelRequest.template` | `model.Reference` | да |

- Одиночная печать и комплект — разные типы запроса: источник запрещает template вместе с
  templates, комплект печатается в pdf без extension. Для части документов в комплект допустимы
  шаблоны связанных документов (отгрузка — счёт, заказ, счёт-фактура и др.) — проверяет сервер.
- Таблицы источника помечают поля запроса «Обязательное при ответе»; SDK считает их обязательными
  в запросе: пустые отклоняются до HTTP. Верхние границы count — бизнес-ограничения, их проверяет API.

## Результат печати и загрузка

Печать выполняется `transport.Client.SendNoRedirect` ([transport.md](transport.md)): одна отправка,
303 не превращается в скрытый GET.

- `Result{StatusCode: 303, Location}` — форма готова; Location — временный адрес файла (5 минут),
  разрешённый относительно адреса запроса.
- `Result{StatusCode: 202, Location}` — адрес статуса печати. Источник называет этот вариант
  будущим и протокол опроса не описывает: SDK не опрашивает и не ждёт (Wait нет).
- Другой статус (200, 302, …), 303/202 без Location или с неразбираемым Location — ошибка
  `moysklad:`; query ссылки в текст ошибки не попадает. Ошибки API — `*sdkerrors.APIError`.
- `Download(ctx, result, w)` — отдельный GET только для 303. Origin ссылки на файл источник не
  фиксирует, поэтому Location не считается заведомо поддержанным: Download сохраняет same-origin
  контракт `transport.Client.Download` и отклоняет ссылку чужого origin до отправки. Такой файл
  вызывающий код скачивает своим HTTP-клиентом без Authorization.

## Противоречия источника

- Пример «Список стандартных шаблонов» (demand): meta.size 6 при 7 строках rows.
- Пример «Список пользовательских ценников и этикеток» побайтно повторяет список стандартных
  (meta.type embeddedtemplate) — фикстурой не взят.
- Атрибуты запроса печати и этикеток помечены «Обязательное при ответе», хотя это поля запроса.
- `template` этикеток помечен `+Expand`, но это поле запроса; expand в SDK нет.
- Ответы 202/303 описаны только заголовками (Location, Content-Type: application/json), тела нет.

Фикстуры — `testdata/custom/print` (README с SHA-256), тела запросов — `*.golden.json`.
