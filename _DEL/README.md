# moysklad-sdk-go

Go SDK для [JSON API 1.2 МойСклад](https://dev.moysklad.ru/doc/api/remap/1.2/).

Версия спеки МойСклад: 0.33.0

Модели генерирует oapi-codegen из [официальной спеки](https://github.com/moysklad/api-remap-1.2-openapi-specification), транспорт и клиент написаны вручную (ADR-0015 в `eteeet/meta`). Сейчас поддержаны Product и ProductFolder: чтение по id, страница списка и обход всех страниц.

## Пример

```go
c := moysklad.New(token,
	moysklad.WithHTTPClient(&http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}),
	moysklad.WithLogger(slog.Default()),
	moysklad.WithLimiter(limiter), // свой бюджет запросов аккаунта
)

p, err := c.Product.Get(ctx, id)
var msErr *moysklad.Error
if errors.As(err, &msErr) && msErr.StatusCode == http.StatusNotFound {
	// товара нет
}

// Все товары, изменённые за сутки, и архивные, и обычные. Время в фильтре SDK переводит в МСК.
filter := moysklad.Filter{}.
	Gte("updated", time.Now().Add(-24*time.Hour)).
	Eq("archived", "true").Eq("archived", "false")
for p, err := range c.Product.All(ctx, moysklad.ListOptions{Filter: filter}) {
	if err != nil {
		return err
	}
	fmt.Println(*p.Name, p.Updated.Time)
}
```

Транспорт:

- ставит Bearer-токен и `Accept-Encoding: gzip`;
- при 429 повторяет запрос через паузу из `X-Lognex-Retry-After`. Если задан ограничитель, паузу выдерживает он (`Limiter.Backoff`);
- при 5xx повторяет GET, PUT и DELETE с растущей паузой, POST не повторяет;
- повторяет в пределах `ctx` и `WithRetry` (по умолчанию 5 попыток);
- возвращает ошибку МойСклад как `*moysklad.Error` с кодом и текстом.

Эндпоинты, которых ещё нет в SDK, вызываются через `c.Do` с тем же транспортом.

## Разработка

| Команда | Что делает |
|---|---|
| `make test` | тесты |
| `make lint` | golangci-lint |
| `make gen` | модели из `api/openapi.yaml` с правками `api/overlay.yaml` |
| `scripts/update-spec.sh 0.34.0 && make gen` | обновить спеку и фикстуры до версии релиза МойСклад |

## Лицензия

Apache-2.0.
