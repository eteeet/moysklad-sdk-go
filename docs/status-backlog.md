# Нереализованное: сущности и технические артефакты

Источник: api-remap-1.2-doc HEAD `0e1b3d61` (2026-10-30). С revision `9a811cef`
изменились только `_labels.md` и `_print.md` (разметка блоков кода, без смены контракта).
Сверено с [status-entities.md](status-entities.md), [status-technical.md](status-technical.md),
кодом пакетов и разделами «Ограничения и TODO» в docs/. Тесты `go test ./...` зелёные,
отметки `[x]` в обоих списках совпадают с кодом.

Пять захваток упорядочены по пользе и зависимостям, expand — последним.
Spec и tickets (DRAFT, needs-triage): [карта](../.scratch/sdk-full/map.md).

Уже покрыто, tickets не нужны: gtd (`positions.GTD`), файлы товара, контрагента и документов
(пакет `files`), запросы ролей сотрудника (`role.Preset`).

## Захватка 1. Справочники, настройки, метаданные — сделано — [spec](../.scratch/sdk-full-dictionaries/spec.md)

Простые CRUD по образцу uom:

- [x] store — склад; вложенные slots и zones (CRUD + массовые)
- [x] productfolder — группа товаров
- [x] group — отдел
- [x] project — проект
- [x] contract — договор
- [x] region — регион
- [x] saleschannel — канал продаж
- [x] saleplatform — площадка для продаж
- [x] organizationbranch — подразделение юрлица
- [x] consignment — партия
- [x] thing — серийный номер

Вложенные ресурсы и настройки существующих сущностей (нет в status-* — найдено по документации):

- [x] counterparty: accounts, contactpersons, notes (события), discounts,
      `mod__requisites__uz/kz`
- [x] counterparty/settings — настройки справочника контрагентов (GET/PUT)
- [x] organization: accounts (CRUD, массовое изменение POST, массовое удаление); fsrarId, utmUrl
- [x] organization: companyVat__ru, advancePaymentVat, chiefAccountSign, directorSign, stamp,
      `mod__requisites__uz/kz`
- [x] assortment/settings — настройки справочника товаров (GET/PUT)
- [x] companysettings — сами настройки компании GET/PUT context/companysettings
      (сейчас есть только metadata)
- [x] bundle: мутации компонентов (добавить/изменить/удалить, массовое удаление)
- [x] employee: активация, деактивация, сброс пароля
- [x] retailstore/{id}/cashiers — кассир
- [x] usersettings — настройки пользователя (context)
- [x] companysubscription — подписка компании

Метаданные и расширения:

- [x] Доп. поля: CRUD описаний, типы файл/справочник
- [x] customentity — пользовательский справочник и его элементы
- [x] Массовое создание/обновление states и массовое создание characteristics

## Захватка 2. Товары, скидки, маркировка — [spec](../.scratch/sdk-full-catalog/spec.md)

Product до полной модели:

- [x] штрихкоды, упаковки, алкоголь, отслеживание, прочие атрибуты
- [x] неснижаемый остаток: модель, fields=minimumStock и storebalances у Product и Variant
- [x] images — изображения (product/variant/bundle/contentcard)
- [x] contentcard — карточка контента

Скидки:

- [x] discount — общий список скидок
- [x] accumulationdiscount — накопительная скидка
- [x] personaldiscount — персональная скидка
- [x] specialpricediscount — скидка на спеццену
- [x] roundoffdiscount — округление копеек (PUT и массовое) — нет в status-entities

Маркировка:

- [x] trackingcode — коды маркировки
- [x] aggregatepack — транспортная упаковка (+ positions)
- [x] emissionorder — заказ кодов маркировки
- [x] retireorder — вывод из оборота

## Захватка 3. Документы и операции документов — [spec](../.scratch/sdk-full-documents/spec.md)

Общие механизмы (делать первыми — новые документы их используют):

