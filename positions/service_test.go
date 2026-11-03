package positions_test

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const readerDocumentID = "77777777-7777-7777-7777-777777777777"

type readerRow struct {
	Label  string  `json:"label"`
	Amount float64 `json:"amount"`
}

func TestNewServiceRequiresClient(t *testing.T) {
	if service, err := positions.NewService[readerRow](nil, "payroll"); err == nil || service != nil {
		t.Fatalf("service = %+v, error = %v, want nil client rejection", service, err)
	}
}

func TestNewServiceDocumentTypeValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	for _, documentType := range []string{
		"", "Supply", " supply", "supply ", "supply\n", "1supply", "supply-order", "supply_order",
		".", "..", "../supply", "supply/positions", "supply?limit=1", "supply#meta",
		"supply%2fmetadata", "supply%2Fmetadata", "supply%252fmetadata", "supply:meta", "счёт", "supply\x00",
	} {
		t.Run(documentType, func(t *testing.T) {
			service, err := positions.NewService[readerRow](client, documentType)
			if err == nil || service != nil {
				t.Fatalf("invalid document type %q: service = %+v, error = %v", documentType, service, err)
			}
		})
	}
	for _, documentType := range []string{"a", "supply", "future2", "a0b1"} {
		if service, err := positions.NewService[readerRow](client, documentType); err != nil || service == nil {
			t.Errorf("safe document type %q: service = %+v, error = %v", documentType, service, err)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("constructor made HTTP calls: %d", calls.Load())
	}
}

func TestServiceGetPage(t *testing.T) {
	for _, compressed := range []bool{false, true} {
		name := "plain"
		if compressed {
			name = "gzip"
		}
		t.Run(name, func(t *testing.T) {
			var calls atomic.Int32
			service := readerNewService(t, "payroll", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodGet ||
					r.URL.Path != "/api/remap/1.2/entity/payroll/"+readerDocumentID+"/positions" ||
					r.URL.RawQuery != "limit=1000&offset=0" {
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
				readerWriteResponse(t, w, `{
			"context":{"employee":{"meta":{"href":"employee-link","type":"employee"}}},
			"meta":{
				"href":"page-link","metadataHref":"metadata-link","type":"payrollposition",
				"mediaType":"application/json","uuidHref":"uuid-link","downloadHref":"download-link",
				"permitted":false,"size":37,"limit":12,"offset":98,
				"nextHref":"https://outside.invalid/next","previousHref":"previous-link"
			},
			"rows":[{"label":"salary","amount":12.5,"unknown":true},{}],
			"unknown":{"href":"https://outside.invalid/unknown"}
		}`, compressed)
			}))

			got, err := service.GetPage(t.Context(), readerDocumentID, positions.PageOptions{})
			if err != nil || got == nil {
				t.Fatalf("page = %+v, error = %v", got, err)
			}
			if got.Context == nil || got.Context.Employee == nil ||
				got.Context.Employee.Meta.Href != "employee-link" || got.Context.Employee.Meta.Type != "employee" {
				t.Errorf("context = %+v", got.Context)
			}
			if got.Meta.Href != "page-link" || got.Meta.MetadataHref != "metadata-link" ||
				got.Meta.Type != "payrollposition" || got.Meta.MediaType != "application/json" ||
				got.Meta.UUIDHref != "uuid-link" || got.Meta.DownloadHref != "download-link" ||
				got.Meta.Permitted == nil || *got.Meta.Permitted || got.Meta.Size != 37 ||
				got.Meta.Limit != 12 || got.Meta.Offset != 98 ||
				got.Meta.NextHref != "https://outside.invalid/next" || got.Meta.PreviousHref != "previous-link" {
				t.Errorf("metadata = %+v", got.Meta)
			}
			if len(got.Rows) != 2 || got.Rows[0] != (readerRow{Label: "salary", Amount: 12.5}) ||
				got.Rows[1] != (readerRow{}) {
				t.Errorf("rows = %+v", got.Rows)
			}
			if calls.Load() != 1 {
				t.Fatalf("calls = %d, want 1 without following metadata links", calls.Load())
			}
		})
	}
}

