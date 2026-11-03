package processingplan_test

import (
	"context"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/processingplan"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const materialID = "120b4591-b0bd-11e6-5bed-427b00000001"

// Пример ответа «Получить материал» источника.
func TestMaterialFixture(t *testing.T) {
	data := read(t, "custom/processingplan/material.json")
	for _, codec := range documenttest.Codecs {
		var got processingplan.Material
		if err := codec.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		if got.ID != materialID || got.AccountID != "d55cbfba-91f1-11e6-5bed-427b00000000" || got.Meta.Type != "processingplanmaterial" || got.Quantity != 1 ||
			got.Assortment == nil || got.Assortment.Meta.Type != "product" || got.Product == nil || got.ProcessingProcessPosition == nil || got.MaterialProcessingPlan != nil {
			t.Fatalf("%s: material = %+v", codec.Name, got)
		}
	}
}

func TestMaterialsHTTP(t *testing.T) {
	const row = `{"meta":{"type":"processingplanmaterial"},"id":"` + materialID + `","quantity":5,"assortment":{"meta":{"type":"variant"}},"materialProcessingPlan":{"meta":{"type":"processingplan"}}}`
	check := func(got *processingplan.Material) {
		if got.ID != materialID || got.Quantity != 5 || got.Assortment == nil || got.Assortment.Meta.Type != "variant" || got.MaterialProcessingPlan == nil {
			t.Fatalf("material = %+v", got)
		}
	}
	t.Run("get materials", func(t *testing.T) {
		pageCases(t, planPath+"/materials", "limit=1000&offset=0", `{"meta":{"href":"h","type":"processingplanmaterial","size":1},"rows":[`+row+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[processingplan.Material], error) {
				return service(c).GetMaterials(ctx, testID, query.PageOptions{})
			}, func(got *model.Collection[processingplan.Material]) { check(&got.Rows[0]) })
	})
	t.Run("get material", func(t *testing.T) {
		httpCases(t, "GET", planPath+"/materials/"+materialID, "", row, func(ctx context.Context, c *transport.Client) (*processingplan.Material, error) {
			return service(c).GetMaterial(ctx, testID, materialID)
		}, check)
	})
	// Пример источника: {"quantity": 5}, ответ — массив из одного объекта.
	t.Run("update material", func(t *testing.T) {
		httpCases(t, "PUT", planPath+"/materials/"+materialID, `{"quantity":5}`, `[`+row+`]`, func(ctx context.Context, c *transport.Client) (*processingplan.Material, error) {
			return service(c).UpdateMaterial(ctx, testID, materialID, processingplan.MaterialRequest{Quantity: model.Value(5.0)})
		}, check)
	})
}

func TestMaterialsWrites(t *testing.T) {
	process := ref("entity/processingprocess/d8da40e9-bbf9-11ed-ac12-0010000000bf/positions/934b5505-ac45-11ed-ac12-000f0000002c", "processingprocessposition")
	runWrites(t, []writeCase{
		// Пример источника: два материала, с привязкой к позиции техпроцесса и без.
		{"create materials", "POST", planPath + "/materials",
			`[{"assortment":` + product + `,"quantity":1,"processingProcessPosition":{"meta":{"href":"` + process.Meta.Href + `","type":"processingprocessposition"}}},{"product":` + product + `,"quantity":2}]`,
			`[{"id":"` + materialID + `","quantity":1},{"id":"` + rowID + `","quantity":2}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateMaterials(ctx, testID, []processingplan.MaterialRequest{
					{Assortment: model.Value(productRef), Quantity: model.Value(1.0), ProcessingProcessPosition: model.Value(process)},
					{Product: model.Value(productRef), Quantity: model.Value(2.0)},
				})
				if err == nil && (r[0].Value.ID != materialID || r[1].Value.Quantity != 2) {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete material", "DELETE", planPath + "/materials/" + materialID, "", "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeleteMaterial(ctx, testID, materialID)
		}},
		{"delete many materials", "POST", planPath + "/materials/delete", `[` + rowMeta("materials", "processingplanmaterial", materialID) + `,` + rowMeta("materials", "processingplanmaterial", rowID) + `]`, "", true,
			func(ctx context.Context, c *transport.Client) error {
				return service(c).DeleteManyMaterials(ctx, testID, []string{materialID, rowID})
			}},
	})
}

func TestMaterialsValidation(t *testing.T) {
	noHTTP(t, func(s *processingplan.Service, ctx context.Context) map[string]func(...transport.RequestOption) error {
		return map[string]func(...transport.RequestOption) error{
			"list bad id": func(o ...transport.RequestOption) error {
				_, err := s.GetMaterials(ctx, "bad", query.PageOptions{}, o...)
				return err
			},
			"list offset": func(o ...transport.RequestOption) error {
				_, err := s.GetMaterials(ctx, testID, query.PageOptions{Offset: -1}, o...)
				return err
			},
			"get bad": func(o ...transport.RequestOption) error {
				_, err := s.GetMaterial(ctx, testID, "bad", o...)
				return err
			},
			"create empty": func(o ...transport.RequestOption) error {
				_, err := s.CreateMaterials(ctx, testID, nil, o...)
				return err
			},
			"create bad plan": func(o ...transport.RequestOption) error {
				_, err := s.CreateMaterials(ctx, "bad", []processingplan.MaterialRequest{{Quantity: model.Value(1.0)}}, o...)
				return err
			},
			"update empty": func(o ...transport.RequestOption) error {
				_, err := s.UpdateMaterial(ctx, testID, materialID, processingplan.MaterialRequest{}, o...)
				return err
			},
			"delete bad": func(o ...transport.RequestOption) error { return s.DeleteMaterial(ctx, "bad", materialID, o...) },
			"delete many empty": func(o ...transport.RequestOption) error {
				return s.DeleteManyMaterials(ctx, testID, nil, o...)
			},
			"delete many 1001": func(o ...transport.RequestOption) error {
				return s.DeleteManyMaterials(ctx, testID, make([]string, 1001), o...)
			},
			"delete many bad row": func(o ...transport.RequestOption) error {
				return s.DeleteManyMaterials(ctx, testID, []string{materialID, "bad"}, o...)
			},
		}
	})
}
