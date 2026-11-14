# Розница и бонусы

Источник контрактов — `api-remap-1.2-doc` revision `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`.
Все методы принимают `opts ...transport.RequestOption` ([ADR 0002](adr/0002-baseline.md)).
Розничные документы и бонусные операции не повторяются автоматически; суммы — копейки, баллы,
курсы и проценты — целые числа или проценты, не деньги.

| Пакет / метод | Маршрут | Обязательно при Create |
| --- | --- | --- |
| `entity/retailstore`: `Get`, `List`, `Create`, `Update` | `entity/retailstore[/{id}]` | name, organization, store, priceType |
| `entity/retailshift`: `Get`, `List`, `Create`, `Update` | `entity/retailshift[/{id}]` | organization, retailStore |
| `entity/retaildemand`: `Get`, `List`, `Create`, `Update`, `GetPositions` | `entity/retaildemand[/{id}[/positions]]` | retailShift |
| `entity/retailsalesreturn`: `Get`, `List`, `Create`, `Update`, `GetPositions` | `entity/retailsalesreturn[/{id}[/positions]]` | name, organization, store, retailStore, retailShift, agent |
| `entity/bonusprogram`: `Get`, `List`, `Create`, `Update` | `entity/bonusprogram[/{id}]` | name, active, allProducts, allAgents, earnRateRoublesToPoint, spendRatePointsToRouble, maxPaidRatePercents |
| `entity/bonustransaction`: `Get`, `List`, `Create`, `Update` | `entity/bonustransaction[/{id}]` | agent, bonusProgram, transactionType |

## Точка продаж и смена

- RetailStore: environment и state — readonly JSON как пришёл (`jsontext.Value`); masterRetailStores и
  productFolders читаются коллекцией или массивом (`Links`), пишутся массивом ссылок. Cashiers — только
  метаданные; ресурс кассиров не входит. Конфликт allowCreateProducts/controlShippingStock, тариф и
  валюту счёта организации проверяет сервер. discountMaxPercent — float64 (таблица: Int, ответы: 100.0).
- RetailShift: позиций нет. Открытие — Create, закрытие и фискальные данные — явные `CloseDate` и
  `Cheque` (`ChequeRequest` на Field) в Update; SDK смену сам не создаёт. Платежи — `PaymentOperations`
  с linkedSum.

## Продажа и возврат

- RetailDemand: метод требует только retailShift (таблица ещё помечает agent и store — противоречие).
  Смена активна и moment позже её moment — проверяет сервер. Поля оплат (cash, noCash, QR, предоплаты,
  аванс, подарочные сертификаты) пишутся как есть: распределение остатка и запреты смешанной оплаты —
  сервер. Фискальный чек, organization, sum и payedSum только читаются. SDK не фискализирует и не
  начисляет бонусы. GetPositions — `fields=stock` или `declaration`.
- RetailSalesReturn: demand — розничная продажа, только для возврата по основанию; без основания
  cashSum и noCashSum обязательны (явный 0 уходит). Update не содержит agent, agentAccount и demand.
  GetPositions — `fields=stock` (у возврата без основания без cost).
- Строки — общая форма `positions.Position` ([positions](positions.md)): cost, pack, things,
  trackingCodes; trackingCodes_1162 и declaration только читаются.

## Бонусы

- BonusProgram: обязательные поля Create — обычные значения, false и 0 всегда отправляются. Диапазоны и
  правила приветственных баллов проверяет сервер; баланс SDK не считает.
- BonusTransaction: тип операции (EARNING/SPENDING) в Update менять нельзя; transactionStatus
  (WAIT_PROCESSING, COMPLETED, CANCELED) и categoryType только читаются. executionDate источник
  описывает для начислений: откладывает обработку (будущая дата — тарифная опция).
- Права: источник не описывает, каким токенам разрешена запись бонусных операций; допускает ли сервер
  токен решения (Vendor), выясняется его ответом. Отказ (например, для токена из `transport.WithAuth`)
  возвращается как APIError, SDK не пробует другой токен.

## Не входит

Шаблоны /new, expand, отдельные операции с позициями, внесения и
выплаты денег (retaildrawercashin/out), фискализация и расчёт бонусов. Кассиры точки продаж — [dictionaries.md](dictionaries.md). Delete и массовые операции — [write-async.md](write-async.md).