func TestServiceGetPagePagination(t *testing.T) {
	for _, tc := range []struct {
		name  string
		id    string
		page  positions.PageOptions
		query string
	}{
		{"default", readerDocumentID, positions.PageOptions{}, "limit=1000&offset=0"},
		{"minimum limit", readerDocumentID, positions.PageOptions{Limit: 1}, "limit=1&offset=0"},
		{"maximum limit and offset", readerDocumentID, positions.PageOptions{Limit: 1000, Offset: 7}, "limit=1000&offset=7"},
		{"default limit with offset", readerDocumentID, positions.PageOptions{Offset: 20}, "limit=1000&offset=20"},
		{"unrestricted UUID version", "00000000-0000-f000-0000-000000000000", positions.PageOptions{Limit: 1, Offset: 2}, "limit=1&offset=2"},
		{"uppercase UUID hex", "ABCDEF12-ABCD-FBCD-ABCD-ABCDEF123456", positions.PageOptions{Limit: 1}, "limit=1&offset=0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			service := readerNewService(t, "a0b1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodGet ||
					r.URL.Path != "/api/remap/1.2/entity/a0b1/"+tc.id+"/positions" || r.URL.RawQuery != tc.query {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				readerWriteResponse(t, w, `{"meta":{"href":"page-link"},"rows":[]}`, false)
			}))
			got, err := service.GetPage(t.Context(), tc.id, tc.page)
			if err != nil || got == nil {
				t.Fatalf("page = %+v, error = %v", got, err)
			}
			if got.Rows == nil || len(got.Rows) != 0 || got.Context != nil ||
				got.Meta.Href != "page-link" || got.Meta.Size != 0 || got.Meta.Limit != 0 ||
				got.Meta.Offset != 0 || got.Meta.NextHref != "" || got.Meta.PreviousHref != "" {
				t.Errorf("empty page or absent metadata changed: %+v", got)
			}
			if calls.Load() != 1 {
				t.Fatalf("calls = %d, want 1", calls.Load())
			}
		})
	}
}

func TestServiceGetPageExplicitPages(t *testing.T) {
	var calls atomic.Int32
	service := readerNewService(t, "payroll", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != http.MethodGet || r.URL.Path != "/api/remap/1.2/entity/payroll/"+readerDocumentID+"/positions" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL)
		}
		switch r.URL.RawQuery {
		case "limit=1&offset=0":
			readerWriteResponse(t, w, `{"meta":{"href":"first-page","nextHref":"https://outside.invalid/next"},"rows":[{"label":"first"}]}`, false)
		case "limit=1&offset=1":
			readerWriteResponse(t, w, `{"meta":{"href":"second-page","previousHref":"https://outside.invalid/previous"},"rows":[{"label":"second"}]}`, false)
		default:
			t.Errorf("query = %s", r.URL.RawQuery)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	first, err := service.GetPage(t.Context(), readerDocumentID, positions.PageOptions{Limit: 1})
	if err != nil || first == nil {
		t.Fatalf("first page = %+v, error = %v", first, err)
	}
	if calls.Load() != 1 {
		t.Fatalf("first GetPage made %d calls", calls.Load())
	}
	second, err := service.GetPage(t.Context(), readerDocumentID, positions.PageOptions{Limit: 1, Offset: 1})
	if err != nil || second == nil {
		t.Fatalf("second page = %+v, error = %v", second, err)
	}
	if first == second || first.Meta.Href != "first-page" || second.Meta.Href != "second-page" ||
		len(first.Rows) != 1 || first.Rows[0].Label != "first" ||
		len(second.Rows) != 1 || second.Rows[0].Label != "second" ||
		first.Meta.NextHref != "https://outside.invalid/next" || second.Meta.PreviousHref != "https://outside.invalid/previous" {
		t.Errorf("first = %+v, second = %+v", first, second)
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d, want one per explicit GetPage", calls.Load())
	}
}

