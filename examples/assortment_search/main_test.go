package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/examples/internal/exampletest"
)

// Две страницы: вторая — GET по meta.nextHref с тем же beta-тегом; отмена ctx останавливает обход до следующей страницы.
func TestFind(t *testing.T) {
	var calls atomic.Int32
	c := exampletest.Client(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		q := r.URL.Query()
		if r.URL.Path != "/api/remap/1.2/entity/assortment" || q.Get("search") != "чай" || q.Get("limit") != "100" ||
			r.Header.Get("X-Lognex-Remap-Beta-Feature") != "assortmentWithoutStock" {
			t.Errorf("request = %s %v", r.URL, r.Header)
		}
		href := "http://" + r.Host + "/api/remap/1.2/entity/assortment?limit=100&search=%D1%87%D0%B0%D0%B9&offset="
		if q.Get("offset") == "" || q.Get("offset") == "0" {
			fmt.Fprintf(w, `{"meta":{"href":%q,"nextHref":%q},"rows":[{"meta":{"type":"product"},"id":"p1","name":"Чай чёрный"}]}`,
				href+"0", href+"100")
			return
		}
		fmt.Fprintf(w, `{"meta":{"href":%q},"rows":[{"meta":{"type":"service"},"id":"s1","name":"Дегустация чая"}]}`, href+"100")
	})
	var out bytes.Buffer
	if err := find(t.Context(), c, "чай", &out); err != nil || out.String() != "product Чай чёрный\nservice Дегустация чая\n" || calls.Load() != 2 {
		t.Fatalf("out = %q, err = %v, calls = %d", out.String(), err, calls.Load())
	}

	ctx, cancel := context.WithCancel(t.Context())
	calls.Store(0)
	out.Reset()
	err := find(ctx, c, "чай", writerFunc(func(p []byte) (int, error) { cancel(); return out.Write(p) }))
	if !errors.Is(err, context.Canceled) || calls.Load() != 1 {
		t.Fatalf("cancel: err = %v, calls = %d", err, calls.Load())
	}
	if strings.Contains(out.String(), "service") {
		t.Fatalf("second page after cancel: %q", out.String())
	}
}

type writerFunc func([]byte) (int, error)

func (f writerFunc) Write(p []byte) (int, error) { return f(p) }
