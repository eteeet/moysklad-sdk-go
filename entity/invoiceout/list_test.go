package invoiceout_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/invoiceout"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/invoice_out.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "invoiceout", row, "", "isDeleted", func(c *transport.Client) listtest.ListFunc[invoiceout.InvoiceOut] {
		s, err := invoiceout.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got invoiceout.InvoiceOut) {
		if got.ID != "6afee657-0888-11e6-9464-e4de00000037" || got.Name != "change name" {
			t.Fatalf("row = %+v", got)
		}
	})
}
