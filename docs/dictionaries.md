# Справочники, настройки и метаданные

Источник контрактов — `api-remap-1.2-doc` revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b`.
Все методы принимают `opts ...transport.RequestOption` ([ADR 0001](adr/0001-request-options.md)).
Ссылки — `model.Reference` без expand; деньги — float64 в копейках, проценты — нет.
Одна операция — один запрос, записи не повторяются.

## Справочники

| Пакет / метод | Маршрут |
| --- | --- |
| `entity/store`: `Get`, `List`, `Create`, `Update`, `Delete`, массовые | `entity/store[/{id}]` |
| `store.Service`: `ListZones`, `GetZone`, `CreateZone`, `UpdateZone`, `DeleteZone`, `CreateManyZones`, `UpdateManyZones`, `DeleteManyZones`; то же для `Slots` | `entity/store/{id}/zones`, `…/slots` |
| `entity/productfolder`, `project`, `contract`, `saleschannel`, `organizationbranch`: `Get`, `List`, `Create`, `Update`, `Delete`, массовые | `entity/{тип}[/{id}]` |
| `entity/group`: `Get`, `List`, `Create`, `Update`, `Delete` | `entity/group[/{id}]` |
| `entity/region`, `saleplatform`, `thing`: `Get`, `List` (только чтение) | `entity/{тип}[/{id}]` |
| `entity/customentity`: `Create`, `Update`, `Delete`; `ListElements`, `GetElement`, `CreateElement`, `UpdateElement`, `DeleteElement` | `entity/customentity/{id}[/{elementID}]` |

- Store: parent и pathName записываются (маркера readonly нет, их передаёт пример PUT). Адрес — строкой
  или `addressFull`; null адреса источник не поддерживает, удаляет адрес `Address: ""`.
- Зоны и ячейки: `slot.zone` только читается — таблица помечает его «только для чтения», хотя
  источник даёт примеры создания и изменения ячейки с зоной; исключение не утверждено.
- Create/Update вложенных коллекций (зоны, ячейки, счета, контактные лица, события) принимают ответ-объект
  и массив из одного объекта: так выглядят примеры источника. Пустой или длиннее массив — ошибка.
- Contract: Create требует name, ownAgent и agent — по списку «Обязательные для создания Договора поля»;
  sum — копейки, rewardPercent — целый процент, rate — `model.DocumentRate`.
- OrganizationBranch: organization обязателен и после создания не меняется (его нет в UpdateRequest);
  КПП — в `mod__requisites__ru`, как в таблице и 8 из 9 примеров.
- SalesChannel: type обязателен; ProductFolder: archived, pathName и effectiveVat* только читаются.
- Group: массовых операций источник не описывает; `Index: 0` ставит отдел в начало списка.
- Пользовательский справочник: список справочников и их доп. поля читает `metadata`
  (`GetCompanySettings`, `GetCustomEntity`); id справочника — id его метаданных.

## Контрагент и юрлицо

| Метод | Маршрут |
| --- | --- |
| `counterparty.Service`: `ListAccounts`, `GetAccount`, `CreateAccount`, `UpdateAccount`, `DeleteAccount`, `DeleteManyAccounts` | `entity/counterparty/{id}/accounts` |
| `counterparty.Service`: `ListContactPersons`, `GetContactPerson`, `CreateContactPerson`, `UpdateContactPerson` | `…/contactpersons` |
| `counterparty.Service`: `ListNotes`, `GetNote`, `CreateNote`, `UpdateNote`, `DeleteNote` | `…/notes` |
| `organization.Service`: `ListAccounts`, `GetAccount`, `CreateAccount`, `UpdateAccount`, `DeleteAccount`, `CreateManyAccounts`, `UpdateManyAccounts`, `DeleteManyAccounts` | `entity/organization/{id}/accounts` |

- Счёт юрлица несёт `currency`; валюту существующего счёта меняет только сервер (ошибка при попытке).
  «Изменить счета Юрлица» — один POST массивом: элемент без meta создаёт счёт, с meta — меняет,
  непереданные поля не меняются.
- Удаления контактного лица источник не описывает. `agent` контактного лица не пишется: родителя
  задаёт путь запроса.
- Counterparty.Discounts — только чтение; `mod__requisites__uz/kz` у контрагента и юрлица —
  `model.RequisitesUZ/KZ`, валидны только у сущностей соответствующего региона (проверяет сервер).
- Organization: `companyVat__ru` (только при payerVat), `advancePaymentVat`, подписи и печать
  (`Image`: чтение — meta файла и миниатюра, запись — `ImageRequest{Filename, Content}` в Base64).
  `fsrarId` и `utmUrl` есть только в примерах — не поддерживаются.

## Настройки

| Метод | Маршрут |
| --- | --- |
| `counterparty.Service.GetSettings`, `UpdateSettings` | `entity/counterparty/settings` |
| `assortment.Service.GetSettings`, `UpdateSettings` | `entity/assortment/settings` |
| `contexts.Service.CompanySettings`, `UpdateCompanySettings` | `context/companysettings` |
| `contexts.Service.UserSettings`, `UpdateUserSettings` | `context/usersettings` |
| `contexts.Service.Subscription` | `accountSettings/subscription` |

- Изменение частичное: передаются только заданные поля, ответ — настройки целиком; пустой запрос —
  ошибка до HTTP.
- Пример настроек товаров пишет `checkUniqueCodeBoolean`; SDK следует таблице — `checkUniqueCode`.
- Настройки компании: пишутся только поля списка редактируемых; типы цен меняет `pricetype.ReplaceAll`.
- Настройки пользователя: `mailFooter` таблица называет Boolean, но пример и фикстура — текст; SDK — string.
  Ссылки Default* в примерах источника — Meta без обёртки, в фикстуре OpenAPI — `{"meta": …}`:
  чтение принимает обе формы (`SettingsRef`), запись — форма источника (`model.Meta`).
- Подписка без meta и id; `subscriptionEndDate` — Unix-время в миллисекундах.

## Сотрудник, комплект, точка продаж

- `employee.Service.Activate` (login нужен, если сотрудник не был активен; ответ
  `mailActivationRequired`), `Deactivate`, `ResetPassword` — `PUT entity/employee/{id}/access/*`.
  Запросы ролей — `role.Service.Preset`.
- `bundle.Service.AddComponent`, `UpdateComponent`, `DeleteComponent`, `DeleteManyComponents`.
- `retailstore.Service.ListCashiers`, `GetCashier` — кассиры, только чтение.

## Метаданные

- `metadata.Service.CreateAttribute` (name и type; для `customentity` — `CustomEntityMeta`),
  `UpdateAttribute` (type после создания не меняется), `DeleteAttribute`, `CreateManyAttributes`,
  `UpdateManyAttributes`, `DeleteManyAttributes`. Доп. поля услуг и комплектов — у `product`.
- `metadata.Service.CreateManyStates`, `UpdateManyStates`; `variant.Service.CreateManyCharacteristics`.

## Не входит

Expand; партия (consignment) — ждёт модели штрихкодов G01; отдельный метод
`context/usersettings/metadata` (источник не описывает); реквизиты юрлица `fsrarId`, `utmUrl`.
