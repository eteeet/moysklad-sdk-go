package model_test

import (
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestDocumentRateJSON(t *testing.T) {
	zero, one, integer, fraction := 0.0, 1.0, 90.0, 93.125
	currency := &model.Reference{Meta: model.Meta{
		Href:         "https://example.test/entity/currency/rub",
		MetadataHref: "https://example.test/entity/currency/metadata",
		Type:         "currency",
		MediaType:    "application/json",
		UUIDHref:     "https://example.test/app/currency/rub",
	}}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				input string
				want  model.DocumentRate
				json  string
			}{
				{
					"currency and fractional value",
					`{"currency":{"meta":{"href":"https://example.test/entity/currency/rub","metadataHref":"https://example.test/entity/currency/metadata","type":"currency","mediaType":"application/json","uuidHref":"https://example.test/app/currency/rub"}},"value":93.125}`,
					model.DocumentRate{Currency: currency, Value: &fraction},
					`{"currency":{"meta":{"href":"https://example.test/entity/currency/rub","metadataHref":"https://example.test/entity/currency/metadata","type":"currency","mediaType":"application/json","uuidHref":"https://example.test/app/currency/rub"}},"value":93.125}`,
				},
				{"integer value", `{"value":90}`, model.DocumentRate{Value: &integer}, `{"value":90}`},
				{"explicit zero", `{"value":0}`, model.DocumentRate{Value: &zero}, `{"value":0}`},
				{"explicit one", `{"value":1}`, model.DocumentRate{Value: &one}, `{"value":1}`},
				{
					"missing value",
					`{"currency":{"meta":{"href":"https://example.test/entity/currency/rub","metadataHref":"https://example.test/entity/currency/metadata","type":"currency","mediaType":"application/json","uuidHref":"https://example.test/app/currency/rub"}}}`,
					model.DocumentRate{Currency: currency},
					`{"currency":{"meta":{"href":"https://example.test/entity/currency/rub","metadataHref":"https://example.test/entity/currency/metadata","type":"currency","mediaType":"application/json","uuidHref":"https://example.test/app/currency/rub"}}}`,
				},
				{"absent fields", `{}`, model.DocumentRate{}, `{}`},
				{"null currency", `{"currency":null,"value":1}`, model.DocumentRate{Value: &one}, `{"value":1}`},
				{"null value", `{"value":null}`, model.DocumentRate{}, `{}`},
				{"unknown field", `{"value":93.125,"future":{"currency":"ignored"}}`, model.DocumentRate{Value: &fraction}, `{"value":93.125}`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var got model.DocumentRate
					if err := codec.unmarshal([]byte(tc.input), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.want) {
						t.Fatalf("rate = %+v, want %+v", got, tc.want)
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

func TestDocumentRateResponseJSON(t *testing.T) {
	type response struct {
		Rate *model.DocumentRate `json:"rate,omitzero"`
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				input string
				want  response
			}{
				{`{}`, response{}},
				{`{"rate":null}`, response{}},
				{`{"rate":{}}`, response{Rate: &model.DocumentRate{}}},
			} {
				t.Run(tc.input, func(t *testing.T) {
					var got response
					if err := codec.unmarshal([]byte(tc.input), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.want) {
						t.Fatalf("response = %+v, want %+v", got, tc.want)
					}
				})
			}
		})
	}
}

func TestDocumentRateRequestJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value model.DocumentRateRequest
				want  string
			}{
				{"unset", model.DocumentRateRequest{}, `{}`},
				{
					"full",
					model.DocumentRateRequest{
						Currency: model.Value(model.Reference{Meta: model.Meta{
							Href:         "https://example.test/entity/currency/usd",
							MetadataHref: "https://example.test/entity/currency/metadata",
							Type:         "currency",
							MediaType:    "application/json",
							UUIDHref:     "https://example.test/app/currency/usd",
						}}),
						Value: model.Value(93.125),
					},
					`{"currency":{"meta":{"href":"https://example.test/entity/currency/usd","metadataHref":"https://example.test/entity/currency/metadata","type":"currency","mediaType":"application/json","uuidHref":"https://example.test/app/currency/usd"}},"value":93.125}`,
				},
				{
					"currency only",
					model.DocumentRateRequest{Currency: model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/entity/currency/rub"}})},
					`{"currency":{"meta":{"href":"https://example.test/entity/currency/rub"}}}`,
				},
				{"value only", model.DocumentRateRequest{Value: model.Value(90.0)}, `{"value":90}`},
				{"explicit zero", model.DocumentRateRequest{Value: model.Value(0.0)}, `{"value":0}`},
				{"explicit one", model.DocumentRateRequest{Value: model.Value(1.0)}, `{"value":1}`},
				{
					"null fields",
					model.DocumentRateRequest{Currency: model.Null[model.Reference](), Value: model.Null[float64]()},
					`{"currency":null,"value":null}`,
				},
			} {
				t.Run(tc.name, func(t *testing.T) {
					encoded, err := codec.marshal(tc.value)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, encoded, tc.want)
					var got model.DocumentRateRequest
					if err := codec.unmarshal([]byte(tc.want), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.value) {
						t.Fatalf("rate request = %+v, want %+v", got, tc.value)
					}
				})
			}
		})
	}
}

func TestDocumentRateRequestFieldJSON(t *testing.T) {
	type request struct {
		Rate model.Field[model.DocumentRateRequest] `json:"rate,omitzero"`
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value request
				want  string
			}{
				{"unset", request{}, `{}`},
				{"null", request{Rate: model.Null[model.DocumentRateRequest]()}, `{"rate":null}`},
				{"empty object", request{Rate: model.Value(model.DocumentRateRequest{})}, `{"rate":{}}`},
				{
					"nested fields",
					request{Rate: model.Value(model.DocumentRateRequest{Currency: model.Null[model.Reference](), Value: model.Value(0.0)})},
					`{"rate":{"currency":null,"value":0}}`,
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
