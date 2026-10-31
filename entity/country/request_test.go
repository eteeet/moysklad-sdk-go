package country_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/country"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestCreateRequestMinimalJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			data, err := codec.marshal(country.CreateRequest{Name: "SDK Country"})
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, data, []byte(`{"name":"SDK Country"}`))
		})
	}
}

func TestCreateRequestJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				input country.CreateRequest
				want  []byte
			}{
				{"zero", country.CreateRequest{}, []byte(`{"name":""}`)},
				{"full", fullCreateRequest(), countryGolden(t, "create")},
				{"empty values", country.CreateRequest{
					Name: "SDK Country", Description: model.Value(""), Code: model.Value(""), ExternalCode: model.Value(""), Shared: model.Value(false),
				}, []byte(`{"name":"SDK Country","description":"","code":"","externalCode":"","shared":false}`)},
				{"shared true", country.CreateRequest{Name: "SDK Country", Shared: model.Value(true)}, []byte(`{"name":"SDK Country","shared":true}`)},
				{"nulls", country.CreateRequest{
					Name: "SDK Country", Description: model.Null[string](), Code: model.Null[string](), ExternalCode: model.Null[string](),
					Shared: model.Null[bool](), Owner: model.Null[model.Reference](), Group: model.Null[model.Reference](),
				}, []byte(`{"name":"SDK Country","description":null,"code":null,"externalCode":null,"shared":null,"owner":null,"group":null}`)},
			} {
				t.Run(tc.name, func(t *testing.T) {
					data, err := codec.marshal(tc.input)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, data, tc.want)
				})
			}
		})
	}
}

func fullCreateRequest() country.CreateRequest {
	return country.CreateRequest{
		Name: "SDK Country", Description: model.Value("Страна SDK"), Code: model.Value("0007"), ExternalCode: model.Value("sdk-country"), Shared: model.Value(false),
		Owner: model.Value(model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/employee/11111111-1111-1111-1111-111111111111"}}),
		Group: model.Value(model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/group/22222222-2222-2222-2222-222222222222"}}),
	}
}

func countryGolden(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/custom/country/" + name + ".golden.json")
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestUpdateRequestJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				input country.UpdateRequest
				want  []byte
			}{
				{"unset", country.UpdateRequest{}, []byte(`{}`)},
				{"partial", country.UpdateRequest{Name: model.Value("renamed")}, []byte(`{"name":"renamed"}`)},
				{"full", fullUpdateRequest(), countryGolden(t, "update")},
				{"false", country.UpdateRequest{Shared: model.Value(false)}, []byte(`{"shared":false}`)},
				{"true", country.UpdateRequest{Shared: model.Value(true)}, []byte(`{"shared":true}`)},
				{"empty strings", country.UpdateRequest{
					Name: model.Value(""), Description: model.Value(""), Code: model.Value(""), ExternalCode: model.Value(""),
				}, []byte(`{"name":"","description":"","code":"","externalCode":""}`)},
				{"nulls", country.UpdateRequest{
					Name: model.Null[string](), Description: model.Null[string](), Code: model.Null[string](), ExternalCode: model.Null[string](),
					Shared: model.Null[bool](), Owner: model.Null[model.Reference](), Group: model.Null[model.Reference](),
				}, []byte(`{"name":null,"description":null,"code":null,"externalCode":null,"shared":null,"owner":null,"group":null}`)},
			} {
				t.Run(tc.name, func(t *testing.T) {
					data, err := codec.marshal(tc.input)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, data, tc.want)
				})
			}
		})
	}
}

func TestRequestsExcludeReadOnlyFields(t *testing.T) {
	for _, typ := range []reflect.Type{reflect.TypeFor[country.CreateRequest](), reflect.TypeFor[country.UpdateRequest]()} {
		for _, name := range []string{"Meta", "ID", "AccountID", "Updated"} {
			if _, exists := typ.FieldByName(name); exists {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func fullUpdateRequest() country.UpdateRequest {
	return country.UpdateRequest{
		Name: model.Value("SDK Country updated"), Description: model.Value(""), Code: model.Value("0008"), ExternalCode: model.Value("sdk-country-updated"), Shared: model.Value(false),
		Owner: model.Value(model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/employee/11111111-1111-1111-1111-111111111111"}}),
		Group: model.Value(model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/group/22222222-2222-2222-2222-222222222222"}}),
	}
}
