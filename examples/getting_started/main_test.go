package main

import (
	"bytes"
	"net/http"
	"testing"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/examples/internal/exampletest"
)

func TestShow(t *testing.T) {
	c := exampletest.Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("authorization = %q", r.Header.Get("Authorization"))
		}
		switch r.URL.Path {
		case "/api/remap/1.2/context/employee":
			w.Write([]byte(`{"id":"e1","name":"Иван"}`))
		case "/api/remap/1.2/context/application":
			w.Write([]byte(`{"id":"a1","appUid":"app.vendor"}`))
		default:
			t.Errorf("path = %s", r.URL.Path)
		}
	})
	for app, want := range map[bool]string{false: "employee id=e1 name=\"Иван\"\n", true: "application id=a1 appUid=app.vendor\n"} {
		var out bytes.Buffer
		if err := show(t.Context(), c, app, &out); err != nil || out.String() != want {
			t.Fatalf("app=%v: out = %q, err = %v", app, out.String(), err)
		}
	}
}

// Ошибка API возвращается из сценария как есть: статус доступен через errors.
func TestShowAPIError(t *testing.T) {
	c := exampletest.Client(t, exampletest.Unauthorized)
	var out bytes.Buffer
	if err := show(t.Context(), c, false, &out); sdkerrors.StatusCode(err) != http.StatusUnauthorized || out.Len() != 0 {
		t.Fatalf("err = %v, out = %q", err, out.String())
	}
}
