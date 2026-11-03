# Webhook

`entity/webhook`: модель `Webhook`, `CreateRequest`, `UpdateRequest`, `NewService(client)`;
методы `Get/List/Create/Update/Delete` и `CreateMany/UpdateMany/DeleteMany`.
[Общие правила методов, JSON и ошибок](documents.md).

- `Get(ctx, id)`, `Create(ctx, CreateRequest)`, `Update(ctx, id, UpdateRequest)` →
  `(*Webhook, error)`; успех — HTTP 200 и непустой `id`, возвращается ответ сервера.
- `List(ctx, query.ListOptions)` — одна страница `GET entity/webhook` без expand и автообхода.
  Фильтры и пагинацию источник для вебхуков не описывает: SDK передаёт заданные параметры как есть.
- `Delete(ctx, id)` — `DELETE entity/webhook/{id}`, безвозвратно.
- `CreateMany`, `UpdateMany` — один `POST entity/webhook` массивом (до `batch.MaxItems`),
  `DeleteMany` — `POST entity/webhook/delete`; ответ по элементам по правилам пакета `batch`.
- id проверяется как UUID до запроса; пустой `UpdateRequest` отклоняется до запроса; записи не повторяются.
- Отключение — `Update(ctx, id, webhook.UpdateRequest{Enabled: model.Value(false)})`, тело `{"enabled":false}`.

Минимальный Create: непустые `url`, `action`, `entityType` (проверяются до запроса, без trim).

```go
hook, err := s.Create(ctx, webhook.CreateRequest{
	URL:        "https://example.com/hooks/moysklad",
	Action:     "UPDATE",
	EntityType: "supply",
	DiffType:   model.Value("FIELDS"),
})
```

## Поля

| JSON | Read Go type | Create | Update |
| --- | --- | --- | --- |
| `meta` | `model.Meta` | — | — |
| `id` | `string` | — | — |
| `accountId` | `string` | — | — |
| `entityType` | `string` (открытая) | required | Field |
| `url` | `string` | required | Field |
| `action` | `string` (открытая) | required | Field |
| `diffType` | `string` (открытая) | Field | Field |
| `enabled` | `bool` | Field | Field |
| `method` | `string` (открытая) | Field | Field |
| `authorApplication` | `*model.Reference` | — | — |

Field — `model.Field[T]` с `omitzero`: не задано — поле не уходит; `model.Value(false)` уходит
явным `false`, `model.Null[T]()` — `null`.

## Ограничения и противоречия

- `action` (`CREATE`, `UPDATE`, `DELETE`, `PROCESSED`), `diffType` (`NONE`, `FIELDS`),
  `entityType` и `method` (`POST`) — открытые строки: SDK не отвергает значения вне перечислений
  ни в ответе, ни в запросе. Связи «`PROCESSED` только для асинхронных задач», «`diffType` только
  для `UPDATE`» и запрет `entityType` `webhook`/`discount`/`notes` проверяет сервер.
- `url` не проверяется (ни схема, ни длина до 255, ни доступность); уникальность сочетания
  `entityType`+`action`+`url`, лимит 5 вебхуков (1 для решения), тариф и права
  администратора проверяет сервер.
- `method` и `authorApplication` в таблице атрибутов без маркера «Только для чтения», их нет ни в
  одном примере запроса. `method` (единственное значение `POST`) пишется в Create/Update — по
  правилу полей проекта. `authorApplication` — метаданные решения, создавшего вебхук (у
  webhookstock он помечен «Только для чтения»): его ставит сервер, SDK только читает.
- `authorApplication` в таблице — `Meta`, в примерах — объект `{"meta": …}`: читается как
  `*model.Reference`; дополнительные поля решения (upstream-фикстура SDK) не сохраняются.
- `enabled` в источнике не помечен «Необходимо при создании»; пример создания webhookstock
  передаёт его, пример webhook — нет. Значение по умолчанию источник не называет.
- Заголовок `X-Lognex-WebHook-DisableByPrefix` относится к любым запросам API, а не к ресурсу
  вебхука: уже есть `transport.WithWebHookDisableByPrefix(prefixes...)` (один заголовок на префикс).
- Приём входящих уведомлений (HTTP-сервер, тело `events`/`auditContext`) — вне SDK клиента.
- Expand и metadata в пакете нет; ticket их не называет.

Источник: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/dictionaries/webhook),
`md/dictionaries/_webhook.md`, revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
