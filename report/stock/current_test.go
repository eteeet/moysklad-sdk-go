package stock_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/report/stock"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	a1 = "12345678-5838-aaeb-0a80-003a003ef439"
	a2 = "12345678-279c-aaeb-0a80-00d6001f847c"
	s1 = "12345678-b123-aaee-0a80-012b0001bb10"
	s2 = "12345678-b123-aaee-0a80-012b0001bb13"
)

// value возвращает значение Field и признаки «задано» и «null».
func value(f model.Field[float64]) (float64, bool, bool) {
	v, ok := f.Get()
	return v, ok, f.IsNull()
}

func TestCurrentRowJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var one stock.CurrentRow
			if err := codec.Unmarshal(read(t, "upstream/stock_current_all.json"), &one); err != nil {
				t.Fatal(err)
			}
			if v, ok, _ := value(one.Stock); one.AssortmentID != a1 || v != 11 || !ok || !one.Quantity.IsZero() || !one.FreeStock.IsZero() {
				t.Fatalf("row = %+v", one)
			}
			var rows []stock.CurrentRow
			if err := codec.Unmarshal(read(t, "custom/stock/all_current.json"), &rows); err != nil {
				t.Fatal(err)
			}
			if v, ok, _ := value(rows[1].Stock); len(rows) != 2 || rows[1].AssortmentID != a2 || v != 46.5 || !ok {
				t.Fatalf("rows = %+v", rows)
			}
			// Поле значения выбирает stockType: читаются все пять, незаданные остаются IsZero.
			for key, get := range map[string]func(stock.CurrentRow) model.Field[float64]{
				"stock":     func(r stock.CurrentRow) model.Field[float64] { return r.Stock },
				"freeStock": func(r stock.CurrentRow) model.Field[float64] { return r.FreeStock },
				"quantity":  func(r stock.CurrentRow) model.Field[float64] { return r.Quantity },
				"reserve":   func(r stock.CurrentRow) model.Field[float64] { return r.Reserve },
				"inTransit": func(r stock.CurrentRow) model.Field[float64] { return r.InTransit },
			} {
				var r stock.CurrentRow
				if err := codec.Unmarshal([]byte(`{"assortmentId":"`+a1+`","`+key+`":-0.5}`), &r); err != nil {
					t.Fatal(err)
				}
				set := 0
				for _, f := range []model.Field[float64]{r.Stock, r.FreeStock, r.Quantity, r.Reserve, r.InTransit} {
					if !f.IsZero() {
						set++
					}
				}
				if v, ok, _ := value(get(r)); v != -0.5 || !ok || set != 1 {
					t.Fatalf("%s: row = %+v", key, r)
				}
			}
		})
	}
}

func TestStoreCurrentRowJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var one stock.StoreCurrentRow
			if err := codec.Unmarshal(read(t, "upstream/stock_current_by_store.json"), &one); err != nil {
				t.Fatal(err)
			}
			if v, ok, _ := value(one.FreeStock); one.AssortmentID != a2 || one.StoreID != "" || v != -5 || !ok || !one.Stock.IsZero() {
				t.Fatalf("row = %+v", one)
			}
			var rows []stock.StoreCurrentRow
			if err := codec.Unmarshal(read(t, "custom/stock/bystore_current.json"), &rows); err != nil {
				t.Fatal(err)
			}
			if v, _, _ := value(rows[3].FreeStock); len(rows) != 4 || rows[0].StoreID != s1 || rows[2].StoreID != "" || v != 46 {
				t.Fatalf("rows = %+v", rows)
			}
			// withRecalculate: пересчитываемая позиция приходит с "stock":null — отличается от 0 и от отсутствия.
			if err := codec.Unmarshal(read(t, "custom/stock/bystore_current_recalculate.json"), &rows); err != nil {
				t.Fatal(err)
			}
			if _, ok, null := value(rows[1].Stock); len(rows) != 3 || ok || !null {
				t.Fatalf("recalculating = %+v", rows[1])
			}
			if v, ok, null := value(rows[0].Stock); v != 1 || !ok || null {
				t.Fatalf("recalculated = %+v", rows[0])
			}
		})
	}
}

func TestSlotRowJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var rows []stock.SlotRow
			if err := codec.Unmarshal(read(t, "custom/stock/byslot_current.json"), &rows); err != nil {
				t.Fatal(err)
			}
			if len(rows) != 3 || rows[0] != (stock.SlotRow{AssortmentID: a1, StoreID: s1, SlotID: "c3b59812-cd5a-11ed-0a80-0142000026a3", Stock: 1.2}) {
				t.Fatalf("rows = %+v", rows)
			}
			// Фикстура официального SDK содержит quantity вместо stock из документации: SDK читает только
			// stock, поэтому остаток здесь 0 — противоречие записано в docs/report-stock.md.
			var upstream stock.SlotRow
			if err := codec.Unmarshal(read(t, "upstream/stock_current_by_slot.json"), &upstream); err != nil ||
				upstream.AssortmentID != a1 || upstream.SlotID != "c3b59812-cd5a-11ed-0a80-0142000026a3" || upstream.Stock != 0 {
				t.Fatalf("upstream = %+v, err = %v", upstream, err)
			}
		})
	}
}

func TestAllCurrentHTTP(t *testing.T) {
	for _, tc := range []struct {
		name    string
		options stock.CurrentOptions
		query   string
	}{
		{"full", stock.CurrentOptions{StockType: stock.StockTypeQuantity, Include: stock.IncludeZeroLines,
			ChangedSince: documenttest.DateTime("2026-11-03 15:00:00.000"), AssortmentIDs: []string{a1, a2}, StoreIDs: []string{s1}},
			"changedSince=2026-11-03+15%3A00%3A00&filter=assortmentId%3D" + a1 + "%2C" + a2 + "%3BstoreId%3D" + s1 + "&include=zeroLines&stockType=quantity"},
		{"defaults", stock.CurrentOptions{}, ""},
		{"unknown stockType and extra", stock.CurrentOptions{StockType: "future", Extra: url.Values{"x": {"1"}}}, "stockType=future&x=1"},
	} {
		invoke := func(ctx context.Context, c *transport.Client) (*[]stock.CurrentRow, error) {
			rows, err := service(c).AllCurrent(ctx, tc.options)
			if err != nil {
				return nil, err
			}
			return &rows, nil
		}
		path := "report/stock/all/current"
		if tc.query != "" {
			path += "?" + tc.query
		}
		for _, rc := range documenttest.ResponseCases(string(read(t, "custom/stock/all_current.json"))) {
			t.Run(tc.name+"/"+rc.Name, func(t *testing.T) {
				documenttest.HTTP(t, rc, "GET", path, nil, invoke, func(got *[]stock.CurrentRow) {
					if v, _ := (*got)[0].Stock.Get(); len(*got) != 2 || (*got)[0].AssortmentID != a1 || v != 11 {
						t.Fatalf("rows = %+v", *got)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
	// Пустой массив — успешный ответ без строк.
	rc := documenttest.ResponseCases(`[]`)[0]
	documenttest.HTTP(t, rc, "GET", "report/stock/all/current", nil, func(ctx context.Context, c *transport.Client) (*[]stock.CurrentRow, error) {
		rows, err := service(c).AllCurrent(ctx, stock.CurrentOptions{})
		return &rows, err
	}, func(got *[]stock.CurrentRow) {
		if *got == nil || len(*got) != 0 {
			t.Fatalf("rows = %#v", *got)
		}
	})
}

func TestByStoreCurrentHTTP(t *testing.T) {
	for _, tc := range []struct {
		name    string
		options stock.CurrentOptions
		query   string
	}{
		{"filter", stock.CurrentOptions{AssortmentIDs: []string{a1, a2}, StoreIDs: []string{s1, s2}},
			"filter=assortmentId%3D" + a1 + "%2C" + a2 + "%3BstoreId%3D" + s1 + "%2C" + s2},
		{"recalculate", stock.CurrentOptions{StockType: stock.StockTypeFreeStock, WithRecalculate: true, StoreIDs: []string{s1}},
			"filter=storeId%3D" + s1 + "&stockType=freeStock&withRecalculate=true"},
	} {
		invoke := func(ctx context.Context, c *transport.Client) (*[]stock.StoreCurrentRow, error) {
			rows, err := service(c).ByStoreCurrent(ctx, tc.options)
			if err != nil {
				return nil, err
			}
			return &rows, nil
		}
		for _, rc := range documenttest.ResponseCases(string(read(t, "custom/stock/bystore_current.json"))) {
			t.Run(tc.name+"/"+rc.Name, func(t *testing.T) {
				documenttest.HTTP(t, rc, "GET", "report/stock/bystore/current?"+tc.query, nil, invoke, func(got *[]stock.StoreCurrentRow) {
					if v, _ := (*got)[2].FreeStock.Get(); len(*got) != 4 || (*got)[2].StoreID != "" || v != -5 {
						t.Fatalf("rows = %+v", *got)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

func TestBySlotCurrentHTTP(t *testing.T) {
	for _, tc := range []struct {
		name    string
		options stock.SlotOptions
		query   string
	}{
		{"assortment and store", stock.SlotOptions{AssortmentIDs: []string{a1, a2}, StoreIDs: []string{s1, s2}},
			"filter=assortmentId%3D" + a1 + "%2C" + a2 + "%3BstoreId%3D" + s1 + "%2C" + s2},
		{"store only", stock.SlotOptions{StoreIDs: []string{s1}}, "filter=storeId%3D" + s1},
		{"assortment and extra", stock.SlotOptions{AssortmentIDs: []string{a1}, Extra: url.Values{"stockType": {"quantity"}}},
			"filter=assortmentId%3D" + a1 + "&stockType=quantity"},
	} {
		invoke := func(ctx context.Context, c *transport.Client) (*[]stock.SlotRow, error) {
			rows, err := service(c).BySlotCurrent(ctx, tc.options)
			if err != nil {
				return nil, err
			}
			return &rows, nil
		}
		for _, rc := range documenttest.ResponseCases(string(read(t, "custom/stock/byslot_current.json"))) {
			t.Run(tc.name+"/"+rc.Name, func(t *testing.T) {
				documenttest.HTTP(t, rc, "GET", "report/stock/byslot/current?"+tc.query, nil, invoke, func(got *[]stock.SlotRow) {
					if len(*got) != 3 || (*got)[2].SlotID != "dc27cf31-212f-11ee-ac12-001000000000" || (*got)[0].Stock != 1.2 {
						t.Fatalf("rows = %+v", *got)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

func TestCurrentValidation(t *testing.T) {
	c, calls := noCalls(t)
	invalid := map[string]stock.CurrentOptions{
		"assortment not UUID": {AssortmentIDs: []string{"p1"}},
		"store not UUID":      {StoreIDs: []string{s1, s1 + ","}},
		"empty store":         {StoreIDs: []string{""}},
		"extra filter":        {Extra: url.Values{"filter": {"storeId=" + s1}}},
		"extra stockType":     {StockType: "stock", Extra: url.Values{"stockType": {"reserve"}}},
		"extra include":       {Extra: url.Values{"include": {"zeroLines"}}},
		"extra changedSince":  {Extra: url.Values{"changedSince": {"2026-11-03 15:00:00"}}},
		"extra recalculate":   {Extra: url.Values{"withRecalculate": {"true"}}},
		"extra async":         {Extra: url.Values{"async": {"true"}}},
	}
	for name, options := range invalid {
		localError(t, name, "get stock all current", func(o transport.RequestOption) error {
			_, err := service(c).AllCurrent(t.Context(), options, o)
			return err
		})
		localError(t, name, "get stock bystore current", func(o transport.RequestOption) error {
			_, err := service(c).ByStoreCurrent(t.Context(), options, o)
			return err
		})
	}
	localError(t, "withRecalculate", "get stock all current", func(o transport.RequestOption) error {
		_, err := service(c).AllCurrent(t.Context(), stock.CurrentOptions{WithRecalculate: true}, o)
		return err
	})
	for name, options := range map[string]stock.SlotOptions{
		"no filter":           {},
		"only extra":          {Extra: url.Values{"stockType": {"stock"}}},
		"assortment not UUID": {AssortmentIDs: []string{"p1"}},
		"extra filter":        {StoreIDs: []string{s1}, Extra: url.Values{"filter": {"storeId=" + s1}}},
	} {
		localError(t, name, "get stock byslot current", func(o transport.RequestOption) error {
			_, err := service(c).BySlotCurrent(t.Context(), options, o)
			return err
		})
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}
