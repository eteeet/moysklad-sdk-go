package salesreturn_test

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/salesreturn"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	fixtureID = "a8b8e1e3-3f85-11e6-8a84-bae50000008d"
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

func service(c *transport.Client) *salesreturn.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := salesreturn.NewService(c)
	return s
}

func row(quantity float64) positions.ItemRequest {
	return positions.ItemRequest{
		Assortment: model.Value(ref("entity/product/88888888-8888-8888-8888-888888888888", "product")),
		Quantity:   model.Value(quantity),
	}
}

var (
	organization = ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization")
	store        = ref("entity/store/33333333-3333-3333-3333-333333333333", "store")
	agent        = ref("entity/counterparty/22222222-2222-2222-2222-222222222222", "counterparty")
	based        = salesreturn.CreateRequest{
		Organization: organization, Store: store, Agent: agent,
		AgentAccount: model.Value(ref("entity/counterparty/22222222-2222-2222-2222-222222222222/accounts/44444444-4444-4444-4444-444444444444", "account")),
		Demand:       model.Value(ref("entity/demand/55555555-5555-5555-5555-555555555555", "demand")),
		Name:         model.Value("SR-1"), VATEnabled: model.Value(true), VATIncluded: model.Value(false),
		Positions: model.Value([]salesreturn.PositionRequest{{
			ItemRequest: row(1), Price: model.Value(10000.0), Discount: model.Value(0.0), VAT: model.Value(20), VATEnabled: model.Value(true),
		}}),
	}
	unbased = salesreturn.CreateRequest{
		Organization: organization, Store: store, Agent: agent,
		Description: model.Value(""), Applicable: model.Value(false),
		Positions: model.Value([]salesreturn.PositionRequest{{
			ItemRequest: row(2), Price: model.Value(5000.5), Cost: model.Value(3000.0),
			Country: model.Value(ref("entity/country/66666666-6666-6666-6666-666666666666", "country")),
			GTD:     model.Value(positions.GTD{Name: "10130000/010120/0000001"}), Things: model.Value([]string{"sn-1", "sn-2"}),
		}}),
	}
	update = salesreturn.UpdateRequest{
		Organization:        model.Value(ref("entity/organization/77777777-7777-7777-7777-777777777777", "organization")),
		OrganizationAccount: model.Value(ref("entity/organization/77777777-7777-7777-7777-777777777777/accounts/99999999-9999-9999-9999-999999999999", "account")),
		Contract:            model.Null[model.Reference](),
		Positions: model.Value([]salesreturn.PositionRequest{{
			ItemRequest: positions.ItemRequest{CoreRequest: positions.CoreRequest{ID: model.Value("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")}, Quantity: model.Value(0.5)},
		}}),
	}
)

func TestSalesReturnFixture(t *testing.T) {
	data := read(t, "upstream/sales_return.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got salesreturn.SalesReturn
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Organization == nil || got.Store == nil || got.Agent == nil || got.Demand == nil ||
				got.Positions == nil || len(got.Positions.Rows) == 0 {
				t.Fatalf("return = %+v", got)
			}
			if r := got.Positions.Rows[0]; r.GTD == nil || r.Country == nil || r.Assortment.Meta.Href == "" {
				t.Fatalf("row = %+v", r)
			}
			var unbased salesreturn.SalesReturn
			if err := codec.Unmarshal([]byte(`{"id":"x","agent":{"meta":{"href":"a"}},"positions":{"meta":{"href":"h","size":1}},"future":1}`), &unbased); err != nil ||
				unbased.Demand != nil || unbased.Losses != nil || unbased.FactureOut != nil || unbased.Positions.Rows != nil {
				t.Fatalf("unbased = %+v, err = %v", unbased, err)
			}
		})
	}
}

func TestSalesReturnRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create based", based, read(t, "custom/salesreturn/create-based.golden.json")},
			{"create unbased", unbased, read(t, "custom/salesreturn/create-unbased.golden.json")},
			{"update", update, read(t, "custom/salesreturn/update.golden.json")},
			{"update false", salesreturn.UpdateRequest{VATEnabled: model.Value(false)}, []byte(`{"vatEnabled":false}`)},
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
	// Запреты Update выражены типом: эти поля менять нельзя или они только для чтения.
	for _, typ := range []reflect.Type{reflect.TypeFor[salesreturn.CreateRequest](), reflect.TypeFor[salesreturn.UpdateRequest]()} {
		names := []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Sum", "VATSum", "PayedSum",
			"Files"}
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

func TestSalesReturnHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса.
	const response = `{"id":"server-id","name":"server-name","sum":10000,"positions":{"meta":{"href":"https://example.test/positions","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*salesreturn.SalesReturn, error)
	}{
		{"get", "GET", "entity/salesreturn/" + testID, nil, func(ctx context.Context, c *transport.Client) (*salesreturn.SalesReturn, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create based", "POST", "entity/salesreturn", read(t, "custom/salesreturn/create-based.golden.json"), func(ctx context.Context, c *transport.Client) (*salesreturn.SalesReturn, error) {
			return service(c).Create(ctx, based)
		}},
		{"create unbased", "POST", "entity/salesreturn", read(t, "custom/salesreturn/create-unbased.golden.json"), func(ctx context.Context, c *transport.Client) (*salesreturn.SalesReturn, error) {
			return service(c).Create(ctx, unbased)
		}},
		{"update", "PUT", "entity/salesreturn/" + testID, read(t, "custom/salesreturn/update.golden.json"), func(ctx context.Context, c *transport.Client) (*salesreturn.SalesReturn, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *salesreturn.SalesReturn) {
						if got.ID != "server-id" || got.Name != "server-name" || got.Sum != 10000 || got.Positions.Rows != nil {
							t.Fatalf("return = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func TestSalesReturnGetPositions(t *testing.T) {
	for _, page := range []positions.PageOptions{{}, {Limit: 1, Offset: 7, Fields: "stock"}} {
		query := "limit=1000&offset=0"
		if page.Fields != "" {
			query = "fields=stock&limit=1&offset=7"
		}
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[salesreturn.Position], error) {
			return service(c).GetPositions(ctx, testID, page)
		}
		// Без основания stock приходит без cost.
		const body = `{"meta":{"href":"https://example.test/positions"},"rows":[{"id":"row","quantity":2,"cost":3000,"stock":{"quantity":5,"reserve":0,"intransit":0,"available":5}}]}`
		for _, tc := range documenttest.ResponseCases(body) {
			t.Run(tc.Name+"?"+query, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", "entity/salesreturn/"+testID+"/positions?"+query, nil, invoke, func(got *model.Collection[salesreturn.Position]) {
					if r := got.Rows[0]; r.Cost != 3000 || r.Stock == nil || r.Stock.Cost != nil || r.Stock.Available != 5 {
						t.Fatalf("row = %+v", r)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

// Ровно 1000 строк уходят; 1001, пустые обязательные ссылки и неподдержанный fields
// отклоняются до запроса. Соответствие строк основанию проверяет сервер.
func TestSalesReturnLimits(t *testing.T) {
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
	rows := func(n int) model.Field[[]salesreturn.PositionRequest] {
		return model.Value(make([]salesreturn.PositionRequest, n))
	}
	if _, err := s.Create(t.Context(), salesreturn.CreateRequest{Organization: organization, Store: store, Agent: agent, Positions: rows(1000)}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(t.Context(), testID, salesreturn.UpdateRequest{Positions: rows(1000)}); err != nil {
		t.Fatal(err)
	}
	for name, err := range map[string]error{
		"create 1001": func() error {
			_, err := s.Create(t.Context(), salesreturn.CreateRequest{Organization: organization, Store: store, Agent: agent, Positions: rows(1001)})
			return err
		}(),
		"update 1001": func() error {
			_, err := s.Update(t.Context(), testID, salesreturn.UpdateRequest{Positions: rows(1001)})
			return err
		}(),
		"no agent": func() error {
			_, err := s.Create(t.Context(), salesreturn.CreateRequest{Organization: organization, Store: store})
			return err
		}(),
		"no store": func() error {
			_, err := s.Create(t.Context(), salesreturn.CreateRequest{Organization: organization, Agent: agent})
			return err
		}(),
		"no organization": func() error {
			_, err := s.Create(t.Context(), salesreturn.CreateRequest{Store: store, Agent: agent})
			return err
		}(),
		"empty update":  func() error { _, err := s.Update(t.Context(), testID, salesreturn.UpdateRequest{}); return err }(),
		"bad id":        func() error { _, err := s.Get(t.Context(), "bad"); return err }(),
		"bad update id": func() error { _, err := s.Update(t.Context(), "bad", update); return err }(),
		"fields": func() error {
			_, err := s.GetPositions(t.Context(), testID, positions.PageOptions{Fields: "declaration"})
			return err
		}(),
		"positions id": func() error { _, err := s.GetPositions(t.Context(), "bad", positions.PageOptions{}); return err }(),
	} {
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "salesreturn") {
			t.Errorf("%s: err = %v", name, err)
		}
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := salesreturn.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestSalesReturnList(t *testing.T) {
	listtest.Run(t, "salesreturn", read(t, "upstream/sales_return.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[salesreturn.SalesReturn] {
		return service(c).List
	}, func(t *testing.T, got salesreturn.SalesReturn) {
		if got.ID != fixtureID || got.Demand == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
