package model_test

import (
	jsonv1 "encoding/json"
	jsonv2 "encoding/json/v2"
	"os"
	"testing"
	"time"
	_ "time/tzdata" // Исторические зоны в тестах не зависят от окружения CI.

	"github.com/eteeet/moysklad-sdk-go/model"
)

var dateTimeCodecs = []struct {
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

func TestDateTimeUnmarshalJSON(t *testing.T) {
	for _, codec := range dateTimeCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name   string
				input  string
				want   time.Time
				offset int
			}{
				{"API seconds", `"2024-01-15 10:30:00"`, time.Date(2024, 1, 15, 7, 30, 0, 0, time.UTC), 3 * 3600},
				{"API milliseconds", `"2024-01-15 10:30:00.123"`, time.Date(2024, 1, 15, 7, 30, 0, 123000000, time.UTC), 3 * 3600},
				{"API short fraction", `"2024-01-15 10:30:00.1"`, time.Date(2024, 1, 15, 7, 30, 0, 100000000, time.UTC), 3 * 3600},
				{"API nanoseconds", `"2024-01-15 10:30:00.123456789"`, time.Date(2024, 1, 15, 7, 30, 0, 123456789, time.UTC), 3 * 3600},
				{"Moscow 2012", `"2012-06-01 12:00:00"`, time.Date(2012, 6, 1, 8, 0, 0, 0, time.UTC), 4 * 3600},
				{"Moscow summer 2008", `"2008-06-01 12:00:00"`, time.Date(2008, 6, 1, 8, 0, 0, 0, time.UTC), 4 * 3600},
				{"Moscow winter 2008", `"2008-01-01 12:00:00"`, time.Date(2008, 1, 1, 9, 0, 0, 0, time.UTC), 3 * 3600},
				{"RFC3339 offset", `"2024-01-15T10:30:00.000+03:00"`, time.Date(2024, 1, 15, 7, 30, 0, 0, time.UTC), 3 * 3600},
				{"RFC3339 UTC", `"2024-01-15T07:30:00Z"`, time.Date(2024, 1, 15, 7, 30, 0, 0, time.UTC), 3 * 3600},
				{"RFC3339 nanoseconds", `"2024-01-15T07:30:00.123456789Z"`, time.Date(2024, 1, 15, 7, 30, 0, 123456789, time.UTC), 3 * 3600},
				{"RFC3339 historical offset", `"2012-06-01T12:00:00+03:00"`, time.Date(2012, 6, 1, 9, 0, 0, 0, time.UTC), 4 * 3600},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var got model.DateTime
					if err := codec.unmarshal([]byte(tc.input), &got); err != nil {
						t.Fatal(err)
					}
					actual := time.Time(got)
					if !actual.Equal(tc.want) {
						t.Fatalf("time = %s, want %s", actual, tc.want)
					}
					if _, offset := actual.Zone(); offset != tc.offset {
						t.Fatalf("offset = %d, want %d", offset, tc.offset)
					}
					if got := actual.Location().String(); got != "Europe/Moscow" {
						t.Fatalf("location = %s, want Europe/Moscow", got)
					}
				})
			}
		})
	}
}

