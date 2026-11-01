package documenttest

import (
	"bytes"
	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"testing"
	"time"

	"github.com/eteeet/moysklad-sdk-go/model"
)

type Codec struct {
	Name      string
	Marshal   func(any) ([]byte, error)
	Unmarshal func([]byte, any) error
}

var Codecs = []Codec{
	{"json/v1", jsonv1.Marshal, jsonv1.Unmarshal},
	{"json/v2",
		func(v any) ([]byte, error) { return jsonv2.Marshal(v) },
		func(data []byte, v any) error { return jsonv2.Unmarshal(data, v) },
	},
}

func JSONEqual(t *testing.T, data, want []byte) {
	t.Helper()
	a, b := jsontext.Value(bytes.Clone(data)), jsontext.Value(bytes.Clone(want))
	if err := a.Canonicalize(jsontext.CanonicalizeRawInts(false)); err != nil {
		t.Error(err)
		return
	}
	if err := b.Canonicalize(jsontext.CanonicalizeRawInts(false)); err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(a, b) {
		t.Errorf("JSON = %s, want %s", data, want)
	}
}

func Pointer[T any](v T) *T { return &v }

func DateTime(v string) model.DateTime {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}
	var parsed time.Time
	if len(v) > 10 && v[10] == 'T' {
		parsed, err = time.Parse(time.RFC3339Nano, v)
	} else {
		parsed, err = time.ParseInLocation("2006-01-02 15:04:05", v, location)
	}
	if err != nil {
		panic(err)
	}
	return model.DateTime(parsed.In(location))
}

func NormalizeAttributeJSON(t *testing.T, attributes []model.Attribute) {
	t.Helper()
	for i := range attributes {
		if raw, ok := attributes[i].Value.Get(); ok {
			data := jsontext.Value(bytes.Clone(raw))
			if err := data.Canonicalize(jsontext.CanonicalizeRawInts(false)); err != nil {
				t.Fatal(err)
			}
			attributes[i].Value = model.Value(data)
		}
	}
}
