package byoperations_test

import (
	"context"
	"fmt"
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
	"github.com/eteeet/moysklad-sdk-go/report/byoperations"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const product = "https://api.moysklad.ru/api/remap/1.2/entity/product/bb1865dc-32e1-11ef-ac16-001100000003"

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func service(c *transport.Client) *byoperations.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := byoperations.NewService(c)
	return s
}

func ref(path, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/" + path, Type: model.MetaType(typ)}}
}

func page(row []byte) string {
	return `{"context":{"employee":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/context/employee"}}},"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/report/byoperations/stock","type":"byoperations","size":1,"limit":1000,"offset":0},"rows":[` + string(row) + "]}"
}

func moment(t *testing.T, got model.DateTime, want string) {
	t.Helper()
	if !time.Time(got).Equal(time.Time(documenttest.DateTime(want))) {
		t.Errorf("moment = %v, want %s", time.Time(got), want)
	}
}

func TestStockRowJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got byoperations.StockRow
			if err := codec.Unmarshal(read(t, "upstream/by_operations_stock.json"), &got); err != nil {
				t.Fatal(err)
			}
			if got.Assortment == nil || got.Assortment.Meta.Href != product || got.Operation == nil || got.Operation.Meta.Type != "supply" ||
				got.Store == nil || got.Store.Meta.Type != "store" || got.Stock != 1 || got.AvgStockDays != 9.72 ||
				got.CostPerUnit == nil || *got.CostPerUnit != 0 || got.SumCost == nil || *got.SumCost != 0 {
				t.Fatalf("row = %+v", got)
			}
			moment(t, got.Moment, "2024-07-02 15:52:00")
			var p model.Collection[byoperations.StockRow]
			if err := codec.Unmarshal(read(t, "custom/byoperations/stock.json"), &p); err != nil {
				t.Fatal(err)
			}
			if p.Context == nil || p.Meta.Type != "byoperations" || p.Meta.Size != 1 || len(p.Rows) != 1 || p.Rows[0].Stock != 1 || p.Rows[0].Operation.Meta.Type != "supply" {
				t.Fatalf("page = %+v", p)
			}
			// Без права видеть себестоимость сервер не присылает costPerUnit и sumCost.
			var hidden byoperations.StockRow
			if err := codec.Unmarshal([]byte(`{"operation":{"meta":{"href":"h","type":"futuredocument"}},"stock":-2.5,"avgStockDays":0,"future":1}`), &hidden); err != nil ||
				hidden.CostPerUnit != nil || hidden.SumCost != nil || hidden.Stock != -2.5 || hidden.Operation.Meta.Type != "futuredocument" || hidden.Store != nil {
				t.Fatalf("hidden = %+v, err = %v", hidden, err)
			}
			for _, money := range []struct {
				json string
				want float64
			}{{"0", 0}, {"0.0", 0}, {"150.5", 150.5}, {"123456789012.25", 123456789012.25}} {
				var row byoperations.StockRow
				body := fmt.Sprintf(`{"costPerUnit":%s,"sumCost":%s}`, money.json, money.json)
				if err := codec.Unmarshal([]byte(body), &row); err != nil || row.CostPerUnit == nil || *row.CostPerUnit != money.want || row.SumCost == nil || *row.SumCost != money.want {
					t.Errorf("%s: row = %+v, err = %v", money.json, row, err)
				}
			}
		})
	}
}

func TestReserveAndInTransitRowJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var reserve byoperations.ReserveRow
			if err := codec.Unmarshal(read(t, "upstream/by_operations_reserve.json"), &reserve); err != nil {
				t.Fatal(err)
			}
			if reserve.Assortment == nil || reserve.Assortment.Meta.Href != product || reserve.Operation == nil || reserve.Operation.Meta.Type != "customerorder" ||
				reserve.Store == nil || reserve.Reserve != 5 {
				t.Fatalf("reserve = %+v", reserve)
			}
			moment(t, reserve.Moment, "2024-07-02 16:45:00")
			var transit byoperations.InTransitRow
			if err := codec.Unmarshal(read(t, "upstream/by_operations_intransit.json"), &transit); err != nil {
				t.Fatal(err)
			}
			if transit.Operation == nil || transit.Operation.Meta.Type != "purchaseorder" || transit.Store == nil || transit.InTransit != 5 {
				t.Fatalf("intransit = %+v", transit)
			}
			moment(t, transit.Moment, "2024-07-12 10:03:00")
			var reserves model.Collection[byoperations.ReserveRow]
			if err := codec.Unmarshal(read(t, "custom/byoperations/reserve.json"), &reserves); err != nil || len(reserves.Rows) != 1 || reserves.Rows[0].Reserve != 5 || reserves.Meta.Type != "byoperations" {
				t.Fatalf("reserves = %+v, err = %v", reserves, err)
			}
			var transits model.Collection[byoperations.InTransitRow]
			if err := codec.Unmarshal(read(t, "custom/byoperations/intransit.json"), &transits); err != nil || len(transits.Rows) != 1 || transits.Rows[0].InTransit != 5 || transits.Context == nil {
				t.Fatalf("transits = %+v, err = %v", transits, err)
			}
			// store у резервов и ожиданий не обязателен.
			var bare byoperations.ReserveRow
			if err := codec.Unmarshal([]byte(`{"operation":{"meta":{"href":"h","type":"customerorder"}},"reserve":0.25}`), &bare); err != nil || bare.Store != nil || bare.Reserve != 0.25 {
				t.Fatalf("bare = %+v, err = %v", bare, err)
			}
		})
	}
}

