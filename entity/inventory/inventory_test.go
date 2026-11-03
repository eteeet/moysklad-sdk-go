package inventory_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/inventory"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	rowID     = "34f6344f-015e-11e6-9464-e4de0000006c"
	fixtureID = "57a0bbfb-ad6d-11e6-8a84-bc52000000e7"
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

func service(c *transport.Client) *inventory.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := inventory.NewService(c)
	return s
}

var (
	organization = ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization")
	store        = ref("entity/store/33333333-3333-3333-3333-333333333333", "store")
	fullCreate   = inventory.CreateRequest{
		Organization: organization, Store: store,
		Name: model.Value("DOC-1"), Description: model.Value(""), SyncID: model.Value("22222222-2222-2222-2222-222222222222"),
		Moment:     model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		Shared:     model.Value(true),
		Attributes: model.Value([]model.Attribute{}),
		Enters:     model.Value([]model.Reference{ref("entity/enter/55555555-5555-5555-5555-555555555555", "enter")}),
		Losses:     model.Value([]model.Reference(nil)),
		Positions: model.Value(inventory.PositionRows{Rows: []inventory.PositionRequest{{
			ItemRequest: positions.ItemRequest{
				Assortment: model.Value(ref("entity/product/88888888-8888-8888-8888-888888888888", "product")),
				Quantity:   model.Value(140.0),
			},
			Price: model.Value(5000.0), CalculatedQuantity: model.Value(200.0), Pack: model.Value(positions.PackRequest{ID: "77777777-7777-7777-7777-777777777777"}),
		}}}),
	}
	fullUpdate = inventory.UpdateRequest{
		Store:  model.Value(ref("entity/store/99999999-9999-9999-9999-999999999999", "store")),
		State:  model.Null[model.Reference](),
		Losses: model.Value([]model.Reference{ref("entity/loss/44444444-4444-4444-4444-444444444444", "loss")}),
		Positions: model.Value(inventory.PositionRows{Rows: []inventory.PositionRequest{{
			ItemRequest:        positions.ItemRequest{CoreRequest: positions.CoreRequest{ID: model.Value("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")}, Quantity: model.Value(3.0)},
			CalculatedQuantity: model.Value(0.0),
		}}}),
	}
)

