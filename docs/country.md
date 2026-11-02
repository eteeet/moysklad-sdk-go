# Country

Пакет `entity/country` поддерживает чтение Country (страны), создание
пользовательской Country и изменение отдельных полей пользовательской Country.
`country.NewService(client)` принимает `*transport.Client`, возвращает
`(*country.Service, error)` и отклоняет nil client.

Методы `Get(ctx, id)`, `Create(ctx, country.CreateRequest)` и
`Update(ctx, id, country.UpdateRequest)` возвращают `(*country.Country, error)`.
Маршруты: GET `entity/country/{id}`, POST `entity/country`, PUT
`entity/country/{id}`. Успешный ответ требует HTTP 200 и непустой ID.
Create и Update возвращают фактические данные ответа сервера.

Модель читает `meta`, `id`, `accountId`, `updated`, `name`, `description`,
`code`, `externalCode`, `owner`, `group`, `shared`. `Code` остаётся строкой,
включая ведущие нули; `Updated` использует `model.DateTime`, ссылки —
`*model.Reference`, без expand. Неизвестные поля ответа игнорируются.
У предустановленных Country могут отсутствовать `accountId`, `owner`, `group`
и `shared`; они дают нулевые значения Go. `Shared == false` не определяет
происхождение Country.

Поля записи: `name`, `description`, `code`, `externalCode`, `owner`, `group`,
`shared`. В Create только `Name` — обязательный string; пустая строка отклоняется,
пробелы не обрезаются. Остальные поля Create и все поля Update используют
`model.Field` с `omitzero`: незаданные поля не отправляются, `model.Value(false)`,
`model.Value("")` и `model.Null[T]()` сохраняются. Полностью незаданный Update
отклоняется. Кодирование null не обещает его принятие или очистку поля сервером.
Get и Update проверяют форму UUID (Universally Unique Identifier), не ограничивая
версию.

Передай context для отмены запроса. При ошибке результат nil; `errors.As`
распознаёт `*sdkerrors.APIError`, `errors.Is` сохраняет ошибки context,
авторизации и транспорта. Предварительного Get, локального слияния, переходов
по ссылкам и автоматических повторов нет, в том числе при HTTP 429.
Документация запрещает изменение предустановленных Country; SDK передаёт
ошибку сервера без локальной классификации Country.

List — [transport.md](transport.md), Delete и массовые операции — [write-async.md](write-async.md);
expand нет.

Контракт сверён с [официальной документацией Country на revision
9a811cef3dc4a3f3a46f070e60452361a06bf7cd](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_country.md).
