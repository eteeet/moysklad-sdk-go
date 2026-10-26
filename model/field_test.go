package model_test

import (
	jsonv1 "encoding/json"
	jsonv2 "encoding/json/v2"
	"reflect"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/model"
)

var fieldCodecs = []struct {
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

type fieldEnvelope[T any] struct {
	Value model.Field[T] `json:"value,omitzero"`
}

func TestFieldStates(t *testing.T) {
	for _, tc := range []struct {
		name     string
		field    model.Field[int]
		isZero   bool
		isNull   bool
		value    int
		hasValue bool
	}{
		{"unset", model.Field[int]{}, true, false, 0, false},
		{"null", model.Null[int](), false, true, 0, false},
		{"zero value", model.Value(0), false, false, 0, true},
		{"value", model.Value(7), false, false, 7, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.field.IsZero(); got != tc.isZero {
				t.Fatalf("IsZero = %v, want %v", got, tc.isZero)
			}
			if got := tc.field.IsNull(); got != tc.isNull {
				t.Fatalf("IsNull = %v, want %v", got, tc.isNull)
			}
			if value, ok := tc.field.Get(); value != tc.value || ok != tc.hasValue {
				t.Fatalf("Get = (%v, %v), want (%v, %v)", value, ok, tc.value, tc.hasValue)
			}
		})
	}
}

func TestFieldMarshalJSON(t *testing.T) {
	type nested struct {
		Name string `json:"name"`
	}
	for _, codec := range fieldCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value any
				want  string
			}{
				{"unset", fieldEnvelope[int]{}, `{}`},
				{"null", fieldEnvelope[int]{model.Null[int]()}, `{"value":null}`},
				{"false", fieldEnvelope[bool]{model.Value(false)}, `{"value":false}`},
				{"zero", fieldEnvelope[int]{model.Value(0)}, `{"value":0}`},
				{"empty string", fieldEnvelope[string]{model.Value("")}, `{"value":""}`},
				{"string", fieldEnvelope[string]{model.Value("product")}, `{"value":"product"}`},
				{"empty slice", fieldEnvelope[[]int]{model.Value([]int{})}, `{"value":[]}`},
				{"nil slice", fieldEnvelope[[]int]{model.Value[[]int](nil)}, `{"value":[]}`},
				{"slice", fieldEnvelope[[]int]{model.Value([]int{1, 2})}, `{"value":[1,2]}`},
				{"nil map", fieldEnvelope[map[string]int]{model.Value[map[string]int](nil)}, `{"value":{}}`},
				{"nested", fieldEnvelope[nested]{model.Value(nested{Name: "product"})}, `{"value":{"name":"product"}}`},
				{"nil pointer", fieldEnvelope[*int]{model.Value[*int](nil)}, `{"value":null}`},
				{"nil interface", fieldEnvelope[any]{model.Value[any](nil)}, `{"value":null}`},
				{"bytes", fieldEnvelope[[]byte]{model.Value([]byte{1, 2})}, `{"value":"AQI="}`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					got, err := codec.marshal(tc.value)
					if err != nil {
						t.Fatal(err)
					}
					if string(got) != tc.want {
						t.Fatalf("JSON = %s, want %s", got, tc.want)
					}
				})
			}
		})
	}
}

func TestFieldMarshalUnsetError(t *testing.T) {
	for _, codec := range fieldCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, value := range []any{
				model.Field[int]{},
				struct {
					Value model.Field[int] `json:"value"`
				}{},
			} {
				if _, err := codec.marshal(value); err == nil {
					t.Fatal("expected error for unset Field without omitzero")
				}
			}
		})
	}
}

func TestFieldUnmarshalJSON(t *testing.T) {
	for _, codec := range fieldCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name     string
				initial  model.Field[int]
				input    string
				isZero   bool
				isNull   bool
				value    int
				hasValue bool
			}{
				{"absent", model.Field[int]{}, `{}`, true, false, 0, false},
				{"null", model.Field[int]{}, `{"value":null}`, false, true, 0, false},
				{"zero", model.Field[int]{}, `{"value":0}`, false, false, 0, true},
				{"value", model.Field[int]{}, `{"value":7}`, false, false, 7, true},
				{"value to null", model.Value(7), "{\"value\": \nnull\t}", false, true, 0, false},
				{"null to value", model.Null[int](), `{"value":7}`, false, false, 7, true},
				{"value to zero", model.Value(7), `{"value":0}`, false, false, 0, true},
				{"absent preserves value", model.Value(7), `{}`, false, false, 7, true},
				{"absent preserves null", model.Null[int](), `{}`, false, true, 0, false},
			} {
				t.Run(tc.name, func(t *testing.T) {
					got := fieldEnvelope[int]{Value: tc.initial}
					if err := codec.unmarshal([]byte(tc.input), &got); err != nil {
						t.Fatal(err)
					}
					if got.Value.IsZero() != tc.isZero || got.Value.IsNull() != tc.isNull {
						t.Fatalf("state: IsZero = %v, IsNull = %v; want %v, %v", got.Value.IsZero(), got.Value.IsNull(), tc.isZero, tc.isNull)
					}
					if value, ok := got.Value.Get(); value != tc.value || ok != tc.hasValue {
						t.Fatalf("Get = (%v, %v), want (%v, %v)", value, ok, tc.value, tc.hasValue)
					}
				})
			}
		})
	}
}

