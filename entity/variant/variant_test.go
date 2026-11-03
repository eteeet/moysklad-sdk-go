package variant_test

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/variant"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "66666666-6666-6666-6666-666666666666"
	fixtureID = "c52d1e78-26b4-11f1-0a83-1f25000001bf"
)

func ref(href, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/" + href, Type: model.MetaType(typ), MediaType: "application/json"}}
}

var (
	product = ref("entity/product/11111111-1111-1111-1111-111111111111", "product")
	full    = variant.CreateRequest{
		Product:         product,
		Characteristics: []variant.CharacteristicValue{{ID: characteristicID, Value: "XL"}, {Name: "Цвет"}},
		Name:            model.Value("Футболка (XL)"), Code: model.Value("TS-XL"), ExternalCode: model.Value("ts-xl"),
		Description: model.Value(""), Article: model.Value("T01 (M01)"), Archived: model.Value(false), DiscountProhibited: model.Value(false),
		BuyPrice:   model.Value(model.Price{Value: 70000}),
		MinPrice:   model.Value(model.Price{Value: 0, Currency: ref("entity/currency/22222222-2222-2222-2222-222222222222", "currency")}),
		SalePrices: model.Value([]model.SalePrice{{Value: 110000.5, PriceType: model.PriceType{Meta: ref("context/companysettings/pricetype/88888888-8888-8888-8888-888888888888", "pricetype").Meta}}}),
		Barcodes:   model.Value([]model.Barcode{{"ean13": "2000000000077"}}),
		Packs: model.Value([]variant.PackRequest{{Barcodes: []model.Barcode{{"ean13": "2000000000114"}},
			ParentPack: ref("entity/product/11111111-1111-1111-1111-111111111111/packs/99999999-9999-9999-9999-999999999999", "goodpack")}}),
	}
	update = variant.UpdateRequest{
		Characteristics: model.Value([]variant.CharacteristicValue{{ID: characteristicID, Value: "L"}}),
		BuyPrice:        model.Null[model.Price](), Barcodes: model.Value([]model.Barcode{}), Packs: model.Value([]variant.PackRequest{}),
		Archived: model.Value(true),
	}
)

func TestVariantFixture(t *testing.T) {
	data := read(t, "upstream/variant.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got variant.Variant
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Name != "1 (v1)" || got.Product.Meta.Type != "product" || len(got.Characteristics) != 1 ||
				got.Characteristics[0].Value != "v1" || got.Characteristics[0].Meta.Type != "attributemetadata" ||
				got.BuyPrice == nil || got.BuyPrice.Value != 200 || got.MinPrice == nil || len(got.SalePrices) != 1 ||
				len(got.Barcodes) != 1 || len(got.Packs) != 1 || got.Packs[0].Barcodes[0]["ean13"] != "2000000000114" ||
				got.Packs[0].ParentPack.Meta.Type != "goodpack" || len(got.Things) != 1 || got.Images == nil {
				t.Fatalf("variant = %+v", got)
			}
			var sparse variant.Variant
			if err := codec.Unmarshal([]byte(`{"id":"x","characteristics":[{"id":"y","name":"Размер","value":"XL"}],"minimumStock":{"type":"FUTURE"},"future":1}`), &sparse); err != nil ||
				len(sparse.Characteristics) != 1 || sparse.MinimumStock == nil || sparse.MinimumStock.Type != "FUTURE" || sparse.BuyPrice != nil || sparse.Images != nil || sparse.Packs != nil {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestVariantRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", full, read(t, "custom/variant/create-full.golden.json")},
			{"create minimal", variant.CreateRequest{Product: product, Characteristics: []variant.CharacteristicValue{{Name: "Размер", Value: "M"}}},
				[]byte(`{"product":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/11111111-1111-1111-1111-111111111111","type":"product","mediaType":"application/json"}},"characteristics":[{"name":"Размер","value":"M"}]}`)},
			{"update", update, read(t, "custom/variant/update.golden.json")},
			{"update zero", variant.UpdateRequest{DiscountProhibited: model.Value(false), Code: model.Value("")}, []byte(`{"discountProhibited":false,"code":""}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[variant.CreateRequest](), reflect.TypeFor[variant.UpdateRequest](), reflect.TypeFor[variant.PackRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated", "Things", "Images"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only or deferred %s", typ.Name(), name)
			}
		}
	}
}

func TestVariantHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/variant/` + testID + `"},"id":"` + testID + `","name":"Футболка (XL)"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*variant.Variant, error)
	}{
		{"get", "GET", "entity/variant/" + testID, nil, func(ctx context.Context, c *transport.Client) (*variant.Variant, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/variant", read(t, "custom/variant/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*variant.Variant, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/variant/" + testID, read(t, "custom/variant/update.golden.json"), func(ctx context.Context, c *transport.Client) (*variant.Variant, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *variant.Variant) {
						if got.ID != testID || got.Name != "Футболка (XL)" {
							t.Fatalf("variant = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Без обязательных product и characteristics запрос не отправляется; товар не загружается скрыто.
func TestVariantValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(c)
	values := []variant.CharacteristicValue{{Name: "Размер", Value: "M"}}
	for name, input := range map[string]variant.CreateRequest{
		"no product":         {Characteristics: values},
		"no characteristics": {Product: product},
	} {
		if _, err := s.Create(t.Context(), input); err == nil || strings.Count(err.Error(), "moysklad:") != 1 ||
			!strings.Contains(err.Error(), "create variant") {
			t.Errorf("create %s: err = %v", name, err)
		}
	}
	if _, err := s.Update(t.Context(), testID, variant.UpdateRequest{}); err == nil || !strings.Contains(err.Error(), "update variant") {
		t.Errorf("empty update: err = %v", err)
	}
	if _, err := s.Get(t.Context(), "bad"); err == nil {
		t.Error("bad id accepted")
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

// Неизвестную характеристику (ни id, ни name не найдены) отклоняет сервер: SDK отправляет
// один запрос и возвращает APIError, описание не создаётся.
func TestVariantUnknownCharacteristic(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusPreconditionFailed)
		w.Write([]byte(`{"errors":[{"error":"characteristic not found","code":3006}]}`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	_, err = service(c).Create(t.Context(), variant.CreateRequest{Product: product, Characteristics: []variant.CharacteristicValue{{Name: "Нет такой", Value: "x"}}})
	if !sdkerrors.HasCode(err, 3006) || calls.Load() != 1 {
		t.Fatalf("err = %v, calls = %d", err, calls.Load())
	}
}

// Модификации товара отбираются фильтром productid; значение уходит без изменений.
func TestVariantListByProduct(t *testing.T) {
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("filter") != "productid=d6877ac4-26ae-11f1-0a83-1f2500000149" {
			t.Errorf("query = %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"meta":{"href":"h"},"rows":[]}`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	filter, err := query.Filter(query.Eq("productid", "d6877ac4-26ae-11f1-0a83-1f2500000149"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service(c).List(t.Context(), query.ListOptions{Filter: filter}); err != nil {
		t.Fatal(err)
	}
}

func TestVariantList(t *testing.T) {
	listtest.Run(t, "variant", read(t, "upstream/variant.json"), "minimumStock", "archived", func(c *transport.Client) listtest.ListFunc[variant.Variant] {
		return service(c).List
	}, func(t *testing.T, got variant.Variant) {
		if got.ID != fixtureID || len(got.Characteristics) != 1 {
			t.Fatalf("row = %+v", got)
		}
	})
}
