package profit_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/async"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/report/profit"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const asyncTask = "498b8673-0308-11e6-9464-e4de00000089"

type submitCall func(context.Context, *transport.Client, ...transport.RequestOption) (*async.Submitted, error)

// asyncCase — метод *Async и точный запрос, который он должен отправить.
type asyncCase struct {
	name, path, query string
	call              submitCall
}

// testAsync: один GET с точными path, query и заголовками, без тела. 202 с Content-Location и
// Location — задача; 200, ошибка API, 503 при включённых повторах и 202 без адресов — ошибка
// после ровно одной отправки; context, auth и transport — documenttest.LocalErrors.
func testAsync(t *testing.T, tc asyncCase) {
	t.Helper()
	var calls, status atomic.Int32
	const noHeaders = -1 // 202 без адресов задачи
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != http.MethodGet || r.URL.Path != "/api/remap/1.2/"+tc.path || r.URL.RawQuery != tc.query {
			t.Errorf("request = %s %s, want query %q", r.Method, r.URL, tc.query)
		}
		for key, want := range map[string]string{"Authorization": "Bearer test-token", "Accept": "application/json;charset=utf-8",
			"Accept-Encoding": "gzip", "Content-Type": "", "User-Agent": "moysklad-sdk-go"} {
			if r.Header.Get(key) != want {
				t.Errorf("header %s = %q", key, r.Header.Get(key))
			}
		}
		if data, err := io.ReadAll(r.Body); err != nil || len(data) != 0 {
			t.Errorf("body = %q, err = %v", data, err)
		}
		task := "http://" + r.Host + "/api/remap/1.2/async/" + asyncTask
		switch code := int(status.Load()); code {
		case http.StatusAccepted:
			w.Header().Set("Content-Location", task)
			w.Header().Set("Location", task+"/result")
			w.WriteHeader(code)
		case noHeaders:
			w.WriteHeader(http.StatusAccepted)
		case http.StatusBadRequest:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			io.WriteString(w, `{"errors":[{"error":"bad","code":1001}]}`)
		default:
			w.WriteHeader(code)
			io.WriteString(w, `{}`)
		}
	}))
	defer server.Close()
	auth, err := transport.BearerToken("test-token")
	if err != nil {
		t.Fatal(err)
	}
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", Auth: auth, Retry: &transport.RetryPolicy{MaxAttempts: 3}})
	if err != nil {
		t.Fatal(err)
	}
	status.Store(http.StatusAccepted)
	info := transport.ResponseInfo{}
	task, err := tc.call(t.Context(), c, transport.WithResponseInfo(&info))
	want := server.URL + "/api/remap/1.2/async/" + asyncTask
	if err != nil || task == nil || task.ID != asyncTask || task.StatusHref != want || task.ResultHref != want+"/result" ||
		info.StatusCode != http.StatusAccepted || calls.Load() != 1 {
		t.Fatalf("task = %+v, err = %v, info = %+v, calls = %d", task, err, info, calls.Load())
	}
	for _, code := range []int{http.StatusOK, http.StatusBadRequest, http.StatusServiceUnavailable, noHeaders} {
		calls.Store(0)
		status.Store(int32(code))
		task, err := tc.call(t.Context(), c)
		if task != nil || err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "async") || calls.Load() != 1 {
			t.Fatalf("status %d: task = %+v, err = %v, calls = %d", code, task, err, calls.Load())
		}
		var apiErr *sdkerrors.APIError
		if errors.As(err, &apiErr) != (code == http.StatusBadRequest || code == http.StatusServiceUnavailable) ||
			code == http.StatusBadRequest && (apiErr.StatusCode != code || len(apiErr.Errors) != 1 || apiErr.Errors[0].Code != 1001) {
			t.Fatalf("status %d: API error = %v", code, err)
		}
	}
	documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*async.Submitted, error) { return tc.call(ctx, c) })
}

