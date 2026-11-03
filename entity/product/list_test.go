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

// fields принимает только minimumStock; другие значения — ошибка до запроса.
func TestServiceListRejectsOtherFields(t *testing.T) {
	client, err := transport.NewClient(transport.Config{BaseURL: "http://example.test"})
	if err != nil {
		t.Fatal(err)
	}
	s, err := product.NewService(client)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := s.List(t.Context(), query.ListOptions{Fields: "stock"}); got != nil || err == nil ||
		!strings.Contains(err.Error(), "stock") {
		t.Fatalf("got = %+v, err = %v", got, err)
	}
}
