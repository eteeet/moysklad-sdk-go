# WebhookStock

`entity/webhookstock`: модель `WebhookStock`, `CreateRequest`, `UpdateRequest`, `NewService(client)`;
методы `Get/List/Create/Update/Delete` и `CreateMany/UpdateMany/DeleteMany`.
[Общие правила методов, JSON и ошибок](documents.md).

- `Get(ctx, id)`, `Create(ctx, CreateRequest)`, `Update(ctx, id, UpdateRequest)` →
  `(*WebhookStock, error)`; успех — HTTP 200 и непустой `id`, возвращается ответ сервера.
- `List(ctx, query.ListOptions)` — одна страница `GET entity/webhookstock` без expand и автообхода.
  Фильтры и пагинацию источник не описывает: SDK передаёт заданные параметры как есть.
- `Delete(ctx, id)` — `DELETE entity/webhookstock/{id}`, безвозвратно.
- `CreateMany`, `UpdateMany` — один `POST entity/webhookstock` массивом (до `batch.MaxItems`),
  `DeleteMany` — `POST entity/webhookstock/delete`; ответ по элементам по правилам пакета `batch`.
- id проверяется как UUID до запроса; пустой `UpdateRequest` отклоняется до запроса; записи не повторяются.
- Отключение — `Update(ctx, id, webhookstock.UpdateRequest{Enabled: model.Value(false)})`.

Минимальный Create: непустые `url`, `stockType`, `reportType` (проверяются до запроса, без trim).

```go
hook, err := s.Create(ctx, webhookstock.CreateRequest{
	URL:        "https://example.com/hooks/stock",
	StockType:  "stock",
	ReportType: "all",
	Enabled:    model.Value(true),
})
```

## Поля

| JSON | Read Go type | Create | Update |
| --- | --- | --- | --- |
| `meta` | `model.Meta` | — | — |
| `id` | `string` | — | — |
| `accountId` | `string` | — | — |
| `stockType` | `string` (открытая) | required | Field |
| `reportType` | `string` (открытая) | required | Field |
| `url` | `string` | required | Field |
| `enabled` | `bool` | Field | Field |
| `authorApplication` | `*model.Reference` | — (Только для чтения) | — |

Field — `model.Field[T]` с `omitzero`: не задано — поле не уходит; `model.Value(false)` уходит
явным `false`, `model.Null[T]()` — `null`.

## Ограничения и противоречия

- `stockType` (`stock`) и `reportType` (`all`, `bystore`) — открытые строки: SDK не
  отвергает значения вне перечислений ни в ответе, ни в запросе.
- `url` не проверяется (схема, длина до 255, доступность); уникальность сочетания
  `stockType`+`reportType`+`url`, лимит 5 вебхуков (1 для решения) и тариф проверяет сервер.
- `enabled` не помечен «Необходимо при создании», но передаётся в примере создания; значение по
  умолчанию источник не называет — SDK отправляет его только если задано.
- Пример массового создания в источнике содержит висячие запятые (невалидный JSON); тела
  `CreateMany`/`UpdateMany` SDK формирует по таблице атрибутов.
- `authorApplication` в таблице — `Meta`, в примерах — объект `{"meta": …}`: читается как
  `*model.Reference`.
- Тело входящего уведомления (`reportUrl` и др.) — вне SDK клиента; сам отчёт по `reportUrl` —
  краткий отчёт об остатках (`report/stock/all/current`).
- Expand и metadata в пакете нет; ticket их не называет.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhookstock),
`md/dictionaries/_webhookstock.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
