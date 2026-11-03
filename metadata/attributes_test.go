package metadata_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/metadata"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	crudAttrID   = "33b2fe47-b465-11e9-7ae5-884b0001562f"
	crudAttrPath = "entity/product/metadata/attributes"
	crudAttrJSON = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product/metadata/attributes/` + crudAttrID + `","type":"attributemetadata","mediaType":"application/json"},` +
		`"id":"` + crudAttrID + `","name":"Наличие CD/DVD-Rom","type":"boolean","required":false}`
)

func attrService(c *transport.Client) *metadata.Service {
	s, _ := metadata.NewService(c)
	return s
}

func TestAttributeRequestJSON(t *testing.T) {
	ce := model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/context/companysettings/metadata/customEntities/0347beb0-a785-11e9-ac12-000800000003", Type: "customentitymetadata"}
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  string
		}{
			{"create string", metadata.AttributeCreateRequest{Name: "Материал корпуса", Type: "string"}, `{"name":"Материал корпуса","type":"string"}`},
			{"create customentity", metadata.AttributeCreateRequest{Name: "Цвет", Type: "customentity", Required: model.Value(false), Show: model.Value(true),
				Description: model.Value(""), CustomEntityMeta: model.Value(ce)},
				`{"name":"Цвет","type":"customentity","required":false,"show":true,"description":"","customEntityMeta":{"href":"` + ce.Href + `","type":"customentitymetadata"}}`},
			{"update", metadata.AttributeUpdateRequest{Name: model.Value("Наличие CD/DVD-Rom"), Required: model.Value(false)}, `{"name":"Наличие CD/DVD-Rom","required":false}`},
		} {
			got, err := codec.Marshal(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, got, []byte(tc.want))
		}
	}
}

