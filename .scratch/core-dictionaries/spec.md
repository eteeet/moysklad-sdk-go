Status: resolved
Причина: все тикеты спеки resolved (5/5). Закрыто 2026-11-14 при сверке здоровья документов.

# Country, ExpenseItem, Counterparty, Organization: Get, Create, Update

## Problem Statement

Разработчик решения уже может получать, создавать и изменять Product и UOM через SDK,
но следующие справочники пока не имеют согласованной реализации.
Нужно расширить поддержку на Country, ExpenseItem, Counterparty и Organization,
сохранив утверждённый образец, проверяемое публичное поведение и явные границы покрытия.

Это следующий срез SDK, а не реализация всего JSON API или всех возможностей этих сущностей.
Service намеренно отложен до завершения ещё нескольких сущностей.

## Solution

Добавить четыре пакета сущностей с моделями чтения, отдельными моделями CreateRequest
и UpdateRequest и обёртками Service по образцу Product/UOM.
Country и ExpenseItem поддержать на полном наборе полей из их таблиц атрибутов.
Counterparty и Organization поддержать на согласованном наборе основных полей,
контактов, реквизитов РФ, адресов, attributes и ссылок.

Для Counterparty и Organization добавить общие Address, AddressRequest и CompanyType
отдельным ticket с самостоятельными JSON-тестами и ревью.
Не менять транспорт, архитектуру Product/UOM или контракт существующих общих типов.

## User Stories

1. Как разработчик решения, я хочу получить Country по ID, чтобы прочитать страну.
2. Как разработчик решения, я хочу читать предустановленные и пользовательские Country, чтобы не выбирать разные модели по их происхождению.
3. Как разработчик решения, я хочу создать Country с одним name, чтобы сервер заполнил остальные значения.
4. Как разработчик решения, я хочу изменить отдельные поля Country, чтобы не отправлять остальные данные повторно.
5. Как разработчик решения, я хочу получить ошибку API при недопустимом изменении предустановленной Country, чтобы обработать ограничение сервера.
6. Как разработчик решения, я хочу получить ExpenseItem по ID, чтобы прочитать статью расходов.
7. Как разработчик решения, я хочу читать archived и operatingExpenses, чтобы учитывать назначение статьи расходов.
8. Как разработчик решения, я хочу создать ExpenseItem без явных необязательных значений, чтобы применялись defaults сервера.
9. Как разработчик решения, я хочу явно передать operatingExpenses=false, чтобы SDK не потерял это значение.
10. Как разработчик решения, я хочу частично обновить ExpenseItem, чтобы изменить только выбранные поля.
11. Как разработчик решения, я хочу получить Counterparty, чтобы прочитать его контакты, реквизиты и связанные справочники.
12. Как разработчик решения, я хочу создать Counterparty с минимальным обязательным набором, чтобы не подставлять defaults за сервер.
13. Как разработчик решения, я хочу передать реквизиты юрлица, ИП или физического лица, чтобы описать соответствующего Counterparty.
14. Как разработчик решения, я хочу использовать именованный строковый CompanyType с открытым набором значений, чтобы известные значения имели константы, а новые не ломали декодирование.
15. Как разработчик решения, я хочу прочитать salesAmount из JSON с числом 0 или 0.0, чтобы оба документальных представления обрабатывались одинаково.
16. Как разработчик решения, я хочу читать bonusPoints и ссылки на bonusProgram, state и priceType, чтобы использовать согласованные данные Counterparty.
17. Как разработчик решения, я хочу передавать tags и discountCardNumber, чтобы менять соответствующие поля Counterparty.
18. Как разработчик решения, я хочу получить Organization, чтобы прочитать данные собственного юрлица аккаунта.
19. Как разработчик решения, я хочу создать и изменить Organization, чтобы передать основные данные, контакты и реквизиты.
20. Как разработчик решения, я хочу передать payerVat, сведения о руководителе и бухгалтере, чтобы заполнить согласованные поля Organization.
21. Как разработчик решения, я хочу читать и записывать дату и номер договора с ЦРПТ, чтобы не терять эти поля Organization.
22. Как разработчик решения, я хочу читать и задавать строковый адрес, чтобы использовать простое представление адреса.
23. Как разработчик решения, я хочу читать и задавать структурированный адрес, чтобы управлять отдельными адресными полями.
24. Как разработчик решения, я хочу различать незаданное вложенное поле адреса и явную пустую строку, чтобы сформировать точное тело запроса.
25. Как разработчик решения, я хочу знать приоритет структурированного адреса и ограничения null, чтобы не считать null командой очистки адреса.
26. Как разработчик решения, я хочу передавать attributes существующими общими типами, чтобы использовать значения и явный сброс без нового API SDK.
27. Как разработчик решения, я хочу передавать owner и group как Reference, чтобы не реализовывать связанные сущности заранее.
28. Как разработчик решения, я хочу указать syncId при создании Counterparty или Organization, чтобы задать идентификатор синхронизации один раз.
29. Как разработчик решения, я хочу различать отсутствие, null и значение поля запроса, чтобы не терять false, 0, пустую строку или массив.
30. Как разработчик решения, я хочу получить фактический объект ответа Create/Update, чтобы не принимать входной request за сохранённое состояние.
31. Как разработчик решения, я хочу отменять операции через context, чтобы управлять временем жизни запроса.
32. Как разработчик решения, я хочу распознавать ошибки через errors.As/errors.Is, чтобы обрабатывать ошибки API и транспорта программно.
33. Как разработчик решения, я хочу получить ошибку для неверного ID, пустого Create.Name и полностью незаданного Update, чтобы не выполнять такие запросы.
34. Как разработчик решения, я хочу видеть точные границы поддержки, чтобы не принять наличие модели за поддержку счетов, expand или всего API сущности.
35. Как сопровождающий SDK, я хочу выполнять каждый согласованный ticket целиком в новой сессии, чтобы проверить воспроизводимость передачи контекста без учебных остановок.
36. Как разработчик решения, я хочу использовать общие адресные типы и CompanyType независимо от сущностей, чтобы их публичный JSON-контракт был проверен отдельно.