func TestHTTP(t *testing.T) {
	const full = "filter=assortment%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fconsignment%2Fc1&limit=2&offset=4&x=1"
	const defaults = "filter=assortment%3Dhttps%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Fproduct%2Fp1&limit=1000&offset=0"
	options := []struct {
		name    string
		options byoperations.Options
		query   string
	}{
		{"full", byoperations.Options{Assortment: ref("entity/consignment/c1", "consignment"), Limit: 2, Offset: 4, Extra: url.Values{"x": {"1"}}}, full},
		{"defaults", byoperations.Options{Assortment: ref("entity/product/p1", "product")}, defaults},
		{"variant", byoperations.Options{Assortment: ref("entity/variant/v1", "variant")}, strings.ReplaceAll(defaults, "product%2Fp1", "variant%2Fv1")},
		// Тип ссылки в запрос не уходит и до HTTP не требуется.
		{"href only", byoperations.Options{Assortment: model.Reference{Meta: model.Meta{Href: ref("entity/product/p1", "product").Meta.Href}}}, defaults},
	}
	for _, tc := range options {
		t.Run("stock/"+tc.name, func(t *testing.T) {
			invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[byoperations.StockRow], error) {
				return service(c).Stock(ctx, tc.options)
			}
			for _, rc := range documenttest.ResponseCases(page(read(t, "upstream/by_operations_stock.json"))) {
				t.Run(rc.Name, func(t *testing.T) {
					documenttest.HTTP(t, rc, "GET", "report/byoperations/stock?"+tc.query, nil, invoke, func(got *model.Collection[byoperations.StockRow]) {
						if len(got.Rows) != 1 || got.Rows[0].Stock != 1 || got.Rows[0].AvgStockDays != 9.72 {
							t.Fatalf("page = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, invoke)
		})
		t.Run("reserve/"+tc.name, func(t *testing.T) {
			invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[byoperations.ReserveRow], error) {
				return service(c).Reserve(ctx, tc.options)
			}
			for _, rc := range documenttest.ResponseCases(page(read(t, "upstream/by_operations_reserve.json"))) {
				t.Run(rc.Name, func(t *testing.T) {
					documenttest.HTTP(t, rc, "GET", "report/byoperations/reserve?"+tc.query, nil, invoke, func(got *model.Collection[byoperations.ReserveRow]) {
						if len(got.Rows) != 1 || got.Rows[0].Reserve != 5 {
							t.Fatalf("page = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, invoke)
		})
		t.Run("intransit/"+tc.name, func(t *testing.T) {
			invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[byoperations.InTransitRow], error) {
				return service(c).InTransit(ctx, tc.options)
			}
			for _, rc := range documenttest.ResponseCases(page(read(t, "upstream/by_operations_intransit.json"))) {
				t.Run(rc.Name, func(t *testing.T) {
					documenttest.HTTP(t, rc, "GET", "report/byoperations/intransit?"+tc.query, nil, invoke, func(got *model.Collection[byoperations.InTransitRow]) {
						if len(got.Rows) != 1 || got.Rows[0].InTransit != 5 {
							t.Fatalf("page = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, invoke)
		})
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
	p := ref("entity/product/p1", "product")
	methods := map[string]func(byoperations.Options, ...transport.RequestOption) error{
		"stock": func(o byoperations.Options, r ...transport.RequestOption) error {
			_, err := service(c).Stock(t.Context(), o, r...)
			return err
		},
		"reserve": func(o byoperations.Options, r ...transport.RequestOption) error {
			_, err := service(c).Reserve(t.Context(), o, r...)
			return err
		},
		"intransit": func(o byoperations.Options, r ...transport.RequestOption) error {
			_, err := service(c).InTransit(t.Context(), o, r...)
			return err
		},
	}
	for method, invoke := range methods {
		for name, options := range map[string]byoperations.Options{
			"no assortment": {},
			"no href":       {Assortment: model.Reference{Meta: model.Meta{Type: "product"}}},
			"extra limit":   {Assortment: p, Extra: url.Values{"limit": {"5"}}},
			"extra offset":  {Assortment: p, Extra: url.Values{"offset": {"5"}}},
			"extra filter":  {Assortment: p, Extra: url.Values{"filter": {"store=x"}}},
			"extra expand":  {Assortment: p, Extra: url.Values{"expand": {"operation"}}},
			"limit -1":      {Assortment: p, Limit: -1},
			"limit 1001":    {Assortment: p, Limit: 1001},
			"offset -1":     {Assortment: p, Offset: -1},
		} {
			info := transport.ResponseInfo{StatusCode: 503}
			err := invoke(options, transport.WithResponseInfo(&info))
			if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "get byoperations "+method) || info.StatusCode != 0 {
				t.Errorf("%s %s: err = %v, info = %+v", method, name, err, info)
			}
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := byoperations.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}
