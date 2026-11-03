# Примеры Webhook для тестов

- Repository: https://github.com/moysklad/api-remap-1.2-doc
- Commit: 0e1b3d61fb8d8fa314469361d96560d6c3f7f57b
- Source: [md/dictionaries/_webhook.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_webhook.md)
  (SHA-256 `bbda00c5ee9ed4ef4f6eb7aebbf7be5e64c45f2817cc5bc53d6299031ab19307`)

| File | Раздел | SHA-256 |
| --- | --- | --- |
| `list.json` | «Получить список вебхуков», ответ (L132–208) | `a12ef5362d979d23a53ce5fa6d5ae9751e0fc7e71800c877c3b723bffea51ba2` |
| `create-fields.json` | «Создать вебхук», ответ с diffType FIELDS (L272–287) | `7e40956e6d1c138daa3682f15a157088e0068f2771917fb690ac23382811ede7` |
| `disabled.json` | «Отключить вебхук», ответ (L469–483) | `ab260be3ea57cb7b5c08005ca0a416e5a0c94575f9ee255ef2a78682089e753c` |

JSON скопирован без изменений; подтверждает форму примеров документации (authorApplication
только у части строк, diffType NONE/FIELDS, явный enabled=false), а не поведение живого API.

`create-full.golden.json` и `update.golden.json` — вручную заданные ожидаемые тела запросов
со всеми полями записи (url, action, entityType, diffType, enabled=false), независимые от
сериализации SDK; не ответы живого аккаунта.