## Implementation Decisions

### Публичный интерфейс и маршруты

| Пакет | Модель | Обёртка | Конструктор | Ресурс относительно API 1.2 |
| --- | --- | --- | --- | --- |
| country | Country | Service | NewService | entity/country |
| expenseitem | ExpenseItem | Service | NewService | entity/expenseitem |
| counterparty | Counterparty | Service | NewService | entity/counterparty |
| organization | Organization | Service | NewService | entity/organization |

- В каждом пакете отдельные CreateRequest и UpdateRequest.
- NewService принимает существующий *transport.Client, возвращает *Service и error;
  nil client отклоняется по образцу Product/UOM.
- Get принимает context и ID, выполняет GET ресурса с суффиксом /{id} без тела и дополнительных query.
- Create принимает context и CreateRequest, выполняет POST на основной ресурс.
- Update принимает context, ID и UpdateRequest, выполняет PUT ресурса с суффиксом /{id}.
- Методы возвращают указатель на модель соответствующей сущности и error.
- Существующий пустой пакет counterparty заполняется, а не дублируется новым пакетом.

### Country: полный набор

| JSON-поле | Поле Go | Тип чтения | Запись |
| --- | --- | --- | --- |
| meta | Meta | model.Meta | нет |
| id | ID | string | нет |
| accountId | AccountID | string | нет |
| updated | Updated | model.DateTime | нет |
| name | Name | string | Create, Update |
| description | Description | string | Create, Update |
| code | Code | string | Create, Update |
| externalCode | ExternalCode | string | Create, Update |
| owner | Owner | *model.Reference | Create, Update |
| group | Group | *model.Reference | Create, Update |
| shared | Shared | bool | Create, Update |

- Все 11 полей чтения и 7 полей записи включены; code сохраняется как строка,
  без числового преобразования и потери ведущих нулей.
- Предустановленные Country доступны для чтения; документация запрещает их изменение.
  SDK не классифицирует Country локально и передаёт ошибку сервера.
- Таблица помечает accountId как обязательное поле ответа, но официальный пример
  предустановленной Country «Марокко» не содержит accountId, owner, group и shared.
  Эти поля не проверяются на обязательность; отсутствие даёт нулевые значения.
- Shared=false не служит признаком происхождения Country.

