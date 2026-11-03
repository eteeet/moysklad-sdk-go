package list_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/list"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

type row struct {
	ID string `json:"id"`
}

func TestPage(t *testing.T) {
	filter, err := query.Filter(query.Eq("name", "a;b"))
	if err != nil {
		t.Fatal(err)
	}
	params, err := query.ListOptions{Limit: 2, Filter: filter}.Values(query.MaxLimit)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, body       string
		status, wantRows int
		success, apiErr  bool
	}{
		{"page", `{"context":{"employee":{"meta":{"href":"e"}}},"meta":{"href":"h","size":3,"limit":2,"offset":0,"nextHref":"n"},"rows":[{"id":"1"},{"id":"2"}]}`, 200, 2, true, false},
		{"empty", `{"meta":{"href":"h","size":0},"rows":[]}`, 200, 0, true, false},
		{"plain array", `[{"id":"1"}]`, 200, 0, false, false},
		{"null rows", `{"meta":{"href":"h"},"rows":null}`, 200, 0, false, false},
		{"no meta", `{"rows":[]}`, 200, 0, false, false},
		{"unexpected status", `{"meta":{"href":"h"},"rows":[]}`, 202, 0, false, false},
		{"api error", `{"errors":[{"error":"bad","code":1014}]}`, 400, 0, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodGet || r.URL.RequestURI() != "/entity/product?filter=name%3Da%5C%3Bb&limit=2" {
					t.Errorf("request = %s %s", r.Method, r.URL.RequestURI())
				}
				w.WriteHeader(tc.status)
				w.Write([]byte(tc.body))
			}))
			defer server.Close()
			client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			var info transport.ResponseInfo
			page, err := list.Page[row](t.Context(), client, "list product", "entity/product", params, transport.WithResponseInfo(&info))
			if tc.success != (err == nil) || calls.Load() != 1 || info.StatusCode != tc.status {
				t.Fatalf("page = %+v, err = %v, calls = %d, info = %+v", page, err, calls.Load(), info)
			}
			if err != nil {
				var apiErr *sdkerrors.APIError
				if errors.As(err, &apiErr) != tc.apiErr || strings.Count(err.Error(), "moysklad:") != 1 {
					t.Fatalf("err = %v", err)
				}
				return
			}
			if page.Meta.Href != "h" || len(page.Rows) != tc.wantRows || tc.wantRows > 0 && (page.Context == nil || page.Meta.NextHref != "n") {
				t.Fatalf("page = %+v", page)
			}
		})
	}
}

func TestPageCanceled(t *testing.T) {
	received := make(chan struct{})
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(received)
		<-r.Context().Done()
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go func() {
		<-received
		cancel()
	}()
	if _, err := list.Page[row](ctx, client, "list product", "entity/product", nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v", err)
	}
}
