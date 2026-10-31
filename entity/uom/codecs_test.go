package uom_test

import (
	jsonv1 "encoding/json"
	jsonv2 "encoding/json/v2"
)

var codecs = []struct {
	name      string
	unmarshal func([]byte, any) error
}{
	{"json/v1", jsonv1.Unmarshal},
	{"json/v2", func(data []byte, v any) error { return jsonv2.Unmarshal(data, v) }},
}
