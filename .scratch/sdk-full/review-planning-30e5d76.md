# Review новой нарезки SDK — 2026-11-03

HEAD: `30e5d769fb27ec7a19ffe47394ffd796120aff0b`.
Scope: untracked `.scratch/sdk-full/map.md`, пять `sdk-full-*/spec.md`,
88 tickets и `docs/status-backlog.md`: 95 Markdown-файлов.
Это review planning, не команда реализации или утверждение DRAFT.

Snapshot SHA-256 по отсортированным path + NUL + bytes + NUL:
`403983c325e0ba6f730f69bd616c997c4fc30d716eef12b06438364ee3266978`.
До создания этого отчёта hash повторно проверен: файлы не изменились.

Метод: независимые Standards/Spec reviewers и дополнительная проверка root.
Git diff пуст: planning ещё не добавлен в Git; проверялись перечисленные untracked files.

## Standards

### S1. P2 — определить владельца Trash для новых документов

`sdk-full-documents/issues/04-trash.md:4,14` задаёт незаблокированный O04
с Trash у всех документов. Карта (`sdk-full/map.md:27`) допускает выполнять
O04 раньше новых ресурсов. Однако новые документы, например O08 и M05,
перечисляют собственные методы без Trash. Получается неучтённый поздний retrofit.

Правило spec: завершённый вертикальный срез, shared dependency раньше зависимого.
Нужно определить единый способ интеграции: O04 для существующих документов плюс
Trash в новых tickets; либо общий механизм с явной обязанностью новых сервисов
подключить его; либо финальный retrofit с зависимостями. Просто наличие O04
не назначает внедрение во всех будущих сервисах.

### S2. P2 — согласовать интеграцию Consignment с G01

`sdk-full-catalog/spec.md:198` включает Consignment в общую модель штрихкодов
и поручает перенос G01. Но `issues/01-barcodes-packs.md:14` перечисляет только
Product/Variant/Service/Bundle; S10 не зависит от G01 и допускает временный TODO.
Можно завершить оба tickets, оставив этот TODO без владельца.

Правило: scope tickets должен обеспечивать обещанный spec результат; shared
dependency раньше зависимого. Нужен явный владелец последующей интеграции и
соответствующий порядок, а не обязательное расширение одних и тех же типов всем сущностям.

Проверено: 88 tickets, все локальные ссылки и dependency targets существуют,
циклов нет, табличные зависимости map/spec/issues совпадают.

## Spec

### P1. P1 — S16 задаёт неверный HTTP-контракт изменения счетов

`sdk-full-dictionaries/issues/16-organization-accounts.md:18`:
«PUT массивом заменяет все счета». В источнике описан POST коллекции accounts,
непереданные поля не изменяются. Неверны и HTTP method, и обещание полного замещения.

