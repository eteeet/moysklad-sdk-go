# Сотрудники, права, метаданные и аудит

Источник контрактов — `api-remap-1.2-doc` revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`.
Все методы принимают `opts ...transport.RequestOption` ([ADR 0001](adr/0001-request-options.md)).

## Сотрудники и права

| Пакет / метод | Маршрут |
| --- | --- |
| `entity/employee`: `Get`, `List`, `Create`, `Update` | `entity/employee[/{id}]` |
| `employee.Service.GetSecurity`, `UpdateSecurity` | `entity/employee/{id}/security` |
| `entity/role`: `Get`, `List`, `Create`, `Update` | `entity/role[/{id}]` |
| `role.Service.Preset(ctx, role.PresetAdmin)` | `entity/role/{owner,admin,individual,cashier,worker,…}` |

- Employee: `lastName` обязателен и в Create, и в Update. salary, owner, group, archived и email меняет
  только администратор — права проверяет сервер. Фото не входит; активация, деактивация и сброс
  пароля — [dictionaries.md](dictionaries.md).
- Security: роль задаётся по meta (`SecurityRoleRequest`); для индивидуальной роли непереданные права
  сервер выключает (кроме view=ALL у currency, country, uom), без `permissions` — сохраняет прежние.
  login и isActive меняются активацией (`Activate`, `Deactivate`), email — через `Employee.Update`.
- Права (`model.Permissions`) — объект «имя → значение»: флаг (`importData: true`), области по действиям
  (`supply: {"view":"ALL"}`, у задач ключ `script` и действие `done`) или иная форма в `Raw` без потерь.
  Таблицы источника называют поле Array(Object), но все примеры — объект.

## Контекст запроса

`contexts.Service.Employee(ctx)` — `GET context/employee`: сотрудник, от имени которого авторизован
клиент, с `code` и правами. `Application(ctx)` — `GET context/application`: id установки и appUid
решения по его токену; это не Vendor API. NewClient эти запросы не делает.

## Метаданные и статусы

| Метод `metadata.Service` | Маршрут |
| --- | --- |
| `Catalog(types)` | `entity/metadata[?filter=type=a;type=b]` |
| `Get(type)` | `entity/{type}/metadata` |
| `GetCustomEntity(id)`, `GetCompanySettings` | `entity/customentity/{id}/metadata`, `context/companysettings/metadata` |
| `ListAttributes(type, page)`, `GetAttribute(type, id)` | `entity/{type}/metadata/attributes[/{id}]` |
| `CreateState(type, …)`, `UpdateState(type, id, …)` | `entity/{type}/metadata/states[/{id}]` |

- `attributes` в метаданных приходит массивом или коллекцией `{meta[, rows]}` — читаются обе формы;
  весь ответ хранится в `Raw`. Создание и изменение описаний доп. полей не входит.
- Статусы читаются в `Get(type).States`; отдельного GET статуса API нет. Create требует name, color
  (ARGB целым, `*int64` — 0 допустим) и stateType. DeleteState удаляет статус
  ([docs/write-async.md](write-async.md)); массовое создание статусов не входит.

## Аудит

`audit.Service`: `ListContexts(ContextsOptions)` — одна страница `GET audit` (limit до 100, по умолчанию
сервера 25, filter по moment, employee, uid, source, eventType, entityType, application);
`GetFilters` — каталог фильтров; `EventsByContext(id)` и `EventsByEntity(type, id)` — события.

- Событие хранит `Diff` и всё событие (`Raw`) без потерь; `FieldChanges()` даёт «атрибут → oldValue/newValue»,
  `DecodeDiff` — диффы регистрации, публикации и писем. Старые значения не раскладываются в текущие модели.
- Параметры пагинации событий источник не описывает: SDK их не передаёт; nextHref можно пройти `pages.All`.
- Аудит — не changelog с курсором: нет гарантий хранения, exactly-once и порядка при равных moment.
- Противоречия источника записаны в godoc: events/info у контекста, регистр ключей каталога фильтров.
