# Move

`entity/move`: модель `Move`, `CreateRequest`, `UpdateRequest`,
`NewService(client)` и `Get/Create/Update`. Дополнительно `GetPositions` и aliases
`Position`/`PositionRequest` общих форм.
[Общие правила методов, JSON, ошибок и позиций](documents.md).

Минимальный Create: `organization`, `sourceStore`, `targetStore`.

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
| `positions` | `*model.Collection[Position]` | Field | Field |
| `files` | `*model.Collection[model.Reference]` | — | — |
| `sourceStore` | `*model.Reference` | required | Field |
| `targetStore` | `*model.Reference` | required | Field |
| `internalOrder` | `*model.Reference` | Field | Field |
| `customerOrder` | `*model.Reference` | Field | Field |
| `demand` | `*model.Reference` | — | — |
| `supply` | `*model.Reference` | — | — |
| `productionTasks` | `[]model.Reference` | Field | Field |
| `overhead` | `*model.DocumentOverhead` | Field | Field |

Required Reference передаются обычными значениями; остальные поля записи —
`model.Field[T]`. Вложенные write-формы: `DocumentRateRequest`, `DocumentOverheadRequest`,
`[]PositionRequest`.
Readonly поля в requests отсутствуют.

## Ограничения и TODO

Files представлены только metadata коллекции, без обещания rows или файловых операций.
Минимальный Create требует два склада и юрлицо. Agent, store, счета сторон,
contract и VAT-поля заголовка не добавлены по аналогии с коммерческими документами.
Нулевые VAT/VATEnabled общей позиции означают отсутствие полей у Move, а не «без НДС».
Неиспользуемые коммерческие Field строки обычно не задаются; допустимость заданных
проверяет сервер. SourceSlot/TargetSlot — ссылки. Overhead документа распределяет
сервер; наличие строк, равенство/различие складов и тип assortment локально не проверяются.
Demand и Supply только читаются (маркер таблицы); ProductionTasks пишется как связь документа.
Таблица допускает service в ассортименте, вводный текст его не перечисляет;
SDK не вводит локальное ограничение по этому расхождению.
Строки — все поля таблицы позиций источника: pack, things, sourceSlot, targetSlot (общая форма);
overhead строки только читается.

Источник полей: [официальная документация](https://dev.moysklad.ru/doc/api/remap/1.2/#/documents/move),
[зафиксированный revision 9a811cef](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_move.md).
