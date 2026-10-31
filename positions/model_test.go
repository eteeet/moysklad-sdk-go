package positions_test

import (
	"bytes"
	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

var jsonCodecs = []struct {
	name      string
	marshal   func(any) ([]byte, error)
	unmarshal func([]byte, any) error
}{
	{"json/v1", jsonv1.Marshal, jsonv1.Unmarshal},
	{"json/v2",
		func(v any) ([]byte, error) { return jsonv2.Marshal(v) },
		func(data []byte, v any) error { return jsonv2.Unmarshal(data, v) },
	},
}

func TestCoreJSON(t *testing.T) {
	const input = `{
		"meta":{"href":"https://example.test/positions/row-1","type":"supplyposition","mediaType":"application/json"},
		"id":"row-1","accountId":"account-1","future":{"enabled":true}
	}`
	want := positions.Core{
		Meta: model.Meta{
			Href: "https://example.test/positions/row-1", Type: "supplyposition", MediaType: "application/json",
		},
		ID: "row-1", AccountID: "account-1",
	}
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			var got positions.Core
			if err := codec.unmarshal([]byte(input), &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("core = %+v, want %+v", got, want)
			}
			encoded, err := codec.marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			jsonAssertEqual(t, encoded, `{"meta":{"href":"https://example.test/positions/row-1","type":"supplyposition","mediaType":"application/json"},"id":"row-1","accountId":"account-1"}`)
		})
	}
}

func TestItemJSON(t *testing.T) {
	const input = `{
		"meta":{"href":"https://example.test/positions/row-2","type":"demandposition","mediaType":"application/json"},
		"id":"row-2","accountId":"account-2",
		"assortment":{"meta":{"href":"https://example.test/entity/product/product-1","type":"product","mediaType":"application/json"}},
		"quantity":2.75,"future":"value"
	}`
	want := positions.Item{
		Core: positions.Core{
			Meta: model.Meta{
				Href: "https://example.test/positions/row-2", Type: "demandposition", MediaType: "application/json",
			},
			ID: "row-2", AccountID: "account-2",
		},
		Assortment: model.Reference{Meta: model.Meta{
			Href: "https://example.test/entity/product/product-1", Type: "product", MediaType: "application/json",
		}},
		Quantity: 2.75,
	}
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			var got positions.Item
			if err := codec.unmarshal([]byte(input), &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("item = %+v, want %+v", got, want)
			}
			encoded, err := codec.marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			jsonAssertEqual(t, encoded, `{
				"meta":{"href":"https://example.test/positions/row-2","type":"demandposition","mediaType":"application/json"},
				"id":"row-2","accountId":"account-2",
				"assortment":{"meta":{"href":"https://example.test/entity/product/product-1","type":"product","mediaType":"application/json"}},
				"quantity":2.75
			}`)
		})
	}
}

func TestCommercialJSON(t *testing.T) {
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				input string
				want  positions.Commercial
			}{
				{
					"fractional values and markup",
					`{
						"meta":{"href":"https://example.test/positions/row-3","type":"invoiceoutposition","mediaType":"application/json"},
						"id":"row-3","accountId":"account-3",
						"assortment":{"meta":{"href":"https://example.test/entity/service/service-1","type":"service","mediaType":"application/json"}},
						"quantity":1.25,"price":12345.67,"discount":-5.5,"vat":20,"vatEnabled":true,"future":{"value":1}
					}`,
					positions.Commercial{
						Item: positions.Item{
							Core: positions.Core{
								Meta: model.Meta{
									Href: "https://example.test/positions/row-3", Type: "invoiceoutposition", MediaType: "application/json",
								},
								ID: "row-3", AccountID: "account-3",
							},
							Assortment: model.Reference{Meta: model.Meta{
								Href: "https://example.test/entity/service/service-1", Type: "service", MediaType: "application/json",
							}},
							Quantity: 1.25,
						},
						Price: 12345.67, Discount: -5.5, VAT: 20, VATEnabled: true,
					},
				},
				{"absent fields", `{}`, positions.Commercial{}},
				{
					"integer zero and disabled VAT",
					`{"id":"","accountId":"","quantity":0,"price":0,"discount":0,"vat":0,"vatEnabled":false}`,
					positions.Commercial{},
				},
				{
					"decimal zero",
					`{"quantity":0.0,"price":0.0,"discount":0.0,"vat":0,"vatEnabled":false,"future":true}`,
					positions.Commercial{},
				},
				{"zero VAT enabled", `{"vat":0,"vatEnabled":true}`, positions.Commercial{VATEnabled: true}},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var got positions.Commercial
					if err := codec.unmarshal([]byte(tc.input), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.want) {
						t.Fatalf("commercial = %+v, want %+v", got, tc.want)
					}
				})
			}
			encoded, err := codec.marshal(positions.Commercial{
				Item: positions.Item{Quantity: 2.5}, Price: 101.25, Discount: -2.5, VATEnabled: true,
			})
			if err != nil {
				t.Fatal(err)
			}
			jsonAssertEqual(t, encoded, `{"meta":{},"id":"","accountId":"","assortment":{"meta":{}},"quantity":2.5,"price":101.25,"discount":-2.5,"vat":0,"vatEnabled":true}`)
		})
	}
}

func jsonAssertEqual(t *testing.T, data []byte, want string) {
	t.Helper()
	gotJSON, wantJSON := jsontext.Value(bytes.Clone(data)), jsontext.Value(want)
	if err := gotJSON.Canonicalize(); err != nil {
		t.Fatal(err)
	}
	if err := wantJSON.Canonicalize(); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotJSON, wantJSON) {
		t.Fatalf("JSON = %s, want %s", data, want)
	}
}
