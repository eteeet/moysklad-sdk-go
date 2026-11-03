package country_test

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/country"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const fixtureCountryID = "44444444-4444-4444-4444-444444444444"
const predefinedCountryID = "000d77a9-3000-4f81-a995-6b9cffdee1d2"

func TestServiceGet(t *testing.T) {
	for _, tc := range []struct {
		name, fixture, id, wantName, wantCode string
		predefined, compressed                bool
	}{
		{"upstream/plain", "../../testdata/upstream/country.json", fixtureCountryID, "Россия", "643", false, false},
		{"upstream/gzip", "../../testdata/upstream/country.json", fixtureCountryID, "Россия", "643", false, true},
		{"predefined/plain", "../../testdata/custom/country/predefined.json", predefinedCountryID, "Марокко", "504", true, false},
		{"predefined/gzip", "../../testdata/custom/country/predefined.json", predefinedCountryID, "Марокко", "504", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture, err := os.ReadFile(tc.fixture)
			if err != nil {
				t.Fatal(err)
			}
			var calls atomic.Int32
			service := newCountryService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				assertRequest(t, r, http.MethodGet, "/api/remap/1.2/entity/country/"+tc.id, nil)
				writeCountryResponse(t, w, fixture, tc.compressed)
			}))
			got, err := service.Get(t.Context(), tc.id)
			if err != nil || got == nil {
				t.Fatalf("country = %+v, error = %v", got, err)
			}
			if got.ID != tc.id || got.Meta.Type != "country" || got.Name != tc.wantName || got.Code != tc.wantCode || got.Shared {
				t.Fatalf("country = %+v", got)
			}
			if tc.predefined {
				if got.AccountID != "" || got.Owner != nil || got.Group != nil {
					t.Fatalf("absent fields were invented: %+v", got)
				}
			} else if got.AccountID != "abcdef12-1234-1234-1234-123456789012" || got.Owner == nil || got.Group == nil {
				t.Fatalf("account or references = %+v", got)
			}
			if calls.Load() != 1 {
				t.Fatalf("calls = %d, want 1", calls.Load())
			}
		})
	}
}

func TestServiceGetValidation(t *testing.T) {
	if service, err := country.NewService(nil); service != nil || err == nil {
		t.Fatal("nil client accepted")
	}
	var calls atomic.Int32
	service := newCountryService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	for _, id := range invalidIDs() {
		if got, err := service.Get(t.Context(), id); got != nil || err == nil {
			t.Fatalf("invalid id %q accepted", id)
		}
	}
	if got, err := service.Get(nil, fixtureCountryID); got != nil || err == nil {
		t.Fatal("nil context accepted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if got, err := service.Get(ctx, fixtureCountryID); got != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("country = %+v, error = %v", got, err)
	}
	if calls.Load() != 0 {
		t.Fatalf("unexpected calls = %d", calls.Load())
	}
}

func TestServiceCreate(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/upstream/country.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name       string
		input      country.CreateRequest
		want       []byte
		compressed bool
	}{
		{"full/plain", fullCreateRequest(), countryGolden(t, "create"), false},
		{"full/gzip", fullCreateRequest(), countryGolden(t, "create"), true},
		{"minimal", country.CreateRequest{Name: "SDK Country"}, []byte(`{"name":"SDK Country"}`), false},
		{"whitespace name", country.CreateRequest{Name: " \t "}, []byte(`{"name":" \t "}`), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			service := newCountryService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				assertRequest(t, r, http.MethodPost, "/api/remap/1.2/entity/country", tc.want)
				writeCountryResponse(t, w, fixture, tc.compressed)
			}))
			got, err := service.Create(t.Context(), tc.input)
			assertUpstreamCountry(t, got, err)
			if calls.Load() != 1 {
				t.Fatalf("calls = %d, want one POST without a preliminary GET", calls.Load())
			}
		})
	}
}

