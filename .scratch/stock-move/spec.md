Status: needs-triage
Stage: DRAFT

# Move — Get/Create/Update

## Problem Statement

Move (перемещение) переносит запасы между двумя складами. Общий коммерческий документ с agent, discount и НДС не выражает его контракт. Разработчику нужны три операции SDK, чтение ссылок/состояния и безопасное представление positions.

## Solution

Предлагается один вертикальный срез Move: model, отдельные CreateRequest/UpdateRequest, Service.Get/Create/Update, JSON/HTTP-тесты и документация. Общие shapes positions.Movement/MovementRequest и отдельный GetPage остаются общей возможностью SDK. Все choices ниже proposed; сохранение DRAFT не разрешает реализацию.

## User Stories

1. Как разработчик решения, я хочу получить Move по UUID и прочитать Organization и два склада.
2. Как разработчик решения, я хочу создать Move без выдуманного agent.
3. Как разработчик решения, я хочу задать sourceSlot/targetSlot ссылками без реализации services ячеек.
4. Как разработчик решения, я хочу читать metadata positions обычным Get одним запросом.
5. Как разработчик решения, я хочу читать строки явными страницами для больших документов.
6. Как разработчик решения, я хочу обновить одно поле без предварительного Get.
7. Как разработчик решения, я хочу явно заменить весь набор строк и отличать omission от [].
8. Как разработчик решения, я хочу сохранить нулевые/пустые значения и получить фактический ответ/ошибку сервера.

## Implementation Decisions

### Предлагаемый полный scope

| Поля JSON | Read: Go type | Create | Update |
| --- | --- | --- | --- |
| meta; id, accountId | model.Meta; string | — | — |
| created, updated, deleted | model.DateTime | — | — |
| name, description, code, externalCode | string | optional Field[string] | Field[string] |
| moment | DateTime | optional Field[DateTime] | Field[DateTime] |
| applicable, shared | bool | optional Field[bool] | Field[bool] |
| organization, sourceStore, targetStore | *Reference | обязательные Reference | Field[Reference] |
| owner, group, project, state | *Reference | optional Field[Reference] | Field[Reference] |
| internalOrder, customerOrder | *Reference | proposed Field[Reference] | proposed Field[Reference] |
| demand, supply | *Reference | — | —, API read-only |
| productionTasks | []Reference | TODO write | TODO write |
| rate | *model.DocumentRate | Field[model.DocumentRateRequest] | Field[model.DocumentRateRequest] |
| overhead | *model.DocumentOverhead | Field[model.DocumentOverheadRequest] | Field[model.DocumentOverheadRequest] |
| positions | *model.Collection[positions.Movement] | Field[[]positions.MovementRequest] | Field[[]positions.MovementRequest] |
| attributes | []model.Attribute | Field[[]model.Attribute] | Field[[]model.Attribute] |
| syncId | string | optional Field[string] | Field[string], write-once проверяет сервер |
| sum | float64, копейки | — | — |
| printed, published | bool | — | — |
| files | *model.Collection[Reference], только metadata | — | — |

