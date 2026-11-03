package productfolder_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/productfolder"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "55555555-5555-5555-5555-555555555555"

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
	full   = productfolder.CreateRequest{Name: "Группа SDK", Description: model.Value(""), Code: model.Value("g-1"), ExternalCode: model.Value("ext"), ProductFolder: model.Value(ref("entity/productfolder/55555555-5555-5555-5555-555555555556", "productfolder")), Vat: model.Value(0), VatEnabled: model.Value(false), UseParentVat: model.Value(false), TaxSystem: model.Value("GENERAL_TAX_SYSTEM"), Shared: model.Value(true), Owner: model.Value(ref("entity/employee/faba7f37-2e58-11e6-8a84-bae500000028", "employee")), Group: model.Value(ref("entity/group/faba7f37-2e58-11e6-8a84-bae500000029", "group"))}
	update = productfolder.UpdateRequest{UseParentVat: model.Value(true), ProductFolder: model.Null[model.Reference](), TaxSystem: model.Value("TAX_SYSTEM_SAME_AS_GROUP")}
)

func TestProductFolderFixture(t *testing.T) {
	data := read(t, "upstream/product_folder.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got productfolder.ProductFolder
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.Name != "Электроника" || got.Vat != 20 || !got.VatEnabled || got.UseParentVat || got.EffectiveVat != 20 || !got.EffectiveVatEnabled || got.TaxSystem != "SIMPLIFIED_TAX_SYSTEM_INCOME" || got.PathName != "Электроника" || got.ProductFolder == nil || got.Updated.IsZero() {
				t.Fatalf("productfolder = %+v", got)
			}
			// Неизвестные поля не ломают разбор, отсутствующие не выдумываются.
			var sparse productfolder.ProductFolder
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
			{"create full", full, read(t, "custom/productfolder/create-full.golden.json")},
			{"update", update, read(t, "custom/productfolder/update.golden.json")},
			{"create minimal", productfolder.CreateRequest{Name: "G"}, []byte(`{"name":"G"}`)},
			{"update vat 0", productfolder.UpdateRequest{Vat: model.Value(0)}, []byte(`{"vat":0}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[productfolder.CreateRequest](), reflect.TypeFor[productfolder.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated", "Archived", "PathName", "EffectiveVat", "EffectiveVatEnabled"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func service(c *transport.Client) *productfolder.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := productfolder.NewService(c)
	return s
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/productfolder/` + testID + `"},"id":"` + testID + `","name":"n"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*productfolder.ProductFolder, error)
	}{
		{"get", "GET", "entity/productfolder/" + testID, nil, func(ctx context.Context, c *transport.Client) (*productfolder.ProductFolder, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/productfolder", read(t, "custom/productfolder/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*productfolder.ProductFolder, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/productfolder/" + testID, read(t, "custom/productfolder/update.golden.json"), func(ctx context.Context, c *transport.Client) (*productfolder.ProductFolder, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *productfolder.ProductFolder) {
						if got.ID != testID || got.Name != "n" {
							t.Fatalf("productfolder = %+v", got)
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
			_, err := s.Create(ctx, productfolder.CreateRequest{Code: model.Value("x")}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, productfolder.UpdateRequest{}, o...)
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

func TestServiceList(t *testing.T) {
	listtest.Run(t, "productfolder", read(t, "upstream/product_folder.json"), "", "archived", func(c *transport.Client) listtest.ListFunc[productfolder.ProductFolder] {
		return service(c).List
	}, func(t *testing.T, got productfolder.ProductFolder) {
		if got.ID != testID || got.Code != "ELEC" {
			t.Fatalf("row = %+v", got)
		}
	})
}
