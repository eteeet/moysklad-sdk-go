# NamedFilter

Пакет `namedfilter` читает сохранённые фильтры списков: `entity/{entityType}/namedfilter`.
`namedfilter.NewService(client)` принимает `*transport.Client`, возвращает
`(*namedfilter.Service, error)` и отклоняет nil client. Пакет только читает: создание,
изменение и удаление фильтров источник не описывает.

| Метод | Запрос | Результат |
| --- | --- | --- |
| `List(ctx, entityType, namedfilter.ListOptions{})` | GET `entity/{entityType}/namedfilter` | `*model.Collection[namedfilter.NamedFilter]` |
| `List(ctx, entityType, namedfilter.ListOptions{Owner: href})` | GET `entity/{entityType}/namedfilter?owner={href}` | то же |
| `Get(ctx, entityType, id)` | GET `entity/{entityType}/namedfilter/{id}` | `*namedfilter.NamedFilter` |

`entityType` — строка (`"product"`, `"supply"`, `"paymentin"` …): до HTTP проверяется только
форма `[a-z][a-z0-9]*`, поддерживает ли тип фильтры — решает сервер. Общего списка фильтров
по всем типам API не даёт. `id` проверяется как UUID любой версии. `Owner` — абсолютный href
сотрудника, иначе ошибка до HTTP; пустой `Owner` — фильтры текущего пользователя. Чужие
фильтры доступны пользователю с правами администратора или решению; права проверяет сервер.

Каждый вызов — один GET без повторов и автообхода страниц. List требует HTTP 200, непустой
`meta.href` и массив `rows`; Get — HTTP 200 и непустой `id`. При ошибке результат nil,
`errors.As` распознаёт `*sdkerrors.APIError`, `errors.Is` сохраняет ошибки context,
авторизации и транспорта; `transport.WithResponseInfo` сбрасывается до проверок аргументов.

## Поля

| JSON | Go | Запись |
| --- | --- | --- |
| `meta` | `Meta model.Meta` | нет |
| `id` | `ID string` | нет |
| `accountId` | `AccountID string` | нет |
| `name` | `Name string` | нет |
| `owner` | `Owner *model.Reference` (сотрудник, без expand) | нет |

Параметры фильтрации API не возвращает. Коллекция — общий `model.Collection` с `context`.

## Применение фильтра

Применение уже есть в списках: `query.ListOptions{NamedFilter: f.Meta.Href}` передаёт
`namedfilter={href}` в List сущности или документа того же типа; проверка — только абсолютный href (`query/options.go`).
Ограничения источника SDK не проверяет локально, их применяет сервер: один фильтр на запрос;
несовместим с `filter`, `search`, `order`; не применяется к `assortment`; на endpoint
применяется только фильтр, полученный через `entity/{entityType}/namedfilter` этого типа;
чужие фильтры — только администратор или решение.

## Ограничения и противоречия источника

- `name` помечен «Необходимо при создании», но создание фильтров не описано: поле только читается.
- `meta` не помечен «Только для чтения», но записи нет вовсе.
- Ответ списка содержит `meta.limit`/`meta.offset`, но параметры `limit`/`offset` и другие
  (кроме `owner`) для этого endpoint не описаны — в `ListOptions` их нет.
- `owner` помечен `+Expand`; expand в SDK не поддерживается.

Контракт сверён с [официальной документацией на revision
0e1b3d61fb8d8fa314469361d96560d6c3f7f57b](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_named_filter.md).
