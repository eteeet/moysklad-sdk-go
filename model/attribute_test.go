package model_test

import (
	"encoding/json/jsontext"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/model"
)

func TestAttributeAsLong(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				input string
				want  int64
			}{
				{`{"meta":{},"type":"long","value":0}`, 0},
				{`{"meta":{},"type":"long","value":9007199254740993}`, 9007199254740993},
				{`{"meta":{},"type":"long","value":9223372036854775807}`, 9223372036854775807},
				{`{"meta":{},"type":"long","value":-9223372036854775808}`, -9223372036854775808},
			} {
				var got model.Attribute
				if err := codec.unmarshal([]byte(tc.input), &got); err != nil {
					t.Fatal(err)
				}
				value, err := got.AsLong()
				if err != nil || value != tc.want {
					t.Fatalf("AsLong = %d, %v; want %d", value, err, tc.want)
				}
				encoded, err := codec.marshal(got)
				if err != nil {
					t.Fatal(err)
				}
				assertJSONEqual(t, encoded, tc.input)
			}
		})
	}
}

func TestAttributeGetters(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				input string
				read  func(model.Attribute) (any, error)
				want  any
			}{
				{`{"type":"string","value":"text"}`, func(a model.Attribute) (any, error) { return a.AsString() }, "text"},
				{`{"type":"text","value":""}`, func(a model.Attribute) (any, error) { return a.AsString() }, ""},
				{`{"type":"link","value":"https://example.com"}`, func(a model.Attribute) (any, error) { return a.AsString() }, "https://example.com"},
				{`{"type":"file","value":"file.txt"}`, func(a model.Attribute) (any, error) { return a.AsString() }, "file.txt"},
				{`{"type":"boolean","value":false}`, func(a model.Attribute) (any, error) { return a.AsBool() }, false},
				{`{"type":"double","value":0.75}`, func(a model.Attribute) (any, error) { return a.AsDouble() }, 0.75},
				{`{"type":"futureentity","value":{"meta":{"href":"future","type":"futureentity"},"name":"new"}}`, func(a model.Attribute) (any, error) { return a.AsReference() }, model.AttributeReference{Meta: model.Meta{Href: "future", Type: "futureentity"}, Name: "new"}},
			} {
				var a model.Attribute
				if err := codec.unmarshal([]byte(tc.input), &a); err != nil {
					t.Fatal(err)
				}
				got, err := tc.read(a)
				if err != nil || !reflect.DeepEqual(got, tc.want) {
					t.Fatalf("input %s: value = %#v, %v; want %#v", tc.input, got, err, tc.want)
				}
			}
			var a model.Attribute
			if err := codec.unmarshal([]byte(`{"type":"time","value":"2024-01-15 10:30:00"}`), &a); err != nil {
				t.Fatal(err)
			}
			got, err := a.AsTime()
			if err != nil || !time.Time(got).Equal(time.Date(2024, 1, 15, 7, 30, 0, 0, time.UTC)) {
				t.Fatalf("AsTime = %s, %v", time.Time(got), err)
			}
		})
	}
}

func TestAttributeGetterErrors(t *testing.T) {
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				input string
				read  func(model.Attribute) error
			}{
				{`{"type":"boolean","value":false}`, func(a model.Attribute) error { _, e := a.AsString(); return e }},
				{`{"type":"string","value":"true"}`, func(a model.Attribute) error { _, e := a.AsBool(); return e }},
				{`{"type":"long","value":1}`, func(a model.Attribute) error { _, e := a.AsDouble(); return e }},
				{`{"type":"double","value":1}`, func(a model.Attribute) error { _, e := a.AsLong(); return e }},
				{`{"type":"long","value":1.5}`, func(a model.Attribute) error { _, e := a.AsLong(); return e }},
				{`{"type":"long","value":9223372036854775808}`, func(a model.Attribute) error { _, e := a.AsLong(); return e }},
				{`{"type":"boolean","value":null}`, func(a model.Attribute) error { _, e := a.AsBool(); return e }},
				{`{"type":"boolean"}`, func(a model.Attribute) error { _, e := a.AsBool(); return e }},
				{`{"type":"time","value":"time"}`, func(a model.Attribute) error { _, e := a.AsTime(); return e }},
				{`{"type":"string","value":"2024-01-15 10:30:00"}`, func(a model.Attribute) error { _, e := a.AsTime(); return e }},
				{`{"type":"double","value":{"meta":{"href":"x"}}}`, func(a model.Attribute) error { _, e := a.AsReference(); return e }},
				{`{"type":"future","value":{"name":"no meta"}}`, func(a model.Attribute) error { _, e := a.AsReference(); return e }},
				{`{"type":"future","value":[1,2]}`, func(a model.Attribute) error { _, e := a.AsReference(); return e }},
				{`{"type":"future","value":"raw"}`, func(a model.Attribute) error { _, e := a.AsString(); return e }},
			} {
				var a model.Attribute
				if err := codec.unmarshal([]byte(tc.input), &a); err != nil {
					t.Fatal(err)
				}
				if err := tc.read(a); err == nil {
					t.Fatalf("getter accepted %s", tc.input)
				}
			}
		})
	}
}

