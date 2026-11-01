package transport_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// recordClock сразу «просыпается» и запоминает запрошенные паузы.
type recordClock struct {
	mu     sync.Mutex
	delays []time.Duration
	cancel context.CancelFunc // если задан, отменяет контекст вместо пробуждения
}

func (c *recordClock) Now() time.Time { return time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC) }
func (c *recordClock) After(d time.Duration) <-chan time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.delays = append(c.delays, d)
	if c.cancel != nil {
		c.cancel()
		return nil
	}
	ch := make(chan time.Time, 1)
	ch <- c.Now()
	return ch
}

type step struct {
	status int
	header string // X-Lognex-Retry-After, мс
	body   string
	hangup bool
}

func retryClient(t *testing.T, clock *recordClock, steps ...step) (*transport.Client, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := int(calls.Add(1))
		// Каждая попытка собрана заново: заголовок вызова ровно один, не накоплен.
		if v := r.Header.Values("X-Call"); len(v) != 1 || v[0] != "1" {
			t.Errorf("attempt %d headers = %v", n, r.Header)
		}
		if n > len(steps) {
			t.Errorf("unexpected attempt %d", n)
			return
		}
		s := steps[n-1]
		if s.hangup {
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
			return
		}
		// Новое соединение на каждую попытку: иначе net/http сам повторит GET после обрыва
		// переиспользованного keep-alive соединения, и попыток станет больше.
		w.Header().Set("Connection", "close")
		w.Header().Set("X-Attempt", strconv.Itoa(n))
		if s.header != "" {
			w.Header().Set(transport.HeaderLognexRetryAfter, s.header)
		}
		w.WriteHeader(s.status)
		w.Write([]byte(s.body))
	}))
	t.Cleanup(server.Close)
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL,
		Retry: &transport.RetryPolicy{Clock: clock, Rand: func() float64 { return 0.5 }}})
	if err != nil {
		t.Fatal(err)
	}
	return client, &calls
}

const apiError = `{"errors":[{"error":"x","code":1049}]}`

func TestSendRetry(t *testing.T) {
	for _, tc := range []struct {
		name     string
		method   string
		query    url.Values
		opts     []transport.RequestOption
		steps    []step
		calls    int32
		status   int
		delays   []time.Duration
		apiError bool
		wantErr  bool
	}{
		{"503 then 429 hint then ok", "GET", nil, nil,
			[]step{{status: 503, body: apiError}, {status: 429, header: "1500", body: apiError}, {status: 200, body: `{}`}},
			3, 200, []time.Duration{250 * time.Millisecond, 1500 * time.Millisecond}, false, false},
		{"429 concurrency 1073", "GET", nil, nil,
			[]step{{status: 429, body: `{"errors":[{"error":"x","code":1073}]}`}, {status: 200, body: `{}`}},
			2, 200, []time.Duration{250 * time.Millisecond}, false, false},
		{"connection drop then ok", "GET", nil, nil, []step{{hangup: true}, {status: 200, body: `{}`}},
			2, 200, []time.Duration{250 * time.Millisecond}, false, false},
		{"final drop clears info", "GET", nil, nil, []step{{status: 503, body: apiError}, {hangup: true}, {hangup: true}},
			3, 0, []time.Duration{250 * time.Millisecond, 500 * time.Millisecond}, false, true},
		{"exhausted keeps last APIError", "HEAD", nil, nil,
			[]step{{status: 502}, {status: 502}, {status: 504}}, 3, 504, []time.Duration{250 * time.Millisecond, 500 * time.Millisecond}, true, true},
		{"POST never", "POST", nil, nil, []step{{status: 503, body: apiError}}, 1, 503, nil, true, true},
		{"PUT never", "PUT", nil, nil, []step{{status: 503, body: apiError}}, 1, 503, nil, true, true},
		{"async GET never", "GET", url.Values{"async": {"true"}}, nil, []step{{status: 503, body: apiError}}, 1, 503, nil, true, true},
		{"WithoutRetry", "GET", nil, []transport.RequestOption{transport.WithoutRetry()}, []step{{status: 503, body: apiError}}, 1, 503, nil, true, true},
		{"decode error after 200", "GET", nil, nil, []step{{status: 200, body: `{"a":`}}, 1, 200, nil, false, true},
		{"400 never", "GET", nil, nil, []step{{status: 400, body: apiError}}, 1, 400, nil, true, true},
		{"hint over MaxWait", "GET", nil, nil, []step{{status: 429, header: "600000", body: apiError}}, 1, 429, nil, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			clock := &recordClock{}
			client, calls := retryClient(t, clock, tc.steps...)
			var info transport.ResponseInfo
			var dst struct{}
			opts := append([]transport.RequestOption{transport.WithHeader("X-Call", "1"), transport.WithResponseInfo(&info)}, tc.opts...)
			_, err := client.Send(t.Context(), tc.method, "entity/product", tc.query, nil, &dst, opts...)
			var apiErr *sdkerrors.APIError
			if calls.Load() != tc.calls || info.Attempts != int(tc.calls) || info.StatusCode != tc.status ||
				errors.As(err, &apiErr) != tc.apiError || (err != nil) != tc.wantErr {
				t.Fatalf("calls = %d, info = %+v, err = %v", calls.Load(), info, err)
			}
			if tc.status == 0 && info.Header != nil || tc.status != 0 && info.Header.Get("X-Attempt") != strconv.Itoa(int(tc.calls)) {
				t.Fatalf("info is not from the last attempt: %v", info.Header)
			}
			if len(clock.delays) != len(tc.delays) {
				t.Fatalf("delays = %v, want %v", clock.delays, tc.delays)
			}
			for i := range tc.delays {
				if clock.delays[i] != tc.delays[i] {
					t.Fatalf("delays = %v, want %v", clock.delays, tc.delays)
				}
			}
		})
	}
}

