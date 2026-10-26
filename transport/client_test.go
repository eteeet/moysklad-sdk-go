package transport

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirectToOtherOriginDropsAuthorization(t *testing.T) {
	var got string
	foreign := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Authorization")
	}))
	defer foreign.Close()

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
