Status: resolved
Stage: IMPLEMENTED
Blocked by: None

# 01: Общий JSON-контракт позиций и embedded массива

Ticket P01 завершён в первой волне. Разделы ниже сохраняют исторический контракт
и результаты на тот момент. Текущий экспорт — Position/PositionRequest; прежние имена
и ограничения набора полей заменены [универсальной формой](../../universal-positions/spec.md).
Актуальные требования — [spec positions](../spec.md). P03 отменён и не зависимость документов.

## What to build

Разработчик решения использует самостоятельные публичные Core/Item/Commercial
и соответствующие request types общего package positions; читает согласованные поля
и сериализует позиции внутри request документа с различием omission и пустого массива.
Это проверяемый JSON-срез общей библиотеки, а не набор деклараций для будущего supply.

Контракт P01 утверждён пользователем: [spec](../spec.md). Ticket завершён;
Movement P03 и entity endpoints в эту волну не входят.

## Acceptance criteria

- [x] Самостоятельный package positions не принадлежит supply или другой сущности.
  Core содержит Meta model.Meta, ID/AccountID string. Item встраивает Core и содержит
  Assortment model.Reference и Quantity float64. Commercial встраивает Item и содержит
  Price/Discount float64, VAT int, VATEnabled bool; JSON плоский.
- [x] CoreRequest содержит только ID model.Field[string] с omitzero; ItemRequest
  встраивает CoreRequest и добавляет Assortment Field[model.Reference], Quantity Field[float64].
  CommercialRequest встраивает ItemRequest и добавляет Price/Discount Field[float64],
  VAT Field[int], VATEnabled Field[bool], все с omitzero.
- [x] Meta и AccountID отсутствуют в request types. ID — matching selector существующей
  строки embedded Update, а не client-generated ID новой строки. Sparse existing-row
  request кодируется без предварительного чтения или достраивания полей.
- [x] Quantity и money сохраняют float64, деньги — в копейках; отрицательный discount,
  VAT=0 и VATEnabled=false не теряются. Nullable Field сохраняет сериализацию null;
  его допустимость для конкретного поля не обещается.
- [x] Общие формы не добавляют overhead, cost, shipped, inTransit, reserve, slots,
  pack/serial/tracking/GTD или return business rules. Новые document-owned типы могут
  расширять согласованный base; CommercialRequest не объявляется write-формой возврата.
- [x] Оба codecs проверяют значения каждого read поля, ссылки, zero values и неизвестное
  поле; pure JSON чтение не вводит валидацию обязательности или непустого ID.
- [x] Независимые golden JSON проверяют все writable поля, минимальную новую строку,
  sparse quantity-only Update с ID, omission, явные 0/false и null. Read-модель не
  сериализуется вместо request, не добавляет readonly поля.
- [x] Публичный container с Positions model.Field[[]CommercialRequest] и omitzero
  обоими codecs различает unset (нет поля), Value(empty slice) ([]), Value(nil slice) ([]),
  непустой массив и Null (null). Response positions-object не смешан с write-array.
- [x] Документация общей библиотеки объясняет full replacement, matching id, удаление
  отсутствующих старых строк и embedded limit 1000; HTTP document tickets проверяют
  1000/1001 и запрет скрытого merge. JSON-модель не создаёт новый validation/transport API.
- [x] Тесты проверяют только публичный JSON-контракт, по образцу таблицы codecs UOM и
  общего AddressRequest; private fields и helpers не являются тестовыми seams.
- [x] Существующие Field, Reference, Collection и transport не перерабатываются;
  новые HTTP-операции или типы всех будущих документов не объявлены реализованными.
- [x] После реализации пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check и отдельные Standards/Spec review относительно зафиксированного BASE.

## Blocked by

Нет. Нужные model.Field, model.Reference, model.Meta и JSON codecs уже существуют.
DocumentRate, document-common, supply, новые справочники и reader 02 этот JSON-срез не блокируют.

