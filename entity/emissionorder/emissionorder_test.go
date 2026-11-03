package emissionorder_test

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
	"time"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/emissionorder"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID     = "12121212-1212-1212-1212-121212121212"
	fixtureID  = "33c8ac0b-a670-11ef-ac12-000d00000111"
	positionID = "34b7de2d-a670-11ef-ac12-000d00000112"
	docPath    = "entity/emissionorder/" + testID
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

func service(c *transport.Client) *emissionorder.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := emissionorder.NewService(c)
	return s
}

var (
	organization = ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization")
	product      = ref("entity/product/77777777-7777-7777-7777-777777777777", "product")
	minimal      = emissionorder.CreateRequest{Organization: organization, TrackingType: "MILK", EmissionType: "LOCAL"}
	fullCreate   = emissionorder.CreateRequest{
		Organization: organization, TrackingType: "MILK", EmissionType: "LOCAL",
		Name: model.Value("2000"), Description: model.Value(""), ExternalCode: model.Value("ext"),
		Moment: model.Value(documenttest.DateTime("2026-03-25 19:03:00")), Shared: model.Value(false),
		Owner:           model.Value(ref("entity/employee/55555555-5555-5555-5555-555555555555", "employee")),
		Group:           model.Value(ref("entity/group/66666666-6666-6666-6666-666666666666", "group")),
		State:           model.Null[model.Reference](),
		Positions:       model.Value([]emissionorder.PositionRequest{{Assortment: model.Value(product), Quantity: model.Value(3.0)}}),
		ProductionTasks: model.Value([]model.Reference{ref("entity/productiontask/88888888-8888-8888-8888-888888888888", "productiontask")}),
	}
	fullUpdate = emissionorder.UpdateRequest{
		EmissionType: model.Value("REMAINS"),
		Name:         model.Value("00001-обновленный"),
		Group:        model.Null[model.Reference](),
		Positions: model.Value([]emissionorder.PositionRequest{
			{CoreRequest: positions.CoreRequest{ID: model.Value(positionID)}, Quantity: model.Value(44.0)},
			{Assortment: model.Value(product), Quantity: model.Value(1.5)},
		}),
		ProductionTasks: model.Value([]model.Reference{}),
	}
)

