package loss_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/loss"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	fixtureID = "b014dab4-4f42-11e6-8a84-bae500000006"
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

func service(c *transport.Client) *loss.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := loss.NewService(c)
	return s
}

var (
	organization = ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization")
	store        = ref("entity/store/33333333-3333-3333-3333-333333333333", "store")
	fullCreate   = loss.CreateRequest{
		Organization: organization, Store: store,
		Name: model.Value("DOC-1"), Description: model.Value(""), SyncID: model.Value("22222222-2222-2222-2222-222222222222"),
		Moment:     model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		Applicable: model.Value(false), Shared: model.Value(true),
		Rate:        model.Value(model.DocumentRateRequest{Currency: model.Value(ref("entity/currency/66666666-6666-6666-6666-666666666666", "currency"))}),
		Attributes:  model.Value([]model.Attribute{}),
		SalesReturn: model.Value(ref("entity/salesreturn/55555555-5555-5555-5555-555555555555", "salesreturn")),
		Positions: model.Value([]loss.PositionRequest{{
			ItemRequest: positions.ItemRequest{
				Assortment: model.Value(ref("entity/product/88888888-8888-8888-8888-888888888888", "product")),
				Quantity:   model.Value(2.0),
			},
			Price: model.Value(12345.5), Reason: model.Value("бой"), Slot: model.Value(ref("entity/store/33333333-3333-3333-3333-333333333333/slots/44444444-4444-4444-4444-444444444444", "slot")), Pack: model.Value(positions.PackRequest{ID: "77777777-7777-7777-7777-777777777777"}),
		}}),
	}
	fullUpdate = loss.UpdateRequest{
		Store:      model.Value(ref("entity/store/99999999-9999-9999-9999-999999999999", "store")),
		Project:    model.Null[model.Reference](),
		Applicable: model.Value(true),
		Positions: model.Value([]loss.PositionRequest{{
			ItemRequest: positions.ItemRequest{CoreRequest: positions.CoreRequest{ID: model.Value("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")}, Quantity: model.Value(3.0)},
			Price:       model.Value(0.0),
		}}),
	}
)

