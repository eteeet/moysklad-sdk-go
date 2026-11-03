# ADR 0001: per-call RequestOption и ResponseInfo

Статус: принято 2026-11-01 (T01, выбор пользователя).

## Контекст

Services теряли headers и status ответа, а per-call headers (регион, язык,
отключение webhooks) передать было нельзя. Решение нужно до retrofit 17 services.

## Решение

- Последний параметр каждого публичного метода — `opts ...transport.RequestOption`.
  Результаты `(*T, error)` не меняются.
- `transport.WithHeader(key, value)` добавляет значение; повтор key даёт несколько
  header values, не CSV.
- `transport.WithResponseInfo(&info)` заполняет `StatusCode` и копию всех headers
  последнего полученного ответа, в том числе при `APIError` и ошибке decode.
  Если ответа нет, `info` обнуляется. Указатель принадлежит вызову; общего
  `LastResponse` нет. Публичная операция обнуляет `info` на входе через
  `transport.ResetResponseInfo(opts)`, до локальной валидации: ошибка до запроса не оставляет
  метаданных прошлого вызова при любой композиции опций.
- Client defaults — `transport.Config.Header`. Приоритет: значения transport
  < `Config.Header` < опции вызова; key из опций заменяет key клиента целиком.
- `Authorization`, `Accept`, `Accept-Encoding`, `Content-Type`, `Content-Length`, `Host`
  задаёт transport. Попытка задать их в Config или опциях — ошибка до отправки.
- `transport.WithAuth(auth)` заменяет `Config.Auth` на один вызов (T04: Basic для
  `security/token` на клиенте с Bearer); заголовок `Authorization` по-прежнему задаёт transport.
- `transport.Client.Send` — seam JSON-вызова: сборка запроса, Do и ResponseInfo.
  Services вызывают только его и его обёртку `SendNoRedirect` (печать: один вызов без следования
  303/202, тело успеха не декодируется; тот же limiter, retry и ResponseInfo, что у Send). Limiter и retry (T08/T09) встраиваются сюда, а не
  в services. Download (F01) не может идти через Send, потому что Do закрывает body:
  это отдельный метод Client с тем же limiter и retry. Решение F01 — writer, а не
  ReadCloser: тело закрывается внутри метода на всех путях, слоты limiter держатся ровно
  до конца копирования, вызывающий код ничего не закрывает. Retry повторяет попытку
  только пока в writer ничего не записано. Исключение из «последнего полученного ответа»: при
  неразбираемом Location net/http ответ не отдаёт, ResponseInfo Download остаётся без статуса;
  обёртку Transport ради этого случая не вводим.
- Policy принадлежит `*transport.Client`: services и будущий фасад, собранные на
  одном Client, делят его настройки, limiter и retry.
- `NewRequest`/`Do` остаются низкоуровневым API без опций; ResponseInfo они не заполняют.
- T09 добавит в ResponseInfo число попыток; StatusCode и Header всегда относятся к
  последней попытке.

## Прототипы для следующих tickets

Компилируются сейчас только методы UOM. Остальные строки фиксируют форму методов,
которых ещё нет; их tickets добавляют компилируемые примеры.

```go
Get(ctx, id, opts...) (*T, error)
Create(ctx, input, opts...) (*T, error)
Update(ctx, id, input, opts...) (*T, error)
List(ctx, query, opts...) (*model.Collection[T], error)        // T06/T11/T12
GetPositions(ctx, id, page, opts...) (*model.Collection[P], error)
Report(ctx, query, opts...) (*R, error)                         // форма R своя у каждого отчёта
Download(ctx, href, w io.Writer, opts...) (int64, error)       // F01: тело копируется в w
```

Query-параметры остаются отдельным аргументом конкретного метода, а не RequestOption:
supported filter/order/fields различаются по endpoint.

## Совместимость

Обычные вызовы `svc.Get(ctx, id)` компилируются без изменений. Сигнатура
`Client.NewRequest` не менялась. Ломается тип method
value: `var f func(context.Context, string) (*uom.UOM, error) = svc.Get` больше не
компилируется, нужна обёртка (см. `entity/uom/example_test.go`).

## Отвергнутая альтернатива: `Result[T]`

`Result[T]{Value, Info}` делает metadata всегда доступной, но меняет сигнатуры всех
существующих методов и заставляет каждого caller разворачивать `.Value` ради
редко нужных headers. Выбран вариант без изменения результатов.