Field означает model.Field с omitzero. Optional read Reference — указатель. Go names — ID/AccountID/SyncID. agent, agentAccount, organizationAccount, store, contract, vatEnabled/vatIncluded/vatSum не добавляются по аналогии: их нет в таблице Move. Минимальный Create — organization, sourceStore, targetStore; name и positions не обязательны. [Таблица и связи](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_move.md#L8-L60), [Create](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_move.md#L317-L323).

Position читает positions.Core(meta/id/accountId), assortment Reference, quantity/price float64, overhead float64 (API Int, копейки), sourceSlot/targetSlot *Reference. Write MovementRequest: matching ID, assortment, quantity, price, sourceSlot/targetSlot через Field; overhead readonly. Discount/vat/vatEnabled не являются частью Move строки. Pack/things — explicit TODO read/write. Assortment — Reference: таблица и Create position включают service, вводное описание — только товары/модификации/партии; это открытая неоднозначность, локального union validator нет. [Позиции](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_move.md#L62-L89), [Create position](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/documents/_move.md#L2044-L2051).

DocumentOverhead имеет sum float64 и distribution weight/volume/price; не смешивается с read-only overhead float64 строки. Без позиций расходы задать нельзя по API; SDK не выполняет Get для проверки. Общий объект не распределяет суммы локально.

### Операции и значения

Get — GET entity/move/{id}; Create — POST entity/move; Update — PUT entity/move/{id}; каждый возвращает (*Move, error). Используется существующий transport.Client, nil отклоняется, UUID и полностью незаданный Update проверяются по образцу UOM. Proposed: обязательные Create refs проверяются только на Meta.Href!="", embedded positions — на максимум 1000. Доступность refs, комбинации складов и наличие строк для overhead проверяет сервер; эти local choices требуют triage. Успех — 200 с непустым id; ошибка — nil-модель. Нет retry, локального merge или скрытых запросов.

Get без body/query получает positions.meta; строки не обещаются без отдельного GetPage. Ни expand, ни переходов по metadata href нет. Если rows уже присутствуют в ответе, shared collection декодирует их, но это не promise маршрута Get.

Omitted Field не отправляется; false/0/""/[] сохраняются. Null передаётся отдельно: общий JSON-контракт не доказывает server-допустимость null для каждого поля. DocumentRate:null сбрасывает валюту на default с пересчётом. Positions omitted сохраняет набор; [] очищает по общему правилу replacement; matching ID обновляется, отсутствующие удаляются, новые добавляются. Positions:null не заменяется на []. Embedded write ограничен API 1000 строками, attributes не подчиняются правилу replacement positions. [Null](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187), [positions](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1031-L1063).

### Зависимости

Единственный [entity ticket](issues/01-move.md) зависит от [DocumentRate](../document-common/issues/01-rate.md), [DocumentOverhead](../document-common/issues/02-overhead.md), [positions JSON](../document-positions/issues/01-json-contracts.md) и [Movement shapes](../document-positions/issues/03-movement-contracts.md). Готовые Supply/Demand и справочники складов/ячеек не нужны: используются Reference. [GetPage](../document-positions/issues/02-get-page.md) — обязательная зависимость entity ticket с типизированным GetPositions.

## Testing Decisions

Публичные JSON-seams encoding/json и encoding/json/v2; HTTP seam — реальный transport.Client + локальный httptest.Server. Обязательный образец структуры — entity/uom/service_test.go: отдельные TestServiceGet/Create/Update с table cases.

Проверить значения всех поддержанных полей неизменённой upstream move fixture, отсутствие/null/unknown, metadata без rows; sourceSlot/targetSlot отдельно и вместе, дробные quantity/price, overhead в формах 0, 0.0 и дробного числа обоими кодеками. Независимые golden минимального и полного Create/Update, omission/false/0/""/[]/null, отсутствие agent/VAT/discount и readonly demand/supply/overhead строки. HTTP-тесты фиксируют method/path/query/body/headers, plain/gzip, один запрос, данные сервера, context/auth/transport/APIError/errors.As/Is, malformed/null JSON, неожиданный status и пустой id.

## Out of Scope

List/query/filter, expand, Delete/batch, templates/metadata, файлы кроме metadata, отдельные position mutations, pack/things/serial/tracking/stock и транспортные механизмы. Эти TODO не объявляются недоступными в API.

## Further Notes

Revision 9a811cef3dc4a3f3a46f070e60452361a06bf7cd проверен исследовательской командой 2026-11-01. [Research](../stock-documents/research.md). Go-код, tests реализации, .env, живой API и commits не выполнялись и не разрешены этим DRAFT.

### Explicit Open Decisions

1. Утвердить таблицу, public names, write internalOrder/customerOrder и proposed SyncID в Create/Update без local write-once validation.
2. Подтвердить service в ассортименте и TODO pack/things; не переносить ограничения вступления на validator.
3. Утвердить локальную валидацию обязательных refs/складов и допустимость null конкретных полей.
4. Подтвердить JSON/HTTP seams и отдельную команду реализации после triage.

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
