package transport

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultBaseURL = "https://api.moysklad.ru/api/remap/1.2"

	// Сервер обрабатывает массовый запрос до 5 минут; клиентский таймаут не должен обрывать его раньше.
	// https://dev.moysklad.ru/doc/api/remap/1.2/#/general#3-sozdanie-i-obnovlenie-neskolkih-obuektov
	DefaultTimeout = 6 * time.Minute
)

// Нулевое значение использует настройки по умолчанию.
type Config struct {
	BaseURL    string
	HTTPClient *http.Client

	// > 0 переопределяет таймаут; 0 сохраняет таймаут HTTPClient (у http.Client{} он не ограничен),
	// а для нового клиента — DefaultTimeout.
	Timeout time.Duration

	Auth Auth
}

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	auth       Auth
}

// HTTPClient копируется поверхностно: Transport и Jar общие с исходным клиентом,
// CheckRedirect оборачивается только в копии.
func NewClient(cfg Config) (*Client, error) {
	if cfg.Timeout < 0 {
		return nil, fmt.Errorf("moysklad: timeout must not be negative")
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("moysklad: invalid base URL: %w", err)
	}
	if (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Hostname() == "" {
		return nil, fmt.Errorf("moysklad: base URL must contain http/https scheme and host")
	}
	if parsedURL.User != nil || parsedURL.RawQuery != "" || parsedURL.ForceQuery || parsedURL.Fragment != "" {
		return nil, fmt.Errorf("moysklad: base URL must not contain credentials, query or fragment")
	}

	httpClient := http.Client{Timeout: DefaultTimeout}
	if cfg.HTTPClient != nil {
		httpClient = *cfg.HTTPClient
	}
	if cfg.Timeout > 0 {
		httpClient.Timeout = cfg.Timeout
	}
	httpClient.CheckRedirect = dropForeignAuthorization(parsedURL, httpClient.CheckRedirect)

	return &Client{
		httpClient: &httpClient,
		baseURL:    parsedURL,
		auth:       cfg.Auth,
	}, nil
}

// Стандартный клиент передаёт Authorization на поддомены и при смене схемы:
// https://pkg.go.dev/net/http#Client
// Заголовок удаляется до next: свой CheckRedirect может осознанно вернуть его доверенному хосту.
func dropForeignAuthorization(base *url.URL, next func(*http.Request, []*http.Request) error) func(*http.Request, []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if !sameOrigin(req.URL, base) {
			req.Header.Del("Authorization")
		}
		if next != nil {
			return next(req, via)
		}
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}
}

func sameOrigin(a, b *url.URL) bool {
	return strings.EqualFold(a.Scheme, b.Scheme) &&
		strings.EqualFold(a.Hostname(), b.Hostname()) &&
		effectivePort(a) == effectivePort(b)
}

func effectivePort(u *url.URL) string {
	if port := u.Port(); port != "" {
		return port
	}
	switch strings.ToLower(u.Scheme) {
	case "https":
		return "443"
	case "http":
		return "80"
	}
	return ""
}
