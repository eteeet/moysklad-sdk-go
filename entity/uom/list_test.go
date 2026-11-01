package uom_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/uom"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/uom.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "uom", row, "", "", func(c *transport.Client) listtest.ListFunc[uom.UOM] {
		s, err := uom.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got uom.UOM) {
		if got.ID != "77777777-7777-7777-7777-777777777777" || got.Name != "шт" {
			t.Fatalf("row = %+v", got)
		}
	})
}
