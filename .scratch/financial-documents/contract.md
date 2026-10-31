Status: needs-triage
Stage: DRAFT

# DRAFT — поимённый контракт финансовых документов

Нормативное продолжение [spec](spec.md). Все решения write scope и новые Go-типы ниже —
proposals; таблицы источника и известные противоречия собраны в [research](research.md).
Для любой из четырёх сущностей применять всю её таблицу и общие правила этого файла.

## Обозначения и общие правила

- `Ref` = model.Reference, при чтении ссылки `*model.Reference` с omitzero.
  `DateTime` = существующий model.DateTime с omitzero; absent/null даёт нулевую дату.
  `Attrs` = []model.Attribute с существующими helpers и семантикой дополнительных полей.
- `DocumentRate`/`DocumentRateRequest` — общий контракт document-common, не повторная финансовая реализация.
  `Op`/`OpRequest` — предложенные FinanceOperation/FinanceOperationRequest из ticket 01.
- `C/U` = CreateRequest/UpdateRequest. `C!` = поле обязательно Create согласно документации.
  `—` = request-поле отсутствует. SyncID доступен C/U, запрет изменения после заполнения проверяет сервер.
  Поля `C!` в Create — обычный Ref; прочие C/U — model.Field[T] с omitzero;
  Update всех изменяемых полей — Field[T]. Имена Go следуют принятому ID/VAT/acronym style.
- `R!` = документация требует поле в ответе. Это не новый контроль полноты декодера:
  отсутствующие поля чтения допустимы, как в Product/UOM; Service требует только ID.
- Столбец null описывает серверный контракт, а не способность Field закодировать null:
  `О` — общий условный сброс необязательного object, если UI допускает пустое значение;
  `Да` — явно документированное rate:null; `?` — конкретный сброс не установлен;
  `нет` — required-ссылка не объявлена очищаемой; `—` — записи нет.
  Field.Null технически передаётся и для `?`; SDK не обещает принятия сервером.
- Read readonly: accountId, id, created, updated, deleted, printed, published;
  собственный meta также только read в предложенном SDK. syncId читается и пишется C/U; SDK не загружает прежнее значение для проверки write-once.
- Files присутствует в атрибутных таблицах как MetaArray/R!, но намеренно исключён
  целиком и из модели, и из requests. Это явный TODO, не пропущенный контракт.

## PaymentIn

Источник: [атрибуты и методы payment-in](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_in.md).

| JSON-поле | Read тип | Ответ | Write | null | Required/write-once/особенность |
| --- | --- | --- | --- | --- | --- |
| accountId | string | R! | — | — | readonly |
| agent | *Ref | R! | C!/U | нет | Counterparty |
| agentAccount | *Ref | optional | C/U | О | счет Counterparty, только ссылка |
| applicable | bool | R! | C/U | ? | false сохраняется |
| attributes | Attrs | optional | C/U | ? | reset отдельного value:null по Attribute |
| code | string | optional | C/U | ? | пустая строка сохраняется |
| contract | *Ref | optional | C/U | О | contract:null явно приведён общим разделом |
| created | DateTime | R! | — | — | readonly |
| deleted | DateTime | optional | — | — | readonly |
| description | string | optional | C/U | ? | без trim |
| externalCode | string | R! | C/U | ? | server default при отсутствии |
| files | исключён | R! API | — | — | TODO целиком |
| group | *Ref | R! | C/U | ? | не обещать очистку required response-поля |
| id | string | R! | — | — | readonly |
| incomingDate | DateTime | optional | C/U | ? | без новой нормализации DateTime |
| incomingNumber | string | optional | C/U | ? | строка, ведущие нули сохраняются |
| meta | model.Meta | R! | — | — | readonly SDK |
| moment | DateTime | R! | C/U | ? | server default при отсутствии |
| name | string | R! | C/U | ? | необязателен Create, автоимя сервера |
| organization | *Ref | R! | C!/U | нет | Organization |
| organizationAccount | *Ref | optional | C/U | О | счет Organization, только ссылка |
| owner | *Ref | optional | C/U | О | Employee |
| paymentPurpose | string | R! | C/U | ? | назначение платежа |
| printed | bool | R! | — | — | readonly |
| project | *Ref | optional | C/U | О | только ссылка |
| published | bool | R! | — | — | readonly |
| rate | *DocumentRate | R! | C/U DocumentRateRequest | Да | null: валюта по умолчанию, пересчёт сервером |
| salesChannel | *Ref | optional | C/U | О | только ссылка |
| shared | bool | R! | C/U | ? | false сохраняется |
| state | *Ref | optional | C/U | О | только ссылка |
| sum | float64 | R! | C/U | ? | Float, без масштабирования |
| syncId | string | optional | C/U | ? | после заполнения нельзя изменять; проверяет сервер |
| updated | DateTime | R! | — | — | readonly |
| vatSum | float64 | R! | C/U | ? | Float, включая 0 |
| factureOut | *Ref | optional | C/U | О | связь со счетом-фактурой выданным |
| operations | []Op | optional | C/U []OpRequest | ? | пропуск, [] и массив — разные JSON |

