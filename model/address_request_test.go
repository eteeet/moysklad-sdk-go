package model_test

import (
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestAddressRequestJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value model.AddressRequest
				want  string
			}{
				{"unset", model.AddressRequest{}, `{}`},
				{
					"full",
					model.AddressRequest{
						AddInfo: model.Value("корпус 2"), Apartment: model.Value("007"),
						City: model.Value("Москва"), Comment: model.Value("  без изменений  "),
						FIASCodeRU: model.Value("001122"), House: model.Value("01Б"),
						PostalCode: model.Value("001234"), Street: model.Value("Тверская"),
						Country: model.Value(model.Reference{Meta: model.Meta{
							Href: "https://api.moysklad.ru/api/remap/1.2/entity/country/11111111-1111-1111-1111-111111111111",
							Type: "country", MediaType: "application/json",
						}}),
						Region: model.Value(model.Reference{Meta: model.Meta{
							Href: "https://api.moysklad.ru/api/remap/1.2/entity/region/00000000-0000-0000-0000-000000000077",
							Type: "region", MediaType: "application/json",
						}}),
					},
					`{
						"addInfo":"корпус 2", "apartment":"007", "city":"Москва", "comment":"  без изменений  ",
						"fiasCode__ru":"001122", "house":"01Б", "postalCode":"001234", "street":"Тверская",
						"country":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/country/11111111-1111-1111-1111-111111111111","type":"country","mediaType":"application/json"}},
						"region":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/region/00000000-0000-0000-0000-000000000077","type":"region","mediaType":"application/json"}}
					}`,
				},
				{"partial", model.AddressRequest{PostalCode: model.Value("001234")}, `{"postalCode":"001234"}`},
				{
					"empty strings",
					model.AddressRequest{
						AddInfo: model.Value(""), Apartment: model.Value(""), City: model.Value(""),
						Comment: model.Value(""), FIASCodeRU: model.Value(""), House: model.Value(""),
						PostalCode: model.Value(""), Street: model.Value(""),
					},
					`{"addInfo":"","apartment":"","city":"","comment":"","fiasCode__ru":"","house":"","postalCode":"","street":""}`,
				},
				{
					"null fields",
					model.AddressRequest{
						AddInfo: model.Null[string](), Apartment: model.Null[string](), City: model.Null[string](),
						Comment: model.Null[string](), FIASCodeRU: model.Null[string](), House: model.Null[string](),
						PostalCode: model.Null[string](), Street: model.Null[string](),
						Country: model.Null[model.Reference](), Region: model.Null[model.Reference](),
					},
					`{"addInfo":null,"apartment":null,"city":null,"comment":null,"fiasCode__ru":null,"country":null,"house":null,"postalCode":null,"region":null,"street":null}`,
				},
			} {
				t.Run(tc.name, func(t *testing.T) {
					encoded, err := codec.marshal(tc.value)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, encoded, tc.want)
					var got model.AddressRequest
					if err := codec.unmarshal([]byte(tc.want), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.value) {
						t.Fatalf("address request = %+v, want %+v", got, tc.value)
					}
				})
			}
		})
	}
}

func TestAddressRequestFieldJSON(t *testing.T) {
	type request struct {
		Address     model.Field[string]               `json:"actualAddress,omitzero"`
		AddressFull model.Field[model.AddressRequest] `json:"actualAddressFull,omitzero"`
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value request
				want  string
			}{
				{"unset", request{}, `{}`},
				{"null", request{AddressFull: model.Null[model.AddressRequest]()}, `{"actualAddressFull":null}`},
				{"empty object", request{AddressFull: model.Value(model.AddressRequest{})}, `{"actualAddressFull":{}}`},
				{
					"nested fields",
					request{AddressFull: model.Value(model.AddressRequest{
						City: model.Value(""), Country: model.Null[model.Reference](),
						Region: model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/region/077"}}),
					})},
					`{"actualAddressFull":{"city":"","country":null,"region":{"meta":{"href":"https://example.test/region/077"}}}}`,
				},
				{"string only", request{Address: model.Value("  строковый адрес  ")}, `{"actualAddress":"  строковый адрес  "}`},
				{"empty string", request{Address: model.Value("")}, `{"actualAddress":""}`},
				{
					"both forms",
					request{
						Address: model.Value("строковый адрес"),
						AddressFull: model.Value(model.AddressRequest{
							City: model.Value("Москва"), AddInfo: model.Value("собственное значение"),
						}),
					},
					`{"actualAddress":"строковый адрес","actualAddressFull":{"city":"Москва","addInfo":"собственное значение"}}`,
				},
			} {
				t.Run(tc.name, func(t *testing.T) {
					encoded, err := codec.marshal(tc.value)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, encoded, tc.want)
					var got request
					if err := codec.unmarshal([]byte(tc.want), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.value) {
						t.Fatalf("request = %+v, want %+v", got, tc.value)
					}
				})
			}
		})
	}
}
