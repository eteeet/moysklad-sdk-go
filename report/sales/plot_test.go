package sales_test

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/report/sales"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func service(c *transport.Client) *sales.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := sales.NewService(c)
	return s
}

type method func(*sales.Service, context.Context, sales.Options, ...transport.RequestOption) (*sales.Plot, error)

var moscow = time.FixedZone("MSK", 3*3600)

func point(hour int, quantity int, sum float64) sales.Point {
	return sales.Point{Date: model.DateTime(time.Date(2018, 9, 6, hour, 0, 0, 0, moscow)), Quantity: quantity, Sum: sum}
}

func equal(got, want []sales.Point) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if !time.Time(got[i].Date).Equal(time.Time(want[i].Date)) || got[i].Quantity != want[i].Quantity || got[i].Sum != want[i].Sum {
			return false
		}
	}
	return true
}

func TestPlotJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for _, tc := range []struct {
				file, metaType string
				want           []sales.Point
			}{
				{"upstream/report_orders_plotseries_list.json", "ordersplotseries", []sales.Point{point(0, 3, 600), point(1, 2, 200)}},
				{"upstream/report_sales_plotseries_list.json", "salesplotseries", []sales.Point{point(0, 3, 900), point(1, 4, 400)}},
				{"custom/sales/orders.json", "ordersplotseries", []sales.Point{point(0, 3, 600), point(1, 2, 200)}},
				{"custom/sales/orders_filter.json", "ordersplotseries", []sales.Point{point(0, 1, 100), point(1, 1, 100)}},
				{"custom/sales/sales.json", "salesplotseries", []sales.Point{point(0, 3, 900), point(1, 4, 400)}},
			} {
				var got sales.Plot
				if err := codec.Unmarshal(read(t, tc.file), &got); err != nil {
					t.Fatalf("%s: %v", tc.file, err)
				}
				if got.Meta.Type != model.MetaType(tc.metaType) || got.Meta.Href == "" || got.Context == nil || got.Context.Employee == nil || !equal(got.Series, tc.want) {
					t.Errorf("%s: plot = %+v", tc.file, got)
				}
			}
			// Сумма — копейки как пришли: 0, 0.0, дробь, большое значение.
			var sparse sales.Plot
			if err := codec.Unmarshal([]byte(`{"meta":{"href":"h"},"series":[{"date":"2026-11-01 00:00:00","quantity":0,"sum":0},`+
				`{"date":"2026-11-02 00:00:00","quantity":1,"sum":0.0},{"date":"2026-11-03 00:00:00","quantity":2,"sum":150.5,"future":1},`+
				`{"date":"2026-11-04 00:00:00","quantity":2147483648,"sum":313495109156366700.0}]}`), &sparse); err != nil {
				t.Fatal(err)
			}
			if len(sparse.Series) != 4 || sparse.Series[0].Sum != 0 || sparse.Series[1].Sum != 0 || sparse.Series[2].Sum != 150.5 ||
				sparse.Series[3].Sum != 313495109156366700.0 || sparse.Series[3].Quantity != 2147483648 || sparse.Context != nil {
				t.Fatalf("sparse = %+v", sparse)
			}
		})
	}
}

