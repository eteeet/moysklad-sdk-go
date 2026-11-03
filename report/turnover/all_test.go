package turnover_test

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
	"github.com/eteeet/moysklad-sdk-go/report/turnover"
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

func service(c *transport.Client) *turnover.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := turnover.NewService(c)
	return s
}

func page(rows ...[]byte) string {
	body := `{"context":{"employee":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/context/employee"}}},"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/report/turnover/all","type":"turnover","size":1,"limit":1000,"offset":0},"rows":[`
	for i, r := range rows {
		if i > 0 {
			body += ","
		}
		body += string(r)
	}
	return body + "]}"
}

func TestRowJSON(t *testing.T) {
	data := read(t, "upstream/report_turnover.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got turnover.Row
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			a := got.Assortment
			if a.Meta.Type != "product" || a.Name != "товар 1" || a.Code != "00001" || a.UOM == nil || a.UOM.Name != "шт" || a.ProductFolder == nil {
				t.Fatalf("row = %+v", got)
			}
			var page model.Collection[turnover.Row]
			if err := codec.Unmarshal(read(t, "custom/turnover/all.json"), &page); err != nil {
				t.Fatal(err)
			}
			if len(page.Rows) != 3 || page.Meta.Type != "turnover" || page.Context == nil || page.Rows[1].Assortment.Meta.Type != "variant" ||
				page.Rows[0].Income != (turnover.Metrics{Quantity: 3}) || page.Rows[0].OnPeriodEnd != (turnover.Metrics{Quantity: 1}) {
				t.Fatalf("page = %+v", page)
			}
			var sparse turnover.Row
			if err := codec.Unmarshal([]byte(`{"assortment":{"meta":{"href":"h","type":"variant"},"name":"v","image":{"meta":{"href":"i"}}},"income":{"sum":150.5,"quantity":0.25},"future":1}`), &sparse); err != nil ||
				sparse.Assortment.Image == nil || sparse.Income.Sum != 150.5 || sparse.Income.Quantity != 0.25 || sparse.Assortment.UOM != nil {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestAllHTTP(t *testing.T) {
	filter, err := query.Filter(query.Eq("store", "https://api.moysklad.ru/api/remap/1.2/entity/store/s1"), query.Eq("withoutturnover", "true"))
	if err != nil {
		t.Fatal(err)
	}
	options := turnover.AllOptions{
		Limit: 2, Offset: 4, GroupBy: turnover.GroupByVariant, Filter: filter,
		MomentFrom: documenttest.DateTime("2026-10-01 00:00:00.000"), MomentTo: documenttest.DateTime("2026-11-01 00:00:00.000"),
	}
	const q = "filter=store%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fstore%2Fs1%3Bwithoutturnover%3Dtrue&groupBy=variant&limit=2&momentFrom=2026-10-01+00%3A00%3A00&momentTo=2026-11-01+00%3A00%3A00&offset=4"
	for _, tc := range []struct {
		name    string
		options turnover.AllOptions
		query   string
	}{{"full", options, q}, {"defaults", turnover.AllOptions{}, "limit=1000&offset=0"},
		{"extra", turnover.AllOptions{Extra: url.Values{"x": {"1"}}}, "limit=1000&offset=0&x=1"}} {
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[turnover.Row], error) {
			return service(c).All(ctx, tc.options)
		}
		path := "report/turnover/all?" + tc.query
		for _, rc := range documenttest.ResponseCases(page(read(t, "upstream/report_turnover.json"))) {
			t.Run(tc.name+"/"+rc.Name, func(t *testing.T) {
				documenttest.HTTP(t, rc, "GET", path, nil, invoke, func(got *model.Collection[turnover.Row]) {
					if len(got.Rows) != 1 || got.Meta.Size != 1 || got.Rows[0].Assortment.Name != "товар 1" {
						t.Fatalf("page = %+v", got)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

// Зона расчёта периода уходит заголовком X-Lognex-Accept-Timezone, ответная — в ResponseInfo.
func TestAllTimezone(t *testing.T) {
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Lognex-Accept-Timezone") != "Sun, 01 Nov 2026 12:00:00 +0300" {
			t.Errorf("timezone = %q", r.Header.Get("X-Lognex-Accept-Timezone"))
		}
		w.Header().Set("X-Lognex-Content-Timezone", "Sun, 01 Nov 2026 12:00:00 +0300")
		w.Write([]byte(`{"meta":{"href":"h"},"rows":[]}`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	var info transport.ResponseInfo
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.FixedZone("", 3*3600))
	if _, err := service(c).All(t.Context(), turnover.AllOptions{}, transport.WithAcceptTimezone(now), transport.WithResponseInfo(&info)); err != nil {
		t.Fatal(err)
	}
	if got, ok := info.ContentTimezone(); !ok || !got.Equal(now) {
		t.Fatalf("content timezone = %v, %v", got, ok)
	}
}

func TestAllValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	for _, options := range []turnover.AllOptions{{Limit: -1}, {Limit: 1001}, {Offset: -1},
		{Extra: url.Values{"limit": {"5"}}}, {GroupBy: "product", Extra: url.Values{"groupBy": {"variant"}}}} {
		if _, err := service(c).All(t.Context(), options); err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "get turnover all") {
			t.Errorf("%+v: err = %v", options, err)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := turnover.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}