Итог proposed write scope: agent, agentAccount, applicable, attributes, code, contract,
description, externalCode, group, incomingDate, incomingNumber, moment, name, organization,
organizationAccount, owner, paymentPurpose, project, rate, salesChannel, shared, state,
sum, syncId, vatSum, factureOut, operations. Остальных request-полей нет.

## PaymentOut

Источник: [атрибуты и методы payment-out](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_out.md).

| JSON-поле | Read тип | Ответ | Write | null | Required/write-once/особенность |
| --- | --- | --- | --- | --- | --- |
| accountId | string | R! | — | — | readonly |
| accrualDate | DateTime | optional | C/U | ? | дата начисления |
| agent | *Ref | R! | C!/U | нет | Counterparty, Employee или Organization |
| agentAccount | *Ref | optional | C/U | О | счет Counterparty/Organization, только ссылка |
| applicable | bool | R! | C/U | ? | false сохраняется |
| attributes | Attrs | optional | C/U | ? | reset отдельного value:null по Attribute |
| code | string | optional | C/U | ? | строка |
| contract | *Ref | optional | C/U | О | contract:null |
| created | DateTime | R! | — | — | readonly |
| deleted | DateTime | optional | — | — | readonly |
| description | string | optional | C/U | ? | без trim |
| expenseItem | *Ref | R! | C!/U | нет | обязательно Create, default не описан |
| externalCode | string | R! | C/U | ? | server default при отсутствии |
| files | исключён | R! API | — | — | TODO целиком |
| group | *Ref | R! | C/U | ? | только ссылка |
| noClosingDocs | bool | R! | C/U | ? | true несовместим с непустым operations |
| id | string | R! | — | — | readonly |
| meta | model.Meta | R! | — | — | readonly SDK |
| moment | DateTime | R! | C/U | ? | server default при отсутствии |
| name | string | R! | C/U | ? | необязателен Create |
| organization | *Ref | R! | C!/U | нет | Organization |
| organizationAccount | *Ref | optional | C/U | О | счет Organization, только ссылка |
| owner | *Ref | optional | C/U | О | Employee |
| paymentPurpose | string | R! | C/U | ? | назначение платежа |
| printed | bool | R! | — | — | readonly |
| project | *Ref | optional | C/U | О | только ссылка |
| published | bool | R! | — | — | readonly |
| rate | *DocumentRate | R! | C/U DocumentRateRequest | Да | default currency и пересчёт сервером |
| salesChannel | *Ref | optional | C/U | О | только ссылка |
| shared | bool | R! | C/U | ? | false сохраняется |
| state | *Ref | optional | C/U | О | только ссылка |
| sum | float64 | R! | C/U | ? | Float, без арифметики SDK |
| syncId | string | optional | C/U | ? | после заполнения нельзя изменять; проверяет сервер |
| updated | DateTime | R! | — | — | readonly |
| vatSum | float64 | R! | C/U | ? | Float |
| factureIn | *Ref | optional | C/U | О | связь со счетом-фактурой полученным |
| operations | []Op | optional | C/U []OpRequest | ? | типы исходящих операций |

Итог proposed write scope: accrualDate, agent, agentAccount, applicable, attributes, code,
contract, description, expenseItem, externalCode, group, noClosingDocs, moment, name,
organization, organizationAccount, owner, paymentPurpose, project, rate, salesChannel,
shared, state, sum, syncId, vatSum, factureIn, operations. IncomingDate и
incomingNumber отсутствуют, несмотря на похожесть имени документа на PaymentIn.

## CashIn

