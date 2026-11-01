package bundle_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/bundle"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID      = "66666666-6666-6666-6666-666666666666"
	fixtureID   = "1545463f-26af-11f1-0a83-1f2500000181"
	componentID = "1545552d-26af-11f1-0a83-1f2500000185"
)

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ref(href, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/" + href, Type: model.MetaType(typ), MediaType: "application/json"}}
}

func price(value float64) model.Price {
	return model.Price{Value: value, Currency: ref("entity/currency/22222222-2222-2222-2222-222222222222", "currency")}
}

var (
	full = bundle.CreateRequest{
		Name: "Набор",
		Components: []bundle.ComponentRequest{
			{Assortment: ref("entity/product/11111111-1111-1111-1111-111111111111", "product"), Quantity: 2},
			{Assortment: ref("entity/variant/33333333-3333-3333-3333-333333333333", "variant"), Quantity: 0.5},
		},
		Code: model.Value("SET-1"), ExternalCode: model.Value("set-1"), Description: model.Value(""), Article: model.Value("A-1"),
		Archived: model.Value(false), VAT: model.Value(0), VATEnabled: model.Value(true), UseParentVAT: model.Value(false),
		DiscountProhibited: model.Value(false), PaymentItemType: model.Value("GOOD"), TaxSystem: model.Value("GENERAL_TAX_SYSTEM"),
		TrackingType: model.Value("NOT_TRACKED"), PartialDisposal: model.Value(false), TNVED: model.Value("1234567890"),
		Weight: model.Value(0.0), Volume: model.Value(1.5),
		UOM:      model.Value(ref("entity/uom/19f1edc0-fc42-4001-94cb-c9ec9c62ec10", "uom")),
		Country:  model.Value(ref("entity/country/44444444-4444-4444-4444-444444444444", "country")),
		MinPrice: model.Value(price(50000)), Overhead: model.Value(price(1050.5)),
		SalePrices: model.Value([]model.SalePrice{{Value: 99900, PriceType: model.PriceType{Meta: ref("context/companysettings/pricetype/88888888-8888-8888-8888-888888888888", "pricetype").Meta}}}),
		Barcodes:   model.Value([]model.Barcode{{"ean13": "2000000000046"}}),
		MarkingUZ:  model.Value(model.MarkingUZ{TrackingType: "WATER"}),
		SyncID:     model.Value(testID),
	}
	update = bundle.UpdateRequest{
		Name:       model.Value("Набор 2"),
		Components: model.Value([]bundle.ComponentRequest{{Assortment: ref("entity/service/55555555-5555-5555-5555-555555555555", "service"), Quantity: 1}, {ID: componentID, Quantity: 3}}),
		Barcodes:   model.Value([]model.Barcode{}), Overhead: model.Null[model.Price](), Country: model.Null[model.Reference](),
		Weight: model.Value(0.0),
	}
)

func TestBundleFixture(t *testing.T) {
	data := read(t, "upstream/bundle.json")
	component := read(t, "upstream/bundle_component.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got bundle.Bundle
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Name != "b1" || got.TrackingType != "TOBACCO" || got.Article != "fsdfds" || got.Volume != 10 ||
				got.MinPrice == nil || got.MinPrice.Value != 2000 || len(got.SalePrices) != 1 || len(got.Barcodes) != 1 ||
				got.Images == nil || got.Images.Meta.Size != 1 || got.Files == nil || got.Components == nil || got.Components.Meta.Size != 1 ||
				len(got.Components.Rows) != 1 || got.Components.Rows[0].ID != componentID || got.Components.Rows[0].Quantity != 1 ||
				got.Components.Rows[0].Assortment.Meta.Type != "product" || got.PathName != "g1" || !got.DiscountProhibited {
				t.Fatalf("bundle = %+v", got)
			}
			var c bundle.Component
			if err := codec.Unmarshal(component, &c); err != nil || c.ID != componentID || c.Meta.Type != "bundlecomponent" || c.Quantity != 1 {
				t.Fatalf("component = %+v, err = %v", c, err)
			}
			var sparse bundle.Bundle
			if err := codec.Unmarshal([]byte(`{"id":"x","name":"y","trackingType":"FUTURE","components":{"meta":{"href":"h","size":2}},"mod__marking__uz":{"trackingType":"BEER"},"future":[1]}`), &sparse); err != nil ||
				sparse.TrackingType != "FUTURE" || sparse.Components == nil || sparse.Components.Rows != nil || sparse.MarkingUZ == nil ||
				sparse.MarkingUZ.TrackingType != "BEER" || sparse.Overhead != nil || sparse.Images != nil {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", full, read(t, "custom/bundle/create-full.golden.json")},
			{"create minimal", bundle.CreateRequest{Name: "Набор", Components: full.Components[:1]},
				[]byte(`{"name":"Набор","components":[{"assortment":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/11111111-1111-1111-1111-111111111111","type":"product","mediaType":"application/json"}},"quantity":2}]}`)},
			{"update", update, read(t, "custom/bundle/update.golden.json")},
			{"update false", bundle.UpdateRequest{PartialDisposal: model.Value(false)}, []byte(`{"partialDisposal":false}`)},
		} {
			t.Run(codec.Name+"/"+tc.name, func(t *testing.T) {
				got, err := codec.Marshal(tc.input)
				if err != nil {
					t.Fatal(err)
				}
				documenttest.JSONEqual(t, got, tc.want)
			})
		}
	}
	for _, typ := range []reflect.Type{reflect.TypeFor[bundle.CreateRequest](), reflect.TypeFor[bundle.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated", "PathName", "EffectiveVAT", "EffectiveVATEnabled", "Files", "Images"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only or unsupported %s", typ.Name(), name)
			}
		}
	}
	if _, ok := reflect.TypeFor[bundle.UpdateRequest]().FieldByName("SyncID"); ok {
		t.Error("UpdateRequest contains create-only SyncID")
	}
}

