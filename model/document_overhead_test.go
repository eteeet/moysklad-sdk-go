package model_test

import (
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestDocumentOverheadJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				input string
				want  model.DocumentOverhead
				json  string
			}{
				{"weight and fractional sum", `{"sum":123.75,"distribution":"weight"}`, model.DocumentOverhead{Sum: 123.75, Distribution: model.OverheadDistributionWeight}, `{"sum":123.75,"distribution":"weight"}`},
				{"volume and zero", `{"sum":0,"distribution":"volume"}`, model.DocumentOverhead{Distribution: model.OverheadDistributionVolume}, `{"sum":0,"distribution":"volume"}`},
				{"price and decimal zero", `{"sum":0.0,"distribution":"price"}`, model.DocumentOverhead{Distribution: model.OverheadDistributionPrice}, `{"sum":0,"distribution":"price"}`},
				{"unknown distribution", `{"sum":42,"distribution":"future"}`, model.DocumentOverhead{Sum: 42, Distribution: "future"}, `{"sum":42,"distribution":"future"}`},
				{"absent fields", `{}`, model.DocumentOverhead{}, `{"sum":0,"distribution":""}`},
				{"empty distribution", `{"sum":0,"distribution":""}`, model.DocumentOverhead{}, `{"sum":0,"distribution":""}`},
				{"unknown field", `{"sum":987.125,"future":{"sum":10}}`, model.DocumentOverhead{Sum: 987.125}, `{"sum":987.125,"distribution":""}`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var got model.DocumentOverhead
					if err := codec.unmarshal([]byte(tc.input), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.want) {
						t.Fatalf("overhead = %+v, want %+v", got, tc.want)
					}
					encoded, err := codec.marshal(tc.want)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, encoded, tc.json)
				})
			}
		})
	}
}

func TestDocumentOverheadRequestJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value model.DocumentOverheadRequest
				want  string
			}{
				{"unset", model.DocumentOverheadRequest{}, `{}`},
				{
					"full",
					model.DocumentOverheadRequest{Sum: model.Value(123.75), Distribution: model.Value(model.OverheadDistributionWeight)},
					`{"sum":123.75,"distribution":"weight"}`,
				},
				{"sum only zero", model.DocumentOverheadRequest{Sum: model.Value(0.0)}, `{"sum":0}`},
				{"distribution only volume", model.DocumentOverheadRequest{Distribution: model.Value(model.OverheadDistributionVolume)}, `{"distribution":"volume"}`},
				{"distribution only price", model.DocumentOverheadRequest{Distribution: model.Value(model.OverheadDistributionPrice)}, `{"distribution":"price"}`},
				{"unknown distribution", model.DocumentOverheadRequest{Distribution: model.Value(model.OverheadDistribution("future"))}, `{"distribution":"future"}`},
				{"empty distribution", model.DocumentOverheadRequest{Distribution: model.Value(model.OverheadDistribution(""))}, `{"distribution":""}`},
				{
					"null fields",
					model.DocumentOverheadRequest{Sum: model.Null[float64](), Distribution: model.Null[model.OverheadDistribution]()},
					`{"sum":null,"distribution":null}`,
				},
			} {
				t.Run(tc.name, func(t *testing.T) {
					encoded, err := codec.marshal(tc.value)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, encoded, tc.want)
					var got model.DocumentOverheadRequest
					if err := codec.unmarshal([]byte(tc.want), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.value) {
						t.Fatalf("overhead request = %+v, want %+v", got, tc.value)
					}
				})
			}
		})
	}
}

func TestDocumentOverheadRequestFieldJSON(t *testing.T) {
	type request struct {
		Overhead model.Field[model.DocumentOverheadRequest] `json:"overhead,omitzero"`
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value request
				want  string
			}{
				{"unset", request{}, `{}`},
				{"null", request{Overhead: model.Null[model.DocumentOverheadRequest]()}, `{"overhead":null}`},
				{"empty object", request{Overhead: model.Value(model.DocumentOverheadRequest{})}, `{"overhead":{}}`},
				{"partial object", request{Overhead: model.Value(model.DocumentOverheadRequest{Sum: model.Value(0.0)})}, `{"overhead":{"sum":0}}`},
				{
					"nested fields",
					request{Overhead: model.Value(model.DocumentOverheadRequest{
						Sum: model.Null[float64](), Distribution: model.Value(model.OverheadDistribution("")),
					})},
					`{"overhead":{"sum":null,"distribution":""}}`,
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
