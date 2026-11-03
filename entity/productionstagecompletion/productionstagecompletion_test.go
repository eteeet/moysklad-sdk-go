package productionstagecompletion_test

import (
	"compress/gzip"
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
	"github.com/eteeet/moysklad-sdk-go/entity/productionstagecompletion"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID    = "7944ef04-f831-11e5-7a69-971500188b19"
	fixtureID = "01ff6808-95de-11ee-0a81-072300000136"
	rowID     = "34f6344f-015e-11e6-9464-e4de0000006c"
	docPath   = "entity/productionstagecompletion/" + testID
	stageHref = "https://api.moysklad.ru/api/remap/1.2/entity/productionstage/3130f7df-660f-11ee-c0a8-100c00000139"
	productH  = "https://api.moysklad.ru/api/remap/1.2/entity/product/9656117b-b717-11ec-0c80-0bba0006dcde"
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

func service(c *transport.Client) *productionstagecompletion.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := productionstagecompletion.NewService(c)
	return s
}

var (
	stage      = model.Reference{Meta: model.Meta{Href: stageHref, Type: "productionstage"}}
	product    = model.Reference{Meta: model.Meta{Href: productH, Type: "product"}}
	productRaw = `{"meta":{"href":"` + productH + `","type":"product"}}`
	minimal    = productionstagecompletion.CreateRequest{ProductionStage: stage, ProductionVolume: 5}
	fullCreate = productionstagecompletion.CreateRequest{
		ProductionStage: stage, ProductionVolume: 5,
		Name: model.Value("000034"), Description: model.Value(""), ExternalCode: model.Value("ext-1"),
		Moment: model.Value(documenttest.DateTime("2016-04-19 13:50:24")), Shared: model.Value(false),
		Defect: model.Value(true), EnableHourAccounting: model.Value(false),
		Owner:              model.Value(ref("entity/employee/faba7f37-2e58-11e6-8a84-bae500000028", "employee")),
		Group:              model.Value(ref("entity/group/f97aa1fb-2e58-11e6-8a84-bae500000002", "group")),
		Performer:          model.Value(ref("entity/counterparty/c8c0f0de-b6c6-11ef-ac12-00130000009e", "counterparty")),
		Service:            model.Value(ref("entity/service/18616ba6-b6c7-11ef-ac12-001300000103", "service")),
		ProcessingUnitCost: model.Value(10.0), LabourUnitCost: model.Value(20.5),
		StandardHourCost: model.Value(0.0), StandardHourUnit: model.Value(30.0),
	}
	fullUpdate = productionstagecompletion.UpdateRequest{
		ProductionVolume: model.Value(22.0), Name: model.Value("000034"), Description: model.Null[string](),
		Shared: model.Value(true), EnableHourAccounting: model.Value(true),
		Performer: model.Null[model.Reference](), Service: model.Null[model.Reference](),
		StandardHourCost: model.Value(1500.0), StandardHourUnit: model.Value(0.25),
	}
)

