package contentcard_test

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/contentcard"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const testID = "529eb5b6-d726-11f0-0a80-073800000329"

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ref(path, typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/" + path, Type: model.MetaType(typ), MediaType: "application/json"}}
}

var (
	full   = contentcard.CreateRequest{Name: "товар", Description: "Описание", CardContentName: "Название", Assortment: ref("entity/product/013c258e-d667-11f0-0a80-073800000148", "product"), SalePlatform: ref("entity/saleplatform/0f9b8b4c-4597-4772-b896-c8aff047a117", "saleplatform"), SalesChannels: []model.Reference{ref("entity/saleschannel/31e8ce0f-c99f-4824-a930-83fe7de68cd6", "saleschannel")}}
	update = contentcard.UpdateRequest{Description: model.Value("Новое описание"), SalesChannels: model.Value([]model.Reference{})}
)

func TestContentCardFixture(t *testing.T) {
	data := read(t, "custom/contentcard/contentcard.json")
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var got contentcard.ContentCard
			if err := codec.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.ID != testID || got.Name != "товар" || got.Description != "Описание" || got.CardContentName != "Название" || got.Assortment == nil || got.SalePlatform == nil || got.SalePlatform.Meta.Type != "saleplatform" || len(got.SalesChannels) != 1 || got.Images == nil {
				t.Fatalf("contentcard = %+v", got)
			}
			// Неизвестные поля не ломают разбор, отсутствующие не выдумываются.
			var sparse contentcard.ContentCard
			if err := codec.Unmarshal([]byte(`{"id":"x","futureField":{"a":1}}`), &sparse); err != nil || sparse.ID != "x" {
				t.Fatalf("sparse = %+v, err = %v", sparse, err)
			}
		})
	}
}

func TestRequestJSON(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			name  string
			input any
			want  []byte
		}{
			{"create full", full, read(t, "custom/contentcard/create-full.golden.json")},
			{"update", update, read(t, "custom/contentcard/update.golden.json")},
		} {
			t.Run(codec.Name+"/"+tc.name, func(t *testing.T) {
				got, err := codec.Marshal(tc.input)
				if err != nil {
					t.Fatal(err)
				}
				documenttest.JSONEqual(t, got, tc.want)
			})
		}
	}
	for _, typ := range []reflect.Type{reflect.TypeFor[contentcard.CreateRequest](), reflect.TypeFor[contentcard.UpdateRequest]()} {
		for _, name := range []string{"ID", "Meta", "AccountID", "Images"} {
			if _, ok := typ.FieldByName(name); ok {
				t.Errorf("%s contains read-only %s", typ.Name(), name)
			}
		}
	}
	if _, ok := reflect.TypeFor[contentcard.UpdateRequest]().FieldByName("Assortment"); ok {
		t.Error("UpdateRequest contains create-only Assortment")
	}
}

func service(c *transport.Client) *contentcard.Service {
	// NewService падает только на nil client, которого тесты не передают.
	s, _ := contentcard.NewService(c)
	return s
}

func TestServiceHTTP(t *testing.T) {
	const response = `{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/contentcard/` + testID + `"},"id":"` + testID + `","name":"n"}`
	for _, op := range []struct {
		name, method, path string
		body               []byte
		invoke             func(context.Context, *transport.Client) (*contentcard.ContentCard, error)
	}{
		{"get", "GET", "entity/contentcard/" + testID, nil, func(ctx context.Context, c *transport.Client) (*contentcard.ContentCard, error) {
			return service(c).Get(ctx, testID)
		}},
		{"create", "POST", "entity/contentcard", read(t, "custom/contentcard/create-full.golden.json"), func(ctx context.Context, c *transport.Client) (*contentcard.ContentCard, error) {
			return service(c).Create(ctx, full)
		}},
		{"update", "PUT", "entity/contentcard/" + testID, read(t, "custom/contentcard/update.golden.json"), func(ctx context.Context, c *transport.Client) (*contentcard.ContentCard, error) {
			return service(c).Update(ctx, testID, update)
		}},
	} {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range documenttest.ResponseCases(response) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, op.method, op.path, op.body, op.invoke, func(got *contentcard.ContentCard) {
						if got.ID != testID || got.Name != "n" {
							t.Fatalf("contentcard = %+v", got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, op.invoke)
		})
	}
}

// Ошибки входа — до HTTP; ResponseInfo прошлого вызова сброшен.
func TestServiceValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, ctx := service(client), t.Context()
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"get bad id": func(o ...transport.RequestOption) error { _, err := s.Get(ctx, "bad", o...); return err },
		"get path":   func(o ...transport.RequestOption) error { _, err := s.Get(ctx, testID+"/x", o...); return err },
		"create no channels": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, contentcard.CreateRequest{Name: "n", Description: "d", CardContentName: "c", Assortment: ref("entity/product/1", "product"), SalePlatform: ref("entity/saleplatform/1", "saleplatform")}, o...)
			return err
		},
		"create no platform": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, contentcard.CreateRequest{Name: "n", Description: "d", CardContentName: "c", Assortment: ref("entity/product/1", "product"), SalesChannels: []model.Reference{ref("entity/saleschannel/1", "saleschannel")}}, o...)
			return err
		},
		"create no card name": func(o ...transport.RequestOption) error {
			_, err := s.Create(ctx, contentcard.CreateRequest{Name: "n", Description: "d", Assortment: ref("entity/product/1", "product"), SalePlatform: ref("entity/saleplatform/1", "saleplatform"), SalesChannels: []model.Reference{ref("entity/saleschannel/1", "saleschannel")}}, o...)
			return err
		},
		"update empty": func(o ...transport.RequestOption) error {
			_, err := s.Update(ctx, testID, contentcard.UpdateRequest{}, o...)
			return err
		},
		"update bad id": func(o ...transport.RequestOption) error { _, err := s.Update(ctx, "bad", update, o...); return err },
		"delete bad id": func(o ...transport.RequestOption) error { return s.Delete(ctx, "bad", o...) },
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

func TestServiceList(t *testing.T) {
	listtest.Run(t, "contentcard", read(t, "custom/contentcard/contentcard.json"), "", "", func(c *transport.Client) listtest.ListFunc[contentcard.ContentCard] {
		return service(c).List
	}, func(t *testing.T, got contentcard.ContentCard) {
		if got.ID != testID || got.CardContentName != "Название" {
			t.Fatalf("row = %+v", got)
		}
	})
}
