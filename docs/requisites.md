# Общие типы реквизитов

Пакет `model` предоставляет самостоятельные `Address`, `AddressRequest` и `CompanyType`.
Наличие этих типов не означает поддержку HTTP-операций Counterparty или Organization.

`CompanyType` — именованный тип поверх `string` с открытым набором значений:

| Константа | Значение JSON |
| --- | --- |
| `CompanyTypeLegal` | `legal` |
| `CompanyTypeEntrepreneur` | `entrepreneur` |
| `CompanyTypeIndividual` | `individual` |
| `CompanyTypeLegalUZ` | `legalUZ` |
| `CompanyTypeEntrepreneurUZ` | `entrepreneurUZ` |
| `CompanyTypeIndividualUZ` | `individualUZ` |
| `CompanyTypeLegalKZ` | `legalKZ` |
| `CompanyTypeEntrepreneurKZ` | `entrepreneurKZ` |
| `CompanyTypeIndividualKZ` | `individualKZ` |

Неизвестная строка сохраняется при чтении и записи; нулевое значение — `""`.
SDK не валидирует тип и не подставляет `legal`. В запросах используй
`Field[CompanyType]` с `omitzero`, чтобы различать отсутствие, пустую строку и `null`.

`Address` описывает чтение десяти адресных полей: `addInfo`, `apartment`, `city`,
`comment`, `fiasCode__ru`, `country`, `house`, `postalCode`, `region`, `street`.
`Country` и `Region` — `*Reference`; отсутствующая или `null`-ссылка даёт `nil`.
Остальные поля — строки, в том числе коды с ведущими нулями. Неизвестные JSON-поля
игнорируются. Отсутствующий или `null`-адрес в поле `*Address` даёт `nil`.

Для записи используй отдельный `AddressRequest`. Те же десять полей имеют типы
`Field[string]` и `Field[Reference]` с `omitzero`. Незаданные части не отправляются,
`Value("")` сохраняет пустую строку, `Null` кодирует явный `null`. Внешний
`Field[AddressRequest]` также различает отсутствие, `null` и значение, включая `{}`:

```go
type request struct {
	Address model.Field[model.AddressRequest] `json:"actualAddressFull,omitzero"`
}

partial := request{Address: model.Value(model.AddressRequest{
	City:       model.Value("Москва"),
	PostalCode: model.Value("001234"),
})}
empty := request{Address: model.Value(model.AddressRequest{})} // {"actualAddressFull":{}}
```

SDK передаёт обе заданные формы адреса без конкатенации, нормализации или слияния.
По документации API при одновременной передаче строковой и структурированной форм
сервер отдаёт приоритет структуре; при передаче только строки сервер отражает её
также в `addInfo`. `null` не очищает адрес: для очистки передай `""` в соответствующее
строковое поле адреса, например `actualAddress`.

Эти ограничения описывают сервер. Локальные тесты проверяют точное кодирование и
чтение через `encoding/json` и `encoding/json/v2`; они не проверяют принятие `{}`,
очистку или правила замены отдельных частей адреса живым API.

Контракт сверён с официальной документацией на revision
`9a811cef3dc4a3f3a46f070e60452361a06bf7cd`:
[адрес Counterparty](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_counterparty.md#L100-L120),
[CompanyType](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_counterparty.md#L192-L206),
[ограничение null](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187).
