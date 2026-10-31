Status: needs-triage
Stage: PARTIALLY IMPLEMENTED

# Общий модуль positions

Первая волна P01/P02 утверждена пользователем 2026-11-01: Core/Item/Commercial JSON
и generic GetPage с правилами validation/ответа ниже. Реализация и локальные коммиты
разрешены только этим двум tickets. Movement P03 и entity tickets остаются DRAFT.
Общие денежные поля — float64, включая overhead/cost независимо от API Int/Float.

P01/P02 завершены и включены в dev: JSON-контракты и reader одной страницы.
[Публичная документация](../../docs/positions.md), [итоги review](../document-workstream/review-wave-1.md).
P03 Movement и entity tickets остаются needs-triage и этим этапом не запускаются.

## Problem Statement

Разработчику решения нужно отдельно читать строки документа и передавать позиции
при создании или изменении документа. Сейчас SDK не предоставляет такого контракта.
Обычный Get документа возвращает metadata positions; количество строк может превышать
1000. Привязка общей логики к supply заставит копировать транспорт в других документах,
а один универсальный write type ошибочно разрешит коммерческие поля у move или возврата.

## Solution

Предоставить самостоятельный package `positions`: reusable JSON-типы и generic reader
одной страницы `/entity/{documentType}/{documentID}/positions`. Документные packages
подключают общие типы и дополняют их только своими полями и правилами. Reader работает
с выбранным вызывающим кодом типом строки; новые документы не требуют расширять allowlist.

Позиции при Create/Update документа передаются массивом. Отсутствие поля, пустой массив
и непустой массив различаются явно. SDK отправляет выбранный набор за один запрос,
сохраняет результат сервера и не выполняет предварительное чтение или локальное слияние.

## User Stories

1. Как разработчик решения, я хочу читать страницу позиций документа, чтобы обработать строки без дополнительных режимов запроса.
2. Как разработчик решения, я хочу задавать limit и offset, чтобы читать документы с количеством строк больше 1000.
3. Как разработчик решения, я хочу получать metadata страницы и строки вместе, чтобы самому выбирать следующую страницу.
4. Как разработчик решения, я хочу получать пустую страницу как успешный результат, чтобы корректно завершить свой обход.
5. Как разработчик решения, я хочу использовать один reader для разных document types, чтобы не дублировать работу с HTTP.
6. Как разработчик решения, я хочу подключать собственный тип строки, чтобы читать документы с полями, отличными от коммерческой позиции.
7. Как разработчик решения, я хочу сохранять ссылки assortment, чтобы работать с поддерживаемыми API типами без полной модели связанной сущности.
8. Как разработчик решения, я хочу читать quantity, price, discount, vat и vatEnabled коммерческой позиции, чтобы сохранять её фактические значения.
9. Как разработчик решения, я хочу различать vatEnabled=false и vat=0, чтобы сохранять отличие «без НДС» от ставки 0%.
10. Как разработчик решения, я хочу передавать дробные quantity и price, чтобы не терять значения API.
11. Как разработчик решения, я хочу передавать отрицательный discount, чтобы задавать документированную наценку.
12. Как разработчик решения, я хочу создать документ с массивом позиций, чтобы выполнить одну согласованную операцию.
13. Как разработчик решения, я хочу не передавать positions при Update, чтобы сохранить существующий набор строк.
14. Как разработчик решения, я хочу передать positions=[], чтобы явно заменить набор позиций пустым.
15. Как разработчик решения, я хочу указать id существующей строки в embedded Update, чтобы сервер сопоставил её с изменяемой позицией.
16. Как разработчик решения, я хочу явно передавать 0, false и null там, где это позволяет JSON-контракт, чтобы кодирование не теряло намерение запроса.
17. Как разработчик решения, я хочу использовать отдельную форму Movement без коммерческих полей, чтобы не отправить неподдерживаемые discount или vat.
18. Как разработчик решения, я хочу читать исходную и целевую ячейки Movement как ссылки, чтобы не загружать справочник склада автоматически.
19. Как разработчик решения, я хочу получать исходные APIError и context errors, чтобы обрабатывать отказ сервера и отмену программно.
20. Как разработчик решения, я хочу видеть ошибку неверного UUID или страницы до HTTP-запроса, чтобы не обращаться к заведомо неверному адресу.
21. Как разработчик решения, я хочу получать именно ответ сервера и ровно одну страницу, чтобы число запросов и бюджет API оставались управляемыми.
22. Как разработчик решения, я хочу видеть границы общих типов, чтобы не считать позиции payroll, inventory и returns полностью реализованными.

## Implementation Decisions

### Владение и публичный интерфейс — PROPOSED

