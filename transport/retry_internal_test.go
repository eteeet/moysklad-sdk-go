package transport

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
)

func TestRetryableRequest(t *testing.T) {
	client, err := NewClient(Config{BaseURL: "https://api.example/api/remap/1.2"})
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		method   string
		endpoint string
		query    url.Values
		opts     RequestOptions
		want     bool
	}{
		{http.MethodGet, "entity/product", nil, RequestOptions{}, true},
		{http.MethodHead, "entity/product", nil, RequestOptions{}, true},
		{http.MethodGet, "report/stock/all?async=false", nil, RequestOptions{}, true},
		{http.MethodGet, "report/stock/all?asyncx=true&filter=a%3Dasync%3Dtrue", nil, RequestOptions{}, true},
		{http.MethodPost, "entity/product", nil, RequestOptions{}, false},
		{http.MethodPut, "entity/product", nil, RequestOptions{}, false},
		{http.MethodDelete, "entity/product", nil, RequestOptions{}, false},
		{http.MethodGet, "report/stock/all", url.Values{"async": {"true"}}, RequestOptions{}, false},
		{http.MethodGet, "report/stock/all?limit=1&async=true", nil, RequestOptions{}, false},
		// Регрессия ревью 46e9d06: закодированные ключ и значение — тоже async.
		{http.MethodGet, "report/stock/all?async=%74rue", nil, RequestOptions{}, false},
		{http.MethodGet, "report/stock/all?%61sync=true", nil, RequestOptions{}, false},
		{http.MethodGet, "report/stock/all?async=TRUE", nil, RequestOptions{}, false},
		{http.MethodGet, "report/stock/all", url.Values{"async": {"false", "true"}}, RequestOptions{}, false},
		// Явный query заменяет одноимённый параметр endpoint: итоговый async=false.
		{http.MethodGet, "report/stock/all?async=true", url.Values{"async": {"false"}}, RequestOptions{}, true},
		{http.MethodGet, "entity/product", nil, RequestOptions{NoRetry: true}, false},
	} {
		req, err := client.NewRequest(t.Context(), tc.method, tc.endpoint, tc.query, nil)
		if err != nil {
			t.Fatal(err)
		}
		if got := retryableRequest(req, tc.opts); got != tc.want {
			t.Errorf("%s %s %v %+v = %v", tc.method, tc.endpoint, tc.query, tc.opts, got)
		}
	}
	if !asyncRequested("a=%zz") {
		t.Error("unparsable query must be treated as async")
	}
}

func TestRetryDelay(t *testing.T) {
	clock := &fakeClock{now: testStart}
	policy := (&RetryPolicy{Clock: clock, Rand: func() float64 { return 0.5 }}).withDefaults()
	response := func(status int, header ...string) *http.Response {
		h := http.Header{}
		for i := 0; i < len(header); i += 2 {
			h.Set(header[i], header[i+1])
		}
		return &http.Response{StatusCode: status, Header: h}
	}
	apiErr := &sdkerrors.APIError{StatusCode: 503}
	for _, tc := range []struct {
		name    string
		attempt int
		resp    *http.Response
		err     error
		delay   time.Duration
		retry   bool
	}{
		{"429 jitter", 1, response(429), apiErr, 250 * time.Millisecond, true},
		{"503 second attempt", 2, response(503), apiErr, 500 * time.Millisecond, true},
		{"attempts exhausted", 3, response(503), apiErr, 0, false},
		{"lognex ms hint", 1, response(429, HeaderLognexRetryAfter, "1500"), apiErr, 1500 * time.Millisecond, true},
		{"reset ms hint", 1, response(429, HeaderReset, "2000", HeaderLognexRetryAfter, "100"), apiErr, 2 * time.Second, true},
		{"retry-after seconds", 1, response(503, HeaderRetryAfter, "40"), apiErr, 40 * time.Second, true},
		{"retry-after date", 1, response(503, HeaderRetryAfter, testStart.Add(10*time.Second).UTC().Format(http.TimeFormat)), apiErr, 10 * time.Second, true},
		{"hint over MaxWait gives up, not truncated", 1, response(429, HeaderLognexRetryAfter, "120000"), apiErr, 2 * time.Minute, false},
		{"invalid hint ignored", 1, response(429, HeaderLognexRetryAfter, "-5", HeaderRetryAfter, "soon"), apiErr, 250 * time.Millisecond, true},
		{"connection reset", 1, nil, &url.Error{Op: "Get", Err: syscall.ECONNRESET}, 250 * time.Millisecond, true},
		{"connection eof", 1, nil, fmt.Errorf("moysklad: send request: %w", io.EOF), 250 * time.Millisecond, true},
		{"timeout", 1, nil, &url.Error{Op: "Get", Err: timeoutError{}}, 250 * time.Millisecond, true},
		{"tls error", 1, nil, &url.Error{Op: "Get", Err: &tls.CertificateVerificationError{}}, 0, false},
		{"unsupported scheme", 1, nil, errors.New("unsupported protocol scheme"), 0, false},
		{"503 with unknown encoding", 1, response(503), errors.Join(apiErr, errors.New("unsupported content encoding")), 0, false},
		{"canceled", 1, nil, context.Canceled, 0, false},
		{"deadline", 1, nil, context.DeadlineExceeded, 0, false},
		{"decode error after 200", 1, response(200), errors.New("decode"), 0, false},
		{"400", 1, response(400), apiErr, 0, false},
		{"401", 1, response(401), apiErr, 0, false},
		{"501", 1, response(501), apiErr, 0, false},
	} {
		delay, retry := policy.retryDelay(tc.attempt, tc.resp, tc.err)
		if retry != tc.retry || retry && delay != tc.delay || !retry && tc.delay != 0 && delay != tc.delay {
			t.Errorf("%s: delay = %v, retry = %v", tc.name, delay, retry)
		}
	}

	// Экспонента ограничена MaxDelay, но не подсказка сервера.
	capped := (&RetryPolicy{Clock: clock, MaxAttempts: 100, MaxWait: time.Hour, Rand: func() float64 { return 0.999999 }}).withDefaults()
	if d, ok := capped.retryDelay(70, response(503), apiErr); !ok || d > DefaultMaxDelay || d < 29*time.Second {
		t.Fatalf("capped delay = %v, %v", d, ok)
	}
	if d, ok := capped.retryDelay(1, response(429, HeaderLognexRetryAfter, "45000"), apiErr); !ok || d != 45*time.Second {
		t.Fatalf("hint delay = %v, %v", d, ok)
	}
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

// Во время паузы между попытками слоты Limiter свободны.
func TestRetryReleasesSlotDuringPause(t *testing.T) {
	limiter, err := NewLimiter(LimiterConfig{Profile: ProfileApp, Concurrency: 1})
	if err != nil {
		t.Fatal(err)
	}
	var held []int
	clock := &hookClock{after: func() { held = append(held, len(limiter.slots), len(limiter.cfg.Account.slots)) }}
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			w.WriteHeader(503)
		}
	}))
	defer server.Close()
	client, err := NewClient(Config{BaseURL: server.URL, Limiter: limiter, Retry: &RetryPolicy{Clock: clock}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Send(t.Context(), http.MethodGet, "x", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	if len(held) != 2 || held[0] != 0 || held[1] != 0 || calls.Load() != 2 {
		t.Fatalf("slots during pause = %v, calls = %d", held, calls.Load())
	}
}

type hookClock struct{ after func() }

func (hookClock) Now() time.Time { return testStart }
func (c *hookClock) After(time.Duration) <-chan time.Time {
	c.after()
	ch := make(chan time.Time, 1)
	ch <- testStart
	return ch
}
