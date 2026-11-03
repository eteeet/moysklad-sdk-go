package processingplan_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/processingplan"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const parametricID = "d6f43c04-5504-11f1-3593-85de000001ad"

// Пример ответа «отдельного параметрического материала» источника: материал и норма зависят от
// характеристик.
func TestParametricFixture(t *testing.T) {
	data := read(t, "custom/processingplan/parametric-material.json")
	for _, codec := range documenttest.Codecs {
		var got processingplan.ParametricMaterial
		if err := codec.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		f, q := got.ParamFeatures, got.ParamQuantities
		if got.ID != parametricID || got.Meta.Type != "processingplanparametricmaterial" || got.AccountID != "976b7d32-5504-11f1-a690-207c00000001" ||
			got.Assortment != nil || got.Quantity != 0 || got.ProcessingProcessPosition == nil ||
			got.ParamFeatureCharacteristic == nil || got.ParamFeatureCharacteristic.Meta.Type != "attributemetadata" || got.ParamQuantityCharacteristic == nil ||
			len(f) != 2 || f[0].CharacteristicValue != "blue" || f[0].Assortment == nil || f[0].Assortment.Meta.Type != "variant" || f[1].CharacteristicValue != "red" ||
			len(q) != 2 || q[0].CharacteristicValue != "L" || q[0].Quantity != 1 || q[1].CharacteristicValue != "M" {
			t.Fatalf("%s: parametric = %+v", codec.Name, got)
		}
		// Ответ источника «параметризировано только количество»: конкретный assortment, без paramFeatures.
		var quantityOnly processingplan.ParametricMaterial
		if err := codec.Unmarshal([]byte(`{"id":"q","assortment":{"meta":{"type":"product"}},"paramQuantityCharacteristic":{"meta":{"type":"attributemetadata"}},"paramQuantities":[{"characteristicValue":"L","quantity":2.0}]}`), &quantityOnly); err != nil ||
			quantityOnly.Assortment == nil || quantityOnly.ParamFeatureCharacteristic != nil || quantityOnly.ParamFeatures != nil || quantityOnly.ParamQuantities[0].Quantity != 2 {
			t.Fatalf("%s: quantity only = %+v, err = %v", codec.Name, quantityOnly, err)
		}
	}
}

// Параметрический материал — отдельная форма: свои поля характеристик, нет product и
// materialProcessingPlan обычного материала.
func TestParametricShape(t *testing.T) {
	for _, typ := range []reflect.Type{reflect.TypeFor[processingplan.ParametricMaterial](), reflect.TypeFor[processingplan.ParametricMaterialRequest]()} {
		for _, name := range []string{"ParamFeatureCharacteristic", "ParamQuantityCharacteristic", "ParamFeatures", "ParamQuantities"} {
			if _, ok := typ.FieldByName(name); !ok {
				t.Errorf("%s has no %s", typ.Name(), name)
			}
		}
		for _, name := range []string{"Product", "MaterialProcessingPlan"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s has %s of a plain material", typ.Name(), name)
			}
		}
	}
	if _, ok := reflect.TypeFor[processingplan.MaterialRequest]().FieldByName("ParamFeatures"); ok {
		t.Error("plain material request has paramFeatures")
	}
}

func TestParametricHTTP(t *testing.T) {
	const row = `{"meta":{"type":"processingplanmaterial"},"id":"` + parametricID + `","paramQuantityCharacteristic":{"meta":{"type":"attributemetadata"}},` +
		`"paramQuantities":[{"characteristicValue":"L","quantity":4.0},{"characteristicValue":"M","quantity":3.5}],"assortment":{"meta":{"type":"product"}}}`
	check := func(got *processingplan.ParametricMaterial) {
		if got.ID != parametricID || len(got.ParamQuantities) != 2 || got.ParamQuantities[1].CharacteristicValue != "M" || got.ParamQuantities[1].Quantity != 3.5 || got.Assortment == nil {
			t.Fatalf("parametric = %+v", got)
		}
	}
	t.Run("get parametric materials", func(t *testing.T) {
		pageCases(t, planPath+"/parametricmaterials", "limit=10&offset=0", `{"meta":{"href":"h","type":"processingplanparametricmaterial","size":1},"rows":[`+row+`]}`,
			func(ctx context.Context, c *transport.Client) (*model.Collection[processingplan.ParametricMaterial], error) {
				return service(c).GetParametricMaterials(ctx, testID, query.PageOptions{Limit: 10})
			}, func(got *model.Collection[processingplan.ParametricMaterial]) { check(&got.Rows[0]) })
	})
	t.Run("get parametric material", func(t *testing.T) {
		httpCases(t, "GET", planPath+"/parametricmaterials/"+parametricID, "", row, func(ctx context.Context, c *transport.Client) (*processingplan.ParametricMaterial, error) {
			return service(c).GetParametricMaterial(ctx, testID, parametricID)
		}, check)
	})
	// Пример источника меняет paramQuantities; characteristicValue SDK не пишет (только чтение в таблице).
	t.Run("update parametric material", func(t *testing.T) {
		httpCases(t, "PUT", planPath+"/parametricmaterials/"+parametricID, `{"paramQuantities":[{"quantity":4},{"quantity":3.5}]}`, row,
			func(ctx context.Context, c *transport.Client) (*processingplan.ParametricMaterial, error) {
				return service(c).UpdateParametricMaterial(ctx, testID, parametricID, processingplan.ParametricMaterialRequest{
					ParamQuantities: model.Value([]processingplan.ParamQuantityRequest{{Quantity: 4}, {Quantity: 3.5}}),
				})
			}, check)
	})
}

