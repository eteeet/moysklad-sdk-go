package query_test

import (
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/query"
)

func TestPageOptionsValues(t *testing.T) {
	for page, want := range map[query.PageOptions]string{
		{}:                       "limit=1000&offset=0",
		{Limit: 1, Offset: 7}:    "limit=1&offset=7",
		{Limit: 1000, Offset: 0}: "limit=1000&offset=0",
		{Fields: "stock"}:        "fields=stock&limit=1000&offset=0",
	} {
		values, err := page.Values()
		if err != nil || values.Encode() != want {
			t.Errorf("%+v: values = %v, err = %v", page, values, err)
		}
	}
	for _, page := range []query.PageOptions{{Limit: -1}, {Limit: 1001}, {Offset: -1}, {Fields: "stock,declaration"}} {
		if _, err := page.Values(); err == nil || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("%+v: err = %v", page, err)
		}
	}
}