// Отмена во время паузы и пауза длиннее дедлайна завершают вызов последним ответом.
func TestSendRetryContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	clock := &recordClock{cancel: cancel}
	client, calls := retryClient(t, clock, step{status: 503, body: apiError})
	_, err := client.Send(ctx, "GET", "entity/product", nil, nil, nil, transport.WithHeader("X-Call", "1"))
	if sdkerrors.StatusCode(err) != 503 || !errors.Is(err, context.Canceled) || calls.Load() != 1 {
		t.Fatalf("err = %v, calls = %d", err, calls.Load())
	}

	clock = &recordClock{}
	client, calls = retryClient(t, clock, step{status: 429, header: "5000", body: apiError})
	deadline, cancelDeadline := context.WithTimeout(t.Context(), time.Second)
	defer cancelDeadline()
	_, err = client.Send(deadline, "GET", "entity/product", nil, nil, nil, transport.WithHeader("X-Call", "1"))
	if sdkerrors.StatusCode(err) != 429 || errors.Is(err, context.DeadlineExceeded) || calls.Load() != 1 || len(clock.delays) != 0 {
		t.Fatalf("err = %v, calls = %d, delays = %v", err, calls.Load(), clock.delays)
	}
}

// Каждая попытка заново проходит Limiter.
func TestSendRetryLimiterPerAttempt(t *testing.T) {
	clock := &recordClock{}
	var admissions atomic.Int32
	limiter, err := transport.NewLimiter(transport.LimiterConfig{Profile: transport.ProfileApp, Clock: clock,
		Weight: func(transport.Profile, string, string, time.Time) int { admissions.Add(1); return 1 }})
	if err != nil {
		t.Fatal(err)
	}
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			w.WriteHeader(503)
		}
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL, Limiter: limiter,
		Retry: &transport.RetryPolicy{Clock: clock, Rand: func() float64 { return 0 }}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Send(t.Context(), "GET", "entity/product", nil, nil, nil); err != nil || admissions.Load() != 2 {
		t.Fatalf("err = %v, admissions = %d", err, admissions.Load())
	}
}

// Отмена, пока вторая попытка ждёт Limiter: итог — первая попытка вместе с отменой.
func TestSendRetryCanceledInLimiter(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	limiter, err := transport.NewLimiter(transport.LimiterConfig{Profile: transport.ProfileApp, Budget: 1,
		Clock: &recordClock{cancel: cancel}})
	if err != nil {
		t.Fatal(err)
	}
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(503)
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL, Limiter: limiter,
		Retry: &transport.RetryPolicy{Clock: &recordClock{}}})
	if err != nil {
		t.Fatal(err)
	}
	var info transport.ResponseInfo
	_, err = client.Send(ctx, "GET", "entity/product", nil, nil, nil, transport.WithResponseInfo(&info))
	if sdkerrors.StatusCode(err) != 503 || !errors.Is(err, context.Canceled) || calls.Load() != 1 ||
		info.Attempts != 1 || info.StatusCode != 503 {
		t.Fatalf("err = %v, calls = %d, info = %+v", err, calls.Load(), info)
	}
}

// Параллельные вызовы делят одну RetryPolicy и её Rand.
func TestSendRetryConcurrent(t *testing.T) {
	var mu sync.Mutex
	seen := map[string]bool{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		first := !seen[r.Header.Get("X-Id")]
		seen[r.Header.Get("X-Id")] = true
		mu.Unlock()
		if first {
			w.WriteHeader(503)
		}
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL, Retry: &transport.RetryPolicy{Clock: &recordClock{}}})
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := range 16 {
		wg.Go(func() {
			var info transport.ResponseInfo
			_, err := client.Send(t.Context(), "GET", "entity/product", nil, nil, nil,
				transport.WithHeader("X-Id", strconv.Itoa(i)), transport.WithResponseInfo(&info))
			if err != nil || info.Attempts != 2 {
				t.Errorf("err = %v, info = %+v", err, info)
			}
		})
	}
	wg.Wait()
}
