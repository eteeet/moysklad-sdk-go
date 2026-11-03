package processingplan_test

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/processingplan"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const fixtureID = "1a18770e-ad9a-11e6-5bed-427b00000064"

var (
	process    = "entity/processingprocess/25980515-923e-11ed-c0a8-30040000002d"
	position   = ref(process+"/positions/d8da461d-bbf9-11ed-ac12-0010000000c1", "processingprocessposition")
	planProd   = ref("entity/product/3b89c8e1-6fae-11f1-c5b5-417f00000153", "product")
	minimal    = processingplan.CreateRequest{Name: "Example", Products: []processingplan.ProductRequest{{Assortment: model.Value(planProd), Quantity: model.Value(1.0)}}}
	fullCreate = processingplan.CreateRequest{
		Name:     "plan_name",
		Products: []processingplan.ProductRequest{{Assortment: model.Value(planProd), Quantity: model.Value(1.0)}},
		Code:     model.Value("c1"), ExternalCode: model.Value(""), Archived: model.Value(false), Shared: model.Value(true),
		Cost:              model.Value(1000.0),
		Owner:             model.Value(ref("entity/employee/d5ad957e-91f1-11e6-5bed-427b0000002a", "employee")),
		Group:             model.Value(ref("entity/group/d55da707-91f1-11e6-5bed-427b00000001", "group")),
		Parent:            model.Null[model.Reference](),
		ProcessingProcess: model.Value(ref(process, "processingprocess")),
		Stages: model.Value([]processingplan.StageRequest{{
			EnableHourAccounting: model.Value(false), LabourCost: model.Value(100.0), StandardHour: model.Value(2.5), ProcessingProcessPosition: model.Value(position),
		}}),
		Materials: model.Value([]processingplan.MaterialRequest{{
			Assortment: model.Value(ref("entity/product/4d2cdff9-6fae-11f1-c5b5-417f00000176", "product")), Quantity: model.Value(1.0), ProcessingProcessPosition: model.Value(position),
		}}),
		ParametricMaterials: model.Value([]processingplan.ParametricMaterialRequest{{
			ParamFeatureCharacteristic: model.Value(ref("entity/variant/metadata/characteristics/58fe2fba-6fae-11f1-c5b5-417f00000183", "attributemetadata")),
			ParamFeatures: model.Value([]processingplan.ParamFeatureRequest{
				{Assortment: ref("entity/product/48a9a623-6fae-11f1-c5b5-417f00000168", "product")},
				{Assortment: ref("entity/product/4af31bf0-6fae-11f1-c5b5-417f0000016f", "product")},
			}),
			ParamQuantityCharacteristic: model.Value(ref("entity/variant/metadata/characteristics/58fe3a69-6fae-11f1-c5b5-417f00000184", "attributemetadata")),
			ParamQuantities:             model.Value([]processingplan.ParamQuantityRequest{{Quantity: 2}, {Quantity: 1}}),
		}}),
	}
	fullUpdate = processingplan.UpdateRequest{
		Name: model.Value("Техкарта"), Cost: model.Value(100000.0), Archived: model.Value(true), Owner: model.Null[model.Reference](),
		Products: model.Value([]processingplan.ProductRequest{
			{ID: model.Value("120b4591-b0bd-11e6-5bed-427b00000001"), Quantity: model.Value(2.0)},
			{Product: model.Value(productRef), Quantity: model.Value(0.5)},
		}),
		Stages:              model.Value([]processingplan.StageRequest{{ID: model.Value("4ca1e0fe-c161-11ed-ac1c-000b00000006"), EnableHourAccounting: model.Value(true)}}),
		Materials:           model.Value([]processingplan.MaterialRequest{}),
		ParametricMaterials: model.Value([]processingplan.ParametricMaterialRequest{{ID: model.Value("d6f43c04-5504-11f1-3593-85de000001ad"), Quantity: model.Value(3.5)}}),
	}
)

