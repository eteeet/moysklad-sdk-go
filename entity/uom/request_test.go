package uom_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/uom"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestCreateRequestJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name    string
				request uom.CreateRequest
				want    []byte
			}{
				{"zero", uom.CreateRequest{}, []byte(`{"name":""}`)},
				{"minimal", uom.CreateRequest{Name: "SDK UOM"}, []byte(`{"name":"SDK UOM"}`)},
				{"full", fullCreateRequest(), uomGolden(t, "create")},
				{"explicit empty values", uom.CreateRequest{
					Name:         "SDK UOM",
					Description:  model.Value(""),
					Code:         model.Value(""),
					ExternalCode: model.Value(""),
					Shared:       model.Value(false),
				}, []byte(`{"name":"SDK UOM","description":"","code":"","externalCode":"","shared":false}`)},
				{"shared true", uom.CreateRequest{
					Name:   "SDK UOM",
					Shared: model.Value(true),
				}, []byte(`{"name":"SDK UOM","shared":true}`)},
				{"explicit nulls", uom.CreateRequest{
					Name:         "SDK UOM",
					Description:  model.Null[string](),
					Code:         model.Null[string](),
					ExternalCode: model.Null[string](),
					Shared:       model.Null[bool](),
					Owner:        model.Null[model.Reference](),
					Group:        model.Null[model.Reference](),
				}, []byte(`{"name":"SDK UOM","description":null,"code":null,"externalCode":null,"shared":null,"owner":null,"group":null}`)},
			} {
				t.Run(tc.name, func(t *testing.T) {
					data, err := codec.marshal(tc.request)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, data, tc.want)
				})
			}
		})
	}
}

func TestCreateRequestReadOnlyFields(t *testing.T) {
	typ := reflect.TypeFor[uom.CreateRequest]()
	for _, name := range []string{"Meta", "ID", "AccountID", "Updated"} {
		if _, exists := typ.FieldByName(name); exists {
			t.Errorf("%s contains read-only %s", typ.Name(), name)
		}
	}
}

func fullCreateRequest() uom.CreateRequest {
	return uom.CreateRequest{
		Name:         "SDK UOM",
		Description:  model.Value("Единица SDK"),
		Code:         model.Value("039"),
		ExternalCode: model.Value("sdk-uom"),
		Shared:       model.Value(false),
		Owner: model.Value(model.Reference{Meta: model.Meta{
			Href: "https://api.moysklad.ru/api/remap/1.2/entity/employee/11111111-1111-1111-1111-111111111111",
		}}),
		Group: model.Value(model.Reference{Meta: model.Meta{
			Href: "https://api.moysklad.ru/api/remap/1.2/entity/group/22222222-2222-2222-2222-222222222222",
		}}),
	}
}

func uomGolden(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/custom/uom/" + name + ".golden.json")
	if err != nil {
		t.Fatal(err)
	}
	return data
}
