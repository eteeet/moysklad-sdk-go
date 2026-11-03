package productiontask_test

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
	"github.com/eteeet/moysklad-sdk-go/entity/productiontask"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/metadata"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID     = "12121212-1212-1212-1212-121212121212"
	fixtureID  = "4d63c1b4-98e0-11ee-0a83-0045000025c7"
	rowID      = "e66de5ff-9a5c-11ee-0a83-0a2e00000910"
	productID  = "ef45d0b2-214e-11ee-c0a8-d00400000069"
	stageID    = "19070bfd-99d6-11ee-0a83-0a2e0000076b"
	materialID = "34f6344f-015e-11e6-9464-e4de0000006c"
	docPath    = "entity/productiontask/" + testID
	stagePath  = "entity/productionstage/" + stageID
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
	return model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/" + path, Type: model.MetaType(typ)}}
}

func refs(path, typ string) model.Field[[]model.Reference] {
	return model.Value([]model.Reference{ref(path, typ)})
}

func service(c *transport.Client) *productiontask.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := productiontask.NewService(c)
	return s
}

var (
	organization   = ref("organization/11111111-1111-1111-1111-111111111111", "organization")
	materialsStore = ref("store/22222222-2222-2222-2222-222222222222", "store")
	productsStore  = ref("store/33333333-3333-3333-3333-333333333333", "store")
	processingPlan = ref("processingplan/66666666-6666-6666-6666-666666666666", "processingplan")
	minimal        = productiontask.CreateRequest{Organization: organization, MaterialsStore: materialsStore, ProductsStore: productsStore}
	fullCreate     = productiontask.CreateRequest{
		Organization: organization, MaterialsStore: materialsStore, ProductsStore: productsStore,
		Name: model.Value("Тестовое Производственное задание"), Description: model.Value(""), Code: model.Value("test-code-123"), ExternalCode: model.Value("954102345"),
		Moment:                model.Value(documenttest.DateTime("2023-12-12 13:39:00")),
		DeliveryPlannedMoment: model.Value(documenttest.DateTime("2023-12-12T13:39:59.951+03:00")),
		ProductionStart:       model.Value(documenttest.DateTime("2023-12-13 09:00:00")),
		Applicable:            model.Value(true), Awaiting: model.Value(false), Reserve: model.Value(false), Shared: model.Value(false),
		Owner:      model.Value(ref("employee/44444444-4444-4444-4444-444444444444", "employee")),
		Group:      model.Value(ref("group/55555555-5555-5555-5555-555555555555", "group")),
		State:      model.Null[model.Reference](),
		Attributes: model.Value([]model.Attribute{}),
		ProductionRows: model.Value([]productiontask.ProductionRowRequest{{
			Name: model.Value("00001-1"), ExternalCode: model.Value("row-ext"), ProcessingPlan: model.Value(processingPlan),
			ProductionVariant: model.Value(ref("variant/77777777-7777-7777-7777-777777777777", "variant")), ProductionVolume: model.Value(10.0),
		}}),
		Products: model.Value([]productiontask.ProductRequest{{
			Assortment: model.Value(ref("product/88888888-8888-8888-8888-888888888888", "product")), PlanQuantity: model.Value(22.0),
		}}),
		CustomerOrders:         refs("customerorder/99999999-9999-9999-9999-999999999991", "customerorder"),
		Demands:                refs("demand/99999999-9999-9999-9999-999999999992", "demand"),
		EmissionOrders:         refs("emissionorder/99999999-9999-9999-9999-999999999993", "emissionorder"),
		InternalOrders:         refs("internalorder/99999999-9999-9999-9999-999999999994", "internalorder"),
		Moves:                  refs("move/99999999-9999-9999-9999-999999999995", "move"),
		ProductionTasks:        refs("productiontask/99999999-9999-9999-9999-999999999996", "productiontask"),
		ProductionTaskSupplies: refs("productiontask/99999999-9999-9999-9999-999999999997", "productiontask"),
		PurchaseOrders:         refs("purchaseorder/99999999-9999-9999-9999-999999999998", "purchaseorder"),
		Supplies:               refs("supply/99999999-9999-9999-9999-999999999999", "supply"),
	}
	fullUpdate = productiontask.UpdateRequest{
		Description:     model.Value("тестовое производственное задание на обновление"),
		ProductionStart: model.Value(documenttest.DateTime("2023-12-14 10:00:00")),
		Awaiting:        model.Value(false), Reserve: model.Value(false),
		State: model.Null[model.Reference](),
		ProductionRows: model.Value([]productiontask.ProductionRowRequest{
			{ID: model.Value(rowID), ProductionVolume: model.Value(15.0)},
			{ProcessingPlan: model.Value(processingPlan), ProductionVolume: model.Value(0.5)},
		}),
		Products: model.Value([]productiontask.ProductRequest{}),
		Supplies: model.Value([]model.Reference{}),
	}
)

