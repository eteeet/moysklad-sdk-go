package stock_test

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
	"github.com/eteeet/moysklad-sdk-go/report/stock"
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

const (
	assortmentID = "12345678-5838-aaeb-0a80-003a003ef439"
	storeID      = "12345678-b123-aaee-0a80-012b0001bb10"
)

// Async-методы всех GET-разделов «Отчета Остатки»: те же параметры, что у синхронных, плюс
// async=true, без limit и offset.
func TestAsyncHTTP(t *testing.T) {
	product := "product=" + base + "entity/product/" + assortmentID
	for _, tc := range []asyncCase{
		{"all", "report/stock/all", "async=true&filter=product%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fproduct%2F12345678-5838-aaeb-0a80-003a003ef439&groupBy=consignment&includeRelated=true&order=name%2Casc&x=1",
			func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
				return service(c).AllAsync(ctx, stock.AllOptions{GroupBy: stock.GroupByConsignment, IncludeRelated: true, Filter: product,
					Order: "name,asc", Extra: url.Values{"x": {"1"}}}, opts...)
			}},
		{"all empty", "report/stock/all", "async=true",
			func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
				return service(c).AllAsync(ctx, stock.AllOptions{}, opts...)
			}},
		{"bystore", "report/stock/bystore", "async=true&filter=stockMode%3Dall&groupBy=product&order=name",
			func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
				return service(c).ByStoreAsync(ctx, stock.ByStoreOptions{GroupBy: stock.GroupByProduct, Filter: "stockMode=all", Order: "name"}, opts...)
			}},
		{"all current", "report/stock/all/current", "async=true&filter=assortmentId%3D12345678-5838-aaeb-0a80-003a003ef439&include=zeroLines&stockType=quantity",
			func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
				return service(c).AllCurrentAsync(ctx, stock.CurrentOptions{StockType: stock.StockTypeQuantity, Include: stock.IncludeZeroLines,
					AssortmentIDs: []string{assortmentID}}, opts...)
			}},
		{"bystore current", "report/stock/bystore/current", "async=true&filter=storeId%3D12345678-b123-aaee-0a80-012b0001bb10&stockType=reserve&withRecalculate=true",
			func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
				return service(c).ByStoreCurrentAsync(ctx, stock.CurrentOptions{StockType: stock.StockTypeReserve, WithRecalculate: true,
					StoreIDs: []string{storeID}}, opts...)
			}},
		{"byslot current", "report/stock/byslot/current", "async=true&filter=assortmentId%3D12345678-5838-aaeb-0a80-003a003ef439%3BstoreId%3D12345678-b123-aaee-0a80-012b0001bb10",
			func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
				return service(c).BySlotCurrentAsync(ctx, stock.SlotOptions{AssortmentIDs: []string{assortmentID}, StoreIDs: []string{storeID}}, opts...)
			}},
		{"byoperation", "report/stock/byoperation", "async=true&operation.id=" + orderID,
			func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
				return service(c).ByOperationAsync(ctx, orderID, opts...)
			}},
	} {
		t.Run(tc.name, func(t *testing.T) { testAsync(t, tc) })
	}
}

// Limit/Offset (полем или в Extra), async в Extra и проверки синхронных методов — до HTTP.
func TestAsyncRejects(t *testing.T) {
	all := func(o stock.AllOptions) submitCall {
		return func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
			return service(c).AllAsync(ctx, o, opts...)
		}
	}
	byStore := func(o stock.ByStoreOptions) submitCall {
		return func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
			return service(c).ByStoreAsync(ctx, o, opts...)
		}
	}
	allCurrent := func(o stock.CurrentOptions) submitCall {
		return func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
			return service(c).AllCurrentAsync(ctx, o, opts...)
		}
	}
	storeCurrent := func(o stock.CurrentOptions) submitCall {
		return func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
			return service(c).ByStoreCurrentAsync(ctx, o, opts...)
		}
	}
	slot := func(o stock.SlotOptions) submitCall {
		return func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
			return service(c).BySlotCurrentAsync(ctx, o, opts...)
		}
	}
	notAllowed := "get stock all async: limit and offset are not allowed with async"
	for name, tc := range map[string]struct {
		call submitCall
		want string
	}{
		"all limit":              {all(stock.AllOptions{Limit: 10}), notAllowed},
		"all offset":             {all(stock.AllOptions{Offset: 10}), notAllowed},
		"all extra limit":        {all(stock.AllOptions{Extra: url.Values{"limit": {"10"}}}), "parameter limit must not be set in Extra"},
		"all extra async":        {all(stock.AllOptions{Extra: url.Values{"async": {"false"}}}), "parameter async is not supported"},
		"all includeRelated":     {all(stock.AllOptions{IncludeRelated: true}), "includeRelated requires a product or variant filter"},
		"bystore limit":          {byStore(stock.ByStoreOptions{Limit: 1}), "get stock bystore async: limit and offset are not allowed with async"},
		"bystore offset":         {byStore(stock.ByStoreOptions{Offset: 1}), "get stock bystore async: limit and offset are not allowed with async"},
		"bystore extra offset":   {byStore(stock.ByStoreOptions{Extra: url.Values{"offset": {"1"}}}), "parameter offset must not be set in Extra"},
		"all current extra":      {allCurrent(stock.CurrentOptions{Extra: url.Values{"limit": {"1"}}}), "get stock all current async: limit and offset are not allowed with async"},
		"all current recalc":     {allCurrent(stock.CurrentOptions{WithRecalculate: true}), "withRecalculate is supported only by bystore/current"},
		"bystore current offset": {storeCurrent(stock.CurrentOptions{Extra: url.Values{"offset": {"1"}}}), "get stock bystore current async: limit and offset are not allowed with async"},
		"bystore current bad id": {storeCurrent(stock.CurrentOptions{StoreIDs: []string{"x"}}), `storeId "x" must be a UUID`},
		"byslot no ids":          {slot(stock.SlotOptions{}), "get stock byslot current async: assortment or store IDs are required"},
		"byslot extra limit":     {slot(stock.SlotOptions{StoreIDs: []string{storeID}, Extra: url.Values{"limit": {"1"}}}), "get stock byslot current async: limit and offset are not allowed with async"},
		"byoperation id": {func(ctx context.Context, c *transport.Client, opts ...transport.RequestOption) (*async.Submitted, error) {
			return service(c).ByOperationAsync(ctx, "x", opts...)
		}, "get stock byoperation async: operation ID must be a UUID"},
	} {
		t.Run(name, func(t *testing.T) { testAsyncRejects(t, tc.call, tc.want) })
	}
}
