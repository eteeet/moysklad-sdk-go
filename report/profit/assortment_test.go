package profit_test

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/report/profit"
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

func service(c *transport.Client) *profit.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := profit.NewService(c)
	return s
}

func decode[T any](t *testing.T, codec documenttest.Codec, data []byte) T {
	t.Helper()
	var v T
	if err := codec.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}
	return v
}

type method[T any] func(*profit.Service, context.Context, profit.Options, ...transport.RequestOption) (*model.Collection[T], error)

// testHTTP проверяет один разрез: точный путь и query с полными и пустыми Options и Extra,
// ответ источника plain/gzip, ошибки API и локальные ошибки, ровно одну отправку.
func testHTTP[T any](t *testing.T, slice, fixture string, call method[T], check func(*testing.T, *model.Collection[T])) {
	filter, err := query.Filter(
		query.Eq("productFolder", "https://api.moysklad.ru/api/remap/1.2/entity/productfolder/f1"), query.Eq("withSubFolders", "false"),
		query.Eq("salesChannel", "https://api.moysklad.ru/api/remap/1.2/entity/saleschannel/c1"),
		query.Eq("salesChannel", "https://api.moysklad.ru/api/remap/1.2/entity/saleschannel/c2"))
	if err != nil {
		t.Fatal(err)
	}
	full := profit.Options{Limit: 2, Offset: 4, Filter: filter,
		MomentFrom: documenttest.DateTime("2026-10-01 00:00:00.000"), MomentTo: documenttest.DateTime("2026-11-01 12:30:00.250")}
	const q = "filter=productFolder%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fproductfolder%2Ff1%3BwithSubFolders%3Dfalse%3BsalesChannel%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fsaleschannel%2Fc1%3BsalesChannel%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fsaleschannel%2Fc2&limit=2&momentFrom=2026-10-01+00%3A00%3A00&momentTo=2026-11-01+12%3A30%3A00.250&offset=4"
	for _, tc := range []struct {
		name    string
		options profit.Options
		query   string
	}{{"full", full, q}, {"defaults", profit.Options{}, "limit=1000&offset=0"},
		{"extra", profit.Options{Extra: url.Values{"x": {"1", "2"}}}, "limit=1000&offset=0&x=1&x=2"}} {
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[T], error) {
			return call(service(c), ctx, tc.options)
		}
		for _, rc := range documenttest.ResponseCases(string(read(t, "custom/profit/"+fixture))) {
			t.Run(tc.name+"/"+rc.Name, func(t *testing.T) {
				documenttest.HTTP(t, rc, "GET", "report/profit/"+slice+"?"+tc.query, nil, invoke, func(got *model.Collection[T]) { check(t, got) })
			})
		}
		t.Run(tc.name+"/local", func(t *testing.T) { documenttest.LocalErrors(t, invoke) })
	}

	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	for _, options := range []profit.Options{{Limit: -1}, {Limit: 1001}, {Offset: -1},
		{Extra: url.Values{"limit": {"5"}}}, {Extra: url.Values{"offset": {"5"}}},
		{Filter: "store=s", Extra: url.Values{"filter": {"x"}}}, {Extra: url.Values{"expand": {"assortment"}}}} {
		var info transport.ResponseInfo
		info.StatusCode = 299
		got, err := call(service(c), t.Context(), options, transport.WithResponseInfo(&info))
		if got != nil || err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "get profit "+slice) {
			t.Errorf("%+v: got = %v, err = %v", options, got, err)
		}
		if info.StatusCode != 0 {
			t.Errorf("%+v: ResponseInfo not reset: %+v", options, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestAssortmentJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			a := decode[profit.Assortment](t, codec, read(t, "upstream/report_profit_assortment.json"))
			if a.Meta.Type != "product" || a.Name != "Курево" || a.Code != "00067" || a.Article != "A-001" ||
				a.UOM == nil || a.UOM.Name != "шт" || a.UOM.Meta.Type != "uom" ||
				!strings.Contains(string(a.Image), "/images/a12f1111-1111-1111-1111-111111111111") {
				t.Fatalf("assortment = %+v", a)
			}

			product := decode[model.Collection[profit.AssortmentRow]](t, codec, read(t, "upstream/report_profit_by_product_list.json"))
			want := profit.AssortmentRow{SellQuantity: 10, SellPrice: 500000, SellCost: 300000, SellSum: 5000000, SellCostSum: 3000000,
				ReturnQuantity: 3, ReturnPrice: 500000, ReturnCost: 300000, ReturnSum: 1500000, ReturnCostSum: 900000,
				Profit: 1400000, Margin: 0.6666666666666666, SalesMargin: 0.4}
			if len(product.Rows) != 1 || product.Meta.Type != "salesbyproduct" || product.Context == nil {
				t.Fatalf("product = %+v", product)
			}
			got := product.Rows[0]
			if got.Assortment.Article != "A-001" || got.Assortment.Image != nil {
				t.Fatalf("product assortment = %+v", got.Assortment)
			}
			got.Assortment = profit.Assortment{}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("product row = %+v", got)
			}

			variant := decode[model.Collection[profit.AssortmentRow]](t, codec, read(t, "upstream/report_profit_by_variant_list.json"))
			if len(variant.Rows) != 1 || variant.Meta.Type != "salesbyvariant" || variant.Rows[0].Assortment.Meta.Type != "variant" ||
				variant.Rows[0].Assortment.Article != "V-001" || variant.Rows[0].ReturnQuantity != 1 || variant.Rows[0].SalesMargin != 0.3003472222222222 {
				t.Fatalf("variant = %+v", variant)
			}

			// Деньги — копейки как пришли: 0, 0.0, дробь и большое значение.
			sparse := decode[profit.AssortmentRow](t, codec, []byte(`{"assortment":{"meta":{"href":"h","type":"bundle"},"name":"b"},`+
				`"sellSum":0,"sellPrice":0.0,"profit":150.5,"sellCostSum":313495109156366700.0,"sellQuantity":0.25,"future":1}`))
			if sparse.SellSum != 0 || sparse.SellPrice != 0 || sparse.Profit != 150.5 || sparse.SellCostSum != 313495109156366700.0 ||
				sparse.SellQuantity != 0.25 || sparse.Assortment.UOM != nil || sparse.Assortment.Meta.Type != "bundle" {
				t.Fatalf("sparse = %+v", sparse)
			}
		})
	}
}

func TestByProductHTTP(t *testing.T) {
	testHTTP(t, "byproduct", "byproduct.json", (*profit.Service).ByProduct, func(t *testing.T, got *model.Collection[profit.AssortmentRow]) {
		if len(got.Rows) != 3 || got.Meta.Size != 5 || got.Meta.Limit != 2 || got.Meta.NextHref == "" ||
			got.Rows[1].SellSum != 476500000000 || got.Rows[2].Assortment.Meta.Type != "service" || got.Rows[2].Margin != 0.6666666666666666 {
			t.Fatalf("page = %+v", got)
		}
	})
}

func TestByVariantHTTP(t *testing.T) {
	testHTTP(t, "byvariant", "byvariant.json", (*profit.Service).ByVariant, func(t *testing.T, got *model.Collection[profit.AssortmentRow]) {
		if len(got.Rows) != 3 || got.Meta.Type != "salesbyvariant" || got.Rows[0].Assortment.Meta.Type != "variant" ||
			got.Rows[0].Assortment.UOM == nil || got.Rows[0].SellPrice != 216174.97396737244 {
			t.Fatalf("page = %+v", got)
		}
	})
}
