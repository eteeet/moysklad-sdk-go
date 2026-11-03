package accumulationdiscount_test

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/discount"

	"github.com/eteeet/moysklad-sdk-go/entity/accumulationdiscount"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "dce08f7f-9a09-11e6-8af5-581e0000007e"

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
	full   = accumulationdiscount.CreateRequest{Name: "discountName", Active: true, AllProducts: false, AllAgents: false, AgentTags: model.Value([]string{"tag1", "tag2"}), Assortment: model.Value([]model.Reference{ref("entity/product/833eac42-b6f4-11ea-ac12-000e00000006", "product")}), ProductFolders: model.Value([]model.Reference{ref("entity/productfolder/39c62b64-a722-11ea-ac12-000d00000015", "productfolder")}), Levels: model.Value([]discount.Level{{Amount: 100, Discount: 10}})}
	update = accumulationdiscount.UpdateRequest{Levels: model.Value([]discount.Level{}), AgentTags: model.Null[[]string](), Active: model.Value(false)}
)

func TestAccumulationDiscountFixture(t *testing.T) {
	data := read(t, "upstream/accumulation_discount.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got accumulationdiscount.AccumulationDiscount
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.Name != "Скидки на сапоги" || !got.Active || !got.AllAgents || got.AllProducts || len(got.Levels) != 2 || got.Levels[0].Amount != 100000 || got.Levels[0].Discount != 10 || len(got.Assortment) == 0 || len(got.ProductFolders) == 0 {
				t.Fatalf("accumulationdiscount = %+v", got)
			}
			// Неизвестные поля не ломают разбор, отсутствующие не выдумываются.
			var sparse accumulationdiscount.AccumulationDiscount
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
			{"create full", full, read(t, "custom/accumulationdiscount/create-full.golden.json")},
			{"update", update, read(t, "custom/accumulationdiscount/update.golden.json")},
			{"create required false", accumulationdiscount.CreateRequest{Name: "n"}, []byte(`{"name":"n","active":false,"allProducts":false,"allAgents":false}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[accumulationdiscount.CreateRequest](), reflect.TypeFor[accumulationdiscount.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func service(c *transport.Client) *accumulationdiscount.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := accumulationdiscount.NewService(c)
	return s
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/accumulationdiscount/` + testID + `"},"id":"` + testID + `","name":"n"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*accumulationdiscount.AccumulationDiscount, error)
	}{
		{"get", "GET", "entity/accumulationdiscount/" + testID, nil, func(ctx context.Context, c *transport.Client) (*accumulationdiscount.AccumulationDiscount, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/accumulationdiscount", read(t, "custom/accumulationdiscount/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*accumulationdiscount.AccumulationDiscount, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/accumulationdiscount/" + testID, read(t, "custom/accumulationdiscount/update.golden.json"), func(ctx context.Context, c *transport.Client) (*accumulationdiscount.AccumulationDiscount, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *accumulationdiscount.AccumulationDiscount) {
						if got.ID != testID || got.Name != "n" {
							t.Fatalf("accumulationdiscount = %+v", got)
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
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
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
			_, err := s.Create(ctx, accumulationdiscount.CreateRequest{Active: true}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, accumulationdiscount.UpdateRequest{}, o...)
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
