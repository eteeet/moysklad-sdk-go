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
	"github.com/eteeet/moysklad-sdk-go/model"
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

				writeUOMResponse(t, w, fixture, tc.gzip)
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

func TestServiceErrorWrapping(t *testing.T) {
	for _, operation := range []struct {
		name    string
		context string
		call    func(*uom.Service, context.Context) (*uom.UOM, error)
	}{
		{"get", "get uom " + fixtureUOMID, func(s *uom.Service, ctx context.Context) (*uom.UOM, error) {
			return s.Get(ctx, fixtureUOMID)
		}},
		{"create", "create uom", func(s *uom.Service, ctx context.Context) (*uom.UOM, error) {
			return s.Create(ctx, uom.CreateRequest{Name: "SDK UOM"})
		}},
		{"update", "update uom " + fixtureUOMID, func(s *uom.Service, ctx context.Context) (*uom.UOM, error) {
			return s.Update(ctx, fixtureUOMID, uom.UpdateRequest{Shared: model.Value(false)})
		}},
	} {
		t.Run(operation.name, func(t *testing.T) {
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
					got, err := operation.call(service, t.Context())
					if got != nil || !errors.Is(err, wantErr) {
						t.Fatalf("uom = %+v, error = %v, want wrapped error", got, err)
					}
					if strings.Count(err.Error(), "moysklad:") != 1 ||
						!strings.HasPrefix(err.Error(), operation.context+": ") {
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

func TestServiceUUIDVersion(t *testing.T) {
	const id = "00000000-0000-f000-0000-000000000000"
	service := newUOMService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/remap/1.2/entity/uom/"+id {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"id":"`+id+`"}`)
	}))
	for _, tc := range []struct {
		name string
		call func(context.Context) (*uom.UOM, error)
	}{
		{"get", func(ctx context.Context) (*uom.UOM, error) { return service.Get(ctx, id) }},
		{"update", func(ctx context.Context) (*uom.UOM, error) {
			return service.Update(ctx, id, uom.UpdateRequest{Shared: model.Value(false)})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.call(t.Context())
			if err != nil || got == nil || got.ID != id {
				t.Fatalf("uom = %+v, error = %v", got, err)
			}
		})
	}
}

func TestServiceCreate(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/upstream/uom.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		input uom.CreateRequest
		want  []byte
		gzip  bool
	}{
		{"full/plain", fullCreateRequest(), uomGolden(t, "create"), false},
		{"full/gzip", fullCreateRequest(), uomGolden(t, "create"), true},
		{"minimal", uom.CreateRequest{Name: "SDK UOM"}, []byte(`{"name":"SDK UOM"}`), false},
		{"whitespace name", uom.CreateRequest{Name: " \t "}, []byte(`{"name":" \t "}`), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			service := newUOMService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodPost || r.URL.Path != "/api/remap/1.2/entity/uom" || r.URL.RawQuery != "" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				if r.Header.Get("Content-Type") != "application/json" ||
					r.Header.Get("Authorization") != "Bearer test-token" ||
					r.Header.Get("Accept") != "application/json;charset=utf-8" ||
					r.Header.Get("Accept-Encoding") != "gzip" ||
					r.Header.Get("User-Agent") != "moysklad-sdk-go" {
					t.Error("required request headers missing")
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				assertJSONEqual(t, body, tc.want)
				writeUOMResponse(t, w, fixture, tc.gzip)
			}))

			got, err := service.Create(t.Context(), tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureUOMID || got.Meta.Type != "uom" ||
				got.Name != "шт" || got.Code != "796" || got.Description != "Штука" ||
				got.ExternalCode != "PCE" || got.Shared || got.Updated.IsZero() ||
				got.Owner == nil || got.Group == nil {
				t.Fatalf("uom = %+v", got)
			}
			if calls.Load() != 1 {
				t.Fatalf("calls = %d, want one POST without a preliminary GET", calls.Load())
			}
		})
	}
}

func TestServiceCreateErrors(t *testing.T) {
	for _, tc := range []struct {
		name    string
		status  int
		body    string
		wantAPI bool
	}{
		{"bad request", 400, `{"errors":[{"error":"invalid name","code":3000}]}`, true},
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
			got, err := service.Create(t.Context(), uom.CreateRequest{Name: "SDK UOM"})
			if err == nil || got != nil {
				t.Fatalf("uom = %+v, error = %v", got, err)
			}
			if strings.Count(err.Error(), "moysklad:") != 1 ||
				!strings.Contains(err.Error(), "create uom") {
				t.Errorf("error context = %v", err)
			}
			var apiErr *sdkerrors.APIError
			if errors.As(err, &apiErr) != tc.wantAPI {
				t.Fatalf("unexpected error type: %v", err)
			}
			if tc.wantAPI && apiErr.StatusCode != tc.status {
				t.Errorf("status = %d, want %d", apiErr.StatusCode, tc.status)
			}
			if tc.status == 400 && (len(apiErr.Errors) != 1 ||
				apiErr.Errors[0].Code != 3000 || apiErr.Errors[0].Message != "invalid name") {
				t.Errorf("API details = %+v", apiErr.Errors)
			}
			if calls.Load() != 1 {
				t.Fatalf("POST retried: calls = %d", calls.Load())
			}
		})
	}
}