`positions` — самостоятельный общий модуль SDK. Его не размещают внутри supply и не
делают владельцем модели конкретного документа. Имя `positions.Service` — имя Go-типа
reader; оно не означает реализацию сущности API Service (услуга), которая отложена.

| Публичный элемент | Контракт |
| --- | --- |
| `Service[T any]` | Reader строки T без требования реализовать interface или встроить общий тип |
| `NewService[T](client *transport.Client, documentType string) (*Service[T], error)` | Сохраняет существующий transport.Client и API-код документа; nil client отклоняется |
| `(*Service[T]).GetPage(ctx context.Context, documentID string, page PageOptions) (*model.Collection[T], error)` | Один GET страницы; отсутствуют тело и побочные запросы |
| `PageOptions.Limit int` | 0 означает default 1000; 1…1000 передаётся без изменения; остальные значения отклоняются |
| `PageOptions.Offset int` | 0 и положительные значения допустимы; отрицательные отклоняются |

- Endpoint: `entity/{documentType}/{documentID}/positions`. Query содержит только явные
  `limit` и `offset`; default кодируется как `limit=1000&offset=0`.
- `documentType` — открытый API-код, PROPOSED syntax `[a-z][a-z0-9]*`. Пустая строка,
  path/query/fragment и закодированные разделители отклоняются; нормализации нет.
  Allowlist первых документов отсутствует. Проверка синтаксиса не обещает, что сервер
  поддерживает такой ресурс: неподдерживаемый код даёт обычный APIError.
- Reader одинаково принимает текущие и будущие документные types, имеющие этот endpoint.
  `T any` позволяет читать специфичные inventory, payroll и return формы без изменения
  reader и без обязательных assortment, quantity, price или id в T.
- Не создавать новые List/collection wrappers: `model.Collection[T]` уже содержит
  optional Context, CollectionMeta и Rows. Общие Field/Reference/transport сохраняются.
- Обычный Get заголовка документа возвращает `positions` как metadata collection;
  документ может использовать `*model.Collection[Position]`. Наличие meta не означает,
  что Rows были прочитаны. Get не вызывает GetPage автоматически.

### Общие формы чтения — PROPOSED

Встраивание типов ниже даёт плоский JSON. `Core` — reusable identity для известных
форм, а не обязательный interface для всех позиций API. `Item` добавляет ассортиментную
часть. `Commercial` подходит коммерческим документам в согласованном наборе полей;
`Movement` имеет отдельный контракт. Это не один закрытый union всех позиций.

| JSON | Go field | Тип чтения | Владелец формы |
| --- | --- | --- | --- |
| meta | Meta | model.Meta | Core |
| id | ID | string | Core |
| accountId | AccountID | string | Core |
| assortment | Assortment | model.Reference | Item, встраивает Core |
| quantity | Quantity | float64 | Item |
| price | Price | float64 | Commercial, встраивает Item |
| discount | Discount | float64 | Commercial |
| vat | VAT | int | Commercial |
| vatEnabled | VATEnabled | bool | Commercial |
| price | Price | float64 | Movement, встраивает Item |
| overhead | Overhead | float64 | Movement, только чтение |
| sourceSlot | SourceSlot | *model.Reference | Movement, optional |
| targetSlot | TargetSlot | *model.Reference | Movement, optional |

- Price сохраняет существующую convention SDK `float64`, в копейках. Все денежные поля, включая API Int overhead/cost, читаются как float64: в официальных
  примерах встречается 0.0. Overhead строки — скаляр, не document-level объект расходов. Нового money type нет.
- В read-моделях отсутствующие скаляры дают zero value; отсутствующие optional ссылки
  дают nil. Неизвестные поля ответа не мешают декодированию.
- Quantity, Price, Discount и VATEnabled не получают проверки обязательности при
  простом JSON-декодировании; успешный HTTP collection не проверяется как одиночная сущность.
- Дополнительные поля принадлежат типу конкретного документа: supply/demand overhead,
  demand cost, purchaseorder shipped/inTransit и специальные поля возврата не добавляются
  в Commercial автоматически. Document Position может встроить Commercial и добавить
  свои read fields. Write-доступность каждое такое расширение устанавливает отдельно.
- Country, GTD, Pack, Things, TrackingCodes, stock enrichment и прочие поля не считаются
  реализованными общим базовым типом; документный контракт может добавить их отдельным срезом.

### Формы embedded записи — PROPOSED

Каждое optional поле ниже имеет `json:",omitzero"`. Встраивание request types плоское;
модели чтения не используются как request. Новая строка коммерческого документа обычно
передаёт assortment и положительный quantity; точные требования определяет документ/API,
а не универсальная бизнес-валидация общей библиотеки.

