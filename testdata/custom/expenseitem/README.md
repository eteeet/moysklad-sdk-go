# Примеры ExpenseItem для тестов

`predefined.json` — объект «Возврат» из `rows` ответа «Получить Статьи расходов»
официальной документации. Значения объекта скопированы без изменений; сохранён
один элемент списка, а не полный ответ. Он содержит `operatingExpenses=false`
и не содержит `accountId`, `owner`, `group`, `shared`.

- Repository: https://github.com/moysklad/api-remap-1.2-doc
- Commit: `9a811cef3dc4a3f3a46f070e60452361a06bf7cd`
- Source: [md/dictionaries/_expenseitem.md, раздел «Получить Статьи расходов»](https://github.com/moysklad/api-remap-1.2-doc/blob/9a811cef3dc4a3f3a46f070e60452361a06bf7cd/md/dictionaries/_expenseitem.md#получить-статьи-расходов).

Upstream `testdata/upstream/expense_item.json` не менялся. Он не содержит
`operatingExpenses`, `owner`, `group`, `shared`; отдельные небольшие JSON в
`entity/expenseitem/model_test.go` проверяют ссылки, bool и строковый `code="003"`.

`create.golden.json` и `update.golden.json` — независимо заданные вручную
ожидаемые тела запросов со всеми девятью полями записи по согласованному контракту.
Это синтетические тестовые запросы, не примеры записи в живой аккаунт.
В Create проверяются `archived=true` и `operatingExpenses=false`, в Update —
`archived=false`, `operatingExpenses=true` и пустое `description`.
Минимальные/частичные запросы и `null` проверяются отдельными JSON-литералами.

| File | SHA-256 |
| --- | --- |
| `predefined.json` | `f3514c57bea922b3da073c867e31a1170696651934ded8df48154ef245495ebb` |
