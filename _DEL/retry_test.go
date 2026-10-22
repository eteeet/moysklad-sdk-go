package moysklad_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/eteeet/moysklad-sdk-go"
)

// tooManyRequests отвечает как МойСклад при превышении частоты запросов.
func tooManyRequests(t *testing.T, w http.ResponseWriter, r *http.Request, retryAfter string) {
	t.Helper()
	w.Header().Set("X-Lognex-Retry-After", retryAfter)
	writeJSON(t, w, r, http.StatusTooManyRequests,
		[]byte(`{"errors":[{"error":"Превышено ограничение на количество запросов в единицу времени","code":1049}]}`))
}

func TestGet_RetriesAfter429(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			tooManyRequests(t, w, r, "50")
			return
		}
		writeJSON(t, w, r, http.StatusOK, fixture(t, "product.json"))
	}))
	t.Cleanup(srv.Close)
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL))

	start := time.Now()
	p, err := c.Product.Get(t.Context(), productID)

	require.NoError(t, err)
	require.Equal(t, "Тестовый товар", *p.Name)
	require.EqualValues(t, 2, calls.Load())
	require.GreaterOrEqual(t, time.Since(start), 50*time.Millisecond)
}

func TestGet_RetriesOn5xx(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		writeJSON(t, w, r, http.StatusOK, fixture(t, "product.json"))
	}))
	t.Cleanup(srv.Close)
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL), moysklad.WithRetry(3, time.Millisecond))

	p, err := c.Product.Get(t.Context(), productID)

	require.NoError(t, err)
	require.Equal(t, "Тестовый товар", *p.Name)
	require.EqualValues(t, 3, calls.Load())
}

func TestDo_PostIsNotRetriedOn5xx(t *testing.T) {
	var calls atomic.Int32
	var gotBody, gotContentType string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		b, _ := io.ReadAll(r.Body)
		gotBody, gotContentType = string(b), r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL), moysklad.WithRetry(3, time.Millisecond))

	err := c.Do(t.Context(), http.MethodPost, "entity/product", nil, map[string]string{"name": "Новый товар"}, nil)

	var msErr *moysklad.Error
	require.ErrorAs(t, err, &msErr)
	require.Equal(t, http.StatusInternalServerError, msErr.StatusCode)
	require.EqualValues(t, 1, calls.Load())
	require.JSONEq(t, `{"name":"Новый товар"}`, gotBody)
	require.Equal(t, "application/json", gotContentType)
}

func TestDo_OnlyGetPutDeleteAreRetriedOn5xx(t *testing.T) {
	tests := map[string]int32{
		http.MethodGet:    2,
		http.MethodPut:    2,
		http.MethodDelete: 2,
		http.MethodPost:   1,
		http.MethodPatch:  1,
		"post":            1,
	}
	for method, wantCalls := range tests {
		t.Run(method, func(t *testing.T) {
			var calls atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				calls.Add(1)
				w.WriteHeader(http.StatusServiceUnavailable)
			}))
			t.Cleanup(srv.Close)
			c := moysklad.New("token", moysklad.WithBaseURL(srv.URL), moysklad.WithRetry(2, time.Millisecond))

			err := c.Do(t.Context(), method, "entity/product", nil, nil, nil)

			require.Error(t, err)
			require.Equal(t, wantCalls, calls.Load())
		})
	}
}

func TestDo_PutIsRetriedWithSameBody(t *testing.T) {
	var bodies []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		bodies = append(bodies, string(b))
		if len(bodies) == 1 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		writeJSON(t, w, r, http.StatusOK, fixture(t, "product.json"))
	}))
	t.Cleanup(srv.Close)
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL), moysklad.WithRetry(3, time.Millisecond))

	var p moysklad.Product
	err := c.Do(t.Context(), http.MethodPut, "entity/product/"+productID, nil, map[string]string{"name": "Тестовый товар"}, &p)

	require.NoError(t, err)
	require.Equal(t, "Тестовый товар", *p.Name)
	require.Len(t, bodies, 2)
	require.JSONEq(t, `{"name":"Тестовый товар"}`, bodies[1])
}

func TestGet_ContextCancelsRetryWait(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		tooManyRequests(t, w, r, "10000")
	}))
	t.Cleanup(srv.Close)
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL))
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	t.Cleanup(cancel)

	start := time.Now()
	_, err := c.Product.Get(ctx, productID)

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Less(t, time.Since(start), time.Second)
	require.EqualValues(t, 1, calls.Load())
}

func TestGet_ReturnsLastErrorWhenAttemptsRunOut(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	t.Cleanup(srv.Close)
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL), moysklad.WithRetry(2, time.Millisecond))

	_, err := c.Product.Get(t.Context(), productID)

	var msErr *moysklad.Error
	require.ErrorAs(t, err, &msErr)
	require.Equal(t, http.StatusServiceUnavailable, msErr.StatusCode)
	require.EqualValues(t, 2, calls.Load())
}
