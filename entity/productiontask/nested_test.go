package productiontask_test

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/productiontask"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Строки наборов из примеров ответа источника, обоими кодеками.
func TestNestedFixtures(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var row productiontask.ProductionRow
			if err := codec.Unmarshal(read(t, "custom/productiontask/productionrow.json"), &row); err != nil {
				t.Fatal(err)
			}
			if row.ID != "85b4d986-99d6-11ee-ac12-000f0000011c" || row.Meta.Type != "productionrow" || row.Name != "example-name" || row.ExternalCode != "3oKBsRwZigCESI13Mncqq1" ||
				row.ProductionVolume != 2 || row.ProcessingPlan == nil || row.ProcessingPlan.Meta.Type != "processingplan" || row.ProductionVariant == nil ||
				row.ProductionVariant.Meta.Type != "variant" || !time.Time(row.Updated).Equal(time.Time(documenttest.DateTime("2023-12-13T19:41:53.990+03:00"))) {
				t.Fatalf("production row = %+v", row)
			}
			var product productiontask.Product
			if err := codec.Unmarshal(read(t, "custom/productiontask/product.json"), &product); err != nil {
				t.Fatal(err)
			}
			if product.ID != productID || product.Meta.Type != "productiontaskresult" || product.PlanQuantity != 2 || product.Assortment == nil ||
				product.Assortment.Meta.Type != "product" || product.ProductionRow == nil || product.ProductionRow.Meta.Type != "productionrow" {
				t.Fatalf("product = %+v", product)
			}
			var stage productiontask.Stage
			if err := codec.Unmarshal(read(t, "custom/productiontask/stage.json"), &stage); err != nil {
				t.Fatal(err)
			}
			if stage.ID != stageID || stage.Meta.Type != "productionstage" || stage.Stage == nil || stage.Stage.Meta.Type != "processingstage" || stage.ProductionRow == nil ||
				stage.Materials == nil || stage.Materials.Meta.Type != "productiontaskmaterial" || stage.Materials.Meta.Size != 2 || stage.Files == nil || stage.Files.Meta.Size != 2 ||
				stage.OrderingPosition != 0 || stage.TotalQuantity != 1 || stage.AvailableQuantity != 1 || stage.CompletedQuantity != 0 || stage.ProcessingUnitCost != 2 ||
				stage.EnableHourAccounting || !strings.HasPrefix(stage.Instruction, "1.Откройте") || stage.MaterialStore != nil || !stage.PlannedEndDate.IsZero() {
				t.Fatalf("stage = %+v", stage)
			}
			var sparse productiontask.Stage
			if err := codec.Unmarshal([]byte(`{"id":"s","enableHourAccounting":true,"orderingPosition":3,"blockedQuantity":1.5,"skippedQuantity":0.5,"labourUnitCost":30.5,`+
				`"standardHourCost":100,"standardHourUnit":43.5,"plannedEndDate":"2023-12-20 18:00:00.000","materialStore":{"meta":{"type":"store"}},"status":"Unknown"}`), &sparse); err != nil {
				t.Fatal(err)
			}
			if !sparse.EnableHourAccounting || sparse.OrderingPosition != 3 || sparse.BlockedQuantity != 1.5 || sparse.SkippedQuantity != 0.5 || sparse.LabourUnitCost != 30.5 ||
				sparse.StandardHourCost != 100 || sparse.StandardHourUnit != 43.5 || sparse.PlannedEndDate.IsZero() || sparse.MaterialStore == nil {
				t.Fatalf("sparse stage = %+v", sparse)
			}
			var material productiontask.Material
			if err := codec.Unmarshal(read(t, "custom/productiontask/material.json"), &material); err != nil {
				t.Fatal(err)
			}
			if material.ID != materialID || material.Meta.Type != "productiontaskmaterial" || material.PlanQuantity != 3 || material.Assortment == nil ||
				material.AccountID != "d55cbfba-91f1-11e6-5bed-427b00000000" {
				t.Fatalf("material = %+v", material)
			}
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

func TestNestedHTTP(t *testing.T) {
	const (
		row      = `{"meta":{"type":"productionrow"},"id":"` + rowID + `","productionVolume":5,"processingPlan":{"meta":{"type":"processingplan"}}}`
		product  = `{"meta":{"type":"productiontaskresult"},"id":"` + productID + `","planQuantity":22,"assortment":{"meta":{"type":"product"}}}`
		stage    = `{"meta":{"type":"productionstage"},"id":"` + stageID + `","processingUnitCost":70,"labourUnitCost":30.5,"enableHourAccounting":true}`
		material = `{"meta":{"type":"productiontaskmaterial"},"id":"` + materialID + `","planQuantity":3,"assortment":{"meta":{"type":"product"}}}`
	)
	assortment := ref("product/57c9a841-99d8-11ee-0a83-0a2e00000773", "product")
	assortmentJSON := `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/57c9a841-99d8-11ee-0a83-0a2e00000773","type":"product"}}`
	checkRow := func(got *productiontask.ProductionRow) {
		if got.ID != rowID || got.ProductionVolume != 5 || got.ProcessingPlan == nil {
			t.Fatalf("row = %+v", got)
		}
	}
	checkProduct := func(got *productiontask.Product) {
		if got.ID != productID || got.PlanQuantity != 22 || got.Assortment == nil {
			t.Fatalf("product = %+v", got)
		}
	}
	checkStage := func(got *productiontask.Stage) {
		if got.ID != stageID || got.ProcessingUnitCost != 70 || got.LabourUnitCost != 30.5 || !got.EnableHourAccounting {
			t.Fatalf("stage = %+v", got)
		}
	}
	checkMaterial := func(got *productiontask.Material) {
		if got.ID != materialID || got.PlanQuantity != 3 || got.Assortment == nil {
			t.Fatalf("material = %+v", got)
		}
	}
	t.Run("get production rows", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/productionrows?limit=2&offset=7", "", `{"meta":{"href":"h"},"rows":[`+row+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[productiontask.ProductionRow], error) {
				return service(c).GetProductionRows(ctx, testID, query.PageOptions{Limit: 2, Offset: 7})
			}, func(got *model.Collection[productiontask.ProductionRow]) { checkRow(&got.Rows[0]) })
	})
	t.Run("get production row", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/productionrows/"+rowID, "", row, func(ctx context.Context, c *transport.Client) (*productiontask.ProductionRow, error) {
			return service(c).GetProductionRow(ctx, testID, rowID)
		}, checkRow)
	})
	t.Run("update production row", func(t *testing.T) {
		httpCases(t, "PUT", docPath+"/productionrows/"+rowID, `{"productionVolume":5}`, row, func(ctx context.Context, c *transport.Client) (*productiontask.ProductionRow, error) {
			return service(c).UpdateProductionRow(ctx, testID, rowID, productiontask.ProductionRowUpdate{ProductionVolume: model.Value(5.0)})
		}, checkRow)
	})
	t.Run("get products", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/products?limit=1000&offset=0", "", `{"meta":{"href":"h"},"rows":[`+product+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[productiontask.Product], error) {
				return service(c).GetProducts(ctx, testID, query.PageOptions{})
			}, func(got *model.Collection[productiontask.Product]) { checkProduct(&got.Rows[0]) })
	})
	t.Run("get product", func(t *testing.T) {
		httpCases(t, "GET", docPath+"/products/"+productID, "", product, func(ctx context.Context, c *transport.Client) (*productiontask.Product, error) {
			return service(c).GetProduct(ctx, testID, productID)
		}, checkProduct)
	})
	t.Run("update product", func(t *testing.T) {
		httpCases(t, "PUT", docPath+"/products/"+productID, `{"assortment":`+assortmentJSON+`,"planQuantity":22}`, product,
			func(ctx context.Context, c *transport.Client) (*productiontask.Product, error) {
				return service(c).UpdateProduct(ctx, testID, productID, productiontask.ProductRequest{Assortment: model.Value(assortment), PlanQuantity: model.Value(22.0)})
			}, checkProduct)
	})
	t.Run("update stage", func(t *testing.T) {
		httpCases(t, "PUT", stagePath, `{"enableHourAccounting":false,"labourUnitCost":30.5}`, stage, func(ctx context.Context, c *transport.Client) (*productiontask.Stage, error) {
			return service(c).UpdateStage(ctx, stageID, productiontask.StageUpdateRequest{EnableHourAccounting: model.Value(false), LabourUnitCost: model.Value(30.5)})
		}, checkStage)
	})
	t.Run("get materials", func(t *testing.T) {
		httpCases(t, "GET", stagePath+"/materials?limit=5&offset=0", "", `{"meta":{"href":"h"},"rows":[`+material+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[productiontask.Material], error) {
				return service(c).GetMaterials(ctx, stageID, query.PageOptions{Limit: 5})
			}, func(got *model.Collection[productiontask.Material]) { checkMaterial(&got.Rows[0]) })
	})
	// Пример ответа источника на добавление — массив из одного объекта.
	for name, response := range map[string]string{"object": material, "array of one": `[` + material + `]`} {
		t.Run("create material/"+name, func(t *testing.T) {
			httpCases(t, "POST", stagePath+"/materials", `{"assortment":`+assortmentJSON+`,"planQuantity":2}`, response,
				func(ctx context.Context, c *transport.Client) (*productiontask.Material, error) {
					return service(c).CreateMaterial(ctx, stageID, productiontask.MaterialRequest{Assortment: model.Value(assortment), PlanQuantity: model.Value(2.0)})
				}, checkMaterial)
		})
	}
	t.Run("update material", func(t *testing.T) {
		httpCases(t, "PUT", stagePath+"/materials/"+materialID, `{"planQuantity":3}`, material, func(ctx context.Context, c *transport.Client) (*productiontask.Material, error) {
			return service(c).UpdateMaterial(ctx, stageID, materialID, productiontask.MaterialRequest{PlanQuantity: model.Value(3.0)})
		}, checkMaterial)
	})
}