Исправить ticket и соответствующий раздел spec; отдельно закрепить неизменяемость
валюты существующего счёта.
Источник: [_organization.md:2323–2346](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_organization.md#L2323).

### P2. P1 — G13 запрещает существующий Update

`sdk-full-catalog/issues/13-retireorder.md:18` утверждает отсутствие Update.
Источник прямо описывает PUT entity/retireorder/{id} и массовое создание/обновление.
Добавить Update/UpdateMany; допустимость documentState проверяет сервер.
Источник: [_retireorder.md:1354–1380](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_retireorder.md#L1354).

### P3. P2 — S01 ошибочно запрещает запись родителя склада

`sdk-full-dictionaries/issues/01-store.md:18`: parent/path и pathName readonly.
У parent/pathName в таблице нет такого маркера; пример PUT передаёт оба поля.
Исправить Write matrix, иначе SDK лишится перемещения склада между группами.
Источник: [_store.md:31–32,1178–1220](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_store.md#L31).

### P4. P2 — G01 навязывает неподтверждённую общую форму packs

`sdk-full-catalog/issues/01-barcodes-packs.md:14` требует barcodes/packs у всех
Product/Variant/Service/Bundle. У Service/Bundle packs не описан. Variant использует
parentpack + barcodes, Product — uom + quantity. Общая Barcode уже есть в SDK.
Нужно разделить capability и JSON формы, переиспользовать существующую Barcode.
Источник: [_variant.md:222–228](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_variant.md#L222),
[_product.md:194–205](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_product.md#L194).

### P5. P2 — G03 теряет операции storebalances

`sdk-full-catalog/issues/03-minimumstock.md:14` ограничивается моделью minimumStock
и fields в List. Источник также описывает List/Get/Create/Update/Delete/DeleteMany
вложенного storebalances; другого владельца в tickets нет.
Источник: [_product.md:343–352](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_product.md#L343),
[_variant.md:211–220](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_variant.md#L211).

### P6. P2 — async отчётов не имеет ticket

`sdk-full-reports/spec.md:92–168` включает синхронные методы; создания async задач
для отчётов в tickets нет. Существующий SDK создаёт задачи только для counterparty
и assortment. Полное покрытие отчётных семейств требует async только у endpoints,
которым источник его разрешает, включая исключение выборочных показателей контрагентов.
Источник: [_async.md:6–15](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/_async.md#L6).

### P7. P2 — I12 требует отсутствующего transport механизма

`sdk-full-reports/issues/12-print.md:18` обещает вернуть Location и скачать через
Download. Текущий Client следует 303 автоматически и делает скрытый GET;
первоначальный foreign URL обычный Download отклоняет. В ticket нужно явно включить
остановку на 303, сохранение status/Location и безопасную явную загрузку результата,
не ослабляя same-origin контракт обычных JSON-вызовов. Источник не фиксирует origin
ссылки на файл: нельзя молча предположить, что любой Location уже поддержан Download.

Кроме того, источник называет вариант 202 будущим и не задаёт протокол опроса;
нельзя обещать готовый Wait по аналогии с async/{id}. Нужно вернуть различимые
результаты 202/303 и обозначить неизвестный контракт polling.
Источник: [_print.md:1–9](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_print.md#L1).
Код: transport/client.go:108–129, transport/download.go:13–17.

### P8. P2 — потеряны операции вложенных ресурсов

Явно назначить scope следующим tickets:

- G04: POST замены/упорядочивания списка images и fields=downloadPermanentHref —
  [_images.md:123,295](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_images.md#L123).
- S19: DeleteMany компонентов —
  [_bundle.md:1998](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_bundle.md#L1998).
- I19: batch create/update/delete notes —
  [_eventfeed.md:299,402](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_eventfeed.md#L299).

В What to build перечислены более узкие наборы; общая фраза «массовые только
если описаны» не добавляет отсутствующую операцию в scope конкретного ticket.

### P9. P2 — I04 задаёт неправильную форму byslot/current

`sdk-full-reports/issues/04-stock-operation-slot.md:18`: slot/store — ссылки.
Ответ byslot/current — plain array со строковыми UUID assortmentId/storeId/slotId,
без meta. Разделить формы двух endpoints, не использовать Reference по аналогии.
Источник: [_report_stock.md:960–979](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/reports/_report_stock.md#L960).

### Дополнительная проверка root

### P10. P2 — O02 не включает шаблон enter

`sdk-full-documents/issues/02-document-template.md:14` явно перечисляет существующие
документы, которым добавляется New, но enter отсутствует. Общая таблица шаблонов
описывает enter на основании inventory и ограничение 500 подходящих позиций.
В _enter.md отдельного раздела шаблона нет: включить операцию по общему источнику
либо явно зафиксировать расхождение для решения; не оставлять её без владельца.
Источник: [_common_info.md:24,48](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_common_info.md#L24).

### P11. P2 — O16 не включает GetPosition предоплаты

`sdk-full-documents/issues/16-prepayment.md:14` включает только GetPositions.
Источник отдельно описывает GET prepayment/{id}/positions/{positionID}.
O01 относится к существующим 13 документам, поэтому будущий Prepayment этим
ticket не закрывается; добавить read-only GetPosition в O16.
Источник: [_prepayment.md:652–669](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/documents/_prepayment.md#L652).

### P12. P2 — заявленная полнота не включает часть известных остатков SDK

Backlog ссылается на TODO существующего SDK, но в 88 tickets нет владельца для:

- Organization: companyVat__ru, advancePaymentVat, chiefAccountSign, directorSign,
  stamp и региональных реквизитов, отложенных в docs/organization.md:105–109.
  S16 явно ограничен accounts/fsrarId/utmUrl. По крайней мере НДС/подписи/печать
  присутствуют в таблице и примерах источника, это не гипотетические поля:
  [_organization.md:45–73,108,616–629](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_organization.md#L45).
- Массовое создание/обновление states и массовое создание characteristics:
  docs/write-async.md:71 оставляет их вне текущего SDK; S24 — только attributes.
  Источник: [_states.md:213](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_states.md#L213),
  [_characteristic.md:105](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_characteristic.md#L105).

Добавить явные tickets/подзадачи либо сузить заявление о полном покрытии.
Это не требование заранее выписать все поля в DRAFT: речь об отсутствующем scope,
который последующий freeze существующих tickets сам по себе не добавит.

## Verification и ограничения

- 95 planning files, 88 tickets; local links, dependency targets, согласованность
  dependency tables и отсутствие циклов — проверены Standards reviewer.
- SHA-256 исходного planning snapshot до отчёта совпал с начальным.
- Первичные тексты проверялись из локального архива revision 9a811cef.
  Через официальный [GitHub compare API](https://api.github.com/repos/moysklad/api-remap-1.2-doc/compare/9a811cef3dc4a3f3a46f070e60452361a06bf7cd...0e1b3d61fb8d8fa314469361d96560d6c3f7f57b)
  подтверждено: среди md изменены лишь labels/print, добавлен язык text в code fences;
  контракт и приведённые выводы от перехода на 0e1b3d61 не меняются.
- Go tests не запускались: новых изменений Go code нет; предыдущий review на HEAD
  проверил async race/vet, полный make check на 565cfa3 был PASS.
- Live API не вызывался. Новый planning, tracker statuses, code и commits не менялись.
  Создан только этот отчёт; .DS_Store не читались и не изменялись.
- Это не утверждение исчерпывающей field/method matrix всего API. Обнаруженных
  расхождений достаточно, чтобы сначала исправить план и сверить оставшийся scope.

Standards: 2 findings P2. Spec: 12 findings, максимум P1.
Разбиение на пять направлений пригодно; tickets требуют правок до реализации.
