# Фикстуры

`product.json`, `product_full.json` и `product_folder.json` — официальные фикстуры МойСклад без правок. Источник — папка `tests/fixtures/` репозитория [moysklad/api-remap-1.2-openapi-specification](https://github.com/moysklad/api-remap-1.2-openapi-specification) на теге релиза. Их обновляет `scripts/update-spec.sh` вместе со спекой.

Версия спеки МойСклад: 0.33.0

Поле `updated` верхнего уровня в этих фикстурах записано в ISO 8601 (`2024-01-15T10:30:00.000+03:00`). Документация API и вложенные даты тех же фикстур используют формат `yyyy-MM-dd HH:mm:ss.SSS`. Поэтому `Timestamp` разбирает оба формата.
