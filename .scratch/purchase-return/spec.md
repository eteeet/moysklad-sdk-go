Status: needs-triage
Stage: DRAFT

# PurchaseReturn — Get/Create/Update

## Problem Statement

Разработчику решения нужны чтение, создание и частичное изменение PurchaseReturn
(возврат поставщику), включая возвраты с основанием Supply и без него. Обычный
коммерческий request не выражает запреты изменения сторон и особенности возвращаемых
строк. Поля и методы ниже — предложения для согласования, не разрешение реализации.

## Solution

Один законченный срез: PurchaseReturn, отдельные CreateRequest/UpdateRequest,
Service.Get/Create/Update, публичные JSON/HTTP-тесты и документация. Общие позиции
используются независимо от Supply; reader строк — общий positions.GetPage.
Expand не разрабатывать до завершения всего остального SDK и только если понадобится.

## User Stories

1. Как разработчик решения, я хочу получить возврат по UUID и прочитать его состояние.
2. Как разработчик решения, я хочу увидеть основание Supply как Reference без его загрузки.
3. Как разработчик решения, я хочу создать возврат по основанию с явными сторонами и складом.
4. Как разработчик решения, я хочу создать возврат без основания, не передавая supply.
5. Как разработчик решения, я хочу передать допустимые возвращаемые строки внутри Create.
6. Как разработчик решения, я хочу изменять разрешённые поля без предварительного Get.
7. Как разработчик решения, я хочу исключить agent/agentAccount/supply из UpdateRequest.
8. Как разработчик решения, я хочу отличать positions metadata от реально загруженных строк.
9. Как разработчик решения, я хочу читать строки явными страницами через общий reader.
10. Как разработчик решения, я хочу различать omission, [], null, 0 и false в запросе.
11. Как разработчик решения, я хочу понимать full replacement и ограничения основания до отправки запроса.
12. Как разработчик решения, я хочу получить именно ответ сервера и распознаваемую ошибку.

## Implementation Decisions

Все новые имена, scope записи и local validation — PROPOSED. R — чтение; C/U —
отдельные requests. Field означает существующий model.Field[T] с omitzero.
Optional ссылки чтения — *model.Reference, даты — существующий model.DateTime.

### Полный предлагаемый scope документа

| JSON-поля | R: тип Go | C | U |
| --- | --- | --- | --- |
| meta; id, accountId | model.Meta; string | — | — |
| created, updated, deleted | model.DateTime | — | — |
| name, description, code, externalCode | string | Field[string] | Field[string] |
| moment | model.DateTime | Field[DateTime] | Field[DateTime] |
| applicable, shared, vatEnabled, vatIncluded | bool | Field[bool] | Field[bool] |
| organization, store | *Reference | обязательные Reference | Field[Reference] |
| agent | *Reference | обязательный Reference | —, запрет Update |
| agentAccount | *Reference | Field[Reference] | —, запрет Update |
| supply | *Reference | Field[Reference], обязательно только для возврата по основанию | —, запрет Update |
| organizationAccount, owner, group, project, contract, state | *Reference | Field[Reference] | Field[Reference] |
| rate | *model.DocumentRate | Field[model.DocumentRateRequest] | Field[model.DocumentRateRequest] |
| positions | *model.Collection[positions.Commercial] | Field[[]CreatePositionRequest] | Field[[]UpdatePositionRequest] |
| attributes | []model.Attribute | Field[[]model.Attribute] | Field[[]model.Attribute] |
| syncId | string | Field[string] | Field[string], write-once проверяет сервер |
| sum, vatSum, payedSum | float64 | — | — |
| printed, published | bool | — | — |
| factureIn, factureOut | *Reference | TODO write | TODO write |
| payments | []Reference | TODO write | TODO write |
| files | *model.Collection[model.Reference], только metadata | — | — |

Go names следуют существующим ID/AccountID/SyncID/VATEnabled/VATIncluded/VATSum.
Минимальный Create без основания: organization, agent, store. Name и positions
не обязательны. Для возврата по основанию дополнительно задаётся supply.
Отдельный Based/Unbased enum, автоматическое создание по шаблону и поиск основания
не нужны. Возможность пропустить supply не доказывает применимость ограничений
основания к возврату без него: эту неоднозначность сохранить явно.

