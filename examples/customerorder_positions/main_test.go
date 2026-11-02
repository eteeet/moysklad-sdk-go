package main

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/examples/internal/exampletest"
)

const orderID = "7944ef04-f831-11e5-7a69-971500188b19"

// Get заказа, затем две страницы позиций: вторая — по meta.nextHref.
func TestShow(t *testing.T) {
	c := exampletest.Client(t, func(w http.ResponseWriter, r *http.Request) {
		base := "http://" + r.Host + "/api/remap/1.2/entity/customerorder/" + orderID + "/positions?limit=100&offset="
		switch r.URL.RequestURI() {
		case "/api/remap/1.2/entity/customerorder/" + orderID:
			fmt.Fprintf(w, `{"id":%q,"name":"00001","sum":150000}`, orderID)
		case "/api/remap/1.2/entity/customerorder/" + orderID + "/positions?limit=100&offset=0":
			fmt.Fprintf(w, `{"meta":{"href":%q,"nextHref":%q},"rows":[{"id":"p1","quantity":2,"price":50000}]}`, base+"0", base+"100")
		case "/api/remap/1.2/entity/customerorder/" + orderID + "/positions?limit=100&offset=100":
			fmt.Fprintf(w, `{"meta":{"href":%q},"rows":[{"id":"p2","quantity":1,"price":50000}]}`, base+"100")
		default:
			t.Errorf("request = %s", r.URL.RequestURI())
			w.WriteHeader(http.StatusNotFound)
		}
	})
	var out bytes.Buffer
	want := "order 00001 sum=150000\nposition p1 quantity=2 price=50000\nposition p2 quantity=1 price=50000\n"
	if err := show(t.Context(), c, orderID, &out); err != nil || out.String() != want {
		t.Fatalf("out = %q, err = %v", out.String(), err)
	}
}

// Ошибка API возвращается из сценария как есть: статус доступен через errors.
func TestShowAPIError(t *testing.T) {
	c := exampletest.Client(t, exampletest.Unauthorized)
	var out bytes.Buffer
	if err := show(t.Context(), c, orderID, &out); sdkerrors.StatusCode(err) != http.StatusUnauthorized || out.Len() != 0 {
		t.Fatalf("err = %v, out = %q", err, out.String())
	}
}