Источник: [атрибуты и методы cashin](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashin.md).

| JSON-поле | Read тип | Ответ | Write | null | Required/write-once/особенность |
| --- | --- | --- | --- | --- | --- |
| accountId | string | R! | — | — | readonly |
| agent | *Ref | R! | C!/U | нет | Counterparty по документации |
| applicable | bool | R! | C/U | ? | false сохраняется |
| attributes | Attrs | optional | C/U | ? | существующий Attribute |
| code | string | optional | C/U | ? | строка |
| contract | *Ref | optional | C/U | О | contract:null |
| created | DateTime | R! | — | — | readonly |
| deleted | DateTime | optional | — | — | readonly |
| description | string | optional | C/U | ? | без trim |
| externalCode | string | R! | C/U | ? | server default при отсутствии |
| files | исключён | R! API | — | — | TODO целиком |
| group | *Ref | R! | C/U | ? | только ссылка |
| id | string | R! | — | — | readonly |
| meta | model.Meta | R! | — | — | readonly SDK |
| moment | DateTime | R! | C/U | ? | server default при отсутствии |
| name | string | R! | C/U | ? | необязателен Create |
| organization | *Ref | R! | C!/U | нет | Organization |
| owner | *Ref | optional | C/U | О | Employee |
| paymentPurpose | string | R! | C/U | ? | основание ордера |
| printed | bool | R! | — | — | readonly |
| project | *Ref | optional | C/U | О | только ссылка |
| published | bool | R! | — | — | readonly |
| rate | *DocumentRate | R! | C/U DocumentRateRequest | Да | default currency и пересчёт сервером |
| salesChannel | *Ref | optional | C/U | О | только ссылка |
| shared | bool | R! | C/U | ? | false сохраняется |
| state | *Ref | optional | C/U | О | только ссылка |
| sum | float64 | R! | C/U | ? | Float |
| syncId | string | optional | C/U | ? | после заполнения нельзя изменять; проверяет сервер |
| updated | DateTime | R! | — | — | readonly |
| vatSum | float64 | R! | C/U | ? | Float |
| factureOut | *Ref | optional | C/U | О | счет-фактура выданный |
| operations | []Op | optional | C/U []OpRequest | ? | типы входящих операций |

Итог proposed write scope: agent, applicable, attributes, code, contract, description,
externalCode, group, moment, name, organization, owner, paymentPurpose, project, rate,
salesChannel, shared, state, sum, syncId, vatSum, factureOut, operations.
AgentAccount/organizationAccount, incomingDate/incomingNumber, expenseItem, noClosingDocs
и accrualDate в CashIn отсутствуют.

## CashOut

Источник: [атрибуты и методы cashout](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashout.md).

| JSON-поле | Read тип | Ответ | Write | null | Required/write-once/особенность |
| --- | --- | --- | --- | --- | --- |
| accountId | string | R! | — | — | readonly |
| accrualDate | DateTime | optional | C/U | ? | дата начисления |
| agent | *Ref | R! | C!/U | нет | Counterparty по документации; multi-type не подтверждён |
| applicable | bool | R! | C/U | ? | false сохраняется |
| attributes | Attrs | optional | C/U | ? | существующий Attribute |
| code | string | optional | C/U | ? | строка |
| contract | *Ref | optional | C/U | О | contract:null |
| created | DateTime | R! | — | — | readonly |
| deleted | DateTime | optional | — | — | readonly |
| description | string | optional | C/U | ? | без trim |
| expenseItem | *Ref | R! | C!/U | нет | обязательно Create, default не описан |
| externalCode | string | R! | C/U | ? | server default при отсутствии |
| files | исключён | R! API | — | — | TODO целиком |
| group | *Ref | R! | C/U | ? | только ссылка |
| noClosingDocs | bool | R! | C/U | ? | true несовместим с непустым operations |
| id | string | R! | — | — | readonly |
| meta | model.Meta | R! | — | — | readonly SDK |
| moment | DateTime | R! | C/U | ? | server default при отсутствии |
| name | string | R! | C/U | ? | необязателен Create |
| organization | *Ref | R! | C!/U | нет | Organization |
| owner | *Ref | optional | C/U | О | Employee |
| paymentPurpose | string | R! | C/U | ? | основание ордера |
| printed | bool | R! | — | — | readonly |
| project | *Ref | optional | C/U | О | только ссылка |
| published | bool | R! | — | — | readonly |
| rate | *DocumentRate | R! | C/U DocumentRateRequest | Да | default currency и пересчёт сервером |
| salesChannel | *Ref | optional | C/U | О | только ссылка |
| shared | bool | R! | C/U | ? | false сохраняется |
| state | *Ref | optional | C/U | О | только ссылка |
| sum | float64 | R! | C/U | ? | Float |
| syncId | string | optional | C/U | ? | после заполнения нельзя изменять; проверяет сервер |
| updated | DateTime | R! | — | — | readonly |
| vatSum | float64 | R! | C/U | ? | Float |
| factureIn | *Ref | optional | C/U | О | счет-фактура полученный |
| operations | []Op | optional | C/U []OpRequest | ? | типы исходящих операций |

