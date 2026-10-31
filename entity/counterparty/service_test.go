package counterparty_test

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

	"github.com/eteeet/moysklad-sdk-go/entity/counterparty"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const fixtureCounterpartyID = "12a8b923-692c-11e6-8a84-bae500000056"

type serviceOperation struct {
	name, method, path, context string
	body                        []byte
	call                        func(*counterparty.Service, context.Context) (*counterparty.Counterparty, error)
}

func serviceOperations(t *testing.T) []serviceOperation {
	t.Helper()
	return []serviceOperation{
		{"get", http.MethodGet, "/api/remap/1.2/entity/counterparty/" + fixtureCounterpartyID, "get counterparty " + fixtureCounterpartyID, nil,
			func(s *counterparty.Service, ctx context.Context) (*counterparty.Counterparty, error) {
				return s.Get(ctx, fixtureCounterpartyID)
			}},
		{"create", http.MethodPost, "/api/remap/1.2/entity/counterparty", "create counterparty", golden(t, "create"),
			func(s *counterparty.Service, ctx context.Context) (*counterparty.Counterparty, error) {
				return s.Create(ctx, fullCreateRequest(t))
			}},
		{"update", http.MethodPut, "/api/remap/1.2/entity/counterparty/" + fixtureCounterpartyID, "update counterparty " + fixtureCounterpartyID, golden(t, "update"),
			func(s *counterparty.Service, ctx context.Context) (*counterparty.Counterparty, error) {
				return s.Update(ctx, fixtureCounterpartyID, fullUpdateRequest(t))
			}},
		{"update sync ID", http.MethodPut, "/api/remap/1.2/entity/counterparty/" + fixtureCounterpartyID, "update counterparty " + fixtureCounterpartyID, []byte(`{"syncId":"11111111-1111-1111-1111-111111111111"}`),
			func(s *counterparty.Service, ctx context.Context) (*counterparty.Counterparty, error) {
				return s.Update(ctx, fixtureCounterpartyID, counterparty.UpdateRequest{SyncID: model.Value("11111111-1111-1111-1111-111111111111")})
			}},
	}
}

func TestServiceOperations(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/upstream/counterparty.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, op := range serviceOperations(t) {
		for _, compressed := range []bool{false, true} {
			name := op.name + "/plain"
			if compressed {
				name = op.name + "/gzip"
			}
			t.Run(name, func(t *testing.T) {
				var calls atomic.Int32
				service := newService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					assertRequest(t, r, op.method, op.path, op.body)
					writeResponse(t, w, fixture, compressed)
				}))
				got, err := op.call(service, t.Context())
				if err != nil {
					t.Fatal(err)
				}
				if got == nil || got.ID != fixtureCounterpartyID || got.Meta.Type != "counterparty" || got.Name != "ООО Тестовый контрагент" ||
					got.Description != "testDescirption" || got.Code != "CP-001" || !got.Shared || got.Archived || got.BonusPoints != 5 || got.SalesAmount != 0 ||
					got.Owner == nil || got.Group == nil || got.State == nil || got.BonusProgram == nil || got.PriceType == nil ||
					got.ActualAddressFull == nil || got.LegalAddressFull == nil || len(got.Attributes) != 3 || len(got.Tags) != 1 ||
					got.Created.IsZero() || got.Updated.IsZero() {
					t.Fatalf("counterparty = %+v", got)
				}
				if calls.Load() != 1 {
					t.Errorf("calls = %d, want one request without follow-up or retry", calls.Load())
				}
			})
		}
	}
}