func TestAttributeHTTP(t *testing.T) {
	href := func(a *metadata.Attribute, err error) (*string, error) {
		if err != nil {
			return nil, err
		}
		return &a.Meta.Href, nil
	}
	for _, o := range []struct {
		name, method, path, body string
		invoke                   func(context.Context, *transport.Client) (*string, error)
	}{
		{"create", "POST", crudAttrPath, `{"name":"Наличие CD-Rom","type":"boolean"}`, func(ctx context.Context, c *transport.Client) (*string, error) {
			return href(attrService(c).CreateAttribute(ctx, "product", metadata.AttributeCreateRequest{Name: "Наличие CD-Rom", Type: "boolean"}))
		}},
		{"update", "PUT", crudAttrPath + "/" + crudAttrID, `{"name":"Наличие CD/DVD-Rom"}`, func(ctx context.Context, c *transport.Client) (*string, error) {
			return href(attrService(c).UpdateAttribute(ctx, "product", crudAttrID, metadata.AttributeUpdateRequest{Name: model.Value("Наличие CD/DVD-Rom")}))
		}},
	} {
		t.Run(o.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(crudAttrJSON) {
				// Критерий успеха — meta.href, не id: ответы {"id":""} и без id и так без meta.
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, o.method, o.path, []byte(o.body), o.invoke, func(got *string) {
						if !strings.HasSuffix(*got, crudAttrID) {
							t.Fatalf("href = %s", *got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, o.invoke)
		})
	}
}

// Пример массового изменения источника: meta существующего поля и два новых в одном POST.
func TestAttributeBatch(t *testing.T) {
	var got []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = append(got, r.Method+" "+r.URL.RequestURI()+" "+string(body))
		switch {
		case strings.HasSuffix(r.URL.Path, "/delete"):
			io.WriteString(w, `[{"info":"ok"}]`)
		case r.Method == http.MethodPost:
			io.WriteString(w, `[`+crudAttrJSON+`]`)
		}
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := attrService(c), t.Context()
	if r, err := s.CreateManyAttributes(ctx, "product", []metadata.AttributeCreateRequest{{Name: "Наличие type-C разъема", Type: "boolean"}}); err != nil || r[0].Value.ID != crudAttrID {
		t.Fatalf("create many = %+v, err = %v", r, err)
	}
	if _, err := s.UpdateManyAttributes(ctx, "product", []batch.Update[metadata.AttributeUpdateRequest]{{ID: crudAttrID,
		Request: metadata.AttributeUpdateRequest{Name: model.Value("Материал корпуса (дополнение)"), Type: model.Value("string")}}}); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteAttribute(ctx, "product", crudAttrID); err != nil {
		t.Fatal(err)
	}
	if r, err := s.DeleteManyAttributes(ctx, "product", []string{crudAttrID}); err != nil || r[0].Info != "ok" {
		t.Fatalf("delete many = %+v, err = %v", r, err)
	}
	meta := `{"href":"` + server.URL + "/" + crudAttrPath + "/" + crudAttrID + `","type":"attributemetadata","mediaType":"application/json"}`
	want := []string{
		"POST /" + crudAttrPath + ` [{"name":"Наличие type-C разъема","type":"boolean"}]`,
		"POST /" + crudAttrPath + ` [{"meta":` + meta + `,"name":"Материал корпуса (дополнение)","type":"string"}]`,
		"DELETE /" + crudAttrPath + "/" + crudAttrID + " ",
		"POST /" + crudAttrPath + `/delete [{"meta":` + meta + `}]`,
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("requests:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestAttributeValidation(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := attrService(c), t.Context()
	ok := metadata.AttributeCreateRequest{Name: "n", Type: "string"}
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"create bad type": func(o ...transport.RequestOption) error {
			_, err := s.CreateAttribute(ctx, "../x", ok, o...)
			return err
		},
		"create no name": func(o ...transport.RequestOption) error {
			_, err := s.CreateAttribute(ctx, "product", metadata.AttributeCreateRequest{Type: "string"}, o...)
			return err
		},
		"create no type": func(o ...transport.RequestOption) error {
			_, err := s.CreateAttribute(ctx, "product", metadata.AttributeCreateRequest{Name: "n"}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateAttribute(ctx, "product", "x", metadata.AttributeUpdateRequest{Name: model.Value("n")}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateAttribute(ctx, "product", crudAttrID, metadata.AttributeUpdateRequest{}, o...)
			return err
		},
		"update bad entity": func(o ...transport.RequestOption) error {
			_, err := s.UpdateAttribute(ctx, "Product", crudAttrID, metadata.AttributeUpdateRequest{Name: model.Value("n")}, o...)
			return err
		},
		"delete bad id": func(o ...transport.RequestOption) error { return s.DeleteAttribute(ctx, "product", "x", o...) },
		"create many bad": func(o ...transport.RequestOption) error {
			_, err := s.CreateManyAttributes(ctx, "x/y", []metadata.AttributeCreateRequest{ok}, o...)
			return err
		},
		"create many no type": func(o ...transport.RequestOption) error {
			_, err := s.CreateManyAttributes(ctx, "product", []metadata.AttributeCreateRequest{ok, {Name: "n"}}, o...)
			return err
		},
		"create many empty": func(o ...transport.RequestOption) error {
			_, err := s.CreateManyAttributes(ctx, "product", nil, o...)
			return err
		},
		"update many bad id": func(o ...transport.RequestOption) error {
			_, err := s.UpdateManyAttributes(ctx, "product", []batch.Update[metadata.AttributeUpdateRequest]{{ID: "x", Request: metadata.AttributeUpdateRequest{Name: model.Value("n")}}}, o...)
			return err
		},
		"delete many bad": func(o ...transport.RequestOption) error {
			_, err := s.DeleteManyAttributes(ctx, "product", []string{"x"}, o...)
			return err
		},
		"delete many bad ent": func(o ...transport.RequestOption) error {
			_, err := s.DeleteManyAttributes(ctx, "", []string{crudAttrID}, o...)
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
