Status: needs-triage
Stage: DRAFT

# 1. Справочники, настройки и метаданные

## Problem Statement

Закрыть оставшиеся справочники JSON API, вложенные ресурсы и настройки уже реализованных сущностей, CRUD описаний доп. полей и пользовательские справочники.

## Solution

Завершённые вертикальные срезы; каждый ticket включает models, methods, тесты и собственное
описание полей. Матрица ниже определяет scope этой захватки.
[Общая карта](../sdk-full/map.md), [нормативные ограничения v1](../sdk-v1/spec.md),
[каталог JSON-контрактов](../sdk-v1/contracts.md), [backlog](../../docs/status-backlog.md).
Источник: api-remap-1.2-doc revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.

## User Stories

1. Как разработчик решения, я могу использовать «Store Get/List/Create/Update + zones/slots» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
2. Как разработчик решения, я могу использовать «ProductFolder Get/List/Create/Update» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
3. Как разработчик решения, я могу использовать «Group Get/List/Create/Update» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
4. Как разработчик решения, я могу использовать «Project Get/List/Create/Update» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
5. Как разработчик решения, я могу использовать «Contract Get/List/Create/Update» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
6. Как разработчик решения, я могу использовать «Region Get/List (read-only)» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
7. Как разработчик решения, я могу использовать «SalesChannel Get/List/Create/Update» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
8. Как разработчик решения, я могу использовать «SalePlatform Get/List (read-only)» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
9. Как разработчик решения, я могу использовать «OrganizationBranch Get/List/Create/Update» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
10. Как разработчик решения, я могу использовать «Consignment Get/List/Create/Update» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
11. Как разработчик решения, я могу использовать «Thing Get/List (read-only)» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
12. Как разработчик решения, я могу использовать «Counterparty accounts» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
13. Как разработчик решения, я могу использовать «Counterparty contactpersons» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
14. Как разработчик решения, я могу использовать «Counterparty notes» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
15. Как разработчик решения, я могу использовать «Counterparty: discounts и региональные реквизиты» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
16. Как разработчик решения, я могу использовать «Organization accounts, fsrarId, utmUrl» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
17. Как разработчик решения, я могу использовать «Настройки справочников: counterparty и assortment» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
18. Как разработчик решения, я могу использовать «CompanySettings GET/PUT» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
19. Как разработчик решения, я могу использовать «Bundle: мутации компонентов» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
20. Как разработчик решения, я могу использовать «Employee: активация, деактивация, сброс пароля» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
21. Как разработчик решения, я могу использовать «RetailStore cashiers (read-only)» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
22. Как разработчик решения, я могу использовать «UserSettings GET/PUT» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
23. Как разработчик решения, я могу использовать «CompanySubscription GET» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
24. Как разработчик решения, я могу использовать «Доп. поля: CRUD описаний» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
25. Как разработчик решения, я могу использовать «CustomEntity и элементы» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
26. Как разработчик решения, я могу использовать «Organization: НДС, подписи, печать, региональные реквизиты» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.
27. Как разработчик решения, я могу использовать «Массовые операции states и characteristics» через публичный SDK
   и получать фактический ответ API с документированными ограничениями.

## Implementation Decisions и scope tickets

