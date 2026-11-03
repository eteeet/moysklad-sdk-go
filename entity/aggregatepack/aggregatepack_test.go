package aggregatepack_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/aggregatepack"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "9497f3a3-bc09-11f1-0a82-14dc000065bb"

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
	full   = aggregatepack.CreateRequest{Barcodes: []model.Barcode{{"ean8": "00000000"}}, ChildrenList: model.Value([]model.Reference{ref("entity/aggregatepack/247d9890-bcd6-11f1-0a83-006800000000", "aggregatepack")}), Shared: model.Value(false)}
	update = aggregatepack.UpdateRequest{Barcodes: model.Value([]model.Barcode{{"gtin": "04600000000001"}}), ChildrenList: model.Value([]model.Reference{}), Owner: model.Null[model.Reference]()}
)

func TestAggregatePackFixture(t *testing.T) {
	data := read(t, "custom/aggregatepack/aggregatepack.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got aggregatepack.AggregatePack
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.Level != 1 || len(got.Barcodes) != 1 || got.Barcodes[0]["ean8"] != "00000000" || got.Moment.IsZero() || got.Positions == nil || got.Positions.Meta.Type != "aggregatepackposition" || !got.Shared || got.Owner == nil {
				t.Fatalf("aggregatepack = %+v", got)
			}
			// Неизвестные поля не ломают разбор, отсутствующие не выдумываются.
			var sparse aggregatepack.AggregatePack
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
			{"create full", full, read(t, "custom/aggregatepack/create-full.golden.json")},
			{"update", update, read(t, "custom/aggregatepack/update.golden.json")},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[aggregatepack.CreateRequest](), reflect.TypeFor[aggregatepack.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated", "Level", "Positions"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func service(c *transport.Client) *aggregatepack.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := aggregatepack.NewService(c)
	return s
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/aggregatepack/` + testID + `"},"id":"` + testID + `","name":"n"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*aggregatepack.AggregatePack, error)
	}{
		{"get", "GET", "entity/aggregatepack/" + testID, nil, func(ctx context.Context, c *transport.Client) (*aggregatepack.AggregatePack, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/aggregatepack", read(t, "custom/aggregatepack/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*aggregatepack.AggregatePack, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/aggregatepack/" + testID, read(t, "custom/aggregatepack/update.golden.json"), func(ctx context.Context, c *transport.Client) (*aggregatepack.AggregatePack, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *aggregatepack.AggregatePack) {
						if got.ID != testID {
							t.Fatalf("aggregatepack = %+v", got)
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
		"create no barcodes": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, aggregatepack.CreateRequest{Shared: model.Value(true)}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, aggregatepack.UpdateRequest{}, o...)
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
	listtest.Run(t, "aggregatepack", read(t, "custom/aggregatepack/aggregatepack.json"), "", "", func(c *transport.Client) listtest.ListFunc[aggregatepack.AggregatePack] {
		return service(c).List
	}, func(t *testing.T, got aggregatepack.AggregatePack) {
		if got.ID != testID || got.Level != 1 {
			t.Fatalf("row = %+v", got)
		}
	})
}
