package expenseitem_test

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

	"github.com/eteeet/moysklad-sdk-go/entity/expenseitem"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const fixtureExpenseItemID = "fb0a4b75-2e58-11e6-8a84-bae500000058"
const predefinedExpenseItemID = "1be2395a-0479-11e5-baee-448a5b426e7e"

func TestServiceResponseErrors(t *testing.T) {
	for _, operation := range []struct {
		name, method, path, context string
		call                        func(*expenseitem.Service, context.Context) (*expenseitem.ExpenseItem, error)
	}{
		{"get", http.MethodGet, "/api/remap/1.2/entity/expenseitem/" + fixtureExpenseItemID, "get expenseitem " + fixtureExpenseItemID, func(s *expenseitem.Service, ctx context.Context) (*expenseitem.ExpenseItem, error) {
			return s.Get(ctx, fixtureExpenseItemID)
		}},
		{"create", http.MethodPost, "/api/remap/1.2/entity/expenseitem", "create expenseitem", func(s *expenseitem.Service, ctx context.Context) (*expenseitem.ExpenseItem, error) {
			return s.Create(ctx, expenseitem.CreateRequest{Name: "SDK ExpenseItem"})
		}},
		{"update", http.MethodPut, "/api/remap/1.2/entity/expenseitem/" + fixtureExpenseItemID, "update expenseitem " + fixtureExpenseItemID, func(s *expenseitem.Service, ctx context.Context) (*expenseitem.ExpenseItem, error) {
			return s.Update(ctx, fixtureExpenseItemID, expenseitem.UpdateRequest{OperatingExpenses: model.Value(false)})
		}},
	} {
		for _, tc := range []struct {
			name   string
			status int
			body   string
			api    bool
		}{
			{"bad request", 400, `{"errors":[{"error":"invalid field","code":3000}]}`, true},
			{"unauthorized", 401, `{"errors":[{"error":"unauthorized","code":1056}]}`, true},
			{"not found", 404, `{"errors":[{"error":"not found","code":1001}]}`, true},
			{"rate limit", 429, `{"errors":[{"error":"slow down","code":1073}]}`, true},
			{"server error", 500, `{"errors":[{"error":"server failed","code":1000}]}`, true},
			{"HTML error", 502, "<html>bad gateway</html>", true},
			{"unexpected status", 201, `{"id":"unexpected"}`, false},
			{"no content", 204, "", false},
			{"missing id", 200, `{"name":"invalid"}`, false},
			{"empty id", 200, `{"id":""}`, false},
			{"null response", 200, "null", false},
			{"invalid JSON", 200, `{"id":`, false},
		} {
			t.Run(operation.name+"/"+tc.name, func(t *testing.T) {
				var calls atomic.Int32
				service := newExpenseItemService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					if r.Method != operation.method || r.URL.Path != operation.path || r.URL.RawQuery != "" {
						t.Errorf("unexpected request: %s %s", r.Method, r.URL)
					}
					w.WriteHeader(tc.status)
					if _, err := io.WriteString(w, tc.body); err != nil {
						t.Error(err)
					}
				}))
				got, err := operation.call(service, t.Context())
				if err == nil || got != nil {
					t.Fatalf("expenseitem = %+v, error = %v", got, err)
				}
				if strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), operation.context) {
					t.Errorf("error context = %v", err)
				}
				var apiErr *sdkerrors.APIError
				if errors.As(err, &apiErr) != tc.api {
					t.Fatalf("unexpected error type: %v", err)
				}
				if tc.api && apiErr.StatusCode != tc.status {
					t.Errorf("status = %d, want %d", apiErr.StatusCode, tc.status)
				}
				if tc.status == 400 && (len(apiErr.Errors) != 1 || apiErr.Errors[0].Code != 3000 || apiErr.Errors[0].Message != "invalid field") {
					t.Errorf("API details = %+v", apiErr.Errors)
				}
				if calls.Load() != 1 {
					t.Fatalf("request retried: calls = %d", calls.Load())
				}
			})
		}
	}
}