| ID | Ticket | Blocked by |
| --- | --- | --- |
| S01 | [Store Get/List/Create/Update + zones/slots](issues/01-store.md) | нет |
| S02 | [ProductFolder Get/List/Create/Update](issues/02-productfolder.md) | нет |
| S03 | [Group Get/List/Create/Update](issues/03-group.md) | нет |
| S04 | [Project Get/List/Create/Update](issues/04-project.md) | нет |
| S05 | [Contract Get/List/Create/Update](issues/05-contract.md) | нет |
| S06 | [Region Get/List (read-only)](issues/06-region.md) | нет |
| S07 | [SalesChannel Get/List/Create/Update](issues/07-saleschannel.md) | нет |
| S08 | [SalePlatform Get/List (read-only)](issues/08-saleplatform.md) | нет |
| S09 | [OrganizationBranch Get/List/Create/Update](issues/09-organizationbranch.md) | нет |
| S10 | [Consignment Get/List/Create/Update](issues/10-consignment.md) | G01 |
| S11 | [Thing Get/List (read-only)](issues/11-thing.md) | нет |
| S12 | [Counterparty accounts](issues/12-counterparty-accounts.md) | нет |
| S13 | [Counterparty contactpersons](issues/13-counterparty-contactpersons.md) | нет |
| S14 | [Counterparty notes](issues/14-counterparty-notes.md) | нет |
| S15 | [Counterparty: discounts и региональные реквизиты](issues/15-counterparty-fields.md) | нет |
| S16 | [Organization accounts, fsrarId, utmUrl](issues/16-organization-accounts.md) | S12 |
| S17 | [Настройки справочников: counterparty и assortment](issues/17-dictionary-settings.md) | нет |
| S18 | [CompanySettings GET/PUT](issues/18-companysettings.md) | нет |
| S19 | [Bundle: мутации компонентов](issues/19-bundle-components.md) | нет |
| S20 | [Employee: активация, деактивация, сброс пароля](issues/20-employee-access.md) | нет |
| S21 | [RetailStore cashiers (read-only)](issues/21-cashier.md) | нет |
| S22 | [UserSettings GET/PUT](issues/22-usersettings.md) | нет |
| S23 | [CompanySubscription GET](issues/23-companysubscription.md) | нет |
| S24 | [Доп. поля: CRUD описаний](issues/24-attributes-crud.md) | нет |
| S25 | [CustomEntity и элементы](issues/25-customentity.md) | S24 |
| S26 | [Organization: НДС, подписи, печать, региональные реквизиты](issues/26-organization-fields.md) | S16 |
| S27 | [Массовые операции states и characteristics](issues/27-metadata-batch.md) | нет |

### S01: Store Get/List/Create/Update + zones/slots

Склад: Get/List/Create/Update/Delete, массовые. Вложенные zones и slots: List/Get/Create/Update/Delete и массовые create/update/delete — отдельные методы сервиса склада.

Проверки и особенности: Address/addressFull как у Organization; parent (группа) и pathName записываются — маркера readonly нет, пример PUT передаёт оба; slot ссылается на zone; ID склада в пути вложенных ресурсов проверяется до HTTP.

Источники: [store](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_store.md).

### S02: ProductFolder Get/List/Create/Update

Группа товаров: Get/List/Create/Update/Delete, массовые.

Проверки и особенности: productFolder (родитель) и pathName; налоговые поля (vat, vatEnabled, useParentVat, effectiveVat, taxSystem) — read/write по таблице; удаление непустой группы проверяет сервер.

Источники: [productfolder](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_productFolder.md).

### S03: Group Get/List/Create/Update

Отдел: Get/List/Create/Update/Delete. Массовых операций в источнике нет.

Проверки и особенности: Ограничения прав на создание отделов — сервер; index/sort по таблице.

Источники: [group](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_group.md).

### S04: Project Get/List/Create/Update

Проект: Get/List/Create/Update/Delete, массовые, attributes.

Проверки и особенности: Attribute через общую модель; archived — явный false.

Источники: [project](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_project.md).

### S05: Contract Get/List/Create/Update

Договор: Get/List/Create/Update/Delete, массовые, attributes.

Проверки и особенности: contractType/rewardType — открытые string; rewardPercent — процент, не деньги; sum — копейки; ownAgent/agent/ownAgentAccount — ссылки без pre-GET; rate через DocumentRate.

Источники: [contract](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_contract.md).

### S06: Region Get/List (read-only)

Регион: Get/List. Записи в источнике нет.

Проверки и особенности: code — строка; без Create/Update.

Источники: [region](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_region.md).

### S07: SalesChannel Get/List/Create/Update

Канал продаж: Get/List/Create/Update/Delete, массовые.

Проверки и особенности: type — открытый enum; system-каналы изменять запрещает сервер.

Источники: [saleschannel](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_saleschannel.md).

### S08: SalePlatform Get/List (read-only)

Площадка для продаж: Get/List.

Проверки и особенности: marketplace/type — открытые значения; без записи.

Источники: [saleplatform](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_saleplatform.md).

### S09: OrganizationBranch Get/List/Create/Update

Подразделение юрлица: Get/List/Create/Update/Delete, массовые.

Проверки и особенности: organization — обязательная ссылка; реквизиты подразделения — по таблице, без копирования полей Organization по аналогии.

Источники: [organizationbranch](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_organizationbranch.md).

### S10: Consignment Get/List/Create/Update

Партия: Get/List/Create/Update/Delete, массовые, attributes.

