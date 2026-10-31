package uom_test

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/uom"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const fixtureUOMID = "77777777-7777-7777-7777-777777777777"

func TestServiceGet(t *testing.T) {
	for _, tc := range []struct {
		name       string
		fixture    string
		id         string
		wantName   string
		wantCode   string
		predefined bool
		gzip       bool
	}{
		{"upstream/plain", "../../testdata/upstream/uom.json", fixtureUOMID, "шт", "796", false, false},
		{"upstream/gzip", "../../testdata/upstream/uom.json", fixtureUOMID, "шт", "796", false, true},
		{"predefined/plain", "../../testdata/custom/uom/predefined.json", "061721df-9197-49a5-b637-7f5b4d3be969", "дюйм", "039", true, false},
		{"predefined/gzip", "../../testdata/custom/uom/predefined.json", "061721df-9197-49a5-b637-7f5b4d3be969", "дюйм", "039", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture, err := os.ReadFile(tc.fixture)
			if err != nil {
				t.Fatal(err)
			}
			var calls atomic.Int32
			service := newUOMService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodGet ||
					r.URL.Path != "/api/remap/1.2/entity/uom/"+tc.id || r.URL.RawQuery != "" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				if r.Header.Get("Authorization") != "Bearer test-token" ||
					r.Header.Get("Accept-Encoding") != "gzip" ||
					r.Header.Get("Accept") != "application/json;charset=utf-8" ||
					r.Header.Get("User-Agent") != "moysklad-sdk-go" {
					t.Error("required request headers missing")
				}
				if r.Header.Get("Content-Type") != "" {
					t.Error("Content-Type set for GET without body")
				}
				body, err := io.ReadAll(r.Body)
				if err != nil || len(body) != 0 {
					t.Errorf("GET body = %q, error = %v", body, err)
				}

				w.Header().Set("Content-Type", "application/json")
				if tc.gzip {
					w.Header().Set("Content-Encoding", "gzip")
					writer := gzip.NewWriter(w)
					if _, err := writer.Write(fixture); err != nil {
						t.Error(err)
					}
					if err := writer.Close(); err != nil {
						t.Error(err)
					}
					return
				}
				if _, err := w.Write(fixture); err != nil {
					t.Error(err)
				}
			}))

			got, err := service.Get(t.Context(), tc.id)
			if err != nil {
				t.Fatal(err)
			}
			if got.ID != tc.id || got.Meta.Type != "uom" ||
				got.Name != tc.wantName || got.Code != tc.wantCode || got.Shared {
				t.Fatalf("uom = %+v", got)
			}
			if tc.predefined {
				if got.AccountID != "" || got.Owner != nil || got.Group != nil {
					t.Fatalf("absent fields were invented: %+v", got)
				}
			} else if got.AccountID != "abcdef12-1234-1234-1234-123456789012" ||
				got.Owner == nil || got.Owner.Meta.Type != "employee" ||
				got.Group == nil || got.Group.Meta.Type != "group" {
				t.Fatalf("account or references = %+v", got)
			}
			if calls.Load() != 1 {
				t.Fatalf("calls = %d, want 1", calls.Load())
			}
		})
	}
}

