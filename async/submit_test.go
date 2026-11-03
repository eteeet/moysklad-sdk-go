package async_test

import (
	"context"
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/async"
	"github.com/eteeet/moysklad-sdk-go/entity/assortment"
	"github.com/eteeet/moysklad-sdk-go/entity/counterparty"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type submitFunc func(context.Context, *transport.Client, query.ListOptions, ...transport.RequestOption) (*async.Submitted, error)

var submits = map[string]submitFunc{
	"counterparty": func(ctx context.Context, c *transport.Client, o query.ListOptions, opts ...transport.RequestOption) (*async.Submitted, error) {
		s, _ := counterparty.NewService(c)
		return s.ListAsync(ctx, o, opts...)
	},
	"assortment": func(ctx context.Context, c *transport.Client, o query.ListOptions, opts ...transport.RequestOption) (*async.Submitted, error) {
		s, _ := assortment.NewService(c)
		return s.ListAsync(ctx, assortment.Options{ListOptions: o}, opts...)
	},
}

// ListAsync: один GET с async=true и фильтром, без limit/offset; у ассортимента — beta-тег.
// 202 с Content-Location и Location даёт задачу; другой ответ — ошибка; 503 при включённых
// повторах не повторяется: создание задачи — побочный эффект.
func TestListAsync(t *testing.T) {
	for entity, submit := range submits {
		t.Run(entity, func(t *testing.T) {
			var calls atomic.Int32
			var status atomic.Int32
			status.Store(http.StatusAccepted)
			server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				q := r.URL.Query()
				if r.Method != http.MethodGet || r.URL.Path != "/entity/"+entity || q.Get("async") != "true" || q.Get("filter") != "name=x" ||
					q.Has("limit") || q.Has("offset") || (entity == "assortment") != (r.Header.Get(transport.HeaderBetaFeature) == assortment.BetaFeature) {
					t.Errorf("request = %s %s %v", r.Method, r.URL, r.Header)
				}
				base := "http://" + r.Host + "/async/" + taskID
				switch status.Load() {
				case http.StatusAccepted:
					w.Header().Set("Content-Location", base)
					w.Header().Set("Location", base+"/result")
				case http.StatusOK:
					w.Write([]byte(`{"meta":{"href":"x"},"rows":[]}`))
					return
				}
				w.WriteHeader(int(status.Load()))
			}))
			defer server.Close()
			c, err := transport.NewClient(transport.Config{BaseURL: server.URL, Retry: &transport.RetryPolicy{MaxAttempts: 3}})
			if err != nil {
				t.Fatal(err)
			}
			var info transport.ResponseInfo
			task, err := submit(t.Context(), c, query.ListOptions{Filter: "name=x"}, transport.WithResponseInfo(&info))
			want := server.URL + "/async/" + taskID
			if err != nil || task.ID != taskID || task.StatusHref != want || task.ResultHref != want+"/result" ||
				info.StatusCode != http.StatusAccepted || calls.Load() != 1 {
				t.Fatalf("task = %+v, err = %v, info = %+v, calls = %d", task, err, info, calls.Load())
			}
			for _, code := range []int{http.StatusOK, http.StatusServiceUnavailable} {
				calls.Store(0)
				status.Store(int32(code))
				if task, err := submit(t.Context(), c, query.ListOptions{Filter: "name=x"}); task != nil || err == nil || calls.Load() != 1 {
					t.Fatalf("status %d: task = %+v, err = %v, calls = %d", code, task, err, calls.Load())
				}
			}
		})
	}
}

// Ответ 202 без адресов задачи — ошибка, а не пустая задача.
func TestListAsyncRequiresTaskHeaders(t *testing.T) {
	for name, headers := range map[string]map[string]string{
		"none":         nil,
		"no location":  {"Content-Location": "https://api.moysklad.ru/api/remap/1.2/async/" + taskID},
		"bad id":       {"Content-Location": "https://api.moysklad.ru/api/remap/1.2/async/x", "Location": "r"},
		"other result": {"Content-Location": "https://api.moysklad.ru/api/remap/1.2/async/" + taskID, "Location": "https://evil.example/x"},
		"not async":    {"Content-Location": "https://api.moysklad.ru/api/remap/1.2/entity/" + taskID, "Location": "https://api.moysklad.ru/api/remap/1.2/entity/" + taskID + "/result"},
		"no status":    {"Location": "https://api.moysklad.ru/api/remap/1.2/async/" + taskID + "/result"},
	} {
		server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for k, v := range headers {
				w.Header().Set(k, v)
			}
			w.WriteHeader(http.StatusAccepted)
		}))
		c, _ := transport.NewClient(transport.Config{BaseURL: server.URL})
		task, err := submits["counterparty"](t.Context(), c, query.ListOptions{})
		server.Close()
		if task != nil || err == nil {
			t.Fatalf("%s: task = %+v, err = %v", name, task, err)
		}
	}
}

// Limit, Offset, fields и async в Extra — ошибка до запроса.
func TestListAsyncRejects(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, _ := transport.NewClient(transport.Config{BaseURL: server.URL})
	for entity, submit := range submits {
		for name, o := range map[string]query.ListOptions{
			"limit":  {Limit: 10},
			"offset": {Offset: 10},
			"fields": {Fields: "stock"},
			"async":  {Extra: url.Values{"async": {"false"}}},
		} {
			info := transport.ResponseInfo{StatusCode: 503}
			if _, err := submit(t.Context(), c, o, transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 {
				t.Fatalf("%s %s: err = %v", entity, name, err)
			}
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}
