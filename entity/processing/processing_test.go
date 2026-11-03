package processing_test

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
	"github.com/eteeet/moysklad-sdk-go/entity/processing"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "12121212-1212-1212-1212-121212121212"
	fixtureID = "493ddf6b-aff9-11e6-5bed-427b00000076"
	rowID     = "34f6344f-015e-11e6-9464-e4de0000006c"
	docPath   = "entity/processing/" + testID
	apiBase   = "https://api.moysklad.ru/api/remap/1.2/"
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
	return model.Reference{Meta: model.Meta{Href: apiBase + path, Type: model.MetaType(typ)}}
}

func service(c *transport.Client) *processing.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := processing.NewService(c)
	return s
}

var (
	organization   = ref("entity/organization/11111111-1111-1111-1111-111111111111", "organization")
	materialsStore = ref("entity/store/22222222-2222-2222-2222-222222222222", "store")
	productsStore  = ref("entity/store/33333333-3333-3333-3333-333333333333", "store")
	product        = ref("entity/product/99999999-9999-9999-9999-999999999999", "product")
	minimal        = processing.CreateRequest{Organization: organization, MaterialsStore: materialsStore, ProductsStore: productsStore}
	minimalJSON    = `{"organization":{"meta":{"href":"` + apiBase + `entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},` +
		`"materialsStore":{"meta":{"href":"` + apiBase + `entity/store/22222222-2222-2222-2222-222222222222","type":"store"}},` +
		`"productsStore":{"meta":{"href":"` + apiBase + `entity/store/33333333-3333-3333-3333-333333333333","type":"store"}}}`
	fullCreate = processing.CreateRequest{
		Organization: organization, MaterialsStore: materialsStore, ProductsStore: productsStore,
		Name: model.Value("Сборка"), Description: model.Value(""), Code: model.Value("c1"), ExternalCode: model.Value("ext1"),
		Moment: model.Value(documenttest.DateTime("2023-04-26 07:23:00")), Applicable: model.Value(false), Shared: model.Value(true),
		OrganizationAccount: model.Value(ref("entity/organization/11111111-1111-1111-1111-111111111111/accounts/44444444-4444-4444-4444-444444444444", "account")),
		Owner:               model.Value(ref("entity/employee/55555555-5555-5555-5555-555555555555", "employee")),
		Group:               model.Value(ref("entity/group/66666666-6666-6666-6666-666666666666", "group")),
		Project:             model.Null[model.Reference](),
		State:               model.Value(ref("entity/processing/metadata/states/77777777-7777-7777-7777-777777777777", "state")),
		Attributes:          model.Value([]model.Attribute{}),
		ProcessingPlan:      model.Value(ref("entity/processingplan/88888888-8888-8888-8888-888888888888", "processingplan")),
		ProcessingSum:       model.Value(150.5), Quantity: model.Value(3.0),
		Materials: model.Value([]processing.PositionRequest{{
			Assortment: model.Value(ref("entity/variant/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "variant")), Quantity: model.Value(20.0),
			Slot: model.Value(ref("entity/store/22222222-2222-2222-2222-222222222222/slots/bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb", "slot")),
		}}),
		Products:        model.Value([]processing.PositionRequest{{Assortment: model.Value(product), Quantity: model.Value(10.0)}}),
		ProcessingOrder: model.Value(ref("entity/processingorder/cccccccc-cccc-cccc-cccc-cccccccccccc", "processingorder")),
	}
	fullUpdate = processing.UpdateRequest{
		MaterialsStore: model.Value(productsStore),
		State:          model.Null[model.Reference](),
		ProcessingPlan: model.Value(ref("entity/processingplan/dddddddd-dddd-dddd-dddd-dddddddddddd", "processingplan")),
		ProcessingSum:  model.Value(0.0), Quantity: model.Value(5.0),
		Materials: model.Value([]processing.PositionRequest{}),
		Products: model.Value([]processing.PositionRequest{
			{ID: model.Value("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"), Quantity: model.Value(400.0)},
			{ID: model.Value("ffffffff-ffff-ffff-ffff-ffffffffffff")},
		}),
	}
)

func TestFixture(t *testing.T) {
	data := read(t, "custom/processing/read.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got processing.Processing
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.Name != "я1" || got.ProcessingSum != 10000 || got.Quantity != 1 || !got.Applicable || !got.Printed || !got.Published || got.Shared ||
				got.ExternalCode != "JhQJY8u1isyuqHyn7B6Wx3" || got.Meta.Type != "processing" {
				t.Fatalf("document = %+v", got)
			}
			if !time.Time(got.Updated).Equal(time.Time(documenttest.DateTime("2016-11-21 17:46:29"))) || !time.Time(got.Moment).Equal(time.Time(documenttest.DateTime("2016-11-21 17:46:00"))) {
				t.Fatalf("dates = %v / %v", got.Updated, got.Moment)
			}
			// processingPlan — ссылка на техкарту; склады материалов и продуктов различаются.
			if got.ProcessingPlan == nil || got.ProcessingPlan.Meta.Type != "processingplan" || got.ProcessingPlan.Meta.Href != apiBase+"entity/processingplan/c38e50b0-acdc-11e6-5bed-427b0000009e" {
				t.Fatalf("processingPlan = %+v", got.ProcessingPlan)
			}
			if got.MaterialsStore == nil || got.MaterialsStore.Meta.Href != apiBase+"entity/store/d5e311c0-91f1-11e6-5bed-427b00000053" ||
				got.ProductsStore == nil || got.ProductsStore.Meta.Href != apiBase+"entity/store/c61ce912-a747-11e6-5bed-427b000000a8" {
				t.Fatalf("stores = %+v / %+v", got.MaterialsStore, got.ProductsStore)
			}
			if got.ProcessingOrder == nil || got.ProcessingOrder.Meta.Type != "processingorder" || got.Organization == nil || got.Owner == nil || got.Group == nil || got.Project == nil {
				t.Fatalf("references = %+v", got)
			}
			if got.Materials == nil || got.Materials.Meta.Type != "processingpositionmaterial" || got.Materials.Meta.Size != 2 || got.Materials.Rows != nil ||
				got.Products == nil || got.Products.Meta.Type != "processingpositionresult" || got.Products.Meta.Size != 2 {
				t.Fatalf("collections = %+v / %+v", got.Materials, got.Products)
			}
			// Поля таблицы, которых нет в примере ответа, строки при expand и строки шаблона без meta.
			var sparse processing.Processing
			if err := codec.Unmarshal([]byte(`{"id":"x","syncId":"s","code":"c","description":"d","created":"2023-04-26 07:23:59.302","deleted":"2023-04-27 07:23:59.302",`+
				`"organizationAccount":{"meta":{"type":"account"}},"state":{"meta":{"type":"state"}},"attributes":[{"id":"a"}],"files":{"meta":{"href":"f","size":0}},`+
				`"materials":{"meta":{"href":"m"},"rows":[{"id":"r","accountId":"acc","quantity":2.5,"assortment":{"meta":{"type":"variant"}},"slot":{"meta":{"type":"slot"}}}]},`+
				`"products":{"rows":[{"quantity":1,"assortment":{"meta":{"type":"product"}}}]},"future":1}`), &sparse); err != nil {
				t.Fatal(err)
			}
			m, p := sparse.Materials.Rows[0], sparse.Products.Rows[0]
			if sparse.SyncID != "s" || sparse.Code != "c" || sparse.Description != "d" || sparse.Created.IsZero() || sparse.Deleted.IsZero() || sparse.OrganizationAccount == nil ||
				sparse.State == nil || len(sparse.Attributes) != 1 || sparse.Files == nil || sparse.Organization != nil || sparse.ProcessingPlan != nil {
				t.Fatalf("sparse = %+v", sparse)
			}
			if m.ID != "r" || m.AccountID != "acc" || m.Quantity != 2.5 || m.Assortment == nil || m.Assortment.Meta.Type != "variant" || m.Slot == nil || m.Slot.Meta.Type != "slot" ||
				p.Quantity != 1 || p.Assortment == nil || p.Slot != nil || sparse.Products.Meta.Href != "" {
				t.Fatalf("rows = %+v / %+v", m, p)
			}
		})
	}
}