## Testing Decisions

согласованный seam — encoding/json и encoding/json/v2 через public types и небольшой
consumer container. Expected JSON задан вручную/независимыми golden, а не маршалингом
того же SDK. Prior art: request/codecs tests UOM и самостоятельный AddressRequest.
Отдельные HTTP document tests в этом ticket не подменяются фиктивным Service.

## Границы

После утверждения — один основной implementer, полный согласованный срез и проверки
по правилам проекта, без учебных остановок. Ticket утверждён пользователем; разрешена первая волна. HTTP reader 02,
Movement 03, document endpoints, expand, live API, .env, push и commits подготовки
в этот ticket не входят. Новые решения контракта требуют отдельного согласования.

## Sources

- [Первичные факты и ссылки](../research.md).
- [Общие правила embedded positions](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1031-L1037).

## Comments

## Утверждение и ownership первой волны

Пользователь 2026-11-01 утвердил первую волну после исправления денежных типов и смысла nil курса.
JSON и HTTP seams этого ticket согласованы. Локальная реализация через TDD и commit разрешены;
независимые Standards/Spec review координирует исполнитель интеграции относительно planning BASE.
Не менять coverage.md/README и общие карты: документацию интегратор сведёт отдельным коммитом.
Файлы и трёхпоточный запуск закреплены в [handoff](../../document-workstream/wave-1.md).

## Answer

Реализованы самостоятельные `positions.Core`, `Item`, `Commercial` и отдельные
request-формы. Встраивание даёт плоский JSON; request содержит только writable поля
и matching ID существующей строки. GoDoc фиксирует full replacement, отсутствие
GET/merge, удаление пропущенных старых строк и embedded limit 1000.

Тесты через `encoding/json` и `encoding/json/v2` проверяют все read-поля, ссылки,
дробные числа, отрицательный discount, price 0/0.0, VAT=0 с true/false, отсутствие
обязательного ID и неизвестные поля. Независимые literal golden проверяют request
omission, новую строку, ID-only/quantity-only Update, все поля, 0/false и null.
Consumer container с `Positions model.Field[[]CommercialRequest]` различает
unset/null/nil slice/empty slice/array. Его тест сразу прошёл на существующем Field;
реализация Field не менялась.

TDD evidence 2026-11-01: первый запуск до реализации Core:

```text
$ go test -count=1 ./positions -run '^TestCoreJSON$'
positions/model_test.go:32:20: undefined: positions.Core
positions/model_test.go:40:22: undefined: positions.Core
FAIL github.com/eteeet/moysklad-sdk-go/positions [build failed]
```

После минимальной реализации тот же тест прошёл (exit 0). Отдельные red → green
циклы затем выполнены для Item, Commercial и трёх request-форм: каждый red был
ошибкой отсутствующего соответствующего public type до добавления реализации.

Проверки прошли (exit 0):

- `gofmt -w positions/doc.go positions/model.go positions/request.go positions/model_test.go positions/request_test.go`;
- `go test -race -count=1 ./positions`;
- `go test -race -count=1 ./...`;
- `go vet ./... ./testdata`;
- `git diff --check`.

Status оставлен claimed до независимых Standards/Spec review интегратором от planning
BASE `a2804283e7aafe3c6ee41c17a4c2c79f7e01f999`. Review ещё не проводились. Movement,
HTTP reader/writer, entity, expand и business validation этим изменением не реализованы.

### Итог интеграции

Включено в dev, commit реализации: 20e00ad.
Независимые Standards: 0 замечаний; Spec: 0 замечаний. Проверен полный срез
от planning BASE a2804283e7aafe3c6ee41c17a4c2c79f7e01f999 до HEAD 1b28603,
включая централизованные docs. [Отчёты](../../document-workstream/review-wave-1.md).
Общий go test -race -count=1 ./..., go vet ./... ./testdata, gofmt и diff-check пройдены.
Публичная документация сведена интегратором после слияния; live API и push не выполнялись.
