package images_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/images"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	ownerID = "7944ef04-f831-11e5-7a69-971500188b19"
	imageID = "f2728180-6afd-4d37-8a13-f3b48069bbb6"
	path    = "/entity/product/" + ownerID + "/images"
)

func fixture(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile("../testdata/upstream/product_image.json")
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func service(t *testing.T, handler http.HandlerFunc) *images.Service {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, _ := images.NewService(c)
	return s
}

func TestImageFixture(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		var got images.Image
		if err := codec.Unmarshal(fixture(t), &got); err != nil {
			t.Fatal(err)
		}
		if got.Title != "birdimage" || got.Filename != "birdimage.png" || got.Size != 14052 || got.Updated.IsZero() ||
			got.Meta.DownloadHref == "" || got.Miniature == nil || got.Miniature.DownloadHref == "" || got.Tiny == nil || got.Tiny.MediaType != "image/png" {
			t.Fatalf("%s: image = %+v", codec.Name, got)
		}
		// Элементы Replace: существующее — только meta, новое — filename и content (Base64).
		data, err := codec.Marshal([]images.Item{images.Keep(got), images.New(images.Upload{Filename: "a.png", Content: []byte{1, 2}})})
		if err != nil {
			t.Fatal(err)
		}
		documenttest.JSONEqual(t, data, []byte(`[{"meta":{"href":"`+got.Meta.Href+`","type":"image","mediaType":"application/json","downloadHref":"`+got.Meta.DownloadHref+`"}},`+
			`{"filename":"a.png","content":"AQI="}]`))
	}
}

func TestOperations(t *testing.T) {
	image := string(fixture(t))
	var got []string
	s := service(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = append(got, r.Method+" "+r.URL.RequestURI()+" "+string(body))
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/download/"):
			io.WriteString(w, "PNG")
		case r.Method == http.MethodGet:
			io.WriteString(w, `{"meta":{"href":"x","size":1},"rows":[`+image+`]}`)
		case strings.HasSuffix(r.URL.Path, "/delete"):
			io.WriteString(w, `[{"info":"ok"}]`)
		case r.Method == http.MethodPost:
			io.WriteString(w, `[`+image+`]`)
		}
	})
	ctx := t.Context()
	l, err := s.List(ctx, "product", ownerID, query.PageOptions{Fields: images.PermanentHref})
	if err != nil || len(l.Rows) != 1 || l.Rows[0].Size != 14052 {
		t.Fatalf("list = %+v, err = %v", l, err)
	}
	if r, err := s.Add(ctx, "product", ownerID, images.Upload{Filename: "birdimageNew.png", Content: []byte("x")}); err != nil || len(r) != 1 {
		t.Fatalf("add = %+v, err = %v", r, err)
	}
	if r, err := s.Replace(ctx, "product", ownerID, []images.Item{images.Keep(l.Rows[0])}); err != nil || len(r) != 1 {
		t.Fatalf("replace = %+v, err = %v", r, err)
	}
	if _, err := s.Replace(ctx, "variant", ownerID, nil); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(ctx, "bundle", ownerID, imageID); err != nil {
		t.Fatal(err)
	}
	if r, err := s.DeleteMany(ctx, "contentcard", ownerID, []string{imageID}); err != nil || r[0].Info != "ok" {
		t.Fatalf("delete many = %+v, err = %v", r, err)
	}
	var buf bytes.Buffer
	l.Rows[0].Meta.DownloadHref = "/download/" + imageID
	if n, err := s.Download(ctx, l.Rows[0], &buf); err != nil || n != 3 || buf.String() != "PNG" {
		t.Fatalf("download = %d %q, err = %v", n, buf.String(), err)
	}
	keep := `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/` + ownerID + `/images/` + imageID + `","type":"image","mediaType":"application/json","downloadHref":"https://api.moysklad.ru/api/remap/1.2/download/` + imageID + `"}}`
	want := []string{
		"GET " + path + "?fields=downloadPermanentHref&limit=1000&offset=0 ",
		"POST " + path + ` {"filename":"birdimageNew.png","content":"eA=="}`,
		"POST " + path + ` [` + keep + `]`,
		"POST /entity/variant/" + ownerID + "/images []",
		"DELETE /entity/bundle/" + ownerID + "/images/" + imageID + " ",
		"", // массовое удаление: href зависит от адреса тестового сервера, проверяется ниже
		"GET /download/" + imageID + " ",
	}
	if len(got) != len(want) {
		t.Fatalf("requests:\n%s", strings.Join(got, "\n"))
	}
	for i := range want {
		if i == 5 {
			if !strings.HasPrefix(got[i], "POST /entity/contentcard/"+ownerID+"/images/delete [{\"meta\":{\"href\":\"http") ||
				!strings.Contains(got[i], "/entity/contentcard/"+ownerID+"/images/"+imageID+`","type":"image"`) {
				t.Errorf("request %d = %s", i, got[i])
			}
			continue
		}
		if got[i] != want[i] {
			t.Errorf("request %d = %s\nwant %s", i, got[i], want[i])
		}
	}
}

