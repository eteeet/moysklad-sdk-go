package demand_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/demand"
	"github.com/eteeet/moysklad-sdk-go/internal/listtest"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func TestServiceList(t *testing.T) {
	row, err := os.ReadFile("../../testdata/upstream/demand.json")
	if err != nil {
		t.Fatal(err)
	}
	listtest.Run(t, "demand", row, "", "isDeleted", func(c *transport.Client) listtest.ListFunc[demand.Demand] {
		s, err := demand.NewService(c)
		if err != nil {
			t.Fatal(err)
		}
		return s.List
	}, func(t *testing.T, got demand.Demand) {
		if got.ID != "c7d8e9f0-2859-11f1-0a83-1f2500000301" || got.Name != "00001" {
			t.Fatalf("row = %+v", got)
		}
	})
}
