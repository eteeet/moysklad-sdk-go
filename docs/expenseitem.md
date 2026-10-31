# ExpenseItem

`expenseitem.NewService(client)` принимает существующий `*transport.Client`
и возвращает `(*Service, error)`; `nil` client отклоняется.
Методы `Get(ctx, id)`, `Create(ctx, CreateRequest)` и
`Update(ctx, id, UpdateRequest)` возвращают `(*ExpenseItem, error)`.

Реализованы `GET entity/expenseitem/{id}`, `POST entity/expenseitem` и
`PUT entity/expenseitem/{id}`. Успешный ответ — HTTP 200 и непустой `id`;
возвращаются фактические данные сервера. При ошибке модель равна `nil`.
Ошибки сохраняют контекст операции; `errors.As` находит `*sdkerrors.APIError`,
`errors.Is` — исходную ошибку транспорта или отмену context.

Модель чтения содержит все 13 полей таблицы атрибутов:
`meta`, `id`, `accountId`, `updated`, `name`, `description`, `code`,
`externalCode`, `owner`, `group`, `shared`, `archived`, `operatingExpenses`.
`updated` использует `model.DateTime`, `owner`/`group` — `*model.Reference`.
`code` сохраняется строкой, включая ведущие нули.

CreateRequest содержит обязательный `Name string`; `Description`, `Code`,
`ExternalCode`, `Owner`, `Group`, `Shared`, `Archived`, `OperatingExpenses`
используют `model.Field` с `omitzero`. UpdateRequest использует `model.Field`
для всех девяти полей записи, включая `Name`. Readonly-полей в requests нет.

Create отклоняет только `Name == ""`, без trim. Минимальный Create отправляет
только `name`: defaults `archived=false` и `operatingExpenses=true` применяет сервер.
Чтобы явно исключить статью из учёта прибыли, задай
`OperatingExpenses: model.Value(false)` — значение `false` уйдёт в JSON.
Отсутствующее `operatingExpenses` в ответе читается как нулевой `bool`, без default.

Get/Update проверяют форму UUID (Universally Unique Identifier), без ограничения версии.
Update отклоняет полностью незаданный request. Явные `false`, `""` и `null`
считаются заданными полями и кодируются без потери; допустимость значений определяет API.
Поддержка кодирования `null` не обещает серверную очистку поля.
Предварительного Get, локального слияния и автоматических повторов нет.
После ошибки Create с неизвестным результатом проверь аккаунт перед повторной отправкой.

Предустановленные ExpenseItem читаются тем же типом. У них могут отсутствовать
`accountId`, `owner`, `group`, `shared`; отсутствие даёт нулевые значения.
`Shared == false` не определяет происхождение статьи. Проверенная документация
ExpenseItem не устанавливает запрет Update предустановленных статей;
SDK не классифицирует их локально и передаёт ошибки API.

TODO: List, поиск, фильтрация, сортировка, пагинация, Delete, expand и массовые операции.
Поддержка ссылок не означает поддержку раскрытых объектов.

JSON проверен через `encoding/json` и `encoding/json/v2`, HTTP — с настоящим
transport и локальным сервером, включая gzip. Живой аккаунт не использовался.
[Происхождение fixtures и golden](../testdata/custom/expenseitem/README.md).

[Официальный контракт, revision 9a811cef](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_expenseitem.md).
