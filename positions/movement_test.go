package positions_test

import (
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

func TestMovementJSON(t *testing.T) {
	const input = `{
		"meta":{"href":"https://example.test/positions/row-1","type":"moveposition","mediaType":"application/json"},
		"id":"row-1","accountId":"account-1",
		"assortment":{"meta":{"href":"https://example.test/entity/product/product-1","type":"product"}},
		"quantity":1.25,"price":12345.67,"overhead":12.5,
		"sourceSlot":{"meta":{"href":"https://example.test/slots/source-1","type":"slot"}},
		"targetSlot":{"meta":{"href":"https://example.test/slots/target-1","type":"slot"}},
		"future":{"enabled":true}
	}`
	want := positions.Movement{
		Item: positions.Item{
			Core: positions.Core{
				Meta: model.Meta{Href: "https://example.test/positions/row-1", Type: "moveposition", MediaType: "application/json"},
				ID:   "row-1", AccountID: "account-1",
			},
			Assortment: model.Reference{Meta: model.Meta{Href: "https://example.test/entity/product/product-1", Type: "product"}},
			Quantity:   1.25,
		},
		Price: 12345.67, Overhead: 12.5,
		SourceSlot: &model.Reference{Meta: model.Meta{Href: "https://example.test/slots/source-1", Type: "slot"}},
		TargetSlot: &model.Reference{Meta: model.Meta{Href: "https://example.test/slots/target-1", Type: "slot"}},
	}
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			var got positions.Movement
			if err := codec.unmarshal([]byte(input), &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("movement = %+v, want %+v", got, want)
			}
			encoded, err := codec.marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			jsonAssertEqual(t, encoded, `{
				"meta":{"href":"https://example.test/positions/row-1","type":"moveposition","mediaType":"application/json"},
				"id":"row-1","accountId":"account-1",
				"assortment":{"meta":{"href":"https://example.test/entity/product/product-1","type":"product"}},
				"quantity":1.25,"price":12345.67,"overhead":12.5,
				"sourceSlot":{"meta":{"href":"https://example.test/slots/source-1","type":"slot"}},
				"targetSlot":{"meta":{"href":"https://example.test/slots/target-1","type":"slot"}}
			}`)
		})
	}
}

func TestMovementOptionalFieldsJSON(t *testing.T) {
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				input string
				want  positions.Movement
			}{
				{"absent fields", `{}`, positions.Movement{}},
				{"integer overhead zero", `{"overhead":0}`, positions.Movement{}},
				{"decimal overhead zero", `{"overhead":0.0}`, positions.Movement{}},
				{"fractional overhead", `{"overhead":0.125}`, positions.Movement{Overhead: 0.125}},
				{"null slots", `{"sourceSlot":null,"targetSlot":null}`, positions.Movement{}},
				{
					"source slot only",
					`{"sourceSlot":{"meta":{"href":"https://example.test/slots/source-1","type":"slot"}},"future":true}`,
					positions.Movement{SourceSlot: &model.Reference{Meta: model.Meta{Href: "https://example.test/slots/source-1", Type: "slot"}}},
				},
				{
					"target slot only",
					`{"targetSlot":{"meta":{"href":"https://example.test/slots/target-1","type":"slot"}}}`,
					positions.Movement{TargetSlot: &model.Reference{Meta: model.Meta{Href: "https://example.test/slots/target-1", Type: "slot"}}},
				},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var got positions.Movement
					if err := codec.unmarshal([]byte(tc.input), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.want) {
						t.Fatalf("movement = %+v, want %+v", got, tc.want)
					}
				})
			}
		})
	}
}
