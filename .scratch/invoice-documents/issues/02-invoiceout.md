Status: resolved
Stage: IMPLEMENTED
Blocked by: ../../document-common/issues/01-rate.md; ../../document-positions/issues/01-json-contracts.md; ../../document-positions/issues/02-get-page.md

> Реализовано по команде пользователя выполнить все 11 tickets. Предложения и
> запреты запуска в исходном DRAFT ниже отражают этап подготовки; команда исполнения
> их заменила. Итоговый контракт — docs/invoiceout.md, результат — в Answer.

# 02: InvoiceOut — Get/Create/Update с позициями

## What to build

Разработчик читает, создаёт и частично обновляет InvoiceOut, включая коммерческие
позиции и ссылки сторон. Завершённый срез: модель, requests, три HTTP-метода,
оба JSON-кодека, независимые golden, локальные HTTP-тесты и документация.
Это DRAFT по [spec](../spec.md); подробный контракт требует утверждения.

## Типизированное чтение позиций

Entity Service предоставляет GetPositions(ctx context.Context, id string, page positions.PageOptions)
(*model.Collection[Position], error). DocumentType и тип строки фиксированы пакетом; вызывающий
код не передаёт строку ресурса и type parameter. Фиксированный documentType — invoiceout.
Один GET через общий reader, без expand,
follow-up и auto-pagination. Ошибки, UUID/page validation и response envelope — как у P02.
Отдельные TestServiceGetPositions с таблицами проверяют фиксированный resource path, правильный
Position, limit/offset, plain/gzip, пустую страницу, ошибки и ровно один запрос через настоящий transport.
Position — публичный alias positions.Position.
Технический blocker P02 обязателен. Реализация сущности остаётся needs-triage до отдельного разрешения.

## Acceptance criteria

- [x] Зависимости завершены: [общая валюта](../../document-common/issues/01-rate.md)
  и [JSON positions](../../document-positions/issues/01-json-contracts.md).
- [x] Модель InvoiceOut читает общий набор полей и различия CustomerOrder/SalesChannel и readonly Demands по таблицам spec;
  readonly поля, суммы и отсутствующие/неизвестные поля обрабатываются без новых defaults.
- [x] CreateRequest требует Name string, Organization/Agent Reference; прочие согласованные
  поля через Field. UpdateRequest содержит изменяемые Field, включая SyncID; write-once проверяет сервер.
- [x] Предлагаемая validation Name!="" без trim и Meta.Href обязательных ссылок проверена;
  UUID проверяется по форме без ограничения версии, пустой Update отклонён до HTTP.
- [x] Get — GET entity/invoiceout/{id}, Create — POST entity/invoiceout,
  Update — PUT entity/invoiceout/{id}; результат (*InvoiceOut,error) с данными ответа.
- [x] Поле Rate использует DocumentRate/DocumentRateRequest, Positions — общие Position/
  PositionRequest; обычный Get допускает meta без rows, expand/query не добавлены.
- [x] Proposed limit: 1000 embedded позиций проходят; 1001 отклоняется до HTTP; omitted/null/[] не схлопываются и не запускают split или дополнительные запросы.
- [x] Readonly/TODO поля исключены из requests; optional account-ссылки передаются как Reference.
- [x] Embedded positions различают unset/[]/null и selector ID; отправляются одним телом
  без Get/merge/follow-up. Полная замена и правила смены сторон описаны как серверные.
- [x] Значения всех поддержанных upstream полей, optional/null/unknown, даты, false/0/""/[] и
  полный/частичный JSON проверены обоими кодеками с независимыми golden.
- [x] Структура HTTP-тестов — entity/uom/service_test.go: отдельные
  TestServiceGet/Create/Update, таблицы plain/gzip, минимального/полного тела и ошибок.
- [x] Успех — HTTP 200 и ID; проверены context, auth/transport/APIError, errors.As/errors.Is,
  неожиданные статусы, malformed/null JSON, пустой ID и отсутствие retry; ошибки возвращают nil.
- [x] Документация фиксирует точные поля и TODO, ссылки на общий GetPage reader,
  ограничения больших документов и отсутствие expand.
- [x] Пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check и независимые Standards/Spec review относительно исходного BASE.
- [x] После разрешённой реализации — локальный commit, Status: resolved и ## Answer;
  исходные изменения планирования и чужие правки не попали в commit.

## Blocked by

- document-common/01-rate: общая JSON-модель валюты.
- document-positions/01-json-contracts: общие коммерческие позиции и request-формы.

Второй счёт и финансовые/складские документы не блокируют этот ticket.
document-positions/02-get-page обязателен для GetPositions и блокирует этот entity ticket;
попутно не реализуется.

## Scope and sources

Точные поля, readonly/TODO, proposal-validation и Testing Decisions —
в [общей spec](../spec.md). Expand отложен до самого конца SDK.
Никаких live API/.env/push/PR/VPS или общей реорганизации существующих пакетов.
[Первичный источник](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_invoice_out.md).

## Triage 2026-11-01

Пользователь поручил выполнить все 11 entity tickets SDK. Предложенные таблицы полей,
local validation и ограничения scope приняты к реализации этой командой. Общие зависимости resolved.
BASE: eba6f75154929790c5b9d2cb1454dd80783b14db; исходные изменения сохранены.

## Answer

Реализован `entity/invoiceout`: InvoiceOut, отдельные CreateRequest/UpdateRequest,
Get/Create/Update и типизированный GetPositions, aliases общей формы позиции.
[Публичный контракт](../../../docs/invoiceout.md). JSON и HTTP проверки, required/readonly
scope и ограничения покрыты; общие типы не дублируются.

Полный race/vet/gofmt/diff-check PASS; независимые Standards 0 / Spec 0.
Upstream fixtures/provenance неизменны.
[Интеграционный отчёт](../../document-workstream/review-entities.md).
