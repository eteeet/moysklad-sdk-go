package stock_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/report/stock"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestStoreRowJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got stock.StoreRow
			if err := codec.Unmarshal(read(t, "upstream/stock_by_store.json"), &got); err != nil {
				t.Fatal(err)
			}
			if got.Meta.Type != "product" || len(got.StockByStore) != 2 ||
				got.StockByStore[0] != (stock.StoreStock{Meta: got.StockByStore[0].Meta, Name: "Не основной склад", Stock: -30}) ||
				got.StockByStore[0].Meta.Type != "store" || got.StockByStore[1].Name != "Основной склад" {
				t.Fatalf("row = %+v", got)
			}
			var all model.Collection[stock.StoreRow]
			if err := codec.Unmarshal(read(t, "custom/stock/bystore.json"), &all); err != nil {
				t.Fatal(err)
			}
			if len(all.Rows) != 7 || all.Meta.Type != "stockbystore" || all.Meta.Size != 225 || all.Context == nil ||
				len(all.Rows[0].StockByStore) != 2 || all.Rows[0].StockByStore[1].Meta.Href == "" {
				t.Fatalf("page = %+v", all)
			}
			var sparse stock.StoreRow
			if err := codec.Unmarshal([]byte(`{"meta":{"href":"h","type":"consignment"},"stockByStore":[{"meta":{"href":"s"},"name":"S","stock":1.5,"reserve":0.5,"inTransit":-2,"future":1}]}`), &sparse); err != nil ||
				sparse.StockByStore[0] != (stock.StoreStock{Meta: model.Meta{Href: "s"}, Name: "S", Stock: 1.5, Reserve: 0.5, InTransit: -2}) {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestByStoreHTTP(t *testing.T) {
	filter, err := query.Filter(query.Eq("store", base+"entity/store/s1"), query.Eq("stockMode", stock.StockModeAll))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name    string
		options stock.ByStoreOptions
		query   string
	}{
		{"full", stock.ByStoreOptions{Limit: 25, Offset: 25, GroupBy: stock.GroupByProduct, Filter: filter, Order: "stockOnAllStores,desc"},
			"filter=store%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fstore%2Fs1%3BstockMode%3Dall&groupBy=product&limit=25&offset=25&order=stockOnAllStores%2Cdesc"},
		{"defaults", stock.ByStoreOptions{}, "limit=1000&offset=0"},
		{"unknown groupBy and extra", stock.ByStoreOptions{GroupBy: "future", Extra: url.Values{"x": {"1"}}}, "groupBy=future&limit=1000&offset=0&x=1"},
	} {
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[stock.StoreRow], error) {
			return service(c).ByStore(ctx, tc.options)
		}
		for _, rc := range documenttest.ResponseCases(page("report/stock/bystore", read(t, "upstream/stock_by_store.json"))) {
			t.Run(tc.name+"/"+rc.Name, func(t *testing.T) {
				documenttest.HTTP(t, rc, "GET", "report/stock/bystore?"+tc.query, nil, invoke, func(got *model.Collection[stock.StoreRow]) {
					if len(got.Rows) != 1 || len(got.Rows[0].StockByStore) != 2 || got.Rows[0].StockByStore[0].Stock != -30 {
						t.Fatalf("page = %+v", got)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

func TestByStoreValidation(t *testing.T) {
	c, calls := noCalls(t)
	for name, options := range map[string]stock.ByStoreOptions{
		"limit over 1000":      {Limit: 1001},
		"negative offset":      {Offset: -1},
		"extra offset":         {Extra: url.Values{"offset": {"5"}}},
		"extra includeRelated": {Extra: url.Values{"includeRelated": {"true"}}},
		"extra filter":         {Filter: "store=s", Extra: url.Values{"filter": {"store=s"}}},
		"extra expand":         {Extra: url.Values{"expand": {"meta"}}},
	} {
		localError(t, name, "get stock bystore", func(o transport.RequestOption) error {
			_, err := service(c).ByStore(t.Context(), options, o)
			return err
		})
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}
