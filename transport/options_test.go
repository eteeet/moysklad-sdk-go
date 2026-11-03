package transport_test

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func newOptionsClient(t *testing.T, header http.Header, handler http.HandlerFunc) *transport.Client {
	t.Helper()
	server := memnet.NewServer(handler)
	t.Cleanup(server.Close)
	auth, err := transport.BearerToken("secret")
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL, Auth: auth, Header: header})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestSendHeaderPrecedence(t *testing.T) {
	var got http.Header
	clientHeader := http.Header{
		"User-Agent":             {"app/1"},
		"X-Lognex-Accept-Region": {"kz"},
		"Accept-Language":        {"en"},
	}
	client := newOptionsClient(t, clientHeader, func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
		w.Write([]byte(`{}`))
	})
	clientHeader.Set("X-Lognex-Accept-Region", "uz") // клиент хранит копию

	_, err := client.Send(t.Context(), http.MethodGet, "entity/uom", nil, nil, nil,
		nil, // nil-опция игнорируется
		transport.WithHeader("accept-language", "ru"),
		transport.WithHeader("X-Lognex-WebHook-DisableByPrefix", "https://a.example"),
		transport.WithHeader("X-Lognex-WebHook-DisableByPrefix", "https://b.example"),
	)
	if err != nil {
		t.Fatal(err)
	}

	for key, want := range map[string]string{
		"User-Agent":                       "app/1",
		"X-Lognex-Accept-Region":           "kz",
		"Accept-Language":                  "ru",
		"X-Lognex-Webhook-Disablebyprefix": "https://a.example|https://b.example",
		"Authorization":                    "Bearer secret",
		"Accept":                           "application/json;charset=utf-8",
		"Accept-Encoding":                  "gzip",
	} {
		if v := strings.Join(got.Values(key), "|"); v != want {
			t.Errorf("%s = %q, want %q", key, v, want)
		}
	}

	// Опции одного вызова не переживают его.
	if _, err := client.Send(t.Context(), http.MethodGet, "entity/uom", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	if got.Get("Accept-Language") != "en" || got.Get("X-Lognex-Webhook-Disablebyprefix") != "" {
		t.Fatalf("options leaked into next call: %v", got)
	}
}

func TestHeaderValues(t *testing.T) {
	c, err := transport.NewClient(transport.Config{Header: http.Header{"x-lognex-remap-beta-feature": {"a"}}})
	if err != nil {
		t.Fatal(err)
	}
	for name, tc := range map[string]struct {
		opts []transport.RequestOption
		want []string
	}{
		"config":   {nil, []string{"a"}},
		"per call": {[]transport.RequestOption{transport.WithBetaFeature("b"), nil}, []string{"b"}},
		"other":    {[]transport.RequestOption{transport.WithAcceptLanguage("en")}, []string{"a"}},
	} {
		if got := c.HeaderValues(transport.HeaderBetaFeature, tc.opts...); !slices.Equal(got, tc.want) {
			t.Errorf("%s: values = %q, want %q", name, got, tc.want)
		}
	}
}

func TestProtectedHeaders(t *testing.T) {
	for _, key := range []string{"Authorization", "accept", "Accept-Encoding", "content-type", "Content-Length", "Host",
		"Idempotency-Key", "x-idempotency-key"} {
		t.Run(key, func(t *testing.T) {
			if _, err := transport.NewClient(transport.Config{Header: http.Header{key: {"x"}}}); err == nil {
				t.Fatal("NewClient accepted protected header")
			}

			var calls atomic.Int32
			client := newOptionsClient(t, nil, func(w http.ResponseWriter, r *http.Request) { calls.Add(1) })
			info := transport.ResponseInfo{StatusCode: 999}
			_, err := client.Send(t.Context(), http.MethodGet, "entity/uom", nil, nil, nil,
				transport.WithResponseInfo(&info), transport.WithHeader(key, "x"))
			if err == nil || !strings.HasPrefix(err.Error(), "moysklad: ") {
				t.Fatalf("err = %v", err)
			}
			if calls.Load() != 0 {
				t.Fatal("request sent")
			}
			if info.StatusCode != 0 || info.Header != nil {
				t.Fatalf("stale info = %+v", info)
			}
		})
	}
}

