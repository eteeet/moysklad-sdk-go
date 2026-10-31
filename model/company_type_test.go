package model_test

import (
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestCompanyTypeJSON(t *testing.T) {
	var zero model.CompanyType
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				value model.CompanyType
				want  string
			}{
				{model.CompanyTypeLegal, `"legal"`},
				{model.CompanyTypeEntrepreneur, `"entrepreneur"`},
				{model.CompanyTypeIndividual, `"individual"`},
				{model.CompanyTypeLegalUZ, `"legalUZ"`},
				{model.CompanyTypeEntrepreneurUZ, `"entrepreneurUZ"`},
				{model.CompanyTypeIndividualUZ, `"individualUZ"`},
				{model.CompanyTypeLegalKZ, `"legalKZ"`},
				{model.CompanyTypeEntrepreneurKZ, `"entrepreneurKZ"`},
				{model.CompanyTypeIndividualKZ, `"individualKZ"`},
				{model.CompanyType("future-region-type"), `"future-region-type"`},
				{zero, `""`},
			} {
				t.Run(tc.want, func(t *testing.T) {
					encoded, err := codec.marshal(tc.value)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, encoded, tc.want)
					var got model.CompanyType
					if err := codec.unmarshal([]byte(tc.want), &got); err != nil {
						t.Fatal(err)
					}
					if got != tc.value {
						t.Fatalf("companyType = %q, want %q", got, tc.value)
					}
				})
			}
		})
	}
}

func TestCompanyTypeFieldJSON(t *testing.T) {
	type request struct {
		CompanyType model.Field[model.CompanyType] `json:"companyType,omitzero"`
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				field model.Field[model.CompanyType]
				want  string
			}{
				{"unset", model.Field[model.CompanyType]{}, `{}`},
				{"empty", model.Value(model.CompanyType("")), `{"companyType":""}`},
				{"null", model.Null[model.CompanyType](), `{"companyType":null}`},
				{"known", model.Value(model.CompanyTypeIndividualKZ), `{"companyType":"individualKZ"}`},
				{"unknown", model.Value(model.CompanyType("future")), `{"companyType":"future"}`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					encoded, err := codec.marshal(request{CompanyType: tc.field})
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, encoded, tc.want)
					var got request
					if err := codec.unmarshal([]byte(tc.want), &got); err != nil {
						t.Fatal(err)
					}
					if got.CompanyType != tc.field {
						t.Fatalf("companyType = %+v, want %+v", got.CompanyType, tc.field)
					}
				})
			}
		})
	}
}
