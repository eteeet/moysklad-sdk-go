package purchaseorder_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/purchaseorder"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/purchase_order.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "purchaseorder", row, "", "isDeleted", func(c *transport.Client) listtest.ListFunc[purchaseorder.PurchaseOrder] {
		s, err := purchaseorder.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got purchaseorder.PurchaseOrder) {
		if got.ID != "a1b2c3d4-1234-1234-1234-a1b2c3d4e5f6" || got.Name != "00001" {
			t.Fatalf("row = %+v", got)
		}
	})
}