func TestServiceGetErrors(t *testing.T) {
	for _, tc := range []struct {
		name    string
		status  int
		body    string
		wantAPI bool
	}{
		{"not found", 404, `{"errors":[{"error":"not found","code":1001}]}`, true},
		{"rate limit", 429, `{"errors":[{"error":"slow down","code":1073}]}`, true},
		{"HTML error", 502, "<html>bad gateway</html>", true},
		{"unexpected status", 201, `{"id":"unexpected"}`, false},
		{"no content", 204, "", false},
		{"missing id", 200, `{"name":"invalid"}`, false},
		{"empty id", 200, `{"id":""}`, false},
		{"null response", 200, "null", false},
		{"invalid JSON", 200, `{"id":`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			service := newUOMService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				w.WriteHeader(tc.status)
				if _, err := io.WriteString(w, tc.body); err != nil {
					t.Error(err)
				}
			}))

			got, err := service.Get(t.Context(), fixtureUOMID)
			if err == nil || got != nil {
				t.Fatalf("uom = %+v, error = %v", got, err)
			}
			if strings.Count(err.Error(), "moysklad:") != 1 ||
				!strings.Contains(err.Error(), "get uom "+fixtureUOMID) {
				t.Errorf("error context = %v", err)
			}
			var apiErr *sdkerrors.APIError
			if errors.As(err, &apiErr) != tc.wantAPI {
				t.Fatalf("unexpected error type: %v", err)
			}
			if tc.wantAPI && apiErr.StatusCode != tc.status {
				t.Errorf("status = %d, want %d", apiErr.StatusCode, tc.status)
			}
			if tc.status == 404 && (len(apiErr.Errors) != 1 ||
				apiErr.Errors[0].Code != 1001 || apiErr.Errors[0].Message != "not found") {
				t.Errorf("API details = %+v", apiErr.Errors)
			}
			if calls.Load() != 1 {
				t.Fatalf("request retried: calls = %d", calls.Load())
			}
		})
	}
}

func TestServiceGetErrorWrapping(t *testing.T) {
	for _, authError := range []bool{true, false} {
		name := "transport"
		if authError {
			name = "authorization"
		}
		t.Run(name, func(t *testing.T) {
			wantErr := errors.New("boundary failed")
			var calls atomic.Int32
			cfg := transport.Config{
				HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
					calls.Add(1)
					return nil, wantErr
				})},
			}
			if authError {
				cfg.Auth = func(*http.Request) error { return wantErr }
			}
			client, err := transport.NewClient(cfg)
			if err != nil {
				t.Fatal(err)
			}
			service, err := uom.NewService(client)
			if err != nil {
				t.Fatal(err)
			}

			got, err := service.Get(t.Context(), fixtureUOMID)
			if got != nil || !errors.Is(err, wantErr) {
				t.Fatalf("uom = %+v, error = %v, want wrapped error", got, err)
			}
			if strings.Count(err.Error(), "moysklad:") != 1 ||
				!strings.HasPrefix(err.Error(), "get uom "+fixtureUOMID+": ") {
				t.Errorf("error context = %v", err)
			}
			wantCalls := int32(1)
			if authError {
				wantCalls = 0
			}
			if calls.Load() != wantCalls {
				t.Errorf("calls = %d, want %d", calls.Load(), wantCalls)
			}
		})
	}
}

func TestServiceGetValidation(t *testing.T) {
	if service, err := uom.NewService(nil); err == nil || service != nil {
		t.Fatal("nil client accepted")
	}
	var calls atomic.Int32
	service := newUOMService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	for _, id := range []string{
		"", "metadata", ".", "..", "../metadata", "x/y", fixtureUOMID + "/metadata",
		"77777777-7777-7777-7777-77777777777g", "777777777777777777777777777777777777",
	} {
		if got, err := service.Get(t.Context(), id); err == nil || got != nil {
			t.Fatalf("invalid id %q accepted", id)
		}
	}
	if got, err := service.Get(nil, fixtureUOMID); err == nil || got != nil {
		t.Fatal("nil context accepted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if got, err := service.Get(ctx, fixtureUOMID); got != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("uom = %+v, error = %v, want cancellation", got, err)
	}
	if calls.Load() != 0 {
		t.Fatalf("unexpected HTTP calls: %d", calls.Load())
	}
}

func TestServiceGetUUIDVersion(t *testing.T) {
	const id = "00000000-0000-f000-0000-000000000000"
	service := newUOMService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/remap/1.2/entity/uom/"+id {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"id":"`+id+`"}`)
	}))
	got, err := service.Get(t.Context(), id)
	if err != nil || got == nil || got.ID != id {
		t.Fatalf("uom = %+v, error = %v", got, err)
	}
}

func newUOMService(t *testing.T, handler http.Handler) *uom.Service {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	auth, err := transport.BearerToken("test-token")
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.NewClient(transport.Config{
		BaseURL: server.URL + "/api/remap/1.2",
		Timeout: 5 * time.Second,
		Auth:    auth,
	})
	if err != nil {
		t.Fatal(err)
	}
	service, err := uom.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	return service
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