func TestFixture(t *testing.T) {
	data := read(t, "custom/productiontask/read.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got productiontask.ProductionTask
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.AccountID != "042cf74c-9358-11ee-0a80-09e700000010" || got.Name != "Тестовое Производственное здание 1" ||
				got.Code != "test-code-123" || got.ExternalCode != "954102345" || got.Description != "тестовое производственное задание" ||
				!got.Shared || !got.Applicable || got.Printed || got.Published || got.Awaiting || got.Reserve || got.Meta.Type != "productiontask" {
				t.Fatalf("document = %+v", got)
			}
			if !time.Time(got.DeliveryPlannedMoment).Equal(time.Time(documenttest.DateTime("2023-12-12 13:39:00"))) ||
				!time.Time(got.Created).Equal(time.Time(documenttest.DateTime("2023-12-12T14:19:23.319+03:00"))) || got.Updated.IsZero() || got.Moment.IsZero() ||
				!got.ProductionStart.IsZero() || !got.ProductionEnd.IsZero() || !got.Deleted.IsZero() {
				t.Fatalf("dates = %+v", got)
			}
			if got.Organization == nil || got.Organization.Meta.Type != "organization" || got.MaterialsStore == nil || got.MaterialsStore.Meta.Type != "store" ||
				got.ProductsStore == nil || got.Owner == nil || got.Group == nil || got.State != nil || len(got.Attributes) != 1 || got.Attributes[0].Name != "дополнительное поле" {
				t.Fatalf("references = %+v", got)
			}
			if got.ProductionRows == nil || got.ProductionRows.Meta.Type != "productionrow" || got.ProductionRows.Meta.Size != 0 ||
				got.Products == nil || got.Products.Meta.Type != "productiontaskresult" || got.Files == nil || got.Files.Meta.Type != "files" {
				t.Fatalf("collections = %+v / %+v / %+v", got.ProductionRows, got.Products, got.Files)
			}
			// Поля таблицы, которых нет в примере ответа, и строки при expand.
			var sparse productiontask.ProductionTask
			if err := codec.Unmarshal([]byte(`{"id":"x","awaiting":true,"reserve":true,"productionStart":"2023-12-13 09:00:00.000","productionEnd":"2023-12-14 18:00:00.000",`+
				`"state":{"meta":{"type":"state"}},"customerOrders":[{"meta":{"type":"customerorder"}}],"demands":[{}],"emissionOrders":[{}],"internalOrders":[{}],"moves":[{}],`+
				`"productionTasks":[{}],"productionTaskSupplies":[{},{}],"purchaseOrders":[{}],"supplies":[{}],`+
				`"productionRows":{"meta":{"href":"h"},"rows":[{"id":"r","productionVolume":2.5}]},"products":{"meta":{"href":"p"},"rows":[{"id":"p","planQuantity":3}]},"future":1}`), &sparse); err != nil {
				t.Fatal(err)
			}
			if !sparse.Awaiting || !sparse.Reserve || sparse.ProductionStart.IsZero() || sparse.ProductionEnd.IsZero() || sparse.State == nil ||
				len(sparse.CustomerOrders) != 1 || len(sparse.Demands) != 1 || len(sparse.EmissionOrders) != 1 || len(sparse.InternalOrders) != 1 || len(sparse.Moves) != 1 ||
				len(sparse.ProductionTasks) != 1 || len(sparse.ProductionTaskSupplies) != 2 || len(sparse.PurchaseOrders) != 1 || len(sparse.Supplies) != 1 ||
				sparse.ProductionRows.Rows[0].ProductionVolume != 2.5 || sparse.Products.Rows[0].PlanQuantity != 3 || sparse.Organization != nil {
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
			{"create full", fullCreate, read(t, "custom/productiontask/create-full.golden.json")},
			{"create minimal", minimal, []byte(`{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},` +
				`"materialsStore":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/store/22222222-2222-2222-2222-222222222222","type":"store"}},` +
				`"productsStore":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/store/33333333-3333-3333-3333-333333333333","type":"store"}}}`)},
			{"update", fullUpdate, read(t, "custom/productiontask/update.golden.json")},
			// Явный false флагов уходит, незаданные поля — нет.
			{"update false flags", productiontask.UpdateRequest{Applicable: model.Value(false), Awaiting: model.Value(false), Reserve: model.Value(false), Shared: model.Value(false)},
				[]byte(`{"applicable":false,"awaiting":false,"reserve":false,"shared":false}`)},
			{"update true flags", productiontask.UpdateRequest{Awaiting: model.Value(true), Reserve: model.Value(true)}, []byte(`{"awaiting":true,"reserve":true}`)},
			{"production row update", productiontask.ProductionRowUpdate{ProductionVolume: model.Value(0.0)}, []byte(`{"productionVolume":0}`)},
			{"product update", productiontask.ProductRequest{PlanQuantity: model.Value(0.5)}, []byte(`{"planQuantity":0.5}`)},
			{"stage update", productiontask.StageUpdateRequest{
				EnableHourAccounting: model.Value(false), ProcessingUnitCost: model.Value(70.0), LabourUnitCost: model.Value(30.5), StandardHourCost: model.Value(0.0),
				StandardHourUnit: model.Value(43.5), PlannedEndDate: model.Value(documenttest.DateTime("2023-12-20 18:00:00")),
				Materials: model.Value([]productiontask.MaterialRequest{
					{ID: model.Value(materialID), PlanQuantity: model.Value(3.0)},
					{Assortment: model.Value(ref("product/57c9a841-99d8-11ee-0a83-0a2e00000773", "product")), PlanQuantity: model.Value(2.0)},
				}),
			}, read(t, "custom/productiontask/stage-update.golden.json")},
			{"stage hour accounting on", productiontask.StageUpdateRequest{EnableHourAccounting: model.Value(true)}, []byte(`{"enableHourAccounting":true}`)},
			{"stage hour accounting off", productiontask.StageUpdateRequest{EnableHourAccounting: model.Value(false)}, []byte(`{"enableHourAccounting":false}`)},
			{"stage unset", productiontask.StageUpdateRequest{}, []byte(`{}`)},
			{"material", productiontask.MaterialRequest{Assortment: model.Value(ref("product/9656117b-b717-11ec-0c80-0bba0006dcde", "product")), PlanQuantity: model.Value(2.0)},
				[]byte(`{"assortment":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/9656117b-b717-11ec-0c80-0bba0006dcde","type":"product"}},"planQuantity":2}`)},
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
	readOnly := map[reflect.Type][]string{
		reflect.TypeFor[productiontask.CreateRequest]():        {"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "ProductionEnd", "Files"},
		reflect.TypeFor[productiontask.UpdateRequest]():        {"ID", "Meta", "AccountID", "Created", "Updated", "Deleted", "Printed", "Published", "ProductionEnd", "Files"},
		reflect.TypeFor[productiontask.ProductionRowRequest](): {"Meta", "AccountID", "Updated"},
		reflect.TypeFor[productiontask.ProductRequest]():       {"Meta", "AccountID", "ProductionRow"},
		reflect.TypeFor[productiontask.MaterialRequest]():      {"Meta", "AccountID"},
		reflect.TypeFor[productiontask.StageUpdateRequest](): {"ID", "Meta", "AccountID", "Files", "Instruction", "MaterialStore", "OrderingPosition", "Stage", "ProductionRow",
			"TotalQuantity", "CompletedQuantity", "AvailableQuantity", "BlockedQuantity", "SkippedQuantity"},
	}
	for typ, names := range readOnly {
		for _, name := range names {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
	if typ := reflect.TypeFor[productiontask.ProductionRowUpdate](); typ.NumField() != 1 || typ.Field(0).Name != "ProductionVolume" {
		t.Error("production row update must write only productionVolume")
	}
}

func TestHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и флаги — только из ответа.
	const response = `{"id":"server-id","name":"server-name","awaiting":true,"reserve":true,"productionRows":{"meta":{"href":"https://example.test/productionrows","size":1}}}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*productiontask.ProductionTask, error)
	}{
		{"get", "GET", docPath, nil, func(ctx context.Context, c *transport.Client) (*productiontask.ProductionTask, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/productiontask", read(t, "custom/productiontask/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*productiontask.ProductionTask, error) {
			return service(c).Create(ctx, fullCreate)
		}},
		{"update", "PUT", docPath, read(t, "custom/productiontask/update.golden.json"), func(ctx context.Context, c *transport.Client) (*productiontask.ProductionTask, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *productiontask.ProductionTask) {
						if got.ID != "server-id" || got.Name != "server-name" || !got.Awaiting || !got.Reserve || got.ProductionRows.Meta.Size != 1 || got.ProductionRows.Rows != nil {
							t.Fatalf("document = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Методы без модели в ответе и массовые: один запрос с точным телом, без повторов при 503.
func TestWrites(t *testing.T) {
	const base = "{base}/entity/productiontask"
	docMeta := `{"href":"` + base + `/` + testID + `","metadataHref":"` + base + `/metadata","type":"productiontask","mediaType":"application/json"}`
	assortmentJSON := `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/88888888-8888-8888-8888-888888888888","type":"product"}}`
	for _, tc := range []struct {
		name, method, path, body, response string
		empty                              bool // ответ без тела: 204 тоже успех
		invoke                             func(context.Context, *transport.Client) error
	}{
		{"delete", "DELETE", docPath, "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Delete(ctx, testID) }},
		{"trash", "POST", docPath + "/trash", "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Trash(ctx, testID) }},
		{"create many", "POST", "entity/productiontask", `[{"organization":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/organization/11111111-1111-1111-1111-111111111111","type":"organization"}},` +
			`"materialsStore":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/store/22222222-2222-2222-2222-222222222222","type":"store"}},` +
			`"productsStore":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/store/33333333-3333-3333-3333-333333333333","type":"store"}}}]`, `[{"id":"x"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateMany(ctx, []productiontask.CreateRequest{minimal})
				if err == nil && r[0].Value.ID != "x" {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"update many", "POST", "entity/productiontask", `[{"meta":` + docMeta + `,"reserve":false}]`, `[{"id":"x"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).UpdateMany(ctx, []batch.Update[productiontask.UpdateRequest]{{ID: testID, Request: productiontask.UpdateRequest{Reserve: model.Value(false)}}})
				if err == nil && r[0].Value.ID != "x" {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete many", "POST", "entity/productiontask/delete", `[{"meta":` + docMeta + `}]`, `[{"info":"ok"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).DeleteMany(ctx, []string{testID})
				if err == nil && r[0].Info != "ok" {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete production row", "DELETE", docPath + "/productionrows/" + rowID, "", "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeleteProductionRow(ctx, testID, rowID)
		}},
		{"create products", "POST", docPath + "/products", `[{"assortment":` + assortmentJSON + `,"planQuantity":22}]`, `[{"id":"` + productID + `","planQuantity":22}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateProducts(ctx, testID, []productiontask.ProductRequest{{Assortment: model.Value(ref("product/88888888-8888-8888-8888-888888888888", "product")), PlanQuantity: model.Value(22.0)}})
				if err == nil && (r[0].Value.ID != productID || r[0].Value.PlanQuantity != 22) {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete product", "DELETE", docPath + "/products/" + productID, "", "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeleteProduct(ctx, testID, productID)
		}},
		{"delete many products", "POST", docPath + "/products/delete", `[{"meta":{"href":"` + base + `/` + testID + `/products/` + productID + `","type":"productiontaskresult","mediaType":"application/json"}}]`, "", true,
			func(ctx context.Context, c *transport.Client) error {
				return service(c).DeleteManyProducts(ctx, testID, []string{productID})
			}},
		{"delete material", "DELETE", stagePath + "/materials/" + materialID, "", "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeleteMaterial(ctx, stageID, materialID)
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
					if r.Method != tc.method || r.URL.Path != "/api/remap/1.2/"+tc.path || r.URL.RawQuery != "" {
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
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "production") ||
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

func TestList(t *testing.T) {
	listtest.Run(t, "productiontask", read(t, "custom/productiontask/read.json"), "", "isDeleted", func(c *transport.Client) listtest.ListFunc[productiontask.ProductionTask] {
		return service(c).List
	}, func(t *testing.T, got productiontask.ProductionTask) {
		if got.ID != fixtureID || got.Products == nil || got.Products.Meta.Type != "productiontaskresult" {
			t.Fatalf("row = %+v", got)
		}
	})
}

// Метаданные — общий пакет metadata: пример ответа источника разбирается для типа productiontask.
func TestMetadata(t *testing.T) {
	data := read(t, "custom/productiontask/metadata.json")
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != http.MethodGet || r.URL.Path != "/api/remap/1.2/entity/productiontask/metadata" {
			t.Errorf("request = %s %s", r.Method, r.URL)
		}
		w.Write(data)
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2"})
	if err != nil {
		t.Fatal(err)
	}
	s, err := metadata.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	got, err := s.Get(t.Context(), "productiontask")
	if err != nil || calls.Load() != 1 {
		t.Fatalf("metadata = %+v, err = %v, calls = %d", got, err, calls.Load())
	}
	if got.CreateShared || got.Meta.Href != "https://api.moysklad.ru/api/remap/1.2/entity/productiontask/metadata" ||
		got.Attributes == nil || got.Attributes.Meta == nil || got.Attributes.Meta.Size != 2 || got.Attributes.Meta.Type != "attributemetadata" {
		t.Fatalf("metadata = %+v", got)
	}
}
