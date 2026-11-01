package positions_test

import (
	"reflect"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

func TestPositionCostAndTransit(t *testing.T) {
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name, input string
				want        positions.Position
			}{
				{"absent", `{}`, positions.Position{}},
				{"integer zero", `{"cost":0,"shipped":0,"inTransit":0}`, positions.Position{}},
				{"decimal zero", `{"cost":0.0,"shipped":0.0,"inTransit":0.0}`, positions.Position{}},
				{"fractions", `{"cost":123.25,"shipped":0.5,"inTransit":1.75}`, positions.Position{Cost: 123.25, Shipped: 0.5, InTransit: 1.75}},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var got positions.Position
					if err := codec.unmarshal([]byte(tc.input), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.want) {
						t.Fatalf("position = %+v, want %+v", got, tc.want)
					}
				})
			}
		})
	}
}

func TestPositionInTransitRequest(t *testing.T) {
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				input positions.PositionRequest
				want  string
			}{
				{"unset", positions.PositionRequest{}, `{}`},
				{"zero", positions.PositionRequest{InTransit: model.Value(0.0)}, `{"inTransit":0}`},
				{"fraction and matching ID", positions.PositionRequest{ItemRequest: positions.ItemRequest{CoreRequest: positions.CoreRequest{ID: model.Value("row")}}, InTransit: model.Value(0.25)}, `{"id":"row","inTransit":0.25}`},
				{"null", positions.PositionRequest{InTransit: model.Null[float64]()}, `{"inTransit":null}`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					data, err := codec.marshal(tc.input)
					if err != nil {
						t.Fatal(err)
					}
					jsonAssertEqual(t, data, tc.want)
					var got positions.PositionRequest
					if err := codec.unmarshal([]byte(tc.want), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.input) {
						t.Fatalf("request = %+v, want %+v", got, tc.input)
					}
				})
			}
			var got positions.PositionRequest
			if err := codec.unmarshal([]byte(`{"shipped":2.5,"overhead":0.0,"inTransit":0,"id":"row"}`), &got); err != nil {
				t.Fatal(err)
			}
			data, err := codec.marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			jsonAssertEqual(t, data, `{"id":"row","inTransit":0}`)
		})
	}
}

func TestPositionCostRequest(t *testing.T) {
	for _, codec := range jsonCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				input positions.PositionRequest
				want  string
			}{
				{"unset", positions.PositionRequest{}, `{}`},
				{"zero", positions.PositionRequest{Cost: model.Value(0.0)}, `{"cost":0}`},
				{"fraction", positions.PositionRequest{Cost: model.Value(123.25)}, `{"cost":123.25}`},
				{"null", positions.PositionRequest{Cost: model.Null[float64]()}, `{"cost":null}`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					data, err := codec.marshal(tc.input)
					if err != nil {
						t.Fatal(err)
					}
					jsonAssertEqual(t, data, tc.want)
					var got positions.PositionRequest
					if err := codec.unmarshal([]byte(tc.want), &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.input) {
						t.Fatalf("request = %+v, want %+v", got, tc.input)
					}
				})
			}
		})
	}
}