// processingSum — копейки float64: 0, 0.0, дробь и большое значение читаются без потерь обоими кодеками.
func TestMoneyJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for raw, want := range map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "123456789.01": 123456789.01} {
			var got processing.Processing
			if err := codec.Unmarshal([]byte(`{"processingSum":`+raw+`}`), &got); err != nil || got.ProcessingSum != want {
				t.Fatalf("%s: processingSum=%s: got %v, err = %v", codec.Name, raw, got.ProcessingSum, err)
			}
			data, err := codec.Marshal(processing.UpdateRequest{ProcessingSum: model.Value(want)})
			if err != nil {
				t.Fatal(err)
			}
			var back struct {
				ProcessingSum *float64 `json:"processingSum"`
			}
			if err := codec.Unmarshal(data, &back); err != nil || back.ProcessingSum == nil || *back.ProcessingSum != want {
				t.Fatalf("%s: request %s", codec.Name, data)
			}
		}
	}
}

func TestRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", fullCreate, read(t, "custom/processing/create-full.golden.json")},
			{"create minimal", minimal, []byte(minimalJSON)},
			{"create without plan", processing.CreateRequest{Organization: organization, MaterialsStore: materialsStore, ProductsStore: productsStore,
				ProcessingSum: model.Value(50.0), Products: model.Value([]processing.PositionRequest{{Assortment: model.Value(product), Quantity: model.Value(10.0)}})},
				[]byte(minimalJSON[:len(minimalJSON)-1] + `,"processingSum":50,"products":[{"assortment":{"meta":{"href":"` + apiBase + `entity/product/99999999-9999-9999-9999-999999999999","type":"product"}},"quantity":10}]}`)},
			{"update", fullUpdate, read(t, "custom/processing/update.golden.json")},
			{"row update", processing.PositionRequest{Quantity: model.Value(0.0), Slot: model.Null[model.Reference]()}, []byte(`{"quantity":0,"slot":null}`)},
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
	for _, typ := range []reflect.Type{reflect.TypeFor[processing.CreateRequest](), reflect.TypeFor[processing.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "SyncID", "Files", "Rate", "Sum"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only or foreign %s", typ.Name(), name)
			}
		}
	}
	for _, name := range []string{"Meta", "AccountID", "Price", "VAT", "Pack", "Things"} {
		if _, ok := reflect.TypeFor[processing.PositionRequest]().FieldByName(name); ok {
			t.Errorf("PositionRequest contains read-only or foreign %s", name)
		}
	}
}

func TestHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя, затраты и объём — только из ответа.
	const response = `{"id":"server-id","name":"server-name","processingSum":24691.5,"quantity":7,"materials":{"meta":{"href":"https://example.test/materials","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*processing.Processing, error)
	}{
		{"get", "GET", docPath, nil, func(ctx context.Context, c *transport.Client) (*processing.Processing, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/processing", read(t, "custom/processing/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*processing.Processing, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", docPath, read(t, "custom/processing/update.golden.json"), func(ctx context.Context, c *transport.Client) (*processing.Processing, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *processing.Processing) {
						if got.ID != "server-id" || got.Name != "server-name" || got.ProcessingSum != 24691.5 || got.Quantity != 7 || got.Materials == nil || got.Materials.Meta.Size != 1 {
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

func TestRowsHTTP(t *testing.T) {
	const assortmentJSON = `{"meta":{"href":"` + apiBase + `entity/product/99999999-9999-9999-9999-999999999999","type":"product"}}`
	create := processing.PositionRequest{Assortment: model.Value(product), Quantity: model.Value(10.5)}
	update := processing.PositionRequest{Quantity: model.Value(111.0)}
	for _, set := range []struct {
		path, metaType string
		list           func(*processing.Service) func(context.Context, string, positions.PageOptions, ...transport.RequestOption) (*model.Collection[processing.Position], error)
		get            func(*processing.Service) func(context.Context, string, string, ...transport.RequestOption) (*processing.Position, error)
		create         func(*processing.Service) func(context.Context, string, processing.PositionRequest, ...transport.RequestOption) (*processing.Position, error)
		update         func(*processing.Service) func(context.Context, string, string, processing.PositionRequest, ...transport.RequestOption) (*processing.Position, error)
	}{
		{"materials", "processingpositionmaterial", func(s *processing.Service) func(context.Context, string, positions.PageOptions, ...transport.RequestOption) (*model.Collection[processing.Position], error) {
			return s.GetMaterials
		}, func(s *processing.Service) func(context.Context, string, string, ...transport.RequestOption) (*processing.Position, error) {
			return s.GetMaterial
		}, func(s *processing.Service) func(context.Context, string, processing.PositionRequest, ...transport.RequestOption) (*processing.Position, error) {
			return s.CreateMaterial
		}, func(s *processing.Service) func(context.Context, string, string, processing.PositionRequest, ...transport.RequestOption) (*processing.Position, error) {
			return s.UpdateMaterial
		}},
		{"products", "processingpositionresult", func(s *processing.Service) func(context.Context, string, positions.PageOptions, ...transport.RequestOption) (*model.Collection[processing.Position], error) {
			return s.GetProducts
		}, func(s *processing.Service) func(context.Context, string, string, ...transport.RequestOption) (*processing.Position, error) {
			return s.GetProduct
		}, func(s *processing.Service) func(context.Context, string, processing.PositionRequest, ...transport.RequestOption) (*processing.Position, error) {
			return s.CreateProduct
		}, func(s *processing.Service) func(context.Context, string, string, processing.PositionRequest, ...transport.RequestOption) (*processing.Position, error) {
			return s.UpdateProduct
		}},
	} {
		row := `{"meta":{"type":"` + set.metaType + `"},"id":"` + rowID + `","accountId":"acc","quantity":45.5,"assortment":` + assortmentJSON + `,"slot":{"meta":{"type":"slot"}}}`
		check := func(got *processing.Position) {
			if got.ID != rowID || got.Meta.Type != model.MetaType(set.metaType) || got.Quantity != 45.5 || got.Assortment == nil || got.Assortment.Meta.Href != product.Meta.Href || got.Slot == nil {
				t.Fatalf("row = %+v", got)
			}
		}
		base := docPath + "/" + set.path
		t.Run(set.path+"/get page", func(t *testing.T) {
			httpCases(t, "GET", base+"?limit=2&offset=7", "", `{"meta":{"href":"h","size":3},"rows":[`+row+`]}`,
				func(ctx context.Context, c *transport.Client) (*model.Collection[processing.Position], error) {
					return set.list(service(c))(ctx, testID, positions.PageOptions{Limit: 2, Offset: 7})
				}, func(got *model.Collection[processing.Position]) {
					if got.Meta.Size != 3 || len(got.Rows) != 1 {
						t.Fatalf("page = %+v", got)
					}
					check(&got.Rows[0])
				})
		})
		t.Run(set.path+"/get default page", func(t *testing.T) {
			httpCases(t, "GET", base+"?limit=1000&offset=0", "", `{"meta":{"href":"h"},"rows":[`+row+`]}`,
				func(ctx context.Context, c *transport.Client) (*model.Collection[processing.Position], error) {
					return set.list(service(c))(ctx, testID, positions.PageOptions{})
				}, func(got *model.Collection[processing.Position]) { check(&got.Rows[0]) })
		})
		t.Run(set.path+"/get", func(t *testing.T) {
			httpCases(t, "GET", base+"/"+rowID, "", row, func(ctx context.Context, c *transport.Client) (*processing.Position, error) {
				return set.get(service(c))(ctx, testID, rowID)
			}, check)
		})
		// Пример источника: запрос — объект, ответ — массив из одного объекта.
		t.Run(set.path+"/create", func(t *testing.T) {
			httpCases(t, "POST", base, `{"assortment":`+assortmentJSON+`,"quantity":10.5}`, `[`+row+`]`, func(ctx context.Context, c *transport.Client) (*processing.Position, error) {
				return set.create(service(c))(ctx, testID, create)
			}, check)
		})
		t.Run(set.path+"/create object response", func(t *testing.T) {
			httpCases(t, "POST", base, `{"assortment":`+assortmentJSON+`,"quantity":10.5}`, row, func(ctx context.Context, c *transport.Client) (*processing.Position, error) {
				return set.create(service(c))(ctx, testID, create)
			}, check)
		})
		t.Run(set.path+"/update", func(t *testing.T) {
			httpCases(t, "PUT", base+"/"+rowID, `{"quantity":111}`, row, func(ctx context.Context, c *transport.Client) (*processing.Position, error) {
				return set.update(service(c))(ctx, testID, rowID, update)
			}, check)
		})
	}
}

// Методы без модели в ответе, шаблон и массовые: один запрос с точным телом, без повторов при 503.
func TestWrites(t *testing.T) {
	const base = "{base}/entity/processing"
	docMeta := `{"href":"` + base + `/` + testID + `","metadataHref":"` + base + `/metadata","type":"processing","mediaType":"application/json"}`
	rowMeta := func(path, typ string) string {
		return `[{"meta":{"href":"` + base + `/` + testID + `/` + path + `/` + rowID + `","type":"` + typ + `","mediaType":"application/json"}}]`
	}
	for _, tc := range []struct {
		name, method, path, body, response string
		empty                              bool // ответ без тела: 204 тоже успех
		invoke                             func(context.Context, *transport.Client) error
	}{
		{"delete", "DELETE", docPath, "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Delete(ctx, testID) }},
		{"trash", "POST", docPath + "/trash", "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Trash(ctx, testID) }},
		{"new", "PUT", "entity/processing/new", "", `{"name":"00001","processingSum":1000,"quantity":1}`, false, func(ctx context.Context, c *transport.Client) error {
			got, err := service(c).New(ctx, nil)
			if err == nil && (got.Name != "00001" || got.ProcessingSum != 1000 || got.ID != "") {
				return errors.New("moysklad: unexpected processing template")
			}
			return err
		}},
		{"new by plan", "PUT", "entity/processing/new", `{"processingPlan":{"meta":{"href":"` + apiBase + `entity/processingplan/88888888-8888-8888-8888-888888888888","type":"processingplan"}}}`,
			`{"processingPlan":{"meta":{"type":"processingplan"}},"products":{"rows":[{"quantity":2,"assortment":{"meta":{"type":"product"}}}]},"materials":{"rows":[]}}`, false,
			func(ctx context.Context, c *transport.Client) error {
				got, err := service(c).New(ctx, map[string]any{"processingPlan": ref("entity/processingplan/88888888-8888-8888-8888-888888888888", "processingplan")})
				if err == nil && (got.ProcessingPlan == nil || got.Products == nil || len(got.Products.Rows) != 1 || got.Products.Rows[0].Quantity != 2 || got.Materials == nil || len(got.Materials.Rows) != 0) {
					return errors.New("moysklad: unexpected processing template")
				}
				return err
			}},
		{"create many", "POST", "entity/processing", `[` + minimalJSON + `]`, `[{"id":"x"}]`, false, func(ctx context.Context, c *transport.Client) error {
			r, err := service(c).CreateMany(ctx, []processing.CreateRequest{minimal})
			if err == nil && r[0].Value.ID != "x" {
				t.Errorf("result = %+v", r)
			}
			return err
		}},
		{"update many", "POST", "entity/processing", `[{"meta":` + docMeta + `,"quantity":5,"processingSum":2000}]`, `[{"id":"x"}]`, false, func(ctx context.Context, c *transport.Client) error {
			_, err := service(c).UpdateMany(ctx, []batch.Update[processing.UpdateRequest]{{ID: testID, Request: processing.UpdateRequest{Quantity: model.Value(5.0), ProcessingSum: model.Value(2000.0)}}})
			return err
		}},
		{"delete many", "POST", "entity/processing/delete", `[{"meta":` + docMeta + `}]`, `[{"info":"ok"}]`, false, func(ctx context.Context, c *transport.Client) error {
			_, err := service(c).DeleteMany(ctx, []string{testID})
			return err
		}},
		{"delete material", "DELETE", docPath + "/materials/" + rowID, "", "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeleteMaterial(ctx, testID, rowID)
		}},
		{"delete many materials", "POST", docPath + "/materials/delete", rowMeta("materials", "processingpositionmaterial"), "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeleteManyMaterials(ctx, testID, []string{rowID})
		}},
		{"delete product", "DELETE", docPath + "/products/" + rowID, "", "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeleteProduct(ctx, testID, rowID)
		}},
		{"delete many products", "POST", docPath + "/products/delete", rowMeta("products", "processingpositionresult"), "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeleteManyProducts(ctx, testID, []string{rowID})
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
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "processing") ||
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

// Ошибки входа — до HTTP, ResponseInfo сброшен; ровно 1000 строк каждого набора уходят.
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
	rows := func(n int) model.Field[[]processing.PositionRequest] {
		return model.Value(make([]processing.PositionRequest, n))
	}
	limit := minimal
	limit.Materials, limit.Products = rows(1000), rows(1000)
	if _, err := s.Create(ctx, limit); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(ctx, testID, processing.UpdateRequest{Materials: rows(1000), Products: rows(1000)}); err != nil {
		t.Fatal(err)
	}
	create := func(change func(*processing.CreateRequest)) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error {
			input := minimal
			change(&input)
			_, err := s.Create(ctx, input, o...)
			return err
		}
	}
	type pageFunc = func(context.Context, string, positions.PageOptions, ...transport.RequestOption) (*model.Collection[processing.Position], error)
	page := func(get pageFunc, id string, p positions.PageOptions) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error { _, err := get(ctx, id, p, o...); return err }
	}
	type rowFunc = func(context.Context, string, processing.PositionRequest, ...transport.RequestOption) (*processing.Position, error)
	add := func(add rowFunc, id string, input processing.PositionRequest) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error { _, err := add(ctx, id, input, o...); return err }
	}
	quantity := processing.PositionRequest{Quantity: model.Value(1.0)}
	assortment := processing.PositionRequest{Assortment: model.Value(product)}
	valid := processing.PositionRequest{Assortment: model.Value(product), Quantity: model.Value(1.0)}
	calls.Store(0)
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id":              func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":                func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no organization":  create(func(r *processing.CreateRequest) { r.Organization = model.Reference{} }),
		"create no materialStore": create(func(r *processing.CreateRequest) { r.MaterialsStore = model.Reference{} }),
		"create no productsStore": create(func(r *processing.CreateRequest) { r.ProductsStore = model.Reference{} }),
		"create 1001 materials":   create(func(r *processing.CreateRequest) { r.Materials = rows(1001) }),
		"create 1001 products":    create(func(r *processing.CreateRequest) { r.Products = rows(1001) }),
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, processing.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", fullUpdate, o...); return err },
		"update 1001 materials": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, processing.UpdateRequest{Materials: rows(1001)}, o...)
			return err
		},
		"update 1001 products": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, processing.UpdateRequest{Products: rows(1001)}, o...)
			return err
		},
		"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
		"trash bad id":  func(o ...transport.RequestOption) error { return s.Trash(ctx, "bad", o...) },
		"create many empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateMany(ctx, nil, o...)
			return err
		},
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[processing.UpdateRequest]{{ID: "bad", Request: fullUpdate}}, o...)
			return err
		},
		"update many empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[processing.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"delete many bad id": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, []string{"bad"}, o...)
			return err
		},
		"materials bad id":           page(s.GetMaterials, "bad", positions.PageOptions{}),
		"materials fields":           page(s.GetMaterials, testID, positions.PageOptions{Fields: "stock"}),
		"materials limit":            page(s.GetMaterials, testID, positions.PageOptions{Limit: 1001}),
		"products bad id":            page(s.GetProducts, "bad", positions.PageOptions{}),
		"products fields":            page(s.GetProducts, testID, positions.PageOptions{Fields: "stock"}),
		"products offset":            page(s.GetProducts, testID, positions.PageOptions{Offset: -1}),
		"create material no qty":     add(s.CreateMaterial, testID, assortment),
		"create material null qty":   add(s.CreateMaterial, testID, processing.PositionRequest{Assortment: model.Value(product), Quantity: model.Null[float64]()}),
		"create material no assort":  add(s.CreateMaterial, testID, quantity),
		"create material empty href": add(s.CreateMaterial, testID, processing.PositionRequest{Assortment: model.Value(model.Reference{}), Quantity: model.Value(1.0)}),
		"create material bad id":     add(s.CreateMaterial, "bad", valid),
		"create product no qty":      add(s.CreateProduct, testID, assortment),
		"create product no assort":   add(s.CreateProduct, testID, quantity),
		"create product bad id":      add(s.CreateProduct, "bad", valid),
		"get material bad": func(o ...transport.RequestOption) error {
			_, err := s.GetMaterial(ctx, testID, "bad", o...)
			return err
		},
		"get product bad": func(o ...transport.RequestOption) error {
			_, err := s.GetProduct(ctx, "bad", rowID, o...)
			return err
		},
		"update material empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMaterial(ctx, testID, rowID, processing.PositionRequest{}, o...)
			return err
		},
		"update product bad": func(o ...transport.RequestOption) error {
			_, err := s.UpdateProduct(ctx, testID, "bad", quantity, o...)
			return err
		},
		"delete material bad":         func(o ...transport.RequestOption) error { return s.DeleteMaterial(ctx, "bad", rowID, o...) },
		"delete product bad":          func(o ...transport.RequestOption) error { return s.DeleteProduct(ctx, testID, "bad", o...) },
		"delete many materials empty": func(o ...transport.RequestOption) error { return s.DeleteManyMaterials(ctx, testID, nil, o...) },
		"delete many materials 1001": func(o ...transport.RequestOption) error {
			return s.DeleteManyMaterials(ctx, testID, make([]string, 1001), o...)
		},
		"delete many materials bad doc": func(o ...transport.RequestOption) error {
			return s.DeleteManyMaterials(ctx, "bad", []string{rowID}, o...)
		},
		"delete many products bad row": func(o ...transport.RequestOption) error {
			return s.DeleteManyProducts(ctx, testID, []string{rowID, "bad"}, o...)
		},
		"delete many products empty": func(o ...transport.RequestOption) error { return s.DeleteManyProducts(ctx, testID, nil, o...) },
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "processing") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := processing.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestList(t *testing.T) {
	listtest.Run(t, "processing", read(t, "custom/processing/read.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[processing.Processing] {
		return service(c).List
	}, func(t *testing.T, got processing.Processing) {
		if got.ID != fixtureID || got.ProcessingSum != 10000 || got.ProcessingPlan == nil || got.Products == nil {
			t.Fatalf("row = %+v", got)
		}
	})
}
