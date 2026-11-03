package files_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/files"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const productID = "7944ef04-f831-11e5-7a69-971500188b19"

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func service(c *transport.Client) *files.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := files.NewService(c)
	return s
}

func TestFileJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var page model.Collection[files.File]
			if err := codec.Unmarshal(read(t, "custom/files/list.json"), &page); err != nil {
				t.Fatal(err)
			}
			image, doc := page.Rows[0], page.Rows[1]
			if len(page.Rows) != 2 || image.Filename != "birdimage.png" || image.Size != 14052 || image.Created.IsZero() ||
				image.CreatedBy == nil || image.Miniature == nil || image.Miniature.Href == "" || image.Tiny == nil ||
				!strings.HasSuffix(image.Meta.DownloadHref, "/download/f2728180-6afd-4d37-8a13-f3b48069bbb6") ||
				doc.Miniature != nil || doc.Tiny != nil {
				t.Fatalf("files = %+v", page.Rows)
			}
			var uploaded []files.File
			if err := codec.Unmarshal(read(t, "custom/files/upload.json"), &uploaded); err != nil || len(uploaded) != 2 {
				t.Fatalf("upload = %+v, err = %v", uploaded, err)
			}
			body, err := codec.Marshal([]files.Upload{{Filename: "doc.pdf", Content: []byte("%PDF\x00\xff")}, {Filename: "empty.txt", Content: []byte{}}})
			if err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, body, []byte(`[{"filename":"doc.pdf","content":"JVBERgD/"},{"filename":"empty.txt","content":""}]`))
		})
	}
}

