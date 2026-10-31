package counterparty_test

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
	{"json/v2", func(v any) ([]byte, error) { return jsonv2.Marshal(v) }, func(data []byte, v any) error { return jsonv2.Unmarshal(data, v) }},
}

func assertJSONEqual(t *testing.T, data, want []byte) {
	t.Helper()
	a, b := jsontext.Value(bytes.Clone(data)), jsontext.Value(bytes.Clone(want))
	if err := a.Canonicalize(jsontext.CanonicalizeRawInts(false)); err != nil {
		t.Fatal(err)
	}
	if err := b.Canonicalize(jsontext.CanonicalizeRawInts(false)); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Errorf("JSON = %s, want %s", data, want)
	}
}
