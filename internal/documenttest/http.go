// Package documenttest содержит JSON helpers и проверки HTTP-контракта документов.
package documenttest

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type ResponseCase struct {
	Name      string
	Status    int
	Body      string
	Gzip      bool
	API       bool
	Success   bool
	ErrorCode int
}

func ResponseCases(success string) []ResponseCase {
	return []ResponseCase{
		{Name: "plain", Status: 200, Body: success, Success: true},
		{Name: "gzip", Status: 200, Body: success, Gzip: true, Success: true},
		{Name: "bad request", Status: 400, Body: `{"errors":[{"error":"constraint","code":22004}]}`, API: true, ErrorCode: 22004},
		{Name: "not found", Status: 404, Body: `{"errors":[{"error":"missing","code":1001}]}`, API: true},
		{Name: "rate limit", Status: 429, Body: `{"errors":[{"error":"slow","code":1073}]}`, API: true},
		{Name: "server error", Status: 500, Body: `{"errors":[{"error":"server","code":1000}]}`, API: true},
		{Name: "html error", Status: 502, Body: `<html>gateway</html>`, API: true},
		{Name: "unexpected status", Status: 201, Body: success},
		{Name: "no content", Status: 204},
		{Name: "malformed", Status: 200, Body: `{"id":`},
		{Name: "null", Status: 200, Body: `null`},
		{Name: "empty body", Status: 200},
		{Name: "empty id", Status: 200, Body: `{"id":""}`},
		{Name: "missing id", Status: 200, Body: `{"name":"missing"}`},
	}
}

// HTTP проверяет одну операцию. Тела ошибок синтетические и не доказывают правила живого API.
func HTTP[T any](t *testing.T, tc ResponseCase, method, path string, body []byte,
	invoke func(context.Context, *transport.Client) (*T, error), check func(*T)) {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != method || r.URL.Path != "/api/remap/1.2/"+strings.Split(path, "?")[0] {
			t.Errorf("request = %s %s", r.Method, r.URL)
		}
		query := ""
		if _, q, ok := strings.Cut(path, "?"); ok {
			query = q
		}
		if r.URL.RawQuery != query {
			t.Errorf("query = %q, want %q", r.URL.RawQuery, query)
		}
		for key, want := range map[string]string{"Authorization": "Bearer test-token", "Accept": "application/json;charset=utf-8", "Accept-Encoding": "gzip", "User-Agent": "moysklad-sdk-go"} {
			if r.Header.Get(key) != want {
				t.Errorf("header %s = %q", key, r.Header.Get(key))
			}
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		if method == http.MethodGet {
			if len(data) != 0 || r.Header.Get("Content-Type") != "" {
				t.Errorf("GET body=%s content-type=%q", data, r.Header.Get("Content-Type"))
			}
		} else {
			JSONEqual(t, data, body)
			if r.Header.Get("Content-Type") != "application/json" {
				t.Error("missing JSON content type")
			}
		}
		if tc.Gzip {
			w.Header().Set("Content-Encoding", "gzip")
		}
		w.WriteHeader(tc.Status)
		if tc.Gzip {
			zw := gzip.NewWriter(w)
			if _, err := zw.Write([]byte(tc.Body)); err != nil {
				t.Error(err)
			}
			if err := zw.Close(); err != nil {
				t.Error(err)
			}
		} else {
			if _, err := io.WriteString(w, tc.Body); err != nil {
				t.Error(err)
			}
		}
	}))
	defer server.Close()
	auth, err := transport.BearerToken("test-token")
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", HTTPClient: server.Client(), Auth: auth})
	if err != nil {
		t.Fatal(err)
	}
	got, err := invoke(t.Context(), client)
	if tc.Success {
		if err != nil || got == nil {
			t.Fatalf("result = %+v, error = %v", got, err)
		}
		check(got)
	} else {
		if err == nil || got != nil {
			t.Fatalf("result = %+v, error = %v", got, err)
		}
		if strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("error context = %v", err)
		}
		operation := strings.ToLower(method)
		if method == http.MethodPost {
			operation = "create"
		}
		if method == http.MethodPut {
			operation = "update"
		}
		entity := strings.Split(path, "/")[1]
		if !strings.Contains(err.Error(), operation+" "+entity) {
			t.Errorf("operation context = %v", err)
		}
		var apiErr *sdkerrors.APIError
		if errors.As(err, &apiErr) != tc.API {
			t.Fatalf("API error = %v", err)
		}
		if tc.API && apiErr.StatusCode != tc.Status {
			t.Errorf("API status = %d", apiErr.StatusCode)
		}
		if tc.ErrorCode != 0 && (len(apiErr.Errors) != 1 || apiErr.Errors[0].Code != tc.ErrorCode) {
			t.Errorf("API details = %+v", apiErr)
		}
	}
	// Сбой соединения тоже даёт ошибку и проходит проверки выше; текст ошибки отличает его.
	if calls.Load() != 1 {
		t.Errorf("calls = %d, want 1; error = %v", calls.Load(), err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func LocalErrors[T any](t *testing.T, invoke func(context.Context, *transport.Client) (*T, error)) {
	t.Helper()
	sentinel := errors.New("sentinel")
	for _, name := range []string{"nil context", "cancelled context", "auth", "transport"} {
		t.Run(name, func(t *testing.T) {
			var calls atomic.Int32
			cfg := transport.Config{HTTPClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				calls.Add(1)
				if err := r.Context().Err(); err != nil {
					return nil, err
				}
				return nil, sentinel
			})}}
			ctx := t.Context()
			var want error
			switch name {
			case "nil context":
				ctx = nil
			case "cancelled context":
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(t.Context())
				cancel()
				want = context.Canceled
			case "auth":
				cfg.Auth = func(*http.Request) error { return sentinel }
				want = sentinel
			case "transport":
				want = sentinel
			}
			client, err := transport.NewClient(cfg)
			if err != nil {
				t.Fatal(err)
			}
			got, err := invoke(ctx, client)
			if err == nil || got != nil {
				t.Fatalf("result=%+v err=%v", got, err)
			}
			if want != nil && !errors.Is(err, want) {
				t.Errorf("lost cause: %v", err)
			}
			if strings.Count(err.Error(), "moysklad:") != 1 {
				t.Errorf("error context = %v", err)
			}
			if (name == "nil context" || name == "auth") && calls.Load() != 0 {
				t.Errorf("unexpected HTTP: %d", calls.Load())
			}
			if calls.Load() > 1 {
				t.Error("request retried")
			}
		})
	}
}
