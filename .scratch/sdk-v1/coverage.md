Status: resolved
Stage: APPROVED

# Покрытие исходного запроса

[Карта](map.md), [общая spec](spec.md). Ни одна строка не считается реализованной
по одному существованию spec/ticket. Приёмка R04 — в конце файла.

| Требование | Tickets | Граница |
| --- | --- | --- |
| audit | [A07](../sdk-v1-access/issues/07-audit-contexts.md), [A08](../sdk-v1-access/issues/08-audit-events.md) | контексты/filters и события/diff; без guarantees durable cursor |
| bundle variant service | [C04](../sdk-v1-catalog/issues/04-service.md), [C05](../sdk-v1-catalog/issues/05-characteristics.md), [C06](../sdk-v1-catalog/issues/06-bundle.md), [C07](../sdk-v1-catalog/issues/07-variant.md) | characteristics — необходимое дополнение к Variant |
| employee role | [A01](../sdk-v1-access/issues/01-employee.md), [A02](../sdk-v1-access/issues/02-roles-security.md) | employee security и presets явно входят; activation/password не входят |
| context/employee | [A03](../sdk-v1-access/issues/03-employee-context.md) | не путать с ResponseContext |
| bonustransaction bonusprogram | [E05](../sdk-v1-retail/issues/05-bonusprogram.md), [E06](../sdk-v1-retail/issues/06-bonustransaction.md) | permissions App/Employee и состояния операции |
| assortment | [C08](../sdk-v1-catalog/issues/08-assortment.md) | read-only mixed List/search; consignment common-read, unknown raw fallback |
| currency | [C01](../sdk-v1-catalog/issues/01-currency.md) | курс и system/custom особенности |
| files | [F01](../sdk-v1-files-reports/issues/01-download-transport.md), [F02](../sdk-v1-files-reports/issues/02-files.md) | upload/list/download; images и Delete отдельны |
| pricetype | [C02](../sdk-v1-catalog/issues/02-pricetype.md) | plain array и полный ReplaceAll, без вымышленного PUT |
| retailstore | [E01](../sdk-v1-retail/issues/01-retailstore.md) | отдельный resource, не settings скрытого клиента |
| states | [A05](../sdk-v1-access/issues/05-entity-metadata.md), [A06](../sdk-v1-access/issues/06-states.md) | metadata resource + Create/Update, admin check сервером |
| taxrate | [C03](../sdk-v1-catalog/issues/03-taxrate.md) | actual resource contract |
| Общая форма новых документных positions | [D01](../sdk-v1-documents/issues/01-position-contracts.md) | shared поля и readonly/write matrix; без generic расширения по аналогии |
| customerorder | [D02](../sdk-v1-documents/issues/02-customerorder.md) | включая typed positions |
| enter loss | [D03](../sdk-v1-documents/issues/03-enter.md), [D04](../sdk-v1-documents/issues/04-loss.md) | без выдуманного agent |
| retaildemand retailsalesreturn retailshift salesreturn | [D05](../sdk-v1-documents/issues/05-salesreturn.md), [E02](../sdk-v1-retail/issues/02-retailshift.md), [E03](../sdk-v1-retail/issues/03-retaildemand.md), [E04](../sdk-v1-retail/issues/04-retailsalesreturn.md) | own Create/update constraints based/unbased |
| turnover all/bystore/byoperations | [F03](../sdk-v1-files-reports/issues/03-turnover-all.md), [F04](../sdk-v1-files-reports/issues/04-turnover-bystore.md), [F05](../sdk-v1-files-reports/issues/05-turnover-byoperations.md) | три разных rows/options, source date contradiction явна |
| security/token; Получение токена POST security/token | [T04](../sdk-v1-foundation/issues/04-security-token.md) | одно требование; no auto-retry, старые токены отзываются |
| entity/metadata | [A05](../sdk-v1-access/issues/05-entity-metadata.md) | global catalogue/per-entity/attribute definitions read |
| Коды ошибок | [T03](../sdk-v1-foundation/issues/03-error-codes.md) | полный pinned numeric catalogue, unknown preserved |
| syncId | [T13](../sdk-v1-foundation/issues/13-product-syncid.md), [C04](../sdk-v1-catalog/issues/04-service.md), [C06](../sdk-v1-catalog/issues/06-bundle.md), [C07](../sdk-v1-catalog/issues/07-variant.md), [D02](../sdk-v1-documents/issues/02-customerorder.md), [D03](../sdk-v1-documents/issues/03-enter.md), [D04](../sdk-v1-documents/issues/04-loss.md), [D05](../sdk-v1-documents/issues/05-salesreturn.md), [E01](../sdk-v1-retail/issues/01-retailstore.md), [E02](../sdk-v1-retail/issues/02-retailshift.md), [E03](../sdk-v1-retail/issues/03-retaildemand.md), [E04](../sdk-v1-retail/issues/04-retailsalesreturn.md), [E05](../sdk-v1-retail/issues/05-bonusprogram.md), [E06](../sdk-v1-retail/issues/06-bonustransaction.md) | только где источник поддерживает; существующие документы не теряют поля |
| Заголовки — все | [T01](../sdk-v1-foundation/issues/01-public-interface.md), [T02](../sdk-v1-foundation/issues/02-headers.md) | полная inventory research: raw+typed, request/response/repeated/units |
| Лимиты — все применимые | [T08](../sdk-v1-foundation/issues/08-rate-concurrency.md), [T09](../sdk-v1-foundation/issues/09-safe-retry.md), [T10](../sdk-v1-foundation/issues/10-size-limits.md), [F02](../sdk-v1-files-reports/issues/02-files.md) | политики rate/concurrency/request-size/files/1000; X02/X03 для batch/async |
| Query-параметры все кроме expand | [T05](../sdk-v1-foundation/issues/05-filter-order.md), [T06](../sdk-v1-foundation/issues/06-query-pages.md), [T07](../sdk-v1-foundation/issues/07-iterator.md) | full grammar/names/raw; endpoint-specific typed capabilities, без обещания иных services |
| list | [T11](../sdk-v1-foundation/issues/11-list-dictionaries.md), [T12](../sdk-v1-foundation/issues/12-list-documents.md), [A01](../sdk-v1-access/issues/01-employee.md), [A02](../sdk-v1-access/issues/02-roles-security.md), [C01](../sdk-v1-catalog/issues/01-currency.md), [C02](../sdk-v1-catalog/issues/02-pricetype.md), [C03](../sdk-v1-catalog/issues/03-taxrate.md), [C04](../sdk-v1-catalog/issues/04-service.md), [C06](../sdk-v1-catalog/issues/06-bundle.md), [C07](../sdk-v1-catalog/issues/07-variant.md), [C08](../sdk-v1-catalog/issues/08-assortment.md), [D02](../sdk-v1-documents/issues/02-customerorder.md), [D03](../sdk-v1-documents/issues/03-enter.md), [D04](../sdk-v1-documents/issues/04-loss.md), [D05](../sdk-v1-documents/issues/05-salesreturn.md), [E01](../sdk-v1-retail/issues/01-retailstore.md), [E02](../sdk-v1-retail/issues/02-retailshift.md), [E03](../sdk-v1-retail/issues/03-retaildemand.md), [E04](../sdk-v1-retail/issues/04-retailsalesreturn.md), [E05](../sdk-v1-retail/issues/05-bonusprogram.md), [E06](../sdk-v1-retail/issues/06-bonustransaction.md), [F02](../sdk-v1-files-reports/issues/02-files.md) | существующие 17 + новые; array-only resources не получают вымышленную pagination |
| context/application | [A04](../sdk-v1-access/issues/04-application-context.md) | не Vendor API и не entity/application по installationId |
| Фасад SDK | [T01](../sdk-v1-foundation/issues/01-public-interface.md), [R01](../sdk-v1-release/issues/01-facade.md) | shared transport и typed handles без HTTP в constructor |
| Готовность первого выпуска | [R02](../sdk-v1-release/issues/02-examples.md), [R03](../sdk-v1-release/issues/03-ci-release.md), [R04](../sdk-v1-release/issues/04-acceptance.md) | examples/CI/CHANGELOG и method-specific status вместо единого CRUD чекбокса |
| Опциональные Delete/batch/async | [X01](../sdk-v1-optional/issues/01-delete.md), [X02](../sdk-v1-optional/issues/02-batch.md), [X03](../sdk-v1-optional/issues/03-async.md) | решение пользователя; needs-info; пока вне основной v1 |

