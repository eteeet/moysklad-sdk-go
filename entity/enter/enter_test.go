package enter_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/enter"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	fixtureID = "7cfff21a-533b-11e6-8a84-bae50000001f"
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

func service(c *transport.Client) *enter.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := enter.NewService(c)
	return s
}

var (
	organization = ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization")
	store        = ref("entity/store/33333333-3333-3333-3333-333333333333", "store")
	fullCreate   = enter.CreateRequest{
		Organization: organization, Store: store,
		Name: model.Value("DOC-1"), Description: model.Value(""), SyncID: model.Value("22222222-2222-2222-2222-222222222222"),
		Moment:     model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		Applicable: model.Value(false), Shared: model.Value(true),
		Rate:       model.Value(model.DocumentRateRequest{Currency: model.Value(ref("entity/currency/66666666-6666-6666-6666-666666666666", "currency"))}),
		Attributes: model.Value([]model.Attribute{}),
		Overhead:   model.Value(model.DocumentOverheadRequest{Sum: model.Value(1050.5), Distribution: model.Value(model.OverheadDistributionWeight)}),
		Positions: model.Value([]enter.PositionRequest{{
			ItemRequest: positions.ItemRequest{
				Assortment: model.Value(ref("entity/product/88888888-8888-8888-8888-888888888888", "product")),
				Quantity:   model.Value(2.0),
			},
			Price: model.Value(12345.5), Reason: model.Value("излишек"), Country: model.Value(ref("entity/country/44444444-4444-4444-4444-444444444444", "country")), GTD: model.Value(positions.GTD{Name: "10130000/010120/0000001"}), Things: model.Value([]string{"sn-1"}),
		}}),
	}
	fullUpdate = enter.UpdateRequest{
		Store:      model.Value(ref("entity/store/99999999-9999-9999-9999-999999999999", "store")),
		Project:    model.Null[model.Reference](),
		Applicable: model.Value(true),
		Positions: model.Value([]enter.PositionRequest{{
			ItemRequest: positions.ItemRequest{CoreRequest: positions.CoreRequest{ID: model.Value("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")}, Quantity: model.Value(3.0)},
			Price:       model.Value(0.0),
		}}),
	}
)

func TestEnterFixture(t *testing.T) {
	data := read(t, "upstream/enter.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got enter.Enter
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Organization == nil || got.Store == nil || got.Positions == nil || len(got.Positions.Rows) == 0 {
				t.Fatalf("document = %+v", got)
			}
			row := got.Positions.Rows[0]
			if got.Overhead == nil || got.Overhead.Distribution != "weight" || row.GTD == nil || row.Country == nil || row.Reason == "" {
				t.Fatalf("document = %+v, row = %+v", got, row)
			}
			var sparse enter.Enter
			if err := codec.Unmarshal([]byte(`{"id":"x","positions":{"meta":{"href":"h","size":2}},"future":1}`), &sparse); err != nil ||
				sparse.Positions.Rows != nil || sparse.Organization != nil || !sparse.Created.IsZero() || !sparse.Deleted.IsZero() || sparse.Overhead != nil || sparse.Inventory != nil {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestEnterRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", fullCreate, read(t, "custom/enter/create-full.golden.json")},
			{"create minimal", enter.CreateRequest{Organization: organization, Store: store}, []byte(`{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},"store":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/store/33333333-3333-3333-3333-333333333333","type":"store"}}}`)},
			{"update", fullUpdate, read(t, "custom/enter/update.golden.json")},
			{"update empty positions", enter.UpdateRequest{Positions: model.Value([]enter.PositionRequest{})}, []byte(`{"positions":[]}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[enter.CreateRequest](), reflect.TypeFor[enter.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Sum", "Files", "Inventory", "Agent"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only, unsupported or foreign %s", typ.Name(), name)
			}
		}
	}
}

func TestEnterHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и сумма — только из ответа.
	const response = `{"id":"server-id","name":"server-name","sum":24691,"positions":{"meta":{"href":"https://example.test/positions","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*enter.Enter, error)
	}{
		{"get", "GET", "entity/enter/" + testID, nil, func(ctx context.Context, c *transport.Client) (*enter.Enter, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/enter", read(t, "custom/enter/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*enter.Enter, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", "entity/enter/" + testID, read(t, "custom/enter/update.golden.json"), func(ctx context.Context, c *transport.Client) (*enter.Enter, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *enter.Enter) {
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

func TestEnterGetPositions(t *testing.T) {
	for _, page := range []positions.PageOptions{{}, {Limit: 1, Offset: 7}} {
		query := "limit=1000&offset=0"
		if page.Limit == 1 {
			query = "limit=1&offset=7"
		}
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[enter.Position], error) {
			return service(c).GetPositions(ctx, testID, page)
		}
		const body = `{"meta":{"href":"https://example.test/positions"},"rows":[{"id":"row","quantity":2,"reason":"излишек","overhead":150,"gtd":{"name":"g"},"slot":{"meta":{"href":"s"}}}]}`
		for _, tc := range documenttest.ResponseCases(body) {
			t.Run(tc.Name+"?"+query, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", "entity/enter/"+testID+"/positions?"+query, nil, invoke, func(got *model.Collection[enter.Position]) {
					row := got.Rows[0]
					if row.Reason != "излишек" || row.Overhead != 150 || row.GTD == nil || row.Slot == nil {
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
func TestEnterLimits(t *testing.T) {
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
	rows := func(n int) model.Field[[]enter.PositionRequest] {
		return model.Value(make([]enter.PositionRequest, n))
	}
	if _, err := s.Create(t.Context(), enter.CreateRequest{Organization: organization, Store: store, Positions: rows(1000)}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(t.Context(), testID, enter.UpdateRequest{Positions: rows(1000)}); err != nil {
		t.Fatal(err)
	}
	for name, err := range map[string]error{
		"create 1001": func() error {
			_, err := s.Create(t.Context(), enter.CreateRequest{Organization: organization, Store: store, Positions: rows(1001)})
			return err
		}(),
		"update 1001": func() error {
			_, err := s.Update(t.Context(), testID, enter.UpdateRequest{Positions: rows(1001)})
			return err
		}(),
		"no store": func() error {
			_, err := s.Create(t.Context(), enter.CreateRequest{Organization: organization})
			return err
		}(),
		"no organization": func() error { _, err := s.Create(t.Context(), enter.CreateRequest{Store: store}); return err }(),
		"empty update":    func() error { _, err := s.Update(t.Context(), testID, enter.UpdateRequest{}); return err }(),
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
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "enter") {
			t.Errorf("%s: err = %v", name, err)
		}
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := enter.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestEnterList(t *testing.T) {
	listtest.Run(t, "enter", read(t, "upstream/enter.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[enter.Enter] {
		return service(c).List
	}, func(t *testing.T, got enter.Enter) {
		if got.ID != fixtureID || got.Positions == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
