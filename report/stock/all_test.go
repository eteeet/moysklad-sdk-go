package stock_test

import (
	"context"
	"errors"
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
	"github.com/eteeet/moysklad-sdk-go/report/stock"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const base = "https://api.moysklad.ru/api/remap/1.2/"

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func service(c *transport.Client) *stock.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := stock.NewService(c)
	return s
}

func page(endpoint string, rows ...[]byte) string {
	body := `{"context":{"employee":{"meta":{"href":"` + base + `context/employee"}}},"meta":{"href":"` + base + endpoint + `","type":"stock","size":1,"limit":1000,"offset":0},"rows":[`
	for i, r := range rows {
		if i > 0 {
			body += ","
		}
		body += string(r)
	}
	return body + "]}"
}

// noCalls — клиент сервера, который считает запросы; локальные ошибки не должны их делать.
func noCalls(t *testing.T) (*transport.Client, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	t.Cleanup(server.Close)
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	return c, &calls
}

// localError проверяет ошибку до HTTP: один префикс moysklad:, контекст операции, сброшенный ResponseInfo.
func localError(t *testing.T, name, operation string, call func(transport.RequestOption) error) {
	t.Helper()
	info := transport.ResponseInfo{StatusCode: 999}
	if err := call(transport.WithResponseInfo(&info)); err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), operation) {
		t.Errorf("%s: err = %v", name, err)
	}
	if info.StatusCode != 0 {
		t.Errorf("%s: ResponseInfo not reset: %+v", name, info)
	}
}

func TestRowJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var product stock.Row
			if err := codec.Unmarshal(read(t, "upstream/stock_all_product.json"), &product); err != nil {
				t.Fatal(err)
			}
			if product.Meta.Type != "product" || product.Name != "Товар" || product.Code != "one1" || product.Article != "Ar23" ||
				product.ExternalCode != "ibWbdtWWhXiIwfZVEal6z2" || product.Stock != -30 || product.Quantity != -30 || product.Price != 0 ||
				product.SalePrice != 346347237000 || product.UOM.Name != "шт" || product.UOM.Meta.Type != "uom" || product.Folder != nil || product.Image != nil {
				t.Fatalf("product = %+v", product)
			}
			var variant stock.Row
			if err := codec.Unmarshal(read(t, "upstream/stock_all_variant.json"), &variant); err != nil {
				t.Fatal(err)
			}
			if variant.Meta.Type != "variant" || variant.Folder == nil || variant.Folder.Name != "Товары Из кастомной группы" || variant.Folder.PathName != "Группа" ||
				variant.Image == nil || variant.Image.Filename != "image-product.jpg" || variant.Image.Size != 7180 || variant.Image.Miniature == nil ||
				variant.Image.Miniature.DownloadHref == "" || variant.Image.Updated.IsZero() {
				t.Fatalf("variant = %+v", variant)
			}
			var consignment stock.Row
			if err := codec.Unmarshal(read(t, "upstream/stock_all_consignment.json"), &consignment); err != nil {
				t.Fatal(err)
			}
			if consignment.Meta.Type != "consignment" || consignment.StockDays != 45 || consignment.Price != 50260.86956521739 || consignment.SalePrice != 123000 ||
				consignment.Stock != 4 || consignment.Reserve != 0 || consignment.InTransit != 0 {
				t.Fatalf("consignment = %+v", consignment)
			}
			var all model.Collection[stock.Row]
			if err := codec.Unmarshal(read(t, "custom/stock/all.json"), &all); err != nil {
				t.Fatal(err)
			}
			last := all.Rows[len(all.Rows)-1]
			if len(all.Rows) != 16 || all.Meta.Type != "stock" || all.Meta.Size != 26 || all.Meta.NextHref == "" || all.Context == nil ||
				last.Image == nil || last.Folder == nil || all.Rows[1].Article != "" {
				t.Fatalf("page = %+v", all)
			}
			// Остатки — количества, могут быть дробными и отрицательными.
			var sparse stock.Row
			if err := codec.Unmarshal([]byte(`{"meta":{"href":"h","type":"future"},"stock":-1.5,"reserve":0.25,"inTransit":2,"quantity":-3.75,"future":1}`), &sparse); err != nil ||
				sparse.Meta.Type != "future" || sparse.Stock != -1.5 || sparse.Reserve != 0.25 || sparse.InTransit != 2 || sparse.Quantity != -3.75 {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

// Себестоимость и цена продажи — копейки float64: 0, 0.0, дробь и большое значение без потерь.
func TestRowMoneyJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for raw, want := range map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "346347237000.01": 346347237000.01} {
			var got stock.Row
			if err := codec.Unmarshal([]byte(`{"price":`+raw+`,"salePrice":`+raw+`}`), &got); err != nil || got.Price != want || got.SalePrice != want {
				t.Fatalf("%s: %s: got %+v, err = %v", codec.Name, raw, got, err)
			}
		}
	}
}

