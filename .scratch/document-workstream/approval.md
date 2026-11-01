Status: resolved
Stage: COMPLETED WAVE 1

# Утверждённая первая волна общих контрактов документов

2026-11-01 пользователь поручил внести согласованные замечания, закоммитить все
накопленные изменения и запустить следующий этап — первую волну из пяти tickets ниже.
Этот файл фиксирует первую волну C01/C02/F01/P01/P02 и последующий пересмотр формы строк,
но не утверждает entity spec. P03 отменён/заменён общей формой.
Подробные таблицы и criteria остаются в связанных spec/tickets. [Handoff](wave-1.md).

## Общие контракты — рекомендованный вариант

| Контракт | Предложение |
| --- | --- |
| Валюта | model.DocumentRate / DocumentRateRequest. Чтение: Currency *Reference, Value *float64; nil Value в присутствующем rate означает курс 1, но 1 в JSON не подставляется. Запись через Field; unset/null/{}/value различаются |
| Накладные расходы | model.DocumentOverhead / DocumentOverheadRequest. Sum float64, открытый OverheadDistribution с weight/volume/price; запись через Field, без распределения сумм в SDK |
| Платёжные связи | model.FinanceOperation / FinanceOperationRequest: плоские meta и linkedSum. Read LinkedSum float64, write Field[float64]; неизвестный Meta.Type не запрещается |
| Форма позиций | positions.Position/PositionRequest для всех товарных документов; Core/Item — общие части. Price/Discount/VAT/VATEnabled, read-only Overhead и sourceSlot/targetSlot. Незаданные Field не отправляются, допустимость заданных проверяет сервер |
| Чтение страниц | Один generic reader positions.Service[T any].GetPage, открытый documentType, существующий model.Collection[T]. Без обязательного interface у T и без allowlist документов |
| Числа | Все денежные поля → float64 в копейках независимо от API Int/Float, ставка vat → int. Без Money/Decimal, масштабирования и вычислений SDK |

Предлагаемый interface reader:

```go
NewService[T any](client *transport.Client, documentType string) (*Service[T], error)
(*Service[T]).GetPage(ctx context.Context, documentID string, page PageOptions) (*model.Collection[T], error)
```

- documentType: один сегмент `[a-z][a-z0-9]*`; без скрытой нормализации.
- Limit=0 → 1000; допустимы 1…1000, Offset>=0. В query всегда limit/offset.
- Успех: HTTP 200, meta.href непустой, rows — non-null массив; [] допустим.
  ID в T и top-level ID не требуются; размер страницы сервером не пересчитываем.
- Reader — один GET, без auto-pagination, перехода по href и retry. Generic reader
  универсален; набор typed полей конкретного документа этим не объявляется полным.
- Общие JSON-типы не знают HTTP и не запускают бизнес-валидацию.
  Embedded positions записывает Create/Update документа: full replacement без Get/merge.
- Unset/null/0/false/""/[] различаются существующим Field. Кодирование null не означает,
  что сервер разрешает очистку конкретного поля. Expand отсутствует во всех срезах.
- Не вводить registry типов, закрытый union документов, новые транспортные interfaces
  или общий CRUD framework: interface reader и JSON-контракты достаточно малы для текущей задачи.

Подробности: [общие вложенные типы](../document-common/spec.md),
[positions](../document-positions/spec.md), [финансовые связи](../financial-documents/issues/01-finance-operation.md).

## Уточнения draft-контрактов сущностей после проверки

1. Demand.Store — Field[Reference] в Create/Update. Не блокировать пропуск локально:
   таблица требует склад при applicable=true, раздел Create перечисляет его без условия.
   Пример со складом остаётся; успех applicable=false без склада живым API не проверен.
2. PurchaseOrder S03 добавит Shipped/InTransit в общую Position, InTransit Field[float64] —
   в общий PositionRequest; собственная форма строки не вводится.
   Запись inTransit есть в официальных примерах; shipped и waitSum предлагаются только
   для чтения как scope SDK, а не как доказанный readonly-запрет API.
3. SyncID — Field[string] в Create и Update всех этих документов. Убираем произвольное
   ограничение Create-only; API запрещает изменение после заполнения, проверяет сервер.
   SDK не делает Get и не определяет, было ли поле заполнено ранее.
4. CashIn/CashOut — account-поля пока не добавлять по одной противоречивой фразе Update.
   Не переносить polymorphic agent из PaymentOut на CashOut по сходству имён.
5. InvoiceIn/InvoiceOut — Name требуется разделами Create; это не универсальное
   требование имени для всех документов. Счета остаются отдельной позиционной группой.
6. PurchaseReturn использует общий PositionRequest; отдельные request-формы не нужны.
   Embedded price/vat/discount и ограничения возврата без основания ещё требуют выбора
   scope. Эти вопросы не блокируют общий reader, JSON positions и финансовую группу.

Проверенные основания и пределы выводов: [повторная сверка](triage-research.md).
Уточнения 1–3 и типизированный GetPositions внесены в entity spec, но целиком entity
tickets ещё не утверждены. Пять tickets первой волны переведены в ready-for-agent
прямо по команде пользователя, не вследствие готовности зависимостей.

## Разрешённый исполнимый набор

Первая волна из пяти независимых tickets:

- [C01 — DocumentRate](../document-common/issues/01-rate.md).
- [C02 — DocumentOverhead](../document-common/issues/02-overhead.md).
- [P01 — Core/Item и общий JSON](../document-positions/issues/01-json-contracts.md).
- [P02 — GetPage](../document-positions/issues/02-get-page.md).
- [F01 — FinanceOperation JSON](../financial-documents/issues/01-finance-operation.md).

[P03 отменён/заменён](../document-positions/issues/03-movement-contracts.md) одной формой
Position/PositionRequest. Повторно не запускать и не добавлять как зависимость M01.
Согласование этих foundations не означает утверждение всех entity spec.

Разрешены три потока в отдельных worktree: один implementer последовательно выполняет
C01/C02/F01 отдельными коммитами; P01 и P02 выполняются независимо. Общие coverage/README
не изменяются исполнителями — интегратор сведёт их отдельным коммитом после слияния.
Имена файлов P01/P02 заранее разведены в handoff. Независимые Standards/Spec review
координирует интегратор относительно planning BASE. После этой волны остановка.
Никаких push/live API/.env. Service (услуга) в этот запуск не входит.

## Завершение

Первая волна реализована, пять tickets resolved. [Коммиты, проверки и независимые review](review-wave-1.md).
Standards 0, Spec 0 первой волны. P03 впоследствии реализован, затем заменён общей формой;
11 entity tickets не запускались. Следующая группа требует команды.
