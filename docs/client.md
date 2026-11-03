# Фасад SDK и примеры

## moysklad.NewClient

```go
auth, err := transport.BearerToken(os.Getenv("MOYSKLAD_TOKEN"))
c, err := moysklad.NewClient(moysklad.Config{Auth: auth})

order, err := c.Entity().CustomerOrder().Get(ctx, id)        // entity/customerorder/{id}
page, err := c.Entity().CustomerOrder().GetPositions(ctx, id, positions.PageOptions{Limit: 100})
me, err := c.Context().Employee(ctx)                           // context/employee
report, err := c.Report().Turnover().All(ctx, turnover.AllOptions{...})
```

`moysklad.Config` — это `transport.Config`. NewClient создаёт один `transport.Client` и больше ничего
не делает: HTTP нет, токен не выпускается, вид токена (сотрудник, решение, черновик) не определяется.
Ошибка Config возвращается как есть.

Группы повторяют путь API:

| Вызов | Что возвращает |
| --- | --- |
| `c.Entity().<Сущность>()` | сервис пакета `entity/<сущность>`: Product, CustomerOrder, Assortment, … UOM; `Service()` возвращает `*service.Client` |
| `c.Entity().Metadata()` | `metadata.Service`: метаданные сущностей и статусы (CreateState/UpdateState/DeleteState) |
| `c.Entity().Files()` | `files.Service`: файлы объектов |
| `c.Report().Turnover()` | `turnover.Service` |
| `c.Context()` | `contexts.Service`: `Employee(ctx)`, `Application(ctx)` |
| `c.Audit()`, `c.Security()` | `audit.Service`, `security.Service` |
| `c.Async()` | `async.Service`: асинхронные задачи ([write-async.md](write-async.md)) |
| `c.Wizard()` | `wizard.Service`: автозаполнение документов ([documents-extra.md](documents-extra.md)) |
| `c.Transport()` | общий `transport.Client`: Download, HeaderValues, standalone сервисы |

Все handles работают через один transport: общие HTTPClient, Auth, Header, Limiter и Retry.
Handle не хранит собственного состояния и создаётся при каждом вызове без HTTP; два вызова дают
разные указатели с одинаковым поведением, поэтому handle можно сохранить в переменную. Глобального
состояния нет: два клиента с разными токенами не влияют друг на друга.

Client создаётся только через NewClient. У нулевого значения нет transport, и получение handle
паникует сообщением NewService.

Прежний путь работает: `product.NewService(c.Transport())` даёт сервис на том же transport. Методы,
опции вызова и результаты сервисов не зависят от способа получения.

## Примеры

Примеры — программы в `examples/`. `main` читает окружение, сценарий — функция с
`*moysklad.Client` и `io.Writer`, её тест идёт через настоящий transport и httptest, без токенов
и сети.

| Пример | Переменные | Что показывает |
| --- | --- | --- |
| [getting_started](../examples/getting_started/main.go) | `MOYSKLAD_TOKEN`, `MOYSKLAD_APP=1` для токена решения | клиент по токену, context/employee или context/application |
| [assortment_search](../examples/assortment_search/main.go) | `MOYSKLAD_SEARCH` | поиск и обход всех страниц через `Assortment().All` |
| [customerorder_positions](../examples/customerorder_positions/main.go) | `MOYSKLAD_ORDER_ID` | Get заказа и страницы GetPositions |
| [audit_events](../examples/audit_events/main.go) | — | страница контекстов аудита и события каждого контекста |
| [turnover](../examples/turnover/main.go) | — | «Обороты» за явный период, все страницы |
| [files](../examples/files/main.go) | `MOYSKLAD_ENTITY_TYPE`, `MOYSKLAD_ENTITY_ID`; загрузка — `MOYSKLAD_UPLOAD_PATH` и `MOYSKLAD_ALLOW_WRITE=1` | список, скачивание без перезаписи и с удалением частичного файла, загрузка |

Общие переменные: `MOYSKLAD_TOKEN` — персональный токен сотрудника или токен решения (оба Bearer),
`MOYSKLAD_BASE_URL` — другой адрес API. Выпуск токена (`c.Security().CreateToken`) в примеры
не входит. Запуск: `go run ./examples/<пример>`; Ctrl+C отменяет запрос.

Обход страниц — `pages.All(ctx, c.Transport(), first, opts...)`: каждая следующая страница —
отдельный GET по nextHref. Опции первой страницы на следующие не переносятся. Ассортимент
обходи через `c.Entity().Assortment().All(ctx, options)`: он сам передаёт beta-тег каждой странице.

## Проверки выпуска

`make check`: gofmt, `go vet ./... ./testdata`, `go test -race -count=1 -p 2 ./...` (примеры,
SHA-256 фикстур, статус сущностей) и `git diff --check` всего дерева и незакоммиченных изменений.
Сеть и токены не нужны. Тестовые HTTP-серверы — `internal/memnet.NewServer`: тот же
`*httptest.Server`, но соединения в памяти. Поток loopback TCP-соединений на macOS изредка давал
«connect: operation timed out»; `httptest.NewServer` допустим только с пометкой `// TCP:` строкой
выше (проверяет `TestNoTCPServers`). CI workflow не настроен (решение для v1); тег и публикация — отдельной командой.
Изменения версии — [CHANGELOG](../CHANGELOG.md).
