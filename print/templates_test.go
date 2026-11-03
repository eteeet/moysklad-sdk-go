package print_test

import (
	"context"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/print"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func read(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../testdata/custom/print/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

const api = "https://api.moysklad.ru/api/remap/1.2/"

func TestTemplateFixtures(t *testing.T) {
	for _, tc := range []struct {
		file string
		rows int            // 0 — одиночный шаблон
		size int            // meta.size: у embedded-list источник пишет 6 при 7 строках
		want print.Template // первая строка списка или сам шаблон
	}{
		{"embedded-list.json", 7, 6, print.Template{
			Meta: model.Meta{Href: api + "entity/demand/metadata/embeddedtemplate/c2645cf0-43f0-3379-adff-99929da86d5e", Type: "embeddedtemplate", MediaType: "application/json"},
			ID:   "c2645cf0-43f0-3379-adff-99929da86d5e", Name: "Транспортная накладная", Type: "entity", Content: api + "download-template/ttn.xls"}},
		{"embedded-get.json", 0, 0, print.Template{
			Meta: model.Meta{Href: api + "entity/demand/metadata/embeddedtemplate/c2645cf0-43f0-3379-adff-99929da86d5e", Type: "embeddedtemplate", MediaType: "application/json"},
			ID:   "c2645cf0-43f0-3379-adff-99929da86d5e", Name: "Транспортная накладная", Type: "entity", Content: api + "download-template/ttn.xls"}},
		{"assortment-embedded-list.json", 2, 2, print.Template{
			Meta: model.Meta{Href: api + "entity/assortment/metadata/embeddedtemplate/a8218d14-6017-3ba4-85c4-254fff93bbc8", Type: "embeddedtemplate", MediaType: "application/json"},
			ID:   "a8218d14-6017-3ba4-85c4-254fff93bbc8", Name: "Термоэтикетка (58х40мм)", Type: "mxtemplate", Content: api + "download-template/termo_58x40.xml"}},
		{"assortment-embedded-get.json", 0, 0, print.Template{
			Meta: model.Meta{Href: api + "entity/assortment/metadata/embeddedtemplate/a8218d14-6017-3ba4-85c4-254fff93bbc8", Type: "embeddedtemplate", MediaType: "application/json"},
			ID:   "a8218d14-6017-3ba4-85c4-254fff93bbc8", Name: "Термоэтикетка (58х40мм)", Type: "mxtemplate", Content: api + "download-template/termo_58x40.xml"}},
		{"custom-list.json", 1, 1, print.Template{
			Meta: model.Meta{Href: api + "entity/customerorder/metadata/customtemplate/38d1c843-1601-11e7-8af5-581e00000009", Type: "customtemplate", MediaType: "application/json"},
			ID:   "38d1c843-1601-11e7-8af5-581e00000009", Name: "tratata", Type: "entity", Content: api + "download/38d1c843-1601-11e7-8af5-581e00000009"}},
		{"custom-get.json", 0, 0, print.Template{
			Meta: model.Meta{Href: api + "entity/customerorder/metadata/customtemplate/38d1c843-1601-11e7-8af5-581e00000009", Type: "customtemplate", MediaType: "application/json"},
			ID:   "38d1c843-1601-11e7-8af5-581e00000009", Name: "достаточно_нестандартный_шаблон", Type: "entity", Content: api + "download/38d1c843-1601-11e7-8af5-581e00000009"}},
		{"assortment-custom-get.json", 0, 0, print.Template{
			Meta: model.Meta{Href: api + "entity/assortment/metadata/customtemplate/a4d36b88-d329-4b73-966a-b66b6a1b1b9a", Type: "customtemplate", MediaType: "application/json"},
			ID:   "a4d36b88-d329-4b73-966a-b66b6a1b1b9a", Name: "NewCustomTemplate", Type: "mxtemplate", Content: api + "download/a4d36b88-d329-4b73-966a-b66b6a1b1b9a"}},
	} {
		for _, codec := range documenttest.Codecs {
			t.Run(tc.file+"/"+codec.Name, func(t *testing.T) {
				var got print.Template
				if tc.rows > 0 {
					var page model.Collection[print.Template]
					if err := codec.Unmarshal(read(t, tc.file), &page); err != nil {
						t.Fatal(err)
					}
					if len(page.Rows) != tc.rows || page.Meta.Size != tc.size {
						t.Fatalf("rows = %d, size = %d", len(page.Rows), page.Meta.Size)
					}
					got = page.Rows[0]
				} else if err := codec.Unmarshal(read(t, tc.file), &got); err != nil {
					t.Fatal(err)
				}
				if got != tc.want {
					t.Errorf("template = %+v\nwant %+v", got, tc.want)
				}
			})
		}
	}
}

const templateID = "c2645cf0-43f0-3379-adff-99929da86d5e"

type templateCall struct {
	name, path string
	fixture    string
	invoke     func(context.Context, *print.Service) (*print.Template, error)
}

func templateCalls() []templateCall {
	first := func(page *model.Collection[print.Template], err error) (*print.Template, error) {
		if err != nil {
			return nil, err
		}
		return &page.Rows[0], nil
	}
	return []templateCall{
		{"list embedded", "entity/demand/metadata/embeddedtemplate", "embedded-list.json", func(ctx context.Context, s *print.Service) (*print.Template, error) {
			return first(s.ListEmbeddedTemplates(ctx, "demand"))
		}},
		{"get embedded", "entity/demand/metadata/embeddedtemplate/" + templateID, "embedded-get.json", func(ctx context.Context, s *print.Service) (*print.Template, error) {
			return s.GetEmbeddedTemplate(ctx, "demand", templateID)
		}},
		{"list custom", "entity/demand/metadata/customtemplate", "custom-list.json", func(ctx context.Context, s *print.Service) (*print.Template, error) {
			return first(s.ListCustomTemplates(ctx, "demand"))
		}},
		{"get custom", "entity/demand/metadata/customtemplate/" + templateID, "custom-get.json", func(ctx context.Context, s *print.Service) (*print.Template, error) {
			return s.GetCustomTemplate(ctx, "demand", templateID)
		}},
	}
}

func service(t *testing.T, c *transport.Client) *print.Service {
	t.Helper()
	s, err := print.NewService(c)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// HTTP-контракт чтения шаблонов: путь, заголовки, plain/gzip, ошибки API, одна отправка.
func TestTemplatesHTTP(t *testing.T) {
	for _, op := range templateCalls() {
		body := read(t, op.fixture)
		var want print.Template
		if strings.HasPrefix(op.name, "list") {
			var page model.Collection[print.Template]
			if err := documenttest.Codecs[1].Unmarshal(body, &page); err != nil {
				t.Fatal(err)
			}
			want = page.Rows[0]
		} else if err := documenttest.Codecs[1].Unmarshal(body, &want); err != nil {
			t.Fatal(err)
		}
		for _, tc := range documenttest.ResponseCases(string(body)) {
			t.Run(op.name+"/"+tc.Name, func(t *testing.T) {
				documenttest.HTTP(t, tc, http.MethodGet, op.path, nil, func(ctx context.Context, c *transport.Client) (*print.Template, error) {
					return op.invoke(ctx, service(t, c))
				}, func(got *print.Template) {
					if *got != want {
						t.Errorf("template = %+v, want %+v", *got, want)
					}
				})
			})
		}
		t.Run(op.name+"/local", func(t *testing.T) {
			documenttest.LocalErrors(t, func(ctx context.Context, c *transport.Client) (*print.Template, error) {
				return op.invoke(ctx, service(t, c))
			})
		})
	}
}

// Неверный тип или id отклоняются до HTTP, ResponseInfo сброшен.
func TestTemplatesValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := service(t, client)
	for name, invoke := range map[string]func(*transport.ResponseInfo) error{
		"list embedded type": func(i *transport.ResponseInfo) error {
			_, err := s.ListEmbeddedTemplates(t.Context(), "Demand", transport.WithResponseInfo(i))
			return err
		},
		"list custom type": func(i *transport.ResponseInfo) error {
			_, err := s.ListCustomTemplates(t.Context(), "demand/../x", transport.WithResponseInfo(i))
			return err
		},
		"get embedded type": func(i *transport.ResponseInfo) error {
			_, err := s.GetEmbeddedTemplate(t.Context(), "", templateID, transport.WithResponseInfo(i))
			return err
		},
		"get embedded id": func(i *transport.ResponseInfo) error {
			_, err := s.GetEmbeddedTemplate(t.Context(), "demand", "x", transport.WithResponseInfo(i))
			return err
		},
		"get custom id": func(i *transport.ResponseInfo) error {
			_, err := s.GetCustomTemplate(t.Context(), "demand", templateID+"/", transport.WithResponseInfo(i))
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 200, Attempts: 1}
		err := invoke(&info)
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || info.StatusCode != 0 || info.Attempts != 0 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
	if calls.Load() != 0 {
		t.Errorf("calls = %d", calls.Load())
	}
}
