package retaildemand_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/retaildemand"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	fixtureID = "cdccba00-0563-11e6-9464-e4de00000090"
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

func service(c *transport.Client) *retaildemand.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := retaildemand.NewService(c)
	return s
}

var (
	shift = ref("entity/retailshift/11111111-1111-1111-1111-111111111111", "retailshift")
	sale  = retaildemand.CreateRequest{
		RetailShift: shift,
		Name:        model.Value("R-1"), Moment: model.Value(documenttest.DateTime("2026-11-02 10:15:00.000")), Applicable: model.Value(true),
		Agent:     model.Value(ref("entity/counterparty/22222222-2222-2222-2222-222222222222", "counterparty")),
		Store:     model.Value(ref("entity/store/33333333-3333-3333-3333-333333333333", "store")),
		TaxSystem: model.Value("PATENT_BASED"), VATEnabled: model.Value(false),
		CashSum: model.Value(15000.0), NoCashSum: model.Value(0.0), PrepaymentCashSum: model.Value(0.0),
		CheckNumber: model.Value(0), SessionNumber: model.Value(17),
		GiftCards: model.Value([]retaildemand.GiftCardRequest{{Name: "123457"}}),
		Positions: model.Value([]retaildemand.PositionRequest{{
			ItemRequest: positions.ItemRequest{
				Assortment: model.Value(ref("entity/product/44444444-4444-4444-4444-444444444444", "product")),
				Quantity:   model.Value(1.0),
			},
			Price: model.Value(15000.0), Discount: model.Value(0.0), VAT: model.Value(0),
			TrackingCodes: model.Value([]positions.TrackingCode{{CIS: "010460406000301021N4N57RTCBUZQZ", Type: "trackingcode"}}),
		}}),
	}
	update = retaildemand.UpdateRequest{
		Description:  model.Value(""),
		Agent:        model.Value(ref("entity/counterparty/55555555-5555-5555-5555-555555555555", "counterparty")),
		AgentAccount: model.Null[model.Reference](),
		QRSum:        model.Value(15000.0),
		Positions: model.Value([]retaildemand.PositionRequest{{
			ItemRequest: positions.ItemRequest{CoreRequest: positions.CoreRequest{ID: model.Value("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")}, Quantity: model.Value(1.0)},
		}}),
	}
)