### ExpenseItem: полный набор

- Все поля и типы Country присутствуют также у ExpenseItem с теми же именами Go и режимами записи.
- Дополнительно присутствуют следующие два поля, доступные для Create и Update:

| JSON-поле | Поле Go | Тип чтения | Default при создании по документации |
| --- | --- | --- | --- |
| archived | Archived | bool | false |
| operatingExpenses | OperatingExpenses | bool | true |

- Итого 13 полей чтения и 9 полей записи. Defaults применяет сервер, не SDK.
- Отсутствующее operatingExpenses при чтении даёт нулевой bool, а не подставленное true:
  default создания не является правилом декодирования ответа.
- Таблица помечает accountId, а также shared обязательными при ответе,
  но официальные примеры предустановленных ExpenseItem не содержат
  accountId, owner, group и shared; отсутствие допускается.
- В разделе ExpenseItem не найден явный запрет изменения предустановленных статей.
  Запрет Country/UOM не переносится сюда как установленный контракт.
  SDK не классифицирует происхождение ExpenseItem и передаёт ошибки API.

### Counterparty и Organization: общие поля

Следующая таблица полностью определяет общий согласованный набор, а не список примеров.

| JSON-поля | Поля Go в том же порядке | Тип чтения | Запись |
| --- | --- | --- | --- |
| meta | Meta | model.Meta | нет |
| id, accountId | ID, AccountID | string | нет |
| created, updated | Created, Updated | model.DateTime | нет |
| name, description, code, externalCode | Name, Description, Code, ExternalCode | string | Create, Update |
| archived, shared | Archived, Shared | bool | Create, Update |
| companyType | CompanyType | model.CompanyType | Create, Update |
| email, phone, fax | Email, Phone, Fax | string | Create, Update |
| legalTitle, legalFirstName, legalLastName, legalMiddleName | LegalTitle, LegalFirstName, LegalLastName, LegalMiddleName | string | Create, Update |
| inn, kpp, ogrn, ogrnip, okpo | INN, KPP, OGRN, OGRNIP, OKPO | string | Create, Update |
| certificateNumber | CertificateNumber | string | Create, Update |
| certificateDate | CertificateDate | model.DateTime | Create, Update |
| actualAddress, legalAddress | ActualAddress, LegalAddress | string | Create, Update |
| actualAddressFull, legalAddressFull | ActualAddressFull, LegalAddressFull | *model.Address | Create, Update через model.AddressRequest |
| attributes | Attributes | []model.Attribute | Create, Update |
| owner, group, bonusProgram | Owner, Group, BonusProgram | *model.Reference | Create, Update |
| syncId | SyncID | string | только Create |
| bonusPoints | BonusPoints | int | нет |

- Реквизиты зависят от companyType; отсутствие неприменимых полей нормально.
  SDK не выводит одни реквизиты из других и не очищает неприменимые поля самостоятельно.
- companyType использует общий именованный строковый model.CompanyType с константами.
  Набор открыт: известные значения не ограничивают декодер и не валидируются локально.
- Региональные значения companyType сохраняются, но вложенные региональные реквизиты
  mod__requisites__uz и mod__requisites__kz в этот срез не входят.
- companyType обязателен в ответе, но не помечен необходимым при создании.
  Примеры без этого поля возвращают legal; нормативный default не установлен.
  SDK не подставляет companyType за сервер.
- syncId после заполнения менять нельзя. В согласованном срезе доступен при чтении
  и в CreateRequest, отсутствует в UpdateRequest; SDK не выполняет предварительный Get.
- Поля только для чтения не включаются в модели запросов.

### Counterparty: дополнительные поля

| JSON-поле | Поле Go | Тип чтения | Запись |
| --- | --- | --- | --- |
| birthDate | BirthDate | model.DateTime | Create, Update |
| sex | Sex | string | Create, Update |
| discountCardNumber | DiscountCardNumber | string | Create, Update |
| tags | Tags | []string | Create, Update |
| state | State | *model.Reference | Create, Update |
| priceType | PriceType | *model.PriceType | Create, Update |
| salesAmount | SalesAmount | float64 | нет |

