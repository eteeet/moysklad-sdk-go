package expenseitem_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/expenseitem"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/expense_item.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "expenseitem", row, "", "archived", func(c *transport.Client) listtest.ListFunc[expenseitem.ExpenseItem] {
		s, err := expenseitem.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got expenseitem.ExpenseItem) {
		if got.ID != "fb0a4b75-2e58-11e6-8a84-bae500000058" || got.Name != "Аренда" {
			t.Fatalf("row = %+v", got)
		}
	})
}