func TestHTTPErrors(t *testing.T) {
	invoke := func(ctx context.Context, c *transport.Client) (*[]images.Image, error) {
		s, _ := images.NewService(c)
		r, err := s.Add(ctx, "product", ownerID, images.Upload{Filename: "a.png", Content: []byte{1}})
		if err != nil {
			return nil, err
		}
		return &r, nil
	}
	for _, tc := range []struct {
		name, body string
		status     int
		ok         bool
	}{
		{"array", `[]`, 200, true},
		{"null", `null`, 200, false},
		{"object", `{}`, 200, false},
		{"api error", `{"errors":[{"error":"too big","code":3006}]}`, 413, false},
		{"no content", ``, 204, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			s := service(t, func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				w.WriteHeader(tc.status)
				io.WriteString(w, tc.body)
			})
			r, err := s.Add(t.Context(), "product", ownerID, images.Upload{Filename: "a.png", Content: []byte{1}})
			if tc.ok {
				// Replace идёт тем же POST и разбирает ответ так же.
				if rr, err := s.Replace(t.Context(), "product", ownerID, []images.Item{images.New(images.Upload{Filename: "b.png", Content: []byte{2}})}); err != nil || rr == nil {
					t.Fatalf("replace = %+v, err = %v", rr, err)
				}
				calls.Add(-1)
			}
			if tc.ok != (err == nil) || calls.Load() != 1 || (tc.ok && r == nil) {
				t.Fatalf("r = %+v, err = %v, calls = %d", r, err, calls.Load())
			}
			if err != nil && strings.Count(err.Error(), "moysklad:") != 1 {
				t.Errorf("error = %v", err)
			}
		})
	}
	documenttest.LocalErrors(t, invoke)
}

func TestValidation(t *testing.T) {
	var calls atomic.Int32
	s := service(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1) })
	ctx := t.Context()
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"list service": func(o ...transport.RequestOption) error {
			_, err := s.List(ctx, "service", ownerID, query.PageOptions{}, o...)
			return err
		},
		"list fields": func(o ...transport.RequestOption) error {
			_, err := s.List(ctx, "product", ownerID, query.PageOptions{Fields: "x"}, o...)
			return err
		},
		"add bad id": func(o ...transport.RequestOption) error {
			_, err := s.Add(ctx, "product", "bad", images.Upload{Filename: "a"}, o...)
			return err
		},
		"add no name": func(o ...transport.RequestOption) error {
			_, err := s.Add(ctx, "product", ownerID, images.Upload{}, o...)
			return err
		},
		"add no content": func(o ...transport.RequestOption) error {
			_, err := s.Add(ctx, "product", ownerID, images.Upload{Filename: "a.png"}, o...)
			return err
		},
		"replace path name": func(o ...transport.RequestOption) error {
			_, err := s.Replace(ctx, "product", ownerID, []images.Item{images.New(images.Upload{Filename: "../a", Content: []byte{1}})}, o...)
			return err
		},
		"replace no content": func(o ...transport.RequestOption) error {
			_, err := s.Replace(ctx, "product", ownerID, []images.Item{images.New(images.Upload{Filename: "a.png"})}, o...)
			return err
		},
		"add path name": func(o ...transport.RequestOption) error {
			_, err := s.Add(ctx, "product", ownerID, images.Upload{Filename: "../a"}, o...)
			return err
		},
		"replace empty": func(o ...transport.RequestOption) error {
			_, err := s.Replace(ctx, "product", ownerID, []images.Item{{}}, o...)
			return err
		},
		"replace both": func(o ...transport.RequestOption) error {
			_, err := s.Replace(ctx, "product", ownerID, []images.Item{{Filename: "a", Meta: &model.Meta{Href: "x"}}}, o...)
			return err
		},
		"replace href": func(o ...transport.RequestOption) error {
			_, err := s.Replace(ctx, "product", ownerID, []images.Item{images.Keep(images.Image{})}, o...)
			return err
		},
		"delete bad": func(o ...transport.RequestOption) error { return s.Delete(ctx, "product", ownerID, "x", o...) },
		"delete many": func(o ...transport.RequestOption) error {
			_, err := s.DeleteMany(ctx, "product", ownerID, nil, o...)
			return err
		},
		"download empty": func(o ...transport.RequestOption) error {
			_, err := s.Download(ctx, images.Image{}, io.Discard, o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d", calls.Load())
	}
}
