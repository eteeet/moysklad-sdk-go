package namedfilter_test

import (
	"context"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/product"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/namedfilter"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	filterID   = "b5863410-ca86-11eb-ac12-000d00000019"
	filterHref = "https://api.moysklad.ru/api/remap/1.2/entity/product/namedfilter/" + filterID
	ownHref    = "https://api.moysklad.ru/api/remap/1.2/entity/employee/234eee6f-c513-11eb-ac12-000d0000003b"
	otherHref  = "https://api.moysklad.ru/api/remap/1.2/entity/employee/25863410-ca86-11eb-ac12-000d00000234"
)

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../testdata/custom/namedfilter/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// checkFilter сравнивает все поля фильтра из примеров источника.
func checkFilter(t *testing.T, got namedfilter.NamedFilter, owner string) {
	t.Helper()
	if got.ID != filterID || got.Name != "filterName" || got.AccountID != "22ef0c54-c513-11eb-ac12-000700000002" ||
		got.Meta.Href != filterHref || got.Meta.Type != "namedfilter" || got.Meta.MediaType != "application/json" ||
		got.Owner == nil || got.Owner.Meta.Href != owner || got.Owner.Meta.Type != "employee" {
		t.Fatalf("filter = %+v, owner = %+v", got, got.Owner)
	}
}

func checkList(t *testing.T, got *model.Collection[namedfilter.NamedFilter], owner string) {
	t.Helper()
	if got.Meta.Href != "https://api.moysklad.ru/api/remap/1.2/entity/product/namedfilter" || got.Meta.Type != "namedfilter" ||
		got.Meta.Size != 1 || got.Meta.Limit != 1000 || got.Meta.Offset != 0 || got.Context == nil || len(got.Rows) != 1 {
		t.Fatalf("collection = %+v", got)
	}
	checkFilter(t, got.Rows[0], owner)
}

func TestFixtures(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for name, owner := range map[string]string{"list.json": ownHref, "list_owner.json": otherHref} {
				var got model.Collection[namedfilter.NamedFilter]
				if err := codec.Unmarshal(fixture(t, name), &got); err != nil {
					t.Fatal(err)
				}
				checkList(t, &got, owner)
			}
			var got namedfilter.NamedFilter
			if err := codec.Unmarshal(fixture(t, "get.json"), &got); err != nil {
				t.Fatal(err)
			}
			checkFilter(t, got, ownHref)
		})
	}
}

func TestServiceHTTP(t *testing.T) {
	service := func(c *transport.Client) *namedfilter.Service {
		s, err := namedfilter.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	list := func(owner string) func(context.Context, *transport.Client) (*model.Collection[namedfilter.NamedFilter], error) {
		return func(ctx context.Context, c *transport.Client) (*model.Collection[namedfilter.NamedFilter], error) {
			return service(c).List(ctx, "product", namedfilter.ListOptions{Owner: owner})
		}
	}
	t.Run("list", func(t *testing.T) {
		for _, tc := range documenttest.ResponseCases(string(fixture(t, "list.json"))) {
			t.Run(tc.Name, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", "entity/product/namedfilter", nil, list(""),
					func(got *model.Collection[namedfilter.NamedFilter]) { checkList(t, got, ownHref) })
			})
		}
		documenttest.LocalErrors(t, list(""))
	})
	t.Run("list owner", func(t *testing.T) {
		path := "entity/product/namedfilter?owner=https%3A%2F%2Fapi.moysklad.ru%2Fapi%2Fremap%2F1.2%2Fentity%2Femployee%2F25863410-ca86-11eb-ac12-000d00000234"
		for _, tc := range documenttest.ResponseCases(string(fixture(t, "list_owner.json"))) {
			t.Run(tc.Name, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", path, nil, list(otherHref),
					func(got *model.Collection[namedfilter.NamedFilter]) { checkList(t, got, otherHref) })
			})
		}
		documenttest.LocalErrors(t, list(otherHref))
	})
	t.Run("get", func(t *testing.T) {
		get := func(ctx context.Context, c *transport.Client) (*namedfilter.NamedFilter, error) {
			return service(c).Get(ctx, "product", filterID)
		}
		for _, tc := range documenttest.ResponseCases(string(fixture(t, "get.json"))) {
			t.Run(tc.Name, func(t *testing.T) {
				documenttest.HTTP(t, tc, "GET", "entity/product/namedfilter/"+filterID, nil, get,
					func(got *namedfilter.NamedFilter) { checkFilter(t, *got, ownHref) })
			})
		}
		documenttest.LocalErrors(t, get)
	})
}

