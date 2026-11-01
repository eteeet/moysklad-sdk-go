# Статус сущностей JSON API 1.2

Источник: api-remap-1.2-doc, revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`.
[x] — есть пакет с Get/Create/Update (документы с позициями — ещё GetPositions).
List/Delete/expand/batch — см. [status-technical.md](status-technical.md).

Итого: справочники 6/59, документы 11/37, отчёты 0/9, прочее 0/8.

## Справочники

- [ ] accumulationdiscount — накопительная скидка
- [ ] aggregatepack — транспортная упаковка
- [ ] application — серверные решения
- [ ] assortment — ассортимент
- [ ] bonusprogram — бонусная программа
- [ ] bonustransaction — бонусная операция
- [ ] bundle — комплект
- [ ] cashier — кассир (retailstore/{id}/cashiers)
- [ ] characteristic — характеристики модификаций
- [ ] companysettings — настройки компании (context)
- [ ] companysubscription — подписка компании (accountSettings)
- [ ] consignment — партия
- [ ] contentcard — карточка контента
- [ ] contract — договор
- [x] counterparty — контрагент
- [x] country — страна
- [ ] currency — валюта
- [ ] customentity — пользовательский справочник
- [ ] discount — скидки (общий список)
- [ ] employee — сотрудник
- [ ] context/employee — контекст сотрудника
- [ ] eventfeed — лента событий (notes)
- [x] expenseitem — статья расходов
- [ ] files — файлы
- [ ] group — отдел
- [ ] gtd — ГТД (объект без endpoint)
- [ ] images — изображения
- [ ] labels — этикетки и ценники
- [ ] namedfilter — сохранённые фильтры
- [x] organization — юрлицо
- [ ] organizationbranch — подразделение юрлица
- [ ] personaldiscount — персональная скидка
- [ ] pricetype — типы цен (context/companysettings)
- [ ] processingplan — техкарта
- [ ] processingplanfolder — группа техкарт
- [ ] processingprocess — техпроцесс
- [ ] processingstage — этап производства
- [x] product — товар
- [ ] productfolder — группа товаров
- [ ] project — проект
- [ ] region — регион
- [ ] retailstore — точка продаж
- [ ] role — пользовательские роли
- [ ] saleplatform — площадка для продаж
- [ ] saleschannel — канал продаж
- [ ] service — услуга
- [ ] specialpricediscount — скидка на спеццену
- [ ] states — статусы (metadata/states)
- [ ] store — склад (+ slots, zones)
- [ ] task — задача
- [ ] taxrate — ставка НДС
- [ ] template — шаблоны печатных форм
- [ ] thing — серийный номер
- [ ] trackingcode — коды маркировки
- [x] uom — единица измерения
- [ ] usersettings — настройки пользователя (context)
- [ ] variant — модификация
- [ ] webhook — вебхук
- [ ] webhookstock — вебхук остатков

## Документы

- [x] cashin — приходный ордер
- [x] cashout — расходный ордер
- [ ] commissionreportin — полученный отчёт комиссионера
- [ ] commissionreportout — выданный отчёт комиссионера
- [ ] counterpartyadjustment — корректировка взаиморасчётов
- [ ] customerorder — заказ покупателя
- [x] demand — отгрузка
- [ ] emissionorder — заказ кодов маркировки
- [ ] enter — оприходование
- [ ] facturein — счёт-фактура полученный
- [ ] factureout — счёт-фактура выданный
- [ ] internalorder — внутренний заказ
- [ ] inventory — инвентаризация
- [x] invoicein — счёт поставщика
- [x] invoiceout — счёт покупателю
- [ ] loss — списание
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
- [ ] retaildemand — розничная продажа
- [ ] retaildrawercashin — внесение денег
- [ ] retaildrawercashout — выплата денег
- [ ] retailsalesreturn — розничный возврат
- [ ] retailshift — розничная смена
- [ ] retireorder — вывод из оборота
- [ ] salesreturn — возврат покупателя
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
- [ ] turnover — all/bystore/byoperations

## Прочие ресурсы

- [ ] audit — аудит (+ events, filters)
- [ ] notification — уведомления (лента, markasread)
- [ ] notification/settings — настройки уведомлений
- [ ] async — асинхронные задачи
- [ ] security/token — получение токена
- [ ] download — скачивание файлов
- [ ] entity/metadata — общие метаданные
- [ ] wizard — автозаполнение (demand, salesreturn)