func TestListHTTP(t *testing.T) {
	invoke := func(ctx context.Context, c *transport.Client) (*model.Collection[files.File], error) {
		return service(c).List(ctx, "product", productID, query.PageOptions{Limit: 2, Offset: 4})
	}
	for _, tc := range documenttest.ResponseCases(string(read(t, "custom/files/list.json"))) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "GET", "entity/product/"+productID+"/files?limit=2&offset=4", nil, invoke, func(got *model.Collection[files.File]) {
				if len(got.Rows) != 2 || got.Meta.Size != 2 {
					t.Fatalf("page = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}

func TestUploadHTTP(t *testing.T) {
	uploads := []files.Upload{{Filename: "birdimageNew.png", Content: []byte("png")}, {Filename: "doc.pdf", Content: []byte("pdf")}}
	invoke := func(ctx context.Context, c *transport.Client) (*[]files.File, error) {
		got, err := service(c).Upload(ctx, "product", productID, uploads)
		if err != nil {
			return nil, err
		}
		return &got, nil
	}
	body := []byte(`[{"filename":"birdimageNew.png","content":"cG5n"},{"filename":"doc.pdf","content":"cGRm"}]`)
	for _, tc := range documenttest.ResponseCases(string(read(t, "custom/files/upload.json"))) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "POST", "entity/product/"+productID+"/files", body, invoke, func(got *[]files.File) {
				if len(*got) != 2 || (*got)[1].Filename != "doc.pdf" {
					t.Fatalf("files = %+v", *got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}

// Содержимое читается по meta.downloadHref одним GET; файл без ссылки не запрашивается.
// Пустое содержимое уходит строкой "" обоими codecs, в том числе nil.
func TestUploadNilContent(t *testing.T) {
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		documenttest.JSONEqual(t, data, []byte(`[{"filename":"empty.txt","content":""}]`))
		w.Write([]byte(`[]`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service(c).Upload(t.Context(), "product", productID, []files.Upload{{Filename: "empty.txt"}}); err != nil {
		t.Fatal(err)
	}
}

// Ошибки загрузки несут операцию и единый префикс; APIError, отмена и транспорт — через errors.As/Is.
func TestDownloadErrors(t *testing.T) {
	file := func(base string) files.File {
		return files.File{Meta: model.Meta{DownloadHref: base + "/api/remap/1.2/download/f1"}}
	}
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"error":"missing","code":1001}]}`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = service(c).Download(t.Context(), file(server.URL), io.Discard)
	if sdkerrors.StatusCode(err) != http.StatusNotFound || !strings.Contains(err.Error(), "download file") || strings.Count(err.Error(), "moysklad:") != 1 {
		t.Fatalf("err = %v", err)
	}
	redirect := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "https://storage.invalid/%zz?signature=top-secret")
		w.WriteHeader(http.StatusFound)
	}))
	defer redirect.Close()
	c, err = transport.NewClient(transport.Config{BaseURL: redirect.URL + "/api/remap/1.2"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = service(c).Download(t.Context(), file(redirect.URL), io.Discard)
	if err == nil || strings.Contains(err.Error(), "top-secret") || !strings.Contains(err.Error(), "download file") {
		t.Fatalf("redirect: err = %v", err)
	}
	documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*int64, error) {
		n, err := service(c).Download(ctx, file("https://api.moysklad.ru"), io.Discard)
		if err != nil {
			return nil, err
		}
		return &n, nil
	})
}

func TestDownload(t *testing.T) {
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/remap/1.2/download/f2728180" || r.Header.Get("X-Lognex-Download-Expiration-Seconds") != "300" {
			t.Errorf("request = %s %v", r.URL, r.Header)
		}
		w.Write([]byte("content"))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2"})
	if err != nil {
		t.Fatal(err)
	}
	file := files.File{Meta: model.Meta{DownloadHref: server.URL + "/api/remap/1.2/download/f2728180"}}
	var out bytes.Buffer
	if n, err := service(c).Download(t.Context(), file, &out, transport.WithDownloadExpiration(300*time.Second)); err != nil || n != 7 || out.String() != "content" {
		t.Fatalf("n = %d, err = %v", n, err)
	}
	if _, err := service(c).Download(t.Context(), files.File{}, &out); err == nil {
		t.Fatal("file without downloadHref accepted")
	}
}

// Ошибки ввода и закодированное тело больше MaxBodySize не уходят на сервер; число файлов
// объекта (до 100) SDK заранее не проверяет.
func TestValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Write([]byte(`[]`))
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL, MaxBodySize: 1000})
	if err != nil {
		t.Fatal(err)
	}
	s := service(c)
	one := []files.Upload{{Filename: "a.txt", Content: []byte("a")}}
	if _, err := s.Upload(t.Context(), "product", productID, []files.Upload{{}}); err == nil {
		t.Error("empty filename accepted")
	}
	tenFiles := make([]files.Upload, files.MaxUpload)
	for i := range tenFiles {
		tenFiles[i] = files.Upload{Filename: "f.txt"}
	}
	if _, err := s.Upload(t.Context(), "product", productID, tenFiles); err != nil {
		t.Fatalf("10 files: %v", err)
	}
	for name, err := range map[string]error{
		"no files": func() error { _, err := s.Upload(t.Context(), "product", productID, nil); return err }(),
		"11 files": func() error {
			_, err := s.Upload(t.Context(), "product", productID, make([]files.Upload, 11))
			return err
		}(),
		"path traversal": func() error {
			_, err := s.Upload(t.Context(), "product", productID, []files.Upload{{Filename: "../a.txt"}})
			return err
		}(),
		"backslash": func() error {
			_, err := s.Upload(t.Context(), "product", productID, []files.Upload{{Filename: `a\b.txt`}})
			return err
		}(),
		"dot dot": func() error {
			_, err := s.Upload(t.Context(), "product", productID, []files.Upload{{Filename: ".."}})
			return err
		}(),
		"encoded body": func() error {
			_, err := s.Upload(t.Context(), "product", productID, []files.Upload{{Filename: "big.bin", Content: make([]byte, 800)}})
			return err
		}(),
		"bad type": func() error { _, err := s.Upload(t.Context(), "../product", productID, one); return err }(),
		"bad id":   func() error { _, err := s.List(t.Context(), "product", "bad", query.PageOptions{}); return err }(),
		"fields": func() error {
			_, err := s.List(t.Context(), "product", productID, query.PageOptions{Fields: "x"})
			return err
		}(),
		"limit": func() error {
			_, err := s.List(t.Context(), "product", productID, query.PageOptions{Limit: 1001})
			return err
		}(),
	} {
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("%s: err = %v", name, err)
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("calls = %d", calls.Load())
	}
	if _, err := files.NewService(nil); err == nil {
		t.Error("nil client accepted")
	}
}
