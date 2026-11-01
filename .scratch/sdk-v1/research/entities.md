# SDK v1: контракты выбранных сущностей

Исследовано по official api-remap-1.2-doc, pinned revision
`9a811cef3dc4a3f3a46f070e60452361a06bf7cd`. Имена полей и явные markers
сведены в [каталог контрактов](../contracts.md); это не готовая Write matrix.
Исходные Markdown получены с official GitHub, upstream fixtures SDK не менялись.
Ниже — факты источника и отдельно рекомендации для v1, без утверждения runtime outcomes.

## Основные методы и исключения

| Ресурс | Подтверждённая форма interface | Create / особенности |
| --- | --- | --- |
| Employee | GET list/id, POST, PUT id | required lastName |
| Role | GET entity/role/id, GET list, POST, PUT id | name; meta.type в примере customrole |
| Currency | GET list/id, POST, PUT id | name/code/isoCode; заданный rate не нулевой |
| PriceType | GET array/list, GET id, GET default; POST whole-list | не обычные Create/Update |
| TaxRate | GET list/id, POST, PUT id | field/required table и метод |
| Service | GET list/id, POST, PUT id | name; отдельный ресурс от Product |
| Bundle | GET list/id, POST, PUT id; собственные components | components не документные positions |
| Variant | GET list/id, POST, PUT id; variant metadata characteristics | product link/characteristics по методу |
| Assortment | GET mixed list/search | не универсальный write endpoint |
| Characteristic | GET variant/metadata и definition/id; POST characteristics | отдельного PUT в источнике нет |
| CustomerOrder | GET list/id, POST, PUT id, GET positions | organization, agent |
| Enter/Loss | GET list/id, POST, PUT id, GET positions | organization, store |
| SalesReturn | GET list/id, POST, PUT id, GET positions | organization/store/agent; demand только based |
| RetailStore | GET list/id, POST, PUT id | name/org/store/priceType |
| RetailShift | GET list/id, POST, PUT id | organization, retailStore; не товарный GetPositions |
| RetailDemand | GET list/id, POST, PUT id, GET positions | минимально retailShift |
| RetailSalesReturn | GET list/id, POST, PUT id, GET positions | own based/unbased required matrix |
| BonusProgram | GET list/id, POST, PUT id | own required bool/rate fields |
| BonusTransaction | GET list/id, POST, PUT id | agent/bonusProgram/transactionType |

Источники: [Employee](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md),
[Role](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_custom_role.md),
[Currency](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_currency.md),
[TaxRate](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_taxrate.md).
Остальные строки ссылаются на собственные источники в разделах ниже и contracts.md.

## PriceType — полный список, а не PUT сущности

`GET context/companysettings/pricetype` возвращает plain JSON array. По ID и по
`default` существуют отдельные GET. Редактирование выполняется POST полного списка:
порядок сохраняется, первый тип становится default, исключение элемента удаляет его.
Разрешено 1…100 типов; все типы удалить нельзя. SDK не должен делать скрытый
read-modify-write или обещать защиту от конкурентного редактирования списка.
Существующий `model.PriceType` уже используется в SalePrice: новый несовместимый
одноимённый публичный тип без необходимости не нужен.
[Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_pricetypes.md).

## Assortment и четыре товарных пакета

Assortment объединяет Product, Service, Bundle, Variant **и Consignment**.
Один search/List имеет серверную общую выдачу и пагинацию; четыре отдельных List
не заменяют её без склейки, нескольких requests и потери общего порядка поиска.
Специализированные пакеты остаются нужны для записи и своих полей.

Облегчённый endpoint использует тот же путь `entity/assortment`, но требует header
`X-Lognex-Remap-Beta-Feature: assortmentWithoutStock`. Без header работает legacy.
Beta не возвращает stock/reserve/inTransit/quantity, меняет default groupBy и
поддерживаемые filters; нельзя молча подменять режим или догружать stock.
Рекомендация v1: explicit mode и lossless unknown MetaType; Consignment common-read
не равен обещанию полноценного Consignment CRUD.
[Источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_assortment.md).

## Variant, Bundle и общие формы

Characteristic definition находится в `entity/variant/metadata/characteristics`;
это не отдельный generic entity/characteristic CRUD. Значение в Variant и metadata
definition имеют разные роли. Bundle components имеют свои JSON fields и маршруты;
они не становятся Position только из-за assortment/quantity.
[Characteristic](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_characteristic.md),
[Variant](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_variant.md),
[Bundle](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bundle.md),
[Service](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_service.md).

## Role и Employee security

Права читаются/меняются через `entity/employee/{id}/security`; preset roles доступны
по role/owner, admin, individual, cashier, worker. Activation/deactivation/password
и назначение прав — разные операции, не side effect обычного Employee.Update.
Role допускает admin или App по описанию источника. Таблица называет permissions
Array(Object), но примеры показывают неоднородный object bool/string/nested scopes.
Это реальное противоречие, которое A02 должен разрешить до field freeze. Рекомендация:
сохранять unknown JSON losslessly, не использовать map[string]any с float64.
[Employee security](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_employee.md#L788),
[Role](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_custom_role.md).

## Документы и required Create

CustomerOrder требует organization/agent, Enter/Loss — organization/store.
SalesReturn требует organization/store/agent; demand нужен только при создании
по основанию. Ограничения совпадения сторон/валюты/строк и Update должны следовать
собственному методу SalesReturn, а не копии PurchaseReturn. В SDK не нужен pre-GET.
[CustomerOrder](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_customerOrder.md),
[Enter](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_enter.md),
[Loss](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_loss.md),
[SalesReturn](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_sales_return.md).

## Розница

RetailStore Create требует name/organization/store/priceType; тариф и несовместимость
allowCreateProducts/controlShippingStock проверяет API. RetailShift Create требует
organization/retailStore. Минимальный RetailDemand Create — retailShift; смена должна
быть активной, moment позже moment смены. Нельзя переносить required response поля
в обязательные Create references по сходству с обычной Demand.

RetailSalesReturn имеет свой набор name/organization/store/retailStore/retailShift/agent,
demand — только для based (ссылка именно на RetailDemand); для unbased источник
дополнительно требует cashSum/noCashSum. Явный 0 не равен omission. Фискальные поля,
суммы и update restrictions требуют отдельной поимённой матрицы, не эвристики CRUD.
[RetailStore](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_retailstore.md),
[RetailShift](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retailshift.md),
[RetailDemand](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retaildemand.md),
[RetailSalesReturn](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_retail_sales_return.md).

## Бонусы и ограничения токена

BonusProgram required Create: name/active/allProducts/allAgents/
earnRateRoublesToPoint/spendRatePointsToRouble/maxPaidRatePercents.
BonusTransaction required Create: agent/bonusProgram/transactionType. Bonus points
не деньги в копейках. Состояния и разрешённые изменения bonusValue/parentDocument
не следует угадывать по похожести с финансовыми документами.
[BonusProgram](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bonus_program.md),
[BonusTransaction](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_bonus_operation.md).

Наличие POST в JSON API не доказывает разрешение токена решения. Vendor ограничения
и фактический ответ такого токена — отдельная проверка capability; SDK не подменяет
его персональным токеном и не имитирует сервер. В экосистеме уже есть отдельный
[Draft/live research ticket](../../../../meta/.scratch/idea-selection/issues/18-mobile-window-test.md).
