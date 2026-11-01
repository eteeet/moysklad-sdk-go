package product_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/product"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestCreateRequestJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				request product.CreateRequest
				want    []byte
			}{
				{product.CreateRequest{}, []byte(`{"name":""}`)},
				{product.CreateRequest{Name: "SDK product"}, []byte(`{"name":"SDK product"}`)},
				{product.CreateRequest{Name: "p", SyncID: model.Value("11111111-1111-1111-1111-111111111111")}, []byte(`{"name":"p","syncId":"11111111-1111-1111-1111-111111111111"}`)},
				{product.CreateRequest{Name: "p", SyncID: model.Value("")}, []byte(`{"name":"p","syncId":""}`)},
				{product.CreateRequest{Name: "p", SyncID: model.Null[string]()}, []byte(`{"name":"p","syncId":null}`)},
				{fullCreateRequest(t), productGolden(t, "create")},
			} {
				data, err := codec.marshal(tc.request)
				if err != nil {
					t.Fatal(err)
				}
				assertJSONEqual(t, data, tc.want)
			}
		})
	}
}

func TestRequestReadOnlyFields(t *testing.T) {
	// SyncID заполняется только при создании.
	if _, exists := reflect.TypeFor[product.UpdateRequest]().FieldByName("SyncID"); exists {
		t.Error("UpdateRequest contains create-only SyncID")
	}
	for _, typ := range []reflect.Type{reflect.TypeFor[product.CreateRequest](), reflect.TypeFor[product.UpdateRequest]()} {
		for _, name := range []string{"ID", "AccountID", "Updated", "PathName", "EffectiveVAT", "EffectiveVATEnabled", "VariantsCount", "Meta"} {
			if _, exists := typ.FieldByName(name); exists {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func TestUpdateRequestJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				request product.UpdateRequest
				want    []byte
			}{
				{product.UpdateRequest{}, []byte(`{}`)},
				{product.UpdateRequest{Name: model.Value("renamed")}, []byte(`{"name":"renamed"}`)},
				{product.UpdateRequest{SalePrices: model.Value([]model.SalePrice{})}, []byte(`{"salePrices":[]}`)},
				{product.UpdateRequest{Attributes: model.Value([]model.Attribute{})}, []byte(`{"attributes":[]}`)},
				{fullUpdateRequest(), productGolden(t, "update")},
			} {
				data, err := codec.marshal(tc.request)
				if err != nil {
					t.Fatal(err)
				}
				assertJSONEqual(t, data, tc.want)
			}
		})
	}
}

func fullUpdateRequest() product.UpdateRequest {
	return product.UpdateRequest{
		Description: model.Value(""), Archived: model.Value(false), Weight: model.Value(0.0), ProductFolder: model.Null[model.Reference](),
		SalePrices: model.Value([]model.SalePrice{{Value: 0, PriceType: model.PriceType{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/context/companysettings/pricetype/33333333-3333-3333-3333-333333333333"}}}}),
		Attributes: model.Value([]model.Attribute{
			model.AttrNull(model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/product/metadata/attributes/11111111-1111-1111-1111-111111111111"}),
			model.AttrFileNull(model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/product/metadata/attributes/22222222-2222-2222-2222-222222222222"}),
		}),
	}
}

func fullCreateRequest(t *testing.T) product.CreateRequest {
	t.Helper()
	attribute, err := model.AttrValue(model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/product/metadata/attributes/11111111-1111-1111-1111-111111111111"}, int64(9007199254740993))
	if err != nil {
		t.Fatal(err)
	}
	ref := func(kind string) model.Field[model.Reference] {
		return model.Value(model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/" + kind + "/11111111-1111-1111-1111-111111111111"}})
	}
	return product.CreateRequest{
		Name: "SDK product", Code: model.Value("sku-1"), ExternalCode: model.Value("sdk-example-test"), Description: model.Value("description"), Article: model.Value("A-1"),
		Archived: model.Value(false), Shared: model.Value(false), Weight: model.Value(0.0), Volume: model.Value(0.0),
		VAT: model.Value(0), VATEnabled: model.Value(false), UseParentVAT: model.Value(false), DiscountProhibited: model.Value(false),
		Owner: ref("employee"), Group: ref("group"), ProductFolder: ref("productfolder"), UOM: ref("uom"), Supplier: ref("counterparty"), Country: ref("country"),
		BuyPrice: model.Value(model.Price{}), MinPrice: model.Value(model.Price{Value: 500}),
		SalePrices: model.Value([]model.SalePrice{{Value: 100000, PriceType: model.PriceType{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/context/companysettings/pricetype/33333333-3333-3333-3333-333333333333"}}}}),
		Attributes: model.Value([]model.Attribute{attribute}),
	}
}

func productGolden(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/custom/product/" + name + ".golden.json")
	if err != nil {
		t.Fatal(err)
	}
	return data
}
