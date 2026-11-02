package transport_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func downloadClient(t *testing.T, base string, cfg transport.Config) *transport.Client {
	t.Helper()
	auth, err := transport.BearerToken("secret-token")
	if err != nil {
		t.Fatal(err)
	}
	cfg.BaseURL, cfg.Auth = base+"/api/remap/1.2", auth
	c, err := transport.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Тело копируется целиком без лимита 20 МБ, файл не принимается за JSON, ResponseInfo заполнен.
func TestDownload(t *testing.T) {
	payload := bytes.Repeat([]byte{0, 1, 2, 0xff}, 21_000_000/4)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/remap/1.2/download/abc" || r.URL.RawQuery != "miniature=true" ||
			r.Header.Get("Authorization") != "Bearer secret-token" || r.Header.Get("Accept") != "*/*" {
			t.Errorf("request = %s %s %v", r.Method, r.URL, r.Header)
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(payload)))
		w.Write(payload)
	}))
	defer server.Close()
	c := downloadClient(t, server.URL, transport.Config{})
	for _, href := range []string{"download/abc?miniature=true", server.URL + "/api/remap/1.2/download/abc?miniature=true"} {
		var out bytes.Buffer
		info := transport.ResponseInfo{StatusCode: 503, Attempts: 3}
		n, err := c.Download(t.Context(), href, &out, transport.WithResponseInfo(&info))
		if err != nil || n != int64(len(payload)) || !bytes.Equal(out.Bytes(), payload) {
			t.Fatalf("%s: n = %d, err = %v", href, n, err)
		}
		if info.StatusCode != 200 || info.Attempts != 1 || info.Header.Get("Content-Type") != "image/png" {
			t.Fatalf("info = %+v", info)
		}
	}
}

// Цепочка redirect: внутри API Authorization сохраняется, на подписанном URL другого origin —
// удаляется; ошибка хранилища — APIError; чужой origin в href — отказ до отправки; подпись не
// попадает в текст ошибки.
func TestDownloadRedirectAndOrigin(t *testing.T) {
	var storageCalls atomic.Int32
	storage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		storageCalls.Add(1)
		if r.Header.Get("Authorization") != "" || r.URL.Query().Get("signature") != "s" {
			t.Errorf("storage request = %s %v", r.URL, r.Header)
		}
		if r.URL.Path == "/files/expired" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("X-Storage", "1")
		w.Write([]byte("file"))
	}))
	defer storage.Close()
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret-token" {
			t.Errorf("api authorization = %q", r.Header.Get("Authorization"))
		}
		switch r.URL.Path {
		case "/api/remap/1.2/download/1":
			http.Redirect(w, r, "/api/remap/1.2/download/2", http.StatusFound)
		case "/api/remap/1.2/download/2":
			http.Redirect(w, r, storage.URL+"/files/1?signature=s", http.StatusFound)
		case "/api/remap/1.2/download/expired":
			http.Redirect(w, r, storage.URL+"/files/expired?signature=s", http.StatusFound)
		}
	}))
	defer api.Close()
	c := downloadClient(t, api.URL, transport.Config{})
	var out bytes.Buffer
	var info transport.ResponseInfo
	if n, err := c.Download(t.Context(), "download/1", &out, transport.WithResponseInfo(&info)); err != nil || n != 4 || out.String() != "file" ||
		storageCalls.Load() != 1 || info.Header.Get("X-Storage") != "1" {
		t.Fatalf("n = %d, err = %v, out = %q, info = %+v", n, err, out.String(), info)
	}
	if _, err := c.Download(t.Context(), "download/expired", &out); sdkerrors.StatusCode(err) != http.StatusForbidden {
		t.Fatalf("expired: err = %v", err)
	}
	if _, err := c.Download(t.Context(), storage.URL+"/files/1", &out); err == nil || storageCalls.Load() != 2 {
		t.Fatalf("foreign href: err = %v, calls = %d", err, storageCalls.Load())
	}
	closed := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	closedURL := closed.URL
	closed.Close()
	leak := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, closedURL+"/files/1?signature=top-secret", http.StatusFound)
	}))
	defer leak.Close()
	_, err := downloadClient(t, leak.URL, transport.Config{}).Download(t.Context(), "download/1", &out)
	if err == nil || strings.Contains(err.Error(), "top-secret") {
		t.Fatalf("err = %v", err)
	}
}

// Неразбираемый Location: подписи нет и во вложенной ошибке, даже после пробела или кавычки в query;
// url.Error доступна. Отказ CheckRedirect: ResponseInfo хранит полученный 302, повтора нет.
func TestDownloadRedirectErrors(t *testing.T) {
	malformed := map[string]string{
		"/api/remap/1.2/download/plain": "https://storage.invalid/%zz?signature=top-secret",
		"/api/remap/1.2/download/space": "https://storage.invalid/%zz?pad=x signature=top-secret",
		"/api/remap/1.2/download/quote": `https://storage.invalid/%zz?pad="&signature=top-secret`,
	}
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("X-Source", "api")
		location, ok := malformed[r.URL.Path]
		if !ok {
			location = "/api/remap/1.2/download/2?signature=top-secret"
		}
		w.Header().Set("Location", location)
		w.WriteHeader(http.StatusFound)
	}))
	defer server.Close()
	rejected := errors.New("redirect rejected")
	c := downloadClient(t, server.URL, transport.Config{
		Retry: &transport.RetryPolicy{BaseDelay: time.Millisecond},
		HTTPClient: &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
			return rejected
		}},
	})
	for path := range malformed {
		info := transport.ResponseInfo{StatusCode: 200}
		_, err := c.Download(t.Context(), strings.TrimPrefix(path, "/api/remap/1.2/"), io.Discard, transport.WithResponseInfo(&info))
		var urlErr *url.Error
		if err == nil || strings.Contains(err.Error(), "top-secret") || !errors.As(err, &urlErr) ||
			info.StatusCode != 0 || info.Attempts != 1 {
			t.Fatalf("%s: err = %v, info = %+v", path, err, info)
		}
	}
	calls.Store(0)
	var info transport.ResponseInfo
	_, err := c.Download(t.Context(), "download/1", io.Discard, transport.WithResponseInfo(&info))
	if !errors.Is(err, rejected) || strings.Contains(err.Error(), "top-secret") || calls.Load() != 1 ||
		info.StatusCode != http.StatusFound || info.Attempts != 1 || info.Header.Get("X-Source") != "api" {
		t.Fatalf("rejected: err = %v, calls = %d, info = %+v", err, calls.Load(), info)
	}
}

type timeoutError struct{}

func (timeoutError) Error() string {
	return "timeout reading https://storage.invalid/f?signature=top-secret"
}
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

type timeoutTransport struct{ calls *atomic.Int32 }

func (rt timeoutTransport) RoundTrip(*http.Request) (*http.Response, error) {
	rt.calls.Add(1)
	return nil, timeoutError{}
}

// Скрытие подписи не меняет классификацию: таймаут со ссылкой в тексте повторяется по Config.Retry.
func TestDownloadRedactedTimeoutRetries(t *testing.T) {
	var calls atomic.Int32
	c := downloadClient(t, "https://api.invalid", transport.Config{
		Retry:      &transport.RetryPolicy{MaxAttempts: 2, BaseDelay: time.Millisecond},
		HTTPClient: &http.Client{Transport: timeoutTransport{&calls}},
	})
	var info transport.ResponseInfo
	_, err := c.Download(t.Context(), "download/1", io.Discard, transport.WithResponseInfo(&info))
	var netErr net.Error
	if !errors.As(err, &netErr) || !netErr.Timeout() || !errors.Is(err, timeoutError{}) ||
		strings.Contains(err.Error(), "top-secret") || calls.Load() != 2 || info.Attempts != 2 {
		t.Fatalf("err = %v, calls = %d, info = %+v", err, calls.Load(), info)
	}
}

// Ошибка API не пишет в w; до записи Config.Retry повторяет попытку, после частичной записи — нет.
func TestDownloadErrors(t *testing.T) {
	var calls atomic.Int32
	var partial atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if partial.Load() {
			w.Header().Set("Content-Length", "10")
			w.Write([]byte("abc"))
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"errors":[{"error":"busy","code":1074}]}`))
	}))
	defer server.Close()
	c := downloadClient(t, server.URL, transport.Config{Retry: &transport.RetryPolicy{BaseDelay: time.Millisecond}})
	var out bytes.Buffer
	var info transport.ResponseInfo
	n, err := c.Download(t.Context(), "download/1", &out, transport.WithResponseInfo(&info))
	if n != 0 || out.Len() != 0 || sdkerrors.StatusCode(err) != 503 || !sdkerrors.HasCode(err, 1074) || calls.Load() != 3 ||
		info.StatusCode != 503 || info.Attempts != 3 || strings.Count(err.Error(), "moysklad:") != 1 {
		t.Fatalf("n = %d, err = %v, calls = %d, info = %+v", n, err, calls.Load(), info)
	}
	calls.Store(0)
	if _, err := c.Download(t.Context(), "download/1", &out, transport.WithoutRetry()); sdkerrors.StatusCode(err) != 503 || calls.Load() != 1 {
		t.Fatalf("without retry: err = %v, calls = %d", err, calls.Load())
	}
	partial.Store(true)
	calls.Store(0)
	out.Reset()
	n, err = c.Download(t.Context(), "download/1", &out, transport.WithResponseInfo(&info))
	if !errors.Is(err, io.ErrUnexpectedEOF) || n != 3 || calls.Load() != 1 || info.Attempts != 1 {
		t.Fatalf("partial: n = %d, err = %v, calls = %d", n, err, calls.Load())
	}
	if _, err := c.Download(t.Context(), "download/1", nil); err == nil {
		t.Fatal("nil writer accepted")
	}
}

// gzip-ответ распаковывается независимо от настроек HTTPClient, в том числе тело ошибки.
func TestDownloadGzip(t *testing.T) {
	gz := func(data string) []byte {
		var b bytes.Buffer
		zw := gzip.NewWriter(&b)
		zw.Write([]byte(data))
		zw.Close()
		return b.Bytes()
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept-Encoding") != "gzip" {
			t.Errorf("accept-encoding = %q", r.Header.Get("Accept-Encoding"))
		}
		w.Header().Set("Content-Encoding", "gzip")
		if r.URL.Path == "/api/remap/1.2/download/missing" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(gz(`{"errors":[{"error":"missing","code":1001}]}`))
			return
		}
		w.Write(gz("unpacked"))
	}))
	defer server.Close()
	c := downloadClient(t, server.URL, transport.Config{HTTPClient: &http.Client{Transport: &http.Transport{DisableCompression: true}}})
	var out bytes.Buffer
	if n, err := c.Download(t.Context(), "download/1", &out); err != nil || n != 8 || out.String() != "unpacked" {
		t.Fatalf("n = %d, err = %v, out = %q", n, err, out.String())
	}
	if _, err := c.Download(t.Context(), "download/missing", &out); !sdkerrors.HasCode(err, 1001) {
		t.Fatalf("err = %v", err)
	}
}

// Ошибка writer прерывает загрузку и возвращается с числом записанных байт.
func TestDownloadWriterError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("data")) }))
	defer server.Close()
	c := downloadClient(t, server.URL, transport.Config{})
	sentinel := errors.New("disk full")
	if n, err := c.Download(t.Context(), "download/1", failWriter{sentinel}); !errors.Is(err, sentinel) || n != 0 {
		t.Fatalf("n = %d, err = %v", n, err)
	}
}

type failWriter struct{ err error }

func (f failWriter) Write([]byte) (int, error) { return 0, f.err }

// Обрыв тела и отмена контекста возвращают ошибку и число уже записанных байт.
func TestDownloadShortReadAndCancel(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "10")
		w.Write([]byte("abcd"))
		w.(http.Flusher).Flush()
		if r.URL.Path == "/api/remap/1.2/download/slow" {
			<-release
		}
		// Соединение закрывается раньше обещанных 10 байт.
	}))
	defer server.Close()
	defer close(release)
	c := downloadClient(t, server.URL, transport.Config{})
	var out bytes.Buffer
	n, err := c.Download(t.Context(), "download/short", &out)
	if err == nil || !errors.Is(err, io.ErrUnexpectedEOF) || n != 4 || out.String() != "abcd" {
		t.Fatalf("short: n = %d, err = %v", n, err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	written := &signalWriter{first: make(chan struct{})}
	go func() {
		<-written.first
		cancel()
	}()
	n, err = c.Download(ctx, "download/slow", written)
	if err == nil || !errors.Is(err, context.Canceled) || n != 4 {
		t.Fatalf("cancel: n = %d, err = %v", n, err)
	}
}

// Слоты Limiter освобождаются после успеха и после ошибки: следующая загрузка не ждёт.
func TestDownloadReleasesLimiter(t *testing.T) {
	var fail atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fail.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte("ok"))
	}))
	defer server.Close()
	limiter, err := transport.NewLimiter(transport.LimiterConfig{Profile: transport.ProfileEmployee, Concurrency: 1})
	if err != nil {
		t.Fatal(err)
	}
	c := downloadClient(t, server.URL, transport.Config{Limiter: limiter})
	for i, f := range []bool{false, true, false} {
		fail.Store(f)
		ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
		_, err := c.Download(ctx, "download/1", io.Discard)
		cancel()
		if f != (err != nil) || errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("download %d: err = %v", i, err)
		}
	}
	// Отмена до ответа тоже освобождает слот.
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := c.Download(ctx, "download/1", io.Discard); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled: err = %v", err)
	}
	ctx, cancel = context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	if _, err := c.Download(ctx, "download/1", io.Discard); err != nil {
		t.Fatalf("after cancel: err = %v", err)
	}
}

// signalWriter закрывает first после первой записи: тест отменяет контекст посреди тела.
type signalWriter struct {
	first chan struct{}
	once  sync.Once
}

func (s *signalWriter) Write(p []byte) (int, error) {
	s.once.Do(func() { close(s.first) })
	return len(p), nil
}