Проверки и особенности: assortment — обязательная ссылка; label; barcodes — существующая model.Barcode, порядок после G01 только ради одного владельца формы штрихкодов.

Источники: [consignment](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_consignment.md).

### S11: Thing Get/List (read-only)

Серийный номер: Get/List.

Проверки и особенности: Только чтение; фильтры по assortment — через существующий query.

Источники: [thing](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_thing.md).

### S12: Counterparty accounts

Счета контрагента: List/Get/Create/Update/Delete и массовое удаление (entity/counterparty/{id}/accounts).

Проверки и особенности: isDefault; accountNumber/bic/bankName; Account — общая модель для S16.

Источники: [counterparty](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_counterparty.md).

### S13: Counterparty contactpersons

Контактные лица: List/Get/Create/Update.

Проверки и особенности: Delete в источнике нет — не добавлять.

Источники: [counterparty](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_counterparty.md).

### S14: Counterparty notes

События контрагента: List/Get/Create/Update/Delete (entity/counterparty/{id}/notes).

Проверки и особенности: author/authorApplication readonly; description обязателен.

Источники: [counterparty](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_counterparty.md).

### S15: Counterparty: discounts и региональные реквизиты

Поля discounts (read) и `mod__requisites__uz`/`mod__requisites__kz` с metadata — чтение и запись по таблице.

Проверки и особенности: discounts — разные типы скидок, открытый meta.type; региональные реквизиты не смешиваются с российскими.

Источники: [counterparty](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_counterparty.md).

### S16: Organization accounts, fsrarId, utmUrl

Счета юрлица: List/Get/Create/Update/Delete, массовое удаление; массовое изменение — POST массивом в entity/organization/{id}/accounts. Решить fsrarId/utmUrl (есть в примерах, нет в таблице).

Проверки и особенности: POST массивом: непереданные поля не изменяются, полного замещения нет. Валюту существующего счёта менять нельзя — ошибку возвращает сервер, SDK не проверяет; accounts-in-currency.

Источники: [organization](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_organization.md).

### S17: Настройки справочников: counterparty и assortment

GET/PUT entity/counterparty/settings и entity/assortment/settings.

Проверки и особенности: Одиночный ресурс без ID; uniqueCodeRules/barcodeRules и т.п. — по таблицам.

Источники: [counterparty](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_counterparty.md), [assortment](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_assortment.md).

### S18: CompanySettings GET/PUT

Настройки компании: GET/PUT context/companysettings (сейчас только metadata).

Проверки и особенности: currency, priceTypes (readonly — есть отдельный pricetype), discountStrategy, globalOperationNumbering и прочие флаги — явный false.

Источники: [companysettings](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_companysettings.md).

### S19: Bundle: мутации компонентов

Добавить/изменить/удалить компонент и DeleteMany компонентов (entity/bundle/{id}/components).

Проверки и особенности: quantity — не деньги; assortment — ссылка; GetComponents остаётся как есть.

Источники: [bundle](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_bundle.md).

### S20: Employee: активация, деактивация, сброс пароля

Activate/Deactivate/ResetPassword — PUT entity/employee/{id}/access/activate|deactivate|resetpassword. Запросы ролей (entity/role/owner и т.д.) уже покрыты role.Preset.

Проверки и особенности: Побочные эффекты (письма, тариф) — сервер; никаких повторов записи; ответы по источнику.

Источники: [employee](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_employee.md).

### S21: RetailStore cashiers (read-only)

Кассиры точки продаж: List/Get (entity/retailstore/{id}/cashiers).

Проверки и особенности: Только чтение; employee и retailStore — ссылки.

Источники: [cashier](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_cashier.md).

### S22: UserSettings GET/PUT

Настройки пользователя: GET/PUT context/usersettings и metadata.

Проверки и особенности: Значения enum открытые; defaultScreen, printFormat, locale по таблице.

Источники: [usersettings](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_usersettings.md).

### S23: CompanySubscription GET

Подписка компании: GET accountSettings/subscription.

Проверки и особенности: Только чтение; даты DateTime; тариф — открытый string.

Источники: [companysubscription](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_companysubscription.md).

### S24: Доп. поля: CRUD описаний

Create/Update/Delete и массовые операции описаний доп. полей entity/{type}/metadata/attributes; типы файл и справочник (customentity).

Проверки и особенности: Тип поля после создания не меняется — сервер; список допустимых типов сущности — открытый.