// Пример ответа «Получить Техкарту» источника.
func TestFixture(t *testing.T) {
	data := read(t, "custom/processingplan/read.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got processingplan.ProcessingPlan
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != fixtureID || got.AccountID != "d55cbfba-91f1-11e6-5bed-427b00000000" || got.Name != "Техкарточка" || got.ExternalCode != "4geOQkq5h7d5w1-tUATmt3" ||
				got.Archived || !got.Shared || got.PathName != "" || got.Cost != 1000 || got.CostDistributionType != "BY_PRICE" || got.Updated.IsZero() || got.Meta.Type != "processingplan" {
				t.Fatalf("plan = %+v", got)
			}
			if got.Owner == nil || got.Group == nil || got.ProcessingProcess == nil || got.ProcessingProcess.Meta.Type != "processingprocess" || got.Parent != nil {
				t.Fatalf("references = %+v", got)
			}
			if got.Stages == nil || got.Stages.Meta.Type != "processingplanstages" || got.Stages.Meta.Size != 1 ||
				got.Materials == nil || got.Materials.Meta.Type != "processingplanmaterial" || got.Materials.Meta.Size != 1 ||
				got.Products == nil || got.Products.Meta.Type != "processingplanresult" || got.Products.Meta.Size != 2 ||
				got.ParametricMaterials == nil || got.ParametricMaterials.Meta.Type != "processingplanparametricmaterial" || got.ParametricMaterials.Meta.Size != 0 {
				t.Fatalf("collections = %+v %+v %+v %+v", got.Stages, got.Materials, got.Products, got.ParametricMaterials)
			}
			// Поля таблицы, которых нет в примере ответа, строки наборов при expand и неизвестные поля.
			var sparse processingplan.ProcessingPlan
			if err := codec.Unmarshal([]byte(`{"id":"x","code":"c","costDistributionType":"BY_PRODUCTION","parent":{"meta":{"type":"processingplanfolder"}},"future":1,`+
				`"materials":{"meta":{"href":"h"},"rows":[{"id":"m","quantity":2.5,"materialProcessingPlan":{"meta":{"type":"processingplan"}}}]},`+
				`"parametricMaterials":{"meta":{"href":"p"},"rows":[{"id":"p","paramQuantities":[{"characteristicValue":"L","quantity":4}]}]}}`), &sparse); err != nil {
				t.Fatal(err)
			}
			m, p := sparse.Materials.Rows[0], sparse.ParametricMaterials.Rows[0]
			if sparse.Code != "c" || sparse.CostDistributionType != "BY_PRODUCTION" || sparse.Parent == nil || m.Quantity != 2.5 || m.MaterialProcessingPlan == nil ||
				len(p.ParamQuantities) != 1 || p.ParamQuantities[0].CharacteristicValue != "L" || p.ParamQuantities[0].Quantity != 4 || sparse.Owner != nil || !sparse.Updated.IsZero() {
				t.Fatalf("sparse = %+v", sparse)
			}
		})
	}
}

// cost техкарты и этапа, labourCost и standardHourCost — копейки float64: 0, 0.0, дробь и большое
// значение читаются без потерь. quantity и standardHour — не деньги: те же значения без пересчёта.
func TestNumbersJSON(t *testing.T) {
	values := map[string]float64{"0": 0, "0.0": 0, "12.5": 12.5, "123456789.01": 123456789.01}
	for _, codec := range documenttest.Codecs {
		for raw, want := range values {
			var plan processingplan.ProcessingPlan
			var stage processingplan.Stage
			var material processingplan.Material
			var parametric processingplan.ParametricMaterial
			var product processingplan.Product
			if err := codec.Unmarshal([]byte(`{"cost":`+raw+`}`), &plan); err != nil || plan.Cost != want {
				t.Fatalf("%s: plan cost=%s: %v, %v", codec.Name, raw, plan.Cost, err)
			}
			if err := codec.Unmarshal([]byte(`{"cost":`+raw+`,"labourCost":`+raw+`,"standardHourCost":`+raw+`,"standardHour":`+raw+`}`), &stage); err != nil ||
				stage.Cost != want || stage.LabourCost != want || stage.StandardHourCost != want || stage.StandardHour != want {
				t.Fatalf("%s: stage=%s: %+v, %v", codec.Name, raw, stage, err)
			}
			if err := codec.Unmarshal([]byte(`{"quantity":`+raw+`}`), &material); err != nil || material.Quantity != want {
				t.Fatalf("%s: material quantity=%s: %v, %v", codec.Name, raw, material.Quantity, err)
			}
			if err := codec.Unmarshal([]byte(`{"quantity":`+raw+`,"paramQuantities":[{"quantity":`+raw+`}]}`), &parametric); err != nil || parametric.Quantity != want || parametric.ParamQuantities[0].Quantity != want {
				t.Fatalf("%s: parametric quantity=%s: %+v, %v", codec.Name, raw, parametric, err)
			}
			if err := codec.Unmarshal([]byte(`{"quantity":`+raw+`}`), &product); err != nil || product.Quantity != want {
				t.Fatalf("%s: product quantity=%s: %v, %v", codec.Name, raw, product.Quantity, err)
			}
		}
		// Запрос: копейки и количество уходят как есть, Value(0) — явный 0.
		data, err := codec.Marshal(processingplan.UpdateRequest{Cost: model.Value(0.0), Stages: model.Value([]processingplan.StageRequest{{Cost: model.Value(123456789.01), StandardHour: model.Value(0.25)}})})
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, data, []byte(`{"cost":0,"stages":[{"cost":123456789.01,"standardHour":0.25}]}`))
	}
}

func TestRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", fullCreate, read(t, "custom/processingplan/create-full.golden.json")},
			{"create minimal", minimal, []byte(`{"name":"Example","products":[{"assortment":{"meta":{"href":"` + apiBase + `entity/product/3b89c8e1-6fae-11f1-c5b5-417f00000153","type":"product"}},"quantity":1}]}`)},
			{"update", fullUpdate, read(t, "custom/processingplan/update.golden.json")},
			{"update source example", processingplan.UpdateRequest{Name: model.Value("Техкарта"), Cost: model.Value(100000.0)}, []byte(`{"name":"Техкарта","cost":100000}`)},
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
		reflect.TypeFor[processingplan.CreateRequest]():             {"ID", "Meta", "AccountID", "Updated", "PathName", "CostDistributionType"},
		reflect.TypeFor[processingplan.UpdateRequest]():             {"ID", "Meta", "AccountID", "Updated", "PathName", "CostDistributionType"},
		reflect.TypeFor[processingplan.StageRequest]():              {"Meta", "AccountID", "StandardHourCost"},
		reflect.TypeFor[processingplan.MaterialRequest]():           {"Meta", "AccountID", "MaterialProcessingPlan"},
		reflect.TypeFor[processingplan.ParametricMaterialRequest](): {"Meta", "AccountID"},
		reflect.TypeFor[processingplan.ParamFeatureRequest]():       {"CharacteristicValue"},
		reflect.TypeFor[processingplan.ParamQuantityRequest]():      {"CharacteristicValue"},
		reflect.TypeFor[processingplan.ProductRequest]():            {"Meta", "AccountID"},
	}
	for typ, names := range readOnly {
		for _, name := range names {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
	// name и products обязательны при создании — обычные значения CreateRequest, не Field.
	for name, want := range map[string]reflect.Type{"Name": reflect.TypeFor[string](), "Products": reflect.TypeFor[[]processingplan.ProductRequest]()} {
		if f, _ := reflect.TypeFor[processingplan.CreateRequest]().FieldByName(name); f.Type != want {
			t.Errorf("CreateRequest.%s = %v, want %v", name, f.Type, want)
		}
	}
}

func TestHTTP(t *testing.T) {
	// Ответ сервера, а не эхо запроса: имя и стоимость — только из ответа.
	const response = `{"id":"server-id","name":"server-name","cost":24691,"materials":{"meta":{"href":"https://example.test/materials","size":1}}}`
	check := func(got *processingplan.ProcessingPlan) {
		if got.ID != "server-id" || got.Name != "server-name" || got.Cost != 24691 || got.Materials.Meta.Size != 1 || got.Materials.Rows != nil {
			t.Fatalf("plan = %+v", got)
		}
	}
	t.Run("get", func(t *testing.T) {
		httpCases(t, "GET", planPath, "", response, func(ctx context.Context, c *transport.Client) (*processingplan.ProcessingPlan, error) {
			return service(c).Get(ctx, testID)
		}, check)
	})
	t.Run("create", func(t *testing.T) {
		httpCases(t, "POST", "entity/processingplan", string(read(t, "custom/processingplan/create-full.golden.json")), response, func(ctx context.Context, c *transport.Client) (*processingplan.ProcessingPlan, error) {
			return service(c).Create(ctx, fullCreate)
		}, check)
	})
	t.Run("update", func(t *testing.T) {
		httpCases(t, "PUT", planPath, string(read(t, "custom/processingplan/update.golden.json")), response, func(ctx context.Context, c *transport.Client) (*processingplan.ProcessingPlan, error) {
			return service(c).Update(ctx, testID, fullUpdate)
		}, check)
	})
}

func TestWrites(t *testing.T) {
	const base = "{base}/entity/processingplan"
	planMeta := `{"href":"` + base + `/` + testID + `","metadataHref":"` + base + `/metadata","type":"processingplan","mediaType":"application/json"}`
	runWrites(t, []writeCase{
		{"delete", "DELETE", planPath, "", "", true, func(ctx context.Context, c *transport.Client) error { return service(c).Delete(ctx, testID) }},
		{"create many", "POST", "entity/processingplan", `[{"name":"Example","products":[{"assortment":{"meta":{"href":"` + apiBase + `entity/product/3b89c8e1-6fae-11f1-c5b5-417f00000153","type":"product"}},"quantity":1}]}]`, `[{"id":"x","cost":1000}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateMany(ctx, []processingplan.CreateRequest{minimal})
				if err == nil && (r[0].Value.ID != "x" || r[0].Value.Cost != 1000) {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"update many", "POST", "entity/processingplan", `[{"meta":` + planMeta + `,"name":"Техкарта","cost":100000}]`, `[{"id":"y"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).UpdateMany(ctx, []batch.Update[processingplan.UpdateRequest]{{ID: testID, Request: processingplan.UpdateRequest{Name: model.Value("Техкарта"), Cost: model.Value(100000.0)}}})
				if err == nil && r[0].Value.ID != "y" {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete many", "POST", "entity/processingplan/delete", `[{"meta":` + planMeta + `},{"meta":` + strings.Replace(planMeta, testID, fixtureID, 1) + `}]`, `[{"info":"a"},{"info":"b"}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).DeleteMany(ctx, []string{testID, fixtureID})
				if err == nil && (r[0].Info != "a" || r[1].Info != "b") {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
	})
}

func TestValidation(t *testing.T) {
	// Ровно 1000 строк каждого набора уходят на сервер.
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, `{"id":"server-id"}`) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	{
		s, ctx := service(client), t.Context()
		input := minimal
		input.Products = make([]processingplan.ProductRequest, 1000)
		input.Stages = model.Value(make([]processingplan.StageRequest, 1000))
		input.Materials = model.Value(make([]processingplan.MaterialRequest, 1000))
		input.ParametricMaterials = model.Value(make([]processingplan.ParametricMaterialRequest, 1000))
		if _, err := s.Create(ctx, input); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Update(ctx, testID, processingplan.UpdateRequest{Products: model.Value(input.Products), Materials: input.Materials}); err != nil {
			t.Fatal(err)
		}
	}
	noHTTP(t, func(s *processingplan.Service, ctx context.Context) map[string]func(...transport.RequestOption) error {
		create := func(change func(*processingplan.CreateRequest)) func(...transport.RequestOption) error {
			return func(o ...transport.RequestOption) error {
				input := minimal
				change(&input)
				_, err := s.Create(ctx, input, o...)
				return err
			}
		}
		return map[string]func(...transport.RequestOption) error{
			"get bad id":         func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
			"get path":           func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
			"create no name":     create(func(r *processingplan.CreateRequest) { r.Name = "" }),
			"create no products": create(func(r *processingplan.CreateRequest) { r.Products = nil }),
			"create empty products": create(func(r *processingplan.CreateRequest) {
				r.Products = []processingplan.ProductRequest{}
			}),
			"create 1001 products": create(func(r *processingplan.CreateRequest) { r.Products = make([]processingplan.ProductRequest, 1001) }),
			"create 1001 stages": create(func(r *processingplan.CreateRequest) {
				r.Stages = model.Value(make([]processingplan.StageRequest, 1001))
			}),
			"create 1001 materials": create(func(r *processingplan.CreateRequest) {
				r.Materials = model.Value(make([]processingplan.MaterialRequest, 1001))
			}),
			"create 1001 parametric": create(func(r *processingplan.CreateRequest) {
				r.ParametricMaterials = model.Value(make([]processingplan.ParametricMaterialRequest, 1001))
			}),
			"update empty": func(o ...transport.RequestOption) error {
				_, err := s.Update(ctx, testID, processingplan.UpdateRequest{}, o...)
				return err
			},
			"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", fullUpdate, o...); return err },
			"update 1001 products": func(o ...transport.RequestOption) error {
				_, err := s.Update(ctx, testID, processingplan.UpdateRequest{Products: model.Value(make([]processingplan.ProductRequest, 1001))}, o...)
				return err
			},
			"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
			"create many empty": func(o ...transport.RequestOption) error {
				_, err := s.CreateMany(ctx, nil, o...)
				return err
			},
			"update many bad id": func(o ...transport.RequestOption) error {
				_, err := s.UpdateMany(ctx, []batch.Update[processingplan.UpdateRequest]{{ID: "bad", Request: fullUpdate}}, o...)
				return err
			},
			"update many empty request": func(o ...transport.RequestOption) error {
				_, err := s.UpdateMany(ctx, []batch.Update[processingplan.UpdateRequest]{{ID: testID}}, o...)
				return err
			},
			"delete many bad id": func(o ...transport.RequestOption) error {
				_, err := s.DeleteMany(ctx, []string{"bad"}, o...)
				return err
			},
		}
	})
	if _, err := processingplan.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}

func TestList(t *testing.T) {
	listtest.Run(t, "processingplan", read(t, "custom/processingplan/read.json"), "", "archived", func(c *transport.Client) listtest.ListFunc[processingplan.ProcessingPlan] {
		return service(c).List
	}, func(t *testing.T, got processingplan.ProcessingPlan) {
		if got.ID != fixtureID || got.Cost != 1000 || got.Products == nil || got.Products.Meta.Size != 2 {
			t.Fatalf("row = %+v", got)
		}
	})
}