func TestServiceCreateValidation(t *testing.T) {
	var calls atomic.Int32
	service := newCountryService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	if got, err := service.Create(t.Context(), country.CreateRequest{}); got != nil || err == nil {
		t.Fatal("empty name accepted")
	}
	input := country.CreateRequest{Name: "SDK Country"}
	if got, err := service.Create(nil, input); got != nil || err == nil {
		t.Fatal("nil context accepted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if got, err := service.Create(ctx, input); got != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("country = %+v, error = %v", got, err)
	}
	if calls.Load() != 0 {
		t.Fatalf("unexpected calls = %d", calls.Load())
	}
}

func TestServiceUpdate(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/upstream/country.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name       string
		input      country.UpdateRequest
		want       []byte
		compressed bool
	}{
		{"full/plain", fullUpdateRequest(), countryGolden(t, "update"), false},
		{"full/gzip", fullUpdateRequest(), countryGolden(t, "update"), true},
		{"name", country.UpdateRequest{Name: model.Value("renamed")}, []byte(`{"name":"renamed"}`), false},
		{"empty name", country.UpdateRequest{Name: model.Value("")}, []byte(`{"name":""}`), false},
		{"code", country.UpdateRequest{Code: model.Value("0007")}, []byte(`{"code":"0007"}`), false},
		{"external code", country.UpdateRequest{ExternalCode: model.Value("updated")}, []byte(`{"externalCode":"updated"}`), false},
		{"empty description", country.UpdateRequest{Description: model.Value("")}, []byte(`{"description":""}`), false},
		{"shared false", country.UpdateRequest{Shared: model.Value(false)}, []byte(`{"shared":false}`), false},
		{"owner", country.UpdateRequest{Owner: model.Value(model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/employee/11111111-1111-1111-1111-111111111111"}})}, []byte(`{"owner":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/employee/11111111-1111-1111-1111-111111111111"}}}`), false},
		{"null", country.UpdateRequest{Group: model.Null[model.Reference]()}, []byte(`{"group":null}`), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			service := newCountryService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				assertRequest(t, r, http.MethodPut, "/api/remap/1.2/entity/country/"+fixtureCountryID, tc.want)
				writeCountryResponse(t, w, fixture, tc.compressed)
			}))
			got, err := service.Update(t.Context(), fixtureCountryID, tc.input)
			assertUpstreamCountry(t, got, err)
			if calls.Load() != 1 {
				t.Fatalf("calls = %d, want one PUT without a preliminary GET", calls.Load())
			}
		})
	}
}

func TestServiceUpdateValidation(t *testing.T) {
	var calls atomic.Int32
	service := newCountryService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	input := country.UpdateRequest{Shared: model.Value(false)}
	for _, id := range invalidIDs() {
		if got, err := service.Update(t.Context(), id, input); got != nil || err == nil {
			t.Fatalf("invalid id %q accepted", id)
		}
	}
	if got, err := service.Update(t.Context(), fixtureCountryID, country.UpdateRequest{}); got != nil || err == nil {
		t.Fatal("unset update accepted")
	}
	if got, err := service.Update(nil, fixtureCountryID, input); got != nil || err == nil {
		t.Fatal("nil context accepted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if got, err := service.Update(ctx, fixtureCountryID, input); got != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("country = %+v, error = %v", got, err)
	}
	if calls.Load() != 0 {
		t.Fatalf("unexpected calls = %d", calls.Load())
	}
}

var operations = []struct {
	name    string
	context string
	call    func(*country.Service, context.Context) (*country.Country, error)
}{
	{"get", "get country " + fixtureCountryID, func(s *country.Service, ctx context.Context) (*country.Country, error) {
		return s.Get(ctx, fixtureCountryID)
	}},
	{"create", "create country", func(s *country.Service, ctx context.Context) (*country.Country, error) {
		return s.Create(ctx, country.CreateRequest{Name: "SDK Country"})
	}},
	{"update", "update country " + fixtureCountryID, func(s *country.Service, ctx context.Context) (*country.Country, error) {
		return s.Update(ctx, fixtureCountryID, country.UpdateRequest{Shared: model.Value(false)})
	}},
}

func TestServiceResponseErrors(t *testing.T) {
	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			for _, tc := range []struct {
				name    string
				status  int
				body    string
				wantAPI bool
			}{
				{"bad request", 400, `{"errors":[{"error":"invalid field","code":3000}]}`, true},
				{"unauthorized", 401, `{"errors":[{"error":"invalid credentials","code":1056}]}`, true},
				{"not found", 404, `{"errors":[{"error":"not found","code":1001}]}`, true},
				{"rate limit", 429, `{"errors":[{"error":"slow down","code":1073}]}`, true},
				{"server error", 500, `{"errors":[{"error":"server failed","code":1000}]}`, true},
				{"HTML error", 502, "<html>bad gateway</html>", true},
				{"unexpected status", 201, `{"id":"unexpected"}`, false},
				{"no content", 204, "", false},
				{"missing id", 200, `{"name":"invalid"}`, false},
				{"empty id", 200, `{"id":""}`, false},
				{"null", 200, "null", false},
				{"malformed JSON", 200, `{"id":`, false},
				{"partial decode", 200, `{"id":"response-id","name":42}`, false},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var calls atomic.Int32
					service := newCountryService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						calls.Add(1)
						w.WriteHeader(tc.status)
						if _, err := io.WriteString(w, tc.body); err != nil {
							t.Error(err)
						}
					}))
					got, err := operation.call(service, t.Context())
					if got != nil || err == nil {
						t.Fatalf("country = %+v, error = %v", got, err)
					}
					if strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), operation.context) {
						t.Errorf("error context = %v", err)
					}
					var apiErr *sdkerrors.APIError
					if errors.As(err, &apiErr) != tc.wantAPI {
						t.Fatalf("unexpected error type: %v", err)
					}
					if tc.wantAPI && apiErr.StatusCode != tc.status {
						t.Errorf("status = %d, want %d", apiErr.StatusCode, tc.status)
					}
					if tc.status == 400 && (len(apiErr.Errors) != 1 || apiErr.Errors[0].Message != "invalid field" || apiErr.Errors[0].Code != 3000) {
						t.Errorf("API details = %+v", apiErr.Errors)
					}
					if calls.Load() != 1 {
						t.Fatalf("request retried: calls = %d", calls.Load())
					}
				})
			}
		})
	}
}

// Синтетический отказ проверяет передачу APIError, а не реальный статус или текст живого API.
func TestServiceUpdatePredefinedAPIError(t *testing.T) {
	var calls atomic.Int32
	service := newCountryService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		assertRequest(t, r, http.MethodPut, "/api/remap/1.2/entity/country/"+predefinedCountryID, []byte(`{"shared":false}`))
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, `{"errors":[{"error":"synthetic predefined country rejection","code":9999}]}`)
	}))
	got, err := service.Update(t.Context(), predefinedCountryID, country.UpdateRequest{Shared: model.Value(false)})
	var apiErr *sdkerrors.APIError
	if got != nil || !errors.As(err, &apiErr) {
		t.Fatalf("country = %+v, error = %v", got, err)
	}
	if apiErr.StatusCode != 403 || len(apiErr.Errors) != 1 || apiErr.Errors[0].Message != "synthetic predefined country rejection" || apiErr.Errors[0].Code != 9999 {
		t.Fatalf("API error = %+v", apiErr)
	}
	if calls.Load() != 1 {
		t.Fatalf("calls = %d, want 1", calls.Load())
	}
}

func TestServiceErrorWrapping(t *testing.T) {
	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			for _, authError := range []bool{true, false} {
				name := "transport"
				if authError {
					name = "authorization"
				}
				t.Run(name, func(t *testing.T) {
					wantErr := errors.New("boundary failed")
					var calls atomic.Int32
					cfg := transport.Config{HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
						calls.Add(1)
						return nil, wantErr
					})}}
					if authError {
						cfg.Auth = func(*http.Request) error { return wantErr }
					}
					client, err := transport.NewClient(cfg)
					if err != nil {
						t.Fatal(err)
					}
					service, err := country.NewService(client)
					if err != nil {
						t.Fatal(err)
					}
					got, err := operation.call(service, t.Context())
					if got != nil || !errors.Is(err, wantErr) {
						t.Fatalf("country = %+v, error = %v", got, err)
					}
					if strings.Count(err.Error(), "moysklad:") != 1 || !strings.HasPrefix(err.Error(), operation.context+": ") {
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

func TestServiceUUIDVersion(t *testing.T) {
	const id = "00000000-0000-f000-0000-000000000000"
	service := newCountryService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/remap/1.2/entity/country/"+id {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"id":"`+id+`"}`)
	}))
	for _, tc := range []struct {
		name string
		call func() (*country.Country, error)
	}{
		{"get", func() (*country.Country, error) { return service.Get(t.Context(), id) }},
		{"update", func() (*country.Country, error) {
			return service.Update(t.Context(), id, country.UpdateRequest{Shared: model.Value(false)})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := tc.call(); err != nil || got == nil || got.ID != id {
				t.Fatalf("country = %+v, error = %v", got, err)
			}
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func assertUpstreamCountry(t *testing.T, got *country.Country, err error) {
	t.Helper()
	if err != nil || got == nil {
		t.Fatalf("country = %+v, error = %v", got, err)
	}
	if got.ID != fixtureCountryID || got.Meta.Type != "country" || got.Name != "Россия" || got.Code != "643" ||
		got.Description != "Российская Федерация" || got.ExternalCode != "RU" || got.Shared || got.Updated.IsZero() ||
		got.Owner == nil || got.Group == nil || got.AccountID != "abcdef12-1234-1234-1234-123456789012" {
		t.Fatalf("country = %+v", got)
	}
}

func invalidIDs() []string {
	return []string{"", "metadata", ".", "..", "../metadata", "x/y", fixtureCountryID + "/metadata", "44444444-4444-4444-4444-44444444444g", "444444444444444444444444444444444444"}
}

func assertRequest(t *testing.T, r *http.Request, method, path string, wantBody []byte) {
	t.Helper()
	if r.Method != method || r.URL.Path != path || r.URL.RawQuery != "" {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	}
	if r.Header.Get("Authorization") != "Bearer test-token" || r.Header.Get("Accept-Encoding") != "gzip" ||
		r.Header.Get("Accept") != "application/json;charset=utf-8" || r.Header.Get("User-Agent") != "moysklad-sdk-go" {
		t.Error("required headers missing")
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
		return
	}
	if wantBody == nil {
		if len(body) != 0 || r.Header.Get("Content-Type") != "" {
			t.Errorf("GET body = %q, Content-Type = %q", body, r.Header.Get("Content-Type"))
		}
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		t.Error("Content-Type missing")
	}
	assertJSONEqual(t, body, wantBody)
}

func writeCountryResponse(t *testing.T, w http.ResponseWriter, fixture []byte, compressed bool) {
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

func newCountryService(t *testing.T, handler http.Handler) *country.Service {
	t.Helper()
	server := memnet.NewServer(handler)
	t.Cleanup(server.Close)
	auth, err := transport.BearerToken("test-token")
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", Timeout: 5 * time.Second, Auth: auth})
	if err != nil {
		t.Fatal(err)
	}
	service, err := country.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	return service
}
