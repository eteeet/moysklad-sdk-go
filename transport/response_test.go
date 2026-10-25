package transport_test

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestDo(t *testing.T) {
	body := `{"name":"Товар"}`
	damaged := []byte(gzipString(t, body))
	damaged[len(damaged)-8] ^= 0xff

	tests := []struct {
		name     string
		status   int
		encoding string
		body     string
		method   string
		discard  bool
		wantName string
		wantErr  bool
	}{
		{name: "JSON", status: 200, body: body, wantName: "Товар"},
		{name: "gzip", status: 200, encoding: "gzip", body: gzipString(t, body), wantName: "Товар"},
		{name: "identity", status: 200, encoding: "identity", body: body, wantName: "Товар"},
		{name: "no content", status: 204, encoding: "gzip"},
		{name: "reset content", status: 205},
		{name: "HEAD", status: 200, encoding: "gzip", method: http.MethodHead},
		{name: "discard", status: 200, body: "ignored", discard: true},
		{name: "discard gzip", status: 200, encoding: "gzip", body: gzipString(t, body), discard: true},
		{name: "empty JSON", status: 200, wantErr: true},
		{name: "invalid JSON", status: 200, body: `{"name":`, wantErr: true},
		{name: "trailing JSON", status: 200, body: body + " {}", wantErr: true},
		{name: "invalid gzip", status: 200, encoding: "gzip", body: "invalid", wantErr: true},
		{name: "gzip checksum", status: 200, encoding: "gzip", body: string(damaged), wantErr: true},
		{name: "unsupported encoding", status: 200, encoding: "br", body: body, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &trackedBody{Reader: strings.NewReader(tt.body)}
			client := responseClient(t, tt.status, tt.encoding, body)

			method := tt.method
			if method == "" {
				method = http.MethodGet
			}
			req := responseRequest(t, method)

			var got struct {
				Name string `json:"name"`
			}
			var dst any = &got
			if tt.discard {
				dst = nil
			}

			resp, err := client.Do(req, dst)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, want error = %t", err, tt.wantErr)
			}
			if resp == nil || resp.StatusCode != tt.status || resp.Header.Get("X-Test") != "kept" {
				t.Fatalf("response metadata lost: %v", resp)
			}
			if !body.closed {
				t.Fatal("response body not closed")
			}
			if !tt.wantErr && got.Name != tt.wantName {
				t.Fatalf("name = %q, want %q", got.Name, tt.wantName)
			}
		})
	}
}

func TestDoAPIError(t *testing.T) {
	body := `{"errors":[{"error":"bad query","code":1000}]}`
	tests := []struct {
		name        string
		status      int
		encoding    string
		body        string
		wantDetails bool
	}{
		{"JSON", 400, "", body, true},
		{"gzip", 429, "gzip", gzipString(t, body), true},
		{"empty", 415, "", "", false},
		{"HTML", 502, "", "<html>bad gateway</html>", false},
		{"invalid JSON", 400, "", `{"errors":[{"error":"partial"}],`, false},
		{"null", 500, "", "null", false},
		{"invalid gzip", 502, "gzip", "invalid", false},
		{"unsupported encoding", 502, "br", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &trackedBody{Reader: strings.NewReader(tt.body)}
			client := responseClient(t, tt.status, tt.encoding, body)

			resp, err := client.Do(responseRequest(t, http.MethodGet), nil)
			var apiErr *sdkerrors.APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("error = %v, want APIError", err)
			}
			if apiErr.StatusCode != tt.status {
				t.Fatalf("status = %d, want %d", apiErr.StatusCode, tt.status)
			}
			if tt.wantDetails {
				if len(apiErr.Errors) != 1 || apiErr.Errors[0].Code != 1000 ||
					apiErr.Errors[0].Message != "bad query" {
					t.Fatalf("details = %+v", apiErr.Errors)
				}
			} else if len(apiErr.Errors) != 0 {
				t.Fatalf("unexpected details = %+v", apiErr.Errors)
			}
			if resp == nil || resp.StatusCode != tt.status || !body.closed {
				t.Fatal("response missing or body not closed")
			}
		})
	}
}

func TestDoErrorBodyLimit(t *testing.T) {
	payload := `{"errors":[{"error":"too big"}],"padding":"` +
		strings.Repeat("x", 2<<20) + `"}`
	body := &trackedBody{Reader: strings.NewReader(payload)}
	client := responseClient(t, http.StatusBadGateway, "", body)

	_, err := client.Do(responseRequest(t, http.MethodGet), nil)
	var apiErr *sdkerrors.APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusBadGateway {
		t.Fatalf("error = %v, want APIError", err)
	}
	if !strings.Contains(err.Error(), "exceeds") || len(apiErr.Errors) != 0 || !body.closed {
		t.Fatalf("error body limit not enforced: %v", err)
	}
}

func TestDoReadError(t *testing.T) {
	wantErr := errors.New("read failed")
	for _, status := range []int{http.StatusOK, http.StatusBadGateway} {
		body := &trackedBody{Reader: failingReader{err: wantErr}}
		client := responseClient(t, status, "", body)
		var dst map[string]any

		_, err := client.Do(responseRequest(t, http.MethodGet), &dst)
		if !errors.Is(err, wantErr) {
			t.Fatalf("status %d: error = %v, want wrapped read error", status, err)
		}
		var apiErr *sdkerrors.APIError
		if status >= 300 && (!errors.As(err, &apiErr) || apiErr.StatusCode != status) {
			t.Fatalf("HTTP status lost: %v", err)
		}
		if !body.closed {
			t.Fatal("response body not closed")
		}
	}
}

func TestDoTransportError(t *testing.T) {
	wantErr := errors.New("connection failed")
	client, err := transport.NewClient(transport.Config{
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, wantErr
			}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Do(responseRequest(t, http.MethodGet), nil); !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapped transport error", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type trackedBody struct {
	io.Reader
	closed bool
}

func (b *trackedBody) Close() error {
	b.closed = true
	return nil
}

type failingReader struct {
	err error
}

func (r failingReader) Read([]byte) (int, error) {
	return 0, r.err
}

func responseClient(t *testing.T, status int, encoding string, body io.ReadCloser) *transport.Client {
	t.Helper()
	client, err := transport.NewClient(transport.Config{
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: status,
					Header: http.Header{
						"Content-Encoding": {encoding},
						"X-Test":           {"kept"},
					},
					Body:    body,
					Request: req,
				}, nil
			}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func responseRequest(t *testing.T, method string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, transport.DefaultBaseURL+"/entity/product", nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func gzipString(t *testing.T, data string) string {
	t.Helper()
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := io.WriteString(w, data); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}
