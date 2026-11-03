package transport

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
)

// SendNoRedirect возвращает 3xx как ответ: скрытого запроса по Location нет, в том числе на
// чужой origin; Send того же ответа следует 303 скрытым GET.
func TestSendNoRedirect(t *testing.T) {
	var foreignCalls atomic.Int32
	foreign := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		foreignCalls.Add(1)
	}))
	defer foreign.Close()

	for _, tc := range []struct {
		name     string
		status   int
		location string
		gzip     bool // заголовок gzip при пустом теле: тело успеха не декодируется
	}{
		{"see other", http.StatusSeeOther, "/download/file.pdf?sig=1", false},
		{"see other gzip empty", http.StatusSeeOther, "/download/file.pdf", true},
		{"see other foreign", http.StatusSeeOther, foreign.URL + "/tmp/file.pdf", false},
		{"accepted", http.StatusAccepted, "/status/1", false},
		{"accepted gzip empty", http.StatusAccepted, "/status/1", true},
		{"temporary", http.StatusTemporaryRedirect, "/moved", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			var gotBody string
			server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if calls.Add(1) > 1 {
					t.Errorf("hidden request %s %s", r.Method, r.URL)
					return
				}
				data, _ := io.ReadAll(r.Body)
				gotBody = string(data)
				if r.Method != http.MethodPost || r.URL.Path != "/entity/demand/x/export" || r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("request = %s %s %v", r.Method, r.URL, r.Header)
				}
				w.Header().Set("Location", tc.location)
				w.Header().Set("Content-Type", "application/json")
				if tc.gzip {
					w.Header().Set("Content-Encoding", "gzip")
				}
				w.WriteHeader(tc.status)
			}))
			defer server.Close()
			client, err := NewClient(Config{BaseURL: server.URL, Retry: &RetryPolicy{}})
			if err != nil {
				t.Fatal(err)
			}
			var info ResponseInfo
			resp, err := client.SendNoRedirect(t.Context(), http.MethodPost, "entity/demand/x/export", map[string]string{"extension": "pdf"}, WithResponseInfo(&info))
			if err != nil {
				t.Fatal(err)
			}
			if calls.Load() != 1 || foreignCalls.Load() != 0 || gotBody != `{"extension":"pdf"}` {
				t.Errorf("calls = %d, foreign = %d, body = %s", calls.Load(), foreignCalls.Load(), gotBody)
			}
			if resp.StatusCode != tc.status || resp.Header.Get("Location") != tc.location ||
				info.StatusCode != tc.status || info.Header.Get("Location") != tc.location || info.Attempts != 1 {
				t.Errorf("resp = %d %q, info = %+v", resp.StatusCode, resp.Header.Get("Location"), info)
			}
		})
	}

	// Обычный Send следует 303 скрытым GET: SendNoRedirect нужен именно поэтому.
	var methods []string
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods = append(methods, r.Method)
		if r.Method == http.MethodPost {
			w.Header().Set("Location", "/file")
			w.WriteHeader(http.StatusSeeOther)
		}
	}))
	defer server.Close()
	client, err := NewClient(Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Send(t.Context(), http.MethodPost, "x", nil, map[string]int{}, nil); err != nil || len(methods) != 2 || methods[1] != http.MethodGet {
		t.Errorf("Send: methods = %v, err = %v", methods, err)
	}
	// Следующий вызов того же клиента снова следует redirect: режим не сохраняется в Client.
	methods = nil
	if _, err := client.SendNoRedirect(t.Context(), http.MethodPost, "x", map[string]int{}); err != nil || len(methods) != 1 {
		t.Errorf("SendNoRedirect: methods = %v, err = %v", methods, err)
	}
	methods = nil
	if _, err := client.Send(t.Context(), http.MethodPost, "x", nil, map[string]int{}, nil); err != nil || len(methods) != 2 {
		t.Errorf("Send after SendNoRedirect: methods = %v, err = %v", methods, err)
	}
}

func TestSendNoRedirectErrors(t *testing.T) {
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"errors":[{"error":"bad template","code":3000}]}`)
	}))
	defer server.Close()
	client, err := NewClient(Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.SendNoRedirect(t.Context(), http.MethodPost, "x", nil)
	var apiErr *sdkerrors.APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusBadRequest || len(apiErr.Errors) != 1 || apiErr.Errors[0].Code != 3000 {
		t.Fatalf("err = %v", err)
	}
	if _, err := client.SendNoRedirect(nil, http.MethodPost, "x", nil); err == nil {
		t.Fatal("nil context accepted")
	}

	// Неразбираемый Location net/http отвергает без ответа; подпись в ошибку не попадает.
	bad := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "http://[::1/f?sig=secret")
		w.WriteHeader(http.StatusSeeOther)
	}))
	defer bad.Close()
	client, err = NewClient(Config{BaseURL: bad.URL})
	if err != nil {
		t.Fatal(err)
	}
	var info ResponseInfo
	resp, err := client.SendNoRedirect(t.Context(), http.MethodPost, "x", nil, WithResponseInfo(&info))
	if err == nil || resp != nil || strings.Contains(err.Error(), "secret") || info.StatusCode != 0 || info.Attempts != 1 {
		t.Fatalf("resp = %v, err = %v, info = %+v", resp, err, info)
	}
}
