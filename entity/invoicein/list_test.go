package invoicein_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/invoicein"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/invoice_in.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "invoicein", row, "", "isDeleted", func(c *transport.Client) listtest.ListFunc[invoicein.InvoiceIn] {
		s, err := invoicein.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got invoicein.InvoiceIn) {
		if got.ID != "36c000ad-3f52-11e6-8a84-bae500000069" || got.Name != "00002" {
			t.Fatalf("row = %+v", got)
		}
	})
}
