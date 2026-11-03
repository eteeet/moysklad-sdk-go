# Примеры ответов eventfeed

- Repository: https://github.com/moysklad/api-remap-1.2-doc
- Commit: 0e1b3d61fb8d8fa314469361d96560d6c3f7f57b
- Source: md/dictionaries/_eventfeed.md (SHA-256 `1af6f3a5079b81ef3a7f6623e9fe9485864e182873dc5694461ee2080fe3abb5`)

| File | Раздел | SHA-256 |
| --- | --- | --- |
| `list.json` | «Получить список Событий документа», ответ (L87–148) | `1ca17ca45ccb60e5333bb01d14fcd07bdda30ce1f0e842d2b650646dcbc382c6` |
| `get.json` | «Получить Событие», ответ (L175–194) | `815c55ca4f7471706dfde6a6c11f5baddf4b0256fef293226e0892f7d7772a7a` |
| `create.json` | «Добавить Событие», ответ — массив из одного события (L224–245) | `98db15b5038a8f1cc5ee4ac90b497aff48388f4e5650797738c8a2dcf74124a9` |
| `batch.json` | «Массовое создание и обновление Событий», ответ (L332–373) | `971847fabfebce0b05320caf50527abe50f6c5139501b3e281b8a82fed8d07c4` |

JSON скопирован без изменений; подтверждает форму примеров, а не поведение живого API.
В примере списка обе строки имеют один id — так в источнике.

`create.golden.json` — вручную заданное тело Create/Update (текст с упоминанием из
примера запроса источника), независимое от сериализации SDK.