func TestParametricWrites(t *testing.T) {
	feature := ref("entity/variant/metadata/characteristics/58fe2fba-6fae-11f1-c5b5-417f00000183", "attributemetadata")
	runWrites(t, []writeCase{
		// Пример источника «параметризирован только продукт»: материал по характеристике, единая норма.
		{"create parametric materials", "POST", planPath + "/parametricmaterials",
			`[{"paramFeatureCharacteristic":{"meta":{"href":"` + feature.Meta.Href + `","type":"attributemetadata"}},"paramFeatures":[{"assortment":` + product + `}],"quantity":3.5}]`,
			`[{"id":"` + parametricID + `","quantity":3.5,"paramFeatures":[{"characteristicValue":"blue","assortment":{"meta":{"type":"variant"}}}]}]`, false,
			func(ctx context.Context, c *transport.Client) error {
				r, err := service(c).CreateParametricMaterials(ctx, testID, []processingplan.ParametricMaterialRequest{{
					ParamFeatureCharacteristic: model.Value(feature),
					ParamFeatures:              model.Value([]processingplan.ParamFeatureRequest{{Assortment: productRef}}),
					Quantity:                   model.Value(3.5),
				}})
				if err == nil && (r[0].Value.ID != parametricID || r[0].Value.Quantity != 3.5 || r[0].Value.ParamFeatures[0].CharacteristicValue != "blue") {
					t.Errorf("result = %+v", r)
				}
				return err
			}},
		{"delete parametric material", "DELETE", planPath + "/parametricmaterials/" + parametricID, "", "", true, func(ctx context.Context, c *transport.Client) error {
			return service(c).DeleteParametricMaterial(ctx, testID, parametricID)
		}},
		// meta.type строк — как в примере массового удаления источника.
		{"delete many parametric materials", "POST", planPath + "/parametricmaterials/delete", `[` + rowMeta("parametricmaterials", "processingplanmaterial", parametricID) + `]`, "", true,
			func(ctx context.Context, c *transport.Client) error {
				return service(c).DeleteManyParametricMaterials(ctx, testID, []string{parametricID})
			}},
	})
}

func TestParametricValidation(t *testing.T) {
	noHTTP(t, func(s *processingplan.Service, ctx context.Context) map[string]func(...transport.RequestOption) error {
		return map[string]func(...transport.RequestOption) error{
			"list bad id": func(o ...transport.RequestOption) error {
				_, err := s.GetParametricMaterials(ctx, "bad", query.PageOptions{}, o...)
				return err
			},
			"get bad": func(o ...transport.RequestOption) error {
				_, err := s.GetParametricMaterial(ctx, testID, "bad", o...)
				return err
			},
			"create empty": func(o ...transport.RequestOption) error {
				_, err := s.CreateParametricMaterials(ctx, testID, nil, o...)
				return err
			},
			"update empty": func(o ...transport.RequestOption) error {
				_, err := s.UpdateParametricMaterial(ctx, testID, parametricID, processingplan.ParametricMaterialRequest{}, o...)
				return err
			},
			"update bad plan": func(o ...transport.RequestOption) error {
				_, err := s.UpdateParametricMaterial(ctx, "bad", parametricID, processingplan.ParametricMaterialRequest{Quantity: model.Value(1.0)}, o...)
				return err
			},
			"delete bad": func(o ...transport.RequestOption) error { return s.DeleteParametricMaterial(ctx, testID, "bad", o...) },
			"delete many empty": func(o ...transport.RequestOption) error {
				return s.DeleteManyParametricMaterials(ctx, testID, nil, o...)
			},
			"delete many bad plan": func(o ...transport.RequestOption) error {
				return s.DeleteManyParametricMaterials(ctx, "bad", []string{parametricID}, o...)
			},
		}
	})
}
