`read.json` — пример ответа «Получить Техпроцесс» раздела
[processingprocess](https://github.com/moysklad/api-remap-1.2-doc/blob/0e1b3d61fb8d8fa314469361d96560d6c3f7f57b/md/dictionaries/_processingprocess.md)
без изменений; upstream-фикстуры нет. Golden: `create-full` — пример «Создать Техпроцесс» с
остальными полями таблицы; `update` — пример «Изменить Техпроцесс» (позиция с meta и новая позиция)
плюс `nextPositions`, null и третья строка — проверяют порядок строк и форму полей SDK.