func assertRequest(t *testing.T, r *http.Request, method, path string, wantBody []byte) {
	t.Helper()
	if r.Method != method || r.URL.Path != path || r.URL.RawQuery != "" {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	}
	if r.Header.Get("Authorization") != "Bearer test-token" || r.Header.Get("Accept-Encoding") != "gzip" ||
		r.Header.Get("Accept") != "application/json;charset=utf-8" || r.Header.Get("User-Agent") != "moysklad-sdk-go" {
		t.Error("required request headers missing")
	}
	wantContentType := "application/json"
	if wantBody == nil {
		wantContentType = ""
	}
	if got := r.Header.Get("Content-Type"); got != wantContentType {
		t.Errorf("Content-Type = %q, want %q", got, wantContentType)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
		return
	}
	if wantBody == nil {
		if len(body) != 0 {
			t.Errorf("GET body = %q", body)
		}
	} else {
		assertJSONEqual(t, body, wantBody)
	}
}

func writeResponse(t *testing.T, w http.ResponseWriter, body []byte, compressed bool) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if compressed {
		w.Header().Set("Content-Encoding", "gzip")
		writer := gzip.NewWriter(w)
		if _, err := writer.Write(body); err != nil {
			t.Error(err)
		}
		if err := writer.Close(); err != nil {
			t.Error(err)
		}
		return
	}
	if _, err := w.Write(body); err != nil {
		t.Error(err)
	}
}

func newService(t *testing.T, handler http.Handler) *counterparty.Service {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	auth, err := transport.BearerToken("test-token")
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", Timeout: 5 * time.Second, Auth: auth})
	if err != nil {
		t.Fatal(err)
	}
	service, err := counterparty.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	return service
}

func TestServiceErrors(t *testing.T) {
	for _, op := range serviceOperations(t) {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range []struct {
				name    string
				status  int
				body    string
				wantAPI bool
			}{
				{"bad request", 400, `{"errors":[{"error":"invalid field","code":3000}]}`, true},
				{"forbidden", 403, `{"errors":[{"error":"access denied","code":1006}]}`, true},
				{"not found", 404, `{"errors":[{"error":"not found","code":1001}]}`, true},
				{"rate limit", 429, `{"errors":[{"error":"slow down","code":1073}]}`, true},
				{"server error", 500, `{"errors":[{"error":"server failed","code":1000}]}`, true},
				{"HTML gateway error", 502, `<html>bad gateway</html>`, true},
				{"unexpected status", 201, `{"id":"unexpected"}`, false},
				{"no content", 204, ``, false},
				{"missing id", 200, `{"name":"missing"}`, false},
				{"empty id", 200, `{"id":""}`, false},
				{"null response", 200, `null`, false},
				{"malformed JSON", 200, `{"id":"partial",`, false},
				{"invalid date", 200, `{"id":"partial","birthDate":"bad-date"}`, false},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var calls atomic.Int32
					service := newService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						calls.Add(1)
						assertRequest(t, r, op.method, op.path, op.body)
						w.WriteHeader(tc.status)
						if _, err := io.WriteString(w, tc.body); err != nil {
							t.Error(err)
						}
					}))
					got, err := op.call(service, t.Context())
					if err == nil || got != nil {
						t.Fatalf("counterparty = %+v, error = %v", got, err)
					}
					if strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), op.context) {
						t.Errorf("error context = %v", err)
					}
					var apiErr *sdkerrors.APIError
					if errors.As(err, &apiErr) != tc.wantAPI {
						t.Fatalf("unexpected error type: %v", err)
					}
					if tc.wantAPI && apiErr.StatusCode != tc.status {
						t.Errorf("status = %d, want %d", apiErr.StatusCode, tc.status)
					}
					if tc.status == 400 && (len(apiErr.Errors) != 1 || apiErr.Errors[0].Code != 3000 || apiErr.Errors[0].Message != "invalid field") {
						t.Errorf("API details = %+v", apiErr.Errors)
					}
					if tc.status == 429 && (len(apiErr.Errors) != 1 || apiErr.Errors[0].Code != 1073) {
						t.Errorf("rate limit details = %+v", apiErr.Errors)
					}
					if calls.Load() != 1 {
						t.Errorf("request retried: calls = %d", calls.Load())
					}
				})
			}
		})
	}
}

