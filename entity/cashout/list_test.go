package cashout_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/cashout"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/cash_out.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "cashout", row, "", "isDeleted", func(c *transport.Client) listtest.ListFunc[cashout.CashOut] {
		s, err := cashout.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got cashout.CashOut) {
		if got.ID != "e446a227-41b1-11e6-8a84-bae500000005" || got.Name != "0721" {
			t.Fatalf("row = %+v", got)
		}
	})
}
