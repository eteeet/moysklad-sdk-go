# Статус сущностей JSON API 1.2

Источник: api-remap-1.2-doc, revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`.
[x] — есть пакет с Get/Create/Update (документы с позициями — ещё GetPositions) и handle фасада
`moysklad.Client` (у download — `Transport().Download`); иной набор методов указан в строке. `status_test.go` проверяет, что отметки и
handles фасада совпадают; набор методов — тесты пакетов.
List/Delete/expand/batch — см. [status-technical.md](status-technical.md) и [write-async.md](write-async.md).

Итого: справочники 29/60, документы 18/37, отчёты 1/9, прочее 5/8.

## Справочники

- [ ] accumulationdiscount — накопительная скидка
- [ ] aggregatepack — транспортная упаковка
- [ ] application — серверные решения
- [x] assortment — ассортимент: List и All (assortmentWithoutStock на каждой странице) и ListLegacy
- [x] bonusprogram — бонусная программа
- [x] bonustransaction — бонусная операция
- [x] bundle — комплект: Get/List/Create/Update, GetComponents/GetComponent
- [ ] cashier — кассир (retailstore/{id}/cashiers)
- [x] characteristic — характеристики модификаций: List/Get/Create
- [x] companysettings — метаданные: пользовательские справочники (metadata.GetCompanySettings)
- [ ] companysubscription — подписка компании (accountSettings)
- [ ] consignment — партия
- [ ] contentcard — карточка контента
- [x] contract — договор
- [x] counterparty — контрагент
- [x] country — страна
- [x] currency — валюта
- [ ] customentity — пользовательский справочник
- [ ] discount — скидки (общий список)
- [x] employee — сотрудник: Get/List/Create/Update, права GetSecurity/UpdateSecurity
- [x] context/employee — контекст сотрудника
- [x] context/application — контекст решения (токен решения)
- [ ] eventfeed — лента событий (notes)
- [x] expenseitem — статья расходов
- [x] files — файлы: List/Upload/Download
- [x] group — отдел (массовых операций нет)
- [ ] gtd — ГТД (объект без endpoint)
- [ ] images — изображения
- [ ] labels — этикетки и ценники
- [ ] namedfilter — сохранённые фильтры
- [x] organization — юрлицо
- [ ] organizationbranch — подразделение юрлица
- [ ] personaldiscount — персональная скидка
- [x] pricetype — типы цен (context/companysettings): List/Get/GetDefault/ReplaceAll
- [ ] processingplan — техкарта
- [ ] processingplanfolder — группа техкарт
- [ ] processingprocess — техпроцесс
- [ ] processingstage — этап производства
- [x] product — товар
- [x] productfolder — группа товаров
- [x] project — проект
- [ ] region — регион
- [x] retailstore — точка продаж
- [x] role — пользовательские роли и предустановленные (Preset)
- [ ] saleplatform — площадка для продаж
- [ ] saleschannel — канал продаж
- [x] service — услуга
- [ ] specialpricediscount — скидка на спеццену
- [x] states — статусы: чтение через metadata, Create/Update
- [x] store — склад: Get/List/Create/Update, зоны и ячейки (List/Get/Create/Update/Delete, массовые)
- [ ] task — задача
- [x] taxrate — ставка НДС
- [ ] template — шаблоны печатных форм
- [ ] thing — серийный номер
- [ ] trackingcode — коды маркировки
- [x] uom — единица измерения
- [ ] usersettings — настройки пользователя (context)
- [x] variant — модификация
- [ ] webhook — вебхук
- [ ] webhookstock — вебхук остатков

## Документы

- [x] cashin — приходный ордер
- [x] cashout — расходный ордер
- [ ] commissionreportin — полученный отчёт комиссионера
- [ ] commissionreportout — выданный отчёт комиссионера
- [ ] counterpartyadjustment — корректировка взаиморасчётов
- [x] customerorder — заказ покупателя
- [x] demand — отгрузка
- [ ] emissionorder — заказ кодов маркировки
- [x] enter — оприходование
- [ ] facturein — счёт-фактура полученный
- [ ] factureout — счёт-фактура выданный
- [ ] internalorder — внутренний заказ
- [ ] inventory — инвентаризация
- [x] invoicein — счёт поставщика
- [x] invoiceout — счёт покупателю
- [x] loss — списание
- [x] move — перемещение
- [x] paymentin — входящий платёж
- [x] paymentout — исходящий платёж
- [ ] payroll — начисление зарплаты
- [ ] prepayment — предоплата
- [ ] prepaymentreturn — возврат предоплаты
- [ ] pricelist — прайс-лист
- [ ] processing — техоперация
- [ ] processingorder — заказ на производство
- [ ] productionstagecompletion — выполнение этапа производства
- [ ] productiontask — производственное задание
- [x] purchaseorder — заказ поставщику
- [x] purchasereturn — возврат поставщику
- [x] retaildemand — розничная продажа
- [ ] retaildrawercashin — внесение денег
- [ ] retaildrawercashout — выплата денег
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
- [ ] wizard — автозаполнение (demand, salesreturn)