func TestServiceCreateValidation(t *testing.T) {
	var calls atomic.Int32
	service := newUOMService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	if got, err := service.Create(t.Context(), uom.CreateRequest{}); err == nil || got != nil {
		t.Fatal("empty name accepted")
	}
	if got, err := service.Create(nil, uom.CreateRequest{Name: "SDK UOM"}); err == nil || got != nil {
		t.Fatal("nil context accepted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if got, err := service.Create(ctx, uom.CreateRequest{Name: "SDK UOM"}); got != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("uom = %+v, error = %v, want cancellation", got, err)
	}
	if calls.Load() != 0 {
		t.Fatalf("unexpected HTTP calls: %d", calls.Load())
	}
}

func TestServiceUpdate(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/upstream/uom.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		input uom.UpdateRequest
		want  []byte
		gzip  bool
	}{
		{"full/plain", fullUpdateRequest(), uomGolden(t, "update"), false},
		{"full/gzip", fullUpdateRequest(), uomGolden(t, "update"), true},
		{"name", uom.UpdateRequest{Name: model.Value("renamed")}, []byte(`{"name":"renamed"}`), false},
		{"empty name", uom.UpdateRequest{Name: model.Value("")}, []byte(`{"name":""}`), false},
		{"code", uom.UpdateRequest{Code: model.Value("039")}, []byte(`{"code":"039"}`), false},
		{"external code", uom.UpdateRequest{ExternalCode: model.Value("updated")}, []byte(`{"externalCode":"updated"}`), false},
		{"empty description", uom.UpdateRequest{Description: model.Value("")}, []byte(`{"description":""}`), false},
		{"shared false", uom.UpdateRequest{Shared: model.Value(false)}, []byte(`{"shared":false}`), false},
		{"owner", uom.UpdateRequest{Owner: model.Value(model.Reference{Meta: model.Meta{
			Href: "https://api.moysklad.ru/api/remap/1.2/entity/employee/11111111-1111-1111-1111-111111111111",
		}})}, []byte(`{"owner":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/employee/11111111-1111-1111-1111-111111111111"}}}`), false},
		{"explicit null", uom.UpdateRequest{Group: model.Null[model.Reference]()}, []byte(`{"group":null}`), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			service := newUOMService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodPut || r.URL.Path != "/api/remap/1.2/entity/uom/"+fixtureUOMID || r.URL.RawQuery != "" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				if r.Header.Get("Content-Type") != "application/json" ||
					r.Header.Get("Authorization") != "Bearer test-token" ||
					r.Header.Get("Accept") != "application/json;charset=utf-8" ||
					r.Header.Get("Accept-Encoding") != "gzip" ||
					r.Header.Get("User-Agent") != "moysklad-sdk-go" {
					t.Error("required request headers missing")
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				assertJSONEqual(t, body, tc.want)
				writeUOMResponse(t, w, fixture, tc.gzip)
			}))

			got, err := service.Update(t.Context(), fixtureUOMID, tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureUOMID || got.Meta.Type != "uom" ||
				got.Name != "шт" || got.Code != "796" || got.Description != "Штука" ||
				got.ExternalCode != "PCE" || got.Shared || got.Updated.IsZero() ||
				got.Owner == nil || got.Group == nil {
				t.Fatalf("uom = %+v", got)
			}
			if calls.Load() != 1 {
				t.Fatalf("calls = %d, want one PUT without a preliminary GET", calls.Load())
			}
		})
	}
}

