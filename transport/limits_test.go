package transport_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/transport"
)

func limitsClient(t *testing.T, cfg transport.Config) *transport.Client {
	t.Helper()
	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://example.test"
	}
	client, err := transport.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// Размер — байты закодированного JSON: Base64 и UTF-8 считаются по факту, граница max/max+1.
func TestBodySizeLimit(t *testing.T) {
	type file struct {
		Content string `json:"content"`
	}
	// {"content":"…"} — 14 байт обёртки; 75 байт файла дают 100 байт Base64.
	body := file{Content: base64.StdEncoding.EncodeToString(make([]byte, 75))}
	for _, tc := range []struct {
		max int
		ok  bool
	}{{114, true}, {113, false}, {-1, true}} {
		_, err := limitsClient(t, transport.Config{MaxBodySize: tc.max}).NewRequest(t.Context(), "POST", "entity/files", nil, body)
		if (err == nil) != tc.ok || err != nil && !strings.HasPrefix(err.Error(), "moysklad: request body is 114 bytes") {
			t.Errorf("max %d: err = %v", tc.max, err)
		}
	}
	// «Ёж» — 2 символа, 4 байта UTF-8: {"name":"Ёж"} = 15 байт.
	name := struct {
		Name string `json:"name"`
	}{"Ёж"}
	if _, err := limitsClient(t, transport.Config{MaxBodySize: 14}).NewRequest(t.Context(), "POST", "x", nil, name); err == nil {
		t.Error("UTF-8 counted as characters")
	}
	big := file{Content: strings.Repeat("a", transport.DefaultMaxBodySize-14)}
	if _, err := limitsClient(t, transport.Config{}).NewRequest(t.Context(), "POST", "x", nil, big); err != nil {
		t.Fatalf("default max rejected: %v", err)
	}
	big.Content += "a"
	if _, err := limitsClient(t, transport.Config{}).NewRequest(t.Context(), "POST", "x", nil, big); err == nil {
		t.Fatal("default max+1 accepted")
	}
}

// Request line, Host и все значения заголовков: у простого GET 138 байт.
func TestHeaderSizeLimit(t *testing.T) {
	for _, tc := range []struct {
		max int
		ok  bool
	}{{138, true}, {137, false}, {-1, true}} {
		_, err := limitsClient(t, transport.Config{MaxHeaderSize: tc.max}).NewRequest(t.Context(), "GET", "entity/uom", nil, nil)
		if (err == nil) != tc.ok {
			t.Errorf("max %d: err = %v", tc.max, err)
		}
	}

	// Default 8000: у GET на /entity/<n букв> метаданные занимают 135+n байт.
	for _, tc := range []struct {
		n  int
		ok bool
	}{{7865, true}, {7866, false}} {
		_, err := limitsClient(t, transport.Config{}).NewRequest(t.Context(), "GET", "entity/"+strings.Repeat("a", tc.n), nil, nil)
		if (err == nil) != tc.ok {
			t.Errorf("default, n=%d: err = %v", tc.n, err)
		}
	}
	// POST без тела: net/http всё равно отправит Content-Length: 0 (+19 байт).
	if _, err := limitsClient(t, transport.Config{MaxHeaderSize: 156}).NewRequest(t.Context(), "POST", "entity/uo", nil, nil); err == nil {
		t.Error("POST Content-Length not counted")
	}

	// Authorization и повторные значения входят в лимит; секрет не попадает в ошибку.
	auth, err := transport.BearerToken(strings.Repeat("s", 7900))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := limitsClient(t, transport.Config{Auth: auth}).NewRequest(t.Context(), "GET", "entity/uom", nil, nil); err == nil ||
		strings.Contains(err.Error(), "sss") {
		t.Fatalf("long Authorization: err = %v", err)
	}
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client := limitsClient(t, transport.Config{BaseURL: server.URL})
	var opts []transport.RequestOption
	for range 100 {
		opts = append(opts, transport.WithWebHookDisableByPrefix("https://"+strings.Repeat("h", 70)+".example"))
	}
	if _, err := client.Send(t.Context(), "GET", "entity/uom", nil, nil, nil, opts...); err == nil || calls.Load() != 0 {
		t.Fatalf("repeated headers: err = %v, calls = %d", err, calls.Load())
	}
}

func TestBodyDepthLimit(t *testing.T) {
	objects := func(depth int) string {
		return strings.Repeat(`{"a":`, depth-1) + `{}` + strings.Repeat(`}`, depth-1)
	}
	for _, tc := range []struct {
		body string
		max  int
		ok   bool
	}{
		{objects(10), 0, true},
		{objects(11), 0, false},
		// Массивы уровнем не считаются: строка документа obj→positions[]→obj→assortment→meta — 4.
		{`{"a":[[[[[[[[[[[` + objects(9) + `]]]]]]]]]]]}`, 0, true},
		{objects(11), -1, true},
		{objects(3), 2, false},
	} {
		_, err := limitsClient(t, transport.Config{MaxDepth: tc.max}).NewRequest(t.Context(), "POST", "x", nil, rawJSON(tc.body))
		if (err == nil) != tc.ok {
			t.Errorf("%s max %d: err = %v", tc.body, tc.max, err)
		}
	}
}

// Успешный ответ больше лимита запроса читается целиком.
func TestNoResponseSizeLimit(t *testing.T) {
	payload := `{"content":"` + strings.Repeat("a", transport.DefaultMaxBodySize+1) + `"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(payload)) }))
	defer server.Close()
	var dst struct {
		Content string `json:"content"`
	}
	if _, err := limitsClient(t, transport.Config{BaseURL: server.URL}).Send(t.Context(), "GET", "x", nil, nil, &dst); err != nil ||
		len(dst.Content) != transport.DefaultMaxBodySize+1 {
		t.Fatalf("len = %d, err = %v", len(dst.Content), err)
	}
}

type rawJSON string

func (r rawJSON) MarshalJSON() ([]byte, error) { return []byte(r), nil }
