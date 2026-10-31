package expenseitem_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/expenseitem"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestCreateRequestJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name    string
				request expenseitem.CreateRequest
				want    []byte
			}{
				{"zero", expenseitem.CreateRequest{}, []byte(`{"name":""}`)},
				{"minimal", expenseitem.CreateRequest{Name: "SDK ExpenseItem"}, []byte(`{"name":"SDK ExpenseItem"}`)},
				{"full", fullCreateRequest(), expenseitemGolden(t, "create")},
				{"explicit empty values", expenseitem.CreateRequest{
					Name:              "SDK ExpenseItem",
					Description:       model.Value(""),
					Code:              model.Value(""),
					ExternalCode:      model.Value(""),
					Shared:            model.Value(false),
					Archived:          model.Value(false),
					OperatingExpenses: model.Value(false),
				}, []byte(`{"name":"SDK ExpenseItem","description":"","code":"","externalCode":"","shared":false,"archived":false,"operatingExpenses":false}`)},
				{"flags true", expenseitem.CreateRequest{
					Name:              "SDK ExpenseItem",
					Shared:            model.Value(true),
					Archived:          model.Value(true),
					OperatingExpenses: model.Value(true),
				}, []byte(`{"name":"SDK ExpenseItem","shared":true,"archived":true,"operatingExpenses":true}`)},
				{"explicit nulls", expenseitem.CreateRequest{
					Name:              "SDK ExpenseItem",
					Description:       model.Null[string](),
					Code:              model.Null[string](),
					ExternalCode:      model.Null[string](),
					Shared:            model.Null[bool](),
					Archived:          model.Null[bool](),
					OperatingExpenses: model.Null[bool](),
					Owner:             model.Null[model.Reference](),
					Group:             model.Null[model.Reference](),
				}, []byte(`{"name":"SDK ExpenseItem","description":null,"code":null,"externalCode":null,"shared":null,"archived":null,"operatingExpenses":null,"owner":null,"group":null}`)},
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

func fullCreateRequest() expenseitem.CreateRequest {
	return expenseitem.CreateRequest{
		Name:              "SDK ExpenseItem",
		Description:       model.Value("Статья SDK"),
		Code:              model.Value("039"),
		ExternalCode:      model.Value("sdk-expenseitem"),
		Shared:            model.Value(false),
		Archived:          model.Value(true),
		OperatingExpenses: model.Value(false),
		Owner: model.Value(model.Reference{Meta: model.Meta{
			Href: "https://api.moysklad.ru/api/remap/1.2/entity/employee/11111111-1111-1111-1111-111111111111",
		}}),
		Group: model.Value(model.Reference{Meta: model.Meta{
			Href: "https://api.moysklad.ru/api/remap/1.2/entity/group/22222222-2222-2222-2222-222222222222",
		}}),
	}
}

func expenseitemGolden(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/custom/expenseitem/" + name + ".golden.json")
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestRequestReadOnlyFields(t *testing.T) {
	for _, typ := range []reflect.Type{reflect.TypeFor[expenseitem.CreateRequest](), reflect.TypeFor[expenseitem.UpdateRequest]()} {
		for _, name := range []string{"Meta", "ID", "AccountID", "Updated"} {
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
				name    string
				request expenseitem.UpdateRequest
				want    []byte
			}{
				{"zero", expenseitem.UpdateRequest{}, []byte(`{}`)},
				{"one field", expenseitem.UpdateRequest{Name: model.Value("renamed")}, []byte(`{"name":"renamed"}`)},
				{"full", fullUpdateRequest(), expenseitemGolden(t, "update")},
				{"shared false", expenseitem.UpdateRequest{Shared: model.Value(false)}, []byte(`{"shared":false}`)},
				{"shared true", expenseitem.UpdateRequest{Shared: model.Value(true)}, []byte(`{"shared":true}`)},
				{"flags false", expenseitem.UpdateRequest{Archived: model.Value(false), OperatingExpenses: model.Value(false)}, []byte(`{"archived":false,"operatingExpenses":false}`)},
				{"flags true", expenseitem.UpdateRequest{Archived: model.Value(true), OperatingExpenses: model.Value(true)}, []byte(`{"archived":true,"operatingExpenses":true}`)},
				{"explicit empty strings", expenseitem.UpdateRequest{
					Name:         model.Value(""),
					Description:  model.Value(""),
					Code:         model.Value(""),
					ExternalCode: model.Value(""),
				}, []byte(`{"name":"","description":"","code":"","externalCode":""}`)},
				{"explicit nulls", expenseitem.UpdateRequest{
					Name:              model.Null[string](),
					Description:       model.Null[string](),
					Code:              model.Null[string](),
					ExternalCode:      model.Null[string](),
					Shared:            model.Null[bool](),
					Archived:          model.Null[bool](),
					OperatingExpenses: model.Null[bool](),
					Owner:             model.Null[model.Reference](),
					Group:             model.Null[model.Reference](),
				}, []byte(`{"name":null,"description":null,"code":null,"externalCode":null,"shared":null,"archived":null,"operatingExpenses":null,"owner":null,"group":null}`)},
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

func fullUpdateRequest() expenseitem.UpdateRequest {
	return expenseitem.UpdateRequest{
		Name:              model.Value("SDK ExpenseItem updated"),
		Description:       model.Value(""),
		Code:              model.Value("expenseitem-updated"),
		ExternalCode:      model.Value("sdk-expenseitem-updated"),
		Shared:            model.Value(false),
		Archived:          model.Value(false),
		OperatingExpenses: model.Value(true),
		Owner: model.Value(model.Reference{Meta: model.Meta{
			Href: "https://api.moysklad.ru/api/remap/1.2/entity/employee/11111111-1111-1111-1111-111111111111",
		}}),
		Group: model.Value(model.Reference{Meta: model.Meta{
			Href: "https://api.moysklad.ru/api/remap/1.2/entity/group/22222222-2222-2222-2222-222222222222",
		}}),
	}
}