func TestServiceErrorWrapping(t *testing.T) {
	for _, operation := range []struct {
		name    string
		context string
		call    func(*expenseitem.Service, context.Context) (*expenseitem.ExpenseItem, error)
	}{
		{"get", "get expenseitem " + fixtureExpenseItemID, func(s *expenseitem.Service, ctx context.Context) (*expenseitem.ExpenseItem, error) {
			return s.Get(ctx, fixtureExpenseItemID)
		}},
		{"create", "create expenseitem", func(s *expenseitem.Service, ctx context.Context) (*expenseitem.ExpenseItem, error) {
			return s.Create(ctx, expenseitem.CreateRequest{Name: "SDK ExpenseItem"})
		}},
		{"update", "update expenseitem " + fixtureExpenseItemID, func(s *expenseitem.Service, ctx context.Context) (*expenseitem.ExpenseItem, error) {
			return s.Update(ctx, fixtureExpenseItemID, expenseitem.UpdateRequest{Shared: model.Value(false)})
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
					service, err := expenseitem.NewService(client)
					if err != nil {
						t.Fatal(err)
					}
					got, err := operation.call(service, t.Context())
					if got != nil || !errors.Is(err, wantErr) {
						t.Fatalf("expenseitem = %+v, error = %v, want wrapped error", got, err)
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
	if service, err := expenseitem.NewService(nil); err == nil || service != nil {
		t.Fatal("nil client accepted")
	}
	var calls atomic.Int32
	service := newExpenseItemService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	for _, id := range []string{
		"", "metadata", ".", "..", "../metadata", "x/y", fixtureExpenseItemID + "/metadata",
		"77777777-7777-7777-7777-77777777777g", "777777777777777777777777777777777777",
	} {
		if got, err := service.Get(t.Context(), id); err == nil || got != nil {
			t.Fatalf("invalid id %q accepted", id)
		}
	}
	if got, err := service.Get(nil, fixtureExpenseItemID); err == nil || got != nil {
		t.Fatal("nil context accepted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if got, err := service.Get(ctx, fixtureExpenseItemID); got != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("expenseitem = %+v, error = %v, want cancellation", got, err)
	}
	if calls.Load() != 0 {
		t.Fatalf("unexpected HTTP calls: %d", calls.Load())
	}
}

func TestServiceUUIDVersion(t *testing.T) {
	const id = "00000000-0000-f000-0000-000000000000"
	service := newExpenseItemService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/remap/1.2/entity/expenseitem/"+id {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"id":"`+id+`"}`)
	}))
	for _, tc := range []struct {
		name string
		call func(context.Context) (*expenseitem.ExpenseItem, error)
	}{
		{"get", func(ctx context.Context) (*expenseitem.ExpenseItem, error) { return service.Get(ctx, id) }},
		{"update", func(ctx context.Context) (*expenseitem.ExpenseItem, error) {
			return service.Update(ctx, id, expenseitem.UpdateRequest{Shared: model.Value(false)})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.call(t.Context())
			if err != nil || got == nil || got.ID != id {
				t.Fatalf("expenseitem = %+v, error = %v", got, err)
			}
		})
	}
}

func TestServiceCreateValidation(t *testing.T) {
	var calls atomic.Int32
	service := newExpenseItemService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	if got, err := service.Create(t.Context(), expenseitem.CreateRequest{}); err == nil || got != nil {
		t.Fatal("empty name accepted")
	}
	if got, err := service.Create(nil, expenseitem.CreateRequest{Name: "SDK ExpenseItem"}); err == nil || got != nil {
		t.Fatal("nil context accepted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if got, err := service.Create(ctx, expenseitem.CreateRequest{Name: "SDK ExpenseItem"}); got != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("expenseitem = %+v, error = %v, want cancellation", got, err)
	}
	if calls.Load() != 0 {
		t.Fatalf("unexpected HTTP calls: %d", calls.Load())
	}
}

func TestServiceUpdateValidation(t *testing.T) {
	var calls atomic.Int32
	service := newExpenseItemService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	input := expenseitem.UpdateRequest{Shared: model.Value(false)}
	for _, id := range []string{
		"", "metadata", ".", "..", "../metadata", "x/y", fixtureExpenseItemID + "/metadata",
		"77777777-7777-7777-7777-77777777777g", "777777777777777777777777777777777777",
	} {
		if got, err := service.Update(t.Context(), id, input); err == nil || got != nil {
			t.Fatalf("invalid id %q accepted", id)
		}
	}
	if got, err := service.Update(t.Context(), fixtureExpenseItemID, expenseitem.UpdateRequest{}); err == nil || got != nil {
		t.Fatal("empty update accepted")
	}
	if got, err := service.Update(nil, fixtureExpenseItemID, input); err == nil || got != nil {
		t.Fatal("nil context accepted")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if got, err := service.Update(ctx, fixtureExpenseItemID, input); got != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("expenseitem = %+v, error = %v, want cancellation", got, err)
	}
	if calls.Load() != 0 {
		t.Fatalf("unexpected HTTP calls: %d", calls.Load())
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestServiceUpdate(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/custom/expenseitem/predefined.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		input expenseitem.UpdateRequest
		want  []byte
		gzip  bool
	}{
		{"full/plain", fullUpdateRequest(), expenseitemGolden(t, "update"), false},
		{"full/gzip", fullUpdateRequest(), expenseitemGolden(t, "update"), true},
		{"name", expenseitem.UpdateRequest{Name: model.Value("renamed")}, []byte(`{"name":"renamed"}`), false},
		{"empty name", expenseitem.UpdateRequest{Name: model.Value("")}, []byte(`{"name":""}`), false},
		{"empty description", expenseitem.UpdateRequest{Description: model.Value("")}, []byte(`{"description":""}`), false},
		{"code", expenseitem.UpdateRequest{Code: model.Value("003")}, []byte(`{"code":"003"}`), false},
		{"external code", expenseitem.UpdateRequest{ExternalCode: model.Value("updated")}, []byte(`{"externalCode":"updated"}`), false},
		{"shared false", expenseitem.UpdateRequest{Shared: model.Value(false)}, []byte(`{"shared":false}`), false},
		{"archived false", expenseitem.UpdateRequest{Archived: model.Value(false)}, []byte(`{"archived":false}`), false},
		{"archived true", expenseitem.UpdateRequest{Archived: model.Value(true)}, []byte(`{"archived":true}`), false},
		{"operating expenses false", expenseitem.UpdateRequest{OperatingExpenses: model.Value(false)}, []byte(`{"operatingExpenses":false}`), false},
		{"operating expenses true", expenseitem.UpdateRequest{OperatingExpenses: model.Value(true)}, []byte(`{"operatingExpenses":true}`), false},
		{"owner", expenseitem.UpdateRequest{Owner: model.Value(model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/employee/11111111-1111-1111-1111-111111111111"}})}, []byte(`{"owner":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/employee/11111111-1111-1111-1111-111111111111"}}}`), false},
		{"explicit null", expenseitem.UpdateRequest{Group: model.Null[model.Reference]()}, []byte(`{"group":null}`), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			service := newExpenseItemService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodPut || r.URL.Path != "/api/remap/1.2/entity/expenseitem/"+fixtureExpenseItemID || r.URL.RawQuery != "" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				assertHeaders(t, r, true)
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
					return
				}
				assertJSONEqual(t, body, tc.want)
				writeExpenseItemResponse(t, w, fixture, tc.gzip)
			}))
			got, err := service.Update(t.Context(), fixtureExpenseItemID, tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got.ID != predefinedExpenseItemID || got.Name != "Возврат" || got.Code != "3" || got.ExternalCode != "3" ||
				got.Meta.Type != "expenseitem" || got.OperatingExpenses || got.Archived || got.Updated.IsZero() {
				t.Fatalf("expenseitem = %+v", got)
			}
			if calls.Load() != 1 {
				t.Fatalf("calls = %d, want one PUT without preliminary GET", calls.Load())
			}
		})
	}
}

func TestServiceCreate(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/upstream/expense_item.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		input expenseitem.CreateRequest
		want  []byte
		gzip  bool
	}{
		{"minimal", expenseitem.CreateRequest{Name: "SDK ExpenseItem"}, []byte(`{"name":"SDK ExpenseItem"}`), false},
		{"full/plain", fullCreateRequest(), expenseitemGolden(t, "create"), false},
		{"full/gzip", fullCreateRequest(), expenseitemGolden(t, "create"), true},
		{"whitespace name", expenseitem.CreateRequest{Name: " \t "}, []byte(`{"name":" \t "}`), false},
		{"false flags", expenseitem.CreateRequest{Name: "SDK ExpenseItem", Archived: model.Value(false), OperatingExpenses: model.Value(false)}, []byte(`{"name":"SDK ExpenseItem","archived":false,"operatingExpenses":false}`), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			service := newExpenseItemService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodPost || r.URL.Path != "/api/remap/1.2/entity/expenseitem" || r.URL.RawQuery != "" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				assertHeaders(t, r, true)
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
					return
				}
				assertJSONEqual(t, body, tc.want)
				writeExpenseItemResponse(t, w, fixture, tc.gzip)
			}))
			got, err := service.Create(t.Context(), tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureExpenseItemID || got.Name != "Аренда" || got.Description != "Аренда" || got.Code != "Аренда" ||
				got.ExternalCode != "IVslr34uhCUuglxPD7Idm0" || got.AccountID != "f976ed28-2e58-11e6-8a84-bae500000001" ||
				got.Meta.Type != "expenseitem" || !got.Archived || got.OperatingExpenses || got.Updated.IsZero() {
				t.Fatalf("expenseitem = %+v", got)
			}
			if calls.Load() != 1 {
				t.Fatalf("calls = %d, want one POST", calls.Load())
			}
		})
	}
}

func TestServiceGet(t *testing.T) {
	for _, fixture := range []struct {
		name, path, id, wantName, wantCode string
		archived                           bool
	}{
		{"upstream", "../../testdata/upstream/expense_item.json", fixtureExpenseItemID, "Аренда", "Аренда", true},
		{"predefined", "../../testdata/custom/expenseitem/predefined.json", predefinedExpenseItemID, "Возврат", "3", false},
	} {
		for _, compressed := range []bool{false, true} {
			name := fixture.name + "/plain"
			if compressed {
				name = fixture.name + "/gzip"
			}
			t.Run(name, func(t *testing.T) {
				data, err := os.ReadFile(fixture.path)
				if err != nil {
					t.Fatal(err)
				}
				var calls atomic.Int32
				service := newExpenseItemService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					if r.Method != http.MethodGet || r.URL.Path != "/api/remap/1.2/entity/expenseitem/"+fixture.id || r.URL.RawQuery != "" {
						t.Errorf("unexpected request: %s %s", r.Method, r.URL)
					}
					assertHeaders(t, r, false)
					body, err := io.ReadAll(r.Body)
					if err != nil || len(body) != 0 {
						t.Errorf("GET body = %q, error = %v", body, err)
					}
					writeExpenseItemResponse(t, w, data, compressed)
				}))
				got, err := service.Get(t.Context(), fixture.id)
				if err != nil {
					t.Fatal(err)
				}
				if got.ID != fixture.id || got.Meta.Type != "expenseitem" || got.Name != fixture.wantName || got.Code != fixture.wantCode ||
					got.Archived != fixture.archived || got.OperatingExpenses || got.Shared || got.Owner != nil || got.Group != nil {
					t.Fatalf("expenseitem = %+v", got)
				}
				if fixture.id == predefinedExpenseItemID && got.AccountID != "" {
					t.Errorf("accountId invented: %q", got.AccountID)
				}
				if calls.Load() != 1 {
					t.Fatalf("calls = %d, want 1", calls.Load())
				}
			})
		}
	}
}

func assertHeaders(t *testing.T, r *http.Request, body bool) {
	t.Helper()
	if r.Header.Get("Authorization") != "Bearer test-token" || r.Header.Get("Accept-Encoding") != "gzip" ||
		r.Header.Get("Accept") != "application/json;charset=utf-8" || r.Header.Get("User-Agent") != "moysklad-sdk-go" {
		t.Error("required request headers missing")
	}
	wantContentType := ""
	if body {
		wantContentType = "application/json"
	}
	if r.Header.Get("Content-Type") != wantContentType {
		t.Errorf("Content-Type = %q", r.Header.Get("Content-Type"))
	}
}

func writeExpenseItemResponse(t *testing.T, w http.ResponseWriter, fixture []byte, compressed bool) {
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

func newExpenseItemService(t *testing.T, handler http.Handler) *expenseitem.Service {
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
	service, err := expenseitem.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	return service
}
