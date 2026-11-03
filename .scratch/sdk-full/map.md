Status: needs-triage
Stage: DRAFT
Type: wayfinder:map

# SDK — полное покрытие JSON API 1.2: карта захваток и тикетов

## Destination

Закрыть всё нереализованное из [status-backlog](../../docs/status-backlog.md), кроме expand.
Источник: api-remap-1.2-doc revision `0e1b3d61fb8d8fa314469361d96560d6c3f7f57b` (от `9a811cef` изменилась только разметка
`_labels.md` и `_print.md`). Нормативные ограничения и тестовый контракт — [sdk-v1](../sdk-v1/spec.md).

## Захватки

| № | Захватка | Tickets |
| --- | --- | --- |
| 1 | [Справочники, настройки и метаданные](../sdk-full-dictionaries/spec.md) | 27 |
| 2 | [Товары, скидки и маркировка](../sdk-full-catalog/spec.md) | 13 |
| 3 | [Документы и операции документов](../sdk-full-documents/spec.md) | 20 |
| 4 | [Производство](../sdk-full-production/spec.md) | 8 |
| 5 | [Отчёты, печать и интеграции](../sdk-full-reports/spec.md) | 23 |

Всего 91 tickets. Expand — отдельно, последним и только по команде.

## Порядок и параллельность

1. Захватки 1, 2 (кроме G12/G13), 5 и общие механизмы O01–O06 независимы друг от друга.
2. O01 (мутации позиций), O02 (шаблон new) и O04 (Trash) раньше новых документов:
   O08–O20, G12/G13, M05–M08.
3. Производство: M01/M02 → M03 → M04 → M05–M07 → M08.
4. Печать: I11 → I12 → I13; I14 после I11. Async отчётов I23 — после синхронных отчётов.
5. S10 после G01 (один владелец формы штрихкодов); S26 после S16 (один пакет organization).
6. Общие файлы (`positions`, `transport`, фасад) имеют одного владельца в волне.

## Decisions-so-far

- gtd и файлы сущностей уже покрыты (`positions.GTD`, пакет `files`) — без tickets.
- O07 needs-info: расхождение источников по счетам кассовых ордеров.
- Захватка 1 реализована 2026-11-03 (BASE 30e5d76 → dd628ca); S10 — после G01 (d89f11a).
- Захватка 3 реализована 2026-11-03 (BASE d21d84d), кроме O07 (needs-info); G12/G13 разблокированы.
- Захватка 5 и хвосты G12/G13 реализованы 2026-11-03 (BASE 3778088 → aa5c8ea); O07 — needs-info; expand не делался.
- Захватка 4 реализована 2026-11-03 (BASE c3bace0 → 3778088); решения пользователя — в review-wave.
- Захватка 2 реализована 2026-11-03 (BASE f2f5f09 → 670f379), кроме G12/G13 (ждут O01, O04).
- [Review 30e5d76](review-planning-30e5d76.md): все 14 замечаний приняты и внесены в spec/tickets;
  добавлены S26, S27, I23.

## Все tickets

### Справочники, настройки и метаданные