func TestFieldUnmarshalValues(t *testing.T) {
	t.Run("false", func(t *testing.T) {
		checkFieldUnmarshalValue(t, `false`, false)
	})
	t.Run("zero", func(t *testing.T) {
		checkFieldUnmarshalValue(t, `0`, 0)
	})
	t.Run("empty string", func(t *testing.T) {
		checkFieldUnmarshalValue(t, `""`, "")
	})
	t.Run("empty slice", func(t *testing.T) {
		checkFieldUnmarshalValue(t, `[]`, []int{})
	})
	t.Run("slice", func(t *testing.T) {
		checkFieldUnmarshalValue(t, `[1,2]`, []int{1, 2})
	})
	t.Run("empty map", func(t *testing.T) {
		checkFieldUnmarshalValue(t, `{}`, map[string]int{})
	})
	t.Run("bytes", func(t *testing.T) {
		checkFieldUnmarshalValue(t, `"AQI="`, []byte{1, 2})
	})
	t.Run("nested", func(t *testing.T) {
		type nested struct {
			Name  string `json:"name"`
			Count int    `json:"count"`
		}
		checkFieldUnmarshalValue(t, `{"name":"product","count":2}`, nested{Name: "product", Count: 2})
	})
}

func checkFieldUnmarshalValue[T any](t *testing.T, input string, want T) {
	t.Helper()
	for _, codec := range fieldCodecs {
		t.Run(codec.name, func(t *testing.T) {
			var got model.Field[T]
			if err := codec.unmarshal([]byte(input), &got); err != nil {
				t.Fatal(err)
			}
			value, ok := got.Get()
			if got.IsZero() || got.IsNull() || !ok {
				t.Fatal("decoded field must be a value")
			}
			if !reflect.DeepEqual(value, want) {
				t.Fatalf("Get = %#v, want %#v", value, want)
			}
			encoded, err := codec.marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			if string(encoded) != input {
				t.Fatalf("JSON = %s, want %s", encoded, input)
			}
		})
	}
}

func TestFieldUnmarshalErrors(t *testing.T) {
	for _, codec := range fieldCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, initial := range []model.Field[int]{
				{}, model.Null[int](), model.Value(7),
			} {
				for _, input := range []string{`"wrong"`, `true`, `{}`, `[]`, `1.5`, ``, `{"broken":`} {
					got := initial
					if err := codec.unmarshal([]byte(input), &got); err == nil {
						t.Fatalf("expected error for %q", input)
					}
					value, ok := got.Get()
					want, wantOK := initial.Get()
					if got.IsZero() != initial.IsZero() || got.IsNull() != initial.IsNull() || value != want || ok != wantOK {
						t.Fatalf("field changed after error for %q", input)
					}
				}
			}
		})
	}
}

func TestFieldUnmarshalInvalidJSON(t *testing.T) {
	for _, codec := range fieldCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, input := range []string{`nullx`, `7 8`, `"unterminated`, `{"broken":`} {
				got := model.Value(7)
				if err := codec.unmarshal([]byte(input), &got); err == nil {
					t.Fatalf("expected invalid JSON error for %q", input)
				}
			}
		})
	}
}

func TestFieldUnmarshalNestedError(t *testing.T) {
	type nested struct {
		Labels map[string]string `json:"labels"`
		Count  int               `json:"count"`
	}
	for _, codec := range fieldCodecs {
		t.Run(codec.name, func(t *testing.T) {
			got := model.Value(nested{Labels: map[string]string{"name": "before"}, Count: 7})
			input := `{"labels":{"name":"after"},"count":"wrong"}`
			if err := codec.unmarshal([]byte(input), &got); err == nil {
				t.Fatal("expected nested value error")
			}
			value, ok := got.Get()
			if !ok || value.Count != 7 || len(value.Labels) != 1 || value.Labels["name"] != "before" {
				t.Fatalf("field changed after error: %+v", value)
			}
		})
	}
}

func TestFieldDateTimeRoundTrip(t *testing.T) {
	for _, codec := range fieldCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name     string
				value    model.Field[model.DateTime]
				wantJSON string
				isZero   bool
				isNull   bool
			}{
				{"unset", model.Field[model.DateTime]{}, `{}`, true, false},
				{"null", model.Null[model.DateTime](), `{"value":null}`, false, true},
				{"zero datetime becomes null", model.Value(model.DateTime{}), `{"value":null}`, false, true},
				{"datetime", model.Value(model.DateTime(time.Date(2024, 1, 15, 7, 30, 0, 123000000, time.UTC))), `{"value":"2024-01-15 10:30:00.123"}`, false, false},
			} {
				t.Run(tc.name, func(t *testing.T) {
					encoded, err := codec.marshal(fieldEnvelope[model.DateTime]{Value: tc.value})
					if err != nil {
						t.Fatal(err)
					}
					if string(encoded) != tc.wantJSON {
						t.Fatalf("JSON = %s, want %s", encoded, tc.wantJSON)
					}
					var got fieldEnvelope[model.DateTime]
					if err := codec.unmarshal(encoded, &got); err != nil {
						t.Fatal(err)
					}
					if got.Value.IsZero() != tc.isZero || got.Value.IsNull() != tc.isNull {
						t.Fatalf("state after roundtrip: IsZero = %v, IsNull = %v", got.Value.IsZero(), got.Value.IsNull())
					}
					value, ok := got.Value.Get()
					if ok != (!tc.isZero && !tc.isNull) {
						t.Fatalf("Get ok = %v", ok)
					}
					if ok && !time.Time(value).Equal(time.Date(2024, 1, 15, 7, 30, 0, 123000000, time.UTC)) {
						t.Fatalf("datetime = %s", time.Time(value))
					}
				})
			}
		})
	}
}

func TestFieldMarshalUnsupportedValue(t *testing.T) {
	for _, codec := range fieldCodecs {
		t.Run(codec.name, func(t *testing.T) {
			if _, err := codec.marshal(fieldEnvelope[chan int]{model.Value[chan int](nil)}); err == nil {
				t.Fatal("expected error for unsupported value")
			}
		})
	}
}