func TestDateTimeMarshalJSON(t *testing.T) {
	moscow, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range dateTimeCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value time.Time
				want  string
			}{
				{"UTC", time.Date(2024, 1, 15, 7, 30, 0, 0, time.UTC), `"2024-01-15 10:30:00.000"`},
				{"milliseconds", time.Date(2024, 1, 15, 7, 30, 0, 123000000, time.UTC), `"2024-01-15 10:30:00.123"`},
				{"truncate nanoseconds", time.Date(2024, 1, 15, 7, 30, 0, 123999999, time.UTC), `"2024-01-15 10:30:00.123"`},
				{"historical offset", time.Date(2012, 6, 1, 8, 0, 0, 0, time.UTC), `"2012-06-01 12:00:00.000"`},
				{"summer 2008", time.Date(2008, 6, 1, 8, 0, 0, 0, time.UTC), `"2008-06-01 12:00:00.000"`},
				{"zero", time.Time{}, "null"},
				{"zero in Moscow", time.Time{}.In(moscow), "null"},
			} {
				t.Run(tc.name, func(t *testing.T) {
					got, err := codec.marshal(model.DateTime(tc.value))
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

func TestDateTimeEmptyJSON(t *testing.T) {
	for _, codec := range dateTimeCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, input := range []string{"null", " \nnull\t", `""`} {
				t.Run(input, func(t *testing.T) {
					got := model.DateTime(time.Date(2024, 1, 15, 7, 30, 0, 0, time.UTC))
					if err := codec.unmarshal([]byte(input), &got); err != nil {
						t.Fatal(err)
					}
					if actual := time.Time(got); !actual.IsZero() {
						t.Fatalf("time = %s, want zero", actual)
					}
					encoded, err := codec.marshal(got)
					if err != nil {
						t.Fatal(err)
					}
					if string(encoded) != "null" {
						t.Fatalf("JSON = %s, want null", encoded)
					}
				})
			}
		})
	}
}

func TestDateTimeOmitZero(t *testing.T) {
	moscow, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range dateTimeCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				name   string
				value  time.Time
				isZero bool
				want   string
			}{
				{"zero", time.Time{}, true, `{}`},
				{"zero in Moscow", time.Time{}.In(moscow), true, `{}`},
				{"nonzero", time.Date(2024, 1, 15, 7, 30, 0, 0, time.UTC), false, `{"date":"2024-01-15 10:30:00.000"}`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					value := model.DateTime(tc.value)
					if got := value.IsZero(); got != tc.isZero {
						t.Fatalf("IsZero = %v, want %v", got, tc.isZero)
					}
					got, err := codec.marshal(struct {
						Date model.DateTime `json:"date,omitzero"`
					}{Date: value})
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

func TestDateTimeUnmarshalErrors(t *testing.T) {
	for _, codec := range dateTimeCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, input := range []string{
				`"invalid"`, `" "`, `"2024-01-15"`, `"2024-01-15 10:30"`,
				`"2024-02-30 10:30:00"`, `"2024-01-15 25:00:00"`,
				`"2024-01-15T10:30:00"`, `"2024-01-15T10:30:00+25:00"`,
				`42`, `true`, `{}`, `[]`, `"unterminated`,
			} {
				t.Run(input, func(t *testing.T) {
					initial := model.DateTime(time.Date(2024, 1, 15, 7, 30, 0, 0, time.UTC))
					got := initial
					if err := codec.unmarshal([]byte(input), &got); err == nil {
						t.Fatal("expected datetime error")
					}
					if got != initial {
						t.Fatalf("time changed after error: %s", time.Time(got))
					}
				})
			}
		})
	}
}

func TestDateTimeMarshalInvalidYear(t *testing.T) {
	for _, codec := range dateTimeCodecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, value := range []time.Time{
				time.Date(-1, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(9999, 12, 31, 23, 0, 0, 0, time.UTC),
			} {
				if _, err := codec.marshal(model.DateTime(value)); err == nil {
					t.Errorf("expected year error for %s", value)
				}
			}
		})
	}
}

func TestDateTimeProductFixture(t *testing.T) {
	data, err := os.ReadFile("../testdata/upstream/product.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range dateTimeCodecs {
		t.Run(codec.name, func(t *testing.T) {
			var got struct {
				Updated model.DateTime `json:"updated"`
			}
			if err := codec.unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			want := time.Date(2024, 1, 15, 7, 30, 0, 0, time.UTC)
			if actual := time.Time(got.Updated); !actual.Equal(want) {
				t.Fatalf("updated = %s, want %s", actual, want)
			}
		})
	}
}
