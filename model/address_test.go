package model_test

import (
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestAddressJSON(t *testing.T) {
	const input = `{
		"addInfo":"корпус 2", "apartment":"007", "city":"Москва", "comment":"  без изменений  ",
		"fiasCode__ru":"001122", "house":"01Б", "postalCode":"001234", "street":"Тверская",
		"country":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/country/11111111-1111-1111-1111-111111111111","type":"country","mediaType":"application/json"}},
		"region":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/region/00000000-0000-0000-0000-000000000077","type":"region","mediaType":"application/json"}}
	}`
	want := model.Address{
		AddInfo: "корпус 2", Apartment: "007", City: "Москва", Comment: "  без изменений  ",
		FIASCodeRU: "001122", House: "01Б", PostalCode: "001234", Street: "Тверская",
		Country: &model.Reference{Meta: model.Meta{
			Href: "https://api.moysklad.ru/api/remap/1.2/entity/country/11111111-1111-1111-1111-111111111111",
			Type: "country", MediaType: "application/json",
		}},
		Region: &model.Reference{Meta: model.Meta{
			Href: "https://api.moysklad.ru/api/remap/1.2/entity/region/00000000-0000-0000-0000-000000000077",
			Type: "region", MediaType: "application/json",
		}},
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			var got model.Address
			if err := codec.unmarshal([]byte(input), &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("address = %+v, want %+v", got, want)
			}
			encoded, err := codec.marshal(want)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, encoded, input)
		})
	}
}

func TestAddressOptionalFieldsJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				input string
				want  model.Address
			}{
				{"absent", `{}`, model.Address{}},
				{"null references", `{"country":null,"region":null}`, model.Address{}},
				{"unknown", `{"city":"Омск","future":{"address":"ignored"}}`, model.Address{City: "Омск"}},
				{"empty strings", `{"addInfo":"","apartment":"","city":"","comment":"","fiasCode__ru":"","house":"","postalCode":"","street":""}`, model.Address{}},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var got model.Address
					if err := codec.unmarshal([]byte(tc.input), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.want) {
						t.Fatalf("address = %+v, want %+v", got, tc.want)
					}
				})
			}
			encoded, err := codec.marshal(model.Address{})
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, encoded, `{"addInfo":"","apartment":"","city":"","comment":"","fiasCode__ru":"","house":"","postalCode":"","street":""}`)
		})
	}
}

func TestAddressResponseContainerJSON(t *testing.T) {
	type response struct {
		Address *model.Address `json:"actualAddressFull,omitzero"`
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				input string
				want  *model.Address
			}{
				{`{}`, nil},
				{`{"actualAddressFull":null}`, nil},
				{`{"actualAddressFull":{}}`, &model.Address{}},
				{`{"actualAddressFull":{"city":"Казань","unknown":true}}`, &model.Address{City: "Казань"}},
			} {
				t.Run(tc.input, func(t *testing.T) {
					var got response
					if err := codec.unmarshal([]byte(tc.input), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got.Address, tc.want) {
						t.Fatalf("address = %+v, want %+v", got.Address, tc.want)
					}
				})
			}
		})
	}
}
