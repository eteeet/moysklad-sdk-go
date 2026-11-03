# Техкарта

- Repository: https://github.com/moysklad/api-remap-1.2-doc
- Commit: 0e1b3d61fb8d8fa314469361d96560d6c3f7f57b
- Source: md/dictionaries/_processingplan.md (SHA-256 `c3624eea2ae1476bce8f0d9bde67dad14eb3f06df47cc15c7191ccbabe5fee81`), примеры ответов

| File | SHA-256 |
| --- | --- |
| `read.json` | `3e0ccfc7ee514f5a882dbfe8192ea0325ad042b5fa2127fe273a2c10eb2787dd` |
| `stage.json` | `f03b621e2909b10eef35bde9894a55bbc19838b4d8e83160317e0c3a7a6743a5` |
| `material.json` | `d770c79336a9a0a21256d7b10c14a9cfd1715a251934ad65f2fb41df99d0471e` |
| `parametric-material.json` | `73929264bd1978235944115f02ee7c50e780d367ff01feba0f2ed9744e11eacd` |
| `product.json` | `191d3fb85f15881c2dd77012577fadd0844556daca3d1807b77c058f8b5b4fbd` |

Разделы: `read.json` — «Получить Техкарту» (первый пример), `stage.json` — «Получить отдельный этап
Техкарты», `material.json` — «Получить материал», `parametric-material.json` — пример «отдельного
параметрического материала» в разделе «Получить параметрические материалы Техкарты», `product.json` —
«Получить продукт» (ответ — массив из одного объекта). JSON скопирован без изменений; upstream-фикстуры
нет. Golden `create-full.golden.json` и `update.golden.json` написаны вручную по примерам создания и
изменения. Подтверждает форму ответов из примеров документации, а не поведение живого API.