// Этапы — GET entity/productionstage с фильтром по href задания от BaseURL клиента; адрес сервера
// заранее неизвестен, поэтому ответы documenttest.ResponseCases отдаёт свой сервер.
func TestGetStages(t *testing.T) {
	stage := `{"meta":{"type":"productionstage"},"id":"` + stageID + `","standardHourCost":0,"completedQuantity":0.5}`
	for _, tc := range documenttest.ResponseCases(`{"meta":{"href":"h","type":"productionstage"},"rows":[` + stage + `]}`) {
		t.Run(tc.Name, func(t *testing.T) {
			var calls atomic.Int32
			var base string
			server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				want := "filter=" + strings.NewReplacer(":", "%3A", "/", "%2F", "=", "%3D").Replace("productionTask="+base+"/entity/productiontask/"+testID) + "&limit=1000&offset=0"
				if r.Method != http.MethodGet || r.URL.Path != "/api/remap/1.2/entity/productionstage" || r.URL.RawQuery != want {
					t.Errorf("request = %s %s, want query %s", r.Method, r.URL, want)
				}
				if r.Header.Get("Authorization") != "Bearer test-token" || r.Header.Get("Accept-Encoding") != "gzip" || r.Header.Get("Content-Type") != "" {
					t.Errorf("headers = %v", r.Header)
				}
				var out io.Writer = w
				if tc.Gzip {
					w.Header().Set("Content-Encoding", "gzip")
				}
				w.WriteHeader(tc.Status)
				if tc.Gzip {
					zw := gzip.NewWriter(w)
					defer zw.Close()
					out = zw
				}
				io.WriteString(out, tc.Body)
			}))
			defer server.Close()
			base = server.URL + "/api/remap/1.2"
			auth, err := transport.BearerToken("test-token")
			if err != nil {
				t.Fatal(err)
			}
			client, err := transport.NewClient(transport.Config{BaseURL: base, Auth: auth})
			if err != nil {
				t.Fatal(err)
			}
			got, err := service(client).GetStages(t.Context(), testID, query.PageOptions{})
			if calls.Load() != 1 {
				t.Fatalf("calls = %d", calls.Load())
			}
			if tc.Success {
				if err != nil || len(got.Rows) != 1 || got.Rows[0].ID != stageID || got.Rows[0].CompletedQuantity != 0.5 || got.Meta.Type != "productionstage" {
					t.Fatalf("stages = %+v, err = %v", got, err)
				}
				return
			}
			var apiErr *sdkerrors.APIError
			if got != nil || err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "get productionstage") || errors.As(err, &apiErr) != tc.API {
				t.Fatalf("stages = %+v, err = %v", got, err)
			}
		})
	}
	documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*model.Collection[productiontask.Stage], error) {
		return service(c).GetStages(ctx, testID, query.PageOptions{})
	})
}

