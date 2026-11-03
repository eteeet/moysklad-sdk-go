package customerorder_test

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/customerorder"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	fixtureID = "22591b08-2859-11f1-0a83-1f2500000247"
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

func service(c *transport.Client) *customerorder.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := customerorder.NewService(c)
	return s
}

var (
	full = customerorder.CreateRequest{
		Organization: ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization"),
		Agent:        ref("entity/counterparty/22222222-2222-2222-2222-222222222222", "counterparty"),
		Name:         model.Value("CO-1"), Description: model.Value(""), ExternalCode: model.Value("co-1"),
		SyncID:                model.Value("33333333-3333-3333-3333-333333333333"),
		Moment:                model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		DeliveryPlannedMoment: model.Value(documenttest.DateTime("2026-11-05 10:00:00.000")),
		Applicable:            model.Value(false), Shared: model.Value(true),
		Store:     model.Value(ref("entity/store/44444444-4444-4444-4444-444444444444", "store")),
		Rate:      model.Value(model.DocumentRateRequest{Currency: model.Value(ref("entity/currency/55555555-5555-5555-5555-555555555555", "currency"))}),
		TaxSystem: model.Value("GENERAL_TAX_SYSTEM"), VATEnabled: model.Value(true), VATIncluded: model.Value(false),
		ShipmentAddressFull: model.Value(model.AddressRequest{City: model.Value("Москва"), Street: model.Value("Тверская"), House: model.Value("1")}),
		Attributes:          model.Value([]model.Attribute{}),
		Positions: model.Value([]customerorder.PositionRequest{{
			ItemRequest: positions.ItemRequest{
				Assortment: model.Value(ref("entity/product/66666666-6666-6666-6666-666666666666", "product")),
				Quantity:   model.Value(2.0),
			},
			Price: model.Value(12345.5), Discount: model.Value(0.0), VAT: model.Value(20), VATEnabled: model.Value(true),
			Reserve: model.Value(1.0), Pack: model.Value(positions.PackRequest{ID: "77777777-7777-7777-7777-777777777777"}),
		}}),
	}
	update = customerorder.UpdateRequest{
		Agent:           model.Value(ref("entity/counterparty/88888888-8888-8888-8888-888888888888", "counterparty")),
		AgentAccount:    model.Value(ref("entity/counterparty/88888888-8888-8888-8888-888888888888/accounts/99999999-9999-9999-9999-999999999999", "account")),
		Contract:        model.Null[model.Reference](),
		ShipmentAddress: model.Value(""),
		Positions: model.Value([]customerorder.PositionRequest{{
			ItemRequest: positions.ItemRequest{CoreRequest: positions.CoreRequest{ID: model.Value("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")}, Quantity: model.Value(3.0)},
			Reserve:     model.Value(0.0),
		}}),
	}
)

