package retailsalesreturn_test

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/retailsalesreturn"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	fixtureID = "c1ab79f5-1381-11e6-9464-e4de00000005"
)

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ref(path, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/" + path, Type: model.MetaType(typ)}}
}

func service(c *transport.Client) *retailsalesreturn.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := retailsalesreturn.NewService(c)
	return s
}

func item(price float64) retailsalesreturn.PositionRequest {
	return retailsalesreturn.PositionRequest{
		ItemRequest: positions.ItemRequest{
			Assortment: model.Value(ref("entity/product/88888888-8888-8888-8888-888888888888", "product")),
			Quantity:   model.Value(1.0),
		},
		Price: model.Value(price),
	}
}

func minimal() retailsalesreturn.CreateRequest {
	return retailsalesreturn.CreateRequest{
		Name:         "RR-1",
		Organization: ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization"),
		Store:        ref("entity/store/33333333-3333-3333-3333-333333333333", "store"),
		RetailStore:  ref("entity/retailstore/44444444-4444-4444-4444-444444444444", "retailstore"),
		RetailShift:  ref("entity/retailshift/55555555-5555-5555-5555-555555555555", "retailshift"),
		Agent:        ref("entity/counterparty/22222222-2222-2222-2222-222222222222", "counterparty"),
	}
}

func based() retailsalesreturn.CreateRequest {
	r := minimal()
	r.Demand = model.Value(ref("entity/retaildemand/66666666-6666-6666-6666-666666666666", "retaildemand"))
	r.Positions = model.Value([]retailsalesreturn.PositionRequest{item(15000)})
	return r
}

// Без основания источник требует cashSum и noCashSum: явный 0 уходит в JSON.
func unbased() retailsalesreturn.CreateRequest {
	r := minimal()
	row := item(12000)
	row.Cost = model.Value(7000.0)
	row.Things = model.Value([]string{"sn-1"})
	r.CashSum, r.NoCashSum = model.Value(0.0), model.Value(12000.0)
	r.Positions = model.Value([]retailsalesreturn.PositionRequest{row})
	return r
}

var update = retailsalesreturn.UpdateRequest{
	Name: model.Value("RR-2"), Contract: model.Null[model.Reference](), QRSum: model.Value(0.0),
	Positions: model.Value([]retailsalesreturn.PositionRequest{{
		ItemRequest: positions.ItemRequest{CoreRequest: positions.CoreRequest{ID: model.Value("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")}, Quantity: model.Value(1.0)},
	}}),
}

func TestRetailSalesReturnFixture(t *testing.T) {
	data := read(t, "upstream/retail_sales_return.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got retailsalesreturn.RetailSalesReturn
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Demand == nil || got.Demand.Meta.Type != "retaildemand" || got.RetailShift == nil ||
				got.RetailStore == nil || got.CashSum != 540000 || got.Positions == nil || len(got.Positions.Rows) == 0 ||
				got.Positions.Rows[0].Pack == nil {
				t.Fatalf("return = %+v", got)
			}
			var sparse retailsalesreturn.RetailSalesReturn
			if err := codec.Unmarshal([]byte(`{"id":"x","positions":{"meta":{"href":"h","size":1}},"future":1}`), &sparse); err != nil ||
				sparse.Demand != nil || sparse.Positions.Rows != nil {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestRetailSalesReturnRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create based", based(), read(t, "custom/retailsalesreturn/create-based.golden.json")},
			{"create based minimal", func() retailsalesreturn.CreateRequest {
				r := minimal()
				r.Demand = model.Value(ref("entity/retaildemand/66666666-6666-6666-6666-666666666666", "retaildemand"))
				return r
			}(), []byte(`{"name":"RR-1","organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},"store":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/store/33333333-3333-3333-3333-333333333333","type":"store"}},"retailStore":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/retailstore/44444444-4444-4444-4444-444444444444","type":"retailstore"}},"retailShift":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/retailshift/55555555-5555-5555-5555-555555555555","type":"retailshift"}},"agent":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/counterparty/22222222-2222-2222-2222-222222222222","type":"counterparty"}},"demand":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/retaildemand/66666666-6666-6666-6666-666666666666","type":"retaildemand"}}}`)},
			{"create unbased zero sums", func() retailsalesreturn.CreateRequest {
				r := minimal()
				r.CashSum, r.NoCashSum = model.Value(0.0), model.Value(0.0)
				return r
			}(), []byte(`{"name":"RR-1","organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},"store":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/store/33333333-3333-3333-3333-333333333333","type":"store"}},"retailStore":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/retailstore/44444444-4444-4444-4444-444444444444","type":"retailstore"}},"retailShift":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/retailshift/55555555-5555-5555-5555-555555555555","type":"retailshift"}},"agent":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/counterparty/22222222-2222-2222-2222-222222222222","type":"counterparty"}},"cashSum":0,"noCashSum":0}`)},
			{"create unbased", unbased(), read(t, "custom/retailsalesreturn/create-unbased.golden.json")},
			{"update", update, read(t, "custom/retailsalesreturn/update.golden.json")},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[retailsalesreturn.CreateRequest](), reflect.TypeFor[retailsalesreturn.UpdateRequest]()} {
		names := []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Sum", "VATSum", "OrganizationBranch"}
		if typ.Name() == "UpdateRequest" {
			names = append(names, "Agent", "AgentAccount", "Demand")
		}
		for _, name := range names {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains %s", typ.Name(), name)
			}
		}
	}
}