Источники: [workbook attributes](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/workbook/_workbook_attributes_crud.md), [general](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/_general.md).

### S25: CustomEntity и элементы

Пользовательский справочник: Create/Update/Delete, metadata; элементы: List/Get/Create/Update/Delete (entity/customentity/{id}).

Проверки и особенности: Связь с metadata.GetCustomEntity; ID справочника в пути элементов проверяется до HTTP.

Источники: [customentity](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_customentity.md).

### S26: Organization: НДС, подписи, печать, региональные реквизиты

Поля companyVat__ru, advancePaymentVat, chiefAccountSign, directorSign, stamp и `mod__requisites__uz`/`mod__requisites__kz` с metadata — по таблице (docs/organization.md «Границы»).

Проверки и особенности: Подписи и печать — объекты изображения (content base64 в записи, meta в ответе); региональные реквизиты — та же форма, что в S15.

Источники: [organization](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_organization.md).

### S27: Массовые операции states и characteristics

Массовое создание/обновление статусов (entity/{type}/metadata/states) и массовое создание характеристик модификаций — сейчас вне SDK (docs/write-async.md).

Проверки и особенности: ≤1000 элементов до HTTP, без разбиения и повторов — как batch X02.

Источники: [states](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_states.md), [characteristic](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_characteristic.md).

## Общий контракт реализации

- Перед кодом — field/method freeze: поимённая Write matrix, readonly/TODO,
  null/required/Create-only/Update-only; status needs-triage сам не разрешает реализацию.
- Образец пакета и тестов указан в каждом ticket; новый пакет повторяет его структуру
  (model.go, request.go, service.go, *_test.go), без новых общих абстракций.
- Delete, CreateMany/UpdateMany/DeleteMany — только если раздел источника их описывает,
  по образцу [optional X01/X02](../sdk-v1-optional/spec.md).
- Новый ресурс получает handle фасада `moysklad.Client`, строку в
  [status-entities](../../docs/status-entities.md) / [status-technical](../../docs/status-technical.md)
  и проходит `status_test.go`; пункт снимается из [status-backlog](../../docs/status-backlog.md).
- JSON: encoding/json v1 и v2, Field+omitzero для optional request, независимые
  golden, неизменённые upstream с SHA-256 provenance; raw JSON long не проходит float64.
- Dates — существующий DateTime; деньги float64 в копейках, counts/проценты/баллы/курс
  не объявляются деньгами. Reference без expand, unknown enum/type не отвергается.
- HTTP: настоящий transport+httptest, exact method/path/query/headers/body,
  plain/gzip, фактический server response, ошибки/auth/context/cancellation,
  errors.Is/As и единый moysklad:. Обычный вызов — одна отправка, если caller
  явно не включил retry/walker; каждый дополнительный запрос доказуем тестами.
- Бизнес-ограничения и permissions оставлены API; SDK не выполняет pre-GET ради
  local validation. Exception — только явный пользовательский multi-call operation.
- Проверки: `make check` (go test -race, go vet ./... ./testdata, gofmt, SHA-256 фикстур,
  статус сущностей, diff --check); независимые Standards/Spec review на фиксированном BASE.
- Коммит по ticket; shared dependency раньше зависимого, публичные docs отдельным
  завершающим commit волны. Planning не коммитить без отдельной команды.
- Live API, credentials, push/tag/release и worktree cleanup не входят.

## Open Decisions

- Утверждено пользователем 2026-11-03: исключение для `slot.zone` (S01) — пишется в Create и Update
  вопреки пометке «только для чтения» таблицы `_store.md`, по отдельным примерам источника. Общее правило
  «readonly сильнее примеров» для остальных полей сохраняется.

- gtd отдельного endpoint не имеет; модель `positions.GTD` уже есть — отдельного ticket нет.
- Файлы контрагента и товара покрывает существующий пакет `files` (entity/{type}/{id}/files); отдельных tickets нет.
- Счета контрагента и юрлица — одна модель Account или две? Решить в S12, S16 переиспользует.
- Read-only справочники (region, saleplatform, thing, cashier, companysubscription) не получают Create/Update по аналогии.

## Out of Scope

Expand; управление типами цен контрагента; администрирование кассира сверх GET. Все не перечисленные в tickets ресурсы/операции.
Этот DRAFT не является командой реализации или публикации.
