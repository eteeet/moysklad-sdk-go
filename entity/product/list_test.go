package product_test

import (
	"os"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/product"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/product.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "product", row, "", "archived", func(c *transport.Client) listtest.ListFunc[product.Product] {
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

// Регрессия ревью 46e9d06: minimumStock не декодируется в Product, поэтому отклоняется до запроса.
func TestServiceListRejectsMinimumStock(t *testing.T) {
	client, err := transport.NewClient(transport.Config{BaseURL: "http://example.test"})
	if err != nil {
		t.Fatal(err)
	}
	s, err := product.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := s.List(t.Context(), query.ListOptions{Fields: "minimumStock"}); got != nil || err == nil ||
		!strings.Contains(err.Error(), "minimumStock") {
		t.Fatalf("got = %+v, err = %v", got, err)
	}
}