func testHTTP(t *testing.T, report, fixture string, call method, filter string, want []sales.Point) {
	options := sales.Options{Interval: sales.IntervalHour, Filter: filter,
		MomentFrom: documenttest.DateTime("2018-09-06 00:00:00.000"), MomentTo: documenttest.DateTime("2018-09-06 01:00:01.500")}
	base := "interval=hour&momentFrom=2018-09-06+00%3A00%3A00&momentTo=2018-09-06+01%3A00%3A01.500"
	for _, tc := range []struct {
		name    string
		options sales.Options
		query   string
	}{
		{"required", sales.Options{Interval: "day", MomentFrom: options.MomentFrom, MomentTo: options.MomentTo},
			"interval=day&momentFrom=2018-09-06+00%3A00%3A00&momentTo=2018-09-06+01%3A00%3A01.500"},
		{"filter", options, "filter=" + url.QueryEscape(filter) + "&" + base},
		{"extra", sales.Options{Interval: "month", MomentFrom: options.MomentFrom, MomentTo: options.MomentTo, Extra: url.Values{"x": {"1"}}},
			"interval=month&momentFrom=2018-09-06+00%3A00%3A00&momentTo=2018-09-06+01%3A00%3A01.500&x=1"},
	} {
		invoke := func(ctx context.Context, c *transport.Client) (*sales.Plot, error) {
			return call(service(c), ctx, tc.options)
		}
		for _, rc := range append(documenttest.ResponseCases(string(read(t, "custom/sales/"+fixture))),
			documenttest.ResponseCase{Name: "null series", Status: 200, Body: `{"meta":{"href":"h"},"series":null}`},
			documenttest.ResponseCase{Name: "no series", Status: 200, Body: `{"meta":{"href":"h"}}`}) {
			t.Run(tc.name+"/"+rc.Name, func(t *testing.T) {
				documenttest.HTTP(t, rc, "GET", "report/"+report+"/plotseries?"+tc.query, nil, invoke, func(got *sales.Plot) {
					if !equal(got.Series, want) || got.Meta.Type != model.MetaType(report+"plotseries") {
						t.Fatalf("plot = %+v", got)
					}
				})
			})
		}
		t.Run(tc.name+"/local", func(t *testing.T) { documenttest.LocalErrors(t, invoke) })
	}

	// Пустой series без meta — успех: обязательным при ответе источник называет только series.
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"series":[]}`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	if got, err := call(service(c), t.Context(), options); err != nil || got.Series == nil || len(got.Series) != 0 {
		t.Fatalf("empty series = %+v, %v", got, err)
	}
}

func TestOrdersHTTP(t *testing.T) {
	filter, err := query.Filter(query.Eq("organization", "https://api.moysklad.ru/api/remap/1.2/entity/organization/o1"),
		query.Eq("project", "https://api.moysklad.ru/api/remap/1.2/entity/project/p1"), query.Eq("store", "https://api.moysklad.ru/api/remap/1.2/entity/store/s1"))
	if err != nil {
		t.Fatal(err)
	}
	testHTTP(t, "orders", "orders.json", (*sales.Service).Orders, filter, []sales.Point{point(0, 3, 600), point(1, 2, 200)})
}

func TestSalesHTTP(t *testing.T) {
	filter, err := query.Filter(query.Eq("retailStore", "https://api.moysklad.ru/api/remap/1.2/entity/retailstore/r1"),
		query.Eq("store", "https://api.moysklad.ru/api/remap/1.2/entity/store/s1"), query.Eq("store", "https://api.moysklad.ru/api/remap/1.2/entity/store/s2"))
	if err != nil {
		t.Fatal(err)
	}
	testHTTP(t, "sales", "sales.json", (*sales.Service).Sales, filter, []sales.Point{point(0, 3, 900), point(1, 4, 400)})
}

// Зона расчёта уходит заголовком X-Lognex-Accept-Timezone, ответная — в ResponseInfo.
func TestTimezone(t *testing.T) {
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Lognex-Accept-Timezone") != "Sun, 01 Nov 2026 12:00:00 +0700" {
			t.Errorf("timezone = %q", r.Header.Get("X-Lognex-Accept-Timezone"))
		}
		w.Header().Set("X-Lognex-Content-Timezone", "Sun, 01 Nov 2026 12:00:00 +0700")
		w.Write([]byte(`{"meta":{"href":"h"},"series":[]}`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	var info transport.ResponseInfo
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.FixedZone("", 7*3600))
	options := sales.Options{Interval: sales.IntervalDay, MomentFrom: documenttest.DateTime("2026-10-01 00:00:00.000"), MomentTo: documenttest.DateTime("2026-11-01 00:00:00.000")}
	if _, err := service(c).Sales(t.Context(), options, transport.WithAcceptTimezone(now), transport.WithResponseInfo(&info)); err != nil {
		t.Fatal(err)
	}
	if got, ok := info.ContentTimezone(); !ok || !got.Equal(now) {
		t.Fatalf("content timezone = %v, %v", got, ok)
	}
}

func TestValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	from, to := documenttest.DateTime("2026-10-01 00:00:00.000"), documenttest.DateTime("2026-11-01 00:00:00.000")
	for _, m := range []struct {
		name string
		call method
	}{{"get orders plotseries", (*sales.Service).Orders}, {"get sales plotseries", (*sales.Service).Sales}} {
		for _, options := range []sales.Options{{}, {MomentTo: to, Interval: "day"}, {MomentFrom: from, Interval: "day"}, {MomentFrom: from, MomentTo: to},
			{MomentFrom: from, MomentTo: to, Interval: "day", Extra: url.Values{"interval": {"hour"}}},
			{MomentFrom: from, MomentTo: to, Interval: "day", Extra: url.Values{"momentFrom": {"2026-01-01 00:00:00"}}},
			{MomentFrom: from, MomentTo: to, Interval: "day", Extra: url.Values{"async": {"true"}}}} {
			var info transport.ResponseInfo
			info.StatusCode = 299
			got, err := m.call(service(c), t.Context(), options, transport.WithResponseInfo(&info))
			if got != nil || err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), m.name) {
				t.Errorf("%s %+v: got = %v, err = %v", m.name, options, got, err)
			}
			if info.StatusCode != 0 {
				t.Errorf("%s %+v: ResponseInfo not reset", m.name, options)
			}
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := sales.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}
