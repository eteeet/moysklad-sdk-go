package product_test

import (
	"compress/gzip"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/product"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const fixtureProductID = "12345678-1234-1234-1234-123456789012"

func TestServiceGet(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/upstream/product.json")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name string
		gzip bool
	}{
		{"plain", false},
		{"gzip", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := newProductService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet ||
					r.URL.Path != "/api/remap/1.2/entity/product/"+fixtureProductID ||
					r.URL.RawQuery != "" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				if r.Header.Get("Authorization") != "Bearer test-token" ||
					r.Header.Get("Accept-Encoding") != "gzip" ||
					r.Header.Get("Accept") != "application/json;charset=utf-8" {
					t.Error("required request headers missing")
				}
				if r.Header.Get("Content-Type") != "" {
					t.Error("Content-Type set for GET without body")
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

			got, err := service.Get(context.Background(), fixtureProductID)
			if err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureProductID || got.Meta.Type != "product" ||
				got.Name != "Тестовый товар" || got.Code != "00001" {
				t.Fatalf("product = %+v", got)
			}
			if got.Archived || !got.Shared || got.VAT != 20 || !got.VATEnabled ||
				got.Weight != 0.5 || got.Volume != 0.001 {
				t.Fatalf("product fields = %+v", got)
			}
			if got.Group == nil || got.Group.Meta.Type != "group" ||
				got.Owner != nil || got.Country != nil ||
				got.Description != "" || got.Article != "" {
				t.Fatalf("references or nullable fields = %+v", got)
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
		{"HTML error", 502, "<html>bad gateway</html>", true},
		{"no content", 204, "", false},
		{"missing id", 200, `{"name":"invalid"}`, false},
		{"null response", 200, "null", false},
		{"invalid JSON", 200, `{"id":`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := newProductService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				if _, err := w.Write([]byte(tc.body)); err != nil {
					t.Error(err)
				}
			}))

			got, err := service.Get(context.Background(), fixtureProductID)
			if err == nil || got != nil {
				t.Fatalf("product = %+v, error = %v", got, err)
			}
			var apiErr *sdkerrors.APIError
			if errors.As(err, &apiErr) != tc.wantAPI {
				t.Fatalf("unexpected error type: %v", err)
			}
			if tc.wantAPI && apiErr.StatusCode != tc.status {
				t.Fatalf("status = %d, want %d", apiErr.StatusCode, tc.status)
			}
		})
	}
}

func TestServiceGetValidation(t *testing.T) {
	if service, err := product.NewService(nil); err == nil || service != nil {
		t.Fatal("nil client accepted")
	}

	var calls atomic.Int32
	service := newProductService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))

	for _, id := range []string{
		"", "metadata", ".", "..", "../metadata", "x/y",
		"12345678-1234-1234-1234-12345678901g",
	} {
		if got, err := service.Get(context.Background(), id); err == nil || got != nil {
			t.Fatalf("invalid id %q accepted", id)
		}
	}

	if _, err := service.Get(nil, fixtureProductID); err == nil {
		t.Fatal("nil context accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := service.Get(ctx, fixtureProductID); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation lost: %v", err)
	}
	if calls.Load() != 0 {
		t.Fatalf("unexpected HTTP calls: %d", calls.Load())
	}
}

func newProductService(t *testing.T, handler http.Handler) *product.Service {
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
	service, err := product.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	return service
}
