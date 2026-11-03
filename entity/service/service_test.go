package service_test

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/service"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "66666666-6666-6666-6666-666666666666"

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func meta(href, typ string) model.Meta {
	return model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/" + href, Type: model.MetaType(typ), MediaType: "application/json"}
}

var (
	full = service.CreateRequest{
		Name: "Доставка",
		Code: model.Value("DLV"), ExternalCode: model.Value("dlv-1"), Description: model.Value(""),
		Archived: model.Value(false), Shared: model.Value(true), VAT: model.Value(0), VATEnabled: model.Value(true),
		UseParentVAT: model.Value(false), DiscountProhibited: model.Value(false),
		PaymentItemType: model.Value("SERVICE"), TaxSystem: model.Value("GENERAL_TAX_SYSTEM"),
		UOM:        model.Value(model.Reference{Meta: meta("entity/uom/19f1edc0-fc42-4001-94cb-c9ec9c62ec10", "uom")}),
		BuyPrice:   model.Value(model.Price{Value: 10000, Currency: model.Reference{Meta: meta("entity/currency/22222222-2222-2222-2222-222222222222", "currency")}}),
		SalePrices: model.Value([]model.SalePrice{{Value: 15050.5, PriceType: model.PriceType{Meta: meta("context/companysettings/pricetype/88888888-8888-8888-8888-888888888888", "pricetype")}}}),
		Barcodes:   model.Value([]model.Barcode{{"ean13": "4600000000001"}, {"code128": "CODE128-001"}}),
		Attributes: model.Value([]model.Attribute{}),
		TasnifUZ:   model.Value(model.TasnifUZ{IKPU: "10199001001000000", PackCode: "1505953"}),
		SyncID:     model.Value(testID),
	}
	update = service.UpdateRequest{Name: model.Value("Монтаж"),
		MinPrice: model.Null[model.Price](), Barcodes: model.Value([]model.Barcode{}),
		ProductFolder: model.Null[model.Reference](), VAT: model.Value(20)}
)

func TestServiceFixture(t *testing.T) {
	data := read(t, "upstream/service.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got service.Service
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.SyncID != testID || got.PaymentItemType != "SERVICE" || got.TaxSystem != "GENERAL_TAX_SYSTEM" ||
				len(got.Barcodes) != 2 || got.Barcodes[1]["code128"] != "CODE128-001" || got.Files == nil || got.Files.Meta.Href == "" ||
				len(got.SalePrices) != 1 || len(got.Attributes) != 3 || got.PathName == "" {
				t.Fatalf("service = %+v", got)
			}
			var sparse service.Service
			if err := codec.Unmarshal([]byte(`{"id":"x","name":"y","paymentItemType":"FUTURE","futureField":[1]}`), &sparse); err != nil ||
				sparse.PaymentItemType != "FUTURE" || sparse.BuyPrice != nil || sparse.Files != nil || sparse.TasnifUZ != nil {
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
			{"create full", full, read(t, "custom/service/create-full.golden.json")},
			{"create minimal", service.CreateRequest{Name: "Услуга"}, []byte(`{"name":"Услуга"}`)},
			{"update", update, read(t, "custom/service/update.golden.json")},
			{"update false", service.UpdateRequest{Archived: model.Value(false)}, []byte(`{"archived":false}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[service.CreateRequest](), reflect.TypeFor[service.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated", "PathName", "EffectiveVAT", "EffectiveVATEnabled", "Files"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
	if _, ok := reflect.TypeFor[service.UpdateRequest]().FieldByName("SyncID"); ok {
		t.Error("UpdateRequest contains create-only SyncID")
	}
}

func client(c *transport.Client) *service.Client {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := service.NewService(c)
	return s
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/service/` + testID + `"},"id":"` + testID + `","name":"Доставка"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*service.Service, error)
	}{
		{"get", "GET", "entity/service/" + testID, nil, func(ctx context.Context, c *transport.Client) (*service.Service, error) {
			return client(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/service", read(t, "custom/service/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*service.Service, error) {
			return client(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/service/" + testID, read(t, "custom/service/update.golden.json"), func(ctx context.Context, c *transport.Client) (*service.Service, error) {
			return client(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *service.Service) {
						if got.ID != testID || got.Name != "Доставка" {
							t.Fatalf("service = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func TestServiceValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := client(c)
	if _, err := s.Create(t.Context(), service.CreateRequest{Code: model.Value("x")}); err == nil {
		t.Error("create without name accepted")
	}
	if _, err := s.Update(t.Context(), testID, service.UpdateRequest{}); err == nil {
		t.Error("empty update accepted")
	}
	if _, err := s.Get(t.Context(), "bad"); err == nil {
		t.Error("bad id accepted")
	}
	if _, err := service.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestServiceList(t *testing.T) {
	listtest.Run(t, "service", read(t, "upstream/service.json"), "", "archived", func(c *transport.Client) listtest.ListFunc[service.Service] {
		return client(c).List
	}, func(t *testing.T, got service.Service) {
		if got.ID != testID || len(got.Barcodes) != 2 {
			t.Fatalf("row = %+v", got)
		}
	})
}
