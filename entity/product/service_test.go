package product_test

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

	"github.com/eteeet/moysklad-sdk-go/entity/product"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/model"
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
			if count := strings.Count(err.Error(), "moysklad:"); count != 1 {
				t.Errorf("moysklad prefix count = %d, want 1: %v", count, err)
			}
			if tc.wantAPI && !strings.HasPrefix(err.Error(), "get product "+fixtureProductID+": ") {
				t.Errorf("product context lost: %v", err)
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

func TestServiceGetRequestErrorWrapping(t *testing.T) {
	wantErr := errors.New("auth failed")
	client, err := transport.NewClient(transport.Config{
		Auth: func(*http.Request) error { return wantErr },
	})
	if err != nil {
		t.Fatal(err)
	}
	service, err := product.NewService(client)
	if err != nil {
		t.Fatal(err)
	}

	got, err := service.Get(t.Context(), fixtureProductID)
	if got != nil || !errors.Is(err, wantErr) {
		t.Fatalf("product = %+v, error = %v, want wrapped auth error", got, err)
	}
	if count := strings.Count(err.Error(), "moysklad:"); count != 1 {
		t.Errorf("moysklad prefix count = %d, want 1: %v", count, err)
	}
	if !strings.HasPrefix(err.Error(), "get product "+fixtureProductID+": ") {
		t.Errorf("product context lost: %v", err)
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

func TestServiceCreate(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/upstream/product_full.json")
	if err != nil {
		t.Fatal(err)
	}
	wantBody := productGolden(t, "create")
	var calls atomic.Int32
	service := newProductService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != http.MethodPost || r.URL.Path != "/api/remap/1.2/entity/product" || r.URL.RawQuery != "" {
			t.Errorf("request = %s %s", r.Method, r.URL)
		}
		if r.Header.Get("Content-Type") != "application/json" || r.Header.Get("Authorization") != "Bearer test-token" || r.Header.Get("Accept") != "application/json;charset=utf-8" || r.Header.Get("Accept-Encoding") != "gzip" {
			t.Error("required headers missing")
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			w.WriteHeader(500)
			return
		}
		assertJSONEqual(t, body, wantBody)
		w.Header().Set("Content-Encoding", "gzip")
		writer := gzip.NewWriter(w)
		if _, err := writer.Write(fixture); err != nil {
			t.Error(err)
		}
		if err := writer.Close(); err != nil {
			t.Error(err)
		}
	}))
	got, err := service.Create(t.Context(), fullCreateRequest(t))
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != fixtureProductID || got.BuyPrice == nil || len(got.Attributes) != 3 {
		t.Fatalf("product = %+v", got)
	}
	if calls.Load() != 1 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestServiceCreateErrors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		api    bool
	}{
		{"bad request", 400, `{"errors":[{"error":"invalid name","code":3000}]}`, true},
		{"rate limit", 429, `{"errors":[{"error":"slow down","code":1073}]}`, true},
		{"gateway", 502, "bad gateway", true},
		{"created status", 201, `{"id":"unexpected"}`, false},
		{"no content", 204, "", false},
		{"missing id", 200, `{"name":"invalid"}`, false},
		{"null", 200, `null`, false},
		{"invalid JSON", 200, `{"id":`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			service := newProductService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			}))
			got, err := service.Create(t.Context(), product.CreateRequest{Name: "test"})
			if err == nil || got != nil {
				t.Fatalf("product = %+v, error = %v", got, err)
			}
			var apiErr *sdkerrors.APIError
			if errors.As(err, &apiErr) != tc.api {
				t.Fatalf("error = %v", err)
			}
			if tc.api && apiErr.StatusCode != tc.status {
				t.Fatalf("status = %d", apiErr.StatusCode)
			}
			if strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "create product") {
				t.Fatalf("error context = %v", err)
			}
			if calls.Load() != 1 {
				t.Fatalf("POST retried: calls = %d", calls.Load())
			}
		})
	}
}

