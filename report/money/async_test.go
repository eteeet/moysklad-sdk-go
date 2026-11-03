package money_test

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
	"github.com/eteeet/moysklad-sdk-go/report/money"
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

func plotSeriesAsync(o money.PlotSeriesOptions) submitCall {
	return func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
		return service(c).PlotSeriesAsync(ctx, o, opts...)
	}
}

func TestAsyncHTTP(t *testing.T) {
	options := money.PlotSeriesOptions{MomentFrom: documenttest.DateTime("2018-09-01 00:00:00.000"), MomentTo: documenttest.DateTime("2018-09-04 00:00:00.000"),
		Interval: money.IntervalDay, Filter: "organization=" + organization, Extra: url.Values{"x": {"1"}}}
	testAsync(t, asyncCase{"plotseries", "report/money/plotseries",
		"async=true&filter=organization%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Forganization%2F850c8195-f504-11e5-8a84-bae50000015e&interval=day&momentFrom=2018-09-01+00%3A00%3A00&momentTo=2018-09-04+00%3A00%3A00&x=1",
		plotSeriesAsync(options)})
	testAsync(t, asyncCase{"byaccount", "report/money/byaccount", "async=true",
		func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
			return service(c).ByAccountAsync(ctx, opts...)
		}})
}

// Обязательные параметры, limit/offset и async в Extra — ошибка до HTTP.
func TestAsyncRejects(t *testing.T) {
	from, to := documenttest.DateTime("2018-09-01 00:00:00.000"), documenttest.DateTime("2018-09-04 00:00:00.000")
	valid := func(extra url.Values) money.PlotSeriesOptions {
		return money.PlotSeriesOptions{MomentFrom: from, MomentTo: to, Interval: money.IntervalDay, Extra: extra}
	}
	notAllowed := "get money plotseries async: limit and offset are not allowed with async"
	for name, tc := range map[string]struct {
		options money.PlotSeriesOptions
		want    string
	}{
		"momentFrom":   {money.PlotSeriesOptions{MomentTo: to, Interval: "day"}, "get money plotseries async: momentFrom is required"},
		"momentTo":     {money.PlotSeriesOptions{MomentFrom: from, Interval: "day"}, "get money plotseries async: momentTo is required"},
		"interval":     {money.PlotSeriesOptions{MomentFrom: from, MomentTo: to}, "get money plotseries async: interval is required"},
		"extra limit":  {valid(url.Values{"limit": {"1"}}), notAllowed},
		"extra offset": {valid(url.Values{"offset": {"1"}}), notAllowed},
		"extra async":  {valid(url.Values{"async": {"false"}}), "parameter async is not supported"},
	} {
		t.Run(name, func(t *testing.T) { testAsyncRejects(t, plotSeriesAsync(tc.options), tc.want) })
	}
}