func TestInventoryFixture(t *testing.T) {
	data := read(t, "custom/inventory/read.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got inventory.Inventory
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Organization == nil || got.Store == nil || got.Owner == nil || got.Group == nil ||
				got.Name != "00004" || !got.Printed || got.Moment.IsZero() || got.Updated.IsZero() ||
				got.Positions == nil || got.Positions.Meta.Size != 1 || got.Positions.Meta.Type != "inventoryposition" {
				t.Fatalf("document = %+v", got)
			}
			var sparse inventory.Inventory
			if err := codec.Unmarshal([]byte(`{"id":"x","positions":{"meta":{"href":"h","size":2}},"enters":[{"meta":{"href":"e","type":"enter"}}],"losses":[],"future":1}`), &sparse); err != nil ||
				sparse.Positions.Rows != nil || sparse.Organization != nil || !sparse.Created.IsZero() ||
				len(sparse.Enters) != 1 || sparse.Enters[0].Meta.Type != "enter" || sparse.Losses == nil || len(sparse.Losses) != 0 {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

// Строки примера «Получить позиции»: calculatedQuantity, correctionAmount и correctionSum читаются.
func TestInventoryPositionsFixture(t *testing.T) {
	data := read(t, "custom/inventory/positions.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got model.Collection[inventory.Position]
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if len(got.Rows) != 2 || got.Meta.Type != "inventoryposition" {
				t.Fatalf("page = %+v", got)
			}
			row := got.Rows[0]
			if row.ID != "201452a2-ad6f-11e6-8a84-bc5200000010" || row.Meta.Type != "inventoryposition" || row.Assortment.Meta.Type != "product" ||
				row.Quantity != 140 || row.Price != 5000 || row.CalculatedQuantity != 200 || row.CorrectionAmount != -60 || row.CorrectionSum != -300000 || row.Pack != nil {
				t.Fatalf("row = %+v", row)
			}
			var packed inventory.Position
			if err := codec.Unmarshal([]byte(`{"pack":{"id":"p","quantity":6},"calculatedQuantity":1.5}`), &packed); err != nil ||
				packed.Pack == nil || packed.Pack.Quantity != 6 || packed.CalculatedQuantity != 1.5 {
				t.Fatalf("packed = %+v, err = %v", packed, err)
			}
		})
	}
}

func TestInventoryRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", fullCreate, read(t, "custom/inventory/create-full.golden.json")},
			{"create minimal", inventory.CreateRequest{Organization: organization, Store: store}, []byte(`{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},"store":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/store/33333333-3333-3333-3333-333333333333","type":"store"}}}`)},
			{"update empty positions", inventory.UpdateRequest{Positions: model.Value(inventory.PositionRows{Rows: []inventory.PositionRequest{}})}, []byte(`{"positions":{"rows":[]}}`)},
			{"update", fullUpdate, read(t, "custom/inventory/update.golden.json")},
			{"position", inventory.PositionRequest{ItemRequest: positions.ItemRequest{Quantity: model.Value(2.0)}, Price: model.Value(0.0), CalculatedQuantity: model.Null[float64](), Pack: model.Value(positions.PackRequest{ID: "p"})},
				[]byte(`{"quantity":2,"price":0,"calculatedQuantity":null,"pack":{"id":"p"}}`)},
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
	// Readonly по таблице и поля, которых у инвентаризации нет.
	for _, typ := range []reflect.Type{reflect.TypeFor[inventory.CreateRequest](), reflect.TypeFor[inventory.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Sum",
			"Files", "Applicable", "Rate", "Project", "Agent"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only or foreign %s", typ.Name(), name)
			}
		}
	}
	for _, name := range []string{"CorrectionAmount", "CorrectionSum", "Meta", "AccountID"} {
		if _, ok := reflect.TypeFor[inventory.PositionRequest]().FieldByName(name); ok {
			t.Errorf("PositionRequest contains read-only %s", name)
		}
	}
}

func TestInventoryHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и сумма — только из ответа.
	const response = `{"id":"server-id","name":"server-name","sum":24691,"positions":{"meta":{"href":"https://example.test/positions","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*inventory.Inventory, error)
	}{
		{"get", "GET", "entity/inventory/" + testID, nil, func(ctx context.Context, c *transport.Client) (*inventory.Inventory, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/inventory", read(t, "custom/inventory/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*inventory.Inventory, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", "entity/inventory/" + testID, read(t, "custom/inventory/update.golden.json"), func(ctx context.Context, c *transport.Client) (*inventory.Inventory, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *inventory.Inventory) {
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

func TestInventoryPositionHTTP(t *testing.T) {
	const row = `{"id":"server-row","quantity":2,"price":100,"calculatedQuantity":5,"correctionAmount":-3,"correctionSum":-300}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*inventory.Position, error)
	}{
		{"get", "GET", "entity/inventory/" + testID + "/positions/" + rowID, nil, func(ctx context.Context, c *transport.Client) (*inventory.Position, error) {
			return service(c).GetPosition(ctx, testID, rowID)
		}},
		{"update", "PUT", "entity/inventory/" + testID + "/positions/" + rowID, []byte(`{"quantity":2,"calculatedQuantity":5}`), func(ctx context.Context, c *transport.Client) (*inventory.Position, error) {
			return service(c).UpdatePosition(ctx, testID, rowID, inventory.PositionRequest{ItemRequest: positions.ItemRequest{Quantity: model.Value(2.0)}, CalculatedQuantity: model.Value(5.0)})
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(row) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *inventory.Position) {
						if got.ID != "server-row" || got.Quantity != 2 || got.Price != 100 || got.CalculatedQuantity != 5 || got.CorrectionAmount != -3 || got.CorrectionSum != -300 {
							t.Fatalf("row = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
	for _, page := range []positions.PageOptions{{}, {Limit: 1, Offset: 7}} {
		query := "limit=1000&offset=0"
		if page.Limit == 1 {
			query = "limit=1&offset=7"
		}
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[inventory.Position], error) {
			return service(c).GetPositions(ctx, testID, page)
		}
		for _, tc := range documenttest.ResponseCases(string(read(t, "custom/inventory/positions.json"))) {
			t.Run("page/"+tc.Name+"?"+query, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", "entity/inventory/"+testID+"/positions?"+query, nil, invoke, func(got *model.Collection[inventory.Position]) {
					if len(got.Rows) != 2 || got.Rows[1].CalculatedQuantity != 100 || got.Rows[1].CorrectionSum != -80000 {
						t.Fatalf("page = %+v", got)
					}
				})
			})
		}
		documenttest.LocalErrors(t, invoke)
	}
}

