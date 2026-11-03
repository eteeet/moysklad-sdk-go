package bundle_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/bundle"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	mutBundleID    = "7944ef04-f831-11e5-7a69-971500188b19"
	mutComponentID = "34f6344f-015e-11e6-9464-e4de0000006c"
	componentPath  = "entity/bundle/" + mutBundleID + "/components"
	componentJSON  = `{"meta":{"type":"bundlecomponent"},"id":"` + mutComponentID + `","assortment":{"meta":{"type":"product"}},"quantity":50}`
	productURL     = "https://api.moysklad.ru/api/remap/1.2/entity/product/010afdea-ea40-11e6-5bed-427b000000b8"
)

func componentService(c *transport.Client) *bundle.Service {
	s, _ := bundle.NewService(c)
	return s
}

// Примеры источника: добавление {assortment, quantity}, изменение {quantity: 50}.
func TestComponentMutationsHTTP(t *testing.T) {
	product := model.Reference{Meta: model.Meta{Href: productURL, Type: "product", MediaType: "application/json"}}
	id := func(c *bundle.Component, err error) (*string, error) {
		if err != nil {
			return nil, err
		}
		return &c.ID, nil
	}
	for _, o := range []struct {
		name, method, path, body string
		invoke                   func(context.Context, *transport.Client) (*string, error)
	}{
		{"add", "POST", componentPath, `{"assortment":{"meta":{"href":"` + productURL + `","type":"product","mediaType":"application/json"}},"quantity":1}`,
			func(ctx context.Context, c *transport.Client) (*string, error) {
				return id(componentService(c).AddComponent(ctx, mutBundleID, bundle.ComponentAddRequest{Assortment: product, Quantity: 1}))
			}},
		{"add zero quantity", "POST", componentPath, `{"assortment":{"meta":{"href":"` + productURL + `","type":"product","mediaType":"application/json"}},"quantity":0}`,
			func(ctx context.Context, c *transport.Client) (*string, error) {
				return id(componentService(c).AddComponent(ctx, mutBundleID, bundle.ComponentAddRequest{Assortment: product}))
			}},
		{"update", "PUT", componentPath + "/" + mutComponentID, `{"quantity":50}`, func(ctx context.Context, c *transport.Client) (*string, error) {
			return id(componentService(c).UpdateComponent(ctx, mutBundleID, mutComponentID, bundle.ComponentUpdateRequest{Quantity: model.Value(50.0)}))
		}},
	} {
		t.Run(o.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(componentJSON) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, o.method, o.path, []byte(o.body), o.invoke, func(got *string) {
						if *got != mutComponentID {
							t.Fatalf("id = %s", *got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, o.invoke)
		})
	}
}

func TestComponentRequestJSON(t *testing.T) {
	product := model.Reference{Meta: model.Meta{Href: productURL, Type: "product"}}
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			input any
			want  string
		}{
			{bundle.ComponentAddRequest{Assortment: product, Quantity: 0.5}, `{"assortment":{"meta":{"href":"` + productURL + `","type":"product"}},"quantity":0.5}`},
			{bundle.ComponentUpdateRequest{Quantity: model.Value(0.0), Assortment: model.Null[model.Reference]()}, `{"assortment":null,"quantity":0}`},
		} {
			got, err := codec.Marshal(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			documenttest.JSONEqual(t, got, []byte(tc.want))
		}
	}
}

func TestComponentDelete(t *testing.T) {
	var got []string
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = append(got, r.Method+" "+r.URL.RequestURI()+" "+string(body))
		if strings.HasSuffix(r.URL.Path, "/delete") {
			io.WriteString(w, `[{"info":"ok"}]`)
		}
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := componentService(c)
	if err := s.DeleteComponent(t.Context(), mutBundleID, mutComponentID); err != nil {
		t.Fatal(err)
	}
	if r, err := s.DeleteManyComponents(t.Context(), mutBundleID, []string{mutComponentID}); err != nil || len(r) != 1 {
		t.Fatalf("delete many = %+v, err = %v", r, err)
	}
	meta := `{"href":"` + server.URL + "/" + componentPath + "/" + mutComponentID + `","type":"bundlecomponent","mediaType":"application/json"}`
	want := "DELETE /" + componentPath + "/" + mutComponentID + " \nPOST /" + componentPath + `/delete [{"meta":` + meta + `}]`
	if strings.Join(got, "\n") != want {
		t.Fatalf("requests:\n%s\nwant:\n%s", strings.Join(got, "\n"), want)
	}
}

func TestComponentValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := componentService(c), t.Context()
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"add no assortment": func(o ...transport.RequestOption) error {
			_, err := s.AddComponent(ctx, mutBundleID, bundle.ComponentAddRequest{Quantity: 1}, o...)
			return err
		},
		"add bad bundle": func(o ...transport.RequestOption) error {
			_, err := s.AddComponent(ctx, "bad", bundle.ComponentAddRequest{Assortment: model.Reference{Meta: model.Meta{Href: productURL}}}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.UpdateComponent(ctx, mutBundleID, mutComponentID, bundle.ComponentUpdateRequest{}, o...)
			return err
		},
		"update bad component": func(o ...transport.RequestOption) error {
			_, err := s.UpdateComponent(ctx, mutBundleID, "x", bundle.ComponentUpdateRequest{Quantity: model.Value(1.0)}, o...)
			return err
		},
		"delete bad": func(o ...transport.RequestOption) error { return s.DeleteComponent(ctx, mutBundleID, "../x", o...) },
		"delete many empty": func(o ...transport.RequestOption) error {
			_, err := s.DeleteManyComponents(ctx, mutBundleID, nil, o...)
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
