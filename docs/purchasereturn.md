# PurchaseReturn

`entity/purchasereturn`: модель `PurchaseReturn`, `CreateRequest`, `UpdateRequest`,
`NewService(client)` и `Get/Create/Update`. Дополнительно `GetPositions` и aliases
`Position`/`PositionRequest` общих форм.
[Общие правила методов, JSON, ошибок и позиций](documents.md).

Минимальный Create: `organization`, `agent`, `store`.

## Поддержанные поля

| JSON | Read Go type | Create | Update |
| --- | --- | --- | --- |
| `id` | `string` | — | — |
| `accountId` | `string` | — | — |
| `name` | `string` | Field | Field |
| `description` | `string` | Field | Field |
| `code` | `string` | Field | Field |
| `externalCode` | `string` | Field | Field |
| `syncId` | `string` | Field | Field |
| `created` | `model.DateTime` | — | — |
| `updated` | `model.DateTime` | — | — |
| `deleted` | `model.DateTime` | — | — |
| `moment` | `model.DateTime` | Field | Field |
| `applicable` | `bool` | Field | Field |
| `shared` | `bool` | Field | Field |
| `printed` | `bool` | — | — |
| `published` | `bool` | — | — |
| `organization` | `*model.Reference` | required | Field |
| `owner` | `*model.Reference` | Field | Field |
| `group` | `*model.Reference` | Field | Field |
| `project` | `*model.Reference` | Field | Field |
| `state` | `*model.Reference` | Field | Field |
| `meta` | `model.Meta` | — | — |
| `attributes` | `[]model.Attribute` | Field | Field |
| `rate` | `*model.DocumentRate` | Field | Field |
| `sum` | `float64` | — | — |
| `agent` | `*model.Reference` | required | — |
| `organizationAccount` | `*model.Reference` | Field | Field |
| `agentAccount` | `*model.Reference` | Field | — |
| `contract` | `*model.Reference` | Field | Field |
| `vatSum` | `float64` | — | — |
| `positions` | `*model.Collection[Position]` | Field | Field |
| `vatEnabled` | `bool` | Field | Field |
| `vatIncluded` | `bool` | Field | Field |
| `payedSum` | `float64` | — | — |
| `store` | `*model.Reference` | required | Field |
| `payments` | `[]model.Reference` | Field | Field |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `supply` | `*model.Reference` | Field | — |
| `factureIn` | `*model.Reference` | Field | Field |
| `factureOut` | `*model.Reference` | Field | Field |

Required Reference передаются обычными значениями; остальные поля записи —
`model.Field[T]`. Вложенные write-формы: `DocumentRateRequest`, `[]PositionRequest`.
Readonly поля в requests отсутствуют.

## Ограничения и TODO

Files представлены только metadata коллекции, без обещания rows или файловых операций.
Create без основания требует organization, agent, store; для возврата по основанию
дополнительно передай Supply. UpdateRequest исключает Agent, AgentAccount и Supply.
API требует совпадения сторон/валюты с Supply, соответствия позиций основанию и
quantity не больше исходного; запрещает смену валюты и новые позиции Update. Применимость
этих формулировок к возврату без основания неоднозначна, SDK не вводит бизнес-валидатор.
Повторённая фраза Update о совместной смене agent/agentAccount противоречит прямому
запрету этих полей: они отсутствуют в UpdateRequest. OrganizationAccount доступен.
Discount помечен readonly в таблице, но встречается в write-примерах. Общий
PositionRequest содержит Discount, это не обещание его изменения у возврата.
Price/VAT доступны в embedded Update; правило «только quantity» источник относит к
отдельному PUT позиции (`UpdatePosition`); SDK его не проверяет, решает сервер. FactureIn/FactureOut/Payments пишутся как связи документов.
Позиции — все поля таблицы позиций источника: pack, slot, things (общая форма).

Источник полей: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/purchase-return),
[зафиксированный revision 9a811cef](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md).
