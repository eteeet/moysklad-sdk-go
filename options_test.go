package moysklad_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/eteeet/moysklad-sdk-go"
)

// roundTripperFunc подменяет транспорт HTTP-клиента потребителя, как это делает otelhttp.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestWithHTTPClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, r, http.StatusOK, fixture(t, "product.json"))
	}))
	t.Cleanup(srv.Close)
	var viaConsumer atomic.Int32
	hc := &http.Client{Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		viaConsumer.Add(1)
		return http.DefaultTransport.RoundTrip(r)
	})}
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL), moysklad.WithHTTPClient(hc))

	_, err := c.Product.Get(t.Context(), productID)

	require.NoError(t, err)
	require.EqualValues(t, 1, viaConsumer.Load())
}

// fakeLimiter — ограничитель потребителя: считает разрешения и запоминает паузы 429.
type fakeLimiter struct {
	mu       sync.Mutex
	err      error
	acquired int
	released int
	backoffs []time.Duration
}

func (l *fakeLimiter) Acquire(context.Context) (func(), error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.err != nil {
		return nil, l.err
	}
	l.acquired++
	return func() {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.released++
	}, nil
}

func (l *fakeLimiter) Backoff(d time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.backoffs = append(l.backoffs, d)
}

func TestWithLimiter_EveryAttemptPassesLimiter(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		writeJSON(t, w, r, http.StatusOK, fixture(t, "product.json"))
	}))
	t.Cleanup(srv.Close)
	lim := &fakeLimiter{}
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL),
		moysklad.WithLimiter(lim), moysklad.WithRetry(3, time.Millisecond))

	_, err := c.Product.Get(t.Context(), productID)

	require.NoError(t, err)
	require.EqualValues(t, 2, calls.Load())
	require.Equal(t, 2, lim.acquired)
	require.Equal(t, 2, lim.released)
}

func TestWithLimiter_RefusalSkipsRequest(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		calls.Add(1)
	}))
	t.Cleanup(srv.Close)
	refused := errors.New("бюджет аккаунта исчерпан")
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL), moysklad.WithLimiter(&fakeLimiter{err: refused}))

	_, err := c.Product.Get(t.Context(), productID)

	require.ErrorIs(t, err, refused)
	require.Zero(t, calls.Load())
}

func TestWithLimiter_429GoesToBackoff(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			tooManyRequests(t, w, r, "30000")
			return
		}
		writeJSON(t, w, r, http.StatusOK, fixture(t, "product.json"))
	}))
	t.Cleanup(srv.Close)
	lim := &fakeLimiter{}
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL), moysklad.WithLimiter(lim))

	start := time.Now()
	_, err := c.Product.Get(t.Context(), productID)

	require.NoError(t, err)
	require.Less(t, time.Since(start), time.Second, "паузу 429 выдерживает ограничитель, а не SDK")
	require.Equal(t, []time.Duration{30 * time.Second}, lim.backoffs)
	require.EqualValues(t, 2, calls.Load())
}
