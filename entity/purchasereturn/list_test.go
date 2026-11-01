package purchasereturn_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/purchasereturn"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/purchase_return.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "purchasereturn", row, "", "isDeleted", func(c *transport.Client) listtest.ListFunc[purchasereturn.PurchaseReturn] {
		s, err := purchasereturn.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got purchasereturn.PurchaseReturn) {
		if got.ID != "0ea5ec6c-58f3-11f1-0a83-1c0f00000817" || got.Name != "ATK-179-PR-00001" {
			t.Fatalf("row = %+v", got)
		}
	})
}
