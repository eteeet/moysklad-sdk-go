# Документы: общий контракт методов

`NewService(*transport.Client)` отклоняет nil. `Get(ctx, id)`,
`Create(ctx, CreateRequest)` и `Update(ctx, id, UpdateRequest)` возвращают
указатель на модель документа и error. Get выполняет GET `entity/<type>/{id}`,
Create — POST `entity/<type>`, Update — PUT `entity/<type>/{id}`.
Успех — HTTP 200 с непустым ID; при ошибке модель nil. Возвращается ответ сервера,
без локального слияния с input. Ошибки сохраняют контекст операции и
`errors.Is`/`errors.As`, в том числе `*errors.APIError`.

Каждый вызов делает один запрос. Get не имеет тела и query; нет expand,
предварительного Get, переходов по href или автоматических повторов.
После неоднозначного сбоя POST проверь результат на сервере до повторного создания.

Get/Update проверяют форму UUID без ограничения версии. Create проверяет только
непустой Meta.Href обязательных ссылок, а у InvoiceIn/InvoiceOut, RetailSalesReturn, InternalOrder и
CommissionReportOut также Name != "" без trim; у CommissionReportIn/Out и Payroll — заданные даты
периода, у PriceList — непустые columns. Update отклоняет полностью незаданный request. Name остальных документов
необязателен; SDK не подставляет суммы, даты, имена и другие defaults.
Типы/доступность ссылок, требуемые сочетания полей и бизнес-правила проверяет API.

Все optional Create и все Update поля используют `model.Field[T]` с `omitzero`:
unset пропускается; false, 0, "", null и [] передаются как заданы; nil slice в Value
кодируется как []. Возможность закодировать null не обещает его принятия сервером.
`contract:null` сбрасывает необязательную связь; `rate:null` запрашивает валюту по
умолчанию с пересчётом сервером. В присутствующем rate отсутствие value означает
курс 1; SDK не подставляет 1 в JSON. `attributes` обновляются по metadata, [] не
обещает очистку всей коллекции; отдельное значение сбрасывается через value:null.

SyncID доступен в Create и Update. После заполнения менять его запрещает API;
SDK передаёт запрос без проверки прежнего значения. Все суммы и денежные поля —
float64 в копейках, без округления, пересчётов и гарантии десятичной точности либо
точности integer за пределами 2^53. Курс, количества и процент Discount не являются
денежными суммами. Отсутствующие поля чтения дают zero values, неизвестные игнорируются.

## Позиции

InvoiceIn, InvoiceOut, Supply, Demand, PurchaseOrder, Move, PurchaseReturn, CustomerOrder, Enter,
Loss, SalesReturn, RetailDemand, RetailSalesReturn ([розница](retail.md)) и InternalOrder используют
общие `positions.Position`/`PositionRequest`. Каждый entity-пакет экспортирует aliases
`Position` и `PositionRequest`, а Service предоставляет типизированный
`GetPositions(ctx, id, positions.PageOptions) (*model.Collection[Position], error)`.
Ресурс и форма позиции фиксированы пакетом. Одна страница, Limit=0 означает 1000,
допустимы 1…1000 и Offset>=0. Обход страниц выполняет вызывающий код.
Prepayment и PrepaymentReturn позиции только читают: экспортируют alias `Position` и тот же
`GetPositions`, `PositionRequest` и операций записи позиций у них нет.

Get документа читает positions.meta, без обещания rows и без вызова reader.
GetPositions делает отдельный GET `entity/<type>/{id}/positions?limit=...&offset=...`;
`PageOptions.Fields` добавляет `fields=stock` или `fields=declaration`, если источник документирует
их для этого документа.
У страницы требуются HTTP 200, непустой meta.href и non-null rows; [] — успешная
пустая страница. Дополнительные правила: [positions](positions.md).

Embedded Create/Update принимает до 1000 позиций; 1001 отклоняется до HTTP. В CreateMany/UpdateMany
лимит позиций каждого элемента проверяет сервер. У Inventory, PriceList, CommissionReportIn/Out, Payroll, EmissionOrder и RetireOrder
своя форма позиции (у PriceList — строки); у Inventory и PriceList positions в теле — объект `{"rows":[…]}`.
Незаданные positions не меняют набор; [] запрашивает замену пустым набором.
Непустой positions в Update полностью заменяет набор: matching ID обновляет позицию,
отсутствующие позиции удаляются, новые добавляются. ID не задают для новой позиции Create.
SDK не выполняет merge, разделение массива на запросы или проверку позиций по основанию.
Null остаётся wire null, не превращается в []; принятие сервером не гарантируется.

Price/Overhead/Cost — float64 в копейках; Quantity/Shipped/InTransit — количества,
Discount — процент (отрицательный означает наценку). Overhead и Shipped только читаются.
Cost записывается для услуг Demand и возвратов без основания; SDK не определяет тип assortment
локально. Reserve пишется у заказа покупателя, Reason — у оприходования и списания.
InTransit доступен в общем request для PurchaseOrder; допустимость заданных полей
общей формы у остальных документов определяет API. SourceSlot/TargetSlot —
optional Reference; Overhead позиции отличается от объекта overhead документа.

Явная пометка «Только для чтения» в таблице атрибутов имеет приоритет над
противоречащими примерами запросов. Так, overhead позиции Supply остаётся readonly,
хотя встречается в примерах записи. Отсутствие пометки само по себе не доказывает
возможность записи: её подтверждают описание метода и применимые примеры
с учётом ограничений поля, как для Cost и InTransit.

## Проверки и границы

Модели/requests проверены encoding/json и encoding/json/v2, включая independent golden,
сырые Attribute long >2^53 и неизменённые upstream fixtures с
[provenance](upstream.md). Методы проверяются настоящим transport.Client и локальными
HTTP-серверами: plain/gzip, exact wire JSON, headers, одна отправка, validation,
context, auth/transport/APIError, malformed/null/empty ответы и ошибочные статусы.
Это не проверка бизнес-поведения живого API. Живой аккаунт не использовался.

Реализованы: List одной страницы с filter/search/order ([transport](transport.md)), явный обход
страниц через `pages.All`, поля позиций общей формы (pack, things, gtd, trackingCodes, declaration,
stock — [positions](positions.md)) и `fields=stock`/`declaration` у GetPositions.
Отдельные операции позиций, шаблоны `New`, `Trash`, связи и wizard —
[documents-extra.md](documents-extra.md). Expand (в том числе stock в Get/List через
expand=positions), автоматический обход страниц и финансовые вычисления остаются вне этих срезов. Delete и массовые операции документов —
[write-async.md](write-async.md).
Точный scope полей и исключений каждого документа указан в его документации.
