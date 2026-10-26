package model_test

import (
	"bytes"
	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"testing"
)

var codecs = []struct {
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

func assertJSONEqual(t *testing.T, data []byte, want string) {
	t.Helper()
	gotJSON, wantJSON := jsontext.Value(bytes.Clone(data)), jsontext.Value(want)
	// RFC 8785 без этой опции округляет long за пределами ±2^53.
	if err := gotJSON.Canonicalize(jsontext.CanonicalizeRawInts(false)); err != nil {
		t.Fatal(err)
	}
	if err := wantJSON.Canonicalize(jsontext.CanonicalizeRawInts(false)); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotJSON, wantJSON) {
		t.Fatalf("JSON = %s, want %s", data, want)
	}
}
