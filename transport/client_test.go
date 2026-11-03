package transport

import (
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
)

func TestRedirectToOtherOriginDropsAuthorization(t *testing.T) {
	var got string
	foreign := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Authorization")
	}))
	defer foreign.Close()

	origin := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, foreign.URL, http.StatusFound)
	}))
	defer origin.Close()

	client, err := NewClient(Config{BaseURL: origin.URL})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, origin.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer secret")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if got != "" {
		t.Fatalf("Authorization leaked to another origin: %q", got)
	}
}

// Редирект 307/308 повторил бы запись: transport возвращает его как ответ, а GET следует ему.
func TestWriteRedirectNotFollowed(t *testing.T) {
	for _, tc := range []struct {
		method string
		status int
		calls  int32
	}{
		{http.MethodPost, http.StatusPermanentRedirect, 1},
		{http.MethodPut, http.StatusTemporaryRedirect, 1},
		{http.MethodGet, http.StatusPermanentRedirect, 2},
	} {
		var calls atomic.Int32
		server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if calls.Add(1) == 1 {
				http.Redirect(w, r, "/moved", tc.status)
			}
		}))
		client, err := NewClient(Config{BaseURL: server.URL})
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Send(t.Context(), tc.method, "security/token", nil, nil, nil)
		if calls.Load() != tc.calls || (err == nil) != (tc.calls == 2) {
			t.Errorf("%s %d: calls = %d, err = %v", tc.method, tc.status, calls.Load(), err)
		}
		server.Close()
	}
}

// Href — адрес от BaseURL как в meta.href; слэш в конце BaseURL не удваивается.
func TestHref(t *testing.T) {
	for base, want := range map[string]string{
		"":                                       "https://api.moysklad.ru/api/remap/1.2/entity/product/x",
		"https://api.moysklad.ru/api/remap/1.2/": "https://api.moysklad.ru/api/remap/1.2/entity/product/x",
		"http://127.0.0.1:8080":                  "http://127.0.0.1:8080/entity/product/x",
	} {
		c, err := NewClient(Config{BaseURL: base})
		if err != nil {
			t.Fatal(err)
		}
		if got := c.Href("entity/product/x"); got != want {
			t.Errorf("%q: %s", base, got)
		}
	}
}
