package transport

import (
	"fmt"
	"net/http"
)

// RequestOption настраивает один вызов; nil игнорируется.
// Опции не сохраняются в Client и не влияют на другие вызовы.
type RequestOption func(*RequestOptions)

// RequestOptions — итог применения опций одного вызова.
type RequestOptions struct {
	// Header добавляется поверх Config.Header; ключи уже канонические.
	Header http.Header

	// ResponseInfo заполняется метаданными последнего полученного ответа.
	ResponseInfo *ResponseInfo

	// Auth заменяет Config.Auth для одного вызова; nil оставляет Auth клиента.
	Auth Auth

	// Limiter заменяет Config.Limiter для одного вызова — нужен вместе с Auth другого principal.
	Limiter *Limiter

	// NoRetry запрещает повтор вызова независимо от Config.Retry.
	NoRetry bool
}

// ResponseInfo принадлежит вызывающему коду: параллельные вызовы должны передавать разные указатели.
type ResponseInfo struct {
	// Всё — о последней отправленной попытке: без ответа StatusCode 0 и Header nil.
	// Attempts — число отправленных попыток.
	StatusCode int
	Attempts   int
	// Копия всех заголовков ответа, включая повторные и неизвестные значения.
	Header http.Header
}

func newRequestOptions(opts ...RequestOption) RequestOptions {
	var options RequestOptions
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}
	return options
}

// WithHeader добавляет значение заголовка; повторный вызов с тем же key даёт несколько значений,
// а не CSV. Значения этого key из Config.Header для вызова заменяются.
func WithHeader(key, value string) RequestOption {
	return func(o *RequestOptions) {
		if o.Header == nil {
			o.Header = make(http.Header)
		}
		o.Header.Add(key, value)
	}
}

// WithResponseInfo заполняет info при получении ответа, в том числе при APIError.
// Публичная операция обнуляет info на входе (ResetResponseInfo), до локальной валидации:
// если ответа нет, info остаётся нулевым.
func WithResponseInfo(info *ResponseInfo) RequestOption {
	return func(o *RequestOptions) {
		o.ResponseInfo = info
	}
}

// ResetResponseInfo обнуляет ResponseInfo из opts, в том числе переданный через составную
// опцию. Публичные операции вызывают её первой строкой, до проверки аргументов.
func ResetResponseInfo(opts []RequestOption) {
	if info := newRequestOptions(opts...).ResponseInfo; info != nil {
		*info = ResponseInfo{}
	}
}

// WithAuth авторизует один вызов иначе, чем клиент: например, Basic для security/token
// на клиенте с Bearer. Config.Auth не меняется.
func WithAuth(auth Auth) RequestOption {
	return func(o *RequestOptions) {
		o.Auth = auth
	}
}

// WithLimiter учитывает вызов в бюджете другого principal, например пользователя, от имени
// которого вызов авторизован через WithAuth.
func WithLimiter(l *Limiter) RequestOption {
	return func(o *RequestOptions) {
		o.Limiter = l
	}
}

// WithoutRetry запрещает повтор этого вызова при любой RetryPolicy клиента.
func WithoutRetry() RequestOption {
	return func(o *RequestOptions) {
		o.NoRetry = true
	}
}

// Эти заголовки задаёт transport: их подмена ломает авторизацию, gzip или JSON-контракт.
var protectedHeaders = map[string]bool{
	"Authorization":   true,
	"Accept":          true,
	"Accept-Encoding": true,
	"Content-Type":    true,
	"Content-Length":  true,
	"Host":            true,
	// net/http повторяет запрос с этими заголовками при обрыве соединения: запись ушла бы дважды.
	"Idempotency-Key":   true,
	"X-Idempotency-Key": true,
}

func checkHeaders(header http.Header) error {
	seen := make(map[string]bool, len(header))
	for key := range header {
		canonical := http.CanonicalHeaderKey(key)
		if protectedHeaders[canonical] {
			return fmt.Errorf("moysklad: header %s is managed by transport", canonical)
		}
		// Иначе итог зависел бы от порядка обхода map.
		if seen[canonical] {
			return fmt.Errorf("moysklad: header %s is set twice with different case", canonical)
		}
		seen[canonical] = true
	}
	return nil
}

// Значения одного key из src заменяют значения dst целиком.
// HeaderValues возвращает значения key, которые получит вызов с opts: Config.Header с заменой
// опциями вызова, как при отправке. Запрос не строится и не отправляется.
func (c *Client) HeaderValues(key string, opts ...RequestOption) []string {
	header := make(http.Header)
	overrideHeaders(header, c.header)
	overrideHeaders(header, newRequestOptions(opts...).Header)
	return header.Values(key)
}

func overrideHeaders(dst, src http.Header) {
	for key, values := range src {
		dst[http.CanonicalHeaderKey(key)] = append([]string(nil), values...)
	}
}
