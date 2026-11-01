package supply_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/supply"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/supply.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "supply", row, "", "isDeleted", func(c *transport.Client) listtest.ListFunc[supply.Supply] {
		s, err := supply.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got supply.Supply) {
		if got.ID != "a3a404e4-2e5d-11e6-8a84-bae5000000fd" || got.Name != "00001" {
			t.Fatalf("row = %+v", got)
		}
	})
}
