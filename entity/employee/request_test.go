package employee_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/employee"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func reference(href, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: href, Type: model.MetaType(typ), MediaType: "application/json"}}
}

func fullCreateRequest() employee.CreateRequest {
	return employee.CreateRequest{
		LastName:     "Мойскладкин",
		FirstName:    model.Value("Петр"),
		MiddleName:   model.Value("Иванович"),
		Description:  model.Value("Описание"),
		ExternalCode: model.Value("EXT-1"),
		Archived:     model.Value(false),
		Shared:       model.Value(true),
		Email:        model.Value("petr@example.test"),
		Phone:        model.Value("+7(999)888-7766"),
		INN:          model.Value("222490425273"),
		Position:     model.Value("Директор"),
		Salary:       model.Value(employee.Salary{Value: 12500050}),
		Owner:        model.Value(reference("https://api.moysklad.ru/api/remap/1.2/entity/employee/00f76fbb-cd3b-11e8-ac12-00080000002d", "employee")),
		Group:        model.Value(reference("https://api.moysklad.ru/api/remap/1.2/entity/group/ffbc0889-cd3a-11e8-ac12-000700000002", "group")),
		Attributes: model.Value([]model.Attribute{{
			Meta:  model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/employee/metadata/attributes/ed14b498-cae3-11e8-9dd2-f3a300000044", Type: "attributemetadata", MediaType: "application/json"},
			Value: model.Value(jsontext.Value(`"Строковое доп поле"`)),
		}}),
	}
}

func fullUpdateRequest() employee.UpdateRequest {
	return employee.UpdateRequest{
		LastName:     "Мойскладкин",
		FirstName:    model.Value(""),
		MiddleName:   model.Value("Петрович"),
		Description:  model.Value("Новое описание"),
		ExternalCode: model.Value("EXT-2"),
		Archived:     model.Value(true),
		Shared:       model.Value(false),
		Email:        model.Value(""),
		Phone:        model.Value("+7(999)000-0000"),
		INN:          model.Value(""),
		Position:     model.Value("Менеджер"),
		Salary:       model.Value(employee.Salary{}),
		Owner:        model.Value(reference("https://api.moysklad.ru/api/remap/1.2/entity/employee/00f76fbb-cd3b-11e8-ac12-00080000002d", "employee")),
		Group:        model.Value(reference("https://api.moysklad.ru/api/remap/1.2/entity/group/ffbc0889-cd3a-11e8-ac12-000700000002", "group")),
		Attributes:   model.Value([]model.Attribute(nil)),
	}
}

func TestRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for _, tc := range []struct {
				name         string
				input        any
				golden, want string
			}{
				{"minimal Create", employee.CreateRequest{LastName: "Мойскладкин"}, "../../testdata/custom/employee/create-minimal.golden.json", ""},
				{"full Create", fullCreateRequest(), "../../testdata/custom/employee/create-full.golden.json", ""},
				{"full Update", fullUpdateRequest(), "../../testdata/custom/employee/update-full.golden.json", ""},
				{"minimal Update", employee.UpdateRequest{LastName: "Л"}, "", `{"lastName":"Л"}`},
				{"archived false", employee.UpdateRequest{LastName: "Л", Archived: model.Value(false)}, "", `{"lastName":"Л","archived":false}`},
				{"salary zero", employee.UpdateRequest{LastName: "Л", Salary: model.Value(employee.Salary{})}, "", `{"lastName":"Л","salary":{"value":0}}`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					want := []byte(tc.want)
					if tc.golden != "" {
						var err error
						if want, err = os.ReadFile(tc.golden); err != nil {
							t.Fatal(err)
						}
					}
					got, err := codec.Marshal(tc.input)
					if err != nil {
						t.Fatal(err)
					}
					documenttest.JSONEqual(t, got, want)
				})
			}
		})
	}
}

// Поля только для чтения и вне scope (фото, роль, логин) в запросах отсутствуют.
func TestRequestReadOnlyFields(t *testing.T) {
	for _, typ := range []reflect.Type{reflect.TypeFor[employee.CreateRequest](), reflect.TypeFor[employee.UpdateRequest]()} {
		for _, name := range []string{"ID", "AccountID", "Meta", "Name", "FullName", "ShortFio", "UID", "Created", "Updated", "Cashiers", "Image"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}
