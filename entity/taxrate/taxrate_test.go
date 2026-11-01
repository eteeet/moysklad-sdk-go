package taxrate_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/taxrate"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "6031a4ab-fec1-11ec-0a80-059200000007"

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ref(href, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: href, Type: model.MetaType(typ), MediaType: "application/json"}}
}

func rate(v float64) *float64 { return &v }

var (
	full = taxrate.CreateRequest{Rate: rate(33), Comment: model.Value("Ставка на продукты с ГМО"), Archived: model.Value(false),
		Shared: model.Value(true),
		Owner:  model.Value(ref("https://api.moysklad.ru/api/remap/1.2/entity/employee/c6f50a9a-fe92-11ec-0a82-09860000027a", "employee")),
		Group:  model.Value(ref("https://api.moysklad.ru/api/remap/1.2/entity/group/c6bcfbf3-fe92-11ec-0a82-062000000019", "group"))}
	update = taxrate.UpdateRequest{Rate: model.Value(28.0), Comment: model.Value(""), Archived: model.Value(true)}
)

func TestTaxRateFixture(t *testing.T) {
	data := read(t, "upstream/tax_rate.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got taxrate.TaxRate
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.Rate != 33 || got.Comment == "" || !got.Shared || got.Owner == nil || got.Group == nil || got.Updated.IsZero() {
				t.Fatalf("taxrate = %+v", got)
			}
			// Предустановленная ставка: без owner, group и shared.
			var system taxrate.TaxRate
			if err := codec.Unmarshal([]byte(`{"id":"x","rate":0,"futureField":1}`), &system); err != nil || system.Owner != nil || system.Shared {
				t.Fatalf("system = %+v, err = %v", system, err)
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
			{"create full", full, read(t, "custom/taxrate/create-full.golden.json")},
			{"create zero rate", taxrate.CreateRequest{Rate: rate(0)}, read(t, "custom/taxrate/create-zero.golden.json")},
			{"update", update, read(t, "custom/taxrate/update.golden.json")},
			{"update zero rate", taxrate.UpdateRequest{Rate: model.Value(0.0)}, []byte(`{"rate":0}`)},
			{"update null owner", taxrate.UpdateRequest{Owner: model.Null[model.Reference]()}, []byte(`{"owner":null}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[taxrate.CreateRequest](), reflect.TypeFor[taxrate.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func service(c *transport.Client) *taxrate.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := taxrate.NewService(c)
	return s
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/taxrate/` + testID + `"},"id":"` + testID + `","rate":33}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*taxrate.TaxRate, error)
	}{
		{"get", "GET", "entity/taxrate/" + testID, nil, func(ctx context.Context, c *transport.Client) (*taxrate.TaxRate, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/taxrate", read(t, "custom/taxrate/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*taxrate.TaxRate, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/taxrate/" + testID, read(t, "custom/taxrate/update.golden.json"), func(ctx context.Context, c *transport.Client) (*taxrate.TaxRate, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *taxrate.TaxRate) {
						if got.ID != testID || got.Rate != 33 {
							t.Fatalf("taxrate = %+v", got)
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(client)
	if _, err := s.Create(t.Context(), taxrate.CreateRequest{Comment: model.Value("x")}); err == nil {
		t.Error("create without rate accepted")
	}
	if _, err := s.Update(t.Context(), testID, taxrate.UpdateRequest{}); err == nil {
		t.Error("empty update accepted")
	}
	if _, err := s.Get(t.Context(), "bad"); err == nil {
		t.Error("bad id accepted")
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestServiceList(t *testing.T) {
	listtest.Run(t, "taxrate", read(t, "upstream/tax_rate.json"), "", "archived", func(c *transport.Client) listtest.ListFunc[taxrate.TaxRate] {
		return service(c).List
	}, func(t *testing.T, got taxrate.TaxRate) {
		if got.ID != testID || got.Rate != 33 {
			t.Fatalf("row = %+v", got)
		}
	})
}
