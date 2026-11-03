# Примеры уведомлений

- Repository: https://github.com/moysklad/api-remap-1.2-doc
- Commit: 0e1b3d61fb8d8fa314469361d96560d6c3f7f57b
- Sources: md/notification/_notification.md (SHA-256 `4caec6a6b5fc9d5a0b8540a15caa84d16c7bbc430609e86086a04a1514ab27da`),
  md/notification/_notifications.md (SHA-256 `55ae914d84aaacbc5f1d68da532c130848b3f8ea091811186562985f35d438f9`)

| File | Раздел | SHA-256 |
| --- | --- | --- |
| `get.json` | _notification.md, «Получить Уведомление», ответ (L24–50) | `7f00a636e19c5cb3a60e54feac68b2781b39063398560dd2df293dc828f50ec3` |
| `list.json` | _notifications.md, «Получить ленту Уведомлений», ответ (L45–159) | `457dcbe0eb8d727c7a9d1ea0e937e6713eb3494fd1341db23f1704c66bf2357a` |

JSON скопирован без изменений; подтверждает форму примеров, а не поведение живого API.
Ответ настроек — upstream-фикстура `testdata/upstream/notification_settings.json`, она совпадает
с примером «Получить настройки уведомлений».

`settings-update.golden.json` — вручную заданное ожидаемое тело PUT notification/settings:
явный `enabled: false`, канал `false`, частичная группа без `enabled` и неизвестные источнику
группа и канал (открытые ключи). Задано независимо от сериализации SDK.
