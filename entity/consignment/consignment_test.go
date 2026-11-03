package consignment_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/consignment"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "c66f4b17-36e7-11e7-8a7f-40d000000113"

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
	full   = consignment.CreateRequest{Label: "Метка", Assortment: ref("entity/variant/7a83c422-3c64-11e6-8a84-bae500000012", "variant"), ExpiryDate: model.Value(model.DateTime(time.Date(2017, 12, 30, 21, 0, 0, 0, time.UTC))), Description: model.Value(""), Code: model.Value("c-1"), ExternalCode: model.Value("ext"), Archived: model.Value(false), Barcodes: model.Value([]model.Barcode{{"ean8": "20000000"}, {"ean13": "2000000000000"}, {"code128": "code128 barcode"}}), Attributes: model.Value([]model.Attribute{})}
	update = consignment.UpdateRequest{Label: model.Value("Метка 2"), ExpiryDate: model.Null[model.DateTime](), Barcodes: model.Value([]model.Barcode{})}
)

func TestConsignmentFixture(t *testing.T) {
	data := read(t, "upstream/consignment.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got consignment.Consignment
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.Name != "product / consignment" || got.Label != "consignment" || got.Code != "1012" || got.Archived || got.ExpiryDate.IsZero() || got.Assortment == nil || len(got.Barcodes) == 0 || got.Images == nil || len(got.Attributes) == 0 || got.Updated.IsZero() {
				t.Fatalf("consignment = %+v", got)
			}
			// Неизвестные поля не ломают разбор, отсутствующие не выдумываются.
			var sparse consignment.Consignment
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
			{"create full", full, read(t, "custom/consignment/create-full.golden.json")},
			{"update", update, read(t, "custom/consignment/update.golden.json")},
			{"create minimal", consignment.CreateRequest{Label: "L", Assortment: ref("entity/product/1", "product")}, []byte(`{"label":"L","assortment":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/1","type":"product","mediaType":"application/json"}}}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[consignment.CreateRequest](), reflect.TypeFor[consignment.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated", "Name", "Images"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func service(c *transport.Client) *consignment.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := consignment.NewService(c)
	return s
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/consignment/` + testID + `"},"id":"` + testID + `","name":"n"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*consignment.Consignment, error)
	}{
		{"get", "GET", "entity/consignment/" + testID, nil, func(ctx context.Context, c *transport.Client) (*consignment.Consignment, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/consignment", read(t, "custom/consignment/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*consignment.Consignment, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/consignment/" + testID, read(t, "custom/consignment/update.golden.json"), func(ctx context.Context, c *transport.Client) (*consignment.Consignment, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *consignment.Consignment) {
						if got.ID != testID || got.Name != "n" {
							t.Fatalf("consignment = %+v", got)
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
		"create no label": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, consignment.CreateRequest{Assortment: ref("entity/product/1", "product")}, o...)
			return err
		},
		"create no assortment": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, consignment.CreateRequest{Label: "L"}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, consignment.UpdateRequest{}, o...)
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
	listtest.Run(t, "consignment", read(t, "upstream/consignment.json"), "", "archived", func(c *transport.Client) listtest.ListFunc[consignment.Consignment] {
		return service(c).List
	}, func(t *testing.T, got consignment.Consignment) {
		if got.ID != testID || got.Label != "consignment" {
			t.Fatalf("row = %+v", got)
		}
	})
}
