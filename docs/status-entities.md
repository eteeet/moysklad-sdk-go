# Статус сущностей JSON API 1.2

Источник: api-remap-1.2-doc, revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`.
[x] — есть пакет с Get/Create/Update (документы с позициями — ещё GetPositions) и handle фасада
`moysklad.Client` (у download — `Transport().Download`); иной набор методов указан в строке. `status_test.go` проверяет, что отметки и
handles фасада совпадают; набор методов — тесты пакетов.
List/Delete/expand/batch — см. [status-technical.md](status-technical.md) и [write-async.md](write-async.md).

Итого: справочники 47/60, документы 31/37, отчёты 1/9, прочее 6/8.

## Справочники

- [x] accumulationdiscount — накопительная скидка: Get/Create/Update/Delete (список — Discount().List)
- [x] aggregatepack — транспортная упаковка и позиции
- [ ] application — серверные решения
- [x] assortment — ассортимент: List и All (assortmentWithoutStock на каждой странице) и ListLegacy
- [x] bonusprogram — бонусная программа
- [x] bonustransaction — бонусная операция
- [x] bundle — комплект: Get/List/Create/Update, GetComponents/GetComponent
- [x] cashier — кассир: RetailStore().ListCashiers/GetCashier (только чтение)
- [x] characteristic — характеристики модификаций: List/Get/Create
- [x] companysettings — настройки компании: Context().CompanySettings/UpdateCompanySettings; метаданные — metadata.GetCompanySettings
- [x] companysubscription — подписка компании: Context().Subscription (только чтение)
- [x] consignment — партия
- [x] contentcard — карточка контента; изображения — Entity.Images
- [x] contract — договор
- [x] counterparty — контрагент
- [x] country — страна
- [x] currency — валюта
- [x] customentity — пользовательский справочник: Create/Update/Delete, элементы List/Get/Create/Update/Delete
- [x] discount — скидки: общий список, изменение и массовое изменение округления копеек
- [x] employee — сотрудник: Get/List/Create/Update, права GetSecurity/UpdateSecurity
- [x] context/employee — контекст сотрудника
- [x] context/application — контекст решения (токен решения)
- [ ] eventfeed — лента событий (notes)
- [x] expenseitem — статья расходов
- [x] files — файлы: List/Upload/Download
- [x] group — отдел (массовых операций нет)
- [ ] gtd — ГТД (объект без endpoint)
- [x] images — изображения: List/Add/Replace/Delete/DeleteMany/Download (product, bundle, variant, contentcard)
- [ ] labels — этикетки и ценники
- [ ] namedfilter — сохранённые фильтры
- [x] organization — юрлицо
- [x] organizationbranch — подразделение юрлица
- [x] personaldiscount — персональная скидка: Get/Create/Update/Delete (список — Discount().List)
- [x] pricetype — типы цен (context/companysettings): List/Get/GetDefault/ReplaceAll
- [ ] processingplan — техкарта
- [ ] processingplanfolder — группа техкарт
- [ ] processingprocess — техпроцесс
- [ ] processingstage — этап производства
- [x] product — товар
- [x] productfolder — группа товаров
- [x] project — проект
- [x] region — регион: Get/List (только чтение)
- [x] retailstore — точка продаж
- [x] role — пользовательские роли и предустановленные (Preset)
- [x] saleplatform — площадка для продаж: Get/List (только чтение)
- [x] saleschannel — канал продаж
- [x] service — услуга
- [x] specialpricediscount — скидка на спеццену: Get/Create/Update/Delete (список — Discount().List)
- [x] states — статусы: чтение через metadata, Create/Update
- [x] store — склад: Get/List/Create/Update/Delete и массовые; зоны и ячейки (List/Get/Create/Update/Delete, массовые)
- [ ] task — задача
- [x] taxrate — ставка НДС
- [ ] template — шаблоны печатных форм
- [x] thing — серийный номер: Get/List (только чтение)
- [x] trackingcode — коды маркировки позиции документа: List/Save/DeleteMany
- [x] uom — единица измерения
- [x] usersettings — настройки пользователя: Context().UserSettings/UpdateUserSettings
- [x] variant — модификация
- [ ] webhook — вебхук
- [ ] webhookstock — вебхук остатков

## Документы

- [x] cashin — приходный ордер
- [x] cashout — расходный ордер
- [x] commissionreportin — полученный отчёт комиссионера
- [x] commissionreportout — выданный отчёт комиссионера
- [x] counterpartyadjustment — корректировка взаиморасчётов
- [x] customerorder — заказ покупателя
- [x] demand — отгрузка
- [ ] emissionorder — заказ кодов маркировки
- [x] enter — оприходование
- [x] facturein — счёт-фактура полученный
- [x] factureout — счёт-фактура выданный
- [x] internalorder — внутренний заказ
- [x] inventory — инвентаризация
- [x] invoicein — счёт поставщика
- [x] invoiceout — счёт покупателю
- [x] loss — списание
- [x] move — перемещение
- [x] paymentin — входящий платёж
- [x] paymentout — исходящий платёж
- [x] payroll — начисление зарплаты
- [x] prepayment — предоплата
- [x] prepaymentreturn — возврат предоплаты
- [x] pricelist — прайс-лист
- [ ] processing — техоперация
- [ ] processingorder — заказ на производство
- [ ] productionstagecompletion — выполнение этапа производства
- [ ] productiontask — производственное задание
- [x] purchaseorder — заказ поставщику
- [x] purchasereturn — возврат поставщику
- [x] retaildemand — розничная продажа
- [x] retaildrawercashin — внесение денег
- [x] retaildrawercashout — выплата денег
- [x] retailsalesreturn — розничный возврат
- [x] retailshift — розничная смена
- [ ] retireorder — вывод из оборота
- [x] salesreturn — возврат покупателя
- [x] supply — приёмка

## Отчёты

- [ ] dashboard — day/week/month
- [ ] byoperations — stock/reserve/intransit
- [ ] counterparty — показатели контрагентов
- [ ] money — plotseries/byaccount
- [ ] profit — byproduct/byvariant/byEmployee/bycounterparty/bysaleschannel
- [ ] orders — plotseries
- [ ] sales — plotseries
- [ ] stock — all, all/current, bystore, bystore/current, byoperation, byslot/current (пакет-заглушка)
- [x] turnover — all/bystore/byoperations

## Прочие ресурсы

- [x] audit — аудит: контексты, события, фильтры
- [ ] notification — уведомления (лента, markasread)
- [ ] notification/settings — настройки уведомлений
- [x] async — асинхронные задачи: Get, List, Wait, Result, Cancel; создание — ListAsync у counterparty и assortment
- [x] security/token — получение токена: CreateToken (Basic на вызов, без повторов)
- [x] download — скачивание файлов (transport.Client.Download)
- [x] entity/metadata — общие метаданные, метаданные типа, описания доп. полей
- [x] wizard — автозаполнение документов (POST wizard/{тип}; demand, salesreturn и другие)
