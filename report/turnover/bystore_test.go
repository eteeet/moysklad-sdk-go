package turnover_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/report/turnover"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func ref(path, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/" + path, Type: model.MetaType(typ)}}
}

func TestStoreRowJSON(t *testing.T) {
	data := read(t, "upstream/report_turnover_by_store.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got turnover.StoreRow
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.Assortment.Meta.Type != "variant" || len(got.StockByStore) != 2 || got.StockByStore[0].Store.Name != "Основной склад" ||
				got.StockByStore[0].OnPeriodEnd != (turnover.Metrics{Sum: 20000, Quantity: 2}) {
				t.Fatalf("row = %+v", got)
			}
			var page model.Collection[turnover.StoreRow]
			if err := codec.Unmarshal(read(t, "custom/turnover/bystore.json"), &page); err != nil {
				t.Fatal(err)
			}
			if len(page.Rows) != 2 || page.Context == nil || len(page.Rows[0].StockByStore) == 0 || page.Rows[0].StockByStore[0].Store.Meta.Href == "" {
				t.Fatalf("page = %+v", page)
			}
			var sparse turnover.StoreRow
			if err := codec.Unmarshal([]byte(`{"assortment":{"meta":{"href":"h"},"name":"x"},"stockByStore":[{"store":{"meta":{"href":"s"},"name":"S"},"income":{"sum":1,"quantity":1},"future":1}]}`), &sparse); err != nil ||
				sparse.StockByStore[0].Income != (turnover.Metrics{Sum: 1, Quantity: 1}) || sparse.StockByStore[0].Outcome != (turnover.Metrics{}) {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

// Зона расчёта периода уходит заголовком и у отчёта по складам.
func TestByStoreTimezone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Lognex-Accept-Timezone") != "Sun, 01 Nov 2026 12:00:00 +0700" {
			t.Errorf("timezone = %q", r.Header.Get("X-Lognex-Accept-Timezone"))
		}
		w.Write([]byte(`{"meta":{"href":"h"},"rows":[]}`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.FixedZone("", 7*3600))
	if _, err := service(c).ByStore(t.Context(), turnover.DetailOptions{Assortment: ref("entity/product/p1", "product")}, transport.WithAcceptTimezone(now)); err != nil {
		t.Fatal(err)
	}
}

func TestByStoreHTTP(t *testing.T) {
	store, err := query.Filter(query.Eq("store", "https://api.moysklad.ru/api/remap/1.2/entity/store/s1"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name    string
		options turnover.DetailOptions
		query   string
	}{
		{"product", turnover.DetailOptions{Assortment: ref("entity/product/p1", "product")},
			"filter=product%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fproduct%2Fp1"},
		{"variant with period and store", turnover.DetailOptions{Assortment: ref("entity/variant/v1", "variant"), Filter: store,
			MomentFrom: documenttest.DateTime("2026-10-01 00:00:00.000"), MomentTo: documenttest.DateTime("2026-11-01 00:00:00.000")},
			"filter=variant%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fvariant%2Fv1%3Bstore%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fstore%2Fs1&momentFrom=2026-10-01+00%3A00%3A00&momentTo=2026-11-01+00%3A00%3A00"},
	} {
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[turnover.StoreRow], error) {
			return service(c).ByStore(ctx, tc.options)
		}
		for _, rc := range documenttest.ResponseCases(page(read(t, "upstream/report_turnover_by_store.json"))) {
			t.Run(tc.name+"/"+rc.Name, func(t *testing.T) {
				documenttest.HTTP(t, rc, "GET", "report/turnover/bystore?"+tc.query, nil, invoke, func(got *model.Collection[turnover.StoreRow]) {
					if len(got.Rows) != 1 || len(got.Rows[0].StockByStore) != 2 {
						t.Fatalf("page = %+v", got)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

// Без товара или модификации запрос не уходит: сервер его не выполнит.
func TestByStoreValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	for name, options := range map[string]turnover.DetailOptions{
		"no assortment":          {},
		"service":                {Assortment: ref("entity/service/s1", "service")},
		"no href":                {Assortment: model.Reference{Meta: model.Meta{Type: "product"}}},
		"filter repeats product": {Assortment: ref("entity/product/p1", "product"), Filter: "store=s;product=p2"},
		"filter repeats variant": {Assortment: ref("entity/product/p1", "product"), Filter: "variant=v1"},
	} {
		if _, err := service(c).ByStore(t.Context(), options); err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "turnover bystore") {
			t.Errorf("%s: err = %v", name, err)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}
