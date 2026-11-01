// Package listtest проверяет HTTP-контракт List одной страницы для любой сущности.
// Тела ответов синтетические: одна строка — настоящая фикстура сущности.
package listtest

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type ListFunc[T any] func(context.Context, query.ListOptions, ...transport.RequestOption) (*model.Collection[T], error)

// Run проверяет List на /entity/{entity}. fields — допустимое значение fields или "", если
// endpoint его не поддерживает; both — поле, по которому true и false вместе отбирают всё
// (archived у справочников с архивом, isDeleted у документов) или ""; check проверяет строку.
func Run[T any](t *testing.T, entity string, row []byte, fields, both string,
	newList func(*transport.Client) ListFunc[T], check func(*testing.T, T)) {
	t.Helper()
	const pageMeta = `"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/x?limit=2","type":"x","size":5,"limit":2,"offset":0,"nextHref":"https://api.moysklad.ru/api/remap/1.2/entity/x?limit=2&offset=2"}`
	page := []byte(`{"context":{"employee":{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/context/employee"}}},` + pageMeta + `,"rows":[` + string(row) + `]}`)

	conditions := []query.Condition{query.Eq("name", "a;b")}
	filterQuery := "filter=name%3Da%5C%3Bb"
	if both != "" {
		conditions = append(conditions, query.Eq(both, "true"), query.Eq(both, "false"))
		filterQuery += "%3B" + both + "%3Dtrue%3B" + both + "%3Dfalse"
	}
	filter, err := query.Filter(conditions...)
	if err != nil {
		t.Fatal(err)
	}
	full := query.ListOptions{Limit: 2, Offset: 4, Filter: filter, Order: "name,desc", Search: "кофе 100%",
		NamedFilter: "https://api.moysklad.ru/api/remap/1.2/entity/" + entity + "/namedfilter/1", Fields: fields}
	fullQuery := filterQuery + "&limit=2" +
		"&namedfilter=https%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2F" + entity + "%2Fnamedfilter%2F1" +
		"&offset=4&order=name%2Cdesc&search=%D0%BA%D0%BE%D1%84%D0%B5+100%25"
	if fields != "" {
		fullQuery = "fields=" + fields + "&" + fullQuery
	}

	for _, tc := range []struct {
		name    string
		options query.ListOptions
		query   string
		status  int
		body    []byte
		gzip    bool
		rows    int // -1 — ожидается ошибка
		apiErr  bool
	}{
		{"no options", query.ListOptions{}, "", 200, page, false, 1, false},
		{"all options", full, fullQuery, 200, page, false, 1, false},
		{"gzip", query.ListOptions{Limit: 1000}, "limit=1000", 200, page, true, 1, false},
		{"empty", query.ListOptions{}, "", 200, []byte(`{` + pageMeta + `,"rows":[]}`), false, 0, false},
		{"null envelope", query.ListOptions{}, "", 200, []byte(`null`), false, -1, false},
		{"null rows", query.ListOptions{}, "", 200, []byte(`{` + pageMeta + `,"rows":null}`), false, -1, false},
		{"plain array", query.ListOptions{}, "", 200, []byte(`[` + string(row) + `]`), false, -1, false},
		{"malformed", query.ListOptions{}, "", 200, []byte(`{"meta":`), false, -1, false},
		{"no content", query.ListOptions{}, "", 204, nil, false, -1, false},
		{"api error", query.ListOptions{}, "", 400, []byte(`{"errors":[{"error":"bad filter","code":1014}]}`), false, -1, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodGet || r.URL.Path != "/api/remap/1.2/entity/"+entity || r.URL.RawQuery != tc.query {
					t.Errorf("request = %s %s", r.Method, r.URL.RequestURI())
				}
				if r.Header.Get("Authorization") != "Bearer test-token" || r.Header.Get("Accept-Encoding") != "gzip" ||
					r.Header.Get("Accept-Language") != "en" || r.Header.Get("Content-Type") != "" {
					t.Errorf("headers = %v", r.Header)
				}
				w.Header().Set("X-Lognex-Content-Region", "kz")
				body := tc.body
				if tc.gzip {
					var buf bytes.Buffer
					zw := gzip.NewWriter(&buf)
					zw.Write(body)
					zw.Close()
					body = buf.Bytes()
					w.Header().Set("Content-Encoding", "gzip")
				}
				w.WriteHeader(tc.status)
				w.Write(body)
			}))
			defer server.Close()
			list := newList(newClient(t, server.URL))
			var info transport.ResponseInfo
			got, err := list(t.Context(), tc.options, transport.WithAcceptLanguage("en"), transport.WithResponseInfo(&info))
			if calls.Load() != 1 || info.StatusCode != tc.status || info.Header.Get("X-Lognex-Content-Region") != "kz" {
				t.Fatalf("calls = %d, info = %+v", calls.Load(), info)
			}
			if tc.rows < 0 {
				var apiErr *sdkerrors.APIError
				if got != nil || err == nil || errors.As(err, &apiErr) != tc.apiErr || strings.Count(err.Error(), "moysklad:") != 1 {
					t.Fatalf("got = %+v, err = %v", got, err)
				}
				return
			}
			if err != nil || len(got.Rows) != tc.rows || got.Meta.Href == "" {
				t.Fatalf("got = %+v, err = %v", got, err)
			}
			if tc.rows > 0 {
				if got.Meta.Size != 5 || got.Meta.NextHref == "" || got.Context == nil {
					t.Fatalf("meta = %+v, context = %+v", got.Meta, got.Context)
				}
				check(t, got.Rows[0])
			}
		})
	}

	// Ошибки опций не доходят до сервера.
	rejected := []query.ListOptions{{Limit: 1001}, {Offset: -1}, {Fields: "stock,declaration"},
		{Extra: map[string][]string{"expand": {"agent"}}}}
	if fields == "" {
		rejected = append(rejected, query.ListOptions{Fields: "stock"})
	} else {
		rejected = append(rejected, query.ListOptions{Fields: fields + "x"})
	}
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	list := newList(newClient(t, server.URL))
	for _, options := range rejected {
		if got, err := list(t.Context(), options); got != nil || err == nil || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("options %+v: got = %+v, err = %v", options, got, err)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("invalid options sent %d requests", calls.Load())
	}

	// Локальные ошибки: nil и отменённый контекст, ошибка Auth, сбой транспорта.
	documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*model.Collection[T], error) {
		return newList(c)(ctx, query.ListOptions{})
	})
}

func newClient(t *testing.T, url string) *transport.Client {
	t.Helper()
	auth, err := transport.BearerToken("test-token")
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.NewClient(transport.Config{BaseURL: url + "/api/remap/1.2", Auth: auth})
	if err != nil {
		t.Fatal(err)
	}
	return client
}
