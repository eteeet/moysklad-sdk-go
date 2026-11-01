package positions_test

import (
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

func TestCoreRequestJSON(t *testing.T) {
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value positions.CoreRequest
				want  string
			}{
				{"unset", positions.CoreRequest{}, `{}`},
				{"existing row selector", positions.CoreRequest{ID: model.Value("existing-row")}, `{"id":"existing-row"}`},
				{"empty ID", positions.CoreRequest{ID: model.Value("")}, `{"id":""}`},
				{"null", positions.CoreRequest{ID: model.Null[string]()}, `{"id":null}`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					encoded, err := codec.marshal(tc.value)
					if err != nil {
						t.Fatal(err)
					}
					jsonAssertEqual(t, encoded, tc.want)
					var got positions.CoreRequest
					if err := codec.unmarshal([]byte(tc.want), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.value) {
						t.Fatalf("core request = %+v, want %+v", got, tc.value)
					}
				})
			}
			var got positions.CoreRequest
			if err := codec.unmarshal([]byte(`{"id":"existing-row","meta":{"href":"https://example.test/positions/existing-row"},"accountId":"account-1"}`), &got); err != nil {
				t.Fatal(err)
			}
			encoded, err := codec.marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			jsonAssertEqual(t, encoded, `{"id":"existing-row"}`)
		})
	}
}

func TestItemRequestJSON(t *testing.T) {
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value positions.ItemRequest
				want  string
			}{
				{"unset", positions.ItemRequest{}, `{}`},
				{
					"new row",
					positions.ItemRequest{
						Assortment: model.Value(model.Reference{Meta: model.Meta{
							Href: "https://example.test/entity/product/product-1", Type: "product", MediaType: "application/json",
						}}),
						Quantity: model.Value(1.5),
					},
					`{"assortment":{"meta":{"href":"https://example.test/entity/product/product-1","type":"product","mediaType":"application/json"}},"quantity":1.5}`,
				},
				{
					"sparse existing row",
					positions.ItemRequest{
						CoreRequest: positions.CoreRequest{ID: model.Value("existing-row")},
						Quantity:    model.Value(0.25),
					},
					`{"id":"existing-row","quantity":0.25}`,
				},
				{"zero quantity", positions.ItemRequest{Quantity: model.Value(0.0)}, `{"quantity":0}`},
				{
					"null fields",
					positions.ItemRequest{Assortment: model.Null[model.Reference](), Quantity: model.Null[float64]()},
					`{"assortment":null,"quantity":null}`,
				},
			} {
				t.Run(tc.name, func(t *testing.T) {
					encoded, err := codec.marshal(tc.value)
					if err != nil {
						t.Fatal(err)
					}
					jsonAssertEqual(t, encoded, tc.want)
					var got positions.ItemRequest
					if err := codec.unmarshal([]byte(tc.want), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.value) {
						t.Fatalf("item request = %+v, want %+v", got, tc.value)
					}
				})
			}
		})
	}
}

