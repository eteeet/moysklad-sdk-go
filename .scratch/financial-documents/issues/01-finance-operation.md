Status: ready-for-agent
Stage: APPROVED
Blocked by: None

# 01: FinanceOperation — самостоятельный read/write JSON-контракт

**What to build:** разработчик решения читает ссылку на оплаченную операцию с linkedSum
и независимо формирует элемент operations для финансового документа: только meta либо
meta с явно заданной суммой. Результат проверяется без реализации какой-либо сущности.

Публичные типы, числовое представление, Field/null-policy и JSON seams утверждены
пользователем для первой волны. Разрешены реализация и локальный commit, без push.

## Контракт

- Нормативные документы: [spec](../spec.md), [contract](../contract.md),
  [исследование](../research.md). Они содержат предложения, не решения из недоступного чата.
- Proposal: model.FinanceOperation с Meta model.Meta и LinkedSum float64;
  model.FinanceOperationRequest с Meta model.Meta и LinkedSum model.Field[float64] с omitzero.
  Meta находится непосредственно в элементе operations, без вложенного operation wrapper.
- Read `linkedSum` = оплаченная данным платежом сумма операции; optional/missing даёт
  нулевой float64. Write linkedSum можно опустить; явно заданный 0 должен присутствовать.
  Значения передаются без округления, пересчёта и конвертации. Точное десятичное вычисление
  и сохранение исходного текстового вида числа не входят в контракт.
- MetaType открытый; неизвестный тип читается и сериализуется без локальной валидации.
  Операция содержит только Meta и LinkedSum; id/name/positions/payments и прочие
  expanded-поля игнорируются. Никакой поддержки expand и получения оснований.
- Field различает unset/null/value. Null кодируется точно, но допустимость сервером
  null для linkedSum не подтверждена и не должна обещаться в документации.
- Тип reusable для PaymentIn/CashIn и PaymentOut/CashOut. Список допустимых операций
  описывает серверные ограничения, а не закрытый union SDK:
  incoming — customerorder, purchasereturn, demand, invoiceout, commissionreportin,
  retailshift; outgoing — salesreturn, supply, invoicein, purchaseorder,
  commissionreportout, payroll. Никаких импортов этих entity packages.
- Здесь нет работы с DocumentRate, поэтому внешний ticket document-common не блокирует этот slice.

## Acceptance criteria

- [x] Имена, read/write shape, float64 и null-policy утверждены пользователем; ticket переведён в ready-for-agent.
- [ ] Публичные типы позволяют прочитать/записать описанный JSON без entity Service и transport.
- [ ] Отдельные JSON-тесты выполняются encoding/json и encoding/json/v2 по имеющемуся
  образцу codec tables. Проверены Meta.Href/MetadataHref/Type/MediaType/UUIDHref и linkedSum.
- [ ] Независимые golden cases: meta-only, linkedSum=0, дробное значение, unset Field,
  Null Field. Для read — отсутствующий linkedSum, 0/0.0, дробь, неизвестный meta.type,
  лишние expanded-поля без превращения их в поддержанный scope.
- [ ] Внешний Field[[]FinanceOperationRequest] проверен через публичную сериализацию:
  unset, null, [], nil slice как [], несколько операций и элементы с разными meta.type.
  Это JSON-проверка, не утверждение server clear/merge поведения.
- [ ] Пустой/нулевой Meta не вызывает новых бизнес-проверок модели; запрет неверных
  ссылок не вводится без отдельного утверждения.
- [ ] Документация общего типа объясняет linkedSum optional write, точность float64,
  отсутствие арифметики, запрет трактовать [] как доказанную очистку и границы expand.
- [ ] Существующие Reference/Meta/Field/Attribute/DateTime и их поведение не изменены.
- [ ] Форматирование, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check; независимые Standards/Spec reviews с фиксированным BASE.

## Out of Scope

Сущности финансовых документов, HTTP, DocumentRate/Currency, positions, банковские счета,
expanded references, закрытые типы всех оснований, Money/Decimal и финансовые расчёты,
живой API/.env. Реализация и локальный commit разрешены командой первой волны; push/live API/.env запрещены.

## Источники

- [Общий контракт привязки платежей](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_common_info.md#L476-L605).
- [PaymentIn.Update: необязательный linkedSum](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_in.md#L1491-L1499).
- [PaymentOut.Update](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_payment_out.md#L1639-L1648).

Master HEAD проверен 2026-11-01 и совпадает с закреплённым SHA.

## Утверждение и ownership первой волны

Пользователь 2026-11-01 утвердил первую волну после исправления денежных типов и смысла nil курса.
JSON и HTTP seams этого ticket согласованы. Локальная реализация через TDD и commit разрешены;
независимые Standards/Spec review координирует исполнитель интеграции относительно planning BASE.
Не менять coverage.md/README и общие карты: документацию интегратор сведёт отдельным коммитом.
Файлы и трёхпоточный запуск закреплены в [handoff](../../document-workstream/wave-1.md).
