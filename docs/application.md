# Application

Пакет `entity/application` читает установленное на аккаунте серверное решение по id.
`application.NewService(client)` принимает `*transport.Client`, возвращает
`(*application.Service, error)` и отклоняет nil client.

`Get(ctx, id)` выполняет один GET `entity/application/{id}` и возвращает
`(*application.Application, error)`. `id` — id установки решения на аккаунте, проверяется как
UUID любой версии до HTTP. Успех — HTTP 200 и непустой `id`. Повторов нет; доступ проверяет
сервер. При ошибке результат nil, `errors.As` распознаёт `*sdkerrors.APIError`, `errors.Is`
сохраняет ошибки context, авторизации и транспорта.

Контекст решения, токеном которого выполняется запрос (GET `context/application`), читает
`contexts.Service.Application` ([пакет contexts](../contexts/contexts.go)); здесь он не дублируется.
Ответы обоих endpoints по источнику одной формы, поэтому `application.Application` — псевдоним
`contexts.ApplicationContext`.

## Поля

| JSON | Go | Запись |
| --- | --- | --- |
| `meta` | `Meta model.Meta` | нет |
| `id` | `ID string` | нет |
| `appUid` | `AppUID string` | нет |

`appUid` — строка вида `<Application-Alias>.<Vendor-Alias>@<Account-Name>`; SDK её не разбирает.
Её используют в `filter=updatedBy=...` списков и `filter=uid=...` аудита, а href решения —
в `filter=application=...` аудита; эти фильтры собираются через `query`.

## Ограничения и противоречия источника

- Таблицы атрибутов у решения нет: поля взяты из примера ответа (как и у `context/application`);
  обязательность и readonly источник не размечает. Записи нет вовсе.
- Список решений, установка, изменение и удаление в JSON API не описаны и не реализованы;
  Vendor API — вне SDK.

Контракт сверён с [официальной документацией на revision
0e1b3d61fb8d8fa314469361d96560d6c3f7f57b](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/_general.md), раздел
«Серверные решения» → «Получение сущности установленного решения».
