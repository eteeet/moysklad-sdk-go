package moysklad

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

const defaultBaseURL = "https://api.moysklad.ru/api/remap/1.2"

// Client — клиент JSON API одного аккаунта. Одна сущность — одно поле.
type Client struct {
	Product       Entity[Product]
	ProductFolder Entity[ProductFolder]

	t *transport
}

// Do выполняет произвольный запрос к API с теми же токеном, gzip, ограничителем и повторами,
// что и методы сущностей. path — путь от корня API, например "entity/product".
// in кодируется в JSON тела запроса, ответ разбирается в out; nil — без тела.
func (c *Client) Do(ctx context.Context, method, path string, query url.Values, in, out any) error {
	err := c.t.do(ctx, method, path, query, in, out)
	if err != nil {
		return fmt.Errorf("transport.do: %w", err)
	}

	return nil
}

// Option настраивает Client.
type Option func(*transport)

// WithBaseURL задаёт адрес API вместо https://api.moysklad.ru/api/remap/1.2.
func WithBaseURL(u string) Option {
	return func(t *transport) { t.baseURL = u }
}

// WithHTTPClient задаёт HTTP-клиент, например с otelhttp. По умолчанию — http.DefaultClient.
func WithHTTPClient(hc *http.Client) Option {
	return func(t *transport) { t.hc = hc }
}

// WithRetry задаёт число попыток на запрос, включая первую, и паузу перед повтором при 5xx.
// Пауза растёт с номером попытки. При 429 пауза берётся из X-Lognex-Retry-After.
// По умолчанию — 5 попыток и 1 с.
func WithRetry(attempts int, pause time.Duration) Option {
	return func(t *transport) { t.attempts, t.pause = max(attempts, 1), pause }
}

// WithLimiter задаёт ограничитель запросов. По умолчанию запросы не ограничиваются.
func WithLimiter(l Limiter) Option {
	return func(t *transport) { t.limiter = l }
}

// WithLogger задаёт логгер. SDK пишет строку на каждый повтор запроса. По умолчанию логов нет.
func WithLogger(l *slog.Logger) Option {
	return func(t *transport) { t.log = l }
}

// New создаёт клиент с токеном доступа аккаунта.
func New(token string, opts ...Option) *Client {
	t := &transport{
		baseURL:  defaultBaseURL,
		token:    token,
		hc:       http.DefaultClient,
		attempts: defaultAttempts,
		pause:    defaultPause,
		log:      slog.New(slog.DiscardHandler),
	}
	for _, o := range opts {
		o(t)
	}

	return &Client{
		Product:       Entity[Product]{t: t, path: "entity/product"},
		ProductFolder: Entity[ProductFolder]{t: t, path: "entity/productfolder"},
		t:             t,
	}
}