Прежние первые пять «открытых замечаний последнего ревью» status-technical.md
уже закрыты кодом/14 коммитами dbe1e01…acbe294. Они не создаются повторно.
Пять dirty core worktree закрыты отдельным [cleanup ticket](../review-followups/issues/01-followups.md)
(resolved: stash -u, branches сохранены); git cherry показывал patch-equivalent commits.

## Приёмка v1 (R04, 2026-11-02)

Все строки, кроме «Опциональные Delete/batch/async», закрыты: каждый перечисленный ticket — resolved,
в его Answer записаны BASE, коммит, тесты и независимые Standards/Spec review; волны прошли
пользовательские review в `.scratch/sdk-v1-*/review-*.md`. X01–X03 вынесены из v1 решением
пользователя (2026-11-01) и остаются needs-info.

Проверяемые доказательства в коде:
- `client_test.go`: каждый пакет `entity/*`, metadata, files, contexts, audit, security и
  report/turnover доступен через фасад и использует один transport;
- `status_test.go`: каждая отметка [x] в docs/status-entities.md ведёт к handle фасада, каждый
  handle отмечен, «Итого» совпадает с отметками (исправлены security/token и context/application);
- `fixtures_test.go`: upstream и скопированные из источника фикстуры совпадают с SHA-256;
- `make check`: gofmt, vet, race-тесты, примеры, пробелы всего дерева.

Нерешённые ограничения (не выдаются за проверенное):
- живой API не использовался ни в одном ticket; поведение сервера подтверждено только примерами
  документации; опциональный read-only smoke — по отдельному согласию;
- противоречия источника: Turnover.ByOperations (даты обязательны в таблице, умолчания в тексте),
  Role.permissions (таблица и примеры), stockByStore и operation.moment — записаны в Answer
  соответствующих tickets;
- `report/stock` — пустой пакет-заглушка вне v1;
- единицы X-RateLimit-Limit/Remaining источник не уточняет; Limiter считает бюджет и параллельность
  только в пределах процесса;
- `fields` в List справочников и документов не принимается; `fields=minimumStock` у Product/Variant
  ждёт модели неснижаемого остатка; изображения в files не поддержаны;
- expand, Delete, batch, async, остальные сущности и отчёты, Vendor API — вне v1
  (docs/status-entities.md, docs/status-technical.md).