func TestPositionRequestJSON(t *testing.T) {
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value positions.PositionRequest
				want  string
			}{
				{"unset", positions.PositionRequest{}, `{}`},
				{
					"source slot only",
					positions.PositionRequest{SourceSlot: model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/slots/source-1", Type: "slot"}})},
					`{"sourceSlot":{"meta":{"href":"https://example.test/slots/source-1","type":"slot"}}}`,
				},
				{
					"target slot only",
					positions.PositionRequest{TargetSlot: model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/slots/target-1", Type: "slot"}})},
					`{"targetSlot":{"meta":{"href":"https://example.test/slots/target-1","type":"slot"}}}`,
				},
				{
					"all writable fields",
					positions.PositionRequest{
						ItemRequest: positions.ItemRequest{
							CoreRequest: positions.CoreRequest{ID: model.Value("existing-row")},
							Assortment: model.Value(model.Reference{Meta: model.Meta{
								Href: "https://example.test/entity/product/product-1", Type: "product", MediaType: "application/json",
							}}),
							Quantity: model.Value(1.25),
						},
						Price: model.Value(12345.67), Discount: model.Value(-5.5),
						VAT: model.Value(20), VATEnabled: model.Value(true),
						SourceSlot: model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/slots/source-1", Type: "slot"}}),
						TargetSlot: model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/slots/target-1", Type: "slot"}}),
					},
					`{"id":"existing-row","assortment":{"meta":{"href":"https://example.test/entity/product/product-1","type":"product","mediaType":"application/json"}},"quantity":1.25,"price":12345.67,"discount":-5.5,"vat":20,"vatEnabled":true,"sourceSlot":{"meta":{"href":"https://example.test/slots/source-1","type":"slot"}},"targetSlot":{"meta":{"href":"https://example.test/slots/target-1","type":"slot"}}}`,
				},
				{
					"minimal new row",
					positions.PositionRequest{ItemRequest: positions.ItemRequest{
						Assortment: model.Value(model.Reference{Meta: model.Meta{
							Href: "https://example.test/entity/service/service-1", Type: "service", MediaType: "application/json",
						}}),
						Quantity: model.Value(2.5),
					}},
					`{"assortment":{"meta":{"href":"https://example.test/entity/service/service-1","type":"service","mediaType":"application/json"}},"quantity":2.5}`,
				},
				{
					"ID only",
					positions.PositionRequest{ItemRequest: positions.ItemRequest{
						CoreRequest: positions.CoreRequest{ID: model.Value("existing-row")},
					}},
					`{"id":"existing-row"}`,
				},
				{
					"sparse quantity update",
					positions.PositionRequest{ItemRequest: positions.ItemRequest{
						CoreRequest: positions.CoreRequest{ID: model.Value("existing-row")},
						Quantity:    model.Value(0.125),
					}},
					`{"id":"existing-row","quantity":0.125}`,
				},
				{
					"explicit zero and false",
					positions.PositionRequest{
						ItemRequest: positions.ItemRequest{Quantity: model.Value(0.0)},
						Price:       model.Value(0.0), Discount: model.Value(0.0),
						VAT: model.Value(0), VATEnabled: model.Value(false),
					},
					`{"quantity":0,"price":0,"discount":0,"vat":0,"vatEnabled":false}`,
				},
				{
					"zero VAT enabled",
					positions.PositionRequest{VAT: model.Value(0), VATEnabled: model.Value(true)},
					`{"vat":0,"vatEnabled":true}`,
				},
				{
					"null fields",
					positions.PositionRequest{
						ItemRequest: positions.ItemRequest{
							CoreRequest: positions.CoreRequest{ID: model.Null[string]()},
							Assortment:  model.Null[model.Reference](), Quantity: model.Null[float64](),
						},
						Price: model.Null[float64](), Discount: model.Null[float64](),
						VAT: model.Null[int](), VATEnabled: model.Null[bool](),
						SourceSlot: model.Null[model.Reference](), TargetSlot: model.Null[model.Reference](),
					},
					`{"id":null,"assortment":null,"quantity":null,"price":null,"discount":null,"vat":null,"vatEnabled":null,"sourceSlot":null,"targetSlot":null}`,
				},
			} {
				t.Run(tc.name, func(t *testing.T) {
					encoded, err := codec.marshal(tc.value)
					if err != nil {
						t.Fatal(err)
					}
					jsonAssertEqual(t, encoded, tc.want)
					var got positions.PositionRequest
					if err := codec.unmarshal([]byte(tc.want), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.value) {
						t.Fatalf("position request = %+v, want %+v", got, tc.value)
					}
				})
			}
			var got positions.PositionRequest
			if err := codec.unmarshal([]byte(`{
				"id":"existing-row","quantity":1.5,"price":2.75,"discount":5,"vat":20,"vatEnabled":true,
				"meta":{"href":"https://example.test/positions/existing-row"},"accountId":"account-1",
				"overhead":12.5,"future":true
			}`), &got); err != nil {
				t.Fatal(err)
			}
			encoded, err := codec.marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			jsonAssertEqual(t, encoded, `{"id":"existing-row","quantity":1.5,"price":2.75,"discount":5,"vat":20,"vatEnabled":true}`)
		})
	}
}

func TestEmbeddedPositionsJSON(t *testing.T) {
	type jsonDocumentRequest struct {
		Positions model.Field[[]positions.PositionRequest] `json:"positions,omitzero"`
	}
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value jsonDocumentRequest
				want  string
			}{
				{"unset", jsonDocumentRequest{}, `{}`},
				{"null", jsonDocumentRequest{Positions: model.Null[[]positions.PositionRequest]()}, `{"positions":null}`},
				{"empty slice", jsonDocumentRequest{Positions: model.Value([]positions.PositionRequest{})}, `{"positions":[]}`},
				{"nil slice", jsonDocumentRequest{Positions: model.Value[[]positions.PositionRequest](nil)}, `{"positions":[]}`},
				{
					"replacement array",
					jsonDocumentRequest{Positions: model.Value([]positions.PositionRequest{
						{ItemRequest: positions.ItemRequest{
							CoreRequest: positions.CoreRequest{ID: model.Value("existing-row")},
							Quantity:    model.Value(0.5),
						}},
						{ItemRequest: positions.ItemRequest{
							Assortment: model.Value(model.Reference{Meta: model.Meta{
								Href: "https://example.test/entity/product/product-2", Type: "product", MediaType: "application/json",
							}}),
							Quantity: model.Value(1.25),
						}, Price: model.Value(99.75), VAT: model.Value(0), VATEnabled: model.Value(false),
							TargetSlot: model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/slots/target-1", Type: "slot"}})},
					})},
					`{"positions":[{"id":"existing-row","quantity":0.5},{"assortment":{"meta":{"href":"https://example.test/entity/product/product-2","type":"product","mediaType":"application/json"}},"quantity":1.25,"price":99.75,"vat":0,"vatEnabled":false,"targetSlot":{"meta":{"href":"https://example.test/slots/target-1","type":"slot"}}}]}`,
				},
			} {
				t.Run(tc.name, func(t *testing.T) {
					encoded, err := codec.marshal(tc.value)
					if err != nil {
						t.Fatal(err)
					}
					jsonAssertEqual(t, encoded, tc.want)
				})
			}
		})
	}
}
