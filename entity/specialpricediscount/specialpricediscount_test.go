package specialpricediscount_test

import (
	"context"
	"github.com/eteeet/moysklad-sdk-go/entity/discount"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/specialpricediscount"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "96673f4d-9f4d-11e6-8af5-581e0000007b"

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ref(path, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/" + path, Type: model.MetaType(typ), MediaType: "application/json"}}
}

var (
	full   = specialpricediscount.CreateRequest{Name: "discountName", Active: true, AllProducts: false, AllAgents: false, AgentTags: model.Value([]string{"tag1", "tag2"}), Assortment: model.Value([]model.Reference{ref("entity/product/833eac42-b6f4-11ea-ac12-000e00000006", "product")}), ProductFolders: model.Value([]model.Reference{ref("entity/productfolder/39c62b64-a722-11ea-ac12-000d00000015", "productfolder")}), UsePriceType: true, SpecialPrice: model.Value(discount.SpecialPriceRequest{PriceType: model.Value(ref("context/companysettings/pricetype/5f441404-a419-11ea-ac12-000a00000078", "pricetype")), Value: model.Value(0.0)})}
	update = specialpricediscount.UpdateRequest{Name: model.Value("updatedName"), UsePriceType: model.Value(false), Discount: model.Value(50.0)}
)

func TestSpecialPriceDiscountFixture(t *testing.T) {
	data := read(t, "upstream/special_price_discount.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got specialpricediscount.SpecialPriceDiscount
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.Name != "Специальная процентная скидка" || got.Active || got.AllAgents || !got.UsePriceType || got.SpecialPrice == nil || got.SpecialPrice.Value != 15 || got.SpecialPrice.PriceType == nil || len(got.AgentTags) != 2 {
				t.Fatalf("specialpricediscount = %+v", got)
			}
			// Неизвестные поля не ломают разбор, отсутствующие не выдумываются.
			var sparse specialpricediscount.SpecialPriceDiscount
			if err := codec.Unmarshal([]byte(`{"id":"x","futureField":{"a":1}}`), &sparse); err != nil || sparse.ID != "x" {
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
			{"create full", full, read(t, "custom/specialpricediscount/create-full.golden.json")},
			{"update", update, read(t, "custom/specialpricediscount/update.golden.json")},
			{"create required false", specialpricediscount.CreateRequest{Name: "n"}, []byte(`{"name":"n","active":false,"allProducts":false,"allAgents":false,"usePriceType":false}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[specialpricediscount.CreateRequest](), reflect.TypeFor[specialpricediscount.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func service(c *transport.Client) *specialpricediscount.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := specialpricediscount.NewService(c)
	return s
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/specialpricediscount/` + testID + `"},"id":"` + testID + `","name":"n"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*specialpricediscount.SpecialPriceDiscount, error)
	}{
		{"get", "GET", "entity/specialpricediscount/" + testID, nil, func(ctx context.Context, c *transport.Client) (*specialpricediscount.SpecialPriceDiscount, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/specialpricediscount", read(t, "custom/specialpricediscount/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*specialpricediscount.SpecialPriceDiscount, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/specialpricediscount/" + testID, read(t, "custom/specialpricediscount/update.golden.json"), func(ctx context.Context, c *transport.Client) (*specialpricediscount.SpecialPriceDiscount, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *specialpricediscount.SpecialPriceDiscount) {
						if got.ID != testID || got.Name != "n" {
							t.Fatalf("specialpricediscount = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Ошибки входа — до HTTP; ResponseInfo прошлого вызова сброшен.
func TestServiceValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(client), t.Context()
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id": func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":   func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no name": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, specialpricediscount.CreateRequest{Active: true}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, specialpricediscount.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", update, o...); return err },
		"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}