func TestServiceCreateValidation(t *testing.T) {
	var calls atomic.Int32
	service := newProductService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1); w.WriteHeader(500) }))
	if got, err := service.Create(t.Context(), product.CreateRequest{}); err == nil || got != nil {
		t.Fatal("empty name accepted")
	}
	if _, err := service.Create(nil, product.CreateRequest{Name: "test"}); err == nil {
		t.Fatal("nil context accepted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := service.Create(ctx, product.CreateRequest{Name: "test"}); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation lost: %v", err)
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestServiceUpdate(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/upstream/product_full.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		input product.UpdateRequest
		want  []byte
	}{
		{"golden", fullUpdateRequest(), productGolden(t, "update")},
		{"one field", product.UpdateRequest{Name: model.Value("new")}, []byte(`{"name":"new"}`)},
		{"empty collections", product.UpdateRequest{SalePrices: model.Value([]model.SalePrice{}), Attributes: model.Value([]model.Attribute{})}, []byte(`{"salePrices":[],"attributes":[]}`)},
		{"nil collections", product.UpdateRequest{SalePrices: model.Value[[]model.SalePrice](nil)}, []byte(`{"salePrices":[]}`)},
		{"explicit null", product.UpdateRequest{ProductFolder: model.Null[model.Reference]()}, []byte(`{"productFolder":null}`)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			service := newProductService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodPut || r.URL.Path != "/api/remap/1.2/entity/product/"+fixtureProductID || r.URL.RawQuery != "" {
					t.Errorf("request = %s %s", r.Method, r.URL)
				}
				if r.Header.Get("Content-Type") != "application/json" || r.Header.Get("Authorization") != "Bearer test-token" {
					t.Error("required headers missing")
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
					w.WriteHeader(500)
					return
				}
				assertJSONEqual(t, body, tc.want)
				_, _ = w.Write(fixture)
			}))
			got, err := service.Update(t.Context(), fixtureProductID, tc.input)
			if err != nil {
				t.Fatal(err)
			}
			// Возвращаем ответ сервера, не заменяем коллекции переданными массивами.
			if got.ID != fixtureProductID || len(got.SalePrices) != 1 || len(got.Attributes) != 3 {
				t.Fatalf("product = %+v", got)
			}
			if calls.Load() != 1 {
				t.Fatalf("calls = %d", calls.Load())
			}
		})
	}
}

func TestServiceUpdateErrors(t *testing.T) {
	for _, tc := range []struct {
		status int
		body   string
		api    bool
	}{
		{404, `{"errors":[{"error":"missing","code":1001}]}`, true},
		{429, `{"errors":[{"error":"slow down","code":1073}]}`, true},
		{502, "bad gateway", true},
		{201, `{"id":"unexpected"}`, false},
		{204, "", false},
		{200, `{}`, false},
		{200, `null`, false},
		{200, `{"id":`, false},
	} {
		var calls atomic.Int32
		service := newProductService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls.Add(1)
			w.WriteHeader(tc.status)
			_, _ = io.WriteString(w, tc.body)
		}))
		got, err := service.Update(t.Context(), fixtureProductID, product.UpdateRequest{Archived: model.Value(false)})
		if err == nil || got != nil {
			t.Fatalf("product = %+v, error = %v", got, err)
		}
		var apiErr *sdkerrors.APIError
		if errors.As(err, &apiErr) != tc.api {
			t.Fatalf("error = %v", err)
		}
		if tc.api && apiErr.StatusCode != tc.status {
			t.Fatalf("status = %d", apiErr.StatusCode)
		}
		if strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "update product "+fixtureProductID) {
			t.Fatalf("error context = %v", err)
		}
		if calls.Load() != 1 {
			t.Fatalf("calls = %d", calls.Load())
		}
	}
}

func TestServiceUpdateValidation(t *testing.T) {
	var calls atomic.Int32
	service := newProductService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1); w.WriteHeader(500) }))
	input := product.UpdateRequest{Archived: model.Value(false)}
	for _, id := range []string{"", "metadata", "../metadata", fixtureProductID + "/positions", "12345678-1234-1234-1234-12345678901g"} {
		if got, err := service.Update(t.Context(), id, input); err == nil || got != nil {
			t.Fatalf("invalid id %q accepted", id)
		}
	}
	if got, err := service.Update(t.Context(), fixtureProductID, product.UpdateRequest{}); err == nil || got != nil {
		t.Fatal("empty update accepted")
	}
	if _, err := service.Update(nil, fixtureProductID, input); err == nil {
		t.Fatal("nil context accepted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := service.Update(ctx, fixtureProductID, input); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation lost: %v", err)
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestServiceWriteRequestErrorWrapping(t *testing.T) {
	wantErr := errors.New("auth failed")
	client, err := transport.NewClient(transport.Config{Auth: func(*http.Request) error { return wantErr }})
	if err != nil {
		t.Fatal(err)
	}
	service, err := product.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	for _, call := range []func() (*product.Product, error){
		func() (*product.Product, error) {
			return service.Create(t.Context(), product.CreateRequest{Name: "test"})
		},
		func() (*product.Product, error) {
			return service.Update(t.Context(), fixtureProductID, product.UpdateRequest{Archived: model.Value(false)})
		},
	} {
		got, err := call()
		if got != nil || !errors.Is(err, wantErr) || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Fatalf("product = %+v, error = %v", got, err)
		}
	}
}