func TestRetailDemandFixture(t *testing.T) {
	data := read(t, "upstream/retail_demand.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got retaildemand.RetailDemand
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.RetailShift == nil || got.RetailStore == nil || got.Organization == nil || got.Cheque == nil ||
				got.Cheque.FNNumber != "110011" || got.Cheque.Time.IsZero() || len(got.GiftCards) != 2 || got.GiftCards[1].PaymentSum != 1000 ||
				got.Positions == nil || len(got.Positions.Rows) == 0 || got.TaxSystem == "" {
				t.Fatalf("sale = %+v", got)
			}
			var sparse retaildemand.RetailDemand
			if err := codec.Unmarshal([]byte(`{"id":"x","positions":{"meta":{"href":"h","size":1}},"future":1}`), &sparse); err != nil ||
				sparse.Cheque != nil || sparse.GiftCards != nil || sparse.Positions.Rows != nil || sparse.Agent != nil {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestRetailDemandRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", sale, read(t, "custom/retaildemand/create-full.golden.json")},
			{"create minimal", retaildemand.CreateRequest{RetailShift: shift}, []byte(`{"retailShift":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/retailshift/11111111-1111-1111-1111-111111111111","type":"retailshift"}}}`)},
			{"update", update, read(t, "custom/retaildemand/update.golden.json")},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[retaildemand.CreateRequest](), reflect.TypeFor[retaildemand.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Sum", "VATSum",
			"PayedSum", "Cheque", "Organization", "OrganizationBranch", "RetailStore", "Files"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func TestRetailDemandHTTP(t *testing.T) {
	// Ответ сервера: суммы и чек — только из ответа.
	const response = `{"id":"server-id","name":"server-name","sum":15000,"cashSum":15000,"positions":{"meta":{"href":"https://example.test/positions","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*retaildemand.RetailDemand, error)
	}{
		{"get", "GET", "entity/retaildemand/" + testID, nil, func(ctx context.Context, c *transport.Client) (*retaildemand.RetailDemand, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/retaildemand", read(t, "custom/retaildemand/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*retaildemand.RetailDemand, error) {
			return service(c).Create(ctx, sale)
		}},
		{"update", "PUT", "entity/retaildemand/" + testID, read(t, "custom/retaildemand/update.golden.json"), func(ctx context.Context, c *transport.Client) (*retaildemand.RetailDemand, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *retaildemand.RetailDemand) {
						if got.ID != "server-id" || got.Sum != 15000 || got.CashSum != 15000 || got.Positions.Rows != nil {
							t.Fatalf("sale = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func TestRetailDemandGetPositions(t *testing.T) {
	for _, page := range []positions.PageOptions{{}, {Limit: 1, Offset: 7, Fields: "stock"}, {Limit: 1, Offset: 7, Fields: "declaration"}} {
		query := "limit=1000&offset=0"
		if page.Fields != "" {
			query = "fields=" + page.Fields + "&limit=1&offset=7"
		}
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[retaildemand.Position], error) {
			return service(c).GetPositions(ctx, testID, page)
		}
		const body = `{"meta":{"href":"https://example.test/positions"},"rows":[{"id":"row","quantity":1,"cost":500,"trackingCodes":[{"cis":"c","type":"trackingcode"}],"trackingCodes_1162":[{"cis_1162":"c1162","type":"trackingcode"}]}]}`
		for _, tc := range documenttest.ResponseCases(body) {
			t.Run(tc.Name+"?"+query, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", "entity/retaildemand/"+testID+"/positions?"+query, nil, invoke, func(got *model.Collection[retaildemand.Position]) {
					if r := got.Rows[0]; r.Cost != 500 || len(r.TrackingCodes) != 1 || r.TrackingCodes1162[0].CIS1162 != "c1162" {
						t.Fatalf("row = %+v", r)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

// Ровно 1000 строк уходят, 1001 и пустая смена отклоняются до запроса. Суммы оплат и
// активность смены проверяет сервер: несходящаяся оплата уходит как есть.
func TestRetailDemandLimits(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Write([]byte(`{"id":"server-id"}`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(c)
	rows := func(n int) model.Field[[]retaildemand.PositionRequest] {
		return model.Value(make([]retaildemand.PositionRequest, n))
	}
	if _, err := s.Create(t.Context(), retaildemand.CreateRequest{RetailShift: shift, Positions: rows(1000), QRSum: model.Value(1.0), CashSum: model.Value(1.0)}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(t.Context(), testID, retaildemand.UpdateRequest{Positions: rows(1000)}); err != nil {
		t.Fatal(err)
	}
	for name, err := range map[string]error{
		"create 1001": func() error {
			_, err := s.Create(t.Context(), retaildemand.CreateRequest{RetailShift: shift, Positions: rows(1001)})
			return err
		}(),
		"update 1001": func() error {
			_, err := s.Update(t.Context(), testID, retaildemand.UpdateRequest{Positions: rows(1001)})
			return err
		}(),
		"no shift":      func() error { _, err := s.Create(t.Context(), retaildemand.CreateRequest{}); return err }(),
		"empty update":  func() error { _, err := s.Update(t.Context(), testID, retaildemand.UpdateRequest{}); return err }(),
		"bad id":        func() error { _, err := s.Get(t.Context(), "bad"); return err }(),
		"bad update id": func() error { _, err := s.Update(t.Context(), "bad", update); return err }(),
		"fields": func() error {
			_, err := s.GetPositions(t.Context(), testID, positions.PageOptions{Fields: "stock,declaration"})
			return err
		}(),
	} {
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "retaildemand") {
			t.Errorf("%s: err = %v", name, err)
		}
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := retaildemand.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestRetailDemandList(t *testing.T) {
	listtest.Run(t, "retaildemand", read(t, "upstream/retail_demand.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[retaildemand.RetailDemand] {
		return service(c).List
	}, func(t *testing.T, got retaildemand.RetailDemand) {
		if got.ID != fixtureID || got.RetailShift == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
