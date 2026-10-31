package positions_test

import (
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

func TestMovementRequestJSON(t *testing.T) {
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value positions.MovementRequest
				want  string
			}{
				{
					"all writable fields",
					positions.MovementRequest{
						ItemRequest: positions.ItemRequest{
							CoreRequest: positions.CoreRequest{ID: model.Value("existing-row")},
							Assortment: model.Value(model.Reference{Meta: model.Meta{
								Href: "https://example.test/entity/product/product-1", Type: "product",
							}}),
							Quantity: model.Value(1.25),
						},
						Price:      model.Value(12345.67),
						SourceSlot: model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/slots/source-1", Type: "slot"}}),
						TargetSlot: model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/slots/target-1", Type: "slot"}}),
					},
					`{"id":"existing-row","assortment":{"meta":{"href":"https://example.test/entity/product/product-1","type":"product"}},"quantity":1.25,"price":12345.67,"sourceSlot":{"meta":{"href":"https://example.test/slots/source-1","type":"slot"}},"targetSlot":{"meta":{"href":"https://example.test/slots/target-1","type":"slot"}}}`,
				},
				{"unset", positions.MovementRequest{}, `{}`},
				{
					"minimal new row without slots",
					positions.MovementRequest{ItemRequest: positions.ItemRequest{
						Assortment: model.Value(model.Reference{Meta: model.Meta{
							Href: "https://example.test/entity/product/product-1", Type: "product",
						}}),
						Quantity: model.Value(2.5),
					}},
					`{"assortment":{"meta":{"href":"https://example.test/entity/product/product-1","type":"product"}},"quantity":2.5}`,
				},
				{
					"ID only",
					positions.MovementRequest{ItemRequest: positions.ItemRequest{
						CoreRequest: positions.CoreRequest{ID: model.Value("existing-row")},
					}},
					`{"id":"existing-row"}`,
				},
				{
					"sparse quantity update",
					positions.MovementRequest{ItemRequest: positions.ItemRequest{
						CoreRequest: positions.CoreRequest{ID: model.Value("existing-row")},
						Quantity:    model.Value(0.125),
					}},
					`{"id":"existing-row","quantity":0.125}`,
				},
				{
					"explicit zero",
					positions.MovementRequest{
						ItemRequest: positions.ItemRequest{Quantity: model.Value(0.0)},
						Price:       model.Value(0.0),
					},
					`{"quantity":0,"price":0}`,
				},
				{
					"source slot only",
					positions.MovementRequest{SourceSlot: model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/slots/source-1", Type: "slot"}})},
					`{"sourceSlot":{"meta":{"href":"https://example.test/slots/source-1","type":"slot"}}}`,
				},
				{
					"target slot only",
					positions.MovementRequest{TargetSlot: model.Value(model.Reference{Meta: model.Meta{Href: "https://example.test/slots/target-1", Type: "slot"}})},
					`{"targetSlot":{"meta":{"href":"https://example.test/slots/target-1","type":"slot"}}}`,
				},
				{
					"null fields",
					positions.MovementRequest{
						ItemRequest: positions.ItemRequest{
							CoreRequest: positions.CoreRequest{ID: model.Null[string]()},
							Assortment:  model.Null[model.Reference](), Quantity: model.Null[float64](),
						},
						Price: model.Null[float64](), SourceSlot: model.Null[model.Reference](), TargetSlot: model.Null[model.Reference](),
					},
					`{"id":null,"assortment":null,"quantity":null,"price":null,"sourceSlot":null,"targetSlot":null}`,
				},
			} {
				t.Run(tc.name, func(t *testing.T) {
					encoded, err := codec.marshal(tc.value)
					if err != nil {
						t.Fatal(err)
					}
					jsonAssertEqual(t, encoded, tc.want)
					var got positions.MovementRequest
					if err := codec.unmarshal([]byte(tc.want), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.value) {
						t.Fatalf("movement request = %+v, want %+v", got, tc.value)
					}
				})
			}
			var got positions.MovementRequest
			if err := codec.unmarshal([]byte(`{
				"id":"existing-row","quantity":1.5,"price":2.75,
				"meta":{"href":"https://example.test/positions/existing-row"},"accountId":"account-1",
				"overhead":12.5,"discount":5,"vat":20,"vatEnabled":true,"future":true
			}`), &got); err != nil {
				t.Fatal(err)
			}
			encoded, err := codec.marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			jsonAssertEqual(t, encoded, `{"id":"existing-row","quantity":1.5,"price":2.75}`)
		})
	}
}

func TestEmbeddedMovementPositionsJSON(t *testing.T) {
	type jsonDocumentRequest struct {
		Positions model.Field[[]positions.MovementRequest] `json:"positions,omitzero"`
	}
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value jsonDocumentRequest
				want  string
			}{
				{"unset", jsonDocumentRequest{}, `{}`},
				{"null", jsonDocumentRequest{Positions: model.Null[[]positions.MovementRequest]()}, `{"positions":null}`},
				{"empty slice", jsonDocumentRequest{Positions: model.Value([]positions.MovementRequest{})}, `{"positions":[]}`},
				{"nil slice", jsonDocumentRequest{Positions: model.Value[[]positions.MovementRequest](nil)}, `{"positions":[]}`},
				{
					"replacement array",
					jsonDocumentRequest{Positions: model.Value([]positions.MovementRequest{
						{ItemRequest: positions.ItemRequest{
							CoreRequest: positions.CoreRequest{ID: model.Value("existing-row")},
							Quantity:    model.Value(0.5),
						}},
						{ItemRequest: positions.ItemRequest{
							Assortment: model.Value(model.Reference{Meta: model.Meta{
								Href: "https://example.test/entity/product/product-2", Type: "product",
							}}),
							Quantity: model.Value(1.25),
						}, Price: model.Value(99.75), TargetSlot: model.Value(model.Reference{Meta: model.Meta{
							Href: "https://example.test/slots/target-1", Type: "slot",
						}})},
					})},
					`{"positions":[{"id":"existing-row","quantity":0.5},{"assortment":{"meta":{"href":"https://example.test/entity/product/product-2","type":"product"}},"quantity":1.25,"price":99.75,"targetSlot":{"meta":{"href":"https://example.test/slots/target-1","type":"slot"}}}]}`,
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
