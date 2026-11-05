package main

import (
	"bytes"
	"net/http"
	"sync"
	"testing"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/examples/internal/exampletest"
)

// Страница контекстов и отдельный GET событий для каждого контекста.
func TestShow(t *testing.T) {
	var mu sync.Mutex
	var paths []string
	c := exampletest.Client(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.RequestURI())
		mu.Unlock()
		switch r.URL.RequestURI() {
		case "/api/remap/1.2/audit?limit=10":
			w.Write([]byte(`{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/audit"},"rows":[` +
				`{"id":"11111111-1111-1111-1111-111111111111","eventType":"update","entityType":"product"},{"id":"22222222-2222-2222-2222-222222222222","eventType":"create","entityType":"demand"}]}`))
		case "/api/remap/1.2/audit/11111111-1111-1111-1111-111111111111/events?limit=100":
			w.Write([]byte(`{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/audit/11111111-1111-1111-1111-111111111111/events"},"rows":[` +
				`{"eventType":"update","entityType":"product","name":"Чай"}]}`))
		case "/api/remap/1.2/audit/22222222-2222-2222-2222-222222222222/events?limit=100":
			w.Write([]byte(`{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/audit/22222222-2222-2222-2222-222222222222/events"},"rows":[]}`))
		default:
			t.Errorf("request = %s", r.URL.RequestURI())
			w.WriteHeader(http.StatusNotFound)
		}
	})
	var out bytes.Buffer
	want := "context 11111111-1111-1111-1111-111111111111 update product\n  event update product \"Чай\"\ncontext 22222222-2222-2222-2222-222222222222 create demand\n"
	err := show(t.Context(), c, &out)
	mu.Lock()
	defer mu.Unlock()
	if err != nil || out.String() != want || len(paths) != 3 {
		t.Fatalf("out = %q, err = %v, paths = %v", out.String(), err, paths)
	}
}

// Ошибка API возвращается из сценария как есть: статус доступен через errors.
func TestShowAPIError(t *testing.T) {
	c := exampletest.Client(t, exampletest.Unauthorized)
	var out bytes.Buffer
	if err := show(t.Context(), c, &out); sdkerrors.StatusCode(err) != http.StatusUnauthorized || out.Len() != 0 {
		t.Fatalf("err = %v, out = %q", err, out.String())
	}
}
