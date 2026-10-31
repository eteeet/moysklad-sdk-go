# Organization (юрлицо)

Пакет `entity/organization` предоставляет частичную модель чтения `Organization`,
`CreateRequest`, `UpdateRequest` и `Service`. `NewService` принимает существующий
`*transport.Client`; `nil` отклоняется.

| Метод | HTTP-запрос |
| --- | --- |
| `Get(ctx, id)` | `GET entity/organization/{id}` |
| `Create(ctx, input)` | `POST entity/organization` |
| `Update(ctx, id, input)` | `PUT entity/organization/{id}` |

Методы возвращают `(*Organization, error)`. Успех требует HTTP 200 и непустого `id`.
Возвращается фактический ответ сервера, включая изменённые или заполненные сервером
значения. Совпадение ID ответа с ID маршрута дополнительно не проверяется.
При ошибке модель равна `nil`; `errors.As` распознаёт `*errors.APIError`,
`errors.Is` сохраняет ошибки авторизации, транспорта и отмены `context`.
Автоматических повторов, предварительного `Get` и переходов по ссылкам нет.

## Поддержанные поля

| Поля JSON | Чтение | Запись |
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
| `certificateDate` | `model.DateTime` | Create, Update |
| `actualAddress`, `legalAddress` | `string` | Create, Update |
| `actualAddressFull`, `legalAddressFull` | `*model.Address` | Create, Update через `model.AddressRequest` |
| `attributes` | `[]model.Attribute` | Create, Update |
| `owner`, `group`, `bonusProgram` | `*model.Reference` | Create, Update |
| `syncId` | `string` | только Create |
| `bonusPoints` | `int` | нет |
| `payerVat` | `bool` (`PayerVAT`) | Create, Update |
| `director`, `directorPosition`, `chiefAccountant`, `trackingContractNumber` | `string` | Create, Update |
| `trackingContractDate` | `model.DateTime` | Create, Update |
| `isEgaisEnable` | `bool` (`IsEGAISEnable`) | нет |

Отсутствующие поля читаются как нулевые значения. Отсутствующие и `null`-адреса
или ссылки дают `nil`, даты — нулевой `DateTime`. Неизвестные поля игнорируются.
Строковые реквизиты и коды сохраняют ведущие нули. `CompanyType` имеет открытый
набор значений: неизвестные строки сохраняются, `legal` не подставляется.
Реквизиты могут отсутствовать в зависимости от `companyType`; SDK не выводит их
из ФИО и не нормализует сочетания полей.

В официальной документации есть противоречие: таблица называет `isEgaisEnable`
всегда `false`, но примеры содержат `true`. SDK читает фактическое значение
без принудительной замены; отсутствие даёт `false`. Запись поля не поддержана.
`payerVat` поддержан отдельно и не означает поддержку всех настроек НДС.

## Запросы

В `CreateRequest` обязателен `Name string`; SDK отклоняет только `Name == ""`,
без trim. Все остальные поля — `model.Field[T]` с `omitzero`.
В `UpdateRequest` все изменяемые поля используют `Field`; полностью незаданный
request отклоняется. Get и Update проверяют форму UUID (Universally Unique
Identifier), без ограничения версии. Допустимость реквизитов, ссылок, `syncId`
и сочетаний полей проверяет сервер.

```go
created, err := service.Create(ctx, organization.CreateRequest{
    Name: "ИП Петров",
    CompanyType: model.Value(model.CompanyTypeEntrepreneur),
    INN: model.Value("001234567890"),
    PayerVAT: model.Value(false),
})

updated, err := service.Update(ctx, id, organization.UpdateRequest{
    ActualAddress: model.Value(""),
    DirectorPosition: model.Value("Директор"),
})
```

Незаданные поля отсутствуют в JSON. `Value(false)`, `Value("")`, `Value([]T{})`
и `Null[T]()` сохраняются; заданный `nil`-слайс кодируется как `[]`.
`syncId` доступен только при создании, поскольку после заполнения его нельзя
менять. Readonly-поля и отложенные поля отсутствуют в обоих request-типах.

Структурированные адреса используют `Field[model.AddressRequest]`.
Их части также используют `Field`, включая ссылки `country` и `region`.
Незаданная часть пропускается; пустая строка, `null` и пустой объект `{}` передаются
точно как заданы. Обе формы адреса могут уйти одновременно. По документации
сервер отдаёт приоритет структурированной форме, а при передаче только строки
отражает её также в `addInfo`. `null` не очищает адрес: для очистки передай `""`
в соответствующее строковое поле адреса. SDK не выполняет конкатенацию или merge
и не обещает серверного слияния отдельных частей. Общие типы описаны в
[requisites.md](requisites.md).

`attributes` использует существующие `model.Attribute` и helpers. Сервер обновляет
значения по метаданным: `[]` не очищает всю коллекцию. Для сброса отдельного значения
используй `model.AttrNull`, файлового значения — `model.AttrFileNull`. SDK не сливает
массив локально. Передача `null` сама по себе не обещает принятия или очистки поля
сервером.

## Границы

Не поддержаны банковские `accounts` целиком, региональные объекты
`mod__requisites__uz`/`mod__requisites__kz`, `companyVat__ru`, `advancePaymentVat`,
`chiefAccountSign`, `directorSign`, `stamp`, `fsrarId`, `utmUrl`.
Последние два поля встречаются в официальных примерах, но отсутствуют в таблицах
атрибутов и намеренно игнорируются. Также не поддержаны List, Delete, expand,
query-параметры, массовые операции и отдельные подресурсы.

Локальные тесты проверяют JSON обоими кодеками и HTTP с настоящим transport.
Upstream Organization используется без изменений; дополнительные случаи
и request golden синтетические. Живой API не проверялся, серверные правила
адресов и реквизитов локальным сервером не имитируются.

Контракт сверён с официальной документацией на revision
`9a811cef3dc4a3f3a46f070e60452361a06bf7cd`:
[поля и адреса Organization](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_organization.md#L14-L104),
[противоречащий таблице пример ЕГАИС](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_organization.md#L1648-L1653),
[attributes](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L314-L367),
[ограничение null для адресов](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187).