- PriceType использует существующий model.PriceType, без управления типами цен.
- sex передаётся строкой без новой локальной предметной валидации.
- Для salesAmount таблица указывает Int, но официальный ответ Get и upstream-пример
  содержат 0.0. Согласован float64, как у существующих цен; тестируются формы 0 и 0.0.
  Это не обещание точной десятичной арифметики или допустимости дробной суммы на сервере.
  Значение читается без пересчёта единиц и без финансовых вычислений.
- tags передаётся ровно как заданный массив. SDK не сливает теги и не обещает
  недокументированного поведения сервера для [], null или замены коллекции.

### Organization: дополнительные поля

| JSON-поле | Поле Go | Тип чтения | Запись |
| --- | --- | --- | --- |
| payerVat | PayerVAT | bool | Create, Update |
| director | Director | string | Create, Update |
| directorPosition | DirectorPosition | string | Create, Update |
| chiefAccountant | ChiefAccountant | string | Create, Update |
| trackingContractNumber | TrackingContractNumber | string | Create, Update |
| trackingContractDate | TrackingContractDate | model.DateTime | Create, Update |
| isEgaisEnable | IsEGAISEnable | bool | нет в этом срезе |

- В документации есть противоречие: таблица описывает isEgaisEnable как всегда false,
  но официальные примеры содержат true. Читаем фактический bool без принудительного false;
  отсутствие поля даёт нулевое значение. Возможность записи в этот срез не входит.
- payerVat поддержан независимо от отложенных companyVat__ru и advancePaymentVat.
  Поддержка флага не означает поддержку всех настроек НДС Organization.

### Общий CompanyType

- model.CompanyType — именованный тип поверх string, не alias и не закрытый enum.
  Специальная проверка значений и кастомный JSON-кодек не нужны.
- Известные значения представлены типизированными константами:

| Константа model | Значение JSON |
| --- | --- |
| CompanyTypeLegal | legal |
| CompanyTypeEntrepreneur | entrepreneur |
| CompanyTypeIndividual | individual |
| CompanyTypeLegalUZ | legalUZ |
| CompanyTypeEntrepreneurUZ | entrepreneurUZ |
| CompanyTypeIndividualUZ | individualUZ |
| CompanyTypeLegalKZ | legalKZ |
| CompanyTypeEntrepreneurKZ | entrepreneurKZ |
| CompanyTypeIndividualKZ | individualKZ |

- Неизвестные строковые значения сохраняются при чтении и записи обоими JSON-кодеками.
  Нулевое значение — пустая строка; значение по умолчанию не подставляется.
- В обоих response-типах CompanyType имеет тип model.CompanyType;
  в обоих CreateRequest/UpdateRequest — model.Field[model.CompanyType] с omitzero.
- Незаданное поле не отправляется; явно заданная пустая строка и null кодируются без потери.
  Допустимость конкретного значения проверяет сервер.
- MetaType остаётся без изменений. Этот ticket не расширяет константами другие типы SDK.

### Общие Address и AddressRequest

Новый общий тип model.Address содержит только согласованные адресные поля:

| JSON-поле | Поле Go | Тип чтения | Тип поля в AddressRequest |
| --- | --- | --- | --- |
| addInfo | AddInfo | string | model.Field[string] |
| apartment | Apartment | string | model.Field[string] |
| city | City | string | model.Field[string] |
| comment | Comment | string | model.Field[string] |
| fiasCode__ru | FIASCodeRU | string | model.Field[string] |
| country | Country | *model.Reference | model.Field[model.Reference] |
| house | House | string | model.Field[string] |
| postalCode | PostalCode | string | model.Field[string] |
| region | Region | *model.Reference | model.Field[model.Reference] |
| street | Street | string | model.Field[string] |

- Address — модель чтения; отдельный AddressRequest описывает только передаваемые поля.
  В CreateRequest/UpdateRequest структурированные адреса имеют внешний
  model.Field[model.AddressRequest], а не Field от модели чтения.
- В AddressRequest все поля используют omitzero; незаданное вложенное поле отсутствует,
  Value("") передаёт пустую строку, Null кодируется явно. Вложенные ссылки тоже используют Field.