Итог proposed write scope: accrualDate, agent, applicable, attributes, code, contract,
description, expenseItem, externalCode, group, noClosingDocs, moment, name, organization,
owner, paymentPurpose, project, rate, salesChannel, shared, state, sum, syncId,
vatSum, factureIn, operations. Account-ссылок, incomingDate и incomingNumber нет.

## Общая валюта и allocation операций

| Контракт | JSON-поле | Read | Write proposal | Семантика |
| --- | --- | --- | --- | --- |
| DocumentRate | currency | *Reference | Field[Reference] внутри DocumentRateRequest | ссылка; обязательна в ответе |
| DocumentRate | value | *float64 | Field[float64] внутри DocumentRateRequest | в ответе отсутствует при курсе 1; отсутствие не подменяется числом |
| FinanceOperation | meta | model.Meta | model.Meta в FinanceOperationRequest | meta лежит непосредственно на элементе operations |
| FinanceOperation | linkedSum | float64 | Field[float64] | необязателен при записи, 0 сохраняется |

DocumentRate требует [внешний ticket](../document-common/issues/01-rate.md); итоговая форма и имена
общего типа утверждаются в нём. Внешнее поле rate — Field[DocumentRateRequest] в requests:
незадано, null, {} или выбранные поля не схлопываются. Пустой объект не объявляется
валидным для API только потому, что SDK может его сериализовать.

FinanceOperation не содержит ID/name/positions/payments и других расширенных данных.
Unknown Meta.Type читается и передаётся без локальной проверки, как в Reference.
Список серверных типов ниже — документация ограничений, а не новый закрытый enum SDK.

| Платежи | Допустимые operations по официальной документации |
| --- | --- |
| PaymentIn, CashIn | customerorder, purchasereturn, demand, invoiceout, commissionreportin, retailshift |
| PaymentOut, CashOut | salesreturn, supply, invoicein, purchaseorder, commissionreportout, payroll |

Описание перечисляет retailShift, но реальные meta.type примеров — retailshift:
proposal использовать фактическую строку meta без нормализации/констант.
LinkedSum — оплаченная данным платежом сумма конкретной операции; read 0 не означает,
что SDK вычислил распределение. Нет локальной проверки суммы allocations, типов операций,
соответствия сторон или пересчёта валют. [] отправляется точно, но гарантированная серверная
очистка/полное замещение operations не доказаны просмотренными разделами.

## Правила сервера, которые SDK не воспроизводит

- PaymentIn/PaymentOut: при изменении organization/agent документация требует одновременно
  обновлять organizationAccount/agentAccount соответственно. Передавать поля явно;
  SDK не делает GET счетов, не сбрасывает их автоматически и не вводит локальную проверку пары.
- Если заполненный organizationAccount банковского платежа имеет валюту, отличную
  от валюты учёта Account, rate.currency должен соответствовать валюте счёта; иначе API 22004.
  Для CashIn/CashOut проверка по счету Organization прямо объявлена неприменимой.
- В CashIn/CashOut текст Update повторяет требование account-полей, отсутствующих
  в атрибутных таблицах. Proposal — поля не добавлять по одной этой фразе;
  сохранить её как неразрешённое противоречие, не блокируя составление DRAFT.
- noClosingDocs=true и непустой operations у исходящих несовместимы.
  Proposal — отправлять запрос и передавать APIError; не добавлять локальный бизнес-валидатор.
- Attributes обновляются существующим механизмом по metadata; [] не является обещанием
  очистки всей коллекции. Поля значений/null следуют общему Attribute, не новой финансовой модели.