// Ошибки входа — до HTTP, ResponseInfo сброшен; ровно 200 позиций и 1000 продуктов уходят.
func TestValidation(t *testing.T) {
	// Корректная новая строка: ошибки в случаях ниже вызывает только изменённая часть.
	newRow := productiontask.ProductRequest{Assortment: model.Value(ref("product/57c9a841-99d8-11ee-0a83-0a2e00000773", "product")), PlanQuantity: model.Value(1.0)}
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
	rows := func(n int) model.Field[[]productiontask.ProductionRowRequest] {
		return model.Value(make([]productiontask.ProductionRowRequest, n))
	}
	products := func(n int) model.Field[[]productiontask.ProductRequest] {
		return model.Value(make([]productiontask.ProductRequest, n))
	}
	limit := minimal
	limit.ProductionRows, limit.Products = rows(200), products(1000)
	if _, err := s.Create(ctx, limit); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(ctx, testID, productiontask.UpdateRequest{ProductionRows: rows(200), Products: products(1000)}); err != nil {
		t.Fatal(err)
	}
	create := func(change func(*productiontask.CreateRequest)) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error {
			input := minimal
			change(&input)
			_, err := s.Create(ctx, input, o...)
			return err
		}
	}
	page := func(get func(context.Context, string, query.PageOptions, ...transport.RequestOption) error, id string, p query.PageOptions) func(...transport.RequestOption) error {
		return func(o ...transport.RequestOption) error { return get(ctx, id, p, o...) }
	}
	getRows := func(ctx context.Context, id string, p query.PageOptions, o ...transport.RequestOption) error {
		_, err := s.GetProductionRows(ctx, id, p, o...)
		return err
	}
	getProducts := func(ctx context.Context, id string, p query.PageOptions, o ...transport.RequestOption) error {
		_, err := s.GetProducts(ctx, id, p, o...)
		return err
	}
	getStages := func(ctx context.Context, id string, p query.PageOptions, o ...transport.RequestOption) error {
		_, err := s.GetStages(ctx, id, p, o...)
		return err
	}
	getMaterials := func(ctx context.Context, id string, p query.PageOptions, o ...transport.RequestOption) error {
		_, err := s.GetMaterials(ctx, id, p, o...)
		return err
	}
	calls.Store(0)
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id":           func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":             func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no org":        create(func(r *productiontask.CreateRequest) { r.Organization = model.Reference{} }),
		"create no materials":  create(func(r *productiontask.CreateRequest) { r.MaterialsStore = model.Reference{} }),
		"create no products":   create(func(r *productiontask.CreateRequest) { r.ProductsStore = model.Reference{} }),
		"create 201 rows":      create(func(r *productiontask.CreateRequest) { r.ProductionRows = rows(201) }),
		"create 1001 products": create(func(r *productiontask.CreateRequest) { r.Products = products(1001) }),
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, productiontask.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", fullUpdate, o...); return err },
		"update 201 rows": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, productiontask.UpdateRequest{ProductionRows: rows(201)}, o...)
			return err
		},
		"update 1001 products": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, productiontask.UpdateRequest{Products: products(1001)}, o...)
			return err
		},
		"delete bad id":     func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
		"trash bad id":      func(o ...transport.RequestOption) error { return s.Trash(ctx, "bad", o...) },
		"create many empty": func(o ...transport.RequestOption) error { _, err := s.CreateMany(ctx, nil, o...); return err },
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[productiontask.UpdateRequest]{{ID: "bad", Request: fullUpdate}}, o...)
			return err
		},
		"update many empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMany(ctx, []batch.Update[productiontask.UpdateRequest]{{ID: testID}}, o...)
			return err
		},
		"delete many bad id": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, []string{"bad"}, o...)
			return err
		},
		"rows bad id": page(getRows, "bad", query.PageOptions{}),
		"rows fields": page(getRows, testID, query.PageOptions{Fields: "stock"}),
		"rows limit":  page(getRows, testID, query.PageOptions{Limit: 1001}),
		"get row bad": func(o ...transport.RequestOption) error {
			_, err := s.GetProductionRow(ctx, testID, "bad", o...)
			return err
		},
		"update row bad doc": func(o ...transport.RequestOption) error {
			_, err := s.UpdateProductionRow(ctx, "bad", rowID, productiontask.ProductionRowUpdate{ProductionVolume: model.Value(1.0)}, o...)
			return err
		},
		"update row empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateProductionRow(ctx, testID, rowID, productiontask.ProductionRowUpdate{}, o...)
			return err
		},
		"delete row bad":  func(o ...transport.RequestOption) error { return s.DeleteProductionRow(ctx, testID, "bad", o...) },
		"products bad id": page(getProducts, "bad", query.PageOptions{}),
		"products fields": page(getProducts, testID, query.PageOptions{Fields: "stock"}),
		"products offset": page(getProducts, testID, query.PageOptions{Offset: -1}),
		"get product bad": func(o ...transport.RequestOption) error {
			_, err := s.GetProduct(ctx, "bad", productID, o...)
			return err
		},
		"create products empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateProducts(ctx, testID, nil, o...)
			return err
		},
		"create products bad id": func(o ...transport.RequestOption) error {
			_, err := s.CreateProducts(ctx, "bad", []productiontask.ProductRequest{newRow}, o...)
			return err
		},
		"create products no assortment": func(o ...transport.RequestOption) error {
			_, err := s.CreateProducts(ctx, testID, []productiontask.ProductRequest{newRow, {PlanQuantity: model.Value(1.0)}}, o...)
			return err
		},
		"create products no quantity": func(o ...transport.RequestOption) error {
			_, err := s.CreateProducts(ctx, testID, []productiontask.ProductRequest{{Assortment: newRow.Assortment}}, o...)
			return err
		},
		"update product empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateProduct(ctx, testID, productID, productiontask.ProductRequest{}, o...)
			return err
		},
		"delete product bad":         func(o ...transport.RequestOption) error { return s.DeleteProduct(ctx, testID, "bad", o...) },
		"delete many products empty": func(o ...transport.RequestOption) error { return s.DeleteManyProducts(ctx, testID, nil, o...) },
		"delete many products 1001": func(o ...transport.RequestOption) error {
			return s.DeleteManyProducts(ctx, testID, make([]string, 1001), o...)
		},
		"delete many products bad": func(o ...transport.RequestOption) error {
			return s.DeleteManyProducts(ctx, testID, []string{productID, "bad"}, o...)
		},
		"delete many products doc": func(o ...transport.RequestOption) error {
			return s.DeleteManyProducts(ctx, "bad", []string{productID}, o...)
		},
		"stages bad id": page(getStages, "bad", query.PageOptions{}),
		"stages fields": page(getStages, testID, query.PageOptions{Fields: "stock"}),
		"stages limit":  page(getStages, testID, query.PageOptions{Limit: -1}),
		"update stage bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateStage(ctx, "bad", productiontask.StageUpdateRequest{LabourUnitCost: model.Value(1.0)}, o...)
			return err
		},
		"update stage empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateStage(ctx, stageID, productiontask.StageUpdateRequest{}, o...)
			return err
		},
		"materials bad id": page(getMaterials, "bad", query.PageOptions{}),
		"materials fields": page(getMaterials, stageID, query.PageOptions{Fields: "stock"}),
		"create material bad id": func(o ...transport.RequestOption) error {
			_, err := s.CreateMaterial(ctx, "bad", productiontask.MaterialRequest(newRow), o...)
			return err
		},
		"create material no assortment": func(o ...transport.RequestOption) error {
			_, err := s.CreateMaterial(ctx, stageID, productiontask.MaterialRequest{PlanQuantity: model.Value(1.0)}, o...)
			return err
		},
		"create material null quantity": func(o ...transport.RequestOption) error {
			_, err := s.CreateMaterial(ctx, stageID, productiontask.MaterialRequest{Assortment: newRow.Assortment, PlanQuantity: model.Null[float64]()}, o...)
			return err
		},
		"update material bad": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMaterial(ctx, stageID, "bad", productiontask.MaterialRequest{PlanQuantity: model.Value(1.0)}, o...)
			return err
		},
		"update material empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateMaterial(ctx, stageID, materialID, productiontask.MaterialRequest{}, o...)
			return err
		},
		"delete material bad": func(o ...transport.RequestOption) error { return s.DeleteMaterial(ctx, "bad", materialID, o...) },
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "production") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := productiontask.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}
