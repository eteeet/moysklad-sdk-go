Status: needs-triage
Stage: DRAFT
Type: research

# Проверка спорных контрактов документов

Проверено 2026-11-01. Источник — официальный `moysklad/api-remap-1.2-doc`, revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`. Demand, PurchaseOrder, PurchaseReturn, InvoiceIn/Out и general сверены с предоставленными полными локальными копиями `/private/tmp/moysklad-document-source.bEEQvf/`; CashIn/Out получены непосредственно с official raw GitHub на этом revision, без сохранения дополнительных файлов. Живой API (Application Programming Interface) аккаунта и `.env` не использовались. Примеры ниже — опубликованные примеры документации, а не результаты нового эксперимента. Выводы и предложения SDK (Software Development Kit) не утверждают публичный контракт автоматически.

## 1. Demand: store

- Факт: раздел «Отгрузки → Атрибуты сущности» требует `store`, если `applicable=true`; у поля нет метки «Необходимо при создании». [Таблица, L43][demand-store].
- Факт: раздел «Создать Отгрузку» без оговорки перечисляет `organization`, `agent`, `store` как обязательные поля. [Create, L557–565][demand-create].
- Вывод: документация содержит расхождение по условию обязательности. Безусловный локальный required-validator `store` отвергал бы сценарий `applicable=false`, допускаемый формулировкой таблицы; успешный новый вызов без склада этой проверкой не установлен. Основание: [таблица][demand-store], [Create][demand-create].
- Предложение SDK: `Store Field[Reference]` в Create/Update; Create локально требует только `organization`/`agent`, не требует `store`; сочетание `store`/`applicable` проверяет сервер. В безопасном минимальном примере передавать склад. Для Update нельзя требовать склад в теле при `applicable=true`: ссылка уже может быть у документа. Это граница локальной validation, а не обещание принять любое сочетание сервером.

## 2. PurchaseOrder: shipped, inTransit, waitSum

- Факт: в таблице позиций `shipped` и `inTransit` имеют «Обязательное при ответе» и `Change-handler`, но не «Только для чтения»; `Update-provider` у обоих также отсутствует. [Позиции, L64–76][po-positions].
- Факт: в таблице документа `waitSum` описан как сумма товаров в пути без read-only/write-маркера; соседний `shippedSum` явно read-only. [Атрибуты, L38–47][po-sums].
- Факт: embedded Create передаёт `inTransit=11/1/216`; документированный ответ возвращает `positions.meta` и `waitSum`, без значений строк. [Запрос, L712–758][po-embedded-create], [ответ, L905–920][po-embedded-response].
- Факт: отдельный Create строки передаёт `shipped=20`, `inTransit=19`; ответ содержит `shipped=0`, `inTransit=19`. [Create строки, L2648–2704][po-row-create].
- Факт: отдельный Update строки передаёт `shipped=10`, `inTransit=200`; ответ содержит `shipped=0`, `inTransit=200`. [Update строки, L2857–2911][po-row-update].
- Факт: документация разрешает embedded создание/обновление коллекции строк с полной заменой. [Позиции, L78–86][po-embedded-rules].
- Вывод: у `inTransit` есть существенная положительная опора для записи; объявлять его неподтверждённым наравне с `shipped` неверно. Наличие `shipped` в запросе само по себе не доказывает изменение: оба опубликованных ответа сохраняют 0. Для writable `waitSum` в проверенном документе нет примера запроса. Основание: [embedded Create][po-embedded-create], [Create строки][po-row-create], [Update строки][po-row-update], [таблица сумм][po-sums].
- Предложение SDK: включить `InTransit Field[float64]` в принадлежащий PurchaseOrder `PositionRequest`; оставить `shipped` и `waitSum` только в response первого среза. Последние два ограничения обозначать как scope SDK, а не как существующую read-only метку API. Совместное использование request для embedded Create/Update — проектное решение по [общему правилу embedded][po-embedded-rules]; отдельный успешный PUT строки не является новым live-тестом embedded Update.

## 3. PurchaseReturn: embedded positions и основание

- Факт: в таблице позиций `discount` явно read-only; `price`, `vat`, `vatEnabled` этой метки не имеют. Общее правило API говорит, что read-only поля игнорируются при Create/Update. [Таблица, L56–68][return-positions], [general, L1152–1155][general-write].
- Факт: embedded Create с `supply`, Create «без основания» и embedded Update передают `price`, `vat`, `discount`; `vatEnabled` в этих request-примерах строк отсутствует. [Create с основанием, L659–725][return-create-based], [Create без основания, L865–998][return-create-unbased], [Update, L2249–2321][return-embedded-update].
- Факт: ответы обоих document Create и embedded Update показывают только `positions.meta`. [Based Create, L843–861][return-based-response], [unbased Create, L1115–1124][return-unbased-response], [Update, L2446–2455][return-update-response].
- Вывод: request-пример с `discount` не отменяет его явную read-only метку; конкретная запись `vatEnabled` примером не проверена. Document-ответы не доказывают фактическое изменение цены/налога строк относительно исходной приёмки. Основание: [таблица][return-positions], [embedded request][return-embedded-update], [document response][return-update-response].
- Факт: только раздел «Изменить Позицию» отдельного ресурса явно допускает изменение лишь `quantity`; предложение о диапазоне содержит пропуск перед «где n». [Отдельный PUT, L2852–2855][return-row-update].
- Факт: embedded Create/Update строк разрешены общим разделом; при Update коллекция полностью заменяется. [Позиции, L73–81][return-embedded-rules].
- Вывод: перенос quantity-only с отдельного PUT на embedded Update не обоснован. Для `price`/`vat`/`vatEnabled` есть таблица без read-only и общее правило обновления, но опубликованные document-ответы не устанавливают все условия их изменения. Основание: [правило embedded][return-embedded-rules], [document Update, L2180–2192][return-update-rules], [отдельный PUT][return-row-update], [таблица][return-positions].
- Факт: `supply` обязателен только для возврата с основанием; без основания его указывать не нужно. Для возврата по документу предписаны совпадение контрагента, валюты и юрлица; embedded Create допускает отличие строк лишь по количеству, не выше исходного. [Создать Возврат поставщику, L569–583][return-create-rules].
- Вывод: ограничения, ссылающиеся на позиции/валюту исходного документа, нельзя превращать в локальный запрет для режима без `supply`. Отсутствие основания документировано отдельно; дополнительных правил unbased Update проверенный текст не устанавливает. Основание: [Create][return-create-rules], [unbased пример][return-create-unbased], [Update][return-update-rules].
- Факт: document Update сначала содержит общую фразу о парных `organizationAccount`/`agentAccount`, но затем прямо запрещает изменение `agentAccount`, `agent`, `supply`. [Изменить Возврат поставщику, L2180–2192][return-update-rules].
- Предложение SDK: отдельные Create/Update requests; `agent`, `agentAccount`, `supply` отсутствуют в Update. В embedded request разрешить передавать `price`, `vat`, `vatEnabled`, исключить `discount`; не обещать принятие изменения цены/налога при основании. Сопоставление с исходной приёмкой и её количествами оставить серверу без предварительного GET. Quantity-only не использовать как общий embedded-validator.

## 4. CashIn/CashOut: account-поля

- Факт: полные таблицы атрибутов CashIn и CashOut содержат `organization`/`agent`, но не `organizationAccount`/`agentAccount`. [CashIn, L5–38][cashin-attributes], [CashOut, L5–41][cashout-attributes].
- Факт: оба раздела Update требуют обновлять эти account-поля при смене соответствующей стороны, иначе обещают ошибку. [CashIn Update, L1523–1531][cashin-update], [CashOut Update, L1694–1702][cashout-update].
- Вывод: конфликт подтверждён первичными исходниками; это не доказательство ни полного отсутствия account-полей в API, ни гарантированной поддержки их записи. Основание: [CashIn таблица][cashin-attributes], [CashIn Update][cashin-update], [CashOut таблица][cashout-attributes], [CashOut Update][cashout-update].
- Предложение SDK: первый срез по таблицам, без account-полей и без локальной проверки парности. Ошибку сервера при смене стороны возвращать как есть; не гарантировать такую смену без дополнительных условий. Поддержку account-полей и успешную смену сторон считать отдельным вопросом, требующим уточнения поставщика API либо разрешённого live-эксперимента.

## 5. InvoiceIn/InvoiceOut: name в Create

- Факт: `name` в «Атрибуты сущности» обоих счетов не имеет метки «Необходимо при создании». [InvoiceIn, L27][invoicein-name], [InvoiceOut, L25][invoiceout-name].
- Факт: раздел «Создать Счет поставщика» прямо называет `name`, `organization`, `agent` обязательными; минимальный пример передаёт `name`. [InvoiceIn Create, L348–367][invoicein-create].
- Факт: раздел «Создать Счет покупателю» содержит то же прямое требование и `name` в минимальном примере. [InvoiceOut Create, L476–495][invoiceout-create].
- Вывод: отсутствие метки в общей таблице не отменяет явного требования конкретной операции. Документированный узкий Create scope требует `name`; поведение сервера при omission/пустой строке здесь не проверено. Основание: [InvoiceIn Create][invoicein-create], [InvoiceOut Create][invoiceout-create].
- Предложение SDK: требовать заданное непустое `Name` в Create, не добавлять требование в Update. Не выводить автонумерацию и конкретный default из отсутствующей метки.

## 6. Дополнение: syncId не обязательно Create-only

- Факт: general описывает `syncId` для предотвращения дублирования при повторном Create; отдельного запрета первого заполнения в Update там нет. Таблицы Demand/PurchaseOrder/PurchaseReturn запрещают изменение **после заполнения**, а не все Update. [Назначение, L1144–1149][syncid-purpose], [Demand, L45][demand-syncid], [PurchaseOrder, L42][po-syncid], [PurchaseReturn, L36][return-syncid].
- Вывод: Create-only — дополнительное сужение SDK, не доказанное требование API; состояние ранее заполненного поля нельзя узнать из одного Update request. Основание: [описание назначения][syncid-purpose], [табличное ограничение][po-syncid].
- Предложение SDK: `SyncID Field[string]` и в Update; сервер проверяет неизменяемость после заполнения. Не выполнять скрытый GET и не обещать очистку/перезапись заполненного значения.

## Что можно закрыть без живого API

- Достаточно documented scope и решения о публичном SDK: optional Demand Store без безусловного validator; запись PurchaseOrder `inTransit`; response-only SDK `shipped`/`waitSum`; исключение PurchaseReturn `discount` и запрещённых Update refs; отдельная semantics embedded/row; обязательный Invoice Name; отсутствие произвольного Create-only у `syncId`. Все решения остаются предложениями до triage.
- Человеку требуется выбрать scope SDK и утвердить публичные типы. Проверка через живой API не является обязательным условием для перечисленного узкого scope.
- Уточнение поставщика API либо отдельно разрешённый эксперимент нужен, только если требуется гарантия спорного поведения: Cash account-поля/смена сторон; запись `shipped`/`waitSum`; изменение price/VAT у based PurchaseReturn; точные ограничения отдельного PUT без основания; реальные defaults/ошибки при Demand без store или Invoice без name. Эти гарантии не заявлены и не должны блокировать независимо определённые общие типы.

[demand-store]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L43
[demand-create]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L557-L565
[po-positions]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L64-L76
[po-sums]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L38-L47
[po-embedded-create]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L712-L758
[po-embedded-response]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L905-L920
[po-row-create]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L2648-L2704
[po-row-update]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L2857-L2911
[po-embedded-rules]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L78-L86
[return-positions]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L56-L68
[general-write]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1152-L1155
[return-create-based]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L659-L725
[return-create-unbased]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L865-L998
[return-embedded-update]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L2249-L2321
[return-based-response]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L843-L861
[return-unbased-response]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L1115-L1124
[return-update-response]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L2446-L2455
[return-row-update]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L2852-L2855
[return-embedded-rules]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L73-L81
[return-create-rules]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L569-L583
[return-update-rules]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L2180-L2192
[cashin-attributes]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashin.md#L5-L38
[cashout-attributes]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashout.md#L5-L41
[cashin-update]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashin.md#L1523-L1531
[cashout-update]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_cashout.md#L1694-L1702
[invoicein-name]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_invoice_in.md#L27
[invoiceout-name]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_invoice_out.md#L25
[invoicein-create]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_invoice_in.md#L348-L367
[invoiceout-create]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_invoice_out.md#L476-L495
[syncid-purpose]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1144-L1149
[demand-syncid]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L45
[po-syncid]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L42
[return-syncid]: https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L36
