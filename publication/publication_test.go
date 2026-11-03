package publication_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	sdkerrors "github.com/eteeet/moysklad-sdk-go/errors"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/publication"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	api           = "https://api.moysklad.ru/api/remap/1.2/"
	documentID    = "0b71daec-055e-11e6-9464-e4de0000007e"
	publicationID = "aec51463-bbd2-11e6-8a84-bae500000003"
)

func read(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../testdata/custom/publication/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func want(href string) publication.Publication {
	return publication.Publication{
		Meta: model.Meta{Href: api + "entity/demand/" + documentID + "/publication/" + publicationID, Type: "operationpublication", MediaType: "application/json"},
		Template: model.Reference{Meta: model.Meta{Href: api + "entity/demand/metadata/customtemplate/38d1c843-1601-11e7-8af5-581e00000009",
			Type: "customtemplate", MediaType: "application/json"}},
		Href: href,
	}
}

func createRequest() publication.CreateRequest {
	return publication.CreateRequest{Template: model.Reference{Meta: model.Meta{
		Href: api + "entity/demand/metadata/customtemplate/38d1c843-1601-11e7-8af5-581e00000009", Type: "customtemplate", MediaType: "application/json"}}}
}

func TestFixturesAndRequest(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var page model.Collection[publication.Publication]
			if err := codec.Unmarshal(read(t, "list.json"), &page); err != nil {
				t.Fatal(err)
			}
			// meta.type списка в источнике — тип документа (demand), а не publication.
			if len(page.Rows) != 1 || page.Meta.Size != 1 || page.Meta.Type != "demand" || page.Rows[0] != want("https://mskld.ru/73NIpnAbPr") {
				t.Errorf("list = %+v", page)
			}
			for file, href := range map[string]string{
				"create-200.json": "https://mskld.ru/reu92ZrjCM",
				"create-201.json": "https://mskld.ru/f00HzRGx8Q",
				"get.json":        "https://mskld.ru/fuXrdd7Uii",
			} {
				var got publication.Publication
				if err := codec.Unmarshal(read(t, file), &got); err != nil || got != want(href) {
					t.Errorf("%s = %+v, err = %v", file, got, err)
				}
			}
			body, err := codec.Marshal(createRequest())
			if err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, body, read(t, "create.golden.json"))
		})
	}
}

func service(t *testing.T, c *transport.Client) *publication.Service {
	t.Helper()
	s, err := publication.NewService(c)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestHTTP(t *testing.T) {
	base := "entity/demand/" + documentID + "/publication"
	for _, op := range []struct {
		name, method, path string
		body               []byte
		cases              []documenttest.ResponseCase
		want               publication.Publication
		invoke             func(context.Context, *publication.Service) (*publication.Publication, error)
	}{
		{"list", http.MethodGet, base, nil, documenttest.ResponseCases(string(read(t, "list.json"))), want("https://mskld.ru/73NIpnAbPr"),
			func(ctx context.Context, s *publication.Service) (*publication.Publication, error) {
				page, err := s.List(ctx, "demand", documentID)
				if err != nil {
					return nil, err
				}
				return &page.Rows[0], nil
			}},
		{"get", http.MethodGet, base + "/" + publicationID, nil, documenttest.ResponseCases(string(read(t, "get.json"))), want("https://mskld.ru/fuXrdd7Uii"),
			func(ctx context.Context, s *publication.Service) (*publication.Publication, error) {
				return s.Get(ctx, "demand", documentID, publicationID)
			}},
		{"create", http.MethodPost, base, read(t, "create.golden.json"), createCases(t), want("https://mskld.ru/f00HzRGx8Q"),
			func(ctx context.Context, s *publication.Service) (*publication.Publication, error) {
				return s.Create(ctx, "demand", documentID, createRequest())
			}},
	} {
		for _, tc := range op.cases {
			t.Run(op.name+"/"+tc.Name, func(t *testing.T) {
				expected := op.want
				if tc.Name == "existing" {
					expected = want("https://mskld.ru/reu92ZrjCM")
				}
				documenttest.HTTP(t, tc, op.method, op.path, op.body, func(ctx context.Context, c *transport.Client) (*publication.Publication, error) {
					return op.invoke(ctx, service(t, c))
				}, func(got *publication.Publication) {
					if *got != expected {
						t.Errorf("publication = %+v, want %+v", *got, expected)
					}
				})
			})
		}
		t.Run(op.name+"/local", func(t *testing.T) {
			documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*publication.Publication, error) {
				return op.invoke(ctx, service(t, c))
			})
		})
	}
}

