Status: resolved
Blocked by: 03

# 05: Organization — Get, Create, Update

## What to build

Разработчик решения получает, создаёт и частично изменяет Organization на согласованном
наборе основных полей, контактов, реквизитов РФ, адресов, attributes и ссылок.
Ticket использует общие типы из 03 и не требует реализации Counterparty.
Результат — законченный HTTP-срез с моделями, тремя методами, тестами и документацией.

Контракт: [общая spec](../spec.md). Зависимость: [общие типы реквизитов](03-shared-requisites.md).
Образцы методов — Product/UOM; Counterparty при наличии можно читать как пример, но не как новый контракт.

## Acceptance criteria

- [x] Ticket 03 завершён (Status: resolved), общие типы доступны и проверены.
  Address/AddressRequest/CompanyType не дублируются и не меняются без согласования.
- [x] Organization читает meta, id, accountId, created, updated; name, description, code,
  externalCode, archived, shared, companyType, email, phone, fax.
- [x] Поддержаны legalTitle, legalFirstName, legalLastName, legalMiddleName, inn, kpp, ogrn, ogrnip, okpo,
  certificateNumber/Date; actualAddress/legalAddress и два *model.Address для *AddressFull;
  attributes ([]model.Attribute), owner/group/bonusProgram (*model.Reference).
- [x] Дополнительно поддержаны payerVat (bool), director, directorPosition, chiefAccountant,
  trackingContractNumber (string), trackingContractDate (model.DateTime).
  Читаются syncId (string), bonusPoints (int), isEgaisEnable (bool).
  Все имена Go и JSON соответствуют таблицам spec.
- [x] CompanyType в модели — model.CompanyType; даты — существующий DateTime;
  неизвестные значения и неприменимые/отсутствующие реквизиты не ломают чтение.
  SDK не подставляет legal и не выводит реквизиты из ФИО/companyType самостоятельно.
- [x] CreateRequest: обязательный Name string, остальные согласованные записываемые поля —
  Field с omitzero. SyncID доступен только при Create. В UpdateRequest все изменяемые поля — Field;
  syncId, isEgaisEnable, bonusPoints, metadata, ID и даты создания/обновления в Update не передаются.
  Readonly-полей и isEgaisEnable нет и в CreateRequest.
- [x] Адресные поля записи используют Field[model.AddressRequest], CompanyType —
  Field[model.CompanyType], ссылки — Field[model.Reference], attributes — Field[[]model.Attribute].
  Незаданные поля отсутствуют; false/""/[]/null/пустой объект не теряются.
- [x] Service/NewService используют существующий transport.Client, nil client отклоняется.
  Get — GET entity/organization/{id}, Create — POST entity/organization,
  Update — PUT entity/organization/{id}; возвращаются (*Organization, error) и данные сервера.
- [x] Валидация UUID без ограничения версии, Name=="" без trim и полностью незаданного Update
  соответствует образцу. Ограничения реквизитов/ссылок/комбинаций полей оставлены серверу.
  Нет предварительного Get, переходов по ссылкам, локального merge/нормализации или retry.
- [x] Официальное противоречие isEgaisEnable («всегда false» в таблице, true в примерах)
  отражено в документации и тестах: оба значения читаются фактически, без принудительного false.
  fsrarId/utmUrl из примеров намеренно игнорируются как согласованные TODO и не попадают в requests.
- [x] Значения всех поддержанных полей проверены по неизменённому upstream и небольшим
  дополнительным случаям для отсутствующих полей, обоими JSON-кодеками.
  Проверены ссылки, даты, optional/null/неизвестные поля, unknown CompanyType,
  payerVat=false, сведения о руководителе/бухгалтере и trackingContractDate/Number.
- [x] Requests проверены обоими кодеками: минимальный Create, полные Create/Update
  по независимым golden, одно поле Update, пропуски, false/true, пустые строки/массивы/null,
  вложенные поля адресов, одновременная передача двух форм адреса и Create-only syncId.
- [x] HTTP-тесты с настоящим transport проверяют метод/путь/заголовки, отсутствие query,
  тела GET и лишних запросов, plain/gzip, точное тело и фактический ответ, отличающийся от input.
- [x] Все методы проверены для context, auth/transport/APIError 4xx/5xx/429,
  errors.As/errors.Is, контекста без двойного moysklad:, неожиданного статуса,
  malformed/null JSON и пустого/отсутствующего ID. Успех — HTTP 200 и непустой ID;
  при ошибке — nil-модель, без автоматических повторов.
- [x] Документация точно отражает поддержку и TODO: accounts целиком, региональные реквизиты,
  companyVat__ru, advancePaymentVat, chiefAccountSign, directorSign, stamp, fsrarId, utmUrl;
  поддержанный payerVat не объявляется полной поддержкой настроек НДС.
- [x] Приватные do/validID соответствуют образцу; общие исправления рассматриваются во всех
  применимых копиях по правилу coverage, без нового общего исполнителя.
- [x] Пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check и встроенное ревью Standards/Spec относительно исходного BASE.

## Blocked by

03 — общие Address, AddressRequest и CompanyType.
04 Counterparty не блокирует этот ticket; Country и ExpenseItem также не являются зависимостями.

## Выполнение и границы

- Один основной агент, одна новая сессия, полный implement/TDD без учебных остановок.
- Прочитать инструкции, всю spec и dependency; проверить исходный Git и записать HEAD как BASE.
  Сохранить начальный dirty diff; не включать предшествующие изменения планирования/чужие правки
  в commit реализации. Свои изменения документации отделять от исходных.
- Reviewer не меняет код. Новые контракты согласовать, даже если это останавливает ticket.
- Не добавлять Counterparty, счета/подресурсы, List/Delete/expand/batch, общие транспортные
  механизмы, фасад, .env или живой API. Наличие примера не расширяет согласованный scope.
- После полного ticket — локальный commit, Status: resolved и ## Answer с SHA/проверками/ревью,
  отчёт и остановка. Push/PR/VPS — только по отдельной команде.

## Sources

[Официальная документация Organization, зафиксированный revision](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_organization.md).
Общий контракт attributes и адресных ограничений приведён в spec.

## Answer

Реализованы Organization, отдельные CreateRequest/UpdateRequest и Get/Create/Update,
весь согласованный набор полей, оба JSON-кодека, локальные HTTP-тесты и независимые
golden. Документация: `docs/organization.md`.

- Commit реализации: `4a119b28c5d1cd019b4be4bf58488833ed218c57`; включён в `dev` как `5e799ae`.
- Проверки: gofmt; `go test -race -count=1 ./...`; `go vet ./... ./testdata`;
  `git diff --cached --check` — пройдены, включая общий прогон после интеграции пяти tickets.
- Независимые ревью: Standards — 0 замечаний; Spec — два пробела тестового покрытия
  устранены (9 CompanyType-констант и Attribute value:0 внутри обоих request-типов,
  оба JSON-кодека), повторное Spec/Standards ревью — 0 открытых замечаний.
- Живой API не использовался; исходные изменения планирования не включены в commit.