- Отсутствующие и null-адреса ответа декодируются в nil; неизвестные адресные поля игнорируются.
- SDK не конкатенирует адрес, не копирует строку в addInfo и не выбирает одну форму за вызывающего.
  Если заданы обе формы, обе уходят в запрос; по документации сервер отдаёт приоритет структуре.
- При передаче только строки сервер отражает её также в addInfo структурированного адреса.
  Это серверное поведение, а не локальная нормализация SDK.
- null не очищает адрес. Для очистки используется пустая строка соответствующего
  строкового адреса; entity-разделы прямо приводят actualAddress, общий раздел null
  говорит о соответствующем строковом поле. SDK не заменяет null на "".
- Для отдельных частей адреса не объявляется серверное слияние или полная замена:
  SDK лишь формирует точное тело без локального merge и предварительного Get.
- Явно заданный пустой AddressRequest считается заданным внешним полем и кодируется как {}.
  SDK не вводит отдельную предметную валидацию адреса или обещание принятия {} сервером.
- Country и Region внутри адреса — ссылки. Реализация SDK для Region не требуется.

### Общие правила моделей чтения и записи

- Использовать существующие Meta, Reference, DateTime, Field, Attribute и PriceType,
  а также согласованные новые Address, AddressRequest и CompanyType.
  Указатели ссылок/вложенных объектов и DateTime оформлять по образцу Product/UOM с omitzero.
- Отсутствующие поля чтения дают нулевые значения; новые/неподдержанные JSON-поля
  не мешают чтению. Нового обязательного контроля полноты ответов, кроме ID, не вводить.
- Не менять глобальные правила DateTime и Attribute; null/отсутствие даты дают нулевую дату
  по существующему контракту. Идентификаторы и реквизиты хранятся строками.
- В каждом CreateRequest Name — обязательный string; все остальные согласованные
  записываемые поля используют model.Field соответствующего типа.
- В каждом UpdateRequest все согласованные изменяемые поля используют model.Field.
- Внешние поля attributes и tags — Field от соответствующего слайса;
  ссылки — Field[model.Reference], priceType — Field[model.PriceType].
- Field использует omitzero, различает отсутствие, null и значение;
  false, 0, пустая строка и [] не теряются. Существующее кодирование nil-слайса как [] сохраняется.
- Кодирование null не обещает принятия или очистки поля сервером.
  Документированные адресные ограничения сохраняются; новой локальной null-валидации не вводить.
- attributes использует существующие model.Attribute и helpers. Сервер обновляет значения
  по метаданным поля; [] не означает очистку всей коллекции. Сброс отдельного значения — value:null,
  файлового значения — file:null. SDK не сливает массив локально.
- Ни в один request не включать отложенные поля и поля только для чтения.

### Валидация, ответы, ошибки и архитектура

- Get/Update проверяют форму UUID по образцу Product/UOM, не ограничивая версию UUID.
- Create отклоняет только Name==""; без trim, нормализации или новой проверки длин.
- Update отклоняет полностью незаданный request. Явные false, пустые строки, [], null
  и заданный пустой объект считаются заданными полями.
- Допустимость реквизитов, ссылок, companyType, syncId и сочетаний значений проверяет API;
  SDK не подставляет defaults и не добавляет новые предметные проверки.
- Успех каждого метода — HTTP 200 и непустой ID декодированного объекта.
  Возвращаются данные ответа сервера, а не входной request.
- Не вводить дополнительную проверку совпадения ID ответа и ID маршрута:
  сохраняется утверждённое поведение Product/UOM. Несогласованные ID в некоторых
  официальных примерах не интерпретируются как отдельное правило API.
- При ошибке возвращать nil, не частично декодированную модель.
- Сохранить контекст операции и ошибки через errors.As/errors.Is без повторного moysklad:.
- Сохранить существующие context, авторизацию, заголовки и gzip через transport.Client.
- Не добавлять автоматические повторы, предварительные Get, follow-up-запросы по ссылкам,
  локальное слияние или новое управление 429.
- Устройство Product/UOM сохраняется. Выполнение запроса и validID остаются приватными
  для пакета сущности; текущие копии допустимы. Общий исполнитель, генератор и рефакторинг
  дублирования не являются предварительным условием этих четырёх сущностей.
