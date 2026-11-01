package cashin_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/cashin"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/cash_in.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "cashin", row, "", "isDeleted", func(c *transport.Client) listtest.ListFunc[cashin.CashIn] {
		s, err := cashin.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got cashin.CashIn) {
		if got.ID != "b0ec3c98-3fac-11e6-8a84-bae500000131" || got.Name != "00001" {
			t.Fatalf("row = %+v", got)
		}
	})
}