func TestCustomerOrderFixture(t *testing.T) {
	data := read(t, "upstream/customer_order.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got customerorder.CustomerOrder
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Organization == nil || got.Agent == nil || got.Positions == nil ||
				len(got.Positions.Rows) == 0 || got.DeliveryPlannedMoment.IsZero() || got.TaxSystem == "" ||
				got.Demands == nil || got.Payments == nil || got.InvoicesOut == nil || got.Moves == nil ||
				got.PurchaseOrders == nil || got.Prepayments == nil || got.ProductionTasks == nil {
				t.Fatalf("order = %+v", got)
			}
			row := got.Positions.Rows[0]
			if row.Assortment.Meta.Href == "" || row.Quantity == 0 {
				t.Fatalf("row = %+v", row)
			}
			var sparse customerorder.CustomerOrder
			if err := codec.Unmarshal([]byte(`{"id":"x","taxSystem":"FUTURE","positions":{"meta":{"href":"h","size":2}},"future":1}`), &sparse); err != nil ||
				sparse.TaxSystem != "FUTURE" || sparse.Positions.Rows != nil || sparse.Agent != nil || sparse.Demands != nil {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestCustomerOrderRequestJSON(t *testing.T) {
	minimal := customerorder.CreateRequest{Organization: full.Organization, Agent: full.Agent}
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", full, read(t, "custom/customerorder/create-full.golden.json")},
			{"create minimal", minimal, []byte(`{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},"agent":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/counterparty/22222222-2222-2222-2222-222222222222","type":"counterparty"}}}`)},
			{"update", update, read(t, "custom/customerorder/update.golden.json")},
			{"update false and empty positions", customerorder.UpdateRequest{Applicable: model.Value(false), Positions: model.Value([]customerorder.PositionRequest{})}, []byte(`{"applicable":false,"positions":[]}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[customerorder.CreateRequest](), reflect.TypeFor[customerorder.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Sum", "VATSum",
			"PayedSum", "ShippedSum", "InvoicedSum", "ReservedSum", "Files"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only or unsupported %s", typ.Name(), name)
			}
		}
	}
}

func TestCustomerOrderHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и сумма — только из ответа.
	const response = `{"id":"server-id","name":"server-name","sum":24691,"positions":{"meta":{"href":"https://example.test/positions","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*customerorder.CustomerOrder, error)
	}{
		{"get", "GET", "entity/customerorder/" + testID, nil, func(ctx context.Context, c *transport.Client) (*customerorder.CustomerOrder, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/customerorder", read(t, "custom/customerorder/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*customerorder.CustomerOrder, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/customerorder/" + testID, read(t, "custom/customerorder/update.golden.json"), func(ctx context.Context, c *transport.Client) (*customerorder.CustomerOrder, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *customerorder.CustomerOrder) {
						if got.ID != "server-id" || got.Name != "server-name" || got.Sum != 24691 || got.Positions.Rows != nil {
							t.Fatalf("order = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func TestCustomerOrderGetPositions(t *testing.T) {
	for _, page := range []positions.PageOptions{{}, {Limit: 1, Offset: 7, Fields: "stock"}} {
		query := "limit=1000&offset=0"
		if page.Fields != "" {
			query = "fields=stock&limit=1&offset=7"
		}
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[customerorder.Position], error) {
			return service(c).GetPositions(ctx, testID, page)
		}
		const body = `{"meta":{"href":"https://example.test/positions"},"rows":[{"id":"row","quantity":2,"reserve":1.5,"shipped":0.5,"taxSystem":"PATENT_BASED","pack":{"id":"p"},"stock":{"cost":100,"quantity":10,"reserve":1,"intransit":0,"available":9}}]}`
		for _, tc := range documenttest.ResponseCases(body) {
			t.Run(tc.Name+"?"+query, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", "entity/customerorder/"+testID+"/positions?"+query, nil, invoke, func(got *model.Collection[customerorder.Position]) {
					row := got.Rows[0]
					if row.Reserve != 1.5 || row.Shipped != 0.5 || row.TaxSystem != "PATENT_BASED" || row.Pack == nil || row.Stock == nil || row.Stock.Available != 9 {
						t.Fatalf("row = %+v", row)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

// Ровно 1000 строк уходят, 1001 и пустые обязательные ссылки отклоняются до запроса.
func TestCustomerOrderLimits(t *testing.T) {
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
	rows := func(n int) model.Field[[]customerorder.PositionRequest] {
		return model.Value(make([]customerorder.PositionRequest, n))
	}
	if _, err := s.Create(t.Context(), customerorder.CreateRequest{Organization: full.Organization, Agent: full.Agent, Positions: rows(1000)}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(t.Context(), testID, customerorder.UpdateRequest{Positions: rows(1000)}); err != nil {
		t.Fatal(err)
	}
	for name, err := range map[string]error{
		"create 1001": func() error {
			_, err := s.Create(t.Context(), customerorder.CreateRequest{Organization: full.Organization, Agent: full.Agent, Positions: rows(1001)})
			return err
		}(),
		"update 1001": func() error {
			_, err := s.Update(t.Context(), testID, customerorder.UpdateRequest{Positions: rows(1001)})
			return err
		}(),
		"no agent": func() error {
			_, err := s.Create(t.Context(), customerorder.CreateRequest{Organization: full.Organization})
			return err
		}(),
		"no organization": func() error {
			_, err := s.Create(t.Context(), customerorder.CreateRequest{Agent: full.Agent})
			return err
		}(),
		"empty update":  func() error { _, err := s.Update(t.Context(), testID, customerorder.UpdateRequest{}); return err }(),
		"bad id":        func() error { _, err := s.Get(t.Context(), "bad"); return err }(),
		"bad update id": func() error { _, err := s.Update(t.Context(), "bad", update); return err }(),
		"fields": func() error {
			_, err := s.GetPositions(t.Context(), testID, positions.PageOptions{Fields: "declaration"})
			return err
		}(),
		"positions id": func() error { _, err := s.GetPositions(t.Context(), "bad", positions.PageOptions{}); return err }(),
		"positions limit": func() error {
			_, err := s.GetPositions(t.Context(), testID, positions.PageOptions{Limit: 1001})
			return err
		}(),
	} {
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "customerorder") {
			t.Errorf("%s: err = %v", name, err)
		}
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := customerorder.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestCustomerOrderList(t *testing.T) {
	listtest.Run(t, "customerorder", read(t, "upstream/customer_order.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[customerorder.CustomerOrder] {
		return service(c).List
	}, func(t *testing.T, got customerorder.CustomerOrder) {
		if got.ID != fixtureID || got.Positions == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