SyncID предлагается в Create/Update. API запрещает изменение после заполнения,
а не первоначальное заполнение существующего документа; SDK оставляет проверку серверу.
Readonly суммы, даты и flags не попадают в requests; массивы связей читаются без
автоматической загрузки. TODO записи связей не означает запрет этих операций API.

### Отдельные формы записи строк — PROPOSED, требуют triage

Чтение использует positions.Commercial: meta/id/accountId, assortment Reference,
quantity/price/discount float64, vat int, vatEnabled bool. Это переиспользование
формы JSON, а не обещание редактирования всех её полей.

- CreatePositionRequest: Assortment Field[Reference], Quantity/Price Field[float64],
  VAT Field[int], VATEnabled Field[bool]. Без ID, Meta, AccountID и Discount.
- UpdatePositionRequest встраивает positions.ItemRequest (matching ID, assortment,
  quantity) и добавляет Price Field[float64], VAT Field[int], VATEnabled Field[bool].
  Нет Meta, AccountID и Discount. JSON плоский, все Field с omitzero.
- CommercialRequest не используется автоматически. Discount в таблице возврата
  readonly, хотя примеры передают его: proposed scope не включает запись Discount.
- Таблица и embedded Update example включают price/vat; правило «только quantity»
  явно относится к отдельному PUT /positions/{positionID}, который вне scope.
  Его нельзя молча переносить на embedded Update. Предложенный write scope
  price/vat/vatEnabled необходимо утвердить с учётом ограничений основания.
- Pack, slot, things, маркировка и stock enrichment — TODO read/write. Ассортимент
  передаётся существующим Reference; services product/service/variant/consignment
  не становятся зависимостями, закрытого union validator нет.

### Операции, ограничения и ошибки

NewService(*transport.Client) (*Service, error), Get(ctx,id), Create(ctx,CreateRequest),
Update(ctx,id,UpdateRequest) возвращают (*PurchaseReturn,error). Nil client — ошибка.
Get — GET entity/purchasereturn/{id}; Create — POST entity/purchasereturn;
Update — PUT entity/purchasereturn/{id}. Без query, скрытых запросов и retry.
Предлагаемый успех — HTTP 200 и непустой ID; при любой ошибке nil-модель.
Errors.As/Is, context, gzip и operation context сохраняются как у UOM.

Proposed local validation: форма UUID без ограничения версии, непустой Meta.Href
трёх обязательных Create References, Update хотя бы с одним заданным Field,
embedded positions не более 1000. Предел 1000 проходит, 1001 отклоняется до HTTP.
Ссылки, остатки, сравнение с основанием, положительность quantity и сочетания
счетов не проверяются локальным чтением. Эти validation choices ещё не утверждены.

При возврате по основанию сервер требует совпадения agent/organization/валюты
с Supply; строки должны соответствовать основанию, quantity не превышает исходное.
Update запрещает agent/agentAccount/supply и другую валюту, не позволяет новые
отсутствующие в документе строки. Смысл «документ» и применение этих фраз к возврату
без основания неоднозначны — SDK не воспроизводит их новым бизнес-валидатором.
Общая фраза Update про одновременную смену agent и agentAccount противоречит их
прямому запрету: в предложенном UpdateRequest обоих полей нет. OrganizationAccount
доступен вместе с Organization; согласованность оставлена серверу.

Обычный Get читает positions.meta одним запросом. Строки — отдельный явный вызов
positions.Service[positions.Commercial] с documentType=purchasereturn; его наличие
не означает полного typed покрытия всех особенностей возврата.

Unset не отправляется; false/0/"" сохраняются. Positions omitted сохраняет набор,
[] запрашивает полную замену пустым, непустой массив — весь новый набор с matching
ID. Локального merge нет. Null остаётся wire null, не преобразуется в [] и не
объявляется разрешённым сервером для каждого поля. Rate:null может сбросить валюту
с пересчётом сервером, но это не отменяет ограничения валюты возврата по основанию.
Attributes следует собственному контракту, не replacement-rule positions.

### Зависимости

