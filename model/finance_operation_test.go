package model_test

import (
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestFinanceOperationJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				input string
				want  model.FinanceOperation
				json  string
			}{
				{
					"full meta and fractional sum",
					`{"meta":{"href":"https://example.test/entity/demand/one","metadataHref":"https://example.test/entity/demand/metadata","type":"demand","mediaType":"application/json","uuidHref":"https://example.test/app/demand/one"},"linkedSum":12345.625}`,
					model.FinanceOperation{
						Meta: model.Meta{
							Href:         "https://example.test/entity/demand/one",
							MetadataHref: "https://example.test/entity/demand/metadata",
							Type:         "demand",
							MediaType:    "application/json",
							UUIDHref:     "https://example.test/app/demand/one",
						},
						LinkedSum: 12345.625,
					},
					`{"meta":{"href":"https://example.test/entity/demand/one","metadataHref":"https://example.test/entity/demand/metadata","type":"demand","mediaType":"application/json","uuidHref":"https://example.test/app/demand/one"},"linkedSum":12345.625}`,
				},
				{"missing sum", `{"meta":{"type":"supply"}}`, model.FinanceOperation{Meta: model.Meta{Type: "supply"}}, `{"meta":{"type":"supply"},"linkedSum":0}`},
				{"explicit zero", `{"meta":{"type":"demand"},"linkedSum":0}`, model.FinanceOperation{Meta: model.Meta{Type: "demand"}}, `{"meta":{"type":"demand"},"linkedSum":0}`},
				{"decimal zero", `{"meta":{"type":"demand"},"linkedSum":0.0}`, model.FinanceOperation{Meta: model.Meta{Type: "demand"}}, `{"meta":{"type":"demand"},"linkedSum":0}`},
				{"unknown type", `{"meta":{"type":"futureoperation"},"linkedSum":42}`, model.FinanceOperation{Meta: model.Meta{Type: "futureoperation"}, LinkedSum: 42}, `{"meta":{"type":"futureoperation"},"linkedSum":42}`},
				{
					"expanded fields ignored",
					`{"meta":{"type":"invoiceout"},"linkedSum":987.125,"id":"ignored","name":"ignored","positions":{"rows":[{"price":100}]},"payments":[{"sum":1}],"future":true}`,
					model.FinanceOperation{Meta: model.Meta{Type: "invoiceout"}, LinkedSum: 987.125},
					`{"meta":{"type":"invoiceout"},"linkedSum":987.125}`,
				},
				{"absent fields", `{}`, model.FinanceOperation{}, `{"meta":{},"linkedSum":0}`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var got model.FinanceOperation
					if err := codec.unmarshal([]byte(tc.input), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.want) {
						t.Fatalf("operation = %+v, want %+v", got, tc.want)
					}
					encoded, err := codec.marshal(got)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, encoded, tc.json)
				})
			}
		})
	}
}

func TestFinanceOperationRequestJSON(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value model.FinanceOperationRequest
				want  string
			}{
				{
					"meta only",
					model.FinanceOperationRequest{Meta: model.Meta{
						Href:         "https://example.test/entity/customerorder/one",
						MetadataHref: "https://example.test/entity/customerorder/metadata",
						Type:         "customerorder",
						MediaType:    "application/json",
						UUIDHref:     "https://example.test/app/customerorder/one",
					}},
					`{"meta":{"href":"https://example.test/entity/customerorder/one","metadataHref":"https://example.test/entity/customerorder/metadata","type":"customerorder","mediaType":"application/json","uuidHref":"https://example.test/app/customerorder/one"}}`,
				},
				{"explicit zero", model.FinanceOperationRequest{Meta: model.Meta{Type: "supply"}, LinkedSum: model.Value(0.0)}, `{"meta":{"type":"supply"},"linkedSum":0}`},
				{"fractional sum", model.FinanceOperationRequest{Meta: model.Meta{Type: "demand"}, LinkedSum: model.Value(12345.625)}, `{"meta":{"type":"demand"},"linkedSum":12345.625}`},
				{"null sum", model.FinanceOperationRequest{Meta: model.Meta{Type: "demand"}, LinkedSum: model.Null[float64]()}, `{"meta":{"type":"demand"},"linkedSum":null}`},
				{"unknown type", model.FinanceOperationRequest{Meta: model.Meta{Type: "futureoperation"}, LinkedSum: model.Value(42.5)}, `{"meta":{"type":"futureoperation"},"linkedSum":42.5}`},
				{"zero meta and unset sum", model.FinanceOperationRequest{}, `{"meta":{}}`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					encoded, err := codec.marshal(tc.value)
					if err != nil {
						t.Fatal(err)
					}
					assertJSONEqual(t, encoded, tc.want)
					var got model.FinanceOperationRequest
					if err := codec.unmarshal([]byte(tc.want), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.value) {
						t.Fatalf("operation request = %+v, want %+v", got, tc.value)
					}
				})
			}
		})
	}
}

func TestFinanceOperationRequestFieldJSON(t *testing.T) {
	type request struct {
		Operations model.Field[[]model.FinanceOperationRequest] `json:"operations,omitzero"`
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value request
				want  string
			}{
				{"unset", request{}, `{}`},
				{"null", request{Operations: model.Null[[]model.FinanceOperationRequest]()}, `{"operations":null}`},
				{"empty slice", request{Operations: model.Value([]model.FinanceOperationRequest{})}, `{"operations":[]}`},
				{
					"multiple operations",
					request{Operations: model.Value([]model.FinanceOperationRequest{
						{Meta: model.Meta{Href: "https://example.test/entity/demand/one", Type: "demand"}, LinkedSum: model.Value(0.0)},
						{Meta: model.Meta{Href: "https://example.test/entity/supply/two", Type: "supply"}},
						{Meta: model.Meta{Type: "futureoperation"}, LinkedSum: model.Null[float64]()},
					})},
					`{"operations":[{"meta":{"href":"https://example.test/entity/demand/one","type":"demand"},"linkedSum":0},{"meta":{"href":"https://example.test/entity/supply/two","type":"supply"}},{"meta":{"type":"futureoperation"},"linkedSum":null}]}`,
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
			t.Run("nil slice", func(t *testing.T) {
				encoded, err := codec.marshal(request{Operations: model.Value[[]model.FinanceOperationRequest](nil)})
				if err != nil {
					t.Fatal(err)
				}
				assertJSONEqual(t, encoded, `{"operations":[]}`)
			})
		})
	}
}
