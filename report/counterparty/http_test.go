package counterparty_test

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/report/counterparty"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// exchange проверяет один вызов List или Selected: точные method, path, query, заголовки и тело,
// ровно одна отправка даже при Retry клиента. documenttest.HTTP не подходит: контекст ошибки
// он выводит из второго сегмента пути и у POST ждёт «create», а здесь путь — report/counterparty
// и POST — чтение.
func exchange(t *testing.T, rc documenttest.ResponseCase, method, rawQuery string, body []byte,
	invoke func(context.Context, *transport.Client, transport.RequestOption) (*model.Collection[counterparty.Row], error),
	operation string, check func(*model.Collection[counterparty.Row])) {
	t.Helper()
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != method || r.URL.Path != "/api/remap/1.2/report/counterparty" || r.URL.RawQuery != rawQuery {
			t.Errorf("request = %s %s", r.Method, r.URL)
		}
		headers := map[string]string{"Authorization": "Bearer test-token", "Accept": "application/json;charset=utf-8",
			"Accept-Encoding": "gzip", "Content-Type": "", "User-Agent": "moysklad-sdk-go"}
		if body != nil {
			headers["Content-Type"] = "application/json"
		}
		for key, want := range headers {
			if r.Header.Get(key) != want {
				t.Errorf("header %s = %q", key, r.Header.Get(key))
			}
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		if body == nil && len(data) != 0 {
			t.Errorf("body = %s", data)
		}
		if body != nil {
			documenttest.JSONEqual(t, data, body)
		}
		if rc.Gzip {
			w.Header().Set("Content-Encoding", "gzip")
		}
		w.WriteHeader(rc.Status)
		var out io.Writer = w
		if rc.Gzip {
			zw := gzip.NewWriter(w)
			defer zw.Close()
			out = zw
		}
		io.WriteString(out, rc.Body)
	}))
	defer server.Close()
	auth, err := transport.BearerToken("test-token")
	if err != nil {
		t.Fatal(err)
	}
	cfg := transport.Config{BaseURL: server.URL + "/api/remap/1.2", Auth: auth}
	if method == http.MethodPost {
		// Retry клиента включён: POST всё равно уходит один раз.
		cfg.Retry = &transport.RetryPolicy{BaseDelay: time.Millisecond}
	}
	c, err := transport.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	var info transport.ResponseInfo
	got, err := invoke(t.Context(), c, transport.WithResponseInfo(&info))
	if calls.Load() != 1 || info.Attempts != 1 {
		t.Fatalf("calls = %d, attempts = %d, err = %v", calls.Load(), info.Attempts, err)
	}
	if rc.Success {
		if err != nil || got == nil {
			t.Fatalf("result = %+v, err = %v", got, err)
		}
		check(got)
		return
	}
	if err == nil || got != nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), operation) {
		t.Fatalf("result = %+v, err = %v", got, err)
	}
	var apiErr *sdkerrors.APIError
	if errors.As(err, &apiErr) != rc.API || rc.API && apiErr.StatusCode != rc.Status {
		t.Fatalf("API error = %v", err)
	}
	if rc.ErrorCode != 0 && (len(apiErr.Errors) != 1 || apiErr.Errors[0].Code != rc.ErrorCode) {
		t.Errorf("API details = %+v", apiErr)
	}
}

func TestListHTTP(t *testing.T) {
	filter, err := query.Filter(query.Eq("counterparty.inn", "7736570902"), query.Greater("balance", "0"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name    string
		options counterparty.ListOptions
		query   string
	}{
		{"full", counterparty.ListOptions{Limit: 2, Offset: 4, Filter: filter}, "filter=counterparty.inn%3D7736570902%3Bbalance%3E0&limit=2&offset=4"},
		{"defaults", counterparty.ListOptions{}, "limit=1000&offset=0"},
		{"extra", counterparty.ListOptions{Extra: url.Values{"x": {"1"}}}, "limit=1000&offset=0&x=1"},
	} {
		invoke := func(ctx context.Context, c *transport.Client, o transport.RequestOption) (*model.Collection[counterparty.Row], error) {
			return service(c).List(ctx, tc.options, o)
		}
		for _, rc := range documenttest.ResponseCases(string(read(t, "custom/reportcounterparty/list.json"))) {
			t.Run(tc.name+"/"+rc.Name, func(t *testing.T) {
				exchange(t, rc, http.MethodGet, tc.query, nil, invoke, "get counterparty report", func(got *model.Collection[counterparty.Row]) {
					if len(got.Rows) != 5 || got.Meta.Size != 5 || got.Rows[4].Counterparty.ExternalCode != "aZBfWOKzj-lcq7c7IWZON3" {
						t.Fatalf("page = %+v", got)
					}
				})
			})
		}
		documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*model.Collection[counterparty.Row], error) {
			return service(c).List(ctx, tc.options)
		})
	}
}

func selectedRefs() []model.Reference {
	var refs []model.Reference
	for _, id := range []string{"49b40d55-66cc-11e7-6adb-ede500000054", "37e4efe0-6ade-11e7-6adb-ede50000001a"} {
		m := meta("entity/counterparty/"+id, "counterparty")
		m.MetadataHref = host + "entity/counterparty/metadata"
		refs = append(refs, model.Reference{Meta: m})
	}
	return refs
}

// Selected — чтение POST-ом: тело по golden, без query, одна отправка.
func TestSelectedHTTP(t *testing.T) {
	golden := read(t, "custom/reportcounterparty/selected.golden.json")
	invoke := func(ctx context.Context, c *transport.Client, o transport.RequestOption) (*model.Collection[counterparty.Row], error) {
		return service(c).Selected(ctx, selectedRefs(), o)
	}
	for _, rc := range documenttest.ResponseCases(string(read(t, "custom/reportcounterparty/selected.json"))) {
		t.Run(rc.Name, func(t *testing.T) {
			exchange(t, rc, http.MethodPost, "", golden, invoke, "select counterparty report", func(got *model.Collection[counterparty.Row]) {
				if len(got.Rows) != 2 || got.Meta.Size != 2 || got.Rows[1].Counterparty.Name != `ООО "Поставщик 2"` ||
					got.Rows[0].Meta.UUIDHref != "https://online.moysklad.ru/app/#company/edit?id=49b40d55-66cc-11e7-6adb-ede500000054" {
					t.Fatalf("page = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*model.Collection[counterparty.Row], error) {
		return service(c).Selected(ctx, selectedRefs())
	})
}

func TestGetHTTP(t *testing.T) {
	const id = "df2fdd2d-6934-11e6-8a84-bae500000049"
	invoke := func(ctx context.Context, c *transport.Client) (*counterparty.Row, error) {
		return service(c).Get(ctx, id)
	}
	for _, rc := range documenttest.ResponseCases(string(read(t, "custom/reportcounterparty/get.json"))) {
		t.Run(rc.Name, func(t *testing.T) {
			documenttest.HTTP(t, rc, "GET", "report/counterparty/"+id, nil, invoke, func(got *counterparty.Row) {
				if got.Counterparty.ID != id || got.Balance != -60000 || got.BonusBalance != 200 || got.LastEventText != "продажа" {
					t.Fatalf("row = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}
