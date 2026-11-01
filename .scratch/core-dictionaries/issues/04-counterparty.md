Status: claimed
Blocked by: 03

# 04: Counterparty — Get, Create, Update

## What to build

Разработчик решения получает, создаёт и частично изменяет Counterparty на согласованном
наборе основных полей, контактов, реквизитов РФ, адресов, attributes и ссылок.
Общие адресные типы и CompanyType уже предоставляет ticket 03; этот ticket подключает их
к законченному HTTP-срезу с тестами и документацией.

Контракт: [общая spec](../spec.md). Зависимость: [общие типы реквизитов](03-shared-requisites.md).
Образцы методов — Product/UOM. Существующий пустой пакет counterparty заполняется без дублирования.

## Acceptance criteria

- [ ] Ticket 03 завершён (Status: resolved), общие типы доступны и их проверки проходят.
  Контракты Address/AddressRequest/CompanyType не переопределяются и не копируются в counterparty.
- [ ] Counterparty читает meta, id, accountId, created, updated и основные поля
  name, description, code, externalCode, archived, shared, companyType, email, phone, fax.
- [ ] Поддержаны legalTitle, legalFirstName, legalLastName, legalMiddleName, inn, kpp, ogrn, ogrnip, okpo,
  certificateNumber/Date, birthDate, sex. Идентификаторы/реквизиты — строки,
  даты — model.DateTime, companyType — model.CompanyType; неприменимые/отсутствующие поля допустимы.
- [ ] Поддержаны actualAddress/legalAddress и *model.Address для двух полей *AddressFull;
  attributes ([]model.Attribute), tags ([]string), discountCardNumber;
  owner/group/bonusProgram/state (*model.Reference), priceType (*model.PriceType).
- [ ] Читаются syncId (string), bonusPoints (int), salesAmount (float64).
  Для salesAmount обе JSON-формы 0 и 0.0 декодируются без пересчёта единиц.
  Финансовая арифметика или точная десятичная модель не добавляются.
- [ ] CreateRequest имеет обязательный Name string; остальные согласованные записываемые поля —
  Field с omitzero. SyncID разрешён только в CreateRequest, отсутствует в UpdateRequest.
  В UpdateRequest все изменяемые поля — Field; metadata, ID, даты создания/обновления,
  bonusPoints и salesAmount отсутствуют в обоих requests.
- [ ] В requests CompanyType — Field[model.CompanyType], адреса — Field[model.AddressRequest],
  ссылки — Field[model.Reference], priceType — Field[model.PriceType], коллекции — Field от слайса.
  Неизвестные CompanyType передаются без валидации; defaults не подставляются.
- [ ] Service/NewService соответствуют образцу: существующий *transport.Client, отказ nil client.
  Get — GET entity/counterparty/{id}; Create — POST entity/counterparty;
  Update — PUT entity/counterparty/{id}; результат (*Counterparty, error) из ответа сервера.
- [ ] UUID проверяется по форме без ограничения версии; Create отклоняет только Name=="";
  полностью незаданный Update отклоняется. false/""/[]/null и явно заданный пустой адрес
  считаются заданными. Новой предметной проверки реквизитов/комбинаций полей нет.
- [ ] SDK не выполняет Get перед записью, не следует по ссылкам, не сливает адреса/attributes/tags,
  не конкатенирует строковый адрес, не меняет null на "" и не добавляет retry.
  Одновременная передача двух форм адреса сохраняется в JSON без локального выбора.
- [ ] Модель проверена по всем поддержанным значениям неизменённых upstream Counterparty и
  CounterpartyLegal; дополнительные небольшие случаи покрывают поля, которых в них нет.
  Проверены optional/null/неизвестные поля, даты, неизвестный CompanyType и сложные ссылки.
  Несогласованность companyType и реквизитов upstream не становится бизнес-правилом SDK.
- [ ] JSON-тесты обоими кодеками покрывают минимальный/полный Create, полный/частичный Update,
  независимые golden, пропуски, false/true, 0 для подходящих значений, пустые строки/массивы/null,
  вложенные поля адреса, оба представления адреса и Create-only syncId.
- [ ] HTTP-тесты с настоящим transport проверяют метод/маршрут, отсутствие query/тела GET/лишних
  запросов, заголовки, plain/gzip, точное тело и серверные данные вместо эха input.
- [ ] Все методы проверены для context, auth/transport/APIError 4xx/5xx/429,
  errors.As/errors.Is, контекста без двойного moysklad:, неожиданного статуса,
  malformed/null JSON, пустого/отсутствующего ID. Успех — HTTP 200 и непустой ID;
  любая ошибка возвращает nil-модель, автоматических повторов нет.
- [ ] Подтверждены границы: accounts/contactpersons/notes/files/discounts, региональные объекты
  реквизитов, их metadata и операции не реализованы и не включены в requests.
  attributes использует существующий общий контракт; для tags SDK не обещает непроверенную
  очистку или серверное слияние коллекции.
- [ ] Документация перечисляет фактически готовые поля/методы и TODO; правило синхронизации
  do/validID соблюдено для обнаруженных общих исправлений без общего рефакторинга.
- [ ] Пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check и встроенное ревью Standards/Spec с зафиксированным исходным BASE.

## Blocked by

03 — общие Address, AddressRequest и CompanyType.
Country, ExpenseItem и Organization не являются блокирующими зависимостями.

## Выполнение и границы

- Один основной агент, одна новая сессия, полный implement/TDD без учебных остановок.
- До правок прочитать инструкции, всю spec и dependency; проверить Git, записать исходный HEAD
  как BASE ревью. Сохранить исходные изменения планирования и чужие изменения;
  не включать их в commit под видом реализации. Свои изменения документации выделить отдельно.
- Reviewer сообщает замечания, исправляет основной агент. Новые контракты требуют согласования.
- Не добавлять Organization, счета, подресурсы, List/Delete/expand/batch, транспортные механизмы,
  фасад, .env или живой API. Не изменять Product/UOM вне подтверждённого общего исправления.
- После полного ticket — локальный commit, Status: resolved и ## Answer с SHA/проверками/ревью,
  отчёт пользователю и остановка. Push/PR/VPS требуют отдельной команды.

## Sources

[Официальная документация Counterparty, зафиксированный revision](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_counterparty.md).
Общий контракт attributes и адресных ограничений приведён в spec.