func TestAttributeUnknownTypeAndFixture(t *testing.T) {
	fixture, err := os.ReadFile("../testdata/upstream/entity_with_attributes.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			var got struct {
				Attributes []model.Attribute `json:"attributes"`
			}
			if err := codec.unmarshal(fixture, &got); err != nil {
				t.Fatal(err)
			}
			if len(got.Attributes) != 18 {
				t.Fatalf("attributes = %d, want 18", len(got.Attributes))
			}
			if value, err := got.Attributes[0].AsString(); err != nil || value != "Значение" {
				t.Fatalf("AsString = %q, %v", value, err)
			}
			if got.Attributes[1].Download == nil || got.Attributes[1].Download.MediaType != "application/octet-stream" {
				t.Fatal("file download metadata lost")
			}
			if _, err := got.Attributes[4].AsTime(); err == nil {
				t.Fatal("invalid fixture date must fail only at getter")
			}
			if value, err := got.Attributes[6].AsBool(); err != nil || value {
				t.Fatalf("AsBool = %v, %v", value, err)
			}
			if value, err := got.Attributes[11].AsReference(); err != nil || value.Name != "customentity" {
				t.Fatalf("AsReference = %+v, %v", value, err)
			}
			const input = `{"meta":{},"type":"future","value":[9007199254740993,{"new":true}]}`
			var unknown model.Attribute
			if err := codec.unmarshal([]byte(input), &unknown); err != nil {
				t.Fatal(err)
			}
			raw, ok := unknown.Value.Get()
			if !ok || raw.Kind() != '[' {
				t.Fatalf("raw = %s", raw)
			}
			encoded, err := codec.marshal(unknown)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, encoded, input)
		})
	}
}

func TestAttributeConstructors(t *testing.T) {
	meta := model.Meta{Href: "attribute"}
	value, err := model.AttrValue(meta, int64(9007199254740993))
	if err != nil {
		t.Fatal(err)
	}
	date, err := model.AttrValue(meta, model.DateTime(time.Date(2024, 1, 15, 7, 30, 45, 123000000, time.UTC)))
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range codecs {
		t.Run(codec.name, func(t *testing.T) {
			for _, tc := range []struct {
				value model.Attribute
				want  string
			}{
				{value, `{"meta":{"href":"attribute"},"value":9007199254740993}`},
				{model.AttrNull(meta), `{"meta":{"href":"attribute"},"value":null}`},
				{model.AttrFile(meta, "file.txt", []byte{1, 2}), `{"meta":{"href":"attribute"},"file":{"filename":"file.txt","content":"AQI="}}`},
				{model.AttrFileNull(meta), `{"meta":{"href":"attribute"},"file":null}`},
				{date, `{"meta":{"href":"attribute"},"value":"2024-01-15 10:30:45.123"}`},
			} {
				encoded, err := codec.marshal(tc.value)
				if err != nil {
					t.Fatal(err)
				}
				assertJSONEqual(t, encoded, tc.want)
			}
		})
	}
	for _, tc := range []struct {
		value any
		want  string
	}{
		{false, `{"meta":{"href":"attribute"},"value":false}`},
		{0, `{"meta":{"href":"attribute"},"value":0}`},
		{"", `{"meta":{"href":"attribute"},"value":""}`},
		{[]int(nil), `{"meta":{"href":"attribute"},"value":[]}`},
	} {
		attribute, err := model.AttrValue(meta, tc.value)
		if err != nil {
			t.Fatal(err)
		}
		for _, codec := range codecs {
			data, err := codec.marshal(attribute)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, data, tc.want)
		}
	}
	for _, invalid := range []any{make(chan int), jsontext.Value(`{"invalid":`)} {
		if _, err := model.AttrValue(meta, invalid); err == nil {
			t.Fatal("expected encoding error")
		}
	}
}