func TestServiceBoundaryErrors(t *testing.T) {
	for _, op := range serviceOperations(t) {
		for _, authError := range []bool{true, false} {
			name := "transport"
			if authError {
				name = "authorization"
			}
			t.Run(op.name+"/"+name, func(t *testing.T) {
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
				service, err := counterparty.NewService(client)
				if err != nil {
					t.Fatal(err)
				}
				got, err := op.call(service, t.Context())
				if got != nil || !errors.Is(err, wantErr) {
					t.Fatalf("counterparty = %+v, error = %v, want wrapped boundary error", got, err)
				}
				if strings.Count(err.Error(), "moysklad:") != 1 || !strings.HasPrefix(err.Error(), op.context+": ") {
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
}

func TestServiceValidationAndContext(t *testing.T) {
	if service, err := counterparty.NewService(nil); service != nil || err == nil {
		t.Fatal("nil client accepted")
	}
	var calls atomic.Int32
	service := newService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	for _, id := range []string{"", "metadata", ".", "..", "../metadata", "x/y", fixtureCounterpartyID + "/metadata", "11111111-1111-1111-1111-11111111111g", "111111111111111111111111111111111111"} {
		if got, err := service.Get(t.Context(), id); got != nil || err == nil {
			t.Errorf("invalid Get id %q accepted", id)
		}
		if got, err := service.Update(t.Context(), id, counterparty.UpdateRequest{Shared: model.Value(false)}); got != nil || err == nil {
			t.Errorf("invalid Update id %q accepted", id)
		}
	}
	if got, err := service.Create(t.Context(), counterparty.CreateRequest{}); got != nil || err == nil {
		t.Error("empty Create.Name accepted")
	}
	if got, err := service.Update(t.Context(), fixtureCounterpartyID, counterparty.UpdateRequest{}); got != nil || err == nil {
		t.Error("unset Update accepted")
	}
	for _, op := range serviceOperations(t) {
		t.Run(op.name, func(t *testing.T) {
			if got, err := op.call(service, nil); got != nil || err == nil {
				t.Error("nil context accepted")
			}
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			got, err := op.call(service, ctx)
			if got != nil || !errors.Is(err, context.Canceled) {
				t.Fatalf("counterparty = %+v, error = %v, want cancellation", got, err)
			}
			if strings.Count(err.Error(), "moysklad:") != 1 || !strings.HasPrefix(err.Error(), op.context+": ") {
				t.Errorf("error context = %v", err)
			}
		})
	}
	if calls.Load() != 0 {
		t.Errorf("unexpected HTTP calls: %d", calls.Load())
	}
}

func TestServiceSparseWrites(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input any
		want  string
	}{
		{"minimal create", counterparty.CreateRequest{Name: "cp"}, `{"name":"cp"}`},
		{"whitespace create name", counterparty.CreateRequest{Name: " \t "}, `{"name":" \t "}`},
		{"create unknown type and syncId", counterparty.CreateRequest{Name: "cp", CompanyType: model.Value(model.CompanyType("future")), SyncID: model.Value("server-validates-sync-id")}, `{"name":"cp","companyType":"future","syncId":"server-validates-sync-id"}`},
		{"create empty address", counterparty.CreateRequest{Name: "cp", ActualAddressFull: model.Value(model.AddressRequest{})}, `{"name":"cp","actualAddressFull":{}}`},
		{"empty name", counterparty.UpdateRequest{Name: model.Value("")}, `{"name":""}`},
		{"sync ID", counterparty.UpdateRequest{SyncID: model.Value("11111111-1111-1111-1111-111111111111")}, `{"syncId":"11111111-1111-1111-1111-111111111111"}`},
		{"empty sync ID", counterparty.UpdateRequest{SyncID: model.Value("")}, `{"syncId":""}`},
		{"null sync ID", counterparty.UpdateRequest{SyncID: model.Null[string]()}, `{"syncId":null}`},
		{"false", counterparty.UpdateRequest{Shared: model.Value(false)}, `{"shared":false}`},
		{"empty tags", counterparty.UpdateRequest{Tags: model.Value([]string{})}, `{"tags":[]}`},
		{"empty attributes", counterparty.UpdateRequest{Attributes: model.Value([]model.Attribute{})}, `{"attributes":[]}`},
		{"null reference", counterparty.UpdateRequest{BonusProgram: model.Null[model.Reference]()}, `{"bonusProgram":null}`},
		{"null address", counterparty.UpdateRequest{ActualAddressFull: model.Null[model.AddressRequest]()}, `{"actualAddressFull":null}`},
		{"empty address", counterparty.UpdateRequest{ActualAddressFull: model.Value(model.AddressRequest{})}, `{"actualAddressFull":{}}`},
		{"empty nested address", counterparty.UpdateRequest{LegalAddressFull: model.Value(model.AddressRequest{Street: model.Value("")})}, `{"legalAddressFull":{"street":""}}`},
		{"both address forms", counterparty.UpdateRequest{ActualAddress: model.Value("строка"), ActualAddressFull: model.Value(model.AddressRequest{City: model.Value("город")})}, `{"actualAddress":"строка","actualAddressFull":{"city":"город"}}`},
		{"empty address string", counterparty.UpdateRequest{LegalAddress: model.Value("")}, `{"legalAddress":""}`},
		{"unknown companyType", counterparty.UpdateRequest{CompanyType: model.Value(model.CompanyType("future"))}, `{"companyType":"future"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			method, path := http.MethodPut, "/api/remap/1.2/entity/counterparty/"+fixtureCounterpartyID
			if _, ok := tc.input.(counterparty.CreateRequest); ok {
				method, path = http.MethodPost, "/api/remap/1.2/entity/counterparty"
			}
			var calls atomic.Int32
			service := newService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				assertRequest(t, r, method, path, []byte(tc.want))
				writeResponse(t, w, []byte(`{"id":"server-id","name":"from server","actualAddress":"from server","tags":["server-tag"]}`), false)
			}))
			var got *counterparty.Counterparty
			var err error
			switch input := tc.input.(type) {
			case counterparty.CreateRequest:
				got, err = service.Create(t.Context(), input)
			case counterparty.UpdateRequest:
				got, err = service.Update(t.Context(), fixtureCounterpartyID, input)
			}
			if err != nil || got == nil || got.ID != "server-id" || got.Name != "from server" || got.ActualAddress != "from server" || len(got.Tags) != 1 || got.Tags[0] != "server-tag" {
				t.Fatalf("counterparty = %+v, error = %v", got, err)
			}
			if calls.Load() != 1 {
				t.Errorf("calls = %d, want one write without preliminary Get", calls.Load())
			}
		})
	}
}

func TestServiceUUIDVersion(t *testing.T) {
	for _, id := range []string{"00000000-0000-f000-0000-000000000000", "ABCDEF12-ABCD-0000-0000-ABCDEF123456"} {
		service := newService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/remap/1.2/entity/counterparty/"+id {
				t.Errorf("path = %s", r.URL.Path)
			}
			writeResponse(t, w, []byte(`{"id":"`+id+`"}`), false)
		}))
		for _, call := range []func() (*counterparty.Counterparty, error){
			func() (*counterparty.Counterparty, error) { return service.Get(t.Context(), id) },
			func() (*counterparty.Counterparty, error) {
				return service.Update(t.Context(), id, counterparty.UpdateRequest{Shared: model.Value(false)})
			},
		} {
			got, err := call()
			if err != nil || got == nil || got.ID != id {
				t.Fatalf("UUID %q: counterparty = %+v, error = %v", id, got, err)
			}
		}
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }
