package security_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/security"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const secret = "6f1a7f3c0e2d4b5a"

func newService(t *testing.T, handler http.HandlerFunc) (*security.Service, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	// Клиент настроен на Bearer: CreateToken всё равно должен отправить Basic.
	bearer, err := transport.BearerToken("old-token")
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL, Auth: bearer})
	if err != nil {
		t.Fatal(err)
	}
	service, err := security.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	return service, &calls
}

func TestCreateToken(t *testing.T) {
	for _, compressed := range []bool{false, true} {
		t.Run(fmt.Sprint("gzip=", compressed), func(t *testing.T) { testCreateToken(t, compressed) })
	}
}

func testCreateToken(t *testing.T, compressed bool) {
	service, calls := newService(t, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		want := "Basic " + base64.StdEncoding.EncodeToString([]byte("admin@shop:pass:word"))
		if r.Method != http.MethodPost || r.URL.RequestURI() != "/security/token" || len(body) != 0 ||
			r.Header.Get("Authorization") != want || r.Header.Get("Content-Type") != "" ||
			r.Header.Get("Accept-Encoding") != "gzip" || r.Header.Get("Accept") != "application/json;charset=utf-8" ||
			r.Header.Get("Accept-Language") != "en" {
			t.Errorf("request = %s %s %v body=%q", r.Method, r.URL, r.Header, body)
		}
		payload := []byte(`{"access_token":"` + secret + `"}`)
		if compressed {
			w.Header().Set("Content-Encoding", "gzip")
			zw := gzip.NewWriter(w)
			zw.Write(payload)
			zw.Close()
			return
		}
		w.Write(payload)
	})

	var info transport.ResponseInfo
	bearer, err := transport.BearerToken("other")
	if err != nil {
		t.Fatal(err)
	}
	token, err := service.CreateToken(t.Context(), "admin@shop", "pass:word",
		transport.WithAcceptLanguage("en"), transport.WithResponseInfo(&info),
		transport.WithAuth(bearer)) // caller не может подменить Basic
	if err != nil || token.AccessToken != secret || info.StatusCode != http.StatusOK || calls.Load() != 1 {
		t.Fatalf("token = %v, err = %v, info = %+v", token, err, info)
	}
	var logs bytes.Buffer
	slog.New(slog.NewJSONHandler(&logs, nil)).Info("issued", "token", token, "value", *token)
	for _, s := range []string{fmt.Sprint(token), fmt.Sprintf("%v %+v %#v", *token, token, *token), logs.String()} {
		if strings.Contains(s, secret) {
			t.Fatalf("token leaked: %s", s)
		}
	}
}

func TestCreateTokenErrors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		api    bool
	}{
		{"unauthorized", 401, `{"errors":[{"error":"auth","code":1056}]}`, true},
		{"rate limit", 429, `{"errors":[{"error":"slow","code":1049}]}`, true},
		{"server error", 500, `<html>oops</html>`, true},
		{"unexpected status", 201, `{"access_token":"` + secret + `"}`, false},
		{"empty token", 200, `{"access_token":""}`, false},
		{"missing token", 200, `{}`, false},
		{"malformed", 200, `{"access_token":`, false},
		{"null", 200, `null`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service, calls := newService(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				w.Write([]byte(tc.body))
			})
			token, err := service.CreateToken(t.Context(), "admin@shop", "password")
			var apiErr *sdkerrors.APIError
			if token != nil || err == nil || errors.As(err, &apiErr) != tc.api || calls.Load() != 1 {
				t.Fatalf("token = %v, err = %v, calls = %d", token, err, calls.Load())
			}
			if strings.Count(err.Error(), "moysklad:") != 1 || strings.Contains(err.Error(), secret) {
				t.Fatalf("err = %v", err)
			}
		})
	}
}

// Потерянный ответ и отмена во время запроса: одна отправка, без повтора.
func TestCreateTokenLostResponse(t *testing.T) {
	received := make(chan struct{}, 1)
	service, calls := newService(t, func(w http.ResponseWriter, r *http.Request) {
		received <- struct{}{}
		if r.Header.Get("X-Case") == "cancel" {
			<-r.Context().Done()
			return
		}
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Error(err)
			return
		}
		conn.Close()
	})
	if _, err := service.CreateToken(t.Context(), "admin@shop", "password"); err == nil || calls.Load() != 1 {
		t.Fatalf("err = %v, calls = %d", err, calls.Load())
	}
	<-received

	ctx, cancel := context.WithCancel(t.Context())
	go func() {
		<-received
		cancel()
	}()
	_, err := service.CreateToken(ctx, "admin@shop", "password", transport.WithHeader("X-Case", "cancel"))
	if !errors.Is(err, context.Canceled) || calls.Load() != 2 {
		t.Fatalf("err = %v, calls = %d", err, calls.Load())
	}
}

func TestCreateTokenNoRequest(t *testing.T) {
	service, calls := newService(t, func(w http.ResponseWriter, r *http.Request) {})
	for _, creds := range [][2]string{{"", "password"}, {"admin@shop", ""}, {"a:b", "password"}} {
		if _, err := service.CreateToken(t.Context(), creds[0], creds[1]); err == nil {
			t.Errorf("credentials %q accepted", creds[0])
		}
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := service.CreateToken(ctx, "admin@shop", "password"); !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v", err)
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := security.NewService(nil); err == nil {
		t.Fatal("nil client accepted")
	}
}

// Даже при RetryPolicy клиента выпуск токена отправляется один раз.
func TestCreateTokenNoRetryWithPolicy(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL, Retry: &transport.RetryPolicy{}})
	if err != nil {
		t.Fatal(err)
	}
	service, err := security.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.CreateToken(t.Context(), "admin@shop", "password"); err == nil || calls.Load() != 1 {
		t.Fatalf("err = %v, calls = %d", err, calls.Load())
	}
}