func TestAllHTTP(t *testing.T) {
	filter, err := query.Filter(query.Eq("product", base+"entity/product/p1"), query.Eq("stockMode", stock.StockModePositiveOnly))
	if err != nil {
		t.Fatal(err)
	}
	order, err := query.Order(query.Desc("stock"))
	if err != nil {
		t.Fatal(err)
	}
	full := stock.AllOptions{Limit: 2, Offset: 4, GroupBy: stock.GroupByConsignment, IncludeRelated: true, Filter: filter, Order: order}
	for _, tc := range []struct {
		name    string
		options stock.AllOptions
		query   string
	}{
		{"full", full, "filter=product%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fproduct%2Fp1%3BstockMode%3DpositiveOnly&groupBy=consignment&includeRelated=true&limit=2&offset=4&order=stock%2Cdesc"},
		{"defaults", stock.AllOptions{}, "limit=1000&offset=0"},
		{"unknown stockMode and extra", stock.AllOptions{Filter: "stockMode=future", Extra: url.Values{"x": {"1"}}}, "filter=stockMode%3Dfuture&limit=1000&offset=0&x=1"},
	} {
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[stock.Row], error) {
			return service(c).All(ctx, tc.options)
		}
		for _, rc := range documenttest.ResponseCases(page("report/stock/all", read(t, "upstream/stock_all_consignment.json"))) {
			t.Run(tc.name+"/"+rc.Name, func(t *testing.T) {
				documenttest.HTTP(t, rc, "GET", "report/stock/all?"+tc.query, nil, invoke, func(got *model.Collection[stock.Row]) {
					if len(got.Rows) != 1 || got.Meta.Size != 1 || got.Rows[0].StockDays != 45 || got.Rows[0].Price != 50260.86956521739 {
						t.Fatalf("page = %+v", got)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

func TestAllValidation(t *testing.T) {
	c, calls := noCalls(t)
	for name, options := range map[string]stock.AllOptions{
		"negative limit":                {Limit: -1},
		"limit over 1000":               {Limit: 1001},
		"negative offset":               {Offset: -1},
		"extra limit":                   {Extra: url.Values{"limit": {"5"}}},
		"extra order":                   {Extra: url.Values{"order": {"name"}}},
		"extra includeRelated":          {Extra: url.Values{"includeRelated": {"true"}}},
		"extra groupBy":                 {GroupBy: "product", Extra: url.Values{"groupBy": {"variant"}}},
		"extra filter":                  {Extra: url.Values{"filter": {"store=s"}}},
		"includeRelated without filter": {IncludeRelated: true},
		"includeRelated with store":     {IncludeRelated: true, Filter: "store=s;productFolder=f"},
		"includeRelated escaped name":   {IncludeRelated: true, Filter: `search=a\;product=p`},
	} {
		localError(t, name, "get stock all", func(o transport.RequestOption) error {
			_, err := service(c).All(t.Context(), options, o)
			return err
		})
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := stock.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

// includeRelated принимает и исключающее условие: filter=variant!=<URL>&includeRelated=true.
func TestAllIncludeRelatedNotEq(t *testing.T) {
	invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[stock.Row], error) {
		return service(c).All(ctx, stock.AllOptions{IncludeRelated: true, Filter: "variant!=v1"})
	}
	rc := documenttest.ResponseCases(page("report/stock/all"))[0]
	documenttest.HTTP(t, rc, "GET", "report/stock/all?filter=variant%21%3Dv1&includeRelated=true&limit=1000&offset=0", nil, invoke,
		func(got *model.Collection[stock.Row]) {
			if len(got.Rows) != 0 {
				t.Fatalf("page = %+v", got)
			}
		})
}

type frozenClock struct{ now time.Time }

func (c frozenClock) Now() time.Time                     { return c.now }
func (frozenClock) After(time.Duration) <-chan time.Time { return nil } // окно не освобождается

// All и ByStore тратят по 5 единиц лимита через transport.DefaultWeight: бюджет 10 они исчерпывают
// вдвоём, и следующий краткий отчёт (вес 1 у решения, 3 у сотрудника в ноябре 2026) ждёт окна;
// при любом другом весе All и ByStore он бы прошёл или заблокировался бы раньше.
func TestLimiterWeight(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if strings.HasSuffix(r.URL.Path, "/current") {
			w.Write([]byte(`[]`))
			return
		}
		w.Write([]byte(`{"meta":{"href":"h"},"rows":[]}`))
	}))
	defer server.Close()
	for _, profile := range []transport.Profile{transport.ProfileApp, transport.ProfileEmployee} {
		limiter, err := transport.NewLimiter(transport.LimiterConfig{Profile: profile, Budget: 10, Clock: frozenClock{time.Date(2026, 11, 3, 12, 0, 0, 0, time.UTC)}})
		if err != nil {
			t.Fatal(err)
		}
		c, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", Limiter: limiter})
		if err != nil {
			t.Fatal(err)
		}
		calls.Store(0)
		// Таймаут не даёт тесту зависнуть, если вес больше 5 и окно не освободится уже здесь.
		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		if _, err := service(c).All(ctx, stock.AllOptions{}); err != nil {
			t.Fatal(err)
		}
		if _, err := service(c).ByStore(ctx, stock.ByStoreOptions{}); err != nil {
			t.Fatal(err)
		}
		cancel()
		ctx, cancel = context.WithTimeout(t.Context(), 50*time.Millisecond)
		_, err = service(c).AllCurrent(ctx, stock.CurrentOptions{})
		cancel()
		if !errors.Is(err, context.DeadlineExceeded) || calls.Load() != 2 {
			t.Fatalf("profile %d: third request err = %v, calls = %d", profile, err, calls.Load())
		}
	}
	// Контроль: краткий отчёт весит 1 у решения — десять запросов помещаются в тот же бюджет.
	limiter, err := transport.NewLimiter(transport.LimiterConfig{Profile: transport.ProfileApp, Budget: 10, Clock: frozenClock{time.Date(2026, 11, 3, 12, 0, 0, 0, time.UTC)}})
	if err != nil {
		t.Fatal(err)
	}
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", Limiter: limiter})
	if err != nil {
		t.Fatal(err)
	}
	for range 10 {
		if _, err := service(c).AllCurrent(t.Context(), stock.CurrentOptions{}); err != nil {
			t.Fatal(err)
		}
	}
}