func TestServiceGetPageValidation(t *testing.T) {
	var calls atomic.Int32
	service := readerNewService(t, "payroll", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	for _, id := range []string{
		"", "metadata", ".", "..", "../metadata", "x/y", readerDocumentID + "/metadata",
		readerDocumentID + "?limit=1", readerDocumentID + "#meta", "77777777-7777-7777-7777-77777777777g",
		"777777777777777777777777777777777777", "777777777777-7777-7777-777777777777",
	} {
		t.Run("id/"+id, func(t *testing.T) {
			got, err := service.GetPage(t.Context(), id, positions.PageOptions{})
			if err == nil || got != nil {
				t.Fatalf("invalid id %q: page = %+v, error = %v", id, got, err)
			}
			readerAssertErrorContext(t, err, "payroll", id)
		})
	}
	for _, page := range []positions.PageOptions{
		{Limit: -1}, {Limit: 1001}, {Offset: -1}, {Limit: 1, Offset: -1},
	} {
		got, err := service.GetPage(t.Context(), readerDocumentID, page)
		if err == nil || got != nil {
			t.Errorf("invalid page %+v: page = %+v, error = %v", page, got, err)
			continue
		}
		readerAssertErrorContext(t, err, "payroll", readerDocumentID)
	}
	if got, err := service.GetPage(nil, readerDocumentID, positions.PageOptions{}); err == nil || got != nil {
		t.Fatalf("nil context: page = %+v, error = %v", got, err)
	} else {
		readerAssertErrorContext(t, err, "payroll", readerDocumentID)
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if got, err := service.GetPage(ctx, readerDocumentID, positions.PageOptions{}); got != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("page = %+v, error = %v, want cancellation", got, err)
	} else {
		readerAssertErrorContext(t, err, "payroll", readerDocumentID)
	}
	expired, cancelExpired := context.WithDeadline(t.Context(), time.Now().Add(-time.Second))
	defer cancelExpired()
	if got, err := service.GetPage(expired, readerDocumentID, positions.PageOptions{}); got != nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("page = %+v, error = %v, want expired deadline", got, err)
	} else {
		readerAssertErrorContext(t, err, "payroll", readerDocumentID)
	}
	if calls.Load() != 0 {
		t.Fatalf("validation made HTTP calls: %d", calls.Load())
	}
}

func TestServiceGetPageErrors(t *testing.T) {
	for _, tc := range []struct {
		name    string
		status  int
		body    string
		wantAPI bool
	}{
		{"unknown document type", 404, `{"errors":[{"error":"not found","code":1001}]}`, true},
		{"bad request", 400, `{"errors":[{"error":"invalid document","code":3000}]}`, true},
		{"unauthorized", 401, `{"errors":[{"error":"unauthorized"}]}`, true},
		{"forbidden", 403, `{"errors":[{"error":"forbidden"}]}`, true},
		{"rate limit", 429, `{"errors":[{"error":"slow down","code":1073}]}`, true},
		{"server error", 500, `{"errors":[{"error":"server failed"}]}`, true},
		{"HTML error", 502, "<html>bad gateway</html>", true},
		{"unexpected status", 201, `{"meta":{"href":"page-link"},"rows":[]}`, false},
		{"no content", 204, "", false},
		{"missing href", 200, `{"meta":{"type":"position"},"rows":[]}`, false},
		{"empty href", 200, `{"meta":{"href":""},"rows":[]}`, false},
		{"missing meta", 200, `{"rows":[]}`, false},
		{"null meta", 200, `{"meta":null,"rows":[]}`, false},
		{"missing rows", 200, `{"meta":{"href":"page-link"}}`, false},
		{"null rows", 200, `{"meta":{"href":"page-link"},"rows":null}`, false},
		{"object rows", 200, `{"meta":{"href":"page-link"},"rows":{}}`, false},
		{"string rows", 200, `{"meta":{"href":"page-link"},"rows":"[]"}`, false},
		{"partly decoded rows", 200, `{"meta":{"href":"page-link"},"rows":[{"label":"partial"},{"amount":"invalid"}]}`, false},
		{"empty object", 200, `{}`, false},
		{"null response", 200, `null`, false},
		{"array response", 200, `[]`, false},
		{"string response", 200, `"page"`, false},
		{"invalid JSON", 200, `{"meta":`, false},
		{"empty response", 200, "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			service := readerNewService(t, "future2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.URL.Path != "/api/remap/1.2/entity/future2/"+readerDocumentID+"/positions" {
					t.Errorf("path = %s", r.URL.Path)
				}
				w.WriteHeader(tc.status)
				if _, err := io.WriteString(w, tc.body); err != nil {
					t.Error(err)
				}
			}))
			got, err := service.GetPage(t.Context(), readerDocumentID, positions.PageOptions{})
			if err == nil || got != nil {
				t.Fatalf("page = %+v, error = %v", got, err)
			}
			readerAssertErrorContext(t, err, "future2", readerDocumentID)
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
				t.Fatalf("calls = %d, want 1 without retry", calls.Load())
			}
		})
	}
}