// writeOp — метод без модели документа в ответе: путь, тело запроса, успешный ответ и имя
// операции в тексте ошибки.
type writeOp struct {
	name, method, path, body, response, context string
	call                                        func(context.Context, *transport.Client) error
}

func writeOps() []writeOp {
	base := "/entity/inventory/" + testID + "/positions"
	return []writeOp{
		{"delete", "DELETE", "/entity/inventory/" + testID, "", "", "delete inventory", func(ctx context.Context, c *transport.Client) error {
			return service(c).Delete(ctx, testID)
		}},
		{"trash", "POST", "/entity/inventory/" + testID + "/trash", "", "", "trash inventory", func(ctx context.Context, c *transport.Client) error {
			return service(c).Trash(ctx, testID)
		}},
		{"new", "PUT", "/entity/inventory/new", "", `{"name":"00001","printed":true}`, "get inventory template", func(ctx context.Context, c *transport.Client) error {
			got, err := service(c).New(ctx, nil)
			if err == nil && (got.Name != "00001" || !got.Printed) {
				return errors.New("unexpected template")
			}
			return err
		}},
		{"new with basis", "PUT", "/entity/inventory/new", `{"customerOrder":{"meta":{"href":"h","type":"customerorder"}}}`, `{"name":"00001"}`, "get inventory template", func(ctx context.Context, c *transport.Client) error {
			_, err := service(c).New(ctx, map[string]any{"customerOrder": model.Reference{Meta: model.Meta{Href: "h", Type: "customerorder"}}})
			return err
		}},
		{"create many", "POST", "/entity/inventory", `[{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},"store":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/store/33333333-3333-3333-3333-333333333333","type":"store"}}}]`, `[{"id":"a"}]`, "create many inventory", func(ctx context.Context, c *transport.Client) error {
			_, err := service(c).CreateMany(ctx, []inventory.CreateRequest{{Organization: organization, Store: store}})
			return err
		}},
		{"update many", "POST", "/entity/inventory", `[{"meta":{"href":"{base}/entity/inventory/` + testID + `","metadataHref":"{base}/entity/inventory/metadata","type":"inventory","mediaType":"application/json"},"losses":[]}]`, `[{"id":"a"}]`, "update many inventory", func(ctx context.Context, c *transport.Client) error {
			_, err := service(c).UpdateMany(ctx, []batch.Update[inventory.UpdateRequest]{{ID: testID, Request: inventory.UpdateRequest{Losses: model.Value([]model.Reference{})}}})
			return err
		}},
		{"delete many", "POST", "/entity/inventory/delete", `[{"meta":{"href":"{base}/entity/inventory/` + testID + `","metadataHref":"{base}/entity/inventory/metadata","type":"inventory","mediaType":"application/json"}}]`, `[{"info":"Deleted"}]`, "delete many inventory", func(ctx context.Context, c *transport.Client) error {
			_, err := service(c).DeleteMany(ctx, []string{testID})
			return err
		}},
		{"create positions", "POST", base, `[{"assortment":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/88888888-8888-8888-8888-888888888888","type":"product"}},"quantity":1,"calculatedQuantity":2}]`, `[{"id":"` + rowID + `","calculatedQuantity":2}]`, "create many inventory " + testID + " positions", func(ctx context.Context, c *transport.Client) error {
			got, err := service(c).CreatePositions(ctx, testID, []inventory.PositionRequest{{ItemRequest: positions.ItemRequest{Assortment: model.Value(ref("entity/product/88888888-8888-8888-8888-888888888888", "product")), Quantity: model.Value(1.0)}, CalculatedQuantity: model.Value(2.0)}})
			if err == nil && (len(got) != 1 || got[0].Value == nil || got[0].Value.ID != rowID || got[0].Value.CalculatedQuantity != 2) {
				return errors.New("unexpected positions")
			}
			return err
		}},
		{"delete position", "DELETE", base + "/" + rowID, "", "", "delete inventory " + testID + " position", func(ctx context.Context, c *transport.Client) error {
			return service(c).DeletePosition(ctx, testID, rowID)
		}},
		{"delete many positions", "POST", base + "/delete", `[{"meta":{"href":"{base}` + base + "/" + rowID + `","type":"inventoryposition","mediaType":"application/json"}}]`, "", "delete inventory " + testID + " positions", func(ctx context.Context, c *transport.Client) error {
			return service(c).DeleteManyPositions(ctx, testID, []string{rowID})
		}},
	}
}

