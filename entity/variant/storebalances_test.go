package variant_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/variant"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	balanceOwner = "88df2fc1-1065-11f0-ac15-000700000012"
	balanceID    = "de922f92-113b-11f0-ac15-001000000009"
	balancePath  = "entity/variant/" + balanceOwner + "/storebalances"
	storeHref    = "https://api.moysklad.ru/api/remap/1.2/entity/store/da1c5548-0ecf-11f0-ac15-001200000125"
)

func balanceService(c *transport.Client) *variant.Service {
	s, _ := variant.NewService(c)
	return s
}

func TestStoreBalancesHTTP(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/upstream/product_storebalance.json")
	if err != nil {
		t.Fatal(err)
	}
	request := model.StoreBalanceRequest{Store: model.Value(model.Reference{Meta: model.Meta{Href: storeHref, Type: "store"}}), Quantity: model.Value(2.5)}
	body := `{"store":{"meta":{"href":"` + storeHref + `","type":"store"}},"quantity":2.5}`
	id := func(b *model.StoreBalance, err error) (*string, error) {
		if err != nil {
			return nil, err
		}
		return &b.ID, nil
	}
	for _, o := range []struct {
		name, method, path, body string
		invoke                   func(context.Context, *transport.Client) (*string, error)
	}{
		{"get", "GET", balancePath + "/" + balanceID, "", func(ctx context.Context, c *transport.Client) (*string, error) {
			return id(balanceService(c).GetStoreBalance(ctx, balanceOwner, balanceID))
		}},
		{"create", "POST", balancePath, body, func(ctx context.Context, c *transport.Client) (*string, error) {
			return id(balanceService(c).CreateStoreBalance(ctx, balanceOwner, request))
		}},
		{"update", "PUT", balancePath + "/" + balanceID, `{"quantity":0}`, func(ctx context.Context, c *transport.Client) (*string, error) {
			return id(balanceService(c).UpdateStoreBalance(ctx, balanceOwner, balanceID, model.StoreBalanceRequest{Quantity: model.Value(0.0)}))
		}},
	} {
		t.Run(o.name, func(t *testing.T) {
			var b []byte
			if o.body != "" {
				b = []byte(o.body)
			}
			for _, tc := range documenttest.ResponseCases(string(fixture)) {
				t.Run(tc.Name, func(t *testing.T) {
					documenttest.HTTP(t, tc, o.method, o.path, b, o.invoke, func(got *string) {
						if *got != balanceID {
							t.Fatalf("id = %s", *got)
						}
					})
				})
			}
			documenttest.LocalErrors(t, o.invoke)
		})
	}
	for _, codec := range documenttest.Codecs {
		var got model.StoreBalance
		if err := codec.Unmarshal(fixture, &got); err != nil || got.Quantity != 2.5 || got.Store == nil || got.Store.Meta.Href != storeHref || got.Meta.Type != "minimumstock" {
			t.Fatalf("%s: balance = %+v, err = %v", codec.Name, got, err)
		}
	}
}

func TestStoreBalancesCollection(t *testing.T) {
	var got []string
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		got = append(got, r.Method+" "+r.URL.RequestURI()+" "+string(b))
		switch {
		case r.Method == http.MethodGet:
			io.WriteString(w, `{"meta":{"href":"x","size":1},"rows":[{"id":"`+balanceID+`","quantity":1}]}`)
		case strings.HasSuffix(r.URL.Path, "/delete"):
			io.WriteString(w, `[{"info":"ok"}]`)
		}
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	s := balanceService(c)
	if l, err := s.ListStoreBalances(t.Context(), balanceOwner, query.PageOptions{Limit: 3}); err != nil || len(l.Rows) != 1 || l.Rows[0].Quantity != 1 {
		t.Fatalf("list = %+v, err = %v", l, err)
	}
	if err := s.DeleteStoreBalance(t.Context(), balanceOwner, balanceID); err != nil {
		t.Fatal(err)
	}
	if r, err := s.DeleteManyStoreBalances(t.Context(), balanceOwner, []string{balanceID}); err != nil || r[0].Info != "ok" {
		t.Fatalf("delete many = %+v, err = %v", r, err)
	}
	meta := `{"href":"` + server.URL + "/" + balancePath + "/" + balanceID + `","type":"minimumstock","mediaType":"application/json"}`
	want := []string{"GET /" + balancePath + "?limit=3&offset=0 ", "DELETE /" + balancePath + "/" + balanceID + " ", "POST /" + balancePath + `/delete [{"meta":` + meta + `}]`}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("requests:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"list bad": func(o ...transport.RequestOption) error {
			_, err := s.ListStoreBalances(t.Context(), "bad", query.PageOptions{}, o...)
			return err
		},
		"create bad": func(o ...transport.RequestOption) error {
			_, err := s.CreateStoreBalance(t.Context(), "../x", model.StoreBalanceRequest{}, o...)
			return err
		},
		"update zero": func(o ...transport.RequestOption) error {
			_, err := s.UpdateStoreBalance(t.Context(), balanceOwner, balanceID, model.StoreBalanceRequest{}, o...)
			return err
		},
		"delete bad": func(o ...transport.RequestOption) error {
			return s.DeleteStoreBalance(t.Context(), balanceOwner, "x", o...)
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
}
