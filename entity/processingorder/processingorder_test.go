package processingorder_test

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
	"github.com/eteeet/moysklad-sdk-go/entity/processingorder"
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
	positionID = "34f6344f-015e-11e6-9464-e4de0000006c"
	docPath    = "entity/processingorder/" + testID
	apiBase    = "https://api.moysklad.ru/api/remap/1.2/"
)

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/custom/processingorder/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ref(path, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: apiBase + path, Type: model.MetaType(typ)}}
}

func service(c *transport.Client) *processingorder.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := processingorder.NewService(c)
	return s
}

var (
	organization = ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization")
	store        = ref("entity/store/22222222-2222-2222-2222-222222222222", "store")
	plan         = ref("entity/processingplan/33333333-3333-3333-3333-333333333333", "processingplan")
	product      = ref("entity/product/44444444-4444-4444-4444-444444444444", "product")
	row          = processingorder.PositionRequest{Assortment: model.Value(product), Quantity: model.Value(111.0)}
	minimal      = processingorder.CreateRequest{Name: "000034", Organization: organization, Store: store, ProcessingPlan: plan, Positions: []processingorder.PositionRequest{row}}
	minimalJSON  = `{"name":"000034","organization":{"meta":{"href":"` + apiBase + `entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},` +
		`"store":{"meta":{"href":"` + apiBase + `entity/store/22222222-2222-2222-2222-222222222222","type":"store"}},` +
		`"processingPlan":{"meta":{"href":"` + apiBase + `entity/processingplan/33333333-3333-3333-3333-333333333333","type":"processingplan"}},` +
		`"positions":[{"assortment":{"meta":{"href":"` + apiBase + `entity/product/44444444-4444-4444-4444-444444444444","type":"product"}},"quantity":111}]}`
	planJSON   = `"processingPlan":{"meta":{"href":"` + apiBase + `entity/processingplan/33333333-3333-3333-3333-333333333333","type":"processingplan"}}`
	fullCreate = processingorder.CreateRequest{
		Name: "000034", Organization: organization, Store: store, ProcessingPlan: plan,
		Positions: []processingorder.PositionRequest{{
			Assortment: model.Value(product), Pack: model.Value(positions.PackRequest{ID: "55555555-5555-5555-5555-555555555555"}),
			Quantity: model.Value(111.0), Reserve: model.Value(0.0),
		}},
		Description: model.Value(""), Code: model.Value("c1"), ExternalCode: model.Value("ext1"),
		Moment:                model.Value(documenttest.DateTime("2016-04-19 13:50:24")),
		DeliveryPlannedMoment: model.Value(documenttest.DateTime("2016-11-23 13:07:00")),
		Applicable:            model.Value(false), Quantity: model.Value(5.0),
		Owner:       model.Value(ref("entity/employee/66666666-6666-6666-6666-666666666666", "employee")),
		Group:       model.Value(ref("entity/group/77777777-7777-7777-7777-777777777777", "group")),
		Project:     model.Null[model.Reference](),
		State:       model.Value(ref("entity/processingorder/metadata/states/88888888-8888-8888-8888-888888888888", "state")),
		Processings: model.Value([]model.Reference{ref("entity/processing/99999999-9999-9999-9999-999999999999", "processing")}),
	}
	fullUpdate = processingorder.UpdateRequest{
		ProcessingPlan: plan, Name: model.Value("000035"),
		Store: model.Value(ref("entity/store/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "store")),
		Positions: model.Value([]processingorder.PositionRequest{
			{ID: model.Value(positionID), Quantity: model.Value(2.5)},
			{Assortment: model.Value(product), Quantity: model.Value(1.0), Reserve: model.Value(0.5)},
		}),
		DeliveryPlannedMoment: model.Null[model.DateTime](),
		Applicable:            model.Value(true), Quantity: model.Value(0.5),
		Processings: model.Value([]model.Reference{}),
	}
)

func moscow(v string) time.Time { return time.Time(documenttest.DateTime(v)) }

// Фикстуры — примеры ответов источника: «Получить Заказ на производство», первая строка списка
// (с processings), «Получить позиции» и шаблон на основе техкарты.
func TestFixtures(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got processingorder.ProcessingOrder
			if err := codec.Unmarshal(read(t, "read.json"), &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != "b5556153-b164-11e6-5bed-427b0000006f" || got.AccountID != "d55cbfba-91f1-11e6-5bed-427b00000000" || got.Name != "Заказ на производство" ||
				got.Description != "Комментарий" || got.ExternalCode != "wWAJXKZFgoOZVKd41Dzzz2" || !got.Applicable || got.Shared || got.Quantity != 5 ||
				got.Meta.Type != "processingorder" {
				t.Fatalf("document = %+v", got)
			}
			if !time.Time(got.Updated).Equal(moscow("2016-11-23 13:07:57")) || !time.Time(got.Moment).Equal(moscow("2016-11-23 13:00:00")) ||
				!time.Time(got.DeliveryPlannedMoment).Equal(moscow("2016-11-23 13:07:00")) || !got.Created.IsZero() {
				t.Fatalf("dates = %v / %v / %v", got.Updated, got.Moment, got.DeliveryPlannedMoment)
			}
			if got.ProcessingPlan == nil || got.ProcessingPlan.Meta.Href != apiBase+"entity/processingplan/c38e50b0-acdc-11e6-5bed-427b0000009e" ||
				got.Store == nil || got.Store.Meta.Type != "store" || got.Organization == nil || got.Organization.Meta.Type != "organization" ||
				got.Owner == nil || got.Group == nil || got.Project == nil || got.Project.Meta.Type != "project" || got.State != nil || got.Processings != nil {
				t.Fatalf("references = %+v", got)
			}
			if got.Positions == nil || got.Positions.Meta.Type != "processingorderposition" || got.Positions.Meta.Size != 2 || got.Positions.Rows != nil {
				t.Fatalf("positions = %+v", got.Positions)
			}

			var listed processingorder.ProcessingOrder
			if err := codec.Unmarshal(read(t, "list-row.json"), &listed); err != nil {
				t.Fatal(err)
			}
			if listed.ID != "5fbf79f4-adac-11e6-5bed-427b0000006a" || !listed.Printed || !listed.Published || listed.Quantity != 1 || len(listed.Processings) != 2 ||
				listed.Processings[0].Meta.Href != apiBase+"entity/processing/bcd6e6b5-b002-11e6-5bed-427b0000000b" || listed.Processings[1].Meta.Type != "processing" {
				t.Fatalf("listed = %+v", listed)
			}

			var page model.Collection[processingorder.Position]
			if err := codec.Unmarshal(read(t, "positions.json"), &page); err != nil {
				t.Fatal(err)
			}
			first := page.Rows[0]
			if len(page.Rows) != 2 || page.Meta.Size != 2 || first.ID != "c7218ccd-afcc-11e6-5bed-427b00000069" || first.Quantity != 45 || first.Reserve != 45 ||
				first.Meta.Type != "processingorderposition" || first.AccountID != "d55cbfba-91f1-11e6-5bed-427b00000000" || first.Assortment == nil ||
				first.Assortment.Meta.Href != apiBase+"entity/product/0de151c1-acdc-11e6-5bed-427b00000080" || first.Pack != nil {
				t.Fatalf("positions = %+v", page)
			}

			var template processingorder.ProcessingOrder
			if err := codec.Unmarshal(read(t, "template.json"), &template); err != nil {
				t.Fatal(err)
			}
			if template.ID != "" || template.Quantity != 1 || template.ProcessingPlan == nil || template.Store == nil || template.Positions == nil ||
				len(template.Positions.Rows) != 1 || template.Positions.Rows[0].Quantity != 1 || template.Positions.Rows[0].Assortment == nil {
				t.Fatalf("template = %+v", template)
			}

			// Поля таблицы, которых нет в примерах, и строка с упаковкой.
			var sparse processingorder.ProcessingOrder
			if err := codec.Unmarshal([]byte(`{"id":"x","code":"c","syncId":"s","deleted":"2016-11-23 13:07:57.000","shared":true,"quantity":0.25,`+
				`"state":{"meta":{"type":"state"}},"attributes":[{"id":"a","name":"n","type":"string","value":"v"}],"files":{"meta":{"href":"f","size":0}},`+
				`"positions":{"meta":{"href":"h"},"rows":[{"id":"p","quantity":0.5,"reserve":0.25,"pack":{"id":"k","quantity":6}}]},"sum":100,"future":1}`), &sparse); err != nil {
				t.Fatal(err)
			}
			pr := sparse.Positions.Rows[0]
			if sparse.Code != "c" || sparse.SyncID != "s" || !sparse.Shared || sparse.Quantity != 0.25 || sparse.State == nil || len(sparse.Attributes) != 1 ||
				sparse.Files == nil || sparse.Deleted.IsZero() || pr.Pack == nil || pr.Pack.ID != "k" || pr.Pack.Quantity != 6 || pr.Reserve != 0.25 || pr.Quantity != 0.5 {
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
			{"create full", fullCreate, read(t, "create-full.golden.json")},
			{"create minimal", minimal, []byte(minimalJSON)},
			{"update", fullUpdate, read(t, "update.golden.json")},
			{"update plan only", processingorder.UpdateRequest{ProcessingPlan: plan}, []byte(`{` + planJSON + `}`)},
			{"position update", processingorder.PositionRequest{Quantity: model.Value(0.0), Reserve: model.Value(0.0), Pack: model.Null[positions.PackRequest]()}, []byte(`{"quantity":0,"reserve":0,"pack":null}`)},
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
	// attributes, shared и syncId таблица помечает «Только для чтения»; sum в таблице нет.
	for _, typ := range []reflect.Type{reflect.TypeFor[processingorder.CreateRequest](), reflect.TypeFor[processingorder.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "Shared", "SyncID", "Attributes", "Files", "Sum"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
	for _, name := range []string{"Meta", "AccountID", "Price", "VAT", "Discount"} {
		if _, ok := reflect.TypeFor[processingorder.PositionRequest]().FieldByName(name); ok {
			t.Errorf("PositionRequest contains %s", name)
		}
	}
}

func TestHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя, объём и связи — только из ответа.
	const response = `{"id":"server-id","name":"server-name","quantity":7.5,"processings":[{"meta":{"href":"https://example.test/processing/p","type":"processing"}}],` +
		`"positions":{"meta":{"href":"https://example.test/positions","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*processingorder.ProcessingOrder, error)
	}{
		{"get", "GET", docPath, nil, func(ctx context.Context, c *transport.Client) (*processingorder.ProcessingOrder, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/processingorder", read(t, "create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*processingorder.ProcessingOrder, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", docPath, read(t, "update.golden.json"), func(ctx context.Context, c *transport.Client) (*processingorder.ProcessingOrder, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *processingorder.ProcessingOrder) {
						if got.ID != "server-id" || got.Name != "server-name" || got.Quantity != 7.5 || len(got.Processings) != 1 ||
							got.Processings[0].Meta.Href != "https://example.test/processing/p" || got.Positions.Meta.Size != 1 || got.Positions.Rows != nil {
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
	const rowJSON = `{"meta":{"type":"processingorderposition"},"id":"` + positionID + `","quantity":111,"reserve":0.5,"pack":{"id":"k"},"assortment":{"meta":{"type":"product"}}}`
	checkRow := func(got *processingorder.Position) {
		if got.ID != positionID || got.Quantity != 111 || got.Reserve != 0.5 || got.Pack == nil || got.Pack.ID != "k" || got.Assortment == nil || got.Meta.Type != "processingorderposition" {
			t.Fatalf("row = %+v", got)
		}
	}
	t.Run("get positions", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/positions?limit=2&offset=7", "", `{"meta":{"href":"h"},"rows":[`+rowJSON+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[processingorder.Position], error) {
				return service(c).GetPositions(ctx, testID, positions.PageOptions{Limit: 2, Offset: 7})
			}, func(got *model.Collection[processingorder.Position]) { checkRow(&got.Rows[0]) })
	})
	t.Run("get positions default page", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/positions?limit=1000&offset=0", "", string(read(t, "positions.json")),
			func(ctx context.Context, c *transport.Client) (*model.Collection[processingorder.Position], error) {
				return service(c).GetPositions(ctx, testID, positions.PageOptions{})
			}, func(got *model.Collection[processingorder.Position]) {
				if len(got.Rows) != 2 || got.Rows[1].ID != "c72196c6-afcc-11e6-5bed-427b0000006a" {
					t.Fatalf("page = %+v", got)
				}
			})
	})
	t.Run("get position", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/positions/"+positionID, "", rowJSON, func(ctx context.Context, c *transport.Client) (*processingorder.Position, error) {
			return service(c).GetPosition(ctx, testID, positionID)
		}, checkRow)
	})
	t.Run("update position", func(t *testing.T) {
		body := `{"assortment":{"meta":{"href":"` + apiBase + `entity/product/44444444-4444-4444-4444-444444444444","type":"product"}},"quantity":111}`
		httpCases(t, "PUT", docPath+"/positions/"+positionID, body, rowJSON, func(ctx context.Context, c *transport.Client) (*processingorder.Position, error) {
			return service(c).UpdatePosition(ctx, testID, positionID, row)
		}, checkRow)
	})
}

// Методы без модели документа в ответе, шаблон и массовые: один запрос с точным телом, без
// повторов при 503.
func TestWrites(t *testing.T) {
	const base = "{base}/entity/processingorder"
	docMeta := `{"href":"` + base + `/` + testID + `","metadataHref":"` + base + `/metadata","type":"processingorder","mediaType":"application/json"}`
	for _, tc := range []struct {
		name, method, path, body, response string
		empty                              bool // ответ без тела: 204 тоже успех
		invoke                             func(context.Context, *transport.Client) error
	}{
		{"delete", "DELETE", docPath, "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Delete(ctx, testID) }},
		{"trash", "POST", docPath + "/trash", "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Trash(ctx, testID) }},
		{"new", "PUT", "entity/processingorder/new", "", `{"name":"00001","quantity":1}`, false, func(ctx context.Context, c *transport.Client) error {
			got, err := service(c).New(ctx, nil)
			if err == nil && (got.Name != "00001" || got.Quantity != 1 || got.ID != "") {
				return errors.New("unexpected template")
			}
			return err
		}},
		{"new by processing plan", "PUT", "entity/processingorder/new", `{` + planJSON + `}`, string(read(t, "template.json")), false, func(ctx context.Context, c *transport.Client) error {
			got, err := service(c).New(ctx, map[string]any{"processingPlan": plan})
			if err == nil && (got.Quantity != 1 || got.Positions == nil || len(got.Positions.Rows) != 1 || got.Positions.Rows[0].Assortment == nil) {
				return errors.New("unexpected template")
			}
			return err
		}},
		{"create many", "POST", "entity/processingorder", `[` + minimalJSON + `]`, `[{"id":"x"}]`, false, func(ctx context.Context, c *transport.Client) error {
			r, err := service(c).CreateMany(ctx, []processingorder.CreateRequest{minimal})
			if err == nil && r[0].Value.ID != "x" {
				return errors.New("unexpected result")
			}
			return err
		}},
		{"update many", "POST", "entity/processingorder", `[{"meta":` + docMeta + `,` + planJSON + `,"processings":[]}]`, `[{"id":"x"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				_, err := service(c).UpdateMany(ctx, []batch.Update[processingorder.UpdateRequest]{{ID: testID, Request: processingorder.UpdateRequest{ProcessingPlan: plan, Processings: model.Value([]model.Reference{})}}})
				return err
			}},
		{"delete many", "POST", "entity/processingorder/delete", `[{"meta":` + docMeta + `}]`, `[{"info":"ok"}]`, false, func(ctx context.Context, c *transport.Client) error {
			_, err := service(c).DeleteMany(ctx, []string{testID})
			return err
		}},
		{"delete position", "DELETE", docPath + "/positions/" + positionID, "", "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeletePosition(ctx, testID, positionID)
		}},
		{"delete many positions", "POST", docPath + "/positions/delete", `[{"meta":{"href":"` + base + `/` + testID + `/positions/` + positionID + `","type":"processingorderposition","mediaType":"application/json"}}]`, "", true,
			func(ctx context.Context, c *transport.Client) error {
				return service(c).DeleteManyPositions(ctx, testID, []string{positionID})
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
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "processingorder") ||
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
	limit := minimal
	limit.Positions = make([]processingorder.PositionRequest, 1000)
	if _, err := s.Create(ctx, limit); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(ctx, testID, processingorder.UpdateRequest{ProcessingPlan: plan, Positions: model.Value(make([]processingorder.PositionRequest, 1000))}); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d", calls.Load())
	}
	create := func(change func(*processingorder.CreateRequest)) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error {
			input := minimal
			change(&input)
			_, err := s.Create(ctx, input, o...)
			return err
		}
	}
	update := func(input processingorder.UpdateRequest) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, input, o...)
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
		"get bad id":          func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":            func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no name":      create(func(r *processingorder.CreateRequest) { r.Name = "" }),
		"create no org":       create(func(r *processingorder.CreateRequest) { r.Organization = model.Reference{} }),
		"create no store":     create(func(r *processingorder.CreateRequest) { r.Store = model.Reference{} }),
		"create no plan":      create(func(r *processingorder.CreateRequest) { r.ProcessingPlan = model.Reference{} }),
		"create no positions": create(func(r *processingorder.CreateRequest) { r.Positions = nil }),
		"create 1001":         create(func(r *processingorder.CreateRequest) { r.Positions = make([]processingorder.PositionRequest, 1001) }),
		"update empty":        update(processingorder.UpdateRequest{}),
		"update no plan":      update(processingorder.UpdateRequest{Name: model.Value("x")}),
		"update 1001":         update(processingorder.UpdateRequest{ProcessingPlan: plan, Positions: model.Value(make([]processingorder.PositionRequest, 1001))}),
		"update bad id":       func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", fullUpdate, o...); return err },
		"delete bad id":       func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
		"trash bad id":        func(o ...transport.RequestOption) error { return s.Trash(ctx, "bad", o...) },
		"create many empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateMany(ctx, nil, o...)
			return err
		},
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[processingorder.UpdateRequest]{{ID: "bad", Request: fullUpdate}}, o...)
			return err
		},
		"update many empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[processingorder.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"delete many bad id": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, []string{"bad"}, o...)
			return err
		},
		"positions bad id": getPositions("bad", positions.PageOptions{}),
		"positions fields": getPositions(testID, positions.PageOptions{Fields: "stock"}),
		"positions limit":  getPositions(testID, positions.PageOptions{Limit: 1001}),
		"positions offset": getPositions(testID, positions.PageOptions{Offset: -1}),
		"get position bad": func(o ...transport.RequestOption) error {
			_, err := s.GetPosition(ctx, testID, "bad", o...)
			return err
		},
		"update position empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdatePosition(ctx, testID, positionID, processingorder.PositionRequest{}, o...)
			return err
		},
		"update position bad doc": func(o ...transport.RequestOption) error {
			_, err := s.UpdatePosition(ctx, "bad", positionID, row, o...)
			return err
		},
		"delete position bad": func(o ...transport.RequestOption) error { return s.DeletePosition(ctx, testID, "bad", o...) },
		"delete many positions empty": func(o ...transport.RequestOption) error {
			return s.DeleteManyPositions(ctx, testID, nil, o...)
		},
		"delete many positions 1001": func(o ...transport.RequestOption) error {
			return s.DeleteManyPositions(ctx, testID, make([]string, 1001), o...)
		},
		"delete many positions bad row": func(o ...transport.RequestOption) error {
			return s.DeleteManyPositions(ctx, testID, []string{positionID, "bad"}, o...)
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "processingorder") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := processingorder.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestList(t *testing.T) {
	listtest.Run(t, "processingorder", read(t, "list-row.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[processingorder.ProcessingOrder] {
		return service(c).List
	}, func(t *testing.T, got processingorder.ProcessingOrder) {
		if got.ID != "5fbf79f4-adac-11e6-5bed-427b0000006a" || len(got.Processings) != 2 || got.ProcessingPlan == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
