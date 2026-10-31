Status: ready-for-agent
Stage: APPROVED
Blocked by: None

# 01: Общий JSON-контракт позиций и embedded массива

## What to build

Разработчик решения использует самостоятельные публичные Core/Item/Commercial
и соответствующие request types общего package positions; читает согласованные поля
и сериализует позиции внутри request документа с различием omission и пустого массива.
Это проверяемый JSON-срез общей библиотеки, а не набор деклараций для будущего supply.

Контракт P01 утверждён пользователем: [spec](../spec.md). Ticket ready-for-agent;
Movement P03 и entity endpoints в эту волну не входят.

## Acceptance criteria

- [ ] Самостоятельный package positions не принадлежит supply или другой сущности.
  Core содержит Meta model.Meta, ID/AccountID string. Item встраивает Core и содержит
  Assortment model.Reference и Quantity float64. Commercial встраивает Item и содержит
  Price/Discount float64, VAT int, VATEnabled bool; JSON плоский.
- [ ] CoreRequest содержит только ID model.Field[string] с omitzero; ItemRequest
  встраивает CoreRequest и добавляет Assortment Field[model.Reference], Quantity Field[float64].
  CommercialRequest встраивает ItemRequest и добавляет Price/Discount Field[float64],
  VAT Field[int], VATEnabled Field[bool], все с omitzero.
- [ ] Meta и AccountID отсутствуют в request types. ID — matching selector существующей
  строки embedded Update, а не client-generated ID новой строки. Sparse existing-row
  request кодируется без предварительного чтения или достраивания полей.
- [ ] Quantity и money сохраняют float64, деньги — в копейках; отрицательный discount,
  VAT=0 и VATEnabled=false не теряются. Nullable Field сохраняет сериализацию null;
  его допустимость для конкретного поля не обещается.
- [ ] Общие формы не добавляют overhead, cost, shipped, inTransit, reserve, slots,
  pack/serial/tracking/GTD или return business rules. Новые document-owned типы могут
  расширять согласованный base; CommercialRequest не объявляется write-формой возврата.
- [ ] Оба codecs проверяют значения каждого read поля, ссылки, zero values и неизвестное
  поле; pure JSON чтение не вводит валидацию обязательности или непустого ID.
- [ ] Независимые golden JSON проверяют все writable поля, минимальную новую строку,
  sparse quantity-only Update с ID, omission, явные 0/false и null. Read-модель не
  сериализуется вместо request, не добавляет readonly поля.
- [ ] Публичный container с Positions model.Field[[]CommercialRequest] и omitzero
  обоими codecs различает unset (нет поля), Value(empty slice) ([]), Value(nil slice) ([]),
  непустой массив и Null (null). Response positions-object не смешан с write-array.
- [ ] Документация общей библиотеки объясняет full replacement, matching id, удаление
  отсутствующих старых строк и embedded limit 1000; HTTP document tickets проверяют
  1000/1001 и запрет скрытого merge. JSON-модель не создаёт новый validation/transport API.
- [ ] Тесты проверяют только публичный JSON-контракт, по образцу таблицы codecs UOM и
  общего AddressRequest; private fields и helpers не являются тестовыми seams.
- [ ] Существующие Field, Reference, Collection и transport не перерабатываются;
  новые HTTP-операции или типы всех будущих документов не объявлены реализованными.
- [ ] После реализации пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata,
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
