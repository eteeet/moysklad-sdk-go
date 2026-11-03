package processingplan_test

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/processingplan"
	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	testID   = "7944ef04-f831-11e5-7a69-971500188b19"
	rowID    = "9560e3e3-9609-11e6-8af5-581e00000008"
	planPath = "entity/processingplan/" + testID
	apiBase  = "https://api.moysklad.ru/api/remap/1.2/"
	product  = `{"meta":{"href":"` + apiBase + `entity/product/0de151c1-acdc-11e6-5bed-427b00000080","type":"product"}}`
)

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ref(path, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: apiBase + path, Type: model.MetaType(typ)}}
}

var productRef = ref("entity/product/0de151c1-acdc-11e6-5bed-427b00000080", "product")

func service(c *transport.Client) *processingplan.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := processingplan.NewService(c)
	return s
}

// httpCases — операция с моделью в ответе: documenttest.HTTP по всем ResponseCases и LocalErrors.
func httpCases[T any](t *testing.T, method, path, body, response string, invoke func(context.Context, *transport.Client) (*T, error), check func(*T)) {
	t.Helper()
	for _, tc := range documenttest.ResponseCases(response) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, method, path, []byte(body), invoke, check)
		})
	}
	documenttest.LocalErrors(t, invoke)
}

// pageCases — страница набора: те же ResponseCases, ровно один GET с точным query. Операция
// страницы — «list …», поэтому documenttest.HTTP не подходит.
func pageCases[T any](t *testing.T, path, rawQuery, response string, invoke func(context.Context, *transport.Client) (*model.Collection[T], error), check func(*model.Collection[T])) {
	t.Helper()
	for _, tc := range documenttest.ResponseCases(response) {
		t.Run(tc.Name, func(t *testing.T) {
			var calls atomic.Int32
			server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodGet || r.URL.Path != "/api/remap/1.2/"+path || r.URL.RawQuery != rawQuery || r.Header.Get("Accept-Encoding") != "gzip" {
					t.Errorf("request = %s %s", r.Method, r.URL)
				}
				if tc.Gzip {
					w.Header().Set("Content-Encoding", "gzip")
					w.WriteHeader(tc.Status)
					zw := gzip.NewWriter(w)
					io.WriteString(zw, tc.Body)
					zw.Close()
					return
				}
				w.WriteHeader(tc.Status)
				io.WriteString(w, tc.Body)
			}))
			defer server.Close()
			client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2"})
			if err != nil {
				t.Fatal(err)
			}
			got, err := invoke(t.Context(), client)
			// Тело-объект без rows («empty id», «missing id») для страницы — ошибка.
			if tc.Success {
				if err != nil || got == nil {
					t.Fatalf("got = %+v, err = %v", got, err)
				}
				check(got)
			} else {
				var apiErr *sdkerrors.APIError
				if err == nil || got != nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "processingplan "+testID) ||
					errors.As(err, &apiErr) != tc.API {
					t.Fatalf("got = %+v, err = %v", got, err)
				}
			}
			if calls.Load() != 1 {
				t.Fatalf("calls = %d", calls.Load())
			}
		})
	}
	documenttest.LocalErrors(t, invoke)
}

type writeCase struct {
	name, method, path, body, response string
	empty                              bool // ответ без тела: 204 тоже успех
	invoke                             func(context.Context, *transport.Client) error
}

// runWrites — методы без модели в ответе и массовые: один запрос с точным телом ({base} — адрес
// тестового сервера), без повторов при 503.
func runWrites(t *testing.T, cases []writeCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, st := range []struct {
				status int
				ok     bool
			}{{200, true}, {204, tc.empty}, {201, false}, {403, false}, {503, false}} {
				var calls atomic.Int32
				var server *httptest.Server
				server = memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					b, _ := io.ReadAll(r.Body)
					if r.Method != tc.method || r.URL.Path != "/api/remap/1.2/"+tc.path {
						t.Errorf("request = %s %s", r.Method, r.URL)
					}
					if tc.body == "" {
						if len(b) != 0 {
							t.Errorf("body = %s", b)
						}
					} else {
						documenttest.JSONEqual(t, b, []byte(strings.ReplaceAll(tc.body, "{base}", server.URL+"/api/remap/1.2")))
					}
					w.WriteHeader(st.status)
					switch st.status {
					case 200:
						io.WriteString(w, tc.response)
					case 403:
						io.WriteString(w, `{"errors":[{"error":"forbidden","code":1016}]}`)
					}
				}))
				client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", Retry: &transport.RetryPolicy{MaxAttempts: 3}})
				if err != nil {
					t.Fatal(err)
				}
				err = tc.invoke(t.Context(), client)
				server.Close()
				if (err == nil) != st.ok || calls.Load() != 1 {
					t.Fatalf("status %d: err = %v, calls = %d", st.status, err, calls.Load())
				}
				var apiErr *sdkerrors.APIError
				if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "processingplan") ||
					st.status == 403 && !sdkerrors.HasCode(err, 1016) || errors.As(err, &apiErr) != (st.status >= 400)) {
					t.Errorf("status %d: error = %v", st.status, err)
				}
			}
			documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*struct{}, error) {
				if err := tc.invoke(ctx, c); err != nil {
					return nil, err
				}
				return &struct{}{}, nil
			})
		})
	}
}

// noHTTP проверяет ошибки входа: до HTTP, ResponseInfo сброшен, ровно один префикс moysklad: и
// контекст processingplan.
func noHTTP(t *testing.T, cases func(*processingplan.Service, context.Context) map[string]func(...transport.RequestOption) error) {
	t.Helper()
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		io.WriteString(w, `{"id":"server-id"}`)
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	for name, invoke := range cases(service(client), t.Context()) {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "processingplan") {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

// rowMeta — meta строки набора в теле массового удаления.
func rowMeta(path, typ, id string) string {
	return `{"meta":{"href":"{base}/` + planPath + `/` + path + `/` + id + `","type":"` + typ + `","mediaType":"application/json"}}`
}
