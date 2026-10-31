# Counterparty

Пакет `entity/counterparty` поддерживает получение, создание и частичное изменение
Counterparty (контрагента) на перечисленном ниже наборе полей.

| Метод | Запрос относительно `/api/remap/1.2/` |
| --- | --- |
| `Get(ctx, id)` | `GET entity/counterparty/{id}` |
| `Create(ctx, CreateRequest)` | `POST entity/counterparty` |
| `Update(ctx, id, UpdateRequest)` | `PUT entity/counterparty/{id}` |

`NewService` принимает существующий `*transport.Client`; `nil` отклоняется.
Все методы возвращают `(*Counterparty, error)`. Успех требует HTTP 200 и непустой
`id` в ответе. Create/Update возвращают данные сервера, а не входной request.
ID ответа дополнительно не сравнивается с ID маршрута.

## Поля

| JSON-поля | Тип чтения | Запись |
| --- | --- | --- |
| `meta` | `model.Meta` | нет |
| `id`, `accountId` | `string` | нет |
| `created`, `updated` | `model.DateTime` | нет |
| `name`, `description`, `code`, `externalCode` | `string` | Create, Update |
| `archived`, `shared` | `bool` | Create, Update |
| `companyType` | `model.CompanyType` | Create, Update |
| `email`, `phone`, `fax` | `string` | Create, Update |
| `legalTitle`, `legalFirstName`, `legalLastName`, `legalMiddleName` | `string` | Create, Update |
| `inn`, `kpp`, `ogrn`, `ogrnip`, `okpo`, `certificateNumber` | `string` | Create, Update |
| `certificateDate`, `birthDate` | `model.DateTime` | Create, Update |
| `sex`, `discountCardNumber` | `string` | Create, Update |
| `actualAddress`, `legalAddress` | `string` | Create, Update |
| `actualAddressFull`, `legalAddressFull` | `*model.Address` | Create, Update через `model.AddressRequest` |
| `attributes` | `[]model.Attribute` | Create, Update |
| `tags` | `[]string` | Create, Update |
| `owner`, `group`, `bonusProgram`, `state` | `*model.Reference` | Create, Update |
| `priceType` | `*model.PriceType` | Create, Update |
| `syncId` | `string` | Create, Update |
| `bonusPoints` | `int` | нет |
| `salesAmount` | `float64` | нет |

Реквизиты и коды сохраняются строками, включая ведущие нули. Отсутствующие или
неприменимые поля допускаются; SDK не выводит и не очищает реквизиты по `companyType`.
`CompanyType` имеет открытый набор значений: неизвестные строки читаются и
записываются без проверки. Defaults, включая `companyType`, SDK не подставляет.
`sex` также передаётся без предметной проверки.

`salesAmount` принимает JSON-числа `0` и `0.0` и читается без пересчёта единиц.
`float64` не обещает точной десятичной финансовой арифметики. Отсутствующие поля
дают нулевые значения, отсутствующие/null-ссылки и адреса — `nil`, даты — нулевой
`model.DateTime`. Неизвестные и отложенные поля ответа игнорируются.

## Запросы

В `CreateRequest` обязателен `Name string`: локально отклоняется только `Name == ""`,
без trim или нормализации. Остальные записываемые поля — `model.Field[T]`.
В `UpdateRequest` все поля имеют тип `model.Field[T]`; полностью незаданный request
отклоняется. `SyncID` доступен в CreateRequest и UpdateRequest через Field[string].
После заполнения менять его нельзя; это проверяет сервер, без предварительного Get.
Пустая строка и null кодируются явно, но допустимость очистки или перезаписи не обещается.

`Field` с `omitzero` различает отсутствие, `model.Value(value)` и `model.Null[T]()`.
Явные `false`, `""`, `[]`, `null` и `model.Value(model.AddressRequest{})` отправляются
и считаются заданными. Значение nil-слайса кодируется как `[]`; для JSON `null`
используй `model.Null`. Допустимость значений и сочетаний проверяет сервер.

```go
service, err := counterparty.NewService(client)
if err != nil {
	return err
}

created, err := service.Create(ctx, counterparty.CreateRequest{Name: "Новый контрагент"})
if err != nil {
	return err
}

updated, err := service.Update(ctx, created.ID, counterparty.UpdateRequest{
	Shared: model.Value(false),
	ActualAddressFull: model.Value(model.AddressRequest{
		City:       model.Value("Москва"),
		PostalCode: model.Value("001234"),
	}),
})
```

Обе формы адреса уходят в JSON, если заданы обе. SDK не конкатенирует адрес,
не копирует строку в `addInfo` и не сливает его части. По документации сервер
отдаёт приоритет структуре; при передаче только строки отражает её также в `addInfo`.
`null` не очищает адрес: для очистки передай `""` в соответствующее строковое поле,
например `ActualAddress: model.Value("")`. Подробности общих типов — в
[requisites.md](requisites.md).

`attributes` использует существующие `model.Attribute`, `AttrValue`, `AttrNull`,
`AttrFile` и `AttrFileNull`. Сервер обновляет значения по метаданным поля;
`[]` не очищает всю коллекцию. Сброс значения — `value:null`, файлового — `file:null`.
SDK не сливает attributes или tags локально. Для tags не обещаются непроверенные
серверные правила очистки, замены или слияния коллекции.

## Ошибки и проверки

Get/Update проверяют форму UUID (Universally Unique Identifier), без ограничения
версии. Методы используют context, авторизацию, заголовки и gzip существующего
transport. При любой ошибке возвращается nil-модель; обёртка сохраняет
`errors.As` для `*sdkerrors.APIError` (пакет SDK `errors`) и `errors.Is` для причин
ошибок транспорта/context.
Автоматических повторов, предварительного Get и запросов по ссылкам нет.

Тесты обоими кодеками `encoding/json` и `encoding/json/v2` проверяют значения
неизменённых upstream Counterparty/CounterpartyLegal, дополнительные optional/null
и неизвестные значения, адреса, даты и независимые golden запросов.
Локальный HTTP-сервер с настоящим transport проверяет маршруты, тела, заголовки,
plain/gzip, серверные ответы, context и ошибки. Живой API не использовался;
серверные правила очистки адресов и коллекций не имитировались.

## TODO

Не реализованы accounts (включая модели и metadata), contactpersons, notes,
files, discounts, региональные `mod__requisites__uz`/`mod__requisites__kz`,
их metadata и операции. Также отсутствуют List/search/filter/sort/pagination,
Delete, expand, batch и управление типами цен. Наличие Reference не означает
поддержку операций связанной сущности.

Контракт сверён с официальной документацией на revision
`9a811cef3dc4a3f3a46f070e60452361a06bf7cd`:
[Counterparty](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_counterparty.md),
[attributes](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L314-L367),
[ограничение null адреса](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187).