func service(c *transport.Client) *bundle.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := bundle.NewService(c)
	return s
}

func TestBundleHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/bundle/` + testID + `"},"id":"` + testID + `","name":"Набор"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*bundle.Bundle, error)
	}{
		{"get", "GET", "entity/bundle/" + testID, nil, func(ctx context.Context, c *transport.Client) (*bundle.Bundle, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/bundle", read(t, "custom/bundle/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*bundle.Bundle, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/bundle/" + testID, read(t, "custom/bundle/update.golden.json"), func(ctx context.Context, c *transport.Client) (*bundle.Bundle, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *bundle.Bundle) {
						if got.ID != testID || got.Name != "Набор" {
							t.Fatalf("bundle = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func TestGetComponents(t *testing.T) {
	row := string(read(t, "upstream/bundle_component.json"))
	for _, page := range []query.PageOptions{{}, {Limit: 1, Offset: 7}} {
		limit := page.Limit
		if limit == 0 {
			limit = 1000
		}
		query := fmt.Sprintf("?limit=%d&offset=%d", limit, page.Offset)
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[bundle.Component], error) {
			return service(c).GetComponents(ctx, testID, page)
		}
		for _, tc := range documenttest.ResponseCases(`{"meta":{"href":"https://example.test/components","type":"bundlecomponent","size":1},"rows":[` + row + `]}`) {
			t.Run(tc.Name+query, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", "entity/bundle/"+testID+"/components"+query, nil, invoke, func(got *model.Collection[bundle.Component]) {
					if len(got.Rows) != 1 || got.Rows[0].ID != componentID || got.Rows[0].Assortment.Meta.Href == "" {
						t.Fatalf("rows = %+v", got.Rows)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

func TestGetComponent(t *testing.T) {
	invoke := func(ctx context.Context, c *transport.Client) (*bundle.Component, error) {
		return service(c).GetComponent(ctx, testID, componentID)
	}
	for _, tc := range documenttest.ResponseCases(string(read(t, "upstream/bundle_component.json"))) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "GET", "entity/bundle/"+testID+"/components/"+componentID, nil, invoke, func(got *bundle.Component) {
				if got.ID != componentID || got.Quantity != 1 {
					t.Fatalf("component = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}

func TestBundleValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(c)
	one := full.Components[:1]
	for name, input := range map[string]bundle.CreateRequest{
		"no name":       {Components: one},
		"no components": {Name: "x"},
	} {
		if _, err := s.Create(t.Context(), input); err == nil {
			t.Errorf("create %s accepted", name)
		}
	}
	if _, err := s.Update(t.Context(), testID, bundle.UpdateRequest{}); err == nil {
		t.Error("empty update accepted")
	}
	if _, err := s.Get(t.Context(), "bad"); err == nil {
		t.Error("bad id accepted")
	}
	if _, err := s.GetComponent(t.Context(), testID, "bad"); err == nil {
		t.Error("bad component id accepted")
	}
	for _, page := range []query.PageOptions{{Limit: -1}, {Limit: 1001}, {Offset: -1}, {Fields: "stock"}} {
		if _, err := s.GetComponents(t.Context(), testID, page); err == nil {
			t.Errorf("page %+v accepted", page)
		}
	}
	if _, err := s.GetComponents(t.Context(), "bad", query.PageOptions{}); err == nil {
		t.Error("bad bundle id accepted")
	}
	if _, err := bundle.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestBundleList(t *testing.T) {
	listtest.Run(t, "bundle", read(t, "upstream/bundle.json"), "", "archived", func(c *transport.Client) listtest.ListFunc[bundle.Bundle] {
		return service(c).List
	}, func(t *testing.T, got bundle.Bundle) {
		if got.ID != fixtureID || len(got.Barcodes) != 1 {
			t.Fatalf("row = %+v", got)
		}
	})
}