func TestSendResponseInfo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newOptionsClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Lognex-Unknown", "a")
			w.Header().Add("X-Lognex-Unknown", "b")
			w.Header().Set("X-RateLimit-Remaining", "not-a-number")
			w.Write([]byte(`{"id":"1"}`))
		})
		var info transport.ResponseInfo
		var dst struct {
			ID string `json:"id"`
		}
		resp, err := client.Send(t.Context(), http.MethodGet, "entity/uom", nil, nil, &dst, transport.WithResponseInfo(&info))
		if err != nil || dst.ID != "1" {
			t.Fatalf("dst = %+v, err = %v", dst, err)
		}
		if info.StatusCode != http.StatusOK ||
			strings.Join(info.Header.Values("X-Lognex-Unknown"), ",") != "a,b" ||
			info.Header.Get("X-RateLimit-Remaining") != "not-a-number" {
			t.Fatalf("info = %+v", info)
		}
		info.Header.Set("X-Lognex-Unknown", "changed")
		if resp.Header.Get("X-Lognex-Unknown") != "a" {
			t.Fatal("info.Header is not a copy")
		}
	})

	t.Run("api error", func(t *testing.T) {
		client := newOptionsClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Lognex-Retry-After", "1500")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"errors":[{"error":"slow","code":1073}]}`))
		})
		var info transport.ResponseInfo
		_, err := client.Send(t.Context(), http.MethodGet, "entity/uom", nil, nil, &struct{}{}, transport.WithResponseInfo(&info))
		var apiErr *sdkerrors.APIError
		if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusTooManyRequests {
			t.Fatalf("err = %v", err)
		}
		if info.StatusCode != http.StatusTooManyRequests || info.Header.Get("X-Lognex-Retry-After") != "1500" {
			t.Fatalf("info = %+v", info)
		}
	})

	t.Run("decode error", func(t *testing.T) {
		client := newOptionsClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Request", "kept")
			w.Write([]byte(`{"id":`))
		})
		var info transport.ResponseInfo
		_, err := client.Send(t.Context(), http.MethodGet, "entity/uom", nil, nil, &struct{}{}, transport.WithResponseInfo(&info))
		if err == nil || info.StatusCode != http.StatusOK || info.Header.Get("X-Request") != "kept" {
			t.Fatalf("info = %+v, err = %v", info, err)
		}
	})

	t.Run("no response resets info", func(t *testing.T) {
		server := memnet.NewServer(http.NotFoundHandler())
		client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
		if err != nil {
			t.Fatal(err)
		}
		server.Close()
		info := transport.ResponseInfo{StatusCode: 200, Header: http.Header{"Stale": {"1"}}}
		if _, err := client.Send(t.Context(), http.MethodGet, "entity/uom", nil, nil, nil, transport.WithResponseInfo(&info)); err == nil {
			t.Fatal("expected transport error")
		}
		if info.StatusCode != 0 || info.Header != nil {
			t.Fatalf("stale info = %+v", info)
		}
	})
}

func TestSendParallelResponseInfo(t *testing.T) {
	client := newOptionsClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Echo", r.Header.Get("X-Call"))
		w.Write([]byte(`{}`))
	})
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Go(func() {
			call := fmt.Sprint(i)
			var info transport.ResponseInfo
			_, err := client.Send(t.Context(), http.MethodGet, "entity/uom", nil, nil, nil,
				transport.WithHeader("X-Call", call), transport.WithResponseInfo(&info))
			if err != nil || info.Header.Get("X-Echo") != call {
				t.Errorf("call %s: info = %+v, err = %v", call, info, err)
			}
		})
	}
	wg.Wait()
}

func TestHeaderCaseDuplicates(t *testing.T) {
	header := http.Header{"accept-language": {"en"}, "Accept-Language": {"ru"}}
	if _, err := transport.NewClient(transport.Config{Header: header}); err == nil {
		t.Fatal("NewClient accepted ambiguous header")
	}
	var calls atomic.Int32
	client := newOptionsClient(t, nil, func(w http.ResponseWriter, r *http.Request) { calls.Add(1) })
	raw := func(o *transport.RequestOptions) { o.Header = header }
	if _, err := client.Send(t.Context(), http.MethodGet, "entity/uom", nil, nil, nil, raw); err == nil || calls.Load() != 0 {
		t.Fatalf("err = %v, calls = %d", err, calls.Load())
	}
}

// Services поверх одного Client делят его настройки; лимитер и retry T08/T09 живут там же.
func TestServicesShareClient(t *testing.T) {
	var regions []string
	var mu sync.Mutex
	client := newOptionsClient(t, http.Header{"X-Lognex-Accept-Region": {"kz"}}, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		regions = append(regions, r.URL.Path+"="+r.Header.Get("X-Lognex-Accept-Region"))
		mu.Unlock()
	})
	for _, endpoint := range []string{"entity/uom", "entity/product"} {
		if _, err := client.Send(t.Context(), http.MethodGet, endpoint, nil, nil, nil); err != nil {
			t.Fatal(err)
		}
	}
	if strings.Join(regions, ",") != "/entity/uom=kz,/entity/product=kz" {
		t.Fatalf("regions = %v", regions)
	}
}

func TestWithAuth(t *testing.T) {
	var got []string
	client := newOptionsClient(t, nil, func(w http.ResponseWriter, r *http.Request) {
		got = append(got, r.Header.Get("Authorization"))
	})
	basic, err := transport.BasicAuth("u", "p")
	if err != nil {
		t.Fatal(err)
	}
	for _, opts := range [][]transport.RequestOption{{transport.WithAuth(basic)}, nil} {
		if _, err := client.Send(t.Context(), http.MethodGet, "entity/uom", nil, nil, nil, opts...); err != nil {
			t.Fatal(err)
		}
	}
	if strings.Join(got, ",") != "Basic dTpw,Bearer secret" {
		t.Fatalf("Authorization = %v", got)
	}
}