func TestServiceGetPageErrorWrapping(t *testing.T) {
	for _, authorization := range []bool{true, false} {
		name := "transport"
		if authorization {
			name = "authorization"
		}
		t.Run(name, func(t *testing.T) {
			wantErr := errors.New("boundary failed")
			var calls atomic.Int32
			cfg := transport.Config{
				BaseURL: "http://boundary.invalid",
				HTTPClient: &http.Client{Transport: readerRoundTripFunc(func(*http.Request) (*http.Response, error) {
					calls.Add(1)
					return nil, wantErr
				})},
			}
			if authorization {
				cfg.Auth = func(*http.Request) error { return wantErr }
			}
			client, err := transport.NewClient(cfg)
			if err != nil {
				t.Fatal(err)
			}
			service, err := positions.NewService[readerRow](client, "payroll")
			if err != nil {
				t.Fatal(err)
			}
			got, err := service.GetPage(t.Context(), readerDocumentID, positions.PageOptions{})
			if got != nil || !errors.Is(err, wantErr) {
				t.Fatalf("page = %+v, error = %v, want wrapped boundary error", got, err)
			}
			readerAssertErrorContext(t, err, "payroll", readerDocumentID)
			wantCalls := int32(1)
			if authorization {
				wantCalls = 0
			}
			if calls.Load() != wantCalls {
				t.Errorf("calls = %d, want %d", calls.Load(), wantCalls)
			}
		})
	}
}

func TestServiceGetPageCancellation(t *testing.T) {
	started := make(chan struct{})
	var calls atomic.Int32
	service := readerNewService(t, "payroll", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			close(started)
		}
		<-r.Context().Done()
	}))
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go func() {
		select {
		case <-started:
			cancel()
		case <-ctx.Done():
		}
	}()
	got, err := service.GetPage(ctx, readerDocumentID, positions.PageOptions{})
	if got != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("page = %+v, error = %v, want in-flight cancellation", got, err)
	}
	readerAssertErrorContext(t, err, "payroll", readerDocumentID)
	if calls.Load() != 1 {
		t.Fatalf("calls = %d, want 1 without retry after cancellation", calls.Load())
	}
}

func readerAssertErrorContext(t *testing.T, err error, documentType, documentID string) {
	t.Helper()
	if strings.Count(err.Error(), "moysklad:") != 1 ||
		!strings.Contains(err.Error(), "get "+documentType+" "+documentID+" positions") {
		t.Errorf("error context = %v", err)
	}
}

func readerWriteResponse(t *testing.T, w http.ResponseWriter, body string, compressed bool) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if compressed {
		w.Header().Set("Content-Encoding", "gzip")
		writer := gzip.NewWriter(w)
		if _, err := io.WriteString(writer, body); err != nil {
			t.Error(err)
		}
		if err := writer.Close(); err != nil {
			t.Error(err)
		}
		return
	}
	if _, err := io.WriteString(w, body); err != nil {
		t.Error(err)
	}
}

func readerNewService(t *testing.T, documentType string, handler http.Handler) *positions.Service[readerRow] {
	t.Helper()
	server := memnet.NewServer(handler)
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
	service, err := positions.NewService[readerRow](client, documentType)
	if err != nil {
		t.Fatal(err)
	}
	return service
}

type readerRoundTripFunc func(*http.Request) (*http.Response, error)

func (f readerRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
