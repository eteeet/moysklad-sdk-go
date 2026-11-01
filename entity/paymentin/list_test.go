package paymentin_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/paymentin"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/payment_in.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "paymentin", row, "", "isDeleted", func(c *transport.Client) listtest.ListFunc[paymentin.PaymentIn] {
		s, err := paymentin.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got paymentin.PaymentIn) {
		if got.ID != "ad1ffc02-3f8e-11e6-8a84-bae5000000fa" || got.Name != "00001" {
			t.Fatalf("row = %+v", got)
		}
	})
}
