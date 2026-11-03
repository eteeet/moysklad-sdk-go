package retailstore_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/retailstore"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const (
	cashierStore = "ea05e0c9-8667-11e7-8a7f-40d000000060"
	cashierID    = "12345678-1234-1234-1234-123456789012"
)

func cashierService(c *transport.Client) *retailstore.Service {
	s, _ := retailstore.NewService(c)
	return s
}

func TestCashierHTTP(t *testing.T) {
	data, err := os.ReadFile("../../testdata/upstream/cashier.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range documenttest.Codecs {
		var got retailstore.Cashier
		if err := codec.Unmarshal(data, &got); err != nil || got.ID != cashierID || got.Meta.Type != "cashier" ||
			got.Employee == nil || got.RetailStore == nil || got.RetailStore.Meta.Type != "retailstore" {
			t.Fatalf("%s: cashier = %+v, err = %v", codec.Name, got, err)
		}
	}
	invoke := func(ctx context.Context, c *transport.Client) (*retailstore.Cashier, error) {
		return cashierService(c).GetCashier(ctx, cashierStore, cashierID)
	}
	for _, tc := range documenttest.ResponseCases(string(data)) {
		t.Run(tc.Name, func(t *testing.T) {
			documenttest.HTTP(t, tc, "GET", "entity/retailstore/"+cashierStore+"/cashiers/"+cashierID, nil, invoke, func(got *retailstore.Cashier) {
				if got.Employee == nil {
					t.Fatalf("cashier = %+v", got)
				}
			})
		})
	}
	documenttest.LocalErrors(t, invoke)
}

func TestListCashiers(t *testing.T) {
	var uri string
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri = r.Method + " " + r.URL.RequestURI()
		io.WriteString(w, `{"meta":{"href":"x","size":1},"rows":[{"id":"`+cashierID+`"}]}`)
	}))
	defer server.Close()
	c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	got, err := cashierService(c).ListCashiers(t.Context(), cashierStore, query.PageOptions{Limit: 50, Offset: 100})
	if err != nil || len(got.Rows) != 1 || uri != "GET /entity/retailstore/"+cashierStore+"/cashiers?limit=50&offset=100" {
		t.Fatalf("got = %+v, err = %v, uri = %s", got, err, uri)
	}
	for name, invoke := range map[string]func(...transport.RequestOption) error{
		"list bad store": func(o ...transport.RequestOption) error {
			_, err := cashierService(c).ListCashiers(t.Context(), "bad", query.PageOptions{}, o...)
			return err
		},
		"get bad id": func(o ...transport.RequestOption) error {
			_, err := cashierService(c).GetCashier(t.Context(), cashierStore, "x", o...)
			return err
		},
	} {
		info := transport.ResponseInfo{StatusCode: 503}
		if err := invoke(transport.WithResponseInfo(&info)); err == nil || info.StatusCode != 0 || strings.Count(err.Error(), "moysklad:") != 1 {
			t.Errorf("%s: err = %v, info = %+v", name, err, info)
		}
	}
}
