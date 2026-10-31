# Общие JSON-типы документов

Типы пакета `model` не выполняют HTTP-запросов и не требуют реализации связанных
сущностей. Чтение и запись используют разные формы; JSON проверяется через
`encoding/json` и `encoding/json/v2`.

## Валюта

`DocumentRate` читает `Currency *Reference` и `Value *float64`.
В корректном присутствующем `rate` значение `Value == nil` означает курс 1:
сервер не передаёт `value`, равный 1. Указатель сохраняет компактный ответ и отличает
отсутствующее поле от явного нуля. SDK не подставляет 1 в структуру или JSON.
Отсутствие всего `rate` не означает известную валюту/курс; незаданный request также
не получает локального default. JSON `null` читается как nil, без обещания отдельного
серверного смысла для `value:null`.

`DocumentRateRequest` содержит `Currency Field[Reference]` и `Value Field[float64]`.
Для поля документа используй `Field[DocumentRateRequest]` с `omitzero`:

```go
type request struct {
	Rate model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
}

currencyOnly := request{Rate: model.Value(model.DocumentRateRequest{
	Currency: model.Value(currencyReference),
})}
reset := request{Rate: model.Null[model.DocumentRateRequest]()}
```

Unset не отправляется, Null отправляет `rate:null`, Value пустой request — `rate:{}`.
Вложенные Field сохраняют пропуски, null, 0 и выбранный курс. Сериализуемость `{}`
или null не означает допустимость для каждого API-поля. `rate:null` может сменить валюту
на валюту по умолчанию с серверным пересчётом цен, сумм и расходов; SDK этого не вычисляет.
Курс — не денежная сумма в копейках.

## Накладные расходы

`DocumentOverhead` содержит `Sum float64` и `Distribution OverheadDistribution`.
Для записи — `DocumentOverheadRequest` с `Sum Field[float64]` и
`Distribution Field[OverheadDistribution]`; внешнее поле документа также использует Field.

`OverheadDistribution` — открытый строковый тип с константами:

| Константа | JSON |
| --- | --- |
| `OverheadDistributionWeight` | `weight` |
| `OverheadDistributionVolume` | `volume` |
| `OverheadDistributionPrice` | `price` |

Неизвестная или пустая строка сохраняется, default и enum-валидация не добавляются.
Сумма передаётся в копейках без округления или распределения на клиенте. Unset, null,
пустой/частичный объект и явный `sum:0` различаются.
Объект расходов документа не является скалярным `overhead` отдельной позиции.

## Связь платежа с операцией

`FinanceOperation` содержит `Meta model.Meta` и `LinkedSum float64`.
`FinanceOperationRequest` — те же `Meta` и `LinkedSum Field[float64]` с `omitzero`.
JSON плоский, без вложенного `operation`:

```json
{"meta":{"href":"https://example.test/entity/invoiceout/id","type":"invoiceout"},"linkedSum":0}
```

LinkedSum — сумма, оплаченная этим платежом по конкретной операции. В request она
необязательна; meta-only, явный 0 и дробное значение различаются. При чтении отсутствующий
linkedSum даёт 0 без вычисления распределения. Неизвестный Meta.Type сохраняется;
расширенные поля связанного документа не становятся поддержанными полями модели.

Внешний `Field[[]FinanceOperationRequest]` различает unset, null, пустой/nil slice как []
и непустой массив. Это wire-контракт, не обещание, что API очищает или полностью заменяет
`operations` по []. Допустимые типы связей, суммы и стороны проверяет сервер.

## Границы

Все денежные значения SDK — float64 в копейках независимо от Int/Float в таблицах API;
точной десятичной арифметики и точности integer за пределами 2^53 нет. Существующие
сырые Attribute long сохраняют свой отдельный целочисленный контракт.

Эти типы не реализуют Get/Create/Update документов, Currency, банковские счета,
финансовые вычисления, retry или expand. Тесты проверяют JSON, не бизнес-поведение
живого Account. Поддержку HTTP-документов утверждают отдельными tickets.

Первичные контракты закреплены на revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`:
[валюта](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_common_info.md#L689-L699),
[расходы документа](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L47-L53),
[привязка платежей](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_common_info.md#L476-L507),
[null](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187).