- Общие исправления do и validID переносить во все применимые копии и проверять тестами;
  обоснованные различия сущностей сохранять. Постоянное напоминание добавлено в
  [известные TODO покрытия](../../docs/coverage.md#известные-todo).
- Имена Service/NewService сохраняются. Конфликт имён доменной сущности Service отложен.

## Testing Decisions

- Пользователь согласовал существующие точки проверки Product/UOM:
  публичный JSON-контракт моделей/requests и локальный HTTP-сервер с настоящим transport.Client.
  Новые transport interfaces, моки приватных функций или общая тестовая инфраструктура не нужны.
- Проверять наблюдаемое поведение, не устройство приватных do/validID или способ проверки пустого Update.
- Для каждой сущности проверять значения всех согласованных полей, а не только успешное
  декодирование JSON. Неизменённый upstream дополнять небольшими независимыми случаями
  для отсутствующих в нём полей; отложенные поля из upstream не объявлять поддержанными.
- Основные upstream-примеры уже доступны: Country, ExpenseItem, Counterparty,
  CounterpartyLegal и Organization. Их происхождение и версия зафиксированы механизмом SDK:
  официальный openapi repository, commit 77a778f41da3c0d7b82de9c2204e11b61a32d563.
  Upstream не редактировать ради тестов.
- Для Country взять небольшой предустановленный пример из официальной документации;
  для ExpenseItem — небольшой пример с operatingExpenses=false, например «Возврат».
  Сохранить источник и revision; не выдавать синтетический пример за официальный.
- ExpenseItem upstream не содержит operatingExpenses, owner, group и shared.
  Дополнительно проверить operatingExpenses=true/false, пропуск поля без default декодера,
  ссылки и явный shared=false.
- Для Country/ExpenseItem проверить отсутствие accountId/owner/group/shared у предустановленных
  примеров и сохранение строковых code, включая самостоятельный случай с ведущими нулями.
- Для Counterparty/Organization проверить применимые реквизиты и отсутствие/null неприменимых,
  даты в принятых форматах, неизвестный companyType, ссылки, attributes и игнорирование TODO-полей.
- Для Counterparty проверить обе JSON-формы salesAmount 0 и 0.0, даты birthDate/certificateDate,
  sex, tags, priceType и bonusPoints. Не превращать противоречивый набор реквизитов upstream
  в требование корректности бизнес-правил сервера.
- Для Organization проверить payerVat=false, сведения о руководителе/бухгалтере,
  trackingContractDate/Number и чтение обоих значений isEgaisEnable без наличия этого поля
  в requests; не воспроизводить ошибочное правило «всегда false» из таблицы документации.
- CompanyType проверить отдельно и внутри requests сущностей: девять типизированных
  констант с точными значениями, сохранение неизвестной строки, отсутствие локальной
  валидации/default, работа с Field для незаданного/пустого/null-значения, оба JSON-кодека.
- Address/AddressRequest проверить через публичную сериализацию отдельно и в запросах сущностей:
  все десять полей, строковые коды, ссылки, отсутствующие поля, пустую строку, null,
  пустой объект и одновременную передачу строковой/структурированной формы.
  Локальный сервер не должен имитировать недоказанные правила слияния адресов.
- Модели и requests проверять с encoding/json и encoding/json/v2 по существующим таблицам кодеков.
  Канонизация JSON не должна скрывать потерю точности целых чисел Attribute.
- Тела Create/Update сверять с независимо заданными golden JSON, не произведёнными SDK.
  Для каждой сущности проверить минимальный Create, полный согласованный Create/Update,
  частичный Update, отсутствие незаданных полей, false/true, 0 там, где применимо,
  пустую строку, [], null и вложенные поля адресов.
- Проверить отсутствие readonly/TODO-полей в обоих request-типах;
  syncId должен присутствовать только в CreateRequest Counterparty/Organization.
- HTTP-тесты проверяют точные метод/маршрут, отсутствие лишних query/запросов,
  заголовки, отсутствие тела GET, обычный/gzip-ответ и серверные данные, отличающиеся от input.
- Для каждого метода проверить nil/отменённый context, локальную валидацию,
  ошибки авторизации/транспорта, APIError для 4xx/5xx/429, errors.As/errors.Is,
  неожиданный статус, malformed JSON, null, отсутствующий/пустой ID и отсутствие retry.
- Отказ Update предустановленной Country проверяется только как передача синтетической APIError;
  тест не доказывает реальный статус/текст живого API. Для ExpenseItem такой запрет не выдумывать.
- Каждый ticket завершается форматированием, go test -race -count=1 ./...,
  go vet ./... ./testdata, git diff --check и встроенным ревью Standards/Spec с фиксированным BASE.
- Не использовать живой аккаунт, .env или реальные операции записи.

## Out of Scope

- Банковские счета Counterparty/Organization целиком: модели чтения, metadata коллекций,
  вложенная запись и отдельные endpoints. Не подменять поддержку accounts урезанным Reference.
- Контактные лица, события, файлы и discounts Counterparty, включая их metadata и операции.
- Региональные объекты реквизитов mod__requisites__uz и mod__requisites__kz.
- companyVat__ru, advancePaymentVat, chiefAccountSign, directorSign и stamp Organization.
- fsrarId и utmUrl Organization: встречаются в официальных примерах, но не описаны
  в таблицах атрибутов; чтение и запись намеренно отложены. Их игнорирование не является
  случайной потерей согласованных полей.
- Изменение syncId через Update; запись isEgaisEnable; определения дополнительных полей через metadata.
- List, search, фильтрация, сортировка, пагинация, Delete, expand и массовые операции.
- Фасад, универсальный CRUD-исполнитель, общий рефакторинг Product/UOM, генерация сущностей,
  новые транспортные механизмы, retry, rate limiting и отдельное управление 429.
- Сущность Service и её конфликт имён; прочие сущности и поля вне этого перечня.
- Живой API, .env, примеры реальной записи, push, публикация PR и действия на VPS.

## Further Notes

- Пользователь утвердил Q1 «полный набор» для Country/ExpenseItem и Q2–Q4 «да»
  для перечисленного scope Counterparty/Organization, TODO и точек тестирования.
- Последующей правкой пользователь утвердил именованный CompanyType, явные TODO
  fsrarId/utmUrl, правило синхронизации копий и разбиение на пять tickets:
  Country; ExpenseItem; общие Address/AddressRequest/CompanyType;
  Counterparty; Organization. Последние два зависят только от общих типов, не друг от друга.
- Эта spec фиксирует готовые требования; tickets сохранены в локальном tracker.
  Статус ready-for-agent не разрешает начинать реализацию в текущей сессии.
- Планирование и передача контекста — здесь; реализация каждого утверждённого ticket —
  один основной агент в новой сессии. Дополнительных учебных остановок внутри ticket нет.
- Implement выполняет TDD на согласованных точках, проверки, встроенное ревью и локальный commit.
  Субагенты разрешены для Standards/Spec review; reviewer не меняет код.
- После ticket сообщить SHA, проверки, ограничения и следующий шаг; остановиться.
  Новые контракты сначала согласовать; отсутствие учебных остановок не отменяет этого правила.
- Локальные коммиты разрешены для реализации ticket. Эта spec, планирование и изменение общего
  glossary не коммитятся автоматически под видом реализации; push требует отдельной команды.
- Архитектурный образец Product принят на 9bca259; UOM.Get/Create/Update завершены на 3c30627.
- Официальные docs проверены на revision 9a811cef3dc4a3f3a46f070e60452361a06bf7cd.
  Документация — основной контракт; OpenAPI и fixtures могут отставать.

Источники:

- [Country: атрибуты и операции](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_country.md).
- [ExpenseItem: атрибуты и операции](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_expenseitem.md).
- [Counterparty: атрибуты и операции](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_counterparty.md).
- [Organization: атрибуты и операции](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_organization.md).
- [Общий контракт attributes](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L314-L367).
- [Адреса и null](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187).
- [Counterparty: замена accounts/contactpersons](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_counterparty.md#L3105-L3115).
- [Organization: правила счетов при создании](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_organization.md#L528-L532).

Accounts намеренно отложены: их read/write-формы различаются, Counterparty описывает
полную замену вложенной коллекции, а Organization имеет собственные валютные ограничения.
Эти контракты не выводятся из Product или из сходства двух сущностей.
