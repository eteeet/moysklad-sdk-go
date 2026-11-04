package vendorapi_test

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/contexts"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
	"github.com/eteeet/moysklad-sdk-go/vendorapi"
)

const (
	fixtureID  = "b0a02321-13e3-11e9-912f-f3d4002516e3"
	contextKey = "1c14e98cd272239c03bf3d9697f167699743292c"
)

// Пример документации Vendor API по форме совпадает с context/employee JSON API.
func TestContextFixture(t *testing.T) {
	data := fixture(t, "context.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got contexts.EmployeeContext
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.UID != "admin@bkozhevnikov" || len(got.Permissions) != 65 ||
				got.Permissions["admin"].Actions["view"] != model.ScopeAll {
				t.Fatalf("context = %+v", got)
			}
		})
	}
}

func TestNewServiceRejectsNil(t *testing.T) {
	if s, err := vendorapi.NewService(nil); err == nil || s != nil {
		t.Fatalf("NewService(nil) = %v, %v", s, err)
	}
}

// vendorServer — Vendor API на memnet; клиент подписывает запросы AppAuth и повторяет чтения.
func vendorServer(t *testing.T, handler http.HandlerFunc) *vendorapi.Service {
	t.Helper()
	server := memnet.NewServer(handler)
	t.Cleanup(server.Close)
	auth, err := vendorapi.AppAuth("app", "secret-key")
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.NewClient(transport.Config{
		BaseURL: server.URL + "/api/vendor/1.0", Auth: auth,
		Retry: &transport.RetryPolicy{MaxAttempts: 3, BaseDelay: 1, MaxDelay: 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	s, err := vendorapi.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// respond пишет ответ случая ResponseCases с заголовком X-Lognex-Auth для проверки ResponseInfo.
func respond(w http.ResponseWriter, tc documenttest.ResponseCase) {
	w.Header().Set("X-Lognex-Auth", "test")
	if tc.Gzip {
		w.Header().Set("Content-Encoding", "gzip")
	}
	w.WriteHeader(tc.Status)
	out := io.Writer(w)
	if tc.Gzip {
		zw := gzip.NewWriter(w)
		defer zw.Close()
		out = zw
	}
	io.WriteString(out, tc.Body)
}

// POST без тела на context/{contextKey}; POST не повторяется даже при RetryPolicy и 5xx.
func TestServiceContext(t *testing.T) {
	for _, tc := range documenttest.ResponseCases(string(fixture(t, "context.json"))) {
		t.Run(tc.Name, func(t *testing.T) {
			var calls atomic.Int32
			s := vendorServer(t, func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodPost || r.URL.RequestURI() != "/api/vendor/1.0/context/"+contextKey {
					t.Errorf("request = %s %s", r.Method, r.URL.RequestURI())
				}
				outgoingClaims(t, r.Header.Get("Authorization"), "secret-key")
				if r.Header.Get("Accept-Encoding") != "gzip" || r.Header.Get("Content-Type") != "" || r.ContentLength != 0 {
					t.Errorf("headers = %v, length = %d", r.Header, r.ContentLength)
				}
				if data, _ := io.ReadAll(r.Body); len(data) != 0 {
					t.Errorf("body = %s", data)
				}
				respond(w, tc)
			})
			var info transport.ResponseInfo
			got, err := s.Context(t.Context(), contextKey, transport.WithResponseInfo(&info))
			if calls.Load() != 1 {
				t.Fatalf("calls = %d", calls.Load())
			}
			if info.StatusCode != tc.Status || info.Header.Get("X-Lognex-Auth") != "test" {
				t.Errorf("info = %+v", info)
			}
			if tc.Success {
				if err != nil || got == nil || got.ID != fixtureID || len(got.Permissions) != 65 {
					t.Fatalf("got = %+v, err = %v", got, err)
				}
				return
			}
			if err == nil || got != nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "vendor context") {
				t.Fatalf("got = %+v, err = %v", got, err)
			}
			var apiErr *sdkerrors.APIError
			if errors.As(err, &apiErr) != tc.API || tc.API && apiErr.StatusCode != tc.Status {
				t.Fatalf("API error = %v", err)
			}
		})
	}
}

// 403 и 404 документированы для этого endpoint: решение без доступа и истёкший contextKey.
func TestServiceContextDocumentedErrors(t *testing.T) {
	for _, status := range []int{http.StatusForbidden, http.StatusNotFound} {
		s := vendorServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			io.WriteString(w, `{"errors":[{"error":"denied","code":1001}]}`)
		})
		got, err := s.Context(t.Context(), contextKey)
		if got != nil || sdkerrors.StatusCode(err) != status {
			t.Errorf("status %d: got = %v, err = %v", status, got, err)
		}
	}
}

// Ошибка ключа — до HTTP и обнуляет ResponseInfo прошлого вызова.
func TestServiceContextRejectsKey(t *testing.T) {
	var calls atomic.Int32
	s := vendorServer(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1) })
	// "." и ".." после JoinPath указывали бы на context или корень API, "user" — на обмен токена.
	for _, key := range []string{"", "a/b", "a?b", "a#b", ".", "..", "user"} {
		info := transport.ResponseInfo{StatusCode: 200, Header: http.Header{"X": {"stale"}}}
		got, err := s.Context(t.Context(), key, transport.WithResponseInfo(&info))
		if err == nil || got != nil || info.StatusCode != 0 || info.Header != nil {
			t.Errorf("key %q: got = %v, err = %v, info = %+v", key, got, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Errorf("calls = %d", calls.Load())
	}
}

func TestServiceContextLocalErrors(t *testing.T) {
	documenttest.LocalErrors(t, func(ctx context.Context, client *transport.Client) (*contexts.EmployeeContext, error) {
		s, err := vendorapi.NewService(client)
		if err != nil {
			return nil, err
		}
		return s.Context(ctx, contextKey)
	})
}
