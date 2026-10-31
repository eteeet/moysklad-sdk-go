Status: resolved
Blocked by: None

# 03: Общие Address, AddressRequest и CompanyType

## What to build

Разработчик решения использует самостоятельные публичные типы model.Address,
model.AddressRequest и model.CompanyType с проверенным JSON-контрактом.
Это законченная возможность общей библиотеки моделей, а не только декларации для будущей сущности.
Counterparty и Organization подключат эти типы в своих независимых tickets.

Контракт: [общая spec](../spec.md), разделы общих типов, правил записи и Testing Decisions.
Предшествующие образцы общих типов — DateTime, Field и PriceType.

## Acceptance criteria

- [x] model.CompanyType — именованный тип поверх string, не alias. Известные значения заданы
  девятью типизированными константами: CompanyTypeLegal="legal", CompanyTypeEntrepreneur="entrepreneur",
  CompanyTypeIndividual="individual" и соответствующими тройками с суффиксами UZ/KZ в имени
  константы и значениях (legalUZ/entrepreneurUZ/individualUZ, legalKZ/entrepreneurKZ/individualKZ).
- [x] Набор CompanyType открыт: неизвестные строки читаются и записываются без изменений;
  нет локальной проверки значений, подстановки legal или другого default, ненужного JSON-кодека.
  Нулевое значение — пустая строка. MetaType не меняется.
- [x] Address читает строковые addInfo, apartment, city, comment, fiasCode__ru, house,
  postalCode, street и *Reference-поля country/region. Поля Go: AddInfo, Apartment,
  City, Comment, FIASCodeRU, House, PostalCode, Street, Country, Region.
- [x] AddressRequest содержит те же десять JSON-полей, но каждое — Field[string]
  либо Field[Reference] с omitzero. Модель чтения не подменяет request-тип.
- [x] Внешний Field[AddressRequest] поддерживает незаданное поле, значение, явный null и {}.
  Внутренние поля различают пропуск, пустую строку, null и ссылку; незаданные не отправляются.
  Отсутствующая/null-ссылка читается в nil, неизвестные поля адреса игнорируются.
- [x] Кодирование не меняет адрес: нет конкатенации, переноса строки в addInfo,
  нормализации кодов, локального merge или приписывания серверу правил замены частей адреса.
  Поддержка Null в Field не объявляет null командой очистки адреса.
- [x] Публичный JSON-контракт обоих типов адреса проверен через encoding/json и encoding/json/v2:
  значения всех десяти полей, ссылки, пропуски, null, строковые коды с ведущими нулями,
  частичный/полный request, явная пустая строка, пустой объект и Field-контейнеры.
- [x] CompanyType проверен обоими кодеками: все константы и точные значения, неизвестная строка,
  пустое значение и Field[CompanyType] с omitzero (незаданное/Value("")/Null).
- [x] Ожидаемый JSON задан независимо от результата SDK. Повторяющийся test harness не выделяется
  в новый production-модуль; существующие Field/Reference/DateTime не меняются без согласования.
- [x] Документация описывает готовые общие типы и серверные ограничения: при передаче обеих
  форм адреса структурированная приоритетнее, null не очищает, очистка — соответствующая строка "".
  Бизнес-поведение сервера не имитируется и не объявляется проверенным локальными тестами.
- [x] Не объявлены реализованными Counterparty/Organization или их методы заранее.
- [x] Пройдены gofmt, go test -race -count=1 ./..., go vet ./... ./testdata,
  git diff --check и отдельное встроенное ревью Standards/Spec относительно исходного BASE.

## Blocked by

Нет. Country, ExpenseItem и Region не требуются: внутри адреса используются Reference.

## Выполнение и границы

- Один основной агент, одна новая сессия, полный implement/TDD без учебных остановок.
  Точка проверки — публичная сериализация типов и Field-контейнеров, не приватное устройство.
- До правок прочитать инструкции и spec, записать исходный HEAD для ревью, проверить dirty diff.
  Сохранить и не включать в commit существовавшие изменения планирования или чужие изменения.
- Не добавлять HTTP-операции сущностей, счета, региональные реквизиты, генератор enum,
  универсальный CRUD, рефакторинг MetaType или новый механизм сериализации.
- Reviewer не меняет код. Новые решения контракта согласовать с пользователем.
- По завершении — локальный commit, Status: resolved и ## Answer с SHA/проверками/ревью,
  затем отчёт пользователю и остановка. Только после фактической готовности открываются 04 и 05.
  Push и живой API не разрешены.

## Sources

- [Counterparty: адреса и companyType](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_counterparty.md).
- [Organization: адреса и companyType](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_organization.md).
- [Адреса и null](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/_general.md#L1170-L1187).

## Answer

Реализованы общие Address, AddressRequest и открытый CompanyType с девятью
типизированными константами, независимыми тестами JSON-контракта обоих кодеков
и документацией `docs/requisites.md`.

- Commit: `386cd668efecf47f70421330f13e01ede00f8bf0`; включён в `dev` как `bfc0f76`.
- Проверки: gofmt; `go test -race -count=1 ./...`;
  `go vet ./... ./testdata`; `git diff --cached --check` — пройдены.
- Независимые ревью: Standards — 0 замечаний; Spec — 0 замечаний.
- HTTP-операции сущностей и живой API в этот ticket не входили.
