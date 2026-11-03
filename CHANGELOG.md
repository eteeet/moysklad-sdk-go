# Changelog

## Следующая версия — не выпущено

### Добавлено

- Справочники: склады (зоны и ячейки), группы товаров, отделы, проекты, договоры, каналы продаж,
  подразделения юрлиц, пользовательские справочники и элементы; только чтение — регионы, площадки
  продаж, серийные номера. Фасад: `c.Entity().Store()` и т.д.
- Контрагент: счета, контактные лица, события, настройки справочника, discounts (чтение),
  `mod__requisites__uz/kz`. Юрлицо: счета, реквизиты Узбекистана и Казахстана, `companyVat__ru`,
  `advancePaymentVat`, подписи и печать.
- Настройки справочника товаров, компании и пользователя, подписка компании (`contexts`).
- Активация, деактивация и сброс пароля сотрудника; мутации компонентов комплекта; кассиры точки
  продаж; CRUD описаний доп. полей; массовые операции статусов и характеристик.
- Подробности: [docs/dictionaries.md](docs/dictionaries.md).
- Товар: штрихкоды, упаковки, алкоголь, маркировка и остальные поля, неснижаемый остаток
  (`fields=minimumStock`, storebalances у товара и модификации); изображения (`images`), карточки
  контента, партии, транспортные упаковки с позициями, скидки всех видов и округление копеек, коды
  маркировки позиций (`trackingcode`). Подробности: [docs/catalog-extra.md](docs/catalog-extra.md).

## v1.0.0 — предложение

Не выпущено: тег не поставлен, публикация — отдельной командой. v1 — выбранный профиль SDK, а не весь
JSON API 1.2. Покрытие по методам: [status-entities](docs/status-entities.md),
[status-technical](docs/status-technical.md); ограничения — в guide каждого семейства.

### Добавлено

- Фасад `moysklad.NewClient(Config)`: группы `Entity()`, `Report()`, `Context()`, `Audit()`,
  `Security()` над одним `transport.Client`; конструктор не делает HTTP.
- Transport: per-call `RequestOption` и `ResponseInfo`, заголовки МойСклад, лимиты запроса,
  общий Limiter, повторы чтений по `Config.Retry`, `Download` в `io.Writer`, коды ошибок API,
  выпуск токена `security/token`.
- Query: filter/order/search/fields, страницы `query.PageOptions`, явный обход `pages.All`;
  у ассортимента — `Assortment().All` с beta-тегом на каждой странице.
- List у справочников и документов; новые справочники Currency, PriceType, TaxRate, Service,
  характеристики, Bundle, Variant, Assortment, Employee, Role, RetailStore, BonusProgram,
  BonusTransaction.
- Документы CustomerOrder, Enter, Loss, SalesReturn, RetailShift, RetailDemand, RetailSalesReturn;
  общая форма позиций с pack, things, gtd, trackingCodes, stock и declaration.
- Контексты сотрудника и решения, метаданные сущностей и статусы, аудит (контексты и события),
  файлы (List/Upload/Download), отчёт «Обороты» (All/ByStore/ByOperations).
- Delete у сервисов `entity/*` (кроме assortment и pricetype), `metadata.DeleteState`,
  `files.Delete`; массовые CreateMany/UpdateMany/DeleteMany (пакет `batch`, до 1000 элементов,
  без разбиения и повторов); асинхронные задачи — пакет `async`, `c.Async()` и `ListAsync`
  контрагентов и ассортимента; `transport.Client.Href`.
- Примеры сценариев в `examples/`, `make check`, проверка SHA-256 фикстур.

### Изменено — совместимость с 7442803

- Get/Create/Update/GetPositions прежних 17 сервисов и `positions.Service.GetPage` получили последний
  параметр `opts ...transport.RequestOption`. Обычные вызовы компилируются без изменений. Ломаются
  method value и интерфейсы со старой сигнатурой:
  `var f func(context.Context, string) (*uom.UOM, error) = svc.Get` — нужна обёртка
  ([ADR 0001](docs/adr/0001-request-options.md), `entity/uom/example_test.go`).
- `positions.PageOptions` стал alias `query.PageOptions` с новым полем `Fields`; литералы с именами
  полей компилируются, литерал без имён `positions.PageOptions{10, 0}` — нет (теперь три поля).
- Удалённых экспортированных имён нет. Прежние `NewService(*transport.Client)` работают и делят
  transport с фасадом через `c.Transport()`.

### Не входит в v1

Expand, корзина документов, операции позиций, async-отчёты, остальные сущности и отчёты, полный
Vendor API. Проверки идут
на httptest и примерах из документации; прогона на живом API нет.