func TestFixture(t *testing.T) {
	data := read(t, "custom/productionstagecompletion/read.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got productionstagecompletion.ProductionStageCompletion
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.AccountID != "a67c68a3-95dd-11ee-0a83-071a00000002" || got.Name != "000034" || got.Description != "Комментарий по выполнению" ||
				got.ExternalCode != "EnuNEmG2jyUF4t9tgPQk72" || !got.Defect || got.EnableHourAccounting || got.Shared || got.Meta.Type != "productionstagecompletion" {
				t.Fatalf("document = %+v", got)
			}
			if got.ProductionVolume != 5 || got.ProcessingUnitCost != 10 || got.LabourUnitCost != 20 || got.StandardHourUnit != 30 || got.StandardHourCost != 0 {
				t.Fatalf("numbers = %+v", got)
			}
			if !time.Time(got.Moment).Equal(time.Time(documenttest.DateTime("2016-04-19 13:50:24"))) || got.Created.IsZero() || got.Updated.IsZero() {
				t.Fatalf("dates = %+v", got)
			}
			if got.Owner == nil || got.Group == nil || got.Performer == nil || got.Performer.Meta.Type != "employee" || got.Service != nil ||
				got.ProductionStage == nil || got.ProductionStage.Meta.Type != "productionstage" {
				t.Fatalf("references = %+v", got)
			}
			if got.Materials == nil || got.Materials.Meta.Type != "productionstagecompletionmaterial" || got.Materials.Meta.Size != 1 || got.Materials.Rows != nil ||
				got.Products == nil || got.Products.Meta.Type != "productionstagecompletionresult" {
				t.Fatalf("collections = %+v / %+v", got.Materials, got.Products)
			}
			// Поля таблицы, которых нет в примере ответа, и строки при expand.
			var sparse productionstagecompletion.ProductionStageCompletion
			if err := codec.Unmarshal([]byte(`{"id":"x","applicable":true,"service":{"meta":{"type":"service"}},"productionVolume":1.1235,`+
				`"materials":{"meta":{"href":"m"},"rows":[{"id":"m","consumedQuantity":0.5,"assortment":{"meta":{"type":"product"}},"things":["A","B"]}]},`+
				`"products":{"meta":{"href":"p"},"rows":[{"id":"p","producedQuantity":0.7,"things":[]}]},"standardHour":0.0,"future":1}`), &sparse); err != nil {
				t.Fatal(err)
			}
			m, p := sparse.Materials.Rows[0], sparse.Products.Rows[0]
			if !sparse.Applicable || sparse.Service == nil || sparse.ProductionVolume != 1.1235 || m.ConsumedQuantity != 0.5 || m.Assortment == nil ||
				!reflect.DeepEqual(m.Things, []string{"A", "B"}) || p.ProducedQuantity != 0.7 || p.Things == nil || len(p.Things) != 0 ||
				sparse.Owner != nil || !sparse.Created.IsZero() {
				t.Fatalf("sparse = %+v", sparse)
			}
		})
	}
}

// Деньги — копейки float64: 0, 0.0, дробь и большое значение читаются без потерь обоими кодеками.
// Количества (объём, нормо-часы, количества строк) — дробные числа без пересчёта, не деньги.
func TestNumbersJSON(t *testing.T) {
	values := map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "123456789.01": 123456789.01}
	for _, codec := range documenttest.Codecs {
		for _, field := range []string{"processingUnitCost", "labourUnitCost", "standardHourCost", "productionVolume", "standardHourUnit"} {
			for raw, want := range values {
				var got productionstagecompletion.ProductionStageCompletion
				if err := codec.Unmarshal([]byte(`{"`+field+`":`+raw+`}`), &got); err != nil {
					t.Fatal(err)
				}
				value := map[string]float64{"processingUnitCost": got.ProcessingUnitCost, "labourUnitCost": got.LabourUnitCost,
					"standardHourCost": got.StandardHourCost, "productionVolume": got.ProductionVolume, "standardHourUnit": got.StandardHourUnit}[field]
				if value != want {
					t.Fatalf("%s: %s=%s: got %v", codec.Name, field, raw, value)
				}
			}
		}
		var m productionstagecompletion.Material
		var p productionstagecompletion.Product
		if err := codec.Unmarshal([]byte(`{"consumedQuantity":0.001}`), &m); err != nil || m.ConsumedQuantity != 0.001 {
			t.Fatalf("%s: material = %+v, err = %v", codec.Name, m, err)
		}
		if err := codec.Unmarshal([]byte(`{"producedQuantity":2.7}`), &p); err != nil || p.ProducedQuantity != 2.7 {
			t.Fatalf("%s: product = %+v, err = %v", codec.Name, p, err)
		}
		// Дробное количество в запросе уходит как есть, без перевода в копейки.
		data, err := codec.Marshal(productionstagecompletion.UpdateRequest{ProductionVolume: model.Value(0.125)})
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, data, []byte(`{"productionVolume":0.125}`))
	}
}

func TestRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", fullCreate, read(t, "custom/productionstagecompletion/create-full.golden.json")},
			{"create minimal", minimal, []byte(`{"productionStage":{"meta":{"href":"` + stageHref + `","type":"productionstage"}},"productionVolume":5}`)},
			{"create zero volume", productionstagecompletion.CreateRequest{ProductionStage: stage}, []byte(`{"productionStage":{"meta":{"href":"` + stageHref + `","type":"productionstage"}},"productionVolume":0}`)},
			{"update", fullUpdate, read(t, "custom/productionstagecompletion/update.golden.json")},
			{"update false", productionstagecompletion.UpdateRequest{Shared: model.Value(false), LabourUnitCost: model.Value(0.0)}, []byte(`{"shared":false,"labourUnitCost":0}`)},
			{"material add", productionstagecompletion.MaterialCreateRequest{Assortment: product, ConsumedQuantity: model.Value(2.0)}, []byte(`{"assortment":` + productRaw + `,"consumedQuantity":2}`)},
			{"material add things", productionstagecompletion.MaterialCreateRequest{Assortment: product, Things: model.Value([]string{"F564X056", "F564X057"})}, []byte(`{"assortment":` + productRaw + `,"things":["F564X056","F564X057"]}`)},
			{"material update", productionstagecompletion.MaterialUpdateRequest{ConsumedQuantity: model.Value(0.0), Assortment: model.Value(product)}, []byte(`{"assortment":` + productRaw + `,"consumedQuantity":0}`)},
			{"product update", productionstagecompletion.ProductUpdateRequest{ProducedQuantity: model.Value(2.7), Things: model.Value([]string{})}, []byte(`{"producedQuantity":2.7,"things":[]}`)},
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
	// Только для чтения — нигде; productionStage и defect после создания не меняются.
	for _, typ := range []reflect.Type{reflect.TypeFor[productionstagecompletion.CreateRequest](), reflect.TypeFor[productionstagecompletion.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Created", "Updated", "Applicable", "Materials", "Products"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
	for _, name := range []string{"ProductionStage", "Defect"} {
		if _, ok := reflect.TypeFor[productionstagecompletion.UpdateRequest]().FieldByName(name); ok {
			t.Errorf("UpdateRequest contains create-only %s", name)
		}
	}
	for _, typ := range []reflect.Type{reflect.TypeFor[productionstagecompletion.MaterialCreateRequest](), reflect.TypeFor[productionstagecompletion.MaterialUpdateRequest](), reflect.TypeFor[productionstagecompletion.ProductUpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func TestHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и объём — только из ответа.
	const response = `{"id":"server-id","name":"server-name","productionVolume":1.6,"labourUnitCost":2500,"materials":{"meta":{"href":"https://example.test/materials","size":2}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*productionstagecompletion.ProductionStageCompletion, error)
	}{
		{"get", "GET", docPath, nil, func(ctx context.Context, c *transport.Client) (*productionstagecompletion.ProductionStageCompletion, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/productionstagecompletion", read(t, "custom/productionstagecompletion/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*productionstagecompletion.ProductionStageCompletion, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", docPath, read(t, "custom/productionstagecompletion/update.golden.json"), func(ctx context.Context, c *transport.Client) (*productionstagecompletion.ProductionStageCompletion, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *productionstagecompletion.ProductionStageCompletion) {
						if got.ID != "server-id" || got.Name != "server-name" || got.ProductionVolume != 1.6 || got.LabourUnitCost != 2500 ||
							got.Materials == nil || got.Materials.Meta.Size != 2 || got.Materials.Rows != nil {
							t.Fatalf("document = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Материалы и продукты: изменение и добавление (примеры источника), ответ — объект.
func TestRowsHTTP(t *testing.T) {
	const (
		material  = `{"meta":{"type":"productionstagecompletionmaterial"},"id":"` + rowID + `","consumedQuantity":3.0,"assortment":{"meta":{"type":"product"}},"things":["F564X056"]}`
		resultRow = `{"meta":{"type":"productionstagecompletionresult"},"id":"` + rowID + `","producedQuantity":2.7,"assortment":{"meta":{"type":"product"}}}`
	)
	checkMaterial := func(got *productionstagecompletion.Material) {
		if got.ID != rowID || got.ConsumedQuantity != 3 || got.Assortment == nil || len(got.Things) != 1 || got.Meta.Type != "productionstagecompletionmaterial" {
			t.Fatalf("material = %+v", got)
		}
	}
	checkProduct := func(got *productionstagecompletion.Product) {
		if got.ID != rowID || got.ProducedQuantity != 2.7 || got.Assortment == nil || got.Meta.Type != "productionstagecompletionresult" {
			t.Fatalf("product = %+v", got)
		}
	}
	t.Run("add material", func(t *testing.T) {
		httpCases(t, "POST", docPath+"/materials", `{"assortment":`+productRaw+`,"consumedQuantity":2}`, material,
			func(ctx context.Context, c *transport.Client) (*productionstagecompletion.Material, error) {
				return service(c).CreateMaterial(ctx, testID, productionstagecompletion.MaterialCreateRequest{Assortment: product, ConsumedQuantity: model.Value(2.0)})
			}, checkMaterial)
	})
	t.Run("update material", func(t *testing.T) {
		httpCases(t, "PUT", docPath+"/materials/"+rowID, `{"consumedQuantity":3}`, material,
			func(ctx context.Context, c *transport.Client) (*productionstagecompletion.Material, error) {
				return service(c).UpdateMaterial(ctx, testID, rowID, productionstagecompletion.MaterialUpdateRequest{ConsumedQuantity: model.Value(3.0)})
			}, checkMaterial)
	})
	t.Run("update product", func(t *testing.T) {
		httpCases(t, "PUT", docPath+"/products/"+rowID, `{"producedQuantity":2.7}`, resultRow,
			func(ctx context.Context, c *transport.Client) (*productionstagecompletion.Product, error) {
				return service(c).UpdateProduct(ctx, testID, rowID, productionstagecompletion.ProductUpdateRequest{ProducedQuantity: model.Value(2.7)})
			}, checkProduct)
	})
	t.Run("get materials", func(t *testing.T) {
		pageCases(t, docPath+"/materials", "limit=2&offset=7", `{"meta":{"href":"h","type":"productionstagecompletionmaterial","size":1},"rows":[`+material+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[productionstagecompletion.Material], error) {
				return service(c).GetMaterials(ctx, testID, query.PageOptions{Limit: 2, Offset: 7})
			}, func(got *model.Collection[productionstagecompletion.Material]) { checkMaterial(&got.Rows[0]) })
	})
	t.Run("get products", func(t *testing.T) {
		pageCases(t, docPath+"/products", "limit=1000&offset=0", `{"meta":{"href":"h","type":"productionstagecompletionresult","size":1},"rows":[`+resultRow+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[productionstagecompletion.Product], error) {
				return service(c).GetProducts(ctx, testID, query.PageOptions{})
			}, func(got *model.Collection[productionstagecompletion.Product]) { checkProduct(&got.Rows[0]) })
	})
}

func httpCases[T any](t *testing.T, method, path, body, response string, invoke func(context.Context, *transport.Client) (*T, error), check func(*T)) {
	t.Helper()
	for _, tc := range documenttest.ResponseCases(response) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, method, path, []byte(body), invoke, check)
		})
	}
	documenttest.LocalErrors(t, invoke)
}

// pageCases — ответы страницы строк: успех plain/gzip, ошибки API, неверный статус и тело; ровно
// один GET с точным query. Операция страницы — «list …», поэтому documenttest.HTTP не подходит.
func pageCases[T any](t *testing.T, path, rawQuery, response string, invoke func(context.Context, *transport.Client) (*model.Collection[T], error), check func(*model.Collection[T])) {
	t.Helper()
	for _, tc := range documenttest.ResponseCases(response) {
		t.Run(tc.Name, func(t *testing.T) {
			var calls atomic.Int32
			server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodGet || r.URL.Path != "/api/remap/1.2/"+path || r.URL.RawQuery != rawQuery || r.Header.Get("Accept-Encoding") != "gzip" {
					t.Errorf("request = %s %s", r.Method, r.URL)
				}
				if tc.Gzip {
					w.Header().Set("Content-Encoding", "gzip")
					w.WriteHeader(tc.Status)
					zw := gzip.NewWriter(w)
					io.WriteString(zw, tc.Body)
					zw.Close()
					return
				}
				w.WriteHeader(tc.Status)
				io.WriteString(w, tc.Body)
			}))
			defer server.Close()
			client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2"})
			if err != nil {
				t.Fatal(err)
			}
			got, err := invoke(t.Context(), client)
			// Тело с объектом без rows («empty id», «missing id») — ошибка для страницы.
			if tc.Success {
				if err != nil || got == nil {
					t.Fatalf("got = %+v, err = %v", got, err)
				}
				check(got)
			} else {
				var apiErr *sdkerrors.APIError
				if err == nil || got != nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "productionstagecompletion "+testID) ||
					errors.As(err, &apiErr) != tc.API {
					t.Fatalf("got = %+v, err = %v", got, err)
				}
			}
			if calls.Load() != 1 {
				t.Fatalf("calls = %d", calls.Load())
			}
		})
	}
	documenttest.LocalErrors(t, invoke)
}

// Методы без модели в ответе и массовые: один запрос с точным телом, без повторов при 503.
func TestWrites(t *testing.T) {
	const base = "{base}/entity/productionstagecompletion"
	docMeta := `{"href":"` + base + `/` + testID + `","metadataHref":"` + base + `/metadata","type":"productionstagecompletion","mediaType":"application/json"}`
	for _, tc := range []struct {
		name, method, path, body, response string
		empty                              bool // ответ без тела: 204 тоже успех
		invoke                             func(context.Context, *transport.Client) error
	}{
		{"delete", "DELETE", docPath, "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Delete(ctx, testID) }},
		{"trash", "POST", docPath + "/trash", "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Trash(ctx, testID) }},
		{"create many", "POST", "entity/productionstagecompletion", `[{"productionStage":{"meta":{"href":"` + stageHref + `","type":"productionstage"}},"productionVolume":15,"name":"000033"}]`, `[{"id":"x","productionVolume":15}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				input := minimal
				input.ProductionVolume, input.Name = 15, model.Value("000033")
				r, err := service(c).CreateMany(ctx, []productionstagecompletion.CreateRequest{input})
				if err == nil && (r[0].Value.ID != "x" || r[0].Value.ProductionVolume != 15) {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"update many", "POST", "entity/productionstagecompletion", `[{"meta":` + docMeta + `,"name":"000034","productionVolume":5}]`, `[{"id":"y"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).UpdateMany(ctx, []batch.Update[productionstagecompletion.UpdateRequest]{{ID: testID, Request: productionstagecompletion.UpdateRequest{Name: model.Value("000034"), ProductionVolume: model.Value(5.0)}}})
				if err == nil && r[0].Value.ID != "y" {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete many", "POST", "entity/productionstagecompletion/delete", `[{"meta":` + docMeta + `},{"meta":` + strings.Replace(docMeta, testID, fixtureID, 1) + `}]`, `[{"info":"a"},{"info":"b"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				// Порядок ids сохраняется: сервер удаляет в указанном порядке.
				r, err := service(c).DeleteMany(ctx, []string{testID, fixtureID})
				if err == nil && (r[0].Info != "a" || r[1].Info != "b") {
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
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "productionstagecompletion") ||
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

// Ошибки входа — до HTTP, ResponseInfo сброшен, ровно один префикс moysklad:.
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
	// Нулевой объём — не ошибка SDK: значение проверяет сервер.
	if _, err := s.Create(ctx, productionstagecompletion.CreateRequest{ProductionStage: stage}); err != nil || calls.Load() != 1 {
		t.Fatalf("zero volume: err = %v, calls = %d", err, calls.Load())
	}
	calls.Store(0)
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id": func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":   func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no stage": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, productionstagecompletion.CreateRequest{ProductionVolume: 5, Name: model.Value("x")}, o...)
			return err
		},
		"create stage without href": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, productionstagecompletion.CreateRequest{ProductionStage: model.Reference{Meta: model.Meta{Type: "productionstage"}}, ProductionVolume: 5}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, productionstagecompletion.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", fullUpdate, o...); return err },
		"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
		"trash bad id":  func(o ...transport.RequestOption) error { return s.Trash(ctx, "bad", o...) },
		"create many empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateMany(ctx, nil, o...)
			return err
		},
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[productionstagecompletion.UpdateRequest]{{ID: "bad", Request: fullUpdate}}, o...)
			return err
		},
		"update many empty request": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[productionstagecompletion.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"delete many bad id": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, []string{"bad"}, o...)
			return err
		},
		"materials bad id": func(o ...transport.RequestOption) error {
			_, err := s.GetMaterials(ctx, "bad", query.PageOptions{}, o...)
			return err
		},
		"materials fields": func(o ...transport.RequestOption) error {
			_, err := s.GetMaterials(ctx, testID, query.PageOptions{Fields: "stock"}, o...)
			return err
		},
		"materials limit": func(o ...transport.RequestOption) error {
			_, err := s.GetMaterials(ctx, testID, query.PageOptions{Limit: 1001}, o...)
			return err
		},
		"add material no assortment": func(o ...transport.RequestOption) error {
			_, err := s.CreateMaterial(ctx, testID, productionstagecompletion.MaterialCreateRequest{ConsumedQuantity: model.Value(1.0)}, o...)
			return err
		},
		"add material bad id": func(o ...transport.RequestOption) error {
			_, err := s.CreateMaterial(ctx, "bad", productionstagecompletion.MaterialCreateRequest{Assortment: product}, o...)
			return err
		},
		"update material empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMaterial(ctx, testID, rowID, productionstagecompletion.MaterialUpdateRequest{}, o...)
			return err
		},
		"update material bad row": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMaterial(ctx, testID, "../x", productionstagecompletion.MaterialUpdateRequest{ConsumedQuantity: model.Value(1.0)}, o...)
			return err
		},
		"products bad id": func(o ...transport.RequestOption) error {
			_, err := s.GetProducts(ctx, "bad", query.PageOptions{}, o...)
			return err
		},
		"products offset": func(o ...transport.RequestOption) error {
			_, err := s.GetProducts(ctx, testID, query.PageOptions{Offset: -1}, o...)
			return err
		},
		"update product empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateProduct(ctx, testID, rowID, productionstagecompletion.ProductUpdateRequest{}, o...)
			return err
		},
		"update product bad doc": func(o ...transport.RequestOption) error {
			_, err := s.UpdateProduct(ctx, "bad", rowID, productionstagecompletion.ProductUpdateRequest{ProducedQuantity: model.Value(1.0)}, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "productionstagecompletion") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := productionstagecompletion.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestList(t *testing.T) {
	listtest.Run(t, "productionstagecompletion", read(t, "custom/productionstagecompletion/read.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[productionstagecompletion.ProductionStageCompletion] {
		return service(c).List
	}, func(t *testing.T, got productionstagecompletion.ProductionStageCompletion) {
		if got.ID != fixtureID || got.ProductionStage == nil || got.ProductionVolume != 5 {
			t.Fatalf("row = %+v", got)
		}
	})
}