// Каждый метод: один запрос с точным методом, путём и телом; 200 — успех; ошибка API и
// неожиданный статус — ошибка с контекстом операции, без повторов.
func TestInventoryWrites(t *testing.T) {
	for _, op := range writeOps() {
		t.Run(op.name, func(t *testing.T) {
			for _, c := range []struct {
				status int
				body   string
				ok     bool
			}{{200, op.response, true}, {400, `{"errors":[{"error":"bad","code":3000}]}`, false}, {202, op.response, false}} {
				var calls atomic.Int32
				var server *httptest.Server
				server = memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					b, _ := io.ReadAll(r.Body)
					want := strings.ReplaceAll(op.body, "{base}", server.URL)
					if r.Method != op.method || r.URL.Path != op.path {
						t.Errorf("request = %s %s", r.Method, r.URL.Path)
					}
					if want == "" && len(b) != 0 {
						t.Errorf("body = %s", b)
					} else if want != "" {
						documenttest.JSONEqual(t, b, []byte(want))
					}
					w.WriteHeader(c.status)
					io.WriteString(w, c.body)
				}))
				client, err := transport.NewClient(transport.Config{BaseURL: server.URL, Retry: &transport.RetryPolicy{MaxAttempts: 3}})
				if err != nil {
					t.Fatal(err)
				}
				err = op.call(t.Context(), client)
				server.Close()
				if (err == nil) != c.ok || calls.Load() != 1 {
					t.Fatalf("status %d: err = %v, calls = %d", c.status, err, calls.Load())
				}
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), op.context)) {
					t.Errorf("error = %v", err)
				}
			}
			documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*struct{}, error) {
				if err := op.call(ctx, c); err != nil {
					return nil, err
				}
				return &struct{}{}, nil
			})
		})
	}
}

