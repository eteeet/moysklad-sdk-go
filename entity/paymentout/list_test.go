package paymentout_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/paymentout"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/payment_out.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "paymentout", row, "", "isDeleted", func(c *transport.Client) listtest.ListFunc[paymentout.PaymentOut] {
		s, err := paymentout.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got paymentout.PaymentOut) {
		if got.ID != "64aa2e8c-3f95-11e6-8a84-bae50000010f" || got.Name != "00001" {
			t.Fatalf("row = %+v", got)
		}
	})
}