// createCases: 201 — новая публикация, 200 — существующая; оба успех, 202 — нет.
func createCases(t *testing.T) []documenttest.ResponseCase {
	created := string(read(t, "create-201.json"))
	var cases []documenttest.ResponseCase
	for _, tc := range documenttest.ResponseCases(created) {
		switch tc.Name {
		case "plain", "gzip":
			tc.Status = http.StatusCreated
		case "unexpected status":
			tc.Status = http.StatusAccepted
		}
		cases = append(cases, tc)
	}
	return append(cases, documenttest.ResponseCase{Name: "existing", Status: http.StatusOK, Body: string(read(t, "create-200.json")), Success: true})
}

func TestDelete(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		ok     bool
		api    bool
	}{
		{"no content", 204, "", true, false},
		{"ok", 200, "", true, false},
		{"not found", 404, `{"errors":[{"error":"missing","code":1001}]}`, false, true},
		{"unexpected", 201, "", false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				body, _ := io.ReadAll(r.Body)
				if r.Method != http.MethodDelete || r.URL.RequestURI() != "/api/remap/1.2/entity/demand/"+documentID+"/publication/"+publicationID ||
					r.Header.Get("Authorization") != "Bearer test-token" || len(body) != 0 || r.Header.Get("Content-Type") != "" {
					t.Errorf("request = %s %s %v %s", r.Method, r.URL.RequestURI(), r.Header, body)
				}
				w.WriteHeader(tc.status)
				io.WriteString(w, tc.body)
			}))
			defer server.Close()
			auth, err := transport.BearerToken("test-token")
			if err != nil {
				t.Fatal(err)
			}
			client, err := transport.NewClient(transport.Config{BaseURL: server.URL + "/api/remap/1.2", Auth: auth})
			if err != nil {
				t.Fatal(err)
			}
			err = service(t, client).Delete(t.Context(), "demand", documentID, publicationID)
			if calls.Load() != 1 {
				t.Fatalf("calls = %d", calls.Load())
			}
			if tc.ok != (err == nil) {
				t.Fatalf("err = %v", err)
			}
			if err != nil && (strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "delete demand") ||
				tc.api != (sdkerrors.StatusCode(err) == tc.status)) {
				t.Errorf("err = %v", err)
			}
		})
	}
	documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*struct{}, error) {
		if err := service(t, c).Delete(ctx, "demand", documentID, publicationID); err != nil {
			return nil, err
		}
		return &struct{}{}, nil
	})
}

// Тип, id и template проверяются до HTTP, ResponseInfo сброшен.
func TestValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(t, client)
	ctx := t.Context()
	for name, invoke := range map[string]func(transport.RequestOption) error{
		"list type": func(o transport.RequestOption) error { _, err := s.List(ctx, "Demand", documentID, o); return err },
		"list id":   func(o transport.RequestOption) error { _, err := s.List(ctx, "demand", "x", o); return err },
		"get id":    func(o transport.RequestOption) error { _, err := s.Get(ctx, "demand", documentID, "", o); return err },
		"create template": func(o transport.RequestOption) error {
			_, err := s.Create(ctx, "demand", documentID, publication.CreateRequest{}, o)
			return err
		},
		"create type": func(o transport.RequestOption) error {
			_, err := s.Create(ctx, "", documentID, createRequest(), o)
			return err
		},
		"delete id": func(o transport.RequestOption) error {
			return s.Delete(ctx, "demand", documentID, publicationID+"/x", o)
		},
		"delete type": func(o transport.RequestOption) error { return s.Delete(ctx, "1demand", documentID, publicationID, o) },
	} {
		info := transport.ResponseInfo{StatusCode: 200, Attempts: 1}
		err := invoke(transport.WithResponseInfo(&info))
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || info.StatusCode != 0 || info.Attempts != 0 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Errorf("calls = %d", calls.Load())
	}
}
