package main

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/examples/internal/exampletest"
)

// Период передаётся явно, московским временем.
func TestShow(t *testing.T) {
	c := exampletest.Client(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if r.URL.Path != "/api/remap/1.2/report/turnover/all" || q.Get("limit") != "100" ||
			q.Get("momentFrom") != "2026-10-01 00:00:00" || q.Get("momentTo") != "2026-10-31 23:59:59" {
			t.Errorf("request = %s", r.URL.RequestURI())
		}
		w.Write([]byte(`{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/report/turnover/all"},"rows":[` +
			`{"assortment":{"name":"Чай"},"income":{"sum":10000,"quantity":2},"outcome":{"sum":5000,"quantity":1}}]}`))
	})
	moscow := time.FixedZone("MSK", 3*60*60)
	var out bytes.Buffer
	err := show(t.Context(), c, time.Date(2026, 10, 1, 0, 0, 0, 0, moscow), time.Date(2026, 10, 31, 23, 59, 59, 0, moscow), &out)
	if err != nil || out.String() != "Чай income=2/10000 outcome=1/5000\n" {
		t.Fatalf("out = %q, err = %v", out.String(), err)
	}
}

// Ошибка API возвращается из сценария как есть: статус доступен через errors.
func TestShowAPIError(t *testing.T) {
	c := exampletest.Client(t, exampletest.Unauthorized)
	var out bytes.Buffer
	if err := show(t.Context(), c, time.Now().AddDate(0, 0, -1), time.Now(), &out); sdkerrors.StatusCode(err) != http.StatusUnauthorized || out.Len() != 0 {
		t.Fatalf("err = %v, out = %q", err, out.String())
	}
}