// Отказ для предустановленной UOM моделируется локально; статус и текст не взяты из живого API.
func TestServiceUpdateErrors(t *testing.T) {
	for _, tc := range []struct {
		name    string
		id      string
		status  int
		body    string
		wantAPI bool
	}{
		{"predefined uom", "061721df-9197-49a5-b637-7f5b4d3be969", 403, `{"errors":[{"error":"predefined uom cannot be changed"}]}`, true},
		{"not found", fixtureUOMID, 404, `{"errors":[{"error":"not found","code":1001}]}`, true},
		{"rate limit", fixtureUOMID, 429, `{"errors":[{"error":"slow down","code":1073}]}`, true},
		{"HTML error", fixtureUOMID, 502, "<html>bad gateway</html>", true},
		{"unexpected status", fixtureUOMID, 201, `{"id":"unexpected"}`, false},
		{"no content", fixtureUOMID, 204, "", false},
		{"missing id", fixtureUOMID, 200, `{"name":"invalid"}`, false},
		{"empty id", fixtureUOMID, 200, `{"id":""}`, false},
		{"null response", fixtureUOMID, 200, "null", false},
		{"invalid JSON", fixtureUOMID, 200, `{"id":`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			service := newUOMService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodPut || r.URL.Path != "/api/remap/1.2/entity/uom/"+tc.id {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				w.WriteHeader(tc.status)
				if _, err := io.WriteString(w, tc.body); err != nil {
					t.Error(err)
				}
			}))
			got, err := service.Update(t.Context(), tc.id, uom.UpdateRequest{Shared: model.Value(false)})
			if err == nil || got != nil {
				t.Fatalf("uom = %+v, error = %v", got, err)
			}
			if strings.Count(err.Error(), "moysklad:") != 1 ||
				!strings.Contains(err.Error(), "update uom "+tc.id) {
				t.Errorf("error context = %v", err)
			}
			var apiErr *sdkerrors.APIError
			if errors.As(err, &apiErr) != tc.wantAPI {
				t.Fatalf("unexpected error type: %v", err)
			}
			if tc.wantAPI && apiErr.StatusCode != tc.status {
				t.Errorf("status = %d, want %d", apiErr.StatusCode, tc.status)
			}
			if tc.name == "predefined uom" && (len(apiErr.Errors) != 1 ||
				apiErr.Errors[0].Message != "predefined uom cannot be changed") {
				t.Errorf("API details = %+v", apiErr.Errors)
			}
			if calls.Load() != 1 {
				t.Fatalf("request retried: calls = %d", calls.Load())
			}
		})
	}
}

func TestServiceUpdateValidation(t *testing.T) {
	var calls atomic.Int32
	service := newUOMService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	input := uom.UpdateRequest{Shared: model.Value(false)}
	for _, id := range []string{
		"", "metadata", ".", "..", "../metadata", "x/y", fixtureUOMID + "/metadata",
		"77777777-7777-7777-7777-77777777777g", "777777777777777777777777777777777777",
	} {
		if got, err := service.Update(t.Context(), id, input); err == nil || got != nil {
			t.Fatalf("invalid id %q accepted", id)
		}
	}
	if got, err := service.Update(t.Context(), fixtureUOMID, uom.UpdateRequest{}); err == nil || got != nil {
		t.Fatal("empty update accepted")
	}
	if got, err := service.Update(nil, fixtureUOMID, input); err == nil || got != nil {
		t.Fatal("nil context accepted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if got, err := service.Update(ctx, fixtureUOMID, input); got != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("uom = %+v, error = %v, want cancellation", got, err)
	}
	if calls.Load() != 0 {
		t.Fatalf("unexpected HTTP calls: %d", calls.Load())
	}
}

func writeUOMResponse(t *testing.T, w http.ResponseWriter, fixture []byte, compressed bool) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if compressed {
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

func TestServiceRequestOptions(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/upstream/uom.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name   string
		status int
		invoke func(*uom.Service, ...transport.RequestOption) (*uom.UOM, error)
	}{
		{"get", http.StatusOK, func(s *uom.Service, opts ...transport.RequestOption) (*uom.UOM, error) {
			return s.Get(t.Context(), fixtureUOMID, opts...)
		}},
		{"create", http.StatusOK, func(s *uom.Service, opts ...transport.RequestOption) (*uom.UOM, error) {
			return s.Create(t.Context(), uom.CreateRequest{Name: "шт"}, opts...)
		}},
		{"update", http.StatusOK, func(s *uom.Service, opts ...transport.RequestOption) (*uom.UOM, error) {
			return s.Update(t.Context(), fixtureUOMID, uom.UpdateRequest{Code: model.Value("796")}, opts...)
		}},
		{"update api error", http.StatusConflict, func(s *uom.Service, opts ...transport.RequestOption) (*uom.UOM, error) {
			return s.Update(t.Context(), fixtureUOMID, uom.UpdateRequest{Code: model.Value("796")}, opts...)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			service := newUOMService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Header.Get("Accept-Language") != "en" || r.Header.Get("Authorization") != "Bearer test-token" {
					t.Errorf("headers = %v", r.Header)
				}
				w.Header().Set("X-Lognex-Content-Region", "kz")
				if tc.status != http.StatusOK {
					w.WriteHeader(tc.status)
					w.Write([]byte(`{"errors":[{"error":"conflict","code":1065}]}`))
					return
				}
				writeUOMResponse(t, w, fixture, false)
			}))

			var info transport.ResponseInfo
			got, err := tc.invoke(service, transport.WithHeader("Accept-Language", "en"), transport.WithResponseInfo(&info))
			if tc.status == http.StatusOK {
				if err != nil || got.ID != fixtureUOMID {
					t.Fatalf("got = %+v, err = %v", got, err)
				}
			} else {
				var apiErr *sdkerrors.APIError
				if got != nil || !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusConflict {
					t.Fatalf("got = %+v, err = %v", got, err)
				}
			}
			if info.StatusCode != tc.status || info.Header.Get("X-Lognex-Content-Region") != "kz" || calls.Load() != 1 {
				t.Fatalf("info = %+v, calls = %d", info, calls.Load())
			}
		})
	}
}
