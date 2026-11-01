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
}

// ResponseInfo принадлежит вызывающему коду: параллельные вызовы должны передавать разные указатели.
type ResponseInfo struct {
	StatusCode int
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
// Если ответа нет, info обнуляется.
func WithResponseInfo(info *ResponseInfo) RequestOption {
	return func(o *RequestOptions) {
		o.ResponseInfo = info
	}
}

// WithAuth авторизует один вызов иначе, чем клиент: например, Basic для security/token
// на клиенте с Bearer. Config.Auth не меняется.
func WithAuth(auth Auth) RequestOption {
	return func(o *RequestOptions) {
		o.Auth = auth
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
func overrideHeaders(dst, src http.Header) {
	for key, values := range src {
		dst[http.CanonicalHeaderKey(key)] = append([]string(nil), values...)
	}
}
