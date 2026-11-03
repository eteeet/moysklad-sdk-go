package customentity_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/customentity"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	entityID  = "7944ef04-f831-11e5-7a69-971500188b19"
	elementID = "01b29035-2656-11e6-8a84-bae500000048"
)

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func service(c *transport.Client) *customentity.Service {
	s, _ := customentity.NewService(c)
	return s
}

func TestFixtures(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		var ce customentity.CustomEntity
		if err := codec.Unmarshal(read(t, "upstream/custom_entity.json"), &ce); err != nil || ce.ID != "3f9a2f30-76af-11e7-6adb-ede50000000b" ||
			ce.Name != "custom dictionary" || ce.Meta.Type != "customentity" {
			t.Fatalf("%s: customentity = %+v, err = %v", codec.Name, ce, err)
		}
		var el customentity.Element
		if err := codec.Unmarshal(read(t, "upstream/custom_entity_element.json"), &el); err != nil {
			t.Fatal(err)
		}
		// description: null в фикстуре — пустая строка.
		if el.ID != elementID || el.Name != "Партнер1" || el.Code != "partner1" || el.ExternalCode != "5434665867876" ||
			el.Description != "" || len(el.Attributes) == 0 || el.Updated.IsZero() {
			t.Fatalf("%s: element = %+v", codec.Name, el)
		}
	}
}

func TestRequestJSON(t *testing.T) {
	owner := model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/employee/9bc5de98-79f2-11e8-1a0d-4e0d0000002b", Type: "employee", MediaType: "application/json"}}
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create", customentity.CreateRequest{Name: "custom dictionary"}, []byte(`{"name":"custom dictionary"}`)},
			{"update", customentity.UpdateRequest{Name: model.Value("new name")}, []byte(`{"name":"new name"}`)},
			{"element create", customentity.ElementCreateRequest{Name: "Партнер 3", Code: model.Value("partner3"), Description: model.Value("Описание"),
				ExternalCode: model.Value("5434665867876"), Shared: model.Value(false), Attributes: model.Value([]model.Attribute{}),
				Owner: model.Value(owner), Group: model.Null[model.Reference]()}, read(t, "custom/customentity/element-create.golden.json")},
			{"element update", customentity.ElementUpdateRequest{Description: model.Value(""), Shared: model.Value(true)}, []byte(`{"description":"","shared":true}`)},
		} {
			got, err := codec.Marshal(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, got, tc.want)
		}
	}
	for _, typ := range []reflect.Type{reflect.TypeFor[customentity.ElementCreateRequest](), reflect.TypeFor[customentity.ElementUpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Updated"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
}

func TestHTTP(t *testing.T) {
	const entityJSON = `{"meta":{"type":"customentity"},"id":"` + entityID + `","name":"n"}`
	const elementJSON = `{"meta":{"type":"customentity"},"id":"` + elementID + `","name":"n"}`
	id := func(v string, err error) (*string, error) {
		if err != nil {
			return nil, err
		}
		return &v, nil
	}
	for _, o := range []struct {
		name, method, path, body, response string
		invoke                             func(context.Context, *transport.Client) (*string, error)
	}{
		{"create", "POST", "entity/customentity", `{"name":"n"}`, entityJSON, func(ctx context.Context, c *transport.Client) (*string, error) {
			r, err := service(c).Create(ctx, customentity.CreateRequest{Name: "n"})
			if r == nil {
				return nil, err
			}
			return id(r.ID, err)
		}},
		{"update", "PUT", "entity/customentity/" + entityID, `{"name":"n"}`, entityJSON, func(ctx context.Context, c *transport.Client) (*string, error) {
			r, err := service(c).Update(ctx, entityID, customentity.UpdateRequest{Name: model.Value("n")})
			if r == nil {
				return nil, err
			}
			return id(r.ID, err)
		}},
		{"get element", "GET", "entity/customentity/" + entityID + "/" + elementID, "", elementJSON, func(ctx context.Context, c *transport.Client) (*string, error) {
			r, err := service(c).GetElement(ctx, entityID, elementID)
			if r == nil {
				return nil, err
			}
			return id(r.ID, err)
		}},
		{"create element", "POST", "entity/customentity/" + entityID, `{"name":"n"}`, elementJSON, func(ctx context.Context, c *transport.Client) (*string, error) {
			r, err := service(c).CreateElement(ctx, entityID, customentity.ElementCreateRequest{Name: "n"})
			if r == nil {
				return nil, err
			}
			return id(r.ID, err)
		}},
		{"update element", "PUT", "entity/customentity/" + entityID + "/" + elementID, `{"code":""}`, elementJSON, func(ctx context.Context, c *transport.Client) (*string, error) {
			r, err := service(c).UpdateElement(ctx, entityID, elementID, customentity.ElementUpdateRequest{Code: model.Value("")})
			if r == nil {
				return nil, err
			}
			return id(r.ID, err)
		}},
	} {
		t.Run(o.name, func(t *testing.T) {
			var body []byte
			if o.body != "" {
				body = []byte(o.body)
			}
			for _, tc := range documenttest.ResponseCases(o.response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, o.method, o.path, body, o.invoke, func(got *string) {
						if *got == "" {
							t.Fatal("empty id")
						}
					})
				})
			}
			documenttest.LocalErrors(t, o.invoke)
		})
	}
}

