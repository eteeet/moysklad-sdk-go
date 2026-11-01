package product_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/product"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/product.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "product", row, "minimumStock", "archived", func(c *transport.Client) listtest.ListFunc[product.Product] {
		s, err := product.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got product.Product) {
		if got.ID != "12345678-1234-1234-1234-123456789012" || got.Name != "Тестовый товар" {
			t.Fatalf("row = %+v", got)
		}
	})
}
