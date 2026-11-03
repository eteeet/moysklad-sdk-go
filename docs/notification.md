# Notification

Пакет `notification` — лента уведомлений текущего пользователя и настройки уведомлений.
Маршруты `notification/...` (не `entity/`). `notification.NewService(client)` принимает
`*transport.Client`, возвращает `(*notification.Service, error)` и отклоняет nil client.

## API

| Метод | HTTP | Результат |
| --- | --- | --- |
| `List(ctx, notification.ListOptions)` | GET `notification` | `(*model.Collection[notification.Notification], error)` |
| `Get(ctx, id)` | GET `notification/{id}` | `(*notification.Notification, error)` |
| `Delete(ctx, id)` | DELETE `notification/{id}` | `error` |
| `MarkAsRead(ctx, id)` | PUT `notification/{id}/markasread` | `error` |
| `MarkAllAsRead(ctx)` | PUT `notification/markasreadall` | `error` |
| `Settings(ctx)` | GET `notification/settings` | `(notification.Settings, error)` |
| `UpdateSettings(ctx, notification.SettingsRequest)` | PUT `notification/settings` | `error` |

Все методы принимают `...transport.RequestOption`, первой строкой сбрасывают `ResponseInfo` и
выполняют ровно одну отправку; PUT и DELETE не повторяются и при `RetryPolicy` клиента
(GET повторяется, только если вызывающий включил `RetryPolicy`). Id проверяется на форму UUID
до HTTP. Ошибка содержит один префикс `moysklad:` и операцию (`get notification {id}`,
`mark all notifications as read`, `update notification settings` …); `errors.As` находит
`*sdkerrors.APIError`, `errors.Is` сохраняет ошибки context, авторизации и транспорта.
Права и принадлежность уведомления текущему пользователю проверяет сервер.

## Лента

`ListOptions`: `Limit` (0 — default сервера 1000, иначе 1…1000), `Offset` (≥ 0). Других
параметров источник не описывает. Одна страница без автообхода; следующие — `pages.All`.
Успех — HTTP 200, непустой `meta.href` и массив `rows`. `Get` — HTTP 200 и непустой `id`.

`Delete`, `MarkAsRead`, `MarkAllAsRead`: источник тела ответа не описывает, SDK его не читает;
успех — HTTP 200 или 204.

### Поля `Notification` (только чтение)

| JSON | Go type | Примечание |
| --- | --- | --- |
| `meta` | `model.Meta` | `meta.type` — тип уведомления, открытая строка |
| `id` | `string` | |
| `accountId` | `string` | |
| `created` | `model.DateTime` | |
| `read` | `bool` | меняют `MarkAsRead`/`MarkAllAsRead` |
| `title` | `string` | |
| `description` | `string` | |
| — | `Raw jsontext.Value` | всё уведомление как пришло; при кодировании не пишется |

Типы уведомлений (30 в сводной таблице источника: `NotificationTaskChanged`,
`NotificationImportCompleted`, `RetireOrderByDemandNotificationCompleted`, `NewMentionInEvent`,
`SalesManagerChatMessage` …) не моделируются отдельными структурами: поля конкретного типа
(`task`, `performedBy`, `diff`, `good`, `actualBalance`, `taskType`, `createdDocsCount` …) читаются
из `Raw` в структуру вызывающего по `Meta.Type` — см. `ExampleNotification_raw`. Неизвестный тип
не отвергается. `Raw` хранит исходные байты, поэтому long-значения не проходят через float64.

## Настройки

`Settings` — `map[string]GroupSettings`, ключ — группа. `GroupSettings{Enabled bool;
ChannelsEnabled map[string]bool}`. Ключи групп (`customerOrder`, `dataExchange`, `invoice`,
`retail`, `scripts`, `stock`, `task`, `mentions`, `onlineStores`, `followedEvents`) и каналов
(`email`, `push`, `interface`) открытые: неизвестные читаются и пишутся как есть; отсутствие
канала и `false` различимы через `v, ok := ChannelsEnabled[k]`. Ответ без `meta`: ошибка, только
если тело не JSON-объект (`null`, массив, пустое тело); `{}` — пустые настройки.

| JSON | Go type (запрос `GroupSettingsRequest`) | Запись |
| --- | --- | --- |
| `enabled` | `model.Field[bool]`, `omitzero` | `model.Value(false)` отправляется явно |
| `channelsEnabled` | `map[string]bool`, `omitzero` | только перечисленные каналы; `false` — явное отключение |

`SettingsRequest` — `map[string]GroupSettingsRequest`; изменение частичное: не переданные группы,
флаги и каналы сервер не меняет. Пустой запрос и группа без полей отклоняются до HTTP.
`UpdateSettings` тело ответа не читает (источник его не описывает); актуальные значения — `Settings`.

## Ограничения и противоречия источника

- Источник: «Отключение уведомлений групп "Сценарии" и "Интернет-магазины" недопустимо.
  Параметр enabled игнорируется» — не ясно, игнорируется `enabled` только у этих групп или
  везде. SDK отправляет заданное значение, решает сервер.
- У полей настроек источник не ставит маркеров «обязательное»/«только для чтения» и не
  описывает null; `model.Null[bool]()` для `enabled` кодируется, но его принятие не обещано.
- `meta` и `read` в таблице общих атрибутов не помечены «Только для чтения», но операций записи
  уведомления источник не описывает: модель только читается.
- В примере `Get` у `performedBy` нет `meta`, в примере ленты — есть; SDK эти поля не разбирает (`Raw`).
- Тело ответов `markasread`, `markasreadall`, PUT `settings` не описано; SDK его не требует.
- Expand, фильтры ленты и прочие операции не входят.

Фикстуры: `testdata/custom/notification` (примеры `get.json`, `list.json` с SHA-256 и golden тела
PUT настроек), `testdata/upstream/notification_settings.json`,
`testdata/upstream/notification_task_changed.json`.

Контракт сверён с официальной документацией на revision 0e1b3d61fb8d8fa314469361d96560d6c3f7f57b:
[_notifications.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/notification/_notifications.md),
[_notification.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/notification/_notification.md),
[_notification_types.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/notification/_notification_types.md),
[_notification_types_details.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/notification/_notification_types_details.md),
[_notification_settings.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/notification/_notification_settings.md).