func TestFixture(t *testing.T) {
	data := read(t, "upstream/emission_order.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got emissionorder.EmissionOrder
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.AccountID != "b9c9f2e4-a66f-11ef-ac12-000f00000001" || got.Name != "00001" ||
				got.ExternalCode != "RIahRnZIjek-CeER27IjF0" || got.Meta.Type != "emissionorder" || got.Shared || got.Printed || got.Published {
				t.Fatalf("document = %+v", got)
			}
			if got.EmissionType != "LOCAL" || got.TrackingType != "MILK" || got.DocumentState != "CREATED" {
				t.Fatalf("enums = %q %q %q", got.EmissionType, got.TrackingType, got.DocumentState)
			}
			if !time.Time(got.Moment).Equal(time.Time(documenttest.DateTime("2024-11-19 15:16:00"))) || got.Created.IsZero() || got.Updated.IsZero() {
				t.Fatalf("dates = %+v", got)
			}
			if got.Organization == nil || got.Organization.Meta.Type != "organization" || got.Owner == nil || got.Group == nil ||
				got.State == nil || got.State.Meta.Type != "state" {
				t.Fatalf("references = %+v", got)
			}
			if len(got.ProductionTasks) != 1 || got.ProductionTasks[0].Meta.Type != "productiontask" {
				t.Fatalf("productionTasks = %+v", got.ProductionTasks)
			}
			if got.Positions == nil || got.Positions.Meta.Type != "emissionorderposition" || got.Positions.Meta.Size != 1 || len(got.Positions.Rows) != 1 {
				t.Fatalf("positions = %+v", got.Positions)
			}
			row := got.Positions.Rows[0]
			if row.ID != positionID || row.Meta.Type != "emissionorderposition" || row.AccountID != got.AccountID || row.Quantity != 1 ||
				row.Status != "EMISSION_NOT_SEND" || row.Assortment.Meta.Type != "product" {
				t.Fatalf("row = %+v", row)
			}
			var sparse emissionorder.EmissionOrder
			if err := codec.Unmarshal([]byte(`{"id":"x","description":"d","trackingType":"FUTURE_TYPE","emissionType":"FUTURE","documentState":"NEW_STATE",`+
				`"positions":{"meta":{"href":"h"},"rows":[{"id":"p","quantity":1000,"status":"EMISSION_FUTURE"}]},"future":1}`), &sparse); err != nil {
				t.Fatal(err)
			}
			if sparse.Description != "d" || sparse.TrackingType != "FUTURE_TYPE" || sparse.EmissionType != "FUTURE" || sparse.DocumentState != "NEW_STATE" ||
				sparse.Positions.Rows[0].Status != "EMISSION_FUTURE" || sparse.Positions.Rows[0].Quantity != 1000 ||
				sparse.Organization != nil || !sparse.Created.IsZero() || sparse.ProductionTasks != nil {
				t.Fatalf("sparse = %+v", sparse)
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
			{"create full", fullCreate, read(t, "custom/emissionorder/create-full.golden.json")},
			{"create minimal", minimal, []byte(`{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},` +
				`"trackingType":"MILK","emissionType":"LOCAL"}`)},
			{"update", fullUpdate, read(t, "custom/emissionorder/update.golden.json")},
			{"update empty positions", emissionorder.UpdateRequest{Positions: model.Value([]emissionorder.PositionRequest{})}, []byte(`{"positions":[]}`)},
			{"position update", emissionorder.PositionRequest{Quantity: model.Value(44.0)}, []byte(`{"quantity":44}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[emissionorder.CreateRequest](), reflect.TypeFor[emissionorder.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Printed", "Published", "DocumentState"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
	for _, name := range []string{"Meta", "AccountID", "Status"} {
		if _, ok := reflect.TypeFor[emissionorder.PositionRequest]().FieldByName(name); ok {
			t.Errorf("PositionRequest contains read-only %s", name)
		}
	}
}

func TestHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и состояние — только из ответа.
	const response = `{"id":"server-id","name":"server-name","documentState":"SUZ_SEND","positions":{"meta":{"href":"https://example.test/positions","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*emissionorder.EmissionOrder, error)
	}{
		{"get", "GET", docPath, nil, func(ctx context.Context, c *transport.Client) (*emissionorder.EmissionOrder, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/emissionorder", read(t, "custom/emissionorder/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*emissionorder.EmissionOrder, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", docPath, read(t, "custom/emissionorder/update.golden.json"), func(ctx context.Context, c *transport.Client) (*emissionorder.EmissionOrder, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *emissionorder.EmissionOrder) {
						if got.ID != "server-id" || got.Name != "server-name" || got.DocumentState != "SUZ_SEND" || got.Positions.Rows != nil {
							t.Fatalf("document = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

func httpCases[T any](t *testing.T, method, path, body, response string, invoke func(context.Context, *transport.Client) (*T, error), check func(*T)) {
	t.Helper()
	var sent []byte
	if body != "" {
		sent = []byte(body)
	}
	for _, tc := range documenttest.ResponseCases(response) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, method, path, sent, invoke, check)
		})
	}
	documenttest.LocalErrors(t, invoke)
}

func TestPositionsHTTP(t *testing.T) {
	const row = `{"meta":{"type":"emissionorderposition"},"id":"` + positionID + `","quantity":44,"assortment":{"meta":{"type":"variant"}},"status":"EMISSION_NOT_SEND"}`
	checkRow := func(got *emissionorder.Position) {
		if got.ID != positionID || got.Quantity != 44 || got.Status != "EMISSION_NOT_SEND" || got.Assortment.Meta.Type != "variant" {
			t.Fatalf("row = %+v", got)
		}
	}
	t.Run("get positions", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/positions?limit=2&offset=7", "", `{"meta":{"href":"h"},"rows":[`+row+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[emissionorder.Position], error) {
				return service(c).GetPositions(ctx, testID, positions.PageOptions{Limit: 2, Offset: 7})
			}, func(got *model.Collection[emissionorder.Position]) { checkRow(&got.Rows[0]) })
	})
	t.Run("get position", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/positions/"+positionID, "", row, func(ctx context.Context, c *transport.Client) (*emissionorder.Position, error) {
			return service(c).GetPosition(ctx, testID, positionID)
		}, checkRow)
	})
	t.Run("update position", func(t *testing.T) {
		httpCases(t, "PUT", docPath+"/positions/"+positionID, `{"quantity":44}`, row, func(ctx context.Context, c *transport.Client) (*emissionorder.Position, error) {
			return service(c).UpdatePosition(ctx, testID, positionID, emissionorder.PositionRequest{Quantity: model.Value(44.0)})
		}, checkRow)
	})
}

// Методы без модели в ответе и массовые: один запрос с точным телом, без повторов при 503.
func TestWrites(t *testing.T) {
	const base = "{base}/entity/emissionorder"
	docMeta := `{"href":"` + base + `/` + testID + `","metadataHref":"` + base + `/metadata","type":"emissionorder","mediaType":"application/json"}`
	productJSON := `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/77777777-7777-7777-7777-777777777777","type":"product"}}`
	for _, tc := range []struct {
		name, method, path, body, response string
		empty                              bool // ответ без тела: 204 тоже успех
		invoke                             func(context.Context, *transport.Client) error
	}{
		{"trash", "POST", docPath + "/trash", "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Trash(ctx, testID) }},
		{"create many", "POST", "entity/emissionorder", `[{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},` +
			`"trackingType":"MILK","emissionType":"LOCAL"}]`, `[{"id":"x"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateMany(ctx, []emissionorder.CreateRequest{minimal})
				if err == nil && r[0].Value.ID != "x" {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"update many", "POST", "entity/emissionorder", `[{"meta":` + docMeta + `,"trackingType":"SHOES"}]`, `[{"id":"x"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				_, err := service(c).UpdateMany(ctx, []batch.Update[emissionorder.UpdateRequest]{{ID: testID, Request: emissionorder.UpdateRequest{TrackingType: model.Value("SHOES")}}})
				return err
			}},
		{"create positions", "POST", docPath + "/positions", `[{"assortment":` + productJSON + `,"quantity":49}]`, `[{"id":"` + positionID + `","status":"EMISSION_NOT_SEND"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreatePositions(ctx, testID, []emissionorder.PositionRequest{{Assortment: model.Value(product), Quantity: model.Value(49.0)}})
				if err == nil && (r[0].Value.ID != positionID || r[0].Value.Status != "EMISSION_NOT_SEND") {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, st := range []struct {
				status int
				ok     bool
			}{{200, true}, {204, tc.empty}, {201, false}, {403, false}, {503, false}} {
				var calls atomic.Int32
				var server *httptest.Server
				server = memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					b, _ := io.ReadAll(r.Body)
					if r.Method != tc.method || r.URL.Path != "/api/remap/1.2/"+tc.path {
						t.Errorf("request = %s %s", r.Method, r.URL)
					}
					if tc.body == "" {
						if len(b) != 0 {
							t.Errorf("body = %s", b)
						}
					} else {
						documenttest.JSONEqual(t, b, []byte(strings.ReplaceAll(tc.body, "{base}", server.URL+"/api/remap/1.2")))
					}
					w.WriteHeader(st.status)
					switch st.status {
					case 200:
						io.WriteString(w, tc.response)
					case 403:
						io.WriteString(w, `{"errors":[{"error":"forbidden","code":1016}]}`)
					}
				}))
				client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", Retry: &transport.RetryPolicy{MaxAttempts: 3}})
				if err != nil {
					t.Fatal(err)
				}
				err = tc.invoke(t.Context(), client)
				server.Close()
				if (err == nil) != st.ok || calls.Load() != 1 {
					t.Fatalf("status %d: err = %v, calls = %d", st.status, err, calls.Load())
				}
				var apiErr *sdkerrors.APIError
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "emissionorder") ||
					st.status == 403 && !sdkerrors.HasCode(err, 1016) || errors.As(err, &apiErr) != (st.status >= 400)) {
					t.Errorf("status %d: error = %v", st.status, err)
				}
			}
			documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*struct{}, error) {
				if err := tc.invoke(ctx, c); err != nil {
					return nil, err
				}
				return &struct{}{}, nil
			})
		})
	}
}

// Источник не описывает удаление документа и строк: таких методов нет.
func TestNoDelete(t *testing.T) {
	typ := reflect.TypeFor[*emissionorder.Service]()
	for _, name := range []string{"Delete", "DeleteMany", "DeletePosition", "DeleteManyPositions"} {
		if _, ok := typ.MethodByName(name); ok {
			t.Errorf("Service has %s", name)
		}
	}
}

// Ошибки входа — до HTTP, ResponseInfo сброшен; ровно 1000 строк уходят.
func TestValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		io.WriteString(w, `{"id":"server-id"}`)
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(client), t.Context()
	rows := func(n int) model.Field[[]emissionorder.PositionRequest] {
		return model.Value(make([]emissionorder.PositionRequest, n))
	}
	limit := minimal
	limit.Positions = rows(1000)
	if _, err := s.Create(ctx, limit); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(ctx, testID, emissionorder.UpdateRequest{Positions: rows(1000)}); err != nil {
		t.Fatal(err)
	}
	// Открытый enum: неизвестные значения уходят на сервер.
	if _, err := s.Create(ctx, emissionorder.CreateRequest{Organization: organization, TrackingType: "FUTURE", EmissionType: "FUTURE"}); err != nil {
		t.Fatal(err)
	}
	create := func(change func(*emissionorder.CreateRequest)) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error {
			input := minimal
			change(&input)
			_, err := s.Create(ctx, input, o...)
			return err
		}
	}
	getPositions := func(id string, p positions.PageOptions) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error {
			_, err := s.GetPositions(ctx, id, p, o...)
			return err
		}
	}
	calls.Store(0)
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id":             func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":               func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no org":          create(func(r *emissionorder.CreateRequest) { r.Organization = model.Reference{} }),
		"create no trackingType": create(func(r *emissionorder.CreateRequest) { r.TrackingType = "" }),
		"create no emissionType": create(func(r *emissionorder.CreateRequest) { r.EmissionType = "" }),
		"create 1001":            create(func(r *emissionorder.CreateRequest) { r.Positions = rows(1001) }),
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, emissionorder.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", fullUpdate, o...); return err },
		"update 1001": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, emissionorder.UpdateRequest{Positions: rows(1001)}, o...)
			return err
		},
		"trash bad id": func(o ...transport.RequestOption) error { return s.Trash(ctx, "bad", o...) },
		"create many empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateMany(ctx, nil, o...)
			return err
		},
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[emissionorder.UpdateRequest]{{ID: "bad", Request: fullUpdate}}, o...)
			return err
		},
		"update many empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[emissionorder.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"positions bad id": getPositions("bad", positions.PageOptions{}),
		"positions fields": getPositions(testID, positions.PageOptions{Fields: "stock"}),
		"positions limit":  getPositions(testID, positions.PageOptions{Limit: 1001}),
		"get position bad": func(o ...transport.RequestOption) error {
			_, err := s.GetPosition(ctx, testID, "bad", o...)
			return err
		},
		"create positions empty": func(o ...transport.RequestOption) error {
			_, err := s.CreatePositions(ctx, testID, nil, o...)
			return err
		},
		"create positions bad id": func(o ...transport.RequestOption) error {
			_, err := s.CreatePositions(ctx, "bad", []emissionorder.PositionRequest{{}}, o...)
			return err
		},
		"update position empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdatePosition(ctx, testID, positionID, emissionorder.PositionRequest{}, o...)
			return err
		},
		"update position bad": func(o ...transport.RequestOption) error {
			_, err := s.UpdatePosition(ctx, "bad", positionID, emissionorder.PositionRequest{Quantity: model.Value(1.0)}, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "emissionorder") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := emissionorder.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestList(t *testing.T) {
	listtest.Run(t, "emissionorder", read(t, "upstream/emission_order.json"), "", "", func(c *transport.Client) listtest.ListFunc[emissionorder.EmissionOrder] {
		return service(c).List
	}, func(t *testing.T, got emissionorder.EmissionOrder) {
		if got.ID != fixtureID || got.TrackingType != "MILK" || len(got.Positions.Rows) != 1 {
			t.Fatalf("row = %+v", got)
		}
	})
}