func TestListAndDelete(t *testing.T) {
	var got []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = append(got, r.Method+" "+r.URL.RequestURI())
		if r.Method == http.MethodGet {
			io.WriteString(w, `{"meta":{"href":"x","size":1},"rows":[{"id":"`+elementID+`","name":"n"}]}`)
		}
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(c)
	filter, err := query.Filter(query.Eq("name", "n"))
	if err != nil {
		t.Fatal(err)
	}
	if l, err := s.ListElements(t.Context(), entityID, query.ListOptions{Limit: 10, Filter: filter}); err != nil || len(l.Rows) != 1 {
		t.Fatalf("list = %+v, err = %v", l, err)
	}
	if err := s.DeleteElement(t.Context(), entityID, elementID); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(t.Context(), entityID); err != nil {
		t.Fatal(err)
	}
	want := []string{
		"GET /entity/customentity/" + entityID + "?filter=name%3Dn&limit=10",
		"DELETE /entity/customentity/" + entityID + "/" + elementID,
		"DELETE /entity/customentity/" + entityID,
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("requests:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(c), t.Context()
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"create no name": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, customentity.CreateRequest{}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, entityID, customentity.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, "bad", customentity.UpdateRequest{Name: model.Value("n")}, o...)
			return err
		},
		"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
		"list bad id": func(o ...transport.RequestOption) error {
			_, err := s.ListElements(ctx, "bad", query.ListOptions{}, o...)
			return err
		},
		"list fields": func(o ...transport.RequestOption) error {
			_, err := s.ListElements(ctx, entityID, query.ListOptions{Fields: "x"}, o...)
			return err
		},
		"get element bad id": func(o ...transport.RequestOption) error { _, err := s.GetElement(ctx, entityID, "x", o...); return err },
		"create element no name": func(o ...transport.RequestOption) error {
			_, err := s.CreateElement(ctx, entityID, customentity.ElementCreateRequest{}, o...)
			return err
		},
		"create element bad id": func(o ...transport.RequestOption) error {
			_, err := s.CreateElement(ctx, "../x", customentity.ElementCreateRequest{Name: "n"}, o...)
			return err
		},
		"update element empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateElement(ctx, entityID, elementID, customentity.ElementUpdateRequest{}, o...)
			return err
		},
		"delete element bad id": func(o ...transport.RequestOption) error { return s.DeleteElement(ctx, entityID, elementID+"/x", o...) },
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