// testAsyncRejects: ошибка до HTTP — calls == 0, ResponseInfo сброшен, один префикс moysklad:,
// текст причины want.
func testAsyncRejects(t *testing.T, call submitCall, want string) {
	t.Helper()
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	info := transport.ResponseInfo{StatusCode: 999}
	task, err := call(t.Context(), c, transport.WithResponseInfo(&info))
	if task != nil || err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), want) ||
		info.StatusCode != 0 || calls.Load() != 0 {
		t.Fatalf("task = %+v, err = %v, info = %+v, calls = %d", task, err, info, calls.Load())
	}
}

// Все пять разрезов: параметры Options и async=true, без limit и offset.
func TestAsyncHTTP(t *testing.T) {
	options := profit.Options{MomentFrom: documenttest.DateTime("2026-10-01 00:00:00.000"), MomentTo: documenttest.DateTime("2026-11-01 12:30:00.250"),
		Filter: "salesChannel=https://api.moysklad.ru/api/remap/1.2/entity/saleschannel/c1", Extra: url.Values{"x": {"1"}}}
	const q = "async=true&filter=salesChannel%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fsaleschannel%2Fc1&momentFrom=2026-10-01+00%3A00%3A00&momentTo=2026-11-01+12%3A30%3A00.250&x=1"
	for slice, method := range methods() {
		t.Run(slice, func(t *testing.T) {
			testAsync(t, asyncCase{slice, "report/profit/" + slice, q, method(options)})
		})
	}
	testAsync(t, asyncCase{"empty", "report/profit/byproduct", "async=true", methods()["byproduct"](profit.Options{})})
}

func methods() map[string]func(profit.Options) submitCall {
	bind := func(m func(*profit.Service) func(context.Context, profit.Options, ...transport.RequestOption) (*async.Submitted, error)) func(profit.Options) submitCall {
		return func(o profit.Options) submitCall {
			return func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
				return m(service(c))(ctx, o, opts...)
			}
		}
	}
	return map[string]func(profit.Options) submitCall{
		"byproduct": bind(func(s *profit.Service) func(context.Context, profit.Options, ...transport.RequestOption) (*async.Submitted, error) {
			return s.ByProductAsync
		}),
		"byvariant": bind(func(s *profit.Service) func(context.Context, profit.Options, ...transport.RequestOption) (*async.Submitted, error) {
			return s.ByVariantAsync
		}),
		"byemployee": bind(func(s *profit.Service) func(context.Context, profit.Options, ...transport.RequestOption) (*async.Submitted, error) {
			return s.ByEmployeeAsync
		}),
		"bycounterparty": bind(func(s *profit.Service) func(context.Context, profit.Options, ...transport.RequestOption) (*async.Submitted, error) {
			return s.ByCounterpartyAsync
		}),
		"bysaleschannel": bind(func(s *profit.Service) func(context.Context, profit.Options, ...transport.RequestOption) (*async.Submitted, error) {
			return s.BySalesChannelAsync
		}),
	}
}

// Limit/Offset полем или в Extra и async в Extra — ошибка до HTTP.
func TestAsyncRejects(t *testing.T) {
	for slice, method := range methods() {
		notAllowed := "get profit " + slice + " async: limit and offset are not allowed with async"
		for name, tc := range map[string]struct {
			options profit.Options
			want    string
		}{
			"limit":        {profit.Options{Limit: 10}, notAllowed},
			"offset":       {profit.Options{Offset: 10}, notAllowed},
			"extra limit":  {profit.Options{Extra: url.Values{"limit": {"10"}}}, notAllowed},
			"extra offset": {profit.Options{Extra: url.Values{"offset": {"10"}}}, notAllowed},
			"extra async":  {profit.Options{Extra: url.Values{"async": {"false"}}}, "parameter async is not supported"},
			"twice":        {profit.Options{Filter: "a=b", Extra: url.Values{"filter": {"c=d"}}}, "parameter filter is set twice"},
		} {
			t.Run(slice+" "+name, func(t *testing.T) { testAsyncRejects(t, method(tc.options), tc.want) })
		}
	}
}
