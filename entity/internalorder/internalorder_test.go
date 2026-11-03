package internalorder_test

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
	"github.com/eteeet/moysklad-sdk-go/entity/internalorder"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	rowID     = "34f6344f-015e-11e6-9464-e4de0000006c"
	fixtureID = "31d58bde-b2fd-11e6-8a84-bae500000068"
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

func service(c *transport.Client) *internalorder.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := internalorder.NewService(c)
	return s
}

var (
	organization = ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization")
	fullCreate   = internalorder.CreateRequest{
		Organization: organization, Name: "DOC-1", Description: model.Value(""),
		Moment:                model.Value(documenttest.DateTime("2026-11-01 12:34:56.000")),
		DeliveryPlannedMoment: model.Value(documenttest.DateTime("2026-11-05 10:00:00.000")),
		Applicable:            model.Value(false),
		Store:                 model.Value(ref("entity/store/33333333-3333-3333-3333-333333333333", "store")),
		Attributes:            model.Value([]model.Attribute{}),
		Rate:                  model.Value(model.DocumentRateRequest{Currency: model.Value(ref("entity/currency/66666666-6666-6666-6666-666666666666", "currency"))}),
		VATEnabled:            model.Value(true), VATIncluded: model.Value(false),
		Moves:           model.Value([]model.Reference{ref("entity/move/55555555-5555-5555-5555-555555555555", "move")}),
		PurchaseOrders:  model.Value([]model.Reference{ref("entity/purchaseorder/77777777-7777-7777-7777-777777777777", "purchaseorder")}),
		ProductionTasks: model.Value([]model.Reference(nil)),
		Positions: model.Value([]internalorder.PositionRequest{{ItemRequest: positions.ItemRequest{Assortment: model.Value(ref("entity/variant/fc635ded-acf7-11e6-8a84-bae500000075", "variant")), Quantity: model.Value(1.0)},
			Price: model.Value(0.0), Discount: model.Value(0.0), VAT: model.Value(0)}}),
	}
	fullUpdate = internalorder.UpdateRequest{
		Name:       model.Value("DOC-2"),
		Store:      model.Value(ref("entity/store/99999999-9999-9999-9999-999999999999", "store")),
		Project:    model.Null[model.Reference](),
		Applicable: model.Value(true),
		Moves:      model.Value([]model.Reference{}),
		Positions:  model.Value([]internalorder.PositionRequest{{ItemRequest: positions.ItemRequest{CoreRequest: positions.CoreRequest{ID: model.Value(rowID)}, Quantity: model.Value(3.0)}}}),
	}
)

