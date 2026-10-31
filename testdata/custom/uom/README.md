# Примеры UOM для тестов

`predefined.json` — JSON ответа «Получить Единицу измерения» из официальной
документации: предустановленная UOM «дюйм», code=039. Данные скопированы без изменений.

- Repository: https://github.com/moysklad/api-remap-1.2-doc
- Commit: d1ac440ecafc2f9299cbcb4cb99927d3f933462f
- Source: [md/dictionaries/_uom.md, раздел «Получить Единицу измерения»](https://github.com/moysklad/api-remap-1.2-doc/blob/d1ac440ecafc2f9299cbcb4cb99927d3f933462f/md/dictionaries/_uom.md).

Пример дополняет upstream-фикстуру: проверяет отсутствие accountId, owner, group
и shared, а также строковый code с ведущим нулём. Автоматический снимок upstream
не изменялся.

`create.golden.json` — ожидаемое тело создания UOM со всеми семью поддержанными
полями. Задано вручную по согласованному контракту, независимо от сериализации SDK;
это тестовый запрос, а не ответ живого аккаунта.