- [x] Отдельные операции позиций: Get, Create, Update, Delete, массовое удаление
- [x] Шаблон документа: PUT entity/{type}/new
- [x] Автозаполнение: wizard/demand, wizard/salesreturn (цены, скидки, НДС, себестоимость)
- [x] Удаление в корзину (entity/{type}/{id}/trash); новые документы — в своих tickets
- [x] Запись связанных документов, сейчас read-only: payments, supplies, demands,
      factureOut, returns, productionTasks, invoicesOut
- [x] Поля позиций, оставшиеся TODO: pack, things, slot, serial/tracking у invoicein,
      invoiceout, move; transport-поля и shippingDocument* у demand
- [ ] agentAccount/organizationAccount у cashin/cashout (расхождение источников)

Торговые и финансовые документы:

- [x] internalorder — внутренний заказ
- [x] inventory — инвентаризация
- [x] facturein — счёт-фактура полученный
- [x] factureout — счёт-фактура выданный
- [x] commissionreportin — полученный отчёт комиссионера
- [x] commissionreportout — выданный отчёт комиссионера
- [x] counterpartyadjustment — корректировка взаиморасчётов
- [x] pricelist — прайс-лист
- [x] prepayment — предоплата
- [x] prepaymentreturn — возврат предоплаты
- [x] retaildrawercashin — внесение денег
- [x] retaildrawercashout — выплата денег
- [x] payroll — начисление зарплаты

## Захватка 4. Производство — [spec](../.scratch/sdk-full-production/spec.md)

- [x] processingplan — техкарта
- [x] processingplanfolder — группа техкарт
- [x] processingprocess — техпроцесс (+ positions)
- [x] processingstage — этап производства
- [x] processing — техоперация
- [x] processingorder — заказ на производство
- [x] productiontask — производственное задание
- [x] productionstagecompletion — выполнение этапа производства

## Захватка 5. Отчёты, печать, интеграции — [spec](../.scratch/sdk-full-reports/spec.md)

Отчёты:

- [x] stock — all, all/current, bystore, bystore/current, byoperation, byslot/current
- [x] Async создания отчётов (остатки, прибыльность, деньги, продажи/заказы, контрагенты, показатели)
- [x] byoperations — stock/reserve/intransit
- [x] dashboard — day/week/month
- [x] report/counterparty — показатели контрагентов
- [x] money — plotseries/byaccount
- [x] profit — byproduct/byvariant/byemployee/bycounterparty/bysaleschannel
- [x] orders — plotseries
- [x] sales — plotseries

Печать и публикации:

- [x] template — шаблоны печатных форм; metadata customtemplate/embeddedtemplate
- [x] Печать документов: export (303 → Location файла без автоследования; 202 — будущий,
      протокол опроса не описан)
- [x] labels — этикетки и ценники
- [x] Публикации документов: List/Create/Get/Delete

Интеграции и уведомления:

- [x] webhook — вебхук
- [x] webhookstock — вебхук остатков
- [x] notification — лента, markasread, markasreadall
- [x] notification/settings — настройки уведомлений
- [x] eventfeed — лента событий
- [x] task — задача (+ notes)
- [x] namedfilter — сохранённые фильтры (сущность; применение уже есть)
- [x] application — GET entity/application/{id} (контекст решения уже есть)

## Expand — последним, только по команде

- [ ] expand: limit≤100, глубина ≤3; раскрытые объекты вместо Reference
- [ ] fields в List справочников и документов (сейчас одно значение, не везде)
- [ ] UOM: expand

## Вне SDK-кода (выпуск и инфраструктура)

- [ ] Live read-only проверка на тестовом аккаунте (только с согласия)
- [ ] Проверка Product.Create/Update на тестовом товаре (только с согласия)
- [ ] CI workflow и проверки безопасности (для v1 не делается)
- [ ] Тег и релиз — отдельной командой
- [ ] Vendor API — отдельный контур
- [ ] Обновить revision источника в status-* до `0e1b3d61` (контракт не менялся)
