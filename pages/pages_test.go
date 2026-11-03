package pages_test

import (
	"context"
	"encoding/json/v2"
	"errors"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/pages"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type row struct {
	ID string `json:"id"`
}

// Позиции документа: путь коллекции общий у meta.href и nextHref, запрос фильтра уже закодирован.
const (
	coll  = "/entity/demand/77777777-7777-7777-7777-777777777777/positions"
	query = "?filter=name%3Da%5C%3Bb&limit=1"
)

// newServer отдаёт тело по RequestURI; {base} в теле заменяется адресом сервера.
// Каждый запрос должен быть GET с заголовком X-Call из опций обхода.
func newServer(t *testing.T, bodies map[string]string) (*transport.Client, string, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	var base string
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		body, ok := bodies[r.URL.RequestURI()]
		if !ok || r.Method != http.MethodGet || r.Header.Get("X-Call") != "walk" {
			t.Errorf("unexpected request %s %s %v", r.Method, r.URL.RequestURI(), r.Header)
		}
		w.Header().Set("X-Page", string(rune('0'+n)))
		w.Write([]byte(strings.ReplaceAll(body, "{base}", base)))
	}))
	t.Cleanup(server.Close)
	base = server.URL
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	return client, base, &calls
}

func firstPage(t *testing.T, base, body string) *model.Collection[row] {
	t.Helper()
	var page model.Collection[row]
	if err := json.Unmarshal([]byte(strings.ReplaceAll(body, "{base}", base)), &page); err != nil {
		t.Fatal(err)
	}
	return &page
}

func walk(ctx context.Context, client *transport.Client, first *model.Collection[row], info *transport.ResponseInfo) (ids []string, err error) {
	for page, e := range pages.All(ctx, client, first, transport.WithHeader("X-Call", "walk"), transport.WithResponseInfo(info)) {
		if e != nil {
			return ids, e
		}
		for _, r := range page.Rows {
			ids = append(ids, r.ID)
		}
	}
	return ids, nil
}

func page(offset, next string, rows string) string {
	body := `{"meta":{"href":"{base}` + coll + query + offset + `"`
	if next != "" {
		body += `,"nextHref":"{base}` + coll + query + next + `"`
	}
	return body + `},"rows":[` + rows + `]}`
}

func TestAll(t *testing.T) {
	client, base, calls := newServer(t, map[string]string{
		coll + query + "&offset=1": page("&offset=1", "&offset=2", `{"id":"2"}`),
		// Строки удалены параллельно: пустая страница с nextHref не останавливает обход.
		coll + query + "&offset=2": page("&offset=2", "&offset=3", ``),
		coll + query + "&offset=3": page("&offset=3", "", `{"id":"3"}`),
	})
	var info transport.ResponseInfo
	ids, err := walk(t.Context(), client, firstPage(t, base, page("", "&offset=1", `{"id":"1"}`)), &info)
	if err != nil || strings.Join(ids, ",") != "1,2,3" || calls.Load() != 3 || info.Header.Get("X-Page") != "3" {
		t.Fatalf("ids = %v, err = %v, calls = %d, info = %v", ids, err, calls.Load(), info.Header)
	}
}

func TestAllEarlyStopAndCancel(t *testing.T) {
	client, base, calls := newServer(t, nil)
	first := firstPage(t, base, page("", "&offset=1", `{"id":"1"}`))
	for range pages.All(t.Context(), client, first) {
		break
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	var err error
	for _, e := range pages.All(ctx, client, first) {
		cancel() // отмена после первой страницы, до запроса второй
		err = e
	}
	if !errors.Is(err, context.Canceled) || calls.Load() != 0 {
		t.Fatalf("err = %v, calls = %d", err, calls.Load())
	}

	n := 0
	for _, err := range pages.All(ctx, client, first) {
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("first page yielded with canceled context: %v", err)
		}
		n++
	}
	if n != 1 || calls.Load() != 0 {
		t.Fatalf("yields = %d, calls = %d", n, calls.Load())
	}
}

func TestAllRejects(t *testing.T) {
	for _, tc := range []struct {
		name, first string
		bodies      map[string]string
		calls       int32
		want        string
	}{
		{"self link", `{"meta":{"href":"{base}` + coll + query + `","nextHref":"{base}` + coll + query + `"},"rows":[{"id":"1"}]}`, nil, 0, "does not advance"},
		{"offset goes back", page("&offset=1", "&offset=2", `{"id":"1"}`), map[string]string{
			coll + query + "&offset=2": page("&offset=2", "&offset=1", `{"id":"2"}`),
		}, 1, "does not advance"},
		{"invalid offset", page("", "&offset=x", `{"id":"1"}`), nil, 0, "invalid offset"},
		{"foreign host", `{"meta":{"href":"{base}` + coll + `","nextHref":"https://evil.example` + coll + `?offset=1"},"rows":[{"id":"1"}]}`, nil, 0, "leaves the collection"},
		{"other collection", `{"meta":{"href":"{base}` + coll + `","nextHref":"{base}/entity/employee?offset=1"},"rows":[{"id":"1"}]}`, nil, 0, "leaves the collection"},
		{"relative nextHref", `{"meta":{"href":"{base}` + coll + `","nextHref":"` + coll + `?offset=1"},"rows":[{"id":"1"}]}`, nil, 0, "leaves the collection"},
		{"missing rows", page("", "&offset=1", `{"id":"1"}`), map[string]string{
			coll + query + "&offset=1": `{"meta":{"href":"{base}/x"}}`,
		}, 1, "non-null array"},
		{"wrong row type", page("", "&offset=1", `{"id":"1"}`), map[string]string{
			coll + query + "&offset=1": `{"meta":{"href":"{base}/x"},"rows":"x"}`,
		}, 1, "decode response"},
		{"no meta.href", `{"meta":{},"rows":[{"id":"1"}]}`, nil, 0, "invalid meta.href"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client, base, calls := newServer(t, tc.bodies)
			var info transport.ResponseInfo
			_, err := walk(t.Context(), client, firstPage(t, base, tc.first), &info)
			if err == nil || !strings.Contains(err.Error(), tc.want) || strings.Count(err.Error(), "moysklad:") != 1 || calls.Load() != tc.calls {
				t.Fatalf("err = %v, calls = %d", err, calls.Load())
			}
		})
	}

	client, _, _ := newServer(t, nil)
	for _, err := range pages.All[row](t.Context(), client, nil) {
		if err == nil {
			t.Fatal("nil first page accepted")
		}
	}
}

// Страницы без size и без nextHref на первой странице: обход без запросов.
func TestAllSinglePage(t *testing.T) {
	client, base, calls := newServer(t, nil)
	var info transport.ResponseInfo
	ids, err := walk(t.Context(), client, firstPage(t, base, page("", "", `{"id":"1"}`)), &info)
	if err != nil || strings.Join(ids, ",") != "1" || calls.Load() != 0 {
		t.Fatalf("ids = %v, err = %v, calls = %d", ids, err, calls.Load())
	}
}
