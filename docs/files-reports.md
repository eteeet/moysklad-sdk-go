# Файлы и отчёт «Обороты»

Источник контрактов — `api-remap-1.2-doc` revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`.
Все методы принимают `opts ...transport.RequestOption` ([ADR 0001](adr/0001-request-options.md)).

## Скачивание

`transport.Client.Download(ctx, href, w io.Writer, opts...) (int64, error)` копирует тело ответа в `w`
и возвращает число записанных байт — и при ошибке.

- href — относительный endpoint или абсолютный URL того же origin, что BaseURL (`meta.downloadHref`);
  чужой origin отклоняется до отправки. Redirect на подписанный URL другого origin выполняется без
  Authorization; query подписанной ссылки в текст ошибки не попадает, в том числе из вложенной ошибки
  неразбираемого Location.
- Размер ответа не ограничен, общий `Config.Timeout` к телу не применяется: срок задаёт ctx.
  gzip распаковывается независимо от настроек HTTPClient.
- Слоты Limiter держатся до конца копирования. `Config.Retry` повторяет попытку, только пока в `w`
  ничего не записано; после частичной записи решает вызывающий код.
- ResponseInfo описывает последний ответ цепочки redirect, в том числе отклонённого своим
  CheckRedirect. При неразбираемом Location net/http ответ не возвращает: StatusCode 0.
  Время жизни ссылки — `transport.WithDownloadExpiration` (по умолчанию сервера 60 с, максимум 86400).
- Решение «writer, а не ReadCloser» записано в ADR 0001: тело всегда закрывается внутри метода.

## Файлы

| Метод `files.Service` | Endpoint |
| --- | --- |
| `List(ctx, type, id, query.PageOptions)` | GET entity/{type}/{id}/files — одна страница |
| `Upload(ctx, type, id, []files.Upload)` | POST entity/{type}/{id}/files — 1–10 файлов, Base64 |
| `Download(ctx, files.File, w)` | GET meta.downloadHref через `Client.Download` |

- Файлы — подресурс объекта (операции, номенклатуры, задачи, контрагента); отдельного CRUD файла нет.
- Upload добавляет файлы и возвращает весь список. Лимит 100 файлов на объект проверяет сервер, SDK не
  считает их заранее. Размер тела — закодированный JSON (`Config.MaxBodySize`); Base64 больше файла
  примерно на треть.
- Имя файла: непустое, без разделителей пути, NUL, «.» и «..»; длину и тип содержимого проверяет сервер.
- `Delete(ctx, type, id, fileID)` — [write-async.md](write-async.md). Изображения товаров и скачивание
  миниатюр не входят.

## Отчёт «Обороты»

| Метод `turnover.Service` | Endpoint | Строка |
| --- | --- | --- |
| `All(ctx, AllOptions)` | report/turnover/all | `Row`: assortment, onPeriodStart, onPeriodEnd, income, outcome |
| `ByStore(ctx, DetailOptions)` | report/turnover/bystore | `StoreRow`: assortment, stockByStore[] по складам |
| `ByOperations(ctx, DetailOptions)` | report/turnover/byoperations | `OperationRow`: assortment, store, operation, quantity, cost, sum |

- Показатели: `sum` и `cost` — себестоимость в копейках, `quantity` — количество.
- `AllOptions`: Limit/Offset (одна страница), период, GroupBy (`product` по умолчанию сервера,
  `variant`), Filter (`query.Filter`), Extra.
- `DetailOptions`: обязательный Assortment — его meta.type (product или variant) задаёт фильтр
  `product=` или `variant=`; Filter — остальные условия, без повтора product/variant. Limit, offset и
  groupBy у этих отчётов источник не описывает.
- Период по умолчанию: без границ — последний месяц, только MomentTo — с начала года, только MomentFrom —
  по текущий день. Для ByOperations таблица называет границы обязательными, а текст допускает
  умолчания — задавайте период явно.
- `query.ReportOptions` передаёт моменты как московское время; с `transport.WithAcceptTimezone` сервер
  считает период в зоне клиента и сообщает её в `ResponseInfo.ContentTimezone`. DateTime моделей от
  заголовка не зависит.
- Противоречия источника: stockByStore — Object в таблице и массив в ответах; moment операции — Object в
  таблице и дата в ответах. SDK следует ответам.