| JSON | Go field | Тип записи | Форма и допустимость |
| --- | --- | --- | --- |
| id | ID | model.Field[string] | CoreRequest: только matching selector существующей строки embedded Update |
| assortment | Assortment | model.Field[model.Reference] | ItemRequest, встраивает CoreRequest |
| quantity | Quantity | model.Field[float64] | ItemRequest |
| price | Price | model.Field[float64] | CommercialRequest, встраивает ItemRequest |
| discount | Discount | model.Field[float64] | CommercialRequest |
| vat | VAT | model.Field[int] | CommercialRequest |
| vatEnabled | VATEnabled | model.Field[bool] | CommercialRequest |
| price | Price | model.Field[float64] | MovementRequest, встраивает ItemRequest |
| sourceSlot | SourceSlot | model.Field[model.Reference] | MovementRequest |
| targetSlot | TargetSlot | model.Field[model.Reference] | MovementRequest |

- Meta и AccountID отсутствуют во всех request types; Overhead отсутствует в
  MovementRequest. ID — явное исключение: таблицы называют его read-only, но общее
  правило embedded Update требует его для сопоставления существующей позиции.
  Это не назначение client-generated ID. При Create новой строки ID не передают.
- ID-only или sparse existing-row JSON кодируется без локального достраивания строки.
  Допустимость конкретного набора полей и ссылки на существующую строку проверяет сервер.
- CommercialRequest нельзя автоматически использовать для purchasereturn: ограничения
  исходного supply требуют отдельного document-owned request. Правило quantity-only
  описано для отдельного PUT позиции и не переносится автоматически на embedded Update.
  У payroll и inventory также могут быть полностью собственные T/request types.
- Field различает unset, Value и Null. Явные 0/false/отрицательный discount сохраняются.
  Возможность закодировать null не обещает, что сервер разрешает null в конкретном поле.

### Передача positions в Create/Update документа

Общий модуль задаёт контракт `model.Field[[]T]` с omitzero в request документа;
конкретный T — утверждённый для этого документа request позиции. Новый HTTP writer
общего модуля не вводится; запись выполняет Create/Update документа.

| Состояние Field | JSON документа | Значение |
| --- | --- | --- |
| unset | positions отсутствует | При Update существующие строки сохраняются |
| Value с пустым или nil slice | `"positions":[]` | Полная замена пустым набором |
| Value с непустым slice | `"positions":[…]` | Передан весь новый набор; matching id обновляется, новые строки добавляются, отсутствующие старые удаляются |
| Null | `"positions":null` | Сериализуемый null; не объявляется поддерживаемой очисткой |

- Field уже сериализует Value(nil slice) как [] обоими поддерживаемыми JSON codecs.
- Embedded write ограничен API 1000 строками; документные tickets должны закрепить
  предел 1000/1001 в своей HTTP-валидации. Самостоятельная JSON-модель не выполняет HTTP
  и не вводит новый Validate/helper API. Тип позиции, обязательность, положительность
  quantity и ограничения исходного документа не превращаются в новые общие проверки.
- Update не выполняет GET для загрузки текущих строк, merge или дописывание пропущенных
  позиций. Локальные HTTP-тесты доказывают отправленное тело, а не бизнес-эффект живого API.
- После Create/Update возвращается состояние ответа сервера. Отсутствующие Rows не
  дополняются из request; повторного Get/GetPage для заполнения ответа нет.

### Валидация и ответы reader — PROPOSED

- documentID проверяется как UUID по существующему образцу: форма 36 символов, hex
  группы и дефисы; версия UUID не ограничивается. Неверные type/ID/page дают nil+error
  до HTTP. Не проверять UUID версии, meta href или допустимость assortment локально.
- Успех reader требует HTTP 200 и collection object с непустым meta.href и rows-array.
  `rows:[]` — успешная пустая страница. Пустое тело, malformed JSON, null, {}, отсутствие
  meta.href, отсутствие rows или rows:null дают nil+error. Это PROPOSED минимальная
  проверка envelope, не проверка id отдельной сущности.
- Top-level id не требуется; T не обязан иметь ID. Не требовать size==len(rows),
  size>=offset, наличия nextHref или равенства серверных limit/offset запросу.
- Возвращать Context, Meta, Rows и href страницы из ответа сервера без вычисления
  отсутствующих полей. NextHref/PreviousHref — данные, reader по ним не ходит.
- Использовать существующий transport.Client, context, заголовки и gzip. nil context,
  отмена, authorization и transport errors сохраняют доступность через errors.Is;
  APIError доступен через errors.As. Контекст операции включает documentType/documentID;
  префикс moysklad: не дублируется. При любой ошибке результат nil, повторов нет.

## Testing Decisions

