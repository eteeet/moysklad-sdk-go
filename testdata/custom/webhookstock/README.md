# Примеры WebhookStock для тестов

- Repository: https://github.com/moysklad/api-remap-1.2-doc
- Commit: 0e1b3d61fb8d8fa314469361d96560d6c3f7f57b
- Source: [md/dictionaries/_webhookstock.md](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_webhookstock.md)
  (SHA-256 `41064194f1761b379b566b1c72ee7f3e3978544dc26e6707155bcb0266f65018`)

| File | Раздел | SHA-256 |
| --- | --- | --- |
| `list.json` | «Получить список вебхуков на изменение остатков», ответ (L67–139) | `c64b80d145ca61ee53165721871f55c61dce107909460848ea7f859791ed7b2c` |
| `disabled.json` | «Отключить вебхук на изменение остатков», ответ (L360–373) | `6d2f6c09d16e5d6f12a0691fa2036842927cce0f7b06297e24e3537d4442acc0` |

JSON скопирован без изменений; подтверждает форму примеров документации (authorApplication
только у первой строки, явный enabled=false), а не поведение живого API.

`create-full.golden.json` и `update.golden.json` — вручную заданные ожидаемые тела запросов
со всеми полями записи (url, stockType, reportType, enabled), независимые от сериализации SDK;
не ответы живого аккаунта.