[Ticket 01](issues/01-purchasereturn.md) зависит только от
[DocumentRate](../document-common/issues/01-rate.md) и
[JSON positions](../document-positions/issues/01-json-contracts.md).
Готовый Supply и другие services не требуются. [Reader](../document-positions/issues/02-get-page.md)
— обязательная зависимость entity ticket с типизированным GetPositions.

## Testing Decisions

Публичные JSON-типы через encoding/json и encoding/json/v2; HTTP — настоящий
transport.Client и локальный httptest.Server. Структура строго как
entity/uom/service_test.go: отдельные TestServiceGet/Create/Update с table cases.
Не закреплять private do/validID или новый mock framework.

- Неизменённая upstream fixture с provenance: значения всех поддержанных полей,
  даты/ссылки/суммы, readonly Discount, отсутствие/null/unknown и positions metadata.
- Независимые golden минимального Create с основанием/без него, полного Create,
  полного и одно-полевого Update, false/0/""/[]/null, matching ID и частичного Rate.
- Отсутствие agent/agentAccount/supply в UpdateRequest, readonly/TODO в записи;
  различия CreatePositionRequest и UpdatePositionRequest, отсутствие Discount.
- Exact method/path/query/body/headers, plain/gzip, один запрос, серверный ответ
  отличается от input; Get не заполняет Rows и не вызывает reader.
- Nil/cancelled context, nil client, malformed UUID, empty required href/Update,
  1000/1001 позиций; auth/transport/APIError 4xx/5xx/429 и errors.Is/As;
  unexpected status, malformed/null JSON, пустой/отсутствующий ID, nil-модель ошибки.
- Синтетические APIError проверяют передачу ошибок основания/валюты, а не доказывают
  бизнес-поведение сервера. Нет обращения к живому Account.
- После реализации: gofmt, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check и независимые Standards/Spec review относительно фиксированного BASE.

## Out of Scope

List/query/filter, expand, Delete/batch, templates и создание /new по основанию;
отдельные HTTP mutations позиций; файлы кроме metadata; pack/slot/serial/tracking/stock;
автоматический обход страниц, финансовая арифметика, локальная имитация возврата,
новый CRUD/transport framework и рефакторинг готовых сущностей.
Текущая работа не включает Go-код, .env, live API, push, PR и commits планирования.

## Further Notes

Официальный revision 9a811cef3dc4a3f3a46f070e60452361a06bf7cd, подтверждён 2026-11-01.
Факты и противоречия также собраны в [исследовании документов](../stock-documents/research.md).

- [Поля, связи и строки](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L4-L83).
- [Create с основанием и без него](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L567-L583).
- [Запреты Update](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L2180-L2192).
- [Embedded Update example](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L2248-L2328).
- [Quantity-only отдельного PUT позиции](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_purchase_return.md#L2852-L2855).

### Explicit Open Decisions

1. Утвердить scope документа и отдельные формы записи позиции; price/vat/vatEnabled
   embedded Update и исключение Discount требуют явного решения, не ready по умолчанию.
2. Уточнить ограничения основания для возврата без supply и конфликт agent/account Update.
3. Утвердить local validation, SyncID в Create/Update и выбранные readonly связи/TODO.
4. Подтвердить общие public types и тестовые seams; затем отдельно разрешить реализацию.

## Типизированное чтение позиций — согласованное уточнение

У каждого entity Service добавить GetPositions(ctx context.Context, id string, page positions.PageOptions)
(*model.Collection[Position], error). Он вызывает общий reader с фиксированными documentType
и Position; вызывающий код не выбирает строку ресурса и generic type. Один HTTP GET, без expand,
автообхода и скрытого Get заголовка. Ошибки/envelope/UUID/page validation остаются контрактом P02.
Move экспортирует alias Position = positions.Movement; InvoiceIn/InvoiceOut/PurchaseReturn —
alias Position = positions.Commercial. Supply/Demand/PurchaseOrder используют собственные Position.
Тестовый seam — публичный GetPositions через transport.Client + httptest.Server;
отдельный TestServiceGetPositions с таблицами правильного path/type/query, plain/gzip, пустой страницы,
ошибок и одной отправки. P02 — обязательная зависимость всего entity ticket.
Это уточнение не разрешает реализацию entity, остальные open decisions остаются needs-triage.