func TestRetailSalesReturnHTTP(t *testing.T) {
	const response = `{"id":"server-id","name":"server-name","sum":12000,"positions":{"meta":{"href":"https://example.test/positions","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*retailsalesreturn.RetailSalesReturn, error)
	}{
		{"get", "GET", "entity/retailsalesreturn/" + testID, nil, func(ctx context.Context, c *transport.Client) (*retailsalesreturn.RetailSalesReturn, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create based", "POST", "entity/retailsalesreturn", read(t, "custom/retailsalesreturn/create-based.golden.json"), func(ctx context.Context, c *transport.Client) (*retailsalesreturn.RetailSalesReturn, error) {
			return service(c).Create(ctx, based())
		}},
		{"create unbased", "POST", "entity/retailsalesreturn", read(t, "custom/retailsalesreturn/create-unbased.golden.json"), func(ctx context.Context, c *transport.Client) (*retailsalesreturn.RetailSalesReturn, error) {
			return service(c).Create(ctx, unbased())
		}},
		{"update", "PUT", "entity/retailsalesreturn/" + testID, read(t, "custom/retailsalesreturn/update.golden.json"), func(ctx context.Context, c *transport.Client) (*retailsalesreturn.RetailSalesReturn, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *retailsalesreturn.RetailSalesReturn) {
						if got.ID != "server-id" || got.Sum != 12000 || got.Positions.Rows != nil {
							t.Fatalf("return = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func TestRetailSalesReturnGetPositions(t *testing.T) {
	for _, page := range []positions.PageOptions{{}, {Limit: 1, Offset: 7, Fields: "stock"}} {
		query := "limit=1000&offset=0"
		if page.Fields != "" {
			query = "fields=stock&limit=1&offset=7"
		}
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[retailsalesreturn.Position], error) {
			return service(c).GetPositions(ctx, testID, page)
		}
		const body = `{"meta":{"href":"https://example.test/positions"},"rows":[{"id":"row","quantity":1,"cost":7000,"stock":{"quantity":3,"reserve":0,"intransit":0,"available":3}}]}`
		for _, tc := range documenttest.ResponseCases(body) {
			t.Run(tc.Name+"?"+query, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", "entity/retailsalesreturn/"+testID+"/positions?"+query, nil, invoke, func(got *model.Collection[retailsalesreturn.Position]) {
					if r := got.Rows[0]; r.Cost != 7000 || r.Stock == nil || r.Stock.Cost != nil {
						t.Fatalf("row = %+v", r)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

func TestRetailSalesReturnLimits(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Write([]byte(`{"id":"server-id"}`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(c)
	rows := func(n int) model.Field[[]retailsalesreturn.PositionRequest] {
		return model.Value(make([]retailsalesreturn.PositionRequest, n))
	}
	full := minimal()
	full.Positions = rows(1000)
	if _, err := s.Create(t.Context(), full); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(t.Context(), testID, retailsalesreturn.UpdateRequest{Positions: rows(1000)}); err != nil {
		t.Fatal(err)
	}
	without := func(clear func(*retailsalesreturn.CreateRequest)) error {
		r := minimal()
		clear(&r)
		_, err := s.Create(t.Context(), r)
		return err
	}
	for name, err := range map[string]error{
		"create 1001": func() error { r := minimal(); r.Positions = rows(1001); _, err := s.Create(t.Context(), r); return err }(),
		"update 1001": func() error {
			_, err := s.Update(t.Context(), testID, retailsalesreturn.UpdateRequest{Positions: rows(1001)})
			return err
		}(),
		"no name":         without(func(r *retailsalesreturn.CreateRequest) { r.Name = "" }),
		"no organization": without(func(r *retailsalesreturn.CreateRequest) { r.Organization = model.Reference{} }),
		"no store":        without(func(r *retailsalesreturn.CreateRequest) { r.Store = model.Reference{} }),
		"no retail store": without(func(r *retailsalesreturn.CreateRequest) { r.RetailStore = model.Reference{} }),
		"no retail shift": without(func(r *retailsalesreturn.CreateRequest) { r.RetailShift = model.Reference{} }),
		"no agent":        without(func(r *retailsalesreturn.CreateRequest) { r.Agent = model.Reference{} }),
		"empty update":    func() error { _, err := s.Update(t.Context(), testID, retailsalesreturn.UpdateRequest{}); return err }(),
		"bad id":          func() error { _, err := s.Get(t.Context(), "bad"); return err }(),
		"fields": func() error {
			_, err := s.GetPositions(t.Context(), testID, positions.PageOptions{Fields: "declaration"})
			return err
		}(),
	} {
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "retailsalesreturn") {
			t.Errorf("%s: err = %v", name, err)
		}
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := retailsalesreturn.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestRetailSalesReturnList(t *testing.T) {
	listtest.Run(t, "retailsalesreturn", read(t, "upstream/retail_sales_return.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[retailsalesreturn.RetailSalesReturn] {
		return service(c).List
	}, func(t *testing.T, got retailsalesreturn.RetailSalesReturn) {
		if got.ID != fixtureID || got.Demand == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