| ID | Ticket | Blocked by | Status |
| --- | --- | --- | --- |
| S01 | [Store Get/List/Create/Update + zones/slots](../sdk-full-dictionaries/issues/01-store.md) | нет | resolved |
| S02 | [ProductFolder Get/List/Create/Update](../sdk-full-dictionaries/issues/02-productfolder.md) | нет | resolved |
| S03 | [Group Get/List/Create/Update](../sdk-full-dictionaries/issues/03-group.md) | нет | resolved |
| S04 | [Project Get/List/Create/Update](../sdk-full-dictionaries/issues/04-project.md) | нет | resolved |
| S05 | [Contract Get/List/Create/Update](../sdk-full-dictionaries/issues/05-contract.md) | нет | resolved |
| S06 | [Region Get/List (read-only)](../sdk-full-dictionaries/issues/06-region.md) | нет | resolved |
| S07 | [SalesChannel Get/List/Create/Update](../sdk-full-dictionaries/issues/07-saleschannel.md) | нет | resolved |
| S08 | [SalePlatform Get/List (read-only)](../sdk-full-dictionaries/issues/08-saleplatform.md) | нет | resolved |
| S09 | [OrganizationBranch Get/List/Create/Update](../sdk-full-dictionaries/issues/09-organizationbranch.md) | нет | resolved |
| S10 | [Consignment Get/List/Create/Update](../sdk-full-dictionaries/issues/10-consignment.md) | G01 | resolved |
| S11 | [Thing Get/List (read-only)](../sdk-full-dictionaries/issues/11-thing.md) | нет | resolved |
| S12 | [Counterparty accounts](../sdk-full-dictionaries/issues/12-counterparty-accounts.md) | нет | resolved |
| S13 | [Counterparty contactpersons](../sdk-full-dictionaries/issues/13-counterparty-contactpersons.md) | нет | resolved |
| S14 | [Counterparty notes](../sdk-full-dictionaries/issues/14-counterparty-notes.md) | нет | resolved |
| S15 | [Counterparty: discounts и региональные реквизиты](../sdk-full-dictionaries/issues/15-counterparty-fields.md) | нет | resolved |
| S16 | [Organization accounts, fsrarId, utmUrl](../sdk-full-dictionaries/issues/16-organization-accounts.md) | S12 | resolved |
| S17 | [Настройки справочников: counterparty и assortment](../sdk-full-dictionaries/issues/17-dictionary-settings.md) | нет | resolved |
| S18 | [CompanySettings GET/PUT](../sdk-full-dictionaries/issues/18-companysettings.md) | нет | resolved |
| S19 | [Bundle: мутации компонентов](../sdk-full-dictionaries/issues/19-bundle-components.md) | нет | resolved |
| S20 | [Employee: активация, деактивация, сброс пароля](../sdk-full-dictionaries/issues/20-employee-access.md) | нет | resolved |
| S21 | [RetailStore cashiers (read-only)](../sdk-full-dictionaries/issues/21-cashier.md) | нет | resolved |
| S22 | [UserSettings GET/PUT](../sdk-full-dictionaries/issues/22-usersettings.md) | нет | resolved |
| S23 | [CompanySubscription GET](../sdk-full-dictionaries/issues/23-companysubscription.md) | нет | resolved |
| S24 | [Доп. поля: CRUD описаний](../sdk-full-dictionaries/issues/24-attributes-crud.md) | нет | resolved |
| S25 | [CustomEntity и элементы](../sdk-full-dictionaries/issues/25-customentity.md) | S24 | resolved |
| S26 | [Organization: НДС, подписи, печать, региональные реквизиты](../sdk-full-dictionaries/issues/26-organization-fields.md) | S16 | resolved |
| S27 | [Массовые операции states и characteristics](../sdk-full-dictionaries/issues/27-metadata-batch.md) | нет | resolved |

### Товары, скидки и маркировка

| ID | Ticket | Blocked by | Status |
| --- | --- | --- | --- |
| G01 | [Product: штрихкоды и упаковки](../sdk-full-catalog/issues/01-barcodes-packs.md) | нет | resolved |
| G02 | [Product: alcohol, отслеживание и прочие поля](../sdk-full-catalog/issues/02-product-attributes.md) | G01 | resolved |
| G03 | [Неснижаемый остаток, storebalances и fields=minimumStock](../sdk-full-catalog/issues/03-minimumstock.md) | нет | resolved |
| G04 | [Images: product/variant/bundle](../sdk-full-catalog/issues/04-images.md) | нет | resolved |
| G05 | [ContentCard Get/List/Create/Update + images](../sdk-full-catalog/issues/05-contentcard.md) | G04 | resolved |
| G06 | [Discount: общий список и округление копеек](../sdk-full-catalog/issues/06-discount.md) | нет | resolved |
| G07 | [AccumulationDiscount CRUD](../sdk-full-catalog/issues/07-accumulationdiscount.md) | G06 | resolved |
| G08 | [PersonalDiscount CRUD](../sdk-full-catalog/issues/08-personaldiscount.md) | G06 | resolved |
| G09 | [SpecialPriceDiscount CRUD](../sdk-full-catalog/issues/09-specialpricediscount.md) | G06 | resolved |
| G10 | [TrackingCode List + массовые](../sdk-full-catalog/issues/10-trackingcode.md) | нет | resolved |
| G11 | [AggregatePack + positions](../sdk-full-catalog/issues/11-aggregatepack.md) | нет | resolved |
| G12 | [EmissionOrder Get/List/Create/Update + positions](../sdk-full-catalog/issues/12-emissionorder.md) | O01, O04 | resolved |
| G13 | [RetireOrder Get/List/Create/Update + positions](../sdk-full-catalog/issues/13-retireorder.md) | O01, O04 | resolved |

### Документы и операции документов

| ID | Ticket | Blocked by | Status |
| --- | --- | --- | --- |
| O01 | [Отдельные операции позиций](../sdk-full-documents/issues/01-position-mutations.md) | нет | resolved |
| O02 | [Шаблон документа PUT entity/{type}/new](../sdk-full-documents/issues/02-document-template.md) | нет | resolved |
| O03 | [Автозаполнение wizard/demand и wizard/salesreturn](../sdk-full-documents/issues/03-wizard.md) | нет | resolved |
| O04 | [Удаление в корзину](../sdk-full-documents/issues/04-trash.md) | нет | resolved |
| O05 | [Запись связанных документов](../sdk-full-documents/issues/05-linked-documents.md) | нет | resolved |
| O06 | [Оставшиеся поля позиций и transport-поля demand](../sdk-full-documents/issues/06-position-fields.md) | нет | resolved |
| O07 | [agentAccount/organizationAccount у cashin/cashout](../sdk-full-documents/issues/07-cash-accounts.md) | нет | needs-info |
| O08 | [InternalOrder + positions](../sdk-full-documents/issues/08-internalorder.md) | O01, O02, O04 | resolved |
| O09 | [Inventory + positions](../sdk-full-documents/issues/09-inventory.md) | O01, O02, O04 | resolved |
| O10 | [FactureIn](../sdk-full-documents/issues/10-facturein.md) | O02, O04 | resolved |
| O11 | [FactureOut](../sdk-full-documents/issues/11-factureout.md) | O02, O04 | resolved |
| O12 | [CommissionReportIn + 2 вида позиций](../sdk-full-documents/issues/12-commissionreportin.md) | O01, O04 | resolved |
| O13 | [CommissionReportOut + positions](../sdk-full-documents/issues/13-commissionreportout.md) | O01, O04 | resolved |
| O14 | [CounterpartyAdjustment](../sdk-full-documents/issues/14-counterpartyadjustment.md) | O04 | resolved |
| O15 | [PriceList + positions](../sdk-full-documents/issues/15-pricelist.md) | O01, O04 | resolved |
| O16 | [Prepayment (read + Delete)](../sdk-full-documents/issues/16-prepayment.md) | O04 | resolved |
| O17 | [PrepaymentReturn (read + Delete)](../sdk-full-documents/issues/17-prepaymentreturn.md) | O04 | resolved |
| O18 | [RetailDrawerCashIn](../sdk-full-documents/issues/18-retaildrawercashin.md) | O02, O04 | resolved |
| O19 | [RetailDrawerCashOut](../sdk-full-documents/issues/19-retaildrawercashout.md) | O02, O04 | resolved |
| O20 | [Payroll + positions](../sdk-full-documents/issues/20-payroll.md) | O01, O04 | resolved |

### Производство