// Неверные тип, id и owner отклоняются до HTTP; ResponseInfo сбрасывается.
func TestLocalValidation(t *testing.T) {
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls.Add(1) }))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s, err := namedfilter.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	for name, call := range map[string]func(transport.RequestOption) error{
		"list empty type": func(o transport.RequestOption) error {
			_, err := s.List(t.Context(), "", namedfilter.ListOptions{}, o)
			return err
		},
		"list path type": func(o transport.RequestOption) error {
			_, err := s.List(t.Context(), "product/../x", namedfilter.ListOptions{}, o)
			return err
		},
		"list upper type": func(o transport.RequestOption) error {
			_, err := s.List(t.Context(), "Product", namedfilter.ListOptions{}, o)
			return err
		},
		"list relative owner": func(o transport.RequestOption) error {
			_, err := s.List(t.Context(), "product", namedfilter.ListOptions{Owner: "entity/employee/1"}, o)
			return err
		},
		"get type": func(o transport.RequestOption) error {
			_, err := s.Get(t.Context(), "a/b", filterID, o)
			return err
		},
		"get id": func(o transport.RequestOption) error {
			_, err := s.Get(t.Context(), "product", "b5863410-ca86-11eb-ac12-000d0000001", o)
			return err
		},
		"get path id": func(o transport.RequestOption) error {
			_, err := s.Get(t.Context(), "product", "../../../entity/product/x-ca86-11eb-ac12-000d0000001", o)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 999, Attempts: 3}
		err := call(transport.WithResponseInfo(&info))
		if err == nil || strings.Count(err.Error(), "moysklad:") != 1 || !strings.Contains(err.Error(), "namedfilter") {
			t.Errorf("%s: err = %v", name, err)
		}
		if info.StatusCode != 0 || info.Attempts != 0 {
			t.Errorf("%s: response info not reset: %+v", name, info)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("invalid arguments sent: %d", calls.Load())
	}
	if _, err := namedfilter.NewService(nil); err == nil {
		t.Fatal("nil client accepted")
	}
}

// Применение: href фильтра из List уходит в namedfilter списка сущности через query.ListOptions.
func TestApplyFilterToList(t *testing.T) {
	list := fixture(t, "list.json")
	var calls atomic.Int32
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		switch r.URL.Path {
		case "/entity/product/namedfilter":
			w.Write(list)
		case "/entity/product":
			if got := r.URL.Query(); len(got) != 1 || got.Get("namedfilter") != filterHref {
				t.Errorf("product query = %q", r.URL.RawQuery)
			}
			w.Write([]byte(`{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product","size":0},"rows":[]}`))
		default:
			t.Errorf("unexpected %s", r.URL)
		}
	}))
	defer server.Close()
	client, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	filters, err := must(namedfilter.NewService(client)).List(t.Context(), "product", namedfilter.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	products, err := must(product.NewService(client)).List(t.Context(), query.ListOptions{NamedFilter: filters.Rows[0].Meta.Href})
	if err != nil || products.Rows == nil || calls.Load() != 2 {
		t.Fatalf("products = %+v, err = %v, calls = %d", products, err, calls.Load())
	}
}

func must[T any](v *T, err error) *T {
	if err != nil {
		panic(err)
	}
	return v
}