JSON/HTTP seams P01/P02 утверждены пользователем; Movement P03 остаётся DRAFT.
Тест проверяет внешнее поведение;
private helpers, внутреннее устройство Service и мок приватного do не закрепляются.

- Самостоятельный JSON-срез проверяется через публичные types и оба codec:
  `encoding/json` и `encoding/json/v2`. Prior art — таблица codecs UOM и проверка
  самостоятельных общих Address/AddressRequest, а не тестирование приватных функций.
- Проверить значения каждого заявленного поля Core/Item/Commercial/Movement; число с
  дробной частью, отрицательный discount, VAT=0+VATEnabled=false и VAT=0+true, optional
  refs absent/present, неизвестное поле. Полного декодирования без assertions недостаточно.
- Golden JSON задаётся независимо от сериализации SDK. Проверить минимальную новую
  строку, existing-row id matching request, sparse quantity-only request, все writable
  поля, явные 0/false, null, omission и отсутствие read-only полей в request API.
- На небольшом тестовом container с Positions Field[[]T] проверить unset/[]/nil slice/
  непустой slice/null обоими codecs. Это завершённый публичный JSON-контракт, который
  проверяем без реализации supply или HTTP writer.
- Прикладной читатель тестировать через публичный GetPage, локальный HTTP-сервер и
  настоящий transport.Client. Структурный образец — отдельные TestServiceGet UOM с
  табличными сценариями; здесь выделить самостоятельные TestServiceGetPage,
  TestServiceGetPageErrors, TestServiceGetPageValidation и проверки error wrapping.
- Проверить точные метод/path/limit/offset, отсутствие тела, заголовки, plain/gzip,
  collection Context/Meta/Rows и пустую страницу; default и границы 1/1000, offset>0.
  Два явно вызванных GetPage возвращают две независимые страницы без auto-pagination.
- Проверить type вне первых document families и пользовательский T, не содержащий
  Core/assortment/id: это доказывает generic contract без обещания бизнес-полей payroll.
- Проверить отказ без HTTP при invalid type, malformed UUID, Limit<0/Limit>1000,
  Offset<0, nil client/context и отмене; допустимость UUID нестандартной версии.
- Проверить 4xx/5xx и APIError, auth/transport/context errors, unexpected 201/204,
  malformed JSON и все invalid envelopes выше. Rows=[] без top-level id принимается.
- Счётчик запросов должен оставаться равен одному даже при nextHref, metadata links,
  неизвестных полях, неполных строках или ошибке; переход по ссылкам отсутствует.
- Upstream fixtures использовать неизменёнными с текущим provenance. Если fixture
  описывает заголовок или иную операцию, не выдавать её за /positions response:
  отдельный небольшой официальный page example или synthetic envelope хранить отдельно.
- Для каждого будущего implementation ticket: gofmt, `go test -race -count=1 ./...`,
  `go vet ./... ./testdata`, `git diff --check` и независимое Standards/Spec review.
  Для P01/P02 реализация и проверки разрешены первой волной; P03 этим не запускается.

## Out of Scope

- Любой expand: возвращаться к нему только после всего остального SDK и если потребуется.
- Автоматический обход страниц, iterators/walker, скрытые followups и чтение ссылок.
- Отдельные HTTP Get by position ID, Create/Update/Delete позиции, batch и positions/delete.
  Существование этих endpoints в API не означает готовность или согласование SDK-методов.
- Реализация Create/Update/Get заголовков документов и return-specific business rules.
- Полное покрытие всех вариантов строк API одним type, закрытый DocumentType enum,
  универсальный CRUD/transport framework и рефакторинг готовых packages.
- Сущность Service, DocumentRate, дополнительные справочники и fields=stock.
- Pack, serial numbers, tracking, GTD и прочие document-specific расширения вне базовых
  таблиц выше; наличие raw JSON неизвестного поля не считается typed поддержкой.
- Retry/429/rate limiting, live API, .env, push, публикация PR и commits подготовки.

## Further Notes

- Подтверждённые первичные факты и отличия документов: [research](research.md).
- Разбиение DRAFT: [01 JSON](issues/01-json-contracts.md), [02 reader](issues/02-get-page.md),
  [03 Movement](issues/03-movement-contracts.md). 01 и 02 независимы; 03 зависит от 01.
- Общий JSON-контракт — отдельный завершённый срез. Reader уже использует существующий
  Collection[T] и не требует завершения Commercial. Movement требует Core/Item из 01.
- Document tickets зависят от нужного JSON-среза и reader 02, поскольку включают
  типизированный GetPositions с фиксированными documentType и Position.
- Export names P01/P02, syntax documentType, PageOptions, envelope и seams утверждены
  для первой волны. P03 и entity tickets требуют отдельного triage/команды.

## Comments
