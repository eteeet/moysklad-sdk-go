# Примеры ответов task

- Repository: https://github.com/moysklad/api-remap-1.2-doc
- Commit: 0e1b3d61fb8d8fa314469361d96560d6c3f7f57b
- Source: md/dictionaries/_task.md (SHA-256 `6978994ed5bbb45aa7b162f857ea8224bdcd4f7145a2840ff2448bc48e33a112`)

| File | Раздел | SHA-256 |
| --- | --- | --- |
| `list.json` | «Получить задачи», ответ (L131–256) | `73f80f4ad593f8c1fd497bff74abc8424f7fb7ca007d7ba9adf3d23b1a2f31dc` |
| `batch.json` | «Массовое создание и обновление Задач», ответ (L426–541) | `e8f04c7348a277aef26993b2d0895f6f9e50886559c7546cdbff28ce56c5365c` |
| `comments.json` | «Получить комментарии Задачи», ответ (L849–906) | `1bbaa188300d2c63250235ec4bbb16f2881071d82c947db7c88e6b925a3efa02` |
| `comment-create.json` | «Создать комментарий Задачи», ответ одного (L940–971) | `e2d0487720f03f8f2a9670113df910ffdcdef5725eef97c8ee0bebf8a22ba807` |
| `comment.json` | «Получить комментарий к Задаче», ответ (L1086–1103) | `521ce660dd3a91c3d989d21543c6ce270152cf4cd4aea4f230f17a0cc33fbd4e` |
| `metadata.json` | «Получить типы задач», ответ (L1232–1267) | `673432e45e0c9650faeaa784579e951dba0c3d9e23540bde17da1fb163011259` |
| `state.json` | «Создать тип задачи», ответ (L1295–1308) | `f604edef2521d34ba2570f6cb07f2604f9f6d9303469e912053b5970946414f1` |
| `states-batch.json` | «Массовое создание и обновление Типов задач», ответ (L1394–1423) | `ce86dc4eb17cdf85b3eb83975926bd93ee102a0a775cfbfc63f01d4719dfa5f1` |

JSON скопирован без изменений; подтверждает форму примеров, а не поведение живого API.
Примеры Get/Update комментария не содержат id (критерий успеха SDK) — `comment.json`
проверяется только разбором. Get задачи проверяется upstream-фикстурой `task.json`.

`create.golden.json`, `update.golden.json`, `comment.golden.json` и
`state-create.golden.json` — вручную заданные ожидаемые тела запросов, независимые от
сериализации SDK. Create содержит явный `done: false`, Update — null у dueToDate и agent
и смену привязки на документ (operation).
