package transport_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestNewRequest(t *testing.T) {
	auth, err := transport.BearerToken("secret")
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.NewClient(transport.Config{Auth: auth})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := url.Values{
		"limit":  {"10"},
		"filter": {"name=Товар;archived=false"},
		"tag":    {"a", "b"},
	}
	before := query.Encode()
	body := struct {
		Name string `json:"name"`
	}{Name: "Товар"}

	req, err := client.NewRequest(ctx, http.MethodPost, "entity/product?limit=1&offset=20", query, body)
	if err != nil {
		t.Fatal(err)
	}
	defer req.Body.Close()

	if req.Method != http.MethodPost || req.Context() != ctx {
		t.Fatal("method or context lost")
	}
	if req.URL.Scheme != "https" || req.URL.Host != "api.moysklad.ru" ||
		req.URL.Path != "/api/remap/1.2/entity/product" {
		t.Fatalf("URL = %s", req.URL)
	}

	params := req.URL.Query()
	if params.Get("limit") != "10" || params.Get("offset") != "20" ||
		params.Get("filter") != query.Get("filter") || strings.Join(params["tag"], ",") != "a,b" {
		t.Fatalf("query = %v", params)
	}
	if query.Encode() != before {
		t.Fatal("input query mutated")
	}

	for name, want := range map[string]string{
		"Accept":          "application/json",
		"Accept-Encoding": "gzip",
		"Content-Type":    "application/json",
		"User-Agent":      "moysklad-sdk-go",
		"Authorization":   "Bearer secret",
	} {
		if got := req.Header.Get(name); got != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"name":"Товар"}` {
		t.Fatalf("body = %s", data)
	}
	if req.GetBody == nil {
		t.Fatal("GetBody is nil")
	}

	replay, err := req.GetBody()
	if err != nil {
		t.Fatal(err)
	}
	defer replay.Close()

	replayed, err := io.ReadAll(replay)
	if err != nil {
		t.Fatal(err)
	}
	if string(replayed) != string(data) {
		t.Fatal("GetBody returned different content")
	}
}

func TestNewRequestURL(t *testing.T) {
	base := transport.DefaultBaseURL
	tests := []struct {
		name     string
		endpoint string
		wantURL  string
	}{
		{"relative", "entity/product", base + "/entity/product"},
		{"leading slash", "/entity/product", base + "/entity/product"},
		{"absolute", base + "/entity/product", base + "/entity/product"},
		{"next page", base + "/entity/product?limit=10&offset=20", base + "/entity/product?limit=10&offset=20"},
		{
			"next page with filter",
			base + "/entity/product?filter=archived=false;type=product&offset=2&limit=1",
			base + "/entity/product?filter=archived=false;type=product&offset=2&limit=1",
		},
		{
			"raw filter",
			"entity/metadata?filter=type=product;type=service",
			base + "/entity/metadata?filter=type=product;type=service",
		},

		{"escaped slash", "entity/product/a%2Fb", base + "/entity/product/a%2Fb"},
		{"default port", "https://api.moysklad.ru:443/api/remap/1.2/entity/product", "https://api.moysklad.ru:443/api/remap/1.2/entity/product"},
		{"empty", "", ""},
		{"query only", "?limit=10", ""},
		{"foreign host", "https://example.com/entity/product", ""},
		{"subdomain", "https://files.api.moysklad.ru/entity/product", ""},
		{"foreign port", "https://api.moysklad.ru:444/entity/product", ""},
		{"downgrade", "http://api.moysklad.ru/entity/product", ""},
		{"scheme relative", "//api.moysklad.ru/entity/product", ""},
		{"credentials", "https://user:pass@api.moysklad.ru/entity/product", ""},
		{"fragment", "entity/product#section", ""},
		{"parent path", "../product", ""},
		{"encoded parent", "entity/%2e%2e/product", ""},
		{"opaque URL", "https:entity/product", ""},
		{"invalid query", "entity/product?value=%zz", ""},
		{"invalid URL", "%zz", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authorized := false
			client, err := transport.NewClient(transport.Config{
				Auth: func(*http.Request) error {
					authorized = true
					return nil
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			req, err := client.NewRequest(context.Background(), http.MethodGet, tt.endpoint, nil, nil)
			if tt.wantURL == "" {
				if err == nil {
					t.Fatal("expected endpoint error")
				}
				if authorized {
					t.Fatal("auth called for invalid endpoint")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if req.URL.String() != tt.wantURL {
				t.Fatalf("URL = %s, want %s", req.URL, tt.wantURL)
			}
			if !authorized {
				t.Fatal("auth was not called")
			}
			if req.Body != nil || req.Header.Get("Content-Type") != "" {
				t.Fatal("body or Content-Type set for request without body")
			}
		})
	}
}

func TestNewRequestErrors(t *testing.T) {
	client, err := transport.NewClient(transport.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.NewRequest(context.Background(), http.MethodPost, "entity/product", nil, make(chan int)); err == nil {
		t.Fatal("expected JSON encoding error")
	}
	if _, err := client.NewRequest(context.Background(), "bad method", "entity/product", nil, nil); err == nil {
		t.Fatal("expected method error")
	}
	if _, err := client.NewRequest(nil, http.MethodGet, "entity/product", nil, nil); err == nil {
		t.Fatal("expected context error")
	}

	wantErr := errors.New("auth failed")
	client, err = transport.NewClient(transport.Config{
		Auth: func(*http.Request) error { return wantErr },
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.NewRequest(context.Background(), http.MethodGet, "entity/product", nil, nil); !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapped auth error", err)
	}
}