// Неверный вход отклоняется до HTTP: ResponseInfo сброшен, запросов нет.
func TestInventoryValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(c)
	for name, call := range map[string]func(transport.RequestOption) error{
		"get id": func(o transport.RequestOption) error { _, err := s.Get(t.Context(), "bad", o); return err },
		"no organization": func(o transport.RequestOption) error {
			_, err := s.Create(t.Context(), inventory.CreateRequest{Store: store}, o)
			return err
		},
		"no store": func(o transport.RequestOption) error {
			_, err := s.Create(t.Context(), inventory.CreateRequest{Organization: organization}, o)
			return err
		},
		"create 1001": func(o transport.RequestOption) error {
			_, err := s.Create(t.Context(), inventory.CreateRequest{Organization: organization, Store: store, Positions: model.Value(inventory.PositionRows{Rows: make([]inventory.PositionRequest, 1001)})}, o)
			return err
		},
		"update 1001": func(o transport.RequestOption) error {
			_, err := s.Update(t.Context(), testID, inventory.UpdateRequest{Positions: model.Value(inventory.PositionRows{Rows: make([]inventory.PositionRequest, 1001)})}, o)
			return err
		},
		"update id": func(o transport.RequestOption) error {
			_, err := s.Update(t.Context(), "bad", fullUpdate, o)
			return err
		},
		"empty update": func(o transport.RequestOption) error {
			_, err := s.Update(t.Context(), testID, inventory.UpdateRequest{}, o)
			return err
		},
		"delete id": func(o transport.RequestOption) error { return s.Delete(t.Context(), "bad", o) },
		"trash id":  func(o transport.RequestOption) error { return s.Trash(t.Context(), testID+"/x", o) },
		"create many empty": func(o transport.RequestOption) error {
			_, err := s.CreateMany(t.Context(), nil, o)
			return err
		},
		"update many id": func(o transport.RequestOption) error {
			_, err := s.UpdateMany(t.Context(), []batch.Update[inventory.UpdateRequest]{{ID: "bad", Request: fullUpdate}}, o)
			return err
		},
		"update many empty request": func(o transport.RequestOption) error {
			_, err := s.UpdateMany(t.Context(), []batch.Update[inventory.UpdateRequest]{{ID: testID}}, o)
			return err
		},
		"delete many id": func(o transport.RequestOption) error {
			_, err := s.DeleteMany(t.Context(), []string{"bad"}, o)
			return err
		},
		"positions id": func(o transport.RequestOption) error {
			_, err := s.GetPositions(t.Context(), "bad", positions.PageOptions{}, o)
			return err
		},
		"positions fields": func(o transport.RequestOption) error {
			_, err := s.GetPositions(t.Context(), testID, positions.PageOptions{Fields: "stock"}, o)
			return err
		},
		"positions limit": func(o transport.RequestOption) error {
			_, err := s.GetPositions(t.Context(), testID, positions.PageOptions{Limit: 1001}, o)
			return err
		},
		"position id": func(o transport.RequestOption) error {
			_, err := s.GetPosition(t.Context(), testID, "bad", o)
			return err
		},
		"position doc id": func(o transport.RequestOption) error {
			_, err := s.GetPosition(t.Context(), "bad", rowID, o)
			return err
		},
		"create positions": func(o transport.RequestOption) error {
			_, err := s.CreatePositions(t.Context(), "bad", []inventory.PositionRequest{{}}, o)
			return err
		},
		"update position empty": func(o transport.RequestOption) error {
			_, err := s.UpdatePosition(t.Context(), testID, rowID, inventory.PositionRequest{}, o)
			return err
		},
		"delete position":        func(o transport.RequestOption) error { return s.DeletePosition(t.Context(), testID, "bad", o) },
		"delete positions empty": func(o transport.RequestOption) error { return s.DeleteManyPositions(t.Context(), testID, nil, o) },
		"delete positions id": func(o transport.RequestOption) error {
			return s.DeleteManyPositions(t.Context(), testID, []string{"bad"}, o)
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		err := call(transport.WithResponseInfo(&info))
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "inventory") || info.StatusCode != 0 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := inventory.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestInventoryList(t *testing.T) {
	listtest.Run(t, "inventory", read(t, "custom/inventory/read.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[inventory.Inventory] {
		return service(c).List
	}, func(t *testing.T, got inventory.Inventory) {
		if got.ID != fixtureID || got.Positions == nil || got.Name != "00004" {
			t.Fatalf("row = %+v", got)
		}
	})
}
