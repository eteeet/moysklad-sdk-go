package transport_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/transport"
)

// stuckClock никогда не срабатывает и сообщает о каждом ожидании.
type stuckClock struct{ waits chan time.Duration }

func (stuckClock) Now() time.Time { return time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC) }
func (c stuckClock) After(d time.Duration) <-chan time.Time {
	c.waits <- d
	return nil
}

// Слот занят, пока ответ читается: в полёте не больше Concurrency запросов.
func TestSendLimiterConcurrency(t *testing.T) {
	var inFlight, peak atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := inFlight.Add(1)
		for p := peak.Load(); n > p && !peak.CompareAndSwap(p, n); p = peak.Load() {
		}
		w.Write([]byte(`{"a":`))
		w.(http.Flusher).Flush()
		inFlight.Add(-1) // декремент до конца тела: клиент ещё держит слот
		w.Write([]byte(`1}`))
	}))
	defer server.Close()
	limiter, err := transport.NewLimiter(transport.LimiterConfig{Profile: transport.ProfileApp, Concurrency: 2})
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL, Limiter: limiter})
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			var dst struct{ A int }
			if _, err := client.Send(t.Context(), http.MethodGet, "entity/product", nil, nil, &dst); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if peak.Load() > 2 {
		t.Fatalf("peak in flight = %d, want <= 2", peak.Load())
	}
}

// Вес считается по пути относительно BaseURL, включая префикс /api/remap/1.2.
func TestSendLimiterPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()
	var paths []string
	limiter, err := transport.NewLimiter(transport.LimiterConfig{Profile: transport.ProfileApp,
		Weight: func(p transport.Profile, method, path string, now time.Time) int {
			paths = append(paths, method+" "+path)
			return transport.DefaultWeight(p, method, path, now)
		}})
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", Limiter: limiter})
	if err != nil {
		t.Fatal(err)
	}
	for _, endpoint := range []string{"report/stock/all", "report/stock/bystore/", server.URL + "/api/remap/1.2/report/stock/all?offset=1"} {
		if _, err := client.Send(t.Context(), http.MethodGet, endpoint, nil, nil, nil); err != nil {
			t.Fatal(err)
		}
	}
	want := []string{"GET report/stock/all", "GET report/stock/bystore", "GET report/stock/all"}
	for i := range want {
		if i >= len(paths) || paths[i] != want[i] {
			t.Fatalf("paths = %v", paths)
		}
	}
}

// Вызов с WithAuth не расходует бюджет клиента; WithLimiter учитывает его в бюджете своего principal.
func TestSendLimiterWithAuth(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	clock := stuckClock{waits: make(chan time.Duration, 10)}
	newLimiter := func() *transport.Limiter {
		l, err := transport.NewLimiter(transport.LimiterConfig{Profile: transport.ProfileApp, Budget: 1, Clock: clock})
		if err != nil {
			t.Fatal(err)
		}
		return l
	}
	appLimiter, employeeLimiter := newLimiter(), newLimiter()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL, Limiter: appLimiter})
	if err != nil {
		t.Fatal(err)
	}
	basic, err := transport.BasicAuth("u", "p")
	if err != nil {
		t.Fatal(err)
	}
	send := func(opts ...transport.RequestOption) error {
		_, err := client.Send(t.Context(), http.MethodGet, "entity/product", nil, nil, nil, opts...)
		return err
	}
	// Бюджет клиента исчерпан первым вызовом; вызовы с WithAuth проходят без ожидания.
	for _, opts := range [][]transport.RequestOption{nil, {transport.WithAuth(basic)},
		{transport.WithAuth(basic), transport.WithLimiter(employeeLimiter)}} {
		if err := send(opts...); err != nil {
			t.Fatal(err)
		}
	}
	// Оба бюджета исчерпаны: запросы встают в ожидание, а не отправляются.
	for _, opts := range [][]transport.RequestOption{nil, {transport.WithAuth(basic), transport.WithLimiter(employeeLimiter)}} {
		ctx, cancel := context.WithCancel(t.Context())
		done := make(chan error, 1)
		go func() {
			_, err := client.Send(ctx, http.MethodGet, "entity/product", nil, nil, nil, opts...)
			done <- err
		}()
		<-clock.waits
		cancel()
		if err := <-done; !errors.Is(err, context.Canceled) {
			t.Fatalf("err = %v", err)
		}
	}
	if calls.Load() != 3 {
		t.Fatalf("calls = %d", calls.Load())
	}
}