| ID | Ticket | Blocked by | Status |
| --- | --- | --- | --- |
| M01 | [ProcessingPlanFolder](../sdk-full-production/issues/01-processingplanfolder.md) | нет | resolved |
| M02 | [ProcessingStage](../sdk-full-production/issues/02-processingstage.md) | нет | resolved |
| M03 | [ProcessingProcess + positions](../sdk-full-production/issues/03-processingprocess.md) | M02 | resolved |
| M04 | [ProcessingPlan + materials/products/stages](../sdk-full-production/issues/04-processingplan.md) | M01, M03 | resolved |
| M05 | [Processing + materials/products](../sdk-full-production/issues/05-processing.md) | M04, O02, O04 | resolved |
| M06 | [ProcessingOrder + positions](../sdk-full-production/issues/06-processingorder.md) | M04, O01, O02, O04 | resolved |
| M07 | [ProductionTask + materials/products/stages](../sdk-full-production/issues/07-productiontask.md) | M04, O04 | resolved |
| M08 | [ProductionStageCompletion](../sdk-full-production/issues/08-productionstagecompletion.md) | M07, O04 | resolved |

### Отчёты, печать и интеграции

| ID | Ticket | Blocked by | Status |
| --- | --- | --- | --- |
| I01 | [Report stock/all (расширенный)](../sdk-full-reports/issues/01-stock-all.md) | нет | resolved |
| I02 | [Report stock/bystore](../sdk-full-reports/issues/02-stock-bystore.md) | нет | resolved |
| I03 | [Краткие отчёты stock/all/current и bystore/current](../sdk-full-reports/issues/03-stock-current.md) | нет | resolved |
| I04 | [stock/byoperation и byslot/current](../sdk-full-reports/issues/04-stock-operation-slot.md) | нет | resolved |
| I05 | [Report byoperations stock/reserve/intransit](../sdk-full-reports/issues/05-byoperations.md) | нет | resolved |
| I06 | [Report dashboard day/week/month](../sdk-full-reports/issues/06-dashboard.md) | нет | resolved |
| I07 | [Report counterparty](../sdk-full-reports/issues/07-report-counterparty.md) | нет | resolved |
| I08 | [Report money plotseries/byaccount](../sdk-full-reports/issues/08-money.md) | нет | resolved |
| I09 | [Report profit (5 разрезов)](../sdk-full-reports/issues/09-profit.md) | нет | resolved |
| I10 | [Report orders/sales plotseries](../sdk-full-reports/issues/10-orders-sales.md) | нет | resolved |
| I11 | [Шаблоны печатных форм (read-only)](../sdk-full-reports/issues/11-templates.md) | нет | resolved |
| I12 | [Печать документов export](../sdk-full-reports/issues/12-print.md) | I11 | resolved |
| I13 | [Печать этикеток и ценников](../sdk-full-reports/issues/13-labels.md) | I11, I12 | resolved |
| I14 | [Публикации документов](../sdk-full-reports/issues/14-publications.md) | I11 | resolved |
| I15 | [Webhook CRUD](../sdk-full-reports/issues/15-webhook.md) | нет | resolved |
| I16 | [WebhookStock CRUD](../sdk-full-reports/issues/16-webhookstock.md) | нет | resolved |
| I17 | [Notification: лента и прочтение](../sdk-full-reports/issues/17-notification.md) | нет | resolved |
| I18 | [Notification settings GET/PUT](../sdk-full-reports/issues/18-notification-settings.md) | нет | resolved |
| I19 | [Лента событий документов (notes)](../sdk-full-reports/issues/19-eventfeed.md) | нет | resolved |
| I20 | [Task + типы задач и комментарии](../sdk-full-reports/issues/20-task.md) | нет | resolved |
| I21 | [NamedFilter (read-only)](../sdk-full-reports/issues/21-namedfilter.md) | нет | resolved |
| I22 | [Application GET](../sdk-full-reports/issues/22-application.md) | нет | resolved |
| I23 | [Async создания отчётов](../sdk-full-reports/issues/23-reports-async.md) | I01, I02, I03, I04, I06, I07, I08, I09, I10 | resolved |