func TestLossFixture(t *testing.T) {
	data := read(t, "upstream/loss.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got loss.Loss
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Organization == nil || got.Store == nil || got.Positions == nil || len(got.Positions.Rows) == 0 {
				t.Fatalf("document = %+v", got)
			}
			row := got.Positions.Rows[0]
			if got.SalesReturn == nil || got.ExpenseItem == nil || got.Inventory == nil || row.Pack == nil || row.Pack.UOM == nil || row.Slot == nil || row.Declaration == nil {
				t.Fatalf("document = %+v, row = %+v", got, row)
			}
			var sparse loss.Loss
			if err := codec.Unmarshal([]byte(`{"id":"x","positions":{"meta":{"href":"h","size":2}},"future":1}`), &sparse); err != nil ||
				sparse.Positions.Rows != nil || sparse.Organization != nil || !sparse.Created.IsZero() || !sparse.Deleted.IsZero() || sparse.SalesReturn != nil || sparse.ExpenseItem != nil {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestLossRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", fullCreate, read(t, "custom/loss/create-full.golden.json")},
			{"create minimal", loss.CreateRequest{Organization: organization, Store: store}, []byte(`{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},"store":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/store/33333333-3333-3333-3333-333333333333","type":"store"}}}`)},
			{"update", fullUpdate, read(t, "custom/loss/update.golden.json")},
			{"update empty positions", loss.UpdateRequest{Positions: model.Value([]loss.PositionRequest{})}, []byte(`{"positions":[]}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[loss.CreateRequest](), reflect.TypeFor[loss.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Sum", "Files", "ExpenseItem", "Agent"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only, unsupported or foreign %s", typ.Name(), name)
			}
		}
	}
}

func TestLossHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и сумма — только из ответа.
	const response = `{"id":"server-id","name":"server-name","sum":24691,"positions":{"meta":{"href":"https://example.test/positions","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*loss.Loss, error)
	}{
		{"get", "GET", "entity/loss/" + testID, nil, func(ctx context.Context, c *transport.Client) (*loss.Loss, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/loss", read(t, "custom/loss/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*loss.Loss, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", "entity/loss/" + testID, read(t, "custom/loss/update.golden.json"), func(ctx context.Context, c *transport.Client) (*loss.Loss, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *loss.Loss) {
						if got.ID != "server-id" || got.Name != "server-name" || got.Sum != 24691 || got.Positions.Rows != nil {
							t.Fatalf("document = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func TestLossGetPositions(t *testing.T) {
	for _, page := range []positions.PageOptions{{}, {Limit: 1, Offset: 7, Fields: "declaration"}} {
		query := "limit=1000&offset=0"
		if page.Limit == 1 {
			query = "fields=declaration&limit=1&offset=7"
		}
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[loss.Position], error) {
			return service(c).GetPositions(ctx, testID, page)
		}
		const body = `{"meta":{"href":"https://example.test/positions"},"rows":[{"id":"row","quantity":2,"price":100,"reason":"бой","declaration":[{"gtd":"10130000/010120/0000001","rnpt":"","quantity":2}]}]}`
		for _, tc := range documenttest.ResponseCases(body) {
			t.Run(tc.Name+"?"+query, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", "entity/loss/"+testID+"/positions?"+query, nil, invoke, func(got *model.Collection[loss.Position]) {
					row := got.Rows[0]
					if row.Reason != "бой" || len(row.Declaration) != 1 || row.Declaration[0].Quantity != 2 {
						t.Fatalf("row = %+v", row)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

// Ровно 1000 строк уходят; 1001, пустые обязательные ссылки и неподдержанный fields
// отклоняются до запроса.
func TestLossLimits(t *testing.T) {
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
	rows := func(n int) model.Field[[]loss.PositionRequest] {
		return model.Value(make([]loss.PositionRequest, n))
	}
	if _, err := s.Create(t.Context(), loss.CreateRequest{Organization: organization, Store: store, Positions: rows(1000)}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(t.Context(), testID, loss.UpdateRequest{Positions: rows(1000)}); err != nil {
		t.Fatal(err)
	}
	for name, err := range map[string]error{
		"create 1001": func() error {
			_, err := s.Create(t.Context(), loss.CreateRequest{Organization: organization, Store: store, Positions: rows(1001)})
			return err
		}(),
		"update 1001": func() error {
			_, err := s.Update(t.Context(), testID, loss.UpdateRequest{Positions: rows(1001)})
			return err
		}(),
		"no store": func() error {
			_, err := s.Create(t.Context(), loss.CreateRequest{Organization: organization})
			return err
		}(),
		"no organization": func() error { _, err := s.Create(t.Context(), loss.CreateRequest{Store: store}); return err }(),
		"empty update":    func() error { _, err := s.Update(t.Context(), testID, loss.UpdateRequest{}); return err }(),
		"bad id":          func() error { _, err := s.Get(t.Context(), "bad"); return err }(),
		"bad update id":   func() error { _, err := s.Update(t.Context(), "bad", fullUpdate); return err }(),
		"fields": func() error {
			_, err := s.GetPositions(t.Context(), testID, positions.PageOptions{Fields: "stock"})
			return err
		}(),
		"positions id": func() error { _, err := s.GetPositions(t.Context(), "bad", positions.PageOptions{}); return err }(),
		"positions limit": func() error {
			_, err := s.GetPositions(t.Context(), testID, positions.PageOptions{Limit: 1001})
			return err
		}(),
	} {
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "loss") {
			t.Errorf("%s: err = %v", name, err)
		}
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := loss.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestLossList(t *testing.T) {
	listtest.Run(t, "loss", read(t, "upstream/loss.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[loss.Loss] {
		return service(c).List
	}, func(t *testing.T, got loss.Loss) {
		if got.ID != fixtureID || got.Positions == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