func TestInternalOrderFixture(t *testing.T) {
	data := read(t, "custom/internalorder/read.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got internalorder.InternalOrder
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Organization == nil || got.Store == nil || got.Project == nil || got.State == nil ||
				got.Rate == nil || got.Positions == nil || got.Positions.Meta.Size != 2 || got.Positions.Meta.Type != "internalorderposition" ||
				!got.VATEnabled || !got.VATIncluded || got.DeliveryPlannedMoment.IsZero() || got.Updated.IsZero() ||
				len(got.Moves) != 2 || got.Moves[0].Meta.Type != "move" || len(got.PurchaseOrders) != 2 || got.PurchaseOrders[1].Meta.Type != "purchaseorder" {
				t.Fatalf("document = %+v", got)
			}
			var sparse internalorder.InternalOrder
			if err := codec.Unmarshal([]byte(`{"id":"x","vatSum":12.5,"positions":{"meta":{"href":"h","size":2}},"productionTasks":[{"meta":{"href":"p","type":"productiontask"}}],"future":1}`), &sparse); err != nil ||
				sparse.Positions.Rows != nil || sparse.Organization != nil || !sparse.Created.IsZero() || sparse.Moves != nil ||
				sparse.VATSum != 12.5 || len(sparse.ProductionTasks) != 1 {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestInternalOrderPositionsFixture(t *testing.T) {
	data := read(t, "custom/internalorder/positions.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got model.Collection[internalorder.Position]
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if len(got.Rows) != 1 || got.Rows[0].Price != 2230 || got.Rows[0].Quantity != 1 || got.Rows[0].VATEnabled ||
				got.Rows[0].Meta.Type != "internalorderposition" || got.Rows[0].Assortment.Meta.Type != "product" {
				t.Fatalf("page = %+v", got)
			}
		})
	}
}

func TestInternalOrderRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", fullCreate, read(t, "custom/internalorder/create-full.golden.json")},
			{"create minimal", internalorder.CreateRequest{Organization: organization, Name: "N"}, []byte(`{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},"name":"N"}`)},
			{"update", fullUpdate, read(t, "custom/internalorder/update.golden.json")},
			{"position", internalorder.PositionRequest{ItemRequest: positions.ItemRequest{Quantity: model.Value(2.0)}, Price: model.Value(0.0), VAT: model.Value(0), VATEnabled: model.Value(false), Pack: model.Value(positions.PackRequest{ID: "p"})},
				[]byte(`{"quantity":2,"price":0,"vat":0,"vatEnabled":false,"pack":{"id":"p"}}`)},
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
	// shared и syncId помечены в таблице «Только для чтения». Positions тоже readonly, но пишется:
	// исключение по примерам Create/Update утверждено пользователем.
	for _, typ := range []reflect.Type{reflect.TypeFor[internalorder.CreateRequest](), reflect.TypeFor[internalorder.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Sum", "VATSum",
			"Files", "Shared", "SyncID", "Agent"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only or foreign %s", typ.Name(), name)
			}
		}
	}
}

func TestInternalOrderHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и сумма — только из ответа.
	const response = `{"id":"server-id","name":"server-name","sum":24691,"positions":{"meta":{"href":"https://example.test/positions","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*internalorder.InternalOrder, error)
	}{
		{"get", "GET", "entity/internalorder/" + testID, nil, func(ctx context.Context, c *transport.Client) (*internalorder.InternalOrder, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/internalorder", read(t, "custom/internalorder/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*internalorder.InternalOrder, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", "entity/internalorder/" + testID, read(t, "custom/internalorder/update.golden.json"), func(ctx context.Context, c *transport.Client) (*internalorder.InternalOrder, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *internalorder.InternalOrder) {
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

func TestInternalOrderPositionHTTP(t *testing.T) {
	const row = `{"id":"server-row","quantity":2,"price":100,"vat":20,"vatEnabled":true}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*internalorder.Position, error)
	}{
		{"get", "GET", "entity/internalorder/" + testID + "/positions/" + rowID, nil, func(ctx context.Context, c *transport.Client) (*internalorder.Position, error) {
			return service(c).GetPosition(ctx, testID, rowID)
		}},
		{"update", "PUT", "entity/internalorder/" + testID + "/positions/" + rowID, []byte(`{"quantity":2,"vat":20}`), func(ctx context.Context, c *transport.Client) (*internalorder.Position, error) {
			return service(c).UpdatePosition(ctx, testID, rowID, internalorder.PositionRequest{ItemRequest: positions.ItemRequest{Quantity: model.Value(2.0)}, VAT: model.Value(20)})
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(row) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *internalorder.Position) {
						if got.ID != "server-row" || got.Quantity != 2 || got.Price != 100 || got.VAT != 20 || !got.VATEnabled {
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
		invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[internalorder.Position], error) {
			return service(c).GetPositions(ctx, testID, page)
		}
		for _, tc := range documenttest.ResponseCases(string(read(t, "custom/internalorder/positions.json"))) {
			t.Run("page/"+tc.Name+"?"+query, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", "entity/internalorder/"+testID+"/positions?"+query, nil, invoke, func(got *model.Collection[internalorder.Position]) {
					if len(got.Rows) != 1 || got.Rows[0].Price != 2230 {
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
	base := "/entity/internalorder/" + testID + "/positions"
	return []writeOp{
		{"delete", "DELETE", "/entity/internalorder/" + testID, "", "", "delete internalorder", func(ctx context.Context, c *transport.Client) error {
			return service(c).Delete(ctx, testID)
		}},
		{"trash", "POST", "/entity/internalorder/" + testID + "/trash", "", "", "trash internalorder", func(ctx context.Context, c *transport.Client) error {
			return service(c).Trash(ctx, testID)
		}},
		{"new", "PUT", "/entity/internalorder/new", "", `{"name":"00001","vatEnabled":true}`, "get internalorder template", func(ctx context.Context, c *transport.Client) error {
			got, err := service(c).New(ctx, nil)
			if err == nil && (got.Name != "00001" || !got.VATEnabled) {
				return errors.New("unexpected template")
			}
			return err
		}},
		{"new with basis", "PUT", "/entity/internalorder/new", `{"customerOrder":{"meta":{"href":"h","type":"customerorder"}}}`, `{"name":"00001"}`, "get internalorder template", func(ctx context.Context, c *transport.Client) error {
			_, err := service(c).New(ctx, map[string]any{"customerOrder": model.Reference{Meta: model.Meta{Href: "h", Type: "customerorder"}}})
			return err
		}},
		{"create many", "POST", "/entity/internalorder", `[{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},"name":"N"}]`, `[{"id":"a"}]`, "create many internalorder", func(ctx context.Context, c *transport.Client) error {
			_, err := service(c).CreateMany(ctx, []internalorder.CreateRequest{{Organization: organization, Name: "N"}})
			return err
		}},
		{"update many", "POST", "/entity/internalorder", `[{"meta":{"href":"{base}/entity/internalorder/` + testID + `","metadataHref":"{base}/entity/internalorder/metadata","type":"internalorder","mediaType":"application/json"},"moves":[]}]`, `[{"id":"a"}]`, "update many internalorder", func(ctx context.Context, c *transport.Client) error {
			_, err := service(c).UpdateMany(ctx, []batch.Update[internalorder.UpdateRequest]{{ID: testID, Request: internalorder.UpdateRequest{Moves: model.Value([]model.Reference{})}}})
			return err
		}},
		{"delete many", "POST", "/entity/internalorder/delete", `[{"meta":{"href":"{base}/entity/internalorder/` + testID + `","metadataHref":"{base}/entity/internalorder/metadata","type":"internalorder","mediaType":"application/json"}}]`, `[{"info":"Deleted"}]`, "delete many internalorder", func(ctx context.Context, c *transport.Client) error {
			_, err := service(c).DeleteMany(ctx, []string{testID})
			return err
		}},
		{"create positions", "POST", base, `[{"assortment":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/88888888-8888-8888-8888-888888888888","type":"product"}},"quantity":1}]`, `[{"id":"` + rowID + `"}]`, "create many internalorder " + testID + " positions", func(ctx context.Context, c *transport.Client) error {
			got, err := service(c).CreatePositions(ctx, testID, []internalorder.PositionRequest{{ItemRequest: positions.ItemRequest{Assortment: model.Value(ref("entity/product/88888888-8888-8888-8888-888888888888", "product")), Quantity: model.Value(1.0)}}})
			if err == nil && (len(got) != 1 || got[0].Value == nil || got[0].Value.ID != rowID) {
				return errors.New("unexpected positions")
			}
			return err
		}},
		{"delete position", "DELETE", base + "/" + rowID, "", "", "delete internalorder " + testID + " position", func(ctx context.Context, c *transport.Client) error {
			return service(c).DeletePosition(ctx, testID, rowID)
		}},
		{"delete many positions", "POST", base + "/delete", `[{"meta":{"href":"{base}` + base + "/" + rowID + `","type":"internalorderposition","mediaType":"application/json"}}]`, "", "delete internalorder " + testID + " positions", func(ctx context.Context, c *transport.Client) error {
			return service(c).DeleteManyPositions(ctx, testID, []string{rowID})
		}},
	}
}

// Каждый метод: один запрос с точным методом, путём и телом; 200 — успех; ошибка API и
// неожиданный статус — ошибка с контекстом операции, без повторов.
func TestInternalOrderWrites(t *testing.T) {
	for _, op := range writeOps() {
		t.Run(op.name, func(t *testing.T) {
			for _, c := range []struct {
				status int
				body   string
				ok     bool
			}{{200, op.response, true}, {400, `{"errors":[{"error":"bad","code":3000}]}`, false}, {202, op.response, false}} {
				var calls atomic.Int32
				var server *httptest.Server
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
func TestInternalOrderValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(c)
	for name, call := range map[string]func(transport.RequestOption) error{
		"get id": func(o transport.RequestOption) error { _, err := s.Get(t.Context(), "bad", o); return err },
		"no organization": func(o transport.RequestOption) error {
			_, err := s.Create(t.Context(), internalorder.CreateRequest{Name: "N"}, o)
			return err
		},
		"no name": func(o transport.RequestOption) error {
			_, err := s.Create(t.Context(), internalorder.CreateRequest{Organization: organization}, o)
			return err
		},
		"update id": func(o transport.RequestOption) error {
			_, err := s.Update(t.Context(), "bad", fullUpdate, o)
			return err
		},
		"empty update": func(o transport.RequestOption) error {
			_, err := s.Update(t.Context(), testID, internalorder.UpdateRequest{}, o)
			return err
		},
		"delete id": func(o transport.RequestOption) error { return s.Delete(t.Context(), "bad", o) },
		"trash id":  func(o transport.RequestOption) error { return s.Trash(t.Context(), testID+"/x", o) },
		"create many empty": func(o transport.RequestOption) error {
			_, err := s.CreateMany(t.Context(), nil, o)
			return err
		},
		"update many id": func(o transport.RequestOption) error {
			_, err := s.UpdateMany(t.Context(), []batch.Update[internalorder.UpdateRequest]{{ID: "bad", Request: fullUpdate}}, o)
			return err
		},
		"update many empty request": func(o transport.RequestOption) error {
			_, err := s.UpdateMany(t.Context(), []batch.Update[internalorder.UpdateRequest]{{ID: testID}}, o)
			return err
		},
		"delete many id": func(o transport.RequestOption) error {
			_, err := s.DeleteMany(t.Context(), []string{"bad"}, o)
			return err
		},
		"create 1001 positions": func(o transport.RequestOption) error {
			_, err := s.Create(t.Context(), internalorder.CreateRequest{Organization: organization, Name: "N", Positions: model.Value(make([]internalorder.PositionRequest, 1001))}, o)
			return err
		},
		"update 1001 positions": func(o transport.RequestOption) error {
			_, err := s.Update(t.Context(), testID, internalorder.UpdateRequest{Positions: model.Value(make([]internalorder.PositionRequest, 1001))}, o)
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
			_, err := s.CreatePositions(t.Context(), "bad", []internalorder.PositionRequest{{}}, o)
			return err
		},
		"update position empty": func(o transport.RequestOption) error {
			_, err := s.UpdatePosition(t.Context(), testID, rowID, internalorder.PositionRequest{}, o)
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
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "internalorder") || info.StatusCode != 0 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := internalorder.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestInternalOrderList(t *testing.T) {
	listtest.Run(t, "internalorder", read(t, "custom/internalorder/read.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[internalorder.InternalOrder] {
		return service(c).List
	}, func(t *testing.T, got internalorder.InternalOrder) {
		if got.ID != fixtureID || got.Positions == nil || len(got.Moves) != 2 {
			t.Fatalf("row = %+v", got)
		}
	})
}
