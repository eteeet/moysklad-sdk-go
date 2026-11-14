Status: resolved
Причина: работа закрыта: `document-positions/spec.md` resolved. Закрыто 2026-11-14 при сверке здоровья документов.
Stage: DRAFT

# Positions: первичные факты для DRAFT

Проверено по официальным Markdown API 1.2, commit
`9a811cef3dc4a3f3a46f070e60452361a06bf7cd`. Повторная проверка общего исследования
2026-11-01 подтвердила тот же master HEAD через
[GitHub API официального repository](https://api.github.com/repos/moysklad/api-remap-1.2-doc/commits/master).
В этом рабочем контексте git ls-remote не разрешил github.com через DNS; повторную
успешную проверку выполнил агент общей группы подготовки. Live API и .env не использовались.

Раннее исследование 2026-10-31 предлагало иной состав групп и исключение для expand.
Его рекомендации не действуют: состав определён текущим решением пользователя,
positions — общий модуль, весь текущий scope выполняется без expand. Ниже только
факты первичных источников и отдельно обозначенные проектные предложения.

## Envelope, страницы и запись

| Проверенный факт | Первичный источник |
| --- | --- |
| Общий раздел перечисляет многие документы с позициями, включая inventory и payroll; работа возможна внутри документа и через специальные ресурсы | [Общий раздел positions](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1022-L1029) |
| Embedded positions в Create/Update — массив; при Update весь набор заменяется, matching id обновляется, новые строки добавляются, отсутствующие старые удаляются | [Embedded write](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1031-L1037) |
| Обычный supply Get без query возвращает positions.meta без rows | [Get request](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L2079-L2094), [positions в ответе](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L2174-L2184) |
| Специальный GET /entity/supply/{id}/positions возвращает meta/context/rows, limit 1…1000 default 1000 и offset default 0 | [Page endpoint и пример](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L2466-L2525) |
| Embedded write имеет предел 1000 строк; отдельный ресурс нужен для больших документов | [Лимит](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L2466-L2467) |
| Отдельные Create/Update/Delete/batch endpoints существуют; их поддержка не следует из реализации документа | [Специальные ресурсы](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1039-L1063) |
| Документация допускает новые поля и новые ресурсы без нарушения совместимости | [Гарантии API](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L37-L42) |

Omission при Update сохраняет строки, а [] заменяет набор пустым — проектный вывод
из partial Update документа и общего правила полной замены. Null не объявляется
эквивалентом []. [Supply partial Update](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L2191-L2196).

## Общие поля и пределы общности

- Supply, demand, purchaseorder, invoicein и invoiceout имеют ассортиментную коммерческую
  основу: accountId/id, assortment, quantity, price, discount, vat и vatEnabled.
  Цена и количество — Float, discount — Float, vat — Int, vatEnabled — Boolean.
  Колонка «обязательно при ответе» сама по себе не делает поле обязательным request.
  [Supply](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L66-L86),
  [Demand](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L86-L107),
  [PurchaseOrder](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L60-L76),
  [InvoiceIn](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_invoice_in.md#L56-L72),
  [InvoiceOut](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_invoice_out.md#L55-L71).
- В таблицах id/accountId read-only; общее правило embedded Update при этом требует
  совпавший id для обновления существующей позиции. Поэтому id в request — узкий
  selector этой операции, а не право назначать ID создаваемой строке.
  [Поля позиции](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L71-L78),
  [Правило matching](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1031-L1037).
- Move не содержит discount/vat/vatEnabled в таблице позиции; содержит price,
  quantity, overhead Int read-only и sourceSlot/targetSlot metadata references.
  Упаковки и серийные номера — дополнительные поля, не основание автоматически
  включать все возможности в общий basic scope.
  [Movement positions](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_move.md#L62-L89).
- Purchaseorder добавляет shipped и inTransit; demand — service cost и другие поля;
  supply/demand имеют read-only overhead. Это различия document-owned форм, не общий
  writable contract всех позиций.
  [PurchaseOrder](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchaseOrder.md#L65-L76),
  [Demand](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_demand.md#L91-L107),
  [Supply](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_supply.md#L71-L86).
- Purchasereturn имеет коммерческие поля в чтении, но discount отмечен read-only;
  основанный на supply возврат и отдельное изменение позиции ограничены правилами
  возврата, включая quantity-only per-position Update. Это не совместимость с общим
  CommercialRequest для записи.
  [Return fields](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L52-L68),
  [Create](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L567-L583),
  [Position Update](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L2852-L2855).

## Проектные предложения, а не факты API

- Package positions, Service[T any]/GetPage, PageOptions и экспортируемые имена
  Core/Item/Commercial/Movement и Requests — PROPOSED интерфейс SDK.
- Generic T any и открытый documentType позволяют использовать другие shape без
  заявления, что их поля или write rules уже реализованы. Inventory/payroll указаны
  как проверенные примеры разнообразия API, их типы в этом срезе не проектируются.
- Минимальная response validation (meta.href и rows-array), UUID/version policy,
  default limit/offset и syntax documentType — выбранные для обсуждения правила SDK.
- Overhead и остальные денежные поля имеют float64 независимо от API Int/Float:
  опубликованные Supply/Demand примеры содержат overhead:0.0. Тестировать 0/0.0/дробь;
  точность денежных integer >2^53 не обещается.
- Старые unknown fields декодируются обычным transport; typed поддержка требует
  отдельного поля и проверок. Модуль не содержит закрытого списка document types.

## Comments
